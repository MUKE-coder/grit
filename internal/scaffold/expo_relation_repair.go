package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// A generated Expo screen for a resource with a belongs_to field imported a
// use<Plural> hook for the related model and showed it as
// (item.x && (item.x.name || item.x.title)). For a relation to the built-in
// User there was no use-users hook, and most models have neither name nor
// title, so the app failed tsc as soon as one such resource was generated.
// Found building the WhatsApp blueprint. The generator writes these two files
// and relationLabel calls now; repairExpoRelations brings screens it already
// wrote up to them.

// ExpoRelationLabel is apps/expo/lib/relation-label.ts.
func ExpoRelationLabel() string {
	return `// The name to show for a related record: its name, title, a person's first and
// last name, or its email, whichever it has. The generated screens call it for
// every belongs_to field, so they never read a property the model lacks.
export function relationLabel(value: unknown): string {
  if (!value || typeof value !== "object") return "";
  const record = value as Record<string, unknown>;
  const text = (key: string) => (typeof record[key] === "string" ? (record[key] as string).trim() : "");
  const person = [text("first_name"), text("last_name")].filter(Boolean).join(" ");
  return text("name") || text("title") || person || text("label") || text("email");
}
`
}

// ExpoUsersHook is apps/expo/hooks/use-users.ts: the list a relationship
// picker needs when a resource belongs to a user. Listing users is an admin
// route, as every generated resource's is; managing users stays in the admin.
func ExpoUsersHook() string {
	return `import { useInfiniteQuery } from "@tanstack/react-query";
import type { User } from "@repo/shared/types";
import { api } from "@/lib/api";

export type { User };

export interface UsersPage {
  data: User[];
  meta: { total: number; page: number; page_size: number; pages: number };
}

// Users for a relationship picker or filter, with the same arguments as a
// generated resource's list hook.
export function useUsers(
  search = "",
  filters: Record<string, string> = {},
  sortBy = "created_at",
  sortOrder: "asc" | "desc" = "desc",
  pageSize = 20,
) {
  return useInfiniteQuery({
    queryKey: ["users", { search, filters, sortBy, sortOrder, pageSize }],
    initialPageParam: 1,
    queryFn: async ({ pageParam }) => {
      const qs = new URLSearchParams({
        page: String(pageParam),
        page_size: String(pageSize),
        sort_by: sortBy,
        sort_order: sortOrder,
      });
      if (search) qs.set("search", search);
      for (const [k, v] of Object.entries(filters)) if (v) qs.set(k, v);
      return (await api.get("/users?" + qs.toString())) as UsersPage;
    },
    getNextPageParam: (last) =>
      last.meta.page < last.meta.pages ? last.meta.page + 1 : undefined,
  });
}
`
}

// relationNameTitle matches the expression the generator used to write,
// (x.rel && (x.rel.name || x.rel.title)). RE2 has no backreferences, so the
// three mentions of x.rel are captured separately and compared in
// replaceRelationNameTitle.
var relationNameTitle = regexp.MustCompile(`\(([A-Za-z_]\w*)\.(\w+) && \(([A-Za-z_]\w*)\.(\w+)\.name \|\| ([A-Za-z_]\w*)\.(\w+)\.title\)\)`)

// replaceRelationNameTitle rewrites each match whose three mentions agree.
func replaceRelationNameTitle(src string) (string, bool) {
	changed := false
	out := relationNameTitle.ReplaceAllStringFunc(src, func(m string) string {
		g := relationNameTitle.FindStringSubmatch(m)
		if g[1] != g[3] || g[1] != g[5] || g[2] != g[4] || g[2] != g[6] {
			return m
		}
		changed = true
		return "relationLabel(" + g[1] + "." + g[2] + ")"
	})
	return out, changed
}

const relationLabelImport = "import { relationLabel } from \"@/lib/relation-label\";\n"

// repairExpoRelationSource replaces the name-or-title expression with
// relationLabel and imports it.
func repairExpoRelationSource(src string) (string, []string, []string) {
	out, changed := replaceRelationNameTitle(src)
	if !changed {
		return src, nil, nil
	}
	if !strings.Contains(out, relationLabelImport) {
		// After the first import, which every generated screen starts with.
		i := strings.Index(out, "\n")
		if !strings.HasPrefix(out, "import ") || i < 0 {
			return src, nil, []string{"this screen does not start with an import: import relationLabel from @/lib/relation-label and use it for related records"}
		}
		out = out[:i+1] + relationLabelImport + out[i+1:]
	}
	return out, []string{"related records are named by relationLabel, not by a .name or .title the model may lack"}, nil
}

// repairExpoRelations fixes the generated screens of an existing Expo app.
func repairExpoRelations(root string) error {
	expo := filepath.Join(root, "apps", "expo")
	if !dirExists(expo) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	needsUsers, repaired := false, false
	for _, dir := range []string{filepath.Join(expo, "app"), filepath.Join(expo, "components", "resource-forms")} {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".tsx") {
				return nil
			}
			raw, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			src := string(raw)
			if strings.Contains(src, `from "@/hooks/use-users"`) {
				needsUsers = true
			}
			if relationNameTitle.MatchString(src) {
				repaired = true
				if err := repairTextFile(root, m, path, repairExpoRelationSource); err != nil {
					return err
				}
			}
			if concatenatedRoute.MatchString(src) {
				return repairTextFile(root, m, path, repairExpoRouteSource)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	if repaired {
		if err := writeIfMissing(filepath.Join(expo, "lib", "relation-label.ts"), ExpoRelationLabel()); err != nil {
			return err
		}
	}
	if needsUsers {
		if err := writeIfMissing(filepath.Join(expo, "hooks", "use-users.ts"), ExpoUsersHook()); err != nil {
			return err
		}
	}
	return nil
}

func writeIfMissing(path, content string) error {
	if fileExists(path) {
		return nil
	}
	return writeFile(path, content)
}

// concatenatedRoute matches router.push("/things/" + item.id). Expo Router's
// typed routes, on in every Grit Expo app, accept a template literal whose
// shape matches a route but not a plain string, so once expo start had written
// the route types, every generated list screen and the roles screen failed tsc.
var concatenatedRoute = regexp.MustCompile(`router\.(push|replace)\("(/[A-Za-z0-9_\-/]*/)" \+ ([A-Za-z_][\w.]*)\)`)

// repairExpoRouteSource turns each concatenated route into a template literal.
func repairExpoRouteSource(src string) (string, []string, []string) {
	out := concatenatedRoute.ReplaceAllString(src, "router.$1(`$2$${$3}`)")
	if out == src {
		return src, nil, nil
	}
	return out, []string{"navigates with a template literal, which Expo Router's typed routes accept"}, nil
}

// DesktopUsersHook is apps/desktop/frontend/src/hooks/use-users.ts: the users a
// relationship picker offers. Users are not in the desktop's offline mirror, so
// this one reads the API, and is empty while offline.
func DesktopUsersHook() string {
	return `import { useQuery } from "@tanstack/react-query";
import type { User } from "@repo/shared/types";
import { apiClient } from "@/lib/api-client";

export type { User };

// Users for a relationship picker. Not in the offline mirror, so read from the
// API: while offline the picker is empty and the form keeps what it had.
export function useUsers() {
  return useQuery<User[]>({
    queryKey: ["users", "options"],
    queryFn: async () => (await apiClient.get("/users?page_size=500")).data.data ?? [],
    staleTime: 60_000,
  });
}
`
}

// desktopFormInputImport finds the resource's input type in a generated
// desktop form: import type { Message, MessageInput } from "@/hooks/use-messages".
var desktopFormInputImport = regexp.MustCompile(`import type \{ \w+, (\w+Input) \} from "@/hooks/use-`)

// desktopFormPayload is the object a generated desktop form submits.
var desktopFormPayload = regexp.MustCompile(`(?s)\n    await onSubmit\(\{\n(.*?)\n    \}\);`)

// repairDesktopFormSource asserts a generated desktop form's payload to the
// resource's input type. The form held a select as a string where the model
// types it as its options, and a file as the desktop's looser FileRef, so a
// resource with either failed tsc.
func repairDesktopFormSource(src string) (string, []string, []string) {
	m := desktopFormInputImport.FindStringSubmatch(src)
	if m == nil || !desktopFormPayload.MatchString(src) || strings.Contains(src, "} as "+m[1]+");") {
		return src, nil, nil
	}
	out := desktopFormPayload.ReplaceAllString(src, "\n    await onSubmit({\n$1\n    } as "+m[1]+");")
	return out, []string{"the form's payload is typed as the resource's input, so a select or file field type-checks"}, nil
}

// repairDesktopForms fixes the generated forms of an existing desktop app and
// adds the users hook a relation to User imports.
func repairDesktopForms(root string) error {
	src := filepath.Join(root, "apps", "desktop", "frontend", "src")
	forms := filepath.Join(src, "components", "resource-forms")
	if !dirExists(forms) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(forms)
	if err != nil {
		return err
	}
	needsUsers := false
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".tsx") {
			continue
		}
		path := filepath.Join(forms, e.Name())
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(raw), `from "@/hooks/use-users"`) {
			needsUsers = true
		}
		if err := repairTextFile(root, m, path, repairDesktopFormSource); err != nil {
			return err
		}
	}
	lists := filepath.Join(src, "routes", "app")
	if listEntries, err := os.ReadDir(lists); err == nil {
		for _, e := range listEntries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".index.tsx") {
				continue
			}
			if err := repairTextFile(root, m, filepath.Join(lists, e.Name()), repairDesktopListSource); err != nil {
				return err
			}
		}
	}
	if needsUsers {
		return writeIfMissing(filepath.Join(src, "hooks", "use-users.ts"), DesktopUsersHook())
	}
	return nil
}

// desktopListRelationField is a generated desktop list's display name for a
// relation, written over the relation's own field: conversation: conversationMap.get(...
var desktopListRelationField = regexp.MustCompile(`([a-z_][a-z0-9_]*): (\w+)Map\.get\(String\(\(r as any\)\.`)

// repairDesktopListSource moves a relation's display name to <field>_label. The
// model types <field> as the related record, so rows carrying a string there
// did not type as the model and the list failed tsc.
func repairDesktopListSource(src string) (string, []string, []string) {
	renamed := map[string]bool{}
	out := desktopListRelationField.ReplaceAllStringFunc(src, func(m string) string {
		g := desktopListRelationField.FindStringSubmatch(m)
		if strings.HasSuffix(g[1], "_label") {
			return m
		}
		renamed[g[1]] = true
		return g[1] + "_label: " + g[2] + "Map.get(String((r as any)."
	})
	if len(renamed) == 0 {
		return src, nil, nil
	}
	for field := range renamed {
		out = strings.ReplaceAll(out, `{ key: "`+field+`", `, `{ key: "`+field+`_label", `)
	}
	return out, []string{"related records' names go in their own column key, so the rows type-check"}, nil
}

// desktopTSConfigNodeOutDir sends what tsc -b emits for vite.config.ts into
// node_modules. A composite project with no outDir wrote vite.config.js and
// vite.config.d.ts beside the source on every pnpm build, and they ended up in
// commits. Found building the WhatsApp blueprint's desktop app.
const desktopTSConfigNodeOutDir = "    \"outDir\": \"node_modules/.tmp/tsconfig-node\",\n"

func repairDesktopTSConfigNodeSource(src string) (string, []string, []string) {
	const anchor = "    \"composite\": true,\n"
	if strings.Contains(src, "\"outDir\"") || strings.Count(src, anchor) != 1 {
		return src, nil, nil
	}
	return strings.Replace(src, anchor, anchor+desktopTSConfigNodeOutDir, 1),
		[]string{"tsc -b writes vite.config's output to node_modules, not beside the source (delete a stray vite.config.js and .d.ts)"}, nil
}

// repairDesktopTSConfigNode applies it to an existing desktop app.
func repairDesktopTSConfigNode(root string) error {
	path := filepath.Join(root, "apps", "desktop", "frontend", "tsconfig.node.json")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairTextFile(root, m, path, repairDesktopTSConfigNodeSource)
}
