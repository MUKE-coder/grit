package scaffold

import "fmt"

// golangci-lint configuration shipped into every scaffolded project.
//
// The enabled set was chosen by measurement, not taste: every linter here
// reports ZERO findings on a freshly generated project, so the first run is
// green. That matters more than breadth — a linter that opens with a hundred
// findings gets switched off within the hour, and then it protects nothing.
//
// The linters that were not green are documented at the bottom of the emitted
// file with their real counts, so a developer opts into them deliberately
// after a cleanup pass instead of meeting them as noise.
func golangciYAML() string {
	return `# golangci-lint configuration.
#
#   golangci-lint run ./...          # check
#   golangci-lint run --fix ./...    # fix what can be fixed automatically
#
# This starts GREEN on a fresh Grit project. Every linter below was verified to
# report zero findings on newly generated code, so anything it reports is
# something you introduced — which is the only way a linter stays switched on.
# The stricter set, and how to adopt it, is documented at the bottom.

version: "2"

run:
  timeout: 5m
  tests: true

linters:
  default: none
  enable:
    # -- Resource leaks ---------------------------------------------------
    # A Go web service dies of these, not of logic errors.
    - bodyclose      # unclosed HTTP response body: leaks an fd and a goroutine
    - sqlclosecheck  # sql.Rows / sql.Stmt never closed
    - rowserrcheck   # sql.Rows.Err() unchecked, hiding a partial result set

    # -- Correctness ------------------------------------------------------
    - govet          # copylocks, printf verbs, loop-variable capture, unreachable
    - ineffassign    # a value assigned and never read, usually a mistyped variable
    - nilerr         # returning nil immediately after checking err != nil
    - durationcheck  # multiplying two durations, which is essentially never meant
    - copyloopvar    # Go 1.22+ loop-variable semantics

    # -- Context discipline -----------------------------------------------
    # This framework already shipped a goroutine-per-request logger that leaked
    # under load. Both of these catch that shape.
    - noctx          # an HTTP request with no context, so nothing can cancel it
    - contextcheck   # a function that should inherit a context but makes its own

  settings:
    govet:
      disable:
        # "Constant reflect.Ptr should be inlined" - pure style, and the single
        # loudest check against generated code (19 hits on an empty project).
        - inline
        # Reorders struct fields to save padding. Real, but it fights both the
        # generator's field order and the readability of a model.
        - fieldalignment

  exclusions:
    generated: lax
    rules:
      # Tests legitimately assert on rejected requests, where there is no body
      # to close and no context to thread.
      - path: _test\.go
        linters:
          - bodyclose
          - noctx

issues:
  # No per-linter cap: a truncated report makes a codebase look cleaner than it
  # is, which is the opposite of the point.
  max-issues-per-linter: 0
  max-same-issues: 0

# ---------------------------------------------------------------------------
# The stricter set - opt in when you are ready.
#
# Each of these is worth enabling, but each reports findings on a brand-new
# project, so they are off rather than greeting you with a wall of output.
# Counts are from a freshly generated API. Enable one, clear it, move on:
#
#   errcheck     (~33)  unchecked error returns. Many are deliberate
#                       fire-and-forget (logging, cache writes) - list those
#                       under settings.errcheck.exclude-functions and fix the
#                       rest. Highest value of the five.
#   staticcheck  (~26)  mostly SA1019 deprecation notices from the AWS SDK
#                       endpoint resolver. Real, but upgrade work with its own
#                       migration path.
#   errorlint    (~22)  %v on an error instead of %w, which breaks errors.Is
#                       and errors.As for every caller downstream.
#   gosec         (~7)  review individually. G404 (math/rand) is fine for
#                       jitter and bucketing; anything that must be
#                       unpredictable already uses crypto/rand.
#   unused        (~3)  seed and mail helpers kept for your own code to call.
#
# Add the name to linters.enable above to turn one on.
# ---------------------------------------------------------------------------
`
}

// lintCIYAML emits .github/workflows/lint.yml.
func lintCIYAML(opts Options) string {
	// --single puts the Go module at the repo root; every other mode nests it
	// under apps/api. Pointing the action at the wrong directory fails with
	// "no go files", which reads like a broken lint config rather than a path
	// mistake.
	workdir := "apps/api"
	if opts.Architecture == ArchSingle {
		workdir = "."
	}
	return fmt.Sprintf(`name: Lint

on:
  push:
    branches: [main, master]
  pull_request:

jobs:
  golangci:
    name: golangci-lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
          cache: true
          cache-dependency-path: %[1]s/go.sum

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          working-directory: %[1]s
          # Only flag findings the PR introduced. The shipped config is green
          # on a new project, so this mainly matters once you enable the
          # stricter linters documented in .golangci.yml.
          only-new-issues: true
`, workdir)
}
