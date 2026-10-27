package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairDevServerListen brings an existing project's cmd/server/main.go up to
// binding the port before anything else.
//
// internal/devserver travels with the framework's own files, so an upgraded
// project gets the package. main.go is not rewritten on upgrade, because it is
// the one file people reasonably edit, so without this the package would
// arrive and nothing would call it: the fix would ship and not apply, which is
// worse than not shipping it, because the changelog would say it was fixed.
//
// Anchored on text Grit wrote, and every step checked before the next. A
// project whose main.go has moved on keeps its main.go and is told.
func repairDevServerListen(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), "cmd", "server", "main.go")
	if !fileExists(path) {
		// A single project keeps main.go at the module root.
		path = filepath.Join(opts.APIRoot(root), "main.go")
		if !fileExists(path) {
			return nil
		}
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	module := opts.Module()
	return repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
		return repairDevServerListenSource(src, module)
	})
}

func repairDevServerListenSource(src, module string) (string, []string, []string) {
	if strings.Contains(src, "devserver.Listen(") {
		return src, nil, nil
	}
	// Nothing to do for a main.go that never served HTTP this way.
	if !strings.Contains(src, "srv.ListenAndServe()") {
		return src, nil, nil
	}

	const loadConfig = `	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
`
	const createServer = `	// Create server
	srv := &http.Server{`
	const gracefulShutdown = `	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit`

	for _, anchor := range []string{loadConfig, createServer, gracefulShutdown} {
		if !strings.Contains(src, anchor) {
			return src, nil, []string{"cmd/server/main.go has moved on from what Grit wrote, " +
				"so the port is still bound at the end of startup: a port already in use is " +
				"reported after the route dump, and on a hot reload the previous build keeps " +
				"answering. Call devserver.Listen before the database by hand, or see grit#110"}
		}
	}

	out := strings.Replace(src, loadConfig, loadConfig+`
	// The port, before anything else.
	//
	// Bound here rather than at ListenAndServe so that a port already in use
	// is the FIRST thing printed. It used to be the last, one line after the
	// route dump, and on a hot reload that line scrolls past while the
	// previous build keeps answering with the code from before the edit
	// (grit#110).
	listener, err := devserver.Listen(cfg.Port, cfg.AppEnv)
	if err != nil {
		log.Fatal(devserver.ExplainBindFailure(cfg.Port, cfg.AppEnv, err))
	}
`, 1)

	// The shutdown channel moves above the server, so the step-aside endpoint
	// can use the same path Ctrl+C does.
	out = strings.Replace(out, gracefulShutdown, `	// Graceful shutdown
	<-quit`, 1)
	out = strings.Replace(out, createServer, `	// Shutdown signalling, set up before the server so the development
	// step-aside endpoint can use the same path Ctrl+C does.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Development only: let the next build of this server ask this one to go,
	// so a rebuild does not end with two servers and the older one winning.
	devserver.MountStepAside(router, cfg.AppEnv, func() {
		quit <- syscall.SIGTERM
	})

`+createServer, 1)

	out = strings.Replace(out, "srv.ListenAndServe()", "srv.Serve(listener)", 1)

	withImport, ok := addImportGroup(out, module+"/internal/devserver")
	if !ok {
		return src, nil, []string{"cmd/server/main.go has no import block to add internal/devserver to"}
	}
	return withImport, []string{"the port is bound before the database and the routes, so a port " +
		"already in use is the first thing you read, and in development a previous build that " +
		"should have exited is asked for it back"}, nil
}

// repairAirKillSettings updates .air.toml, for what little air can do here.
//
// kill_delay and send_interrupt are a Unix-only pair: air reads kill_delay in
// its Unix kill path and only when send_interrupt is on, and on Windows it
// logs that send_interrupt is unsupported and runs TASKKILL. So this is a
// graceful shutdown on macOS and Linux and nothing at all on Windows, which is
// why the rest of grit#110 is in the server.
//
// Only when both lines are exactly what Grit wrote. A project that has tuned
// its own reload keeps its own numbers.
func repairAirKillSettings(root string, opts Options) error {
	path := filepath.Join(opts.APIRoot(root), ".air.toml")
	if !fileExists(path) {
		return nil
	}
	// Not through repairSourceFile: that one runs the Go formatter over its
	// output and, when the result does not parse, says so and leaves the file
	// alone. For a TOML file it never parses, so the repair would report
	// nothing and do nothing, every time, in silence.
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %s: %w", path, err)
	}
	crlf := strings.Contains(string(raw), "\r\n")
	out, fixed, _ := repairAirKillSettingsSource(strings.ReplaceAll(string(raw), "\r\n", "\n"))
	if len(fixed) == 0 {
		return nil
	}
	if crlf {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	if rel, rerr := filepath.Rel(root, path); rerr == nil {
		fmt.Printf("  ✓ %s: %s\n", filepath.ToSlash(rel), fixed[0])
	}
	return nil
}

func repairAirKillSettingsSource(src string) (string, []string, []string) {
	const old = `  kill_delay = "0s"
  send_interrupt = false`
	if !strings.Contains(src, old) {
		return src, nil, nil
	}
	out := strings.Replace(src, old, `  # Stop the old server before starting the new one, and let it shut down
  # rather than being shot. Unix-only: air reads kill_delay in its Unix kill
  # path and only when send_interrupt is on, and on Windows it logs
  # "send_interrupt is not supported" and runs TASKKILL /F /T. The Windows half
  # of this is handled in internal/devserver (grit#110).
  kill_delay = "1s"
  send_interrupt = true`, 1)
	return out, []string{fmt.Sprintf("a reload now interrupts the old server instead of killing it, " +
		"so it closes its listener on the way out")}, nil
}
