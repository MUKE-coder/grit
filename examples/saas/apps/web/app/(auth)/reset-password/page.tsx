"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { api } from "@/lib/api";
import { getApiErrorMessage } from "@/lib/api-core";
import { AuthShell } from "@/components/auth/AuthShell";
import { PasswordStrength } from "@/components/password-strength";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";

function ResetPasswordForm() {
  const router = useRouter();
  const token = useSearchParams().get("token") ?? "";
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [error, setError] = useState("");
  const [done, setDone] = useState(false);
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (password.length < 8) {
      setError("Password must be at least 8 characters");
      return;
    }
    if (password !== confirm) {
      setError("Passwords do not match");
      return;
    }

    setLoading(true);
    try {
      await api.post("/api/auth/reset-password", { token, password });
      setDone(true);
      // The reset signed every device out, so there is nowhere to go but login.
      setTimeout(() => router.push("/login"), 2500);
    } catch (err: unknown) {
      setError(getApiErrorMessage(err, "Something went wrong. Please try again."));
    } finally {
      setLoading(false);
    }
  };

  // A link with no token was mangled in transit or typed by hand. Say so rather
  // than collecting a password that could never be saved.
  if (!token) {
    return (
      <AuthShell
        mode="reset"
        title="This link is incomplete"
        subtitle="Request a new password reset link and use the most recent email."
        showSocial={false}
      >
        <a
          href="/forgot-password"
          className="block w-full rounded-[var(--auth-radius)] py-3 text-center font-medium transition-colors"
          style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
        >
          Request a new link
        </a>
      </AuthShell>
    );
  }

  return (
    <AuthShell
      mode="reset"
      title={done ? "Password updated" : "Choose a new password"}
      subtitle={
        done
          ? "You have been signed out everywhere. Redirecting you to sign in..."
          : "Pick something at least 8 characters long. This also signs you out of every device."
      }
      errorMessage={error}
      showSocial={false}
    >
      {!done && (
        <form onSubmit={onSubmit} className="space-y-5">
          <div className="space-y-2">
            <label htmlFor="password" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              New password
            </label>
            <input
              id="password"
              type="password"
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={inputOk}
              placeholder="At least 8 characters"
              aria-describedby="password-strength"
              autoFocus
            />
            <PasswordStrength value={password} />
          </div>

          <div className="space-y-2">
            <label htmlFor="confirm" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              Confirm password
            </label>
            <input
              id="confirm"
              type="password"
              autoComplete="new-password"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              className={inputOk}
              placeholder="Re-enter your password"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-[var(--auth-radius)] py-3 font-medium disabled:opacity-50 transition-colors"
            style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
          >
            {loading ? "Updating..." : "Update password"}
          </button>
        </form>
      )}
    </AuthShell>
  );
}

// useSearchParams needs a Suspense boundary or the route opts out of static
// rendering and the build warns.
export default function ResetPasswordPage() {
  return (
    <Suspense fallback={null}>
      <ResetPasswordForm />
    </Suspense>
  );
}
