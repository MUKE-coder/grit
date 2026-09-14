package scaffold

import "strings"

// The CI a generated project gets: Dependabot, the security scans and a build
// and test run on every pull request (H21 in the contact-app review).
//
// Every one of these used to assume the single-app layout, a go.mod at the root
// and the frontend in frontend/, whatever the project was. In a monorepo that
// meant govulncheck ran where there is no module, the audit ran in a directory
// that does not exist, Dependabot watched neither, and no workflow ran the tests
// on a pull request at all.

// ciLayout is where a project's code lives, for the workflows.
type ciLayout struct {
	apiDir      string // the Go module, as a workflow path: "apps/api" or "."
	jsDir       string // where pnpm install runs: "." or "frontend"
	hasFrontend bool
	single      bool
	desktop     bool
}

func ciLayoutFor(opts Options) ciLayout {
	l := ciLayout{apiDir: "apps/api", jsDir: ".", hasFrontend: opts.Architecture != ArchAPI, desktop: opts.ShouldIncludeDesktop()}
	if opts.Architecture == ArchSingle {
		l.apiDir, l.jsDir, l.single = ".", "frontend", true
	}
	return l
}

// dependabotDir is a workflow path as Dependabot writes it: rooted, no dot.
func dependabotDir(dir string) string {
	if dir == "." {
		return "/"
	}
	return "/" + dir
}

// ciEmbedPlaceholderStep lets a single-app project's Go compile in CI. The server
// embeds frontend/dist, which .gitignore leaves out of the repository, and
// //go:embed refuses a pattern that matches nothing.
const ciEmbedPlaceholderStep = `      # The server embeds frontend/dist, which is built rather than committed.
      # A placeholder is enough for Go to compile.
      - name: Placeholder for the embedded frontend
        run: mkdir -p frontend/dist && [ -e frontend/dist/index.html ] || echo '<!doctype html><title>Build the frontend</title>' > frontend/dist/index.html
`

// ciPnpmSteps installs the frontend. The project's package.json names pnpm for a
// monorepo; a single app's frontend does not, so the version is given here.
func ciPnpmSteps(l ciLayout) string {
	version := ""
	if l.single {
		version = "        with:\n          version: 10\n"
	}
	return `      - uses: pnpm/action-setup@v4
` + version + `      - uses: actions/setup-node@v4
        with:
          node-version: '22'
      # The committed lockfile is what gets installed and audited. Until one is
      # committed, the install resolves one.
      - name: Install
        run: if [ -f pnpm-lock.yaml ]; then pnpm install --frozen-lockfile; else pnpm install; fi
`
}

func ciReplace(l ciLayout, s string) string {
	return strings.NewReplacer(
		"{{API_DIR}}", l.apiDir,
		"{{JS_DIR}}", l.jsDir,
		"{{GO_DEPENDABOT_DIR}}", dependabotDir(l.apiDir),
		"{{JS_DEPENDABOT_DIR}}", dependabotDir(l.jsDir),
	).Replace(s)
}

// dependabotYAML is .github/dependabot.yml.
func dependabotYAML(opts Options) string {
	l := ciLayoutFor(opts)
	npm := ""
	if l.hasFrontend {
		npm = `
  # JavaScript packages. From the pnpm workspace, which covers every app and
  # package in it.
  - package-ecosystem: npm
    directory: "{{JS_DEPENDABOT_DIR}}"
    schedule:
      interval: weekly
    open-pull-requests-limit: 10
    labels:
      - dependencies
      - javascript
`
	}
	return ciReplace(l, `# Dependabot: automated dependency updates.
# Surfaces vulnerable and outdated dependencies as pull requests you can review
# and merge (OWASP Top 10:2025 A03, Software Supply Chain Failures).

version: 2

updates:
  # Go modules, where the API's go.mod is.
  - package-ecosystem: gomod
    directory: "{{GO_DEPENDABOT_DIR}}"
    schedule:
      interval: weekly
    open-pull-requests-limit: 10
    labels:
      - dependencies
      - go
`+npm+`
  # GitHub Actions workflows
  - package-ecosystem: github-actions
    directory: "/"
    schedule:
      interval: weekly
    labels:
      - dependencies
      - ci
`)
}

// securityCIYAML is .github/workflows/security.yml.
func securityCIYAML(opts Options) string {
	l := ciLayoutFor(opts)
	placeholder := ""
	if l.single {
		placeholder = ciEmbedPlaceholderStep
	}
	audit, jsMatrix := "", ""
	if l.hasFrontend {
		audit = `
  pnpm-audit:
    name: JavaScript dependency audit (pnpm audit)
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: {{JS_DIR}}
    steps:
      - uses: actions/checkout@v4
` + ciPnpmSteps(l) + `      # High and critical, in production dependencies: a moderate issue in a
      # dev tool should not block every pull request.
      - name: Audit
        run: pnpm audit --prod --audit-level=high
`
		jsMatrix = `          - language: javascript-typescript
            build-mode: none
`
	}
	return ciReplace(l, `name: security

# Security scans on every pull request, and weekly for the advisories published
# about code that has not changed (OWASP Top 10:2025 A03 and A02). Pair with
# Dependabot.

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]
  schedule:
    - cron: '0 6 * * 1'

permissions:
  contents: read

jobs:
  govulncheck:
    name: Go vulnerability scan (govulncheck)
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: {{API_DIR}}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          # The Go the API builds with. Scanned on an older one, the report lists
          # standard library issues the build does not have.
          go-version-file: {{API_DIR}}/go.mod
          cache-dependency-path: {{API_DIR}}/go.sum
`+placeholder+`      - name: Run govulncheck
        run: go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./...
`+audit+`
  codeql:
    name: CodeQL static analysis
    runs-on: ubuntu-latest
    permissions:
      actions: read
      contents: read
      security-events: write
    strategy:
      fail-fast: false
      matrix:
        include:
          # Go is built by hand: its module is in {{API_DIR}}.
          - language: go
            build-mode: manual
`+jsMatrix+`    steps:
      - uses: actions/checkout@v4
      - if: matrix.language == 'go'
        uses: actions/setup-go@v5
        with:
          go-version-file: {{API_DIR}}/go.mod
`+placeholder+`      - uses: github/codeql-action/init@v3
        with:
          languages: ${{ matrix.language }}
          build-mode: ${{ matrix.build-mode }}
      - if: matrix.build-mode == 'manual'
        working-directory: {{API_DIR}}
        run: go build ./...
      - uses: github/codeql-action/analyze@v3
        with:
          category: "/language:${{ matrix.language }}"
`)
}

// ciYAML is .github/workflows/ci.yml: the tests, on every push and pull request.
func ciYAML(opts Options) string {
	l := ciLayoutFor(opts)
	placeholder := ""
	if l.single {
		placeholder = ciEmbedPlaceholderStep
	}
	web := ""
	if l.hasFrontend {
		steps := `      - name: Type-check
        run: pnpm type-check
      - name: Test
        run: pnpm test
      - name: Build
        run: pnpm build
`
		name := "Frontend (type-check, test, build)"
		if l.single {
			// The single app's build is tsr generate, tsc and vite: it
			// type-checks as it builds, and the app has no test suite of its own.
			steps = `      - name: Type-check and build
        run: pnpm build
`
			name = "Frontend (type-check, build)"
		}
		web = `
  web:
    name: ` + name + `
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: {{JS_DIR}}
    steps:
      - uses: actions/checkout@v4
` + ciPnpmSteps(l) + steps
	}
	return ciReplace(l, `name: CI

# The tests, on every push and pull request, so a change that breaks them is
# seen before it is merged rather than when it is released.

on:
  push:
    branches: [main, master]
  pull_request:

permissions:
  contents: read

jobs:
  api:
    name: API (vet, test)
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: {{API_DIR}}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: {{API_DIR}}/go.mod
          cache-dependency-path: {{API_DIR}}/go.sum
`+placeholder+`      - name: Vet
        run: go vet ./...
      - name: Test
        run: go test -race ./...
`+web)
}
