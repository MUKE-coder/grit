"use client";

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
/** The columns the matrix draws for you. */
export type VariantColumnKey = "sku" | "stock" | "price" | "override" | "active";

/** What a cell renderer is handed. */
export interface VariantCellContext {
  /** This row's unsaved edits, if it has any. */
  draft?: VariantDraft;
  /**
   * Stage an edit on this row.
   *
   * It joins the same Save button the built-in cells feed, so a custom column
   * is not a second way to write: one click still sends one PATCH per row.
   */
  patch: (patch: VariantDraft) => void;
  /** What this combination costs with no override, already resolved. */
  resolved: number;
}

/**
 * One column, added or overridden.
 *
 * Keyed the way the resource definition's own column overrides are keyed: a
 * built-in key patches that column, any other key adds one. Deliberately the
 * same shape, so this is a pattern you already know rather than a second one
 * to learn.
 *
 *   <VariantMatrix
 *     {...props}
 *     columns={{
 *       // add: a photo of this combination, which the variant already stores
 *       images: {
 *         label: "Photo",
 *         after: "sku",
 *         cell: (variant) => <VariantThumb images={variant.images} />,
 *       },
 *       // override: same cell, different header
 *       sku: { label: "Barcode" },
 *       // drop one you do not use
 *       override: { hidden: true },
 *     }}
 *   />
 */
export interface VariantColumn {
  /** Header text. Defaults to the built-in label, or to the key. */
  label?: string;
  /** Drop the column. Only meaningful on a built-in. */
  hidden?: boolean;
  /** Render the cell. Without one, a built-in keeps its own renderer. */
  cell?: (variant: AdminVariant, ctx: VariantCellContext) => React.ReactNode;
  /** Place an added column after this one. Appended when absent. */
  after?: VariantColumnKey | string;
}

const BUILT_IN_COLUMNS: VariantColumnKey[] = ["sku", "stock", "price", "override", "active"];

const BUILT_IN_LABELS: Record<VariantColumnKey, string> = {
  sku: "SKU",
  stock: "Stock",
  price: "Price",
  override: "Override",
  active: "Active",
};

export function VariantMatrix<T>({
  resource,
  id,
  controller,
  columns,
}: ResourceDetailPartProps<T> & {
  /** Add columns, or patch the built-in ones. See VariantColumn. */
  columns?: Record<string, VariantColumn>;
}) {
  const matrix = useVariantMatrix(resource.slug, id);
  const library = useOptions();
  const setOptions = useSetResourceOptions(resource.slug, id);
  const generate = useGenerateMatrix(resource.slug, id);
  const save = useSaveVariants(resource.name, resource.slug, id);

  const [picking, setPicking] = useState(false);
  const [chosen, setChosen] = useState<string[]>([]);
  // The values of those axes this record offers. Empty for an axis means all
  // of them, so a record nobody has narrowed stores nothing at all.
  const [chosenValues, setChosenValues] = useState<string[]>([]);
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

  /**
   * The columns to draw: the built-ins, patched, with any added ones slotted
   * in. Option columns are not in here because they are data rather than
   * configuration, and they always come first.
   */
  const columnList = useMemo(() => {
    const out: Array<{ key: string; label: string; custom?: VariantColumn }> = [];
    for (const key of BUILT_IN_COLUMNS) {
      const custom = columns?.[key];
      if (custom?.hidden) continue;
      out.push({ key, label: custom?.label ?? BUILT_IN_LABELS[key], custom });
    }
    for (const [key, custom] of Object.entries(columns ?? {})) {
      if ((BUILT_IN_COLUMNS as string[]).includes(key) || custom.hidden) continue;
      const entry = { key, label: custom.label ?? key, custom };
      const at = custom.after ? out.findIndex((col) => col.key === custom.after) : -1;
      if (at >= 0) out.splice(at + 1, 0, entry);
      else out.push(entry);
    }
    return out;
  }, [columns]);

  /** Price is the one built-in that is text rather than a control. */
  function cellClass(key: string) {
    return key === "price"
      ? "px-3 py-2.5 whitespace-nowrap tabular-nums text-text-secondary"
      : "px-3 py-2.5";
  }

  /**
   * The stock renderers.
   *
   * A column with no cell of its own falls through to here, and an added
   * column with no cell renders nothing rather than throwing, because a
   * half-written customisation should not take the page down.
   */
  function builtInCell(key: string, variant: AdminVariant, draft?: VariantDraft) {
    switch (key) {
      case "sku":
        return (
          // Wide enough for a real SKU. The generated ones carry the record's
          // slug and every option value, and a box showing the first fifteen
          // characters of that is a box you cannot check anything against.
          <div className="w-56">
            <input
              className={inputClasses({ inputSize: "sm", className: "font-mono text-xs" })}
              placeholder="unset"
              value={draft?.sku ?? variant.sku ?? ""}
              onChange={(e) => patchDraft(variant.id, { sku: e.target.value })}
            />
          </div>
        );
      case "stock":
        return (
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
        );
      case "price":
        return formatCurrency(effectivePrice(variant, draft, computed));
      case "override":
        return (
          <div className="w-24">
            <input
              className={inputClasses({ inputSize: "sm", className: "tabular-nums" })}
              inputMode="decimal"
              /* Empty means resolved, and the placeholder says what that
                 resolves to, so clearing the box is never a guess about what
                 the price becomes. */
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
        );
      case "active":
        return (
          <button
            type="button"
            aria-label={(draft?.active ?? variant.active) ? "Deactivate" : "Activate"}
            onClick={() => patchDraft(variant.id, { active: !(draft?.active ?? variant.active) })}
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
        );
      default:
        return null;
    }
  }

  function patchDraft(variantID: string, patch: VariantDraft) {
    setDrafts((prev) => ({ ...prev, [variantID]: { ...prev[variantID], ...patch } }));
  }

  function openPicker() {
    setChosen(options.map((option) => option.id));
    // options here are already narrowed to what this record offers, so the
    // ticks start where the record is rather than at "everything".
    setChosenValues(options.flatMap((option) => (option.values ?? []).map((value) => value.id)));
    setPicking(true);
  }

  function toggleValue(valueID: string) {
    setChosenValues((prev) =>
      prev.includes(valueID) ? prev.filter((each) => each !== valueID) : [...prev, valueID],
    );
  }

  function toggleOption(optionID: string) {
    // Selection order is the display order the storefront gets, so a click
    // appends rather than slotting the option back into library order.
    setChosen((prev) =>
      prev.includes(optionID) ? prev.filter((each) => each !== optionID) : [...prev, optionID],
    );
    // Ticking an axis offers all of it until somebody says otherwise;
    // unticking one takes its values with it.
    const option = (library.data ?? []).find((each) => each.id === optionID);
    const valueIDs = (option?.values ?? []).map((value) => value.id);
    setChosenValues((prev) =>
      chosen.includes(optionID)
        ? prev.filter((each) => !valueIDs.includes(each))
        : [...prev, ...valueIDs.filter((each) => !prev.includes(each))],
    );
  }

  function applyOptions() {
    const offeredNow = options.flatMap((option) => (option.values ?? []).map((value) => value.id));
    const sameValues =
      chosenValues.length === offeredNow.length &&
      [...chosenValues].sort().every((each, i) => [...offeredNow].sort()[i] === each);
    const unchanged =
      chosen.length === options.length &&
      chosen.every((each, i) => options[i]?.id === each) &&
      sameValues;
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
    await setOptions.mutateAsync({ optionIDs: chosen, valueIDs: chosenValues });
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
          chosenValues={chosenValues}
          onToggle={toggleOption}
          onToggleValue={toggleValue}
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
                {columnList.map((col) => (
                  <th key={col.key} className="px-3 py-3 font-medium">
                    {col.label}
                  </th>
                ))}
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

                    {columnList.map((col) => (
                      <td key={col.key} className={cellClass(col.key)}>
                        {col.custom?.cell
                          ? col.custom.cell(variant, {
                              draft,
                              patch: (patch) => patchDraft(variant.id, patch),
                              resolved: computed(variant),
                            })
                          : builtInCell(col.key, variant, draft)}
                      </td>
                    ))}

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
  chosenValues,
  onToggle,
  onToggleValue,
  onCancel,
  onApply,
  saving,
}: {
  library: AdminOption[];
  loading: boolean;
  chosen: string[];
  chosenValues: string[];
  onToggle: (id: string) => void;
  onToggleValue: (id: string) => void;
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

      {chosen.length > 0 && (
        <div className="mt-5 space-y-4">
          <p className="text-xs text-text-muted">
            And which of their values this record offers. Options are shared by every record, so
            this is how one shirt comes in ecru and navy while another comes in black and sand.
            Leave an axis fully ticked to offer all of it, including values added later.
          </p>
          {chosen.map((optionID) => {
            const option = library.find((each) => each.id === optionID);
            if (!option) return null;
            const values = option.values ?? [];
            const picked = values.filter((value) => chosenValues.includes(value.id));
            return (
              <div key={option.id}>
                <div className="mb-2 flex items-center gap-2 text-xs">
                  <span className="font-medium text-foreground">{option.name}</span>
                  <span className="text-text-muted">
                    {picked.length === values.length
                      ? "all " + values.length
                      : picked.length + " of " + values.length}
                  </span>
                </div>
                <div className="flex flex-wrap gap-2">
                  {values.map((value) => {
                    const on = chosenValues.includes(value.id);
                    return (
                      <button
                        key={value.id}
                        type="button"
                        onClick={() => onToggleValue(value.id)}
                        className={
                          "inline-flex items-center gap-2 rounded-lg border px-2.5 py-1.5 text-xs transition-colors " +
                          (on
                            ? "border-accent bg-accent/10 text-foreground"
                            : "border-border text-text-muted hover:border-accent/40 hover:text-foreground")
                        }
                      >
                        {option.kind === "swatch" && value.swatch && (
                          <span
                            className="h-3 w-3 shrink-0 rounded-full border border-border"
                            style={{ backgroundColor: value.swatch }}
                          />
                        )}
                        {value.label}
                      </button>
                    );
                  })}
                  {values.length === 0 && (
                    <span className="text-xs text-text-muted">
                      This axis has no values yet, so no combination can include it.
                    </span>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}

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
