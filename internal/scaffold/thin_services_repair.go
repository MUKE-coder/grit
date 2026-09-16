package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/codefmt"
)

// Contact-app review M29 and M30, in the API files upgrade does not deliver
// whole.
//
// M29: the ticket system had no service at all, so every query and every rule
// it has lived in handlers/ticket.go and could not be called from anywhere that
// is not an HTTP request. services/ticket.go is new and arrives whole;
// handlers/ticket.go is rewritten to call it.
//
// M30: mail left the auth handlers from bare goroutines. handlers/auth.go and
// routes.go are files people edit, so those two are patched by anchor; the two
// auth flow files are framework-owned and arrive whole.
//
// Every write goes through the manifest guard, so a file somebody edited is
// left alone and reported rather than clobbered.

// repairThinServices applies M29 and M30 to an existing project.
func repairThinServices(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)

	// The ticket handler and service are written by refreshSupportAndDashboards,
	// which rewrites handlers/ticket.go whole and runs earlier in the upgrade.
	repairQueuedAuthMail(apiRoot, opts)
	repairServiceWiring(apiRoot)
	return nil
}

// ─── M30: the auth emails go on the queue ────────────────────────────────────

const (
	authHandlerMailerField = "\tMailer *mail.Mailer\n}"
	authHandlerJobsField   = `	Mailer *mail.Mailer
	// Jobs is optional too: it is nil when the project runs without Redis.
	// When it is there, every email this handler sends goes onto the queue
	// instead of a goroutine, so it is retried and survives a restart.
	Jobs *jobs.Client
}`

	verifyGoroutineCall = "\tgo h.deliverVerificationEmail(user)"
	verifyDirectCall    = "\th.deliverVerificationEmail(c.Request.Context(), user)"

	verifyFuncOld = "func (h *AuthHandler) deliverVerificationEmail(user models.User) {"
	verifyFuncNew = "func (h *AuthHandler) deliverVerificationEmail(ctx context.Context, user models.User) {"

	resetGoroutineCall = "\tgo h.deliverPasswordReset(user, c.ClientIP())"
	resetDetachCall    = `	clientIP := c.ClientIP()
	detach(func(ctx context.Context) { h.deliverPasswordReset(ctx, user, clientIP) })`

	resetFuncOld = "func (h *AuthHandler) deliverPasswordReset(user models.User, clientIP string) {"
	resetFuncNew = "func (h *AuthHandler) deliverPasswordReset(ctx context.Context, user models.User, clientIP string) {"
)

// repairQueuedAuthMail stops the auth handlers sending mail from goroutines.
//
// Three files move together. auth.go holds the handler struct, which gains the
// Jobs field, and one of the calls; the verification and reset files hold the
// other calls and the functions they call, which now read h.Jobs. Half of this
// change does not compile, so nothing is written unless all of it can be.
//
// auth.go is written exactly once. The manifest guard compares a file with what
// Grit last recorded for it, and a second write in the same upgrade reads as an
// edit and is refused, which is how an earlier version of this repair left the
// struct without its field and the project without a build.
func repairQueuedAuthMail(apiRoot string, opts Options) {
	auth := filepath.Join(apiRoot, "internal", "handlers", "auth.go")
	verify := filepath.Join(apiRoot, "internal", "handlers", "auth_email_verification.go")
	reset := filepath.Join(apiRoot, "internal", "handlers", "auth_password_reset.go")
	if !fileExists(auth) {
		return
	}
	authSrc, err := readText(auth)
	if err != nil {
		fmt.Printf("  ⚠ reading internal/handlers/auth.go: %v\n", err)
		return
	}

	// The field first: without it neither flow can be queued.
	newAuth := authSrc
	if !strings.Contains(newAuth, "Jobs *jobs.Client") {
		if strings.Count(newAuth, authHandlerMailerField) != 1 {
			fmt.Printf("  ⚠ AuthHandler is not the struct Grit wrote, so its emails are still sent from goroutines: give it a Jobs *jobs.Client field set from svc.Jobs, and send with dispatchMail\n")
			return
		}
		newAuth = strings.Replace(newAuth, authHandlerMailerField, authHandlerJobsField, 1)
		newAuth = addGoImport(newAuth, opts.Module()+"/internal/jobs")
	}

	// The verification flow.
	var verifySrc, newVerify string
	if fileExists(verify) {
		if verifySrc, err = readText(verify); err != nil {
			fmt.Printf("  ⚠ reading internal/handlers/auth_email_verification.go: %v\n", err)
			return
		}
		newVerify = verifySrc
		if !strings.Contains(verifySrc, verifyFuncNew) {
			if !strings.Contains(verifySrc, verifyFuncOld) {
				fmt.Printf("  ⚠ the verification email is not sent the way Grit wrote it: send it with dispatchMail rather than from a `go` statement, so it is queued, retried and survives a restart\n")
				return
			}
			newAuth = strings.ReplaceAll(newAuth, verifyGoroutineCall, verifyDirectCall)
			newVerify = strings.ReplaceAll(newVerify, verifyGoroutineCall, verifyDirectCall)
			newVerify = strings.Replace(newVerify, verifyFuncOld, verifyFuncNew, 1)
			newVerify = replaceMailerSend(newVerify, `"verify:"+user.ID+":"+token`)
		}
	}

	// The reset flow.
	var resetSrc, newReset string
	if fileExists(reset) {
		if resetSrc, err = readText(reset); err != nil {
			fmt.Printf("  ⚠ reading internal/handlers/auth_password_reset.go: %v\n", err)
			return
		}
		newReset = resetSrc
		if !strings.Contains(resetSrc, resetFuncNew) {
			if !strings.Contains(resetSrc, resetFuncOld) || !strings.Contains(resetSrc, resetGoroutineCall) {
				fmt.Printf("  ⚠ the password-reset email is not sent the way Grit wrote it: run the delivery through detach rather than a bare `go` statement, and send it with dispatchMail so it is queued and retried\n")
				return
			}
			newReset = strings.Replace(newReset, resetGoroutineCall, resetDetachCall, 1)
			newReset = strings.Replace(newReset, resetFuncOld, resetFuncNew, 1)
			newReset = replaceMailerSend(newReset, `"password-reset:"+user.ID+":"+token`)
		}
	}

	if newAuth == authSrc && newVerify == verifySrc && newReset == resetSrc {
		return
	}
	writes := []pendingWrite{{auth, newAuth, authSrc}}
	if newVerify != verifySrc {
		writes = append(writes, pendingWrite{verify, newVerify, verifySrc})
	}
	if newReset != resetSrc {
		writes = append(writes, pendingWrite{reset, newReset, resetSrc})
	}
	if writeAll(writes) {
		fmt.Printf("  ✓ the auth emails are queued rather than sent from goroutines, and the reset delivery is capped\n")
	}
}

// mailerSendOpen is the block that sent the mail inline, from its own context,
// inside the goroutine that called it.
const mailerSendOpen = `	if h.Mailer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := h.Mailer.Send(ctx, mail.SendOptions{`

// replaceMailerSend swaps that inline send for a queued one.
//
// Only the two lines that open the block change: the recipient, the subject,
// the template and the data in between are untouched, so a project that has
// rewritten the copy of its own emails keeps it. key is the Go expression that
// deduplicates the enqueue.
func replaceMailerSend(src, key string) string {
	if !strings.Contains(src, mailerSendOpen) {
		return src
	}
	replacement := "\tif h.Jobs != nil || h.Mailer != nil {\n" +
		"\t\tif err := dispatchMail(ctx, h.Mailer, h.Jobs, " + key + ", mail.SendOptions{"
	return strings.Replace(src, mailerSendOpen, replacement, 1)
}

// ─── wiring ──────────────────────────────────────────────────────────────────

// repairServiceWiring adds the middleware and the two handler fields routes.go
// has to set for any of the above to take effect.
func repairServiceWiring(apiRoot string) {
	routes := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	src, err := readText(routes)
	if err != nil {
		return
	}
	out := src
	var fixed []string

	if !strings.Contains(out, "middleware.RequestMeta()") && strings.Contains(out, "r.Use(middleware.RequestID())") {
		out = strings.Replace(out, "r.Use(middleware.RequestID())",
			"r.Use(middleware.RequestID())\n\t// The client IP, the user agent and the request id, on the request's\n"+
				"\t// context, so a service can read them without taking a *gin.Context.\n"+
				"\tr.Use(middleware.RequestMeta())", 1)
		fixed = append(fixed, "services read the request's IP, user agent and id from its context")
	}

	const ticketOld = "&handlers.TicketHandler{DB: db, Mail: svc.Mailer}"
	if strings.Contains(out, ticketOld) {
		out = strings.Replace(out, ticketOld, "&handlers.TicketHandler{DB: db, Mail: svc.Mailer, Jobs: svc.Jobs}", 1)
		fixed = append(fixed, "the new-ticket email goes on the job queue")
	}

	const authOld = "\t\tMailer:      svc.Mailer,\n\t}"
	if strings.Contains(out, "handlers.AuthHandler{") && strings.Count(out, authOld) == 1 &&
		!strings.Contains(out, "Jobs:        svc.Jobs") {
		out = strings.Replace(out, authOld, "\t\tMailer:      svc.Mailer,\n\t\tJobs:        svc.Jobs,\n\t}", 1)
		fixed = append(fixed, "the auth emails go on the job queue")
	}

	if out == src {
		return
	}
	if wrote, err := writeGuarded(routes, "", out); err != nil {
		fmt.Printf("  ⚠ internal/routes/routes.go: %v\n", err)
	} else if wrote {
		fmt.Printf("  ✓ internal/routes/routes.go: %s\n", strings.Join(fixed, "; "))
	} else {
		fmt.Printf("  ⚠ internal/routes/routes.go is yours: add r.Use(middleware.RequestMeta()) after RequestID, and Jobs: svc.Jobs to the auth and ticket handlers\n")
	}
}

// ─── small helpers ───────────────────────────────────────────────────────────

// readText reads a file with its line endings normalised.
func readText(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(raw), "\r\n", "\n"), nil
}

// writeGuarded writes content through the manifest guard, formatted the way the
// scaffold would have written it, with the module placeholder filled in.
// Reports whether it wrote.
//
// module is passed rather than assumed empty: a template written raw leaves
// {{MODULE}} in its own import paths, and the project then fails to build with
// "invalid import path: {{MODULE}}/internal/mail". Pass "" for a file that
// carries no placeholder.
func writeGuarded(path, module, content string) (bool, error) {
	if module != "" {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
	}
	return guardedWrite(path, codefmt.File(path, content))
}

// pendingWrite is one file of a change that has to land whole.
type pendingWrite struct {
	path, content, previous string
}

// writeAll writes every file or none: part of a change that spans several
// files does not compile. When one is refused, the ones already written are put
// back the way they were.
func writeAll(writes []pendingWrite) bool {
	for i, w := range writes {
		wrote, err := writeGuarded(w.path, "", w.content)
		if err == nil && wrote {
			continue
		}
		for _, done := range writes[:i] {
			if _, revertErr := writeGuarded(done.path, "", done.previous); revertErr != nil {
				fmt.Printf("  ⚠ %s: %v\n", filepath.Base(done.path), revertErr)
			}
		}
		if err != nil {
			fmt.Printf("  ⚠ %s: %v\n", filepath.Base(w.path), err)
		} else {
			fmt.Printf("  ⚠ %s is yours, so the auth emails are still sent from goroutines: give AuthHandler a Jobs field and send with dispatchMail\n", filepath.Base(w.path))
		}
		return false
	}
	return true
}

// addGoImport adds one import path to a Go file's first import block.
//
// Built as a line inserted before an existing import rather than as a string
// replacement on an anchor line: a Replace that matches nothing fails silently,
// and the file then does not compile for a reason nothing names.
func addGoImport(src, path string) string {
	quoted := "\t\"" + path + "\"\n"
	if strings.Contains(src, quoted) {
		return src
	}
	open := strings.Index(src, "\nimport (\n")
	if open < 0 {
		return src
	}
	close := strings.Index(src[open:], "\n)\n")
	if close < 0 {
		return src
	}
	block := src[open+len("\nimport (\n") : open+close+1]
	lines := strings.Split(block, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "\t\"") && line > quoted[:len(line)] {
			lines = append(lines[:i], append([]string{strings.TrimRight(quoted, "\n")}, lines[i:]...)...)
			return src[:open+len("\nimport (\n")] + strings.Join(lines, "\n") + src[open+close+1:]
		}
	}
	return src[:open+len("\nimport (\n")] + block + quoted + src[open+close+1:]
}
