package scaffold

// Extra components to reach 100 total in the Grit UI registry.
// 7 marketing + 2 layout = 9 new components (91 existing + 9 = 100).

// ── Marketing ──────────────────────────────────────────────────────────────────

func mktSocialProofBar01() string {
	return `// social-proof-bar-01 — Avatar stack + trust bar (common SaaS landing pattern)
export function SocialProofBar({
  count = "50,000+",
  label = "developers ship with Grit",
  rating = "4.9",
  reviewSource = "Product Hunt",
  avatars = ["SC", "JW", "PM", "AR", "TK"],
}: {
  count?: string; label?: string; rating?: string;
  reviewSource?: string; avatars?: string[];
}) {
  return (
    <div className="flex flex-col sm:flex-row items-center gap-4 sm:gap-6">
      <div className="flex items-center gap-3">
        <div className="flex -space-x-2">
          {avatars.slice(0, 5).map((initials, i) => (
            <div
              key={i}
              className="w-8 h-8 rounded-full border-2 border-background bg-accent/20 flex items-center justify-center flex-shrink-0"
              style={{ zIndex: avatars.length - i }}
            >
              <span className="text-xs font-semibold text-accent">{initials}</span>
            </div>
          ))}
        </div>
        <p className="text-sm text-text-secondary">
          <span className="font-semibold text-foreground">{count}</span> {label}
        </p>
      </div>
      <div className="hidden sm:block w-px h-6 bg-border" />
      <div className="flex items-center gap-1.5">
        <div className="flex items-center gap-0.5">
          {Array.from({ length: 5 }).map((_, i) => (
            <svg key={i} className="w-3.5 h-3.5 text-yellow-400" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
            </svg>
          ))}
        </div>
        <span className="text-sm font-semibold text-foreground">{rating}</span>
        <span className="text-sm text-text-muted">on {reviewSource}</span>
      </div>
    </div>
  );
}
`
}

func mktFeatureComparison01() string {
	return `// feature-comparison-01 — Feature comparison table (Linear pricing page-inspired)
export function FeatureComparison({
  plans = ["Starter", "Pro", "Enterprise"],
  categories = [
    {
      name: "Core",
      features: [
        { label: "Projects", values: ["5", "Unlimited", "Unlimited"] },
        { label: "Team members", values: ["3", "15", "Unlimited"] },
        { label: "Storage", values: ["5 GB", "100 GB", "1 TB"] },
      ],
    },
    {
      name: "Features",
      features: [
        { label: "Analytics", values: [false, true, true] },
        { label: "Custom domain", values: [false, true, true] },
        { label: "Priority support", values: [false, false, true] },
        { label: "SSO / SAML", values: [false, false, true] },
        { label: "Audit log", values: [false, false, true] },
      ],
    },
  ],
}: {
  plans?: string[];
  categories?: Array<{
    name: string;
    features: Array<{ label: string; values: (string | boolean)[] }>;
  }>;
}) {
  const Cell = ({ v }: { v: string | boolean }) => {
    if (typeof v === "boolean") {
      return v ? (
        <svg className="w-5 h-5 text-success mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
        </svg>
      ) : (
        <svg className="w-4 h-4 text-text-muted mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
        </svg>
      );
    }
    return <span className="text-sm text-foreground font-medium">{v}</span>;
  };
  return (
    <div className="w-full overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-border">
            <th className="text-left py-4 pr-4 text-sm font-medium text-text-muted w-48">Features</th>
            {plans.map((plan) => (
              <th key={plan} className="text-center py-4 px-4 text-sm font-semibold text-foreground">{plan}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {categories.map((cat) => (
            <>
              <tr key={cat.name}>
                <td colSpan={plans.length + 1} className="pt-6 pb-2">
                  <span className="text-xs font-semibold text-text-muted uppercase tracking-wider">{cat.name}</span>
                </td>
              </tr>
              {cat.features.map((feat) => (
                <tr key={feat.label} className="border-b border-border/50 hover:bg-bg-hover transition-colors">
                  <td className="py-3 pr-4 text-sm text-text-secondary">{feat.label}</td>
                  {feat.values.map((v, i) => (
                    <td key={i} className="py-3 px-4 text-center"><Cell v={v} /></td>
                  ))}
                </tr>
              ))}
            </>
          ))}
        </tbody>
      </table>
    </div>
  );
}
`
}

func mktAnnouncementBar01() string {
	return `// announcement-bar-01 — Dismissible top-of-page announcement banner (Vercel-style)
"use client";
import { useState } from "react";
export function AnnouncementBar({
  message = "Grit v1.0 is here",
  linkLabel = "Read the announcement",
  linkHref = "/blog/v1",
  badge,
}: {
  message?: string; linkLabel?: string; linkHref?: string; badge?: string;
}) {
  const [visible, setVisible] = useState(true);
  if (!visible) return null;
  return (
    <div className="bg-accent/10 border-b border-accent/20 px-4 py-2.5 flex items-center justify-center gap-3 relative">
      {badge && (
        <span className="text-xs font-bold bg-accent text-white px-2 py-0.5 rounded-full">{badge}</span>
      )}
      <p className="text-sm text-text-secondary text-center">
        {message}{" "}
        <a href={linkHref} className="text-accent font-semibold hover:underline ml-1">
          {linkLabel} &rarr;
        </a>
      </p>
      <button
        onClick={() => setVisible(false)}
        className="absolute right-4 p-1 rounded text-text-muted hover:text-foreground transition-colors"
        aria-label="Dismiss"
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  );
}
`
}

func mktWaitlistForm01() string {
	return `// waitlist-form-01 — Waitlist signup with position counter (common pre-launch pattern)
"use client";
import { useState } from "react";
export function WaitlistForm({
  title = "Join the waitlist",
  description = "Be the first to know when we launch. We'll send you an invite.",
  currentCount = 2_847,
  onJoin,
}: {
  title?: string; description?: string; currentCount?: number;
  onJoin?: (email: string) => Promise<{ position: number }>;
}) {
  const [email, setEmail] = useState("");
  const [state, setState] = useState<"idle" | "loading" | "done">("idle");
  const [position, setPosition] = useState(0);
  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || state !== "idle") return;
    setState("loading");
    try {
      const result = await onJoin?.(email);
      setPosition(result?.position ?? currentCount + 1);
      setState("done");
    } catch {
      setState("idle");
    }
  };
  return (
    <div className="bg-bg-elevated border border-border rounded-2xl p-6 flex flex-col gap-5 w-full max-w-md">
      {state === "done" ? (
        <div className="flex flex-col items-center gap-3 py-4 text-center">
          <div className="w-14 h-14 rounded-full bg-success/20 flex items-center justify-center">
            <svg className="w-7 h-7 text-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <h3 className="text-base font-semibold text-foreground">You're #{position} on the list!</h3>
          <p className="text-sm text-text-secondary">We'll email you at <strong>{email}</strong> when your spot is ready.</p>
        </div>
      ) : (
        <>
          <div>
            <h3 className="text-base font-semibold text-foreground">{title}</h3>
            <p className="text-sm text-text-secondary mt-1">{description}</p>
          </div>
          <div className="flex items-center gap-2 text-xs text-text-muted bg-background border border-border rounded-lg px-3 py-2">
            <div className="flex -space-x-1">
              {["SC", "JW", "PM"].map((i) => (
                <div key={i} className="w-5 h-5 rounded-full bg-accent/30 border border-border flex items-center justify-center">
                  <span className="text-xs font-bold text-accent" style={{ fontSize: 8 }}>{i[0]}</span>
                </div>
              ))}
            </div>
            <span><strong className="text-foreground">{currentCount.toLocaleString()}</strong> people already waiting</span>
          </div>
          <form onSubmit={submit} className="flex gap-2">
            <input
              type="email" value={email} onChange={(e) => setEmail(e.target.value)}
              placeholder="you@example.com" required
              className="flex-1 bg-background border border-border rounded-lg px-3 py-2.5 text-sm text-foreground placeholder-text-muted outline-none focus:border-accent transition-colors"
            />
            <button
              type="submit" disabled={state === "loading"}
              className="px-4 py-2.5 rounded-lg bg-accent hover:bg-accent/90 text-white text-sm font-semibold transition-colors disabled:opacity-60 flex-shrink-0"
            >
              {state === "loading" ? "..." : "Join"}
            </button>
          </form>
          <p className="text-xs text-text-muted text-center">No spam. Unsubscribe anytime.</p>
        </>
      )}
    </div>
  );
}
`
}

func mktCaseStudyCard01() string {
	return `// case-study-card-01 — Customer case study card with metric callout (Stripe/Vercel-style)
export function CaseStudyCard({
  companyName = "Acme Corp",
  companyLogo = "A",
  logoColor = "bg-blue-600",
  industry = "E-commerce",
  metric = "3x",
  metricLabel = "faster deployment",
  quote = "Grit cut our time-to-market in half. We shipped our MVP in 2 weeks instead of 2 months.",
  authorName = "Sarah Chen",
  authorRole = "CTO",
  href = "#",
}: {
  companyName?: string; companyLogo?: string; logoColor?: string;
  industry?: string; metric?: string; metricLabel?: string;
  quote?: string; authorName?: string; authorRole?: string; href?: string;
}) {
  return (
    <a
      href={href}
      className="group block bg-bg-elevated border border-border rounded-2xl p-6 hover:border-accent/40 transition-all flex flex-col gap-5"
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-3">
          <div className={"w-10 h-10 rounded-xl " + logoColor + " flex items-center justify-center flex-shrink-0"}>
            <span className="text-white font-bold text-lg">{companyLogo}</span>
          </div>
          <div>
            <p className="text-sm font-semibold text-foreground">{companyName}</p>
            <p className="text-xs text-text-muted">{industry}</p>
          </div>
        </div>
        <svg className="w-4 h-4 text-text-muted group-hover:text-accent group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-all" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
      </div>
      <div className="flex items-baseline gap-2">
        <span className="text-4xl font-black text-accent">{metric}</span>
        <span className="text-sm text-text-secondary">{metricLabel}</span>
      </div>
      <blockquote className="text-sm text-text-secondary leading-relaxed italic border-l-2 border-accent/40 pl-3">
        "{quote}"
      </blockquote>
      <p className="text-xs text-text-muted mt-auto">
        — <span className="text-text-secondary font-medium">{authorName}</span>, {authorRole}
      </p>
    </a>
  );
}
`
}

func mktRoadmapItem01() string {
	return `// roadmap-item-01 — Public roadmap item card with status and upvote (Canny/Linear-style)
"use client";
import { useState } from "react";
export function RoadmapItem({
  title = "AI-powered code review",
  description = "Automatically review pull requests using Claude to catch bugs, suggest improvements, and enforce coding standards.",
  status = "planned",
  votes = 142,
  tags = ["AI", "Developer Tools"],
  eta,
}: {
  title?: string; description?: string;
  status?: "planned" | "in-progress" | "done" | "considering";
  votes?: number; tags?: string[]; eta?: string;
}) {
  const [upvoted, setUpvoted] = useState(false);
  const [count, setCount] = useState(votes);
  const toggle = () => {
    setUpvoted(!upvoted);
    setCount(upvoted ? count - 1 : count + 1);
  };
  const statusConfig = {
    "planned": { label: "Planned", classes: "text-accent bg-accent/10 border-accent/20" },
    "in-progress": { label: "In Progress", classes: "text-yellow-400 bg-yellow-500/10 border-yellow-500/20" },
    "done": { label: "Shipped", classes: "text-success bg-success/10 border-success/20" },
    "considering": { label: "Considering", classes: "text-text-muted bg-bg-hover border-border" },
  };
  const s = statusConfig[status] ?? statusConfig["planned"];
  return (
    <div className="bg-bg-elevated border border-border rounded-xl p-5 flex items-start gap-4 hover:border-accent/30 transition-colors">
      <button
        onClick={toggle}
        className={"flex flex-col items-center gap-1 p-2 rounded-lg border transition-all min-w-[3rem] " + (upvoted ? "border-accent/30 bg-accent/10 text-accent" : "border-border text-text-muted hover:text-foreground hover:border-accent/30")}
      >
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
        </svg>
        <span className="text-xs font-bold">{count}</span>
      </button>
      <div className="flex-1 min-w-0">
        <div className="flex items-start justify-between gap-3 flex-wrap">
          <h3 className="text-sm font-semibold text-foreground">{title}</h3>
          <div className="flex items-center gap-2 flex-shrink-0">
            {eta && <span className="text-xs text-text-muted">ETA: {eta}</span>}
            <span className={"text-xs border px-2 py-0.5 rounded-full font-medium " + s.classes}>{s.label}</span>
          </div>
        </div>
        <p className="text-sm text-text-secondary mt-1 leading-relaxed">{description}</p>
        {tags.length > 0 && (
          <div className="flex items-center gap-1.5 mt-3 flex-wrap">
            {tags.map((tag) => (
              <span key={tag} className="text-xs bg-background border border-border text-text-muted px-2 py-0.5 rounded">{tag}</span>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
`
}

func mktPartnerCard01() string {
	return `// partner-card-01 — Technology partner / integration showcase card
export function PartnerCard({
  name = "Vercel",
  description = "Deploy your Grit web and admin apps to Vercel with zero configuration.",
  logoChar = "V",
  logoColor = "bg-black",
  category = "Deployment",
  href = "#",
  isOfficial = false,
}: {
  name?: string; description?: string; logoChar?: string;
  logoColor?: string; category?: string; href?: string; isOfficial?: boolean;
}) {
  return (
    <a
      href={href}
      target="_blank"
      rel="noopener noreferrer"
      className="group bg-bg-elevated border border-border rounded-2xl p-5 flex flex-col gap-4 hover:border-accent/40 hover:bg-bg-hover transition-all"
    >
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-3">
          <div className={"w-10 h-10 rounded-xl " + logoColor + " flex items-center justify-center border border-white/10 flex-shrink-0"}>
            <span className="text-white font-bold text-base">{logoChar}</span>
          </div>
          <div>
            <div className="flex items-center gap-1.5">
              <span className="text-sm font-semibold text-foreground">{name}</span>
              {isOfficial && (
                <svg className="w-3.5 h-3.5 text-accent" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              )}
            </div>
            <span className="text-xs text-text-muted">{category}</span>
          </div>
        </div>
        <svg className="w-4 h-4 text-text-muted group-hover:text-accent transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
        </svg>
      </div>
      <p className="text-sm text-text-secondary leading-relaxed">{description}</p>
      <span className="text-xs text-accent font-medium group-hover:underline">Learn more &rarr;</span>
    </a>
  );
}
`
}

// ── Layout ──────────────────────────────────────────────────────────────────────

func layoutSettingsNav01() string {
	return `// settings-nav-01 — Settings sidebar with grouped navigation items (Vercel/GitHub-style)
export function SettingsNav({
  groups = [
    {
      label: "Account",
      items: [
        { label: "Profile", href: "/settings/profile", active: true },
        { label: "Security", href: "/settings/security" },
        { label: "Notifications", href: "/settings/notifications" },
      ],
    },
    {
      label: "Workspace",
      items: [
        { label: "General", href: "/settings/workspace" },
        { label: "Members", href: "/settings/members" },
        { label: "Billing", href: "/settings/billing" },
        { label: "Integrations", href: "/settings/integrations" },
        { label: "API Keys", href: "/settings/api-keys" },
      ],
    },
    {
      label: "Danger Zone",
      items: [
        { label: "Delete Workspace", href: "/settings/delete", danger: true },
      ],
    },
  ],
}: {
  groups?: Array<{
    label: string;
    items: Array<{ label: string; href: string; active?: boolean; danger?: boolean }>;
  }>;
}) {
  return (
    <nav className="flex flex-col gap-6 w-48">
      {groups.map((group) => (
        <div key={group.label}>
          <p className="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 px-3">{group.label}</p>
          <div className="flex flex-col gap-0.5">
            {group.items.map((item) => (
              <a
                key={item.href}
                href={item.href}
                className={"flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-colors " + (item.active ? "bg-bg-elevated text-foreground" : item.danger ? "text-danger hover:bg-danger/10" : "text-text-secondary hover:text-foreground hover:bg-bg-hover")}
              >
                {item.label}
              </a>
            ))}
          </div>
        </div>
      ))}
    </nav>
  );
}
`
}

func layoutCommandPaletteShell01() string {
	return `// command-palette-shell-01 — Full command palette modal with search + groups (Linear-style)
"use client";
import { useEffect, useState, useRef } from "react";
export function CommandPaletteShell({
  groups = [
    {
      label: "Recent",
      items: [
        { id: "1", label: "GRT-42 — Implement OAuth login flow", meta: "Issue" },
        { id: "2", label: "main branch", meta: "Branch" },
      ],
    },
    {
      label: "Actions",
      items: [
        { id: "3", label: "Create new issue", shortcut: ["C"], meta: "Action" },
        { id: "4", label: "Go to settings", shortcut: ["G", "S"], meta: "Navigate" },
        { id: "5", label: "Switch workspace", shortcut: ["O"], meta: "Navigate" },
      ],
    },
  ],
  placeholder = "Search or type a command...",
  onClose,
  onSelect,
}: {
  groups?: Array<{ label: string; items: Array<{ id: string; label: string; shortcut?: string[]; meta?: string }> }>;
  placeholder?: string; onClose?: () => void;
  onSelect?: (id: string) => void;
}) {
  const [query, setQuery] = useState("");
  const [activeId, setActiveId] = useState(groups[0]?.items[0]?.id ?? "");
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    inputRef.current?.focus();
    const handler = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose?.();
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [onClose]);

  const allItems = groups.flatMap((g) => g.items);
  const filtered = query
    ? allItems.filter((i) => i.label.toLowerCase().includes(query.toLowerCase()))
    : allItems;

  return (
    <div className="fixed inset-0 z-50 flex items-start justify-center pt-[15vh] px-4">
      <div className="absolute inset-0 bg-background/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-lg bg-bg-elevated border border-border rounded-2xl shadow-2xl overflow-hidden">
        {/* Search input */}
        <div className="flex items-center gap-3 px-4 py-3.5 border-b border-border">
          <svg className="w-5 h-5 text-text-muted flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            ref={inputRef}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={placeholder}
            className="flex-1 bg-transparent text-sm text-foreground placeholder-text-muted outline-none"
          />
          <kbd className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">ESC</kbd>
        </div>
        {/* Results */}
        <div className="max-h-72 overflow-y-auto py-2">
          {query ? (
            filtered.length === 0 ? (
              <p className="text-sm text-text-muted text-center py-8">No results for "{query}"</p>
            ) : (
              <div className="px-2">
                {filtered.map((item) => (
                  <button
                    key={item.id}
                    onClick={() => { onSelect?.(item.id); onClose?.(); }}
                    className={"w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-left transition-colors " + (item.id === activeId ? "bg-accent/10" : "hover:bg-bg-hover")}
                  >
                    <span className="text-sm text-foreground flex-1 truncate">{item.label}</span>
                    {item.meta && <span className="text-xs text-text-muted">{item.meta}</span>}
                    {item.shortcut && (
                      <div className="flex gap-1">
                        {item.shortcut.map((k) => (
                          <kbd key={k} className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">{k}</kbd>
                        ))}
                      </div>
                    )}
                  </button>
                ))}
              </div>
            )
          ) : (
            groups.map((group) => (
              <div key={group.label} className="mb-2">
                <p className="text-xs font-semibold text-text-muted uppercase tracking-wider px-5 py-1.5">{group.label}</p>
                <div className="px-2">
                  {group.items.map((item) => (
                    <button
                      key={item.id}
                      onClick={() => { onSelect?.(item.id); onClose?.(); }}
                      className={"w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors " + (item.id === activeId ? "bg-accent/10" : "hover:bg-bg-hover")}
                    >
                      <span className="text-sm text-foreground flex-1 truncate">{item.label}</span>
                      {item.meta && <span className="text-xs text-text-muted">{item.meta}</span>}
                      {item.shortcut && (
                        <div className="flex gap-1">
                          {item.shortcut.map((k) => (
                            <kbd key={k} className="text-xs bg-background border border-border text-text-muted px-1.5 py-0.5 rounded font-mono">{k}</kbd>
                          ))}
                        </div>
                      )}
                    </button>
                  ))}
                </div>
              </div>
            ))
          )}
        </div>
        {/* Footer hint */}
        <div className="border-t border-border px-4 py-2.5 flex items-center gap-4 text-xs text-text-muted">
          <span className="flex items-center gap-1"><kbd className="bg-background border border-border px-1 py-0.5 rounded">↑↓</kbd> navigate</span>
          <span className="flex items-center gap-1"><kbd className="bg-background border border-border px-1 py-0.5 rounded">↵</kbd> select</span>
          <span className="flex items-center gap-1"><kbd className="bg-background border border-border px-1 py-0.5 rounded">ESC</kbd> close</span>
        </div>
      </div>
    </div>
  );
}
`
}
