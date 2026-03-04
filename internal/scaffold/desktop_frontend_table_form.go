package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDesktopFrontendTableFormFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "frontend", "src", "components", "data-table.tsx"):   desktopDataTable(),
		filepath.Join(root, "frontend", "src", "components", "form-builder.tsx"): desktopFormBuilder(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func desktopDataTable() string {
	return `import { ChevronUp, ChevronDown, ChevronLeft, ChevronRight, Inbox } from "lucide-react";
import { cn } from "@/lib/utils";

interface Column {
  key: string;
  label: string;
  sortable?: boolean;
  render?: (value: any, row: any) => React.ReactNode;
}

interface DataTableProps {
  columns: Column[];
  data: any[];
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number) => void;
  onSort?: (key: string, direction: "asc" | "desc") => void;
  sortKey?: string;
  sortDirection?: "asc" | "desc";
  isLoading?: boolean;
  onRowClick?: (row: any) => void;
  actions?: (row: any) => React.ReactNode;
}

export default function DataTable({
  columns,
  data,
  total,
  page,
  pageSize,
  onPageChange,
  onSort,
  sortKey,
  sortDirection,
  isLoading,
  onRowClick,
  actions,
}: DataTableProps) {
  const totalPages = Math.ceil(total / pageSize);
  const startItem = (page - 1) * pageSize + 1;
  const endItem = Math.min(page * pageSize, total);

  const handleSort = (col: Column) => {
    if (!col.sortable || !onSort) return;
    const nextDir = sortKey === col.key && sortDirection === "asc" ? "desc" : "asc";
    onSort(col.key, nextDir);
  };

  // Loading skeleton rows
  if (isLoading) {
    return (
      <div className="rounded-xl border border-border overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="bg-bg-elevated">
              {columns.map((col) => (
                <th
                  key={col.key}
                  className="px-4 py-3 text-left text-xs font-medium text-text-secondary uppercase tracking-wider"
                >
                  {col.label}
                </th>
              ))}
              {actions && (
                <th className="px-4 py-3 text-right text-xs font-medium text-text-secondary uppercase tracking-wider">
                  Actions
                </th>
              )}
            </tr>
          </thead>
          <tbody>
            {Array.from({ length: pageSize }).map((_, i) => (
              <tr key={i} className="border-t border-border">
                {columns.map((col) => (
                  <td key={col.key} className="px-4 py-3">
                    <div className="h-4 bg-bg-hover rounded animate-pulse" />
                  </td>
                ))}
                {actions && (
                  <td className="px-4 py-3">
                    <div className="h-4 w-16 bg-bg-hover rounded animate-pulse ml-auto" />
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  // Empty state
  if (!data || data.length === 0) {
    return (
      <div className="rounded-xl border border-border overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="bg-bg-elevated">
              {columns.map((col) => (
                <th
                  key={col.key}
                  className="px-4 py-3 text-left text-xs font-medium text-text-secondary uppercase tracking-wider"
                >
                  {col.label}
                </th>
              ))}
              {actions && (
                <th className="px-4 py-3 text-right text-xs font-medium text-text-secondary uppercase tracking-wider">
                  Actions
                </th>
              )}
            </tr>
          </thead>
        </table>
        <div className="flex flex-col items-center justify-center py-16 text-text-muted">
          <Inbox size={40} className="mb-3 opacity-50" />
          <p className="text-sm">No records found</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="rounded-xl border border-border overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="bg-bg-elevated">
              {columns.map((col) => (
                <th
                  key={col.key}
                  onClick={() => handleSort(col)}
                  className={cn(
                    "px-4 py-3 text-left text-xs font-medium text-text-secondary uppercase tracking-wider",
                    col.sortable && "cursor-pointer select-none hover:text-foreground transition-colors"
                  )}
                >
                  <div className="flex items-center gap-1">
                    <span>{col.label}</span>
                    {col.sortable && sortKey === col.key && (
                      sortDirection === "asc"
                        ? <ChevronUp size={14} className="text-accent" />
                        : <ChevronDown size={14} className="text-accent" />
                    )}
                    {col.sortable && sortKey !== col.key && (
                      <ChevronUp size={14} className="opacity-0 group-hover:opacity-30" />
                    )}
                  </div>
                </th>
              ))}
              {actions && (
                <th className="px-4 py-3 text-right text-xs font-medium text-text-secondary uppercase tracking-wider">
                  Actions
                </th>
              )}
            </tr>
          </thead>
          <tbody>
            {data.map((row, rowIndex) => (
              <tr
                key={row.id || rowIndex}
                onClick={() => onRowClick && onRowClick(row)}
                className={cn(
                  "border-t border-border transition-colors",
                  onRowClick && "cursor-pointer",
                  "hover:bg-bg-hover/50"
                )}
              >
                {columns.map((col) => (
                  <td key={col.key} className="px-4 py-3 text-sm text-foreground">
                    {col.render ? col.render(row[col.key], row) : row[col.key]}
                  </td>
                ))}
                {actions && (
                  <td className="px-4 py-3 text-right">
                    <div onClick={(e) => e.stopPropagation()}>
                      {actions(row)}
                    </div>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-between text-sm">
        <span className="text-text-secondary">
          {"Showing " + startItem + "-" + endItem + " of " + total}
        </span>
        <div className="flex items-center gap-2">
          <button
            onClick={() => onPageChange(page - 1)}
            disabled={page <= 1}
            className={cn(
              "flex items-center gap-1 px-3 py-1.5 rounded-lg border border-border text-sm transition-colors",
              page <= 1
                ? "opacity-40 cursor-not-allowed text-text-muted"
                : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
            )}
          >
            <ChevronLeft size={14} />
            <span>Prev</span>
          </button>
          <span className="px-3 py-1.5 text-text-secondary">
            {"Page " + page + " of " + totalPages}
          </span>
          <button
            onClick={() => onPageChange(page + 1)}
            disabled={page >= totalPages}
            className={cn(
              "flex items-center gap-1 px-3 py-1.5 rounded-lg border border-border text-sm transition-colors",
              page >= totalPages
                ? "opacity-40 cursor-not-allowed text-text-muted"
                : "text-text-secondary hover:bg-bg-hover hover:text-foreground"
            )}
          >
            <span>Next</span>
            <ChevronRight size={14} />
          </button>
        </div>
      </div>
    </div>
  );
}
`
}

func desktopFormBuilder() string {
	return `import { useForm, Controller } from "react-hook-form";
import { cn } from "@/lib/utils";

interface FormField {
  name: string;
  label: string;
  type: "text" | "textarea" | "number" | "select" | "date" | "toggle";
  placeholder?: string;
  required?: boolean;
  options?: { label: string; value: string }[];
}

interface FormBuilderProps {
  fields: FormField[];
  onSubmit: (data: any) => void;
  defaultValues?: Record<string, any>;
  isLoading?: boolean;
  submitLabel?: string;
  onCancel?: () => void;
}

const inputClasses =
  "w-full bg-background border border-border rounded-lg px-3 py-2 text-foreground text-sm placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent focus:outline-none transition-colors";

const errorClasses = "text-xs text-danger mt-1";

function ToggleSwitch({ value, onChange }: { value: boolean; onChange: (val: boolean) => void }) {
  return (
    <button
      type="button"
      onClick={() => onChange(!value)}
      className={cn(
        "relative inline-flex h-6 w-11 items-center rounded-full transition-colors",
        value ? "bg-accent" : "bg-border"
      )}
    >
      <span
        className={cn(
          "inline-block h-4 w-4 rounded-full bg-white transition-transform",
          value ? "translate-x-6" : "translate-x-1"
        )}
      />
    </button>
  );
}

export default function FormBuilder({
  fields,
  onSubmit,
  defaultValues = {},
  isLoading,
  submitLabel = "Save",
  onCancel,
}: FormBuilderProps) {
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm({ defaultValues });

  const renderField = (field: FormField) => {
    const error = errors[field.name];
    const validation = field.required ? { required: field.label + " is required" } : {};

    switch (field.type) {
      case "text":
        return (
          <div key={field.name} className="space-y-1.5">
            <label className="block text-sm font-medium text-foreground">
              {field.label}
              {field.required && <span className="text-danger ml-1">*</span>}
            </label>
            <input
              type="text"
              placeholder={field.placeholder || ""}
              className={cn(inputClasses, error && "border-danger focus:border-danger focus:ring-danger")}
              {...register(field.name, validation)}
            />
            {error && <p className={errorClasses}>{String(error.message)}</p>}
          </div>
        );

      case "textarea":
        return (
          <div key={field.name} className="space-y-1.5">
            <label className="block text-sm font-medium text-foreground">
              {field.label}
              {field.required && <span className="text-danger ml-1">*</span>}
            </label>
            <textarea
              rows={4}
              placeholder={field.placeholder || ""}
              className={cn(inputClasses, "resize-none", error && "border-danger focus:border-danger focus:ring-danger")}
              {...register(field.name, validation)}
            />
            {error && <p className={errorClasses}>{String(error.message)}</p>}
          </div>
        );

      case "number":
        return (
          <div key={field.name} className="space-y-1.5">
            <label className="block text-sm font-medium text-foreground">
              {field.label}
              {field.required && <span className="text-danger ml-1">*</span>}
            </label>
            <input
              type="number"
              placeholder={field.placeholder || ""}
              className={cn(inputClasses, error && "border-danger focus:border-danger focus:ring-danger")}
              {...register(field.name, {
                ...validation,
                valueAsNumber: true,
              })}
            />
            {error && <p className={errorClasses}>{String(error.message)}</p>}
          </div>
        );

      case "select":
        return (
          <div key={field.name} className="space-y-1.5">
            <label className="block text-sm font-medium text-foreground">
              {field.label}
              {field.required && <span className="text-danger ml-1">*</span>}
            </label>
            <select
              className={cn(inputClasses, "appearance-none", error && "border-danger focus:border-danger focus:ring-danger")}
              {...register(field.name, validation)}
            >
              <option value="">{field.placeholder || "Select..."}</option>
              {(field.options || []).map((opt) => (
                <option key={opt.value} value={opt.value}>
                  {opt.label}
                </option>
              ))}
            </select>
            {error && <p className={errorClasses}>{String(error.message)}</p>}
          </div>
        );

      case "date":
        return (
          <div key={field.name} className="space-y-1.5">
            <label className="block text-sm font-medium text-foreground">
              {field.label}
              {field.required && <span className="text-danger ml-1">*</span>}
            </label>
            <input
              type="date"
              className={cn(inputClasses, error && "border-danger focus:border-danger focus:ring-danger")}
              {...register(field.name, validation)}
            />
            {error && <p className={errorClasses}>{String(error.message)}</p>}
          </div>
        );

      case "toggle":
        return (
          <div key={field.name} className="space-y-1.5">
            <div className="flex items-center justify-between">
              <label className="block text-sm font-medium text-foreground">
                {field.label}
              </label>
              <Controller
                name={field.name}
                control={control}
                render={({ field: f }) => (
                  <ToggleSwitch value={!!f.value} onChange={f.onChange} />
                )}
              />
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
      {fields.map(renderField)}

      <div className="flex items-center gap-3 pt-2">
        <button
          type="submit"
          disabled={isLoading}
          className={cn(
            "px-4 py-2 rounded-lg text-sm font-medium text-white bg-accent hover:bg-accent/90 transition-colors",
            isLoading && "opacity-60 cursor-not-allowed"
          )}
        >
          {isLoading ? (
            <span className="flex items-center gap-2">
              <span className="h-4 w-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
              Saving...
            </span>
          ) : (
            submitLabel
          )}
        </button>
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 rounded-lg text-sm font-medium text-text-secondary border border-border hover:bg-bg-hover hover:text-foreground transition-colors"
          >
            Cancel
          </button>
        )}
      </div>
    </form>
  );
}
`
}
