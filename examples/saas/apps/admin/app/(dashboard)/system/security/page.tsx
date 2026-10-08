"use client";

import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { apiClient } from "@/lib/api-client";
import {
  Shield, AlertTriangle, AlertCircle, ExternalLink, Activity as ActivityIcon, Clock,
} from "@/lib/icons";

// The Sentinel UI is mounted on the Go API, not on this admin host. Use
// the API base so "Open Sentinel" works whether the admin is on :3001
// in dev or a different origin in prod.
import { API_URL } from "@/lib/api-core";

// DegradedBanner says which panels have no data behind them.
//
// The dashboard used to answer a dead upstream with a page of zeros and a 200,
// which reads as good news. The API now names the calls that failed, and this
// puts that on the screen above the numbers it applies to.
function DegradedBanner({ source, degraded }: { source: string; degraded?: string[] }) {
  if (!degraded || degraded.length === 0) return null;
  return (
    <div className="mt-4 flex items-start gap-2 rounded-lg border border-warning/30 bg-warning/5 px-4 py-3 text-sm text-warning">
      <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
      <span>
        {source} did not answer for {degraded.join(", ")}. Those panels are empty rather than zero.
      </span>
    </div>
  );
}

interface SecuritySummary {
  banned_ips_now: number;
  auto_bans_24h: number;
  rate_limited_last_hour: number;
  active_bans?: Array<{ ip: string; reason: string; expires_at: string; level: number }>;
  rate_limit_hits_5min?: Array<{ ip: string; hits: number; last_hit: string }>;
  recent_threats?: Array<{ id: string; type: string; ip: string; description: string; created_at: string }>;
  // The calls to Sentinel that did not answer. Empty is the healthy case.
  degraded?: string[];
}

export default function SecurityPage() {
  const { data, isLoading } = useQuery<SecuritySummary>({
    queryKey: ["security", "summary"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<SecuritySummary>("/api/admin/security/summary");
        return data;
      } catch {
        // Nothing came back. A page of zeros would read as "no IP is banned and
        // nothing has attacked you", so say what actually happened instead.
        return { banned_ips_now: 0, auto_bans_24h: 0, rate_limited_last_hour: 0, degraded: ["the whole summary"] };
      }
    },
    refetchInterval: 60_000,
  });

  return (
    <div>
      <PageHeader
        title="Security"
        subtitle="IP bans, rate-limit pressure, recent threats — powered by Sentinel."
        actions={
          <a
            href={`${API_URL}/sentinel/ui`}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex shrink-0 items-center gap-2 whitespace-nowrap rounded-lg border border-border bg-bg-elevated px-3 py-2 text-sm font-medium text-foreground hover:bg-bg-hover"
          >
            <ExternalLink className="h-4 w-4" />
            Open Sentinel
          </a>
        }
      />

      <DegradedBanner source="Sentinel" degraded={data?.degraded} />

      {/* KPI row */}
      {isLoading ? (
        <SkeletonCards count={3} />
      ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <KPI
            label="Currently Banned IPs"
            value={data?.banned_ips_now ?? 0}
            icon={<Shield className="h-4 w-4" />}
            tone={(data?.banned_ips_now ?? 0) > 0 ? "danger" : "default"}
          />
          <KPI
            label="Auto-bans (last 24h)"
            value={data?.auto_bans_24h ?? 0}
            icon={<AlertCircle className="h-4 w-4" />}
            tone={(data?.auto_bans_24h ?? 0) > 0 ? "warning" : "default"}
          />
          <KPI
            label="Rate-limited IPs (last hour)"
            value={data?.rate_limited_last_hour ?? 0}
            icon={<AlertTriangle className="h-4 w-4" />}
            tone={(data?.rate_limited_last_hour ?? 0) > 0 ? "warning" : "default"}
          />
        </div>
      )}

      {/* Auto-ban escalation policy — explainer card. The schedule is
          surfaced here so operators don't have to grep Sentinel config to
          understand what's about to happen to a re-offender. */}
      <section className="mt-6 rounded-xl border border-accent/20 bg-accent/5 p-5">
        <div className="flex items-start gap-3">
          <Clock className="mt-0.5 h-4 w-4 shrink-0 text-accent" />
          <div className="min-w-0">
            <p className="text-sm font-semibold text-foreground">Escalating auto-ban policy</p>
            <p className="mt-1 text-sm text-text-secondary">
              When an IP trips the brute-force rate limit, Sentinel auto-bans it.
              Re-offenders escalate quickly so a bot can't simply wait out the cooldown.
            </p>
            <ul className="mt-3 grid grid-cols-1 gap-2 text-xs text-text-secondary sm:grid-cols-4">
              <li className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
                <span className="block text-[10px] font-semibold uppercase tracking-wide text-text-muted">1st offence</span>
                <span className="text-foreground font-mono">5 hours</span>
              </li>
              <li className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
                <span className="block text-[10px] font-semibold uppercase tracking-wide text-text-muted">2nd offence</span>
                <span className="text-foreground font-mono">8 hours</span>
              </li>
              <li className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
                <span className="block text-[10px] font-semibold uppercase tracking-wide text-text-muted">3rd offence</span>
                <span className="text-foreground font-mono">24 hours</span>
              </li>
              <li className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
                <span className="block text-[10px] font-semibold uppercase tracking-wide text-text-muted">4th+ offence</span>
                <span className="text-foreground font-mono">7 days</span>
              </li>
            </ul>
          </div>
        </div>
      </section>

      {/* Active bans */}
      <Section title="Active IP bans" icon={<Shield className="h-4 w-4" />}>
        {(data?.active_bans?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No IPs are currently banned.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.active_bans!.map((b) => (
              <li key={b.ip} className="flex items-center justify-between gap-3 px-5 py-3">
                <div className="min-w-0">
                  <p className="font-mono text-sm text-foreground">{b.ip}</p>
                  <p className="text-xs text-text-muted">{b.reason}</p>
                </div>
                <div className="text-right text-xs">
                  <span className="rounded bg-danger/10 px-1.5 py-0.5 font-semibold uppercase text-danger">
                    Level {b.level}
                  </span>
                  <p className="mt-1 text-text-muted">expires {new Date(b.expires_at).toLocaleString()}</p>
                </div>
              </li>
            ))}
          </ul>
        )}
      </Section>

      {/* Rate-limit pressure */}
      <Section title="IPs hitting rate limits (last 5 min)" icon={<AlertTriangle className="h-4 w-4" />}>
        {(data?.rate_limit_hits_5min?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No rate-limit caps in the last 5 minutes.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.rate_limit_hits_5min!.map((r) => (
              <li key={r.ip} className="flex items-center justify-between gap-3 px-5 py-3">
                <p className="font-mono text-sm text-foreground">{r.ip}</p>
                <span className="text-xs text-text-muted">
                  {r.hits} hits · last {new Date(r.last_hit).toLocaleTimeString()}
                </span>
              </li>
            ))}
          </ul>
        )}
      </Section>

      {/* Recent threats */}
      <Section title="Recent threats" icon={<ActivityIcon className="h-4 w-4" />}>
        {(data?.recent_threats?.length ?? 0) === 0 ? (
          <p className="px-5 py-8 text-center text-sm text-text-muted">No threats detected recently.</p>
        ) : (
          <ul className="divide-y divide-border">
            {data!.recent_threats!.map((t) => (
              <li key={t.id} className="px-5 py-3">
                <div className="flex items-center justify-between">
                  <p className="text-sm font-semibold text-foreground">{t.type}</p>
                  <p className="text-xs text-text-muted">{new Date(t.created_at).toLocaleString()}</p>
                </div>
                <p className="mt-1 text-xs text-text-secondary">{t.description}</p>
                <p className="mt-1 font-mono text-xs text-text-muted">{t.ip}</p>
              </li>
            ))}
          </ul>
        )}
      </Section>
    </div>
  );
}

function KPI({ label, value, icon, tone }: { label: string; value: number; icon: React.ReactNode; tone: "default" | "warning" | "danger" }) {
  const toneClass = {
    default: "border-border bg-bg-elevated",
    warning: "border-warning/30 bg-warning/5",
    danger:  "border-danger/30 bg-danger/5",
  }[tone];
  const iconClass = { default: "text-text-secondary", warning: "text-warning", danger: "text-danger" }[tone];
  return (
    <div className={"rounded-xl border p-4 " + toneClass}>
      <div className="flex items-center justify-between">
        <p className="text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</p>
        <span className={iconClass}>{icon}</span>
      </div>
      <p className="mt-2 text-3xl font-bold text-foreground">{value}</p>
    </div>
  );
}

function Section({ title, icon, children }: { title: string; icon: React.ReactNode; children: React.ReactNode }) {
  return (
    <section className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-elevated">
      <header className="flex items-center gap-2 border-b border-border px-5 py-3">
        <span className="text-text-secondary">{icon}</span>
        <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">{title}</p>
      </header>
      {children}
    </section>
  );
}
