package docscheck

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Go code samples on the docs site, checked against the rules the framework
// states everywhere else.
//
// A product review in October 2026 found the homepage hero running a query
// inside a handler, a float64 price, a dropped bind error and a hand-built
// error envelope. Every one of those is a thing the agent skill tells an agent
// never to do, and the homepage is the page an agent reads first: it was
// teaching the wrong pattern from the opening screen and had been for months.
//
// The commands in the docs are already checked against the real command tree.
// The code was not checked against anything, which is why the drift was found
// by a reader and not by the build.
//
// Deliberately narrow. Only the four rules Grit states in its own README, skill
// and stability page, each with an unambiguous signature, so the check cannot
// cry wolf and get switched off. A rule that needs judgement belongs in review,
// not here.

// SampleRule is one thing a Go sample on the docs site must not do.
type SampleRule struct {
	// ID is what a failure names, so the fix is searchable.
	ID string
	// Match finds the offence. Anchored at a line.
	Match *regexp.Regexp
	// Why the rule exists, and what to write instead.
	Message string
}

var sampleRules = []SampleRule{
	{
		ID: "handler-query",
		// h.DB.Where(, h.DB.Find(, h.DB.First( and friends. The receiver is the
		// tell: a handler reaching for the database at all.
		Match:   regexp.MustCompile(`\bh\.DB\.\s*(Where|Find|First|Create|Save|Delete|Model|Preload|Raw)\b`),
		Message: "a handler runs no queries: call the service (h.service().X(h.ctx(c), ...)). Stability says \"handlers run no queries\" and grit doctor checks it",
	},
	{
		ID:      "float-money",
		Match:   regexp.MustCompile(`(?i)\b(price|amount|total|cost|fee|balance|salary)\b\s+float(32|64)\b|\b(price|amount|total|cost|fee|balance|salary):float\b`),
		Message: "money, never float: use the money type (integer minor units plus a currency column), or price:money on the command line",
	},
	{
		ID:      "dropped-bind",
		Match:   regexp.MustCompile(`^\s*c\.Should(BindJSON|Bind|BindQuery|BindUri)\(`),
		Message: "a bind error is the 422 the client needs: if err := c.ShouldBindJSON(&req); err != nil { respond.ValidationError(c, err); return }",
	},
	{
		ID: "hand-built-error",
		// c.JSON(401, gin.H{"error": "some string"}): an envelope with no code,
		// which no client can branch on and which is not in the catalogue.
		//
		// Not the nested form. A page explaining the envelope is entitled to
		// write gin.H{"error": gin.H{"code": ..., "message": ...}} out in full,
		// and flagging that taught nothing and would have got the check
		// switched off inside a week.
		Match:   regexp.MustCompile(`c\.JSON\(\s*(4\d\d|5\d\d|http\.Status\w+)\s*,\s*gin\.H\{\s*"error"\s*:\s*"`),
		Message: "an error needs a code a client can branch on: respond.Error(c, err), or the full envelope {\"error\":{\"code\",\"message\"}}",
	},
}

// SampleProblem is one rule broken in one place.
type SampleProblem struct {
	File    string // relative to the docs root
	Line    int    // 1-indexed
	Rule    string
	Text    string // the offending line, trimmed
	Message string
}

func (p SampleProblem) String() string {
	return fmt.Sprintf("%s:%d [%s] %s\n      %s", p.File, p.Line, p.Rule, p.Message, p.Text)
}

// goSampleStart marks a Go code block: a CodeBlock with language="go", or a
// const NAME_CODE = `...` holding Go. Both shapes are on the site.
var (
	goLangProp  = regexp.MustCompile(`language=["']go["']`)
	goCodeConst = regexp.MustCompile("^const [A-Z_]*CODE[A-Z_]* = `|icon: 'go'")
)

// CheckSamples reads every docs page and reports the Go samples that break a
// rule the framework states elsewhere.
//
// The scan is line-based rather than parsed. A docs page is JSX wrapping
// template literals, so there is no Go file to parse, and the four signatures
// above are specific enough that context is not needed to recognise them.
func CheckSamples(root string) ([]SampleProblem, error) {
	var out []SampleProblem
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if skipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if ext := filepath.Ext(path); ext != ".tsx" && ext != ".ts" && ext != ".mdx" {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		out = append(out, checkSampleFile(filepath.ToSlash(rel), string(data))...)
		return nil
	})
	return out, err
}

// skipDir keeps the walk out of build output and dependencies.
func skipDir(name string) bool {
	switch name {
	case "node_modules", ".next", "out", ".git", "__pycache__":
		return true
	}
	return false
}

// checkSampleFile scans one page, tracking whether the current line is inside a
// Go sample. Outside one, the same text is prose about the rule rather than a
// breach of it: this file itself would otherwise fail its own check.
func checkSampleFile(rel, src string) []SampleProblem {
	allowed := allowedRules(src)
	var out []SampleProblem
	inGo := false
	for i, line := range strings.Split(src, "\n") {
		switch {
		case goLangProp.MatchString(line) || goCodeConst.MatchString(line):
			inGo = true
		case inGo && (strings.Contains(line, "language=\"ts\"") ||
			strings.Contains(line, "language=\"tsx\"") ||
			strings.Contains(line, "language=\"bash\"") ||
			strings.Contains(line, "icon: 'ts'") ||
			strings.Contains(line, "icon: 'tsx'")):
			inGo = false
		}
		if !inGo || isProse(line) {
			continue
		}
		for _, rule := range sampleRules {
			if allowed[rule.ID] || !rule.Match.MatchString(line) {
				continue
			}
			out = append(out, SampleProblem{
				File:    rel,
				Line:    i + 1,
				Rule:    rule.ID,
				Text:    strings.TrimSpace(line),
				Message: rule.Message,
			})
		}
	}
	return out
}

// isProse reports whether a line is page copy rather than code. A docs page
// explaining a rule quotes the thing the rule forbids, which is the point of
// the page: "Changing price:float to price:money replaces one column with two".
// Flagging that is how a checker earns its reputation for crying wolf.
func isProse(line string) bool {
	return strings.Contains(line, "<code>") ||
		strings.Contains(line, "</p>") ||
		strings.Contains(line, "<strong>") ||
		strings.Contains(line, "<li>")
}

// allowRe is a page opting out of one rule, with its reason on the same line:
//
//	{/* docscheck:allow float-money - a Go language primer, not Grit guidance */}
//
// Named, not blanket. A page that needs to show what the rule forbids says
// which rule and why, and keeps every other rule. The alternative is a
// skip-this-file switch, which is how a check stops covering the pages that
// drifted most.
var allowRe = regexp.MustCompile(`docscheck:allow\s+([a-z-]+)`)

// allowedRules reads the opt-outs a page declares.
func allowedRules(src string) map[string]bool {
	out := map[string]bool{}
	for _, m := range allowRe.FindAllStringSubmatch(src, -1) {
		out[m[1]] = true
	}
	return out
}
