package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

// ImportUnfillableColumns lists the columns of a related model that a CSV
// importer cannot supply when it has to CREATE that record, which is what
// decides whether creating one is safe at all.
//
// The importer resolves a belongs_to column by the related record's natural
// key, and used to create a missing record from that key alone. For a lookup
// table whose only column is a name that is exactly right. For a model that
// requires anything else it is a constraint violation per row: Company.OwnerID
// is binding:"required" with a foreign key to users, so an imported contact
// naming a company that does not exist yet failed, and in a real contact list
// that is nearly every row. See grit#129.
//
// Required means binding:"required", which is what the generator writes and
// what the API validates against. A column with a gorm default can be left
// out, and so can the natural key itself, which the importer does fill. An
// empty result means creating the record is safe.
//
// Read from the model file on disk rather than from the resource definition,
// because the related model was generated earlier, possibly by hand, possibly
// edited since. The file is the fact.
func ImportUnfillableColumns(apiRoot, relModel, naturalKeyJSON string) []string {
	body, err := os.ReadFile(filepath.Join(apiRoot, "internal", "models", toSnake(relModel)+".go"))
	if err != nil {
		return nil
	}
	return unfillableRequiredColumns(string(body), relModel, naturalKeyJSON)
}

// unfillableRequiredColumns is the part of the above that reads source, split
// out so it can be tested without a project on disk.
func unfillableRequiredColumns(src, relModel, naturalKeyJSON string) []string {
	start := strings.Index(src, "type "+relModel+" struct {")
	if start < 0 {
		return nil
	}
	body := src[start:]
	// The struct ends at the first closing brace in the first column.
	if end := strings.Index(body, "\n}"); end > 0 {
		body = body[:end]
	}

	// Columns the importer fills itself, or the database does.
	skip := map[string]bool{
		"":           true,
		"-":          true,
		"id":         true,
		"created_at": true,
		"updated_at": true,
		"deleted_at": true,
	}
	if naturalKeyJSON != "" {
		skip[naturalKeyJSON] = true
	}

	var out []string
	for _, line := range strings.Split(body, "\n") {
		tag := importStructTag(line)
		if tag == "" || !strings.Contains(importTagValue(tag, "binding"), "required") {
			continue
		}
		if strings.Contains(importTagValue(tag, "gorm"), "default:") {
			continue
		}
		name := importTagValue(tag, "json")
		if i := strings.Index(name, ","); i >= 0 {
			name = name[:i]
		}
		if skip[name] {
			continue
		}
		out = append(out, name)
	}
	return out
}

// importStructTag returns the backquoted tag of a struct field line, or "".
func importStructTag(line string) string {
	open := strings.Index(line, "`")
	if open < 0 {
		return ""
	}
	close := strings.LastIndex(line, "`")
	if close <= open {
		return ""
	}
	return line[open+1 : close]
}

// importTagValue reads one key out of a struct tag.
func importTagValue(tag, key string) string {
	at := strings.Index(tag, key+`:"`)
	if at < 0 {
		return ""
	}
	rest := tag[at+len(key)+2:]
	if end := strings.Index(rest, `"`); end >= 0 {
		return rest[:end]
	}
	return ""
}

// ImportNameResolverFor writes the resolver for a belongs_to column, choosing
// between creating a missing record and refusing to.
//
// The choice is here rather than at each call site so the generator and the
// upgrade cannot answer it differently, which is how two copies of generated
// code start to disagree.
func ImportNameResolverFor(base, relModel, keyJSON, keyGo string, unfillable []string) string {
	if len(unfillable) == 0 {
		return ImportNameResolver(base, relModel, keyJSON, keyGo)
	}
	return ImportNameResolverNoCreate(base, relModel, keyJSON, unfillable)
}

// ImportNameResolverNoCreate is the resolver for a related model the importer
// cannot create: it looks the record up and fails the row when it is missing,
// naming the columns it could not have supplied.
//
// It would rather refuse than guess. The only value an importer could invent
// for an owner column is "whoever is running the import", which is right for
// an owner and wrong for anything else, and a wrong foreign key is worse than
// a refused row because nobody finds out. See grit#129.
func ImportNameResolverNoCreate(base, relModel, keyJSON string, unfillable []string) string {
	return strings.NewReplacer(
		"{{name}}", importLowerCamel(base),
		"{{fn}}", "resolve"+importUpperCamel(base),
		"{{base}}", base,
		"{{Rel}}", relModel,
		"{{key}}", keyJSON,
		"{{cols}}", strings.Join(unfillable, ", "),
	).Replace(`	// {{name}}IDs remembers the {{Rel}} each {{key}} in the CSV resolved to, so a
	// distinct {{key}} costs one lookup rather than one per row.
	//
	// A missing {{Rel}} is NOT created here, because {{Rel}} requires {{cols}}
	// and a {{base}} column holds only the {{key}}. Creating it from the {{key}}
	// alone fails on the constraint, once per row, which is what grit#129 was.
	// The row is refused instead, with something a person can act on.
	{{name}}IDs := map[string]string{}
	{{fn}} := func(v string) (string, error) {
		if id, ok := {{name}}IDs[v]; ok {
			return id, nil
		}
		var rel models.{{Rel}}
		if err := db.Where("{{key}} = ?", v).First(&rel).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", fmt.Errorf("{{base}} %q does not exist. An import cannot create one, because {{Rel}} requires {{cols}}, which this file has no column for. Add the {{Rel}} first, or import {{base}} records before this one", v)
			}
			return "", fmt.Errorf("{{base}} %q: %w", v, err)
		}
		{{name}}IDs[v] = rel.ID
		return rel.ID, nil
	}

`)
}

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

// ImportUserLookup resolves a CSV cell to an existing user and fails the row
// when there is none. It never creates one: a users row is an account, and an
// importer that made one per unrecognised email let any signed-in caller mint
// accounts around registration, with no password and nothing verified.
// It writes the row code that resolves a User column.
//
// scoped narrows the lookup to the active organization's members, which is
// what the multitenant plugin needs and what a project without it must not
// have. Unnarrowed, this resolved an email against every user on the
// platform: an import running in one organization could file rows under
// somebody who belongs only to another, and the row error said whether the
// address exists at all. See grit#139, and GHSA-4j45-98gm-mv9v for the shape.
func ImportUserLookup(col, where, fkGo, cond string, ptr bool, scoped bool) string {
	// An optional relationship is a nullable column, so the importer has
	// to assign a pointer to it.
	assign := "rel.ID"
	if ptr {
		assign = "&rel.ID"
	}
	// Read at runtime, not decided here: an administrator importing
	// platform-wide has no active organization and must still be able to.
	// Membership is checked unscoped on purpose, since the members table is
	// how the scope itself is resolved.
	query, scope := "db", ""
	if scoped {
		query = "q"
		scope = "\t\t\tq, inOrg := db, \"\"\n" +
			"\t\t\tif orgID, ok := tenant.FromContext(ctx); ok {\n" +
			"\t\t\t\tq = db.Where(\"id IN (?)\", tenant.Unscoped(db).\n" +
			"\t\t\t\t\tModel(&models.OrganizationMember{}).\n" +
			"\t\t\t\t\tSelect(\"user_id\").\n" +
			"\t\t\t\t\tWhere(\"organization_id = ?\", orgID))\n" +
			"\t\t\t\tinOrg = \" in this organization\"\n" +
			"\t\t\t}\n"
	}
	message := "\"no user with %s %%q\", v"
	if scoped {
		message = "\"no user with %s %%q%%s\", v, inOrg"
	}
	return fmt.Sprintf("\t\tif v, ok := get(rec, %q); ok && v != \"\"%s {\n"+
		"\t\t\tvar rel models.User\n"+
		scope+
		"\t\t\tif err := "+query+".Where(%q, v).First(&rel).Error; err != nil {\n"+
		"\t\t\t\tmessage := fmt.Sprintf("+message+")\n"+
		"\t\t\t\tif !errors.Is(err, gorm.ErrRecordNotFound) {\n"+
		"\t\t\t\t\tmessage = err.Error()\n"+
		"\t\t\t\t}\n"+
		"\t\t\t\tfailed++\n"+
		"\t\t\t\tif len(rowErrors) < 50 {\n"+
		"\t\t\t\t\trowErrors = append(rowErrors, map[string]interface{}{\"row\": rowNum, \"message\": message})\n"+
		"\t\t\t\t}\n"+
		"\t\t\t\tcontinue\n"+
		"\t\t\t}\n"+
		"\t\t\titem.%s = "+assign+"\n"+
		"\t\t}\n", col, cond, where+" = ?", where, fkGo)
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
