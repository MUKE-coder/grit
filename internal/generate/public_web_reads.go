package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// The web half of --public.
//
// --public generated the API side and stopped there. The hooks beside it still
// called the authenticated routes, and nothing in apps/web ever sent the
// publishable key the seeder had already written into its .env.local, so a
// storefront built the obvious way got 401s from every request and the developer
// had to write the fetch layer by hand to find out why. Reported as grit#91.
//
// What this writes is the read layer for the public endpoints: the API key on
// every request, the published shape as a TypeScript type, and one function per
// endpoint the resource actually has.

// writePublicWebReads emits apps/web/lib/<plural>-public.ts, or nothing.
//
// Never overwritten, like the public handler it calls: the allowlist on the Go
// side is the developer's to edit, and a client typed against last week's
// allowlist is worse than no file.
func (g *Generator) writePublicWebReads(names Names) (string, error) {
	if !g.Definition.Public {
		return "", nil
	}
	dir, next := g.publicWebLibDir()
	if dir == "" {
		return "", nil
	}
	path := filepath.Join(dir, names.PluralKebab+"-public.ts")
	if fileExists(path) {
		return "", nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(g.publicWebReadsSource(names, next)), 0o644); err != nil {
		return "", err
	}
	return relativeToRoot(g.Root, path), nil
}

// publicWebLibDir is where the web app keeps its library code, and whether it is
// a Next.js app. Empty when the project has no web app, which is every --api and
// every single-app project.
func (g *Generator) publicWebLibDir() (dir string, next bool) {
	web := filepath.Join(g.Root, "apps", "web")
	if dirExists(filepath.Join(web, "src", "routes")) {
		return filepath.Join(web, "src", "lib"), false
	}
	if dirExists(filepath.Join(web, "lib")) {
		return filepath.Join(web, "lib"), true
	}
	return "", false
}

// publicTSFields is the published shape, as TypeScript, in the order the Go view
// struct declares it.
func (g *Generator) publicTSFields(included []Field) []string {
	out := []string{"  id: string;"}
	if g.Definition.Tree {
		out = append(out,
			"  parent_id: string;",
			"  depth: number;",
			"  descendant_ids?: string[];")
	}
	for _, f := range g.Definition.Fields {
		if !f.IsBelongsTo() || f.RelatedModelName() == toPascalCase(g.Definition.Name) {
			continue
		}
		out = append(out, fmt.Sprintf("  %s: string;", f.FKColumnName()))
	}
	for _, f := range included {
		out = append(out, fmt.Sprintf("  %s: %s;", toSnakeCase(f.Name), publicTSType(f)))
	}
	return out
}

// publicTSType maps a published field to the type the JSON actually carries.
//
// It has to agree with the Go view struct, which carries the model's own types
// rather than plausible guesses at them: an uploaded image is a FileRef object
// and not a URL string, and a client typed against the guess compiles and then
// renders "[object Object]".
func publicTSType(f Field) string {
	switch FieldType(f.Type) {
	case FieldInt, FieldUint, FieldFloat, FieldPercent, FieldRating:
		return "number"
	case FieldBool, FieldToggle:
		return "boolean"
	case FieldStringArray:
		return "string[]"
	case FieldFile:
		return "FileRef | null"
	case FieldFiles:
		return "FileRef[]"
	default:
		return "string"
	}
}

// publicNeedsFileRef reports whether the published shape carries an upload.
func publicNeedsFileRef(included []Field) bool {
	for _, f := range included {
		if t := FieldType(f.Type); t == FieldFile || t == FieldFiles {
			return true
		}
	}
	return false
}

func (g *Generator) publicWebReadsSource(names Names, next bool) string {
	included, _ := PublicFields(g.Definition.Fields)
	view := "Public" + names.Pascal
	base := "/api/" + apiVersion + "/public/" + names.Plural

	// Next reads these on the server, where the container talks to the API by
	// its compose service name; a Vite app only ever runs in a browser.
	apiURL := `const API_URL = (
  import.meta.env.VITE_API_URL ||
  "http://localhost:8080"
).replace(/\/+$/, "");

const API_KEY = import.meta.env.VITE_API_KEY ?? "";`
	fetchOpts := ""
	cacheOpen, cacheClose := "", ""
	imports := ""
	revalidate := ""
	if next {
		imports = "import { cache } from \"react\";\n\n"
		apiURL = `const API_URL = (
  process.env.API_INTERNAL_URL ||
  process.env.NEXT_PUBLIC_API_URL ||
  "http://localhost:8080"
).replace(/\/+$/, "");

// Publishable by design: it reaches the read-only public endpoints and nothing
// else, which is why the seeder writes it as NEXT_PUBLIC_.
const API_KEY = process.env.NEXT_PUBLIC_API_KEY ?? "";`
		revalidate = `
// Pages revalidate on this interval, so an edit in the admin shows within a
// minute without a deploy.
const REVALIDATE_SECONDS = 60;
`
		fetchOpts = ",\n      next: { revalidate: REVALIDATE_SECONDS }"
		cacheOpen, cacheClose = "cache(", ")"
	}

	sharedTypes := "PaginatedResponse"
	if publicNeedsFileRef(included) {
		sharedTypes = "FileRef, PaginatedResponse"
	}

	src := imports + `import type { ` + sharedTypes + ` } from "@repo/shared/types";

// Reads of the public ` + names.Plural + ` API, for a client with no signed-in user.
//
// Generated by grit with --public, and not overwritten when the resource is
// regenerated: the Go side's allowlist is yours to edit, and this file is typed
// against it.
//
// These call /public/` + names.Plural + `, not the authenticated routes the hooks in
// hooks/use-` + names.PluralKebab + `.ts call. The authenticated ones answer 401
// without a session, which is correct for the admin panel and useless here.

` + apiURL + `
` + revalidate + `
// ` + view + ` is the published shape: the allowlist in
// internal/handlers/` + names.Snake + `_public.go, not the model. Add a field
// there and add it here.
export interface ` + view + ` {
` + strings.Join(g.publicTSFields(included), "\n") + `
}

export type ` + view + `Page = {
  data: ` + view + `[];
  meta: PaginatedResponse<` + view + `>["meta"] | undefined;
};

async function publicGet<T>(path: string): Promise<{ status: number; body: T | null }> {
  try {
    const res = await fetch(API_URL + path, {
      headers: API_KEY ? { "X-API-Key": API_KEY } : {}` + fetchOpts + `,
    });
    if (!res.ok) return { status: res.status, body: null };
    return { status: res.status, body: (await res.json()) as T };
  } catch {
    // The API is unreachable, as it is while a production build runs in CI. A
    // list renders empty and fills in on the next request rather than failing
    // the build.
    return { status: 0, body: null };
  }
}

// One page of ` + names.Plural + `. Every search, sort, filter and range the
// public handler allows is a query parameter; anything it does not allow is
// ignored rather than rejected.
export const getPublic` + names.PluralPascal + ` = ` + cacheOpen + `async (
  params: Record<string, string | number | undefined> = {}
): Promise<` + view + `Page> => {
  const query = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== "") query.set(key, String(value));
  }
  const suffix = query.toString() ? "?" + query.toString() : "";
  const { body } = await publicGet<PaginatedResponse<` + view + `>>("` + base + `" + suffix);
  return { data: body?.data ?? [], meta: body?.meta };
}` + cacheClose + `;

// One ` + names.Lower + ` by its public key, or null when there is none.
//
// A row that is archived or switched off is a 404 here, the same as one that
// never existed: whether it exists is itself not public. An API that fails is
// an error rather than a missing row, so it is not reported as null.
export const getPublic` + names.Pascal + ` = ` + cacheOpen + `async (key: string): Promise<` + view + ` | null> => {
  const { status, body } = await publicGet<{ data: ` + view + ` }>(
    "` + base + `/" + encodeURIComponent(key)
  );
  if (status === 404) return null;
  if (!body) throw new Error("the public ` + names.Lower + ` API answered " + (status || "nothing") + " for " + key);
  return body.data;
}` + cacheClose + `;
`

	if _, ok := relatedParent(names, g.Definition); ok {
		src += `
// The "similar items" strip on a detail page. The server caps the limit, so
// asking for more than it allows returns what it allows.
export const getRelated` + names.PluralPascal + ` = ` + cacheOpen + `async (
  key: string,
  limit = 8
): Promise<` + view + `[]> => {
  const { body } = await publicGet<{ data: ` + view + `[] }>(
    "` + base + `/" + encodeURIComponent(key) + "/related?limit=" + limit
  );
  return body?.data ?? [];
}` + cacheClose + `;
`
	}

	if g.Definition.Tree {
		src += `
// The whole tree, assembled by the server: parents before children, siblings in
// the order the admin arranged. One request rather than one per level.
export interface ` + view + `Node extends ` + view + ` {
  children: ` + view + `Node[];
}

export const getPublic` + names.PluralPascal + `Tree = ` + cacheOpen + `async (): Promise<` + view + `Node[]> => {
  const { body } = await publicGet<{ data: ` + view + `Node[] }>("` + base + `/tree");
  return body?.data ?? [];
}` + cacheClose + `;
`
	}

	return src
}
