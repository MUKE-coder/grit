package setup

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeDefaults(t *testing.T) {
	opts := Options{ProjectName: "my-app", Yes: true}
	if err := opts.Normalize(); err != nil {
		t.Fatalf("Normalize: %v", err)
	}

	checks := map[string]string{
		"architecture": opts.Architecture,
		"frontend":     opts.Frontend,
		"theme":        opts.Theme,
		"db":           opts.DBProvider,
		"starter":      string(opts.Starter),
		"redis":        string(opts.Redis.Mode),
		"storage":      string(opts.Storage.Mode),
		"mail":         string(opts.Mail.Mode),
		"mail from":    opts.Mail.From,
	}
	want := map[string]string{
		"architecture": "triple",
		"frontend":     "next",
		"theme":        "atlas",
		"db":           "postgres",
		"starter":      "blank",
		"redis":        "local",
		"storage":      "minio",
		"mail":         "mailhog",
		"mail from":    "noreply@my-app.dev",
	}
	for name, got := range checks {
		if got != want[name] {
			t.Errorf("%s = %q, want %q", name, got, want[name])
		}
	}

	if opts.Start == nil {
		t.Fatal("Yes should set a start answer")
	}
	if *opts.Start {
		t.Error("Yes should default to not starting")
	}
}

func TestNormalizeForcesBlankWithoutAWebFrontend(t *testing.T) {
	for _, arch := range []string{"api", "mobile"} {
		opts := Options{
			ProjectName:  "my-app",
			Architecture: arch,
			Starter:      StarterBlog,
			Yes:          true,
		}
		if err := opts.Normalize(); err != nil {
			t.Fatalf("%s Normalize: %v", arch, err)
		}
		if opts.Starter != StarterBlank {
			t.Errorf("%s starter = %q, want blank", arch, opts.Starter)
		}
	}
}

func TestNormalizeRejectsIncompatiblePlugin(t *testing.T) {
	opts := Options{
		ProjectName:  "my-app",
		Architecture: "api",
		Frontend:     "next",
		Plugins:      []string{"kanban"},
		Yes:          true,
	}
	err := opts.Normalize()
	if err == nil {
		t.Fatal("kanban should be rejected on an API-only project")
	}
	if !strings.Contains(err.Error(), "admin") {
		t.Errorf("error should say why: %v", err)
	}
}

func TestPluginCompatibility(t *testing.T) {
	tests := []struct {
		name   string
		opts   Options
		plugin string
		want   bool
	}{
		{"api rejects admin", Options{Architecture: "api"}, "kanban", false},
		{"api accepts webhooks", Options{Architecture: "api"}, "webhooks", true},
		{"api rejects stripe", Options{Architecture: "api", Frontend: "next"}, "stripe", false},
		{"single rejects webhooks", Options{Architecture: "single"}, "webhooks", false},
		{"triple accepts admin", Options{Architecture: "triple"}, "kanban", true},
		{"triple accepts stripe", Options{Architecture: "triple", Frontend: "next"}, "stripe", true},
		{"triple rejects push", Options{Architecture: "triple"}, "push", false},
		{"full accepts push", Options{Architecture: "full"}, "push", true},
		{"tanstack rejects stripe", Options{Architecture: "triple", Frontend: "tanstack"}, "stripe", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.opts.PluginIncompatibility(tc.plugin) == ""
			if got != tc.want {
				t.Errorf("compatible = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestLocalServicesOnlyStartsWhatWasChosen(t *testing.T) {
	none := Options{
		StartDocker: true,
		DBProvider:  "sqlite",
		Redis:       RedisConfig{Mode: RedisOff},
		Storage:     StorageConfig{Mode: StorageS3},
		Mail:        MailConfig{Mode: MailLog},
	}
	if got := none.LocalServices(); len(got) != 0 {
		t.Errorf("services = %v, want none", got)
	}

	all := Options{
		StartDocker: true,
		DBProvider:  "postgres",
		Redis:       RedisConfig{Mode: RedisLocal},
		Storage:     StorageConfig{Mode: StorageMinIO},
		Mail:        MailConfig{Mode: MailMailhog},
	}
	got := all.LocalServices()
	want := []string{"postgres", "redis", "minio", "mailhog"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("services = %v, want %v", got, want)
	}
}

func TestSummaryDoesNotLeakSecrets(t *testing.T) {
	opts := Options{
		ProjectName:  "my-app",
		Architecture: "triple",
		Frontend:     "next",
		Starter:      StarterBlank,
		DBProvider:   "postgres",
		Redis:        RedisConfig{Mode: RedisExternal, URL: "rediss://default:supersecret@redis.example.com"},
		Storage:      StorageConfig{Mode: StorageS3, AccessKey: "AKIASECRET", SecretKey: "shhh", Bucket: "uploads"},
		Mail:         MailConfig{Mode: MailResend, From: "noreply@example.com", APIKey: "re_supersecret"},
	}
	var text strings.Builder
	for _, item := range opts.Summary() {
		text.WriteString(item.Label + " " + item.Value + " " + item.Detail + "\n")
	}
	for _, secret := range []string{"supersecret", "AKIASECRET", "shhh", "re_supersecret"} {
		if strings.Contains(text.String(), secret) {
			t.Errorf("summary leaked %q:\n%s", secret, text.String())
		}
	}
}

func TestStarterRemovalNeeded(t *testing.T) {
	root := t.TempDir()
	model := filepath.Join(root, "api", "internal", "models", "blog.go")
	if err := os.MkdirAll(filepath.Dir(model), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(model, []byte("package models"), 0o644); err != nil {
		t.Fatal(err)
	}

	blank := Options{Architecture: "single", Starter: StarterBlank}
	if !blank.StarterRemovalNeeded(root) {
		t.Error("blank should remove the Blog model")
	}
	blog := Options{Architecture: "single", Starter: StarterBlog}
	if blog.StarterRemovalNeeded(root) {
		t.Error("blog should keep the Blog model")
	}
	if err := os.Remove(model); err != nil {
		t.Fatal(err)
	}
	if blank.StarterRemovalNeeded(root) {
		t.Error("blank should not ask to remove a model that is already gone")
	}
}
