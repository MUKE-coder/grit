"use client";

import { useId, useEffect, useMemo, useRef, useState, useCallback } from "react";
import { createPortal } from "react-dom";
import type { FieldDefinition } from "@/lib/resource";
import { Calendar, ChevronLeft, ChevronRight } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";

interface DateFieldProps {
  field: FieldDefinition;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

/* ─── Date maths, all local-time ──────────────────────────────────────────

   Every helper below builds dates from year/month/day NUMBERS and formats
   them by hand. Nothing here goes near new Date("2024-03-15").

   That string form is parsed as UTC midnight, so anywhere west of Greenwich
   it reads back as the 14th. A date of birth that shifts by a day depending
   on the viewer's timezone is the kind of bug that surfaces months later in
   a report nobody can reconcile.                                          */

const MONTHS = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];
const WEEKDAYS = ["Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"];

const pad = (n: number) => String(n).padStart(2, "0");

interface Parts {
  year: number;
  month: number; // 0-indexed, matching Date
  day: number;
  hour: number;
  minute: number;
}

/** "2024-03-15" or "2024-03-15T14:30" → parts. Null when unparseable. */
function parseValue(value: string): Parts | null {
  if (!value) return null;
  const [datePart, timePart = ""] = value.split("T");
  const [y, m, d] = datePart.split("-").map(Number);
  if (!y || !m || !d) return null;
  const [hh = 0, mm = 0] = timePart.split(":").map(Number);
  return { year: y, month: m - 1, day: d, hour: hh || 0, minute: mm || 0 };
}

/** Back to the exact string shape the API and the old native input used. */
function formatValue(parts: Parts, withTime: boolean): string {
  const date = `${parts.year}-${pad(parts.month + 1)}-${pad(parts.day)}`;
  return withTime ? `${date}T${pad(parts.hour)}:${pad(parts.minute)}` : date;
}

function todayParts(): Parts {
  const now = new Date();
  return {
    year: now.getFullYear(),
    month: now.getMonth(),
    day: now.getDate(),
    hour: now.getHours(),
    minute: now.getMinutes(),
  };
}

const daysInMonth = (year: number, month: number) => new Date(year, month + 1, 0).getDate();
const firstWeekday = (year: number, month: number) => new Date(year, month, 1).getDay();

/** Sortable yyyymmdd, so a range check never allocates a Date. */
const ord = (y: number, m: number, d: number) => y * 10000 + (m + 1) * 100 + d;

function ordOf(iso: string | undefined): number | null {
  if (!iso) return null;
  const p = parseValue(iso);
  return p ? ord(p.year, p.month, p.day) : null;
}

export function DateField({ field, value, onChange, error }: DateFieldProps) {
  const fieldId = useId();
  const withTime = field.type === "datetime";
  const selected = parseValue(value);
  const today = todayParts();

  const [open, setOpen] = useState(false);
  // The month on screen, which is NOT the selection — you browse away from the
  // selected date all the time without picking anything.
  const [view, setView] = useState(() => ({
    year: selected?.year ?? today.year,
    month: selected?.month ?? today.month,
  }));
  const [pos, setPos] = useState({ top: 0, left: 0, width: 0 });
  const triggerRef = useRef<HTMLButtonElement>(null);
  const panelRef = useRef<HTMLDivElement>(null);

  const minOrd = ordOf(field.minDate);
  const maxOrd = ordOf(field.maxDate);

  // Year list. Wide enough for a date of birth by default — reaching 1985
  // through a native date input means holding an arrow key, which is the whole
  // reason this component exists. minDate/maxDate narrow it when the field
  // knows better.
  // biome-ignore lint/correctness/useExhaustiveDependencies: only the selected year changes the list, not the rest of the selection.
  const years = useMemo(() => {
    const min = field.minDate ? parseValue(field.minDate)!.year : today.year - 100;
    const max = field.maxDate ? parseValue(field.maxDate)!.year : today.year + 10;
    const list: number[] = [];
    for (let y = max; y >= min; y--) list.push(y);
    // A stored value outside the configured window still has to appear, or
    // opening the picker on an old record silently changes its year.
    if (selected && !list.includes(selected.year)) {
      list.push(selected.year);
      list.sort((a, b) => b - a);
    }
    return list;
  }, [field.minDate, field.maxDate, today.year, selected?.year]);

  const place = useCallback(() => {
    const el = triggerRef.current;
    if (!el) return;
    const rect = el.getBoundingClientRect();
    const panelHeight = withTime ? 400 : 350;
    // Flip above when there is no room below. Forms live inside scrolling
    // modals, and a panel that opens off the bottom edge is the same as no
    // panel at all.
    const below = window.innerHeight - rect.bottom;
    const top =
      below < panelHeight && rect.top > panelHeight ? rect.top - panelHeight - 4 : rect.bottom + 4;
    // Fixed width rather than the field's. A seven-column grid stretched across
    // a full-width form field puts an inch of dead space between each number,
    // which reads as a broken layout rather than a calendar.
    const width = 320;
    const left = Math.max(8, Math.min(rect.left, window.innerWidth - width - 8));
    setPos({ top, left, width });
  }, [withTime]);

  useEffect(() => {
    if (!open) return;
    place();
    function onPointerDown(e: PointerEvent) {
      const t = e.target as Node;
      if (
        triggerRef.current && !triggerRef.current.contains(t) &&
        panelRef.current && !panelRef.current.contains(t)
      ) {
        setOpen(false);
      }
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key !== "Escape") return;
      // Capture + stop, so closing the picker does not also close the modal
      // the form is sitting in and throw away everything typed so far.
      e.stopPropagation();
      setOpen(false);
      triggerRef.current?.focus();
    }
    function onScroll() { place(); }
    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown, true);
    window.addEventListener("scroll", onScroll, true);
    window.addEventListener("resize", onScroll);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown, true);
      window.removeEventListener("scroll", onScroll, true);
      window.removeEventListener("resize", onScroll);
    };
  }, [open, place]);

  // Re-centre on the selection each time it opens, so browsing to 1990 and
  // closing without picking does not strand the next visit there.
  // biome-ignore lint/correctness/useExhaustiveDependencies: runs when the picker opens, not as the selection changes.
  useEffect(() => {
    if (!open) return;
    setView({ year: selected?.year ?? today.year, month: selected?.month ?? today.month });
  }, [open]);

  function commit(day: number, from = view) {
    const next: Parts = {
      year: from.year,
      month: from.month,
      day,
      // Keep the time that was already stored: picking a day should not
      // silently reset an appointment to midnight.
      hour: selected?.hour ?? 0,
      minute: selected?.minute ?? 0,
    };
    onChange(formatValue(next, withTime));
    // A datetime field still needs the time row, so it stays open.
    if (!withTime) setOpen(false);
  }

  function setTime(hour: number, minute: number) {
    const base = selected ?? { ...today, hour: 0, minute: 0 };
    onChange(formatValue({ ...base, hour, minute }, true));
  }

  const shiftMonth = (delta: number) =>
    setView((v) => {
      const m = v.month + delta;
      return { year: v.year + Math.floor(m / 12), month: ((m % 12) + 12) % 12 };
    });

  const isDisabled = (day: number) => {
    const o = ord(view.year, view.month, day);
    if (minOrd !== null && o < minOrd) return true;
    if (maxOrd !== null && o > maxOrd) return true;
    return false;
  };

  const label = selected
    ? `${MONTHS[selected.month]} ${selected.day}, ${selected.year}` +
      (withTime ? ` · ${pad(selected.hour)}:${pad(selected.minute)}` : "")
    : "";

  const leading = firstWeekday(view.year, view.month);
  const total = daysInMonth(view.year, view.month);

  const panel = open ? createPortal(
    <div
      ref={panelRef}
      role="dialog"
      aria-label={`Choose ${field.label}`}
      className="fixed z-[9999] rounded-xl border border-border bg-bg-elevated p-3 shadow-2xl"
      style={{ top: pos.top, left: pos.left, width: pos.width, backgroundColor: "var(--bg-elevated, #22222e)" }}
    >
      {/* Month and year are dropdowns, not just arrows. Clicking a chevron 480
          times to reach a birth year is the problem this replaces. */}
      <div className="flex items-center gap-1.5">
        <button
          type="button"
          onClick={() => shiftMonth(-1)}
          aria-label="Previous month"
          className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-secondary hover:bg-bg-hover hover:text-foreground"
        >
          <ChevronLeft className="h-4 w-4" />
        </button>

        <select
          value={view.month}
          onChange={(e) => setView((v) => ({ ...v, month: Number(e.target.value) }))}
          aria-label="Month"
          className="min-w-0 flex-1 rounded-lg border border-border bg-bg-secondary px-2 py-1.5 text-sm text-foreground outline-none focus:border-accent"
          style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
        >
          {MONTHS.map((m, i) => (
            <option key={m} value={i}>{m}</option>
          ))}
        </select>

        <select
          value={view.year}
          onChange={(e) => setView((v) => ({ ...v, year: Number(e.target.value) }))}
          aria-label="Year"
          className="w-[5.5rem] shrink-0 rounded-lg border border-border bg-bg-secondary px-2 py-1.5 text-sm text-foreground outline-none focus:border-accent"
          style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
        >
          {years.map((y) => (
            <option key={y} value={y}>{y}</option>
          ))}
        </select>

        <button
          type="button"
          onClick={() => shiftMonth(1)}
          aria-label="Next month"
          className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-text-secondary hover:bg-bg-hover hover:text-foreground"
        >
          <ChevronRight className="h-4 w-4" />
        </button>
      </div>

      <div className="mt-3 grid grid-cols-7 gap-0.5">
        {WEEKDAYS.map((w) => (
          <div key={w} className="py-1 text-center text-[11px] font-medium text-text-muted">
            {w}
          </div>
        ))}

        {Array.from({ length: leading }).map((_, i) => (
          <div key={"pad" + i} />
        ))}

        {Array.from({ length: total }).map((_, i) => {
          const day = i + 1;
          const isSelected =
            !!selected && selected.year === view.year && selected.month === view.month && selected.day === day;
          const isToday =
            today.year === view.year && today.month === view.month && today.day === day;
          const off = isDisabled(day);
          return (
            <button
              key={day}
              type="button"
              disabled={off}
              onClick={() => commit(day)}
              aria-current={isSelected ? "date" : undefined}
              className={
                "flex h-9 items-center justify-center rounded-lg text-sm transition-colors " +
                (isSelected
                  ? "bg-accent font-semibold text-white"
                  : off
                    ? "cursor-not-allowed text-text-muted opacity-40"
                    : isToday
                      ? "text-foreground ring-1 ring-accent/60 hover:bg-bg-hover"
                      : "text-foreground hover:bg-bg-hover")
              }
            >
              {day}
            </button>
          );
        })}
      </div>

      {withTime && (
        <div className="mt-3 flex items-center gap-2 border-t border-border pt-3">
          <label className="text-xs text-text-secondary">Time</label>
          <input
          id={fieldId}
            type="time"
            value={selected ? `${pad(selected.hour)}:${pad(selected.minute)}` : ""}
            onChange={(e) => {
              const [h, m] = e.target.value.split(":").map(Number);
              if (!Number.isNaN(h) && !Number.isNaN(m)) setTime(h, m);
            }}
            className={inputClasses({ inputSize: "sm", className: "flex-1" })}
            style={{ backgroundColor: "var(--bg-secondary, #111118)" }}
          />
        </div>
      )}

      <div className="mt-3 flex items-center justify-between border-t border-border pt-3">
        <button
          type="button"
          onClick={() => {
            const t = todayParts();
            setView({ year: t.year, month: t.month });
            commit(t.day, { year: t.year, month: t.month });
          }}
          className="rounded-lg px-2.5 py-1.5 text-xs font-medium text-accent hover:bg-bg-hover"
        >
          Today
        </button>
        <button
          type="button"
          onClick={() => { onChange(""); setOpen(false); }}
          className="rounded-lg px-2.5 py-1.5 text-xs font-medium text-text-secondary hover:bg-bg-hover hover:text-foreground"
        >
          Clear
        </button>
      </div>
    </div>,
    document.body
  ) : null;

  return (
    <div className="space-y-1.5">
      <label htmlFor={fieldId} className="block text-sm font-medium text-foreground">
        {field.label}
        {field.required && <span className="text-danger ml-1">*</span>}
      </label>

      <button
        ref={triggerRef}
        type="button"
        onClick={() => { if (!open) place(); setOpen(!open); }}
        aria-haspopup="dialog"
        aria-expanded={open}
        className={
          "flex w-full items-center justify-between rounded-lg border bg-bg-tertiary px-4 py-2.5 text-sm transition-colors focus:outline-none focus:ring-1 focus:ring-accent " +
          (error ? "border-danger " : "border-border ") +
          (open ? "border-accent " : "") +
          (selected ? "text-foreground" : "text-text-secondary")
        }
      >
        <span>{selected ? label : field.placeholder || "Select a date..."}</span>
        <Calendar className="h-4 w-4 shrink-0 opacity-60" />
      </button>

      {panel}

      {field.description && !error && (
        <p className="text-xs text-text-muted">{field.description}</p>
      )}
      {error && <p className="text-xs text-danger">{error}</p>}
    </div>
  );
}
