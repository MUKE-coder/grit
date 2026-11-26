package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// maxRetention is the longest window --retain accepts. Past it the value is
// not a retention policy but "keep forever" spelled as a number, and a large
// enough one overflows time.Duration (about 292 years) into a negative window,
// which would delete every row on the first sweep.
const maxRetention = 100 * 365 * 24 * time.Hour

// ParseRetention reads a retention window as written on the command line or in
// a definition file: 90d, 2w, 36h.
//
// Not time.ParseDuration, because it has no unit longer than an hour, and a
// privacy policy is written in days. Months and years are refused rather than
// guessed: a month is 28 to 31 days, and a policy that promises "6 months"
// and deletes at 180 days has deleted early, which is the one direction a
// retention promise cannot be wrong in.
func ParseRetention(s string) (time.Duration, error) {
	v := strings.ToLower(strings.TrimSpace(s))
	if v == "" {
		return 0, fmt.Errorf("--retain needs a window, such as 90d, 2w or 36h")
	}
	unitAt := len(v) - 1
	n, err := strconv.Atoi(v[:unitAt])
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("--retain %q is not a window: use a whole number of days, "+
			"weeks or hours, such as 90d, 2w or 36h", s)
	}
	var unit time.Duration
	switch v[unitAt] {
	case 'h':
		unit = time.Hour
	case 'd':
		unit = 24 * time.Hour
	case 'w':
		unit = 7 * 24 * time.Hour
	case 'm', 'y':
		return 0, fmt.Errorf("--retain %q: months and years have no fixed length, so "+
			"write the window in days (180d, 365d) and it means exactly what the policy says", s)
	default:
		return 0, fmt.Errorf("--retain %q has an unknown unit: use d, w or h, such as 90d, 2w or 36h", s)
	}
	// Compared before multiplying, so a huge n cannot wrap past the check.
	if time.Duration(n) > maxRetention/unit {
		return 0, fmt.Errorf("--retain %q is longer than 100 years: to keep rows forever, "+
			"leave --retain off", s)
	}
	return time.Duration(n) * unit, nil
}

// RetentionWindow is how long this resource's rows are kept, or zero when the
// resource keeps them until somebody deletes them.
func (d *ResourceDefinition) RetentionWindow() (time.Duration, error) {
	if d.Retain == "" {
		return 0, nil
	}
	return ParseRetention(d.Retain)
}

// RetentionSpec is the window written the way it was asked for, tidied: " 90D"
// and "090d" are both "90d". It is what the model comment and the admin
// definition show, so neither repeats a typo back at the reader. Empty when
// the resource has no window or it does not parse.
func (d *ResourceDefinition) RetentionSpec() string {
	window, err := d.RetentionWindow()
	if err != nil || window == 0 {
		return ""
	}
	v := strings.ToLower(strings.TrimSpace(d.Retain))
	unit := map[byte]time.Duration{'h': time.Hour, 'd': 24 * time.Hour, 'w': 7 * 24 * time.Hour}[v[len(v)-1]]
	return fmt.Sprintf("%d%c", window/unit, v[len(v)-1])
}

// checkRetention refuses a retention window the generated code could not keep,
// before anything is written.
func (d *ResourceDefinition) checkRetention() error {
	if d.Retain == "" {
		return nil
	}
	if _, err := d.RetentionWindow(); err != nil {
		return err
	}
	// Append-only rows are refused deletion by a GORM guard and a database
	// trigger, so a sweep would fail every night while the policy claimed the
	// rows were going. Which of the two wins is the project's decision, not
	// one for the generator to make quietly, which is also why erasure leaves
	// append-only resources alone (see erasesWithOwner).
	if d.AppendOnly {
		return fmt.Errorf("--retain and --append-only cannot be combined: append-only " +
			"rows can never be deleted, so the retention sweep could not remove them. " +
			"Generate the resource with one or the other")
	}
	return nil
}

// retentionModelInit registers the model's window with the retention package
// from its own file, so removing the resource removes its policy with it.
func (g *Generator) retentionModelInit(names Names) string {
	window, err := g.Definition.RetentionWindow()
	if err != nil || window == 0 {
		return ""
	}
	// What the options mean, said where somebody reading the model will see
	// them: lines that are younger than the window still go, and a parent past
	// it may stay.
	var shape string
	if g.Definition.Items != nil {
		shape += fmt.Sprintf("\n// Its %s rows go with it, whatever their own age.", BuildNames(g.Definition.Items).Pascal)
	}
	if g.Definition.Tree {
		shape += "\n// A row stays while any row below it is still inside the window."
	}
	return fmt.Sprintf(`
// init keeps %[1]s rows for %[2]s after they are created (--retain %[2]s): the
// nightly retention sweep removes older ones, and releases their files to the
// orphan upload sweep. See internal/retention.%[5]s
func init() { retention.Register(&%[1]s{}, %[3]s%[4]s) }
`, names.Pascal, g.Definition.RetentionSpec(), retentionGoExpr(window), g.retentionOptions(names), shape)
}

// retentionOptions is what the resource's shape adds to its policy. Lines from
// --items go with their record, by the foreign key the parent's has-many is
// declared on; a --tree goes leaf first. See internal/retention.
func (g *Generator) retentionOptions(names Names) string {
	var opts string
	if g.Definition.Items != nil {
		opts += fmt.Sprintf(", retention.WithLines(&%s{}, %q)", BuildNames(g.Definition.Items).Pascal, names.Snake+"_id")
	}
	if g.Definition.Tree {
		opts += ", retention.AsTree()"
	}
	return opts
}

// retentionGoExpr writes a window as Go a person can read back as the policy:
// 90*24*time.Hour rather than 2160*time.Hour. Weeks become days, because
// time has no unit for either and days is what a policy is written in.
func retentionGoExpr(window time.Duration) string {
	day := 24 * time.Hour
	if window%day == 0 {
		return fmt.Sprintf("%d*24*time.Hour", window/day)
	}
	return fmt.Sprintf("%d*time.Hour", window/time.Hour)
}

// prepareRetention makes sure the retention package is there for the model to
// import, and says so when the project has no sweep to enforce the window.
func (g *Generator) prepareRetention() error {
	api := g.APIRoot()
	if err := scaffold.WriteRetentionPackage(api, g.Module, false); err != nil {
		return err
	}
	// A project from before --items and --tree could be retained has the
	// package without the options, and WriteRetentionPackage leaves an
	// existing one alone. The model would then not compile ("too many
	// arguments in call to retention.Register"), so this refuses before
	// anything is written, which is the one moment it is cheap to.
	needs := map[string]bool{"WithLines": g.Definition.Items != nil, "AsTree": g.Definition.Tree}
	read := func(rel ...string) string {
		b, _ := os.ReadFile(filepath.Join(append([]string{api, "internal"}, rel...)...))
		return string(b)
	}
	registry, sweep := read("retention", "retention.go"), read("services", "retention.go")
	for option, needed := range needs {
		if needed && !strings.Contains(registry, "func "+option+"(") {
			return fmt.Errorf("--retain with --items or --tree needs retention.%s, which this "+
				"project's internal/retention predates: run grit upgrade, then generate again", option)
		}
	}
	// The registry can be new while the sweep is a copy the project edited and
	// upgrade left alone. It compiles, and then ignores the lines and the tree,
	// so the database refuses the delete every night.
	if (needs["WithLines"] && !strings.Contains(sweep, "p.Lines")) || (needs["AsTree"] && !strings.Contains(sweep, "p.Tree")) {
		if fileExists(filepath.Join(api, "internal", "services", "retention.go")) {
			fmt.Println("  ⚠ This project's retention sweep predates --items and --tree, so it will")
			fmt.Println("    fail on these rows every night until it is updated: see grit upgrade --diff.")
		}
	}
	// The model compiles with only the package, so nothing else would notice
	// that no job reads it. A window nobody enforces is worse than none,
	// because the privacy policy says the rows are going.
	if !fileExists(filepath.Join(api, "internal", "services", "retention.go")) {
		fmt.Println("  ⚠ This project has no retention sweep yet, so nothing will remove these")
		fmt.Println("    rows until you run grit upgrade.")
	}
	return nil
}
