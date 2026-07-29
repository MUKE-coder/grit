"use client";

import { useState } from "react";
import { ArrowLeft, Mail, CheckCircle, Loader2 } from "lucide-react";

export default function ForgotPasswordCard01() {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [sent, setSent] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    await new Promise((r) => setTimeout(r, 1000));
    setLoading(false);
    setSent(true);
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <div className="w-full max-w-sm">
        {/* Logo */}
        <div className="flex justify-center mb-8">
          <div className="w-12 h-12 rounded-2xl bg-accent flex items-center justify-center">
            <span className="text-white font-black text-xl tracking-tight">G</span>
          </div>
        </div>

        <div className="rounded-2xl border border-border bg-bg-elevated p-8">
          {sent ? (
            <div className="text-center flex flex-col items-center gap-4">
              <div className="w-14 h-14 rounded-full bg-bg-secondary border border-border flex items-center justify-center">
                <CheckCircle size={28} className="text-success" />
              </div>
              <div>
                <h2 className="text-xl font-bold text-foreground tracking-tight">
                  Check your email
                </h2>
                <p className="mt-2 text-text-secondary text-sm leading-relaxed">
                  We sent a password reset link to{" "}
                  <span className="text-foreground font-medium">{email}</span>.
                  It expires in 15 minutes.
                </p>
              </div>
              <div className="flex flex-col gap-3 w-full mt-2">
                <button
                  onClick={() => { setSent(false); setEmail(""); }}
                  className="w-full py-2.5 rounded-lg border border-border bg-bg-secondary hover:bg-bg-hover text-foreground font-semibold text-sm transition-colors"
                >
                  Send again
                </button>
                <a
                  href="/login"
                  className="w-full inline-flex items-center justify-center gap-2 py-2.5 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors"
                >
                  <ArrowLeft size={15} />
                  Back to login
                </a>
              </div>
            </div>
          ) : (
            <>
              <div className="mb-7">
                <h1 className="text-2xl font-bold text-foreground tracking-tight">
                  Forgot your password?
                </h1>
                <p className="mt-2 text-text-secondary text-sm leading-relaxed">
                  No worries. Enter your email address and we&rsquo;ll send you a reset link.
                </p>
              </div>

              <form onSubmit={handleSubmit} className="flex flex-col gap-5">
                <div className="flex flex-col gap-1.5">
                  <label className="text-sm font-medium text-foreground">Email address</label>
                  <div className="relative">
                    <Mail size={15} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
                    <input
                      type="email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="you@company.com"
                      required
                      className="w-full pl-9 pr-3 py-2.5 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
                    />
                  </div>
                </div>

                <button
                  type="submit"
                  disabled={loading}
                  className="w-full py-2.5 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors duration-150 disabled:opacity-60 flex items-center justify-center gap-2"
                >
                  {loading ? <Loader2 size={16} className="animate-spin" /> : "Send reset link"}
                </button>
              </form>

              <div className="mt-6 text-center">
                <a
                  href="/login"
                  className="inline-flex items-center gap-1.5 text-sm text-text-secondary hover:text-foreground transition-colors"
                >
                  <ArrowLeft size={14} />
                  Back to login
                </a>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
