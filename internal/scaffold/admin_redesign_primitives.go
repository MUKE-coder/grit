package scaffold

// This file holds the primitives the redesigned pages are assembled from.
//
// They exist so that a status in the users table and a status on the health
// page are the same object rather than two things that happen to look similar
// this week. Every one of them had three or four hand-rolled versions across
// the panel before this, differing in padding, in radius, and in whether the
// state was a word or a colour.

// adminPillPrimitive is components/ui/pill.tsx.
//
// A status, as a word, with a dot. The word is the part that matters: colour
// alone fails for the eight percent of men who cannot separate the red one
// from the green one, and it fails for everybody in a screenshot pasted into
// a ticket.
func adminPillPrimitive() string {
	return `import type { ReactNode } from "react";

export type PillTone = "good" | "warn" | "bad" | "accent" | "neutral" | "info";

const TONES: Record<PillTone, { fill: string; text: string; dot: string }> = {
  good: { fill: "bg-success-soft", text: "text-success", dot: "bg-success" },
  warn: { fill: "bg-warning-soft", text: "text-warning", dot: "bg-warning" },
  bad: { fill: "bg-danger-soft", text: "text-danger", dot: "bg-danger" },
  accent: { fill: "bg-accent-soft", text: "text-accent", dot: "bg-accent" },
  info: { fill: "bg-info-soft", text: "text-info", dot: "bg-info" },
  neutral: { fill: "bg-bg-hover", text: "text-text-secondary", dot: "bg-text-muted" },
};

export interface PillProps {
  children: ReactNode;
  tone?: PillTone;
  /** A dot before the label. Use it for a state, not for a count. */
  dot?: boolean;
  className?: string;
  title?: string;
}

/**
 * A status, as a word.
 *
 * Active, Inactive, Published, Draft, Healthy, Off, Admin, Editor. 22px tall
 * so it sits inside a table row without changing the row's height.
 */
export function Pill({ children, tone = "neutral", dot = false, className = "", title }: PillProps) {
  const t = TONES[tone];
  return (
    <span
      title={title}
      className={
        "inline-flex h-[22px] shrink-0 items-center gap-1.5 rounded-md px-2 text-xs font-medium " +
        t.fill +
        " " +
        t.text +
        " " +
        className
      }
    >
      {dot && <span className={"h-1.5 w-1.5 shrink-0 rounded-full " + t.dot} aria-hidden="true" />}
      {children}
    </span>
  );
}
`
}

// adminEmptyStatePrimitive is components/ui/empty-state.tsx.
//
// "No data" is the most common sentence in an admin panel and it is the least
// useful one. It tells a reader that the request succeeded and nothing else:
// not whether this is normal, not whether they did something wrong, not what
// would put something here.
//
// So the component requires a title and a description, and offers an action
// and a mono hint. A table that is empty because nobody has made one yet reads
// differently from a table that is empty because a filter excluded everything,
// and the component's job is to make the page say which.
func adminEmptyStatePrimitive() string {
	return `import type { ReactNode } from "react";

export interface EmptyStateProps {
  /** A Lucide icon element. Decorative: the title carries the meaning. */
  icon?: ReactNode;
  /** What is not here. "No API keys yet", not "No data". */
  title: string;
  /** One sentence: what will appear here, and what puts it there. */
  description: string;
  /** The one thing to do about it. */
  action?: ReactNode;
  /** A command or a snippet, in mono, when the answer is something to run. */
  hint?: string;
  className?: string;
}

export function EmptyState({
  icon,
  title,
  description,
  action,
  hint,
  className = "",
}: EmptyStateProps) {
  return (
    <div className={"flex flex-col items-center px-6 py-10 text-center " + className}>
      {icon && (
        <div
          className="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-bg-hover text-text-muted"
          aria-hidden="true"
        >
          {icon}
        </div>
      )}
      <p className="text-[15px] font-semibold text-foreground">{title}</p>
      <p className="mt-1.5 max-w-[46ch] text-sm text-text-secondary">{description}</p>
      {action && <div className="mt-4">{action}</div>}
      {hint && (
        <code className="mt-4 rounded-md border border-border bg-bg-tertiary px-3 py-1.5 font-mono text-[12.5px] text-text-secondary">
          {hint}
        </code>
      )}
    </div>
  );
}
`
}

// adminSegmentedPrimitive is components/ui/segmented.tsx.
//
// A one-of choice with few enough options to show them all: a time range, a
// status filter, a layout. A real radio group, because it is one, and because
// a group of buttons with aria-pressed tells a screen reader there are four
// independent toggles here rather than one choice with four answers.
func adminSegmentedPrimitive() string {
	return `"use client";

import { useId } from "react";

export interface SegmentedOption<T extends string> {
  value: T;
  label: string;
  /** A count beside the label, in mono. */
  count?: number | string;
}

export interface SegmentedProps<T extends string> {
  options: ReadonlyArray<SegmentedOption<T>>;
  value: T;
  onChange: (value: T) => void;
  /** Names the choice for a screen reader. "Time range", "Status". */
  label: string;
  className?: string;
}

/**
 * The selected option is a raised chip on a recessed rail, which is the one
 * arrangement that reads as "this one of these" rather than as four buttons.
 */
export function Segmented<T extends string>({
  options,
  value,
  onChange,
  label,
  className = "",
}: SegmentedProps<T>) {
  const name = useId();
  return (
    <div
      role="radiogroup"
      aria-label={label}
      className={
        "inline-flex h-9 shrink-0 items-center gap-0.5 rounded-lg bg-bg-hover p-0.5 " + className
      }
    >
      {options.map((option) => {
        const selected = option.value === value;
        return (
          <label
            key={option.value}
            className={
              "inline-flex h-8 cursor-pointer items-center gap-1.5 rounded-md px-3 text-[13px] transition-grit " +
              (selected
                ? "bg-bg-elevated font-medium text-foreground shadow-sm"
                : "text-text-secondary hover:text-foreground")
            }
          >
            <input
              type="radio"
              name={name}
              value={option.value}
              checked={selected}
              onChange={() => onChange(option.value)}
              className="sr-only"
            />
            {option.label}
            {option.count !== undefined && (
              <span className="font-mono text-[11px] text-text-muted">{option.count}</span>
            )}
          </label>
        );
      })}
    </div>
  );
}
`
}

// adminUnderlineTabsPrimitive is components/ui/tabs-underline.tsx.
//
// Sections of one page, with a count each. The count is the point: "Flagged 0"
// answers the question the tab exists for without anybody clicking it, and a
// tab you do not have to open is the cheapest answer an interface can give.
//
// The selection lives in the URL. A tab that resets on refresh cannot be
// linked to, cannot be bookmarked, and loses its place every time somebody
// follows a link out and comes back.
func adminUnderlineTabsPrimitive() string {
	return `"use client";

export interface UnderlineTab<T extends string> {
  value: T;
  label: string;
  count?: number | string;
}

export interface UnderlineTabsProps<T extends string> {
  tabs: ReadonlyArray<UnderlineTab<T>>;
  value: T;
  onChange: (value: T) => void;
  label: string;
  className?: string;
}

export function UnderlineTabs<T extends string>({
  tabs,
  value,
  onChange,
  label,
  className = "",
}: UnderlineTabsProps<T>) {
  return (
    <div
      role="tablist"
      aria-label={label}
      className={"flex items-center gap-1 overflow-x-auto border-b border-border " + className}
      onKeyDown={(event) => {
        // Arrow keys move between tabs, which is what the tablist pattern
        // promises and what a keyboard user will try first.
        const index = tabs.findIndex((t) => t.value === value);
        if (event.key === "ArrowRight") {
          event.preventDefault();
          onChange(tabs[(index + 1) % tabs.length].value);
        } else if (event.key === "ArrowLeft") {
          event.preventDefault();
          onChange(tabs[(index - 1 + tabs.length) % tabs.length].value);
        }
      }}
    >
      {tabs.map((tab) => {
        const selected = tab.value === value;
        return (
          <button
            key={tab.value}
            type="button"
            role="tab"
            aria-selected={selected}
            tabIndex={selected ? 0 : -1}
            onClick={() => onChange(tab.value)}
            className={
              "-mb-px inline-flex shrink-0 items-center gap-2 border-b-2 px-3 py-2.5 text-sm transition-grit " +
              (selected
                ? "border-foreground font-medium text-foreground"
                : "border-transparent text-text-secondary hover:text-foreground")
            }
          >
            {tab.label}
            {tab.count !== undefined && (
              <span className="font-mono text-[11px] text-text-muted">{tab.count}</span>
            )}
          </button>
        );
      })}
    </div>
  );
}
`
}

// adminRadioCardPrimitive is components/ui/radio-card.tsx.
//
// For a one-of choice where the options need explaining: which kind of API
// key, how often to back up. A select hides the explanation behind a click,
// and these are exactly the choices somebody should not make without reading
// one.
func adminRadioCardPrimitive() string {
	return `"use client";

import type { ReactNode } from "react";

export interface RadioCardProps {
  name: string;
  value: string;
  checked: boolean;
  onChange: (value: string) => void;
  title: string;
  /** A machine-ish detail under the title: a key prefix, a cron expression. */
  meta?: string;
  description: string;
  icon?: ReactNode;
  disabled?: boolean;
}

export function RadioCard({
  name,
  value,
  checked,
  onChange,
  title,
  meta,
  description,
  icon,
  disabled = false,
}: RadioCardProps) {
  return (
    <label
      className={
        "flex cursor-pointer gap-3 rounded-xl border p-4 transition-grit " +
        (disabled ? "cursor-not-allowed opacity-60 " : "") +
        (checked
          ? "border-accent-line bg-accent-soft"
          : "border-border bg-bg-elevated hover:border-border-strong")
      }
    >
      <input
        type="radio"
        name={name}
        value={value}
        checked={checked}
        disabled={disabled}
        onChange={() => onChange(value)}
        className="mt-0.5 h-4 w-4 shrink-0 accent-[var(--accent)]"
      />
      <span className="min-w-0 flex-1">
        <span className="flex items-center gap-2">
          {icon && <span className="shrink-0 text-text-muted">{icon}</span>}
          <span className="text-sm font-medium text-foreground">{title}</span>
          {meta && <code className="font-mono text-[11.5px] text-text-muted">{meta}</code>}
        </span>
        <span className="mt-1 block text-[13px] text-text-secondary">{description}</span>
      </span>
    </label>
  );
}
`
}

// adminHubTilePrimitive is components/ui/hub-tile.tsx.
//
// One destination on the system hub. The status pill is the reason this is a
// component rather than a link: a tile that says "Off, needs Redis" answers
// the question before the click, and a grid of tiles that all say the same
// static sentence is a menu pretending to be a dashboard.
func adminHubTilePrimitive() string {
	return `"use client";

import Link from "next/link";
import type { ReactNode } from "react";
import { ArrowRight } from "@/lib/icons";
import { Pill, type PillTone } from "./pill";

export interface HubTileProps {
  href: string;
  icon: ReactNode;
  title: string;
  description: string;
  /** Live, from the backend. "6 healthy, 2 off", "Off, needs Redis". */
  status?: string;
  statusTone?: PillTone;
}

export function HubTile({
  href,
  icon,
  title,
  description,
  status,
  statusTone = "neutral",
}: HubTileProps) {
  return (
    <Link
      href={href}
      className="group flex flex-col rounded-xl border border-border bg-bg-elevated p-4 transition-grit hover:border-border-strong"
    >
      <div className="flex items-start justify-between gap-3">
        <span
          className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-bg-hover text-text-secondary"
          aria-hidden="true"
        >
          {icon}
        </span>
        {status && (
          <Pill tone={statusTone} dot>
            {status}
          </Pill>
        )}
      </div>
      <p className="mt-3 text-[15px] font-semibold text-foreground">{title}</p>
      <p className="mt-1 flex-1 text-[13px] text-text-secondary">{description}</p>
      <span className="mt-3 inline-flex items-center gap-1 text-[13px] font-medium text-accent">
        Open
        <ArrowRight className="h-3.5 w-3.5 transition-grit group-hover:translate-x-0.5" />
      </span>
    </Link>
  );
}
`
}

// adminNotePrimitive is components/ui/note.tsx.
//
// The footnote under a section: what a header format is, what the difference
// between delete and erase is, which of two similar buttons is the
// irreversible one. Small, quiet, and in the place somebody looks after
// reading the thing it explains rather than in a tooltip they have to find.
func adminNotePrimitive() string {
	return `import type { ReactNode } from "react";
import { Info } from "@/lib/icons";

export function Note({
  children,
  className = "",
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <p className={"flex items-start gap-2 text-[13px] text-text-secondary " + className}>
      <Info className="mt-0.5 h-3.5 w-3.5 shrink-0 text-text-muted" aria-hidden="true" />
      <span>{children}</span>
    </p>
  );
}
`
}
