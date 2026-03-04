package scaffold

// gritUIHeroSplit01TSX returns a split hero section: left text+CTAs, right dark code block.
// Inspired by Vercel homepage.
func gritUIHeroSplit01TSX() string {
	return `import { ArrowRight, Terminal, Zap, Shield } from "lucide-react";

export default function HeroSplit01() {
  return (
    <section className="min-h-screen bg-background flex items-center">
      <div className="container mx-auto px-6 py-24 grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
        {/* Left: Text + CTAs */}
        <div className="flex flex-col gap-8">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border border-border bg-bg-elevated w-fit">
            <Zap size={12} className="text-accent" />
            <span className="text-xs text-text-secondary font-medium tracking-wide">Now in public beta</span>
          </div>

          <div className="flex flex-col gap-4">
            <h1 className="text-5xl lg:text-6xl font-bold text-foreground leading-tight tracking-tight">
              Ship faster with{" "}
              <span className="text-accent">full-stack</span>{" "}
              confidence
            </h1>
            <p className="text-lg text-text-secondary leading-relaxed max-w-md">
              Grit combines Go and Next.js into a single CLI. Generate APIs, admin panels,
              and frontends in seconds — not days.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row gap-3">
            <a
              href="/get-started"
              className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors duration-150"
            >
              Get started free
              <ArrowRight size={16} />
            </a>
            <a
              href="/docs"
              className="inline-flex items-center justify-center gap-2 px-6 py-3 rounded-lg border border-border bg-bg-elevated hover:bg-bg-hover text-foreground font-semibold text-sm transition-colors duration-150"
            >
              <Terminal size={16} className="text-text-secondary" />
              View docs
            </a>
          </div>

          <div className="flex items-center gap-6 pt-2">
            <div className="flex items-center gap-2">
              <Shield size={14} className="text-success" />
              <span className="text-xs text-text-secondary">SOC 2 compliant</span>
            </div>
            <div className="flex items-center gap-2">
              <Zap size={14} className="text-warning" />
              <span className="text-xs text-text-secondary">Deploy in 60 seconds</span>
            </div>
          </div>
        </div>

        {/* Right: Code block */}
        <div className="relative">
          <div className="absolute -inset-1 rounded-2xl bg-accent opacity-10 blur-xl" />
          <div className="relative rounded-2xl border border-border bg-bg-secondary overflow-hidden">
            {/* Window chrome */}
            <div className="flex items-center gap-2 px-4 py-3 border-b border-border bg-bg-elevated">
              <div className="w-3 h-3 rounded-full bg-danger opacity-70" />
              <div className="w-3 h-3 rounded-full bg-warning opacity-70" />
              <div className="w-3 h-3 rounded-full bg-success opacity-70" />
              <span className="ml-3 text-xs text-text-muted font-mono">terminal</span>
            </div>

            {/* Code content */}
            <div className="p-6 font-mono text-sm leading-relaxed">
              <div className="flex gap-3">
                <span className="text-text-muted select-none">$</span>
                <span className="text-success">npx create-grit-app</span>
                <span className="text-text-secondary">my-saas</span>
              </div>
              <div className="mt-4 flex flex-col gap-1">
                <div className="text-text-muted">{"  "}✔ Scaffolding Go API...</div>
                <div className="text-text-muted">{"  "}✔ Setting up Next.js admin...</div>
                <div className="text-text-muted">{"  "}✔ Generating TypeScript types...</div>
                <div className="text-text-muted">{"  "}✔ Wiring Docker compose...</div>
                <div className="mt-2 text-success">  Project ready in 1.4s</div>
              </div>
              <div className="mt-4 flex gap-3">
                <span className="text-text-muted select-none">$</span>
                <span className="text-foreground">cd my-saas</span>
              </div>
              <div className="mt-2 flex gap-3">
                <span className="text-text-muted select-none">$</span>
                <span className="text-success">grit dev</span>
              </div>
              <div className="mt-2 flex flex-col gap-1">
                <div className="text-text-muted">{"  "}▲ API running on :8080</div>
                <div className="text-text-muted">{"  "}▲ Web running on :3000</div>
                <div className="text-text-muted">{"  "}▲ Admin running on :3001</div>
              </div>
              <div className="mt-3 flex items-center gap-2">
                <span className="inline-block w-2 h-4 bg-accent animate-pulse rounded-sm" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
`
}

// gritUILogoCLoud01TSX returns a "trusted by 500+ teams" infinite scroll marquee.
// Inspired by Stripe.
func gritUILogoCLoud01TSX() string {
	return `import React from "react";

const logos = [
  "Acme Corp",
  "Veritas",
  "Nexus Labs",
  "Orbit AI",
  "Stackline",
  "Pulsar Co",
  "Zenith",
  "Forge",
  "Altus",
  "Meridian",
  "Apex IO",
  "Crest",
];

export default function LogoCloud01() {
  const doubled = [...logos, ...logos];

  return (
    <section className="bg-background border-y border-border py-12 overflow-hidden">
      <div className="container mx-auto px-6 mb-8 text-center">
        <p className="text-sm text-text-muted font-medium tracking-widest uppercase">
          Trusted by 500+ teams worldwide
        </p>
      </div>

      <div className="relative flex overflow-hidden">
        {/* Fade masks */}
        <div
          className="absolute left-0 top-0 bottom-0 w-24 z-10 pointer-events-none"
          style={{ background: "linear-gradient(to right, var(--color-background, #0a0a0f), transparent)" }}
        />
        <div
          className="absolute right-0 top-0 bottom-0 w-24 z-10 pointer-events-none"
          style={{ background: "linear-gradient(to left, var(--color-background, #0a0a0f), transparent)" }}
        />

        {/* Scrolling track */}
        <div
          className="flex gap-16 items-center whitespace-nowrap"
          style={{
            animation: "marquee 28s linear infinite",
          }}
        >
          {doubled.map((name, i) => (
            <span
              key={i}
              className="text-text-muted font-semibold text-lg tracking-tight opacity-50 hover:opacity-80 transition-opacity duration-200 select-none"
            >
              {name}
            </span>
          ))}
        </div>
      </div>

      <style dangerouslySetInnerHTML={{ __html: "@keyframes marquee { from { transform: translateX(0); } to { transform: translateX(-50%); } }" }} />
    </section>
  );
}
`
}

// gritUITestimonialCard01TSX returns a single testimonial card with avatar, stars, quote.
// Inspired by Linear.
func gritUITestimonialCard01TSX() string {
	return `import { Star } from "lucide-react";

interface TestimonialCard01Props {
  quote: string;
  name: string;
  role: string;
  company: string;
  avatar?: string;
  stars?: number;
}

export default function TestimonialCard01({
  quote,
  name,
  role,
  company,
  avatar,
  stars = 5,
}: TestimonialCard01Props) {
  const initials = name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <div className="flex flex-col gap-5 p-6 rounded-2xl border border-border bg-bg-elevated hover:bg-bg-hover transition-colors duration-200">
      {/* Stars */}
      <div className="flex items-center gap-1">
        {Array.from({ length: 5 }).map((_, i) => (
          <Star
            key={i}
            size={14}
            className={i < stars ? "text-warning fill-warning" : "text-border fill-border"}
          />
        ))}
      </div>

      {/* Quote */}
      <blockquote className="text-foreground text-sm leading-relaxed">
        &ldquo;{quote}&rdquo;
      </blockquote>

      {/* Author */}
      <div className="flex items-center gap-3 pt-2 border-t border-border">
        {avatar ? (
          <img
            src={avatar}
            alt={name}
            className="w-9 h-9 rounded-full object-cover border border-border"
          />
        ) : (
          <div className="w-9 h-9 rounded-full bg-accent flex items-center justify-center text-white text-xs font-bold shrink-0">
            {initials}
          </div>
        )}
        <div className="flex flex-col min-w-0">
          <span className="text-foreground text-sm font-semibold truncate">{name}</span>
          <span className="text-text-muted text-xs truncate">
            {role} at {company}
          </span>
        </div>
      </div>
    </div>
  );
}
`
}

// gritUITestimonialGrid01TSX returns a 3-column grid of pre-filled testimonial cards.
// Inspired by Vercel's testimonial page.
func gritUITestimonialGrid01TSX() string {
	return `import { Star } from "lucide-react";

interface Testimonial {
  quote: string;
  name: string;
  role: string;
  company: string;
  stars: number;
}

const testimonials: Testimonial[] = [
  {
    quote:
      "Grit cut our backend setup time from a week to under an hour. The generated types and handlers are exactly what we would have written by hand — just faster.",
    name: "Sofia Ramirez",
    role: "CTO",
    company: "Stackline",
    stars: 5,
  },
  {
    quote:
      "The admin panel generation alone is worth it. We launched a fully featured dashboard to our clients in two days. Previously that was a month of work.",
    name: "James Whitfield",
    role: "Lead Engineer",
    company: "Orbit AI",
    stars: 5,
  },
  {
    quote:
      "Finally a Go framework that doesn't make me choose between productivity and control. The grit generate command is genuinely magical.",
    name: "Priya Nambiar",
    role: "Senior Developer",
    company: "Nexus Labs",
    stars: 5,
  },
];

function TestimonialCard({ t }: { t: Testimonial }) {
  const initials = t.name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <div className="flex flex-col gap-5 p-6 rounded-2xl border border-border bg-bg-elevated hover:bg-bg-hover transition-colors duration-200 h-full">
      <div className="flex items-center gap-1">
        {Array.from({ length: 5 }).map((_, i) => (
          <Star
            key={i}
            size={14}
            className={i < t.stars ? "text-warning fill-warning" : "text-border fill-border"}
          />
        ))}
      </div>
      <blockquote className="text-foreground text-sm leading-relaxed flex-1">
        &ldquo;{t.quote}&rdquo;
      </blockquote>
      <div className="flex items-center gap-3 pt-2 border-t border-border">
        <div className="w-9 h-9 rounded-full bg-accent flex items-center justify-center text-white text-xs font-bold shrink-0">
          {initials}
        </div>
        <div className="flex flex-col min-w-0">
          <span className="text-foreground text-sm font-semibold truncate">{t.name}</span>
          <span className="text-text-muted text-xs truncate">
            {t.role} at {t.company}
          </span>
        </div>
      </div>
    </div>
  );
}

export default function TestimonialGrid01() {
  return (
    <section className="bg-background py-24">
      <div className="container mx-auto px-6">
        <div className="text-center mb-16">
          <h2 className="text-4xl font-bold text-foreground tracking-tight">
            Loved by developers
          </h2>
          <p className="mt-4 text-text-secondary text-lg max-w-xl mx-auto">
            Teams across the world trust Grit to ship production-ready apps faster.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {testimonials.map((t) => (
            <TestimonialCard key={t.name} t={t} />
          ))}
        </div>
      </div>
    </section>
  );
}
`
}

// gritUIStatsBanner01TSX returns a dark stats strip with 4 big numbers.
// Inspired by Stripe.
func gritUIStatsBanner01TSX() string {
	return `interface Stat {
  value: string;
  label: string;
}

interface StatsBanner01Props {
  stats: Stat[];
}

export default function StatsBanner01({ stats }: StatsBanner01Props) {
  const displayStats =
    stats && stats.length > 0
      ? stats
      : [
          { value: "10K+", label: "Developers" },
          { value: "$2M+", label: "Saved in dev hours" },
          { value: "500+", label: "Companies" },
          { value: "99.9%", label: "Uptime" },
        ];

  return (
    <section className="bg-bg-secondary border-y border-border py-16 relative overflow-hidden">
      {/* Subtle radial glow */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background:
            "radial-gradient(ellipse 60% 100% at 50% 50%, rgba(108,92,231,0.07) 0%, transparent 70%)",
        }}
      />

      <div className="container mx-auto px-6 relative z-10">
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-8 lg:gap-0 lg:divide-x lg:divide-border">
          {displayStats.map((stat, i) => (
            <div
              key={i}
              className="flex flex-col items-center justify-center gap-2 py-4"
            >
              <span className="text-5xl font-bold text-foreground tracking-tight">
                {stat.value}
              </span>
              <span className="text-sm text-text-secondary font-medium text-center">
                {stat.label}
              </span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
`
}

// gritUIFaqAccordion01TSX returns a FAQ accordion with animated chevron.
// Inspired by Intercom FAQ page.
func gritUIFaqAccordion01TSX() string {
	return `"use client";

import { useState } from "react";
import { ChevronDown } from "lucide-react";

interface Faq {
  question: string;
  answer: string;
}

interface FaqAccordion01Props {
  faqs: Faq[];
}

const defaultFaqs: Faq[] = [
  {
    question: "What is Grit?",
    answer:
      "Grit is a full-stack meta-framework that combines a Go (Gin + GORM) backend with a Next.js frontend. A single CLI scaffolds your entire project, generates resources, and keeps your TypeScript types in sync with your Go models.",
  },
  {
    question: "Do I need to know both Go and TypeScript?",
    answer:
      "Basic Go knowledge is helpful, but Grit is designed to be beginner-friendly. The scaffolded code follows strict conventions so you can learn as you build. The frontend is pure TypeScript/React, so frontend developers can start there.",
  },
  {
    question: "How is Grit different from other Go frameworks?",
    answer:
      "Most Go frameworks handle only the backend. Grit generates the full stack — Go API, Next.js admin panel, web frontend, shared types, Docker setup, and tests — all wired together and production-ready from day one.",
  },
  {
    question: "Is there a self-hosted option?",
    answer:
      "Yes. Grit is open-source and fully self-hostable. The included Docker Compose file starts PostgreSQL, Redis, MinIO, and Mailhog locally. You can deploy to any cloud provider that runs Docker.",
  },
];

export default function FaqAccordion01({ faqs }: FaqAccordion01Props) {
  const items = faqs && faqs.length > 0 ? faqs : defaultFaqs;
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  return (
    <section className="bg-background py-24">
      <div className="container mx-auto px-6 max-w-3xl">
        <div className="text-center mb-14">
          <h2 className="text-4xl font-bold text-foreground tracking-tight">
            Frequently asked questions
          </h2>
          <p className="mt-4 text-text-secondary text-lg">
            Everything you need to know before getting started.
          </p>
        </div>

        <div className="flex flex-col divide-y divide-border rounded-2xl border border-border overflow-hidden">
          {items.map((faq, i) => {
            const isOpen = openIndex === i;
            return (
              <div key={i} className="bg-bg-elevated">
                <button
                  onClick={() => setOpenIndex(isOpen ? null : i)}
                  className="w-full flex items-center justify-between gap-4 px-6 py-5 text-left hover:bg-bg-hover transition-colors duration-150"
                  aria-expanded={isOpen}
                >
                  <span className="text-foreground font-semibold text-sm leading-snug">
                    {faq.question}
                  </span>
                  <ChevronDown
                    size={18}
                    className="text-text-secondary shrink-0 transition-transform duration-200"
                    style={{ transform: isOpen ? "rotate(180deg)" : "rotate(0deg)" }}
                  />
                </button>

                <div
                  className="overflow-hidden transition-all duration-200"
                  style={{ maxHeight: isOpen ? "400px" : "0px" }}
                >
                  <p className="px-6 pb-5 text-text-secondary text-sm leading-relaxed">
                    {faq.answer}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
`
}

// gritUICtaBanner01TSX returns a full-width gradient CTA with email input.
// Inspired by Vercel's homepage bottom CTA.
func gritUICtaBanner01TSX() string {
	return `"use client";

import { useState } from "react";
import { ArrowRight, Loader2 } from "lucide-react";

interface CtaBanner01Props {
  title: string;
  subtitle: string;
  buttonText: string;
  buttonHref: string;
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
`
}

// gritUINewsletterSignup01TSX returns a newsletter signup card with subscriber count.
// Inspired by Substack/Ghost.
func gritUINewsletterSignup01TSX() string {
	return `"use client";

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
`
}

// gritUIChangelogItem01TSX returns a changelog entry with version, date, bullet list.
// Inspired by Vercel changelog.
func gritUIChangelogItem01TSX() string {
	return `import { Tag, Calendar, CheckCircle2 } from "lucide-react";

interface ChangelogItem01Props {
  version: string;
  date: string;
  title: string;
  description: string;
  changes: string[];
}

export default function ChangelogItem01({
  version = "v1.4.0",
  date = "March 3, 2026",
  title = "Improved resource generation and dark UI polish",
  description = "This release focuses on speed improvements to the code generator and a visual refresh of the admin scaffold with better dark mode contrast.",
  changes = [
    "grit generate now runs 3x faster with parallel file writes",
    "Admin panel sidebar: collapsible groups and animated transitions",
    "New StatsBanner and TestimonialGrid registry components",
    "Fixed race condition in grit sync when multiple models changed",
    "Upgraded to Next.js 15 in all scaffolded frontend templates",
  ],
}: ChangelogItem01Props) {
  return (
    <article className="relative flex gap-8 py-10 border-b border-border last:border-0">
      {/* Timeline dot */}
      <div className="hidden md:flex flex-col items-center pt-1">
        <div className="w-3 h-3 rounded-full bg-accent shrink-0 ring-4 ring-background" />
        <div className="flex-1 w-px bg-border mt-3" />
      </div>

      <div className="flex-1 min-w-0">
        {/* Meta */}
        <div className="flex flex-wrap items-center gap-3 mb-4">
          <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-accent text-white text-xs font-bold tracking-wide">
            <Tag size={11} />
            {version}
          </span>
          <span className="inline-flex items-center gap-1.5 text-text-muted text-xs">
            <Calendar size={12} />
            {date}
          </span>
        </div>

        {/* Title + description */}
        <h3 className="text-foreground text-xl font-bold tracking-tight mb-3">{title}</h3>
        <p className="text-text-secondary text-sm leading-relaxed mb-6">{description}</p>

        {/* Changes list */}
        <div className="rounded-xl border border-border bg-bg-elevated p-5">
          <p className="text-xs text-text-muted font-semibold uppercase tracking-widest mb-4">
            What&rsquo;s new
          </p>
          <ul className="flex flex-col gap-3">
            {changes.map((change, i) => (
              <li key={i} className="flex items-start gap-3">
                <CheckCircle2 size={15} className="text-success shrink-0 mt-0.5" />
                <span className="text-foreground text-sm leading-snug">{change}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </article>
  );
}
`
}

// gritUIPricingToggle01TSX returns a full pricing section with monthly/annual toggle.
// Inspired by Stripe pricing page.
func gritUIPricingToggle01TSX() string {
	return `"use client";

import { useState } from "react";
import { Check, Zap } from "lucide-react";

interface Plan {
  name: string;
  monthlyPrice: number;
  annualPrice: number;
  description: string;
  features: string[];
  cta: string;
  popular: boolean;
  enterprise: boolean;
}

const plans: Plan[] = [
  {
    name: "Starter",
    monthlyPrice: 0,
    annualPrice: 0,
    description: "For solo developers and side projects.",
    features: [
      "Up to 3 projects",
      "grit generate (all resource types)",
      "Docker Compose setup",
      "Community support",
      "1 team member",
    ],
    cta: "Start for free",
    popular: false,
    enterprise: false,
  },
  {
    name: "Pro",
    monthlyPrice: 29,
    annualPrice: 23,
    description: "For professional developers and small teams.",
    features: [
      "Unlimited projects",
      "Everything in Starter",
      "Grit UI component registry",
      "Priority email support",
      "Up to 10 team members",
      "Audit logs",
    ],
    cta: "Get started",
    popular: true,
    enterprise: false,
  },
  {
    name: "Enterprise",
    monthlyPrice: 99,
    annualPrice: 79,
    description: "For growing teams with advanced needs.",
    features: [
      "Everything in Pro",
      "SSO (SAML / OIDC)",
      "Custom deployment support",
      "SLA + dedicated Slack",
      "Unlimited team members",
      "Custom billing",
    ],
    cta: "Contact sales",
    popular: false,
    enterprise: true,
  },
];

export default function PricingToggle01() {
  const [annual, setAnnual] = useState(false);

  return (
    <section className="bg-background py-24">
      <div className="container mx-auto px-6">
        {/* Header */}
        <div className="text-center mb-12">
          <h2 className="text-4xl font-bold text-foreground tracking-tight">
            Simple, transparent pricing
          </h2>
          <p className="mt-4 text-text-secondary text-lg max-w-xl mx-auto">
            Start free, upgrade when you&rsquo;re ready. No hidden fees.
          </p>

          {/* Toggle */}
          <div className="mt-8 inline-flex items-center gap-3 p-1 rounded-xl border border-border bg-bg-elevated">
            <button
              onClick={() => setAnnual(false)}
              className={"px-5 py-2 rounded-lg text-sm font-semibold transition-colors duration-150 " +
                (!annual
                  ? "bg-accent text-white"
                  : "text-text-secondary hover:text-foreground")}
            >
              Monthly
            </button>
            <button
              onClick={() => setAnnual(true)}
              className={"px-5 py-2 rounded-lg text-sm font-semibold transition-colors duration-150 flex items-center gap-2 " +
                (annual
                  ? "bg-accent text-white"
                  : "text-text-secondary hover:text-foreground")}
            >
              Annual
              <span className="text-xs px-2 py-0.5 rounded-full bg-success text-background font-bold">
                Save 20%
              </span>
            </button>
          </div>
        </div>

        {/* Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto">
          {plans.map((plan) => (
            <div
              key={plan.name}
              className={"relative flex flex-col rounded-2xl border p-7 transition-all duration-200 " +
                (plan.popular
                  ? "border-accent bg-bg-elevated ring-2 ring-accent ring-offset-2 ring-offset-background scale-105"
                  : "border-border bg-bg-elevated hover:bg-bg-hover")}
            >
              {plan.popular && (
                <div className="absolute -top-4 left-1/2 -translate-x-1/2">
                  <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-accent text-white text-xs font-bold tracking-wide shadow-lg">
                    <Zap size={11} />
                    Most popular
                  </span>
                </div>
              )}

              <div className="mb-6">
                <h3 className="text-foreground font-bold text-lg">{plan.name}</h3>
                <p className="text-text-muted text-sm mt-1">{plan.description}</p>
              </div>

              <div className="mb-6">
                <div className="flex items-end gap-1">
                  <span className="text-5xl font-bold text-foreground tracking-tight">
                    {plan.enterprise
                      ? "Custom"
                      : (annual ? "$" + plan.annualPrice : "$" + plan.monthlyPrice)}
                  </span>
                  {!plan.enterprise && (
                    <span className="text-text-muted text-sm mb-2">/mo</span>
                  )}
                </div>
                {annual && !plan.enterprise && plan.annualPrice > 0 && (
                  <p className="text-text-muted text-xs mt-1">
                    Billed annually (${plan.annualPrice * 12}/yr)
                  </p>
                )}
              </div>

              <ul className="flex flex-col gap-3 mb-8 flex-1">
                {plan.features.map((feature, i) => (
                  <li key={i} className="flex items-center gap-2.5">
                    <Check size={15} className="text-success shrink-0" />
                    <span className="text-text-secondary text-sm">{feature}</span>
                  </li>
                ))}
              </ul>

              <a
                href={plan.enterprise ? "/contact" : "/signup"}
                className={"w-full inline-flex items-center justify-center py-3 rounded-lg font-semibold text-sm transition-colors duration-150 " +
                  (plan.popular
                    ? "bg-accent hover:bg-accent-hover text-white"
                    : "border border-border bg-bg-secondary hover:bg-bg-hover text-foreground")}
              >
                {plan.cta}
              </a>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
`
}

// gritUILoginCard01TSX returns a premium dark login card.
// Inspired by Linear's login page.
func gritUILoginCard01TSX() string {
	return `"use client";

import { useState } from "react";
import { Eye, EyeOff, Github, Loader2 } from "lucide-react";

export default function LoginCard01() {
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [remember, setRemember] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    await new Promise((r) => setTimeout(r, 1000));
    setLoading(false);
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

        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-foreground tracking-tight">
            Sign in to your account
          </h1>
          <p className="mt-2 text-sm text-text-secondary">
            Don&rsquo;t have an account?{" "}
            <a href="/signup" className="text-accent hover:underline font-medium">
              Sign up
            </a>
          </p>
        </div>

        {/* Social buttons */}
        <div className="flex flex-col gap-3 mb-6">
          <button className="w-full flex items-center justify-center gap-3 px-4 py-2.5 rounded-lg border border-border bg-bg-elevated hover:bg-bg-hover text-foreground text-sm font-medium transition-colors duration-150">
            <svg width="18" height="18" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            Continue with Google
          </button>
          <button className="w-full flex items-center justify-center gap-3 px-4 py-2.5 rounded-lg border border-border bg-bg-elevated hover:bg-bg-hover text-foreground text-sm font-medium transition-colors duration-150">
            <Github size={18} />
            Continue with GitHub
          </button>
        </div>

        {/* OR divider */}
        <div className="relative mb-6">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-border" />
          </div>
          <div className="relative flex justify-center text-xs">
            <span className="bg-background px-3 text-text-muted">or continue with email</span>
          </div>
        </div>

        {/* Email/password form */}
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1.5">
            <label className="text-sm font-medium text-foreground">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@company.com"
              required
              className="px-3 py-2.5 rounded-lg border border-border bg-bg-elevated text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
            />
          </div>

          <div className="flex flex-col gap-1.5">
            <div className="flex items-center justify-between">
              <label className="text-sm font-medium text-foreground">Password</label>
              <a href="/forgot-password" className="text-xs text-accent hover:underline">
                Forgot password?
              </a>
            </div>
            <div className="relative">
              <input
                type={showPassword ? "text" : "password"}
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
                required
                className="w-full px-3 py-2.5 pr-10 rounded-lg border border-border bg-bg-elevated text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-text-muted hover:text-foreground"
              >
                {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
              </button>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <input
              id="remember"
              type="checkbox"
              checked={remember}
              onChange={(e) => setRemember(e.target.checked)}
              className="w-4 h-4 rounded border-border accent-accent"
            />
            <label htmlFor="remember" className="text-sm text-text-secondary cursor-pointer">
              Remember me for 30 days
            </label>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-2.5 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors duration-150 disabled:opacity-60 flex items-center justify-center gap-2"
          >
            {loading ? <Loader2 size={16} className="animate-spin" /> : "Sign in"}
          </button>
        </form>
      </div>
    </div>
  );
}
`
}

// gritUISignupCard01TSX returns a sign up card with feature list on the left.
// Inspired by Vercel sign-up.
func gritUISignupCard01TSX() string {
	return `"use client";

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
`
}

// gritUIForgotPasswordCard01TSX returns a compact forgot password card with success state.
// Inspired by GitHub.
func gritUIForgotPasswordCard01TSX() string {
	return `"use client";

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
`
}

// gritUIOtpInput01TSX returns a 6-digit OTP input with auto-advance and resend timer.
// Inspired by Vercel's MFA flow.
func gritUIOtpInput01TSX() string {
	return `"use client";

import { useState, useRef, useEffect, KeyboardEvent, ClipboardEvent } from "react";
import { RefreshCw } from "lucide-react";

interface OtpInput01Props {
  email: string;
  onComplete: (code: string) => void;
}

export default function OtpInput01({ email = "you@company.com", onComplete }: OtpInput01Props) {
  const [digits, setDigits] = useState<string[]>(Array(6).fill(""));
  const [countdown, setCountdown] = useState(60);
  const [canResend, setCanResend] = useState(false);
  const inputRefs = useRef<(HTMLInputElement | null)[]>([]);

  useEffect(() => {
    if (countdown <= 0) {
      setCanResend(true);
      return;
    }
    const t = setTimeout(() => setCountdown((c) => c - 1), 1000);
    return () => clearTimeout(t);
  }, [countdown]);

  const handleResend = () => {
    setCountdown(60);
    setCanResend(false);
    setDigits(Array(6).fill(""));
    inputRefs.current[0]?.focus();
  };

  const handleChange = (index: number, value: string) => {
    const char = value.replace(/\D/g, "").slice(-1);
    const next = [...digits];
    next[index] = char;
    setDigits(next);
    if (char && index < 5) {
      inputRefs.current[index + 1]?.focus();
    }
    if (next.every((d) => d !== "")) {
      onComplete(next.join(""));
    }
  };

  const handleKeyDown = (index: number, e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Backspace") {
      if (digits[index]) {
        const next = [...digits];
        next[index] = "";
        setDigits(next);
      } else if (index > 0) {
        inputRefs.current[index - 1]?.focus();
      }
    }
    if (e.key === "ArrowLeft" && index > 0) inputRefs.current[index - 1]?.focus();
    if (e.key === "ArrowRight" && index < 5) inputRefs.current[index + 1]?.focus();
  };

  const handlePaste = (e: ClipboardEvent<HTMLInputElement>) => {
    e.preventDefault();
    const pasted = e.clipboardData.getData("text").replace(/\D/g, "").slice(0, 6);
    if (!pasted) return;
    const next = Array(6).fill("");
    pasted.split("").forEach((char, i) => { next[i] = char; });
    setDigits(next);
    const lastFilled = Math.min(pasted.length, 5);
    inputRefs.current[lastFilled]?.focus();
    if (pasted.length === 6) onComplete(pasted);
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <div className="w-full max-w-sm text-center">
        <div className="flex justify-center mb-8">
          <div className="w-12 h-12 rounded-2xl bg-accent flex items-center justify-center">
            <span className="text-white font-black text-xl tracking-tight">G</span>
          </div>
        </div>

        <h1 className="text-2xl font-bold text-foreground tracking-tight mb-2">
          Check your email
        </h1>
        <p className="text-text-secondary text-sm mb-1">
          We sent a 6-digit code to
        </p>
        <p className="text-foreground text-sm font-semibold mb-8">{email}</p>

        {/* OTP boxes */}
        <div className="flex justify-center gap-3 mb-8">
          {digits.map((digit, i) => (
            <input
              key={i}
              ref={(el) => { inputRefs.current[i] = el; }}
              type="text"
              inputMode="numeric"
              maxLength={1}
              value={digit}
              onChange={(e) => handleChange(i, e.target.value)}
              onKeyDown={(e) => handleKeyDown(i, e)}
              onPaste={handlePaste}
              className={"w-12 h-14 text-center text-xl font-bold rounded-xl border transition-colors duration-150 bg-bg-elevated text-foreground focus:outline-none " +
                (digit
                  ? "border-accent text-accent"
                  : "border-border focus:border-accent")}
            />
          ))}
        </div>

        {/* Resend */}
        <div className="flex items-center justify-center gap-2 text-sm">
          {canResend ? (
            <button
              onClick={handleResend}
              className="inline-flex items-center gap-1.5 text-accent hover:underline font-medium"
            >
              <RefreshCw size={14} />
              Resend code
            </button>
          ) : (
            <p className="text-text-muted">
              Resend code in{" "}
              <span className="text-text-secondary font-semibold">{countdown}s</span>
            </p>
          )}
        </div>

        <p className="mt-6 text-xs text-text-muted">
          Didn&rsquo;t get the email? Check your spam folder.
        </p>
      </div>
    </div>
  );
}
`
}

// gritUISocialLoginButtons01TSX returns Google + GitHub OAuth buttons with divider.
// Inspired by Clerk auth components.
func gritUISocialLoginButtons01TSX() string {
	return `import { Github } from "lucide-react";

interface SocialLoginButtons01Props {
  googleHref?: string;
  githubHref?: string;
}

export default function SocialLoginButtons01({
  googleHref = "/auth/google",
  githubHref = "/auth/github",
}: SocialLoginButtons01Props) {
  return (
    <div className="flex flex-col gap-4 w-full">
      {/* Divider */}
      <div className="relative flex items-center gap-4">
        <div className="flex-1 h-px bg-border" />
        <span className="text-text-muted text-xs font-medium whitespace-nowrap">
          or continue with
        </span>
        <div className="flex-1 h-px bg-border" />
      </div>

      {/* Buttons */}
      <a
        href={googleHref}
        className="w-full flex items-center justify-center gap-3 px-4 py-2.5 rounded-lg border border-border bg-bg-elevated hover:bg-bg-hover text-foreground text-sm font-medium transition-colors duration-150"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
          <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
          <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
          <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
          <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
        </svg>
        Continue with Google
      </a>

      <a
        href={githubHref}
        className="w-full flex items-center justify-center gap-3 px-4 py-2.5 rounded-lg border border-border bg-bg-elevated hover:bg-bg-hover text-foreground text-sm font-medium transition-colors duration-150"
      >
        <Github size={18} />
        Continue with GitHub
      </a>
    </div>
  );
}
`
}

// gritUIOnboardingWizard01TSX returns a 3-step onboarding wizard.
// Inspired by Notion's new workspace setup.
func gritUIOnboardingWizard01TSX() string {
	return `"use client";

import { useState } from "react";
import { ArrowRight, ArrowLeft, Check, User, Building2, Zap } from "lucide-react";

interface OnboardingWizard01Props {
  onComplete?: () => void;
}

const roles = ["Developer", "Designer", "Product Manager", "Founder", "Other"];
const teamSizes = ["Just me", "2–5", "6–15", "16–50", "50+"];

const plans = [
  {
    name: "Starter",
    price: "Free",
    description: "Perfect to explore Grit",
    features: ["3 projects", "Community support"],
  },
  {
    name: "Pro",
    price: "$29/mo",
    description: "For professionals",
    features: ["Unlimited projects", "Priority support", "Grit UI registry"],
    popular: true,
  },
  {
    name: "Enterprise",
    price: "Custom",
    description: "For larger teams",
    features: ["SSO", "SLA", "Dedicated support"],
  },
];

export default function OnboardingWizard01({ onComplete }: OnboardingWizard01Props) {
  const [step, setStep] = useState(0);
  const [fullName, setFullName] = useState("");
  const [role, setRole] = useState("");
  const [company, setCompany] = useState("");
  const [teamSize, setTeamSize] = useState("");
  const [selectedPlan, setSelectedPlan] = useState("Pro");

  const totalSteps = 3;

  const stepMeta = [
    { icon: User, title: "Tell us about you", subtitle: "Help us personalize your experience." },
    { icon: Building2, title: "Set up your workspace", subtitle: "A few details about where you work." },
    { icon: Zap, title: "Choose your plan", subtitle: "Start free, upgrade anytime." },
  ];

  const canAdvance = () => {
    if (step === 0) return fullName.trim() && role;
    if (step === 1) return company.trim() && teamSize;
    return true;
  };

  const handleNext = () => {
    if (step < totalSteps - 1) {
      setStep(step + 1);
    } else {
      onComplete?.();
    }
  };

  const Icon = stepMeta[step].icon;

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-lg">
        {/* Step dots */}
        <div className="flex items-center justify-center gap-2 mb-10">
          {Array.from({ length: totalSteps }).map((_, i) => (
            <div
              key={i}
              className={"transition-all duration-300 rounded-full " +
                (i === step
                  ? "w-6 h-2 bg-accent"
                  : i < step
                  ? "w-2 h-2 bg-accent opacity-60"
                  : "w-2 h-2 bg-border")}
            />
          ))}
        </div>

        <div className="rounded-2xl border border-border bg-bg-elevated p-8">
          {/* Header */}
          <div className="flex items-center gap-4 mb-8">
            <div className="w-11 h-11 rounded-xl bg-accent flex items-center justify-center shrink-0">
              <Icon size={20} className="text-white" />
            </div>
            <div>
              <h2 className="text-xl font-bold text-foreground tracking-tight">
                {stepMeta[step].title}
              </h2>
              <p className="text-text-secondary text-sm mt-0.5">{stepMeta[step].subtitle}</p>
            </div>
          </div>

          {/* Step content */}
          {step === 0 && (
            <div className="flex flex-col gap-5">
              <div className="flex flex-col gap-1.5">
                <label className="text-sm font-medium text-foreground">Full name</label>
                <input
                  type="text"
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  placeholder="Alex Johnson"
                  className="px-3 py-2.5 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
                />
              </div>
              <div className="flex flex-col gap-2">
                <label className="text-sm font-medium text-foreground">Your role</label>
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                  {roles.map((r) => (
                    <button
                      key={r}
                      type="button"
                      onClick={() => setRole(r)}
                      className={"px-3 py-2 rounded-lg border text-sm font-medium transition-colors duration-150 " +
                        (role === r
                          ? "border-accent bg-accent text-white"
                          : "border-border bg-bg-secondary text-text-secondary hover:bg-bg-hover hover:text-foreground")}
                    >
                      {r}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          )}

          {step === 1 && (
            <div className="flex flex-col gap-5">
              <div className="flex flex-col gap-1.5">
                <label className="text-sm font-medium text-foreground">Company name</label>
                <input
                  type="text"
                  value={company}
                  onChange={(e) => setCompany(e.target.value)}
                  placeholder="Acme Corp"
                  className="px-3 py-2.5 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
                />
              </div>
              <div className="flex flex-col gap-2">
                <label className="text-sm font-medium text-foreground">Team size</label>
                <div className="grid grid-cols-3 sm:grid-cols-5 gap-2">
                  {teamSizes.map((size) => (
                    <button
                      key={size}
                      type="button"
                      onClick={() => setTeamSize(size)}
                      className={"px-3 py-2 rounded-lg border text-sm font-medium transition-colors duration-150 " +
                        (teamSize === size
                          ? "border-accent bg-accent text-white"
                          : "border-border bg-bg-secondary text-text-secondary hover:bg-bg-hover hover:text-foreground")}
                    >
                      {size}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          )}

          {step === 2 && (
            <div className="flex flex-col gap-3">
              {plans.map((plan) => (
                <button
                  key={plan.name}
                  type="button"
                  onClick={() => setSelectedPlan(plan.name)}
                  className={"relative w-full flex items-start gap-4 p-4 rounded-xl border text-left transition-all duration-150 " +
                    (selectedPlan === plan.name
                      ? "border-accent bg-bg-secondary ring-1 ring-accent"
                      : "border-border bg-bg-secondary hover:bg-bg-hover")}
                >
                  <div className={"w-5 h-5 rounded-full border-2 flex items-center justify-center shrink-0 mt-0.5 transition-colors " +
                    (selectedPlan === plan.name ? "border-accent bg-accent" : "border-border")}
                  >
                    {selectedPlan === plan.name && <Check size={11} className="text-white" />}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className="text-foreground font-semibold text-sm">{plan.name}</span>
                      {(plan as { popular?: boolean }).popular && (
                        <span className="px-2 py-0.5 rounded-full bg-accent text-white text-xs font-bold">
                          Popular
                        </span>
                      )}
                      <span className="ml-auto text-foreground font-bold text-sm">{plan.price}</span>
                    </div>
                    <p className="text-text-muted text-xs mt-0.5">{plan.description}</p>
                    <ul className="mt-2 flex flex-wrap gap-x-4 gap-y-1">
                      {plan.features.map((f) => (
                        <li key={f} className="flex items-center gap-1.5 text-text-secondary text-xs">
                          <Check size={11} className="text-success" />
                          {f}
                        </li>
                      ))}
                    </ul>
                  </div>
                </button>
              ))}
            </div>
          )}

          {/* Navigation */}
          <div className="flex items-center justify-between mt-8 pt-6 border-t border-border">
            <button
              type="button"
              onClick={() => setStep(step - 1)}
              disabled={step === 0}
              className="inline-flex items-center gap-2 px-4 py-2.5 rounded-lg border border-border bg-bg-secondary hover:bg-bg-hover text-foreground text-sm font-medium transition-colors duration-150 disabled:opacity-0 disabled:pointer-events-none"
            >
              <ArrowLeft size={15} />
              Back
            </button>

            <div className="text-text-muted text-xs">
              Step {step + 1} of {totalSteps}
            </div>

            <button
              type="button"
              onClick={handleNext}
              disabled={!canAdvance()}
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors duration-150 disabled:opacity-50"
            >
              {step < totalSteps - 1 ? (
                <>Next <ArrowRight size={15} /></>
              ) : (
                <>Get started <Zap size={15} /></>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
`
}
