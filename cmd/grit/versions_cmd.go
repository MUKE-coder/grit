package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/MUKE-coder/grit/v3/internal/plugin"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// versionsCmd answers "how far behind is this project, and in what".
//
// grit doctor already warns about a library below its floor, but a warning
// only fires when something is wrong. Somebody deciding whether to adopt a
// newer Grit in an existing project needs the opposite: the whole picture,
// including the parts that are fine, because what they are weighing is the
// size of the move rather than any single fault.
func versionsCmd() *cobra.Command {
	var asJSON bool

	cmd := &cobra.Command{
		Use:     "versions",
		Aliases: []string{"version-report"},
		Short:   "What this project is running, against what this CLI ships",
		Long: "Reports the project's Grit version, the libraries every generated API\n" +
			"mounts, and the installed plugins, each against the version this CLI\n" +
			"expects.\n\n" +
			"Nothing is changed. `grit upgrade` is what raises them.",
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := scaffold.FindProjectRoot()
			if err != nil {
				return err
			}
			report := buildVersionReport(root)

			if asJSON {
				out, err := json.MarshalIndent(report, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(out))
				return nil
			}
			printVersionReport(report)
			return nil
		},
	}

	cmd.Flags().BoolVar(&asJSON, "json", false, "Print the report as JSON")
	return cmd
}

// VersionLine is one row of the report.
type VersionLine struct {
	Name      string `json:"name"`
	Installed string `json:"installed"`
	Expected  string `json:"expected"`
	// Status is current, behind, ahead or absent.
	Status string `json:"status"`
	Why    string `json:"why,omitempty"`
}

// VersionReport is what `grit versions` found.
type VersionReport struct {
	ProjectVersion string        `json:"project_version"`
	CLIVersion     string        `json:"cli_version"`
	Architecture   string        `json:"architecture,omitempty"`
	Batteries      []VersionLine `json:"batteries"`
	Plugins        []VersionLine `json:"plugins"`
	Behind         int           `json:"behind"`
}

func buildVersionReport(root string) VersionReport {
	report := VersionReport{
		ProjectVersion: scaffold.ProjectVersion(root),
		CLIVersion:     version,
		Architecture:   architectureOf(root),
	}

	goMod := readTextFile(filepath.Join(apiRootOf(root), "go.mod"))
	for _, dep := range scaffold.FrameworkDepFloors() {
		line := VersionLine{Name: dep.Path, Expected: dep.Floor, Why: dep.Why}

		re := regexp.MustCompile(`(?m)^\s*(?:require\s+)?` + regexp.QuoteMeta(dep.Path) + `\s+(v\S+)`)
		m := re.FindStringSubmatch(goMod)
		switch {
		case m == nil:
			// Not an error. A floor says "if you have this, do not run it
			// below this"; a project that never mounted the library is
			// entitled not to have it.
			line.Status = "absent"
		case scaffold.VersionBelow(m[1], dep.Floor):
			line.Installed, line.Status = m[1], "behind"
			report.Behind++
		case m[1] == dep.Floor:
			line.Installed, line.Status = m[1], "current"
		default:
			line.Installed, line.Status = m[1], "ahead"
		}
		report.Batteries = append(report.Batteries, line)
	}

	if lock, err := plugin.LoadLock(root); err == nil {
		for _, installed := range lock.Plugins {
			line := VersionLine{Name: installed.Name, Installed: installed.Version}
			shipped, err := plugin.Get(installed.Name)
			switch {
			case err != nil:
				// A community plugin, or one installed from a directory. This
				// CLI has no version to compare it against and says so rather
				// than implying it is current.
				line.Expected, line.Status = "", "unknown"
			case shipped.Version == installed.Version:
				line.Expected, line.Status = shipped.Version, "current"
			case scaffold.VersionBelow("v"+installed.Version, "v"+shipped.Version):
				line.Expected, line.Status = shipped.Version, "behind"
				report.Behind++
			default:
				line.Expected, line.Status = shipped.Version, "ahead"
			}
			report.Plugins = append(report.Plugins, line)
		}
	}

	if report.ProjectVersion != "" && scaffold.VersionBelow("v"+report.ProjectVersion, "v"+version) {
		report.Behind++
	}
	return report
}

func printVersionReport(r VersionReport) {
	printLogo()
	faint := color.New(color.Faint)

	fmt.Println("  Grit")
	project := r.ProjectVersion
	if project == "" {
		// Projects generated before the version stamp existed. Saying so is
		// more useful than printing nothing, because it dates the project.
		project = "not recorded"
	}
	fmt.Printf("    %-26s %s\n", "this project", project)
	fmt.Printf("    %-26s %s\n", "this CLI", r.CLIVersion)
	if r.Architecture != "" {
		fmt.Printf("    %-26s %s\n", "architecture", faint.Sprint(r.Architecture))
	}
	if r.ProjectVersion != "" && r.ProjectVersion != r.CLIVersion {
		fmt.Printf("    %s\n", faint.Sprint("grit upgrade moves the framework files to "+r.CLIVersion))
	}
	fmt.Println()

	fmt.Println("  Batteries")
	fmt.Println("  " + faint.Sprint("the libraries every generated API mounts") + "\n")
	for _, line := range r.Batteries {
		printVersionLine(line, shortModulePath(line.Name))
	}
	fmt.Println()

	fmt.Println("  Plugins")
	if len(r.Plugins) == 0 {
		fmt.Println("  " + faint.Sprint("none installed") + "\n")
	} else {
		fmt.Println()
		for _, line := range r.Plugins {
			printVersionLine(line, line.Name)
		}
		fmt.Println()
	}

	switch {
	case r.Behind == 0:
		fmt.Println("  " + color.GreenString("Everything is at or above what this CLI expects."))
	case r.Behind == 1:
		fmt.Println("  " + color.YellowString("1 thing is behind.") + " " + faint.Sprint("grit upgrade"))
	default:
		fmt.Println("  " + color.YellowString(fmt.Sprintf("%d things are behind.", r.Behind)) +
			" " + faint.Sprint("grit upgrade"))
	}
	fmt.Println()
	fmt.Println("  " + faint.Sprint("grit doctor   what is actually wrong, rather than merely old"))
	fmt.Println()
}

func printVersionLine(line VersionLine, label string) {
	faint := color.New(color.Faint)

	var status, detail string
	switch line.Status {
	case "behind":
		status = color.YellowString("behind")
		detail = line.Installed + " -> " + line.Expected
	case "current":
		status = color.GreenString("current")
		detail = line.Installed
	case "ahead":
		status = color.New(color.FgCyan).Sprint("ahead")
		detail = line.Installed + " (expects " + line.Expected + ")"
	case "absent":
		status = faint.Sprint("absent")
		detail = faint.Sprint("not in go.mod")
	default:
		status = faint.Sprint("unknown")
		detail = line.Installed + faint.Sprint(" (not a plugin this CLI ships)")
	}

	fmt.Printf("    %-26s %-9s %s\n", label, status, detail)

	// The reason only earns its space when it is a reason to act.
	if line.Status == "behind" && line.Why != "" {
		fmt.Printf("    %-26s %s\n", "", faint.Sprint(firstSentence(line.Why)))
	}
}

// shortModulePath drops the host and owner: github.com/MUKE-coder/pulse reads
// as pulse in a column of names.
func shortModulePath(path string) string {
	parts := strings.Split(path, "/")
	name := parts[len(parts)-1]
	// A major-version suffix is not the name.
	if regexp.MustCompile(`^v\d+$`).MatchString(name) && len(parts) > 1 {
		name = parts[len(parts)-2]
	}
	return name
}

// firstSentence keeps a one-line reason out of a paragraph-long one.
func firstSentence(s string) string {
	for i, r := range s {
		if r == ',' || r == ':' {
			return strings.TrimSpace(s[:i])
		}
	}
	if len(s) > 90 {
		return strings.TrimSpace(s[:90]) + "..."
	}
	return s
}

func architectureOf(root string) string {
	body, err := os.ReadFile(filepath.Join(root, "grit.json"))
	if err != nil {
		return ""
	}
	var cfg struct {
		Architecture string `json:"architecture"`
	}
	if json.Unmarshal(body, &cfg) != nil {
		return ""
	}
	return cfg.Architecture
}

func apiRootOf(root string) string {
	mono := filepath.Join(root, "apps", "api")
	if _, err := os.Stat(filepath.Join(mono, "go.mod")); err == nil {
		return mono
	}
	return root
}

func readTextFile(path string) string {
	body, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(body)
}
