package scaffold

// Email verification, admin side.
//
// Two pieces, because the flow has two halves that happen in different places:
//
//   /verify-email      the page the mail link lands on. Public, because the
//                      user is usually not signed in when they open mail.
//   EmailVerifiedBanner  the nudge for a signed-in user who never clicked it,
//                      with a resend button.
//
// The page consumes the token on mount rather than behind a "Confirm" button.
// A confirmation step would be the right call if this were destructive, but it
// is not — and an extra click between the mail and the outcome is a place for
// people to drop out.

func adminVerifyEmailPage() string {
	return `"use client";

import { Suspense, useEffect, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { apiClient } from "@/lib/api-client";
import { AuthShell } from "@/components/auth/AuthShell";
import { Check, AlertTriangle, Loader2 } from "@/lib/icons";

type State = "working" | "done" | "failed";

function VerifyEmailInner() {
  const params = useSearchParams();
  const token = params.get("token") ?? "";
  const [state, setState] = useState<State>("working");
  const [message, setMessage] = useState("");

  useEffect(() => {
    if (!token) {
      setState("failed");
      setMessage("That link is missing its token. Copy the whole URL from the email.");
      return;
    }

    let cancelled = false;
    (async () => {
      try {
        await apiClient.post("/api/auth/verify-email", { token });
        if (!cancelled) setState("done");
      } catch (e) {
        if (cancelled) return;
        const msg = (e as { response?: { data?: { error?: { message?: string } } } })
          ?.response?.data?.error?.message;
        setState("failed");
        setMessage(msg ?? "That link is invalid or has expired.");
      }
    })();

    // The token is single-use: a re-run under React 18 StrictMode would spend
    // it twice and show a failure for a verification that actually worked.
    return () => { cancelled = true; };
  }, [token]);

  return (
    <AuthShell
      mode="login"
      title={
        state === "working" ? "Confirming your email…"
        : state === "done" ? "Email confirmed"
        : "We could not confirm that link"
      }
      subtitle={
        state === "done"
          ? "Thanks — your address is verified."
          : state === "failed"
            ? message
            : "One moment."
      }
    >
      <div className="flex flex-col items-center gap-5 py-4">
        {state === "working" && (
          <Loader2 className="h-8 w-8 animate-spin" style={{ color: "var(--auth-primary)" }} />
        )}

        {state === "done" && (
          <>
            <div
              className="flex h-12 w-12 items-center justify-center rounded-full"
              style={{ background: "var(--auth-primary)" }}
            >
              <Check className="h-6 w-6 text-white" />
            </div>
            <Link
              href="/login"
              className="w-full rounded-[var(--auth-radius)] py-3 text-center font-medium text-white"
              style={{ background: "var(--auth-primary)" }}
            >
              Continue to sign in
            </Link>
          </>
        )}

        {state === "failed" && (
          <>
            <AlertTriangle className="h-8 w-8 text-red-500" />
            <p className="text-center text-sm" style={{ color: "var(--auth-muted)" }}>
              Sign in and use the banner at the top to send yourself a fresh link.
            </p>
            <Link
              href="/login"
              className="w-full rounded-[var(--auth-radius)] py-3 text-center font-medium text-white"
              style={{ background: "var(--auth-primary)" }}
            >
              Go to sign in
            </Link>
          </>
        )}
      </div>
    </AuthShell>
  );
}

export default function VerifyEmailPage() {
  return (
    <Suspense
      fallback={
        <AuthShell mode="login" title="Confirming your email…" subtitle="One moment.">
          <div className="flex justify-center py-4">
            <Loader2 className="h-8 w-8 animate-spin" style={{ color: "var(--auth-primary)" }} />
          </div>
        </AuthShell>
      }
    >
      <VerifyEmailInner />
    </Suspense>
  );
}
`
}

func adminEmailVerifiedBanner() string {
	return `"use client";

import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useMe } from "@/hooks/use-auth";
import { apiClient } from "@/lib/api-client";
import { AlertTriangle, Check, Loader2 } from "@/lib/icons";

/**
 * Shown to a signed-in user whose address is still unconfirmed.
 *
 * Renders nothing at all when there is nothing to say — including while the
 * user query is still loading, so it never flashes in and out on every page
 * load for people who verified months ago.
 */
export function EmailVerifiedBanner() {
  const { data: user, isLoading } = useMe();
  const [sent, setSent] = useState(false);

  const resend = useMutation({
    mutationFn: async () => {
      await apiClient.post("/api/auth/verify-email/send", {});
    },
    onSuccess: () => setSent(true),
  });

  if (isLoading || !user || user.email_verified_at) return null;

  return (
    <div className="mx-6 mt-4 flex flex-wrap items-center gap-3 rounded-lg border border-warning/40 bg-warning/[0.07] px-4 py-3">
      <AlertTriangle className="h-4 w-4 shrink-0 text-warning" />
      <p className="min-w-0 flex-1 text-sm text-foreground">
        Confirm your email address.{" "}
        <span className="text-foreground-secondary">
          We sent a link to {user.email} when you signed up.
        </span>
      </p>

      {sent ? (
        <span className="inline-flex items-center gap-1.5 text-sm text-success">
          <Check className="h-3.5 w-3.5" />
          Sent — check your inbox
        </span>
      ) : (
        <button
          type="button"
          onClick={() => resend.mutate()}
          disabled={resend.isPending}
          className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-3 py-1.5 text-sm hover:bg-bg-hover disabled:opacity-50"
        >
          {resend.isPending && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
          Resend link
        </button>
      )}
    </div>
  );
}
`
}
