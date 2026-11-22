package scaffold

import "strings"

// grit migrate plan: what AutoMigrate is about to do, before it does it.
//
// A migration is the highest-blast-radius thing a deploy runs, and the answer
// up to now was "run it and read what it says afterwards" (grit#58).
//
// The prediction is not a guess from the model source. The planned shape comes
// from GORM's own parse of each model, which is the parse AutoMigrate uses, and
// the current shape from Snapshot, which reads the database. The plan is Diff
// between them: the same function that records what a run actually changed, so
// a plan and the run it describes are compared by one rule and cannot disagree
// about what counts as a change.

// apiMigratePlanGo returns internal/migrate/plan.go.
func apiMigratePlanGo() string {
	return `package migrate

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// What AutoMigrate is about to do, worked out without doing it.
//
// A migration is the highest-blast-radius thing a deploy does, and the
// framework's answer up to now was "run it and read what it says afterwards".
// This answers the same question before the lock is taken (grit#58).
//
// ── How the prediction is made ────────────────────────────────────────────
//
// Not by guessing from the model source. The planned shape comes from GORM's
// own parse of the model, which is the same parse AutoMigrate uses to decide
// what to create, so the two cannot disagree about a column name or an index.
// The current shape comes from Snapshot, which reads the database. The plan is
// Diff between them: exactly the function that records what a run actually
// changed, so a plan and the run it describes are compared by one rule.
//
// ── What it deliberately does not say ─────────────────────────────────────
//
// Seconds. Lock duration depends on the table, the hardware, the concurrent
// load and the plan the server picks, and a number here would be believed. It
// reports the row count and says what holds the lock, which is the input to a
// judgement rather than a substitute for one.
//
// It also does not warn about dropping anything, because AutoMigrate never
// drops a column, a table or an index. There is no data-loss case to block in
// an AutoMigrate plan, and saying so is more useful than a flag guarding
// something that cannot happen. The destructive command is migrate down,
// which has its own confirmation.

// How many rows make an index build worth warning about.
//
// Not a cliff. An index on ten thousand rows is milliseconds on any machine
// built this decade; on ten million it is long enough that somebody should
// have decided to take it deliberately. The number is low enough to catch the
// second case early and high enough that a dev database stays quiet.
const indexWarningRows = 50_000

// Risk is one thing about a planned change that is worth knowing first.
type Risk struct {
	Change Change
	// What could go wrong, in one sentence.
	Says string
	// What to do instead, as something to type or to decide.
	Instead string
	// Rows in the table at the time of the plan, where that is the point.
	// Negative when it was not counted.
	Rows int64
}

// Plan is what a migration would change, and what to know before it does.
type Plan struct {
	Changes []Change
	Risks   []Risk
	// fields is the parsed model field behind each "table.column", so the
	// rendered statement can carry the real column definition rather than
	// just the name.
	fields map[string]*schema.Field
}

// Empty reports a plan with nothing in it, which is the usual answer.
func (p Plan) Empty() bool { return len(p.Changes) == 0 }

// BuildPlan works out what AutoMigrate would change, without changing anything.
//
// models is the same slice the migration passes to AutoMigrate.
func BuildPlan(db *gorm.DB, models []any) (Plan, error) {
	db = quietly(db)

	current, err := Snapshot(db)
	if err != nil {
		return Plan{}, err
	}
	planned, fields, err := plannedSchema(db, models)
	if err != nil {
		return Plan{}, err
	}

	// Diff in the direction a migration runs: what the models have that the
	// database does not. The reverse direction is a column somebody dropped
	// from a model, which AutoMigrate leaves alone, so it is not a change.
	changes := Diff(current, planned)

	plan := Plan{Changes: changes, fields: fields}
	for _, change := range changes {
		if risk, ok := riskOf(db, current, fields, change); ok {
			plan.Risks = append(plan.Risks, risk)
		}
	}
	return plan, nil
}

// quietly returns a session whose log output goes nowhere.
//
// The plan is a report. It should print what it found and nothing else, and
// reading a schema is chatty: the sqlite driver logs a line per table while
// listing indexes, which on a plain project is 43 lines of PRAGMA before the
// first word of the answer.
//
// A Silent log level does not do it, and the reason is worth writing down.
// glebarez/sqlite GetIndexes calls m.DB.Debug().Raw(...), and Debug() forces
// the level back to Info on whatever logger it is given, so a level set here
// loses. The writer does not: Debug() promotes a logger that writes to
// io.Discard and the output still goes nowhere.
func quietly(db *gorm.DB) *gorm.DB {
	silent := logger.New(log.New(io.Discard, "", 0), logger.Config{
		LogLevel: logger.Silent,
	})
	return db.Session(&gorm.Session{Logger: silent})
}

// plannedSchema is what the database would look like with every model applied.
//
// Built from GORM's parse of each model rather than from the Go source, so the
// column names, the index names and the dialect's own naming rules are the
// ones AutoMigrate will use.
func plannedSchema(db *gorm.DB, models []any) (Schema, map[string]*schema.Field, error) {
	out := Schema{}
	fields := map[string]*schema.Field{}

	for _, model := range models {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return nil, nil, fmt.Errorf("parsing %T: %w", model, err)
		}
		s := stmt.Schema
		if historyTables[s.Table] || internalTable(s.Table) {
			continue
		}

		shape, seen := out[s.Table]
		if !seen {
			shape = Table{Columns: map[string]bool{}, Indexes: map[string]bool{}}
		}
		for _, field := range s.Fields {
			// A field with no column is an association, not a column.
			if field.DBName == "" || field.IgnoreMigration {
				continue
			}
			shape.Columns[field.DBName] = true
			fields[s.Table+"."+field.DBName] = field
		}
		for _, index := range s.ParseIndexes() {
			shape.Indexes[index.Name] = true
		}
		out[s.Table] = shape
	}
	return out, fields, nil
}

// riskOf decides whether one planned change is worth a warning.
func riskOf(db *gorm.DB, current Schema, fields map[string]*schema.Field, change Change) (Risk, bool) {
	switch change.Kind {
	case KindTable:
		// A new table is empty. Nothing to lock and nothing to backfill.
		return Risk{}, false

	case KindColumn:
		field := fields[change.OnTable+"."+change.Name]
		if field == nil || !field.NotNull || field.HasDefaultValue || field.PrimaryKey {
			return Risk{}, false
		}
		rows := countRows(db, change.OnTable)
		if rows == 0 {
			// Adding a NOT NULL column to an empty table is free, and warning
			// about it on a dev database is how a warning gets ignored.
			return Risk{}, false
		}
		return Risk{
			Change: change,
			Rows:   rows,
			Says: fmt.Sprintf(
				"NOT NULL with no default, and %s already has %s. Postgres refuses this outright; "+
					"the ones that accept it write a value you did not choose into every row.",
				change.OnTable, plural(rows, "row")),
			Instead: "Add it nullable, backfill it, then add the constraint. " +
				"grit generate migration gives you the pair to do that in.",
		}, true

	default:
		if !hasRows(current, change.OnTable) {
			return Risk{}, false
		}
		rows := countRows(db, change.OnTable)
		if rows < indexWarningRows {
			return Risk{}, false
		}
		return Risk{
			Change: change,
			Rows:   rows,
			Says: fmt.Sprintf(
				"building an index on %s, which has %s. AutoMigrate does not use CONCURRENTLY, "+
					"so writes to the table wait until it finishes.",
				change.OnTable, plural(rows, "row")),
			Instead: "Create it by hand with CREATE INDEX CONCURRENTLY and let AutoMigrate " +
				"find it already there, or take the lock at a time you have chosen.",
		}, true
	}
}

// hasRows reports whether the table exists at all in the current schema. A
// table about to be created cannot have rows, and counting one that does not
// exist is an error nobody needs to see.
func hasRows(current Schema, table string) bool {
	_, exists := current[table]
	return exists
}

// countRows is an exact count.
//
// An estimate would be cheaper, and every dialect spells it differently; the
// count is only taken for a table a warning already applies to, which is at
// most a handful per plan. On a table big enough for the count to be slow, the
// warning is the one you most want to be right.
func countRows(db *gorm.DB, table string) int64 {
	var n int64
	if err := db.Table(table).Count(&n).Error; err != nil {
		return 0
	}
	return n
}

// wrap breaks a sentence at 72 columns under a fixed indent.
//
// A warning is only read if it can be read. These run to two hundred
// characters, and a terminal folds them at whatever width it has, which puts
// the second half of a sentence under the left margin.
func wrap(text, indent string) string {
	const width = 72
	var out strings.Builder
	line := indent
	for _, word := range strings.Fields(text) {
		if len(line)+len(word)+1 > len(indent)+width && line != indent {
			out.WriteString(line + "\n")
			line = indent
		}
		if line != indent {
			line += " "
		}
		line += word
	}
	if line != indent {
		out.WriteString(line + "\n")
	}
	return out.String()
}

func plural(n int64, word string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}

// planStatement is the statement a change will run, as closely as it can be
// said before AutoMigrate builds it.
//
// Not Statements(), which is next door and renders the INVERSE of a change
// because it exists for the rollback. Using it here printed
// ALTER TABLE "orders" DROP COLUMN "owner_name" for a column being added,
// which is the one thing a plan must never do.
//
// A column gets its real definition through the migrator, so the type, the
// NOT NULL and the default are visible: those are what the warnings are
// about. A table is named rather than spelled out, because the CREATE for a
// generated model is forty lines nobody reads in a plan.
func planStatement(db *gorm.DB, plan Plan, change Change) string {
	table := quote(db, change.OnTable)
	switch change.Kind {
	case KindTable:
		return "CREATE TABLE " + table

	case KindColumn:
		column := quote(db, change.Name)
		field := plan.fields[change.OnTable+"."+change.Name]
		if field == nil {
			return "ALTER TABLE " + table + " ADD COLUMN " + column
		}
		definition := db.Migrator().FullDataTypeOf(field)
		rendered := strings.TrimSpace(definition.SQL)
		for _, value := range definition.Vars {
			rendered = strings.Replace(rendered, "?", fmt.Sprintf("%v", value), 1)
		}
		if rendered == "" {
			return "ALTER TABLE " + table + " ADD COLUMN " + column
		}
		return "ALTER TABLE " + table + " ADD COLUMN " + column + " " + rendered

	default:
		return "CREATE INDEX " + quote(db, change.Name) + " ON " + table
	}
}

// Format renders a plan the way the command prints it.
//
// Here rather than in cmd/migrate so a test can read what a person would.
func Format(db *gorm.DB, plan Plan) string {
	var b strings.Builder

	if plan.Empty() {
		b.WriteString("\n  Nothing to migrate: the database already matches the models.\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("\n  %s:\n\n", plural(int64(len(plan.Changes)), "change")))
	for _, change := range plan.Changes {
		b.WriteString("    " + planStatement(db, plan, change) + "\n")
	}

	if len(plan.Risks) == 0 {
		b.WriteString("\n  Nothing here takes a lock worth planning around.\n")
		return b.String()
	}

	risks := make([]Risk, len(plan.Risks))
	copy(risks, plan.Risks)
	sort.Slice(risks, func(i, j int) bool { return risks[i].Rows > risks[j].Rows })

	b.WriteString(fmt.Sprintf("\n  %s to read first:\n\n",
		plural(int64(len(risks)), "warning")))
	for _, risk := range risks {
		b.WriteString("    " + risk.Change.String() + "\n")
		b.WriteString(wrap(risk.Says, "      "))
		b.WriteString(wrap("Instead: "+risk.Instead, "      ") + "\n")
	}
	b.WriteString("  AutoMigrate never drops a column, a table or an index, so nothing\n")
	b.WriteString("  in this plan loses data. The command that can is migrate down.\n")
	return b.String()
}
`
}

// apiMigratePlanTestGo returns internal/migrate/plan_test.go.
//
// The test that matters is the last one: it builds a plan, runs the
// migration, and requires the plan to equal what the run actually changed.
// Two implementations of "what counts as a change" would drift, and the
// drift would be invisible, because a plan nobody checks is reassuring
// either way.
//
// @@ stands in for a backtick, because the body is Go with struct tags in it
// and a raw string cannot hold one.
func apiMigratePlanTestGo() string {
	return strings.ReplaceAll(`package migrate

import (
	"strings"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// The plan, against a real database.
//
// SQLite rather than a mock, because the question the plan answers is what the
// database already has, and a mock of that would be a mock of the only input.

type planWidget struct {
	ID        uint @@gorm:"primarykey"@@
	Name      string
	CreatedAt time.Time
}

func (planWidget) TableName() string { return "plan_widgets" }

// The same table with a column added, which is the change a plan describes.
type planWidgetGrown struct {
	ID        uint @@gorm:"primarykey"@@
	Name      string
	Colour    string
	CreatedAt time.Time
}

func (planWidgetGrown) TableName() string { return "plan_widgets" }

// And with a NOT NULL column and no default, which is the one worth warning
// about on a table that already has rows.
type planWidgetRequired struct {
	ID        uint @@gorm:"primarykey"@@
	Name      string
	Owner     string @@gorm:"not null"@@
	CreatedAt time.Time
}

func (planWidgetRequired) TableName() string { return "plan_widgets" }

func planDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("opening the database: %v", err)
	}
	// One connection: every :memory: connection is its own empty database, so
	// a pool makes the schema appear and disappear between statements.
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("taking the pool: %v", err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })
	return db
}

// A database that already matches its models has nothing to plan.
func TestPlanIsEmptyWhenTheSchemaMatches(t *testing.T) {
	db := planDB(t)
	if err := db.AutoMigrate(&planWidget{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	plan, err := BuildPlan(db, []any{&planWidget{}})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if !plan.Empty() {
		t.Errorf("a matching schema planned %d changes: %v", len(plan.Changes), plan.Changes)
	}
	if out := Format(db, plan); !strings.Contains(out, "Nothing to migrate") {
		t.Errorf("the empty plan does not say so:\n%s", out)
	}
}

// A table the database does not have is a create.
func TestPlanSeesANewTable(t *testing.T) {
	db := planDB(t)

	plan, err := BuildPlan(db, []any{&planWidget{}})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if len(plan.Changes) != 1 || plan.Changes[0].Kind != KindTable {
		t.Fatalf("got %v, want one table", plan.Changes)
	}
	if plan.Changes[0].OnTable != "plan_widgets" {
		t.Errorf("the table is %q", plan.Changes[0].OnTable)
	}
	// A new table is empty, so there is nothing to lock or backfill.
	if len(plan.Risks) != 0 {
		t.Errorf("a new table produced warnings: %v", plan.Risks)
	}
}

// A column the model has and the database does not is an add.
func TestPlanSeesANewColumn(t *testing.T) {
	db := planDB(t)
	if err := db.AutoMigrate(&planWidget{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	plan, err := BuildPlan(db, []any{&planWidgetGrown{}})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if len(plan.Changes) != 1 || plan.Changes[0].Kind != KindColumn {
		t.Fatalf("got %v, want one column", plan.Changes)
	}
	if plan.Changes[0].Name != "colour" {
		t.Errorf("the column is %q, want colour", plan.Changes[0].Name)
	}
}

// The warning this exists for: NOT NULL, no default, on a table with rows.
func TestPlanWarnsAboutARequiredColumnOnAFullTable(t *testing.T) {
	db := planDB(t)
	if err := db.AutoMigrate(&planWidget{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}
	if err := db.Create(&planWidget{Name: "one"}).Error; err != nil {
		t.Fatalf("inserting a row: %v", err)
	}

	plan, err := BuildPlan(db, []any{&planWidgetRequired{}})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if len(plan.Risks) != 1 {
		t.Fatalf("got %d warnings, want one about owner: %v", len(plan.Risks), plan.Risks)
	}
	risk := plan.Risks[0]
	if risk.Change.Name != "owner" {
		t.Errorf("the warning is about %q", risk.Change.Name)
	}
	if risk.Rows != 1 {
		t.Errorf("the warning reports %d rows, want 1", risk.Rows)
	}
	if !strings.Contains(risk.Instead, "nullable") {
		t.Errorf("the warning does not say what to do instead: %q", risk.Instead)
	}
	if out := Format(db, plan); !strings.Contains(out, "owner") ||
		!strings.Contains(out, "Instead:") {
		t.Errorf("the rendered plan does not carry the warning:\n%s", out)
	}
}

// The same column on an empty table is free, and warning about it is how a
// warning gets ignored.
func TestPlanIsQuietAboutARequiredColumnOnAnEmptyTable(t *testing.T) {
	db := planDB(t)
	if err := db.AutoMigrate(&planWidget{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	plan, err := BuildPlan(db, []any{&planWidgetRequired{}})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if len(plan.Changes) != 1 {
		t.Fatalf("got %v, want the column", plan.Changes)
	}
	if len(plan.Risks) != 0 {
		t.Errorf("an empty table produced warnings: %v", plan.Risks)
	}
}

// The plan has to describe the migration that then runs.
//
// Two implementations of "what counts as a change" would drift, and the drift
// would be invisible: the plan would be reassuring and the run would do
// something else. So this runs the plan, then the migration, then compares the
// plan against what the run actually changed.
func TestThePlanMatchesWhatTheMigrationDoes(t *testing.T) {
	db := planDB(t)
	if err := db.AutoMigrate(&planWidget{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	plan, err := BuildPlan(db, []any{&planWidgetGrown{}})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}

	before, err := Snapshot(db)
	if err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if err := db.AutoMigrate(&planWidgetGrown{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}
	after, err := Snapshot(db)
	if err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	actual := Diff(before, after)

	if len(actual) != len(plan.Changes) {
		t.Fatalf("the plan said %d changes and the run made %d:\nplan %v\nrun  %v",
			len(plan.Changes), len(actual), plan.Changes, actual)
	}
	for i := range actual {
		if actual[i] != plan.Changes[i] {
			t.Errorf("change %d: the plan said %v and the run made %v",
				i, plan.Changes[i], actual[i])
		}
	}

	// And having run it, the plan is empty.
	again, err := BuildPlan(db, []any{&planWidgetGrown{}})
	if err != nil {
		t.Fatalf("BuildPlan: %v", err)
	}
	if !again.Empty() {
		t.Errorf("after migrating, the plan still wants %v", again.Changes)
	}
}
`, "@@", string(rune(96)))
}
