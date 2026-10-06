package scaffold

import (
	"strings"
	"testing"
)

// air rebuilds the package this project actually has.
//
// The config said ./cmd/server, which a monorepo has and a single does not: its
// main.go is at the top of its module. air failed with "directory not found",
// the API never started, and `grit start` shut the frontend down with it. What
// the user saw was a Vite dev server proxying to nothing and a login that did
// nothing, with the real error twenty lines up a log they had no reason to read.
func TestAirBuildsThePackageThisLayoutHas(t *testing.T) {
	single := airConfig(Options{ProjectName: "app", Architecture: ArchSingle})
	if !strings.Contains(single, `cmd = "go build -o ./tmp/server.exe ."`) {
		t.Errorf("a single project's air builds the wrong package:\n%s", single)
	}

	triple := airConfig(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext})
	if !strings.Contains(triple, `cmd = "go build -o ./tmp/server.exe ./cmd/server"`) {
		t.Errorf("a monorepo's air stopped building cmd/server:\n%s", triple)
	}

	// A project from before the api/ layout keeps its entry point under
	// cmd/server, like a monorepo.
	legacy := airConfig(Options{ProjectName: "app", Architecture: ArchSingle, LegacySingleFlat: true})
	if !strings.Contains(legacy, `cmd = "go build -o ./tmp/server.exe ./cmd/server"`) {
		t.Errorf("an existing single project's air changed under it:\n%s", legacy)
	}

	// And the Go expression is interpolated, not printed. The first version of
	// this put the concatenation inside a raw string, so every project got the
	// literal text `" + airMainPackage(opts) + "` in its config.
	for _, cfg := range []string{single, triple, legacy} {
		if strings.Contains(cfg, "airMainPackage") {
			t.Errorf("the config carries Go source rather than its result:\n%s", cfg)
		}
	}
}

// The dev proxy and the API agree on a port.
//
// The proxy named 8080 outright while the API reads APP_PORT from .env. Change
// that variable, which is a thing people do, and every call through the proxy is
// refused: the browser shows a login that does nothing, and the Vite log says
// ECONNREFUSED for a target nobody chose.
func TestViteProxyFollowsTheConfiguredPort(t *testing.T) {
	cfg := singleFrontendViteConfig()

	if !strings.Contains(cfg, "loadEnv") {
		t.Error("the config does not read .env, so it cannot know which port the API is on")
	}
	if !strings.Contains(cfg, "env.APP_PORT") {
		t.Error("the proxy does not read APP_PORT, the variable the Go binary reads")
	}
	if strings.Contains(cfg, "target: 'http://localhost:8080'") {
		t.Error("a proxy target is still pinned to 8080")
	}

	// Every path the binary serves goes through the same target: the admin
	// panel links to GORM Studio, Pulse and Sentinel, and in development those
	// are on the API, not on the dev server.
	for _, path := range []string{"'/api'", "'/studio'", "'/pulse'", "'/sentinel'", "'/docs'"} {
		if !strings.Contains(cfg, path+": { target: apiTarget") {
			t.Errorf("%s does not proxy to the API", path)
		}
	}
}
