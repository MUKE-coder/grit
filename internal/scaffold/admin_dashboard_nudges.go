package scaffold

// The one thing on the dashboard that is about the reader's own account.
//
// An account with no second factor is one leaked password from gone, and that
// fact is visible nowhere the person looks: it lives on a settings page they
// have no reason to open. So it is put in front of them once, on the screen
// they land on, and then it goes away. The card can be dismissed, and it
// disappears on its own the moment two-factor is on. A banner that cannot be
// dismissed is a banner people learn to look past, which costs the next one its
// attention too.
//
// Deliberately not also an "confirm your email" card. EmailVerifiedBanner has
// done that since long before this file existed, and it does it on every page
// rather than only this one, so a second copy here would be two prompts for one
// job stacked on top of each other. It shipped that way in v3.320.0 and was
// caught by looking at a screenshot of a running admin, which is the argument
// for looking at one.

func adminDashboardNudgesTSX() string {
	return `"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { AlertTriangle, X } from "@/lib/icons";
import { useTOTPStatus } from "@/hooks/use-auth";

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
  const [dismissed, setDismissed] = useState<string[]>([]);

  // Read after mount: localStorage does not exist while this renders on the
  // server, and reading it during render would mismatch the hydration.
  useEffect(() => setDismissed(readDismissed()), []);

  const dismiss = (id: string) => {
    const next = Array.from(new Set([...dismissed, id]));
    setDismissed(next);
    writeDismissed(next);
  };

  // Only once the server has answered. Flashing "turn on two-factor" at
  // somebody who already has it, for the half second before the query lands,
  // is worse than showing nothing.
  const needsTwoFactor = totp !== undefined && !totp.enabled;
  if (!needsTwoFactor || dismissed.includes("two-factor")) return null;

  return (
    <div className="mb-6">
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
    </div>
  );
}
`
}
