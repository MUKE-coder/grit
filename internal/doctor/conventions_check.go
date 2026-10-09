package doctor

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// The conventions a machine can check.
//
// Most of what a Grit developer gets wrong is written down in the project's
// AGENTS.md and read by a person or an agent. Three of those rules are
// mechanical, and all three fail the same way: the project builds, nothing
// logs an error, and a screen is quietly broken. A check is worth more than a
// paragraph for exactly those.

// checkNodeLinker catches a workspace missing the hoisted node linker.
//
// pnpm 10 ignores `node-linker` in .npmrc when a pnpm-workspace.yaml exists,
// and the setting has to be in the workspace file instead. Without it, Next
// cannot resolve @swc/helpers and every page throws
// _interop_require_wildcard: a complete outage that looks like a Next bug and
// has nothing to do with Next.
func checkNodeLinker(p *project) []Finding {
	workspace := filepath.Join(p.root, "pnpm-workspace.yaml")
	body, err := os.ReadFile(workspace)
	if err != nil {
		// No workspace file: .npmrc is read and there is nothing to warn
		// about.
		return nil
	}
	if strings.Contains(string(body), "nodeLinker:") || strings.Contains(string(body), "node-linker:") {
		return nil
	}

	// Only a problem when there is a Next app to break.
	if !hasNextApp(p.root) {
		return nil
	}

	return []Finding{{
		Level:    "error",
		Resource: "pnpm-workspace.yaml",
		Message: "pnpm-workspace.yaml does not set nodeLinker, and pnpm 10 ignores " +
			"node-linker in .npmrc when a workspace file exists. Without the hoisted " +
			"layout, Next cannot resolve @swc/helpers and every page throws " +
			"_interop_require_wildcard.",
		Fix: "Add `nodeLinker: hoisted` to pnpm-workspace.yaml, then reinstall.",
	}}
}

func hasNextApp(root string) bool {
	for _, app := range []string{"web", "admin"} {
		for _, name := range []string{"next.config.mjs", "next.config.js", "next.config.ts"} {
			if _, err := os.Stat(filepath.Join(root, "apps", app, name)); err == nil {
				return true
			}
		}
	}
	for _, name := range []string{"next.config.mjs", "next.config.js", "next.config.ts"} {
		if _, err := os.Stat(filepath.Join(root, name)); err == nil {
			return true
		}
	}
	return false
}

var portInURL = regexp.MustCompile(`:(\d{2,5})(/|$)`)

// checkAppURLPort catches APP_URL pointing at a port nothing is listening on.
//
// APP_URL is what the API hands a client for a stored file. Move APP_PORT and
// leave APP_URL behind and every image loads from the wrong port: the upload
// succeeded, the row is right, and the picture is broken, which sends people
// looking at the storage driver.
func checkAppURLPort(p *project) []Finding {
	port := p.env["APP_PORT"]
	appURL := p.env["APP_URL"]
	if port == "" || appURL == "" {
		return nil
	}
	// Only local URLs can be compared. A deployed APP_URL is behind a proxy on
	// 443 and has no business matching APP_PORT.
	if !strings.Contains(appURL, "localhost") && !strings.Contains(appURL, "127.0.0.1") {
		return nil
	}

	match := portInURL.FindStringSubmatch(appURL)
	if match == nil {
		// No port in the URL means port 80, which APP_PORT is not.
		if port != "80" {
			return []Finding{{
				Level:    "error",
				Resource: "APP_URL",
				Message: "APP_URL has no port but APP_PORT is " + port + ". Anything the API " +
					"hands a client for a stored file will point at port 80.",
				Fix: "Set APP_URL to include :" + port + ".",
			}}
		}
		return nil
	}
	if match[1] == port {
		return nil
	}

	return []Finding{{
		Level:    "error",
		Resource: "APP_URL",
		Message: "APP_URL points at port " + match[1] + " and APP_PORT is " + port + ". " +
			"Every URL the API hands a client for a stored file will point at a port " +
			"nothing is listening on: the upload succeeds, the row is correct, and the " +
			"image is broken.",
		Fix: "Change APP_URL to use :" + port + ". They are one setting in two places.",
	}}
}

var nextDependency = regexp.MustCompile(`"next"\s*:\s*"([^"]+)"`)

// checkNextPinned catches a caret range on Next.
//
// It is pinned exactly on purpose. A caret once adopted a release a day old
// whose layout router threw on every navigation, in a project that had changed
// nothing. A range on the framework that renders every page is a change you
// did not make arriving on a morning you did not choose.
func checkNextPinned(p *project) []Finding {
	var findings []Finding

	for _, dir := range []string{
		filepath.Join(p.root, "apps", "web"),
		filepath.Join(p.root, "apps", "admin"),
		p.root,
	} {
		path := filepath.Join(dir, "package.json")
		body, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		match := nextDependency.FindSubmatch(body)
		if match == nil {
			continue
		}
		version := string(match[1])
		if !strings.HasPrefix(version, "^") && !strings.HasPrefix(version, "~") &&
			!strings.HasPrefix(version, ">") {
			continue
		}
		rel, err := filepath.Rel(p.root, path)
		if err != nil {
			rel = path
		}
		findings = append(findings, Finding{
			Level:    "warning",
			Resource: filepath.ToSlash(rel),
			Message: "Next is at \"" + version + "\", a range rather than an exact version. " +
				"A caret once adopted a one-day-old release whose layout router threw on " +
				"every navigation, in a project that had changed nothing.",
			Fix: "Pin it exactly: drop the " + version[:1] + ".",
		})
	}

	return findings
}
