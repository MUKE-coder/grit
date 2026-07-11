package scaffold

// Source for desktop frontend primitives shipped in v3.15:
//   format.ts            — formatCurrency / formatDate / humanize / initials   (#33)
//   currency-field.tsx   — CurrencyField with live comma formatting             (#19)
//   searchable-select.tsx — SearchableSelect with typeahead + keyboard nav      (#20)
//   date-field.tsx       — DateField + DateRangeFilter with preset bar          (#21)
//   drawer.tsx           — right-side slide-in panel for forms                  (#22)
//   status-badge.tsx     — status -> colour map with extension hook             (#34)
//   nav-config.ts        — grouped sidebar section config                       (#23)
//   app-shell.tsx        — composes TitleBar + Sidebar + Topbar + content       (#23)

// ═══════════════════════════════════════════════════════════════════
// #33 — Format helpers (lib/format.ts)
// ═══════════════════════════════════════════════════════════════════

func desktopClientFormatLib() string {
	return `// Tiny formatting library every business app rebuilds. Locale + default
// currency are configurable via setFormatConfig({...}) — call once at app
// boot from main.tsx.

interface FormatConfig {
  locale: string;
  currency: string;
}

let config: FormatConfig = {
  locale: "en-US",
  currency: "USD",
};

export function setFormatConfig(next: Partial<FormatConfig>) {
  config = { ...config, ...next };
}

// formatCurrency(4500000) -> "$4,500,000.00"
// formatCurrency(4500000, "UGX") -> "UGX 4,500,000"
//
// We pick "currency" style for known 2-decimal currencies and "decimal +
// prefix" for shilling-style currencies where the trailing .00 is noise.
const NO_DECIMAL_CURRENCIES = new Set(["UGX", "JPY", "KRW", "RWF", "TZS", "VND"]);

export function formatCurrency(amount: number, currency?: string): string {
  const cur = currency || config.currency;
  if (NO_DECIMAL_CURRENCIES.has(cur)) {
    const n = new Intl.NumberFormat(config.locale, { maximumFractionDigits: 0 }).format(amount);
    return cur + " " + n;
  }
  return new Intl.NumberFormat(config.locale, {
    style: "currency",
    currency: cur,
  }).format(amount);
}

// formatDate(value, fmt?) — fmt defaults to "MMM d, yyyy".
// Accepts Date, ISO string, or millis. Empty/null returns "".
const MONTHS_SHORT = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
const MONTHS_LONG = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];

export function formatDate(value: Date | string | number | null | undefined, fmt = "MMM d, yyyy"): string {
  if (value === null || value === undefined || value === "") return "";
  const d = value instanceof Date ? value : new Date(value);
  if (isNaN(d.getTime())) return "";
  // Minimal token formatter — enough for the standard cases. Uses a
  // single-pass replace so "MMMM" (long month) wins over "MMM" (short).
  return fmt
    .replace(/yyyy/g, String(d.getFullYear()))
    .replace(/yy/g, String(d.getFullYear()).slice(-2))
    .replace(/MMMM/g, MONTHS_LONG[d.getMonth()])
    .replace(/MMM/g, MONTHS_SHORT[d.getMonth()])
    .replace(/MM/g, String(d.getMonth() + 1).padStart(2, "0"))
    .replace(/dd/g, String(d.getDate()).padStart(2, "0"))
    .replace(/d/g, String(d.getDate()))
    .replace(/HH/g, String(d.getHours()).padStart(2, "0"))
    .replace(/mm/g, String(d.getMinutes()).padStart(2, "0"))
    .replace(/ss/g, String(d.getSeconds()).padStart(2, "0"));
}

// formatDateTime("2026-05-02T14:30:00Z") -> "May 2, 2026 · 2:30 PM"
export function formatDateTime(value: Date | string | number | null | undefined): string {
  if (value === null || value === undefined || value === "") return "";
  const d = value instanceof Date ? value : new Date(value);
  if (isNaN(d.getTime())) return "";
  const date = formatDate(d, "MMM d, yyyy");
  const time = new Intl.DateTimeFormat(config.locale, {
    hour: "numeric",
    minute: "2-digit",
    hour12: true,
  }).format(d);
  return date + " · " + time;
}

// humanize("checked_in") -> "Checked in". Snake/kebab-case aware.
export function humanize(s: string | null | undefined): string {
  if (!s) return "";
  const tokens = s
    .replace(/[_-]+/g, " ")
    .replace(/([a-z])([A-Z])/g, "$1 $2")
    .toLowerCase()
    .trim();
  if (!tokens) return "";
  return tokens.charAt(0).toUpperCase() + tokens.slice(1);
}

// initials("Abu Seal") -> "AS". Up to 2 characters.
export function initials(name: string | null | undefined): string {
  if (!name) return "";
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "";
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}
`
}

// ═══════════════════════════════════════════════════════════════════
// #19 — CurrencyField
// ═══════════════════════════════════════════════════════════════════

func desktopClientCurrencyField() string {
	return `import { forwardRef, useState, useEffect } from "react";
import { cn } from "@/lib/utils";

// CurrencyInput is the controlled primitive: shows comma-separated
// digits in the box, emits the raw number to onChange. Pasting
// "1,234.56" or "$3,000" both work.
interface CurrencyInputProps {
  value: number | null | undefined;
  onChange: (value: number | null) => void;
  prefix?: string;        // "USD", "UGX", etc.
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  className?: string;
  allowDecimals?: boolean; // default true
}

function formatDisplay(n: number | null | undefined, allowDecimals: boolean): string {
  if (n === null || n === undefined || isNaN(n as number)) return "";
  if (!allowDecimals) {
    return Math.round(n).toLocaleString();
  }
  return (n as number).toLocaleString(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  });
}

function parseRaw(raw: string): number | null {
  if (!raw) return null;
  // Strip everything that isn't a digit, decimal point, or minus sign.
  const cleaned = raw.replace(/[^0-9.\\-]/g, "");
  if (!cleaned || cleaned === "-" || cleaned === ".") return null;
  const n = parseFloat(cleaned);
  return isNaN(n) ? null : n;
}

export function CurrencyInput({
  value,
  onChange,
  prefix,
  placeholder,
  disabled,
  required,
  className,
  allowDecimals = true,
}: CurrencyInputProps) {
  const [display, setDisplay] = useState(() => formatDisplay(value, allowDecimals));
  const [focused, setFocused] = useState(false);

  // Keep display in sync when the parent updates value externally
  // (e.g. RHF reset, server fetch). Skip while focused so we don't
  // fight the user's typing.
  useEffect(() => {
    if (!focused) setDisplay(formatDisplay(value, allowDecimals));
  }, [value, focused, allowDecimals]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value;
    setDisplay(raw);
    onChange(parseRaw(raw));
  };

  const handleBlur = () => {
    setFocused(false);
    setDisplay(formatDisplay(value, allowDecimals));
  };

  const handleFocus = () => {
    setFocused(true);
    // Show raw digits while editing — easier to delete characters.
    if (value !== null && value !== undefined) {
      setDisplay(allowDecimals ? String(value) : String(Math.round(value)));
    }
  };

  return (
    <div className="relative">
      {prefix && (
        <span className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-[12px] font-mono text-foreground-muted">
          {prefix}
        </span>
      )}
      <input
        type="text"
        inputMode={allowDecimals ? "decimal" : "numeric"}
        value={display}
        onChange={handleChange}
        onFocus={handleFocus}
        onBlur={handleBlur}
        placeholder={placeholder}
        disabled={disabled}
        required={required}
        className={cn(
          "w-full h-10 px-3 rounded-lg border border-border bg-surface text-[13.5px] text-foreground placeholder:text-foreground-muted focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15 disabled:bg-surface-2 disabled:text-foreground-muted",
          prefix && "pl-12",
          className,
        )}
      />
    </div>
  );
}

// CurrencyField is the labelled wrapper. Pair with react-hook-form's
// Controller, or use directly as a controlled value+onChange pair.
interface CurrencyFieldProps extends CurrencyInputProps {
  label?: string;
  hint?: string;
  error?: string;
}

export const CurrencyField = forwardRef<HTMLDivElement, CurrencyFieldProps>(
  ({ label, hint, error, required, ...rest }, ref) => (
    <div ref={ref} className="block space-y-1">
      {label && (
        <span className="block text-[12.5px] font-medium text-foreground-secondary">
          {label}
          {required && <span className="text-danger ml-0.5">*</span>}
        </span>
      )}
      <CurrencyInput required={required} {...rest} />
      {error ? (
        <span className="block text-[11.5px] text-danger">{error}</span>
      ) : hint ? (
        <span className="block text-[11.5px] text-foreground-muted">{hint}</span>
      ) : null}
    </div>
  ),
);
CurrencyField.displayName = "CurrencyField";
`
}

// ═══════════════════════════════════════════════════════════════════
// #20 — SearchableSelect
// ═══════════════════════════════════════════════════════════════════

func desktopClientSearchableSelect() string {
	return `import { useState, useRef, useEffect, useMemo, useCallback } from "react";
import { createPortal } from "react-dom";
import { ChevronDown, Check, X, Search } from "lucide-react";
import { cn } from "@/lib/utils";

export interface SelectOption {
  value: string;
  label: string;
  hint?: string; // optional secondary text shown next to the label
}

interface SearchableSelectProps {
  value: string | null | undefined;
  onChange: (value: string | null) => void;
  options: SelectOption[];
  placeholder?: string;
  searchPlaceholder?: string;
  disabled?: boolean;
  required?: boolean;
  clearable?: boolean;
  className?: string;
}

// SearchableSelect is a typeahead combobox that replaces native <select>
// for FK and large-enum fields. The dropdown is portal'd to document.body
// so it isn't clipped by overflow:hidden ancestors (drawers, modals).
//
// Keyboard:
//   Down/Up   - move highlight
//   Enter     - select highlighted
//   Esc       - close
//   Tab       - close + commit nothing
export function SearchableSelect({
  value,
  onChange,
  options,
  placeholder = "Select...",
  searchPlaceholder = "Search...",
  disabled,
  required,
  clearable,
  className,
}: SearchableSelectProps) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [highlight, setHighlight] = useState(0);
  const triggerRef = useRef<HTMLButtonElement>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [pos, setPos] = useState({ top: 0, left: 0, width: 0 });

  const updatePosition = useCallback(() => {
    if (triggerRef.current) {
      const rect = triggerRef.current.getBoundingClientRect();
      setPos({ top: rect.bottom + 4, left: rect.left, width: rect.width });
    }
  }, []);

  useEffect(() => {
    if (!open) return;
    updatePosition();
    setSearch("");
    setHighlight(0);
    // Focus the search input after the portal mounts.
    requestAnimationFrame(() => inputRef.current?.focus());

    const handleClickOutside = (e: MouseEvent) => {
      const target = e.target as Node;
      if (
        triggerRef.current && !triggerRef.current.contains(target) &&
        dropdownRef.current && !dropdownRef.current.contains(target)
      ) {
        setOpen(false);
      }
    };
    const handleScroll = () => updatePosition();
    document.addEventListener("mousedown", handleClickOutside);
    window.addEventListener("scroll", handleScroll, true);
    window.addEventListener("resize", updatePosition);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      window.removeEventListener("scroll", handleScroll, true);
      window.removeEventListener("resize", updatePosition);
    };
  }, [open, updatePosition]);

  const filtered = useMemo(() => {
    if (!search) return options;
    const q = search.toLowerCase();
    return options.filter(
      (o) =>
        o.label.toLowerCase().includes(q) ||
        (o.hint || "").toLowerCase().includes(q),
    );
  }, [options, search]);

  // Keep highlight in range as the filter shrinks.
  useEffect(() => {
    if (highlight >= filtered.length) setHighlight(Math.max(0, filtered.length - 1));
  }, [filtered, highlight]);

  const selected = useMemo(() => options.find((o) => o.value === value), [options, value]);

  const commit = (opt: SelectOption | undefined) => {
    if (!opt) return;
    onChange(opt.value);
    setOpen(false);
  };

  const handleKey = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setHighlight((h) => Math.min(filtered.length - 1, h + 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setHighlight((h) => Math.max(0, h - 1));
    } else if (e.key === "Enter") {
      e.preventDefault();
      commit(filtered[highlight]);
    } else if (e.key === "Escape") {
      e.preventDefault();
      setOpen(false);
    }
  };

  const dropdown = open && createPortal(
    <div
      ref={dropdownRef}
      className="fixed z-[9999] rounded-lg border border-border bg-surface shadow-xl"
      style={{ top: pos.top, left: pos.left, width: pos.width }}
    >
      <div className="p-2 border-b border-border-subtle relative">
        <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-foreground-muted" />
        <input
          ref={inputRef}
          type="search"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          onKeyDown={handleKey}
          placeholder={searchPlaceholder}
          className="w-full h-8 pl-8 pr-2.5 rounded-md border border-border bg-surface-2 text-[12.5px] placeholder:text-foreground-muted focus:border-accent focus:outline-none"
        />
      </div>
      <div className="max-h-60 overflow-y-auto py-1">
        {filtered.length === 0 ? (
          <div className="px-3 py-4 text-center text-[12.5px] text-foreground-muted">
            No matches
          </div>
        ) : (
          filtered.map((opt, idx) => (
            <button
              key={opt.value}
              type="button"
              onMouseEnter={() => setHighlight(idx)}
              onClick={() => commit(opt)}
              className={cn(
                "w-full flex items-center justify-between px-3 py-1.5 text-left text-[13px] transition-colors",
                idx === highlight ? "bg-accent/10 text-foreground" : "text-foreground-secondary",
              )}
            >
              <span className="flex-1 truncate">
                {opt.label}
                {opt.hint && (
                  <span className="ml-2 text-[11.5px] text-foreground-muted">{opt.hint}</span>
                )}
              </span>
              {opt.value === value && <Check className="h-3.5 w-3.5 text-accent" />}
            </button>
          ))
        )}
      </div>
    </div>,
    document.body,
  );

  return (
    <>
      <button
        ref={triggerRef}
        type="button"
        disabled={disabled}
        onClick={() => setOpen((o) => !o)}
        className={cn(
          "w-full h-10 flex items-center justify-between px-3 rounded-lg border border-border bg-surface text-[13.5px] text-left transition-colors",
          "hover:border-foreground-muted focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15",
          "disabled:bg-surface-2 disabled:text-foreground-muted disabled:cursor-not-allowed",
          className,
        )}
      >
        <span className={cn("truncate", !selected && "text-foreground-muted")}>
          {selected ? selected.label : placeholder}
        </span>
        <span className="flex items-center gap-1 shrink-0">
          {clearable && selected && !disabled && (
            <span
              role="button"
              tabIndex={0}
              onClick={(e) => {
                e.stopPropagation();
                onChange(null);
              }}
              className="h-5 w-5 rounded-full hover:bg-surface-hover inline-flex items-center justify-center text-foreground-muted hover:text-foreground"
              aria-label="Clear"
            >
              <X className="h-3 w-3" />
            </span>
          )}
          <ChevronDown className="h-4 w-4 text-foreground-muted" />
        </span>
      </button>
      {dropdown}
      {required && !value && (
        // Hidden anchor so the form's HTML5 validation can require this.
        <input
          tabIndex={-1}
          aria-hidden
          className="sr-only"
          required
          value=""
          onChange={() => {}}
        />
      )}
    </>
  );
}

// SearchableSelectField: labelled wrapper for forms.
interface SearchableSelectFieldProps extends SearchableSelectProps {
  label?: string;
  hint?: string;
  error?: string;
}

export function SearchableSelectField({
  label,
  hint,
  error,
  required,
  ...rest
}: SearchableSelectFieldProps) {
  return (
    <div className="block space-y-1">
      {label && (
        <span className="block text-[12.5px] font-medium text-foreground-secondary">
          {label}
          {required && <span className="text-danger ml-0.5">*</span>}
        </span>
      )}
      <SearchableSelect required={required} {...rest} />
      {error ? (
        <span className="block text-[11.5px] text-danger">{error}</span>
      ) : hint ? (
        <span className="block text-[11.5px] text-foreground-muted">{hint}</span>
      ) : null}
    </div>
  );
}
`
}

// ═══════════════════════════════════════════════════════════════════
// #21 — DateField + DateRangeFilter
// ═══════════════════════════════════════════════════════════════════

func desktopClientDateField() string {
	return `import { forwardRef } from "react";
import { Calendar } from "lucide-react";
import { cn } from "@/lib/utils";

// DateField wraps a native <input type="date"> with the standard
// label/hint/error chrome. We use the native picker on purpose — it's
// the only one with reliable accessibility, RTL, and i18n.
//
// Value is a YYYY-MM-DD string (the format <input type="date"> emits).
// Empty string = no date selected.
interface DateFieldProps {
  value: string;
  onChange: (value: string) => void;
  label?: string;
  hint?: string;
  error?: string;
  required?: boolean;
  disabled?: boolean;
  min?: string;
  max?: string;
  className?: string;
}

export const DateField = forwardRef<HTMLInputElement, DateFieldProps>(
  ({ value, onChange, label, hint, error, required, disabled, min, max, className }, ref) => (
    <div className={cn("block space-y-1", className)}>
      {label && (
        <span className="block text-[12.5px] font-medium text-foreground-secondary">
          {label}
          {required && <span className="text-danger ml-0.5">*</span>}
        </span>
      )}
      <div className="relative">
        <Calendar className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-foreground-muted" />
        <input
          ref={ref}
          type="date"
          value={value || ""}
          onChange={(e) => onChange(e.target.value)}
          required={required}
          disabled={disabled}
          min={min}
          max={max}
          className="w-full h-10 pl-10 pr-3 rounded-lg border border-border bg-surface text-[13.5px] text-foreground focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15 disabled:bg-surface-2 disabled:text-foreground-muted"
        />
      </div>
      {error ? (
        <span className="block text-[11.5px] text-danger">{error}</span>
      ) : hint ? (
        <span className="block text-[11.5px] text-foreground-muted">{hint}</span>
      ) : null}
    </div>
  ),
);
DateField.displayName = "DateField";

// ─── DateRangeFilter ───────────────────────────────────────────────

export interface DateRange {
  from: string; // YYYY-MM-DD or ""
  to: string;   // YYYY-MM-DD or ""
}

export type DatePreset =
  | "today"
  | "last7"
  | "last30"
  | "thisMonth"
  | "lastMonth"
  | "last90"
  | "thisYear"
  | "all"
  | "custom";

const PRESETS: { key: DatePreset; label: string }[] = [
  { key: "today", label: "Today" },
  { key: "last7", label: "Last 7" },
  { key: "last30", label: "Last 30" },
  { key: "thisMonth", label: "This month" },
  { key: "lastMonth", label: "Last month" },
  { key: "last90", label: "Last 90" },
  { key: "thisYear", label: "This year" },
  { key: "all", label: "All" },
];

function pad(n: number) {
  return String(n).padStart(2, "0");
}

function isoDate(d: Date): string {
  return d.getFullYear() + "-" + pad(d.getMonth() + 1) + "-" + pad(d.getDate());
}

// presetRange returns { from, to } ISO strings for a named preset.
// "all" returns empty strings ⇒ filter not applied.
export function presetRange(preset: DatePreset): DateRange {
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());

  switch (preset) {
    case "today":
      return { from: isoDate(today), to: isoDate(today) };
    case "last7": {
      const from = new Date(today);
      from.setDate(today.getDate() - 6);
      return { from: isoDate(from), to: isoDate(today) };
    }
    case "last30": {
      const from = new Date(today);
      from.setDate(today.getDate() - 29);
      return { from: isoDate(from), to: isoDate(today) };
    }
    case "thisMonth": {
      const from = new Date(today.getFullYear(), today.getMonth(), 1);
      return { from: isoDate(from), to: isoDate(today) };
    }
    case "lastMonth": {
      const from = new Date(today.getFullYear(), today.getMonth() - 1, 1);
      const to = new Date(today.getFullYear(), today.getMonth(), 0); // day 0 = last of prev month
      return { from: isoDate(from), to: isoDate(to) };
    }
    case "last90": {
      const from = new Date(today);
      from.setDate(today.getDate() - 89);
      return { from: isoDate(from), to: isoDate(today) };
    }
    case "thisYear": {
      const from = new Date(today.getFullYear(), 0, 1);
      return { from: isoDate(from), to: isoDate(today) };
    }
    case "all":
    default:
      return { from: "", to: "" };
  }
}

// matchesPreset works out which preset (if any) the current value
// corresponds to. Lets the chip bar highlight the right preset when
// the parent feeds in a known range.
function matchesPreset(value: DateRange): DatePreset {
  for (const p of PRESETS) {
    const r = presetRange(p.key);
    if (r.from === value.from && r.to === value.to) return p.key;
  }
  return "custom";
}

interface DateRangeFilterProps {
  value: DateRange;
  onChange: (value: DateRange) => void;
}

export function DateRangeFilter({ value, onChange }: DateRangeFilterProps) {
  const active = matchesPreset(value);
  return (
    <div className="flex flex-wrap items-center gap-1.5">
      {PRESETS.map((p) => {
        const isActive = active === p.key;
        return (
          <button
            key={p.key}
            type="button"
            onClick={() => onChange(presetRange(p.key))}
            className={cn(
              "h-7 px-2.5 rounded-full text-[12px] font-medium border transition-colors",
              isActive
                ? "bg-accent/15 border-accent/30 text-accent"
                : "bg-surface border-border text-foreground-secondary hover:bg-surface-hover hover:text-foreground",
            )}
          >
            {p.label}
          </button>
        );
      })}
      <div className="flex items-center gap-1">
        <input
          type="date"
          value={value.from}
          onChange={(e) => onChange({ ...value, from: e.target.value })}
          className="h-7 px-2 rounded-md border border-border bg-surface text-[12px] text-foreground"
        />
        <span className="text-[11px] text-foreground-muted">to</span>
        <input
          type="date"
          value={value.to}
          onChange={(e) => onChange({ ...value, to: e.target.value })}
          className="h-7 px-2 rounded-md border border-border bg-surface text-[12px] text-foreground"
        />
      </div>
    </div>
  );
}
`
}

// ═══════════════════════════════════════════════════════════════════
// #22 — Drawer
// ═══════════════════════════════════════════════════════════════════

func desktopClientDrawer() string {
	return `import { useEffect } from "react";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";

// Drawer is the right-edge slide-in panel. Used for create/edit forms,
// review panels, anything that's heavier than a popover but lighter
// than a full page route. Closes on Esc + backdrop click + the X button.
interface DrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title?: string;
  description?: string;
  width?: "sm" | "md" | "lg" | "xl";
  children: React.ReactNode;
  // Footer sticks to the bottom (typical "Cancel + Save" row). When
  // omitted, the children are responsible for their own actions.
  footer?: React.ReactNode;
}

const WIDTHS = {
  sm: "w-[360px]",
  md: "w-[480px]",
  lg: "w-[640px]",
  xl: "w-[860px]",
};

export function Drawer({
  open,
  onOpenChange,
  title,
  description,
  width = "md",
  children,
  footer,
}: DrawerProps) {
  // Close on Esc.
  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onOpenChange(false);
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [open, onOpenChange]);

  if (!open) return null;

  return (
    <>
      <div
        className="fixed inset-0 z-40 bg-black/40"
        onClick={() => onOpenChange(false)}
      />
      <aside
        className={cn(
          "fixed right-0 top-0 z-50 h-full bg-surface border-l border-border flex flex-col",
          WIDTHS[width],
        )}
      >
        {(title || description) && (
          <header className="flex items-start justify-between px-5 py-4 border-b border-border">
            <div>
              {title && <h2 className="text-[15px] font-semibold text-foreground">{title}</h2>}
              {description && (
                <p className="text-[12.5px] text-foreground-muted mt-0.5">{description}</p>
              )}
            </div>
            <button
              type="button"
              onClick={() => onOpenChange(false)}
              className="text-foreground-muted hover:text-foreground"
              aria-label="Close"
            >
              <X className="h-4 w-4" />
            </button>
          </header>
        )}
        <div className="flex-1 overflow-y-auto px-5 py-4">{children}</div>
        {footer && (
          <footer className="border-t border-border px-5 py-3">{footer}</footer>
        )}
      </aside>
    </>
  );
}
`
}

// ═══════════════════════════════════════════════════════════════════
// #34 — StatusBadge
// ═══════════════════════════════════════════════════════════════════

func desktopClientStatusBadge() string {
	return `import { cn } from "@/lib/utils";
import { humanize } from "@/lib/format";

// StatusBadge maps a status string to a coloured pill. Defaults cover
// the common cases (paid, pending, overdue, draft, active, etc).
// Apps extend the map via setStatusVariants({...}).

export type StatusVariant = "success" | "warning" | "danger" | "info" | "muted" | "accent";

const DEFAULT_MAP: Record<string, StatusVariant> = {
  // success — done / good state
  paid: "success",
  active: "success",
  completed: "success",
  approved: "success",
  delivered: "success",
  resolved: "success",
  // warning — needs attention but not broken
  pending: "warning",
  partial: "warning",
  processing: "warning",
  in_review: "warning",
  // danger — bad state
  overdue: "danger",
  cancelled: "danger",
  failed: "danger",
  rejected: "danger",
  expired: "danger",
  // muted — neutral / archived
  draft: "muted",
  archived: "muted",
  inactive: "muted",
  closed: "muted",
  // accent — informational, in-flight
  checked_in: "accent",
  in_progress: "accent",
  scheduled: "accent",
  new: "accent",
};

let statusMap: Record<string, StatusVariant> = { ...DEFAULT_MAP };

// setStatusVariants extends or overrides the global status map. Call
// once at app boot from main.tsx if you have domain-specific statuses.
export function setStatusVariants(extra: Record<string, StatusVariant>) {
  statusMap = { ...statusMap, ...Object.fromEntries(
    Object.entries(extra).map(([k, v]) => [k.toLowerCase(), v]),
  ) };
}

const VARIANT_CLASSES: Record<StatusVariant, string> = {
  success: "bg-success/15 text-success",
  warning: "bg-warning/15 text-warning",
  danger:  "bg-danger/15 text-danger",
  info:    "bg-info/15 text-info",
  muted:   "bg-surface-2 text-foreground-muted",
  accent:  "bg-accent/15 text-accent",
};

interface StatusBadgeProps {
  status: string | null | undefined;
  // Override: pass a variant to ignore the map and use this colour.
  variant?: StatusVariant;
  // Override the rendered label (defaults to humanize(status)).
  label?: string;
  className?: string;
}

export function StatusBadge({ status, variant, label, className }: StatusBadgeProps) {
  const v: StatusVariant =
    variant || statusMap[(status || "").toLowerCase()] || "muted";
  const text = label ?? humanize(status);
  return (
    <span
      className={cn(
        "inline-flex items-center h-5 px-2 rounded-full text-[11px] font-medium",
        VARIANT_CLASSES[v],
        className,
      )}
    >
      {text || "—"}
    </span>
  );
}
`
}

// ═══════════════════════════════════════════════════════════════════
// #23 — nav-config.ts + AppShell + grouped sidebar
// ═══════════════════════════════════════════════════════════════════

func desktopClientNavConfig() string {
	return `import {
  Home,
  User,
  Settings,
  Box,
  RefreshCw,
  Users,
  FileText,
  Activity,
  MessageSquare,
  Bell,
  TrendingUp,
  Shield,
  LayoutGrid,
  type LucideIcon,
} from "lucide-react";

// Single source of truth for the desktop sidebar nav. Add new sections
// or items here; <Sidebar> reads this and renders grouped nav.
//
// A "section" is a logical cluster of related items (Properties,
// Operations, Admin, etc). The first section is rendered without a
// header so the dashboard / home item sits clean at the top.

export interface NavItem {
  to: string;
  label: string;
  icon: LucideIcon;
  badge?: string;
}

export interface NavSection {
  title?: string; // omit on the first section for the unbranded "main" group
  items: NavItem[];
}

export const NAV_SECTIONS: NavSection[] = [
  {
    items: [
      { to: "/app", label: "Dashboard", icon: Home },
    ],
  },
  {
    title: "Manage",
    items: [
      { to: "/app/system/users", label: "Users", icon: Users },
      // grit generate resource injects generated resources here. Box is a
      // shared icon so no per-resource import is needed.
      // grit:nav
    ],
  },
  {
    title: "Internal",
    items: [
      { to: "/app/system/activity", label: "Activity", icon: Activity },
      { to: "/app/system/support", label: "Support", icon: MessageSquare },
      { to: "/app/system/notifications", label: "Notifications", icon: Bell },
    ],
  },
  {
    title: "System",
    items: [
      { to: "/app/system/blogs", label: "Blogs", icon: FileText },
      { to: "/app/system/dashboard-settings", label: "Dashboard settings", icon: Settings },
      { to: "/app/system/health", label: "System Health", icon: Activity },
      { to: "/app/system/performance", label: "Performance", icon: TrendingUp },
      { to: "/app/system/security", label: "Security", icon: Shield },
      { to: "/app/system", label: "System Hub", icon: LayoutGrid },
    ],
  },
  {
    title: "Account",
    items: [
      { to: "/app/sync", label: "Sync", icon: RefreshCw },
      { to: "/app/profile", label: "Profile", icon: User },
      { to: "/app/settings", label: "Settings", icon: Settings },
    ],
  },
];

void Box;
`
}

// desktopClientSidebarV2 is the grouped, COLLAPSIBLE sidebar — mirrors the
// admin panel's CollapsibleSidebar: a brand header with a chevron toggle,
// section labels, active states, and an icon-only collapsed mode (with title
// tooltips). Collapse state persists in localStorage. Sections + items come
// from nav-config.ts so `grit generate resource` only edits one file.
func desktopClientSidebarV2() string {
	return `import { useEffect, useState } from "react";
import { Link, useRouterState } from "@tanstack/react-router";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { brand } from "@repo/shared/brand.config";
import { cn } from "@/lib/utils";
import { NAV_SECTIONS } from "@/lib/nav-config";

export function Sidebar() {
  const pathname = useRouterState({ select: (s) => s.location.pathname });
  const [collapsed, setCollapsed] = useState<boolean>(() => {
    if (typeof window !== "undefined") {
      return localStorage.getItem("grit-sidebar-collapsed") === "1";
    }
    return false;
  });

  useEffect(() => {
    localStorage.setItem("grit-sidebar-collapsed", collapsed ? "1" : "0");
  }, [collapsed]);

  return (
    <aside
      // min-w-0 + overflow-hidden so the collapsed width (w-16) actually
      // takes effect — otherwise flexbox floors the item at its content's
      // min-content width and the collapse doesn't visibly shrink.
      className={cn(
        "shrink-0 min-w-0 overflow-hidden border-r border-border-subtle bg-surface flex flex-col transition-[width] duration-200",
        collapsed ? "w-16" : "w-sidebar",
      )}
    >
      {/* Brand + collapse toggle */}
      <div className="flex h-14 shrink-0 items-center gap-2 border-b border-border-subtle px-3">
        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-accent text-[13px] font-bold text-white">
          {brand.logo.text}
        </div>
        {!collapsed && (
          <span className="flex-1 truncate text-[14px] font-semibold text-foreground">{brand.name}</span>
        )}
        <button
          type="button"
          onClick={() => setCollapsed((v) => !v)}
          title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
          className={cn(
            "flex h-6 w-6 items-center justify-center rounded-md text-foreground-secondary hover:bg-surface-hover hover:text-foreground",
            collapsed && "mx-auto",
          )}
        >
          {collapsed ? <ChevronRight className="h-4 w-4" /> : <ChevronLeft className="h-4 w-4" />}
        </button>
      </div>

      <nav className="flex-1 overflow-y-auto p-3">
        {NAV_SECTIONS.filter((s) => s.items.length > 0).map((section, idx) => (
          <div key={idx} className={cn(idx > 0 && "mt-4")}>
            {section.title && !collapsed && (
              <h3 className="mb-1 px-3 text-[10px] font-semibold uppercase tracking-wider text-foreground-muted">
                {section.title}
              </h3>
            )}
            {section.title && collapsed && idx > 0 && (
              <div className="my-2 mx-auto h-px w-6 bg-border-subtle" />
            )}
            <div className="space-y-0.5">
              {section.items.map((item) => {
                const Icon = item.icon;
                const active =
                  pathname === item.to ||
                  (item.to !== "/app" && pathname.startsWith(item.to));
                return (
                  <Link
                    key={item.to}
                    to={item.to}
                    title={collapsed ? item.label : undefined}
                    className={cn(
                      "flex items-center gap-3 rounded-lg text-[13px] font-medium transition-colors",
                      collapsed ? "justify-center px-0 py-2" : "justify-between px-3 py-2",
                      active
                        ? "bg-accent/10 text-accent"
                        : "text-foreground-secondary hover:bg-surface-hover hover:text-foreground",
                    )}
                  >
                    <span className={cn("flex items-center gap-3 min-w-0", collapsed && "gap-0")}>
                      <Icon className="h-4 w-4 shrink-0" />
                      {!collapsed && <span className="truncate">{item.label}</span>}
                    </span>
                    {!collapsed && item.badge && (
                      <span className="shrink-0 inline-flex h-5 min-w-[20px] items-center justify-center rounded-full bg-accent/15 px-1.5 text-[10px] font-semibold text-accent">
                        {item.badge}
                      </span>
                    )}
                  </Link>
                );
              })}
            </div>
          </div>
        ))}
      </nav>

      {!collapsed && (
        <div className="border-t border-border-subtle px-4 py-3 text-[11px] text-foreground-muted">
          Built with Grit
        </div>
      )}
    </aside>
  );
}
`
}

func desktopClientAppShell() string {
	return `import { useState } from "react";
import { TitleBar } from "@/components/layout/title-bar";
import { Sidebar } from "@/components/layout/sidebar";
import { Topbar } from "@/components/layout/topbar";
import { CommandPalette } from "@/components/layout/command-palette";
import { useShortcuts } from "@/lib/use-shortcuts";

// AppShell is the standard authenticated dashboard layout. Composes
// TitleBar + Sidebar + Topbar + scrollable content + Command Palette
// — including Cmd/Ctrl-K binding.
//
// Wrap your dashboard route Outlet with this. Sections in the sidebar
// come from lib/nav-config.ts so adding a section is one config edit.
//
// Usage in routes/_app.tsx:
//   <AppShell><Outlet /></AppShell>
export function AppShell({ children }: { children: React.ReactNode }) {
  const [paletteOpen, setPaletteOpen] = useState(false);

  useShortcuts({
    "mod+k": () => setPaletteOpen(true),
  });

  return (
    <div className="h-screen flex flex-col bg-background">
      <TitleBar />
      <div className="flex-1 flex overflow-hidden">
        <Sidebar />
        <main className="flex-1 flex flex-col overflow-hidden">
          <Topbar onOpenPalette={() => setPaletteOpen(true)} />
          <div className="flex-1 overflow-auto p-content">
            {children}
          </div>
        </main>
      </div>
      <CommandPalette open={paletteOpen} onOpenChange={setPaletteOpen} />
    </div>
  );
}
`
}
