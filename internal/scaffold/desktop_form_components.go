package scaffold

// Shared desktop form/dialog components used by generated resource pages:
// a promise-based confirm dialog (replaces the native window.confirm, which in
// Wails shows an ugly "wails.localhost says" box), a comma-formatting number
// input, and a searchable select (combobox) for relationship pickers.

// desktopClientConfirmDialog is a styled, promise-based confirm modal + provider.
// Usage: const confirm = useConfirm(); if (await confirm({ message, danger: true })) { ... }
func desktopClientConfirmDialog() string {
	return `import { createContext, useCallback, useContext, useEffect, useState } from "react";
import { AlertTriangle } from "lucide-react";

interface ConfirmOptions {
  title?: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  danger?: boolean;
}
type Resolver = (v: boolean) => void;

const ConfirmContext = createContext<(o: ConfirmOptions) => Promise<boolean>>(
  async () => false,
);

export function ConfirmProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<{ opts: ConfirmOptions; resolve: Resolver } | null>(null);

  const confirm = useCallback(
    (opts: ConfirmOptions) => new Promise<boolean>((resolve) => setState({ opts, resolve })),
    [],
  );

  const close = (v: boolean) => {
    state?.resolve(v);
    setState(null);
  };

  useEffect(() => {
    if (!state) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") close(false);
      if (e.key === "Enter") close(true);
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [state]);

  const o = state?.opts;

  return (
    <ConfirmContext.Provider value={confirm}>
      {children}
      {state && o && (
        <div className="fixed inset-0 z-[70] flex items-center justify-center p-4" onClick={() => close(false)}>
          <div className="fixed inset-0 bg-black/50 backdrop-blur-sm" />
          <div onClick={(e) => e.stopPropagation()} className="relative z-10 w-full max-w-md rounded-2xl border border-border bg-surface p-6 shadow-2xl">
            <div className="flex items-start gap-3">
              <span className={"inline-flex h-10 w-10 shrink-0 items-center justify-center rounded-full " + (o.danger ? "bg-danger/10 text-danger" : "bg-accent/10 text-accent")}>
                <AlertTriangle className="h-5 w-5" />
              </span>
              <div className="min-w-0 flex-1">
                <h2 className="text-[16px] font-semibold text-foreground">{o.title ?? "Are you sure?"}</h2>
                <p className="mt-1 text-[13px] text-foreground-secondary">{o.message}</p>
              </div>
            </div>
            <div className="mt-6 flex justify-end gap-3">
              <button onClick={() => close(false)} className="rounded-lg border border-border px-4 py-2 text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover">
                {o.cancelLabel ?? "Cancel"}
              </button>
              <button
                onClick={() => close(true)}
                className={"rounded-lg px-4 py-2 text-[13px] font-semibold text-white " + (o.danger ? "bg-danger hover:bg-danger/90" : "bg-accent hover:bg-accent-hover")}
              >
                {o.confirmLabel ?? "Confirm"}
              </button>
            </div>
          </div>
        </div>
      )}
    </ConfirmContext.Provider>
  );
}

export function useConfirm() {
  return useContext(ConfirmContext);
}
`
}

// desktopClientNumberInput is a text input that live-formats with thousands
// separators (1000 -> 1,000) while exposing a clean numeric value. numberKind
// controls decimals/negatives: uint (none), int (negatives, no decimals),
// float (both).
func desktopClientNumberInput() string {
	return `import { useEffect, useRef, useState } from "react";

interface NumberInputProps {
  value: number;
  onChange: (value: number) => void;
  kind?: "int" | "uint" | "float";
  className?: string;
  placeholder?: string;
}

function format(n: number, kind: string): string {
  if (n === 0) return "0";
  if (!Number.isFinite(n)) return "";
  const opts = kind === "float" ? { maximumFractionDigits: 6 } : { maximumFractionDigits: 0 };
  return n.toLocaleString("en-US", opts);
}

// Keep only characters valid for the domain, then group the integer part with
// commas. Preserves a trailing "." and trailing zeros while typing a decimal.
function reformat(raw: string, kind: string): { display: string; value: number } {
  let s = raw.replace(/,/g, "");
  if (kind !== "float") s = s.replace(/\./g, "");
  if (kind === "uint") s = s.replace(/-/g, "");
  s = s.replace(/[^0-9.\-]/g, "");
  if (s === "" || s === "-" || s === ".") return { display: s, value: 0 };
  const neg = s.startsWith("-");
  if (neg) s = s.slice(1);
  const [intPart, ...rest] = s.split(".");
  const decPart = rest.join("");
  const groupedInt = intPart.replace(/^0+(?=\d)/, "").replace(/\B(?=(\d{3})+(?!\d))/g, ",") || "0";
  let display = (neg ? "-" : "") + groupedInt;
  const hadDot = kind === "float" && raw.replace(/,/g, "").includes(".");
  if (hadDot) display += "." + decPart;
  const value = parseFloat((neg ? "-" : "") + intPart + (decPart ? "." + decPart : "")) || 0;
  return { display, value };
}

export function NumberInput({ value, onChange, kind = "float", className, placeholder }: NumberInputProps) {
  const [display, setDisplay] = useState(() => format(value, kind));
  const focused = useRef(false);

  // Sync external value changes (e.g. edit prefill) unless the user is typing.
  useEffect(() => {
    if (!focused.current) setDisplay(format(value, kind));
  }, [value, kind]);

  return (
    <input
      type="text"
      inputMode={kind === "float" ? "decimal" : "numeric"}
      value={display}
      placeholder={placeholder}
      className={className}
      onFocus={() => { focused.current = true; }}
      onBlur={() => { focused.current = false; setDisplay(format(value, kind)); }}
      onChange={(e) => {
        const { display: d, value: v } = reformat(e.target.value, kind);
        setDisplay(d);
        onChange(v);
      }}
    />
  );
}
`
}

