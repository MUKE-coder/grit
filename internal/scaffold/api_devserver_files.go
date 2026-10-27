package scaffold

// On Windows, a rebuild could leave the previous server running and still
// answering, with the new one dead.
//
// air starts the binary through PowerShell, so the process it tracks is
// PowerShell and not the server. It stops the old one with
// TASKKILL /F /T on that PowerShell pid, which walks the process tree: when
// PowerShell has already gone, the server underneath it is an orphan with no
// parent to be found through, and TASKKILL misses it. air's kill_delay is no
// help, whatever the docs imply -- it is read only by the Unix kill path, and
// only when send_interrupt is on, which Windows does not support either. air
// then rebuilds happily, because Windows lets the linker replace the image of
// a running process, and the new server binds a port the orphan still holds:
//
//	Server failed: listen tcp :8080: bind: Only one usage of each socket
//	address (protocol/network address/port) is normally permitted.
//
// One line, after the route dump, and then the old code answers every request
// as though nothing happened. The developer's edit is simply not live, and the
// only symptom is that it does not work (grit#110).
//
// Two things fix it, both in the server rather than in air, because air cannot
// do this on Windows:
//
//  1. The port is bound before anything else runs, so a failure is the first
//     thing printed and not the twentieth, and it says what is actually wrong.
//
//  2. In development the new server asks whoever holds the port to step
//     aside, over HTTP. Nothing is killed and no pid is guessed: a stale Grit
//     server recognises the request and shuts down gracefully, and anything
//     else answers 404 and is left alone.

// apiDevServerGo is internal/devserver/devserver.go in a generated project.
func apiDevServerGo() string {
	return `// Package devserver holds the parts of starting up that only matter while
// you are developing: taking the port back from a server that should have
// exited, and explaining it when that is not possible.
package devserver

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// StepAsidePath is the endpoint a newer instance calls on an older one.
//
// Under /api/dev/ so it reads as what it is, and mounted only in development.
const StepAsidePath = "/api/dev/step-aside"

// stepAsideToken is sent and checked so that an unrelated service on this port
// cannot be mistaken for one of ours. It is not a secret and does not need to
// be: the only thing it protects against is ambiguity.
const stepAsideToken = "grit-dev-step-aside"

// Listen binds the port the server will serve on.
//
// Called before the database, the routes and everything else, so that a port
// which is already taken is the first thing reported rather than one line
// after two hundred lines of route registration.
//
// In development, a port held by an older instance of this same server is
// taken back: see MountStepAside. In any other environment the error is
// returned as it came, because a port conflict in production is a deployment
// problem and quietly stopping the other process is the last thing anyone
// wants.
func Listen(port, env string) (net.Listener, error) {
	addr := ":" + port

	ln, err := net.Listen("tcp", addr)
	if err == nil || !isAddrInUse(err) || !isDev(env) {
		return ln, err
	}

	if !askToStepAside(port) {
		return nil, err
	}

	// It said yes. Give it a moment to finish shutting down, then take the
	// port. Graceful shutdown drains in-flight requests, so this is not
	// instant.
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(200 * time.Millisecond)
		if ln, err = net.Listen("tcp", addr); err == nil {
			fmt.Printf("devserver: took port %s back from the previous server, which had not exited\n", port)
			return ln, nil
		}
		if !isAddrInUse(err) {
			return nil, err
		}
	}
	return nil, err
}

// MountStepAside registers the endpoint that lets a newer instance of this
// server ask this one to exit. Development only, and localhost only.
//
// An HTTP request rather than a signal or a kill, because this is the one
// mechanism that cannot hurt an innocent process: a server that is not ours
// answers 404 and keeps running. Working out which pid holds a port and
// killing it means being wrong occasionally, and being wrong means killing
// something somebody else was using.
//
// stop should begin the same graceful shutdown that Ctrl+C does.
func MountStepAside(router *gin.Engine, env string, stop func()) {
	if !isDev(env) || router == nil || stop == nil {
		return
	}
	router.POST(StepAsidePath, func(c *gin.Context) {
		if !isLoopback(c.ClientIP()) {
			c.Status(http.StatusNotFound)
			return
		}
		if c.GetHeader("X-Grit-Dev") != stepAsideToken {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, gin.H{"stepping_aside": true})
		// After the response, so the caller knows the answer was yes before
		// the connection goes away with the server.
		go func() {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("devserver: a newer build asked for the port, shutting down")
			stop()
		}()
	})
}

// askToStepAside asks whoever is on the port to leave, and reports whether
// they agreed.
func askToStepAside(port string) bool {
	url := "http://127.0.0.1:" + port + StepAsidePath
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("X-Grit-Dev", stepAsideToken)

	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// ExplainBindFailure is what the developer reads when the port could not be
// taken, which on a rebuild means their change is not live.
//
// Long on purpose. The one-line version of this message was
// indistinguishable from a route that had not been registered, and cost
// somebody an afternoon.
func ExplainBindFailure(port, env string, err error) string {
	if !isAddrInUse(err) {
		return fmt.Sprintf("Server failed to start: %v", err)
	}

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(box("PORT " + port + " IS ALREADY IN USE"))
	b.WriteString("  Another process is serving on this port, so this server did not\n")
	b.WriteString("  start. If you are running a hot reload, that other process is\n")
	b.WriteString("  probably the PREVIOUS build of this same server, which means:\n\n")
	b.WriteString("    your API still answers, with the code from before your edit.\n\n")
	if isDev(env) {
		b.WriteString("  It was asked to step aside and did not answer, so it is either\n")
		b.WriteString("  not a Grit server or it is stuck.\n\n")
	}
	b.WriteString("  Find it and stop it:\n\n")
	b.WriteString(findAndStopHint(port))
	b.WriteString("\n  Then save any .go file to rebuild.\n")
	return b.String()
}

// findAndStopHint is the two commands that actually answer "what do I type".
//
// Per platform, because the whole point is that the developer does not have to
// go and look it up while their edit is silently not running.
func findAndStopHint(port string) string {
	if runtime.GOOS == "windows" {
		return "    netstat -ano | findstr :" + port + "\n" +
			"    taskkill /F /PID <the pid in the last column>\n"
	}
	return "    lsof -nP -iTCP:" + port + " -sTCP:LISTEN\n" +
		"    kill <the pid>\n"
}

// box draws the heading, measured rather than padded by hand. The first
// version counted the spaces into the string literal and was two short, which
// nobody notices until a five-digit port makes it four short.
func box(title string) string {
	const width = 62
	line := strings.Repeat("─", width)
	inner := "  " + title
	if n := width - runeLen(inner); n > 0 {
		inner += strings.Repeat(" ", n)
	}
	return "  ┌" + line + "┐\n" +
		"  │" + inner + "│\n" +
		"  └" + line + "┘\n\n"
}

// runeLen counts characters, not bytes: the box is drawn in runes.
func runeLen(s string) int { return len([]rune(s)) }

func isDev(env string) bool {
	// Not the empty string. config.Load defaults APP_ENV to "production"
	// precisely so that an unset environment is the safe one, and this has to
	// agree with it: taking a port back is a development convenience and must
	// never be something an unconfigured deployment does.
	switch strings.ToLower(strings.TrimSpace(env)) {
	case "dev", "development", "local":
		return true
	}
	return false
}

func isLoopback(ip string) bool {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	return parsed != nil && parsed.IsLoopback()
}

// isAddrInUse reports the one error this package exists for, on every
// platform. Windows does not use EADDRINUSE, it uses WSAEADDRINUSE, and
// neither is reachable by name from portable code without the syscall package
// per GOOS, so the check goes through errors.Is on the portable sentinel and
// falls back to the message.
func isAddrInUse(err error) bool {
	if err == nil {
		return false
	}
	var se syscall.Errno
	if errors.As(err, &se) {
		// 48 is EADDRINUSE on BSD/macOS, 98 on Linux, 10048 is WSAEADDRINUSE.
		switch uintptr(se) {
		case 48, 98, 10048:
			return true
		}
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "address already in use") ||
		strings.Contains(msg, "only one usage of each socket address")
}
`
}
