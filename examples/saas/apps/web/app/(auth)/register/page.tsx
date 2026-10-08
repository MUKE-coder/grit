"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api";
import { getApiErrorMessage } from "@/lib/api-core";
import { AuthShell } from "@/components/auth/AuthShell";
import { setWebSessionMarker } from "@/lib/web-session";
import { PasswordStrength } from "@/components/password-strength";

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
      setError(getApiErrorMessage(err, "Registration failed"));
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
            autoComplete="new-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className={inputOk}
            placeholder="At least 8 characters"
            aria-describedby="password-strength"
            minLength={8}
            required
          />
          <PasswordStrength value={password} about={[email, firstName, lastName]} />
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
