"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { apiClient } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-core";
import { AuthShell } from "@/components/auth/AuthShell";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";
const inputErr = inputBase + " border-red-400 focus:border-red-500 focus:ring-red-400/30";

// Confirmation lives here rather than in the shared schema: the API only wants
// token + password, and a mistyped confirmation should never reach it.
const FormSchema = z
  .object({
    password: z.string().min(8, "Password must be at least 8 characters"),
    confirm_password: z.string().min(1, "Please confirm your password"),
  })
  .refine((d) => d.password === d.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });
type FormValues = z.infer<typeof FormSchema>;

function ResetPasswordForm() {
  const router = useRouter();
  const token = useSearchParams().get("token") ?? "";
  const [error, setError] = useState("");
  const [done, setDone] = useState(false);
  const [loading, setLoading] = useState(false);

  const { register, handleSubmit, formState: { errors } } = useForm<FormValues>({
    resolver: zodResolver(FormSchema),
  });

  const onSubmit = async (data: FormValues) => {
    setError("");
    setLoading(true);
    try {
      await apiClient.post("/api/auth/reset-password", {
        token,
        password: data.password,
      });
      setDone(true);
      // Every device was signed out by the reset, so there is nothing to
      // return to except a fresh login.
      setTimeout(() => router.push("/login"), 2500);
    } catch (err: unknown) {
      setError(getApiErrorMessage(err, "Something went wrong. Please try again."));
    } finally {
      setLoading(false);
    }
  };

  // No token in the URL means the link was mangled or typed by hand. Say so
  // before asking for a password that could never be saved.
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
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          <div className="space-y-2">
            <label htmlFor="password" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              New password
            </label>
            <input
              id="password"
              type="password"
              autoComplete="new-password"
              {...register("password")}
              className={errors.password ? inputErr : inputOk}
              placeholder="At least 8 characters"
              autoFocus
            />
            {errors.password && <p className="text-sm text-red-500">{errors.password.message}</p>}
          </div>

          <div className="space-y-2">
            <label htmlFor="confirm_password" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              Confirm password
            </label>
            <input
              id="confirm_password"
              type="password"
              autoComplete="new-password"
              {...register("confirm_password")}
              className={errors.confirm_password ? inputErr : inputOk}
              placeholder="Re-enter your password"
            />
            {errors.confirm_password && <p className="text-sm text-red-500">{errors.confirm_password.message}</p>}
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

// useSearchParams needs a Suspense boundary or the whole route opts out of
// static rendering and the build warns.
export default function ResetPasswordPage() {
  return (
    <Suspense fallback={null}>
      <ResetPasswordForm />
    </Suspense>
  );
}
