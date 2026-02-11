package scaffold

import "fmt"

// adminResourceTypes returns the resource type system (lib/resource.ts).
func adminResourceTypes() string {
	return `// Resource Definition Types — The foundation of Grit Admin Panel
// Define resources with defineResource() and get full CRUD pages automatically.

// ─── Column Definitions ─────────────────────────────────────────────

export type ColumnFormat = "text" | "badge" | "currency" | "date" | "relative" | "boolean" | "image";

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

export type TableAction = "create" | "edit" | "delete" | "export";
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

export type FieldType = "text" | "textarea" | "number" | "select" | "date" | "datetime" | "toggle" | "checkbox" | "radio" | "image";

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
}

export interface FormDefinition {
  fields: FieldDefinition[];
  layout?: "single" | "two-column";
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
  table: TableDefinition;
  form: FormDefinition;
  dashboard?: DashboardDefinition;
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
      actions: config.table.actions ?? ["create", "edit", "delete"],
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
// grit:resources

import type { ResourceDefinition } from "@/lib/resource";

export const resources: ResourceDefinition[] = [
  usersResource,
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
      { key: "name", label: "Name", sortable: true, searchable: true },
      { key: "email", label: "Email", sortable: true, searchable: true },
      {
        key: "role",
        label: "Role",
        sortable: true,
        format: "badge",
        badge: {
          admin: { color: "accent", label: "Admin" },
          editor: { color: "info", label: "Editor" },
          user: { color: "muted", label: "User" },
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
          { label: "Admin", value: "admin" },
          { label: "Editor", value: "editor" },
          { label: "User", value: "user" },
        ],
      },
      { key: "active", label: "Status", type: "boolean" },
    ],
    searchable: true,
    searchPlaceholder: "Search by name or email...",
    actions: ["create", "edit", "delete"],
    bulkActions: ["delete"],
    defaultSort: { key: "created_at", direction: "desc" },
    pageSize: 20,
  },

  form: {
    layout: "two-column",
    fields: [
      {
        key: "name",
        label: "Full Name",
        type: "text",
        required: true,
        placeholder: "Enter full name",
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
        key: "role",
        label: "Role",
        type: "select",
        required: true,
        options: [
          { label: "Admin", value: "admin" },
          { label: "Editor", value: "editor" },
          { label: "User", value: "user" },
        ],
        defaultValue: "user",
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
import type { ResourceDefinition } from "@/lib/resource";
import { useResource, useDeleteResource, useBulkDeleteResource } from "@/hooks/use-resource";
import { DataTable } from "@/components/tables/data-table";
import { TableToolbar } from "@/components/tables/table-toolbar";
import { TablePagination } from "@/components/tables/table-pagination";
import { TableFilters } from "@/components/tables/table-filters";
import { FormModal } from "@/components/forms/form-modal";

interface ResourcePageProps {
  resource: ResourceDefinition;
}

export function ResourcePage({ resource }: ResourcePageProps) {
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

  // Data fetching
  const { data, isLoading } = useResource(resource.endpoint, {
    page,
    pageSize,
    search,
    sortBy,
    sortOrder,
    filters,
  });

  const { mutate: deleteItem } = useDeleteResource(resource.endpoint);
  const { mutate: bulkDelete } = useBulkDeleteResource(resource.endpoint);

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

  const handleEdit = useCallback((item: Record<string, unknown>) => {
    setEditingItem(item);
    setFormOpen(true);
  }, []);

  const handleCreate = useCallback(() => {
    setEditingItem(null);
    setFormOpen(true);
  }, []);

  const handleDelete = useCallback(
    (id: number) => {
      if (confirm(` + "`" + `Delete this ${resource.label?.singular ?? resource.name}?` + "`" + `)) {
        deleteItem(id);
      }
    },
    [deleteItem, resource]
  );

  const handleBulkDelete = useCallback(() => {
    if (
      selectedRows.length > 0 &&
      confirm(` + "`" + `Delete ${selectedRows.length} ${resource.label?.plural ?? resource.slug}?` + "`" + `)
    ) {
      bulkDelete(selectedRows);
      setSelectedRows([]);
    }
  }, [bulkDelete, selectedRows, resource]);

  const handleFormClose = useCallback(() => {
    setFormOpen(false);
    setEditingItem(null);
  }, []);

  const actions = resource.table.actions ?? ["create", "edit", "delete"];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">
          {resource.label?.plural ?? resource.slug}
        </h1>
        <p className="text-text-secondary mt-1">
          Manage {(resource.label?.plural ?? resource.slug).toLowerCase()}
        </p>
      </div>

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

      {formOpen && (
        <FormModal
          resource={resource}
          item={editingItem}
          onClose={handleFormClose}
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
func adminUseResource() string {
	return `import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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

export function useResourceItem<T = Record<string, unknown>>(endpoint: string, id: number | null) {
  return useQuery<{ data: T }>({
    queryKey: [endpoint, id],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `${endpoint}/${id}` + "`" + `);
      return data;
    },
    enabled: id !== null,
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
    },
  });
}
`
}

// adminDashboardPage returns the enhanced dashboard page.
func adminDashboardPage() string {
	return fmt.Sprintf(`"use client";

import { resources } from "@/resources";
import { StatsCard } from "@/components/widgets/stats-card";
import { WidgetGrid } from "@/components/widgets/widget-grid";

export default function AdminDashboard() {
  // Collect all dashboard widgets from registered resources
  const allWidgets = resources.flatMap((r) => r.dashboard?.widgets ?? []);

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Dashboard</h1>
        <p className="text-text-secondary mt-1">Overview of your application</p>
      </div>

      {allWidgets.length > 0 ? (
        <WidgetGrid widgets={allWidgets} />
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard label="Total Resources" value="—" icon="Database" color="accent" />
          <StatsCard label="Registered" value={String(resources.length)} icon="Layers" color="success" />
        </div>
      )}

      <div className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-xl border border-border bg-bg-secondary p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Quick Actions</h2>
          <div className="grid gap-3 sm:grid-cols-2">
            {resources.map((r) => (
              <a
                key={r.slug}
                href={%s/resources/${r.slug}%s}
                className="rounded-lg border border-border bg-bg-tertiary p-4 hover:bg-bg-hover transition-colors"
              >
                <h3 className="font-medium text-foreground">{r.label?.plural ?? r.name}</h3>
                <p className="text-xs text-text-muted mt-1">
                  Manage {(r.label?.plural ?? r.slug).toLowerCase()}
                </p>
              </a>
            ))}
            <a
              href="http://localhost:8080/studio"
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-lg border border-border bg-bg-tertiary p-4 hover:bg-bg-hover transition-colors"
            >
              <h3 className="font-medium text-foreground">GORM Studio</h3>
              <p className="text-xs text-text-muted mt-1">Browse database</p>
            </a>
          </div>
        </div>

        <div className="rounded-xl border border-border bg-bg-secondary p-6">
          <h2 className="text-lg font-semibold text-foreground mb-4">Recent Activity</h2>
          <div className="space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="flex items-center gap-3 text-sm">
                <div className="h-2 w-2 rounded-full bg-accent" />
                <span className="text-text-secondary">Activity placeholder #{i}</span>
                <span className="ml-auto text-text-muted text-xs">Just now</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
`, "`", "`")
}
