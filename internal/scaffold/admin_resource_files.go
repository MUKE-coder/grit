package scaffold

import "fmt"

// adminResourceTypes returns the resource type system (lib/resource.ts).
func adminResourceTypes() string {
	return `// Resource Definition Types — The foundation of Grit Admin Panel
// Define resources with defineResource() and get full CRUD pages automatically.

// ─── Column Definitions ─────────────────────────────────────────────

export type ColumnFormat = "text" | "badge" | "currency" | "date" | "relative" | "boolean" | "image" | "video" | "link" | "email" | "color" | "richtext";

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
}

// ─── Filter Definitions ─────────────────────────────────────────────

export type FilterType = "select" | "date-range" | "number-range" | "boolean";

export interface FilterOption {
  label: string;
  value: string;
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

export interface TableDefinition {
  columns: ColumnDefinition[];
  filters?: FilterDefinition[];
  searchable?: boolean;
  searchPlaceholder?: string;
  actions?: TableAction[];
  bulkActions?: BulkAction[];
  defaultSort?: { key: string; direction: "asc" | "desc" };
  pageSize?: number;
}

// ─── Form Field Definitions ─────────────────────────────────────────

export type FieldType = "text" | "textarea" | "number" | "select" | "date" | "datetime" | "toggle" | "checkbox" | "radio" | "richtext" | "image" | "images" | "video" | "videos" | "file" | "files" | "relationship-select" | "multi-relationship-select";

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
}

export interface StepDefinition {
  title: string;
  description?: string;
  fields: string[];
}

export interface FormDefinition {
  fields: FieldDefinition[];
  layout?: "single" | "two-column";
  steps?: StepDefinition[];
  fieldsPerStep?: number;
  stepVariant?: "horizontal" | "vertical";
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
  widgets: WidgetDefinition[];
}

// ─── Resource Definition ────────────────────────────────────────────

export interface ResourceDefinition {
  name: string;
  slug: string;
  endpoint: string;
  icon: string;
  label?: { singular: string; plural: string };
  formView?: "modal" | "page" | "modal-steps" | "page-steps";
  table: TableDefinition;
  form: FormDefinition;
  dashboard?: DashboardDefinition;
  stats?: StatsConfig | boolean;
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

export function defineResource(config: ResourceDefinition): ResourceDefinition {
  return {
    ...config,
    label: config.label ?? {
      singular: config.name,
      plural: config.slug.charAt(0).toUpperCase() + config.slug.slice(1),
    },
    table: {
      ...config.table,
      pageSize: config.table.pageSize ?? 20,
      actions: config.table.actions ?? ["create", "view", "edit", "delete"],
      searchable: config.table.searchable ?? true,
    },
    form: {
      ...config.form,
      layout: config.form.layout ?? "single",
    },
  };
}
`
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

export const usersResource = defineResource({
  name: "User",
  slug: "users",
  endpoint: "/api/users",
  icon: "Users",
  label: { singular: "User", plural: "Users" },

  table: {
    columns: [
      { key: "id", label: "ID", sortable: true, width: "80px" },
      { key: "first_name", label: "First Name", sortable: true, searchable: true },
      { key: "last_name", label: "Last Name", sortable: true, searchable: true },
      { key: "email", label: "Email", sortable: true, searchable: true },
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
});
`
}

// adminResourcePage returns the generic resource page component.
func adminResourcePage() string {
	return `"use client";

import { useState, useCallback, useMemo } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import dynamic from "next/dynamic";
import type { ResourceDefinition } from "@/lib/resource";
import { useResource, useDeleteResource, useBulkDeleteResource } from "@/hooks/use-resource";
import { PageHeader, type StatCard } from "@/components/layout/page-header";
import { DataTable } from "@/components/tables/data-table";
import { TableToolbar } from "@/components/tables/table-toolbar";
import { TablePagination } from "@/components/tables/table-pagination";
import { TableFilters } from "@/components/tables/table-filters";

// Lazy-load modal/form components — they are only shown conditionally and
// would otherwise inflate the initial page bundle for every admin resource.
const FormModal = dynamic(() =>
  import("@/components/forms/form-modal").then((m) => m.FormModal)
);
const FormPage = dynamic(() =>
  import("@/components/forms/form-page").then((m) => m.FormPage)
);
const FormModalSteps = dynamic(() =>
  import("@/components/forms/form-modal-steps").then((m) => m.FormModalSteps)
);
const FormPageSteps = dynamic(() =>
  import("@/components/forms/form-page-steps").then((m) => m.FormPageSteps)
);
const ViewModal = dynamic(() =>
  import("@/components/resource/view-modal").then((m) => m.ViewModal)
);
const ConfirmModal = dynamic(() =>
  import("@/components/ui/confirm-modal").then((m) => m.ConfirmModal)
);

interface ResourcePageProps {
  resource: ResourceDefinition;
}

export function ResourcePage({ resource }: ResourcePageProps) {
  const router = useRouter();
  const searchParams = useSearchParams();
  const isFormPage = resource.formView === "page" || resource.formView === "page-steps";
  const isSteps = resource.formView === "modal-steps" || resource.formView === "page-steps";
  const formAction = searchParams.get("action");

  // If formView is "page" or "page-steps" and we have an action param, show the form page
  if (isFormPage && (formAction === "create" || formAction === "edit")) {
    return isSteps ? <FormPageSteps resource={resource} /> : <FormPage resource={resource} />;
  }

  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(resource.table.pageSize ?? 20);
  const [search, setSearch] = useState("");
  const [sortBy, setSortBy] = useState(resource.table.defaultSort?.key ?? "");
  const [sortOrder, setSortOrder] = useState<"asc" | "desc">(
    resource.table.defaultSort?.direction ?? "desc"
  );
  const [selectedRows, setSelectedRows] = useState<number[]>([]);
  const [filters, setFilters] = useState<Record<string, string>>({});
  const [hiddenColumns, setHiddenColumns] = useState<string[]>([]);

  // Form modal state
  const [formOpen, setFormOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<Record<string, unknown> | null>(null);

  // View modal state
  const [viewingItem, setViewingItem] = useState<Record<string, unknown> | null>(null);

  // Confirm modal state
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [deletingId, setDeletingId] = useState<number | null>(null);
  const [bulkConfirmOpen, setBulkConfirmOpen] = useState(false);

  // Data fetching
  const { data, isLoading } = useResource(resource.endpoint, {
    page,
    pageSize,
    search,
    sortBy,
    sortOrder,
    filters,
  });

  const { mutate: deleteItem, isPending: isDeleting } = useDeleteResource(resource.endpoint);
  const { mutate: bulkDelete, isPending: isBulkDeleting } = useBulkDeleteResource(resource.endpoint);

  // Visible columns
  const visibleColumns = useMemo(
    () => resource.table.columns.filter((col) => !col.hidden && !hiddenColumns.includes(col.key)),
    [resource.table.columns, hiddenColumns]
  );

  // Handlers
  const handleSort = useCallback(
    (key: string) => {
      if (sortBy === key) {
        setSortOrder((prev) => (prev === "asc" ? "desc" : "asc"));
      } else {
        setSortBy(key);
        setSortOrder("asc");
      }
      setPage(1);
    },
    [sortBy]
  );

  const handleSearch = useCallback((value: string) => {
    setSearch(value);
    setPage(1);
  }, []);

  const handleFilter = useCallback((key: string, value: string) => {
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

  const handleView = useCallback((item: Record<string, unknown>) => {
    setViewingItem(item);
  }, []);

  const handleEdit = useCallback((item: Record<string, unknown>) => {
    if (isFormPage) {
      router.push(` + "`" + `/resources/${resource.slug}?action=edit&edit=${item.id}` + "`" + `);
    } else {
      setEditingItem(item);
      setFormOpen(true);
    }
  }, [isFormPage, router, resource.slug]);

  const handleCreate = useCallback(() => {
    if (isFormPage) {
      router.push(` + "`" + `/resources/${resource.slug}?action=create` + "`" + `);
    } else {
      setEditingItem(null);
      setFormOpen(true);
    }
  }, [isFormPage, router, resource.slug]);

  const handleDelete = useCallback((id: number) => {
    setDeletingId(id);
    setConfirmOpen(true);
  }, []);

  const confirmDelete = useCallback(() => {
    if (deletingId !== null) {
      deleteItem(deletingId, {
        onSuccess: () => {
          setConfirmOpen(false);
          setDeletingId(null);
        },
      });
    }
  }, [deleteItem, deletingId]);

  const handleBulkDelete = useCallback(() => {
    if (selectedRows.length > 0) {
      setBulkConfirmOpen(true);
    }
  }, [selectedRows]);

  const confirmBulkDelete = useCallback(() => {
    bulkDelete(selectedRows, {
      onSuccess: () => {
        setBulkConfirmOpen(false);
        setSelectedRows([]);
      },
    });
  }, [bulkDelete, selectedRows]);

  const handleFormClose = useCallback(() => {
    setFormOpen(false);
    setEditingItem(null);
  }, []);

  const actions = resource.table.actions ?? ["create", "view", "edit", "delete"];
  const singularName = resource.label?.singular ?? resource.name;
  const pluralName = resource.label?.plural ?? resource.slug;

  // Build stats cards: either from resource.stats.cards config, or auto-defaults.
  // Set resource.stats = false (or { enabled: false }) to disable stats on this page.
  const statsConfig = resource.stats;
  const statsEnabled =
    statsConfig === undefined ||
    statsConfig === true ||
    (typeof statsConfig === "object" && statsConfig !== null && statsConfig.enabled !== false);

  const statsCards: StatCard[] | undefined = useMemo(() => {
    if (!statsEnabled) return undefined;
    // Explicit cards override auto-defaults
    if (
      typeof statsConfig === "object" &&
      statsConfig !== null &&
      Array.isArray(statsConfig.cards) &&
      statsConfig.cards.length > 0
    ) {
      return statsConfig.cards;
    }
    // Auto-defaults: 4 cards based on the resource endpoint
    const ep = resource.endpoint;
    return [
      { label: "Total", endpoint: ` + "`" + `${ep}?page_size=1` + "`" + `, field: "meta.total", icon: resource.icon || "Package" },
      { label: "This Week", endpoint: ` + "`" + `${ep}?page_size=1&created_since=7d` + "`" + `, field: "meta.total", icon: "TrendingUp", color: "success" },
      { label: "This Month", endpoint: ` + "`" + `${ep}?page_size=1&created_since=30d` + "`" + `, field: "meta.total", icon: "Calendar", color: "info" },
      { label: "Updated Recently", endpoint: ` + "`" + `${ep}?page_size=1&updated_since=7d` + "`" + `, field: "meta.total", icon: "RefreshCw" },
    ];
  }, [statsEnabled, statsConfig, resource.endpoint, resource.icon]);

  const headerActions = actions.includes("create") ? (
    <button
      onClick={handleCreate}
      className="inline-flex items-center gap-1.5 rounded-lg bg-accent px-3 h-9 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
    >
      <span className="text-base leading-none">+</span>
      New {singularName}
    </button>
  ) : undefined;

  return (
    <div>
      <PageHeader
        title={pluralName}
        description={` + "`" + `Manage ${pluralName.toLowerCase()}` + "`" + `}
        actions={headerActions}
        stats={statsCards}
      />

      <div className="rounded-xl border border-border bg-bg-secondary">
        <TableToolbar
          resource={resource}
          search={search}
          onSearch={handleSearch}
          selectedCount={selectedRows.length}
          onBulkDelete={handleBulkDelete}
          onCreate={actions.includes("create") ? handleCreate : undefined}
          allColumns={resource.table.columns}
          hiddenColumns={hiddenColumns}
          onToggleColumn={(key) =>
            setHiddenColumns((prev) =>
              prev.includes(key) ? prev.filter((k) => k !== key) : [...prev, key]
            )
          }
          data={data?.data}
        />

        {resource.table.filters && resource.table.filters.length > 0 && (
          <TableFilters
            filters={resource.table.filters}
            values={filters}
            onChange={handleFilter}
          />
        )}

        <DataTable
          columns={visibleColumns}
          data={data?.data ?? []}
          isLoading={isLoading}
          sortBy={sortBy}
          sortOrder={sortOrder}
          onSort={handleSort}
          selectedRows={selectedRows}
          onSelectRows={setSelectedRows}
          onView={actions.includes("view") ? handleView : undefined}
          onEdit={actions.includes("edit") ? handleEdit : undefined}
          onDelete={actions.includes("delete") ? handleDelete : undefined}
        />

        <TablePagination
          page={page}
          pageSize={pageSize}
          total={data?.meta?.total ?? 0}
          totalPages={data?.meta?.pages ?? 1}
          onPageChange={setPage}
          onPageSizeChange={(size) => {
            setPageSize(size);
            setPage(1);
          }}
        />
      </div>

      {!isFormPage && formOpen && (
        isSteps ? (
          <FormModalSteps
            resource={resource}
            item={editingItem}
            onClose={handleFormClose}
          />
        ) : (
          <FormModal
            resource={resource}
            item={editingItem}
            onClose={handleFormClose}
          />
        )
      )}

      {viewingItem && (
        <ViewModal
          resource={resource}
          item={viewingItem}
          onClose={() => setViewingItem(null)}
          onEdit={actions.includes("edit") ? handleEdit : undefined}
        />
      )}

      <ConfirmModal
        open={confirmOpen}
        onConfirm={confirmDelete}
        onCancel={() => { setConfirmOpen(false); setDeletingId(null); }}
        title={` + "`" + `Delete ${singularName}` + "`" + `}
        description={` + "`" + `Are you sure you want to delete this ${singularName.toLowerCase()}? This action cannot be undone.` + "`" + `}
        confirmLabel="Delete"
        variant="danger"
        loading={isDeleting}
      />

      <ConfirmModal
        open={bulkConfirmOpen}
        onConfirm={confirmBulkDelete}
        onCancel={() => setBulkConfirmOpen(false)}
        title={` + "`" + `Delete ${selectedRows.length} ${pluralName.toLowerCase()}` + "`" + `}
        description={` + "`" + `Are you sure you want to delete ${selectedRows.length} ${pluralName.toLowerCase()}? This action cannot be undone.` + "`" + `}
        confirmLabel="Delete All"
        variant="danger"
        loading={isBulkDeleting}
      />
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
  const { page = 1, pageSize = 20, search, sortBy, sortOrder, filters } = params;

  return useQuery<PaginatedResponse<T>>({
    queryKey: [endpoint, { page, pageSize, search, sortBy, sortOrder, filters }],
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

      const { data } = await apiClient.get(` + "`" + `${endpoint}?${searchParams}` + "`" + `);
      return data;
    },
  });
}

export function useResourceItem<T = Record<string, unknown>>(
  endpoint: string,
  id: number,
  options?: { enabled?: boolean }
) {
  return useQuery<{ data: T }>({
    queryKey: [endpoint, id],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `${endpoint}/${id}` + "`" + `);
      return data;
    },
    enabled: (options?.enabled ?? true) && id > 0,
  });
}

export function useCreateResource(endpoint: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (body: Record<string, unknown>) => {
      const { data } = await apiClient.post(endpoint, body);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success("Created successfully");
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(axiosErr?.response?.data?.error?.message || "Failed to create");
    },
  });
}

export function useUpdateResource(endpoint: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ id, body }: { id: number; body: Record<string, unknown> }) => {
      const { data } = await apiClient.put(` + "`" + `${endpoint}/${id}` + "`" + `, body);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success("Updated successfully");
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(axiosErr?.response?.data?.error?.message || "Failed to update");
    },
  });
}

export function useDeleteResource(endpoint: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (id: number) => {
      await apiClient.delete(` + "`" + `${endpoint}/${id}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success("Deleted successfully");
    },
    onError: (err: unknown) => {
      const axiosErr = err as { response?: { data?: { error?: { message?: string } } } };
      toast.error(axiosErr?.response?.data?.error?.message || "Failed to delete");
    },
  });
}

export function useBulkDeleteResource(endpoint: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (ids: number[]) => {
      await Promise.all(ids.map((id) => apiClient.delete(` + "`" + `${endpoint}/${id}` + "`" + `)));
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [endpoint] });
      toast.success("Deleted successfully");
    },
    onError: () => {
      toast.error("Failed to delete some items");
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
              href="http://localhost:8080/studio"
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
              href="http://localhost:8080/api/health"
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

export const blogsResource = defineResource({
  name: "Blog",
  slug: "blogs",
  endpoint: "/api/admin/blogs",
  icon: "FileText",
  label: { singular: "Blog", plural: "Blogs" },

  table: {
    columns: [
      { key: "id", label: "ID", sortable: true, width: "80px" },
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
});
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
