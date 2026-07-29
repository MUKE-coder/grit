"use client";

// waitlist-form-01 — Waitlist signup with position counter (common pre-launch pattern)
import { useState } from "react";
export function WaitlistForm({
  title = "Join the waitlist",
  description = "Be the first to know when we launch. We'll send you an invite.",
  currentCount = 2_847,
  onJoin,
}: {
  title?: string; description?: string; currentCount?: number;
  onJoin?: (email: string) => Promise<{ position: number }>;
}) {
  const [email, setEmail] = useState("");
  const [state, setState] = useState<"idle" | "loading" | "done">("idle");
  const [position, setPosition] = useState(0);
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || state !== "idle") return;
    setState("loading");
    try {
      const result = await onJoin?.(email);
      setPosition(result?.position ?? currentCount + 1);
      setState("done");
    } catch {
      setState("idle");
    }
  };
  return (
    <div className="bg-bg-elevated border border-border rounded-2xl p-6 flex flex-col gap-5 w-full max-w-md">
      {state === "done" ? (
        <div className="flex flex-col items-center gap-3 py-4 text-center">
          <div className="w-14 h-14 rounded-full bg-success/20 flex items-center justify-center">
            <svg className="w-7 h-7 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h3 className="text-base font-semibold text-foreground">You're #{position} on the list!</h3>
          <p className="text-sm text-text-secondary">We'll email you at <strong>{email}</strong> when your spot is ready.</p>
        </div>
      ) : (
        <>
          <div>
            <h3 className="text-base font-semibold text-foreground">{title}</h3>
            <p className="text-sm text-text-secondary mt-1">{description}</p>
          </div>
          <div className="flex items-center gap-2 text-xs text-text-muted bg-background border border-border rounded-lg px-3 py-2">
            <div className="flex -space-x-1">
              {["SC", "JW", "PM"].map((i) => (
                <div key={i} className="w-5 h-5 rounded-full bg-accent/30 border border-border flex items-center justify-center">
                  <span className="text-xs font-bold text-accent" style={{ fontSize: 8 }}>{i[0]}</span>
                </div>
              ))}
            </div>
            <span><strong className="text-foreground">{currentCount.toLocaleString()}</strong> people already waiting</span>
          </div>
          <form onSubmit={submit} className="flex gap-2">
            <input
              type="email" value={email} onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com" required
              className="flex-1 bg-background border border-border rounded-lg px-3 py-2.5 text-sm text-foreground placeholder-text-muted outline-none focus:border-accent transition-colors"
            />
            <button
              type="submit" disabled={state === "loading"}
              className="px-4 py-2.5 rounded-lg bg-accent hover:bg-accent/90 text-white text-sm font-semibold transition-colors disabled:opacity-60 flex-shrink-0"
            >
              {state === "loading" ? "..." : "Join"}
            </button>
          </form>
          <p className="text-xs text-text-muted text-center">No spam. Unsubscribe anytime.</p>
        </>
      )}
    </div>
  );
}

export default WaitlistForm;
