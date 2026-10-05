import { useState } from "react";
import { Sparkles, Loader2 } from "@/lib/icons";

interface GenerateButtonProps {
  /** Runs the field's generate() and applies the result. May be async. */
  onGenerate: () => void | Promise<void>;
}

// A compact affordance next to a field's label. Clicking it runs the
// developer-defined generate() for that field; the button shows a spinner
// until the (possibly async) generator resolves, and never submits the form.
export function GenerateButton({ onGenerate }: GenerateButtonProps) {
  const [busy, setBusy] = useState(false);

  const run = async () => {
    if (busy) return;
    try {
      setBusy(true);
      await onGenerate();
    } finally {
      setBusy(false);
    }
  };

  return (
    <button
      type="button"
      onClick={run}
      disabled={busy}
      title="Generate"
      className="inline-flex items-center gap-1 rounded-md border border-border bg-bg-tertiary px-2 py-1 text-xs font-medium text-text-secondary transition-colors hover:border-accent/40 hover:text-accent disabled:opacity-50"
    >
      {busy ? <Loader2 className="h-3 w-3 animate-spin" /> : <Sparkles className="h-3 w-3" />}
      Generate
    </button>
  );
}
