package scaffold

// The two things on the dashboard that are about the reader's own account.
//
// Both are one-time jobs that nobody does unprompted. An unverified address
// means a password reset goes nowhere, which is discovered at the worst
// possible moment; an account with no second factor is one leaked password from
// gone. Neither fact is visible anywhere the person looks, because both live on
// a settings page they have no reason to open.
//
// So they are put in front of them once, on the screen they land on, and then
// they go away: each card can be dismissed, and both disappear on their own the
// moment the underlying thing is done. A banner that cannot be dismissed is a
// banner people learn to look past, which costs the next one its attention too.

func adminDashboardNudgesTSX() string {
	return `"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { AlertTriangle, Mail, X } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";
import { useTOTPStatus } from "@/hooks/use-auth";
import { useSecurityOverview } from "@/hooks/use-security";

const DISMISSED_KEY = "grit.dashboard.dismissed-nudges";

function readDismissed(): string[] {
  try {
    const raw = window.localStorage.getItem(DISMISSED_KEY);
    return raw ? (JSON.parse(raw) as string[]) : [];
  } catch {
    // Private windows, blocked site data, and a value somebody hand-edited.
    // A nudge that shows again is a much smaller problem than a crash here.
    return [];
  }
}

function writeDismissed(ids: string[]) {
  try {
    window.localStorage.setItem(DISMISSED_KEY, JSON.stringify(ids));
  } catch {
    // Nothing to do. The card closes for this view either way.
  }
}

interface NudgeProps {
  id: string;
  icon: React.ReactNode;
  title: string;
  body: string;
  action: React.ReactNode;
  onDismiss: (id: string) => void;
}

function Nudge({ id, icon, title, body, action, onDismiss }: NudgeProps) {
  return (
    <div className="relative flex items-start gap-3 rounded-xl border border-warning/25 bg-warning/[0.06] p-4 sm:p-5">
      <div className="mt-0.5 shrink-0">{icon}</div>
      <div className="min-w-0 flex-1">
        <h3 className="font-semibold text-text-primary">{title}</h3>
        <p className="mt-1 text-sm leading-relaxed text-text-secondary">{body}</p>
        <div className="mt-3">{action}</div>
      </div>
      <button
        type="button"
        onClick={() => onDismiss(id)}
        aria-label={"Dismiss: " + title}
        className="-m-1.5 shrink-0 rounded-md p-1.5 text-text-muted transition-colors hover:bg-bg-hover hover:text-text-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
      >
        <X className="h-4 w-4" aria-hidden="true" />
      </button>
    </div>
  );
}

/**
 * Nudges for the signed-in user's own account.
 *
 * Renders nothing at all when there is nothing to say, which is the common
 * case after the first week. Nothing here blocks the dashboard: these sit above
 * the statistics, not in place of them.
 */
export function DashboardSecurityNudges() {
  const { data: totp } = useTOTPStatus();
  const { data: overview } = useSecurityOverview();
  const [dismissed, setDismissed] = useState<string[]>([]);
  const [sending, setSending] = useState(false);
  const [sent, setSent] = useState(false);
  const [sendError, setSendError] = useState("");

  // Read after mount: localStorage does not exist while this renders on the
  // server, and reading it during render would mismatch the hydration.
  useEffect(() => setDismissed(readDismissed()), []);

  const dismiss = (id: string) => {
    const next = Array.from(new Set([...dismissed, id]));
    setDismissed(next);
    writeDismissed(next);
  };

  const resend = async () => {
    setSending(true);
    setSendError("");
    try {
      await apiClient.post("/api/auth/verify-email/send", {});
      setSent(true);
    } catch (err: unknown) {
      setSendError(
        (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error
          ?.message ?? "Could not send the email. Try again in a moment.",
      );
    } finally {
      setSending(false);
    }
  };

  // Only once the server has answered. Flashing "turn on two-factor" at
  // somebody who already has it, for the half second before the query lands,
  // is worse than showing nothing.
  const needsEmail = overview !== undefined && !overview.email_verified;
  const needsTwoFactor = totp !== undefined && !totp.enabled;

  const show = [
    needsEmail && !dismissed.includes("verify-email") ? "verify-email" : null,
    needsTwoFactor && !dismissed.includes("two-factor") ? "two-factor" : null,
  ].filter(Boolean) as string[];

  if (show.length === 0) return null;

  return (
    <div className="mb-6 space-y-3">
      {show.includes("verify-email") && (
        <Nudge
          id="verify-email"
          onDismiss={dismiss}
          icon={<Mail className="h-5 w-5 text-warning" aria-hidden="true" />}
          title="Confirm your email address"
          body={
            "We have not confirmed " +
            (overview?.email ?? "your address") +
            " yet. Until you do, a password reset has nowhere to go, which is a bad thing to find out on the day you need it."
          }
          action={
            sent ? (
              <p role="status" className="text-sm text-success">
                Sent. Open the link in that email.
              </p>
            ) : (
              <div className="flex flex-wrap items-center gap-3">
                <button
                  type="button"
                  onClick={resend}
                  disabled={sending}
                  className="inline-flex h-9 items-center rounded-lg bg-accent px-4 text-sm font-medium text-accent-fg transition-colors hover:bg-accent-hover disabled:opacity-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
                >
                  {sending ? "Sending..." : "Send me the link"}
                </button>
                {sendError && <p className="text-sm text-danger">{sendError}</p>}
              </div>
            )
          }
        />
      )}

      {show.includes("two-factor") && (
        <Nudge
          id="two-factor"
          onDismiss={dismiss}
          icon={<AlertTriangle className="h-5 w-5 text-warning" aria-hidden="true" />}
          title="Turn on two-factor authentication"
          body="Right now your password is the only thing between this account and anybody who has it. Two minutes with an authenticator app, or codes by email if you would rather not install one."
          action={
            <Link
              href="/system/account?tab=security"
              className="inline-flex h-9 items-center rounded-lg bg-accent px-4 text-sm font-medium text-accent-fg transition-colors hover:bg-accent-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
            >
              Set it up
            </Link>
          }
        />
      )}
    </div>
  );
}
`
}
