"use client";

import { useState } from "react";
import { Mail, Users, ArrowRight, CheckCircle } from "lucide-react";

export default function NewsletterSignup01() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email) return;
    setLoading(true);
    await new Promise((r) => setTimeout(r, 800));
    setLoading(false);
    setSubmitted(true);
  };

  return (
    <div className="rounded-2xl border border-border bg-bg-elevated p-8 max-w-lg mx-auto">
      <div className="flex items-center gap-3 mb-6">
        <div className="w-10 h-10 rounded-xl bg-accent flex items-center justify-center shrink-0">
          <Mail size={18} className="text-white" />
        </div>
        <div>
          <h3 className="text-foreground font-bold text-lg leading-tight">
            The Grit Newsletter
          </h3>
          <p className="text-text-muted text-xs mt-0.5">Weekly Go + React tips</p>
        </div>
      </div>

      <div className="flex items-center gap-2 mb-5 px-3 py-2 rounded-lg bg-bg-secondary border border-border w-fit">
        <Users size={13} className="text-accent" />
        <span className="text-text-secondary text-xs font-medium">
          Join <span className="text-foreground font-bold">5,000+</span> developers
        </span>
      </div>

      <p className="text-text-secondary text-sm leading-relaxed mb-6">
        Get practical tutorials, framework updates, and behind-the-scenes content — straight
        to your inbox. No spam, unsubscribe anytime.
      </p>

      {submitted ? (
        <div className="flex items-center gap-3 px-4 py-3 rounded-lg border border-border bg-bg-secondary">
          <CheckCircle size={18} className="text-success shrink-0" />
          <div>
            <p className="text-foreground text-sm font-semibold">You&rsquo;re subscribed!</p>
            <p className="text-text-muted text-xs mt-0.5">Check your inbox for a confirmation.</p>
          </div>
        </div>
      ) : (
        <form onSubmit={handleSubmit} className="flex flex-col sm:flex-row gap-3">
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@example.com"
            required
            className="flex-1 px-4 py-2.5 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
          />
          <button
            type="submit"
            disabled={loading}
            className="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors duration-150 disabled:opacity-60 whitespace-nowrap"
          >
            {loading ? "Subscribing..." : (
              <>Subscribe <ArrowRight size={14} /></>
            )}
          </button>
        </form>
      )}
    </div>
  );
}
