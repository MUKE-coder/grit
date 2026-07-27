package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// stripUseClient removes "use client" directives from component code.
// TanStack Router apps don't need this since everything is client-side.
func stripUseClient(code string) string {
	code = strings.Replace(code, "\"use client\"\n\n", "", 1)
	code = strings.Replace(code, "'use client'\n\n", "", 1)
	code = strings.Replace(code, "\"use client\";\n\n", "", 1)
	code = strings.Replace(code, "'use client';\n\n", "", 1)
	code = strings.Replace(code, "\"use client\"\n", "", 1)
	code = strings.Replace(code, "'use client'\n", "", 1)
	return code
}

// nextToTanStack adapts a component originally written for the Next.js admin so
// it runs unchanged under Vite + TanStack Router. The shared admin component
// templates (sidebar, navbar, tables, forms, ...) are reused by BOTH admins, so
// we can't edit their source — instead we rewrite their Next.js imports to a
// small compat shim (@/lib/next-compat, see adminNextCompatShim) that maps
// next/link, next/image, next/navigation and next/dynamic onto TanStack Router
// equivalents. Without this, the Vite admin ships components that import
// `usePathname`/`next/link`/`next/image` and fail to build (grit issue #69).
func nextToTanStack(code string) string {
	code = stripUseClient(code)

	// Default imports → named imports from the compat shim.
	replacements := []struct{ from, to string }{
		{`import Link from "next/link";`, `import { Link } from "@/lib/next-compat";`},
		{`import Link from 'next/link';`, `import { Link } from "@/lib/next-compat";`},
		{`import Image from "next/image";`, `import { Image } from "@/lib/next-compat";`},
		{`import Image from 'next/image';`, `import { Image } from "@/lib/next-compat";`},
		{`import dynamic from "next/dynamic";`, `import { dynamic } from "@/lib/next-compat";`},
		{`import dynamic from 'next/dynamic';`, `import { dynamic } from "@/lib/next-compat";`},
	}
	for _, r := range replacements {
		code = strings.ReplaceAll(code, r.from, r.to)
	}

	// next/navigation named imports (useRouter, usePathname, useSearchParams,
	// ReadonlyURLSearchParams) → the shim, which re-exports compatible versions.
	code = strings.ReplaceAll(code, `from "next/navigation"`, `from "@/lib/next-compat"`)
	code = strings.ReplaceAll(code, `from 'next/navigation'`, `from "@/lib/next-compat"`)

	// process.env → import.meta.env.
	//
	// Vite does not polyfill `process`, so a transformed component reading
	// process.env.NEXT_PUBLIC_X throws "process is not defined" in the browser.
	// The build wouldn't catch it either — vite build uses esbuild and does no
	// type checking, so this surfaced only as a typecheck failure and would
	// otherwise have reached users as a blank screen.
	code = nextPublicEnvPattern.ReplaceAllString(code, "import.meta.env.VITE_$1")
	code = strings.ReplaceAll(code, "process.env.NODE_ENV", "import.meta.env.MODE")

	return code
}

// nextPublicEnvPattern matches NEXT_PUBLIC_* env reads so they can be rewritten
// to Vite's equivalent. Vite only exposes VITE_-prefixed vars to client code,
// so the prefix has to change as well as the accessor.
var nextPublicEnvPattern = regexp.MustCompile(`process\.env\.NEXT_PUBLIC_([A-Z0-9_]+)`)

// adminTanStackAPIClient adapts the shared Next.js admin api-client for Vite:
//   - process.env.NEXT_PUBLIC_API_URL → import.meta.env.VITE_API_URL (Vite has no
//     process.env at build or runtime).
//   - adds an `api` alias that the TanStack route templates import (the Next.js
//     admin refers to it as `apiClient`).
func adminTanStackAPIClient() string {
	code := adminAPIClient()
	code = strings.ReplaceAll(code,
		`process.env.NEXT_PUBLIC_API_URL`,
		`(import.meta as any).env?.VITE_API_URL`)
	code += "\n// TanStack route templates import { api }; alias it to the configured instance.\nexport const api = apiClient;\n"
	return code
}

// adminTanStackPageRoute returns a thin TanStack route that renders a page
// component reused verbatim from the Next.js admin.
//
// The Vite admin previously hand-wrote each route as a placeholder ("System
// page content will be loaded here", dashboards with "--" stats), so it drifted
// far from the Next.js admin. Routing is the ONLY thing that should differ
// between the two admins — the pages themselves are shared, transformed by
// nextToTanStack. Keep new routes thin like this.
func adminTanStackPageRoute(routeID, importPath string) string {
	return fmt.Sprintf(`import { createFileRoute } from '@tanstack/react-router'
import Page from '%s'

export const Route = createFileRoute('%s')({
  component: Page,
})
`, importPath, routeID)
}

// adminTanStackGlobalCSS adapts the shared admin globals.css for Tailwind v4 +
// the @tailwindcss/vite plugin (the Vite admin uses v4; the Next.js admin stays
// on v3). It swaps the v3 @tailwind directives for the v4 @import, registers the
// design-token CSS variables as @theme colors/fonts, and keeps every
// [data-theme] block below unchanged so runtime theme switching still works.
//
// Regular @theme (NOT @theme inline) is deliberate: the utilities compile to
// var(--color-*), whose values are var(--bg-*) indirections that re-resolve per
// element — so overriding --bg-* under [data-theme="…"] repaints the whole UI.
func adminTanStackGlobalCSS() string {
	v4Header := `@import "tailwindcss";
/* Successor to tailwindcss-animate for v4: provides animate-in / fade-in /
 * zoom-in used by the dialog, sheet and dropdown components. */
@import "tw-animate-css";

/* Map the design-token CSS variables (defined per [data-theme] below) onto
 * Tailwind utility colors + fonts. Regular @theme so the var() indirection
 * re-resolves per element and [data-theme] switching repaints at runtime. */
@theme {
  --color-background: var(--bg-primary);
  --color-bg-secondary: var(--bg-secondary);
  --color-bg-tertiary: var(--bg-tertiary);
  --color-bg-elevated: var(--bg-elevated);
  --color-bg-hover: var(--bg-hover);
  --color-border: var(--border);
  --color-foreground: var(--text-primary);
  --color-text-secondary: var(--text-secondary);
  --color-text-muted: var(--text-muted);
  --color-accent: var(--accent);
  --color-accent-hover: var(--accent-hover);
  --color-success: var(--success);
  --color-danger: var(--danger);
  --color-warning: var(--warning);
  --color-info: var(--info);

  --font-sans: var(--font-display), system-ui, sans-serif;
  --font-mono: var(--font-mono), ui-monospace, monospace;
  --font-serif: var(--font-serif), Georgia, serif;
}

/* v3 parity: v4's default border color is currentColor, but the admin relies
 * on the themed border. Point bare borders at --border. */
@layer base {
  *,
  ::after,
  ::before {
    border-color: var(--border);
  }
}`
	css := adminGlobalCSS()
	css = strings.Replace(css,
		"@tailwind base;\n@tailwind components;\n@tailwind utilities;",
		v4Header, 1)
	return css
}

// adminTanStackUseAuth reuses the shared auth hooks and adds the useAuth() hook
// the TanStack route templates expect ({ user }), backed by the shared useMe().
func adminTanStackUseAuth() string {
	code := nextToTanStack(adminUseAuth())
	code += "\n// next-style hook the TanStack route templates consume: { user }.\n" +
		"export function useAuth() {\n" +
		"  const { data: user } = useMe();\n" +
		"  return { user: user ?? null };\n" +
		"}\n"
	return code
}

// adminNextCompatShim is the small module the Vite admin's rewritten imports
// point at. It re-implements the handful of Next.js APIs the shared admin
// components use (next/link, next/image, next/navigation, next/dynamic) on top
// of TanStack Router + plain DOM, so those components compile and run unchanged.
func adminNextCompatShim() string {
	return `// next/* compatibility shim for the Vite + TanStack Router admin.
//
// The admin's shared components were authored for the Next.js admin and import
// from next/link, next/image, next/navigation and next/dynamic. The Vite build
// rewrites those imports to this module (see nextToTanStack in the Grit
// scaffold) so the same components run under TanStack Router with no source
// changes. The Next.js admin keeps using the real next/* modules.
import * as React from "react";
import {
  Link as RouterLink,
  useNavigate,
  useRouterState,
  useParams as useRouterParams,
} from "@tanstack/react-router";

// RouterLink is heavily generic over the typed route tree; the shared components
// pass dynamic string hrefs, so we present a loosely-typed Link.
const AnyLink = RouterLink as unknown as React.ComponentType<any>;

// next/link: <Link href="..."> -> TanStack <Link to="...">.
export function Link({
  href,
  children,
  ...props
}: { href: string; children?: React.ReactNode } & Record<string, unknown>) {
  return (
    <AnyLink to={href} {...props}>
      {children}
    </AnyLink>
  );
}

// next/image: <Image src alt width height .../> -> a plain <img>. Next-only
// props (width/height/priority/fill/quality/placeholder/loader/sizes) are
// dropped; src may be a string or a static-import object ({ src }).
export function Image({
  src,
  alt,
  width,
  height,
  priority,
  fill,
  quality,
  placeholder,
  loader,
  sizes,
  unoptimized,
  ...rest
}: any) {
  const resolved = typeof src === "string" ? src : src?.src;
  return <img src={resolved} alt={alt ?? ""} {...rest} />;
}

// next/navigation: useRouter().push/replace -> TanStack navigate.
export function useRouter() {
  const navigate = useNavigate();
  return {
    push: (to: string) => navigate({ to: to as string }),
    replace: (to: string) => navigate({ to: to as string, replace: true }),
    back: () => window.history.back(),
    forward: () => window.history.forward(),
    refresh: () => {},
    prefetch: () => {},
  };
}

// next/navigation: usePathname().
export function usePathname(): string {
  return useRouterState({ select: (s: any) => s.location.pathname as string });
}

// next/navigation: useParams(). Non-strict so it resolves against whatever
// route is active, matching Next.js's untyped-by-default behaviour. Detail
// pages read params.id — TanStack names the segment $id, which yields the
// same { id } shape.
export function useParams<T = Record<string, string>>(): T {
  // Cast through any: useParams is generic over the generated route tree, so a
  // route-agnostic helper cannot satisfy its option type.
  return (useRouterParams as any)({ strict: false }) as T;
}

export type ReadonlyURLSearchParams = URLSearchParams;

// next/navigation: useSearchParams() -> a URLSearchParams over the current query
// string. Read-only usage (.get/.has/.getAll) works exactly as on the web.
export function useSearchParams(): ReadonlyURLSearchParams {
  const searchStr = useRouterState({ select: (s: any) => (s.location.searchStr as string) ?? "" });
  return new URLSearchParams(searchStr);
}

// next/dynamic(loader, { ssr:false }) -> React.lazy (Vite has no SSR, so the
// ssr flag is a no-op). Supports both { default: C } and bare-component loaders.
export function dynamic<T extends React.ComponentType<any>>(
  loader: () => Promise<{ default: T } | T>,
  _options?: { ssr?: boolean; loading?: React.ComponentType },
): React.LazyExoticComponent<T> {
  return React.lazy(async () => {
    const mod = (await loader()) as any;
    return mod && mod.default ? { default: mod.default as T } : { default: mod as T };
  });
}
`
}

func writeAdminTanStackFiles(root string, opts Options) error {
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{
		// Config files
		filepath.Join(adminRoot, "package.json"):   adminTanStackPackageJSON(opts),
		filepath.Join(adminRoot, "vite.config.ts"): adminTanStackViteConfig(),
		filepath.Join(adminRoot, "index.html"):     adminTanStackIndexHTML(opts),
		// Tailwind v4: styling is driven by the @tailwindcss/vite plugin +
		// @theme/@import in globals.css — no tailwind.config or postcss.config.
		filepath.Join(adminRoot, "tsconfig.json"):      adminTanStackTSConfig(),
		filepath.Join(adminRoot, "src", "main.tsx"):    adminTanStackMain(),
		filepath.Join(adminRoot, "src", "globals.css"): adminTanStackGlobalCSS(),

		// Pages — reused verbatim from the Next.js admin (transformed by
		// nextToTanStack). These are the real dashboards/forms/tables; the
		// routes below are thin wrappers that render them. Adding a page here
		// plus a route below is how you surface a new admin screen.
		filepath.Join(adminRoot, "src", "pages", "login.tsx"):                    nextToTanStack(adminThemedLoginPage()),
		filepath.Join(adminRoot, "src", "pages", "sign-up.tsx"):                  nextToTanStack(adminThemedSignUpPage()),
		filepath.Join(adminRoot, "src", "pages", "forgot-password.tsx"):          nextToTanStack(adminThemedForgotPasswordPage()),
		filepath.Join(adminRoot, "src", "pages", "reset-password.tsx"):           nextToTanStack(adminThemedResetPasswordPage()),
		filepath.Join(adminRoot, "src", "pages", "callback.tsx"):                 nextToTanStack(adminAuthCallbackPage()),
		filepath.Join(adminRoot, "src", "pages", "not-found.tsx"):                nextToTanStack(adminNotFoundPage()),
		filepath.Join(adminRoot, "src", "pages", "dashboard.tsx"):                nextToTanStack(adminDashboardPageForStyle(opts.Style)),
		filepath.Join(adminRoot, "src", "pages", "profile.tsx"):                  nextToTanStack(adminCaptivatingProfile()),
		filepath.Join(adminRoot, "src", "pages", "settings", "dashboard.tsx"):    nextToTanStack(adminDashboardSettingsPageTS()),
		filepath.Join(adminRoot, "src", "pages", "resources", "users.tsx"):       nextToTanStack(adminUsersPage()),
		filepath.Join(adminRoot, "src", "pages", "resources", "blogs.tsx"):       nextToTanStack(adminBlogsListPage()),
		filepath.Join(adminRoot, "src", "pages", "resources", "blog-detail.tsx"): nextToTanStack(adminBlogDetailPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "index.tsx"):          nextToTanStack(adminSystemHubPageV2()),
		filepath.Join(adminRoot, "src", "pages", "system", "activity.tsx"):       nextToTanStack(adminWalkieActivityPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "health.tsx"):         nextToTanStack(adminSystemHealthPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "notifications.tsx"):  nextToTanStack(adminNotificationsPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "performance.tsx"):    nextToTanStack(adminPerformancePageV2()),
		filepath.Join(adminRoot, "src", "pages", "system", "security.tsx"):       nextToTanStack(adminSecurityPageV2()),
		filepath.Join(adminRoot, "src", "pages", "system", "access-reviews.tsx"): nextToTanStack(adminAccessReviewPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "gdpr.tsx"): nextToTanStack(adminGDPRPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "roles.tsx"):          nextToTanStack(adminRolesPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "support.tsx"):        nextToTanStack(adminSupportListPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "ticket.tsx"):         nextToTanStack(adminTicketThreadPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "jobs.tsx"):           nextToTanStack(adminJobsPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "files.tsx"):          nextToTanStack(adminFilesPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "cron.tsx"):           nextToTanStack(adminCronPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "mail.tsx"):           nextToTanStack(adminMailPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "backups.tsx"):        nextToTanStack(adminBackupsPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "observability.tsx"):  nextToTanStack(adminObservabilityPage()),
		filepath.Join(adminRoot, "src", "pages", "system", "form-shares.tsx"):    nextToTanStack(adminFormSharesPage()),

		// Routes — thin wrappers around the pages above. Route ids mirror the
		// Next.js admin's URLs exactly so the shared sidebar's links resolve.
		filepath.Join(adminRoot, "src", "routes", "__root.tsx"):                                    adminTanStackRootRoute(),
		filepath.Join(adminRoot, "src", "routes", "index.tsx"):                                     adminTanStackRedirectRoute(),
		filepath.Join(adminRoot, "src", "routes", "_auth.tsx"):                                     adminTanStackAuthLayout(),
		filepath.Join(adminRoot, "src", "routes", "_auth", "login.tsx"):                            adminTanStackPageRoute("/_auth/login", "@/pages/login"),
		filepath.Join(adminRoot, "src", "routes", "_auth", "sign-up.tsx"):                          adminTanStackPageRoute("/_auth/sign-up", "@/pages/sign-up"),
		filepath.Join(adminRoot, "src", "routes", "_auth", "forgot-password.tsx"):                  adminTanStackPageRoute("/_auth/forgot-password", "@/pages/forgot-password"),
		filepath.Join(adminRoot, "src", "routes", "_auth", "reset-password.tsx"):                   adminTanStackPageRoute("/_auth/reset-password", "@/pages/reset-password"),
		filepath.Join(adminRoot, "src", "routes", "_auth", "callback.tsx"):                         adminTanStackPageRoute("/_auth/callback", "@/pages/callback"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard.tsx"):                                adminTanStackDashboardLayout(),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "dashboard.tsx"):                   adminTanStackPageRoute("/_dashboard/dashboard", "@/pages/dashboard"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "profile.tsx"):                     adminTanStackPageRoute("/_dashboard/profile", "@/pages/profile"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "settings", "dashboard.tsx"):       adminTanStackPageRoute("/_dashboard/settings/dashboard", "@/pages/settings/dashboard"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "resources", "users", "index.tsx"): adminTanStackPageRoute("/_dashboard/resources/users/", "@/pages/resources/users"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "resources", "users", "$id.tsx"):   adminResourceDetailRouteTanStack("users", "users"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "resources", "blogs", "index.tsx"): adminTanStackPageRoute("/_dashboard/resources/blogs/", "@/pages/resources/blogs"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "resources", "blogs", "$id.tsx"):   adminTanStackPageRoute("/_dashboard/resources/blogs/$id", "@/pages/resources/blog-detail"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "index.tsx"):             adminTanStackPageRoute("/_dashboard/system/", "@/pages/system"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "activity.tsx"):          adminTanStackPageRoute("/_dashboard/system/activity", "@/pages/system/activity"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "health.tsx"):            adminTanStackPageRoute("/_dashboard/system/health", "@/pages/system/health"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "notifications.tsx"):     adminTanStackPageRoute("/_dashboard/system/notifications", "@/pages/system/notifications"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "performance.tsx"):       adminTanStackPageRoute("/_dashboard/system/performance", "@/pages/system/performance"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "security.tsx"):          adminTanStackPageRoute("/_dashboard/system/security", "@/pages/system/security"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "access-reviews.tsx"):    adminTanStackPageRoute("/_dashboard/system/access-reviews", "@/pages/system/access-reviews"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "gdpr.tsx"):             adminTanStackPageRoute("/_dashboard/system/gdpr", "@/pages/system/gdpr"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "roles.tsx"):             adminTanStackPageRoute("/_dashboard/system/roles", "@/pages/system/roles"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "support", "index.tsx"):  adminTanStackPageRoute("/_dashboard/system/support/", "@/pages/system/support"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "support", "$id.tsx"):    adminTanStackPageRoute("/_dashboard/system/support/$id", "@/pages/system/ticket"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "jobs.tsx"):              adminTanStackPageRoute("/_dashboard/system/jobs", "@/pages/system/jobs"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "files.tsx"):             adminTanStackPageRoute("/_dashboard/system/files", "@/pages/system/files"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "cron.tsx"):              adminTanStackPageRoute("/_dashboard/system/cron", "@/pages/system/cron"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "mail.tsx"):              adminTanStackPageRoute("/_dashboard/system/mail", "@/pages/system/mail"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "backups.tsx"):           adminTanStackPageRoute("/_dashboard/system/backups", "@/pages/system/backups"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "observability.tsx"):     adminTanStackPageRoute("/_dashboard/system/observability", "@/pages/system/observability"),
		filepath.Join(adminRoot, "src", "routes", "_dashboard", "system", "form-shares.tsx"):       adminTanStackPageRoute("/_dashboard/system/form-shares", "@/pages/system/form-shares"),

		// Lib (same as Next.js versions)
		filepath.Join(adminRoot, "src", "lib", "api-client.ts"):   adminTanStackAPIClient(),
		filepath.Join(adminRoot, "src", "lib", "query-client.ts"): adminQueryClient(),
		filepath.Join(adminRoot, "src", "lib", "utils.ts"):        adminUtils(),
		filepath.Join(adminRoot, "src", "lib", "resource.ts"):     adminResourceTypes(),
		filepath.Join(adminRoot, "src", "lib", "icons.ts"):        adminIconMap(),
		filepath.Join(adminRoot, "src", "lib", "formatters.ts"):   adminFormatters(),
		// Compat shim: maps next/link, next/image, next/navigation, next/dynamic
		// onto TanStack Router so the shared admin components run under Vite (#69).
		filepath.Join(adminRoot, "src", "lib", "next-compat.tsx"): adminNextCompatShim(),
		// Lib modules imported by reused table/form components.
		filepath.Join(adminRoot, "src", "lib", "excel-utils.ts"):  adminExcelUtils(),
		filepath.Join(adminRoot, "src", "lib", "file-accepts.ts"): adminFileAcceptsLib(),

		// Hooks (same as Next.js versions)
		filepath.Join(adminRoot, "src", "hooks", "use-auth.ts"):             adminTanStackUseAuth(),
		filepath.Join(adminRoot, "src", "hooks", "use-resource.ts"):         nextToTanStack(adminUseResource()),
		filepath.Join(adminRoot, "src", "hooks", "use-system.ts"):           nextToTanStack(adminUseSystem()),
		filepath.Join(adminRoot, "src", "hooks", "use-profile.ts"):          nextToTanStack(adminUseProfile()),
		filepath.Join(adminRoot, "src", "hooks", "use-roles.ts"):            nextToTanStack(adminUseRoles()),
		filepath.Join(adminRoot, "src", "hooks", "use-permissions.ts"):      nextToTanStack(adminUsePermissions()),
		filepath.Join(adminRoot, "src", "hooks", "use-modules.ts"):          nextToTanStack(adminUseModules()),
		filepath.Join(adminRoot, "src", "hooks", "use-backups.ts"):          nextToTanStack(adminUseBackups()),
		filepath.Join(adminRoot, "src", "hooks", "use-dashboard-layout.ts"): nextToTanStack(adminUseDashboardLayoutTS()),
		filepath.Join(adminRoot, "src", "hooks", "use-toasted-mutation.ts"): nextToTanStack(adminToastHook()),

		// Lib modules the reused pages import.
		filepath.Join(adminRoot, "src", "lib", "dashboard-catalog.ts"): nextToTanStack(adminDashboardCatalogTS()),
		filepath.Join(adminRoot, "src", "lib", "export.ts"):            nextToTanStack(adminExportLib()),

		// Dashboard widgets. The Next.js admin gets these from
		// writeAdminResourceDashboardWidgets + writeAdminCustomChartFiles, which
		// write to a non-src root — so they're mirrored here instead of calling
		// those writers.
		filepath.Join(adminRoot, "src", "components", "dashboard", "ResourceStatCard.tsx"):    nextToTanStack(adminResourceStatCardTSX()),
		filepath.Join(adminRoot, "src", "components", "dashboard", "ResourceLatestTable.tsx"): nextToTanStack(adminResourceLatestTableTSX()),
		filepath.Join(adminRoot, "src", "components", "dashboard", "ResourceWidgetsRow.tsx"):  nextToTanStack(adminResourceWidgetsRowTSX()),
		filepath.Join(adminRoot, "src", "components", "dashboard", "CustomChartCard.tsx"):     nextToTanStack(adminCustomChartCardTSX()),
		filepath.Join(adminRoot, "src", "components", "dashboard", "ChartBuilderForm.tsx"):    nextToTanStack(adminChartBuilderFormTSX()),

		// Auth shells — the themed login/sign-up/forgot chrome. Without these the
		// Vite admin's auth pages looked nothing like the Next.js admin's.
		filepath.Join(adminRoot, "src", "components", "auth", "AuthShell.tsx"):         nextToTanStack(adminAuthShellDispatcher()),
		filepath.Join(adminRoot, "src", "components", "auth", "AtlasAuthShell.tsx"):    nextToTanStack(adminAtlasAuthShell()),
		filepath.Join(adminRoot, "src", "components", "auth", "AuroraAuthShell.tsx"):   nextToTanStack(adminAuroraAuthShell()),
		filepath.Join(adminRoot, "src", "components", "auth", "PulseAuthShell.tsx"):    nextToTanStack(adminPulseAuthShell()),
		filepath.Join(adminRoot, "src", "components", "auth", "SocialAuthButtons.tsx"): nextToTanStack(adminAuthSocialButtons()),

		// UI primitives the reused pages import.
		filepath.Join(adminRoot, "src", "components", "ui", "Skeleton.tsx"):        nextToTanStack(adminSkeletonComponent()),
		filepath.Join(adminRoot, "src", "components", "ui", "IconButton.tsx"):      nextToTanStack(adminIconButtonComponent()),
		filepath.Join(adminRoot, "src", "components", "ui", "UserCell.tsx"):        nextToTanStack(adminUserCellComponent()),
		filepath.Join(adminRoot, "src", "components", "ui", "CurrencyInput.tsx"):   nextToTanStack(adminCurrencyInputComponent()),
		filepath.Join(adminRoot, "src", "components", "ui", "ResponsiveSheet.tsx"): nextToTanStack(adminResponsiveSheetComponent()),
		filepath.Join(adminRoot, "src", "components", "ui", "ResponsiveTable.tsx"): nextToTanStack(adminResponsiveTableComponent()),
		filepath.Join(adminRoot, "src", "components", "forms", "word-editor.tsx"):  nextToTanStack(adminWordEditor()),

		// Shared components
		filepath.Join(adminRoot, "src", "components", "shared", "providers.tsx"):      nextToTanStack(adminProviders()),
		filepath.Join(adminRoot, "src", "components", "shared", "theme-provider.tsx"): nextToTanStack(adminThemeProvider()),

		// Layout components (reuse with stripped "use client")
		filepath.Join(adminRoot, "src", "components", "layout", "admin-layout.tsx"): nextToTanStack(adminLayoutComponent()),
		filepath.Join(adminRoot, "src", "components", "layout", "sidebar.tsx"):      nextToTanStack(adminSidebar()),
		filepath.Join(adminRoot, "src", "components", "layout", "navbar.tsx"):       nextToTanStack(adminNavbar()),
		filepath.Join(adminRoot, "src", "components", "layout", "page-header.tsx"):  nextToTanStack(adminPageHeader()),

		// Chrome components imported by the admin layout (collapsible sidebar,
		// quick-access menu, session watchdog) — reused from the Next.js admin.
		filepath.Join(adminRoot, "src", "components", "chrome", "CollapsibleSidebar.tsx"): nextToTanStack(adminCollapsibleSidebarComponent(opts)),
		filepath.Join(adminRoot, "src", "components", "chrome", "QuickAccess.tsx"):        nextToTanStack(adminQuickAccessComponent()),
		filepath.Join(adminRoot, "src", "components", "chrome", "SessionWatchdog.tsx"):    nextToTanStack(adminSessionWatchdogComponent()),
		// UserMenu owns the sign-out action; PageHeader is imported by every
		// system page. Both were missing, so logout did nothing and the system
		// pages could not render.
		filepath.Join(adminRoot, "src", "components", "chrome", "UserMenu.tsx"):         nextToTanStack(adminUserMenuComponent()),
		filepath.Join(adminRoot, "src", "components", "chrome", "PageHeader.tsx"):       nextToTanStack(adminPageHeaderComponent()),
		filepath.Join(adminRoot, "src", "components", "chrome", "NotificationBell.tsx"): nextToTanStack(adminNotificationBellComponent()),
		filepath.Join(adminRoot, "src", "components", "chrome", "DarkModeToggle.tsx"):   nextToTanStack(adminDarkModeToggleComponent()),

		// Table components (pure React — strip "use client")
		filepath.Join(adminRoot, "src", "components", "tables", "data-table.tsx"):        nextToTanStack(adminDataTable()),
		filepath.Join(adminRoot, "src", "components", "tables", "column-header.tsx"):     nextToTanStack(adminColumnHeader()),
		filepath.Join(adminRoot, "src", "components", "tables", "cell-renderers.tsx"):    nextToTanStack(adminCellRenderers()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-filters.tsx"):     nextToTanStack(adminTableFilters()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-toolbar.tsx"):     nextToTanStack(adminTableToolbar()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-pagination.tsx"):  nextToTanStack(adminTablePagination()),
		filepath.Join(adminRoot, "src", "components", "tables", "date-filter.tsx"):       nextToTanStack(adminDateFilter()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-skeleton.tsx"):    nextToTanStack(adminTableSkeleton()),
		filepath.Join(adminRoot, "src", "components", "tables", "table-empty-state.tsx"): nextToTanStack(adminTableEmptyState()),
		filepath.Join(adminRoot, "src", "components", "tables", "import-modal.tsx"):      nextToTanStack(adminImportModal()),
		filepath.Join(adminRoot, "src", "components", "tables", "export-menu.tsx"):       nextToTanStack(adminExportMenu()),
		// Imported by generated resource definitions when the name/email
		// column-pack heuristic fires.
		filepath.Join(adminRoot, "src", "components", "tables", "stacked-cell.tsx"): nextToTanStack(adminStackedCell()),

		// Form components (pure React — strip "use client")
		filepath.Join(adminRoot, "src", "components", "forms", "form-builder.tsx"):                              nextToTanStack(adminFormBuilder()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-modal.tsx"):                                nextToTanStack(adminFormModal()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-page.tsx"):                                 nextToTanStack(adminFormPage()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-stepper.tsx"):                              nextToTanStack(adminFormStepper()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-modal-steps.tsx"):                          nextToTanStack(adminFormModalSteps()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-page-steps.tsx"):                           nextToTanStack(adminFormPageSteps()),
		filepath.Join(adminRoot, "src", "components", "forms", "form-sheet.tsx"):                                nextToTanStack(adminFormSheet()),
		filepath.Join(adminRoot, "src", "components", "forms", "update-groups.tsx"):                             nextToTanStack(adminUpdateGroups()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "text-field.tsx"):                      nextToTanStack(adminTextField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "textarea-field.tsx"):                  nextToTanStack(adminTextareaField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "number-field.tsx"):                    nextToTanStack(adminNumberField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "select-field.tsx"):                    nextToTanStack(adminSelectField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "date-field.tsx"):                      nextToTanStack(adminDateField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "toggle-field.tsx"):                    nextToTanStack(adminToggleField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "checkbox-field.tsx"):                  nextToTanStack(adminCheckboxField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "radio-field.tsx"):                     nextToTanStack(adminRadioField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "checkbox-group-field.tsx"):           nextToTanStack(adminCheckboxGroupField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "image-field.tsx"):                     nextToTanStack(adminImageField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "images-field.tsx"):                    nextToTanStack(adminImagesField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "video-field.tsx"):                     nextToTanStack(adminVideoField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "videos-field.tsx"):                    nextToTanStack(adminVideosField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "file-field.tsx"):                      nextToTanStack(adminFileField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "files-field.tsx"):                     nextToTanStack(adminFilesField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "relationship-select-field.tsx"):       nextToTanStack(adminRelationshipSelectField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "multi-relationship-select-field.tsx"): nextToTanStack(adminMultiRelationshipSelectField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "line-items-field.tsx"):                nextToTanStack(adminLineItemsField()),
		filepath.Join(adminRoot, "src", "components", "forms", "fields", "rich-text-field.tsx"):                 nextToTanStack(adminRichTextField()),

		// UI components
		filepath.Join(adminRoot, "src", "components", "ui", "dropzone.tsx"):      nextToTanStack(adminDropzone()),
		filepath.Join(adminRoot, "src", "components", "ui", "confirm-modal.tsx"): nextToTanStack(adminConfirmModal()),

		// Widget components
		filepath.Join(adminRoot, "src", "components", "widgets", "stats-card.tsx"):      nextToTanStack(adminStatsCard()),
		filepath.Join(adminRoot, "src", "components", "widgets", "chart-widget.tsx"):    nextToTanStack(adminChartWidget()),
		filepath.Join(adminRoot, "src", "components", "widgets", "activity-widget.tsx"): nextToTanStack(adminActivityWidget()),
		filepath.Join(adminRoot, "src", "components", "widgets", "widget-grid.tsx"):     nextToTanStack(adminWidgetGrid()),

		// Resource components
		filepath.Join(adminRoot, "src", "components", "resource", "resource-page.tsx"):        nextToTanStack(adminResourcePage()),
		filepath.Join(adminRoot, "src", "components", "resource", "resource-detail-page.tsx"): nextToTanStack(adminResourceDetailPage()),
		filepath.Join(adminRoot, "src", "components", "resource", "view-modal.tsx"):    nextToTanStack(adminViewModal()),

		// Resource definitions (same as Next.js)
		filepath.Join(adminRoot, "src", "resources", "index.ts"): adminResourceRegistry(),
		filepath.Join(adminRoot, "src", "resources", "users.ts"): adminUsersResource(),
		filepath.Join(adminRoot, "src", "resources", "blogs.ts"): adminBlogsResource(),

		// Profile
		filepath.Join(adminRoot, "src", "components", "profile", "delete-account-dialog.tsx"): nextToTanStack(adminDeleteAccountDialog()),
		filepath.Join(adminRoot, "src", "components", "profile", "active-sessions.tsx"):       nextToTanStack(adminActiveSessions()),

		filepath.Join(adminRoot, "public", ".gitkeep"): "",
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func adminTanStackPackageJSON(opts Options) string {
	return fmt.Sprintf(`{
  "name": "@%s/admin",
  "version": "0.1.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite --port 3001",
    "build": "vite build",
    "typecheck": "tsc -b",
    "preview": "vite preview",
    "lint": "eslint .",
    "test": "vitest run",
    "test:watch": "vitest"
  },
  "dependencies": {
    "@hookform/resolvers": "^3.9.1",
    "@react-pdf/renderer": "^4.1.5",
    "@tanstack/react-query": "^5.62.0",
    "@tanstack/react-router": "^1.93.0",
    "@tanstack/react-table": "^8.20.6",
    "@tiptap/extension-link": "^2.1.0",
    "@tiptap/extension-text-align": "^2.1.0",
    "@tiptap/extension-text-style": "^2.1.0",
    "@tiptap/extension-color": "^2.1.0",
    "@tiptap/extension-highlight": "^2.1.0",
    "@tiptap/extension-underline": "^2.1.0",
    "@tiptap/extension-image": "^2.1.0",
    "@tiptap/extension-table": "^2.1.0",
    "@tiptap/extension-table-row": "^2.1.0",
    "@tiptap/extension-table-cell": "^2.1.0",
    "@tiptap/extension-table-header": "^2.1.0",
    "@tiptap/extension-placeholder": "^2.1.0",
    "@tiptap/pm": "^2.1.0",
    "@tiptap/react": "^2.1.0",
    "@tiptap/starter-kit": "^2.1.0",
    "axios": "^1.7.9",
    "clsx": "^2.1.1",
    "lucide-react": "^0.468.0",
    "react": "19.2.7",
    "react-dom": "19.2.7",
    "react-dropzone": "^14.2.0",
    "react-hook-form": "^7.54.1",
    "recharts": "^2.15.0",
    "sonner": "^1.3.0",
    "tailwind-merge": "^2.6.0",
    "xlsx": "^0.18.5",
    "zod": "^3.24.1",
    "@repo/shared": "workspace:*"
  },
  "devDependencies": {
    "@tanstack/react-router-devtools": "^1.93.0",
    "@tanstack/router-vite-plugin": "^1.93.0",
    "@types/react": "^19.0.0",
    "@types/react-dom": "^19.0.0",
    "tailwindcss": "^4.1.13",
    "@tailwindcss/vite": "^4.1.13",
    "tw-animate-css": "^1.4.0",
    "typescript": "~5.7.0",
    "vite": "^6.0.0",
    "@vitejs/plugin-react": "^4.3.4",
    "vitest": "^2.1.0",
    "@testing-library/react": "^16.1.0"
  }
}`, opts.ProjectName)
}

func adminTanStackViteConfig() string {
	return `import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-vite-plugin'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
` + viteSecurityHeaders() + `
export default defineConfig({
  plugins: [
    TanStackRouterVite(),
    react(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  // Tailwind v4 is handled by @tailwindcss/vite above. Pin an empty inline
  // PostCSS config so Vite does NOT walk up and pick up the monorepo-root
  // postcss.config.mjs (v3-style, meant for the Next.js apps), which would
  // try to load tailwindcss as a PostCSS plugin and fail under v4.
  css: {
    postcss: {},
  },
  // Security headers, mirroring the Next.js admin + the Go API's
  // middleware.SecurityHeaders. Vite only applies these when IT serves the
  // files (dev + preview) — a production deploy serves dist/ from a reverse
  // proxy or CDN, so set the same headers there too.
  preview: {
    headers: securityHeaders,
  },
  server: {
    headers: securityHeaders,
    port: 3001,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
`
}

// adminTanStackIndexHTML mirrors the Next.js admin's root layout so both admins
// render identically for a given theme.
//
// Two things must match, and previously neither did:
//   - data-theme on <html> selects the palette in globals.css. The Vite admin
//     used to hard-code class="dark" and set no data-theme, which pinned it to
//     the .dark override forever — so an atlas project (a LIGHT theme) rendered
//     dark and looked nothing like the Next.js admin. .dark is applied at
//     runtime by DarkModeToggle, never baked in.
//   - the theme's fonts populate --font-display / --font-mono / --font-serif,
//     which globals.css uses for body text. The Next.js admin loads them with
//     next/font; Vite has no equivalent, so we link Google Fonts and set the
//     same variables. Without this the admin silently fell back to system-ui.
func adminTanStackIndexHTML(opts Options) string {
	theme := opts.Theme
	if theme == "" {
		theme = "atlas"
	}

	var fontLink, fontVars string
	switch theme {
	case "aurora":
		fontLink = `<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=Geist+Mono:wght@400;500;600&display=swap" rel="stylesheet" />`
		fontVars = "--font-display: 'Geist', system-ui, sans-serif;\n      --font-mono: 'Geist Mono', ui-monospace, monospace;"
	case "pulse":
		// Cloudflare-inspired: clean sans throughout, no serif display face.
		fontLink = `<link href="https://fonts.googleapis.com/css2?family=Onest:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet" />`
		fontVars = "--font-display: 'Onest', system-ui, sans-serif;\n      --font-mono: 'JetBrains Mono', ui-monospace, monospace;"
	default: // atlas
		fontLink = `<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500;600&display=swap" rel="stylesheet" />`
		fontVars = "--font-display: 'Inter', system-ui, sans-serif;\n      --font-mono: 'JetBrains Mono', ui-monospace, monospace;"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en" data-theme="%s">
  <head>
    <meta charset="UTF-8" />
    <link rel="icon" type="image/svg+xml" href="/vite.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    %s
    <style>
      :root {
      %s
      }
    </style>
    <title>%s — Admin</title>
  </head>
  <body class="min-h-screen font-sans antialiased">
    <div id="root"></div>
    <script type="module" src="/src/main.tsx"></script>
  </body>
</html>
`, theme, fontLink, fontVars, opts.ProjectName)
}

func adminTanStackTSConfig() string {
	return webTanStackTSConfig() // Same config
}

func adminTanStackMain() string {
	return `import React from 'react'
import ReactDOM from 'react-dom/client'
import { RouterProvider, createRouter } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { routeTree } from './routeTree.gen'
import './globals.css'

// Mirrors the Next.js admin's NEXT_PUBLIC_THEME: index.html bakes in the theme
// picked at scaffold time, and VITE_THEME overrides it at runtime, so changing
// the theme in .env repaints the dashboard without re-scaffolding.
const envTheme = (import.meta as any).env?.VITE_THEME
if (envTheme) {
  document.documentElement.setAttribute('data-theme', envTheme)
}

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      retry: 1,
    },
  },
})

const router = createRouter({ routeTree })

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>
  </React.StrictMode>,
)
`
}

func adminTanStackRootRoute() string {
	return `import { createRootRoute, Outlet } from '@tanstack/react-router'
import NotFound from '@/pages/not-found'

export const Route = createRootRoute({
  component: () => <Outlet />,
  // Unmatched URLs render the same branded 404 the Next.js admin ships,
  // instead of TanStack's bare "Not Found" text.
  notFoundComponent: NotFound,
})
`
}

func adminTanStackRedirectRoute() string {
	return `import { createFileRoute, redirect } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  beforeLoad: () => {
    throw redirect({ to: '/dashboard' })
  },
})
`
}

func adminTanStackAuthLayout() string {
	return `import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'
import { apiClient } from '@/lib/api-client'

export const Route = createFileRoute('/_auth')({
  // Auth tokens live in HttpOnly cookies, so JS cannot read them — a
  // localStorage check here is always empty and would show the login form to
  // an already-signed-in user. Ask the API instead: /api/auth/me succeeds
  // when the cookie is valid and 401s when it isn't.
  beforeLoad: async () => {
    try {
      await apiClient.get('/api/auth/me')
    } catch {
      return // not signed in — show the auth screens
    }
    throw redirect({ to: '/dashboard' })
  },
  component: () => (
    <div className="min-h-screen flex items-center justify-center bg-background">
      <Outlet />
    </div>
  ),
})
`
}

func adminTanStackDashboardLayout() string {
	return `import { createFileRoute, Outlet, redirect } from '@tanstack/react-router'
import { AdminLayout } from '@/components/layout/admin-layout'
import { apiClient } from '@/lib/api-client'

export const Route = createFileRoute('/_dashboard')({
  // Auth tokens live in HttpOnly cookies, so JS cannot read them. Reading
  // localStorage here always returned null, which bounced every signed-in
  // user straight back to /login and made the whole admin unreachable.
  // Ask the API instead — the cookie rides along on the request.
  beforeLoad: async () => {
    try {
      await apiClient.get('/api/auth/me')
    } catch {
      throw redirect({ to: '/login' })
    }
  },
  component: () => (
    <AdminLayout>
      <Outlet />
    </AdminLayout>
  ),
})
`
}
