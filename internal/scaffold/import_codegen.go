package scaffold

import "strings"

// The CSV importer grit generate writes for a resource resolves a belongs_to
// column by the related record's natural key. These functions write that code,
// for the generator and for the upgrade that brings an older importer up to it,
// so the two cannot drift apart (H19 in the contact-app review).

// ImportLimitGate goes before an importer opens its file: imports take turns.
const ImportLimitGate = `	// Imports take turns: each holds a database connection and writes in
	// batches for its whole run, and a burst of them drained the pool every
	// request shares. See internal/imports.
	release := imports.Wait(func() {
		record(map[string]interface{}{"message": "Waiting for another import to finish"})
	})
	defer release()

`

// ImportNameResolver is the closure that turns a natural key from the CSV into
// the related record's id: base is the column ("group"), relModel the model
// ("Group"), keyJSON and keyGo its natural key ("name", "Name").
func ImportNameResolver(base, relModel, keyJSON, keyGo string) string {
	return strings.NewReplacer(
		"{{name}}", importLowerCamel(base),
		"{{fn}}", "resolve"+importUpperCamel(base),
		"{{base}}", base,
		"{{Rel}}", relModel,
		"{{key}}", keyJSON,
		"{{KeyGo}}", keyGo,
	).Replace(`	// {{name}}IDs remembers the {{Rel}} each {{key}} in the CSV resolved to, so a
	// distinct {{key}} costs one lookup rather than one per row. A missing {{Rel}}
	// is created; any other error fails the row instead of passing for missing.
	{{name}}IDs := map[string]string{}
	{{fn}} := func(v string) (string, error) {
		if id, ok := {{name}}IDs[v]; ok {
			return id, nil
		}
		var rel models.{{Rel}}
		err := db.Where("{{key}} = ?", v).First(&rel).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			rel = models.{{Rel}}{{{KeyGo}}: v}
			if err = db.Create(&rel).Error; err != nil {
				// Another import may have created it in the meantime.
				if again := db.Where("{{key}} = ?", v).First(&rel).Error; again == nil {
					err = nil
				}
			}
		}
		if err != nil {
			return "", fmt.Errorf("{{base}} %q: %w", v, err)
		}
		{{name}}IDs[v] = rel.ID
		return rel.ID, nil
	}

`)
}

// ImportNameAssign is the row code that uses the resolver: a lookup that fails
// fails the row. pointer is for a self-reference, whose column is nullable.
func ImportNameAssign(base, fkGo string, pointer bool) string {
	amp := ""
	if pointer {
		amp = "&"
	}
	return strings.NewReplacer(
		"{{fn}}", "resolve"+importUpperCamel(base),
		"{{base}}", base,
		"{{FK}}", fkGo,
		"{{amp}}", amp,
	).Replace(`		if v, ok := get(rec, "{{base}}"); ok && v != "" {
			id, err := {{fn}}(v)
			if err != nil {
				failed++
				if len(rowErrors) < 50 {
					rowErrors = append(rowErrors, map[string]interface{}{"row": rowNum, "message": err.Error()})
				}
				continue
			}
			item.{{FK}} = {{amp}}id
		}
`)
}

func importUpperCamel(s string) string {
	var b strings.Builder
	for _, part := range strings.Split(s, "_") {
		if part != "" {
			b.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return b.String()
}

func importLowerCamel(s string) string {
	upper := importUpperCamel(s)
	if upper == "" {
		return upper
	}
	return strings.ToLower(upper[:1]) + upper[1:]
}

// importsLimitGo emits internal/imports/limit.go.
func importsLimitGo() string {
	return `// Package imports keeps CSV imports from all running at once.
package imports

import (
	"os"
	"strconv"
	"sync"
)

// Concurrency is how many imports run at the same time, from IMPORT_CONCURRENCY
// (default 2). The others wait their turn.
//
// Each import holds a database connection for its whole run and writes in
// batches, so a burst of them took the connections every request shares.
var Concurrency = concurrencyFromEnv()

var slots = make(chan struct{}, Concurrency)

func concurrencyFromEnv() int {
	if n, err := strconv.Atoi(os.Getenv("IMPORT_CONCURRENCY")); err == nil && n > 0 {
		return n
	}
	return 2
}

// Wait blocks until an import may run, and returns the function that ends its
// turn. onQueued runs first when the import has to wait, so its job can say so.
func Wait(onQueued func()) (release func()) {
	select {
	case slots <- struct{}{}:
	default:
		if onQueued != nil {
			onQueued()
		}
		slots <- struct{}{}
	}
	var once sync.Once
	return func() { once.Do(func() { <-slots }) }
}
`
}

// importsLimitTestGo emits internal/imports/limit_test.go.
func importsLimitTestGo() string {
	return `package imports

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestImportsTakeTurns(t *testing.T) {
	var running, peak, queued atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < Concurrency+3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			release := Wait(func() { queued.Add(1) })
			defer release()
			n := running.Add(1)
			for {
				p := peak.Load()
				if n <= p || peak.CompareAndSwap(p, n) {
					break
				}
			}
			time.Sleep(20 * time.Millisecond)
			running.Add(-1)
		}()
	}
	wg.Wait()
	if int(peak.Load()) > Concurrency {
		t.Errorf("%d imports ran at once; the limit is %d", peak.Load(), Concurrency)
	}
	if queued.Load() == 0 {
		t.Error("no import was told it was waiting")
	}
}

func TestReleasingTwiceFreesOneTurn(t *testing.T) {
	release := Wait(nil)
	release()
	release()
	if got := len(slots); got != 0 {
		t.Errorf("%d turns held after one import released twice", got)
	}
}
`
}
