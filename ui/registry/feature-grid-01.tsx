import { Zap, Shield, Layers, Database, Bot, GitBranch } from "lucide-react";

const FEATURES = [
  {
    icon: Zap,
    title: "One command",
    body: "Scaffold the API, the admin panel and the frontend together — wired, typed and ready to run.",
  },
  {
    icon: Shield,
    title: "Hardened by default",
    body: "CSRF, strict CSP, rate limiting and an audit log ship in the scaffold, not in a checklist.",
  },
  {
    icon: Layers,
    title: "Five architectures",
    body: "Embedded SPA, monorepo, API-only, mobile or desktop — the same generators throughout.",
  },
  {
    icon: Database,
    title: "Typed end to end",
    body: "Go models generate TypeScript types and Zod schemas, so the two halves cannot drift.",
  },
  {
    icon: Bot,
    title: "Agent friendly",
    body: "One way to do each thing means an AI assistant never has to guess which pattern you chose.",
  },
  {
    icon: GitBranch,
    title: "Code you own",
    body: "Everything lands in your repo as ordinary Go and TypeScript. No runtime magic to fight.",
  },
];

export default function FeatureGrid01() {
  return (
    <section className="bg-background py-24">
      <div className="container mx-auto px-6">
        <div className="max-w-2xl mb-14">
          <span className="text-xs font-mono uppercase tracking-wider text-accent">
            Features
          </span>
          <h2 className="mt-3 text-3xl md:text-4xl font-bold text-foreground tracking-tight leading-tight">
            Everything the boring part of your app needs
          </h2>
          <p className="mt-4 text-base text-text-secondary leading-relaxed">
            The work that produces no product differentiation, already done — so the
            first week goes into the thing people actually pay for.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {FEATURES.map(({ icon: Icon, title, body }) => (
            <div
              key={title}
              className="group rounded-2xl border border-border bg-bg-secondary p-6 hover:border-accent/40 hover:bg-bg-elevated transition-colors"
            >
              <span className="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-accent/10 text-accent mb-4">
                <Icon size={18} />
              </span>
              <h3 className="text-sm font-semibold text-foreground mb-2">{title}</h3>
              <p className="text-sm text-text-muted leading-relaxed">{body}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
