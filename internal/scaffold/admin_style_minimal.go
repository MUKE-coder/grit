package scaffold

import "fmt"

// minimalLoginPage returns a Vercel/Linear-inspired minimal login page.
// Full-screen bg-primary, form centered, no cards or decorative backgrounds.
func minimalLoginPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { Eye, EyeOff } from "@/lib/icons";
import { useLogin } from "@/hooks/use-auth";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginSchema, type LoginInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";
const errorInputClass = "w-full rounded-lg border border-danger bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors";

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  const { mutate: login, isPending, error: serverError } = useLogin();

  const { register, handleSubmit, formState: { errors } } = useForm<LoginInput>({
    resolver: zodResolver(LoginSchema),
  });

  const onSubmit = (data: LoginInput) => {
    login(data);
  };

  return (
    <div className="flex min-h-screen flex-col bg-background">
      {/* Subtle logo top-left */}
      <div className="px-8 pt-8">
        <span className="text-sm font-medium text-text-muted tracking-wide">grit</span>
      </div>

      {/* Centered form */}
      <div className="flex flex-1 items-center justify-center px-6">
        <div className="w-full max-w-sm space-y-8">
          <div>
            <h1 className="text-xl font-semibold text-foreground">Sign in</h1>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {serverError && (
              <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
                {(serverError as unknown as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || "Invalid credentials"}
              </div>
            )}

            <div className="space-y-1.5">
              <label htmlFor="email" className="block text-sm text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                {...register("email")}
                className={errors.email ? errorInputClass : inputClass}
                placeholder="you@example.com"
                autoFocus
              />
              {errors.email && <p className="text-sm text-danger">{errors.email.message}</p>}
            </div>

            <div className="space-y-1.5">
              <label htmlFor="password" className="block text-sm text-text-secondary">
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  {...register("password")}
                  className={errors.password ? errorInputClass + " pr-12" : inputClass + " pr-12"}
                  placeholder="Enter your password"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.password && <p className="text-sm text-danger">{errors.password.message}</p>}
            </div>

            <div className="flex justify-end">
              <Link href="/forgot-password" className="text-sm text-text-muted hover:text-foreground transition-colors">
                Forgot password?
              </Link>
            </div>

            <button
              type="submit"
              disabled={isPending}
              className="w-full rounded-lg bg-accent py-3 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
            >
              {isPending ? "Signing in..." : "Sign in"}
            </button>
          </form>

          <p className="text-center text-sm text-text-muted">
            Don&apos;t have an account?{" "}
            <Link href="/sign-up" className="text-foreground hover:text-accent transition-colors">
              Sign up
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
`
}

// minimalSignUpPage returns a Vercel/Linear-inspired minimal sign-up page.
// Same clean layout — no cards, generous whitespace, centered form.
func minimalSignUpPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { Eye, EyeOff } from "@/lib/icons";
import { useRegister } from "@/hooks/use-auth";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { RegisterSchema, type RegisterInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";
const errorInputClass = "w-full rounded-lg border border-danger bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors";

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

  return (
    <div className="flex min-h-screen flex-col bg-background">
      {/* Subtle logo top-left */}
      <div className="px-8 pt-8">
        <span className="text-sm font-medium text-text-muted tracking-wide">grit</span>
      </div>

      {/* Centered form */}
      <div className="flex flex-1 items-center justify-center px-6">
        <div className="w-full max-w-sm space-y-8">
          <div>
            <h1 className="text-xl font-semibold text-foreground">Create an account</h1>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            {serverError && (
              <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
                {(serverError as unknown as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || "Registration failed"}
              </div>
            )}

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <label htmlFor="firstName" className="block text-sm text-text-secondary">
                  First name
                </label>
                <input
                  id="firstName"
                  type="text"
                  {...register("firstName")}
                  className={errors.firstName ? errorInputClass : inputClass}
                  placeholder="John"
                  autoFocus
                />
                {errors.firstName && <p className="text-sm text-danger">{errors.firstName.message}</p>}
              </div>
              <div className="space-y-1.5">
                <label htmlFor="lastName" className="block text-sm text-text-secondary">
                  Last name
                </label>
                <input
                  id="lastName"
                  type="text"
                  {...register("lastName")}
                  className={errors.lastName ? errorInputClass : inputClass}
                  placeholder="Doe"
                />
                {errors.lastName && <p className="text-sm text-danger">{errors.lastName.message}</p>}
              </div>
            </div>

            <div className="space-y-1.5">
              <label htmlFor="email" className="block text-sm text-text-secondary">
                Email
              </label>
              <input
                id="email"
                type="email"
                {...register("email")}
                className={errors.email ? errorInputClass : inputClass}
                placeholder="you@example.com"
              />
              {errors.email && <p className="text-sm text-danger">{errors.email.message}</p>}
            </div>

            <div className="space-y-1.5">
              <label htmlFor="password" className="block text-sm text-text-secondary">
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  {...register("password")}
                  className={errors.password ? errorInputClass + " pr-12" : inputClass + " pr-12"}
                  placeholder="Min. 8 characters"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.password && <p className="text-sm text-danger">{errors.password.message}</p>}
            </div>

            <div className="space-y-1.5">
              <label htmlFor="confirmPassword" className="block text-sm text-text-secondary">
                Confirm password
              </label>
              <div className="relative">
                <input
                  id="confirmPassword"
                  type={showConfirmPassword ? "text" : "password"}
                  {...register("confirmPassword")}
                  className={errors.confirmPassword ? errorInputClass + " pr-12" : inputClass + " pr-12"}
                  placeholder="Repeat your password"
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-text-secondary transition-colors"
                >
                  {showConfirmPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.confirmPassword && <p className="text-sm text-danger">{errors.confirmPassword.message}</p>}
            </div>

            <button
              type="submit"
              disabled={isPending}
              className="w-full rounded-lg bg-accent py-3 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
            >
              {isPending ? "Creating account..." : "Create account"}
            </button>
          </form>

          <p className="text-center text-sm text-text-muted">
            Already have an account?{" "}
            <Link href="/login" className="text-foreground hover:text-accent transition-colors">
              Sign in
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
`
}

// minimalForgotPasswordPage returns a Vercel/Linear-inspired minimal forgot password page.
// Same clean layout with sent state showing a success message.
func minimalForgotPasswordPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { apiClient } from "@/lib/api-client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ForgotPasswordSchema, type ForgotPasswordInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-border bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors";
const errorInputClass = "w-full rounded-lg border border-danger bg-bg-tertiary px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors";

export default function ForgotPasswordPage() {
  const [sent, setSent] = useState(false);
  const [loading, setLoading] = useState(false);

  const { register, handleSubmit, formState: { errors } } = useForm<ForgotPasswordInput>({
    resolver: zodResolver(ForgotPasswordSchema),
  });

  const onSubmit = async (data: ForgotPasswordInput) => {
    setLoading(true);
    try {
      await apiClient.post("/api/auth/forgot-password", data);
      setSent(true);
    } catch {
      setSent(true); // Always show success for security
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen flex-col bg-background">
      {/* Subtle logo top-left */}
      <div className="px-8 pt-8">
        <span className="text-sm font-medium text-text-muted tracking-wide">grit</span>
      </div>

      {/* Centered form */}
      <div className="flex flex-1 items-center justify-center px-6">
        <div className="w-full max-w-sm space-y-8">
          <div>
            <h1 className="text-xl font-semibold text-foreground">Reset password</h1>
            <p className="mt-2 text-sm text-text-muted">
              Enter your email and we&apos;ll send you a reset link.
            </p>
          </div>

          {sent ? (
            <div className="space-y-6">
              <div className="space-y-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-full bg-success/10">
                  <svg className="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <p className="text-sm text-text-secondary">
                  If an account with that email exists, we&apos;ve sent a password reset link. Check your inbox.
                </p>
              </div>
              <Link
                href="/login"
                className="inline-block text-sm text-foreground hover:text-accent transition-colors"
              >
                &larr; Back to sign in
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="space-y-1.5">
                <label htmlFor="email" className="block text-sm text-text-secondary">
                  Email
                </label>
                <input
                  id="email"
                  type="email"
                  {...register("email")}
                  className={errors.email ? errorInputClass : inputClass}
                  placeholder="you@example.com"
                  autoFocus
                />
                {errors.email && <p className="text-sm text-danger">{errors.email.message}</p>}
              </div>

              <button
                type="submit"
                disabled={loading}
                className="w-full rounded-lg bg-accent py-3 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
              >
                {loading ? "Sending..." : "Send reset link"}
              </button>
            </form>
          )}

          {!sent && (
            <p className="text-center text-sm text-text-muted">
              <Link href="/login" className="text-foreground hover:text-accent transition-colors">
                &larr; Back to sign in
              </Link>
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
`
}

// minimalDashboardPage returns a compact, data-dense dashboard inspired by Vercel/Linear.
// No welcome banner — just a small greeting, stats row, compact resource table, and inline quick links.
func minimalDashboardPage() string {
	return fmt.Sprintf(`"use client";

import { useMe } from "@/hooks/use-auth";
import { resources } from "@/resources";
import { StatsCard } from "@/components/widgets/stats-card";
import { WidgetGrid } from "@/components/widgets/widget-grid";
import { getIcon } from "@/lib/icons";

export default function AdminDashboard() {
  const { data: user } = useMe();
  const allWidgets = resources.flatMap((r) => r.dashboard?.widgets ?? []);

  const greeting = () => {
    const hour = new Date().getHours();
    if (hour < 12) return "Good morning";
    if (hour < 17) return "Good afternoon";
    return "Good evening";
  };

  return (
    <div className="space-y-6">
      {/* Compact header */}
      <div>
        <h1 className="text-lg font-semibold text-foreground">Dashboard</h1>
        <p className="text-sm text-text-muted">
          {greeting()}, {user?.first_name || "Admin"}
        </p>
      </div>

      {/* Stats row */}
      {allWidgets.length > 0 ? (
        <WidgetGrid widgets={allWidgets} />
      ) : (
        <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
          <StatsCard label="Resources" value={String(resources.length)} icon="Layers" color="accent" />
          <StatsCard label="Total Resources" value="—" icon="Database" color="success" />
        </div>
      )}

      {/* Resources as compact table */}
      <div>
        <h2 className="text-sm font-medium text-text-secondary mb-3">Resources</h2>
        <div className="rounded-lg border border-border bg-bg-secondary overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border text-left">
                <th className="px-4 py-2.5 text-xs font-medium text-text-muted uppercase tracking-wider">Name</th>
                <th className="px-4 py-2.5 text-xs font-medium text-text-muted uppercase tracking-wider">Endpoint</th>
                <th className="px-4 py-2.5 text-right text-xs font-medium text-text-muted uppercase tracking-wider">Action</th>
              </tr>
            </thead>
            <tbody>
              {resources.map((r) => {
                const Icon = getIcon(r.icon);
                return (
                  <tr key={r.slug} className="border-b border-border last:border-0 hover:bg-bg-hover/50 transition-colors">
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-3">
                        <Icon className="h-4 w-4 text-text-muted" />
                        <span className="font-medium text-foreground">{r.label?.plural ?? r.name}</span>
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <span className="font-mono text-xs text-text-muted">{r.endpoint}</span>
                    </td>
                    <td className="px-4 py-3 text-right">
                      <a
                        href={%s/resources/${r.slug}%s}
                        className="text-xs font-medium text-accent hover:text-accent-hover transition-colors"
                      >
                        Manage
                      </a>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>

      {/* Quick links as inline pills */}
      <div>
        <h2 className="text-sm font-medium text-text-secondary mb-3">Quick links</h2>
        <div className="flex flex-wrap gap-2">
          <a
            href="http://localhost:8080/studio"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1.5 rounded-full border border-border bg-bg-secondary px-3.5 py-1.5 text-xs font-medium text-text-secondary hover:text-foreground hover:border-accent/40 transition-colors"
          >
            <span className="text-info">DB</span>
            GORM Studio
          </a>
          <a
            href="http://localhost:8080/api/health"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-1.5 rounded-full border border-border bg-bg-secondary px-3.5 py-1.5 text-xs font-medium text-text-secondary hover:text-foreground hover:border-accent/40 transition-colors"
          >
            <span className="text-success">OK</span>
            API Health
          </a>
          <a
            href="/system/jobs"
            className="inline-flex items-center gap-1.5 rounded-full border border-border bg-bg-secondary px-3.5 py-1.5 text-xs font-medium text-text-secondary hover:text-foreground hover:border-accent/40 transition-colors"
          >
            <span className="text-warning">Q</span>
            Job Queue
          </a>
          <a
            href="/system/files"
            className="inline-flex items-center gap-1.5 rounded-full border border-border bg-bg-secondary px-3.5 py-1.5 text-xs font-medium text-text-secondary hover:text-foreground hover:border-accent/40 transition-colors"
          >
            <span className="text-accent">S3</span>
            File Storage
          </a>
          <a
            href="/system/cron"
            className="inline-flex items-center gap-1.5 rounded-full border border-border bg-bg-secondary px-3.5 py-1.5 text-xs font-medium text-text-secondary hover:text-foreground hover:border-accent/40 transition-colors"
          >
            <span className="text-info">CR</span>
            Cron Jobs
          </a>
          <a
            href="/system/mail"
            className="inline-flex items-center gap-1.5 rounded-full border border-border bg-bg-secondary px-3.5 py-1.5 text-xs font-medium text-text-secondary hover:text-foreground hover:border-accent/40 transition-colors"
          >
            <span className="text-danger">ML</span>
            Mail Preview
          </a>
        </div>
      </div>
    </div>
  );
}
`, "`", "`")
}
