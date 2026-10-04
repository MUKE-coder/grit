"use client";

import { useEffect, useState } from "react";
import { Calendar } from "@/lib/icons";

/**
 * The date and time, for the dashboard header.
 *
 * Client-only on purpose. Formatting a date on the server and again in the
 * browser disagrees on two counts, the clock and the locale, and React's
 * response to a hydration mismatch is to discard the subtree and rebuild it.
 * That is a steep price for a decoration, so this renders nothing until it is
 * mounted and there is nothing for the server to disagree with.
 *
 * The tick is every second but the state is the formatted string, so React
 * bails out of the re-render on the 59 ticks a minute that produce the same
 * text. A 60-second interval would have been cheaper still and would show a
 * minute that is up to a minute stale, which on a clock is the one thing that
 * cannot be wrong.
 */
export function LiveClock() {
  const [now, setNow] = useState<{ date: string; time: string } | null>(null);

  useEffect(() => {
    const read = () => {
      const d = new Date();
      setNow((prev) => {
        const next = {
          date: d.toLocaleDateString(undefined, {
            weekday: "long",
            day: "numeric",
            month: "long",
            year: "numeric",
          }),
          time: d.toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit" }),
        };
        // Same text, same object: React skips the render.
        return prev && prev.date === next.date && prev.time === next.time ? prev : next;
      });
    };
    read();
    const id = setInterval(read, 1000);
    return () => clearInterval(id);
  }, []);

  if (!now) {
    // Reserves the row so the header does not jump when the clock arrives.
    return <div className="h-9" aria-hidden="true" />;
  }

  return (
    <div className="flex items-center gap-2.5 rounded-lg border border-border bg-bg-elevated px-3 py-1.5">
      <Calendar className="h-4 w-4 shrink-0 text-accent" aria-hidden="true" />
      <div className="leading-tight">
        <p className="text-sm font-medium text-foreground">{now.date}</p>
        <p className="font-mono text-xs text-text-muted">{now.time}</p>
      </div>
    </div>
  );
}
