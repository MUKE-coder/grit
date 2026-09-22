package plugin

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVideoPluginIsRegistered(t *testing.T) {
	p, err := Get("video")
	if err != nil {
		t.Fatalf("the video plugin is not registered: %v", err)
	}
	if len(p.NodeDeps) != 1 || p.NodeDeps[0].Name != "expo-video" || p.NodeDeps[0].Workspace != "apps/expo" {
		t.Errorf("NodeDeps = %+v, want expo-video in apps/expo", p.NodeDeps)
	}
}

// The API half always installs; each client only where the app is.
func TestVideoFilesFollowTheProject(t *testing.T) {
	root := t.TempDir()
	ctx := Context{Root: root, Module: "shop/apps/api", Architecture: "triple"}
	files := videoFiles(ctx)
	for _, want := range []string{"apps/api/internal/models/video.go", "apps/api/internal/video/ffmpeg.go", "apps/api/internal/video/service.go", "apps/api/internal/video/video_test.go", "apps/api/internal/handlers/video.go"} {
		if _, ok := files[want]; !ok {
			t.Errorf("missing %s", want)
		}
	}
	for path := range files {
		if strings.HasPrefix(path, "apps/web/") || strings.HasPrefix(path, "apps/expo/") {
			t.Errorf("%s was written to a project without that app", path)
		}
	}
	for _, app := range []string{"web", "expo"} {
		if err := os.MkdirAll(filepath.Join(root, "apps", app), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	files = videoFiles(ctx)
	for _, want := range []string{"apps/web/components/video-player.tsx", "apps/web/hooks/use-video.ts", "apps/expo/components/video-player.tsx", "apps/expo/lib/video.ts"} {
		if _, ok := files[want]; !ok {
			t.Errorf("missing %s once the app exists", want)
		}
	}
}

// Every Go template is gofmt-clean once the module is filled in, so a project
// stays gofmt-clean after the plugin, and none keeps the placeholder.
func TestVideoTemplatesAreFormattedGo(t *testing.T) {
	root := t.TempDir()
	for _, app := range []string{"web", "expo"} {
		if err := os.MkdirAll(filepath.Join(root, "apps", app), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	ctx := Context{Root: root, Module: "shop/apps/api", Architecture: "triple"}
	for path, src := range videoFiles(ctx) {
		if strings.Contains(src, "{{MODULE}}") {
			t.Errorf("%s still has {{MODULE}}", path)
		}
		if !strings.HasSuffix(path, ".go") {
			continue
		}
		formatted, err := format.Source([]byte(src))
		if err != nil {
			t.Errorf("%s does not parse: %v", path, err)
			continue
		}
		if string(formatted) != src {
			t.Errorf("%s is not gofmt-clean", path)
		}
	}
}

// The Dockerfile injection targets the marker the scaffold writes.
func TestVideoInjectsFFmpegAtTheRuntimeMarker(t *testing.T) {
	ctx := Context{Root: t.TempDir(), Module: "shop/apps/api", Architecture: "triple"}
	for _, inj := range videoInjections(ctx) {
		if strings.HasSuffix(inj.File, "Dockerfile") {
			if inj.Marker != "# grit:runtime-packages" || !strings.Contains(inj.Code, "ffmpeg") || !inj.Optional {
				t.Errorf("Dockerfile injection = %+v", inj)
			}
			return
		}
	}
	t.Error("no Dockerfile injection")
}
