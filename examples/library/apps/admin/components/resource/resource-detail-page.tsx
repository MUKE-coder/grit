"use client";

import { Suspense, useState } from "react";
import Link from "next/link";
import dynamic from "next/dynamic";
import { useRouter } from "next/navigation";
import type { ResourceDefinition, ColumnDefinition, FieldDefinition } from "@/lib/resource";
import { useResource } from "@/hooks/use-resource";
import { useResourceDetailController } from "@/hooks/use-resource-detail-controller";
import { renderCell } from "@/components/tables/cell-renderers";
import { DataTable } from "@/components/tables/data-table";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { ArrowLeft, Pencil, Trash2, Loader2, Printer, Plus, FileText } from "@/lib/icons";
import { useLocalizedResource } from "@/lib/i18n";
import { buttonClasses } from "@/components/ui/button";

// The form stack loads when a record is edited or a related one created, the way
// the list page loads it. Suspense covers the Vite admin, where dynamic() is
// React.lazy.
const FormSheet = dynamic(() =>
  import("@/components/forms/form-sheet").then((m) => m.FormSheet)
);

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
export function ResourceDetailPage({ resource: definition, id }: ResourceDetailPageProps) {
  const resource = useLocalizedResource(definition);
  const CustomPage = resource.components?.DetailPage;
  if (CustomPage) return <CustomPage resource={resource} id={id} />;
  return <ResourceDetailView resource={resource} id={id} />;
}

// The default detail view. Every piece of state it uses comes from
// useResourceDetailController, so this component is markup and nothing else.
// That is the proof the hook is complete enough to build your own on.
function ResourceDetailView({ resource, id }: ResourceDetailPageProps) {
  const c = useResourceDetailController(resource, id);

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
    // space-y-6 rather than a margin on each block. A DetailAside or
    // DetailFields component is somebody else's, and it should not have to know
    // this page's rhythm to sit correctly in it.
    <div id="print-area" className="space-y-6">
      {/* Header */}
      {CustomHeader ? (
        <CustomHeader resource={resource} id={id} controller={c} />
      ) : (
      <div className="flex flex-wrap items-start justify-between gap-4">
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
          {c.can("edit") && (
          <button
            onClick={() => setEditing(true)}
            className={buttonClasses()}
          >
            <Pencil className="h-4 w-4" /> Edit
          </button>
          )}
          {c.can("delete") && (
          <button
            disabled={isDeleting}
            onClick={() => setConfirmDelete(true)}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-2 text-sm text-text-secondary hover:border-danger/40 hover:text-danger disabled:opacity-50 transition-colors"
          >
            <Trash2 className="h-4 w-4" /> Delete
          </button>
          )}
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
      {related.length > 0 && (
        <div className="no-print space-y-6">
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
      )}

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

      {editing && (
        <Suspense fallback={null}>
          <FormSheet resource={resource} item={record} onClose={() => setEditing(false)} />
        </Suspense>
      )}
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
    <div className="rounded-xl border border-border bg-bg-elevated">
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
        <Suspense fallback={null}>
          <FormSheet
            resource={createResource}
            item={null}
            defaults={{ [fk]: parentId }}
            onClose={() => setCreating(false)}
          />
        </Suspense>
      )}
    </div>
  );
}
