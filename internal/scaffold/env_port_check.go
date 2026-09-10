package scaffold

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// envPortDrift reports service URLs in a project's .env that name localhost on
// a different port from the one docker-compose publishes that service on.
//
// v3.205.0 moved the host ports into .env and wrote REDIS_URL and
// MINIO_ENDPOINT beside them with the default ports baked in. A project that
// moved REDIS_PORT kept dialling 6380, which with a second Grit project running
// is that project's Redis, and nothing fails: the wrong Redis answers. New
// projects build both addresses from the ports. An existing .env still names
// the old ones, and grit upgrade does not rewrite .env, so it says so here.
func envPortDrift(root string) []string {
	vals := readDotEnv(filepath.Join(root, ".env"))
	checks := []struct{ urlVar, portVar, scheme string }{
		{"REDIS_URL", "REDIS_PORT", "redis"},
		{"MINIO_ENDPOINT", "MINIO_PORT", "http"},
	}

	var out []string
	for _, c := range checks {
		raw, port := vals[c.urlVar], vals[c.portVar]
		if raw == "" || port == "" {
			continue
		}
		u, err := url.Parse(raw)
		if err != nil {
			continue
		}
		host := u.Hostname()
		if host != "localhost" && host != "127.0.0.1" {
			continue
		}
		if got := u.Port(); got != "" && got != port {
			out = append(out, fmt.Sprintf(
				"%s in .env points at port %s, but %s is %s: the API is probably using another project's service. Set %s=%s://%s:%s",
				c.urlVar, got, c.portVar, port, c.urlVar, c.scheme, host, port))
		}
	}
	return out
}

// readDotEnv reads KEY=VALUE lines. Comments, blank lines and export prefixes
// are ignored, an inline " # comment" is dropped, and surrounding quotes are
// removed. A missing file is an empty map, not an error: plenty of projects
// keep their environment elsewhere.
func readDotEnv(path string) map[string]string {
	vals := map[string]string{}
	f, err := os.Open(path)
	if err != nil {
		return vals
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if i := strings.Index(value, " #"); i >= 0 {
			value = value[:i]
		}
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		vals[strings.TrimSpace(key)] = value
	}
	return vals
}
