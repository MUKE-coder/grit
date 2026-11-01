package setup

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GitResult says what the git step did, so the caller can print one truthful
// line instead of guessing from a nil error.
type GitResult struct {
	Initialised bool
	Committed   bool
	Skipped     string
}

// InitGit creates a repository at the project root and makes the first commit.
//
// It refuses to touch a folder that is already inside a repository: `grit
// setup my-app` inside a monorepo should not commit the monorepo, and a
// project that already has history has a first commit.
func InitGit(ctx context.Context, root string) (GitResult, error) {
	if _, err := exec.LookPath("git"); err != nil {
		return GitResult{Skipped: "git is not installed"}, nil
	}

	abs, err := filepath.Abs(root)
	if err != nil {
		return GitResult{}, err
	}
	if top, ok := gitTopLevel(ctx, root); ok {
		if samePath(top, abs) {
			return GitResult{Skipped: "this folder is already a git repository"}, nil
		}
		return GitResult{Skipped: "this folder is inside another git repository"}, nil
	}

	if out, err := runGit(ctx, root, "init"); err != nil {
		return GitResult{}, fmt.Errorf("git init: %w\n%s", err, strings.TrimSpace(out))
	}
	// -b main is not on every git in the wild. Initialising and then pointing
	// HEAD at main works on all of them and leaves the same result.
	_, _ = runGit(ctx, root, "symbolic-ref", "HEAD", "refs/heads/main")

	name, _ := runGit(ctx, root, "config", "user.name")
	email, _ := runGit(ctx, root, "config", "user.email")
	if strings.TrimSpace(name) == "" || strings.TrimSpace(email) == "" {
		return GitResult{
			Initialised: true,
			Skipped:     "git has no user.name and user.email; set them, then commit when you are ready",
		}, nil
	}

	// The scaffold's .gitignore already ignores .env, but this step is the one
	// that would put a generated JWT secret in the first commit if that ever
	// changed. Check rather than assume.
	if err := ensureIgnored(ctx, root, ".env"); err != nil {
		return GitResult{}, err
	}

	if out, err := runGit(ctx, root, "add", "-A"); err != nil {
		return GitResult{}, fmt.Errorf("git add: %w\n%s", err, strings.TrimSpace(out))
	}
	if out, err := runGit(ctx, root, "commit", "-m", "Initial commit"); err != nil {
		return GitResult{}, fmt.Errorf("git commit: %w\n%s", err, strings.TrimSpace(out))
	}
	return GitResult{Initialised: true, Committed: true}, nil
}

// StarGrit stars the Grit repository with the GitHub CLI. The caller falls
// back to a link when gh is missing or signed out; a star is a favour, not a
// setup dependency.
func StarGrit(ctx context.Context) error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("the GitHub CLI (gh) is not installed")
	}
	cmd := exec.CommandContext(ctx, "gh", "repo", "star", "MUKE-coder/grit")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}
	text := strings.ToLower(string(out))
	if strings.Contains(text, "already") && strings.Contains(text, "star") {
		return nil
	}
	if strings.TrimSpace(string(out)) == "" {
		return err
	}
	return fmt.Errorf("%s", strings.TrimSpace(string(out)))
}

func gitTopLevel(ctx context.Context, root string) (string, bool) {
	out, err := runGit(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", false
	}
	top := strings.TrimSpace(out)
	if top == "" {
		return "", false
	}
	return top, true
}

func runGit(ctx context.Context, root string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func ensureIgnored(ctx context.Context, root, name string) error {
	cmd := exec.CommandContext(ctx, "git", "check-ignore", "-q", name)
	cmd.Dir = root
	if err := cmd.Run(); err == nil {
		return nil
	}

	path := filepath.Join(root, ".gitignore")
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("reading .gitignore: %w", err)
	}
	text := string(data)
	if !strings.HasSuffix(text, "\n") && text != "" {
		text += "\n"
	}
	text += "# Added by grit setup: never commit generated secrets.\n" + name + "\n"
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}
	return nil
}

func samePath(a, b string) bool {
	aa, errA := filepath.EvalSymlinks(a)
	bb, errB := filepath.EvalSymlinks(b)
	if errA == nil {
		a = aa
	}
	if errB == nil {
		b = bb
	}
	return filepath.Clean(a) == filepath.Clean(b)
}
