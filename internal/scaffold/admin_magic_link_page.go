package scaffold

// The page an emailed sign-in link opens: apps/admin/app/(auth)/magic-link/page.tsx.
//
// It spends the token with a POST as soon as it loads, rather than the server
// spending it on the GET that opened it. Corporate mail scanners follow every
// link in a message before anybody reads it, and a token spent by a GET is a
// token the scanner burns: the person then clicks their own link and is told
// it has already been used. Nothing about that failure suggests its cause,
// which is why it survives in so many products.
//
// It is written beside the other auth pages rather than through
// adminPageFiles, because those put a page inside the dashboard layout and
// this one is reached by somebody who is not signed in yet.

func adminMagicLinkPage() string {
	return `"use client";

import { Suspense, useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { apiClient } from "@/lib/api-client";
import { AuthShell } from "@/components/auth/AuthShell";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-center font-mono text-lg tracking-[0.4em] text-[var(--auth-fg)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";

const errText = (e: unknown) =>
  (e as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;

function MagicLink() {
  const params = useSearchParams();
  const router = useRouter();
  const token = params.get("token") ?? "";
  const [error, setError] = useState("");
  const [pendingToken, setPendingToken] = useState<string | null>(null);
  const [factorMethod, setFactorMethod] = useState("app");
  const [code, setCode] = useState("");
  const [verifying, setVerifying] = useState(false);
  // Spent once, even under React's double-invoked effects in development: a
  // link works once, and firing twice would report the second as used.
  const spent = useRef(false);

  useEffect(() => {
    if (!token || spent.current) return;
    spent.current = true;
    apiClient
      .post("/api/auth/magic-link/consume", { token })
      .then((res) => {
        const data = (res.data as { data?: { totp_required?: boolean; pending_token?: string; method?: string } })?.data;
        if (data?.totp_required && data.pending_token) {
          // A link replaces the password, not the second factor.
          setPendingToken(data.pending_token);
          setFactorMethod(data.method ?? "app");
          return;
        }
        router.replace("/dashboard");
      })
      .catch((err: unknown) => {
        setError(errText(err) ?? "That link did not work. Ask for a new one.");
      });
  }, [token, router]);

  const onVerify = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!pendingToken) return;
    setVerifying(true);
    setError("");
    try {
      await apiClient.post("/api/auth/totp/verify", { pending_token: pendingToken, code });
      router.replace("/dashboard");
    } catch (err: unknown) {
      setError(errText(err) ?? "That code was not accepted.");
      setVerifying(false);
    }
  };

  // A link with no token is one that was truncated by a mail client or typed
  // by hand. Say that, rather than showing a spinner that never resolves.
  if (!token) {
    return (
      <AuthShell
        mode="login"
        title="This link is incomplete"
        subtitle="Use the most recent email, or sign in with your password."
        showSocial={false}
      >
        <Link
          href="/login"
          className="block w-full rounded-[var(--auth-radius)] py-3 text-center font-medium transition-colors"
          style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
        >
          Back to sign in
        </Link>
      </AuthShell>
    );
  }

  if (pendingToken) {
    return (
      <AuthShell
        mode="login"
        title="Two-factor authentication"
        subtitle={
          factorMethod === "email"
            ? "We emailed you a 6-digit code. It expires in a few minutes."
            : "Enter the 6-digit code from your authenticator app"
        }
        errorMessage={error}
        showSocial={false}
      >
        <form onSubmit={onVerify} className="space-y-5">
          <div className="space-y-2">
            <label htmlFor="totp-code" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              Authentication code
            </label>
            <input
              id="totp-code"
              inputMode="numeric"
              autoComplete="one-time-code"
              maxLength={6}
              value={code}
              onChange={(e) => setCode(e.target.value.replace(/[^0-9]/g, ""))}
              className={inputOk}
              autoFocus
            />
          </div>
          <button
            type="submit"
            disabled={code.length !== 6 || verifying}
            className="w-full rounded-[var(--auth-radius)] py-3 font-medium disabled:opacity-50 transition-colors"
            style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
          >
            {verifying ? "Verifying..." : "Verify and sign in"}
          </button>
        </form>
      </AuthShell>
    );
  }

  if (error) {
    return (
      <AuthShell mode="login" title="That link did not work" subtitle={error} showSocial={false}>
        <Link
          href="/login"
          className="block w-full rounded-[var(--auth-radius)] py-3 text-center font-medium transition-colors"
          style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
        >
          Ask for a new one
        </Link>
      </AuthShell>
    );
  }

  return (
    <AuthShell mode="login" title="Signing you in" subtitle="One moment." showSocial={false}>
      <p role="status" aria-live="polite" className="text-center text-sm" style={{ color: "var(--auth-muted)" }}>
        Checking your link...
      </p>
    </AuthShell>
  );
}

export default function MagicLinkPage() {
  // useSearchParams needs a Suspense boundary or the whole route opts out of
  // static rendering and the build says so.
  return (
    <Suspense fallback={null}>
      <MagicLink />
    </Suspense>
  );
}
`
}
