package scaffold

// Linear-school auth pages: single centered card, minimal chrome, generous whitespace.
// Activated with `grit new <name> --style centered`.
// Follows GRIT_STYLE_GUIDE §7.12 (Auth Pages).

// centeredLoginPage returns the Linear-style centered login page.
func centeredLoginPage() string {
	return `"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useLogin } from "@/hooks/use-auth";
import { LoginSchema, type LoginInput } from "@repo/shared/schemas";
import { Eye, EyeOff, Loader2 } from "@/lib/icons";

export default function LoginPage() {
  const router = useRouter();
  const { mutate: login, isPending, error } = useLogin();
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginInput>({
    resolver: zodResolver(LoginSchema),
  });

  const onSubmit = (data: LoginInput) => {
    login(data, {
      onSuccess: () => router.push("/dashboard"),
    });
  };

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const base = "w-full h-11 rounded-lg border bg-bg-tertiary px-3.5 text-[14px] text-foreground placeholder:text-text-muted transition-colors focus:outline-none focus:ring-2 focus:ring-accent/20";
  const inputClass = ` + "`" + `${base} border-border focus:border-accent` + "`" + `;
  const errorInputClass = ` + "`" + `${base} border-danger focus:border-danger focus:ring-danger/20` + "`" + `;

  return (
    <div className="relative min-h-screen flex items-center justify-center bg-background px-4 py-8">
      {/* Subtle radial glow at top */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-x-0 top-0 h-[400px] opacity-60"
        style={{
          background: "radial-gradient(600px at 50% 0%, rgba(108, 92, 231, 0.10), transparent 70%)",
        }}
      />

      <div className="relative w-full max-w-[420px]">
        {/* Card */}
        <div className="rounded-2xl border border-border bg-bg-secondary p-8 sm:p-10 shadow-sm">
          {/* Logo mark */}
          <div className="flex justify-center">
            <div className="h-10 w-10 rounded-xl bg-accent flex items-center justify-center shadow-sm">
              <span className="text-[17px] font-bold text-white">G</span>
            </div>
          </div>

          {/* Heading */}
          <div className="mt-5 text-center">
            <h1 className="text-[22px] font-semibold text-foreground tracking-tight">
              Welcome back
            </h1>
            <p className="mt-1.5 text-[13px] text-text-secondary">
              Sign in to your account to continue
            </p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit(onSubmit)} className="mt-7 space-y-5">
            {error && (
              <div className="rounded-lg border border-danger/30 bg-danger/10 px-3.5 py-2.5 text-[13px] text-danger">
                {(error as Error).message || "Invalid email or password"}
              </div>
            )}

            <div className="space-y-1.5">
              <label htmlFor="email" className="block text-[13px] font-medium text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="you@example.com"
                className={errors.email ? errorInputClass : inputClass}
                {...register("email")}
              />
              {errors.email && (
                <p className="text-[12px] font-medium text-danger">{errors.email.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <div className="flex items-center justify-between">
                <label htmlFor="password" className="text-[13px] font-medium text-text-secondary">
                  Password
                </label>
                <Link
                  href="/forgot-password"
                  className="text-[12px] font-medium text-accent hover:underline"
                >
                  Forgot password?
                </Link>
              </div>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="current-password"
                  placeholder="Enter your password"
                  className={errors.password ? errorInputClass : inputClass}
                  {...register("password")}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-text-muted hover:text-foreground transition-colors"
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              {errors.password && (
                <p className="text-[12px] font-medium text-danger">{errors.password.message}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={isPending}
              className="w-full h-11 rounded-lg bg-accent text-[14px] font-medium text-white hover:bg-accent-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
              {isPending ? "Signing in..." : "Sign in"}
            </button>
          </form>

          {/* OAuth divider */}
          <div className="relative my-7">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-border" />
            </div>
            <div className="relative flex justify-center">
              <span className="bg-bg-secondary px-3 text-[12px] text-text-muted">
                or continue with
              </span>
            </div>
          </div>

          {/* OAuth buttons */}
          <div className="grid grid-cols-2 gap-2.5">
            <a
              href={` + "`" + `${apiUrl}/api/auth/oauth/google` + "`" + `}
              className="flex items-center justify-center gap-2 h-11 rounded-lg border border-border bg-bg-tertiary text-[13px] font-medium text-foreground hover:bg-bg-hover transition-colors"
            >
              <svg className="h-4 w-4" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
              </svg>
              Google
            </a>
            <a
              href={` + "`" + `${apiUrl}/api/auth/oauth/github` + "`" + `}
              className="flex items-center justify-center gap-2 h-11 rounded-lg border border-border bg-bg-tertiary text-[13px] font-medium text-foreground hover:bg-bg-hover transition-colors"
            >
              <svg className="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 .5C5.37.5 0 5.78 0 12.29c0 5.2 3.44 9.61 8.21 11.17.6.11.82-.26.82-.57v-2.01c-3.34.72-4.04-1.61-4.04-1.61-.55-1.35-1.34-1.72-1.34-1.72-1.09-.74.08-.73.08-.73 1.21.08 1.84 1.22 1.84 1.22 1.07 1.81 2.81 1.29 3.5.99.11-.77.42-1.29.76-1.59-2.67-.3-5.47-1.31-5.47-5.83 0-1.29.47-2.34 1.23-3.17-.12-.3-.54-1.51.12-3.14 0 0 1.01-.32 3.3 1.2.96-.26 1.99-.39 3.02-.4 1.02.01 2.06.14 3.02.4 2.29-1.52 3.3-1.2 3.3-1.2.66 1.63.25 2.84.12 3.14.77.83 1.23 1.88 1.23 3.17 0 4.53-2.81 5.53-5.49 5.82.43.36.81 1.08.81 2.18v3.23c0 .31.22.69.82.57C20.57 21.89 24 17.48 24 12.29 24 5.78 18.63.5 12 .5z" />
              </svg>
              GitHub
            </a>
          </div>

          {/* Footer link */}
          <p className="mt-7 text-center text-[13px] text-text-secondary">
            Don&apos;t have an account?{" "}
            <Link href="/sign-up" className="font-medium text-accent hover:underline">
              Sign up
            </Link>
          </p>
        </div>

        {/* Fine print below card */}
        <p className="mt-6 text-center text-[11px] text-text-muted">
          By signing in, you agree to our{" "}
          <a href="#" className="hover:text-text-secondary underline underline-offset-2">Terms</a>
          {" "}and{" "}
          <a href="#" className="hover:text-text-secondary underline underline-offset-2">Privacy Policy</a>
        </p>
      </div>
    </div>
  );
}
`
}

// centeredSignUpPage returns the Linear-style centered sign-up page.
func centeredSignUpPage() string {
	return `"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRegister } from "@/hooks/use-auth";
import { RegisterSchema, type RegisterInput } from "@repo/shared/schemas";
import { Eye, EyeOff, Loader2 } from "@/lib/icons";

export default function SignUpPage() {
  const router = useRouter();
  const { mutate: registerUser, isPending, error } = useRegister();
  const [showPassword, setShowPassword] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterInput>({
    resolver: zodResolver(RegisterSchema),
  });

  const onSubmit = (data: RegisterInput) => {
    registerUser(data, {
      onSuccess: () => router.push("/dashboard"),
    });
  };

  const apiUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const base = "w-full h-11 rounded-lg border bg-bg-tertiary px-3.5 text-[14px] text-foreground placeholder:text-text-muted transition-colors focus:outline-none focus:ring-2 focus:ring-accent/20";
  const inputClass = ` + "`" + `${base} border-border focus:border-accent` + "`" + `;
  const errorInputClass = ` + "`" + `${base} border-danger focus:border-danger focus:ring-danger/20` + "`" + `;

  return (
    <div className="relative min-h-screen flex items-center justify-center bg-background px-4 py-8">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-x-0 top-0 h-[400px] opacity-60"
        style={{
          background: "radial-gradient(600px at 50% 0%, rgba(108, 92, 231, 0.10), transparent 70%)",
        }}
      />

      <div className="relative w-full max-w-[420px]">
        <div className="rounded-2xl border border-border bg-bg-secondary p-8 sm:p-10 shadow-sm">
          <div className="flex justify-center">
            <div className="h-10 w-10 rounded-xl bg-accent flex items-center justify-center shadow-sm">
              <span className="text-[17px] font-bold text-white">G</span>
            </div>
          </div>

          <div className="mt-5 text-center">
            <h1 className="text-[22px] font-semibold text-foreground tracking-tight">
              Create your account
            </h1>
            <p className="mt-1.5 text-[13px] text-text-secondary">
              Start building with Grit in minutes
            </p>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="mt-7 space-y-5">
            {error && (
              <div className="rounded-lg border border-danger/30 bg-danger/10 px-3.5 py-2.5 text-[13px] text-danger">
                {(error as Error).message || "Registration failed"}
              </div>
            )}

            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-1.5">
                <label htmlFor="first_name" className="block text-[13px] font-medium text-text-secondary">
                  First name
                </label>
                <input
                  id="first_name"
                  type="text"
                  autoComplete="given-name"
                  placeholder="Jane"
                  className={errors.first_name ? errorInputClass : inputClass}
                  {...register("first_name")}
                />
                {errors.first_name && (
                  <p className="text-[12px] font-medium text-danger">{errors.first_name.message}</p>
                )}
              </div>
              <div className="space-y-1.5">
                <label htmlFor="last_name" className="block text-[13px] font-medium text-text-secondary">
                  Last name
                </label>
                <input
                  id="last_name"
                  type="text"
                  autoComplete="family-name"
                  placeholder="Doe"
                  className={errors.last_name ? errorInputClass : inputClass}
                  {...register("last_name")}
                />
                {errors.last_name && (
                  <p className="text-[12px] font-medium text-danger">{errors.last_name.message}</p>
                )}
              </div>
            </div>

            <div className="space-y-1.5">
              <label htmlFor="email" className="block text-[13px] font-medium text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                autoComplete="email"
                placeholder="you@example.com"
                className={errors.email ? errorInputClass : inputClass}
                {...register("email")}
              />
              {errors.email && (
                <p className="text-[12px] font-medium text-danger">{errors.email.message}</p>
              )}
            </div>

            <div className="space-y-1.5">
              <label htmlFor="password" className="block text-[13px] font-medium text-text-secondary">
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="new-password"
                  placeholder="At least 8 characters"
                  className={errors.password ? errorInputClass : inputClass}
                  {...register("password")}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 p-1 text-text-muted hover:text-foreground transition-colors"
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              </div>
              {errors.password && (
                <p className="text-[12px] font-medium text-danger">{errors.password.message}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={isPending}
              className="w-full h-11 rounded-lg bg-accent text-[14px] font-medium text-white hover:bg-accent-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
              {isPending ? "Creating account..." : "Create account"}
            </button>
          </form>

          <div className="relative my-7">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t border-border" />
            </div>
            <div className="relative flex justify-center">
              <span className="bg-bg-secondary px-3 text-[12px] text-text-muted">
                or continue with
              </span>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-2.5">
            <a
              href={` + "`" + `${apiUrl}/api/auth/oauth/google` + "`" + `}
              className="flex items-center justify-center gap-2 h-11 rounded-lg border border-border bg-bg-tertiary text-[13px] font-medium text-foreground hover:bg-bg-hover transition-colors"
            >
              <svg className="h-4 w-4" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
              </svg>
              Google
            </a>
            <a
              href={` + "`" + `${apiUrl}/api/auth/oauth/github` + "`" + `}
              className="flex items-center justify-center gap-2 h-11 rounded-lg border border-border bg-bg-tertiary text-[13px] font-medium text-foreground hover:bg-bg-hover transition-colors"
            >
              <svg className="h-4 w-4" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 .5C5.37.5 0 5.78 0 12.29c0 5.2 3.44 9.61 8.21 11.17.6.11.82-.26.82-.57v-2.01c-3.34.72-4.04-1.61-4.04-1.61-.55-1.35-1.34-1.72-1.34-1.72-1.09-.74.08-.73.08-.73 1.21.08 1.84 1.22 1.84 1.22 1.07 1.81 2.81 1.29 3.5.99.11-.77.42-1.29.76-1.59-2.67-.3-5.47-1.31-5.47-5.83 0-1.29.47-2.34 1.23-3.17-.12-.3-.54-1.51.12-3.14 0 0 1.01-.32 3.3 1.2.96-.26 1.99-.39 3.02-.4 1.02.01 2.06.14 3.02.4 2.29-1.52 3.3-1.2 3.3-1.2.66 1.63.25 2.84.12 3.14.77.83 1.23 1.88 1.23 3.17 0 4.53-2.81 5.53-5.49 5.82.43.36.81 1.08.81 2.18v3.23c0 .31.22.69.82.57C20.57 21.89 24 17.48 24 12.29 24 5.78 18.63.5 12 .5z" />
              </svg>
              GitHub
            </a>
          </div>

          <p className="mt-7 text-center text-[13px] text-text-secondary">
            Already have an account?{" "}
            <Link href="/login" className="font-medium text-accent hover:underline">
              Sign in
            </Link>
          </p>
        </div>

        <p className="mt-6 text-center text-[11px] text-text-muted">
          By creating an account, you agree to our{" "}
          <a href="#" className="hover:text-text-secondary underline underline-offset-2">Terms</a>
          {" "}and{" "}
          <a href="#" className="hover:text-text-secondary underline underline-offset-2">Privacy Policy</a>
        </p>
      </div>
    </div>
  );
}
`
}

// centeredForgotPasswordPage returns the Linear-style centered forgot-password page.
func centeredForgotPasswordPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { apiClient } from "@/lib/api-client";
import { ForgotPasswordSchema, type ForgotPasswordInput } from "@repo/shared/schemas";
import { ArrowLeft, Loader2, Mail } from "@/lib/icons";

export default function ForgotPasswordPage() {
  const [sent, setSent] = useState(false);
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ForgotPasswordInput>({
    resolver: zodResolver(ForgotPasswordSchema),
  });

  const onSubmit = async (data: ForgotPasswordInput) => {
    setLoading(true);
    setEmail(data.email);
    try {
      await apiClient.post("/api/auth/forgot-password", data);
    } catch {
      // Always show success for security (prevent email enumeration)
    } finally {
      setLoading(false);
      setSent(true);
    }
  };

  const isPending = loading;

  const base = "w-full h-11 rounded-lg border bg-bg-tertiary px-3.5 text-[14px] text-foreground placeholder:text-text-muted transition-colors focus:outline-none focus:ring-2 focus:ring-accent/20";
  const inputClass = ` + "`" + `${base} border-border focus:border-accent` + "`" + `;
  const errorInputClass = ` + "`" + `${base} border-danger focus:border-danger focus:ring-danger/20` + "`" + `;

  return (
    <div className="relative min-h-screen flex items-center justify-center bg-background px-4 py-8">
      <div
        aria-hidden
        className="pointer-events-none absolute inset-x-0 top-0 h-[400px] opacity-60"
        style={{
          background: "radial-gradient(600px at 50% 0%, rgba(108, 92, 231, 0.10), transparent 70%)",
        }}
      />

      <div className="relative w-full max-w-[420px]">
        <div className="rounded-2xl border border-border bg-bg-secondary p-8 sm:p-10 shadow-sm">
          {sent ? (
            <>
              <div className="flex justify-center">
                <div className="h-12 w-12 rounded-full bg-accent/10 flex items-center justify-center">
                  <Mail className="h-5 w-5 text-accent" />
                </div>
              </div>

              <div className="mt-5 text-center">
                <h1 className="text-[22px] font-semibold text-foreground tracking-tight">
                  Check your email
                </h1>
                <p className="mt-1.5 text-[13px] text-text-secondary">
                  If an account exists for <span className="text-foreground">{email}</span>, we sent a reset link.
                </p>
              </div>

              <Link
                href="/login"
                className="mt-7 inline-flex w-full items-center justify-center gap-2 h-11 rounded-lg border border-border bg-bg-tertiary text-[14px] font-medium text-foreground hover:bg-bg-hover transition-colors"
              >
                <ArrowLeft className="h-4 w-4" />
                Back to sign in
              </Link>

              <p className="mt-5 text-center text-[12px] text-text-muted">
                Didn&apos;t get an email?{" "}
                <button
                  onClick={() => setSent(false)}
                  className="font-medium text-accent hover:underline"
                >
                  Try again
                </button>
              </p>
            </>
          ) : (
            <>
              <div className="flex justify-center">
                <div className="h-10 w-10 rounded-xl bg-accent flex items-center justify-center shadow-sm">
                  <span className="text-[17px] font-bold text-white">G</span>
                </div>
              </div>

              <div className="mt-5 text-center">
                <h1 className="text-[22px] font-semibold text-foreground tracking-tight">
                  Reset your password
                </h1>
                <p className="mt-1.5 text-[13px] text-text-secondary">
                  We&apos;ll email you a link to reset your password
                </p>
              </div>

              <form onSubmit={handleSubmit(onSubmit)} className="mt-7 space-y-5">
                <div className="space-y-1.5">
                  <label htmlFor="email" className="block text-[13px] font-medium text-text-secondary">
                    Email
                  </label>
                  <input
                    id="email"
                    type="email"
                    autoComplete="email"
                    placeholder="you@example.com"
                    className={errors.email ? errorInputClass : inputClass}
                    {...register("email")}
                  />
                  {errors.email && (
                    <p className="text-[12px] font-medium text-danger">{errors.email.message}</p>
                  )}
                </div>

                <button
                  type="submit"
                  disabled={isPending}
                  className="w-full h-11 rounded-lg bg-accent text-[14px] font-medium text-white hover:bg-accent-hover transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                >
                  {isPending && <Loader2 className="h-4 w-4 animate-spin" />}
                  {isPending ? "Sending..." : "Send reset link"}
                </button>
              </form>

              <p className="mt-7 text-center text-[13px] text-text-secondary">
                Remembered it?{" "}
                <Link href="/login" className="font-medium text-accent hover:underline">
                  Back to sign in
                </Link>
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
`
}
