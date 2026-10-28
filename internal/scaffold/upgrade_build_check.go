package scaffold

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
)

/*
 * Does the project still compile?
 *
 * An upgrade writes framework files wholesale and patches the project's own
 * files with anchored repairs. config.go is the one that matters here: it is
 * not in the file map, so every field added to it since the project was
 * scaffolded arrives through a repair, and a repair whose anchor is missing
 * warns and moves on. The file that USES the field is written anyway.
 *
 * Upgrading the WhatsApp blueprint from 3.306.0 produced four of these at
 * once, from three different repairs:
 *
 *	cfg.SecurityAutoBanEnabled undefined
 *	cfg.SentinelAutoBlockExemptIPs undefined
 *	cfg.SentinelLockoutExemptIPs undefined
 *	undefined: services
 *
 * Each of those is now fixed where it was caused. This exists because the next
 * one will not be: a repair anchor is a guess about a file somebody else owns,
 * and guesses go stale. `go build ./...` is the one check that catches the
 * whole class, so the upgrade ends by running it. A broken build is loud now
 * instead of being discovered by the next command somebody types.
 *
 * It does not try to fix anything. Telling the truth at the end of an upgrade
 * is the entire job.
 */

// buildCheckTimeout is generous: a cold build of a large API after an upgrade
// that added modules compiles a lot. Exceeding it is reported as "could not
// check", never as a failure, because a slow machine is not a broken project.
const buildCheckTimeout = 4 * time.Minute

// tidyModules resolves go.sum for whatever the upgrade added to go.mod.
//
// Reported as "could not check" when it fails rather than as a broken build,
// because the usual reason is no network, and an offline machine says nothing
// about the project's code.
func tidyModules(apiDir string) bool {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = apiDir

	done := make(chan struct{})
	var err error
	go func() {
		_, err = cmd.CombinedOutput()
		close(done)
	}()

	select {
	case <-done:
		return err == nil
	case <-time.After(buildCheckTimeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return false
	}
}

// ReportBuild runs the check against a project root, for the caller that owns
// the last word.
//
// `grit upgrade` ends in the cobra command, not in Upgrade(): plugins are
// updated after it returns, and a plugin rewrites routes.go. Running the check
// inside Upgrade() therefore reported the storefront as compiling and the
// Stripe plugin then gave routes.go a second `paymentService :=`.
func ReportBuild(root string) bool {
	opts := Options{Architecture: detectArchitecture(root)}
	return reportBuildAfterUpgrade(opts.APIRoot(root))
}

// reportBuildAfterUpgrade compiles the API and says what broke.
//
// Returns true when the project builds, or when the check could not be run at
// all (no Go toolchain, no API directory, timed out). Never returns an error:
// the upgrade has already happened, and failing it here would help nobody.
func reportBuildAfterUpgrade(apiDir string) bool {
	if !dirExists(apiDir) {
		return true
	}
	if _, err := exec.LookPath("go"); err != nil {
		return true
	}

	// go.sum first. An upgrade adds modules to go.mod, so without this every
	// run that touched a dependency reports a wall of "missing go.sum entry"
	// as though the project were broken. It is not: it is the step the Next
	// steps block names, and a check that cries wolf is one nobody reads.
	if !tidyModules(apiDir) {
		color.New(color.FgHiBlack).Printf("  … could not run `go mod tidy`, so the build was not checked\n")
		color.New(color.FgHiCyan).Printf("    Run `cd %s && go mod tidy` yourself, then `go build ./...`\n", apiDir)
		return true
	}

	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = apiDir

	done := make(chan struct{})
	var out []byte
	var err error
	go func() {
		out, err = cmd.CombinedOutput()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(buildCheckTimeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		color.New(color.FgHiBlack).Printf("  … the build check timed out; run `go build ./...` yourself\n")
		return true
	}

	if err == nil {
		color.New(color.FgHiGreen).Printf("  ✓ go.sum is current and the API still compiles\n")
		return true
	}

	red := color.New(color.FgHiRed)
	cyan := color.New(color.FgHiCyan)
	fmt.Println()
	red.Printf("  ✗ The API does not compile after this upgrade.\n")
	fmt.Println()
	for _, line := range buildErrorLines(string(out), 12) {
		fmt.Printf("      %s\n", line)
	}
	fmt.Println()
	cyan.Printf("    Almost always a field a repair could not add, because the file it\n")
	cyan.Printf("    anchors on is one you have edited. The ⚠ lines above say which, and\n")
	cyan.Printf("    what to add by hand. If none of them look related, this is a Grit\n")
	cyan.Printf("    bug: https://github.com/MUKE-coder/grit/issues\n")
	cyan.Printf("    Your previous files are in .grit/backups.\n")
	fmt.Println()
	return false
}

// buildErrorLines picks the lines worth showing: the compiler's own, without
// the package headers, capped so a project with fifty errors prints twelve and
// a count rather than a wall.
func buildErrorLines(out string, limit int) []string {
	var lines []string
	total := 0
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "# ") {
			continue
		}
		total++
		if len(lines) < limit {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	if total > len(lines) {
		lines = append(lines, fmt.Sprintf("… and %d more", total-len(lines)))
	}
	return lines
}
