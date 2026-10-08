"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { api } from "@/lib/api";
import { getApiErrorMessage } from "@/lib/api-core";
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
      setError(getApiErrorMessage(err, "Invalid email or password"));
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
              autoComplete="current-password"
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
