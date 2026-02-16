package scaffold

import "fmt"

// glassLoginPage returns the frosted-glass login page for the "glass" style.
func glassLoginPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { Eye, EyeOff } from "@/lib/icons";
import { useLogin } from "@/hooks/use-auth";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginSchema, type LoginInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-white/10 bg-white/5 px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors backdrop-blur-sm";
const errorInputClass = "w-full rounded-lg border border-danger/50 bg-white/5 px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors backdrop-blur-sm";

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
    <div className="relative min-h-screen flex items-center justify-center overflow-hidden">
      {/* Background layers */}
      <div className="absolute inset-0 bg-gradient-to-br from-accent/20 via-background to-info/10" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-accent/15 via-transparent to-transparent" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_bottom_left,_var(--tw-gradient-stops))] from-info/10 via-transparent to-transparent" />

      {/* Floating decorative orbs */}
      <div className="pointer-events-none absolute top-1/4 left-1/4 h-64 w-64 rounded-full bg-accent/10 blur-[100px]" />
      <div className="pointer-events-none absolute bottom-1/4 right-1/4 h-48 w-48 rounded-full bg-info/10 blur-[80px]" />

      {/* Frosted glass card */}
      <div className="relative z-10 w-full max-w-md mx-4">
        <div className="rounded-2xl border border-white/10 bg-bg-secondary/60 backdrop-blur-2xl p-8 shadow-2xl shadow-black/30">
          {/* Logo */}
          <div className="mb-8 text-center">
            <div className="inline-flex items-center gap-2">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-accent/20 border border-accent/30 backdrop-blur-sm">
                <span className="text-accent font-mono font-bold text-lg">G</span>
              </div>
              <span className="text-3xl font-bold text-foreground">rit</span>
              <span className="rounded-md bg-white/10 backdrop-blur-sm px-2 py-0.5 text-xs font-medium text-accent">
                Admin
              </span>
            </div>
          </div>

          {/* Heading */}
          <div className="mb-6 text-center">
            <h2 className="text-2xl font-bold text-foreground">Welcome back</h2>
            <p className="mt-2 text-foreground/60">Sign in to your admin account</p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            {serverError && (
              <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
                {(serverError as unknown as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || "Invalid credentials"}
              </div>
            )}

            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-foreground/80">
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

            <div className="space-y-2">
              <label htmlFor="password" className="block text-sm font-medium text-foreground/80">
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
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-foreground/80 transition-colors"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.password && <p className="text-sm text-danger">{errors.password.message}</p>}
            </div>

            <div className="flex items-center justify-between">
              <label className="flex items-center gap-2 text-sm text-foreground/60 cursor-pointer">
                <input type="checkbox" className="h-4 w-4 rounded border-white/10 bg-white/5 accent-accent" />
                Remember me
              </label>
              <Link href="/forgot-password" className="text-sm text-accent hover:text-accent-hover transition-colors">
                Forgot password?
              </Link>
            </div>

            <button
              type="submit"
              disabled={isPending}
              className="w-full rounded-lg bg-accent/90 backdrop-blur-sm py-3 font-medium text-white hover:bg-accent disabled:opacity-50 transition-colors shadow-lg shadow-accent/20"
            >
              {isPending ? "Signing in..." : "Sign In"}
            </button>
          </form>

          {/* Footer link */}
          <p className="mt-6 text-center text-sm text-foreground/60">
            Don&apos;t have an account?{" "}
            <Link href="/sign-up" className="text-accent hover:text-accent-hover font-medium transition-colors">
              Create one
            </Link>
          </p>
        </div>

        {/* Below-card branding */}
        <p className="mt-6 text-center text-xs text-foreground/40">
          Built with Grit — Go + React framework
        </p>
      </div>
    </div>
  );
}
`
}

// glassSignUpPage returns the frosted-glass sign-up page for the "glass" style.
func glassSignUpPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { Eye, EyeOff } from "@/lib/icons";
import { useRegister } from "@/hooks/use-auth";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { RegisterSchema, type RegisterInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-white/10 bg-white/5 px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors backdrop-blur-sm";
const errorInputClass = "w-full rounded-lg border border-danger/50 bg-white/5 px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors backdrop-blur-sm";

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
    <div className="relative min-h-screen flex items-center justify-center overflow-hidden">
      {/* Background layers */}
      <div className="absolute inset-0 bg-gradient-to-br from-accent/20 via-background to-info/10" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-accent/15 via-transparent to-transparent" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_bottom_left,_var(--tw-gradient-stops))] from-info/10 via-transparent to-transparent" />

      {/* Floating decorative orbs */}
      <div className="pointer-events-none absolute top-1/4 left-1/4 h-64 w-64 rounded-full bg-accent/10 blur-[100px]" />
      <div className="pointer-events-none absolute bottom-1/4 right-1/4 h-48 w-48 rounded-full bg-info/10 blur-[80px]" />

      {/* Frosted glass card */}
      <div className="relative z-10 w-full max-w-md mx-4 py-12">
        <div className="rounded-2xl border border-white/10 bg-bg-secondary/60 backdrop-blur-2xl p-8 shadow-2xl shadow-black/30">
          {/* Logo */}
          <div className="mb-8 text-center">
            <div className="inline-flex items-center gap-2">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-accent/20 border border-accent/30 backdrop-blur-sm">
                <span className="text-accent font-mono font-bold text-lg">G</span>
              </div>
              <span className="text-3xl font-bold text-foreground">rit</span>
              <span className="rounded-md bg-white/10 backdrop-blur-sm px-2 py-0.5 text-xs font-medium text-accent">
                Admin
              </span>
            </div>
          </div>

          {/* Heading */}
          <div className="mb-6 text-center">
            <h2 className="text-2xl font-bold text-foreground">Create account</h2>
            <p className="mt-2 text-foreground/60">Sign up to get started</p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            {serverError && (
              <div className="rounded-lg bg-danger/10 border border-danger/20 px-4 py-3 text-sm text-danger">
                {(serverError as unknown as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message || "Registration failed"}
              </div>
            )}

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <label htmlFor="firstName" className="block text-sm font-medium text-foreground/80">
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
              <div className="space-y-2">
                <label htmlFor="lastName" className="block text-sm font-medium text-foreground/80">
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

            <div className="space-y-2">
              <label htmlFor="email" className="block text-sm font-medium text-foreground/80">
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

            <div className="space-y-2">
              <label htmlFor="password" className="block text-sm font-medium text-foreground/80">
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
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-foreground/80 transition-colors"
                >
                  {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.password && <p className="text-sm text-danger">{errors.password.message}</p>}
            </div>

            <div className="space-y-2">
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-foreground/80">
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
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-foreground/80 transition-colors"
                >
                  {showConfirmPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
                </button>
              </div>
              {errors.confirmPassword && <p className="text-sm text-danger">{errors.confirmPassword.message}</p>}
            </div>

            <button
              type="submit"
              disabled={isPending}
              className="w-full rounded-lg bg-accent/90 backdrop-blur-sm py-3 font-medium text-white hover:bg-accent disabled:opacity-50 transition-colors shadow-lg shadow-accent/20"
            >
              {isPending ? "Creating account..." : "Create Account"}
            </button>
          </form>

          {/* Footer link */}
          <p className="mt-6 text-center text-sm text-foreground/60">
            Already have an account?{" "}
            <Link href="/login" className="text-accent hover:text-accent-hover font-medium transition-colors">
              Sign in
            </Link>
          </p>
        </div>

        {/* Below-card branding */}
        <p className="mt-6 text-center text-xs text-foreground/40">
          Built with Grit — Go + React framework
        </p>
      </div>
    </div>
  );
}
`
}

// glassForgotPasswordPage returns the frosted-glass forgot password page for the "glass" style.
func glassForgotPasswordPage() string {
	return `"use client";

import { useState } from "react";
import Link from "next/link";
import { apiClient } from "@/lib/api-client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ForgotPasswordSchema, type ForgotPasswordInput } from "@repo/shared/schemas";

const inputClass = "w-full rounded-lg border border-white/10 bg-white/5 px-4 py-3 text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent transition-colors backdrop-blur-sm";
const errorInputClass = "w-full rounded-lg border border-danger/50 bg-white/5 px-4 py-3 text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger transition-colors backdrop-blur-sm";

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
    <div className="relative min-h-screen flex items-center justify-center overflow-hidden">
      {/* Background layers */}
      <div className="absolute inset-0 bg-gradient-to-br from-accent/20 via-background to-info/10" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-accent/15 via-transparent to-transparent" />
      <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_bottom_left,_var(--tw-gradient-stops))] from-info/10 via-transparent to-transparent" />

      {/* Floating decorative orbs */}
      <div className="pointer-events-none absolute top-1/4 left-1/4 h-64 w-64 rounded-full bg-accent/10 blur-[100px]" />
      <div className="pointer-events-none absolute bottom-1/4 right-1/4 h-48 w-48 rounded-full bg-info/10 blur-[80px]" />

      {/* Frosted glass card */}
      <div className="relative z-10 w-full max-w-md mx-4">
        <div className="rounded-2xl border border-white/10 bg-bg-secondary/60 backdrop-blur-2xl p-8 shadow-2xl shadow-black/30">
          {/* Logo */}
          <div className="mb-8 text-center">
            <div className="inline-flex items-center gap-2">
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-accent/20 border border-accent/30 backdrop-blur-sm">
                <span className="text-accent font-mono font-bold text-lg">G</span>
              </div>
              <span className="text-3xl font-bold text-foreground">rit</span>
              <span className="rounded-md bg-white/10 backdrop-blur-sm px-2 py-0.5 text-xs font-medium text-accent">
                Admin
              </span>
            </div>
          </div>

          {/* Heading */}
          <div className="mb-6 text-center">
            <h2 className="text-2xl font-bold text-foreground">Forgot password?</h2>
            <p className="mt-2 text-foreground/60">No worries, we&apos;ll send you reset instructions.</p>
          </div>

          {sent ? (
            <div className="rounded-xl bg-white/5 border border-white/10 backdrop-blur-sm p-8 text-center space-y-4">
              <div className="flex justify-center">
                <div className="h-12 w-12 rounded-full bg-success/10 border border-success/20 flex items-center justify-center">
                  <svg className="h-6 w-6 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </div>
              </div>
              <h3 className="text-lg font-medium text-foreground">Check your email</h3>
              <p className="text-foreground/60 text-sm">
                If an account with that email exists, we&apos;ve sent a password reset link.
              </p>
              <Link
                href="/login"
                className="inline-block text-accent hover:text-accent-hover font-medium text-sm transition-colors"
              >
                Back to login
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
              <div className="space-y-2">
                <label htmlFor="email" className="block text-sm font-medium text-foreground/80">
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
                className="w-full rounded-lg bg-accent/90 backdrop-blur-sm py-3 font-medium text-white hover:bg-accent disabled:opacity-50 transition-colors shadow-lg shadow-accent/20"
              >
                {loading ? "Sending..." : "Send Reset Link"}
              </button>
            </form>
          )}

          {/* Footer link */}
          <p className="mt-6 text-center text-sm text-foreground/60">
            <Link href="/login" className="text-accent hover:text-accent-hover font-medium transition-colors">
              Back to login
            </Link>
          </p>
        </div>

        {/* Below-card branding */}
        <p className="mt-6 text-center text-xs text-foreground/40">
          Built with Grit — Go + React framework
        </p>
      </div>
    </div>
  );
}
`
}

// glassDashboardPage returns the hero-banner dashboard page for the "glass" style.
func glassDashboardPage() string {
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

  const today = new Date().toLocaleDateString("en-US", {
    weekday: "long",
    year: "numeric",
    month: "long",
    day: "numeric",
  });

  return (
    <div className="space-y-6">
      {/* Hero banner */}
      <div className="relative rounded-2xl overflow-hidden">
        {/* Hero background layers */}
        <div className="absolute inset-0 bg-gradient-to-r from-accent/20 via-info/10 to-bg-secondary" />
        <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top_right,_var(--tw-gradient-stops))] from-accent/15 via-transparent to-transparent" />
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_bottom_left,_var(--tw-gradient-stops))] from-info/10 via-transparent to-transparent" />

        {/* Decorative mesh pattern */}
        <div className="absolute inset-0 opacity-[0.03]" style={{ backgroundImage: "url(\"data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='1'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E\")" }} />

        {/* Floating orbs */}
        <div className="pointer-events-none absolute -top-10 -right-10 h-40 w-40 rounded-full bg-accent/15 blur-[60px]" />
        <div className="pointer-events-none absolute -bottom-10 -left-10 h-32 w-32 rounded-full bg-info/15 blur-[50px]" />

        {/* Hero content */}
        <div className="relative px-8 py-10 sm:px-10 sm:py-12">
          <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-6">
            <div>
              <p className="text-foreground/50 text-sm mb-1">{today}</p>
              <h1 className="text-3xl sm:text-4xl font-bold text-foreground">
                {greeting()}, {user?.first_name || "Admin"}
              </h1>
              <p className="text-foreground/60 mt-2 text-lg">
                Here&apos;s an overview of your application.
              </p>
            </div>

            {/* Inline stats in hero */}
            <div className="flex gap-6">
              <div className="rounded-xl border border-white/10 bg-white/5 backdrop-blur-xl px-5 py-3 text-center">
                <p className="text-2xl font-bold text-foreground">{resources.length}</p>
                <p className="text-xs text-foreground/50 mt-0.5">Resources</p>
              </div>
              <div className="rounded-xl border border-white/10 bg-white/5 backdrop-blur-xl px-5 py-3 text-center">
                <p className="text-2xl font-bold text-success">OK</p>
                <p className="text-xs text-foreground/50 mt-0.5">Status</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Stats widgets */}
      {allWidgets.length > 0 ? (
        <WidgetGrid widgets={allWidgets} />
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <StatsCard label="Total Resources" value="—" icon="Database" color="accent" />
          <StatsCard label="Registered" value={String(resources.length)} icon="Layers" color="success" />
        </div>
      )}

      {/* Resources as glass cards */}
      <div>
        <h2 className="text-lg font-semibold text-foreground mb-4">Resources</h2>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {resources.map((r) => {
            const Icon = getIcon(r.icon);
            return (
              <a
                key={r.slug}
                href={%s/resources/${r.slug}%s}
                className="group relative rounded-2xl border border-white/5 bg-bg-secondary/60 backdrop-blur-xl p-6 hover:border-accent/30 hover:bg-bg-secondary/80 transition-all overflow-hidden"
              >
                {/* Card glow on hover */}
                <div className="pointer-events-none absolute -top-10 -right-10 h-24 w-24 rounded-full bg-accent/0 group-hover:bg-accent/10 blur-2xl transition-all duration-500" />

                <div className="relative flex items-start gap-4">
                  <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-accent/10 border border-accent/20 group-hover:bg-accent/20 group-hover:border-accent/30 transition-all">
                    <Icon className="h-6 w-6 text-accent" />
                  </div>
                  <div>
                    <h3 className="font-semibold text-foreground group-hover:text-accent transition-colors">
                      {r.label?.plural ?? r.name}
                    </h3>
                    <p className="text-sm text-foreground/50 mt-1">
                      Manage {(r.label?.plural ?? r.slug).toLowerCase()}
                    </p>
                  </div>
                </div>
              </a>
            );
          })}
        </div>
      </div>

      {/* Quick Links — glass panel */}
      <div className="rounded-2xl border border-white/5 bg-bg-secondary/60 backdrop-blur-xl p-6">
        <h2 className="text-lg font-semibold text-foreground mb-4">Quick Links</h2>
        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
          <a
            href="http://localhost:8080/studio"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-3 rounded-xl border border-white/5 bg-white/5 backdrop-blur-sm px-4 py-3 hover:border-accent/20 hover:bg-white/10 transition-all group"
          >
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-info/10 border border-info/20">
              <span className="text-info text-sm font-bold">DB</span>
            </div>
            <div>
              <p className="text-sm font-medium text-foreground group-hover:text-accent transition-colors">GORM Studio</p>
              <p className="text-xs text-foreground/40">Browse database</p>
            </div>
          </a>
          <a
            href="http://localhost:8080/api/health"
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-3 rounded-xl border border-white/5 bg-white/5 backdrop-blur-sm px-4 py-3 hover:border-accent/20 hover:bg-white/10 transition-all group"
          >
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-success/10 border border-success/20">
              <span className="text-success text-sm font-bold">OK</span>
            </div>
            <div>
              <p className="text-sm font-medium text-foreground group-hover:text-accent transition-colors">API Health</p>
              <p className="text-xs text-foreground/40">Check status</p>
            </div>
          </a>
          <a
            href="/system/jobs"
            className="flex items-center gap-3 rounded-xl border border-white/5 bg-white/5 backdrop-blur-sm px-4 py-3 hover:border-accent/20 hover:bg-white/10 transition-all group"
          >
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-warning/10 border border-warning/20">
              <span className="text-warning text-sm font-bold">Q</span>
            </div>
            <div>
              <p className="text-sm font-medium text-foreground group-hover:text-accent transition-colors">Job Queue</p>
              <p className="text-xs text-foreground/40">Background jobs</p>
            </div>
          </a>
          <a
            href="/system/files"
            className="flex items-center gap-3 rounded-xl border border-white/5 bg-white/5 backdrop-blur-sm px-4 py-3 hover:border-accent/20 hover:bg-white/10 transition-all group"
          >
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent/10 border border-accent/20">
              <span className="text-accent text-sm font-bold">S3</span>
            </div>
            <div>
              <p className="text-sm font-medium text-foreground group-hover:text-accent transition-colors">File Storage</p>
              <p className="text-xs text-foreground/40">Manage uploads</p>
            </div>
          </a>
        </div>
      </div>
    </div>
  );
}
`, "`", "`")
}
