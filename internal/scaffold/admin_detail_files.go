package scaffold

import "strings"

// adminResourceDetailPage returns the generic record detail page. Every "view"
// action navigates here instead of opening a modal: it presents the record's
// fields, lets you edit in place (reusing the FormSheet), and — the important
// part — fetches every RELATED table (a resource's own inline line-items, plus
// any registry resource that belongs_to this one) so an Invoice shows its items
// and payments without a hand-written page.
func adminResourceDetailPage() string {
	return `"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import type { ResourceDefinition, ColumnDefinition, FieldDefinition } from "@/lib/resource";
import { resources } from "@/resources";
import { useResourceItem, useResource, useDeleteResource } from "@/hooks/use-resource";
import { renderCell } from "@/components/tables/cell-renderers";
import { DataTable } from "@/components/tables/data-table";
import { FormSheet } from "@/components/forms/form-sheet";
import { ArrowLeft, Pencil, Trash2, Loader2, Printer } from "@/lib/icons";

interface ResourceDetailPageProps {
  resource: ResourceDefinition;
  id: string;
}

// A readable title for the record — the first present human field, else the
// resource's singular label.
function titleOf(resource: ResourceDefinition, record: Record<string, unknown>): string {
  for (const k of ["number", "title", "name", "label", "reference", "slug", "email"]) {
    const v = record[k];
    if (typeof v === "string" && v) return v;
  }
  return resource.label?.singular ?? resource.name;
}

// Convert a line-items field's itemFields into table columns for the detail
// view (read-only). renderCell handles the value formatting.
function itemColumns(itemFields: FieldDefinition[]): ColumnDefinition[] {
  return itemFields.map((f) => ({ key: f.key, label: f.label }));
}

export function ResourceDetailPage({ resource, id }: ResourceDetailPageProps) {
  const router = useRouter();
  const { data, isLoading } = useResourceItem<Record<string, unknown>>(resource.endpoint, id);
  const record = data?.data;
  const [editing, setEditing] = useState(false);
  const { mutate: deleteItem, isPending: isDeleting } = useDeleteResource(resource.endpoint);

  // Inline line-items declared on THIS resource (rendered from the field's own
  // itemEndpoint + itemFields).
  const lineItemFields = (resource.form?.fields ?? []).filter((f) => f.type === "line-items");

  // Other registry resources that belongs_to this one — discovered by a
  // relationship-select field whose relatedEndpoint matches this endpoint.
  const related = useMemo(() => {
    const out: { resource: ResourceDefinition; fk: string }[] = [];
    // Endpoints already rendered as inline line-items — don't show them twice.
    const inlineEndpoints = new Set(
      (resource.form?.fields ?? [])
        .filter((f) => f.type === "line-items" && f.itemEndpoint)
        .map((f) => f.itemEndpoint as string)
    );
    for (const r of resources) {
      if (r.slug === resource.slug || r.hidden || inlineEndpoints.has(r.endpoint)) continue;
      const fkField = (r.form?.fields ?? []).find(
        (f) => f.type === "relationship-select" && f.relatedEndpoint === resource.endpoint
      );
      if (fkField) out.push({ resource: r, fk: fkField.key });
    }
    return out;
  }, [resource]);

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 p-8 text-sm text-text-muted">
        <Loader2 className="h-4 w-4 animate-spin" /> Loading…
      </div>
    );
  }
  if (!record) {
    return (
      <div className="p-8">
        <Link href={"/resources/" + resource.slug} className="text-sm text-accent hover:underline">
          ← Back to {resource.label?.plural ?? resource.name}
        </Link>
        <p className="mt-4 text-sm text-text-muted">This record could not be found.</p>
      </div>
    );
  }

  const cols = resource.table.columns.filter((c) => !c.hidden);

  return (
    <div id="print-area">
      {/* Header */}
      <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
        <div>
          <Link
            href={"/resources/" + resource.slug}
            className="no-print mb-2 inline-flex items-center gap-1 text-xs text-text-muted hover:text-foreground"
          >
            <ArrowLeft className="h-3.5 w-3.5" /> Back to {resource.label?.plural ?? resource.name}
          </Link>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">{titleOf(resource, record)}</h1>
          <p className="text-sm text-text-muted">{resource.label?.singular ?? resource.name} details</p>
        </div>
        <div className="no-print flex items-center gap-2">
          <button
            onClick={() => window.print()}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm text-text-secondary hover:border-accent/40 hover:text-foreground transition-colors"
          >
            <Printer className="h-4 w-4" /> Print
          </button>
          <button
            onClick={() => setEditing(true)}
            className="inline-flex items-center gap-1.5 rounded-lg bg-accent px-3 py-2 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
          >
            <Pencil className="h-4 w-4" /> Edit
          </button>
          <button
            disabled={isDeleting}
            onClick={() => {
              if (window.confirm("Delete this " + (resource.label?.singular ?? resource.name) + "?")) {
                deleteItem(id, { onSuccess: () => router.push("/resources/" + resource.slug) });
              }
            }}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm text-text-secondary hover:border-danger/40 hover:text-danger disabled:opacity-50 transition-colors"
          >
            <Trash2 className="h-4 w-4" /> Delete
          </button>
        </div>
      </div>

      {/* Details */}
      <div className="rounded-xl border border-border bg-bg-elevated p-6">
        <h2 className="mb-5 text-sm font-semibold text-foreground">Details</h2>
        <dl className="grid grid-cols-1 gap-x-8 gap-y-5 sm:grid-cols-2 lg:grid-cols-3">
          {cols.map((col) => (
            <div key={col.key} className="min-w-0">
              <dt className="text-xs font-medium uppercase tracking-wide text-text-muted">{col.label}</dt>
              <dd className="mt-1 break-words text-sm text-foreground">
                {renderCell(col, record[col.key], record)}
              </dd>
            </div>
          ))}
        </dl>
      </div>

      {/* Inline line-items on this resource */}
      {lineItemFields.map((f) =>
        f.itemEndpoint && f.foreignKey ? (
          <RelatedTable
            key={f.key}
            title={f.label}
            endpoint={f.itemEndpoint}
            fk={f.foreignKey}
            parentId={id}
            columns={itemColumns(f.itemFields ?? [])}
          />
        ) : null
      )}

      {/* Related registry resources — not part of the printed record */}
      <div className="no-print">
      {related.map(({ resource: r, fk }) => (
        <RelatedTable
          key={r.slug}
          title={r.label?.plural ?? r.name}
          endpoint={r.endpoint}
          fk={fk}
          parentId={id}
          columns={r.table.columns.filter((c) => !c.hidden)}
          slug={r.slug}
        />
      ))}
      </div>

      {editing && <FormSheet resource={resource} item={record} onClose={() => setEditing(false)} />}
    </div>
  );
}

function RelatedTable({
  title,
  endpoint,
  fk,
  parentId,
  columns,
  slug,
}: {
  title: string;
  endpoint: string;
  fk: string;
  parentId: string;
  columns: ColumnDefinition[];
  slug?: string;
}) {
  const router = useRouter();
  const { data, isLoading } = useResource<Record<string, unknown>>(endpoint, {
    filters: { [fk]: parentId },
    pageSize: 100,
  });
  const rows = data?.data ?? [];

  return (
    <div className="mt-6 rounded-xl border border-border bg-bg-elevated">
      <div className="flex items-center gap-2 border-b border-border px-6 py-4">
        <h2 className="text-sm font-semibold text-foreground">{title}</h2>
        <span className="rounded-full bg-bg-hover px-2 py-0.5 text-xs text-text-muted">{rows.length}</span>
      </div>
      <div className="p-2">
        <DataTable
          columns={columns}
          data={rows}
          isLoading={isLoading}
          onView={slug ? (item) => router.push("/resources/" + slug + "/" + String(item.id)) : undefined}
        />
      </div>
    </div>
  );
}
`
}

// adminResourceDetailRoute returns the thin per-resource [id] route wrapper that
// renders <ResourceDetailPage> for a given resource. Emitted for every resource
// so "view" always has a page to land on.
func adminResourceDetailRoute(camelName, pluralKebab, pascalName string) string {
	src := `"use client";

import { use } from "react";
import { ResourceDetailPage } from "@/components/resource/resource-detail-page";
import { {{CAMEL}}Resource } from "@/resources/{{KEBAB}}";

export default function {{PASCAL}}DetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  return <ResourceDetailPage resource={{{CAMEL}}Resource} id={id} />;
}
`
	src = strings.ReplaceAll(src, "{{CAMEL}}", camelName)
	src = strings.ReplaceAll(src, "{{KEBAB}}", pluralKebab)
	src = strings.ReplaceAll(src, "{{PASCAL}}", pascalName)
	return src
}

// adminResourceDetailRouteTanStack is the TanStack equivalent of the Next detail
// route wrapper — a $id route that reads the param and renders <ResourceDetailPage>.
func adminResourceDetailRouteTanStack(camelName, pluralKebab string) string {
	src := `import { createFileRoute } from '@tanstack/react-router'
import { ResourceDetailPage } from '@/components/resource/resource-detail-page'
import { {{CAMEL}}Resource } from '@/resources/{{KEBAB}}'

export const Route = createFileRoute('/_dashboard/resources/{{KEBAB}}/$id')({
  component: RouteComponent,
})

function RouteComponent() {
  const { id } = Route.useParams()
  return <ResourceDetailPage resource={{{CAMEL}}Resource} id={id} />
}
`
	src = strings.ReplaceAll(src, "{{CAMEL}}", camelName)
	src = strings.ReplaceAll(src, "{{KEBAB}}", pluralKebab)
	return src
}
