package scaffold

// Themed auth page generators (v3.28). Each function returns a thin page
// that wraps its form in the shared <AuthShell> dispatcher — the shell
// reads NEXT_PUBLIC_THEME at the call site and renders the matching
// layout (Atlas / Aurora / Pulse). The form markup itself stays identical
// across themes; only the chrome around it changes.
//
// Pages live at:
//   apps/admin/app/(auth)/login/page.tsx
//   apps/admin/app/(auth)/sign-up/page.tsx
//   apps/admin/app/(auth)/forgot-password/page.tsx
//
// Inputs use CSS variables (--auth-border, --auth-primary, --auth-card,
// --auth-muted, --auth-radius) set by the shell at the layout root so the
// same JSX picks up theme colors without per-theme branches.

func adminThemedLoginPage() string {
	return `"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Eye, EyeOff } from "@/lib/icons";
import { useLogin, useMe } from "@/hooks/use-auth";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginSchema, type LoginInput } from "@repo/shared/schemas";
import { AuthShell } from "@/components/auth/AuthShell";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";
const inputErr = inputBase + " border-red-400 focus:border-red-500 focus:ring-red-400/30";

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const { mutate: login, isPending, error: serverError } = useLogin();
  const { data: existingUser, isLoading: meLoading } = useMe();
  const router = useRouter();
  const { register, handleSubmit, formState: { errors } } = useForm<LoginInput>({
    resolver: zodResolver(LoginSchema),
  });

  // v3.31.15: if the session cookie is still valid, don't show the
  // login form — bounce straight to the dashboard.
  useEffect(() => {
    if (!meLoading && existingUser) {
      router.replace(existingUser.role === "USER" ? "/profile" : "/dashboard");
    }
  }, [meLoading, existingUser, router]);

  const onSubmit = (data: LoginInput) => login(data);

  const message = (serverError as unknown as { response?: { data?: { error?: { message?: string } } } })
    ?.response?.data?.error?.message;

  return (
    <AuthShell
      mode="login"
      title="Welcome back"
      subtitle="Sign in to your account"
      errorMessage={message}
    >
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

        <div className="space-y-2">
          <label htmlFor="password" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Password
          </label>
          <div className="relative">
            <input
              id="password"
              type={showPassword ? "text" : "password"}
              {...register("password")}
              className={(errors.password ? inputErr : inputOk) + " pr-12"}
              placeholder="Enter your password"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2"
              style={{ color: "var(--auth-muted)" }}
              aria-label={showPassword ? "Hide password" : "Show password"}
            >
              {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
            </button>
          </div>
          {errors.password && <p className="text-sm text-red-500">{errors.password.message}</p>}
        </div>

        <div className="flex items-center justify-between text-sm">
          <label className="flex items-center gap-2 cursor-pointer" style={{ color: "var(--auth-muted)" }}>
            <input type="checkbox" className="h-4 w-4 rounded border-[var(--auth-border)]" />
            Remember me
          </label>
          <Link href="/forgot-password" style={{ color: "var(--auth-primary)" }}>
            Forgot password?
          </Link>
        </div>

        <button
          type="submit"
          disabled={isPending}
          className="w-full rounded-[var(--auth-radius)] py-3 font-medium disabled:opacity-50 transition-colors"
          style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
        >
          {isPending ? "Signing in..." : "Sign In"}
        </button>
      </form>
    </AuthShell>
  );
}
`
}

func adminThemedSignUpPage() string {
	return `"use client";

import { useState } from "react";
import { Eye, EyeOff } from "@/lib/icons";
import { useRegister } from "@/hooks/use-auth";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { RegisterSchema, type RegisterInput } from "@repo/shared/schemas";
import { AuthShell } from "@/components/auth/AuthShell";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";
const inputErr = inputBase + " border-red-400 focus:border-red-500 focus:ring-red-400/30";

export default function SignUpPage() {
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirmPassword, setShowConfirmPassword] = useState(false);
  const { mutate: registerUser, isPending, error: serverError } = useRegister();

  const { register, handleSubmit, formState: { errors } } = useForm<RegisterInput>({
    resolver: zodResolver(RegisterSchema),
  });

  const onSubmit = (data: RegisterInput) => {
    registerUser({
      first_name: data.firstName,
      last_name: data.lastName,
      email: data.email,
      password: data.password,
    });
  };

  const message = (serverError as unknown as { response?: { data?: { error?: { message?: string } } } })
    ?.response?.data?.error?.message;

  return (
    <AuthShell
      mode="sign-up"
      title="Create your account"
      subtitle="Sign up to get started"
      errorMessage={message}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-2">
            <label htmlFor="firstName" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              First name
            </label>
            <input
              id="firstName"
              type="text"
              {...register("firstName")}
              className={errors.firstName ? inputErr : inputOk}
              placeholder="Jane"
            />
            {errors.firstName && <p className="text-sm text-red-500">{errors.firstName.message}</p>}
          </div>
          <div className="space-y-2">
            <label htmlFor="lastName" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              Last name
            </label>
            <input
              id="lastName"
              type="text"
              {...register("lastName")}
              className={errors.lastName ? inputErr : inputOk}
              placeholder="Doe"
            />
            {errors.lastName && <p className="text-sm text-red-500">{errors.lastName.message}</p>}
          </div>
        </div>

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
          />
          {errors.email && <p className="text-sm text-red-500">{errors.email.message}</p>}
        </div>

        <div className="space-y-2">
          <label htmlFor="password" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Password
          </label>
          <div className="relative">
            <input
              id="password"
              type={showPassword ? "text" : "password"}
              {...register("password")}
              className={(errors.password ? inputErr : inputOk) + " pr-12"}
              placeholder="At least 8 characters"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2"
              style={{ color: "var(--auth-muted)" }}
              aria-label={showPassword ? "Hide password" : "Show password"}
            >
              {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
            </button>
          </div>
          {errors.password && <p className="text-sm text-red-500">{errors.password.message}</p>}
        </div>

        <div className="space-y-2">
          <label htmlFor="confirmPassword" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Confirm password
          </label>
          <div className="relative">
            <input
              id="confirmPassword"
              type={showConfirmPassword ? "text" : "password"}
              {...register("confirmPassword")}
              className={(errors.confirmPassword ? inputErr : inputOk) + " pr-12"}
              placeholder="Re-enter your password"
            />
            <button
              type="button"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2"
              style={{ color: "var(--auth-muted)" }}
              aria-label={showConfirmPassword ? "Hide password" : "Show password"}
            >
              {showConfirmPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
            </button>
          </div>
          {errors.confirmPassword && <p className="text-sm text-red-500">{errors.confirmPassword.message}</p>}
        </div>

        <button
          type="submit"
          disabled={isPending}
          className="w-full rounded-[var(--auth-radius)] py-3 font-medium disabled:opacity-50 transition-colors"
          style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
        >
          {isPending ? "Creating account..." : "Create account"}
        </button>
      </form>
    </AuthShell>
  );
}
`
}

// adminThemedResetPasswordPage is where the emailed link lands. Without it the
// reset flow dead-ends on a 404, so it ships with the API side that issues the
// link — a working token nobody can use is not a feature.
func adminThemedResetPasswordPage() string {
	return `"use client";

import { Suspense, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { apiClient } from "@/lib/api-client";
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
      const msg = (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
      setError(msg || "Something went wrong. Please try again.");
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
`
}

func adminThemedForgotPasswordPage() string {
	return `"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ForgotPasswordSchema, type ForgotPasswordInput } from "@repo/shared/schemas";
import { apiClient } from "@/lib/api-client";
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
`
}
