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
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { ShieldCheck, Loader2, Copy, Check, Monitor, X } from "@/lib/icons";
import { describeDevice } from "@/components/profile/active-sessions";
import { Button, buttonClasses } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";

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
      <p className="text-xs text-text-secondary mb-3">
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
          className={buttonClasses({ size: "sm", className: "ml-auto" })}
          title={saved ? undefined : "Copy or download them first"}
        >
          I have saved them
        </button>
      </div>
    </div>
  );
}

/**
 * Ask for a code, to prove the mailbox works before email becomes the second
 * factor. Turning on a factor you cannot receive is how an account locks
 * itself out, so it is deliberately two steps.
 */
export function useSendEmailSetupCode() {
  return useMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post<{ data: { sent_to: string } }>(
        "/api/auth/totp/email/send",
        {},
      );
      return data.data;
    },
  });
}

/** Confirm that code, and turn it on. */
export function useEnableEmailTwoFactor() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (code: string) => {
      const { data } = await apiClient.post<{ data: { backup_codes: string[]; method: string } }>(
        "/api/auth/totp/email/enable",
        { code: code.trim() },
      );
      return data.data;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["totp-status"] }),
  });
}

export function TwoFactorCard() {
  const { data: status, isLoading } = useTOTPStatus();
  const { data: devices } = useTrustedDevices(!!status?.enabled);
  const setup = useTOTPSetup();
  const enable = useEnableTOTP();
  const sendEmailCode = useSendEmailSetupCode();
  const enableEmail = useEnableEmailTwoFactor();
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
  // Mid-setup for the email method: the address it went to, and the digits.
  const [emailSetupSentTo, setEmailSetupSentTo] = useState<string | null>(null);
  const [emailCode, setEmailCode] = useState("");

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

  const startEmailSetup = () => {
    setError(null);
    sendEmailCode.mutate(undefined, {
      onSuccess: (d) => setEmailSetupSentTo(d.sent_to),
      onError: (e) => setError(errText(e)),
    });
  };

  const confirmEmailSetup = () => {
    setError(null);
    enableEmail.mutate(emailCode, {
      onSuccess: (d) => {
        setCodes(d.backup_codes);
        setEmailSetupSentTo(null);
        setEmailCode("");
      },
      onError: (e) => setError(errText(e)),
    });
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
      <div className="rounded-xl border border-border bg-bg-elevated p-6">
        <Loader2 className="h-4 w-4 animate-spin text-text-muted" />
      </div>
    );
  }

  return (
    <div className="rounded-xl border border-border bg-bg-elevated p-6">
      <div className="flex items-start gap-3">
        <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent/10 text-accent">
          <ShieldCheck className="h-4 w-4" aria-hidden="true" />
        </span>
        <div className="min-w-0">
          <h2 className="text-base font-semibold text-foreground">Two-factor authentication</h2>
          <p className="mt-1 text-sm leading-relaxed text-text-muted">
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

      <div className="mt-5 space-y-5">
        {error && (
          <p className="rounded-lg bg-danger/10 px-3 py-2 text-sm text-danger">{error}</p>
        )}

        {codes && <BackupCodes codes={codes} onDone={() => setCodes(null)} />}

        {/* ── Off, and not mid-setup ── */}
        {!status?.enabled && !secret && !codes && !emailSetupSentTo && (
          <div className="flex flex-wrap gap-2">
            {/* The app first, and not only for looks: a code in an
                authenticator never leaves the device, while a code by email is
                only as safe as the mailbox, which is also where a password
                reset goes. */}
            <Button onClick={startSetup} loading={setup.isPending}>
              Set up an authenticator app
            </Button>
            <Button variant="outline" onClick={startEmailSetup} loading={sendEmailCode.isPending}>
              Use codes by email
            </Button>
          </div>
        )}

        {/* ── Mid-setup by email: prove the mailbox works ── */}
        {emailSetupSentTo && !codes && (
          <div className="space-y-3">
            <p className="text-sm text-text-secondary">
              We sent a code to <span className="font-medium text-foreground">{emailSetupSentTo}</span>.
              Enter it to turn this on, so nothing is switched on that you cannot receive.
            </p>
            <div className="flex flex-wrap items-end gap-2">
              <div className="space-y-1.5">
                <label htmlFor="two-factor-email-code" className="block text-xs font-medium text-text-secondary">
                  The code from your email
                </label>
                <input
                  id="two-factor-email-code"
                  inputMode="numeric"
                  autoComplete="one-time-code"
                  maxLength={6}
                  value={emailCode}
                  onChange={(e) => setEmailCode(e.target.value.replace(/[^0-9]/g, ""))}
                  className={inputClasses({ className: "w-40 font-mono text-lg tracking-[0.3em]" })}
                />
              </div>
              <Button onClick={confirmEmailSetup} loading={enableEmail.isPending} disabled={emailCode.length !== 6}>
                Turn it on
              </Button>
              <Button
                variant="ghost"
                onClick={() => {
                  setEmailSetupSentTo(null);
                  setEmailCode("");
                }}
              >
                Cancel
              </Button>
            </div>
          </div>
        )}

        {/* ── Mid-setup: scan, then confirm with a live code ── */}
        {secret && qr && (
          <div className="grid gap-6 sm:grid-cols-[auto_1fr]">
            <div className="rounded-lg border border-border bg-white p-3">
              {/* Rendered by the API, so no QR library ships in the client. */}
              <img src={qr} alt="Two-factor setup QR code" width={176} height={176} />
            </div>

            <div className="min-w-0 space-y-3">
              <p className="text-sm text-text-secondary">
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
                  className="w-32 rounded-lg border border-border bg-background px-3 py-2 text-center font-mono tracking-[0.3em] text-foreground"
                />
                <Button
                  onClick={confirmEnable}
                  disabled={code.trim().length !== 6}
                  loading={enable.isPending}
                >
                  Verify and turn on
                </Button>
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
              <span className="text-sm text-text-secondary">
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
                    className="w-56 rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground"
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
