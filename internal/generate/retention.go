package generate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
