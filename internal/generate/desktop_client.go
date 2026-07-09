package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// desktopClientRoot returns apps/desktop (the monorepo Wails client). The
// generated screens live under its frontend/ tree.
func (g *Generator) desktopClientRoot() string {
	return filepath.Join(g.Root, "apps", "desktop")
}

// writeDesktopClientResourceFiles generates OFFLINE-FIRST CRUD screens for the
// monorepo Wails desktop client. Every read/write goes through the sync
// engine's Local* bindings (local SQLite mirror + outbox), so the screens work
// with no connection and the background sync loop reconciles with the server.
//
// It writes:
//
//	frontend/src/hooks/use-<plural>.ts                    — React Query over Local*
//	frontend/src/components/resource-forms/<plural>-form.tsx
//	frontend/src/routes/_app/<plural>.index.tsx           — list
//	frontend/src/routes/_app/<plural>.new.tsx
//	frontend/src/routes/_app/<plural>.$id.edit.tsx
//
// …and injects a sidebar nav entry.
func (g *Generator) writeDesktopClientResourceFiles(names Names) error {
	root := g.desktopClientRoot()

	files := map[string]string{
		filepath.Join(root, "frontend", "src", "hooks", "use-"+names.PluralKebab+".ts"):                       g.desktopClientHook(names),
		filepath.Join(root, "frontend", "src", "components", "resource-forms", names.PluralKebab+"-form.tsx"): g.desktopClientForm(names),
		filepath.Join(root, "frontend", "src", "routes", "_app", names.Plural+".index.tsx"):                   g.desktopClientListRoute(names),
		filepath.Join(root, "frontend", "src", "routes", "_app", names.Plural+".new.tsx"):                     g.desktopClientNewRoute(names),
		filepath.Join(root, "frontend", "src", "routes", "_app", names.Plural+".$id.edit.tsx"):                g.desktopClientEditRoute(names),
	}
	for path, content := range files {
		if err := writeFileWithDirs(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return g.injectDesktopClientNav(names)
}

// desktopClientHook wraps the Local* sync bindings in React Query hooks. Lists
// refetch on an interval so the background pull's mirror updates surface.
func (g *Generator) desktopClientHook(names Names) string {
	p := names.Plural
	Pascal := names.Pascal
	Plural := names.PluralPascal
	return `import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { localList, localGet, localCreate, localUpdate, localDelete } from "@/lib/sync-client";

export type ` + Pascal + ` = Record<string, unknown> & { id: string };

export function use` + Plural + `() {
  return useQuery<` + Pascal + `[]>({
    queryKey: ["` + p + `"],
    queryFn: async () => (await localList("` + p + `")) as ` + Pascal + `[],
    // The background sync loop keeps the local mirror fresh; re-read it so the
    // list reflects server changes without a manual refresh.
    refetchInterval: 3000,
  });
}

export function use` + Pascal + `(id: string) {
  return useQuery<` + Pascal + ` | null>({
    queryKey: ["` + p + `", id],
    queryFn: async () => (await localGet("` + p + `", id)) as ` + Pascal + ` | null,
    enabled: !!id,
  });
}

export function useCreate` + Pascal + `() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (data: Record<string, unknown>) => localCreate("` + p + `", "", data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["` + p + `"] }),
  });
}

export function useUpdate` + Pascal + `() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Record<string, unknown> }) =>
      localUpdate("` + p + `", id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["` + p + `"] }),
  });
}

export function useDelete` + Pascal + `() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => localDelete("` + p + `", id),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["` + p + `"] }),
  });
}
`
}

// desktopClientFormField describes how one field renders in the form.
type desktopClientFieldParts struct {
	imports  string
	state    string
	prefill  string
	jsx      string
	payload  string
	optionLd string // option-loading hooks (belongs_to)
}

// buildDesktopClientForm assembles the per-field pieces for the form component.
func (g *Generator) buildDesktopClientForm(names Names) desktopClientFieldParts {
	var parts desktopClientFieldParts
	var imp, st, pf, jsx, pay, opt strings.Builder
	seenRelImport := map[string]bool{}

	for _, f := range g.Definition.Fields {
		if f.IsSlug() || f.IsFile() || f.IsFiles() || f.IsManyToMany() || f.IsStringArray() {
			// slug is auto; file/array/m2m aren't rendered in the offline form.
			continue
		}
		camel := lowerCamel(f.Name)
		setter := "set" + toPascalCase(f.Name)
		json := toSnakeCase(f.Name)
		label := humanizeLabel(f.Name)
		ft := FieldType(f.Type)

		switch {
		case f.IsBelongsTo():
			rel := MakeNames(f.RelatedModelName())
			hook := "use" + rel.PluralPascal
			fk := f.FKColumnName()
			fkCamel := lowerCamel(fk)
			fkSetter := "set" + toPascalCase(fk)
			optsVar := lowerCamel(rel.Plural) + "Opts"
			if !seenRelImport[hook] {
				imp.WriteString("import { " + hook + " } from \"@/hooks/use-" + rel.PluralKebab + "\";\n")
				seenRelImport[hook] = true
			}
			st.WriteString("  const [" + fkCamel + ", " + fkSetter + "] = useState<string>(\"\");\n")
			opt.WriteString("  const " + optsVar + " = " + hook + "().data ?? [];\n")
			pf.WriteString("      " + fkSetter + "(String(record." + fk + " ?? \"\"));\n")
			jsx.WriteString("        <div>\n")
			jsx.WriteString("          <label className=\"block text-[13px] font-medium text-foreground mb-1.5\">" + label + "</label>\n")
			jsx.WriteString("          <select value={" + fkCamel + "} onChange={(e) => " + fkSetter + "(e.target.value)} className={inputCls}>\n")
			jsx.WriteString("            <option value=\"\">Select " + strings.ToLower(label) + "…</option>\n")
			jsx.WriteString("            {" + optsVar + ".map((o: any) => (<option key={o.id} value={o.id}>{o.name || o.title || o.id}</option>))}\n")
			jsx.WriteString("          </select>\n")
			jsx.WriteString("        </div>\n")
			pay.WriteString("      " + fk + ": " + fkCamel + ",\n")

		case ft == FieldInt || ft == FieldUint || ft == FieldFloat:
			step := ""
			if ft == FieldFloat {
				step = " step=\"any\""
			}
			st.WriteString("  const [" + camel + ", " + setter + "] = useState<number>(0);\n")
			pf.WriteString("      " + setter + "(Number(record." + json + " ?? 0));\n")
			jsx.WriteString("        <div>\n")
			jsx.WriteString("          <label className=\"block text-[13px] font-medium text-foreground mb-1.5\">" + label + "</label>\n")
			jsx.WriteString("          <input type=\"number\"" + step + " value={" + camel + "} onChange={(e) => " + setter + "(Number(e.target.value))} className={inputCls} />\n")
			jsx.WriteString("        </div>\n")
			pay.WriteString("      " + json + ": " + camel + ",\n")

		case ft == FieldBool:
			st.WriteString("  const [" + camel + ", " + setter + "] = useState<boolean>(false);\n")
			pf.WriteString("      " + setter + "(Boolean(record." + json + "));\n")
			jsx.WriteString("        <label className=\"flex items-center gap-2 text-[13px] text-foreground\">\n")
			jsx.WriteString("          <input type=\"checkbox\" checked={" + camel + "} onChange={(e) => " + setter + "(e.target.checked)} className=\"h-4 w-4\" />\n")
			jsx.WriteString("          " + label + "\n")
			jsx.WriteString("        </label>\n")
			pay.WriteString("      " + json + ": " + camel + ",\n")

		case ft == FieldDate || ft == FieldDatetime:
			inputType := "date"
			slice := "10"
			if ft == FieldDatetime {
				inputType = "datetime-local"
				slice = "16"
			}
			st.WriteString("  const [" + camel + ", " + setter + "] = useState<string>(\"\");\n")
			pf.WriteString("      " + setter + "(record." + json + " ? String(record." + json + ").slice(0, " + slice + ") : \"\");\n")
			jsx.WriteString("        <div>\n")
			jsx.WriteString("          <label className=\"block text-[13px] font-medium text-foreground mb-1.5\">" + label + "</label>\n")
			jsx.WriteString("          <input type=\"" + inputType + "\" value={" + camel + "} onChange={(e) => " + setter + "(e.target.value)} className={inputCls} />\n")
			jsx.WriteString("        </div>\n")
			pay.WriteString("      " + json + ": " + camel + ",\n")

		case ft == FieldText || ft == FieldRichtext:
			st.WriteString("  const [" + camel + ", " + setter + "] = useState<string>(\"\");\n")
			pf.WriteString("      " + setter + "(String(record." + json + " ?? \"\"));\n")
			jsx.WriteString("        <div>\n")
			jsx.WriteString("          <label className=\"block text-[13px] font-medium text-foreground mb-1.5\">" + label + "</label>\n")
			jsx.WriteString("          <textarea value={" + camel + "} onChange={(e) => " + setter + "(e.target.value)} rows={4} className={inputCls} />\n")
			jsx.WriteString("        </div>\n")
			pay.WriteString("      " + json + ": " + camel + ",\n")

		default: // string
			st.WriteString("  const [" + camel + ", " + setter + "] = useState<string>(\"\");\n")
			pf.WriteString("      " + setter + "(String(record." + json + " ?? \"\"));\n")
			jsx.WriteString("        <div>\n")
			jsx.WriteString("          <label className=\"block text-[13px] font-medium text-foreground mb-1.5\">" + label + "</label>\n")
			jsx.WriteString("          <input type=\"text\" value={" + camel + "} onChange={(e) => " + setter + "(e.target.value)} className={inputCls} />\n")
			jsx.WriteString("        </div>\n")
			pay.WriteString("      " + json + ": " + camel + ",\n")
		}
	}

	parts.imports = imp.String()
	parts.state = st.String()
	parts.prefill = pf.String()
	parts.jsx = jsx.String()
	parts.payload = pay.String()
	parts.optionLd = opt.String()
	return parts
}

// desktopClientForm is the shared create/edit form component. It takes an
// optional record (edit mode) and an onSubmit that receives the typed payload.
func (g *Generator) desktopClientForm(names Names) string {
	parts := g.buildDesktopClientForm(names)
	Pascal := names.Pascal

	return `import { useEffect, useState } from "react";
` + parts.imports + `
const inputCls =
  "w-full bg-background border border-border rounded-lg px-3 py-2 text-[13px] text-foreground focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors";

interface ` + Pascal + `FormProps {
  record?: Record<string, unknown> | null;
  submitting?: boolean;
  submitLabel: string;
  onSubmit: (values: Record<string, unknown>) => void | Promise<void>;
}

export function ` + Pascal + `Form({ record, submitting, submitLabel, onSubmit }: ` + Pascal + `FormProps) {
` + parts.state + parts.optionLd + `
  useEffect(() => {
    if (!record) return;
` + parts.prefill + `  }, [record]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit({
` + parts.payload + `    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4 max-w-2xl">
` + parts.jsx + `
      <button
        type="submit"
        disabled={submitting}
        className="rounded-lg bg-accent px-4 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
      >
        {submitting ? "Saving…" : submitLabel}
      </button>
    </form>
  );
}
`
}

// desktopClientListRoute is the table view. Reads the local mirror, supports a
// simple in-memory search, and links to create/edit; delete is local-first.
func (g *Generator) desktopClientListRoute(names Names) string {
	p := names.Plural
	Pascal := names.Pascal
	Plural := names.PluralPascal
	title := humanizeLabel(names.Plural)

	// Choose up to 3 scalar columns for the table.
	var cols []string
	for _, f := range g.Definition.Fields {
		if f.IsSlug() || f.IsFile() || f.IsFiles() || f.IsManyToMany() || f.IsStringArray() || f.IsBelongsTo() {
			continue
		}
		cols = append(cols, toSnakeCase(f.Name))
		if len(cols) >= 3 {
			break
		}
	}
	if len(cols) == 0 {
		cols = []string{"id"}
	}
	var headJSX, cellJSX strings.Builder
	for _, c := range cols {
		headJSX.WriteString("              <th className=\"text-left font-medium px-4 py-2\">" + humanizeLabel(c) + "</th>\n")
		cellJSX.WriteString("                <td className=\"px-4 py-2\">{String(item." + c + " ?? \"\")}</td>\n")
	}
	searchCol := cols[0]

	return `import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Plus, Pencil, Trash2 } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { Card, CardContent } from "@/components/ui/card";
import { use` + Plural + `, useDelete` + Pascal + ` } from "@/hooks/use-` + names.PluralKebab + `";

export const Route = createFileRoute("/_app/` + p + `/")({
  component: ` + Plural + `Page,
});

function ` + Plural + `Page() {
  const navigate = useNavigate();
  const { data: items = [], isLoading } = use` + Plural + `();
  const del = useDelete` + Pascal + `();
  const [search, setSearch] = useState("");

  const filtered = items.filter((i) =>
    !search || String(i.` + searchCol + ` ?? "").toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <div>
      <PageHeader title="` + title + `" description="Manage your ` + strings.ToLower(title) + `" />

      <div className="mt-6 flex items-center justify-between gap-3">
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search…"
          className="w-72 bg-background border border-border rounded-lg px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent"
        />
        <Link
          to="/app/` + p + `/new"
          className="flex items-center gap-2 rounded-lg bg-accent px-3 py-2 text-[13px] font-semibold text-white hover:bg-accent-hover transition-colors"
        >
          <Plus className="h-4 w-4" /> New
        </Link>
      </div>

      <Card className="mt-4">
        <CardContent className="p-0">
          <table className="w-full text-[13px]">
            <thead className="text-foreground-secondary border-b border-border-subtle">
              <tr>
` + headJSX.String() + `                <th className="px-4 py-2"></th>
              </tr>
            </thead>
            <tbody>
              {isLoading ? (
                <tr><td className="px-4 py-6 text-foreground-muted" colSpan={99}>Loading…</td></tr>
              ) : filtered.length === 0 ? (
                <tr><td className="px-4 py-6 text-foreground-muted" colSpan={99}>No ` + strings.ToLower(title) + ` yet.</td></tr>
              ) : (
                filtered.map((item) => (
                  <tr key={String(item.id)} className="border-b border-border-subtle last:border-0 hover:bg-surface-hover">
` + cellJSX.String() + `                    <td className="px-4 py-2 text-right whitespace-nowrap">
                      <button
                        onClick={() => navigate({ to: "/app/` + p + `/$id/edit", params: { id: String(item.id) } })}
                        className="p-1.5 rounded hover:bg-surface-2 text-foreground-secondary"
                        title="Edit"
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        onClick={() => {
                          if (confirm("Delete this ` + strings.ToLower(names.Lower) + `?")) del.mutate(String(item.id));
                        }}
                        className="p-1.5 rounded hover:bg-danger/10 text-danger"
                        title="Delete"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </CardContent>
      </Card>
    </div>
  );
}
`
}

// desktopClientNewRoute is the create screen.
func (g *Generator) desktopClientNewRoute(names Names) string {
	p := names.Plural
	Pascal := names.Pascal
	return `import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { PageHeader } from "@/components/layout/page-header";
import { ` + Pascal + `Form } from "@/components/resource-forms/` + names.PluralKebab + `-form";
import { useCreate` + Pascal + ` } from "@/hooks/use-` + names.PluralKebab + `";

export const Route = createFileRoute("/_app/` + p + `/new")({
  component: New` + Pascal + `Page,
});

function New` + Pascal + `Page() {
  const navigate = useNavigate();
  const create = useCreate` + Pascal + `();

  return (
    <div>
      <PageHeader title="New ` + names.Pascal + `" description="Create a new ` + strings.ToLower(names.Lower) + `" />
      <div className="mt-6">
        <` + Pascal + `Form
          submitting={create.isPending}
          submitLabel="Create ` + names.Pascal + `"
          onSubmit={async (values) => {
            await create.mutateAsync(values);
            navigate({ to: "/app/` + p + `" });
          }}
        />
      </div>
    </div>
  );
}
`
}

// desktopClientEditRoute is the edit screen.
func (g *Generator) desktopClientEditRoute(names Names) string {
	p := names.Plural
	Pascal := names.Pascal
	return `import { createFileRoute, useNavigate, useParams } from "@tanstack/react-router";
import { PageHeader } from "@/components/layout/page-header";
import { ` + Pascal + `Form } from "@/components/resource-forms/` + names.PluralKebab + `-form";
import { use` + Pascal + `, useUpdate` + Pascal + ` } from "@/hooks/use-` + names.PluralKebab + `";

export const Route = createFileRoute("/_app/` + p + `/$id/edit")({
  component: Edit` + Pascal + `Page,
});

function Edit` + Pascal + `Page() {
  const navigate = useNavigate();
  const { id } = useParams({ from: "/_app/` + p + `/$id/edit" });
  const { data: record, isLoading } = use` + Pascal + `(id);
  const update = useUpdate` + Pascal + `();

  if (isLoading) return <div className="text-[13px] text-foreground-muted">Loading…</div>;

  return (
    <div>
      <PageHeader title="Edit ` + names.Pascal + `" description="Update this ` + strings.ToLower(names.Lower) + `" />
      <div className="mt-6">
        <` + Pascal + `Form
          record={record}
          submitting={update.isPending}
          submitLabel="Save Changes"
          onSubmit={async (values) => {
            await update.mutateAsync({ id, data: values });
            navigate({ to: "/app/` + p + `" });
          }}
        />
      </div>
    </div>
  );
}
`
}

// injectDesktopClientNav adds a sidebar entry for the resource to nav-config.ts
// (the single source of truth the sidebar reads). Idempotent and scoped to the
// // grit:nav marker.
func (g *Generator) injectDesktopClientNav(names Names) error {
	path := filepath.Join(g.desktopClientRoot(), "frontend", "src", "lib", "nav-config.ts")
	if !fileExists(path) {
		return nil
	}
	label := humanizeLabel(names.Plural)
	entry := "      { to: \"/app/" + names.Plural + "\", label: \"" + label + "\", icon: Box },"
	return injectBefore(path, "// grit:nav", entry)
}
