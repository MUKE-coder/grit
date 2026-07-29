"use client";

import { useState } from "react";
import { Check, Eye, EyeOff, Loader2 } from "lucide-react";

const features = [
  "Generate full-stack apps in seconds",
  "Go API + Next.js admin panel",
  "Type-safe from database to UI",
  "Built-in auth, jobs, and storage",
  "Deploy anywhere with Docker",
];

export default function SignupCard01() {
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [agreed, setAgreed] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    await new Promise((r) => setTimeout(r, 1000));
    setLoading(false);
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-4xl rounded-2xl border border-border overflow-hidden grid grid-cols-1 md:grid-cols-2">
        {/* Left: Feature list */}
        <div className="bg-bg-secondary p-10 flex flex-col justify-center">
          <div className="w-10 h-10 rounded-xl bg-accent flex items-center justify-center mb-8">
            <span className="text-white font-black text-xl">G</span>
          </div>
          <h2 className="text-2xl font-bold text-foreground tracking-tight mb-2">
            Build your next project with Grit
          </h2>
          <p className="text-text-secondary text-sm mb-8 leading-relaxed">
            The full-stack framework for developers who value both speed and control.
          </p>
          <ul className="flex flex-col gap-4">
            {features.map((f, i) => (
              <li key={i} className="flex items-center gap-3">
                <div className="w-5 h-5 rounded-full bg-accent flex items-center justify-center shrink-0">
                  <Check size={11} className="text-white" />
                </div>
                <span className="text-text-secondary text-sm">{f}</span>
              </li>
            ))}
          </ul>
        </div>

        {/* Right: Form */}
        <div className="bg-bg-elevated p-10 flex flex-col justify-center">
          <div className="mb-8">
            <h1 className="text-2xl font-bold text-foreground tracking-tight">
              Create your account
            </h1>
            <p className="mt-1 text-sm text-text-secondary">
              Already have an account?{" "}
              <a href="/login" className="text-accent hover:underline font-medium">
                Sign in
              </a>
            </p>
          </div>

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium text-foreground">Full name</label>
              <input
                type="text"
                placeholder="Alex Johnson"
                required
                className="px-3 py-2.5 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium text-foreground">Email</label>
              <input
                type="email"
                placeholder="you@company.com"
                required
                className="px-3 py-2.5 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-sm font-medium text-foreground">Password</label>
              <div className="relative">
                <input
                  type={showPassword ? "text" : "password"}
                  placeholder="••••••••"
                  required
                  minLength={8}
                  className="w-full px-3 py-2.5 pr-10 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-foreground"
                >
                  {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
              <p className="text-text-muted text-xs">Minimum 8 characters</p>
            </div>

            <div className="flex items-start gap-2 mt-1">
              <input
                id="terms"
                type="checkbox"
                checked={agreed}
                onChange={(e) => setAgreed(e.target.checked)}
                required
                className="mt-0.5 w-4 h-4 rounded border-border accent-accent"
              />
              <label htmlFor="terms" className="text-xs text-text-secondary cursor-pointer leading-relaxed">
                I agree to the{" "}
                <a href="/terms" className="text-accent hover:underline">Terms of Service</a>
                {" "}and{" "}
                <a href="/privacy" className="text-accent hover:underline">Privacy Policy</a>
              </label>
            </div>

            <button
              type="submit"
              disabled={loading || !agreed}
              className="w-full py-2.5 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors duration-150 disabled:opacity-60 flex items-center justify-center gap-2 mt-2"
            >
              {loading ? <Loader2 size={16} className="animate-spin" /> : "Create account"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
