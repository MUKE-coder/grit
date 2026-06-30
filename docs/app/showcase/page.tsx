import Link from "next/link";
import type { Metadata } from "next";
import {
  ExternalLink,
  Globe,
  Database,
  Layers,
  ArrowRight,
  Sparkles,
  Box,
  Zap,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { ShowcaseGrid, type ShowcaseProject } from "@/components/showcase-grid";

export const metadata: Metadata = {
  title: "Showcase | Grit",
  description:
    "Explore real-world applications built with the Grit framework. See what's possible with Go + React full-stack development.",
  openGraph: {
    title: "Showcase | Grit",
    description:
      "Explore real-world applications built with the Grit framework.",
    url: "https://gritframework.dev/showcase",
  },
};

// ── Showcase data ──────────────────────────────────────────────

const projects: ShowcaseProject[] = [
  {
    name: "GritCMS",
    url: "https://gritcms.com",
    description: "All-in-one business platform",
    longDescription:
      "A comprehensive content management and business automation platform with email marketing, course hosting, community spaces, funnels, booking, affiliate management, and more. Built entirely with the Grit framework.",
    image: "/showcase/gritcms.png",
    tags: ["SaaS", "CMS", "E-Commerce", "Email Marketing", "LMS"],
    techStack: ["Go", "Gin", "GORM", "React", "Next.js", "Tailwind CSS", "PostgreSQL", "Redis", "Cloudflare R2"],
    stats: {
      tables: "67",
      models: "64",
      goroutines: "27",
      modules: "8+",
      highlights: [
        "Website & blog builder",
        "Email campaigns & sequences",
        "Course creation & LMS",
        "Product & order management",
        "Community spaces",
        "Sales & opt-in funnels",
        "Booking & scheduling",
        "Automation workflows",
        "Affiliate management",
        "Contact management",
      ],
    },
    featured: true,
  },
  {
    name: "Socially Pro",
    url: "https://sociallypro-web.devsamk.com",
    description: "Social media management platform",
    longDescription:
      "A comprehensive social media management platform that empowers teams to plan, publish, and analyze content across 7+ social platforms from a unified dashboard. Features workflow approvals, AI-powered content assistance, and role-based access control.",
    image: "/showcase/sociallypro.png",
    tags: ["SaaS", "Social Media", "AI", "Analytics"],
    techStack: ["Go", "Gin", "GORM", "React", "Next.js", "Tailwind CSS", "PostgreSQL", "Redis", "Cloudflare R2", "Google OAuth", "Gemini", "Claude API"],
    stats: {
      tables: "20",
      models: "18",
      goroutines: "",
      modules: "10+",
      highlights: [
        "Multi-platform publishing (7+ platforms)",
        "Content scheduling & queue management",
        "Team collaboration & approvals",
        "AI content assistant (Gemini, Claude, OpenAI)",
        "Unified inbox for engagement",
        "Media library with R2 storage",
        "Cross-channel analytics",
        "Role-based access control",
        "Multi-workspace support",
        "Google OAuth authentication",
      ],
    },
  },
  {
    name: "DGateway",
    url: "https://dgateway.desispay.com",
    description: "Unified payment aggregation for East Africa",
    longDescription:
      "A unified payment aggregation platform for East Africa. Businesses integrate once through a single REST API and instantly gain access to MTN Mobile Money, Airtel Money, card payments via Stripe, and PesaPal — without managing multiple provider integrations. Includes a full merchant admin dashboard for managing apps, API keys, transactions, withdrawals, webhooks, and teams.",
    image: "/showcase/dgateway-wide.png",
    tags: ["FinTech", "Payments", "SaaS", "API Platform"],
    techStack: ["Go", "Gin", "GORM", "React", "Next.js 15", "Tailwind CSS", "PostgreSQL", "Turborepo", "Stripe", "PesaPal", "MTN MoMo", "Airtel Money", "Resend", "Dokploy"],
    stats: {
      tables: "11",
      models: "11",
      goroutines: "",
      modules: "12",
      highlights: [
        "Payment collection (mobile money & card)",
        "Disbursements to mobile wallets",
        "Multi-provider routing (MTN, Airtel, Stripe, PesaPal)",
        "App management with isolated API keys",
        "Webhook delivery with signed payloads",
        "Transaction ledger & status tracking",
        "Withdrawal requests & payout management",
        "Team invitations with RBAC",
        "CORS domain management per app",
        "Blog & content management",
        "Public landing page & API docs",
      ],
    },
  },
];

// ── Page ────────────────────────────────────────────────────────

const featured = projects.filter((p) => p.featured);
const rest = projects.filter((p) => !p.featured);

export default function ShowcasePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      {/* Hero */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 -z-10">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[500px] w-[500px] rounded-full bg-primary/[0.06] blur-[120px]" />
        </div>
        <div className="absolute inset-0 -z-20 bg-grit-grid mask-fade-center" />
        <span className="crosshair absolute top-20 left-[15%] text-foreground/20 hidden md:block" style={{ width: 14, height: 14 }} />
        <span className="crosshair absolute top-32 right-[17%] text-primary/30 hidden md:block" style={{ width: 14, height: 14 }} />
        <div className="container max-w-screen-xl py-20 md:py-32 px-6">
          <div className="max-w-3xl mx-auto text-center">
            <div className="flex items-center justify-center gap-2 rounded-full border border-border/60 bg-accent/40 px-4 py-1.5 mb-6 mx-auto w-fit">
              <Sparkles className="h-3.5 w-3.5 text-primary" />
              <span className="text-xs font-mono uppercase tracking-wider text-muted-foreground">
                Built with Grit
              </span>
            </div>
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight mb-6">
              <span className="gradient-text">Showcase</span>
            </h1>
            <p className="text-lg md:text-xl text-muted-foreground leading-relaxed mb-8 max-w-2xl mx-auto">
              Real-world applications built with the Grit framework.
              From SaaS platforms to content management systems &mdash; see
              what&apos;s possible with Go + React.
            </p>
            <div className="flex items-center justify-center gap-3">
              <Button
                size="sm"
                className="h-9 px-4 text-sm bg-primary/90 hover:bg-primary"
                asChild
              >
                <Link href="/docs/getting-started/quick-start">
                  Start Building
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
              <Button
                variant="outline"
                size="sm"
                className="h-9 px-4 text-sm"
                asChild
              >
                <Link
                  href="https://github.com/MUKE-coder/grit/issues"
                  target="_blank"
                  rel="noreferrer"
                >
                  Submit Your Project
                  <ExternalLink className="ml-2 h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* Featured project */}
      {featured.map((project) => (
        <section key={project.name} className="container max-w-screen-xl px-6 pb-16">
          <div className="rounded-2xl border border-border/60 bg-card/50 overflow-hidden">
            {/* Screenshot */}
            <div className="relative aspect-[16/8] bg-accent/30 border-b border-border/40 overflow-hidden group">
              <img
                src={project.image}
                alt={`${project.name} screenshot`}
                className="w-full h-full object-cover object-top transition-transform duration-500 group-hover:scale-[1.02]"
              />
              <div className="absolute inset-0 bg-gradient-to-t from-background/80 via-transparent to-transparent" />
              <div className="absolute bottom-6 left-6 right-6 flex items-end justify-between">
                <div>
                  <div className="flex items-center gap-3 mb-2">
                    <h2 className="text-2xl md:text-3xl font-bold text-foreground">
                      {project.name}
                    </h2>
                    <span className="inline-flex items-center rounded-full bg-primary/20 border border-primary/30 px-2.5 py-0.5 text-[10px] font-mono font-medium text-primary uppercase tracking-wider">
                      Featured
                    </span>
                  </div>
                  <p className="text-sm text-foreground/70">
                    {project.description}
                  </p>
                </div>
                <Button
                  size="sm"
                  className="hidden md:inline-flex h-8 px-3 text-xs bg-primary/90 hover:bg-primary shrink-0"
                  asChild
                >
                  <Link href={project.url} target="_blank" rel="noreferrer">
                    Visit Site
                    <ExternalLink className="ml-1.5 h-3 w-3" />
                  </Link>
                </Button>
              </div>
            </div>

            {/* Details */}
            <div className="p-6 md:p-8">
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Left */}
                <div className="lg:col-span-2 space-y-6">
                  <p className="text-muted-foreground leading-relaxed">
                    {project.longDescription}
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {project.tags.map((tag) => (
                      <span
                        key={tag}
                        className="inline-flex items-center rounded-md bg-accent/60 border border-border/40 px-2.5 py-1 text-xs font-medium text-muted-foreground"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                  <div>
                    <h3 className="text-sm font-semibold text-foreground mb-3">
                      Modules &amp; Capabilities
                    </h3>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                      {project.stats.highlights.map((item) => (
                        <div
                          key={item}
                          className="flex items-center gap-2 text-sm text-muted-foreground"
                        >
                          <div className="h-1.5 w-1.5 rounded-full bg-primary/60 shrink-0" />
                          {item}
                        </div>
                      ))}
                    </div>
                  </div>
                  <div className="md:hidden">
                    <Button
                      size="sm"
                      className="h-9 px-4 text-sm bg-primary/90 hover:bg-primary w-full"
                      asChild
                    >
                      <Link href={project.url} target="_blank" rel="noreferrer">
                        Visit {project.name}
                        <ExternalLink className="ml-1.5 h-3.5 w-3.5" />
                      </Link>
                    </Button>
                  </div>
                </div>

                {/* Right: Stats */}
                <div className="space-y-4">
                  <div className="rounded-xl border border-border/60 bg-accent/30 p-5">
                    <h3 className="text-sm font-semibold text-foreground mb-4">
                      Project Scale
                    </h3>
                    <div className="space-y-4">
                      <div className="flex items-center gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/20">
                          <Database className="h-5 w-5 text-primary" />
                        </div>
                        <div>
                          <div className="text-2xl font-bold text-foreground">{project.stats.tables}</div>
                          <div className="text-xs text-muted-foreground">Database Tables</div>
                        </div>
                      </div>
                      <div className="h-px bg-border/40" />
                      <div className="flex items-center gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/20">
                          <Box className="h-5 w-5 text-primary" />
                        </div>
                        <div>
                          <div className="text-2xl font-bold text-foreground">{project.stats.models}</div>
                          <div className="text-xs text-muted-foreground">Registered Models</div>
                        </div>
                      </div>
                      {project.stats.goroutines && (
                        <>
                          <div className="h-px bg-border/40" />
                          <div className="flex items-center gap-3">
                            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-amber-500/10 border border-amber-500/20">
                              <Zap className="h-5 w-5 text-amber-500" />
                            </div>
                            <div>
                              <div className="text-2xl font-bold text-foreground">{project.stats.goroutines}</div>
                              <div className="text-xs text-muted-foreground">Active Goroutines</div>
                            </div>
                          </div>
                        </>
                      )}
                      <div className="h-px bg-border/40" />
                      <div className="flex items-center gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/20">
                          <Layers className="h-5 w-5 text-primary" />
                        </div>
                        <div>
                          <div className="text-2xl font-bold text-foreground">{project.stats.modules}</div>
                          <div className="text-xs text-muted-foreground">Feature Modules</div>
                        </div>
                      </div>
                      <div className="h-px bg-border/40" />
                      <div className="flex items-center gap-3">
                        <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-emerald-500/10 border border-emerald-500/20">
                          <Globe className="h-5 w-5 text-emerald-500" />
                        </div>
                        <div>
                          <div className="text-sm font-semibold text-foreground">Live in Production</div>
                          <div className="text-xs text-muted-foreground">
                            <Link href={project.url} target="_blank" rel="noreferrer" className="text-primary hover:underline">
                              {project.url.replace("https://", "")}
                            </Link>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="rounded-xl border border-border/60 bg-accent/30 p-5">
                    <h3 className="text-sm font-semibold text-foreground mb-2">Tech Stack</h3>
                    <div className="flex flex-wrap gap-1.5">
                      {project.techStack.map((tech) => (
                        <span
                          key={tech}
                          className="inline-flex items-center rounded-md bg-background/60 border border-border/40 px-2 py-0.5 text-[11px] font-mono text-muted-foreground"
                        >
                          {tech}
                        </span>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      ))}

      {/* Community projects grid */}
      {rest.length > 0 && (
        <section className="container max-w-screen-xl px-6 pb-16">
          <h2 className="text-2xl font-bold text-foreground mb-8">
            Community Projects
          </h2>
          <ShowcaseGrid projects={rest} />
        </section>
      )}

      {/* CTA */}
      <section className="container max-w-screen-xl px-6 pb-24">
        <div className="rounded-2xl border border-border/60 bg-accent/20 p-8 md:p-12 text-center">
          <h2 className="text-2xl md:text-3xl font-bold tracking-tight mb-3">
            <span className="gradient-text">Built something with Grit?</span>
          </h2>
          <p className="text-muted-foreground mb-6 max-w-lg mx-auto">
            We&apos;d love to feature your project. Open an issue on GitHub or
            reach out to get your app added to the showcase.
          </p>
          <div className="flex items-center justify-center gap-3">
            <Button
              size="sm"
              className="h-9 px-4 text-sm bg-primary/90 hover:bg-primary"
              asChild
            >
              <Link
                href="https://github.com/MUKE-coder/grit/issues"
                target="_blank"
                rel="noreferrer"
              >
                Submit on GitHub
                <ExternalLink className="ml-2 h-3.5 w-3.5" />
              </Link>
            </Button>
            <Button
              variant="outline"
              size="sm"
              className="h-9 px-4 text-sm"
              asChild
            >
              <Link href="/docs/getting-started/quick-start">
                Get Started
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </section>
    </div>
  );
}
