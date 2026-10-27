package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// expertPack is the repository holding the Grit expert skills.
const expertPack = "MUKE-coder/grit-expert-skills"

// externalSkills are the ones the experts use when they are present and work
// around when they are not. Separate projects with their own licences and
// release cadence, so they are installed from source rather than vendored.
var externalSkills = []struct{ repo, skill, why string }{
	{"https://github.com/vercel-labs/skills", "find-skills",
		"finding a skill for something the pack does not cover"},
	{"https://github.com/mattpocock/skills", "grill-me",
		"stress-testing a plan before it becomes code"},
	{"https://github.com/vercel-labs/agent-browser", "agent-browser",
		"standing the app up and looking at it"},
	{"https://github.com/ihlamury/design-skills", "dribbble-ui-skills",
		"researching a design pattern instead of recalling one"},
	{"https://github.com/lakr233/dribblecli", "dribbble-reference-harvester",
		"collecting design references to work from"},
}

// skillsCmd installs the expert skills an agent uses to build with Grit.
//
// This exists because an agent cannot install them itself any other way. The
// skills ship as a Claude Code plugin, and `/plugin` is a slash command typed
// by a person: it is not reachable from a tool call, and in the VS Code
// extension it does not exist at all, which answers
// "/plugin isn't available in this environment".
//
// An ordinary CLI command is reachable. So the project brief can say "run
// grit skills install" and the agent runs it, instead of stopping to ask
// somebody to type three lines.
func skillsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Install the Grit expert skills an AI assistant builds with",
		Long: `The Grit expert skills: a router that reads your project and assigns the work,
one expert per architecture, and specialists for resources, the admin panel,
design, data, security, performance, deployment and review, plus one that reads
what the work cost in tokens and one that writes up what was worth sharing.

They read the docs out of this binary rather than from memory, interview you
once and remember the answers, and file what Grit got wrong as an issue while
they work.

  grit skills install          the pack, at user level
  grit skills install --all    the pack and the five skills it works better with
  grit skills status           what is installed`,
	}
	cmd.AddCommand(skillsInstallCmd())
	cmd.AddCommand(skillsStatusCmd())
	return cmd
}

func skillsInstallCmd() *cobra.Command {
	var all, project bool
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install the expert skills",
		RunE: func(cmd *cobra.Command, _ []string) error {
			dim := color.New(color.FgHiBlack)
			green := color.New(color.FgHiGreen)
			cyan := color.New(color.FgHiCyan)

			if err := requireNpx(); err != nil {
				return err
			}

			scope := []string{"-g"}
			where := "user level, so every project sees them"
			if project {
				scope = nil
				where = "this project only"
			}

			fmt.Println()
			dim.Printf("  → Installing the Grit experts (%s)...\n", where)
			dim.Printf("    The installer writes to every agent it finds. A ✗ naming an\n")
			dim.Printf("    agent you do not use, such as PromptScript, is not your problem.\n")
			if err := runSkillsAdd(append([]string{expertPack}, scope...)); err != nil {
				return err
			}

			// Count what actually landed rather than trusting the exit code.
			// The installer targets several agents and succeeds overall while
			// failing for one of them, so "it exited 0" is not the same as
			// "the skills are there".
			dirs := skillDirs()
			var absent []string
			for _, name := range expertSkillNames {
				if !skillInstalled(dirs, name) {
					absent = append(absent, name)
				}
			}
			switch {
			case len(absent) == 0:
				green.Printf("  ✓ %d experts and their references\n", len(expertSkillNames))
			case len(absent) < len(expertSkillNames):
				color.New(color.FgHiYellow).Printf("  ! %d of %d installed. Missing: %s\n",
					len(expertSkillNames)-len(absent), len(expertSkillNames),
					strings.Join(absent, ", "))
			default:
				return fmt.Errorf("the installer reported success and nothing was installed\n\n"+
					"Looked in: %s\n\nTry by hand:\n    npx skills add %s -g -y",
					strings.Join(dirs, ", "), expertPack)
			}

			if all {
				for _, s := range externalSkills {
					dim.Printf("  → %s (%s)...\n", s.skill, s.why)
					args := append([]string{s.repo, "--skill", s.skill}, scope...)
					if err := runSkillsAdd(args); err != nil {
						// Not fatal. Every expert that wants one of these says
						// so and falls back to a named alternative, so a
						// failure here is a weaker pack and not a broken one.
						color.New(color.FgHiYellow).Printf("  ! %s did not install: %v\n", s.skill, err)
						continue
					}
					green.Printf("  ✓ %s\n", s.skill)
				}
			}

			fmt.Println()
			cyan.Printf("  Start a new session, then:\n\n")
			fmt.Printf("      /grit-setup      check what the experts depend on\n")
			fmt.Printf("      /grit-build      read the project and assign the work\n\n")
			if !all {
				dim.Printf("  The experts also use five skills from other projects for design\n")
				dim.Printf("  research and browser testing. Add them with --all.\n\n")
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&all, "all", false,
		"also install the five external skills the experts use")
	cmd.Flags().BoolVar(&project, "project", false,
		"install into this project rather than at user level")
	return cmd
}

func skillsStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show which expert skills are installed",
		RunE: func(cmd *cobra.Command, _ []string) error {
			green := color.New(color.FgHiGreen)
			dim := color.New(color.FgHiBlack)

			dirs := skillDirs()
			if len(dirs) == 0 {
				fmt.Println("\n  No skill directory found. Run: grit skills install")
				return nil
			}

			fmt.Println()
			missing := 0
			for _, name := range expertSkillNames {
				if skillInstalled(dirs, name) {
					green.Printf("  ✓ %s\n", name)
					continue
				}
				dim.Printf("  · %s\n", name)
				missing++
			}
			fmt.Println()
			for _, s := range externalSkills {
				if skillInstalled(dirs, s.skill) {
					green.Printf("  ✓ %s\n", s.skill)
				} else {
					dim.Printf("  · %s  (%s)\n", s.skill, s.why)
				}
			}
			fmt.Println()
			if missing > 0 {
				fmt.Printf("  %d of the experts are not installed. Run: grit skills install\n\n", missing)
			}
			return nil
		},
	}
}

// expertSkillNames is the pack, for reporting. Kept here rather than fetched
// so `status` works with no network.
var expertSkillNames = []string{
	"grit-build", "grit-core", "grit-resource", "grit-data",
	"grit-triple", "grit-double", "grit-single", "grit-desktop", "grit-mobile",
	"grit-offline", "grit-admin", "grit-ui-ux", "grit-security", "grit-perf",
	"grit-deploy", "grit-reviewer", "grit-feedback", "grit-setup", "grit-report",
	"grit-tokens", "grit-advocate",
	"grit-references",
}

// skillDirs is every place a skill may be installed, newest convention first.
func skillDirs() []string {
	var out []string
	if home, err := os.UserHomeDir(); err == nil {
		for _, rel := range [][]string{
			{".claude", "skills"},
			{".agents", "skills"},
		} {
			d := filepath.Join(append([]string{home}, rel...)...)
			if dirExistsAt(d) {
				out = append(out, d)
			}
		}
	}
	if cwd, err := os.Getwd(); err == nil {
		if d := filepath.Join(cwd, ".claude", "skills"); dirExistsAt(d) {
			out = append(out, d)
		}
	}
	return out
}

func skillInstalled(dirs []string, name string) bool {
	for _, d := range dirs {
		// Lstat, not Stat: the skills CLI installs a symlink, and a broken one
		// is still a thing somebody installed and should see reported.
		if _, err := os.Lstat(filepath.Join(d, name)); err == nil {
			return true
		}
	}
	return false
}

func dirExistsAt(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

// runSkillsAdd shells out to the skills CLI, which is what actually knows how
// to install into every agent's directory layout.
func runSkillsAdd(args []string) error {
	full := append([]string{"-y", "skills", "add"}, args...)
	full = append(full, "-y")
	c := exec.Command(npxName(), full...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("npx skills add %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

func npxName() string {
	if runtime.GOOS == "windows" {
		return "npx.cmd"
	}
	return "npx"
}

// requireNpx fails with something actionable rather than "exec: not found".
func requireNpx() error {
	if _, err := exec.LookPath(npxName()); err == nil {
		return nil
	}
	// On Windows the plain name resolves through PATHEXT in a shell but not
	// through LookPath, so try both before giving up.
	if _, err := exec.LookPath("npx"); err == nil {
		return nil
	}
	return fmt.Errorf("npx was not found, and the skills are installed through it\n\n" +
		"Install Node.js (which ships npx), then run this again:\n" +
		"    https://nodejs.org\n\n" +
		"Or install them by hand:\n" +
		"    npx skills add " + expertPack + " -g -y")
}
