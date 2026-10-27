package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// The port is bound before the database and the routes.
//
// It used to be bound at ListenAndServe, after everything: the database
// connected, two hundred routes registered and printed themselves, the
// dashboards announced their URLs, and then one line said the server had not
// started. On a hot reload that line scrolls past while the previous build
// keeps answering with the code from before the edit, so the only symptom is
// that the change does nothing (grit#110).
func TestMainBindsThePortFirst(t *testing.T) {
	src := apiMainGo(Options{})

	listen := strings.Index(src, "devserver.Listen(cfg.Port, cfg.AppEnv)")
	if listen < 0 {
		t.Fatal("main.go does not bind the port through devserver.Listen")
	}
	connect := strings.Index(src, "database.Connect(")
	if connect < 0 {
		t.Fatal("main.go no longer connects to the database")
	}
	if listen > connect {
		t.Error("the port is bound after the database, so a port already in use " +
			"is still reported at the end of startup")
	}
	if strings.Contains(src, "srv.ListenAndServe()") {
		t.Error("main.go still binds at ListenAndServe, which undoes the above")
	}
	if !strings.Contains(src, "srv.Serve(listener)") {
		t.Error("main.go does not serve the listener it bound")
	}
	// And the shutdown channel is in place before the server, because the
	// step-aside endpoint shuts down through it.
	if strings.Index(src, "quit := make(chan os.Signal, 1)") >
		strings.Index(src, "devserver.MountStepAside(") {
		t.Error("MountStepAside is wired before the channel it sends on")
	}
}

// Taking a port back is a development convenience. An unconfigured
// deployment must never do it, and config.Load defaults APP_ENV to
// "production" precisely so that unset is the safe answer.
func TestDevServerIsDevelopmentOnly(t *testing.T) {
	src := apiDevServerGo()

	for _, want := range []string{
		`case "dev", "development", "local":`,
		"func isLoopback(",
		"if !isLoopback(c.ClientIP()) {",
		`if !isDev(env) || router == nil || stop == nil {`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the devserver template is missing %q", want)
		}
	}
	// The empty string must not be development: it would make an API with no
	// APP_ENV set willing to stop another process.
	if strings.Contains(src, `case "", "dev", "development", "local":`) {
		t.Error("an empty APP_ENV counts as development, which disagrees with " +
			"config.Load defaulting it to production")
	}
	// Nothing is killed and no pid is guessed: that is the whole reason this
	// is an HTTP request rather than a signal. The word taskkill does appear,
	// in the hint printed for the person reading the failure, so the check is
	// on what the code can do rather than on the text.
	for _, forbidden := range []string{"os/exec", "exec.Command", "syscall.Kill", ".Process.Kill"} {
		if strings.Contains(src, forbidden) {
			t.Errorf("the devserver template reaches for %q; a process that is "+
				"not ours must be left alone", forbidden)
		}
	}
}

// air's two reload settings are Unix-only, and the config says so, because the
// report that found this proposed them as the Windows fix.
func TestAirConfigSaysWhichPlatformItHelps(t *testing.T) {
	cfg := airConfig(Options{})

	for _, want := range []string{
		`kill_delay = "1s"`,
		"send_interrupt = true",
		"Unix-only",
		"internal/devserver",
	} {
		if !strings.Contains(cfg, want) {
			t.Errorf("the air config is missing %q", want)
		}
	}
}

// A project scaffolded before this gets it too: internal/devserver travels
// with the framework's files, but main.go is never rewritten on upgrade, so
// without the repair the package would arrive and nothing would call it.
func TestDevServerRepairRewritesAnOlderMain(t *testing.T) {
	out, fixed, warnings := repairDevServerListenSource(olderMainGo, "demo/apps/api")
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(fixed) == 0 {
		t.Fatal("the repair reported nothing")
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("the repaired main.go is not valid Go: %v\n%s", err, out)
	}
	for _, want := range []string{
		`"demo/apps/api/internal/devserver"`,
		"devserver.Listen(cfg.Port, cfg.AppEnv)",
		"devserver.ExplainBindFailure(",
		"devserver.MountStepAside(router, cfg.AppEnv,",
		"srv.Serve(listener)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired main.go is missing %q", want)
		}
	}
	if strings.Contains(out, "srv.ListenAndServe()") {
		t.Error("the repaired main.go still binds at ListenAndServe")
	}
	// The channel moved up; it must not now be declared twice.
	if n := strings.Count(out, "quit := make(chan os.Signal, 1)"); n != 1 {
		t.Errorf("quit is declared %d times", n)
	}

	// Twice is a no-op: an upgrade runs on an upgraded project more often
	// than on one that has never been upgraded.
	again, fixedAgain, _ := repairDevServerListenSource(out, "demo/apps/api")
	if again != out || len(fixedAgain) != 0 {
		t.Error("the second run changed the file again")
	}
}

// A main.go that has moved on keeps its main.go, and is told why it was left.
func TestDevServerRepairLeavesAChangedMainAlone(t *testing.T) {
	src := `package main

func main() {
	srv.ListenAndServe()
}
`
	out, fixed, warnings := repairDevServerListenSource(src, "demo/apps/api")
	if out != src {
		t.Error("the file was rewritten from anchors it does not have")
	}
	if len(fixed) != 0 {
		t.Errorf("it claimed a fix: %v", fixed)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "grit#110") {
		t.Errorf("expected one warning pointing at the issue, got %v", warnings)
	}
}

// And a project that has tuned its own reload keeps its own numbers.
func TestAirRepairOnlyTouchesTheDefaults(t *testing.T) {
	tuned := "[build]\n  kill_delay = \"3s\"\n  send_interrupt = true\n"
	out, fixed, _ := repairAirKillSettingsSource(tuned)
	if out != tuned || len(fixed) != 0 {
		t.Error("a hand-tuned .air.toml was overwritten")
	}

	stock := "[build]\n  kill_delay = \"0s\"\n  send_interrupt = false\n"
	out, fixed, _ = repairAirKillSettingsSource(stock)
	if len(fixed) == 0 {
		t.Fatal("the stock settings were left alone")
	}
	if !strings.Contains(out, `kill_delay = "1s"`) || !strings.Contains(out, "send_interrupt = true") {
		t.Errorf("the repaired config is wrong:\n%s", out)
	}
}

// olderMainGo is the shape main.go had before v3.411.19, cut down to the
// anchors the repair works from.
const olderMainGo = `package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"demo/apps/api/internal/config"
	"demo/apps/api/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	router := routes.Setup(db, cfg)

	// Create server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
`
