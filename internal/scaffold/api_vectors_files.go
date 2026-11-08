package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeVectorFiles writes internal/vectors: the column an embedding goes in.
//
// grit#191. The AI battery produces embeddings and there was nowhere
// first-class to put them, so the retrieval half of an AI feature was four
// hand-built pieces in every project that tried: the extension, the column, the
// index, and a Go type implementing sql.Scanner and driver.Valuer. Each of them
// slightly different each time.
//
// Shipped to every project whether or not it has a vector field. It is one
// small file, it costs nothing unused, and it means `grit generate resource
// X --fields "embedding:vector[1536]"` needs no second step.
func writeVectorFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "vectors", "vectors.go"):      vectorsGo(),
		filepath.Join(apiRoot, "internal", "vectors", "vectors_test.go"): vectorsTestGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func vectorsGo() string {
	return `// Package vectors is an embedding, and the column it lives in.
//
// An embedding is a fixed-length list of float32. Postgres with the pgvector
// extension has a real type for it, ` + "`vector(n)`" + `, with distance operators and
// an index that makes nearest-neighbour search fast. Nothing else Grit supports
// does, so this type stores itself as text everywhere else and the queries that
// need the operators say so.
//
//	type Document struct {
//		Embedding vectors.Vector ` + "`" + `vector:"1536" json:"embedding"` + "`" + `
//	}
//
//	// The 10 nearest, by cosine distance.
//	db.Model(&Document{}).
//		Clauses(vectors.Nearest("embedding", query)).
//		Limit(10).
//		Find(&out)
//
// Clauses, not Order. gorm's Order() switches on the type of its argument and
// handles only a string and a clause.OrderByColumn: an expression falls
// through and nothing is added, so the query runs with no ORDER BY and returns
// rows in scan order. It looks like it worked.
//
// # The dimension is a schema decision
//
// It is part of the column type, so changing it is a migration and a full
// re-embed of every row. That follows the embedding model: 1536 is OpenAI's
// text-embedding-3-small, 3072 its large, 768 many open models. Picking the
// model is picking the number, and picking it later is expensive. It is cheap
// now.
package vectors

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Vector is an embedding.
//
// float32 rather than float64 because that is what pgvector stores and what
// every embedding API returns. float64 would double the column for no
// precision anything uses.
type Vector []float32

// New copies the values into a Vector.
func New(values ...float32) Vector {
	out := make(Vector, len(values))
	copy(out, values)
	return out
}

// FromFloat64 converts what a JSON decoder hands you. An AI gateway returns
// JSON numbers, which decode to float64.
func FromFloat64(values []float64) Vector {
	out := make(Vector, len(values))
	for i, v := range values {
		out[i] = float32(v)
	}
	return out
}

// Dim is the number of dimensions.
func (v Vector) Dim() int { return len(v) }

// String is pgvector's text format, which is also what Value sends.
func (v Vector) String() string {
	if v == nil {
		return ""
	}
	var b strings.Builder
	b.WriteByte('[')
	for i, f := range v {
		if i > 0 {
			b.WriteByte(',')
		}
		// -1 gives the shortest representation that round-trips, so a value
		// read back is the value written.
		b.WriteString(strconv.FormatFloat(float64(f), 'f', -1, 32))
	}
	b.WriteByte(']')
	return b.String()
}

// Value stores the vector.
//
// The text format on purpose, not the binary one: it is what pgvector accepts
// on input, it is readable in a query log, and it is the same bytes the text
// column on SQLite and MySQL holds. One representation, every driver.
func (v Vector) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}
	return v.String(), nil
}

// Scan reads it back, from any driver.
func (v *Vector) Scan(src any) error {
	switch s := src.(type) {
	case nil:
		*v = nil
		return nil
	case string:
		return v.parse(s)
	case []byte:
		return v.parse(string(s))
	default:
		return fmt.Errorf("vectors: cannot scan %T into a Vector", src)
	}
}

// parse reads the text format, and tolerates the JSON array it is otherwise
// identical to.
func (v *Vector) parse(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		*v = nil
		return nil
	}
	if !strings.HasPrefix(s, "[") || !strings.HasSuffix(s, "]") {
		return fmt.Errorf("vectors: %q is not a vector: expected [a,b,c]", s)
	}
	inner := strings.TrimSpace(s[1 : len(s)-1])
	if inner == "" {
		*v = Vector{}
		return nil
	}
	parts := strings.Split(inner, ",")
	out := make(Vector, 0, len(parts))
	for _, p := range parts {
		f, err := strconv.ParseFloat(strings.TrimSpace(p), 32)
		if err != nil {
			return fmt.Errorf("vectors: %q is not a number in %q", strings.TrimSpace(p), s)
		}
		out = append(out, float32(f))
	}
	*v = out
	return nil
}

// GormDataType is the logical type, which is what GORM's schema parser asks
// for first.
//
// Without it, parsing the model fails outright:
//
//	failed to parse field: Embedding, error: unsupported data type:
//
// GORM cannot infer a data type for a named slice type, and GormDBDataType is
// consulted later, by the migrator, for the DDL. "string" because that is what
// Value sends and what Scan reads: the physical column is decided below.
func (Vector) GormDataType() string { return "string" }

// GormDBDataType is the column type, per database.
//
// The dimension comes from the field's own ` + "`vector:\"n\"`" + ` tag rather than from
// the Go type, because Go has no way to put it there: a length is not part of a
// slice type. The generator writes the tag; nothing else reads it.
//
// Only Postgres gets the real type. Everywhere else it is text, which the
// format above is designed to survive: a project developing on SQLite stores
// and reads the same values, and loses only the operators and the index, which
// are what you move to Postgres for.
func (Vector) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	if db == nil || db.Dialector == nil {
		return "text"
	}
	if db.Dialector.Name() != "postgres" {
		return "text"
	}
	if dim := DimensionOf(field); dim > 0 {
		return fmt.Sprintf("vector(%d)", dim)
	}
	// No dimension declared: pgvector allows it, and refuses an index until
	// every row agrees. Worth having rather than refusing to migrate.
	return "vector"
}

// DimensionOf reads the declared dimension, or 0.
func DimensionOf(field *schema.Field) int {
	if field == nil {
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(field.Tag.Get("vector")))
	if err != nil || n <= 0 {
		return 0
	}
	return n
}

// MaxDimension is pgvector's ceiling for a ` + "`vector`" + ` column.
const MaxDimension = 16000

// MaxIndexDimension is the ceiling for an hnsw or ivfflat index, which is
// lower than the column's. A wider column still stores and still searches, by
// reading every row.
const MaxIndexDimension = 2000

// columnName is a bare identifier. The distance helpers interpolate the column
// into SQL, because an operand of an operator cannot be a bind parameter, so
// the name is checked rather than trusted.
var columnName = regexp.MustCompile(` + "`" + `^[a-z_][a-z0-9_]*$` + "`" + `)

// Nearest orders by cosine distance: 0 is identical, 2 is opposite. The usual
// choice for text embeddings, because it ignores magnitude, and the one to use
// unless you know you want another.
//
//	db.Model(&Document{}).Clauses(vectors.Nearest("embedding", query)).Limit(10)
//
// Pass it to Clauses. Order() would ignore it: see the package comment.
func Nearest(column string, v Vector) clause.OrderBy {
	return clause.OrderBy{Expression: CosineDistance(column, v)}
}

// NearestL2 orders by straight-line distance.
func NearestL2(column string, v Vector) clause.OrderBy {
	return clause.OrderBy{Expression: L2Distance(column, v)}
}

// NearestByInnerProduct orders by negative inner product.
func NearestByInnerProduct(column string, v Vector) clause.OrderBy {
	return clause.OrderBy{Expression: InnerProduct(column, v)}
}

// CosineDistance is the distance itself, for a Select that returns it or a
// Where that puts a ceiling on it. For ordering, use Nearest.
//
//	db.Model(&Document{}).
//		Select("*, ?  AS distance", vectors.CosineDistance("embedding", query)).
//		Clauses(vectors.Nearest("embedding", query))
func CosineDistance(column string, v Vector) clause.Expression {
	return distance(column, "<=>", v)
}

// L2Distance is straight-line distance, for embeddings whose magnitude carries
// meaning.
func L2Distance(column string, v Vector) clause.Expression {
	return distance(column, "<->", v)
}

// InnerProduct is the negative inner product, which pgvector spells <#>.
func InnerProduct(column string, v Vector) clause.Expression {
	return distance(column, "<#>", v)
}

// rejectedColumn is the expression a bad column name produces: a call to a
// function no database has.
//
// It has to be something that cannot run. The obvious choice, (SELECT 1/0),
// does not work: the planner drops an ORDER BY whose value is constant, so the
// query succeeded and returned rows in whatever order the scan produced, which
// is the exact failure this is here to prevent. A function name has to resolve
// at parse time, so it cannot be optimised away, and the error names itself.
const rejectedColumn = "grit_vectors_rejected_column_name()"

// distance builds the ordering expression, or one that fails loudly.
//
// A bad column name returns an expression that cannot run rather than silently
// ordering by nothing: a nearest-neighbour query that quietly returns arbitrary
// rows is worse than one that errors.
func distance(column, operator string, v Vector) clause.Expression {
	if !columnName.MatchString(column) {
		return clause.Expr{SQL: rejectedColumn}
	}
	return clause.Expr{
		SQL:  column + " " + operator + " ?",
		Vars: []any{v},
	}
}

// Extension is the Postgres extension a vector column needs.
const Extension = "vector"

// EnsureExtension creates the extension when some model actually has a vector
// column, and does nothing otherwise.
//
// Before AutoMigrate, not in a SQL migration pair. The pairs run after
// AutoMigrate, and AutoMigrate is what creates the column: a CREATE EXTENSION
// in a pair runs too late and the migration fails with
// ` + "`type \"vector\" does not exist`" + `.
//
// Conditional on a model using one, because the extension has to be available
// on the server. Running it unconditionally would fail every Postgres project
// whose server does not carry pgvector, which is most of them, including the
// stock ` + "`postgres:16-alpine`" + ` image.
func EnsureExtension(db *gorm.DB, models ...any) error {
	if db == nil || db.Dialector == nil || db.Dialector.Name() != "postgres" {
		return nil
	}
	if !AnyUses(models...) {
		return nil
	}
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS " + Extension).Error; err != nil {
		return fmt.Errorf("vectors: enabling the %s extension, which a model needs for its "+
			"embedding column: %w.\n"+
			"The stock postgres image does not carry pgvector. Use pgvector/pgvector:pg16 in "+
			"docker-compose.yml, or install it on the server", Extension, err)
	}
	return nil
}

// AnyUses reports whether any of these models has a Vector field.
//
// Reflection rather than a list the generator injects into: a list is a second
// place that has to agree with the models, and it drifts the first time
// somebody adds a column by hand.
func AnyUses(models ...any) bool {
	want := reflect.TypeOf(Vector{})
	for _, m := range models {
		if m == nil {
			continue
		}
		t := reflect.TypeOf(m)
		for t != nil && t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t == nil || t.Kind() != reflect.Struct {
			continue
		}
		if structHas(t, want, 0) {
			return true
		}
	}
	return false
}

// structHas walks a struct for the type, through embedded fields.
func structHas(t, want reflect.Type, depth int) bool {
	// Embedded structs nest; a model is not deep, and a cycle would not end.
	if depth > 4 {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type == want {
			return true
		}
		if f.Anonymous && f.Type.Kind() == reflect.Struct && structHas(f.Type, want, depth+1) {
			return true
		}
	}
	return false
}
`
}

func vectorsTestGo() string {
	return `package vectors

import (
	"reflect"
	"strings"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// The round trip is the whole contract: a value written and read back is the
// value written. Everything else in this package is in service of that.
func TestValueScanRoundTrip(t *testing.T) {
	cases := []Vector{
		nil,
		{},
		{0},
		{1, 2, 3},
		{-0.5, 0.25, 1e-7},
		// A float32 that is not exactly representable in decimal, which is
		// what an embedding is made of.
		{0.1, 0.2, 0.30000001},
	}

	for _, want := range cases {
		stored, err := want.Value()
		if err != nil {
			t.Fatalf("Value() on %v: %v", want, err)
		}
		var got Vector
		if err := got.Scan(stored); err != nil {
			t.Fatalf("Scan(%v): %v", stored, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("round trip changed the value: wrote %v, read %v", want, got)
		}
	}
}

// And from the bytes a driver may hand back instead of a string.
func TestScanFromBytes(t *testing.T) {
	var v Vector
	if err := v.Scan([]byte("[1,2.5,-3]")); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if want := (Vector{1, 2.5, -3}); !reflect.DeepEqual(v, want) {
		t.Errorf("got %v, want %v", v, want)
	}
}

func TestScanRefusesNonsense(t *testing.T) {
	for _, bad := range []any{"1,2,3", "[1,two]", "{1,2}", 42} {
		var v Vector
		if err := v.Scan(bad); err == nil {
			t.Errorf("Scan(%#v) was accepted, and a silently empty embedding "+
				"makes every search return arbitrary rows", bad)
		}
	}
}

func TestDimensionOfReadsTheTag(t *testing.T) {
	cases := []struct {
		tag  string
		want int
	}{
		{` + "`vector:\"1536\"`" + `, 1536},
		{` + "`vector:\"768\" json:\"e\"`" + `, 768},
		{` + "`json:\"e\"`" + `, 0},
		{` + "`vector:\"0\"`" + `, 0},
		{` + "`vector:\"nope\"`" + `, 0},
	}
	for _, c := range cases {
		field := &schema.Field{Tag: reflect.StructTag(c.tag)}
		if got := DimensionOf(field); got != c.want {
			t.Errorf("DimensionOf(%s) = %d, want %d", c.tag, got, c.want)
		}
	}
	if got := DimensionOf(nil); got != 0 {
		t.Errorf("DimensionOf(nil) = %d, want 0", got)
	}
}

// A column name is interpolated into SQL, because an operand of an operator
// cannot be a bind parameter. So it is checked.
func TestDistanceRefusesAnInjectedColumn(t *testing.T) {
	for _, bad := range []string{
		"embedding; DROP TABLE documents",
		"embedding)",
		"Embedding",
		"",
		"1embedding",
	} {
		expr := CosineDistance(bad, Vector{1})
		if sql := exprSQL(expr); !isRefusal(sql) {
			t.Errorf("CosineDistance(%q) built %q, which runs", bad, sql)
		}
	}
	if sql := exprSQL(CosineDistance("embedding", Vector{1})); sql != "embedding <=> ?" {
		t.Errorf("a good column built %q", sql)
	}
	if sql := exprSQL(L2Distance("embedding", Vector{1})); sql != "embedding <-> ?" {
		t.Errorf("L2Distance built %q", sql)
	}
	if sql := exprSQL(InnerProduct("embedding", Vector{1})); sql != "embedding <#> ?" {
		t.Errorf("InnerProduct built %q", sql)
	}
}

type vectorModel struct {
	ID        string
	Embedding Vector
}

type plainModel struct {
	ID   string
	Name string
}

type embedsVector struct {
	vectorModel
	Extra string
}

func TestAnyUses(t *testing.T) {
	if !AnyUses(&vectorModel{}) {
		t.Error("a model with a Vector field was not detected, so the extension is never created")
	}
	if !AnyUses(&plainModel{}, &vectorModel{}) {
		t.Error("a Vector field in the second model was missed")
	}
	if !AnyUses(&embedsVector{}) {
		t.Error("a Vector reached through an embedded struct was missed")
	}
	if AnyUses(&plainModel{}) {
		t.Error("a model with no Vector was reported as using one, which would run " +
			"CREATE EXTENSION on a server that may not have pgvector")
	}
	if AnyUses() || AnyUses(nil) {
		t.Error("no models reported as using one")
	}
}

// exprSQL is the SQL a distance expression would contribute.
func exprSQL(e clause.Expression) string {
	expr, ok := e.(clause.Expr)
	if !ok {
		return ""
	}
	return expr.SQL
}

// isRefusal reports the expression distance() builds for a column it will not
// interpolate: one that cannot run, rather than one that orders by nothing.
func isRefusal(sql string) bool {
	return strings.Contains(sql, "grit_vectors_rejected_column_name")
}

// The ordering has to reach the query.
//
// This is the test that matters, and the one whose absence let a real bug look
// like a working feature. gorm's Order() switches on its argument's type and
// silently ignores an expression, so
//
//	db.Model(&M{}).Order(vectors.CosineDistance("embedding", q))
//
// built a query with no ORDER BY at all and returned rows in scan order, which
// for a few freshly inserted rows is insertion order and reads exactly like a
// working nearest-neighbour search. A test on the helper's own SQL string
// cannot see that. This one asks gorm what it would send.
func TestNearestReachesTheQuery(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening an in-memory database: %v", err)
	}

	q := Vector{1, 0, 0}
	cases := []struct {
		name     string
		clause   clause.OrderBy
		operator string
	}{
		{"Nearest", Nearest("embedding", q), "<=>"},
		{"NearestL2", NearestL2("embedding", q), "<->"},
		{"NearestByInnerProduct", NearestByInnerProduct("embedding", q), "<#>"},
	}

	for _, c := range cases {
		var out []vectorRow
		sql := db.Session(&gorm.Session{DryRun: true}).
			Model(&vectorRow{}).
			Clauses(c.clause).
			Limit(10).
			Find(&out).Statement.SQL.String()

		if !strings.Contains(strings.ToUpper(sql), "ORDER BY") {
			t.Errorf("%s produced a query with no ORDER BY, so it returns rows in scan "+
				"order and looks like it worked:\n  %s", c.name, sql)
			continue
		}
		if !strings.Contains(sql, "embedding "+c.operator) {
			t.Errorf("%s did not order by %s:\n  %s", c.name, c.operator, sql)
		}
	}
}

// And the refused column reaches the query as something that cannot run,
// rather than as nothing at all. (SELECT 1/0) was the first attempt and the
// planner dropped it, because an ORDER BY of a constant is a no-op.
func TestARefusedColumnCannotRun(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening an in-memory database: %v", err)
	}
	var out []vectorRow
	sql := db.Session(&gorm.Session{DryRun: true}).
		Model(&vectorRow{}).
		Clauses(Nearest("embedding; DROP TABLE documents", Vector{1})).
		Find(&out).Statement.SQL.String()

	if strings.Contains(sql, "DROP TABLE") {
		t.Fatalf("the column name was interpolated into the query:\n  %s", sql)
	}
	if !isRefusal(sql) {
		t.Errorf("a refused column produced a query that runs:\n  %s", sql)
	}
}

// vectorRow is a model for the statement-building tests above.
type vectorRow struct {
	ID        uint
	Embedding Vector ` + "`" + `vector:"3"` + "`" + `
}
`
}
