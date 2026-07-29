"use client";

// two-factor-setup-01 — 2FA setup card with QR + backup codes (GitHub / Clerk-inspired)
import { useState } from "react";
export function TwoFactorSetup({
  qrCodeUrl = "",
  secret = "JBSWY3DPEHPK3PXP",
  backupCodes = ["a1b2-c3d4", "e5f6-g7h8", "i9j0-k1l2", "m3n4-o5p6", "q7r8-s9t0", "u1v2-w3x4"],
  onComplete,
}: {
  qrCodeUrl?: string; secret?: string; backupCodes?: string[];
  onComplete?: (code: string) => void;
}) {
  const [step, setStep] = useState<"qr" | "verify" | "backup">("qr");
  const [code, setCode] = useState("");
  return (
    <div className="bg-bg-elevated border border-border rounded-2xl p-6 w-full max-w-md flex flex-col gap-5">
      <div className="flex items-center gap-3">
        {["qr", "verify", "backup"].map((s, i) => (
          <div key={s} className="flex items-center gap-2">
            <div className={"w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold " + (s === step ? "bg-accent text-white" : ["qr","verify","backup"].indexOf(step) > i ? "bg-success text-white" : "bg-bg-hover text-text-muted border border-border")}>{i + 1}</div>
            {i < 2 && <div className={"h-0.5 w-8 " + (["qr","verify","backup"].indexOf(step) > i ? "bg-success" : "bg-border")} />}
          </div>
        ))}
      </div>
      {step === "qr" && (
        <div className="flex flex-col gap-4">
          <h3 className="text-base font-semibold text-foreground">Scan QR Code</h3>
          <p className="text-sm text-text-secondary">Scan this code with your authenticator app (Google Authenticator, Authy, etc.)</p>
          <div className="w-40 h-40 bg-white rounded-xl flex items-center justify-center self-center">
            {qrCodeUrl ? <img src={qrCodeUrl} alt="QR Code" className="w-full h-full" /> : <span className="text-xs text-gray-400">QR Code</span>}
          </div>
          <div className="flex flex-col gap-1">
            <p className="text-xs text-text-muted">Can't scan? Enter this code manually:</p>
            <code className="text-sm font-mono text-foreground bg-background border border-border rounded-lg px-3 py-2 select-all">{secret}</code>
          </div>
          <button onClick={() => setStep("verify")} className="w-full py-2.5 rounded-lg bg-accent hover:bg-accent/90 text-white font-semibold transition-colors">Continue</button>
        </div>
      )}
      {step === "verify" && (
        <div className="flex flex-col gap-4">
          <h3 className="text-base font-semibold text-foreground">Verify Setup</h3>
          <p className="text-sm text-text-secondary">Enter the 6-digit code from your authenticator app to confirm setup.</p>
          <input
            value={code} onChange={(e) => setCode(e.target.value.slice(0, 6))}
            placeholder="000000" maxLength={6}
            className="text-center text-2xl font-mono tracking-widest bg-background border border-border rounded-xl px-4 py-3 text-foreground outline-none focus:border-accent transition-colors"
          />
          <button onClick={() => { if (code.length === 6) { onComplete?.(code); setStep("backup"); } }} disabled={code.length < 6} className="w-full py-2.5 rounded-lg bg-accent hover:bg-accent/90 text-white font-semibold transition-colors disabled:opacity-50">Verify</button>
        </div>
      )}
      {step === "backup" && (
        <div className="flex flex-col gap-4">
          <h3 className="text-base font-semibold text-foreground">Save backup codes</h3>
          <p className="text-sm text-text-secondary">Store these codes somewhere safe. Each can be used once if you lose access to your device.</p>
          <div className="grid grid-cols-2 gap-2 bg-background border border-border rounded-xl p-4">
            {backupCodes.map((c) => <code key={c} className="text-xs font-mono text-text-secondary">{c}</code>)}
          </div>
          <button onClick={() => navigator.clipboard.writeText(backupCodes.join("\n")).catch(() => {})} className="w-full py-2 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover text-sm transition-colors">Copy all codes</button>
        </div>
      )}
    </div>
  );
}

export default TwoFactorSetup;
