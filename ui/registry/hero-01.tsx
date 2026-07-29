import { ArrowRight, Sparkles } from "lucide-react";

export default function Hero01() {
  return (
    <section className="relative min-h-screen bg-background flex items-center overflow-hidden">
      {/* Radial accent glow behind the headline */}
      <div
        aria-hidden
        className="pointer-events-none absolute inset-0"
        style={{
          background:
            "radial-gradient(ellipse 60% 50% at 50% 0%, rgba(108,92,231,0.18), transparent 60%)",
        }}
      />

      <div className="relative container mx-auto px-6 py-24 flex flex-col items-center text-center gap-8">
        {/* Announcement badge */}
        <a
          href="#"
          className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-border bg-bg-elevated hover:bg-bg-hover transition-colors"
        >
          <Sparkles size={13} className="text-accent" />
          <span className="text-xs font-medium text-text-secondary tracking-wide">
            Introducing Grit UI — 100 components
          </span>
          <ArrowRight size={12} className="text-text-muted" />
        </a>

        <h1 className="max-w-4xl text-5xl md:text-7xl font-bold text-foreground leading-[1.05] tracking-tight">
          Build your next idea{" "}
          <span className="bg-gradient-to-r from-accent via-info to-accent bg-clip-text text-transparent">
            remarkably fast
          </span>
        </h1>

        <p className="max-w-2xl text-lg text-text-secondary leading-relaxed">
          Production-ready components, a typed API, and an admin panel — generated
          in one command. Spend your time on the product, not the plumbing.
        </p>

        <div className="flex flex-col sm:flex-row items-center gap-3 mt-2">
          <a
            href="#"
            className="group inline-flex items-center gap-2 h-11 px-7 rounded-full bg-accent text-accent-foreground font-semibold text-sm hover:bg-accent-hover transition-colors"
          >
            Get started
            <ArrowRight
              size={15}
              className="transition-transform group-hover:translate-x-0.5"
            />
          </a>
          <a
            href="#"
            className="inline-flex items-center h-11 px-7 rounded-full border border-border text-foreground font-medium text-sm hover:bg-bg-elevated transition-colors"
          >
            Read the docs
          </a>
        </div>

        <p className="text-xs text-text-muted font-mono mt-4">
          MIT licensed · No signup required
        </p>
      </div>
    </section>
  );
}
