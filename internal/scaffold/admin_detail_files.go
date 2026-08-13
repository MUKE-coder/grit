package scaffold

import "strings"

// adminResourceDetailPage returns the generic record detail page. Every "view"
// action navigates here instead of opening a modal: it presents the record's
// fields, lets you edit in place (reusing the FormSheet), and — the important
// part — fetches every RELATED table (a resource's own inline line-items, plus
// any registry resource that belongs_to this one) so an Invoice shows its items
// and payments without a hand-written page.
// adminUseResourceDetailController emits hooks/use-resource-detail-controller.ts.
func adminUseResourceDetailController() string {
	return `"use client";

import { useCallback, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import type {
  RelatedResource,
  ResourceDefinition,
  ResourceDetailController,
} from "@/lib/resource";
import { useResourceItem, useDeleteResource } from "@/hooks/use-resource";
import { apiClient } from "@/lib/api-client";
import { resources } from "@/resources";

/*
 * Everything a resource DETAIL page needs except the markup.
 *
 * The list page got this treatment first, and the reasoning is the same: the
 * data was never the hard part, the rest of the page was. Loading one record
 * is a hook call. Working out which other resources point at this one, pulling
 * the inline line-item fields out of the form definition, fetching a
 * server-rendered PDF through the auth interceptor rather than a bare link,
 * and routing back to the list after a delete are not.
 *
 * const c = useResourceDetailController(resource, id)
 * <MyDetail record={c.record} onEdit={c.edit} sections={c.related} />
 */

/** The first present human-readable field, else the resource's own label. */
function titleOf(resource: ResourceDefinition, record: Record<string, unknown> | undefined): string {
  if (!record) return resource.label?.singular ?? resource.name;
  for (const key of ["number", "title", "name", "label", "reference", "slug", "email"]) {
    const value = record[key];
    if (typeof value === "string" && value) return value;
  }
  return resource.label?.singular ?? resource.name;
}

export function useResourceDetailController<T = Record<string, unknown>>(
  resource: ResourceDefinition,
  id: string,
): ResourceDetailController<T> {
  const router = useRouter();
  const { data, isLoading } = useResourceItem<T>(resource.endpoint, id);
  const record = data?.data;

  const [editing, setEditing] = useState(false);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [isPdfBusy, setPdfBusy] = useState(false);

  const singular = resource.label?.singular ?? resource.name;
  const { mutate: deleteItem, isPending: isDeleting } = useDeleteResource(resource.endpoint, singular);

  // Through apiClient, not a bare <a href>: that way the auth cookies, the
  // CSRF header and the 401-refresh interceptor all apply. A link gets none
  // of them and silently downloads a login page.
  const downloadPdf = useCallback(async () => {
    setPdfBusy(true);
    try {
      const res = await apiClient.get(resource.endpoint + "/" + id + "/pdf", {
        responseType: "blob",
      });
      const url = URL.createObjectURL(new Blob([res.data], { type: "application/pdf" }));
      window.open(url, "_blank", "noopener,noreferrer");
      // Revoked late: revoking straight away can race the new tab's load.
      setTimeout(() => URL.revokeObjectURL(url), 60000);
    } finally {
      setPdfBusy(false);
    }
  }, [resource.endpoint, id]);

  const lineItemFields = useMemo(
    () => (resource.form?.fields ?? []).filter((f) => f.type === "line-items"),
    [resource.form?.fields],
  );

  // Other resources in the registry that belong to this one, found by a
  // relationship-select whose relatedEndpoint is this endpoint. Anything
  // already rendered inline as line items is skipped, so it is not shown twice.
  const related = useMemo<RelatedResource[]>(() => {
    const inlineEndpoints = new Set(
      (resource.form?.fields ?? [])
        .filter((f) => f.type === "line-items" && f.itemEndpoint)
        .map((f) => f.itemEndpoint as string),
    );
    const out: RelatedResource[] = [];
    for (const other of resources) {
      if (other.slug === resource.slug || other.hidden || inlineEndpoints.has(other.endpoint)) continue;
      const fkField = (other.form?.fields ?? []).find(
        (f) => f.type === "relationship-select" && f.relatedEndpoint === resource.endpoint,
      );
      if (fkField) out.push({ resource: other, fk: fkField.key });
    }
    return out;
  }, [resource]);

  const columns = useMemo(
    () => resource.table.columns.filter((col) => !col.hidden),
    [resource.table.columns],
  );

  const back = useCallback(() => {
    router.push("/resources/" + resource.slug);
  }, [router, resource.slug]);

  return {
    resource,
    id,

    record,
    isLoading,
    notFound: !isLoading && !record,
    title: titleOf(resource, record as Record<string, unknown> | undefined),

    columns,
    lineItemFields,
    related,

    edit: () => setEditing(true),
    remove: () => setConfirmOpen(true),
    print: () => window.print(),
    downloadPdf,
    back,
    isPdfBusy,
    isDeleting,

    form: { open: editing, item: (record ?? null) as T | null, close: () => setEditing(false) },
    confirmDelete: {
      open: confirmOpen,
      confirm: () => {
        setConfirmOpen(false);
        // Back to the list: staying on the detail page of a record that no
        // longer exists shows "could not be found", which reads as an error
        // rather than as the delete having worked.
        deleteItem(id, { onSuccess: back });
      },
      cancel: () => setConfirmOpen(false),
    },
  };
}
`
}

func adminResourceDetailPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import type { ResourceDefinition, ColumnDefinition, FieldDefinition } from "@/lib/resource";
import { useResource } from "@/hooks/use-resource";
import { useResourceDetailController } from "@/hooks/use-resource-detail-controller";
import { renderCell } from "@/components/tables/cell-renderers";
import { DataTable } from "@/components/tables/data-table";
import { FormSheet } from "@/components/forms/form-sheet";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { apiClient } from "@/lib/api-client";
import { ArrowLeft, Pencil, Trash2, Loader2, Printer, Plus, FileText } from "@/lib/icons";
import { buttonClasses } from "@/components/ui/button";

interface ResourceDetailPageProps {
  resource: ResourceDefinition;
  id: string;
}

// Convert a line-items field's itemFields into table columns for the detail
// view (read-only). renderCell handles the value formatting.
function itemColumns(itemFields: FieldDefinition[]): ColumnDefinition[] {
  return itemFields.map((f) => ({ key: f.key, label: f.label }));
}

// v3.145.0: a thin router, the same split ResourcePage has. A DetailPage slot
// replaces everything below, so it is checked first and unconditionally:
// somebody who has supplied a whole page owns its header and its dialogs too.
export function ResourceDetailPage({ resource, id }: ResourceDetailPageProps) {
  const CustomPage = resource.components?.DetailPage;
  if (CustomPage) return <CustomPage resource={resource} id={id} />;
  return <ResourceDetailView resource={resource} id={id} />;
}

// The default detail view. Every piece of state it uses comes from
// useResourceDetailController, so this component is markup and nothing else.
// That is the proof the hook is complete enough to build your own on.
function ResourceDetailView({ resource, id }: ResourceDetailPageProps) {
  const c = useResourceDetailController(resource, id);
  const router = useRouter();

  const CustomHeader = resource.components?.DetailHeader;
  const CustomFields = resource.components?.DetailFields;
  const CustomAside = resource.components?.DetailAside;

  const record = c.record;
  const isLoading = c.isLoading;
  const lineItemFields = c.lineItemFields;
  const related = c.related;
  const pdfBusy = c.isPdfBusy;
  const isDeleting = c.isDeleting;
  const downloadPdf = c.downloadPdf;
  const setEditing = (open: boolean) => (open ? c.edit() : c.form.close());
  const setConfirmDelete = (open: boolean) => (open ? c.remove() : c.confirmDelete.cancel());
  const editing = c.form.open;
  const confirmDelete = c.confirmDelete.open;

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

  const cols = c.columns;

  return (
    <div id="print-area">
      {/* Header */}
      {CustomHeader ? (
        <CustomHeader resource={resource} id={id} controller={c} />
      ) : (
      <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
        <div>
          <Link
            href={"/resources/" + resource.slug}
            className="no-print mb-2 inline-flex items-center gap-1 text-xs text-text-muted hover:text-foreground"
          >
            <ArrowLeft className="h-3.5 w-3.5" /> Back to {resource.label?.plural ?? resource.name}
          </Link>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">{c.title}</h1>
          <p className="text-sm text-text-muted">{resource.label?.singular ?? resource.name} details</p>
        </div>
        <div className="no-print flex items-center gap-2">
          {/* The PDF is rendered server-side (GET <endpoint>/:id/pdf) so it
              looks the same everywhere and can be emailed or archived —
              unlike the browser's print dialog, which only reproduces the
              page. Fetched through apiClient rather than opened as a bare
              link: auth rides HttpOnly cookies, and a cross-origin top-level
              navigation (admin :3001 → api :8080) would not reliably carry
              them. The blob is opened in a new tab for viewing/saving. */}
          <button
            onClick={downloadPdf}
            disabled={pdfBusy}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm text-text-secondary hover:border-accent/40 hover:text-foreground transition-colors disabled:opacity-50"
          >
            {pdfBusy ? <Loader2 className="h-4 w-4 animate-spin" /> : <FileText className="h-4 w-4" />}
            PDF
          </button>
          <button
            onClick={() => window.print()}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm text-text-secondary hover:border-accent/40 hover:text-foreground transition-colors"
          >
            <Printer className="h-4 w-4" /> Print
          </button>
          <button
            onClick={() => setEditing(true)}
            className={buttonClasses()}
          >
            <Pencil className="h-4 w-4" /> Edit
          </button>
          <button
            disabled={isDeleting}
            onClick={() => setConfirmDelete(true)}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm text-text-secondary hover:border-danger/40 hover:text-danger disabled:opacity-50 transition-colors"
          >
            <Trash2 className="h-4 w-4" /> Delete
          </button>
        </div>
      </div>

      )}

      {/* Details */}
      {CustomFields ? (
        <CustomFields resource={resource} id={id} controller={c} />
      ) : (
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
      )}

      {/* Anything of your own that belongs between the fields and the related
          records: a status timeline, an activity feed, a map. */}
      {CustomAside && <CustomAside resource={resource} id={id} controller={c} />}

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
          createResource={r}
        />
      ))}
      </div>

      <ConfirmModal
        open={confirmDelete}
        title={"Delete this " + (resource.label?.singular ?? resource.name) + "?"}
        description="This cannot be undone."
        confirmLabel="Delete"
        variant="danger"
        loading={isDeleting}
        onCancel={c.confirmDelete.cancel}
        onConfirm={c.confirmDelete.confirm}
      />

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
  createResource,
}: {
  title: string;
  endpoint: string;
  fk: string;
  parentId: string;
  columns: ColumnDefinition[];
  slug?: string;
  // When set, the table gets a "New <child>" button that opens the child's
  // create form pre-scoped to this parent (its belongs_to FK is pre-filled).
  createResource?: ResourceDefinition;
}) {
  const router = useRouter();
  const [creating, setCreating] = useState(false);
  const { data, isLoading } = useResource<Record<string, unknown>>(endpoint, {
    filters: { [fk]: parentId },
    pageSize: 100,
  });
  const rows = data?.data ?? [];
  const childLabel = createResource?.label?.singular ?? createResource?.name ?? "item";

  return (
    <div className="mt-6 rounded-xl border border-border bg-bg-elevated">
      <div className="flex items-center justify-between gap-2 border-b border-border px-6 py-4">
        <div className="flex items-center gap-2">
          <h2 className="text-sm font-semibold text-foreground">{title}</h2>
          <span className="rounded-full bg-bg-hover px-2 py-0.5 text-xs text-text-muted">{rows.length}</span>
        </div>
        {createResource && (
          <button
            onClick={() => setCreating(true)}
            className={buttonClasses({ size: "sm", className: "no-print" })}
          >
            <Plus className="h-3.5 w-3.5" /> New {childLabel}
          </button>
        )}
      </div>
      <div className="p-2">
        <DataTable
          columns={columns}
          data={rows}
          isLoading={isLoading}
          onView={slug ? (item) => router.push("/resources/" + slug + "/" + String(item.id)) : undefined}
        />
      </div>
      {creating && createResource && (
        <FormSheet
          resource={createResource}
          item={null}
          defaults={{ [fk]: parentId }}
          onClose={() => setCreating(false)}
        />
      )}
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
import { {{CAMEL}}Resource } from "@/resources/{{KEBAB}}/{{KEBAB}}";

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
import { {{CAMEL}}Resource } from '@/resources/{{KEBAB}}/{{KEBAB}}'

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
