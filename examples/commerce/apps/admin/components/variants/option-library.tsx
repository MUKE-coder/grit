"use client";

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

/**
 * The placeholder for a new value.
 *
 * Deliberately something the seeded library does NOT contain. The first
 * version used "Black", which is also the first chip on the Colour card, so
 * the add row looked like a duplicate of a value already there rather than an
 * empty field waiting for input.
 */
function newValueHint(kind: string): string {
  if (kind === "size") return "e.g. XXL";
  if (kind === "swatch") return "e.g. Red";
  return "e.g. 512GB";
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
              title={"Delete " + value.label + ". Refused while a variant is built on it."}
              disabled={removeValue.isPending}
              onClick={() => removeValue.mutate(value.id)}
              className="rounded p-1.5 text-text-muted transition-colors hover:bg-bg-hover hover:text-danger disabled:opacity-50"
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </span>
        ))}
      </div>

      <form
        onSubmit={submitValue}
        className="flex flex-wrap items-end gap-2 border-t border-border px-6 py-3"
      >
        <div className="w-44">
          <label
            htmlFor={"add-value-" + option.id}
            className="mb-1 block text-[11px] font-medium uppercase tracking-wide text-text-muted"
          >
            Add a value
          </label>
          <input
            id={"add-value-" + option.id}
            className={inputClasses({ inputSize: "sm" })}
            placeholder={newValueHint(option.kind)}
            value={label}
            onChange={(e) => setLabel(e.target.value)}
          />
        </div>
        {option.kind === "swatch" && (
          <div>
            <label
              htmlFor={"add-swatch-" + option.id}
              className="mb-1 block text-[11px] font-medium uppercase tracking-wide text-text-muted"
            >
              Swatch
            </label>
            <input
              id={"add-swatch-" + option.id}
              type="color"
              aria-label="Swatch colour"
              className="h-8 w-12 cursor-pointer rounded-lg border border-border bg-bg-secondary p-1"
              value={swatch}
              onChange={(e) => setSwatch(e.target.value)}
            />
          </div>
        )}
        {option.affects_price && (
          <div className="w-28">
            <label
              htmlFor={"add-delta-" + option.id}
              className="mb-1 block text-[11px] font-medium uppercase tracking-wide text-text-muted"
            >
              Price change
            </label>
            <input
              id={"add-delta-" + option.id}
              className={inputClasses({ inputSize: "sm", className: "tabular-nums" })}
              inputMode="decimal"
              placeholder="+ / -"
              value={delta}
              onChange={(e) => setDelta(e.target.value)}
            />
          </div>
        )}
        <button
          type="submit"
          disabled={addValue.isPending || !label.trim()}
          title={label.trim() ? "Add this value" : "Type a name first"}
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
