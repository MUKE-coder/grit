package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PlanFile is where the runbook is written, at the project root.
//
// A file rather than terminal output because the whole point is that somebody
// hands it to a coding agent, and a path is easier to hand over than a screen
// of scrollback. It is also the thing to read on a Monday after upgrading on a
// Friday.
const PlanFile = "UPGRADE-PLAN.md"

// writeUpgradePlan emits a runbook for the part of an upgrade that could not be
// done mechanically.
//
// # Why this exists
//
// `grit upgrade` rewrites the project: manifest-guarded delivery, a repair
// function per release, and a per-file status so a file somebody edited is
// never clobbered. That last protection is also the limit. A file you have
// edited is left exactly as it is, and what the new version would have changed
// in it goes unapplied, reported as one line of terminal output that scrolls
// away.
//
// There is no mechanical answer to that: the whole reason the file was skipped
// is that applying the change needs judgment about code somebody wrote. But a
// change that needs judgment is precisely the kind an agent can carry out when
// told what the change is, what to look at and how to check the result, and
// that is what this writes.
//
// # The shape, and why it is this shape
//
// Each entry is detect, change, verify. Not prose: a command that finds the
// affected lines, a diff showing what the new version does, and a command whose
// success means it worked. An instruction without a verification step is one an
// agent will report as done whether or not it is.
func writeUpgradePlan(root, fromVersion, toVersion string, skipped []SkippedFile, delivered int) (string, error) {
	sorted := make([]SkippedFile, len(skipped))
	copy(sorted, skipped)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Rel < sorted[j].Rel })

	var b strings.Builder

	fmt.Fprintf(&b, "# Finishing the upgrade to v%s\n\n", toVersion)
	if fromVersion != "" {
		fmt.Fprintf(&b, "This project was on v%s.\n\n", fromVersion)
	}

	b.WriteString("`grit upgrade` has already run and rewritten ")
	fmt.Fprintf(&b, "%d %s.\n", delivered, plural(delivered, "file", "files"))
	b.WriteString("What is left is the part it will not do on its own: files you have edited\n")
	b.WriteString("since Grit wrote them. Grit leaves those exactly as they are, because\n")
	b.WriteString("overwriting somebody's work to deliver a template is the one thing the\n")
	b.WriteString("upgrade machinery exists not to do.\n\n")

	if len(sorted) == 0 {
		b.WriteString("**Nothing is outstanding.** Every file either took the new version or was\n")
		b.WriteString("already current. You can delete this file.\n")
		return writePlan(root, b.String())
	}

	b.WriteString("## How to use this\n\n")
	b.WriteString("Paste this whole file into a coding agent, or work through it yourself.\n")
	b.WriteString("Each section is the same three steps: find the affected lines, make the\n")
	b.WriteString("change, run something that proves it worked.\n\n")
	b.WriteString("Before you start, commit or stash what you have, and check the project\n")
	b.WriteString("builds now, so a failure later belongs to this work and not to something\n")
	b.WriteString("that was already broken.\n\n")
	b.WriteString("```shell\ngit status --short\n")
	b.WriteString("cd apps/api && go build ./... && cd ../..\n```\n\n")

	b.WriteString("Each diff below is **what Grit would have written**, against **what you\n")
	b.WriteString("have**. It is not a patch to apply blindly: your edits are in there for a\n")
	b.WriteString("reason, and the job is to carry the new behaviour into them rather than to\n")
	b.WriteString("replace one with the other. Where the two genuinely conflict, keep yours\n")
	b.WriteString("and note why.\n\n")

	fmt.Fprintf(&b, "## %d %s to finish\n\n", len(sorted), plural(len(sorted), "file", "files"))

	for i, file := range sorted {
		fmt.Fprintf(&b, "### %d. `%s`\n\n", i+1, file.Rel)

		fmt.Fprintf(&b, "**Detect.** Open it, and see what you changed:\n\n")
		fmt.Fprintf(&b, "```shell\ngit log --oneline -- %s\ngit diff HEAD -- %s\n```\n\n", file.Rel, file.Rel)

		b.WriteString("**Change.** This is what the new version does to it:\n\n")
		b.WriteString("```diff\n")
		b.WriteString(trimDiff(file.Diff()))
		b.WriteString("\n```\n\n")

		if verify := verifyCommandFor(file.Rel); verify != "" {
			fmt.Fprintf(&b, "**Verify.**\n\n```shell\n%s\n```\n\n", verify)
		}
	}

	b.WriteString("## When every section is done\n\n")
	b.WriteString("```shell\ncd apps/api && go build ./... && go vet ./... && go test ./... && cd ../..\n")
	b.WriteString("grit doctor\n```\n\n")
	b.WriteString("Then delete this file: it describes one upgrade and is wrong after the next\n")
	b.WriteString("one.\n")

	return writePlan(root, b.String())
}

func writePlan(root, body string) (string, error) {
	path := filepath.Join(root, PlanFile)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return "", fmt.Errorf("writing %s: %w", PlanFile, err)
	}
	return path, nil
}

// verifyCommandFor is how you check that one file is still right.
//
// Per file type rather than one command for everything, because "the whole
// project builds" is slow enough that nobody runs it between edits, and a check
// nobody runs is not a check.
func verifyCommandFor(rel string) string {
	rel = filepath.ToSlash(rel)

	switch {
	case strings.HasSuffix(rel, ".go"):
		if dir := goModuleDirFor(rel); dir != "" {
			return "cd " + dir + " && gofmt -l . && go build ./... && go vet ./..."
		}
		return "gofmt -l . && go build ./... && go vet ./..."

	case strings.HasSuffix(rel, ".tsx"), strings.HasSuffix(rel, ".ts"):
		if app := frontendAppFor(rel); app != "" {
			return "cd " + app + " && npx tsc --noEmit"
		}
		return "npx tsc --noEmit"

	case strings.HasSuffix(rel, "docker-compose.yml"), strings.HasSuffix(rel, "docker-compose.prod.yml"),
		strings.HasSuffix(rel, "docker-compose.shared-network.yml"):
		return "docker compose -f " + rel + " config > /dev/null"

	case strings.HasSuffix(rel, ".json"):
		// A malformed package.json fails at install time, which is late and
		// loud. Parsing it here is immediate and quiet.
		return "node -e \"JSON.parse(require('fs').readFileSync('" + rel + "','utf8'))\" && echo ok"

	case strings.HasSuffix(rel, ".yml"), strings.HasSuffix(rel, ".yaml"):
		return "python -c \"import sys,yaml; yaml.safe_load(open(sys.argv[1]))\" " + rel + " && echo ok"

	case strings.HasSuffix(rel, ".env"), strings.HasSuffix(rel, ".env.example"):
		return "grit doctor"

	default:
		return ""
	}
}

// goModuleDirFor is the directory to run go build from for a Go file, which is
// the app it belongs to in a monorepo and the root in a single-app project.
func goModuleDirFor(rel string) string {
	if strings.HasPrefix(rel, "apps/api/") {
		return "apps/api"
	}
	if strings.HasPrefix(rel, "apps/desktop/") {
		return "apps/desktop"
	}
	return ""
}

// frontendAppFor is the Next.js app a .tsx file belongs to.
func frontendAppFor(rel string) string {
	for _, app := range []string{"apps/admin", "apps/web", "apps/expo", "packages/shared"} {
		if strings.HasPrefix(rel, app+"/") {
			return app
		}
	}
	return ""
}

// maxDiffLines caps one entry, so a runbook stays something a person or a model
// reads rather than something they scroll past.
//
// A file whose diff runs past this is one where the diff was never going to be
// the useful part: the entry says so and points at the command that shows the
// whole thing.
const maxDiffLines = 200

func trimDiff(diff string) string {
	diff = strings.TrimRight(diff, "\n")
	lines := strings.Split(diff, "\n")
	if len(lines) <= maxDiffLines {
		return diff
	}
	kept := strings.Join(lines[:maxDiffLines], "\n")
	return kept + fmt.Sprintf("\n\n... %d more lines. Run `grit upgrade --diff` to see all of it.",
		len(lines)-maxDiffLines)
}
