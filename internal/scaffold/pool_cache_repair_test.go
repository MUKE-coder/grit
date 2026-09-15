package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPoolSizedForReplicasAndSentinel(t *testing.T) {
	db := apiDatabaseGo()
	if strings.Contains(db, dbPoolBlockOld) || !strings.Contains(db, dbPoolBlockNew) || !strings.Contains(db, "func SentinelPoolSize()") {
		t.Error("database.go still defaults the pool to 100 or has no SentinelPoolSize")
	}
	mustFormatGo(t, "database.go", db)

	routes := apiRoutesGo()
	if !strings.Contains(routes, sentinelPoolAnchor+sentinelPoolLine) {
		t.Error("routes.go does not cap Sentinel's storage pool")
	}

	env := envExampleFile(Options{ProjectName: "app", Architecture: ArchDouble})
	if strings.Contains(env, "DB_MAX_OPEN_CONNS=100") || !strings.Contains(env, "SENTINEL_DB_MAX_OPEN_CONNS=5") {
		t.Error(".env.example does not document the replica-sized pools")
	}
}

func TestPoolRepairs(t *testing.T) {
	fresh := apiDatabaseGo()
	old := strings.Replace(fresh, dbPoolBlockNew, dbPoolBlockOld, 1)
	old = strings.Replace(old, sentinelPoolFunc, "", 1)
	out, changes, warnings := repairDBPoolSource(old)
	if out != fresh || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("repairing the old database.go did not give the template (changes %v, warnings %v)", changes, warnings)
	}
	if again, changes, _ := repairDBPoolSource(out); again != out || len(changes) != 0 {
		t.Error("the database.go repair is not idempotent")
	}
	if _, _, warnings := repairDBPoolSource("package database\n\nfunc Connect() {}\n"); len(warnings) != 1 {
		t.Error("an edited database.go should get a note")
	}

	routes := apiRoutesGo()
	oldRoutes := strings.Replace(routes, sentinelPoolLine, "", 1)
	out, changes, warnings = repairSentinelPoolSource(oldRoutes)
	if out != routes || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("repairing the old routes.go did not give the template (changes %v, warnings %v)", changes, warnings)
	}
	if again, changes, _ := repairSentinelPoolSource(out); again != out || len(changes) != 0 {
		t.Error("the routes.go repair is not idempotent")
	}

	env := envExampleFile(Options{ProjectName: "app", Architecture: ArchDouble})
	out, _, _ = repairEnvPoolSource(strings.Replace(env, envPoolNew, envPoolOld, 1))
	if out != env {
		t.Error("repairing the old .env.example did not give the template")
	}
}

func TestProductionAPIUsesPgbouncer(t *testing.T) {
	for _, arch := range []Architecture{ArchDouble, ArchSingle, ArchTriple, ArchAPI} {
		fresh := dockerComposeProd(Options{ProjectName: "app", Architecture: arch})
		start, end, ok := serviceBlock(fresh, "api")
		if !ok {
			t.Fatalf("%s: no api service", arch)
		}
		block := fresh[start:end]
		if !strings.Contains(block, "POSTGRES_HOST: pgbouncer") || !strings.Contains(block, composeAPIDependsNew) {
			t.Errorf("%s: the production api does not connect through pgbouncer", arch)
		}
		if strings.Contains(fresh, "DB_HOST=pgbouncer") {
			t.Errorf("%s: the pgbouncer note still names DB_HOST, which the config never reads", arch)
		}
		if out, changes, _ := repairComposePgbouncerSource(fresh); out != fresh || len(changes) != 0 {
			t.Errorf("%s: a fresh production compose needs no repair: %v", arch, changes)
		}

		old := strings.Replace(fresh, composeAPICommentNew, composeAPICommentOld, 1)
		old = strings.Replace(old, composeAPIPostgresNew, composeAPIPostgresOld, 1)
		old = strings.Replace(old, composeAPIDependsNew, composeAPIDependsOld, 1)
		old = strings.Replace(old, pgbouncerUseNew, pgbouncerUseOld, 1)
		old = strings.Replace(old, pgbouncerListenNew, pgbouncerListenOld, 1)
		if !strings.Contains(fresh, pgbouncerListenNew) {
			t.Errorf("%s: pgbouncer does not listen on 6432, where its health check looks", arch)
		}
		out, changes, warnings := repairComposePgbouncerSource(old)
		if out != fresh || len(warnings) != 0 || len(changes) != 2 {
			t.Errorf("%s: repairing the previous compose file did not give the fresh one (changes %v, warnings %v)", arch, changes, warnings)
		}
	}
}

func TestCacheMiddlewareCollapsesMisses(t *testing.T) {
	src := cacheMiddlewareGo()
	for _, want := range []string{"cacheFlights.Do(key,", "resp.encode()", "\"http:v2:\""} {
		if !strings.Contains(src, want) {
			t.Errorf("the cache middleware is missing %q", want)
		}
	}
	if strings.Contains(src, "json:\"body\"") {
		t.Error("the cache middleware still stores the body as JSON")
	}
	mustFormatGo(t, "middleware/cache.go", src)
	mustFormatGo(t, "the old middleware/cache.go", cacheMiddlewareOld)
	if !strings.Contains(apiGoMod(Options{ProjectName: "app"}), "golang.org/x/sync "+xSyncVersion) {
		t.Error("go.mod does not require golang.org/x/sync")
	}

	old := strings.ReplaceAll(cacheMiddlewareOld, "{{MODULE}}", "example.com/app")
	out, changes, warnings := repairCacheMiddlewareSource(old)
	if out != strings.ReplaceAll(cacheMiddlewareNew, "{{MODULE}}", "example.com/app") || len(changes) != 1 || len(warnings) != 0 {
		t.Errorf("the old cache.go was not replaced (changes %v, warnings %v)", changes, warnings)
	}
	if again, changes, _ := repairCacheMiddlewareSource(out); again != out || len(changes) != 0 {
		t.Error("the cache.go repair is not idempotent")
	}
	if _, _, warnings := repairCacheMiddlewareSource(old + "\n// mine\n"); len(warnings) != 1 {
		t.Error("an edited cache.go should get a note, not a rewrite")
	}
}

func TestEnsureXSync(t *testing.T) {
	dir := t.TempDir()
	var fetched []string
	saved := goGet
	goGet = func(_ string, specs ...string) error { fetched = append(fetched, specs...); return nil }
	defer func() { goGet = saved }()

	write := func(body string) {
		if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("module x\n\nrequire (\n\tgolang.org/x/sync v0.19.0 // indirect\n)\n")
	if err := ensureXSync(dir); err != nil || len(fetched) != 0 {
		t.Errorf("an indirect x/sync should be enough: %v %v", err, fetched)
	}
	write("module x\n\nrequire (\n\tgolang.org/x/crypto v0.57.0\n)\n")
	if err := ensureXSync(dir); err != nil || len(fetched) != 1 || fetched[0] != "golang.org/x/sync@"+xSyncVersion {
		t.Errorf("a missing x/sync should be fetched: %v %v", err, fetched)
	}
}
