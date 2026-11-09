package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// grit#197: regenerating a resource wrote a second copy of its index pair.
//
// The existence check asked whether <this run's timestamp>_<slug>.up.sql was
// there, and the timestamp is new every run, so it only caught a collision
// inside the same second.
//
// The up side is harmless: CREATE INDEX CONCURRENTLY IF NOT EXISTS. The down
// side is not. Both pairs record as applied and either one's down drops the
// index they share, so a rollback leaves the history saying the older pair is
// applied while the index is gone. On a table big enough to want an hnsw
// index, that is a retrieval path quietly becoming a sequential scan.
func TestVectorIndexMigrationIsWrittenOnce(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)

	// The pair is Postgres-only, so the project has to say so.
	writePostgresEnv(t, root)

	def, err := ParseInlineFields("SourceChunk", "text:text,embedding:vector[1536]")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}

	first := newTestGenerator(root, module, def)
	if err := first.Run(); err != nil {
		t.Fatalf("first Run(): %v", err)
	}
	after := vectorPairs(t, root)
	if len(after) != 2 {
		t.Fatalf("the first generate wrote %d migration files, want 2 (one up, one down): %v",
			len(after), after)
	}

	// Again, which is what --force does.
	second := newTestGenerator(root, module, def)
	if err := second.Run(); err != nil {
		t.Fatalf("second Run(): %v", err)
	}
	again := vectorPairs(t, root)
	if len(again) != 2 {
		t.Errorf("regenerating wrote %d migration files, want the same 2.\n\n"+
			"    Two pairs for one index both record as applied, and either one's down\n"+
			"    drops the index they share, so the recorded state and the real state\n"+
			"    disagree with nothing to say so.\n\n    %v", len(again), again)
	}
	// And it is the same pair, not a replacement under a new name.
	for i := range after {
		if after[i] != again[i] {
			t.Errorf("the pair was rewritten under a different name: %s became %s",
				after[i], again[i])
		}
	}
}

// And the reporting: --force promises to rewrite this resource's files, and
// the migrations directory is not one of them, so leaving it alone is said out
// loud rather than done silently.
func TestVectorIndexMigrationReportsWhatItDid(t *testing.T) {
	const module = "myapp/apps/api"
	root := setupMinimalProject(t, module)
	writePostgresEnv(t, root)

	def, err := ParseInlineFields("SourceChunk", "text:text,embedding:vector[1536]")
	if err != nil {
		t.Fatal(err)
	}
	names := MakeNames("SourceChunk")
	g := newTestGenerator(root, module, def)

	name, wrote, err := g.writeVectorIndexMigration(names, def.Fields[1])
	if err != nil {
		t.Fatalf("first call: %v", err)
	}
	if name == "" || !wrote {
		t.Fatalf("the first call reported name=%q wrote=%v, want a name and true", name, wrote)
	}

	again, wroteAgain, err := g.writeVectorIndexMigration(names, def.Fields[1])
	if err != nil {
		t.Fatalf("second call: %v", err)
	}
	if wroteAgain {
		t.Error("the second call wrote a second pair")
	}
	if again != name {
		t.Errorf("the second call named %q, want the existing %q: the caller prints this, "+
			"so it has to be the pair that is really there", again, name)
	}
}

// vectorPairs is the migration files in the project, sorted by ReadDir.
func vectorPairs(t *testing.T, root string) []string {
	t.Helper()
	dir := filepath.Join(root, "apps", "api", "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			out = append(out, e.Name())
		}
	}
	return out
}

// writePostgresEnv makes projectDBProvider answer postgres, which is what the
// index pair is gated on.
func writePostgresEnv(t *testing.T, root string) {
	t.Helper()
	path := filepath.Join(root, ".env")
	if err := os.WriteFile(path, []byte("DB_PROVIDER=postgres\n"), 0o644); err != nil {
		t.Fatalf("writing .env: %v", err)
	}
}
