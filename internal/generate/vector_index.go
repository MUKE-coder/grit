package generate

// The index a vector column needs, as a checked-in SQL pair.
//
// grit#191 asked for an index flag, on the grounds that "HNSW and IVFFlat take
// parameters and an operator class (vector_cosine_ops and friends) that are
// easy to get subtly wrong". A pair written for you is better than a flag:
// nothing to look up, nothing to remember, and the alternatives are in the file
// as comments, so changing the method or the distance is editing one line.
//
// A pair rather than AutoMigrate because AutoMigrate cannot express it, and a
// pair is correctly ordered for it: pairs run after AutoMigrate, so the table
// and the column already exist. The extension is the opposite case and is
// handled before AutoMigrate, in internal/vectors.
//
// Postgres only. The pairs run on whatever database the project uses, and
// `USING hnsw` is a syntax error everywhere else, so on SQLite or MySQL no pair
// is written at all: there is no index to write, and the search still works by
// reading every row, which is what a development database does anyway.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// writeVectorIndexMigration writes the pair for one embedding column.
//
// Returns the pair's name and whether this call wrote it. A name with
// wrote=false is a pair an earlier run left, which is the --force case.
func (g *Generator) writeVectorIndexMigration(names Names, f Field) (name string, wrote bool, err error) {
	if !f.IsVector() {
		return "", false, nil
	}
	if projectDBProvider(g.Root) != "postgres" {
		return "", false, nil
	}

	table := names.PluralSnake
	column := toSnakeCase(f.Name)
	slug := migrationSlug(fmt.Sprintf("index_%s_%s", table, column))

	dir := filepath.Join(g.APIRoot(), "migrations")

	// Already there, from an earlier run of this generator, under whatever
	// timestamp that run used.
	//
	// Asking whether <this run's version>_<slug>.up.sql exists is not the same
	// question and was the bug: the version is new every run, so it only ever
	// caught a collision inside the same second, and `--force` wrote a second
	// copy of the same pair. Both then record as applied and either one's down
	// drops the shared index, leaving the history and the database disagreeing
	// with nothing to say so (grit#197).
	existing, err := existingMigrationFor(dir, slug)
	if err != nil {
		return "", false, err
	}
	if existing != "" {
		return existing + "_" + slug, false, nil
	}

	version, err := freeVersion(dir, time.Now())
	if err != nil {
		return "", false, err
	}
	up := filepath.Join(dir, version+"_"+slug+".up.sql")
	down := filepath.Join(dir, version+"_"+slug+".down.sql")

	indexName := fmt.Sprintf("idx_%s_%s", table, column)
	if err := writeFileWithDirs(up, vectorIndexUpSQL(version, slug, table, column, indexName, f.VectorDim())); err != nil {
		return "", false, fmt.Errorf("writing the vector index migration: %w", err)
	}
	if err := writeFileWithDirs(down, vectorIndexDownSQL(version, slug, indexName)); err != nil {
		return "", false, fmt.Errorf("writing the vector index rollback: %w", err)
	}
	return version + "_" + slug, true, nil
}

// vectorIndexUpSQL builds the index, or explains why it cannot.
func vectorIndexUpSQL(version, slug, table, column, indexName string, dim int) string {
	header := NoTransactionMarkerText + "\n" + migrationHeader(version, slug, "up")

	// hnsw and ivfflat both stop at 2000 dimensions. The column is fine wider
	// than that and the search still works by reading every row, so this is a
	// pair that explains itself rather than one that fails.
	if dim > vectorIndexMaxDimension {
		return header + fmt.Sprintf(`--
-- No index: %s has %d dimensions and both pgvector index types stop at %d.
--
-- The column is valid and the search works; it reads every row, which is exact
-- and gets slower as the table grows. The usual answer is to reduce the
-- dimension at the embedding model (text-embedding-3-small is 1536, and
-- OpenAI's API takes a "dimensions" parameter to shorten a larger one), because
-- that is cheaper than the scan and more accurate than no index.
--
-- Nothing to run. This file is here so the decision is recorded.
`, column, dim, vectorIndexMaxDimension)
	}

	return header + fmt.Sprintf(`--
-- The index that makes nearest-neighbour search fast on %s.%s.
--
-- CONCURRENTLY takes no write lock, so the table stays readable and writable
-- while the index builds, which is why this file opens with
-- %s. An hnsw build is not quick on a large table.
--
-- hnsw with vector_cosine_ops is the default for a reason: cosine distance
-- ignores magnitude, which is what you want for text embeddings, and hnsw
-- gives better recall per unit of query time than ivfflat without needing to
-- be told how many rows you have.
--
-- The operator class has to match the distance operator your query uses, or
-- the index is not used and nothing says so:
--
--     vector_cosine_ops  <=>  vectors.CosineDistance   (this file)
--     vector_l2_ops      <->  vectors.L2Distance
--     vector_ip_ops      <#>  vectors.InnerProduct
--
-- To use ivfflat instead, which builds faster and uses less memory at some
-- recall, replace the statement below with:
--
--     CREATE INDEX CONCURRENTLY IF NOT EXISTS %s
--         ON %s USING ivfflat (%s vector_cosine_ops) WITH (lists = 100);
--
-- ivfflat needs rows before it can build a useful index, so build it after the
-- table is populated, and set lists to about rows/1000.
CREATE INDEX CONCURRENTLY IF NOT EXISTS %s
    ON %s USING hnsw (%s vector_cosine_ops);
`, table, column, NoTransactionMarkerText, indexName, table, column, indexName, table, column)
}

// vectorIndexDownSQL drops it.
func vectorIndexDownSQL(version, slug, indexName string) string {
	return NoTransactionMarkerText + "\n" + migrationHeader(version, slug, "down") +
		fmt.Sprintf(`--
-- CONCURRENTLY here too: dropping an index takes a lock that blocks reads of
-- the table without it.
DROP INDEX CONCURRENTLY IF EXISTS %s;
`, indexName)
}

// vectorIndexMaxDimension is the ceiling for an hnsw or ivfflat index, which is
// lower than the column's own limit of 16000.
const vectorIndexMaxDimension = 2000

// existingMigrationFor finds a pair already written for this slug, whatever
// timestamp it carries, and returns its version or "".
//
// By suffix rather than by parsing the SQL: the slug is derived from the table
// and the column, so two pairs for the same index always share it.
func existingMigrationFor(dir, slug string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// No migrations directory yet, so nothing to collide with.
			return "", nil
		}
		return "", fmt.Errorf("reading %s: %w", dir, err)
	}
	suffix := "_" + slug + ".up.sql"
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), suffix) {
			return strings.TrimSuffix(e.Name(), suffix), nil
		}
	}
	return "", nil
}
