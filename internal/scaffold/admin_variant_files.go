package scaffold

// The admin half of the variant system.
//
// Every file here is written once, however many resources offer variants, and
// none of them is generated per resource. That is not a shortcut: the matrix
// editor takes the resource off the detail page's own props, so the same
// component drives Products, Courses and anything else without a second copy
// going stale the moment one of them is improved.
//
// It hangs off the DetailAside slot rather than living on a page of its own,
// because variants are a fact about one product and a shop owner looks for them
// where the product is. The option library, being shop-wide, is the opposite: a
// top-level entry in the sidebar, registered as a resource with a custom Page so
// it gets the nav item and the route for free.

// AdminVariantHooksTS emits apps/admin/hooks/use-variants.ts.
func AdminVariantHooksTS() string {
	return `"use client";

import { useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";

/**
 * Data access for options and variants.
 *
 * The endpoints are built from a resource slug rather than hard-coded, because
 * the same matrix editor serves every resource that offers variants. The option
 * endpoints have no slug in them at all: options are shop-wide, and nesting
 * them under a product would suggest each product owns its colours.
 */

export interface AdminOptionValue {
  id: string;
  option_id: string;
  label: string;
  slug: string;
  swatch?: string;
  /** Added to the base price, and only where the OPTION affects price. */
  price_delta: number;
  position: number;
}

export interface AdminOption {
  id: string;
  name: string;
  slug: string;
  /** How a storefront draws it: "swatch", "size" or "select". */
  kind: string;
  affects_price: boolean;
  position: number;
  values?: AdminOptionValue[];
}

export interface AdminVariant {
  id: string;
  sku: string;
  /** null when the price is resolved rather than pinned. */
  price_override: number | null;
  /** Resolved by the server: the override, or the base plus the deltas. */
  price: number;
  stock: number;
  active: boolean;
  position: number;
  option_values?: AdminOptionValue[];
}

export interface VariantMatrixData {
  options: AdminOption[];
  variants: AdminVariant[];
}

/** What a PATCH to one variant may carry. */
export interface VariantPatch {
  sku?: string;
  stock?: number;
  active?: boolean;
  price_override?: number;
  /**
   * Removes an override. A missing price_override cannot mean "clear it",
   * because that is also what a partial update sends for a field it is not
   * touching, so the server takes a separate flag.
   */
  clear_price?: boolean;
}

export const optionsKey = ["variant-options"];

export function matrixKey(slug: string, id: string) {
  return ["variant-matrix", slug, id];
}

/**
 * The path segment the per-variant update endpoint lives under: "Product"
 * becomes "product-variants".
 *
 * Not nested under the record as /<plural>/:id/variants/:variant, because gin
 * routes a static segment beside a parameter at the same position by panicking
 * at boot, and a bare /variants/:id collides the moment a second resource
 * offers variants. Derived from the resource NAME rather than its slug, because
 * depluralising "products" back to "product" on the client is a guess and
 * kebab-casing the Go name is not.
 */
export function variantPathFor(resourceName: string): string {
  return resourceName.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase() + "-variants";
}

export function useOptions() {
  return useQuery({
    queryKey: optionsKey,
    queryFn: async () => {
      const res = await apiClient.get("/api/options");
      return (res.data.data ?? []) as AdminOption[];
    },
  });
}

export function useVariantMatrix(slug: string, id: string) {
  return useQuery({
    queryKey: matrixKey(slug, id),
    enabled: Boolean(slug) && Boolean(id),
    queryFn: async () => {
      const res = await apiClient.get("/api/" + slug + "/" + id + "/variants");
      const data = res.data.data ?? {};
      return {
        options: data.options ?? [],
        variants: data.variants ?? [],
      } as VariantMatrixData;
    },
  });
}

export function useCreateOption() {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (input: { name: string; kind: string; affects_price: boolean }) => {
      const res = await apiClient.post("/api/options", input);
      return res.data.data as AdminOption;
    },
    successMessage: (option) => option.name + " added",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: optionsKey });
    },
  });
}

/**
 * Deletes an option and its values.
 *
 * The server refuses while anything is built on it and says which way, so there
 * is no pre-flight check here. A client that guessed would be guessing about
 * every other resource in the shop too.
 */
export function useDeleteOption() {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (optionID: string) => {
      await apiClient.delete("/api/options/" + optionID);
    },
    successMessage: "Option deleted",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: optionsKey });
    },
  });
}

export function useCreateOptionValue() {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (input: {
      optionID: string;
      label: string;
      swatch?: string;
      price_delta?: number;
    }) => {
      const { optionID, ...body } = input;
      const res = await apiClient.post("/api/options/" + optionID + "/values", body);
      return res.data.data as AdminOptionValue;
    },
    successMessage: (value) => value.label + " added",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: optionsKey });
    },
  });
}

export function useDeleteOptionValue() {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (valueID: string) => {
      await apiClient.delete("/api/option-values/" + valueID);
    },
    successMessage: "Value deleted",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: optionsKey });
    },
  });
}

/**
 * Sets which options a record offers, in order.
 *
 * Changing the set clears the existing matrix on the server, so the caller is
 * expected to have asked first. The response says how many rows went, which is
 * what the toast reports rather than a flat "saved" that hides the damage.
 */
export function useSetResourceOptions(slug: string, id: string) {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (optionIDs: string[]) => {
      const res = await apiClient.put("/api/" + slug + "/" + id + "/options", {
        option_ids: optionIDs,
      });
      return {
        cleared: (res.data.data?.variants_cleared ?? 0) as number,
        message: (res.data.message ?? "Options set") as string,
      };
    },
    successMessage: (result) =>
      result.cleared > 0
        ? result.message + ", " + result.cleared + " existing combinations cleared"
        : result.message,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: matrixKey(slug, id) });
    },
  });
}

export function useGenerateMatrix(slug: string, id: string) {
  const qc = useQueryClient();
  return useToastedMutation({
    mutationFn: async (limit?: number) => {
      const res = await apiClient.post(
        "/api/" + slug + "/" + id + "/variants/generate",
        limit ? { limit } : {},
      );
      return (res.data.data?.created ?? 0) as number;
    },
    successMessage: (created) =>
      created === 0 ? "Nothing to add, every combination exists" : created + " combinations added",
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: matrixKey(slug, id) });
    },
  });
}

/**
 * Saves every edited row in one go.
 *
 * Sequential rather than parallel, and one toast rather than one per row. A
 * matrix edit is a single act of intent even when it touched nine rows, and
 * nine toasts for it is a notification tray nobody reads.
 */
export function useSaveVariants(resourceName: string, slug: string, id: string) {
  const qc = useQueryClient();
  const path = variantPathFor(resourceName);
  return useToastedMutation({
    mutationFn: async (edits: Array<{ id: string; patch: VariantPatch }>) => {
      for (const edit of edits) {
        await apiClient.patch("/api/" + path + "/" + edit.id, edit.patch);
      }
      return edits.length;
    },
    successMessage: (count) => count + (count === 1 ? " variant saved" : " variants saved"),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: matrixKey(slug, id) });
    },
  });
}
`
}

// AdminVariantMatrixTSX emits apps/admin/components/variants/variant-matrix.tsx.
//
// The DetailAside slot on any resource that offers variants. It reads the slug
// off the resource it is handed and the base price off the record the detail
// controller already loaded, so it is the same component everywhere and there
// is no generated copy to drift.
func AdminVariantMatrixTSX() string {
	return `"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import type { ResourceDetailPartProps } from "@/lib/resource";
import {
  useVariantMatrix,
  useOptions,
  useSetResourceOptions,
  useGenerateMatrix,
  useSaveVariants,
  type AdminOption,
  type AdminOptionValue,
  type AdminVariant,
  type VariantPatch,
} from "@/hooks/use-variants";
import { Button, buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { formatCurrency } from "@/lib/formatters";
import { AlertCircle, Check, Loader2, Save, Settings2, Sparkles, X } from "@/lib/icons";

/**
 * The variant matrix, on the record's own detail page.
 *
 * Rendered through the DetailAside slot, which is why it takes the resource and
 * the detail controller rather than a product id: the slug tells it which
 * endpoints to call and the loaded record tells it the base price, so one
 * component serves every resource that offers variants.
 *
 * The base price matters more than it looks. A variant's price is resolved and
 * not stored, so this table has to show what each combination WOULD cost with
 * no override, and that number is the record's price plus the deltas of the
 * values whose option affects price. Computing it here mirrors what the server
 * does, and it is the only way an empty override box can be honest about what
 * clearing it would mean.
 */

/** One row's unsaved edits. An absent key means the field was not touched. */
interface VariantDraft {
  sku?: string;
  stock?: number;
  active?: boolean;
  /** null clears the override; a number pins one. */
  priceOverride?: number | null;
}

// Generic in the row type rather than pinned to Record<string, unknown>.
//
// A customisation file is typed against its own model — ResourceCustomisation
// <Product> — so a slot component fixed to the erased row type is not
// assignable to it, and attaching this to a real resource would not compile.
// The parameter is never used for anything but that assignability.
export function VariantMatrix<T>({ resource, id, controller }: ResourceDetailPartProps<T>) {
  const matrix = useVariantMatrix(resource.slug, id);
  const library = useOptions();
  const setOptions = useSetResourceOptions(resource.slug, id);
  const generate = useGenerateMatrix(resource.slug, id);
  const save = useSaveVariants(resource.name, resource.slug, id);

  const [picking, setPicking] = useState(false);
  const [chosen, setChosen] = useState<string[]>([]);
  const [confirmSwap, setConfirmSwap] = useState(false);
  const [drafts, setDrafts] = useState<Record<string, VariantDraft>>({});

  const options = matrix.data?.options ?? [];
  const variants = matrix.data?.variants ?? [];
  const singular = (resource.label?.singular ?? resource.name).toLowerCase();

  // Through unknown, because the row type is the caller's and this only ever
  // asks it one question.
  const record = controller.record as unknown as { price?: unknown } | undefined;
  const basePrice = typeof record?.price === "number" ? record.price : 0;

  /** What the cartesian product would come to, for the generate button's label. */
  const combinations = useMemo(() => {
    if (options.length === 0) return 0;
    return options.reduce((total, option) => total * Math.max(option.values?.length ?? 0, 1), 1);
  }, [options]);

  /** The resolved price with no override: base plus the deltas that count. */
  const computed = useMemo(() => {
    const byID = new Map(options.map((option) => [option.id, option]));
    return (variant: AdminVariant) => {
      let price = basePrice;
      for (const value of variant.option_values ?? []) {
        const option = byID.get(value.option_id);
        // Read off the option, never the value. That is what stops a shop
        // charging extra for a colour because a delta was typed on one swatch.
        if (option?.affects_price) price += value.price_delta;
      }
      return price;
    };
  }, [options, basePrice]);

  const edits = useMemo(() => collectEdits(variants, drafts), [variants, drafts]);

  function patchDraft(variantID: string, patch: VariantDraft) {
    setDrafts((prev) => ({ ...prev, [variantID]: { ...prev[variantID], ...patch } }));
  }

  function openPicker() {
    setChosen(options.map((option) => option.id));
    setPicking(true);
  }

  function toggleOption(optionID: string) {
    // Selection order is the display order the storefront gets, so a click
    // appends rather than slotting the option back into library order.
    setChosen((prev) =>
      prev.includes(optionID) ? prev.filter((each) => each !== optionID) : [...prev, optionID],
    );
  }

  function applyOptions() {
    const unchanged =
      chosen.length === options.length && chosen.every((each, i) => options[i]?.id === each);
    if (unchanged) {
      setPicking(false);
      return;
    }
    // Changing the set destroys the matrix server-side, so it is asked about
    // whenever there is one to lose.
    if (variants.length > 0) {
      setConfirmSwap(true);
      return;
    }
    void commitOptions();
  }

  async function commitOptions() {
    setConfirmSwap(false);
    await setOptions.mutateAsync(chosen);
    setDrafts({});
    setPicking(false);
  }

  async function saveEdits() {
    if (edits.length === 0) return;
    await save.mutateAsync(edits);
    setDrafts({});
  }

  if (matrix.isLoading) {
    return (
      <section className="rounded-xl border border-border bg-bg-elevated p-6">
        <div className="flex items-center gap-2 text-sm text-text-muted">
          <Loader2 className="h-4 w-4 animate-spin" /> Loading variants...
        </div>
      </section>
    );
  }

  return (
    <section className="rounded-xl border border-border bg-bg-elevated">
      <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-6 py-4">
        <div className="min-w-0">
          <h2 className="text-sm font-semibold text-foreground">Variants</h2>
          <p className="mt-0.5 text-xs text-text-muted">
            {options.length === 0
              ? "This " + singular + " offers no options yet."
              : options.map((option) => option.name).join(" x ") +
                " — " +
                variants.length +
                " of " +
                combinations +
                " combinations"}
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          <Link href="/resources/options" className={buttonClasses({ variant: "ghost", size: "sm" })}>
            Option library
          </Link>
          <Button variant="outline" size="sm" onClick={openPicker}>
            <Settings2 className="h-4 w-4" /> Choose options
          </Button>
          {options.length > 0 && variants.length < combinations && (
            <Button
              size="sm"
              disabled={generate.isPending}
              onClick={() => generate.mutate(undefined)}
            >
              {generate.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Sparkles className="h-4 w-4" />
              )}
              Generate {combinations - variants.length} missing
            </Button>
          )}
          {edits.length > 0 && (
            <Button size="sm" disabled={save.isPending} onClick={saveEdits}>
              {save.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Save className="h-4 w-4" />
              )}
              Save {edits.length}
            </Button>
          )}
        </div>
      </header>

      {picking && (
        <OptionPicker
          library={library.data ?? []}
          loading={library.isLoading}
          chosen={chosen}
          onToggle={toggleOption}
          onCancel={() => setPicking(false)}
          onApply={applyOptions}
          saving={setOptions.isPending}
        />
      )}

      {options.length === 0 && !picking && (
        <EmptyPanel
          title={"No options on this " + singular}
          body={
            "A variant is a combination of choices, so pick the axes first: Colour, " +
            "Size, Memory. They come from the shop-wide library, which is what keeps " +
            "one spelling of Colour across the whole catalogue."
          }
          action={
            <Button size="sm" onClick={openPicker}>
              <Settings2 className="h-4 w-4" /> Choose options
            </Button>
          }
        />
      )}

      {options.length > 0 && variants.length === 0 && !picking && (
        <EmptyPanel
          title="No combinations yet"
          body={
            "Generating writes one row per combination and leaves anything already " +
            "there alone, so it is safe to run again after adding a value."
          }
          action={
            <Button size="sm" disabled={generate.isPending} onClick={() => generate.mutate(undefined)}>
              {generate.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Sparkles className="h-4 w-4" />
              )}
              Generate {combinations} combinations
            </Button>
          }
        />
      )}

      {variants.length > 0 && (
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border text-left text-xs uppercase tracking-wide text-text-muted">
                {options.map((option) => (
                  <th key={option.id} className="px-6 py-3 font-medium">
                    {option.name}
                  </th>
                ))}
                <th className="px-3 py-3 font-medium">SKU</th>
                <th className="px-3 py-3 font-medium">Stock</th>
                <th className="px-3 py-3 font-medium">Price</th>
                <th className="px-3 py-3 font-medium">Override</th>
                <th className="px-3 py-3 font-medium">Active</th>
              </tr>
            </thead>
            <tbody>
              {variants.map((variant) => {
                const draft = drafts[variant.id];
                const dirty = isDirty(variant, draft);
                return (
                  <tr
                    key={variant.id}
                    className={
                      "border-b border-border/60 last:border-0 " +
                      (dirty ? "bg-accent/5" : "")
                    }
                  >
                    {options.map((option) => (
                      <td key={option.id} className="px-6 py-2.5">
                        <ValueCell option={option} value={valueFor(variant, option.id)} />
                      </td>
                    ))}

                    <td className="px-3 py-2.5">
                      {/* Wide enough for a real SKU. The generated ones carry
                          the record's slug and both option values, and a box
                          that shows the first fifteen characters of that is a
                          box you cannot check anything against. */}
                      <div className="w-56">
                        <input
                          className={inputClasses({ inputSize: "sm", className: "font-mono text-xs" })}
                          placeholder="unset"
                          value={draft?.sku ?? variant.sku ?? ""}
                          onChange={(e) => patchDraft(variant.id, { sku: e.target.value })}
                        />
                      </div>
                    </td>

                    <td className="px-3 py-2.5">
                      <div className="w-20">
                        <input
                          className={inputClasses({ inputSize: "sm", className: "tabular-nums" })}
                          inputMode="numeric"
                          value={String(draft?.stock ?? variant.stock)}
                          onChange={(e) => {
                            const next = parseInt(e.target.value, 10);
                            patchDraft(variant.id, { stock: Number.isNaN(next) ? 0 : next });
                          }}
                        />
                      </div>
                    </td>

                    <td className="px-3 py-2.5 whitespace-nowrap tabular-nums text-text-secondary">
                      {formatCurrency(effectivePrice(variant, draft, computed))}
                    </td>

                    <td className="px-3 py-2.5">
                      <div className="w-24">
                        <input
                        className={inputClasses({ inputSize: "sm", className: "tabular-nums" })}
                        inputMode="decimal"
                        /* Empty means resolved, and the placeholder says what
                           that resolves to, so clearing the box is never a
                           guess about what the price becomes. */
                        placeholder={formatCurrency(computed(variant))}
                        value={overrideInput(variant, draft)}
                        onChange={(e) => {
                          const raw = e.target.value.trim();
                          if (raw === "") {
                            patchDraft(variant.id, { priceOverride: null });
                            return;
                          }
                          const next = parseFloat(raw);
                          if (!Number.isNaN(next)) patchDraft(variant.id, { priceOverride: next });
                        }}
                        />
                      </div>
                    </td>

                    <td className="px-3 py-2.5">
                      <button
                        type="button"
                        aria-label={(draft?.active ?? variant.active) ? "Deactivate" : "Activate"}
                        onClick={() =>
                          patchDraft(variant.id, { active: !(draft?.active ?? variant.active) })
                        }
                        className={
                          "inline-flex h-7 w-7 items-center justify-center rounded-lg border transition-colors " +
                          ((draft?.active ?? variant.active)
                            ? "border-success/40 bg-success/10 text-success"
                            : "border-border text-text-muted hover:text-foreground")
                        }
                      >
                        {(draft?.active ?? variant.active) ? (
                          <Check className="h-4 w-4" />
                        ) : (
                          <X className="h-4 w-4" />
                        )}
                      </button>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      )}

      <ConfirmModal
        open={confirmSwap}
        title="Replace the options?"
        description={
          "A variant is a combination of the options this " +
          singular +
          " offers, so changing them clears all " +
          variants.length +
          " existing combinations along with their SKUs, stock and prices. This cannot be undone."
        }
        confirmLabel="Replace and clear"
        variant="danger"
        loading={setOptions.isPending}
        onCancel={() => setConfirmSwap(false)}
        onConfirm={() => void commitOptions()}
      />
    </section>
  );
}

// ─── pieces ──────────────────────────────────────────────────────────

function OptionPicker({
  library,
  loading,
  chosen,
  onToggle,
  onCancel,
  onApply,
  saving,
}: {
  library: AdminOption[];
  loading: boolean;
  chosen: string[];
  onToggle: (id: string) => void;
  onCancel: () => void;
  onApply: () => void;
  saving: boolean;
}) {
  return (
    <div className="border-b border-border bg-bg-secondary px-6 py-5">
      <p className="mb-3 text-xs text-text-muted">
        Tick the axes this record offers. The order you tick them is the order a storefront draws
        them in.
      </p>

      {loading && (
        <div className="flex items-center gap-2 text-sm text-text-muted">
          <Loader2 className="h-4 w-4 animate-spin" /> Loading the library...
        </div>
      )}

      {!loading && library.length === 0 && (
        <div className="flex items-start gap-2 rounded-lg border border-border bg-bg-elevated p-4 text-sm text-text-secondary">
          <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-warning" />
          <span>
            The option library is empty. Add Colour or Size in{" "}
            <Link href="/resources/options" className="text-accent hover:underline">
              Options
            </Link>{" "}
            first — they are shared by every product, which is what keeps one spelling of each.
          </span>
        </div>
      )}

      <div className="flex flex-wrap gap-2">
        {library.map((option) => {
          const index = chosen.indexOf(option.id);
          const selected = index >= 0;
          return (
            <button
              key={option.id}
              type="button"
              onClick={() => onToggle(option.id)}
              className={
                "inline-flex items-center gap-2 rounded-lg border px-3 py-2 text-sm transition-colors " +
                (selected
                  ? "border-accent bg-accent/10 text-foreground"
                  : "border-border text-text-secondary hover:border-accent/40 hover:text-foreground")
              }
            >
              {selected && (
                <span className="inline-flex h-5 w-5 items-center justify-center rounded-md bg-accent text-[11px] font-semibold text-white">
                  {index + 1}
                </span>
              )}
              {option.name}
              <span className="text-xs text-text-muted">{option.values?.length ?? 0}</span>
              {option.affects_price && (
                <span className="rounded bg-bg-tertiary px-1.5 py-0.5 text-[10px] text-text-muted">
                  priced
                </span>
              )}
            </button>
          );
        })}
      </div>

      <div className="mt-4 flex items-center gap-2">
        <Button size="sm" disabled={saving} onClick={onApply}>
          {saving && <Loader2 className="h-4 w-4 animate-spin" />} Apply
        </Button>
        <Button size="sm" variant="ghost" onClick={onCancel}>
          Cancel
        </Button>
      </div>
    </div>
  );
}

function ValueCell({ option, value }: { option: AdminOption; value?: AdminOptionValue }) {
  if (!value) return <span className="text-text-muted">—</span>;
  return (
    <span className="inline-flex items-center gap-2 whitespace-nowrap">
      {option.kind === "swatch" && value.swatch && (
        <span
          className="h-4 w-4 shrink-0 rounded-full border border-border"
          style={{ backgroundColor: value.swatch }}
        />
      )}
      <span className="text-foreground">{value.label}</span>
    </span>
  );
}

function EmptyPanel({
  title,
  body,
  action,
}: {
  title: string;
  body: string;
  action: React.ReactNode;
}) {
  return (
    <div className="px-6 py-10 text-center">
      <h3 className="text-sm font-medium text-foreground">{title}</h3>
      <p className="mx-auto mt-1.5 max-w-md text-xs leading-relaxed text-text-muted">{body}</p>
      <div className="mt-4 flex justify-center">{action}</div>
    </div>
  );
}

// ─── draft arithmetic ────────────────────────────────────────────────

function valueFor(variant: AdminVariant, optionID: string) {
  return (variant.option_values ?? []).find((value) => value.option_id === optionID);
}

/** What the override box shows: the draft if touched, otherwise what is stored. */
function overrideInput(variant: AdminVariant, draft?: VariantDraft): string {
  if (draft?.priceOverride !== undefined) {
    return draft.priceOverride === null ? "" : String(draft.priceOverride);
  }
  return variant.price_override === null || variant.price_override === undefined
    ? ""
    : String(variant.price_override);
}

/**
 * The price the row displays.
 *
 * An untouched row shows the server's figure. A touched one is recomputed here,
 * including the case that matters: clearing the override falls back to the
 * resolved price rather than leaving the old pinned one on screen.
 */
function effectivePrice(
  variant: AdminVariant,
  draft: VariantDraft | undefined,
  computed: (variant: AdminVariant) => number,
): number {
  if (draft?.priceOverride !== undefined) {
    return draft.priceOverride === null ? computed(variant) : draft.priceOverride;
  }
  return variant.price;
}

function isDirty(variant: AdminVariant, draft?: VariantDraft): boolean {
  if (!draft) return false;
  if (draft.sku !== undefined && draft.sku !== (variant.sku ?? "")) return true;
  if (draft.stock !== undefined && draft.stock !== variant.stock) return true;
  if (draft.active !== undefined && draft.active !== variant.active) return true;
  if (
    draft.priceOverride !== undefined &&
    draft.priceOverride !== (variant.price_override ?? null)
  ) {
    return true;
  }
  return false;
}

/**
 * Turns the drafts into PATCH bodies, dropping fields that were typed back to
 * what they already were.
 *
 * Sending those anyway would work, and would also bump the version column on
 * every row somebody clicked into, which turns an audit trail into noise.
 */
function collectEdits(
  variants: AdminVariant[],
  drafts: Record<string, VariantDraft>,
): Array<{ id: string; patch: VariantPatch }> {
  const out: Array<{ id: string; patch: VariantPatch }> = [];
  for (const variant of variants) {
    const draft = drafts[variant.id];
    if (!isDirty(variant, draft)) continue;

    const patch: VariantPatch = {};
    if (draft.sku !== undefined && draft.sku !== (variant.sku ?? "")) patch.sku = draft.sku;
    if (draft.stock !== undefined && draft.stock !== variant.stock) patch.stock = draft.stock;
    if (draft.active !== undefined && draft.active !== variant.active) patch.active = draft.active;
    if (draft.priceOverride !== undefined && draft.priceOverride !== (variant.price_override ?? null)) {
      if (draft.priceOverride === null) patch.clear_price = true;
      else patch.price_override = draft.priceOverride;
    }
    out.push({ id: variant.id, patch });
  }
  return out;
}
`
}

// AdminOptionLibraryTSX emits apps/admin/components/variants/option-library.tsx.
//
// The shop-wide half. It is a page rather than a panel because options belong
// to the shop and not to any one product, and putting them on a product page is
// the per-product-options mistake the schema exists to avoid: fourteen
// spellings of Colour and a filter that can only match one.
func AdminOptionLibraryTSX() string {
	return `"use client";

import { useState } from "react";
import type { ResourcePageSlotProps } from "@/lib/resource";
import {
  useOptions,
  useCreateOption,
  useCreateOptionValue,
  useDeleteOption,
  useDeleteOptionValue,
  type AdminOption,
} from "@/hooks/use-variants";
import { PageHeader } from "@/components/chrome/PageHeader";
import { Button, buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { Loader2, Plus, Trash2, X } from "@/lib/icons";

/**
 * The shared option library.
 *
 * Two things about the shape of an option are decisions rather than details,
 * and the form says so where somebody is about to make them:
 *
 *   kind          how a storefront draws it. Deciding that on the client from
 *                 the option's name is how one shop renders Colour as a
 *                 dropdown because somebody spelled it Color.
 *
 *   affects_price a fact about the axis, not about one value on it. Memory
 *                 changes a laptop's price; colour does not. Per-value would
 *                 allow "32GB is priced but 16GB is not", which nobody means.
 */

const KINDS = [
  { value: "select", label: "Dropdown", hint: "A list. The safe default for anything with words in it." },
  { value: "swatch", label: "Swatch", hint: "Colour dots. Each value carries a CSS colour." },
  { value: "size", label: "Size boxes", hint: "A row of short labels: S, M, L, 42." },
];

export function OptionLibraryPage(_: ResourcePageSlotProps) {
  const options = useOptions();
  const createOption = useCreateOption();
  const removeOption = useDeleteOption();

  const [adding, setAdding] = useState(false);
  const [name, setName] = useState("");
  const [kind, setKind] = useState("select");
  const [affectsPrice, setAffectsPrice] = useState(false);
  const [confirmRemove, setConfirmRemove] = useState<AdminOption | null>(null);

  async function submitOption(e: React.FormEvent) {
    e.preventDefault();
    if (!name.trim()) return;
    await createOption.mutateAsync({ name: name.trim(), kind, affects_price: affectsPrice });
    setName("");
    setKind("select");
    setAffectsPrice(false);
    setAdding(false);
  }

  return (
    <div className="space-y-6">
      <PageHeader
        title="Options"
        subtitle="Colour, Size, Memory — shared by every product that offers variants."
        refreshKeys={["variant-options"]}
        backHref={null}
        actions={
          <Button size="sm" onClick={() => setAdding((open) => !open)}>
            <Plus className="h-4 w-4" /> New option
          </Button>
        }
      />

      {adding && (
        <form
          onSubmit={submitOption}
          className="rounded-xl border border-border bg-bg-elevated p-6"
        >
          <div className="grid gap-4 sm:grid-cols-2">
            <label className="block">
              <span className="mb-1.5 block text-xs font-medium text-text-secondary">Name</span>
              <input
                autoFocus
                className={inputClasses()}
                placeholder="Colour"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </label>

            <label className="block">
              <span className="mb-1.5 block text-xs font-medium text-text-secondary">
                How a storefront draws it
              </span>
              <select
                className={inputClasses()}
                value={kind}
                onChange={(e) => setKind(e.target.value)}
              >
                {KINDS.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
              <span className="mt-1 block text-xs text-text-muted">
                {KINDS.find((option) => option.value === kind)?.hint}
              </span>
            </label>
          </div>

          <label className="mt-4 flex items-start gap-2.5">
            <input
              type="checkbox"
              className="mt-0.5 h-4 w-4 rounded border-border accent-accent"
              checked={affectsPrice}
              onChange={(e) => setAffectsPrice(e.target.checked)}
            />
            <span className="text-sm text-text-secondary">
              Choosing on this axis changes the price
              <span className="mt-0.5 block text-xs text-text-muted">
                Turn this on for Memory or Capacity, off for Colour. The per-value amounts are
                ignored entirely while it is off, which is what stops a stray number on one swatch
                charging a customer extra for red.
              </span>
            </span>
          </label>

          <div className="mt-5 flex items-center gap-2">
            <Button type="submit" size="sm" disabled={createOption.isPending || !name.trim()}>
              {createOption.isPending && <Loader2 className="h-4 w-4 animate-spin" />} Create
            </Button>
            <Button type="button" size="sm" variant="ghost" onClick={() => setAdding(false)}>
              Cancel
            </Button>
          </div>
        </form>
      )}

      {options.isLoading && (
        <div className="flex items-center gap-2 rounded-xl border border-border bg-bg-elevated p-6 text-sm text-text-muted">
          <Loader2 className="h-4 w-4 animate-spin" /> Loading the library...
        </div>
      )}

      {!options.isLoading && (options.data?.length ?? 0) === 0 && (
        <div className="rounded-xl border border-border bg-bg-elevated px-6 py-12 text-center">
          <h2 className="text-sm font-medium text-foreground">No options yet</h2>
          <p className="mx-auto mt-1.5 max-w-md text-xs leading-relaxed text-text-muted">
            An option is one axis of choice. Add Colour and Size here once, then attach them to any
            product from its detail page. Keeping them shared is what stops a catalogue accumulating
            four spellings of the same thing.
          </p>
          <div className="mt-4 flex justify-center">
            <Button size="sm" onClick={() => setAdding(true)}>
              <Plus className="h-4 w-4" /> New option
            </Button>
          </div>
        </div>
      )}

      <div className="grid gap-4">
        {(options.data ?? []).map((option) => (
          <OptionCard key={option.id} option={option} onRemove={() => setConfirmRemove(option)} />
        ))}
      </div>

      <ConfirmModal
        open={confirmRemove !== null}
        title={"Delete " + (confirmRemove?.name ?? "this option") + "?"}
        description="Its values go with it. Anything already built on them has to be cleared first, and the server will say so rather than letting it happen."
        confirmLabel="Delete"
        variant="danger"
        loading={removeOption.isPending}
        onCancel={() => setConfirmRemove(null)}
        onConfirm={async () => {
          if (!confirmRemove) return;
          await removeOption.mutateAsync(confirmRemove.id);
          setConfirmRemove(null);
        }}
      />
    </div>
  );
}

function OptionCard({ option, onRemove }: { option: AdminOption; onRemove: () => void }) {
  const addValue = useCreateOptionValue();
  const removeValue = useDeleteOptionValue();

  const [label, setLabel] = useState("");
  const [swatch, setSwatch] = useState("#111118");
  const [delta, setDelta] = useState("");

  async function submitValue(e: React.FormEvent) {
    e.preventDefault();
    if (!label.trim()) return;
    const parsed = parseFloat(delta);
    await addValue.mutateAsync({
      optionID: option.id,
      label: label.trim(),
      swatch: option.kind === "swatch" ? swatch : undefined,
      price_delta: Number.isNaN(parsed) ? 0 : parsed,
    });
    setLabel("");
    setDelta("");
  }

  return (
    <section className="rounded-xl border border-border bg-bg-elevated">
      <header className="flex items-center justify-between gap-3 border-b border-border px-6 py-4">
        <div className="flex items-center gap-2.5">
          <h2 className="text-sm font-semibold text-foreground">{option.name}</h2>
          <span className="rounded bg-bg-tertiary px-2 py-0.5 text-[11px] text-text-muted">
            {KINDS.find((each) => each.value === option.kind)?.label ?? option.kind}
          </span>
          {option.affects_price && (
            <span className="rounded bg-accent/10 px-2 py-0.5 text-[11px] text-accent">
              changes the price
            </span>
          )}
        </div>
        <button
          type="button"
          aria-label={"Delete " + option.name}
          onClick={onRemove}
          className="rounded-lg p-2 text-text-muted transition-colors hover:bg-bg-hover hover:text-danger"
        >
          <Trash2 className="h-4 w-4" />
        </button>
      </header>

      <div className="flex flex-wrap gap-2 px-6 py-4">
        {(option.values ?? []).length === 0 && (
          <p className="text-xs text-text-muted">
            No values yet. An option with none of them cannot be part of any combination.
          </p>
        )}
        {(option.values ?? []).map((value) => (
          <span
            key={value.id}
            className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary py-1.5 pl-2.5 pr-1.5 text-sm"
          >
            {option.kind === "swatch" && value.swatch && (
              <span
                className="h-4 w-4 rounded-full border border-border"
                style={{ backgroundColor: value.swatch }}
              />
            )}
            <span className="text-foreground">{value.label}</span>
            {option.affects_price && value.price_delta !== 0 && (
              <span className="tabular-nums text-xs text-text-muted">
                {value.price_delta > 0 ? "+" : ""}
                {value.price_delta}
              </span>
            )}
            <button
              type="button"
              aria-label={"Delete " + value.label}
              disabled={removeValue.isPending}
              onClick={() => removeValue.mutate(value.id)}
              className="rounded p-1 text-text-muted transition-colors hover:bg-bg-hover hover:text-danger disabled:opacity-50"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </span>
        ))}
      </div>

      <form
        onSubmit={submitValue}
        className="flex flex-wrap items-center gap-2 border-t border-border px-6 py-3"
      >
        <div className="w-40">
          <input
            className={inputClasses({ inputSize: "sm" })}
            placeholder={option.kind === "size" ? "XL" : "Black"}
            value={label}
            onChange={(e) => setLabel(e.target.value)}
          />
        </div>
        {option.kind === "swatch" && (
          <input
            type="color"
            aria-label="Swatch colour"
            className="h-8 w-10 cursor-pointer rounded-lg border border-border bg-bg-secondary p-1"
            value={swatch}
            onChange={(e) => setSwatch(e.target.value)}
          />
        )}
        {option.affects_price && (
          <div className="w-28">
            <input
              className={inputClasses({ inputSize: "sm", className: "tabular-nums" })}
              inputMode="decimal"
              placeholder="+ / - price"
              value={delta}
              onChange={(e) => setDelta(e.target.value)}
            />
          </div>
        )}
        <button
          type="submit"
          disabled={addValue.isPending || !label.trim()}
          className={buttonClasses({ variant: "outline", size: "sm" })}
        >
          {addValue.isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Plus className="h-4 w-4" />
          )}
          Add value
        </button>
      </form>
    </section>
  );
}
`
}

// AdminOptionsResourceTS emits apps/admin/resources/options/options.ts.
//
// Registered as a resource purely for the sidebar entry and the route: the
// generic list view could not serve this one, because the values nested under
// an option are a second level the table and form know nothing about. The Page
// slot replaces it wholesale, which is what that slot is for.
func AdminOptionsResourceTS() string {
	return `import { defineResource } from "@/lib/resource";
import { OptionLibraryPage } from "@/components/variants/option-library";

/**
 * The shared option library, as a sidebar entry.
 *
 * The table and form below are never rendered — components.Page replaces the
 * whole list view — but they are not decoration either. The registry reads
 * endpoint and label to resolve relationships and breadcrumbs, so they describe
 * the resource honestly rather than being left empty.
 */
export const optionResource = defineResource(
  {
    name: "Option",
    slug: "options",
    endpoint: "/api/options",
    icon: "Layers",
    label: { singular: "Option", plural: "Options" },
    adminOnly: true,
    // The stat cards are counts of a resource the server knows how to count,
    // and /admin/dashboard/resource-stats only knows the resources the
    // generator registered. Options were registered here by hand, so asking is
    // two 400s on every visit. There is nothing to count anyway: the page shows
    // the whole library.
    stats: false,
    table: {
      columns: [
        { key: "name", label: "Name", sortable: true, searchable: true },
        { key: "kind", label: "Kind" },
        { key: "affects_price", label: "Changes price", format: "boolean" },
      ],
      searchable: true,
    },
    form: {
      fields: [
        { key: "name", label: "Name", type: "text", required: true },
        {
          key: "kind",
          label: "How a storefront draws it",
          type: "select",
          options: [
            { value: "select", label: "Dropdown" },
            { value: "swatch", label: "Swatch" },
            { value: "size", label: "Size boxes" },
          ],
        },
        { key: "affects_price", label: "Changes the price", type: "toggle" },
      ],
    },
  },
  { components: { Page: OptionLibraryPage } },
);
`
}

// AdminOptionsPageTSX emits the route file for the option library.
func AdminOptionsPageTSX() string {
	return `"use client";

import { ResourcePage } from "@/components/resource/resource-page";
import { optionResource } from "@/resources/options/options";

export default function OptionsPage() {
  return <ResourcePage resource={optionResource} />;
}
`
}
