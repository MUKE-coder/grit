package scaffold

// Two-factor management UI for the admin profile page.
//
// The API has shipped TOTP for a while — setup, enable, disable, backup codes,
// trusted devices, and the login challenge — but nothing in the panel could
// drive any of it, so the feature was unreachable for anyone not writing curl
// by hand. This is that missing surface.
//
// Two ordering rules matter here and are easy to get wrong:
//
//  1. The login page must understand the challenge BEFORE anyone can enable
//     2FA, or the first person who turns it on locks themselves out. That
//     lives in adminThemedLoginPage.
//  2. Backup codes are shown exactly once, at the moment they are generated —
//     the server stores only hashes. The dialog therefore refuses to close
//     until they have been copied or downloaded, because "I'll grab them
//     later" is not an option the API can honour.

func adminTwoFactorCard() string {
	return `"use client";

import { useState } from "react";
import {
  useTOTPStatus,
  useTrustedDevices,
  useTOTPSetup,
  useEnableTOTP,
  useDisableTOTP,
  useRegenerateBackupCodes,
  useRevokeTrustedDevice,
} from "@/hooks/use-auth";
import { ShieldCheck, Loader2, Copy, Check, Monitor, X } from "@/lib/icons";
import { describeDevice } from "@/components/profile/active-sessions";

/** Recovery codes are unrecoverable once dismissed — see the file header. */
function BackupCodes({ codes, onDone }: { codes: string[]; onDone: () => void }) {
  const [saved, setSaved] = useState(false);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(codes.join("\n"));
    } catch {
      // clipboard.writeText rejects on an insecure origin or when the document
      // is not focused. Mark them saved anyway — refusing would trap the user
      // in a panel with no exit, and Download is still there as the reliable
      // path.
    }
    setSaved(true);
  };

  const download = () => {
    const blob = new Blob([codes.join("\n") + "\n"], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "backup-codes.txt";
    a.click();
    URL.revokeObjectURL(url);
    setSaved(true);
  };

  return (
    <div className="rounded-lg border border-warning/40 bg-warning/[0.06] p-4">
      <p className="text-sm font-medium text-foreground mb-1">Save your backup codes</p>
      <p className="text-xs text-foreground-secondary mb-3">
        Each code works once, if you lose your authenticator. This is the only time they are
        shown — the server keeps only hashes.
      </p>

      <div className="grid grid-cols-2 gap-2 mb-3 font-mono text-sm">
        {codes.map((c) => (
          <div key={c} className="rounded bg-bg-tertiary px-3 py-1.5 text-center tracking-wider">
            {c}
          </div>
        ))}
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <button
          type="button"
          onClick={copy}
          className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-sm hover:bg-bg-hover"
        >
          {saved ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
          Copy
        </button>
        <button
          type="button"
          onClick={download}
          className="rounded-lg border border-border px-3 py-1.5 text-sm hover:bg-bg-hover"
        >
          Download
        </button>
        <button
          type="button"
          onClick={onDone}
          disabled={!saved}
          className="ml-auto rounded-lg bg-accent px-3 py-1.5 text-sm text-white disabled:opacity-50"
          title={saved ? undefined : "Copy or download them first"}
        >
          I have saved them
        </button>
      </div>
    </div>
  );
}

export function TwoFactorCard() {
  const { data: status, isLoading } = useTOTPStatus();
  const { data: devices } = useTrustedDevices(!!status?.enabled);
  const setup = useTOTPSetup();
  const enable = useEnableTOTP();
  const disable = useDisableTOTP();
  const regenerate = useRegenerateBackupCodes();
  const revoke = useRevokeTrustedDevice();

  const [secret, setSecret] = useState<string | null>(null);
  const [qr, setQr] = useState<string | null>(null);
  const [code, setCode] = useState("");
  const [codes, setCodes] = useState<string[] | null>(null);
  const [password, setPassword] = useState("");
  const [disabling, setDisabling] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const errText = (e: unknown) =>
    (e as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error
      ?.message ?? "Something went wrong";

  const startSetup = () => {
    setError(null);
    setup.mutate(undefined, {
      onSuccess: (d) => {
        setSecret(d.secret);
        setQr(d.qr_code);
      },
      onError: (e) => setError(errText(e)),
    });
  };

  const confirmEnable = () => {
    if (!secret) return;
    setError(null);
    enable.mutate(
      { secret, code },
      {
        onSuccess: (d) => {
          setCodes(d.backup_codes);
          setSecret(null);
          setQr(null);
          setCode("");
        },
        onError: (e) => setError(errText(e)),
      }
    );
  };

  const confirmDisable = () => {
    setError(null);
    disable.mutate(password, {
      onSuccess: () => {
        setDisabling(false);
        setPassword("");
      },
      onError: (e) => setError(errText(e)),
    });
  };

  if (isLoading) {
    return (
      <div className="rounded-xl border border-border bg-bg-secondary p-6">
        <Loader2 className="h-4 w-4 animate-spin text-text-muted" />
      </div>
    );
  }

  return (
    <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
      <div className="flex items-start gap-3 border-b border-border px-6 py-4">
        <ShieldCheck className="mt-0.5 h-4 w-4 text-accent" />
        <div className="min-w-0">
          <h2 className="font-semibold text-foreground">Two-factor authentication</h2>
          <p className="text-sm text-foreground-secondary">
            A code from your authenticator app, on top of your password.
          </p>
        </div>
        <span
          className={
            "ml-auto shrink-0 rounded-full px-2.5 py-1 text-xs font-medium " +
            (status?.enabled
              ? "bg-success/15 text-success"
              : "bg-bg-tertiary text-text-muted")
          }
        >
          {status?.enabled ? "On" : "Off"}
        </span>
      </div>

      <div className="space-y-5 p-6">
        {error && (
          <p className="rounded-lg bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
        )}

        {codes && <BackupCodes codes={codes} onDone={() => setCodes(null)} />}

        {/* ── Off, and not mid-setup ── */}
        {!status?.enabled && !secret && !codes && (
          <button
            type="button"
            onClick={startSetup}
            disabled={setup.isPending}
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50"
          >
            {setup.isPending && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
            Turn on two-factor
          </button>
        )}

        {/* ── Mid-setup: scan, then confirm with a live code ── */}
        {secret && qr && (
          <div className="grid gap-6 sm:grid-cols-[auto_1fr]">
            <div className="rounded-lg border border-border bg-white p-3">
              {/* Rendered by the API, so no QR library ships in the client. */}
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src={qr} alt="Two-factor setup QR code" width={176} height={176} />
            </div>

            <div className="min-w-0 space-y-3">
              <p className="text-sm text-foreground-secondary">
                Scan this with Google Authenticator, 1Password, Authy or similar. Cannot scan?
                Enter this key by hand:
              </p>
              <code className="block break-all rounded bg-bg-tertiary px-3 py-2 font-mono text-xs text-foreground">
                {secret}
              </code>

              <div className="flex flex-wrap items-center gap-2">
                <input
                  value={code}
                  onChange={(e) => setCode(e.target.value)}
                  placeholder="000000"
                  inputMode="numeric"
                  maxLength={6}
                  className="w-32 rounded-lg border border-border bg-bg-primary px-3 py-2 text-center font-mono tracking-[0.3em] text-foreground"
                />
                <button
                  type="button"
                  onClick={confirmEnable}
                  disabled={enable.isPending || code.trim().length !== 6}
                  className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50"
                >
                  {enable.isPending && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
                  Verify and turn on
                </button>
                <button
                  type="button"
                  onClick={() => { setSecret(null); setQr(null); setCode(""); }}
                  className="rounded-lg border border-border px-3 py-2 text-sm hover:bg-bg-hover"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        )}

        {/* ── On: codes left, regenerate, trusted devices, disable ── */}
        {status?.enabled && !codes && (
          <>
            <div className="flex flex-wrap items-center gap-3">
              <span className="text-sm text-foreground-secondary">
                {status.backup_codes_remaining} backup code
                {status.backup_codes_remaining === 1 ? "" : "s"} left
              </span>
              <button
                type="button"
                onClick={() =>
                  regenerate.mutate(undefined, {
                    onSuccess: (d) => setCodes(d.backup_codes),
                    onError: (e) => setError(errText(e)),
                  })
                }
                disabled={regenerate.isPending}
                className="rounded-lg border border-border px-3 py-1.5 text-sm hover:bg-bg-hover disabled:opacity-50"
              >
                Generate new codes
              </button>
            </div>

            <div>
              <div className="mb-2 flex items-center justify-between">
                <h3 className="text-sm font-medium text-foreground">Trusted devices</h3>
                {!!devices?.length && (
                  <button
                    type="button"
                    onClick={() => revoke.mutate(undefined)}
                    className="text-xs text-danger hover:underline"
                  >
                    Revoke all
                  </button>
                )}
              </div>

              {devices?.length ? (
                <ul className="space-y-2">
                  {devices.map((d) => (
                    <li
                      key={d.id}
                      className="flex items-center gap-3 rounded-lg border border-border px-3 py-2"
                    >
                      <Monitor className="h-4 w-4 shrink-0 text-text-muted" />
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-sm text-foreground">
                          {describeDevice(d.user_agent).label}
                          {d.current && (
                            <span className="ml-2 rounded bg-accent/15 px-1.5 py-0.5 text-[10px] text-accent">
                              this device
                            </span>
                          )}
                        </p>
                        <p className="text-xs text-text-muted">
                          {d.ip_address} · trusted until{" "}
                          {new Date(d.expires_at).toLocaleDateString()}
                        </p>
                      </div>
                      <button
                        type="button"
                        onClick={() => revoke.mutate(d.id)}
                        aria-label="Revoke this device"
                        className="rounded p-1 text-text-muted hover:bg-bg-hover hover:text-danger"
                      >
                        <X className="h-4 w-4" />
                      </button>
                    </li>
                  ))}
                </ul>
              ) : (
                <p className="text-sm text-text-muted">
                  No trusted devices. Every sign-in asks for a code.
                </p>
              )}
            </div>

            <div className="border-t border-border pt-4">
              {disabling ? (
                <div className="flex flex-wrap items-center gap-2">
                  <input
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="Confirm your password"
                    className="w-56 rounded-lg border border-border bg-bg-primary px-3 py-2 text-sm text-foreground"
                  />
                  <button
                    type="button"
                    onClick={confirmDisable}
                    disabled={disable.isPending || !password}
                    className="rounded-lg bg-danger px-3 py-2 text-sm text-white disabled:opacity-50"
                  >
                    Turn off two-factor
                  </button>
                  <button
                    type="button"
                    onClick={() => { setDisabling(false); setPassword(""); }}
                    className="rounded-lg border border-border px-3 py-2 text-sm hover:bg-bg-hover"
                  >
                    Cancel
                  </button>
                </div>
              ) : (
                <button
                  type="button"
                  onClick={() => setDisabling(true)}
                  className="text-sm text-danger hover:underline"
                >
                  Turn off two-factor authentication
                </button>
              )}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
`
}
