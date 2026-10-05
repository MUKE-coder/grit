"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient, API_URL } from "@/lib/api-client";
import { inputClasses } from "@/components/ui/input";
import { Search, Copy, Check, Server, Loader2 } from "@/lib/icons";

interface Access {
  group: string;
  roles?: string[];
  permissions?: string[];
  authenticated: boolean;
  api_key_only?: boolean;
  staff?: boolean;
  summary: string;
}

interface Route {
  method: string;
  path: string;
  handler?: string;
  group: string;
  access_known: boolean;
  access?: Access;
}

interface RoutesResponse {
  data: Route[];
  count: number;
  covered?: { routes: number; note?: string };
}

const METHOD_TONE: Record<string, string> = {
  GET: "border-info/30 bg-info/10 text-info",
  POST: "border-success/30 bg-success/10 text-success",
  PUT: "border-warning/30 bg-warning/10 text-warning",
  PATCH: "border-warning/30 bg-warning/10 text-warning",
  DELETE: "border-danger/30 bg-danger/10 text-danger",
};

const METHODS = ["GET", "POST", "PUT", "PATCH", "DELETE"] as const;

// The four states a route's access can be in, which is three more than a
// boolean and one more than feels necessary.
//
// "unknown" is the one that earns its place. The reference, the profiler and
// the database browser mount their own routes, and this project's route files
// never saw them, so nothing can say what they require. Calling those public
// would be the one wrong answer worth avoiding: a reader takes this column as a
// security statement, and /auth/me is not public.
type AccessKind = "unknown" | "public" | "session" | "guarded";

function kindOf(route: Route): AccessKind {
  if (!route.access_known || !route.access) return "unknown";
  const { roles, permissions, authenticated } = route.access;
  if ((roles && roles.length > 0) || (permissions && permissions.length > 0)) return "guarded";
  if (authenticated) return "session";
  return "public";
}

const ACCESS_TONE: Record<AccessKind, string> = {
  unknown: "border-dashed border-border text-text-muted",
  public: "border-border text-text-secondary",
  session: "border-info/30 bg-info/10 text-info",
  guarded: "border-accent/30 bg-accent/10 text-accent",
};

const ACCESS_LABEL: Record<AccessKind, string> = {
  unknown: "not ours",
  public: "public",
  session: "signed in",
  guarded: "guarded",
};

function AccessBadge({ route }: { route: Route }) {
  const kind = kindOf(route);
  const text = kind === "unknown" ? "unknown" : (route.access?.summary ?? "public");
  return (
    <span
      className={
        "inline-flex whitespace-nowrap rounded border px-1.5 py-0.5 font-mono text-[11px] " +
        ACCESS_TONE[kind]
      }
      title={
        kind === "unknown"
          ? "This route was registered by a mounted dashboard, not by this project's route files, so what it requires is not knowable here."
          : undefined
      }
    >
      {text}
    </span>
  );
}

// A curl for the route, with the path parameters left as they are written. A
// command with :id still in it is obviously a template; one with a plausible
// fake id in it is the kind of thing people paste without reading.
function curlFor(route: Route): string {
  const url = API_URL.replace(/\/$/, "") + route.path;
  const parts = ["curl -X " + route.method, '"' + url + '"'];
  if (route.method !== "GET" && route.method !== "DELETE") {
    parts.push('-H "Content-Type: application/json"');
    parts.push("-d '{}'");
  }
  // Omitted for a route that needs nothing: a public endpoint's curl should be
  // the command that actually works, not one carrying a header the reader then
  // has to decide about.
  if (kindOf(route) !== "public") {
    parts.push('-H "Authorization: Bearer $TOKEN"');
  }
  return parts.join(" \\\n  ");
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <button
      onClick={async () => {
        try {
          await navigator.clipboard.writeText(text);
          setCopied(true);
          setTimeout(() => setCopied(false), 1500);
        } catch {
          // A browser that refuses the clipboard, usually because the page is
          // not on https. Saying nothing is better than an error for something
          // the person can still select by hand.
        }
      }}
      title="Copy as curl"
      className="inline-flex items-center gap-1 rounded border border-border px-1.5 py-1 text-[11px] text-text-muted hover:bg-bg-hover hover:text-foreground"
    >
      {copied ? <Check className="h-3 w-3 text-success" /> : <Copy className="h-3 w-3" />}
      {copied ? "copied" : "curl"}
    </button>
  );
}

export default function RoutesPage() {
  const [query, setQuery] = useState("");
  const [method, setMethod] = useState<string>("");
  const [kind, setKind] = useState<AccessKind | "">("");

  const { data, isLoading } = useQuery<RoutesResponse>({
    queryKey: ["system-routes"],
    queryFn: async () => {
      const { data } = await apiClient.get<RoutesResponse>("/api/admin/routes");
      return data;
    },
    // The table changes when the binary does, which is not while somebody is
    // looking at it.
    staleTime: 5 * 60_000,
  });

  const all = data?.data ?? [];

  // Filtered in the browser rather than by refetching: it is a few hundred rows
  // of three short strings, and typing should not wait on a round trip.
  const routes = useMemo(() => {
    const q = query.trim().toLowerCase();
    return all.filter((r) => {
      if (method && r.method !== method) return false;
      if (kind && kindOf(r) !== kind) return false;
      if (!q) return true;
      return (
        r.path.toLowerCase().includes(q) ||
        (r.handler ?? "").toLowerCase().includes(q) ||
        (r.access?.summary ?? "").toLowerCase().includes(q) ||
        r.group.toLowerCase().includes(q)
      );
    });
  }, [all, query, method, kind]);

  // The counts are the reason this screen is worth opening rather than running
  // the command: "how many of my endpoints are reachable with no credentials"
  // is one number, and it was not available anywhere before.
  const counts = useMemo(() => {
    const out: Record<AccessKind, number> = { unknown: 0, public: 0, session: 0, guarded: 0 };
    for (const route of all) out[kindOf(route)] += 1;
    return out;
  }, [all]);

  const groups = useMemo(() => {
    const out = new Map<string, Route[]>();
    for (const route of routes) {
      const list = out.get(route.group);
      if (list) list.push(route);
      else out.set(route.group, [route]);
    }
    return [...out.entries()];
  }, [routes]);

  return (
    <div className="space-y-6">
      <PageHeader
        title="Routes"
        subtitle="Every endpoint this API serves, and what each one asks of a caller."
      />

      <div className="grid gap-3 sm:grid-cols-4">
        {(["public", "session", "guarded", "unknown"] as const).map((k) => (
          <button
            key={k}
            onClick={() => setKind(kind === k ? "" : k)}
            className={
              "rounded-xl border p-4 text-left transition-colors " +
              (kind === k ? "border-accent bg-accent/5" : "border-border bg-bg-secondary hover:bg-bg-hover")
            }
          >
            <div className="font-mono text-2xl font-semibold text-foreground">{counts[k]}</div>
            <div className="mt-1 text-xs text-text-secondary">{ACCESS_LABEL[k]}</div>
          </button>
        ))}
      </div>

      <div className="flex flex-wrap items-center gap-2">
        <div className="relative flex-1 min-w-[16rem]">
          <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-text-muted" />
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Filter by path, handler, group or permission"
            className={inputClasses({ className: "pl-9" })}
          />
        </div>
        {METHODS.map((m) => (
          <button
            key={m}
            onClick={() => setMethod(method === m ? "" : m)}
            className={
              "rounded-full border px-2.5 py-1 text-xs font-semibold " +
              (method === m ? METHOD_TONE[m] : "border-border text-text-secondary hover:bg-bg-hover")
            }
          >
            {m}
          </button>
        ))}
        <span className="text-xs text-text-muted">
          {routes.length} of {data?.count ?? 0}
        </span>
      </div>

      {isLoading ? (
        <div className="flex items-center gap-2 text-sm text-text-muted">
          <Loader2 className="h-4 w-4 animate-spin" /> Reading the route table
        </div>
      ) : groups.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-secondary p-8 text-center">
          <Server className="mx-auto h-8 w-8 text-text-muted" />
          <p className="mt-3 text-sm font-semibold text-foreground">Nothing matches</p>
        </div>
      ) : (
        <div className="space-y-5">
          {groups.map(([group, rows]) => (
            <section key={group}>
              <h2 className="mb-2 text-xs font-semibold uppercase tracking-wider text-text-muted">
                {group} <span className="font-mono normal-case opacity-70">{rows.length}</span>
              </h2>
              <div className="overflow-x-auto rounded-xl border border-border">
                <table className="w-full text-sm">
                  <tbody className="divide-y divide-border">
                    {rows.map((route) => (
                      <tr key={route.method + route.path} className="hover:bg-bg-hover">
                        <td className="w-24 px-3 py-2">
                          <span
                            className={
                              "inline-flex rounded border px-1.5 py-0.5 font-mono text-[11px] font-semibold " +
                              (METHOD_TONE[route.method] ?? "border-border text-text-secondary")
                            }
                          >
                            {route.method}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-3 py-2 font-mono text-xs text-foreground">
                          {route.path}
                        </td>
                        <td className="w-px px-3 py-2">
                          <AccessBadge route={route} />
                        </td>
                        <td className="hidden px-3 py-2 xl:table-cell">
                          <div className="ml-auto max-w-[10rem] truncate text-right font-mono text-[11px] text-text-muted">
                            {route.handler}
                          </div>
                        </td>
                        <td className="w-20 px-3 py-2 text-right">
                          <CopyButton text={curlFor(route)} />
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </section>
          ))}
        </div>
      )}

      <p className="max-w-3xl text-xs text-text-muted">
        Method, path and handler come from the router. Access comes from a table{" "}
        <code className="font-mono">grit</code> builds out of{" "}
        <code className="font-mono">routes.go</code> and the per-resource route files, because
        the router does not carry it: Gin records a route&apos;s last handler and nothing about
        the middleware in front of it. The table covers{" "}
        {data?.covered?.routes ?? 0} routes, and anything outside it reads{" "}
        <span className="font-mono">unknown</span> rather than public, which is what the
        reference, the profiler and the database browser are. Run{" "}
        <code className="font-mono">grit routes</code> for the same data in a terminal, or{" "}
        <code className="font-mono">grit doctor</code> if you have added a route by hand and
        want to know whether the table is still current.
      </p>
    </div>
  );
}
