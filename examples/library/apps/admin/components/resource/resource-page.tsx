"use client";

import { useState } from "react";
import { useSearchParams } from "next/navigation";
import dynamic from "next/dynamic";
import type { ResourceDefinition } from "@/lib/resource";
import { useResourceController } from "@/hooks/use-resource-controller";
import type { ResourceController } from "@/hooks/use-resource-controller";
import { PageHeader } from "@/components/chrome/PageHeader";
import { DataTable } from "@/components/tables/data-table";
// Lazy: only resources declaring tree: true ever render this, and an eager
// import would put the drag-and-drop tree in every admin page's bundle.
const ResourceTree = dynamic(() =>
  import("@/components/resource/resource-tree").then((m) => m.ResourceTree)
);
import { TableToolbar } from "@/components/tables/table-toolbar";
import { TablePagination } from "@/components/tables/table-pagination";
import { TableFilters } from "@/components/tables/table-filters";
import { TableTabs } from "@/components/tables/table-tabs";
import { BulkActionBar } from "@/components/tables/bulk-action-bar";
import { BulkEditModal } from "@/components/tables/bulk-edit-modal";
import { exportToFile } from "@/lib/excel-utils";
// grit:resource:imports
import { buttonClasses } from "@/components/ui/button";
import { useLocalizedResource, useT } from "@/lib/i18n";

// Lazy-load modal/form components — they are only shown conditionally and
// would otherwise inflate the initial page bundle for every admin resource.
const FormModal = dynamic(() =>
  import("@/components/forms/form-modal").then((m) => m.FormModal)
);
const FormSheet = dynamic(() =>
  import("@/components/forms/form-sheet").then((m) => m.FormSheet)
);
const FormPage = dynamic(() =>
  import("@/components/forms/form-page").then((m) => m.FormPage)
);
const UpdateGroups = dynamic(() =>
  import("@/components/forms/update-groups").then((m) => m.UpdateGroups)
);
const FormModalSteps = dynamic(() =>
  import("@/components/forms/form-modal-steps").then((m) => m.FormModalSteps)
);
const FormPageSteps = dynamic(() =>
  import("@/components/forms/form-page-steps").then((m) => m.FormPageSteps)
);
const ConfirmModal = dynamic(() =>
  import("@/components/ui/confirm-modal").then((m) => m.ConfirmModal)
);
// v3.31.35 — Excel import modal, lazy-loaded so the xlsx parser
// only joins the bundle when the user actually clicks "Import".
const ImportModal = dynamic(() =>
  import("@/components/tables/import-modal").then((m) => m.ImportModal)
);

// NewLabel is the create button's text. Its own component, so the translated
// "New" costs the list view no hook of its own.
function NewLabel({ name }: { name: string }) {
  const t = useT();
  return <>{t("table.new", "New {name}", { name })}</>;
}

interface ResourcePageProps {
  resource: ResourceDefinition;
}

// v3.31.27: ResourcePage is a thin router. It picks between four possible
// views (UpdateGroups, FormPageSteps, FormPage, ResourceListView) based on
// formView + the ?action param. Before this split, the list-mode hooks all
// sat below the form-mode early returns — meaning the hook count varied
// between renders, which React 19 strict mode errors on. Splitting into two
// components keeps each function\'s hook list stable.
export function ResourcePage({ resource: definition }: ResourcePageProps) {
  // Labels translated once, here, so the table, forms and headings below read
  // them without knowing about translation.
  const resource = useLocalizedResource(definition);
  const searchParams = useSearchParams();

  // A Page slot replaces this entire component. It is checked first and
  // unconditionally: someone who has supplied a whole page owns the routing
  // inside it too, including whatever it wants to do with ?action=create.
  const CustomPage = resource.components?.Page;
  const isFormPage = resource.formView === "page" || resource.formView === "page-steps";
  const isSteps = resource.formView === "modal-steps" || resource.formView === "page-steps";
  const formAction = searchParams.get("action");

  // v3.31.18: editing + form has groups → render per-group cards with
  // PATCH-per-group saves. Falls back to the standard FormPage when no
  // groups are defined.
  const editId = searchParams.get("edit");
  const hasUpdateGroups = (resource.form.groups ?? []).some(
    (g) => !g.scope || g.scope === "update" || g.scope === "both"
  );

  if (CustomPage) {
    return <CustomPage resource={resource} />;
  }
  if (isFormPage && formAction === "edit" && editId && hasUpdateGroups) {
    return <UpdateGroups resource={resource} id={editId} />;
  }

  // If formView is "page" or "page-steps" and we have an action param, show the form page
  if (isFormPage && (formAction === "create" || formAction === "edit")) {
    return isSteps ? <FormPageSteps resource={resource} /> : <FormPage resource={resource} />;
  }

  return <ResourceListView resource={resource} />;
}

// Exports the ticked rows rather than the page or the whole table. The rows
// are already in memory, so this needs no request: the point of "export
// selection" is the selection.
function exportSelection(c: ResourceController) {
  if (c.selectedRows.length === 0) return;
  void exportToFile(c.selectedRows, c.columns, c.resource.slug, "csv");
  c.announce(c.selectedRows.length + " rows exported.");
}

// The default list view. Every piece of state and behaviour it uses comes from
// useResourceController — this component is markup and nothing else. That is
// deliberate: it is the proof that the hook is complete enough for someone to
// build their own page on, because the stock page is built on it too.
//
// Porting a bought template? Copy this file, keep the useResourceController
// line, and replace the JSX.
function ResourceListView({ resource }: ResourcePageProps) {
  const c = useResourceController(resource);
  // Tree resources open on the tree, because somebody who asked for a
  // hierarchy is looking for the hierarchy. The table is one click away and
  // keeps every filter, tab and bulk action it had.
  const [view, setView] = useState<"tree" | "table">(resource.tree ? "tree" : "table");

  // Slots, each falling back to the stock component. The props handed to a
  // custom Table are exactly DataTable's, so a replacement can also wrap the
  // original: (props) => <Card><DataTable {...props} /></Card>.
  const Table = resource.components?.Table ?? DataTable;
  const CustomForm = resource.components?.Form;
  const EmptyState = resource.components?.EmptyState;
  const CustomBulkBar = resource.components?.BulkBar;
  const showEmptyState = Boolean(EmptyState) && !c.isLoading && c.rows.length === 0;

  // Archive is a view, not a filter chip: the rows in it cannot be edited the
  // same way and the actions on them differ, so it gets its own tab. Shown
  // only when the resource actually has somewhere to archive to.
  const hasArchive =
    (resource.table.bulkActions ?? []).includes("archive") ||
    (resource.table.bulkActions ?? []).includes("restore");

  const headerActions = c.can("create") ? (
    <button onClick={c.create} className={buttonClasses({ size: "sm" })}>
      <span className="text-base leading-none">+</span>
      <NewLabel name={c.singularName} />
    </button>
  ) : undefined;

  return (
    <div>
      <PageHeader
        title={c.pluralName}
        subtitle={`Manage ${c.pluralName.toLowerCase()}`}
        actions={headerActions}
        stats={c.stats}
      />

      {/* Bulk actions change the table without moving focus, so every one of
          them is spoken here. */}
      <p role="status" aria-live="polite" className="sr-only">
        {c.liveMessage}
      </p>

      {hasArchive && (
        <div className="mb-3 flex w-fit gap-1 rounded-lg border border-border bg-bg-secondary p-1">
          {[
            { label: "Published", archived: false },
            { label: "Archived", archived: true },
          ].map((tab) => (
            <button
              key={tab.label}
              type="button"
              onClick={() => c.setShowArchived(tab.archived)}
              aria-pressed={c.showArchived === tab.archived}
              className={
                "min-h-9 rounded-md px-3 text-sm transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent " +
                (c.showArchived === tab.archived
                  ? "bg-accent/15 font-medium text-accent"
                  : "text-text-secondary hover:bg-bg-hover hover:text-foreground")
              }
            >
              {tab.label}
            </button>
          ))}
        </div>
      )}

      <div className="rounded-xl border border-border bg-bg-secondary">
        <TableTabs
          tabs={c.tabs}
          active={c.activeTab}
          onChange={c.setActiveTab}
          endpoint={resource.endpoint}
          baseFilters={c.filters}
        />

        <TableToolbar
          resource={resource}
          search={c.search}
          onSearch={c.setSearch}
          selectedCount={c.selection.length}
          onBulkDelete={c.bulkRemove}
          onCreate={c.can("create") ? c.create : undefined}
          allColumns={c.allColumns}
          hiddenColumns={c.hiddenColumns}
          onToggleColumn={c.toggleColumn}
          data={c.rows}
          dateRange={c.dateRange}
          onDateRangeChange={c.setDateRange}
          apiSearchParams={c.apiSearchParams}
          onImport={resource.table.import !== false ? () => c.importer.setOpen(true) : undefined}
        />

        {/* grit:table:toolbar */}

        {resource.table.filters && resource.table.filters.length > 0 && (
          <TableFilters
            filters={resource.table.filters}
            values={c.filters}
            onChange={c.setFilter}
          />
        )}

        {resource.tree && (
          <div className="mb-3 inline-flex rounded-lg border border-border p-0.5">
            {(["tree", "table"] as const).map((option) => (
              <button
                key={option}
                type="button"
                onClick={() => setView(option)}
                aria-pressed={view === option}
                className={
                  "rounded-md px-3 py-1 text-xs font-medium capitalize transition-colors " +
                  (view === option
                    ? "bg-accent text-accent-fg"
                    : "text-text-muted hover:text-foreground")
                }
              >
                {option}
              </button>
            ))}
          </div>
        )}

        {/* The tabs control this region, so it is their panel. Without the
            pairing a reader hears a tablist and never learns what it filters. */}
        {/* biome-ignore lint/a11y/useAriaPropsSupportedByRole: aria-labelledby is set only while the role is tabpanel. */}
        <div
          id="table-panel"
          role={c.tabs.length > 0 ? "tabpanel" : undefined}
          aria-labelledby={c.activeTab ? "table-tab-" + c.activeTab : undefined}
        >
        {resource.tree && view === "tree" ? (
          <ResourceTree
            resource={resource}
            // The page's own edit handler, so a node edited from the tree opens
            // the form that is already here rather than a second one that
            // drifts from it. It takes the row, and a tree node is the row.
            onEdit={
              c.can("edit")
                ? (node) => c.edit(node as unknown as Record<string, unknown>)
                : undefined
            }
            // createWith rather than create, so the row you clicked is already
            // chosen as the parent when the form opens.
            onAddChild={
              c.can("create")
                ? (parentID) => c.createWith({ parent_id: parentID })
                : undefined
            }
          />
        ) : showEmptyState && EmptyState ? (
          <EmptyState resource={resource} />
        ) : (
          <Table
            columns={c.columns}
            data={c.rows}
            isLoading={c.isLoading}
            isFetching={c.isFetching}
            sortBy={c.sortBy}
            sortOrder={c.sortOrder}
            onSort={c.setSort}
            selectedRows={c.selection}
            onSelectRows={c.setSelection}
            onView={c.can("view") ? c.view : undefined}
            onEdit={c.can("edit") ? c.edit : undefined}
            onDelete={c.can("delete") ? c.remove : undefined}
            rowActions={resource.table.rowActions}
          />
        )}
        </div>

        {/* Room under the table for the floating pill, only while it is there.
            Without it the last row sits behind the bar with nowhere to scroll,
            which is the objection that kept the bar in the flow to begin
            with. */}
        {c.selection.length > 0 && <div aria-hidden="true" className="h-20" />}

        {CustomBulkBar ? (
          <CustomBulkBar resource={resource} />
        ) : (
          <BulkActionBar
            count={c.selection.length}
            actions={c.bulkActions}
            custom={c.customBulkActions}
            pending={c.isBulkPending}
            singularName={c.singularName}
            pluralName={c.pluralName}
            onEdit={c.bulkEdit}
            onArchive={c.bulkArchive}
            onRestore={c.bulkRestore}
            onDelete={c.bulkRemove}
            onExport={() => exportSelection(c)}
            onCustom={c.runBulkAction}
            onClear={c.clearSelection}
          />
        )}

        <TablePagination
          page={c.page}
          pageSize={c.pageSize}
          total={c.total}
          totalPages={c.totalPages}
          onPageChange={c.setPage}
          onPageSizeChange={c.setPageSize}
        />
      </div>

      {!c.isFormPage && c.form.open && CustomForm && (
        <CustomForm resource={resource} item={c.form.item} onClose={c.form.close} />
      )}

      {!c.isFormPage && c.form.open && !CustomForm && (
        c.isSteps ? (
          <FormModalSteps
            resource={resource}
            item={c.form.item}
            onClose={c.form.close}
          />
        ) : resource.formView === "modal" ? (
          <FormModal
            resource={resource}
            item={c.form.item}
            defaults={c.form.defaults}
            onClose={c.form.close}
          />
        ) : (
          // Default + explicit "sheet" — right drawer / bottom sheet.
          <FormSheet
            resource={resource}
            item={c.form.item}
            defaults={c.form.defaults}
            onClose={c.form.close}
          />
        )
      )}

      <ConfirmModal
        open={c.confirmDelete.open}
        onConfirm={c.confirmDelete.confirm}
        onCancel={c.confirmDelete.cancel}
        title={`Delete ${c.singularName}`}
        description={`Are you sure you want to delete this ${c.singularName.toLowerCase()}? This action cannot be undone.`}
        confirmLabel="Delete"
        variant="danger"
        loading={c.isDeleting}
      />

      <ConfirmModal
        open={c.confirmBulkDelete.open}
        onConfirm={c.confirmBulkDelete.confirm}
        onCancel={c.confirmBulkDelete.cancel}
        title={`Delete ${c.selection.length} ${c.pluralName.toLowerCase()}`}
        description={`Are you sure you want to delete ${c.selection.length} ${c.pluralName.toLowerCase()}? This action cannot be undone.`}
        confirmLabel="Delete All"
        variant="danger"
        loading={c.isBulkDeleting}
      />

      <ConfirmModal
        open={c.confirmBulkArchive.open}
        onConfirm={c.confirmBulkArchive.confirm}
        onCancel={c.confirmBulkArchive.cancel}
        title={"Archive " + c.selection.length + " " + c.pluralName.toLowerCase()}
        description="Archived rows leave this list but keep their data. You can restore them from the Archived tab."
        confirmLabel="Archive"
        loading={c.isBulkPending}
      />

      {c.confirmCustom.open && c.confirmCustom.action && (
        <ConfirmModal
          open
          onConfirm={c.confirmCustom.confirm}
          onCancel={c.confirmCustom.cancel}
          title={c.confirmCustom.action.label}
          description={c.confirmCustom.action.confirm ?? ""}
          confirmLabel={c.confirmCustom.action.label}
          variant={c.confirmCustom.action.variant === "danger" ? "danger" : undefined}
          loading={c.isBulkPending}
        />
      )}

      {c.bulkEditor.open && (
        <BulkEditModal
          resource={resource}
          count={c.selection.length}
          pending={c.isBulkPending}
          onApply={c.applyBulkEdit}
          onClose={c.bulkEditor.close}
        />
      )}

      {c.importer.open && (
        <ImportModal
          resource={resource}
          onClose={() => c.importer.setOpen(false)}
        />
      )}
    </div>
  );
}
