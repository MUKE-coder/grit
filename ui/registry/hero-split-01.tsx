import { ArrowRight, Terminal, Zap, Shield } from "lucide-react";

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
