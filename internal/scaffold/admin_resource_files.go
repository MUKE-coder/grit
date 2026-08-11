package scaffold

import "fmt"

// adminResourceTypes returns the resource type system (lib/resource.ts).
func adminResourceTypes() string {
	return `// Resource Definition Types — The foundation of Grit Admin Panel
// Define resources with defineResource() and get full CRUD pages automatically.

import type { ComponentType, ReactNode } from "react";

// ─── Column Definitions ─────────────────────────────────────────────

export type ColumnFormat = "text" | "badge" | "currency" | "date" | "relative" | "boolean" | "image" | "video" | "file" | "files" | "link" | "email" | "color" | "richtext" | "user";

export interface BadgeConfig {
  [value: string]: { color: string; label: string };
}

export interface ColumnDefinition {
  key: string;
  label: string;
  sortable?: boolean;
  searchable?: boolean;
  hidden?: boolean;
  width?: string;
  format?: ColumnFormat;
  badge?: BadgeConfig;
  currencyPrefix?: string;
  className?: string;
  // v3.31.15: optional custom cell renderer. Lets you pack multiple
  // fields into one column (Name + email stacked, price + currency
  // badge, status pill + relative date) without dropping out to a
  // hand-written page. Receives the full row so dotted keys aren't
  // necessary. When defined, takes precedence over format / badge.
  cell?: (row: Record<string, unknown>) => ReactNode;
  // v3.101.0: make a cell's value clickable. Two behaviors are built in —
  // "link" opens the row's detail page, "copy" copies the cell value to the
  // clipboard (with a brief check-mark) — or pass your own function to do
  // anything (open a modal, fire a mutation, deep-link elsewhere). It gets the
  // cell value and the full row. The click never triggers the row's other
  // actions. Generated resources set onClick: "link" on their first column so
  // the primary identifier is click-to-open out of the box.
  onClick?: ColumnClick;
}

// ColumnClick is a table cell's click behavior: a built-in ("link" → open the
// detail page, "copy" → copy the value) or a custom handler.
export type ColumnClick =
  | "link"
  | "copy"
  | ((value: unknown, row: Record<string, unknown>) => void);

// ─── Filter Definitions ─────────────────────────────────────────────

export type FilterType = "select" | "date-range" | "number-range" | "boolean";

export interface FilterOption {
  label: string;
  value: string;
  // Optional extras used by the card-style radio control: a secondary line
  // under the label, and a short right-aligned hint (e.g. "Days" / "Weeks").
  description?: string;
  hint?: string;
}

export interface FilterDefinition {
  key: string;
  label: string;
  type: FilterType;
  options?: FilterOption[];
  placeholder?: string;
}

// ─── Table Definitions ──────────────────────────────────────────────

export type TableAction = "create" | "view" | "edit" | "delete" | "export";
export type BulkAction = "delete" | "export";

// v3.104.0 — extra per-row actions rendered after the built-in view/edit/
// delete controls. Either link somewhere (href) or run a handler (onClick);
// both receive the row. Used by the Users resource to offer "Erase (GDPR)",
// which deep-links to the GDPR page with the subject pre-selected.
export interface RowActionDefinition {
  label: string;
  /** Link target. Takes precedence over onClick when both are set. */
  href?: (row: Record<string, unknown>) => string;
  onClick?: (row: Record<string, unknown>) => void;
  /** "danger" renders the label in the danger color. */
  variant?: "default" | "danger";
  /** Hide the action for rows where this returns false. */
  visible?: (row: Record<string, unknown>) => boolean;
}

export interface TableDefinition {
  columns: ColumnDefinition[];
  filters?: FilterDefinition[];
  searchable?: boolean;
  searchPlaceholder?: string;
  actions?: TableAction[];
  /** Extra per-row actions rendered after view/edit/delete. */
  rowActions?: RowActionDefinition[];
  bulkActions?: BulkAction[];
  defaultSort?: { key: string; direction: "asc" | "desc" };
  pageSize?: number;
  // v3.31.34 — date-window filter on this resource's list page.
  // Defaults to enabled with field="created_at", label="Created".
  // Set enabled:false to hide; override field to filter on a domain
  // column (e.g. "scheduled_for" for a Booking resource).
  dateFilter?: {
    enabled?: boolean;
    field?: string;
    label?: string;
  };
  // v3.31.35 — client-side export formats offered in the toolbar's
  // download menu. Defaults to all three on. Set the whole field to
  // false to hide the menu entirely; flip individual flags to hide a
  // single format. allPages (default true) means the menu fetches
  // every page from the API before building the file -- otherwise
  // only the rows currently on screen get exported.
  export?: false | {
    csv?: boolean;
    json?: boolean;
    excel?: boolean;
    allPages?: boolean;
  };
  // v3.31.35 — Excel import button + modal flow. Defaults to enabled.
  // Set to false to hide. fields restricts which form fields are
  // accepted in the upload (useful for excluding computed columns or
  // user-supplied IDs); defaults to every form field.
  import?: false | {
    excel?: boolean;
    fields?: string[];
  };
}

// ─── Form Field Definitions ─────────────────────────────────────────

export type FieldType = "text" | "textarea" | "number" | "select" | "date" | "datetime" | "toggle" | "checkbox" | "checkbox-group" | "radio" | "richtext" | "image" | "images" | "video" | "videos" | "file" | "files" | "relationship-select" | "multi-relationship-select" | "line-items";

export interface FieldDefinition {
  key: string;
  label: string;
  type: FieldType;
  required?: boolean;
  placeholder?: string;
  description?: string;
  defaultValue?: unknown;
  options?: FilterOption[];
  min?: number;
  max?: number;
  step?: number;
  prefix?: string;
  suffix?: string;
  rows?: number;
  colSpan?: 1 | 2;
  accept?: string;
  maxSize?: number;
  relatedEndpoint?: string;
  displayField?: string;
  relationshipKey?: string;

  // v3.113.0 — relationship-select / multi-relationship-select only. The
  // dropdown offers a "New <Related>" row that opens the related resource's own
  // form in a nested dialog and selects the record it creates. Requires the
  // related model to be a registered resource (the row is looked up by
  // relatedEndpoint) and the caller to hold <slug>.create. Set false to hide
  // the row on a field where creating on the fly is not appropriate.
  allowCreate?: boolean;

  // v3.114.0 — date / datetime only. Bounds the picker: days outside the range
  // are unselectable and the year dropdown only lists years inside it. ISO
  // "YYYY-MM-DD". Without them the year list runs 100 years back to 10 forward,
  // which covers a date of birth and a scheduling field alike.
  minDate?: string;
  maxDate?: string;

  // select field: load options from an endpoint at render time, on top of any
  // static options. optionsLabelKey/optionsValueKey default to "name".
  optionsUrl?: string;
  optionsLabelKey?: string;
  optionsValueKey?: string;

  // v3.31.30 — file / files field knobs. Set by the resource generator
  // from the CLI :file:<accepts> / :files:<accepts> syntax, but can be
  // overridden by hand in the resource definition.
  /** Accept-alias list ("image", "all", or e.g. ["pdf","doc"]). */
  accepts?: string[];
  /** Per-field max size in megabytes. Defaults: 5MB, 300MB for video. */
  maxSizeMB?: number;
  // v3.31.31 — visual knobs for the FileField / FilesField.
  /** Dropzone visual variant. "default" boxed-dashed, "compact" inline,
   *  "minimal" link, "avatar" circular for profile pics,
   *  "inline" tag-style. */
  dropzone?: "default" | "compact" | "minimal" | "avatar" | "inline";
  /** Progress indicator variant. "bar" (default linear), "circular"
   *  (donut with % inside), "pulse" (three dots + %, minimal). */
  progress?: "bar" | "circular" | "pulse";
  /** Allow up/down arrow reordering of files in the preview list.
   *  Multi-file (:files:) only. Defaults to true. */
  reorderable?: boolean;

  // v3.31.38 — number-input behaviour. Only applies when type === "number".
  /** Domain of the underlying Go column. Controls comma formatting:
   *  "int" allows negatives, no decimals; "uint" disallows negatives
   *  + decimals; "float" allows both. The generator sets this from
   *  the Go field type. Unset = "float" (legacy permissive). */
  numberKind?: "int" | "uint" | "float";

  // v3.103.0 — a visible field with a small "Generate" button in its label
  // row. Unlike an auto field (which is server-filled and hidden from the
  // form), this keeps the input visible and editable; clicking Generate runs
  // YOUR function with the current form values and fills the field with what it
  // returns (sync or async — e.g. call an endpoint, derive from another field,
  // mint a code). text / number fields only. You define this by hand in the
  // resource definition; the generator never emits it.
  generate?: (values: Record<string, unknown>) => string | number | Promise<string | number>;

  // ── Inline line-items (type === "line-items") ──────────────────────
  // Renders a child resource as an editable table INSIDE the parent form
  // (e.g. an Invoice's items). The rows are submitted as an array under
  // this field's key and saved atomically by the parent's create/update
  // handler (GORM has-many). Generated by "grit generate resource Parent
  // --items Child:fields", but hand-tunable.
  /** Columns of the inline table — the child's editable fields. Supports
   *  text / number / select / relationship-select / date per row. */
  itemFields?: FieldDefinition[];
  /** The child endpoint, used by the detail page's related table. */
  itemEndpoint?: string;
  /** The child's foreign-key column pointing back at the parent
   *  (e.g. "invoice_id"). */
  foreignKey?: string;
  /** Singular noun for the add-row button, e.g. "item" → "Add item". */
  itemNoun?: string;
}

export interface StepDefinition {
  title: string;
  description?: string;
  fields: string[];
}

// v3.31.18: groups unify the Create wizard and the Update cards view.
// On Create (sheet/modal/page) they render as a stepped wizard with
// Next/Back. On Update they render as per-group cards, each with its
// own Save button that PATCHes only that group's fields — so editing
// "Address" doesn't rewrite "Pricing".
//
// scope picks which contexts the group appears in:
//   "create"  — wizard step on Create only; hidden on Update
//   "update"  — card on Update only; hidden on Create
//   "both"    — both contexts (default)
//
// Useful pattern: minimal Create with title + price (scope: "create"),
// the rest deferred to Update cards (scope: "update").
export interface GroupDefinition {
  title: string;
  description?: string;
  fields: string[];
  scope?: "create" | "update" | "both";
}

export interface FormDefinition {
  fields: FieldDefinition[];
  layout?: "single" | "two-column";
  steps?: StepDefinition[];
  groups?: GroupDefinition[];
  fieldsPerStep?: number;
  stepVariant?: "horizontal" | "vertical";
  // v3.113.0 — on EDIT, give every step its own Update button that PATCHes only
  // that step's fields. Disabled until the step is actually changed, and back to
  // disabled once it saves. Defaults to on for stepped forms; set false to keep
  // the old behaviour of one submit at the end that rewrites every field.
  perStepSave?: boolean;
  // Drawer width for formView: "sheet". "half" (default) opens at 50% of the
  // viewport; "wide" opens at 80%. Either way the maximize button toggles to 80%.
  sheetWidth?: "half" | "wide";
}

// ─── Widget Definitions ─────────────────────────────────────────────

export type WidgetType = "stat" | "chart" | "activity";
export type ChartType = "line" | "bar" | "pie";
export type WidgetFormat = "number" | "currency" | "percentage";

export interface WidgetDefinition {
  type: WidgetType;
  label: string;
  endpoint?: string;
  icon?: string;
  color?: string;
  format?: WidgetFormat;
  chartType?: ChartType;
  limit?: number;
  colSpan?: 1 | 2 | 3 | 4;
}

export interface DashboardDefinition {
  // v3.31.44 -- set to false to hide the per-resource preset widgets
  // (Total + sparkline + Latest N) from the main dashboard. The
  // widgets are opt-in disabled, not opt-in enabled: every newly
  // generated resource gets them by default.
  enabled?: boolean;
  // Reserved for the custom widget builder (v3.31.40 dashboard
  // layout work). Existing resources may already declare widgets[];
  // the preset Total + Latest N widgets render even when this is
  // empty.
  widgets?: WidgetDefinition[];
}

// ─── Custom Components ──────────────────────────────────────────────
//
// A resource is config, and config cannot hold JSX: resources/<name>.ts is a
// .ts file, and the generator rewrites it. So anything with a component in it
// lives next door in resources/<name>.custom.tsx, which is written once and
// never touched again. defineResource() merges the two.
//
// The props below are deliberately the same shape as the components they
// replace. DataTable already satisfies ResourceTableProps, which means a swap
// is a swap — no adapter, and you can wrap the original by rendering it inside
// your own component.

/** Props a replacement table receives. Identical to DataTable's own props. */
export interface ResourceTableProps<T = Record<string, unknown>> {
  columns: ColumnDefinition[];
  data: T[];
  isLoading?: boolean;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
  onSort?: (key: string) => void;
  selectedRows?: string[];
  onSelectRows?: (rows: string[]) => void;
  onView?: (item: T) => void;
  onEdit?: (item: T) => void;
  onDelete?: (id: string) => void;
  rowActions?: RowActionDefinition[];
}

/** Props a replacement form receives. Identical to FormSheet / FormModal. */
export interface ResourceFormProps<T = Record<string, unknown>> {
  resource: ResourceDefinition;
  item: T | null;
  onClose: () => void;
}

/**
 * Props a replacement page receives. It gets only the resource — call
 * useResourceController(resource) inside to get the rest, exactly as the
 * stock page does.
 */
export interface ResourcePageSlotProps {
  resource: ResourceDefinition;
}

export interface ResourceComponents<T = Record<string, unknown>> {
  /** Replaces the whole list view. The last resort, and the most freedom. */
  Page?: ComponentType<ResourcePageSlotProps>;
  /** Replaces the table, keeping the header, toolbar, filters and pagination. */
  Table?: ComponentType<ResourceTableProps<T>>;
  /** Replaces the create / edit form in whichever container it opens in. */
  Form?: ComponentType<ResourceFormProps<T>>;
  /** Rendered instead of the table when there are no rows and none are loading. */
  EmptyState?: ComponentType<ResourcePageSlotProps>;
}

/**
 * The contents of resources/<name>.custom.tsx.
 *
 * columns and fields are patched by key rather than replaced wholesale, so
 * grit sync can keep adding new columns from the Go model without wiping
 * your renderers.
 */
export interface ResourceCustomisation<T = Record<string, unknown>> {
  components?: ResourceComponents<T>;
  /** Per-column overrides, keyed by column key. Merged over the generated column. */
  columns?: Record<string, Partial<ColumnDefinition>>;
  /** Per-field overrides, keyed by field key. Merged over the generated field. */
  fields?: Record<string, Partial<FieldDefinition>>;
}

// ─── Resource Definition ────────────────────────────────────────────

export interface ResourceDefinition {
  /** Set by defineResource() from resources/<name>.custom.tsx. */
  components?: ResourceComponents;
  name: string;
  slug: string;
  endpoint: string;
  icon: string;
  label?: { singular: string; plural: string };
  // How the Create / Edit form is presented:
  //   "sheet"        — right-drawer on desktop, bottom-sheet on mobile (default)
  //   "modal"        — centered dialog, best for short forms (1-6 fields)
  //   "page"         — a dedicated route at /resources/<slug>?action=create|edit
  //   "modal-steps"  — sheet/drawer with multi-step wizard
  //   "page-steps"   — dedicated page with multi-step wizard
  // Leave undefined to inherit the "sheet" default. (Pre-v3.31.17 the
  // bare "modal" value also rendered as a sheet — now "modal" is a
  // proper centered dialog. Switch to "sheet" if you preferred the
  // old behavior.)
  formView?: "sheet" | "modal" | "page" | "modal-steps" | "page-steps";
  table: TableDefinition;
  form: FormDefinition;
  dashboard?: DashboardDefinition;
  stats?: StatsConfig | boolean;
  // Optional sidebar nav grouping. Resources sharing the same group key
  // render under a collapsible group header in the sidebar.
  group?: string;
  // Hide this resource from the sidebar for users without ADMIN/EDITOR role.
  adminOnly?: boolean;
  // Hide this resource from the sidebar entirely (still routable + usable via
  // relationships). Set on inline --items children — you manage them through
  // the parent's form and detail page, not a top-level nav entry.
  hidden?: boolean;
}

// Stats cards shown above the data table on every resource page.
// See GRIT_STYLE_GUIDE §7.8 (Page Header).
// Set stats: false to disable stats on this resource page.
// Omit stats to get 4 auto-generated default cards (Total, This Week, This Month, Updated Recently).
// Provide stats: { cards: [...] } to fully customize.
export interface StatsConfig {
  enabled?: boolean;
  cards?: StatCardConfig[];
}

export interface StatCardConfig {
  label: string;
  icon?: string;
  color?: "default" | "success" | "warning" | "danger" | "info";
  value?: string | number;
  endpoint?: string;
  field?: string;
  trend?: { value: number; direction: "up" | "down" };
}

// ─── defineResource Helper ──────────────────────────────────────────

/**
 * Build a resource, optionally merged with the customisation sitting next to it.
 *
 * The second argument is the default export of resources/<name>.custom.tsx.
 * Splitting them is what makes both halves safe: the generator owns the config
 * file and rewrites it freely, while the custom file is written once and never
 * touched again.
 *
 * Columns and fields are patched per key rather than replaced, so grit sync can
 * go on adding new ones from the Go model without discarding your renderers.
 */
export function defineResource<T = Record<string, unknown>>(
  config: ResourceDefinition,
  custom?: ResourceCustomisation<T>,
): ResourceDefinition {
  const columnPatches = custom?.columns ?? {};
  const fieldPatches = custom?.fields ?? {};

  const columns = config.table.columns.map((col) =>
    columnPatches[col.key] ? { ...col, ...columnPatches[col.key] } : col,
  );

  const fields = config.form.fields.map((field) =>
    fieldPatches[field.key] ? { ...field, ...fieldPatches[field.key] } : field,
  );

  return {
    ...config,
    label: config.label ?? {
      singular: config.name,
      plural: config.slug.charAt(0).toUpperCase() + config.slug.slice(1),
    },
    components: custom?.components as ResourceComponents | undefined,
    table: {
      ...config.table,
      columns,
      pageSize: config.table.pageSize ?? 20,
      actions: config.table.actions ?? ["create", "view", "edit", "delete"],
      searchable: config.table.searchable ?? true,
    },
    form: {
      ...config.form,
      fields,
      layout: config.form.layout ?? "single",
    },
  };
}
`
}

// AdminResourceCustomStub returns resources/<slug>.custom.tsx — the half of a
// resource that the generator writes once and never touches again.
//
// It exists because the config half cannot hold components. resources/<slug>.ts
// is a .ts file, so JSX will not compile in it, and grit generate rewrites it
// whole. Anything with a component in it — a custom cell, a replacement table,
// a whole page — goes here instead and survives every regeneration.
func AdminResourceCustomStub(pascal string) string {
	return fmt.Sprintf(`import type { ResourceCustomisation } from "@/lib/resource";

/**
 * Customisations for the %s resource.
 *
 * This file is yours. The generator creates it once and never writes to it
 * again, so anything you put here survives grit generate and grit sync.
 *
 * Three things you can do:
 *
 *   1. Override a single cell, keeping everything else
 *
 *      columns: {
 *        status: { cell: (row) => <StatusPill value={String(row.status)} /> },
 *      }
 *
 *   2. Replace the table, keeping the header, toolbar, filters and pagination.
 *      The props are the same ones DataTable takes, so you can also wrap it:
 *
 *      components: {
 *        Table: (props) => <MyTemplateTable rows={props.data} onSort={props.onSort} />,
 *      }
 *
 *   3. Replace the whole page. Call useResourceController(resource) inside to
 *      get rows, sorting, paging, filters, selection and the CRUD actions:
 *
 *      components: { Page: MyProductsPage }
 *
 * See /docs/admin/custom-pages for the full list of props.
 */
const custom: ResourceCustomisation = {};

export default custom;
`, pascal)
}

// adminResourceRegistry returns the resource registry (resources/index.ts).
func adminResourceRegistry() string {
	return `import { usersResource } from "./users";
import { blogsResource } from "./blogs";
// grit:resources

import type { ResourceDefinition } from "@/lib/resource";

export const resources: ResourceDefinition[] = [
  usersResource,
  blogsResource,
  // grit:resource-list
];

export function getResource(slug: string): ResourceDefinition | undefined {
  return resources.find((r) => r.slug === slug);
}

export function getResourceByEndpoint(endpoint: string): ResourceDefinition | undefined {
  return resources.find((r) => r.endpoint === endpoint);
}
`
}

// adminUsersResource returns the users resource definition (resources/users.ts).
func adminUsersResource() string {
	return `import { defineResource } from "@/lib/resource";
import custom from "./users.custom";

export const usersResource = defineResource({
  name: "User",
  slug: "users",
  endpoint: "/api/users",
  icon: "Users",
  label: { singular: "User", plural: "Users" },

  table: {
    columns: [
      // v3.31.5: dropped the raw UUID column and packed first+last+email
      // into a single "user" cell so the table reads cleanly on small
      // screens. The "user" format renders avatar + name + email together.
      { key: "first_name", label: "Name", sortable: true, searchable: true, format: "user" },
      {
        key: "role",
        label: "Role",
        sortable: true,
        format: "badge",
        badge: {
          ADMIN: { color: "accent", label: "Admin" },
          EDITOR: { color: "info", label: "Editor" },
          USER: { color: "muted", label: "User" },
          // grit:role-badges
        },
      },
      { key: "job_title", label: "Job Title" },
      {
        key: "provider",
        label: "Provider",
        format: "badge",
        badge: {
          local: { color: "muted", label: "Email" },
          google: { color: "info", label: "Google" },
          github: { color: "accent", label: "GitHub" },
        },
      },
      { key: "active", label: "Status", format: "boolean" },
      { key: "created_at", label: "Created", format: "relative", sortable: true },
    ],
    filters: [
      {
        key: "role",
        label: "Role",
        type: "select",
        options: [
          { label: "Admin", value: "ADMIN" },
          { label: "Editor", value: "EDITOR" },
          { label: "User", value: "USER" },
          // grit:role-filters
        ],
      },
      { key: "active", label: "Status", type: "boolean" },
      {
        key: "provider",
        label: "Provider",
        type: "select",
        options: [
          { label: "Email", value: "local" },
          { label: "Google", value: "google" },
          { label: "GitHub", value: "github" },
        ],
      },
    ],
    searchable: true,
    searchPlaceholder: "Search by name or email...",
    actions: ["create", "view", "edit", "delete"],
    // Delete is an ordinary, reversible soft delete — it keeps the row and its
    // PII, and is deliberately NOT written to the GDPR journal. A real Art. 17
    // request needs an erasure, so link to the GDPR page with this user already
    // selected rather than leaving the two surfaces unconnected.
    rowActions: [
      {
        label: "Erase (GDPR)",
        variant: "danger",
        href: (row) => "/system/gdpr?user=" + String(row.id),
      },
    ],
    bulkActions: ["delete"],
    defaultSort: { key: "created_at", direction: "desc" },
    pageSize: 20,
  },

  form: {
    layout: "two-column",
    fields: [
      {
        key: "first_name",
        label: "First Name",
        type: "text",
        required: true,
        placeholder: "Enter first name",
        colSpan: 1,
      },
      {
        key: "last_name",
        label: "Last Name",
        type: "text",
        required: true,
        placeholder: "Enter last name",
        colSpan: 1,
      },
      {
        key: "email",
        label: "Email",
        type: "text",
        required: true,
        placeholder: "user@example.com",
        colSpan: 1,
      },
      {
        key: "password",
        label: "Password",
        type: "text",
        placeholder: "Enter password",
        description: "Required when creating a new user",
        colSpan: 1,
      },
      {
        key: "role",
        label: "Role",
        type: "select",
        required: true,
        // Loads every role from the database (built-in and custom), so a role
        // created at runtime through Roles & permissions is assignable here.
        // The static list stays as an offline fallback + the CLI injection point.
        optionsUrl: "/api/roles",
        options: [
          { label: "Admin", value: "ADMIN" },
          { label: "Editor", value: "EDITOR" },
          { label: "User", value: "USER" },
          // grit:role-options
        ],
        defaultValue: "USER",
        colSpan: 1,
      },
      {
        key: "job_title",
        label: "Job Title",
        type: "text",
        placeholder: "e.g. Software Engineer",
        colSpan: 1,
      },
      {
        key: "avatar",
        label: "Avatar",
        type: "image",
        description: "Profile picture",
        colSpan: 2,
      },
      {
        key: "active",
        label: "Active",
        type: "toggle",
        defaultValue: true,
        description: "Whether this user can log in",
        colSpan: 1,
      },
    ],
  },

  dashboard: {
    widgets: [
      {
        type: "stat",
        label: "Total Users",
        icon: "Users",
        color: "accent",
        endpoint: "/api/users?page_size=1",
        format: "number",
        colSpan: 1,
      },
      {
        type: "stat",
        label: "Active Users",
        icon: "UserCheck",
        color: "success",
        endpoint: "/api/users?active=true&page_size=1",
        format: "number",
        colSpan: 1,
      },
    ],
  },
}, custom);
`
}

// adminResourcePage returns the generic resource page component.
func adminResourcePage() string {
	return `"use client";

import { useSearchParams } from "next/navigation";
import dynamic from "next/dynamic";
import type { ResourceDefinition } from "@/lib/resource";
import { useResourceController } from "@/hooks/use-resource-controller";
import { PageHeader } from "@/components/layout/page-header";
import { DataTable } from "@/components/tables/data-table";
import { TableToolbar } from "@/components/tables/table-toolbar";
import { TablePagination } from "@/components/tables/table-pagination";
import { TableFilters } from "@/components/tables/table-filters";
// grit:resource:imports
import { buttonClasses } from "@/components/ui/button";

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

interface ResourcePageProps {
  resource: ResourceDefinition;
}

// v3.31.27: ResourcePage is a thin router. It picks between four possible
// views (UpdateGroups, FormPageSteps, FormPage, ResourceListView) based on
// formView + the ?action param. Before this split, the list-mode hooks all
// sat below the form-mode early returns — meaning the hook count varied
// between renders, which React 19 strict mode errors on. Splitting into two
// components keeps each function\'s hook list stable.
export function ResourcePage({ resource }: ResourcePageProps) {
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

// The default list view. Every piece of state and behaviour it uses comes from
// useResourceController — this component is markup and nothing else. That is
// deliberate: it is the proof that the hook is complete enough for someone to
// build their own page on, because the stock page is built on it too.
//
// Porting a bought template? Copy this file, keep the useResourceController
// line, and replace the JSX.
function ResourceListView({ resource }: ResourcePageProps) {
  const c = useResourceController(resource);

  // Slots, each falling back to the stock component. The props handed to a
  // custom Table are exactly DataTable's, so a replacement can also wrap the
  // original: (props) => <Card><DataTable {...props} /></Card>.
  const Table = resource.components?.Table ?? DataTable;
  const CustomForm = resource.components?.Form;
  const EmptyState = resource.components?.EmptyState;
  const showEmptyState = Boolean(EmptyState) && !c.isLoading && c.rows.length === 0;

  const headerActions = c.can("create") ? (
    <button onClick={c.create} className={buttonClasses({ size: "sm" })}>
      <span className="text-base leading-none">+</span>
      New {c.singularName}
    </button>
  ) : undefined;

  return (
    <div>
      <PageHeader
        title={c.pluralName}
        description={` + "`" + `Manage ${c.pluralName.toLowerCase()}` + "`" + `}
        actions={headerActions}
        stats={c.stats}
      />

      <div className="rounded-xl border border-border bg-bg-secondary">
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

        {showEmptyState && EmptyState ? (
          <EmptyState resource={resource} />
        ) : (
          <Table
            columns={c.columns}
            data={c.rows}
            isLoading={c.isLoading}
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
            onClose={c.form.close}
          />
        ) : (
          // Default + explicit "sheet" — right drawer / bottom sheet.
          <FormSheet
            resource={resource}
            item={c.form.item}
            onClose={c.form.close}
          />
        )
      )}

      <ConfirmModal
        open={c.confirmDelete.open}
        onConfirm={c.confirmDelete.confirm}
        onCancel={c.confirmDelete.cancel}
        title={` + "`" + `Delete ${c.singularName}` + "`" + `}
        description={` + "`" + `Are you sure you want to delete this ${c.singularName.toLowerCase()}? This action cannot be undone.` + "`" + `}
        confirmLabel="Delete"
        variant="danger"
        loading={c.isDeleting}
      />

      <ConfirmModal
        open={c.confirmBulkDelete.open}
        onConfirm={c.confirmBulkDelete.confirm}
        onCancel={c.confirmBulkDelete.cancel}
        title={` + "`" + `Delete ${c.selection.length} ${c.pluralName.toLowerCase()}` + "`" + `}
        description={` + "`" + `Are you sure you want to delete ${c.selection.length} ${c.pluralName.toLowerCase()}? This action cannot be undone.` + "`" + `}
        confirmLabel="Delete All"
        variant="danger"
        loading={c.isBulkDeleting}
      />

      {c.importer.open && (
        <ImportModal
          resource={resource}
          onClose={() => c.importer.setOpen(false)}
        />
      )}
    </div>
  );
}
`
}

// adminUsersPage returns the thin users page wrapper.
func adminUsersPage() string {
	return `"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { usersResource } from "@/resources/users";

export default function UsersPage() {
  return <ResourcePage resource={usersResource} />;
}
`
}

// adminUseResource returns the generic resource data hooks.

// adminUseResourceController returns hooks/use-resource-controller.ts — the
// headless brain behind a resource list page.
//
// Everything ResourcePage does apart from rendering lives here: query state and
// its URL round-trip, selection, column visibility, the create/edit/delete
// flows and the stat-card endpoints. ResourcePage is the first consumer, which
// is the point — if the default page cannot be rebuilt on this hook, the hook
// is missing something.
//
// The reason it exists is porting. Someone who has bought an admin template
// wants their own table and their own page shell but not to reimplement
// URL-synced sorting, paging, filters, bulk delete and toasts. They call this,
// render whatever they like, and keep all of it.
func adminUseResourceController() string {
	return `"use client";

import { useCallback, useMemo, useState } from "react";
import {
  usePathname,
  useRouter,
  useSearchParams,
  type ReadonlyURLSearchParams,
} from "next/navigation";
import type { ColumnDefinition, ResourceDefinition, TableAction } from "@/lib/resource";
import {
  useBulkDeleteResource,
  useDeleteResource,
  useResource,
} from "@/hooks/use-resource";
import type { StatCard } from "@/components/layout/page-header";
import { dateRangeToQueryParams, type DateRange } from "@/components/tables/date-filter";

// Read the date filter back out of the address bar so a refresh or a shared
// link rehydrates the same view.
function readDateRangeFromURL(sp: ReadonlyURLSearchParams | null): DateRange {
  if (!sp) return {};
  const preset = sp.get("date") as DateRange["preset"] | null;
  if (preset === "custom") {
    return {
      preset: "custom",
      from: sp.get("date_from") ?? undefined,
      to: sp.get("date_to") ?? undefined,
    };
  }
  if (preset === "today" || preset === "7d" || preset === "30d" || preset === "month") {
    return { preset };
  }
  return {};
}

// replace, not push: the back button should not collect one entry per filter
// tweak.
function writeDateRangeToURL(
  router: ReturnType<typeof useRouter>,
  pathname: string,
  current: ReadonlyURLSearchParams | null,
  range: DateRange,
) {
  const params = new URLSearchParams(current?.toString() ?? "");
  params.delete("date");
  params.delete("date_from");
  params.delete("date_to");
  if (range.preset) {
    params.set("date", range.preset);
    if (range.preset === "custom") {
      if (range.from) params.set("date_from", range.from);
      if (range.to) params.set("date_to", range.to);
    }
  }
  const qs = params.toString();
  router.replace(qs ? pathname + "?" + qs : pathname, { scroll: false });
}

export interface ResourceControllerOptions {
  /** Start on a page other than 1. */
  initialPage?: number;
  /** Override the resource's configured page size. */
  initialPageSize?: number;
}

export interface ResourceController<T = Record<string, unknown>> {
  resource: ResourceDefinition;

  // ── data ────────────────────────────────────────────────────────────
  rows: T[];
  meta: { total: number; page: number; page_size: number; pages: number } | undefined;
  total: number;
  totalPages: number;
  isLoading: boolean;

  // ── query state (all of it URL- or server-aware) ────────────────────
  page: number;
  pageSize: number;
  search: string;
  sortBy: string;
  sortOrder: "asc" | "desc";
  filters: Record<string, string>;
  dateRange: DateRange;
  setPage: (page: number) => void;
  setPageSize: (size: number) => void;
  setSearch: (value: string) => void;
  /** Toggles direction when the same key is passed twice. */
  setSort: (key: string) => void;
  setFilter: (key: string, value: string) => void;
  setDateRange: (range: DateRange) => void;

  // ── columns ─────────────────────────────────────────────────────────
  /** resource.table.columns minus hidden ones — what a table should render. */
  columns: ColumnDefinition[];
  allColumns: ColumnDefinition[];
  hiddenColumns: string[];
  toggleColumn: (key: string) => void;

  // ── selection ───────────────────────────────────────────────────────
  selection: string[];
  setSelection: (ids: string[]) => void;
  clearSelection: () => void;

  // ── actions ─────────────────────────────────────────────────────────
  actions: TableAction[];
  can: (action: TableAction) => boolean;
  create: () => void;
  edit: (row: T) => void;
  view: (row: T) => void;
  /** Opens the confirm dialog; deletion happens on confirm. */
  remove: (id: string) => void;
  bulkRemove: () => void;
  isDeleting: boolean;
  isBulkDeleting: boolean;

  // ── dialog state, for anyone rendering their own ────────────────────
  form: { open: boolean; item: T | null; close: () => void };
  confirmDelete: { open: boolean; confirm: () => void; cancel: () => void };
  confirmBulkDelete: { open: boolean; confirm: () => void; cancel: () => void };
  importer: { open: boolean; setOpen: (open: boolean) => void };

  // ── odds and ends the default page needs ────────────────────────────
  /** Same query the table ran, for an export that matches what is on screen. */
  apiSearchParams: URLSearchParams;
  stats: StatCard[] | undefined;
  singularName: string;
  pluralName: string;
  isFormPage: boolean;
  isSteps: boolean;
}

/**
 * Everything a resource list page needs except the markup.
 *
 * const c = useResourceController(productsResource)
 * <MyTable rows={c.rows} onSort={c.setSort} onRowClick={c.edit} />
 */
export function useResourceController<T = Record<string, unknown>>(
  resource: ResourceDefinition,
  options: ResourceControllerOptions = {},
): ResourceController<T> {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const isFormPage = resource.formView === "page" || resource.formView === "page-steps";
  const isSteps = resource.formView === "modal-steps" || resource.formView === "page-steps";

  const [page, setPage] = useState(options.initialPage ?? 1);
  const [pageSize, setPageSizeState] = useState(
    options.initialPageSize ?? resource.table.pageSize ?? 20,
  );
  const [search, setSearchState] = useState("");
  const [sortBy, setSortBy] = useState(resource.table.defaultSort?.key ?? "");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">(
    resource.table.defaultSort?.direction ?? "desc",
  );
  const [selection, setSelection] = useState<string[]>([]);
  const [filters, setFilters] = useState<Record<string, string>>({});
  const [hiddenColumns, setHiddenColumns] = useState<string[]>([]);

  const [dateRange, setDateRangeState] = useState<DateRange>(() =>
    readDateRangeFromURL(searchParams),
  );
  const dateParams = useMemo(() => dateRangeToQueryParams(dateRange), [dateRange]);
  const setDateRange = useCallback(
    (next: DateRange) => {
      setDateRangeState(next);
      writeDateRangeToURL(router, pathname, searchParams, next);
      setPage(1);
    },
    [router, pathname, searchParams],
  );

  const [formOpen, setFormOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<T | null>(null);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [bulkConfirmOpen, setBulkConfirmOpen] = useState(false);
  const [importOpen, setImportOpen] = useState(false);

  // Mirrors the query useResource builds, so an export applies the same
  // filter and sort the operator is looking at.
  const apiSearchParams = useMemo(() => {
    const sp = new URLSearchParams();
    if (search) sp.set("search", search);
    if (sortBy) {
      sp.set("sort_by", sortBy);
      sp.set("sort_order", sortOrder);
    }
    Object.entries(filters).forEach(([k, v]) => {
      if (v) sp.set(k, v);
    });
    Object.entries(dateParams).forEach(([k, v]) => {
      if (v) sp.set(k, v);
    });
    const df = resource.table.dateFilter?.field;
    if (df && df !== "created_at") sp.set("date_field", df);
    return sp;
  }, [search, sortBy, sortOrder, filters, dateParams, resource.table.dateFilter?.field]);

  const { data, isLoading } = useResource<T>(resource.endpoint, {
    page,
    pageSize,
    search,
    sortBy,
    sortOrder,
    filters,
    dateParams,
    dateField: resource.table.dateFilter?.field,
  });

  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  const { mutate: deleteItem, isPending: isDeleting } = useDeleteResource(
    resource.endpoint,
    singularName,
  );
  const { mutate: bulkDelete, isPending: isBulkDeleting } = useBulkDeleteResource(
    resource.endpoint,
    pluralName,
  );

  const columns = useMemo(
    () => resource.table.columns.filter((col) => !col.hidden && !hiddenColumns.includes(col.key)),
    [resource.table.columns, hiddenColumns],
  );

  const toggleColumn = useCallback((key: string) => {
    setHiddenColumns((prev) =>
      prev.includes(key) ? prev.filter((k) => k !== key) : [...prev, key],
    );
  }, []);

  // Any change to what is being queried resets to page 1 — otherwise a
  // search from page 7 lands on an empty page 7 of two results.
  const setSearch = useCallback((value: string) => {
    setSearchState(value);
    setPage(1);
  }, []);

  const setPageSize = useCallback((size: number) => {
    setPageSizeState(size);
    setPage(1);
  }, []);

  const setSort = useCallback(
    (key: string) => {
      if (sortBy === key) {
        setSortOrder((prev) => (prev === "asc" ? "desc" : "asc"));
      } else {
        setSortBy(key);
        setSortOrder("asc");
      }
      setPage(1);
    },
    [sortBy],
  );

  const setFilter = useCallback((key: string, value: string) => {
    setFilters((prev) => {
      if (!value) {
        const next = { ...prev };
        delete next[key];
        return next;
      }
      return { ...prev, [key]: value };
    });
    setPage(1);
  }, []);

  const clearSelection = useCallback(() => setSelection([]), []);

  const view = useCallback(
    (row: T) => {
      const id = String((row as Record<string, unknown>).id);
      router.push("/resources/" + resource.slug + "/" + id);
    },
    [router, resource.slug],
  );

  const edit = useCallback(
    (row: T) => {
      if (isFormPage) {
        const id = String((row as Record<string, unknown>).id);
        router.push("/resources/" + resource.slug + "?action=edit&edit=" + id);
      } else {
        setEditingItem(row);
        setFormOpen(true);
      }
    },
    [isFormPage, router, resource.slug],
  );

  const create = useCallback(() => {
    if (isFormPage) {
      router.push("/resources/" + resource.slug + "?action=create");
    } else {
      setEditingItem(null);
      setFormOpen(true);
    }
  }, [isFormPage, router, resource.slug]);

  const remove = useCallback((id: string) => {
    setDeletingId(id);
    setConfirmOpen(true);
  }, []);

  const doDelete = useCallback(() => {
    if (deletingId !== null) {
      deleteItem(deletingId, {
        onSuccess: () => {
          setConfirmOpen(false);
          setDeletingId(null);
        },
      });
    }
  }, [deleteItem, deletingId]);

  const bulkRemove = useCallback(() => {
    if (selection.length > 0) setBulkConfirmOpen(true);
  }, [selection]);

  const doBulkDelete = useCallback(() => {
    bulkDelete(selection, {
      onSuccess: () => {
        setBulkConfirmOpen(false);
        setSelection([]);
      },
    });
  }, [bulkDelete, selection]);

  const closeForm = useCallback(() => {
    setFormOpen(false);
    setEditingItem(null);
  }, []);

  const actions = resource.table.actions ?? ["create", "view", "edit", "delete"];
  const can = useCallback((action: TableAction) => actions.includes(action), [actions]);

  const statsConfig = resource.stats;
  const statsEnabled =
    statsConfig === undefined ||
    statsConfig === true ||
    (typeof statsConfig === "object" && statsConfig !== null && statsConfig.enabled !== false);

  const stats: StatCard[] | undefined = useMemo(() => {
    if (!statsEnabled) return undefined;

    // When a date range is active, every stat endpoint gets it too —
    // otherwise "Total: 10,000" sits above a table showing 142 matches.
    const applyDateParams = (cards: StatCard[]): StatCard[] => {
      if (Object.keys(dateParams).length === 0) return cards;
      return cards.map((card) => {
        if (!card.endpoint) return card;
        const sep = card.endpoint.includes("?") ? "&" : "?";
        const qs = new URLSearchParams(dateParams).toString();
        return { ...card, endpoint: card.endpoint + sep + qs };
      });
    };

    if (
      typeof statsConfig === "object" &&
      statsConfig !== null &&
      Array.isArray(statsConfig.cards) &&
      statsConfig.cards.length > 0
    ) {
      return applyDateParams(statsConfig.cards);
    }

    const ep = resource.endpoint;
    const defaults: StatCard[] = [
      { label: "Total", endpoint: ep + "?page_size=1", field: "meta.total", icon: resource.icon || "Package" },
      { label: "This Week", endpoint: ep + "?page_size=1&created_since=7d", field: "meta.total", icon: "TrendingUp", color: "success" },
      { label: "This Month", endpoint: ep + "?page_size=1&created_since=30d", field: "meta.total", icon: "Calendar", color: "info" },
      { label: "Updated Recently", endpoint: ep + "?page_size=1&updated_since=7d", field: "meta.total", icon: "RefreshCw" },
    ];
    return applyDateParams(defaults);
  }, [statsEnabled, statsConfig, resource.endpoint, resource.icon, dateParams]);

  return {
    resource,

    rows: data?.data ?? [],
    meta: data?.meta,
    total: data?.meta?.total ?? 0,
    totalPages: data?.meta?.pages ?? 1,
    isLoading,

    page,
    pageSize,
    search,
    sortBy,
    sortOrder,
    filters,
    dateRange,
    setPage,
    setPageSize,
    setSearch,
    setSort,
    setFilter,
    setDateRange,

    columns,
    allColumns: resource.table.columns,
    hiddenColumns,
    toggleColumn,

    selection,
    setSelection,
    clearSelection,

    actions,
    can,
    create,
    edit,
    view,
    remove,
    bulkRemove,
    isDeleting,
    isBulkDeleting,

    form: { open: formOpen, item: editingItem, close: closeForm },
    confirmDelete: {
      open: confirmOpen,
      confirm: doDelete,
      cancel: () => {
        setConfirmOpen(false);
        setDeletingId(null);
      },
    },
    confirmBulkDelete: {
      open: bulkConfirmOpen,
      confirm: doBulkDelete,
      cancel: () => setBulkConfirmOpen(false),
    },
    importer: { open: importOpen, setOpen: setImportOpen },

    apiSearchParams,
    stats,
    singularName,
    pluralName,
    isFormPage,
    isSteps,
  };
}
`
}

func adminUseResource() string {
	return `import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { apiClient } from "@/lib/api-client";

interface ResourceQueryParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: "asc" | "desc";
  filters?: Record<string, string>;
  // v3.31.34 — date-window filter. dateParams comes from
  // dateRangeToQueryParams(); dateField overrides the server's
  // default "created_at" target column when set.
  dateParams?: Record<string, string>;
  dateField?: string;
}

interface PaginatedResponse<T = Record<string, unknown>> {
  data: T[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

export function useResource<T = Record<string, unknown>>(
  endpoint: string,
  params: ResourceQueryParams = {}
) {
  const { page = 1, pageSize = 20, search, sortBy, sortOrder, filters, dateParams, dateField } = params;

  return useQuery<PaginatedResponse<T>>({
    // v3.31.34: dateParams + dateField included in key so a date
    // filter change invalidates the cache and the list refetches.
    queryKey: [endpoint, { page, pageSize, search, sortBy, sortOrder, filters, dateParams, dateField }],
    queryFn: async () => {
      const searchParams = new URLSearchParams({
        page: String(page),
        page_size: String(pageSize),
      });

      if (search) searchParams.set("search", search);
      if (sortBy) {
        searchParams.set("sort_by", sortBy);
        searchParams.set("sort_order", sortOrder ?? "desc");
      }
      if (filters) {
        Object.entries(filters).forEach(([key, value]) => {
          if (value) searchParams.set(key, value);
        });
      }
      if (dateParams) {
        Object.entries(dateParams).forEach(([key, value]) => {
          if (value) searchParams.set(key, value);
        });
      }
      if (dateField && dateField !== "created_at") {
        searchParams.set("date_field", dateField);
      }

      const { data } = await apiClient.get(` + "`" + `${endpoint}?${searchParams}` + "`" + `);
      return data;
    },
  });
}

export function useResourceItem<T = Record<string, unknown>>(
  endpoint: string,
  id: string,
  options?: { enabled?: boolean }
) {
  return useQuery<{ data: T }>({
    queryKey: [endpoint, id],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `${endpoint}/${id}` + "`" + `);
      return data;
    },
    enabled: (options?.enabled ?? true) && !!id,
  });
}

// Every mutation hook takes an optional resource label (the singular, e.g.
// "Invoice") so toasts name what actually happened — "Invoice created
// successfully" rather than a bare "Created successfully". Omitting it keeps
// the old generic wording, so existing call sites still compile.
function said(label: string | undefined, verb: string) {
  return label ? label + " " + verb : verb.charAt(0).toUpperCase() + verb.slice(1);
}

function failed(label: string | undefined, verb: string) {
  return label ? "Failed to " + verb + " " + label.toLowerCase() : "Failed to " + verb;
}

export function useCreateResource(endpoint: string, label?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (body: Record<string, unknown>) => {
      const { data } = await apiClient.post(endpoint, body);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success(said(label, "created successfully"));
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(axiosErr?.response?.data?.error?.message || failed(label, "create"));
    },
  });
}

export function useUpdateResource(endpoint: string, label?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, body }: { id: string; body: Record<string, unknown> }) => {
      const { data } = await apiClient.put(` + "`" + `${endpoint}/${id}` + "`" + `, body);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success(said(label, "updated successfully"));
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(axiosErr?.response?.data?.error?.message || failed(label, "update"));
    },
  });
}

// v3.31.18: partial updates for the grouped update view. Each group's
// Save button calls patch() with only the fields it owns. The Go-side
// Patch handler whitelists writable columns and silently drops anything
// else, so it's safe to send only a subset.
export function usePatchResource(endpoint: string, label?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, body }: { id: string; body: Record<string, unknown> }) => {
      const { data } = await apiClient.patch(` + "`" + `${endpoint}/${id}` + "`" + `, body);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success(label ? label + " saved" : "Saved");
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(axiosErr?.response?.data?.error?.message || failed(label, "save"));
    },
  });
}

export function useDeleteResource(endpoint: string, label?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete(` + "`" + `${endpoint}/${id}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success(said(label, "deleted successfully"));
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(axiosErr?.response?.data?.error?.message || failed(label, "delete"));
    },
  });
}

// Bulk delete reports a count, so it takes the PLURAL label ("Invoices")
// and names how many rows went — "3 Invoices deleted successfully".
export function useBulkDeleteResource(endpoint: string, pluralLabel?: string) {
  const queryClient = useQueryClient();
  const [count, setCount] = useState(0);

  return useMutation({
    // ids are strings because Grit's models use UUID primary keys
    // (the User.ID column in packages/shared/types/user.ts is 'string',
    // and the same is true for every grit generate'd model).
    mutationFn: async (ids: string[]) => {
      setCount(ids.length);
      await Promise.all(ids.map((id) => apiClient.delete(` + "`" + `${endpoint}/${id}` + "`" + `)));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success(
        pluralLabel ? count + " " + pluralLabel + " deleted successfully" : "Deleted successfully"
      );
    },
    onError: () => {
      toast.error(
        pluralLabel ? "Failed to delete some " + pluralLabel.toLowerCase() : "Failed to delete some items"
      );
    },
  });
}
`
}

// adminDashboardPage returns the enhanced dashboard page.
func adminDashboardPage() string {
	return fmt.Sprintf(`"use client";

import { useMe } from "@/hooks/use-auth";
import { resources } from "@/resources";
import { StatsCard } from "@/components/widgets/stats-card";
import { WidgetGrid } from "@/components/widgets/widget-grid";
import { getIcon } from "@/lib/icons";

// The API origin the browser talks to. Hardcoding localhost:8080 here meant
// the Quick Links pointed at the wrong port whenever the API moved.
const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function AdminDashboard() {
  const { data: user } = useMe();
  const allWidgets = resources.flatMap((r) => r.dashboard?.widgets ?? []);

  const greeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return "Good morning";
    if (hour < 17) return "Good afternoon";
    return "Good evening";
  };

  return (
    <div className="space-y-8">
      {/* Welcome header */}
      <div className="rounded-xl border border-border bg-gradient-to-r from-accent/10 via-bg-secondary to-bg-secondary p-6 sm:p-8">
        <h1 className="text-2xl font-bold text-foreground">
          {greeting()}, {user?.first_name || "Admin"}
        </h1>
        <p className="text-text-secondary mt-1">
          Here&apos;s an overview of your application.
        </p>
      </div>

      {/* Stats widgets */}
      {allWidgets.length > 0 ? (
        <WidgetGrid widgets={allWidgets} />
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard label="Total Resources" value="—" icon="Database" color="accent" />
          <StatsCard label="Registered" value={String(resources.length)} icon="Layers" color="success" />
        </div>
      )}

      {/* Quick Actions + System */}
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Resources */}
        <div className="lg:col-span-2 rounded-xl border border-border bg-bg-secondary p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Resources</h2>
          <div className="grid gap-3 sm:grid-cols-2">
            {resources.map((r) => {
              const Icon = getIcon(r.icon);
              return (
                <a
                  key={r.slug}
                  href={%s/resources/${r.slug}%s}
                  className="flex items-center gap-4 rounded-lg border border-border bg-bg-tertiary p-4 hover:border-accent/30 hover:bg-bg-hover transition-all group"
                >
                  <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-accent/10 group-hover:bg-accent/20 transition-colors">
                    <Icon className="h-5 w-5 text-accent" />
                  </div>
                  <div>
                    <h3 className="font-medium text-foreground group-hover:text-accent transition-colors">
                      {r.label?.plural ?? r.name}
                    </h3>
                    <p className="text-xs text-text-muted">
                      Manage {(r.label?.plural ?? r.slug).toLowerCase()}
                    </p>
                  </div>
                </a>
              );
            })}
          </div>
        </div>

        {/* Quick Links */}
        <div className="rounded-xl border border-border bg-bg-secondary p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Quick Links</h2>
          <div className="space-y-2">
            <a
              href={API_URL + "/studio"}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 rounded-lg border border-border bg-bg-tertiary px-4 py-3 hover:border-accent/30 hover:bg-bg-hover transition-all group"
            >
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-info/10">
                <span className="text-info text-sm font-bold">DB</span>
              </div>
              <div>
                <p className="text-sm font-medium text-foreground group-hover:text-accent transition-colors">GORM Studio</p>
                <p className="text-xs text-text-muted">Browse database</p>
              </div>
            </a>
            <a
              href={API_URL + "/api/health"}
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-3 rounded-lg border border-border bg-bg-tertiary px-4 py-3 hover:border-accent/30 hover:bg-bg-hover transition-all group"
            >
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-success/10">
                <span className="text-success text-sm font-bold">OK</span>
              </div>
              <div>
                <p className="text-sm font-medium text-foreground group-hover:text-accent transition-colors">API Health</p>
                <p className="text-xs text-text-muted">Check status</p>
              </div>
            </a>
            <a
              href="/system/jobs"
              className="flex items-center gap-3 rounded-lg border border-border bg-bg-tertiary px-4 py-3 hover:border-accent/30 hover:bg-bg-hover transition-all group"
            >
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-warning/10">
                <span className="text-warning text-sm font-bold">Q</span>
              </div>
              <div>
                <p className="text-sm font-medium text-foreground group-hover:text-accent transition-colors">Job Queue</p>
                <p className="text-xs text-text-muted">Background jobs</p>
              </div>
            </a>
            <a
              href="/system/files"
              className="flex items-center gap-3 rounded-lg border border-border bg-bg-tertiary px-4 py-3 hover:border-accent/30 hover:bg-bg-hover transition-all group"
            >
              <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-accent/10">
                <span className="text-accent text-sm font-bold">S3</span>
              </div>
              <div>
                <p className="text-sm font-medium text-foreground group-hover:text-accent transition-colors">File Storage</p>
                <p className="text-xs text-text-muted">Manage uploads</p>
              </div>
            </a>
          </div>
        </div>
      </div>
    </div>
  );
}
`, "`", "`")
}

// adminConfirmModal returns the reusable confirm modal component.
func adminConfirmModal() string {
	return `"use client";

import { AlertCircle, Loader2 } from "@/lib/icons";

interface ConfirmModalProps {
  open: boolean;
  onConfirm: () => void;
  onCancel: () => void;
  title?: string;
  description?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: "danger" | "default";
  loading?: boolean;
}

export function ConfirmModal({
  open,
  onConfirm,
  onCancel,
  title = "Are you sure?",
  description = "This action cannot be undone.",
  confirmLabel = "Confirm",
  cancelLabel = "Cancel",
  variant = "default",
  loading = false,
}: ConfirmModalProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={onCancel} />
      <div className="relative z-10 w-full max-w-sm rounded-xl border border-border bg-bg-secondary p-6 shadow-2xl animate-in fade-in zoom-in-95 duration-200">
        <div className="flex items-start gap-4">
          <div className={` + "`" + `flex h-10 w-10 shrink-0 items-center justify-center rounded-full ${
            variant === "danger" ? "bg-danger/10" : "bg-accent/10"
          }` + "`" + `}>
            <AlertCircle className={` + "`" + `h-5 w-5 ${
              variant === "danger" ? "text-danger" : "text-accent"
            }` + "`" + `} />
          </div>
          <div className="space-y-2">
            <h3 className="text-lg font-semibold text-foreground">{title}</h3>
            <p className="text-sm text-text-secondary">{description}</p>
          </div>
        </div>

        <div className="mt-6 flex items-center justify-end gap-3">
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover transition-colors disabled:opacity-50"
          >
            {cancelLabel}
          </button>
          <button
            type="button"
            onClick={onConfirm}
            disabled={loading}
            className={` + "`" + `flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium text-white transition-colors disabled:opacity-50 ${
              variant === "danger"
                ? "bg-danger hover:bg-danger/90"
                : "bg-accent hover:bg-accent-hover"
            }` + "`" + `}
          >
            {loading && <Loader2 className="h-4 w-4 animate-spin" />}
            {confirmLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
`
}

// adminViewModal returns the resource view modal component.
func adminViewModal() string {
	return `"use client";

import type { ResourceDefinition } from "@/lib/resource";
import { renderCell } from "@/components/tables/cell-renderers";
import { X, Pencil } from "@/lib/icons";

interface ViewModalProps {
  resource: ResourceDefinition;
  item: Record<string, unknown>;
  onClose: () => void;
  onEdit?: (item: Record<string, unknown>) => void;
}

export function ViewModal({ resource, item, onClose, onEdit }: ViewModalProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/50" onClick={onClose} />
      <div className="relative z-10 w-full max-w-lg max-h-[90vh] overflow-y-auto rounded-xl border border-border bg-bg-secondary shadow-2xl">
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <h2 className="text-lg font-semibold text-foreground">
            {resource.label?.singular ?? resource.name} Details
          </h2>
          <div className="flex items-center gap-2">
            {onEdit && (
              <button
                onClick={() => { onClose(); onEdit(item); }}
                className="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-sm font-medium text-accent hover:bg-accent/10 transition-colors"
              >
                <Pencil className="h-3.5 w-3.5" />
                Edit
              </button>
            )}
            <button
              onClick={onClose}
              className="rounded-lg p-1 text-text-secondary hover:bg-bg-hover hover:text-foreground transition-colors"
            >
              <X className="h-5 w-5" />
            </button>
          </div>
        </div>

        <div className="p-6">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-5">
            {resource.table.columns.map((col) => {
              const value = item[col.key];

              return (
                <div key={col.key} className="space-y-1.5">
                  <p className="text-xs font-medium text-text-muted uppercase tracking-wider">
                    {col.label}
                  </p>
                  <div className="text-sm text-foreground">
                    {value !== null && value !== undefined
                      ? renderCell(col, value, item)
                      : <span className="text-text-muted">—</span>
                    }
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        <div className="flex items-center justify-end gap-3 border-t border-border px-6 py-4">
          <button
            onClick={onClose}
            className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
`
}

// adminBlogsResource returns the blogs resource definition (resources/blogs.ts).
func adminBlogsResource() string {
	return `import { defineResource } from "@/lib/resource";
import custom from "./blogs.custom";

export const blogsResource = defineResource({
  name: "Blog",
  slug: "blogs",
  endpoint: "/api/admin/blogs",
  icon: "Newspaper",
  label: { singular: "Blog", plural: "Blogs" },

  table: {
    columns: [
      // v3.31.5: dropped the raw UUID column. Title + status + author
      // already identify a blog row clearly; the ID lives in the URL when
      // you open the detail view.
      { key: "title", label: "Title", sortable: true, searchable: true },
      { key: "slug", label: "Slug" },
      { key: "image", label: "Image", format: "image" },
      {
        key: "published",
        label: "Status",
        format: "badge",
        badge: {
          true: { color: "success", label: "Published" },
          false: { color: "muted", label: "Draft" },
        },
      },
      { key: "published_at", label: "Published At", format: "relative", sortable: true },
      { key: "created_at", label: "Created", format: "relative", sortable: true },
    ],
    filters: [
      {
        key: "published",
        label: "Status",
        type: "select",
        options: [
          { label: "Published", value: "true" },
          { label: "Draft", value: "false" },
        ],
      },
    ],
    searchable: true,
    searchPlaceholder: "Search blogs by title...",
    actions: ["create", "view", "edit", "delete"],
    bulkActions: ["delete"],
    defaultSort: { key: "created_at", direction: "desc" },
    pageSize: 20,
  },

  form: {
    layout: "single",
    fields: [
      {
        key: "title",
        label: "Title",
        type: "text",
        required: true,
        placeholder: "Enter blog title",
      },
      {
        key: "excerpt",
        label: "Excerpt",
        type: "textarea",
        placeholder: "Brief summary of the blog post",
      },
      {
        key: "content",
        label: "Content",
        type: "richtext",
      },
      {
        key: "image",
        label: "Cover Image",
        type: "image",
      },
      {
        key: "published",
        label: "Published",
        type: "toggle",
      },
    ],
  },
}, custom);
`
}

// adminBlogsPage returns the blogs resource page.
func adminBlogsPage() string {
	return `"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { blogsResource } from "@/resources/blogs";

export default function BlogsPage() {
  return <ResourcePage resource={blogsResource} />;
}
`
}
