// Package clierr is the catalogue of problems the CLI can explain.
//
// A message that says what went wrong is worth more than one that says it went
// wrong, and a message that says where to read more is worth more again. Every
// entry here has a slug, and every slug has a page at
// https://gritframework.dev/e/<slug>.
//
// The pages and this catalogue are two things that must agree, so a test reads
// both: TestEveryProblemHasAPage in the docs checks that every slug here has a
// page on disk and that every page has an entry here. Without it, a renamed
// slug leaves the CLI pointing at a 404, which is worse than printing no link
// at all.
//
// The rule for an entry: it must be a problem a person can actually fix, and
// the Fix must be the thing to type or do, not a description of the category.
// "Check your Docker installation" is not a fix. "Open Docker Desktop and wait
// for the whale to stop animating" is.
package clierr

import "fmt"

// BaseURL is where the explanation pages live.
const BaseURL = "https://gritframework.dev/e/"

// Problem is one thing that goes wrong, with the words to say about it.
type Problem struct {
	// Slug is the URL segment and the stable identifier. Never renamed
	// without a redirect, because it is printed into people's terminals and
	// pasted into issues.
	Slug string
	// Title is the one line above the explanation. Written as the thing that
	// happened rather than as a category.
	Title string
	// Cause is why it happened, in a sentence.
	Cause string
	// Fix is what to do. A command where there is one.
	Fix string
}

// URL is where to read more about this problem.
func (p Problem) URL() string { return BaseURL + p.Slug }

// Error renders the problem as an error, so it can be returned from a command
// and still carry its explanation.
func (p Problem) Error() string {
	return fmt.Sprintf("%s: %s", p.Title, p.Cause)
}

// All is every problem the CLI knows how to explain, in the order the
// explanation index lists them.
//
// Keep this sorted by slug. The test that pairs it with the pages reports
// differences as sets, so order does not matter to it, but a sorted list is
// one somebody can scan.
var All = []Problem{
	{
		Slug:  "docker-not-running",
		Title: "Docker is not running",
		Cause: "The Docker daemon did not answer. On macOS and Windows that almost always means Docker Desktop is closed.",
		Fix:   "Open Docker Desktop and wait for its icon to stop animating. On Linux: sudo systemctl start docker",
	},
	{
		Slug:  "go-not-found",
		Title: "Go is not installed, or not on PATH",
		Cause: "The go command could not be run. Installing Go changes PATH, and a terminal that was already open keeps the old one.",
		Fix:   "Install Go from https://go.dev/dl, then open a NEW terminal and run: go version",
	},
	{
		Slug:  "go-too-old",
		Title: "Your Go is older than Grit needs",
		Cause: "A generated project's go.mod asks for a newer Go than the one on PATH, and older patches carry standard-library vulnerabilities this code reaches.",
		Fix:   "Install the current Go from https://go.dev/dl, then: go version",
	},
	{
		Slug:  "node-not-found",
		Title: "Node is not installed, or not on PATH",
		Cause: "The node command could not be run. The frontends need it; an API-only project does not.",
		Fix:   "Install Node from https://nodejs.org, or with a version manager: fnm install 22",
	},
	{
		Slug:  "node-too-old",
		Title: "Your Node is older than the frontend needs",
		Cause: "Next.js declares a minimum Node version, and an older one fails during the build rather than at install.",
		Fix:   "fnm install 22 && fnm use 22, or nvm install 22 && nvm use 22",
	},
	{
		Slug:  "not-a-grit-project",
		Title: "This is not a Grit project",
		Cause: "The command looks for grit.json at the project root and did not find it here or above.",
		Fix:   "cd into the directory grit new created, or run: grit new <name>",
	},
	{
		Slug:  "pnpm-not-found",
		Title: "pnpm is not installed",
		Cause: "Generated projects install with pnpm, not npm, because the workspace layout depends on it.",
		Fix:   "corepack enable && corepack prepare pnpm@latest --activate",
	},
	{
		Slug:  "port-in-use",
		Title: "A port the project needs is already taken",
		Cause: "Another process is listening on it: usually a server from an earlier run whose terminal was closed without stopping it, or another project's containers.",
		Fix:   "macOS and Linux: lsof -ti:PORT | xargs kill. Windows: Get-NetTCPConnection -LocalPort PORT -State Listen | Select-Object -Expand OwningProcess | ForEach-Object { Stop-Process -Id $_ -Force }",
	},
}

// BySlug returns the problem with that slug.
func BySlug(slug string) (Problem, bool) {
	for _, p := range All {
		if p.Slug == slug {
			return p, true
		}
	}
	return Problem{}, false
}

// Slugs is every slug, for the test that pairs this with the pages.
func Slugs() []string {
	out := make([]string, 0, len(All))
	for _, p := range All {
		out = append(out, p.Slug)
	}
	return out
}
