package doctor

import (
	"os"
	"path/filepath"
	"strings"
)

// One address, written down four times.
//
// A Grit project says where its API is in four places:
//
//	APP_PORT                  the port the Go server binds
//	APP_URL                   what the API hands a client for a stored file
//	API_URL                   what the frontends call, baked into their bundles
//	apps/*/.env.local         NEXT_PUBLIC_API_URL, written once by grit seed
//
// The comment above APP_PORT invites you to change it. Two of the four move
// with it if you are careful; the fourth is written by the seeder with the
// header "Only ever created, never overwritten", so a project seeded before
// the port moved keeps pointing at the old one forever.
//
// Every symptom of the mismatch points somewhere else. A frontend calling a
// port nothing is listening on reports a network error that looks like the API
// being down. A frontend calling a port something else is listening on reports
// a CORS failure that looks like a CORS misconfiguration. Neither names the
// setting that is wrong.
//
// Found by moving APP_PORT to 8099 on a generated triple-tier project. The
// existing APP_URL check passed, and sign-in still did nothing.

// checkAPIAddressAgrees compares every place the API's address is written
// against APP_PORT.
func checkAPIAddressAgrees(p *project) []Finding {
	port := p.env["APP_PORT"]
	if port == "" {
		return nil
	}

	var findings []Finding

	// API_URL, from the root .env: what the frontends are built to call.
	if url := p.env["API_URL"]; url != "" && isLocalURL(url) {
		if got, ok := urlPort(url); ok && got != port {
			findings = append(findings, Finding{
				Level:    "error",
				Resource: "API_URL",
				Message: "API_URL points at port " + got + " and the API listens on " + port +
					". Every call the web app and the admin panel make goes to the wrong " +
					"port, and the browser reports that as a network or CORS failure rather " +
					"than as the setting it is.",
				Fix: "Change API_URL to use :" + port + ".",
			})
		}
	}

	// The per-app file, written once by grit seed and never rewritten, so this
	// is the one that goes stale without anybody editing anything.
	//
	// Each bundler exposes only its own prefix, so the name differs per app and
	// all of them are checked: a file naming the wrong one is its own problem,
	// and a file naming the right one with the wrong port is this one.
	for _, app := range []string{"web", "admin", "expo", "desktop"} {
		path := filepath.Join(p.root, "apps", app, ".env.local")
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		body := string(raw)
		rel := "apps/" + app + "/.env.local"

		// An Expo app is given the port on its own, because a phone cannot
		// reach a localhost URL: that address is the phone.
		if got := envValueIn(body, "EXPO_PUBLIC_API_PORT"); got != "" && got != port {
			findings = append(findings, Finding{
				Level:    "error",
				Resource: rel,
				Message: "EXPO_PUBLIC_API_PORT is " + got + " and the API listens on " + port +
					". The app derives the dev machine's address from Metro and takes the " +
					"port from here, so every request goes to a port nothing answers.",
				Fix: "Change EXPO_PUBLIC_API_PORT in " + rel + " to " + port + ".",
			})
		}

		for _, key := range []string{"NEXT_PUBLIC_API_URL", "VITE_API_URL", "EXPO_PUBLIC_API_URL"} {
			url := envValueIn(body, key)
			if url == "" || !isLocalURL(url) {
				continue
			}
			got, ok := urlPort(url)
			if !ok || got == port {
				continue
			}
			findings = append(findings, Finding{
				Level:    "error",
				Resource: rel,
				Message: key + " points at port " + got + " and the API listens on " + port +
					". This file is written once by grit seed and never overwritten, so it " +
					"keeps the port the project had when it was first seeded, and it wins " +
					"over API_URL in the root .env.",
				Fix: "Change " + key + " in " + rel + " to use :" + port +
					", or delete the line and let the root API_URL answer.",
			})
		}
	}

	return findings
}

// isLocalURL says whether a URL is one whose port can be compared to APP_PORT.
// A deployed origin is behind a proxy on 443 and has no business matching.
func isLocalURL(url string) bool {
	return strings.Contains(url, "localhost") || strings.Contains(url, "127.0.0.1")
}

// urlPort returns the port in a URL, if it names one.
func urlPort(url string) (string, bool) {
	match := portInURL.FindStringSubmatch(url)
	if match == nil {
		return "", false
	}
	return match[1], true
}

// envValueIn reads one KEY=value out of an env file's text.
func envValueIn(body, key string) string {
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, key+"=") {
			return strings.TrimSpace(strings.TrimPrefix(line, key+"="))
		}
	}
	return ""
}
