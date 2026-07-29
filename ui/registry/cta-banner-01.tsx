"use client";

import { useState } from "react";
import { ArrowRight, Loader2 } from "lucide-react";

interface CtaBanner01Props {
  title?: string;
  subtitle?: string;
  buttonText?: string;
  buttonHref?: string;
}

export default function CtaBanner01({
  title = "Start building today",
  subtitle = "Join thousands of developers shipping faster with Grit.",
  buttonText = "Get started free",
  buttonHref = "/get-started",
}: CtaBanner01Props) {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;
    setLoading(true);
    await new Promise((r) => setTimeout(r, 900));
    setLoading(false);
    setSubmitted(true);
  };

  return (
    <section className="bg-background py-24 relative overflow-hidden">
      {/* Radial purple glow */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background:
            "radial-gradient(ellipse 70% 80% at 50% 60%, rgba(108,92,231,0.18) 0%, transparent 65%)",
        }}
      />

      <div className="container mx-auto px-6 relative z-10 max-w-2xl text-center">
        <h2 className="text-5xl font-bold text-foreground tracking-tight leading-tight">
          {title}
        </h2>
        <p className="mt-5 text-lg text-text-secondary leading-relaxed">{subtitle}</p>

        {submitted ? (
          <div className="mt-10 inline-flex items-center gap-3 px-6 py-4 rounded-xl border border-border bg-bg-elevated">
            <span className="text-success font-semibold">You&rsquo;re on the list!</span>
            <span className="text-text-secondary text-sm">We&rsquo;ll be in touch soon.</span>
          </div>
        ) : (
          <form
            onSubmit={handleSubmit}
            className="mt-10 flex flex-col sm:flex-row gap-3 justify-center"
          >
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@company.com"
              required
              className="flex-1 max-w-sm px-4 py-3 rounded-lg border border-border bg-bg-elevated text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
            />
            <a
              href={!email ? undefined : buttonHref}
              onClick={email ? undefined : (e) => e.preventDefault()}
              className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors duration-150 whitespace-nowrap cursor-pointer"
            >
              {loading ? (
                <Loader2 size={16} className="animate-spin" />
              ) : (
                <>
                  {buttonText}
                  <ArrowRight size={16} />
                </>
              )}
            </a>
          </form>
        )}

        <p className="mt-4 text-text-muted text-xs">
          No credit card required. Free forever on the starter plan.
        </p>
      </div>
    </section>
  );
}
