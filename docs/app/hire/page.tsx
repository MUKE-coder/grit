import Link from "next/link";
import type { Metadata } from "next";
import {
  ArrowRight,
  Rocket,
  Code2,
  Server,
  Shield,
  Clock,
  Users,
  CheckCircle2,
  Mail,
  ExternalLink,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { GridFrame } from "@/components/grid-frame";

export const metadata: Metadata = {
  title: "Hire Grit Developers | Grit",
  description:
    "Hire experienced Go + React developers to build your next production application with the Grit framework. Fast delivery, full-stack expertise, battle-tested stack.",
};

export default function HirePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      {/* Hero */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 -z-10">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[500px] w-[500px] rounded-full bg-primary/[0.06] blur-[120px]" />
        </div>
        <div className="absolute inset-0 -z-20 bg-grit-grid mask-fade-center" />
        <span className="crosshair absolute top-20 left-[15%] text-foreground/20 hidden md:block" style={{ width: 14, height: 14 }} />
        <span className="crosshair absolute top-32 right-[17%] text-primary/30 hidden md:block" style={{ width: 14, height: 14 }} />
        <div className="max-w-6xl mx-auto py-24 md:py-36 px-6">
          <div className="max-w-3xl mx-auto text-center">
            <div className="flex items-center justify-center gap-2 rounded-full border border-border/60 bg-accent/40 px-4 py-1.5 mb-6 mx-auto w-fit">
              <Users className="h-3.5 w-3.5 text-primary" />
              <span className="text-xs font-mono uppercase tracking-wider text-muted-foreground">
                Professional Development Services
              </span>
            </div>
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight mb-6">
              <span className="gradient-text">Hire Grit Developers</span>
            </h1>
            <p className="text-lg md:text-xl text-muted-foreground leading-relaxed mb-8 max-w-2xl mx-auto">
              Need a production-ready application built with Go + React? Our
              team knows the Grit framework inside and out &mdash; because we
              built it. Get a full-stack app shipped fast, with the same
              quality we put into the framework itself.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-3">
              <Button size="lg" className="glow-purple-sm" asChild>
                <Link href="mailto:desishub.info@gmail.com">
                  <Mail className="mr-2 h-4 w-4" />
                  Get in Touch
                </Link>
              </Button>
              <Button variant="outline" size="lg" className="border-border/60 bg-transparent hover:bg-accent/50" asChild>
                <Link href="https://gritcms.com" target="_blank" rel="noreferrer">
                  See GritCMS
                  <ExternalLink className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

      {/* What We Build */}
      <section className="border-t border-border/30">
        <div className="max-w-6xl mx-auto py-20 md:py-28 px-6">
          <div className="text-center mb-14">
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">
              What We Build
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              From SaaS products to internal tools, we deliver full-stack
              applications built on a battle-tested stack.
            </p>
          </div>
          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3 max-w-5xl mx-auto">
            {[
              {
                icon: <Rocket className="h-5 w-5" />,
                title: "SaaS Applications",
                desc: "Multi-tenant SaaS platforms with authentication, billing, dashboards, and admin panels.",
              },
              {
                icon: <Server className="h-5 w-5" />,
                title: "APIs & Backends",
                desc: "High-performance Go APIs with RBAC, caching, background jobs, file storage, and email.",
              },
              {
                icon: <Shield className="h-5 w-5" />,
                title: "Admin Panels",
                desc: "Rich admin dashboards with data tables, forms, charts, and role-based access.",
              },
              {
                icon: <Code2 className="h-5 w-5" />,
                title: "Internal Tools",
                desc: "Custom business tools, CRMs, inventory systems, and workflow automation apps.",
              },
              {
                icon: <Users className="h-5 w-5" />,
                title: "Community Platforms",
                desc: "Creator platforms, course sites, forums, and membership apps with payments.",
              },
              {
                icon: <Clock className="h-5 w-5" />,
                title: "MVP & Prototypes",
                desc: "Go from idea to working product fast. Grit's code generation means rapid delivery.",
              },
            ].map((item) => (
              <div
                key={item.title}
                className="card-grit group rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/40 hover:bg-card/80"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 border border-primary/15 text-primary mb-4">
                  {item.icon}
                </div>
                <h3 className="text-base font-semibold mb-2">{item.title}</h3>
                <p className="text-sm text-muted-foreground/70 leading-relaxed">
                  {item.desc}
                </p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Why Hire Us */}
      <section className="border-t border-border/30 bg-accent/20">
        <div className="max-w-6xl mx-auto py-20 md:py-28 px-6">
          <div className="text-center mb-14">
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">
              Why Hire the Grit Team
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              We don&apos;t just use Grit &mdash; we created it. That means
              deep expertise, faster delivery, and code that follows every
              best practice.
            </p>
          </div>
          <div className="grid gap-4 sm:grid-cols-2 max-w-3xl mx-auto">
            {[
              "We built the framework — nobody knows it better",
              "Full-stack: Go API + React frontend + admin panel",
              "Code generation means faster delivery and lower cost",
              "Battle-tested stack: Gin, GORM, Next.js, PostgreSQL, Redis",
              "Production-ready from day one: Docker, CI/CD, monitoring",
              "Ongoing support and maintenance available",
              "Clean, well-documented, maintainable code",
              "Your code, your servers — no vendor lock-in",
            ].map((point) => (
              <div
                key={point}
                className="flex items-start gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-3"
              >
                <CheckCircle2 className="h-4 w-4 text-primary mt-0.5 shrink-0" />
                <span className="text-sm text-foreground/80">{point}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Tech Stack */}
      <section className="border-t border-border/30">
        <div className="max-w-6xl mx-auto py-20 md:py-28 px-6">
          <div className="text-center mb-14">
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">
              The Stack
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Every project is built on the same proven Grit technology stack.
            </p>
          </div>
          <div className="grid gap-3 grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 max-w-4xl mx-auto">
            {[
              { name: "Go", role: "Backend" },
              { name: "Gin", role: "Web Framework" },
              { name: "GORM", role: "ORM" },
              { name: "PostgreSQL", role: "Database" },
              { name: "Redis", role: "Cache & Queues" },
              { name: "Next.js", role: "Frontend" },
              { name: "React", role: "UI Library" },
              { name: "TypeScript", role: "Type Safety" },
              { name: "Tailwind CSS", role: "Styling" },
              { name: "shadcn/ui", role: "Components" },
              { name: "Docker", role: "Deployment" },
              { name: "TanStack Query", role: "Data Fetching" },
            ].map((tech) => (
              <div
                key={tech.name}
                className="rounded-lg border border-border/30 bg-card/30 px-4 py-3 text-center"
              >
                <span className="text-sm font-semibold block">
                  {tech.name}
                </span>
                <span className="text-xs text-muted-foreground/50">
                  {tech.role}
                </span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="border-t border-border/30">
        <div className="max-w-6xl mx-auto py-20 md:py-28 px-6">
          <div className="relative rounded-2xl border border-primary/20 bg-primary/[0.04] overflow-hidden">
            <div className="absolute inset-0 -z-10">
              <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[300px] w-[300px] rounded-full bg-primary/[0.08] blur-[100px]" />
            </div>
            <div className="text-center py-16 md:py-24 px-6">
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">
                Ready to Build?
              </h2>
              <p className="text-lg text-muted-foreground max-w-xl mx-auto mb-8">
                Tell us about your project and we&apos;ll get back to you
                within 24 hours with a plan and estimate.
              </p>
              <div className="flex flex-col sm:flex-row items-center justify-center gap-3">
                <Button size="lg" className="glow-purple-sm" asChild>
                  <Link href="mailto:desishub.info@gmail.com">
                    <Mail className="mr-2 h-4 w-4" />
                    desishub.info@gmail.com
                  </Link>
                </Button>
              </div>
              <p className="text-xs text-muted-foreground/40 mt-4">
                Or reach out on{" "}
                <Link
                  href="https://www.linkedin.com/company/grit-framework"
                  target="_blank"
                  rel="noreferrer"
                  className="text-primary/70 hover:underline"
                >
                  LinkedIn
                </Link>{" "}
                or{" "}
                <Link
                  href="https://github.com/MUKE-coder"
                  target="_blank"
                  rel="noreferrer"
                  className="text-primary/70 hover:underline"
                >
                  GitHub
                </Link>
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border/30 py-8 px-6">
        <div className="max-w-6xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-muted-foreground/50">
          <span>Grit Framework &mdash; Go + React. Built with Grit.</span>
          <div className="flex items-center gap-4">
            <Link href="/docs" className="hover:text-foreground transition-colors">
              Docs
            </Link>
            <Link href="https://gritcms.com" target="_blank" rel="noreferrer" className="hover:text-foreground transition-colors">
              GritCMS
            </Link>
            <Link href="/donate" className="hover:text-foreground transition-colors">
              Support
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
