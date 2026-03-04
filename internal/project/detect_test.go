package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsDesktop(t *testing.T) {
	dir := t.TempDir()

	if IsDesktop(dir) {
		t.Error("empty dir should not be desktop")
	}

	os.WriteFile(filepath.Join(dir, "wails.json"), []byte(`{"name":"test"}`), 0644)

	if !IsDesktop(dir) {
		t.Error("dir with wails.json should be desktop")
	}
}

func TestIsWeb(t *testing.T) {
	dir := t.TempDir()

	if IsWeb(dir) {
		t.Error("empty dir should not be web")
	}

	// turbo.json alone is not enough
	os.WriteFile(filepath.Join(dir, "turbo.json"), []byte(`{}`), 0644)
	if IsWeb(dir) {
		t.Error("turbo.json without apps/api should not be web")
	}

	// turbo.json + apps/api = web
	os.MkdirAll(filepath.Join(dir, "apps", "api"), 0755)
	if !IsWeb(dir) {
		t.Error("dir with turbo.json + apps/api should be web")
	}
}

func TestDetectProjectFrom_Desktop(t *testing.T) {
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, "wails.json"), []byte(`{"name":"myapp"}`), 0644)
	os.WriteFile(filepath.Join(root, "go.mod"), []byte("module myapp\n\ngo 1.21\n"), 0644)

	info, err := DetectProjectFrom(root)
	if err != nil {
		t.Fatalf("DetectProjectFrom: %v", err)
	}

	if info.Type != ProjectDesktop {
		t.Errorf("expected desktop, got %s", info.Type)
	}
	if info.Module != "myapp" {
		t.Errorf("expected module myapp, got %s", info.Module)
	}
	if info.Root != root {
		t.Errorf("expected root %s, got %s", root, info.Root)
	}
}

func TestDetectProjectFrom_Web(t *testing.T) {
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, "turbo.json"), []byte(`{}`), 0644)
	os.MkdirAll(filepath.Join(root, "apps", "api"), 0755)
	os.WriteFile(filepath.Join(root, "apps", "api", "go.mod"), []byte("module myapp/apps/api\n\ngo 1.21\n"), 0644)

	info, err := DetectProjectFrom(root)
	if err != nil {
		t.Fatalf("DetectProjectFrom: %v", err)
	}

	if info.Type != ProjectWeb {
		t.Errorf("expected web, got %s", info.Type)
	}
	if info.Module != "myapp/apps/api" {
		t.Errorf("expected module myapp/apps/api, got %s", info.Module)
	}
}

func TestDetectProjectFrom_Subdirectory(t *testing.T) {
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, "wails.json"), []byte(`{"name":"myapp"}`), 0644)
	os.WriteFile(filepath.Join(root, "go.mod"), []byte("module myapp\n\ngo 1.21\n"), 0644)

	subDir := filepath.Join(root, "internal", "service")
	os.MkdirAll(subDir, 0755)

	info, err := DetectProjectFrom(subDir)
	if err != nil {
		t.Fatalf("DetectProjectFrom from subdir: %v", err)
	}

	if info.Type != ProjectDesktop {
		t.Errorf("expected desktop from subdir, got %s", info.Type)
	}
	if info.Root != root {
		t.Errorf("expected root %s from subdir, got %s", root, info.Root)
	}
}

func TestDetectProjectFrom_NotFound(t *testing.T) {
	dir := t.TempDir()

	_, err := DetectProjectFrom(dir)
	if err == nil {
		t.Error("expected error for empty directory")
	}
}

func TestDesktopTakesPrecedence(t *testing.T) {
	// If a directory has both wails.json and turbo.json, desktop wins
	root := t.TempDir()

	os.WriteFile(filepath.Join(root, "wails.json"), []byte(`{"name":"myapp"}`), 0644)
	os.WriteFile(filepath.Join(root, "turbo.json"), []byte(`{}`), 0644)
	os.MkdirAll(filepath.Join(root, "apps", "api"), 0755)
	os.WriteFile(filepath.Join(root, "go.mod"), []byte("module myapp\n\ngo 1.21\n"), 0644)

	info, err := DetectProjectFrom(root)
	if err != nil {
		t.Fatalf("DetectProjectFrom: %v", err)
	}

	if info.Type != ProjectDesktop {
		t.Errorf("expected desktop to take precedence, got %s", info.Type)
	}
}
