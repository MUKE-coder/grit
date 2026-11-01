package setup

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ErrDockerUnavailable means Docker or Docker Compose is not on PATH.
var ErrDockerUnavailable = errors.New("docker compose is not available")

// ErrPnpmUnavailable means pnpm is not on PATH.
var ErrPnpmUnavailable = errors.New("pnpm is not on PATH")

// CommandError is a grit command that ran and failed. Output is what it wrote
// before exiting; setup shows the tail of it, because the last lines of a
// migration failure are the part that says what went wrong.
type CommandError struct {
	Err    error
	Output string
}

func (e *CommandError) Error() string { return e.Err.Error() }
func (e *CommandError) Unwrap() error { return e.Err }

// selfExecutable returns the running grit binary. os.Executable is the binary
// the person actually invoked, including a rebuilt one from `go build`; falling
// back to PATH keeps `go run` usable when the temporary executable has already
// been cleaned up.
func selfExecutable() string {
	exe, err := os.Executable()
	if err != nil || exe == "" {
		exe, _ = exec.LookPath("grit")
	}
	return exe
}

// selfCommand returns a grit command with the terminal attached, for the one
// step that should own the terminal: grit start.
func selfCommand(ctx context.Context, root string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, selfExecutable(), args...)
	cmd.Dir = root
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// runSelfOutput runs a grit command with its output captured. The commands
// setup calls this way print their own wordmark and a list of every file they
// touched, which is useful when a person runs them and noise in the middle of
// an installer. The output is kept for the failure path.
func runSelfOutput(ctx context.Context, root string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, selfExecutable(), args...)
	cmd.Dir = root
	// NO_COLOR keeps the captured text readable when it is printed inside a
	// failure block; DB_LOG_LEVEL=silent removes GORM's per-statement SQL.
	cmd.Env = append(os.Environ(), "NO_COLOR=1", "DB_LOG_LEVEL=silent")
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// tailOutput returns the last max lines of a command's output.
func tailOutput(out string, max int) string {
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) > max {
		lines = lines[len(lines)-max:]
	}
	return strings.Join(lines, "\n")
}

// startForeground runs `grit start` and hands it the terminal. It deliberately
// does not use CommandContext: Ctrl+C should reach the dev servers, not kill a
// parent that is waiting to print a summary nobody will read.
func startForeground(root string) error {
	cmd := exec.Command(selfExecutable(), "start")
	cmd.Dir = root
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if errors.As(err, &exit) && exit.ExitCode() == 130 {
			// Ctrl+C is how somebody stops a dev server. It is not a setup
			// failure, and saying "exit status 130" after it helps nobody.
			return nil
		}
		return err
	}
	return nil
}

// composeBase finds the compose command. `docker compose` is the current one;
// `docker-compose` is still installed on plenty of machines.
func composeBase(ctx context.Context, root string) ([]string, error) {
	for _, candidate := range [][]string{{"docker", "compose"}, {"docker-compose"}} {
		cmd := exec.CommandContext(ctx, candidate[0], append(candidate[1:], "version")...)
		cmd.Dir = root
		if err := cmd.Run(); err == nil {
			return candidate, nil
		}
	}
	return nil, ErrDockerUnavailable
}

// StartLocalServices starts only the services this project chose. `--wait`
// makes the next step (migrate) wait for Postgres rather than race its first
// connection; older Compose builds do not know the flag, so the fallback is
// the plain `up -d` they do understand.
func StartLocalServices(ctx context.Context, root string, services []string) error {
	if len(services) == 0 {
		return nil
	}
	base, err := composeBase(ctx, root)
	if err != nil {
		return err
	}

	args := append(append([]string{}, base[1:]...), "up", "-d", "--wait")
	args = append(args, services...)
	cmd := exec.CommandContext(ctx, base[0], args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err == nil {
		return nil
	}

	// A failure here is either the services or the flag. Retrying without
	// --wait distinguishes the two cheaply: if the second run succeeds, the
	// first one was an old Compose.
	args = append(append([]string{}, base[1:]...), "up", "-d")
	args = append(args, services...)
	cmd = exec.CommandContext(ctx, base[0], args...)
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// InstallFrontendDeps runs pnpm where the project put its package.json.
func InstallFrontendDeps(ctx context.Context, root string, ci bool) error {
	if _, err := exec.LookPath("pnpm"); err != nil {
		return ErrPnpmUnavailable
	}
	cmd := exec.CommandContext(ctx, "pnpm", "install")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if ci {
		cmd.Env = append(os.Environ(), "CI=1")
	}
	return cmd.Run()
}

// TidyGoDeps resolves the API's go.sum when the scaffold could not.
func TidyGoDeps(ctx context.Context, apiDir string) error {
	if _, err := exec.LookPath("go"); err != nil {
		return fmt.Errorf("go is not on PATH: install Go, then run `go mod tidy` in %s", apiDir)
	}
	cmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	cmd.Dir = apiDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunMigrate runs the project's own migration command through the CLI, so the
// version marker and the command's output are exactly what a person sees when
// they run it themselves. The output is captured because the installer already
// has a line for the result; it is shown only when the command fails.
func RunMigrate(ctx context.Context, root string) error {
	out, err := runSelfOutput(ctx, root, "migrate")
	if err != nil {
		return &CommandError{Err: err, Output: out}
	}
	return nil
}

// RunSeed is optional: a seed failure is reported, not fatal.
func RunSeed(ctx context.Context, root string) error {
	out, err := runSelfOutput(ctx, root, "seed")
	if err != nil {
		return &CommandError{Err: err, Output: out}
	}
	return nil
}

// RemoveBlog removes the built-in Blog example for the blank starter. It goes
// through `grit remove resource` for the same reason the scaffold does: the
// removal list is the generator's, not a second list setup has to maintain.
func RemoveBlog(ctx context.Context, root string) error {
	out, err := runSelfOutput(ctx, root, "remove", "resource", "Blog", "--force")
	if err != nil {
		return fmt.Errorf("%w\n\n%s", err, tailOutput(out, 30))
	}
	return nil
}

// toolMissing reports whether err means a tool is absent rather than broken.
func toolMissing(err error) bool {
	return errors.Is(err, ErrDockerUnavailable) || errors.Is(err, ErrPnpmUnavailable)
}

// composeServicesLabel is the list shown in a step line.
func composeServicesLabel(services []string) string {
	return strings.Join(services, ", ")
}
