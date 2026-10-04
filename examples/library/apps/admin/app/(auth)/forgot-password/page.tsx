"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ForgotPasswordSchema, type ForgotPasswordInput } from "@repo/shared/schemas";
import { apiClient } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-core";
import { AuthShell } from "@/components/auth/AuthShell";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";
const inputErr = inputBase + " border-red-400 focus:border-red-500 focus:ring-red-400/30";

export default function ForgotPasswordPage() {
  const [submitted, setSubmitted] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const { register, handleSubmit, formState: { errors } } = useForm<ForgotPasswordInput>({
    resolver: zodResolver(ForgotPasswordSchema),
  });

  const onSubmit = async (data: ForgotPasswordInput) => {
    setError("");
    setLoading(true);
    try {
      await apiClient.post("/api/auth/forgot-password", data);
      setSubmitted(true);
    } catch (err: unknown) {
      setError(getApiErrorMessage(err, "Something went wrong. Please try again."));
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
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          <div className="space-y-2">
            <label htmlFor="email" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              Email
            </label>
            <input
              id="email"
              type="email"
              {...register("email")}
              className={errors.email ? inputErr : inputOk}
              placeholder="you@example.com"
              autoFocus
            />
            {errors.email && <p className="text-sm text-red-500">{errors.email.message}</p>}
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
