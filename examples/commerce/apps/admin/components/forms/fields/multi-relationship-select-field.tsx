"use client";

import { useState, useRef, useEffect, useMemo, useCallback, useId, lazy, Suspense } from "react";
import { createPortal } from "react-dom";
import { getResourceByEndpoint } from "@/resources";
import { usePermissions } from "@/hooks/use-permissions";
import { useDebouncedValue, useRelationshipOptions } from "@/hooks/use-resource";
import type { FieldDefinition } from "@/lib/resource";
import { Plus, X, Check } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";

// Lazy for the same reason as the single select — see adminInlineCreateDialog.
const InlineCreateDialog = lazy(() =>
  import("./inline-create-dialog").then((m) => ({ default: m.InlineCreateDialog }))
);

/**
 * The many-to-many picker: tags on a post, categories on a product.
 *
 * Built as a real combobox rather than a div that opens another div. The
 * version this replaces could not be reached with Tab, could not be opened
 * without a mouse, offered no way to move through the options with the
 * keyboard, and gave a screen reader a pile of unlabelled buttons. All of that
 * is fixed here, and none of it is decoration: this is the control somebody
 * uses fifty times while tagging a catalogue.
 *
 * What it does now:
 *
 *   - the trigger is a button with role="combobox", so Tab reaches it and
 *     Enter, Space or ArrowDown open it
 *   - the list is role="listbox" with aria-multiselectable, options carry
 *     aria-selected, and the active one is tracked with aria-activedescendant
 *     so a screen reader announces movement without focus ever leaving the
 *     search box
 *   - Arrow keys move, Enter toggles, Escape closes and returns focus to the
 *     trigger, Home and End jump
 *   - every remove button has a name: "Remove Electronics", not a bare X
 *
 * The relationship is not required to create the parent. You can save a
 * product with no categories and attach them afterwards, which is the normal
 * way a catalogue actually gets built: the products arrive first and the
 * taxonomy settles later.
 */

interface MultiRelationshipSelectFieldProps {
  field: FieldDefinition;
  value: string[];
  onChange: (value: string[]) => void;
  error?: string;
}

export function MultiRelationshipSelectField({
  field,
  value = [],
  onChange,
  error,
}: MultiRelationshipSelectFieldProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [active, setActive] = useState(0);
  const [creating, setCreating] = useState(false);
  const [position, setPosition] = useState({ top: 0, left: 0, width: 0 });

  // Records created inline, kept until the refetch returns them — otherwise a
  // record you just made vanishes from the list for a moment.
  const [justCreated, setJustCreated] = useState<Record<string, unknown>[]>([]);

  const triggerRef = useRef<HTMLButtonElement>(null);
  const searchRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLDivElement>(null);
  const listboxId = useId();
  const { can } = usePermissions();

  const endpoint = field.relatedEndpoint ?? "";
  const displayField = field.displayField ?? "name";

  const debouncedSearch = useDebouncedValue(search.trim(), 250);
  // The first page, which also names the chips. Fetched when the list opens, or
  // straight away when there are chips to label. Under the related endpoint's
  // own key, so a record created anywhere in the admin shows up here: the key
  // this used before was one no save ever invalidated.
  const base = useRelationshipOptions(endpoint, { enabled: open || value.length > 0 });
  // What was typed, searched on the server.
  const searched = useRelationshipOptions(endpoint, {
    search: debouncedSearch,
    enabled: open && debouncedSearch !== "",
  });
  const options = useMemo(
    () => (debouncedSearch && searched.data) || base.data || [],
    [debouncedSearch, searched.data, base.data]
  );
  const isLoading = options.length === 0 && (base.isLoading || searched.isLoading);
  // Once the server has answered for exactly what is in the box, its results
  // stand as they are. Until then the last results are narrowed here, so Enter
  // never toggles an option that does not match what was typed.
  const serverAnswered =
    debouncedSearch !== "" && debouncedSearch === search.trim() &&
    searched.isSuccess && !searched.isPlaceholderData;

  const allOptions = useMemo(() => {
    if (justCreated.length === 0) return options;
    const known = new Set(options.map((o) => String(o.id ?? "")));
    const extra = justCreated.filter((r) => !known.has(String(r.id ?? "")));
    return [...extra, ...options];
  }, [options, justCreated]);

  const labelOf = useCallback(
    (item: Record<string, unknown>) =>
      String(item[displayField] || item.name || item.title || item.id || ""),
    [displayField]
  );

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q || serverAnswered) return allOptions;
    return allOptions.filter((o) => labelOf(o).toLowerCase().includes(q));
  }, [allOptions, search, serverAnswered, labelOf]);

  const selectedLabels = useMemo(
    () =>
      value
        .map((id) => {
          const found = allOptions.find((o) => String(o.id) === String(id));
          return found ? { id: String(id), label: labelOf(found) } : null;
        })
        .filter(Boolean) as { id: string; label: string }[],
    [value, allOptions, labelOf]
  );

  const relatedResource = endpoint ? getResourceByEndpoint(endpoint) : undefined;
  const canCreate =
    field.allowCreate !== false && !!relatedResource && can(relatedResource.slug + ".create");

  // Carry the typed search into the new record, but only when the related form
  // actually has the field being searched.
  const prefill = useMemo(() => {
    const q = search.trim();
    if (!q || !relatedResource) return undefined;
    const target = relatedResource.form?.fields?.find((f) => f.key === displayField);
    return target ? { [displayField]: q } : undefined;
  }, [search, relatedResource, displayField]);

  const updatePosition = useCallback(() => {
    const el = triggerRef.current;
    if (!el) return;
    const r = el.getBoundingClientRect();
    setPosition({ top: r.bottom + window.scrollY + 4, left: r.left + window.scrollX, width: r.width });
  }, []);

  const openList = useCallback(() => {
    updatePosition();
    setOpen(true);
    setActive(0);
  }, [updatePosition]);

  const closeList = useCallback(
    (returnFocus = true) => {
      setOpen(false);
      setSearch("");
      // Focus goes back where it came from. Without this, closing the list
      // drops the caret at the top of the document and the next Tab starts
      // over from the beginning of the page.
      if (returnFocus) triggerRef.current?.focus();
    },
    []
  );

  const toggleItem = useCallback(
    (id: string) => {
      onChange(value.includes(id) ? value.filter((v) => v !== id) : [...value, id]);
    },
    [value, onChange]
  );

  // Reposition on scroll and resize while open, or the list detaches from its
  // trigger the moment the page moves under it.
  useEffect(() => {
    if (!open) return;
    const onMove = () => updatePosition();
    window.addEventListener("scroll", onMove, true);
    window.addEventListener("resize", onMove);
    return () => {
      window.removeEventListener("scroll", onMove, true);
      window.removeEventListener("resize", onMove);
    };
  }, [open, updatePosition]);

  useEffect(() => {
    if (open) searchRef.current?.focus();
  }, [open]);

  // Keep the active option in view when the arrows walk past the fold.
  // biome-ignore lint/correctness/useExhaustiveDependencies: active is read from the DOM, so the effect has to rerun when it moves.
  useEffect(() => {
    if (!open) return;
    const el = listRef.current?.querySelector('[data-active="true"]');
    el?.scrollIntoView({ block: "nearest" });
  }, [active, open]);

  useEffect(() => {
    if (!open) return;
    const onDown = (e: MouseEvent) => {
      const t = e.target as Node;
      if (triggerRef.current?.contains(t) || listRef.current?.contains(t)) return;
      setOpen(false);
      setSearch("");
    };
    document.addEventListener("mousedown", onDown);
    return () => document.removeEventListener("mousedown", onDown);
  }, [open]);

  function onSearchKeyDown(e: React.KeyboardEvent) {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActive((i) => Math.min(i + 1, filtered.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActive((i) => Math.max(i - 1, 0));
    } else if (e.key === "Home") {
      e.preventDefault();
      setActive(0);
    } else if (e.key === "End") {
      e.preventDefault();
      setActive(Math.max(filtered.length - 1, 0));
    } else if (e.key === "Enter") {
      e.preventDefault();
      const item = filtered[active];
      // Enter on an empty result set creates, because that is what somebody
      // typing a tag that does not exist yet is trying to do.
      if (item) toggleItem(String(item.id));
      else if (canCreate && search.trim()) {
        setOpen(false);
        setCreating(true);
      }
    } else if (e.key === "Escape") {
      e.preventDefault();
      closeList();
    } else if (e.key === "Backspace" && search === "" && value.length > 0) {
      // Backspace on an empty box removes the last chip, the way every tag
      // input people have used before this one behaves.
      onChange(value.slice(0, -1));
    }
  }

  const dropdown =
    open && typeof document !== "undefined"
      ? createPortal(
          <div
            ref={listRef}
            style={{ position: "absolute", top: position.top, left: position.left, width: position.width, zIndex: 9999 }}
            className="overflow-hidden rounded-md border border-border bg-bg-elevated shadow-lg"
          >
            <div className="border-b border-border p-2">
              <input
                ref={searchRef}
                value={search}
                onChange={(e) => {
                  setSearch(e.target.value);
                  setActive(0);
                }}
                onKeyDown={onSearchKeyDown}
                placeholder={"Search " + (field.label ?? "") + "..."}
                className={inputClasses({ inputSize: "sm" })}
                role="combobox"
                aria-expanded="true"
                aria-controls={listboxId}
                aria-autocomplete="list"
                aria-activedescendant={filtered[active] ? listboxId + "-" + String(filtered[active].id) : undefined}
                aria-label={"Search " + (field.label ?? "options")}
              />
            </div>

            <div
              id={listboxId}
              role="listbox"
              aria-multiselectable="true"
              aria-label={field.label}
              className="max-h-60 overflow-y-auto p-1"
            >
              {isLoading ? (
                <div className="px-3 py-2 text-sm text-text-secondary">Loading...</div>
              ) : filtered.length === 0 ? (
                <div className="px-3 py-2 text-sm text-text-secondary">
                  {search.trim() ? "Nothing matches “" + search.trim() + "”" : "Nothing to choose from yet"}
                </div>
              ) : (
                filtered.map((item, i) => {
                  const id = String(item.id);
                  const label = labelOf(item);
                  const isSelected = value.includes(id);
                  const isActive = i === active;
                  return (
                    <div
                      key={id}
                      id={listboxId + "-" + id}
                      role="option"
                      aria-selected={isSelected}
                      data-active={isActive}
                      onMouseEnter={() => setActive(i)}
                      onClick={() => toggleItem(id)}
                      className={
                        "flex cursor-pointer items-center gap-2 rounded-sm px-3 py-2 text-sm text-foreground " +
                        (isActive ? "bg-bg-hover" : "")
                      }
                    >
                      <span
                        className={
                          "flex h-4 w-4 shrink-0 items-center justify-center rounded border " +
                          (isSelected ? "border-accent bg-accent text-accent-fg" : "border-border")
                        }
                      >
                        {isSelected && <Check className="h-3 w-3" aria-hidden="true" />}
                      </span>
                      <span className="truncate">{label}</span>
                    </div>
                  );
                })
              )}
            </div>

            {value.length > 0 && (
              <div className="border-t border-border p-1">
                <button
                  type="button"
                  onClick={() => onChange([])}
                  className="flex w-full items-center rounded-sm px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover"
                >
                  Clear all
                </button>
              </div>
            )}

            {/* Outside the empty branch on purpose: "nothing matches" is exactly
                when somebody needs to create the record. */}
            {canCreate && relatedResource && (
              <div className="border-t border-border p-1">
                <button
                  type="button"
                  onClick={() => {
                    setOpen(false);
                    setCreating(true);
                  }}
                  className="flex w-full items-center gap-2 rounded-sm px-3 py-2 text-sm font-medium text-accent hover:bg-bg-hover"
                >
                  <Plus className="h-4 w-4 shrink-0" aria-hidden="true" />
                  <span className="truncate">
                    New {relatedResource.label?.singular ?? field.label}
                    {search.trim() ? " “" + search.trim() + "”" : ""}
                  </span>
                </button>
              </div>
            )}
          </div>,
          document.body
        )
      : null;

  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground" id={listboxId + "-label"}>
        {field.label}
        {field.required && <span className="ml-1 text-danger">*</span>}
      </label>

      {/* A real button, so Tab reaches it and Enter opens it. The chips live
          beside it rather than inside, because a button containing buttons is
          invalid and the remove controls have to be reachable on their own. */}
      <div
        className={
          "flex min-h-10 w-full flex-wrap items-center gap-1 rounded-md border bg-bg-secondary px-2 py-1.5 text-sm transition-colors " +
          (error ? "border-red-500 " : "border-border ") +
          (open ? "ring-2 ring-accent" : "")
        }
      >
        {selectedLabels.map(({ id, label }) => (
          <span
            key={id}
            className="inline-flex items-center gap-1 rounded-md bg-accent/20 px-2 py-0.5 text-xs font-medium text-accent"
          >
            {label}
            <button
              type="button"
              onClick={() => toggleItem(id)}
              aria-label={"Remove " + label}
              className="ml-0.5 rounded-full p-0.5 hover:bg-danger/20 hover:text-danger focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
            >
              <X className="h-3 w-3" aria-hidden="true" />
            </button>
          </span>
        ))}

        <button
          ref={triggerRef}
          type="button"
          role="combobox"
          aria-expanded={open}
          aria-haspopup="listbox"
          aria-controls={open ? listboxId : undefined}
          aria-labelledby={listboxId + "-label"}
          onClick={() => (open ? closeList(false) : openList())}
          onKeyDown={(e) => {
            if (e.key === "ArrowDown" || e.key === "Enter" || e.key === " ") {
              e.preventDefault();
              openList();
            }
          }}
          className="flex-1 rounded px-1 py-0.5 text-left text-text-secondary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
        >
          {selectedLabels.length > 0 ? "Add more..." : "Select " + (field.label ?? "") + "..."}
        </button>
      </div>

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
                setJustCreated((prev) => [record, ...prev]);
                // Append rather than replace: this select holds a list, and the
                // whole point of creating inline is to add to what is picked.
                const next = String(id);
                if (!value.includes(next)) onChange([...value, next]);
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
