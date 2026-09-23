package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M21 and M22. The new text lives here once: templates
// splice it in, and repairPoolAndCache puts it into an existing project,
// anchored on what Grit generated before.

// ─── M21: connection pools sized for replicas and Sentinel ─────────────────

// dbPoolBlockOld is the pool block database.go carried before M21.
const dbPoolBlockOld = `	// Connection pool settings. SQLite ignores most of these — single-writer
	// semantics mean MaxOpenConns above 1 only helps concurrent reads, and
	// SQLite serialises writes internally. Postgres uses every knob.
	//
	// Idle defaults to Open, and that default matters more than it looks. When
	// idle is lower, a request past the idle limit returns its connection to a
	// full pool, so the connection is CLOSED — and the next request opens a new
	// one, which makes Postgres fork a backend process. Under concurrency that
	// is a connection storm, and it surfaces as database CPU rather than as
	// anything you would think to look for in the application.
	//
	// Measured with k6 at 50 VUs, 4 CPUs per container, single-row reads:
	// idle=10 gave ~810 req/s with Postgres pinned near 840% while the API used
	// 196%; idle=100 gave ~2,720 req/s with both around 300%. Same binary, same
	// query.
	//
	// Both are tunable because the right answer depends on the workload. If
	// your queries are heavy enough to saturate the database — an unindexed
	// COUNT over a large table on every request, say — a smaller pool acts as
	// admission control and can measure faster, because queueing in the app is
	// cheaper than thrashing in Postgres. Start here, then measure.
	maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 100)
	if maxOpen < 1 {
		maxOpen = 1
	}
	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", maxOpen)
	if maxIdle < 1 || maxIdle > maxOpen {
		maxIdle = maxOpen
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
`

// dbPoolBlockNew sizes the app's pool for a server shared by every replica and
// by Sentinel's own pool in each of them.
const dbPoolBlockNew = `	// Connection pool settings. Postgres and MySQL use every knob; SQLite is
	// handled at the bottom, and not the way this comment used to claim.
	//
	// The number that matters is the total across processes, not this one.
	// Every replica opens its own pool, and Sentinel opens a second one in each
	// (SENTINEL_DB_MAX_OPEN_CONNS, see SentinelPoolSize), all against the
	// server's max_connections: 100 on a default Postgres, 151 on MySQL, a few
	// of them reserved for superusers. The default here was 100, so one replica
	// could take the whole server and three asked for 330 connections. Size it
	// as
	//
	//	DB_MAX_OPEN_CONNS = (max_connections - reserve) / replicas - SENTINEL_DB_MAX_OPEN_CONNS
	//
	// The defaults, 25 here and 5 for Sentinel, are that sum for three replicas
	// on a default Postgres with 10 left for migrations and psql. Behind
	// pgbouncer, as docker-compose.prod.yml runs it, Postgres only ever sees
	// pgbouncer's DEFAULT_POOL_SIZE, and these bound how much one process queues
	// there instead.
	//
	// Idle defaults to Open, and that default matters more than it looks. When
	// idle is lower, a request past the idle limit returns its connection to a
	// full pool, so the connection is closed, and the next request opens a new
	// one, which makes Postgres fork a backend process. Under concurrency that
	// is a connection storm, and it surfaces as database CPU rather than as
	// anything you would think to look for in the application.
	//
	// If your queries are heavy enough to saturate the database (an unindexed
	// COUNT over a large table on every request, say) a smaller pool acts as
	// admission control and can measure faster, because queueing in the app is
	// cheaper than thrashing in Postgres. Start here, then measure.
	maxOpen := getEnvInt("DB_MAX_OPEN_CONNS", 25)
	if maxOpen < 1 {
		maxOpen = 1
	}
	maxIdle := getEnvInt("DB_MAX_IDLE_CONNS", maxOpen)
	if maxIdle < 1 || maxIdle > maxOpen {
		maxIdle = maxOpen
	}

	// SQLite takes one connection, whatever the numbers above say. It was long
	// written here that SQLite "ignores most of these", and it does not: a pool
	// of several connections is several writers competing for one file lock, so
	// a transaction holding the write lock on one connection and a statement
	// wanting it on another wait for each other until busy_timeout gives up.
	// What comes back is "database is locked" on a machine doing almost
	// nothing: a background worker ticking while a webhook writes was enough to
	// lose the write, and it reads as a load problem it is not.
	//
	// One connection makes that a queue inside the process, which is what
	// SQLite wants. It is also the only correct setting for :memory:, where a
	// second connection is a second, empty database.
	if strings.HasPrefix(dsn, "sqlite:") {
		maxOpen, maxIdle = 1, 1
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	if maxOpen == 1 && strings.HasPrefix(dsn, "sqlite:") {
		log.Println("Database pool: one connection, because SQLite writes one at a time")
	} else {
		log.Printf("Database pool: up to %d connections from this process (DB_MAX_OPEN_CONNS)", maxOpen)
	}
`

// sentinelPoolFunc is appended to database.go.
const sentinelPoolFunc = `
// SentinelPoolSize is the connection pool Sentinel opens for its own storage.
// On Postgres that is a second pool on the app's server, in every replica, so
// it counts against max_connections the same way DB_MAX_OPEN_CONNS does.
// Sentinel writes from a batching background pipeline and needs few.
func SentinelPoolSize() (maxOpen, maxIdle int) {
	maxOpen = getEnvInt("SENTINEL_DB_MAX_OPEN_CONNS", 5)
	if maxOpen < 1 {
		maxOpen = 1
	}
	maxIdle = getEnvInt("SENTINEL_DB_MAX_IDLE_CONNS", maxOpen)
	if maxIdle < 1 || maxIdle > maxOpen {
		maxIdle = maxOpen
	}
	return maxOpen, maxIdle
}
`

const (
	sentinelPoolAnchor = "\t\tsentinelStorage.AuditKey = cfg.SentinelAuditKey\n"
	sentinelPoolLine   = "\t\t// Sentinel's storage is a second connection pool on the same server, in\n" +
		"\t\t// every replica, so it is capped like the app's (database.SentinelPoolSize).\n" +
		"\t\tsentinelStorage.MaxOpenConns, sentinelStorage.MaxIdleConns = database.SentinelPoolSize()\n"
)

// The production api service connects through pgbouncer.
const (
	composeAPICommentOld = `    # POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_DB come from .env
    # (env_file above). We override POSTGRES_HOST so the API hits the
    # postgres container on the Docker network, and POSTGRES_PORT back
    # to 5432 because the postgres container listens on 5432 internally
    # — only the dev host port mapping uses 5434. The Go config builds
    # DATABASE_URL from these parts at startup.
`
	composeAPICommentNew = `    # POSTGRES_USER / POSTGRES_PASSWORD / POSTGRES_DB come from .env
    # (env_file above). POSTGRES_HOST and POSTGRES_PORT are overridden to
    # reach pgbouncer on the Docker network, whatever host port .env names
    # for development. The Go config builds DATABASE_URL from these parts at
    # startup.
`
	composeAPIPostgresOld = "      POSTGRES_HOST: postgres\n      POSTGRES_PORT: \"5432\"\n"
	composeAPIPostgresNew = "      # Through pgbouncer, so replicas share its DEFAULT_POOL_SIZE server\n" +
		"      # connections instead of each holding its own pool open on Postgres.\n" +
		"      POSTGRES_HOST: pgbouncer\n" +
		"      POSTGRES_PORT: \"6432\"\n"
	composeAPIDependsOld = "    depends_on:\n      postgres:\n        condition: service_healthy\n"
	composeAPIDependsNew = "    depends_on:\n      pgbouncer:\n        condition: service_healthy\n"

	// The image listens on 5432 unless LISTEN_PORT says otherwise, so the
	// health check on 6432 never passed.
	pgbouncerListenOld = "      POOL_MODE: transaction\n"
	pgbouncerListenNew = "      # The image listens on 5432 unless told otherwise, and the health check\n" +
		"      # and the api use 6432.\n" +
		"      LISTEN_PORT: 6432\n" +
		"      POOL_MODE: transaction\n"

	pgbouncerUseOld = "  # To use it, set in .env:  DB_HOST=pgbouncer  DB_PORT=6432\n"
	pgbouncerUseNew = "  # The api service connects through it. Grit's queries use the simple\n" +
		"  # protocol, and pgbouncer 1.21 and later carries the prepared statements\n" +
		"  # Sentinel's driver sends (max_prepared_statements, 200 by default).\n" +
		"  # MAX_CLIENT_CONN covers DB_MAX_OPEN_CONNS + SENTINEL_DB_MAX_OPEN_CONNS\n" +
		"  # (30 by default) per replica.\n"
)

const (
	envPoolOld = `# Connection pool. Idle defaults to Open, which is what you want: when idle is
# lower, connections returned to a full pool get closed and Postgres forks a new
# backend for the next request. Under load that is a connection storm, and it
# shows up as database CPU rather than as anything obvious in the app.
#
# Measured at 50 concurrent users, 4 CPUs a side: idle=10 gave ~810 req/s on a
# single-row read; idle=100 gave ~2,720. Lower them only if your queries are
# heavy enough to saturate the database, where a smaller pool acts as admission
# control — then queueing happens in the app instead of thrashing in Postgres.
# DB_MAX_OPEN_CONNS=100
# DB_MAX_IDLE_CONNS=100
`
	envPoolNew = `# Connection pools. Each API process opens one for the app and one for Sentinel,
# and every replica does the same against one max_connections (100 on a default
# Postgres). Size them so the total fits:
#   (DB_MAX_OPEN_CONNS + SENTINEL_DB_MAX_OPEN_CONNS) x replicas < max_connections - 10
# The defaults, 25 + 5, fit three replicas. Behind pgbouncer (the production
# compose file) Postgres sees only pgbouncer's DEFAULT_POOL_SIZE.
# Idle defaults to Open: a lower idle closes connections under load and Postgres
# forks a new backend for each one it reopens.
# DB_MAX_OPEN_CONNS=25
# DB_MAX_IDLE_CONNS=25
# SENTINEL_DB_MAX_OPEN_CONNS=5
# SENTINEL_DB_MAX_IDLE_CONNS=5
`
)

// ─── M22: response cache without a stampede or base64 bodies ───────────────

// xSyncVersion is golang.org/x/sync, for singleflight. Already in every API's
// module graph as an indirect dependency; the cache middleware imports it.
const xSyncVersion = "v0.23.0"

// cacheMiddlewareOld is middleware/cache.go before M22.
const cacheMiddlewareOld = `package middleware

import (
	"bytes"
	"hash/fnv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"{{MODULE}}/internal/cache"
)

// CacheResponse caches GET request responses in Redis for the given duration.
// Only caches 200 OK responses. Skips caching if no cache service is available.
func CacheResponse(cacheService *cache.Cache, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cacheService == nil || c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Build cache key from URL + query params. We use FNV-1a (non-
		// cryptographic) instead of SHA-256 because cache keys don't
		// need cryptographic strength and FNV is ~50x faster on the
		// hot path of every cacheable request.
		h := fnv.New64a()
		h.Write([]byte(c.Request.URL.String()))
		key := "http:" + strconv.FormatUint(h.Sum64(), 16)

		// Try to serve from cache
		var cached cachedResponse
		found, err := cacheService.Get(c.Request.Context(), key, &cached)
		if err == nil && found {
			c.Header("X-Cache", "HIT")
			c.Data(cached.Status, cached.ContentType, cached.Body)
			c.Abort()
			return
		}

		// Capture the response. bytes.Buffer grows in chunks (vs []byte
		// append which can reallocate on every Write), so a 100 KB
		// response only takes ~3 allocations instead of one per chunk.
		writer := &responseCapture{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
		c.Writer = writer
		c.Header("X-Cache", "MISS")

		c.Next()

		// Cache successful responses
		if writer.status == http.StatusOK && writer.body.Len() > 0 {
			resp := cachedResponse{
				Status:      writer.status,
				ContentType: writer.Header().Get("Content-Type"),
				Body:        writer.body.Bytes(),
			}
			_ = cacheService.Set(c.Request.Context(), key, resp, ttl)
		}
	}
}

type cachedResponse struct {
	Status      int    ` + "`" + `json:"status"` + "`" + `
	ContentType string ` + "`" + `json:"content_type"` + "`" + `
	Body        []byte ` + "`" + `json:"body"` + "`" + `
}

type responseCapture struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseCapture) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseCapture) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
`

// cacheMiddlewareNew is middleware/cache.go.
const cacheMiddlewareNew = `package middleware

import (
	"bytes"
	"context"
	"hash/fnv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"

	"{{MODULE}}/internal/cache"
)

// cacheFlights collapses concurrent misses on one key into one handler run.
// Without it an entry expiring on a busy endpoint sent every request that
// arrived before the first one finished to the database at once, and each of
// them wrote the entry again.
var cacheFlights singleflight.Group

// CacheResponse caches GET request responses in Redis for the given duration.
// Only caches 200 OK responses. Skips caching if no cache service is available.
func CacheResponse(cacheService *cache.Cache, ttl time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cacheService == nil || c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Build cache key from URL + query params. FNV-1a rather than SHA-256:
		// cache keys don't need cryptographic strength and FNV is much faster
		// on the hot path of every cacheable request. v2 is the raw layout
		// below; entries in the earlier JSON layout are never read, and expire.
		h := fnv.New64a()
		h.Write([]byte(c.Request.URL.String()))
		key := "http:v2:" + strconv.FormatUint(h.Sum64(), 16)
		ctx := c.Request.Context()

		if cached, ok := loadCachedResponse(ctx, cacheService, key); ok {
			serveCachedResponse(c, cached)
			return
		}

		// One request per key runs the handler; the others wait for its
		// response. Only the one that ran it has written to its client.
		ran := false
		shared, _, _ := cacheFlights.Do(key, func() (interface{}, error) {
			// A flight that ended between the lookup above and this one has
			// stored the entry already.
			if cached, ok := loadCachedResponse(ctx, cacheService, key); ok {
				return cached, nil
			}
			ran = true
			return captureAndStore(ctx, c, cacheService, key, ttl), nil
		})
		if ran {
			return
		}
		if cached, ok := shared.(*cachedResponse); ok && cached != nil {
			serveCachedResponse(c, cached)
			return
		}
		// The handler's answer was not cacheable (an error, a redirect), so this
		// request gets its own.
		c.Next()
	}
}

// captureAndStore runs the rest of the chain, writing to the client as usual,
// and stores a 200 response. It returns the response, or nil when it was not
// one to cache.
func captureAndStore(ctx context.Context, c *gin.Context, cacheService *cache.Cache, key string, ttl time.Duration) *cachedResponse {
	// bytes.Buffer grows in chunks, so a 100 KB response takes a few
	// allocations instead of one per Write.
	writer := &responseCapture{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
	c.Writer = writer
	c.Header("X-Cache", "MISS")

	c.Next()

	if writer.Status() != http.StatusOK || writer.body.Len() == 0 {
		return nil
	}
	resp := &cachedResponse{
		Status:      http.StatusOK,
		ContentType: writer.Header().Get("Content-Type"),
		Body:        writer.body.Bytes(),
	}
	// A failed write only costs the next request a miss.
	_ = cacheService.Client().Set(ctx, key, resp.encode(), ttl).Err()
	return resp
}

func serveCachedResponse(c *gin.Context, cached *cachedResponse) {
	c.Header("X-Cache", "HIT")
	c.Data(cached.Status, cached.ContentType, cached.Body)
	c.Abort()
}

// loadCachedResponse reads an entry. A missing key, a Redis error and an entry
// it cannot parse all mean the request is served uncached.
func loadCachedResponse(ctx context.Context, cacheService *cache.Cache, key string) (*cachedResponse, bool) {
	raw, err := cacheService.Client().Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	return decodeCachedResponse(raw)
}

type cachedResponse struct {
	Status      int
	ContentType string
	Body        []byte
}

// encode lays an entry out as one line of status and content type, then the
// body as it was sent: "200 application/json; charset=utf-8\n{...}". JSON with a
// []byte field stored the body as base64, a third larger in Redis, and cost an
// encode on every miss and a decode on every hit.
func (r *cachedResponse) encode() []byte {
	head := strconv.Itoa(r.Status) + " " + r.ContentType + "\n"
	out := make([]byte, 0, len(head)+len(r.Body))
	out = append(out, head...)
	return append(out, r.Body...)
}

func decodeCachedResponse(raw []byte) (*cachedResponse, bool) {
	nl := bytes.IndexByte(raw, '\n')
	if nl < 0 {
		return nil, false
	}
	sp := bytes.IndexByte(raw[:nl], ' ')
	if sp < 0 {
		return nil, false
	}
	status, err := strconv.Atoi(string(raw[:sp]))
	if err != nil {
		return nil, false
	}
	return &cachedResponse{Status: status, ContentType: string(raw[sp+1 : nl]), Body: raw[nl+1:]}, true
}

type responseCapture struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseCapture) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString is captured too: gin's ResponseWriter has it, so a handler that
// writes a string would otherwise bypass Write and reach the client uncached.
func (w *responseCapture) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
`

var cacheImportModule = regexp.MustCompile(`(?m)^\t"([^"\n]+)/internal/cache"$`)

// repairPoolAndCache applies M21 and M22 to an existing project.
func repairPoolAndCache(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	internal := filepath.Join(apiRoot, "internal")
	if path := filepath.Join(internal, "database", "database.go"); fileExists(path) {
		if err := repairSourceFile(root, m, path, repairDBPoolSource); err != nil {
			return err
		}
	}
	if path := filepath.Join(internal, "routes", "routes.go"); fileExists(path) && fileContains(filepath.Join(internal, "database", "database.go"), "func SentinelPoolSize(") {
		if err := repairSourceFile(root, m, path, repairSentinelPoolSource); err != nil {
			return err
		}
	}
	if path := filepath.Join(root, "docker-compose.prod.yml"); fileExists(path) {
		if err := repairTextFile(root, m, path, repairComposePgbouncerSource); err != nil {
			return err
		}
	}
	if path := filepath.Join(root, ".env.example"); fileExists(path) {
		if err := repairTextFile(root, m, path, repairEnvPoolSource); err != nil {
			return err
		}
	}

	path := filepath.Join(internal, "middleware", "cache.go")
	if !fileExists(path) {
		return nil
	}
	if err := ensureXSync(apiRoot); err != nil {
		return err
	}
	return repairSourceFile(root, m, path, repairCacheMiddlewareSource)
}

func repairDBPoolSource(src string) (string, []string, []string) {
	var changes []string
	out := src
	if !strings.Contains(src, "SENTINEL_DB_MAX_OPEN_CONNS") {
		if strings.Count(src, dbPoolBlockOld) != 1 {
			return src, nil, []string{"the connection pool block is not the one Grit wrote: default DB_MAX_OPEN_CONNS to (max_connections - reserve) / replicas - Sentinel's pool, and add SentinelPoolSize for routes.go"}
		}
		out = strings.Replace(out, dbPoolBlockOld, dbPoolBlockNew, 1)
		out = strings.TrimRight(out, "\n") + "\n" + sentinelPoolFunc
		changes = append(changes, "the pool defaults to 25 connections, and Sentinel's to 5 (SENTINEL_DB_MAX_OPEN_CONNS)")
	}
	return out, changes, nil
}

func repairSentinelPoolSource(src string) (string, []string, []string) {
	if strings.Contains(src, "database.SentinelPoolSize()") || !strings.Contains(src, "sentinel.StorageConfig{") {
		return src, nil, nil
	}
	if strings.Count(src, sentinelPoolAnchor) != 1 || !strings.Contains(src, "/internal/database\"") {
		return src, nil, []string{"the Sentinel storage block is not the one Grit wrote: set sentinelStorage.MaxOpenConns and MaxIdleConns from database.SentinelPoolSize(), or Sentinel opens 10 more connections per replica"}
	}
	return strings.Replace(src, sentinelPoolAnchor, sentinelPoolAnchor+sentinelPoolLine, 1),
		[]string{"Sentinel's connection pool is capped by SENTINEL_DB_MAX_OPEN_CONNS"}, nil
}

func repairComposePgbouncerSource(src string) (string, []string, []string) {
	if !strings.Contains(src, "\n  pgbouncer:\n") {
		return src, nil, nil
	}
	out := strings.Replace(src, pgbouncerUseOld, pgbouncerUseNew, 1)
	var changes []string

	// LISTEN_PORT first: the api is only pointed at 6432 once pgbouncer listens there.
	pbStart, pbEnd, ok := serviceBlock(out, "pgbouncer")
	if !ok {
		return src, nil, nil
	}
	pb := out[pbStart:pbEnd]
	if !strings.Contains(pb, "LISTEN_PORT:") {
		if strings.Count(pb, pgbouncerListenOld) != 1 {
			return src, nil, []string{"the pgbouncer service is not the one Grit wrote: set LISTEN_PORT: 6432, or it listens on 5432 and its health check on 6432 never passes"}
		}
		out = out[:pbStart] + strings.Replace(pb, pgbouncerListenOld, pgbouncerListenNew, 1) + out[pbEnd:]
		changes = append(changes, "pgbouncer listens on 6432, where its health check looks")
	}

	start, end, ok := serviceBlock(out, "api")
	if !ok {
		return src, nil, nil
	}
	block := out[start:end]
	if !strings.Contains(block, "POSTGRES_HOST: pgbouncer") {
		if strings.Count(block, composeAPIPostgresOld) != 1 || strings.Count(block, composeAPIDependsOld) != 1 {
			return out, changes, []string{"the api service is not the one Grit wrote: set POSTGRES_HOST: pgbouncer and POSTGRES_PORT: \"6432\" so replicas share pgbouncer's server connections"}
		}
		block = strings.Replace(block, composeAPICommentOld, composeAPICommentNew, 1)
		block = strings.Replace(block, composeAPIPostgresOld, composeAPIPostgresNew, 1)
		block = strings.Replace(block, composeAPIDependsOld, composeAPIDependsNew, 1)
		out = out[:start] + block + out[end:]
		changes = append(changes, "the production api connects to Postgres through pgbouncer")
	}
	if out == src {
		return src, nil, nil
	}
	if len(changes) == 0 {
		changes = append(changes, "the pgbouncer notes describe how the api uses it")
	}
	return out, changes, nil
}

func repairEnvPoolSource(src string) (string, []string, []string) {
	if strings.Contains(src, "SENTINEL_DB_MAX_OPEN_CONNS") || !strings.Contains(src, envPoolOld) {
		return src, nil, nil
	}
	return strings.Replace(src, envPoolOld, envPoolNew, 1), []string{"documents DB_MAX_OPEN_CONNS and SENTINEL_DB_MAX_OPEN_CONNS for replicas"}, nil
}

func repairCacheMiddlewareSource(src string) (string, []string, []string) {
	if strings.Contains(src, "singleflight") {
		// The first singleflight version took the context from c inside
		// captureAndStore, which golangci-lint's contextcheck reports on
		// every new project. The same context, passed in.
		out := src
		for old, updated := range map[string]string{
			"return captureAndStore(c, cacheService, key, ttl), nil":                            "return captureAndStore(ctx, c, cacheService, key, ttl), nil",
			"func captureAndStore(c *gin.Context, cacheService":                                 "func captureAndStore(ctx context.Context, c *gin.Context, cacheService",
			"_ = cacheService.Client().Set(c.Request.Context(), key, resp.encode(), ttl).Err()": "_ = cacheService.Client().Set(ctx, key, resp.encode(), ttl).Err()",
		} {
			if strings.Count(out, old) != 1 {
				return src, nil, nil
			}
			out = strings.Replace(out, old, updated, 1)
		}
		return out, []string{"cache.go passes the request context to captureAndStore"}, nil
	}
	mod := cacheImportModule.FindStringSubmatch(src)
	if mod == nil || src != strings.ReplaceAll(cacheMiddlewareOld, "{{MODULE}}", mod[1]) {
		return src, nil, []string{"cache.go is not the file Grit wrote: collapse concurrent misses with singleflight, and store the body raw rather than as JSON with a base64 body"}
	}
	return strings.ReplaceAll(cacheMiddlewareNew, "{{MODULE}}", mod[1]),
		[]string{"concurrent misses on one key run the handler once, and entries store the body raw"}, nil
}

// ensureXSync adds golang.org/x/sync to go.mod when it is missing. Every API
// has it as an indirect dependency, so this is for one that somehow does not.
func ensureXSync(apiRoot string) error {
	data, err := os.ReadFile(filepath.Join(apiRoot, "go.mod"))
	if err != nil {
		return nil
	}
	if regexp.MustCompile(`(?m)^\s*(?:require\s+)?golang\.org/x/sync\s+v`).Match(data) {
		return nil
	}
	if err := goGet(apiRoot, "golang.org/x/sync@"+xSyncVersion); err != nil {
		return fmt.Errorf("could not add golang.org/x/sync for the cache middleware, so run `go get golang.org/x/sync@%s` in %s: %v", xSyncVersion, apiRoot, err)
	}
	return nil
}
