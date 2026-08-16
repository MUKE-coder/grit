package scaffold

import (
	"regexp"
	"strings"
)

// syncRegisterRe matches the registration the resource generator injects into
// routes.go: syncRegistry.Register("products", &models.Product{}).
var syncRegisterRe = regexp.MustCompile(`syncRegistry\.Register\(\s*"([a-z0-9_]+)"`)

// quotedList renders a TS string array literal from a slice of model names.
func quotedList(models []string) string {
	quoted := make([]string, len(models))
	for i, m := range models {
		quoted[i] = "\"" + m + "\""
	}
	return strings.Join(quoted, ", ")
}

// webSyncSetupTS emits apps/web/lib/sync.ts.
func webSyncSetupTS(project string, models []string) string {
	return `import { SyncEngine } from "@` + project + `/sync";
import { IndexedDBAdapter } from "@` + project + `/sync/adapters/indexed";
import { MemoryAdapter } from "@` + project + `/sync/adapters/memory";

/**
 * The offline mirror for the web app.
 *
 * The adapter is chosen at module load because the answer never changes
 * within a runtime: a browser has IndexedDB, a Next.js server render does
 * not, and reaching for it there throws during render rather than degrading.
 * The memory adapter gives the server pass something to read from that is
 * simply always empty.
 */
const adapter =
  typeof indexedDB === "undefined"
    ? new MemoryAdapter()
    : new IndexedDBAdapter("` + project + `-sync");

export const syncEngine = new SyncEngine({
  apiUrl: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080/api",
  adapter,
  // Plural snake_case, matching what the API registered in routes.go.
  models: [` + quotedList(models) + `],
  getToken: () =>
    typeof localStorage === "undefined"
      ? null
      : localStorage.getItem("access_token"),
});
`
}

// expoSyncSetupTS emits apps/expo/lib/sync.ts.
func expoSyncSetupTS(project string, models []string) string {
	return `import * as SecureStore from "expo-secure-store";
import * as SQLite from "expo-sqlite";

import { SyncEngine } from "@` + project + `/sync";
import { SQLiteAdapter } from "@` + project + `/sync/adapters/sqlite";

/**
 * The offline mirror for the mobile app.
 *
 * Built rather than exported, because opening the database is async and the
 * app wants that to happen behind a splash screen. Call this once at startup
 * and hand the result to <SyncProvider>.
 */
export async function createSyncEngine(): Promise<SyncEngine> {
  const db = await SQLite.openDatabaseAsync("` + project + `-sync.db");

  return new SyncEngine({
    apiUrl: process.env.EXPO_PUBLIC_API_URL ?? "http://localhost:8080/api",
    adapter: new SQLiteAdapter(db),
    // Plural snake_case, matching what the API registered in routes.go.
    models: [` + quotedList(models) + `],
    getToken: () => SecureStore.getItemAsync("access_token"),
    // Phones lose signal often enough that a slower loop just means a longer
    // window where the screen is confidently wrong.
    autoSyncIntervalMs: 15000,
  });
}
`
}
