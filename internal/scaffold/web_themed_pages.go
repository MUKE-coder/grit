package scaffold

// Themed web auth pages (v3.28.1). Web uses a simpler auth flow than
// admin — direct api.post calls + useRouter, no react-hook-form — but the
// same <AuthShell> shell so all three themes render correctly. Pages stay
// thin: the shell owns the chrome; pages own the form.

func webThemedLoginPage() string {
	return `"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { api } from "@/lib/api";
import { AuthShell } from "@/components/auth/AuthShell";
import { setWebSessionMarker } from "@/lib/web-session";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";

function LoginForm() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      // API sets HttpOnly grit_access + grit_refresh via Set-Cookie.
      // withCredentials on the axios client carries them on subsequent
      // requests automatically.
      await api.post("/api/auth/login", { email, password });
      // v3.31.42: stamp the web-origin marker so the middleware
      // recognises this session as a web-app session and lets
      // protected pages render. The API cookie alone isn't enough
      // because it's scoped to the API origin -- the web middleware
      // can't see it.
      setWebSessionMarker();
      // Honour ?next= so a bounce off a protected page returns the
      // user to where they were headed.
      const next = searchParams?.get("next");
      router.push(next && next.startsWith("/") ? next : "/");
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
      setError(msg || "Invalid email or password");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthShell
      mode="login"
      title="Welcome back"
      subtitle="Sign in to your account"
      errorMessage={error}
    >
      <form onSubmit={onSubmit} className="space-y-5">
        <div className="space-y-2">
          <label htmlFor="email" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Email
          </label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className={inputOk}
            placeholder="you@example.com"
            required
            autoFocus
          />
        </div>

        <div className="space-y-2">
          <label htmlFor="password" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Password
          </label>
          <div className="relative">
            <input
              id="password"
              type={showPassword ? "text" : "password"}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={inputOk + " pr-16"}
              placeholder="Enter your password"
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-sm"
              style={{ color: "var(--auth-muted)" }}
            >
              {showPassword ? "Hide" : "Show"}
            </button>
          </div>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-[var(--auth-radius)] py-3 font-medium disabled:opacity-50 transition-colors"
          style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
        >
          {loading ? "Signing in..." : "Sign In"}
        </button>
      </form>
    </AuthShell>
  );
}

// useSearchParams reads the URL on the client. Without a Suspense boundary the
// whole route bails out of prerendering and "next build" fails outright.
export default function LoginPage() {
  return (
    <Suspense fallback={null}>
      <LoginForm />
    </Suspense>
  );
}
`
}

func webThemedRegisterPage() string {
	return `"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { AuthShell } from "@/components/auth/AuthShell";
import { setWebSessionMarker } from "@/lib/web-session";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";

export default function RegisterPage() {
  const router = useRouter();
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await api.post("/api/auth/register", {
        first_name: firstName,
        last_name: lastName,
        email,
        password,
      });
      // v3.31.42: stamp the web-origin marker. See lib/web-session.ts.
      setWebSessionMarker();
      router.push("/");
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
      setError(msg || "Registration failed");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthShell
      mode="sign-up"
      title="Create your account"
      subtitle="Sign up to get started"
      errorMessage={error}
    >
      <form onSubmit={onSubmit} className="space-y-5">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <label htmlFor="firstName" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              First name
            </label>
            <input
              id="firstName"
              type="text"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              className={inputOk}
              placeholder="Jane"
              required
            />
          </div>
          <div className="space-y-2">
            <label htmlFor="lastName" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              Last name
            </label>
            <input
              id="lastName"
              type="text"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
              className={inputOk}
              placeholder="Doe"
              required
            />
          </div>
        </div>

        <div className="space-y-2">
          <label htmlFor="email" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Email
          </label>
          <input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className={inputOk}
            placeholder="you@example.com"
            required
          />
        </div>

        <div className="space-y-2">
          <label htmlFor="password" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Password
          </label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className={inputOk}
            placeholder="At least 8 characters"
            minLength={8}
            required
          />
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-[var(--auth-radius)] py-3 font-medium disabled:opacity-50 transition-colors"
          style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
        >
          {loading ? "Creating account..." : "Create Account"}
        </button>
      </form>
    </AuthShell>
  );
}
`
}

// webThemedResetPasswordPage is where the emailed reset link lands when
// OAUTH_FRONTEND_URL points at the web app rather than the admin. Both need the
// page or the flow dead-ends on a 404 depending on how the project is configured.
func webThemedResetPasswordPage() string {
	return `"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { api } from "@/lib/api";
import { AuthShell } from "@/components/auth/AuthShell";

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
      const msg = (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
      setError(msg || "Something went wrong. Please try again.");
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
              autoFocus
            />
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
`
}

func webThemedForgotPasswordPage() string {
	return `"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import { AuthShell } from "@/components/auth/AuthShell";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await api.post("/api/auth/forgot-password", { email });
      setSubmitted(true);
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
      setError(msg || "Something went wrong. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthShell
      mode="forgot"
      title={submitted ? "Check your email" : "Reset your password"}
      subtitle={
        submitted
          ? "If an account exists for that email, a reset link is on its way."
          : "Enter your email and we'll send you a link to reset your password."
      }
      errorMessage={error}
      showSocial={false}
    >
      {!submitted && (
        <form onSubmit={onSubmit} className="space-y-5">
          <div className="space-y-2">
            <label htmlFor="email" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              Email
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className={inputOk}
              placeholder="you@example.com"
              required
              autoFocus
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-[var(--auth-radius)] py-3 font-medium disabled:opacity-50 transition-colors"
            style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
          >
            {loading ? "Sending..." : "Send reset link"}
          </button>
        </form>
      )}
    </AuthShell>
  );
}
`
}
