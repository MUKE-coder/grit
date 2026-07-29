"use client";

// flash-sale-timer-01 — Flash sale countdown timer (ecommerce urgency pattern)
import { useEffect, useState } from "react";
export function FlashSaleTimer({
  endsAt = new Date(Date.now() + 4 * 60 * 60 * 1000),
  label = "Flash Sale ends in",
}: {
  endsAt?: Date; label?: string;
}) {
  const calc = () => {
    const diff = Math.max(0, Math.floor((endsAt.getTime() - Date.now()) / 1000));
    return { h: Math.floor(diff / 3600), m: Math.floor((diff % 3600) / 60), s: diff % 60 };
  };
  const [time, setTime] = useState(calc());
  useEffect(() => {
    const t = setInterval(() => setTime(calc()), 1000);
    return () => clearInterval(t);
  }, [endsAt]);
  const pad = (n: number) => String(n).padStart(2, "0");
  return (
    <div className="flex items-center gap-3 p-4 bg-danger/10 border border-danger/30 rounded-xl">
      <svg className="w-5 h-5 text-danger flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
      <span className="text-sm font-medium text-foreground">{label}</span>
      <div className="flex items-center gap-1 ml-auto">
        {[{ v: time.h, l: "h" }, { v: time.m, l: "m" }, { v: time.s, l: "s" }].map(({ v, l }, i) => (
          <div key={l} className="flex items-center gap-1">
            <div className="bg-danger text-white font-bold font-mono text-sm px-2.5 py-1 rounded-lg min-w-[2.5rem] text-center">
              {pad(v)}
            </div>
            <span className="text-text-muted text-xs">{l}</span>
            {i < 2 && <span className="text-danger font-bold">:</span>}
          </div>
        ))}
      </div>
    </div>
  );
}

export default FlashSaleTimer;
