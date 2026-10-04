// Resource Definition Types — The foundation of Grit Admin Panel
// Define resources with defineResource() and get full CRUD pages automatically.

import type { ComponentType, ReactNode } from "react";

// ─── Column Definitions ─────────────────────────────────────────────

export type ColumnFormat = "text" | "badge" | "currency" | "money" | "date" | "relative" | "boolean" | "image" | "video" | "file" | "files" | "link" | "email" | "color" | "richtext" | "user" | "domain" | "tel" | "country" | "percent" | "rating" | "time" | "json";

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
  /** How many stars a "rating" column is out of. Defaults to 5. */
  ratingMax?: number;
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

export type FieldType = "text" | "textarea" | "number" | "money" | "select" | "date" | "datetime" | "toggle" | "checkbox" | "checkbox-group" | "radio" | "richtext" | "image" | "images" | "video" | "videos" | "file" | "files" | "relationship-select" | "multi-relationship-select" | "line-items" | "email" | "url" | "domain" | "tel" | "country" | "color" | "percent" | "rating" | "time" | "json";

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

  /** Currency codes a money field offers. Unset shows a short default list.
   *  A shop that only ever trades in one currency should name that one: the
   *  picker then has a single option and nobody can pick the wrong one. */
  currencies?: string[];
  /** Which of them a new record starts on. Unset uses the first. */
  defaultCurrency?: string;

  /** tel only: the ISO 3166-1 country the picker starts on and a local number
   *  is read in, e.g. "UG". Set by the generator from phone:tel:UG. Unset uses
   *  the browser's locale. A rating field's stars come from max (default 5). */
  defaultCountry?: string;

  // v3.103.0 — a visible field with a small "Generate" button in its label
  // row. Unlike an auto field (which is server-filled and hidden from the
  // form), this keeps the input visible and editable; clicking Generate runs
  // YOUR function with the current form values and fills the field with what it
  // returns (sync or async — e.g. call an endpoint, derive from another field,
  // mint a code). text / number fields only. You define this by hand in the
  // resource definition; the generator never emits it.
  generate?: (values: Record<string, unknown>) => string | number | Promise<string | number>;

  /**
   * Shown, never sent. The value comes from the record, a column the server
   * computes or a status another process owns, and the form displays it the way
   * the table would. It is left out of what the form submits, so saving cannot
   * write the displayed value back over the server's.
   */
  readOnly?: boolean;
  /**
   * Derived from the other values in the form and recomputed as they change: a
   * line total, a balance, the difference between debits and credits. Implies
   * readOnly. Formatted by the field's type, so type: "money" shows money.
   */
  compute?: (values: Record<string, unknown>) => unknown;

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
  isFetching?: boolean;
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
  /** Whether the resource's table allows this action. The detail page asks
   *  before offering Edit or Delete, the same question the list asks. */
  can: (action: TableAction) => boolean;
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
