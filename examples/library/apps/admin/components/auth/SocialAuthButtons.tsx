"use client";

import { useState } from "react";
import { isSocialAuthEnabled } from "@repo/shared/themes";
import { API_URL, apiUrl } from "@/lib/api-core";

// apiRedirect turns a path the API answered with into an address on the API,
// or null when it is not a path or resolves to another origin. Appending the
// value to the API's address let "@evil.example" become
// http://api.example.com@evil.example, which is evil.example.
function apiRedirect(path: unknown): string | null {
  if (typeof path !== "string" || !path.startsWith("/") || path.startsWith("//") || path.startsWith("/\\")) {
    return null;
  }
  try {
    const api = new URL(API_URL, window.location.href);
    const target = new URL(apiUrl(path), api);
    return target.origin === api.origin ? target.href : null;
  } catch {
    return null;
  }
}

// SSOSignIn is the enterprise entry point. Rather than listing every customer's
// identity provider on a public page (which leaks your customer list), the user
// types their work address and the server decides where it belongs. An address
// with no connection falls back to the password form, which is what happens for
// most users of most apps.
export function SSOSignIn() {
  const [open, setOpen] = useState(false);
  const [email, setEmail] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");

  const go = async () => {
    if (!email.trim() || busy) return;
    setBusy(true);
    setError("");
    try {
      const res = await fetch(apiUrl("/api/auth/sso/discover"), {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: email.trim() }),
      });
      const body = await res.json();
      // Follow the server's redirect_url rather than assembling one here: OIDC
      // and SAML start at different paths, and the server already knows which
      // protocol this connection uses. Building the URL client-side is how a
      // SAML user ends up at the OIDC endpoint and is told their sign-in method
      // doesn't exist.
      const target = res.ok && body?.data?.sso ? apiRedirect(body.data.redirect_url) : null;
      if (target) {
        window.location.href = target;
        return;
      }
      setError("No single sign-on is set up for that address. Use your password below.");
    } catch {
      setError("Could not check that address. Please try again.");
    } finally {
      setBusy(false);
    }
  };

  if (!open) {
    return (
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="w-full rounded-[var(--auth-radius)] border border-[var(--auth-border)] px-4 py-2.5 text-sm font-medium text-[var(--auth-text)] transition-colors hover:bg-black/5"
      >
        Sign in with SSO
      </button>
    );
  }

  return (
    <div className="space-y-2">
      <input
        type="email"
        autoFocus
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        onKeyDown={(e) => { if (e.key === "Enter") go(); }}
        placeholder="you@company.com"
        className="w-full rounded-[var(--auth-radius)] border border-[var(--auth-border)] bg-transparent px-4 py-2.5 text-sm text-[var(--auth-text)] outline-none focus:border-[var(--auth-accent)]"
      />
      <button
        type="button"
        onClick={go}
        disabled={busy}
        className="w-full rounded-[var(--auth-radius)] bg-[var(--auth-accent)] px-4 py-2.5 text-sm font-semibold text-white disabled:opacity-50"
      >
        {busy ? "Checking…" : "Continue"}
      </button>
      {error && <p className="text-xs text-red-500">{error}</p>}
    </div>
  );
}

export function SocialAuthButtons() {
  if (!isSocialAuthEnabled()) return null;

  return (
    <div className="flex gap-3">
      <a
        href={apiUrl("/api/auth/oauth/google")}
        className="flex flex-1 items-center justify-center gap-2 rounded-[var(--auth-radius)] border border-[var(--auth-border)] bg-white px-4 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 transition-colors"
      >
        <svg className="h-5 w-5" viewBox="0 0 24 24">
          <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4" />
          <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
          <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05" />
          <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
        </svg>
        Google
      </a>
      <a
        href={apiUrl("/api/auth/oauth/github")}
        className="flex flex-1 items-center justify-center gap-2 rounded-[var(--auth-radius)] border border-[var(--auth-border)] bg-[#24292f] px-4 py-2.5 text-sm font-medium text-white hover:bg-[#2f363d] transition-colors"
      >
        <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
          <path fillRule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z" clipRule="evenodd" />
        </svg>
        GitHub
      </a>
    </div>
  );
}

export function SocialAuthDivider() {
  if (!isSocialAuthEnabled()) return null;
  return (
    <div className="relative my-6">
      <div className="absolute inset-0 flex items-center">
        <div className="w-full border-t border-[var(--auth-border)]" />
      </div>
      <div className="relative flex justify-center text-xs">
        <span className="bg-[var(--auth-card)] px-3 text-[var(--auth-muted)]">
          or continue with
        </span>
      </div>
    </div>
  );
}
