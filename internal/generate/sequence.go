package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SequenceOptions describes a single sequence generation request.
//
// Name is the resource being numbered (e.g. "Invoice"). Prefix is the
// alphabetic prefix in the output (e.g. "INV"); when empty it defaults to
// the first three uppercase letters of Name. Reset is "monthly" / "yearly"
// / "never" — anything else is treated as "monthly". Width is the
// zero-padded width of the numeric portion (default 4 → INV-202605-0001).
type SequenceOptions struct {
	Name   string
	Prefix string
	Reset  string
	Width  int
}

// GenerateSequence runs the sequence generator. On the first invocation
// across the project it creates the shared internal/sequence package and
// registers the Counter model with AutoMigrate. Every invocation also
// writes a per-resource convenience wrapper at internal/services/<name>_sequence.go.
func GenerateSequence(opts SequenceOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("sequence name is required (e.g. grit generate sequence Invoice)")
	}
	if opts.Width <= 0 {
		opts.Width = 4
	}
	if opts.Reset == "" {
		opts.Reset = "monthly"
	}
	switch opts.Reset {
	case "monthly", "yearly", "never":
	default:
		return fmt.Errorf("--reset must be monthly, yearly, or never (got %q)", opts.Reset)
	}
	if opts.Prefix == "" {
		opts.Prefix = defaultPrefix(opts.Name)
	}

	root, err := findProjectRoot()
	if err != nil {
		return err
	}
	arch, _ := readGritJSON(root)
	module, err := readModulePath(root, arch)
	if err != nil {
		return err
	}

	apiRoot := root
	apiPrefix := ""
	if arch != "single" {
		apiRoot = filepath.Join(root, "apps", "api")
		apiPrefix = "apps/api/"
	}

	names := MakeNames(opts.Name)

	fmt.Printf("\n  Generating sequence: %s (prefix=%s, reset=%s)\n\n", names.Pascal, opts.Prefix, opts.Reset)

	sequencePkgPath := filepath.Join(apiRoot, "internal", "sequence", "sequence.go")
	wrapperPath := filepath.Join(apiRoot, "internal", "services", strings.ToLower(names.Snake)+"_sequence.go")

	// Sequence package — write only on first invocation.
	if !fileExists(sequencePkgPath) {
		if err := writeFileWithDirs(sequencePkgPath, sequencePackageGo()); err != nil {
			return fmt.Errorf("writing sequence package: %w", err)
		}
		fmt.Printf("  ✓ %sinternal/sequence/sequence.go\n", apiPrefix)

		// Inject the Counter model into AutoMigrate.
		if err := injectSequenceMigrate(apiRoot, module); err != nil {
			return fmt.Errorf("injecting AutoMigrate: %w", err)
		}
		fmt.Printf("  ✓ Registered sequence.Counter with AutoMigrate\n")
	} else {
		fmt.Printf("  • internal/sequence/sequence.go already exists (skipped)\n")
	}

	// Per-resource convenience wrapper. Always overwrite — the wrapper is
	// derived purely from opts and is safe to regenerate if the user
	// changes prefix / reset / width.
	wrapper := sequenceWrapperGo(module, names, opts)
	if err := writeFileWithDirs(wrapperPath, wrapper); err != nil {
		return fmt.Errorf("writing wrapper: %w", err)
	}
	fmt.Printf("  ✓ %sinternal/services/%s_sequence.go\n", apiPrefix, names.Snake)

	fmt.Printf("\n  ✅ Sequence generated.\n\n")
	fmt.Printf("  Use it from a handler:\n\n")
	fmt.Printf("      number, err := services.Next%sNumber(h.DB, time.Now())\n", names.Pascal)
	fmt.Printf("      // → %s-202605-%s\n\n", opts.Prefix, strings.Repeat("0", opts.Width-1)+"1")

	return nil
}

func defaultPrefix(name string) string {
	upper := strings.ToUpper(name)
	if len(upper) >= 3 {
		return upper[:3]
	}
	return upper
}

// sequencePackageGo returns the source for the shared sequence package.
// The package is written once per project; subsequent generator calls
// only add per-resource wrappers in services/.
func sequencePackageGo() string {
	return `// Package sequence generates atomic, gap-free sequential numbers like
// INV-202605-0001 backed by a counter table. Each (name, bucket) pair
// gets its own row; "bucket" is the period key for resets — "202605"
// for monthly, "2026" for yearly, "" for no reset.
//
// Use it through a generated wrapper in services/, e.g.:
//
//	number, err := services.NextInvoiceNumber(db, time.Now())
//
// Concurrency note: Next() runs inside a transaction with a row-level
// SELECT FOR UPDATE on Postgres so concurrent callers serialize on the
// counter row. SQLite serializes writes globally. For high-throughput
// numbering (>100 writes/sec on a single counter), consider a dedicated
// sequence service.
package sequence

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Reset controls the bucket key used to scope a counter.
type Reset string

const (
	ResetMonthly Reset = "monthly"
	ResetYearly  Reset = "yearly"
	ResetNever   Reset = "never"
)

// Counter is the row backing a single (name, bucket) sequence. Add it
// to AutoMigrate once: db.AutoMigrate(&sequence.Counter{}).
type Counter struct {
	Name      string    ` + "`" + `gorm:"primaryKey;size:50"` + "`" + `
	Bucket    string    ` + "`" + `gorm:"primaryKey;size:20"` + "`" + `
	NextValue uint64    ` + "`" + `gorm:"not null;default:0"` + "`" + `
	UpdatedAt time.Time
}

// TableName pins the counter table to a stable name across resources
// that might have a model with the same default plural ("counters").
func (Counter) TableName() string {
	return "sequence_counters"
}

// Config describes how to format a sequence number.
type Config struct {
	Name   string // logical sequence name (e.g. "invoice")
	Prefix string // alphabetic prefix in the output (e.g. "INV")
	Reset  Reset  // monthly | yearly | never
	Width  int    // zero-padded width of the numeric portion (default 4)
}

// Next returns the next number for cfg, atomically incrementing the
// counter row inside a transaction. The bucket is derived from t and
// cfg.Reset.
func Next(db *gorm.DB, cfg Config, t time.Time) (string, error) {
	width := cfg.Width
	if width <= 0 {
		width = 4
	}
	bucket := bucketKey(t, cfg.Reset)

	var next uint64
	err := db.Transaction(func(tx *gorm.DB) error {
		// Lock or create the counter row, then increment in one go.
		// The row-level lock ensures concurrent callers serialize.
		var c Counter
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("name = ? AND bucket = ?", cfg.Name, bucket).
			First(&c).Error
		if err == gorm.ErrRecordNotFound {
			c = Counter{Name: cfg.Name, Bucket: bucket, NextValue: 1}
			next = 1
			return tx.Create(&c).Error
		}
		if err != nil {
			return err
		}
		c.NextValue++
		next = c.NextValue
		return tx.Save(&c).Error
	})
	if err != nil {
		return "", err
	}
	return formatNumber(cfg.Prefix, bucket, width, next), nil
}

func bucketKey(t time.Time, reset Reset) string {
	switch reset {
	case ResetMonthly:
		return fmt.Sprintf("%04d%02d", t.Year(), int(t.Month()))
	case ResetYearly:
		return fmt.Sprintf("%04d", t.Year())
	default:
		return ""
	}
}

func formatNumber(prefix, bucket string, width int, n uint64) string {
	parts := []string{prefix}
	if bucket != "" {
		parts = append(parts, bucket)
	}
	parts = append(parts, fmt.Sprintf("%0*d", width, n))
	return strings.Join(parts, "-")
}
`
}

// sequenceWrapperGo emits a tiny Next<Name>Number(db, t) helper that
// closes over the configured prefix/reset/width so handlers don't need
// to know any of that.
func sequenceWrapperGo(module string, names Names, opts SequenceOptions) string {
	return fmt.Sprintf(`package services

import (
	"time"

	"gorm.io/gorm"

	"%s/internal/sequence"
)

// Next%sNumber returns the next %s identifier (e.g. %s-202605-%s).
// The counter resets %s and is generated atomically.
func Next%sNumber(db *gorm.DB, t time.Time) (string, error) {
	return sequence.Next(db, sequence.Config{
		Name:   "%s",
		Prefix: "%s",
		Reset:  sequence.Reset%s,
		Width:  %d,
	}, t)
}
`,
		module,
		names.Pascal,
		names.Lower,
		opts.Prefix,
		strings.Repeat("0", opts.Width-1)+"1",
		opts.Reset,
		names.Pascal,
		strings.ToLower(names.Snake),
		opts.Prefix,
		titleReset(opts.Reset),
		opts.Width,
	)
}

func titleReset(r string) string {
	switch r {
	case "monthly":
		return "Monthly"
	case "yearly":
		return "Yearly"
	default:
		return "Never"
	}
}

// injectSequenceMigrate adds &sequence.Counter{} to the Models() slice
// in internal/models/user.go (where the // grit:models marker lives).
// Idempotent: skipped if Counter is already in the slice.
func injectSequenceMigrate(apiRoot, module string) error {
	modelFile := filepath.Join(apiRoot, "internal", "models", "user.go")
	data, err := os.ReadFile(modelFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", modelFile, err)
	}
	content := string(data)

	if strings.Contains(content, "sequence.Counter{}") {
		return nil
	}

	// Add the import if it's missing. user.go uses a single import block
	// with the standard imports first, then external. We slot the sequence
	// import after the existing module-internal imports if any, otherwise
	// after gorm.
	importLine := fmt.Sprintf("\t\"%s/internal/sequence\"", module)
	if !strings.Contains(content, importLine) {
		anchor := "\"gorm.io/gorm\""
		if idx := strings.Index(content, anchor); idx >= 0 {
			lineEnd := strings.IndexByte(content[idx:], '\n')
			if lineEnd > 0 {
				insertAt := idx + lineEnd + 1
				content = content[:insertAt] + "\n" + importLine + "\n" + content[insertAt:]
			}
		}
	}

	if err := os.WriteFile(modelFile, []byte(content), 0o644); err != nil {
		return err
	}

	return injectBefore(modelFile, "// grit:models", "\t\t&sequence.Counter{},")
}
