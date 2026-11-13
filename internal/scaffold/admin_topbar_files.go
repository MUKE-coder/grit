package scaffold

// adminTopBarComponent is components/chrome/TopBar.tsx: the strip above every
// page.
//
// It takes over the chrome that used to be repeated inside each page's
// PageHeader. That put the theme toggle, the notification bell and the user
// menu on the same line as the page's own title and its own primary action,
// so every screen had two unrelated groups of controls fighting for the top
// right corner, and the page's one real action was the least prominent thing
// in it.
//
// Four things live here, and they are the four that are true of the whole
// panel rather than of the page: where you are, how the system is doing, and
// the two controls (refresh, theme) that apply to whatever is on screen.
//
// Breadcrumbs are derived from the path rather than passed in. A crumb a page
// has to declare is a crumb that goes stale the day the route moves, and the
// path already knows the answer.
func adminTopBarComponent() string {
	return `"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";
import { resources } from "@/resources";
import { useSystemHealth } from "@/hooks/use-system-health";
import { DarkModeToggle } from "./DarkModeToggle";
import { NotificationBell } from "./NotificationBell";
import { RefreshCw, Menu } from "@/lib/icons";

// Queries that belong to the panel rather than to any page: the signed-in
// user, their permissions and the notification bell. Refresh leaves them
// alone; they have their own lifetimes, and refetching permissions re-gates
// the sidebar for a frame.
const CHROME_QUERY_ROOTS = ["me", "my-permissions", "notifications"];

// Words for the path segments that are not a resource slug.
//
// A segment with no entry here is title-cased from the slug, so a page added
// tomorrow gets a readable crumb without being registered anywhere. The map
// is for the ones where that guess reads badly: "Api keys", "Gdpr".
const SEGMENT_LABELS: Record<string, string> = {
  dashboard: "Dashboard",
  resources: "Content",
  system: "System hub",
  account: "Account",
  profile: "Profile",
  settings: "Settings",
  "api-keys": "API keys",
  gdpr: "GDPR",
  sso: "SSO",
  "access-reviews": "Access reviews",
  "deleted-accounts": "Deleted accounts",
  "form-shares": "Form sharing",
  cron: "Cron schedules",
  mail: "Mail preview",
  jobs: "Background jobs",
  files: "File storage",
  audit: "Audit log",
  health: "System health",
  performance: "Operations",
  observability: "Observability",
  sagas: "Sagas",
  trash: "Trash",
  backups: "Backups",
  security: "Security",
  activity: "Activity",
  roles: "Roles",
  routes: "Routes",
  support: "Support",
  notifications: "Notifications",
  "dashboard-settings": "Dashboard settings",
};

function titleCase(slug: string) {
  return slug
    .split("-")
    .map((word) => (word ? word[0].toUpperCase() + word.slice(1) : word))
    .join(" ");
}

interface Crumb {
  label: string;
  href?: string;
}

/**
 * Turns a path into crumbs.
 *
 * A record page under a resource keeps the resource's own plural as its
 * parent, so /resources/blogs/12 reads "Content / Blogs / 12" rather than
 * repeating the slug. An id segment is left as it is: it is the one part of
 * the path that is data, and guessing a title for it would mean a request.
 */
function crumbsFor(pathname: string): Crumb[] {
  const parts = pathname.split("/").filter(Boolean);
  if (parts.length === 0) return [{ label: "Dashboard" }];

  const out: Crumb[] = [];
  let href = "";

  for (let i = 0; i < parts.length; i += 1) {
    const part = parts[i];
    href += "/" + part;

    if (part === "resources") {
      // Not a page of its own. It names the group the next crumb is in.
      out.push({ label: SEGMENT_LABELS.resources });
      continue;
    }

    if (i > 0 && parts[i - 1] === "resources") {
      const resource = resources.find((r) => r.slug === part);
      out.push({
        label: resource ? (resource.label?.plural ?? resource.name) : titleCase(part),
        href,
      });
      continue;
    }

    out.push({ label: SEGMENT_LABELS[part] ?? titleCase(part), href });
  }

  // The last crumb is where you already are, so it is not a link.
  if (out.length > 0) out[out.length - 1] = { label: out[out.length - 1].label };
  return out;
}

export function TopBar({ onOpenMenu }: { onOpenMenu: () => void }) {
  const pathname = usePathname();
  const queryClient = useQueryClient();
  const { summary, isFetching } = useSystemHealth();
  const crumbs = crumbsFor(pathname);

  const dotClass =
    summary.overall === "degraded"
      ? "bg-danger"
      : summary.overall === "unknown"
        ? "bg-text-muted"
        : "bg-success";

  // Refresh reloads what the page shows, not the chrome around it.
  const onRefresh = () => {
    queryClient.invalidateQueries({
      predicate: (query) => !CHROME_QUERY_ROOTS.includes(String(query.queryKey[0])),
    });
  };

  return (
    <header className="sticky top-0 z-30 flex h-[60px] shrink-0 items-center gap-3 border-b border-border bg-background px-4 md:px-6">
      <button
        type="button"
        onClick={onOpenMenu}
        aria-label="Open menu"
        className="-ml-1 inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-md text-text-secondary transition-grit hover:bg-bg-hover hover:text-foreground md:hidden"
      >
        <Menu className="h-5 w-5" />
      </button>

      <nav aria-label="Breadcrumb" className="min-w-0 flex-1">
        <ol className="flex items-center gap-2 text-sm">
          {crumbs.map((crumb, index) => (
            <li key={index} className="flex min-w-0 items-center gap-2">
              {index > 0 && (
                <span aria-hidden="true" className="shrink-0 text-text-muted">
                  /
                </span>
              )}
              {crumb.href ? (
                <Link
                  href={crumb.href}
                  className="truncate text-text-secondary transition-grit hover:text-foreground"
                >
                  {crumb.label}
                </Link>
              ) : (
                <span
                  aria-current={index === crumbs.length - 1 ? "page" : undefined}
                  className={
                    index === crumbs.length - 1
                      ? "truncate font-medium text-foreground"
                      : "truncate text-text-muted"
                  }
                >
                  {crumb.label}
                </span>
              )}
            </li>
          ))}
        </ol>
      </nav>

      {/* How the system is doing, in the one place that is on every screen.
          A link, because the next question after "something is degraded" is
          always "which one". */}
      <Link
        href="/system/health"
        className="hidden shrink-0 items-center gap-2 rounded-full border border-border px-3 py-1.5 font-mono text-[12px] text-text-secondary transition-grit hover:border-border-strong hover:text-foreground sm:inline-flex"
      >
        <span className={"h-2 w-2 shrink-0 rounded-full " + dotClass} aria-hidden="true" />
        {summary.label}
      </Link>

      <div className="flex shrink-0 items-center gap-1">
        <button
          type="button"
          onClick={onRefresh}
          aria-label="Refresh this page's data"
          className="inline-flex h-9 w-9 items-center justify-center rounded-md text-text-secondary transition-grit hover:bg-bg-hover hover:text-foreground"
        >
          <RefreshCw className={"h-4 w-4 " + (isFetching ? "animate-spin" : "")} />
        </button>
        <DarkModeToggle />
        <NotificationBell />
      </div>
    </header>
  );
}
`
}
