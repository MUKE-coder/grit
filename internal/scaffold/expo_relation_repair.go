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
