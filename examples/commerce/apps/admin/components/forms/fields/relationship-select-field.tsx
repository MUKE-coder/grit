"use client";

import { useState, useRef, useEffect, useMemo, useCallback, lazy, Suspense } from "react";
import { createPortal } from "react-dom";
import { getResourceByEndpoint } from "@/resources";
import { usePermissions } from "@/hooks/use-permissions";
import { useDebouncedValue, useRelationshipOptions } from "@/hooks/use-resource";
import type { FieldDefinition } from "@/lib/resource";
import { Plus } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";

// Lazy on purpose — see the note on adminInlineCreateDialog. A static import
// here closes a cycle back through form-builder and the nested form silently
// renders as nothing.
const InlineCreateDialog = lazy(() =>
  import("./inline-create-dialog").then((m) => ({ default: m.InlineCreateDialog }))
);

interface RelationshipSelectFieldProps {
  field: FieldDefinition;
  value: string | null;
  onChange: (value: string | null) => void;
  error?: string;
}

export function RelationshipSelectField({ field, value, onChange, error }: RelationshipSelectFieldProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [creating, setCreating] = useState(false);
  // The record just created inline. Held locally so its label shows the instant
  // it is selected: the options query is refetching at that moment, and without
  // this the field displays a raw UUID until the network settles.
  const [justCreated, setJustCreated] = useState<Record<string, unknown> | null>(null);
  const triggerRef = useRef<HTMLButtonElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const [pos, setPos] = useState({ top: 0, left: 0, width: 0 });

  const endpoint = field.relatedEndpoint ?? "";
  const debouncedSearch = useDebouncedValue(search.trim(), 250);
  // The first page, which also names the selected record. Fetched when the list
  // opens, or straight away when there is a value to show a label for.
  const base = useRelationshipOptions(endpoint, { enabled: open || !!value });
  // What was typed, searched on the server.
  const searched = useRelationshipOptions(endpoint, {
    search: debouncedSearch,
    enabled: open && debouncedSearch !== "",
  });
  const options = (debouncedSearch && searched.data) || base.data || [];
  const isLoading = options.length === 0 && (base.isLoading || searched.isLoading);
  // Once the server has answered for exactly what is in the box, its results
  // stand as they are, matches on other columns included. Until then the last
  // results are narrowed here, so what is on screen always fits the box.
  const serverAnswered =
    debouncedSearch !== "" && debouncedSearch === search.trim() &&
    searched.isSuccess && !searched.isPlaceholderData;

  const updatePosition = useCallback(() => {
    if (triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      setPos({ top: rect.bottom + 4, left: rect.left, width: rect.width });
    }
  }, []);

  useEffect(() => {
    if (!open) return;
    updatePosition();
    function handleClickOutside(e: MouseEvent) {
      const target = e.target as Node;
      if (
        triggerRef.current && !triggerRef.current.contains(target) &&
        dropdownRef.current && !dropdownRef.current.contains(target)
      ) {
        setOpen(false);
      }
    }
    function handleScroll() { updatePosition(); }
    document.addEventListener("mousedown", handleClickOutside);
    window.addEventListener("scroll", handleScroll, true);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      window.removeEventListener("scroll", handleScroll, true);
    };
  }, [open, updatePosition]);

  const displayField = field.displayField || "name";

  // Options plus anything created inline that the refetch has not returned yet.
  const allOptions = useMemo(() => {
    const list = options as Record<string, unknown>[];
    if (!justCreated) return list;
    const id = String(justCreated.id ?? "");
    if (!id || list.some((o) => String(o.id) === id)) return list;
    return [justCreated, ...list];
  }, [options, justCreated]);

  // The related resource, looked up by the endpoint the field already points
  // at. Undefined when the related model has no registered admin resource — in
  // which case there is no form to open and the button must not appear.
  const relatedResource = useMemo(
    () => (field.relatedEndpoint ? getResourceByEndpoint(field.relatedEndpoint) : undefined),
    [field.relatedEndpoint]
  );

  // can() is false while permissions load, which is the right default here:
  // a button that appears and then vanishes reads as a bug, and this one opens
  // a form the API would reject anyway.
  const { can } = usePermissions();
  const canCreate =
    field.allowCreate !== false && !!relatedResource && can(relatedResource.slug + ".create");

  // Carry the typed search into the new record, but only when the related form
  // actually has that field. Guessing at the first field instead would drop the
  // text into whatever happens to be declared first.
  const prefill = useMemo(() => {
    const typed = search.trim();
    if (!typed || !relatedResource) return undefined;
    return relatedResource.form.fields.some((f) => f.key === displayField)
      ? { [displayField]: typed }
      : undefined;
  }, [search, relatedResource, displayField]);

  const filtered = useMemo(() =>
    allOptions.filter((item) => {
      if (!search || serverAnswered) return true;
      const label = String(item[displayField] || item.name || item.title || item.id || "");
      return label.toLowerCase().includes(search.toLowerCase());
    }),
    [allOptions, search, serverAnswered, displayField]
  );

  const selectedLabel = useMemo(() => {
    if (!value) return "";
    const found = allOptions.find((item) => String(item.id) === String(value));
    if (!found) return String(value);
    return String(found[displayField] || found.name || found.title || found.id || "");
  }, [value, allOptions, displayField]);

  const dropdown = open ? createPortal(
    <div
      ref={dropdownRef}
      className="fixed z-[9999] rounded-md border border-border bg-bg-elevated shadow-lg"
      style={{ top: pos.top, left: pos.left, width: pos.width, backgroundColor: "var(--bg-elevated, #22222e)" }}
    >
      <div className="p-2">
        <input
          type="text"
          placeholder="Search..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className={inputClasses({ inputSize: "sm", className: "flex" })}
          style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
          autoFocus
        />
      </div>
      <div className="max-h-60 overflow-y-auto p-1">
        {isLoading ? (
          <div className="px-3 py-2 text-sm text-text-secondary">Loading...</div>
        ) : filtered.length === 0 ? (
          <div className="px-3 py-2 text-sm text-text-secondary">No results found</div>
        ) : (
          <>
            {value && (
              <button
                type="button"
                onClick={() => { onChange(null); setOpen(false); setSearch(""); }}
                className="flex w-full items-center rounded-sm px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover"
              >
                Clear selection
              </button>
            )}
            {filtered.map((item) => {
              const id = String(item.id);
              const label = String(item[displayField] || item.name || item.title || item.id || "");
              return (
                <button
                  key={id}
                  type="button"
                  onClick={() => { onChange(id); setOpen(false); setSearch(""); }}
                  className={`flex w-full items-center rounded-sm px-3 py-2 text-sm text-foreground hover:bg-bg-hover
                    ${value === id ? "bg-bg-hover font-medium" : ""}`}
                >
                  {label}
                </button>
              );
            })}
          </>
        )}
      </div>

      {/* Deliberately OUTSIDE the empty/loading branch above. "No results
          found" is exactly the moment someone needs to create the record, and
          nesting this inside the populated branch would hide it precisely
          then. */}
      {canCreate && relatedResource && (
        <div className="border-t border-border p-1">
          <button
            type="button"
            onClick={() => { setOpen(false); setCreating(true); }}
            className="flex w-full items-center gap-2 rounded-sm px-3 py-2 text-sm font-medium text-accent hover:bg-bg-hover"
          >
            <Plus className="h-4 w-4 shrink-0" />
            <span className="truncate">
              New {relatedResource.label?.singular ?? field.label}
              {search.trim() ? ' “' + search.trim() + '”' : ""}
            </span>
          </button>
        </div>
      )}
    </div>,
    document.body
  ) : null;

  return (
    <div className="space-y-1.5">
      {/* Every other field type renders its label; this one did not, so a
          relationship was the one control on the form with nothing above it.
          Visible on any generated form with a belongs_to, not just a tree. */}
      <label className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>
      <button
        ref={triggerRef}
        type="button"
        onClick={() => { if (!open) updatePosition(); setOpen(!open); }}
        className={`flex h-10 w-full items-center justify-between rounded-md border bg-bg-secondary px-3 py-2 text-sm text-foreground transition-colors
          ${error ? "border-red-500" : "border-border"}
          ${open ? "ring-2 ring-accent" : ""}`}
      >
        <span className={value ? "text-foreground" : "text-text-secondary"}>
          {value ? selectedLabel : `Select ${field.label}...`}
        </span>
        <svg className="h-4 w-4 opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="m6 9 6 6 6-6" />
        </svg>
      </button>
      {dropdown}
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}

      {creating && relatedResource && (
        <Suspense fallback={null}>
          <InlineCreateDialog
            resource={relatedResource}
            defaults={prefill}
            onCreated={(record) => {
              const id = record?.id;
              if (id !== undefined && id !== null) {
                setJustCreated(record);
                onChange(String(id));
              }
              setCreating(false);
              setSearch("");
            }}
            onClose={() => setCreating(false)}
          />
        </Suspense>
      )}
    </div>
  );
}
