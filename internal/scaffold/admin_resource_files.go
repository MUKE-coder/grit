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

export interface ColumnDefinition<T = Record<string, unknown>> {
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
  cell?: (row: T) => ReactNode;
  // v3.101.0: make a cell's value clickable. Two behaviors are built in —
  // "link" opens the row's detail page, "copy" copies the cell value to the
  // clipboard (with a brief check-mark) — or pass your own function to do
  // anything (open a modal, fire a mutation, deep-link elsewhere). It gets the
  // cell value and the full row. The click never triggers the row's other
  // actions. Generated resources set onClick: "link" on their first column so
  // the primary identifier is click-to-open out of the box.
  onClick?: ColumnClick<T>;
}

// ColumnClick is a table cell's click behavior: a built-in ("link" → open the
// detail page, "copy" → copy the value) or a custom handler.
export type ColumnClick<T = Record<string, unknown>> =
  | "link"
  | "copy"
  | ((value: unknown, row: T) => void);

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

/**
 * Built-in bulk actions, offered once rows are selected.
 *
 *   edit     one field, one value, written to every selected row
 *   archive  put away without destroying: still listable under Archived,
 *            still exportable, restorable in one click
 *   restore  the inverse, shown only while the Archived view is open
 *   delete   soft delete, the same as the per-row Delete
 *   export   download the selection rather than the whole table
 *
 * archive and restore need the resource's model to carry archived_at, which
 * every model from grit generate resource has since v3.142.0. A resource without
 * the column should not list them.
 */
export type BulkAction = "edit" | "archive" | "restore" | "delete" | "export";

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

/**
 * A bulk action of your own, supplied from resources/<name>.custom.tsx.
 *
 * The built-in five cover put-away and delete. Everything domain-shaped is
 * yours: "Send invoices", "Assign to rep", "Mark as shipped". It goes in the
 * overlay rather than the resource definition because it holds a function,
 * and the resource definition is a .ts file the generator rewrites.
 */
export interface CustomBulkAction<T = Record<string, unknown>> {
  /** Stable key, used for the React key and for keeping order deterministic. */
  key: string;
  label: string;
  /** Any name from lib/icons. Rendered before the label. */
  icon?: string;
  /** "danger" colours it red. Reserve it for the irreversible. */
  variant?: "default" | "danger";
  /**
   * Ask first. A string is the dialog's body; the title and buttons come from
   * the label. Omit for actions that do not need it: a confirm on everything
   * trains people to dismiss confirms.
   */
  confirm?: string;
  /**
   * Runs the action. Receives the selected ids and the rows behind them, so
   * you can act without refetching. Return a promise and the bar shows a
   * pending state until it settles.
   *
   * The second argument carries what the page can do for you: refresh the
   * list, clear the selection, and announce a result to screen readers.
   */
  onSelect: (
    ids: string[],
    rows: T[],
    helpers: {
      refresh: () => void;
      clearSelection: () => void;
      announce: (message: string) => void;
    },
  ) => void | Promise<unknown>;
  /** Hide the action for some selections, e.g. only when exactly one row is on. */
  visible?: (rows: T[]) => boolean;
}

/** One filter preset in the tab strip above a table. */
export interface TableTab {
  /** Stable key. Also the value written to the URL, so keep it URL-safe. */
  key: string;
  label: string;
  /**
   * Query parameters this tab applies. Merged over the resource's own filters,
   * and cleared when another tab is chosen, so tabs never accumulate.
   * Omit for an "All" tab.
   */
  filters?: Record<string, string>;
  /**
   * Fetch and show a count on this tab. One extra request per tab that asks
   * for it, which is why it is not the default.
   */
  count?: boolean;
  /** Any name from lib/icons, rendered before the label. */
  icon?: string;
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
  /**
   * Filter presets shown as a tab strip above the table.
   *
   * A tab is a named set of query parameters. "Unpaid" is not a different
   * page, it is this page with status=pending, and a tab says that more
   * plainly than a dropdown someone has to open to discover.
   *
   *   tabs: [
   *     { key: "all", label: "All" },
   *     { key: "unpaid", label: "Unpaid", filters: { status: "pending" } },
   *     { key: "overdue", label: "Overdue", filters: { status: "pending", overdue: "true" } },
   *   ]
   *
   * The first tab is selected on load, and a tab with no filters clears them,
   * which is what makes "All" work without a special case.
   *
   * Counts are opt-in per tab because each one costs a request. Set
   * count: true and the tab fetches its own total with page_size=1; the badge
   * appears when it arrives rather than reserving space for a number that may
   * never come.
   *
   * These are config, so they live in the resource definition. Anything that
   * needs a function or JSX belongs in the overlay instead.
   */
  tabs?: TableTab[];
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
  /**
   * The column carries a unique constraint. Set by the generator from the
   * :unique field modifier. Bulk edit reads it to leave the field out: writing
   * one SKU to forty rows is either a constraint violation or, worse, not one.
   */
  unique?: boolean;
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
  columns: ColumnDefinition<T>[];
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

/** A resource whose belongs_to points at the one being viewed. */
export interface RelatedResource {
  resource: ResourceDefinition;
  /** Its foreign-key field, already resolved from the registry. */
  fk: string;
}

/**
 * What useResourceDetailController(resource, id) returns.
 *
 * Declared here rather than beside the hook so a slot's props can name it
 * without the types file and the hooks file importing each other.
 */
export interface ResourceDetailController<T = Record<string, unknown>> {
  resource: ResourceDefinition;
  id: string;

  // ── data ────────────────────────────────────────────────────────────
  record: T | undefined;
  isLoading: boolean;
  /** True once loading has finished and there is still nothing. */
  notFound: boolean;
  /** The first human-readable field on the record, for a page title. */
  title: string;

  // ── what to show ────────────────────────────────────────────────────
  /** Visible table columns, reused as the detail field list. */
  columns: ColumnDefinition[];
  /** line-items fields declared on this resource, rendered inline. */
  lineItemFields: FieldDefinition[];
  /** Resources whose belongs_to points here, discovered from the registry. */
  related: RelatedResource[];

  // ── actions ─────────────────────────────────────────────────────────
  edit: () => void;
  /** Opens the confirm dialog; deletion happens on confirm. */
  remove: () => void;
  print: () => void;
  /** Fetches the server-rendered PDF and opens it. */
  downloadPdf: () => Promise<void>;
  back: () => void;
  isPdfBusy: boolean;
  isDeleting: boolean;

  // ── dialog state, for anyone rendering their own ────────────────────
  //
  // form carries the item as well as the flag, mirroring the list
  // controller. Without it every caller writes item={c.record} and hits the
  // difference between "still loading" (undefined) and "creating" (null),
  // which the stock form distinguishes and a query result does not.
  form: {
    open: boolean;
    item: T | null;
    /** Pre-filled values for create mode, set by createWith(). */
    defaults?: Record<string, unknown>;
    close: () => void;
  };
  confirmDelete: { open: boolean; confirm: () => void; cancel: () => void };
}

/**
 * Props the whole-page detail slot receives. It gets only the resource and the
 * id: call useResourceDetailController(resource, id) inside for the rest, as
 * the stock page does. A page that replaces everything owns its dialogs too.
 */
export interface ResourceDetailSlotProps {
  resource: ResourceDefinition;
  id: string;
}

/**
 * Props the PART slots receive: the controller itself, already built.
 *
 * This is the important difference from the whole-page slot, and it is not a
 * convenience. Every call to useResourceDetailController creates its own
 * state, so a header that built its own controller would open an edit sheet
 * the page around it never reads: you press Edit and nothing happens. Sharing
 * one controller is what makes the parts able to drive the page.
 */
export interface ResourceDetailPartProps<T = Record<string, unknown>> {
  resource: ResourceDefinition;
  id: string;
  controller: ResourceDetailController<T>;
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
  /**
   * Replaces the bar that appears when rows are selected. It gets only the
   * resource; call useResourceController(resource) inside for the selection,
   * the actions and the pending state, exactly as the stock bar does.
   */
  BulkBar?: ComponentType<ResourcePageSlotProps>;

  // ── the detail page ──────────────────────────────────────────────────
  //
  // Same three tiers as the list view: swap a piece, or take the whole page.
  // These receive the resource and the record id, and call
  // useResourceDetailController(resource, id) inside for the rest, exactly as
  // the stock detail page does.

  /** Replaces the entire detail view, including its header and dialogs. */
  DetailPage?: ComponentType<ResourceDetailSlotProps>;
  /** Replaces the title block and its actions, keeping the body below. */
  DetailHeader?: ComponentType<ResourceDetailPartProps<T>>;
  /** Replaces the field list, keeping the header and the related sections. */
  DetailFields?: ComponentType<ResourceDetailPartProps<T>>;
  /** Rendered after the fields and before the related tables. */
  DetailAside?: ComponentType<ResourceDetailPartProps<T>>;
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
  /**
   * Per-column overrides, keyed by column key. Merged over the generated column.
   *
   * Typed against T, so a cell renderer gets the real row:
   *   columns: { total: { cell: (row) => <b>{row.total.toFixed(2)}</b> } }
   * with row.total known to be a number rather than unknown.
   */
  columns?: Record<string, Partial<ColumnDefinition<T>>>;
  /** Per-field overrides, keyed by field key. Merged over the generated field. */
  fields?: Record<string, Partial<FieldDefinition>>;
  /**
   * Bulk actions of your own, appended after the built-in ones.
   *
   * Here rather than in the resource definition because they hold functions,
   * and the definition is a .ts file the generator rewrites in full.
   */
  bulkActions?: CustomBulkAction<T>[];
}

// ─── Resource Definition ────────────────────────────────────────────

export interface ResourceDefinition {
  /** Set by defineResource() from resources/<name>.custom.tsx. */
  components?: ResourceComponents;
  /** Set by defineResource() from resources/<name>.custom.tsx. */
  customBulkActions?: CustomBulkAction[];
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
  // Set by "grit generate resource --tree". Adds a Table / Tree toggle to the
  // list page, where the tree view can reparent and reorder by dragging.
  //
  // It needs the endpoints --tree generates (/tree, /:id/move, /reorder,
  // /rebuild-tree), so setting it by hand on a resource that has no parent
  // column gives you a view that cannot load.
  tree?: boolean;
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

  // The patches are typed against the caller's row type; the registry stores
  // the erased form, and a (row: Product) => ReactNode is not assignable to a
  // (row: Record<string, unknown>) => ReactNode under strictFunctionTypes. One
  // cast here is what buys a typed authoring surface everywhere else.
  const columns = config.table.columns.map((col) => {
    const patch = columnPatches[col.key] as Partial<ColumnDefinition> | undefined;
    return patch ? { ...col, ...patch } : col;
  });

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
    // Erased the same way the components are, and for the same reason: the
    // registry holds every resource in one array, so it cannot be generic.
    customBulkActions: custom?.bulkActions as CustomBulkAction[] | undefined,
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
//
// typeName is the shared row type (e.g. "Product", exported from
// @repo/shared/types by grit sync). Pass "" when there is no matching type and
// the stub falls back to the untyped form, which still compiles.
func AdminResourceCustomStub(pascal, typeName string) string {
	typeImport := ""
	generic := ""
	rowNote := "Rows are Record<string, unknown> here. Generate or sync the resource to\n * pick up a typed row from @repo/shared/types."
	if typeName != "" {
		typeImport = fmt.Sprintf("import type { %s } from \"@repo/shared/types\";\n", typeName)
		generic = fmt.Sprintf("<%s>", typeName)
		rowNote = fmt.Sprintf("Rows are typed as %s, so row.id and friends autocomplete and a renamed\n * column fails at compile time rather than at runtime.", typeName)
	}

	return fmt.Sprintf(`import type { ResourceCustomisation } from "@/lib/resource";
%s
/**
 * Customisations for the %s resource.
 *
 * This file is yours. The generator creates it once and never writes to it
 * again, so anything you put here survives grit generate and grit sync.
 *
 * %s
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
const custom: ResourceCustomisation%s = {};

export default custom;
`, typeImport, pascal, rowNote, generic)
}

// adminResourceRegistry returns the resource registry (resources/index.ts).
func adminResourceRegistry() string {
	return `import { usersResource } from "./users/users";
import { blogsResource } from "./blogs/blogs";
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
    // No "archive": the scaffold's own models predate archived_at. A
    // generated resource gets the column and the full set.
    bulkActions: ["edit", "export", "delete"],
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

import { useState } from "react";
import { useSearchParams } from "next/navigation";
import dynamic from "next/dynamic";
import type { ResourceDefinition } from "@/lib/resource";
import { useResourceController } from "@/hooks/use-resource-controller";
import type { ResourceController } from "@/hooks/use-resource-controller";
import { PageHeader } from "@/components/layout/page-header";
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

// Exports the ticked rows rather than the page or the whole table. The rows
// are already in memory, so this needs no request: the point of "export
// selection" is the selection.
function exportSelection(c: ResourceController) {
  if (c.selectedRows.length === 0) return;
  exportToFile(c.selectedRows, c.columns, c.resource.slug, "csv");
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
                  : "text-text-secondary hover:bg-bg-hover hover:text-text-primary")
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
                    ? "bg-accent text-white"
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
`
}

// adminUsersPage returns the thin users page wrapper.
func adminUsersPage() string {
	return `"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { usersResource } from "@/resources/users/users";

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
import { useQueryClient } from "@tanstack/react-query";
import type {
  BulkAction,
  ColumnDefinition,
  CustomBulkAction,
  ResourceDefinition,
  TableAction,
  TableTab,
} from "@/lib/resource";
import {
  useBulkResource,
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
  /** Create with fields pre-filled, e.g. createWith({ parent_id: id }). */
  createWith: (defaults: Record<string, unknown>) => void;
  edit: (row: T) => void;
  view: (row: T) => void;
  /** Opens the confirm dialog; deletion happens on confirm. */
  remove: (id: string) => void;
  bulkRemove: () => void;
  isDeleting: boolean;
  isBulkDeleting: boolean;

  // ── bulk actions ────────────────────────────────────────────────────
  /** Built-ins the resource has switched on, minus any that make no sense
   *  in the current view (restore only appears while Archived is open). */
  bulkActions: BulkAction[];
  /** Custom ones from resources/<name>.custom.tsx, already filtered by visible(). */
  customBulkActions: CustomBulkAction<T>[];
  /** The rows behind the current selection, readable without a refetch. */
  selectedRows: T[];
  /** Opens the confirm dialog; archiving happens on confirm. */
  bulkArchive: () => void;
  /** Restores immediately: putting something back is not destructive. */
  bulkRestore: () => void;
  /** Opens the bulk edit dialog. */
  bulkEdit: () => void;
  /** Writes one field to every selected row and closes the dialog. */
  applyBulkEdit: (patch: Record<string, unknown>) => void;
  /** Runs a custom action, handing it the ids, the rows and the helpers. */
  runBulkAction: (action: CustomBulkAction<T>) => void;
  isBulkPending: boolean;
  /** Re-runs the list query. Handed to custom actions so they can refresh. */
  refresh: () => void;
  /** Speaks to the page's live region. Bulk changes never move focus. */
  announce: (message: string) => void;
  liveMessage: string;

  // ── tabs ────────────────────────────────────────────────────────────
  /** The resource's filter presets, or an empty array when it has none. */
  tabs: TableTab[];
  /** Key of the tab currently applied. "" when the resource has no tabs. */
  activeTab: string;
  setActiveTab: (key: string) => void;

  // ── archived view ───────────────────────────────────────────────────
  /** True while the Archived tab is open. */
  showArchived: boolean;
  setShowArchived: (value: boolean) => void;

  // ── dialog state, for anyone rendering their own ────────────────────
  form: {
    open: boolean;
    item: T | null;
    /** Pre-filled values for create mode, set by createWith(). */
    defaults?: Record<string, unknown>;
    close: () => void;
  };
  confirmDelete: { open: boolean; confirm: () => void; cancel: () => void };
  confirmBulkDelete: { open: boolean; confirm: () => void; cancel: () => void };
  confirmBulkArchive: { open: boolean; confirm: () => void; cancel: () => void };
  bulkEditor: { open: boolean; close: () => void };
  /** Set when a custom action asked to confirm first. */
  confirmCustom: {
    open: boolean;
    action: CustomBulkAction<T> | null;
    confirm: () => void;
    cancel: () => void;
  };
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
  // Starting values for the next create. Used by "add a child here" in the tree
  // view, and by anything else that opens a form already scoped to a parent.
  const [formDefaults, setFormDefaults] = useState<Record<string, unknown> | undefined>(undefined);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [bulkConfirmOpen, setBulkConfirmOpen] = useState(false);
  const [bulkArchiveOpen, setBulkArchiveOpen] = useState(false);
  const [bulkEditOpen, setBulkEditOpen] = useState(false);
  const [pendingCustom, setPendingCustom] = useState<CustomBulkAction<T> | null>(null);
  const [importOpen, setImportOpen] = useState(false);
  const tabs = useMemo(() => resource.table.tabs ?? [], [resource.table.tabs]);
  // First tab on load. A tab strip where nothing is selected reads as broken,
  // and the first tab is conventionally the unfiltered one.
  const [activeTab, setActiveTabState] = useState(() => tabs[0]?.key ?? "");
  const [showArchived, setShowArchivedState] = useState(false);
  const [liveMessage, setLiveMessage] = useState("");

  const queryClient = useQueryClient();

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
    // Tab filters, then the operator's own, then the archived flag. The
    // operator's win: picking "Unpaid" and then filtering by customer should
    // narrow the tab, not silently leave it.
    filters: {
      ...(tabs.find((t) => t.key === activeTab)?.filters ?? {}),
      ...filters,
      ...(showArchived ? { archived: "true" } : {}),
    },
    dateParams,
    dateField: resource.table.dateFilter?.field,
  });

  const rows = useMemo(() => data?.data ?? [], [data]);

  // Switching views changes which rows exist, so a selection made in the
  // other one is stale. Keeping it is how you archive something you cannot
  // see.
  // Switching tabs changes which rows exist, so a selection made under the
  // other one is stale, the same reasoning as the archived view.
  const setActiveTab = useCallback((key: string) => {
    setActiveTabState(key);
    setSelection([]);
    setPage(1);
  }, []);

  const setShowArchived = useCallback((value: boolean) => {
    setShowArchivedState(value);
    setSelection([]);
    setPage(1);
  }, []);

  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  const { mutate: deleteItem, isPending: isDeleting } = useDeleteResource(
    resource.endpoint,
    singularName,
  );
  const { mutate: runBulk, isPending: isBulkPending } = useBulkResource(
    resource.endpoint,
    pluralName,
    singularName,
  );
  // Kept as its own name because the delete confirm dialog shows a spinner
  // for delete specifically, not for any bulk action in flight.
  const isBulkDeleting = isBulkPending;

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
    setFormDefaults(undefined);
    if (isFormPage) {
      router.push("/resources/" + resource.slug + "?action=create");
    } else {
      setEditingItem(null);
      setFormOpen(true);
    }
  }, [isFormPage, router, resource.slug]);

  /**
   * Create, with some fields already filled in.
   *
   * createWith({ parent_id: id }) is how the tree view adds a child to the row
   * you clicked. In page mode the values ride along as query params, because a
   * route change is the only state that survives the navigation.
   */
  const createWith = useCallback(
    (defaults: Record<string, unknown>) => {
      if (isFormPage) {
        const params = new URLSearchParams({ action: "create" });
        for (const [key, value] of Object.entries(defaults)) {
          if (value !== undefined && value !== null) params.set(key, String(value));
        }
        router.push("/resources/" + resource.slug + "?" + params.toString());
        return;
      }
      setFormDefaults(defaults);
      setEditingItem(null);
      setFormOpen(true);
    },
    [isFormPage, router, resource.slug],
  );

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
    runBulk(
      { action: "delete", ids: selection },
      {
        onSuccess: () => {
          setBulkConfirmOpen(false);
          setSelection([]);
        },
      },
    );
  }, [runBulk, selection]);

  // ── the rest of the bulk surface ──────────────────────────────────────

  // The rows behind the selection. Custom actions get these so "email the
  // people I ticked" does not need a second round trip for data already here.
  const selectedRows = useMemo(
    () => rows.filter((row) => selection.includes(String((row as Record<string, unknown>).id))),
    [rows, selection],
  );

  const announce = useCallback((message: string) => {
    // Cleared first: setting the same string twice is not a change, and a
    // live region that has not changed says nothing. Two identical bulk
    // actions in a row would be announced once.
    setLiveMessage("");
    requestAnimationFrame(() => setLiveMessage(message));
  }, []);

  const refresh = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: [resource.endpoint] });
  }, [queryClient, resource.endpoint]);

  const bulkArchive = useCallback(() => {
    if (selection.length > 0) setBulkArchiveOpen(true);
  }, [selection]);

  const doBulkArchive = useCallback(() => {
    runBulk(
      { action: "archive", ids: selection },
      {
        onSuccess: () => {
          setBulkArchiveOpen(false);
          setSelection([]);
          announce(selection.length + " archived.");
        },
      },
    );
  }, [runBulk, selection, announce]);

  // No confirm: putting something back is not destructive, and a dialog in
  // front of an undo is a dialog nobody reads.
  const bulkRestore = useCallback(() => {
    if (selection.length === 0) return;
    runBulk(
      { action: "restore", ids: selection },
      {
        onSuccess: () => {
          setSelection([]);
          announce(selection.length + " restored.");
        },
      },
    );
  }, [runBulk, selection, announce]);

  const bulkEdit = useCallback(() => {
    if (selection.length > 0) setBulkEditOpen(true);
  }, [selection]);

  const applyBulkEdit = useCallback(
    (patch: Record<string, unknown>) => {
      runBulk(
        { action: "patch", ids: selection, patch },
        {
          onSuccess: () => {
            setBulkEditOpen(false);
            setSelection([]);
            announce(selection.length + " updated.");
          },
        },
      );
    },
    [runBulk, selection, announce],
  );

  const runCustom = useCallback(
    (action: CustomBulkAction<T>) => {
      void action.onSelect(selection, selectedRows, {
        refresh,
        clearSelection: () => setSelection([]),
        announce,
      });
    },
    [selection, selectedRows, refresh, announce],
  );

  const runBulkAction = useCallback(
    (action: CustomBulkAction<T>) => {
      if (action.confirm) {
        setPendingCustom(action);
        return;
      }
      runCustom(action);
    },
    [runCustom],
  );

  // Restore only makes sense on rows that are archived, and archive only on
  // rows that are not, so the two never appear together. Offering both is how
  // an operator ends up archiving what they meant to bring back.
  const bulkActions = useMemo(() => {
    // ["edit", "export", "delete"] rather than ["delete"] alone: a resource
    // that predates bulkActions still gets the three that work against any
    // API. Archive and restore are opt-in because they need both the column
    // and the endpoint.
    const configured = resource.table.bulkActions ?? ["edit", "export", "delete"];
    return configured.filter((action) => {
      if (action === "restore") return showArchived;
      if (action === "archive") return !showArchived;
      return true;
    });
  }, [resource.table.bulkActions, showArchived]);

  const customBulkActions = useMemo(() => {
    const all = (resource.customBulkActions ?? []) as CustomBulkAction<T>[];
    return all.filter((action) => !action.visible || action.visible(selectedRows));
  }, [resource.customBulkActions, selectedRows]);

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

    // Every stat endpoint gets whatever narrows the table, or "Total: 10,000"
    // sits above a table showing 142 matches. The archived view counts as
    // narrowing: without it the Archived tab reads "Total 9" over two rows.
    const applyViewParams = (cards: StatCard[]): StatCard[] => {
      const extra: Record<string, string> = { ...dateParams };
      if (showArchived) extra.archived = "true";
      if (Object.keys(extra).length === 0) return cards;
      return cards.map((card) => {
        if (!card.endpoint) return card;
        const sep = card.endpoint.includes("?") ? "&" : "?";
        const qs = new URLSearchParams(extra).toString();
        return { ...card, endpoint: card.endpoint + sep + qs };
      });
    };

    if (
      typeof statsConfig === "object" &&
      statsConfig !== null &&
      Array.isArray(statsConfig.cards) &&
      statsConfig.cards.length > 0
    ) {
      return applyViewParams(statsConfig.cards);
    }

    const ep = resource.endpoint;
    const defaults: StatCard[] = [
      { label: "Total", endpoint: ep + "?page_size=1", field: "meta.total", icon: resource.icon || "Package" },
      { label: "This Week", endpoint: ep + "?page_size=1&created_since=7d", field: "meta.total", icon: "TrendingUp", color: "success" },
      { label: "This Month", endpoint: ep + "?page_size=1&created_since=30d", field: "meta.total", icon: "Calendar", color: "info" },
      { label: "Updated Recently", endpoint: ep + "?page_size=1&updated_since=7d", field: "meta.total", icon: "RefreshCw" },
    ];
    return applyViewParams(defaults);
  }, [statsEnabled, statsConfig, resource.endpoint, resource.icon, dateParams, showArchived]);

  return {
    resource,

    rows,
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
    createWith,
    edit,
    view,
    remove,
    bulkRemove,
    isDeleting,
    isBulkDeleting,

    bulkActions,
    customBulkActions,
    selectedRows,
    bulkArchive,
    bulkRestore,
    bulkEdit,
    applyBulkEdit,
    runBulkAction,
    isBulkPending,
    refresh,
    announce,
    liveMessage,

    tabs,
    activeTab,
    setActiveTab,

    showArchived,
    setShowArchived,

    form: { open: formOpen, item: editingItem, defaults: formDefaults, close: closeForm },
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
    confirmBulkArchive: {
      open: bulkArchiveOpen,
      confirm: doBulkArchive,
      cancel: () => setBulkArchiveOpen(false),
    },
    bulkEditor: {
      open: bulkEditOpen,
      close: () => setBulkEditOpen(false),
    },
    confirmCustom: {
      open: pendingCustom !== null,
      action: pendingCustom,
      confirm: () => {
        if (pendingCustom) runCustom(pendingCustom);
        setPendingCustom(null);
      },
      cancel: () => setPendingCustom(null),
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
	return `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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

export type BulkOperation = "delete" | "archive" | "restore" | "patch";

export interface BulkPayload {
  action: BulkOperation;
  ids: string[];
  /** Only read for "patch". */
  patch?: Record<string, unknown>;
}

const BULK_PAST: Record<BulkOperation, string> = {
  delete: "deleted",
  archive: "archived",
  restore: "restored",
  patch: "updated",
};

/**
 * One request for the whole selection, against POST <endpoint>/bulk.
 *
 * This used to be N parallel DELETEs, which is N transactions and N audit
 * entries, and leaves a half-applied result when the eleventh fails: the
 * operator is told it failed while ten rows are already gone. The server does
 * it in one transaction now, so the answer is all or nothing.
 *
 * Takes the PLURAL label ("Invoices") because the message counts rows, and
 * reports what the server actually did rather than what was asked: archiving
 * twelve rows of which three were already archived says nine.
 */
export function useBulkResource(endpoint: string, pluralLabel?: string, singularLabel?: string) {
  const queryClient = useQueryClient();

  return useMutation({
    // ids are strings because Grit's models use UUID primary keys.
    mutationFn: async (payload: BulkPayload) => {
      try {
        const { data } = await apiClient.post(` + "`" + `${endpoint}/bulk` + "`" + `, payload);
        return { ...data, action: payload.action } as {
          data?: { affected: number; requested: number };
          action: BulkOperation;
        };
      } catch (err) {
        // No /bulk route on this endpoint. That is the normal state of an
        // upgraded project: grit upgrade replaces the admin but never
        // regenerates API handlers, so the browser gets the new code and the
        // server keeps the old routes. Falling back per row keeps the button
        // working instead of 404ing on every existing install.
        //
        // The fallback is genuinely worse: N requests, N transactions, and a
        // partial result if one fails. Run grit generate for the resource to
        // get the real endpoint.
        const status = (err as { response?: { status?: number } })?.response?.status;
        if (status !== 404) throw err;

        const results = await Promise.allSettled(
          payload.ids.map((id) =>
            payload.action === "delete"
              ? apiClient.delete(` + "`" + `${endpoint}/${id}` + "`" + `)
              : apiClient.patch(` + "`" + `${endpoint}/${id}` + "`" + `, payload.patch ?? {}),
          ),
        );
        const affected = results.filter((r) => r.status === "fulfilled").length;
        if (affected === 0) throw err;
        return {
          data: { affected, requested: payload.ids.length },
          action: payload.action,
        };
      }
    },
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      const affected = result?.data?.affected ?? 0;
      const requested = result?.data?.requested ?? affected;
      const noun = affected === 1 ? (singularLabel ?? "item") : (pluralLabel ?? "items");

      if (affected === 0) {
        // Not a success worth celebrating and not an error either. Saying
        // "0 archived" beats a green tick over a table that did not change.
        toast("Nothing to " + result.action + ": no matching rows");
        return;
      }
      const skipped = requested - affected;
      toast.success(
        affected + " " + noun + " " + BULK_PAST[result.action] +
          (skipped > 0 ? " (" + skipped + " already were)" : "")
      );
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(
        axiosErr?.response?.data?.error?.message ||
          "Bulk action failed. Nothing was changed."
      );
    },
  });
}

/**
 * Kept so existing call sites and hand-written pages keep working. Delegates
 * to the bulk endpoint rather than firing one request per row.
 *
 * @deprecated Use useBulkResource, which also archives, restores and patches.
 */
export function useBulkDeleteResource(endpoint: string, pluralLabel?: string) {
  const bulk = useBulkResource(endpoint, pluralLabel);
  return {
    ...bulk,
    mutate: (ids: string[], options?: Parameters<typeof bulk.mutate>[1]) =>
      bulk.mutate({ action: "delete", ids }, options),
  };
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
// adminTableTabs emits components/tables/table-tabs.tsx.
func adminTableTabs() string {
	return `"use client";

import { useEffect, useRef, useState } from "react";
import type { TableTab } from "@/lib/resource";
import { apiClient } from "@/lib/api-client";
import { getIcon } from "@/lib/icons";

/*
 * Filter presets as a tab strip.
 *
 * A real tablist, not a row of buttons that happen to filter. That means
 * arrow keys move between tabs and Tab leaves the group, which matters here
 * more than it looks: without roving focus a keyboard user walks through every
 * filter on the way to the table, and with six tabs that is six stops before
 * reaching the thing being filtered.
 *
 * The panel these control is the table, so the table carries the tabpanel role
 * and is labelled by the active tab.
 *
 * Counts are opt-in and arrive late. The badge is rendered only once its number
 * is known rather than showing a spinner or a zero, because a tab that says 0
 * and then says 47 is worse than a tab that said nothing for a moment. Each
 * count is one request with page_size=1, reading meta.total.
 */

export interface TableTabsProps {
  tabs: TableTab[];
  active: string;
  onChange: (key: string) => void;
  /** Endpoint the counts are fetched from, when a tab asks for one. */
  endpoint: string;
  /** Filters already applied outside the tabs, so a count matches the table. */
  baseFilters?: Record<string, string>;
}

export function TableTabs({ tabs, active, onChange, endpoint, baseFilters }: TableTabsProps) {
  const refs = useRef<(HTMLButtonElement | null)[]>([]);
  const [counts, setCounts] = useState<Record<string, number>>({});

  const wanted = tabs.filter((tab) => tab.count).map((tab) => tab.key).join(",");
  const base = JSON.stringify(baseFilters ?? {});

  useEffect(() => {
    if (!wanted) return;
    let cancelled = false;

    const load = async () => {
      const results = await Promise.all(
        tabs
          .filter((tab) => tab.count)
          .map(async (tab) => {
            const params = new URLSearchParams({
              ...(JSON.parse(base) as Record<string, string>),
              ...(tab.filters ?? {}),
              page_size: "1",
            });
            try {
              const { data } = await apiClient.get(endpoint + "?" + params.toString());
              return [tab.key, Number(data?.meta?.total ?? 0)] as const;
            } catch {
              // A count that fails is a missing badge, not a broken page.
              return null;
            }
          }),
      );
      if (cancelled) return;
      setCounts(Object.fromEntries(results.filter(Boolean) as (readonly [string, number])[]));
    };

    void load();
    return () => {
      cancelled = true;
    };
  }, [wanted, endpoint, base, tabs]);

  if (tabs.length === 0) return null;

  // Roving focus: Left and Right move between tabs, Home and End jump to the
  // ends, and Tab leaves the strip entirely.
  function onKeyDown(event: React.KeyboardEvent, index: number) {
    const last = tabs.length - 1;
    let next = index;
    if (event.key === "ArrowRight") next = index === last ? 0 : index + 1;
    else if (event.key === "ArrowLeft") next = index === 0 ? last : index - 1;
    else if (event.key === "Home") next = 0;
    else if (event.key === "End") next = last;
    else return;
    event.preventDefault();
    onChange(tabs[next].key);
    refs.current[next]?.focus();
  }

  return (
    <div
      role="tablist"
      aria-label="Filter presets"
      className="flex flex-wrap gap-1 border-b border-border px-3 pt-3"
    >
      {tabs.map((tab, index) => {
        const selected = tab.key === active;
        const Icon = tab.icon ? getIcon(tab.icon) : null;
        const count = counts[tab.key];
        return (
          <button
            key={tab.key}
            ref={(node) => {
              refs.current[index] = node;
            }}
            type="button"
            role="tab"
            id={"table-tab-" + tab.key}
            aria-selected={selected}
            aria-controls="table-panel"
            tabIndex={selected ? 0 : -1}
            onClick={() => onChange(tab.key)}
            onKeyDown={(event) => onKeyDown(event, index)}
            className={
              "inline-flex min-h-10 items-center gap-2 rounded-t-lg border-b-2 px-3 text-sm transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent " +
              (selected
                ? "border-accent font-medium text-accent"
                : "border-transparent text-text-secondary hover:bg-bg-hover hover:text-text-primary")
            }
          >
            {Icon && <Icon className="h-3.5 w-3.5" aria-hidden="true" />}
            {tab.label}
            {typeof count === "number" && (
              <span
                className={
                  "rounded-full px-1.5 py-0.5 text-xs tabular-nums " +
                  (selected ? "bg-accent/15 text-accent" : "bg-bg-tertiary text-text-secondary")
                }
              >
                {count}
                {/* The number alone is ambiguous next to a label. */}
                <span className="sr-only"> matching</span>
              </span>
            )}
          </button>
        );
      })}
    </div>
  );
}
`
}

// adminBulkActionBar emits components/tables/bulk-action-bar.tsx.
func adminBulkActionBar() string {
	return `"use client";

import { useState } from "react";
import type { ResourceDefinition, CustomBulkAction } from "@/lib/resource";
import { getIcon, Archive, ArchiveRestore, Download, Pencil, Trash2, X } from "@/lib/icons";

/*
 * The bar that appears once rows are ticked.
 *
 * Fixed to the bottom of the viewport, centred. It sat in the flow at the foot
 * of the table first, on the reasoning that a floating bar covers the rows it
 * acts on. That reasoning only holds for a table that fits on screen: with
 * twenty rows you tick something near the top, the bar appears eight hundred
 * pixels below the fold, and as far as the operator can tell nothing happened.
 * A control that responds to a selection has to be where the selection is
 * being made.
 *
 * The original worry is answered rather than ignored. It is a centred pill
 * rather than a full-width bar, so the table is visible either side of it, and
 * the page reserves space underneath while it is shown, so the last rows can
 * still be scrolled clear of it.
 *
 * It is a labelled region so it turns up in a landmark list, and its arrival is
 * announced through the page's live region, because ticking a checkbox does not
 * move focus and a bar that silently appears is a bar a keyboard user never
 * learns about.
 *
 * Delete is the only red control. If everything is red then nothing is.
 */

export interface BulkActionBarProps {
  count: number;
  /** Built-ins the resource switched on, already filtered for the view. */
  actions: string[];
  custom: CustomBulkAction[];
  pending: boolean;
  singularName: string;
  pluralName: string;
  onEdit: () => void;
  onArchive: () => void;
  onRestore: () => void;
  onDelete: () => void;
  onExport: () => void;
  onCustom: (action: CustomBulkAction) => void;
  onClear: () => void;
}

export function BulkActionBar({
  count,
  actions,
  custom,
  pending,
  singularName,
  pluralName,
  onEdit,
  onArchive,
  onRestore,
  onDelete,
  onExport,
  onCustom,
  onClear,
}: BulkActionBarProps) {
  if (count === 0) return null;

  const noun = count === 1 ? singularName.toLowerCase() : pluralName.toLowerCase();
  const base =
    "inline-flex min-h-9 items-center gap-1.5 rounded-lg border px-3 text-sm font-medium transition-colors disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent";
  const neutral = base + " border-border bg-bg-secondary text-text-primary hover:bg-bg-hover";
  const danger = base + " border-danger/40 bg-transparent text-danger hover:bg-danger/10";

  return (
    <section
      aria-label="Bulk actions"
      className="fixed bottom-6 left-1/2 z-40 flex max-w-[calc(100vw-2rem)] -translate-x-1/2 flex-wrap items-center gap-3 rounded-xl border border-border bg-bg-elevated px-4 py-2.5 shadow-2xl shadow-black/25"
    >
      <p className="text-sm font-medium text-text-primary">
        {count} {noun} selected
      </p>

      <div className="flex flex-wrap items-center gap-2">
        {actions.includes("edit") && (
          <button type="button" onClick={onEdit} disabled={pending} className={neutral}>
            <Pencil className="h-3.5 w-3.5" aria-hidden="true" />
            Edit
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}

        {actions.includes("archive") && (
          <button type="button" onClick={onArchive} disabled={pending} className={neutral}>
            <Archive className="h-3.5 w-3.5" aria-hidden="true" />
            Archive
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}

        {actions.includes("restore") && (
          <button type="button" onClick={onRestore} disabled={pending} className={neutral}>
            <ArchiveRestore className="h-3.5 w-3.5" aria-hidden="true" />
            Restore
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}

        {actions.includes("export") && (
          <button type="button" onClick={onExport} disabled={pending} className={neutral}>
            <Download className="h-3.5 w-3.5" aria-hidden="true" />
            Export
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}

        {custom.map((action) => {
          const Icon = action.icon ? getIcon(action.icon) : null;
          return (
            <button
              key={action.key}
              type="button"
              onClick={() => onCustom(action)}
              disabled={pending}
              className={action.variant === "danger" ? danger : neutral}
            >
              {Icon && <Icon className="h-3.5 w-3.5" aria-hidden="true" />}
              {action.label}
              <span className="sr-only"> for the {count} selected {noun}</span>
            </button>
          );
        })}

        {actions.includes("delete") && (
          <button type="button" onClick={onDelete} disabled={pending} className={danger}>
            <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
            Delete
            <span className="sr-only"> the {count} selected {noun}</span>
          </button>
        )}
      </div>

      <button
        type="button"
        onClick={onClear}
        className="inline-flex min-h-9 items-center gap-1.5 rounded-lg px-3 text-sm text-text-secondary hover:bg-bg-hover hover:text-text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
      >
        <X className="h-3.5 w-3.5" aria-hidden="true" />
        Clear selection
      </button>
    </section>
  );
}
`
}

// adminBulkEditModal emits components/tables/bulk-edit-modal.tsx.
func adminBulkEditModal() string {
	return `"use client";

import { useMemo, useState } from "react";
import type { FieldDefinition, ResourceDefinition } from "@/lib/resource";
import { Loader2, X } from "@/lib/icons";

/*
 * Bulk edit: one field, one value, written to every selected row.
 *
 * Deliberately one field rather than a whole form. Bulk editing every field at
 * once means deciding what an empty input means, and there is no good answer:
 * "clear it" destroys data the operator never looked at, and "ignore it" makes
 * it impossible to clear anything. One field sidesteps the question entirely,
 * and it is what the job actually is nine times out of ten: set the status,
 * change the owner, move the category.
 *
 * Only fields that can carry the same value for many rows are offered. A
 * unique field is excluded, because writing one SKU to forty products is
 * either a constraint violation or, worse, not one.
 */

const UNSUITABLE: FieldDefinition["type"][] = [
  "file",
  "files",
  "image",
  "images",
  "video",
  "videos",
  "line-items",
];

export interface BulkEditModalProps {
  resource: ResourceDefinition;
  count: number;
  pending: boolean;
  onApply: (patch: Record<string, unknown>) => void;
  onClose: () => void;
}

export function BulkEditModal({ resource, count, pending, onApply, onClose }: BulkEditModalProps) {
  const fields = useMemo(
    () =>
      resource.form.fields.filter(
        (field) => !UNSUITABLE.includes(field.type) && !field.unique,
      ),
    [resource.form.fields],
  );

  const [key, setKey] = useState(fields[0]?.key ?? "");
  const [value, setValue] = useState<string>("");

  const field = fields.find((f) => f.key === key);
  const noun = count === 1 ? resource.label?.singular ?? resource.name : resource.label?.plural ?? resource.slug;

  function submit(event: React.FormEvent) {
    event.preventDefault();
    if (!field) return;
    // Cast at the boundary: an <input> always hands back a string, and the
    // API expects the column's real type.
    let parsed: unknown = value;
    if (field.type === "number") parsed = value === "" ? null : Number(value);
    if (field.type === "toggle" || field.type === "checkbox") parsed = value === "true";
    onApply({ [field.key]: parsed });
  }

  const inputClass =
    "w-full rounded-lg border border-border bg-bg-primary px-3 py-2 text-sm text-text-primary focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent";

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" onClick={onClose} />
      <form
        onSubmit={submit}
        className="relative w-full max-w-md rounded-xl border border-border bg-bg-secondary shadow-xl"
      >
        <div className="flex items-center justify-between border-b border-border px-5 py-4">
          <h2 className="text-sm font-semibold text-text-primary">
            Edit {count} {noun.toLowerCase()}
          </h2>
          <button
            type="button"
            onClick={onClose}
            className="rounded-md p-1 text-text-secondary hover:bg-bg-hover hover:text-text-primary focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            <X className="h-4 w-4" aria-hidden="true" />
            <span className="sr-only">Close</span>
          </button>
        </div>

        <div className="space-y-4 px-5 py-4">
          {fields.length === 0 ? (
            <p className="text-sm text-text-secondary">
              No fields on this resource can be set in bulk.
            </p>
          ) : (
            <>
              <label className="block">
                <span className="mb-1 block text-xs font-medium text-text-secondary">Field</span>
                <select
                  value={key}
                  onChange={(e) => {
                    setKey(e.target.value);
                    setValue("");
                  }}
                  className={inputClass}
                >
                  {fields.map((f) => (
                    <option key={f.key} value={f.key}>
                      {f.label}
                    </option>
                  ))}
                </select>
              </label>

              {field && (
                <label className="block">
                  <span className="mb-1 block text-xs font-medium text-text-secondary">
                    {field.label}
                  </span>
                  {field.type === "select" || field.type === "radio" ? (
                    <select
                      value={value}
                      onChange={(e) => setValue(e.target.value)}
                      className={inputClass}
                    >
                      <option value="">Choose...</option>
                      {(field.options ?? []).map((opt) => (
                        <option key={opt.value} value={opt.value}>
                          {opt.label}
                        </option>
                      ))}
                    </select>
                  ) : field.type === "toggle" || field.type === "checkbox" ? (
                    <select
                      value={value}
                      onChange={(e) => setValue(e.target.value)}
                      className={inputClass}
                    >
                      <option value="">Choose...</option>
                      <option value="true">Yes</option>
                      <option value="false">No</option>
                    </select>
                  ) : field.type === "textarea" || field.type === "richtext" ? (
                    <textarea
                      rows={4}
                      value={value}
                      onChange={(e) => setValue(e.target.value)}
                      className={inputClass}
                    />
                  ) : (
                    <input
                      type={field.type === "number" ? "number" : field.type === "date" ? "date" : "text"}
                      value={value}
                      onChange={(e) => setValue(e.target.value)}
                      className={inputClass}
                    />
                  )}
                </label>
              )}

              <p className="rounded-lg bg-warning/10 px-3 py-2 text-xs text-text-secondary">
                This writes the same value to all {count} selected {noun.toLowerCase()}. It cannot be
                undone in one step.
              </p>
            </>
          )}
        </div>

        <div className="flex justify-end gap-2 border-t border-border px-5 py-4">
          <button
            type="button"
            onClick={onClose}
            className="inline-flex min-h-10 items-center rounded-lg border border-border px-4 text-sm font-medium text-text-secondary hover:bg-bg-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={pending || fields.length === 0 || value === ""}
            className="inline-flex min-h-10 items-center gap-2 rounded-lg bg-accent px-4 text-sm font-semibold text-accent-fg hover:bg-accent-hover disabled:opacity-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
          >
            {pending && <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />}
            Apply to {count}
          </button>
        </div>
      </form>
    </div>
  );
}
`
}

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
      // grit:cols:auto-start
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
      // grit:cols:auto-end
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
    // No "archive": the scaffold's own models predate archived_at. A
    // generated resource gets the column and the full set.
    bulkActions: ["edit", "export", "delete"],
    defaultSort: { key: "created_at", direction: "desc" },
    pageSize: 20,
  },

  form: {
    layout: "single",
    fields: [
      // grit:fields:auto-start
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
      // grit:fields:auto-end
    ],
  },
}, custom);
`
}

// adminBlogsPage returns the blogs resource page.
func adminBlogsPage() string {
	return `"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { blogsResource } from "@/resources/blogs/blogs";

export default function BlogsPage() {
  return <ResourcePage resource={blogsResource} />;
}
`
}
