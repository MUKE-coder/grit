package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// A sentinel error compared with == is a check that works until somebody wraps
// the error.
//
// "err == gorm.ErrRecordNotFound" was how twenty-three places in the generated
// API decided whether a row was missing, and every one of them was one
// fmt.Errorf("loading user: %w", err) away from turning a 404 into a 500. The
// wrap is the normal thing to do when adding context to an error, which is what
// makes this a trap rather than a preference: the code that breaks it is the
// code a careful developer writes.
//
// errors.Is walks the wrap chain, so it keeps answering the same question after
// somebody adds context. The rewrite is mechanical and it is applied to the
// whole of apps/api, including code the project wrote itself, because the
// transformation preserves meaning exactly: errors.Is(x, target) is true
// wherever x == target was.

var (
	// A comparison against an exported package-level sentinel: err,
	// result.Error and friends against gorm.ErrRecordNotFound,
	// services.ErrUserNotFound, http.ErrServerClosed, redis.Nil and the
	// same-package ErrX a test in that package compares against.
	sentinelEqRe = regexp.MustCompile(`\b((?:[a-zA-Z_]\w*)(?:\.[A-Z]\w*)?)\s+(==|!=)\s+((?:[a-z]\w*\.)?Err[A-Z]\w*|redis\.Nil)\b`)

	// The left-hand side has to be something that holds an error. Matching any
	// identifier would rewrite "name != models.ErrorLabel" into nonsense.
	sentinelLHS = map[string]bool{
		"err": true, "e": true, "cerr": true, "derr": true, "ferr": true,
		"serr": true, "terr": true, "txErr": true, "rerr": true, "werr": true,
		"result.Error": true, "res.Error": true, "tx.Error": true, "db.Error": true,
	}
)

// repairSentinelErrors rewrites == and != against a sentinel error as
// errors.Is across every Go file under apps/api.
func repairSentinelErrors(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	paths := []string{
		// cmd/server/main.go has the one non-database case: the
		// http.ErrServerClosed that graceful shutdown checks for. A
		// single-file project keeps it at the module root instead.
		filepath.Join(apiRoot, "cmd", "server", "main.go"),
		filepath.Join(apiRoot, "main.go"),
	}
	// Every package under internal, walked rather than listed. A list would
	// have to be kept in step with each new package, and the one left out is
	// the one that keeps the defect.
	internalRoot := filepath.Join(apiRoot, "internal")
	entries, derr := os.ReadDir(internalRoot)
	if derr != nil && !os.IsNotExist(derr) {
		return fmt.Errorf("reading %s: %w", internalRoot, derr)
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		files, ferr := goSourcesIncludingTests(filepath.Join(internalRoot, e.Name()))
		if ferr != nil {
			return ferr
		}
		paths = append(paths, files...)
	}
	for _, path := range paths {
		if !fileExists(path) {
			continue
		}
		if err := repairSourceFile(root, m, path, repairSentinelErrorSource); err != nil {
			return err
		}
	}
	return nil
}

// repairSentinelErrorSource rewrites one file's sentinel comparisons.
func repairSentinelErrorSource(src string) (string, []string, []string) {
	out, n := rewriteErrorSwitches(src)
	out = sentinelEqRe.ReplaceAllStringFunc(out, func(match string) string {
		parts := sentinelEqRe.FindStringSubmatch(match)
		lhs, op, sentinel := parts[1], parts[2], parts[3]
		if !sentinelLHS[lhs] {
			return match
		}
		n++
		call := "errors.Is(" + lhs + ", " + sentinel + ")"
		if op == "!=" {
			return "!" + call
		}
		return call
	})
	if n == 0 {
		return src, nil, nil
	}
	withImport, ok := addImportGroup(out, "errors")
	if !ok {
		return src, nil, []string{`a sentinel error is compared with == here and the file has no import block to add "errors" to: compare with errors.Is by hand, or a wrapped error turns a 404 into a 500`}
	}
	return withImport, []string{fmt.Sprintf("%d sentinel %s compared with errors.Is, so a wrapped error still reads as itself",
		n, plural(n, "error is", "errors are"))}, nil
}

// goSourcesIncludingTests is goSources with the _test.go files kept.
//
// Tests are linted too, and the comparison in a test is the same mistake: a
// test that asserts err == ErrPasskeyChallengeGone starts failing the moment
// the code it tests wraps that error, and the failure points at the test
// rather than at the wrap.
func goSourcesIncludingTests(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		out = append(out, filepath.Join(dir, e.Name()))
	}
	return out, nil
}

// rewriteErrorSwitches turns "switch err { case X:" into
// "switch { case errors.Is(err, X):", and reports how many cases it changed.
//
// A switch on a value compares with ==, so this is the same defect wearing
// different punctuation, and golangci-lint's errorlint reports it as such.
// Only "switch err {" is touched, and only cases that name a single sentinel:
// anything else is left exactly as it is.
func rewriteErrorSwitches(src string) (string, int) {
	lines := strings.Split(src, "\n")
	changed := 0
	for i := 0; i < len(lines); i++ {
		indent, ok := switchOnErrIndent(lines[i])
		if !ok {
			continue
		}
		body, end := errorSwitchCases(lines, i+1, indent)
		if end < 0 || len(body) == 0 {
			continue
		}
		lines[i] = indent + "switch {"
		for _, at := range body {
			sentinel := strings.TrimSuffix(strings.TrimSpace(lines[at])[len("case "):], ":")
			lines[at] = indent + "case errors.Is(err, " + sentinel + "):"
			changed++
		}
		i = end
	}
	if changed == 0 {
		return src, 0
	}
	return strings.Join(lines, "\n"), changed
}

// switchOnErrIndent reports whether a line opens a switch on err, and its indent.
func switchOnErrIndent(line string) (string, bool) {
	trimmed := strings.TrimLeft(line, "\t")
	if trimmed != "switch err {" {
		return "", false
	}
	return line[:len(line)-len(trimmed)], true
}

// errorSwitchCases lists the case lines of the switch that opened at indent,
// and the line its closing brace is on. It reports -1 for a switch holding a
// case this cannot rewrite, so a mixed switch is left whole.
func errorSwitchCases(lines []string, from int, indent string) ([]int, int) {
	var cases []int
	for i := from; i < len(lines); i++ {
		line := lines[i]
		if line == indent+"}" {
			return cases, i
		}
		if !strings.HasPrefix(line, indent+"case ") {
			continue
		}
		expr := strings.TrimSuffix(strings.TrimPrefix(line, indent+"case "), ":")
		if !sentinelCaseRe.MatchString(expr) {
			return nil, -1
		}
		cases = append(cases, i)
	}
	return nil, -1
}

// sentinelCaseRe is one sentinel and nothing else: no comma list, no call.
var sentinelCaseRe = regexp.MustCompile(`^(?:[a-z]\w*\.)?Err[A-Z]\w*$`)
