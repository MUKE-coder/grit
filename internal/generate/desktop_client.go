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
//	frontend/src/routes/app/<plural>.index.tsx           — list
//	frontend/src/routes/app/<plural>.new.tsx
//	frontend/src/routes/app/<plural>.$id.edit.tsx
//
// …and injects a sidebar nav entry.
func (g *Generator) writeDesktopClientResourceFiles(names Names) error {
	root := g.desktopClientRoot()

	files := map[string]string{
		filepath.Join(root, "frontend", "src", "hooks", "use-"+names.PluralKebab+".ts"):                       g.desktopClientHook(names),
		filepath.Join(root, "frontend", "src", "components", "resource-forms", names.PluralKebab+"-form.tsx"): g.desktopClientForm(names),
		filepath.Join(root, "frontend", "src", "routes", "app", names.Plural+".index.tsx"):                   g.desktopClientListRoute(names),
		filepath.Join(root, "frontend", "src", "routes", "app", names.Plural+".new.tsx"):                     g.desktopClientNewRoute(names),
		filepath.Join(root, "frontend", "src", "routes", "app", names.Plural+".$id.edit.tsx"):                g.desktopClientEditRoute(names),
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
	seenFileImport := false

	for _, f := range g.Definition.Fields {
		if f.IsSlug() || f.IsManyToMany() || f.IsStringArray() {
			// slug is auto; array/m2m aren't rendered in the offline form.
			continue
		}
		camel := lowerCamel(f.Name)
		setter := "set" + toPascalCase(f.Name)
		json := toSnakeCase(f.Name)
		label := humanizeLabel(f.Name)
		ft := FieldType(f.Type)

		switch {
		case f.IsFile() || f.IsFiles():
			if !seenFileImport {
				imp.WriteString("import { FileDropzone } from \"@/components/file-dropzone\";\n")
				imp.WriteString("import type { FileRef } from \"@/lib/api-client\";\n")
				seenFileImport = true
			}
			accept := desktopFileAccept(f.FileAccepts)
			acceptAttr := ""
			if accept != "" {
				acceptAttr = " accept=\"" + accept + "\""
			}
			if f.IsFiles() {
				st.WriteString("  const [" + camel + ", " + setter + "] = useState<FileRef[]>([]);\n")
				pf.WriteString("      " + setter + "((record." + json + " as FileRef[]) ?? []);\n")
				jsx.WriteString("        <FileDropzone label=\"" + label + "\" value={" + camel + "} onChange={(v) => " + setter + "((v as FileRef[]) ?? [])} multiple" + acceptAttr + " />\n")
			} else {
				st.WriteString("  const [" + camel + ", " + setter + "] = useState<FileRef | null>(null);\n")
				pf.WriteString("      " + setter + "((record." + json + " as FileRef | null) ?? null);\n")
				jsx.WriteString("        <FileDropzone label=\"" + label + "\" value={" + camel + "} onChange={(v) => " + setter + "((v as FileRef) ?? null)}" + acceptAttr + " />\n")
			}
			pay.WriteString("      " + json + ": " + camel + ",\n")

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

// desktopFileAccept turns a file field's accept-aliases (image, video, pdf, …)
// into an <input accept="…"> string. "all"/empty means no restriction.
func desktopFileAccept(accepts []string) string {
	aliases := map[string]string{
		"image": "image/*",
		"video": "video/*",
		"audio": "audio/*",
		"pdf":   ".pdf",
		"doc":   ".doc,.docx",
		"docx":  ".doc,.docx",
		"xls":   ".xls,.xlsx",
		"xlsx":  ".xls,.xlsx",
		"csv":   ".csv",
		"zip":   ".zip",
	}
	var parts []string
	for _, a := range accepts {
		if a == "all" || a == "" {
			return ""
		}
		if v, ok := aliases[a]; ok {
			parts = append(parts, v)
		}
	}
	return strings.Join(parts, ",")
}

// desktopClientForm is the shared create/edit form component. It takes an
// optional record (edit mode) and an onSubmit that receives the typed payload.
func (g *Generator) desktopClientForm(names Names) string {
	parts := g.buildDesktopClientForm(names)
	Pascal := names.Pascal

	return `import { useEffect, useState } from "react";
import { Loader2 } from "lucide-react";
` + parts.imports + `
const inputCls =
  "w-full rounded-lg border border-border bg-surface-2 px-4 py-2.5 text-[13px] text-foreground placeholder:text-foreground-muted outline-none transition-colors focus:border-accent focus:ring-1 focus:ring-accent";

interface ` + Pascal + `FormProps {
  record?: Record<string, unknown> | null;
  submitting?: boolean;
  submitLabel: string;
  onSubmit: (values: Record<string, unknown>) => void | Promise<void>;
  onCancel?: () => void;
}

export function ` + Pascal + `Form({ record, submitting, submitLabel, onSubmit, onCancel }: ` + Pascal + `FormProps) {
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
    <form onSubmit={handleSubmit} className="flex h-full flex-col">
      <div className="flex-1 space-y-4">
` + parts.jsx + `      </div>

      <div className="mt-6 flex justify-end gap-3 border-t border-border pt-4">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="rounded-lg border border-border px-4 py-2 text-[13px] font-medium text-foreground-secondary transition-colors hover:bg-surface-hover"
          >
            Cancel
          </button>
        )}
        <button
          type="submit"
          disabled={submitting}
          className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-[13px] font-semibold text-white transition-colors hover:bg-accent-hover disabled:opacity-50"
        >
          {submitting && <Loader2 className="h-4 w-4 animate-spin" />}
          {submitting ? "Saving…" : submitLabel}
        </button>
      </div>
    </form>
  );
}
`
}

// desktopClientListRoute is the rich list view: a shared client-side DataTable
// (stat cards + search + date filter + column visibility + CSV export +
// sortable headers + pagination) over the offline mirror, with create/edit in
// a right slide-over drawer — mirroring the admin panel's DataTable/FormSheet.
func (g *Generator) desktopClientListRoute(names Names) string {
	p := names.Plural
	Pascal := names.Pascal
	Plural := names.PluralPascal
	kebab := names.PluralKebab
	title := humanizeLabel(names.Plural)
	lower := strings.ToLower(names.Lower)

	// Build table columns from the fields (skip auto/file/relation-array
	// fields the offline form doesn't render). belongs_to shows the FK value.
	var colLines []string
	var searchKeys []string
	for _, f := range g.Definition.Fields {
		// m2m / string arrays are too noisy for a table cell; everything else
		// (incl. slug and file/image fields) gets a column.
		if f.IsManyToMany() || f.IsStringArray() {
			continue
		}
		var key, label, format string
		switch {
		case f.IsBelongsTo():
			key = f.FKColumnName()
			label = humanizeLabel(strings.TrimSuffix(f.Name, "_id"))
			format = "text"
		case f.IsFile() || f.IsFiles():
			key = toSnakeCase(f.Name)
			label = humanizeLabel(f.Name)
			format = "image"
		case f.IsSlug():
			key = toSnakeCase(f.Name)
			label = humanizeLabel(f.Name)
			format = "text"
		default:
			key = toSnakeCase(f.Name)
			label = humanizeLabel(f.Name)
			format = f.ColumnFormat()
			switch FieldType(f.Type) {
			case FieldString, FieldText:
				searchKeys = append(searchKeys, key)
			}
		}
		colLines = append(colLines, "  { key: \""+key+"\", label: \""+label+"\", format: \""+format+"\" },")
	}
	colLines = append(colLines, "  { key: \"created_at\", label: \"Created\", format: \"relative\" },")
	if len(searchKeys) == 0 {
		// fall back so the search box still filters on something.
		searchKeys = []string{"id"}
	}
	quotedKeys := make([]string, 0, len(searchKeys))
	for _, k := range searchKeys {
		quotedKeys = append(quotedKeys, "\""+k+"\"")
	}

	return `import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { PageHeader } from "@/components/layout/page-header";
import { DataTable, type DataColumn } from "@/components/tables/data-table";
import { ResourceDrawer } from "@/components/resource-drawer";
import { ` + Pascal + `Form } from "@/components/resource-forms/` + kebab + `-form";
import {
  use` + Plural + `,
  useCreate` + Pascal + `,
  useUpdate` + Pascal + `,
  useDelete` + Pascal + `,
  type ` + Pascal + `,
} from "@/hooks/use-` + kebab + `";

export const Route = createFileRoute("/app/` + p + `/")({
  component: ` + Plural + `Page,
});

const COLUMNS: DataColumn[] = [
` + strings.Join(colLines, "\n") + `
];

const SEARCH_KEYS = [` + strings.Join(quotedKeys, ", ") + `];

function ` + Plural + `Page() {
  const { data: items = [], isLoading } = use` + Plural + `();
  const create = useCreate` + Pascal + `();
  const update = useUpdate` + Pascal + `();
  const del = useDelete` + Pascal + `();
  const [drawer, setDrawer] = useState<{ open: boolean; record: ` + Pascal + ` | null }>({
    open: false,
    record: null,
  });

  const closeDrawer = () => setDrawer({ open: false, record: null });

  return (
    <div>
      <PageHeader title="` + title + `" description="Manage your ` + strings.ToLower(title) + `" />

      <div className="mt-6">
        <DataTable<` + Pascal + `>
          title="` + title + `"
          singular="` + Pascal + `"
          columns={COLUMNS}
          rows={items}
          loading={isLoading}
          searchKeys={SEARCH_KEYS}
          onNew={() => setDrawer({ open: true, record: null })}
          onEdit={(row) => setDrawer({ open: true, record: row })}
          onDelete={(row) => {
            if (confirm("Delete this ` + lower + `?")) del.mutate(String(row.id));
          }}
          onBulkDelete={(rows) => {
            if (confirm("Delete " + rows.length + " ` + lower + `(s)?")) rows.forEach((r) => del.mutate(String(r.id)));
          }}
          onImport={(records) => records.forEach((rec) => create.mutate(rec))}
        />
      </div>

      <ResourceDrawer
        open={drawer.open}
        title={drawer.record ? "Edit ` + Pascal + `" : "New ` + Pascal + `"}
        description={drawer.record ? "Update this ` + lower + `" : "Create a new ` + lower + `"}
        onClose={closeDrawer}
      >
        <` + Pascal + `Form
          record={drawer.record}
          submitting={create.isPending || update.isPending}
          submitLabel={drawer.record ? "Save changes" : "Create ` + Pascal + `"}
          onCancel={closeDrawer}
          onSubmit={async (values) => {
            if (drawer.record) {
              await update.mutateAsync({ id: String(drawer.record.id), data: values });
            } else {
              await create.mutateAsync(values);
            }
            closeDrawer();
          }}
        />
      </ResourceDrawer>
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

export const Route = createFileRoute("/app/` + p + `/new")({
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

export const Route = createFileRoute("/app/` + p + `/$id/edit")({
  component: Edit` + Pascal + `Page,
});

function Edit` + Pascal + `Page() {
  const navigate = useNavigate();
  const { id } = useParams({ from: "/app/` + p + `/$id/edit" });
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
