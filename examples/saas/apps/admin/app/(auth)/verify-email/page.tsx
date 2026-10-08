"use client";

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
