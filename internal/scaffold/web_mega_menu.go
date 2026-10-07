package scaffold

import "fmt"

// The marketing navbar's mega menu.
//
// Built on Base UI's NavigationMenu rather than by hand, which is what
// CLAUDE.md asks for when a new primitive needs behaviour rather than
// appearance. This one needs all of it: a portalled popup that survives the
// sticky header's stacking context, a viewport that animates between two
// panels of different sizes, hover intent so crossing one trigger on the way to
// another does not flicker, Escape and click-outside to dismiss, and arrow-key
// movement through the links. Every one of those is a thing hand-rolled menus
// get wrong for keyboard and screen-reader users specifically.
//
// Base UI is already a dependency of all four frontends (the country picker in
// the admin uses its combobox), so this adds a component and no package.
//
// The content is data, in lib/nav-menu.ts, because a mega menu full of
// plausible-looking links nobody wrote is the thing that tells every visitor
// the site is a template. The default below names only pages the project
// actually has.

// webNavMenuConfig emits lib/nav-menu.ts: what is in the menus.
func webNavMenuConfig(opts Options) string {
	// Where the panel is, resolved here rather than left as {{ADMIN_HREF}}:
	// the Vite writers do not run the substitution the web writer does, so a
	// single would have shipped the placeholder as a literal URL.
	//
	// In-app when the panel is a route of this same application, which is a
	// double (a route group) or a single (a section of the SPA). Then it is a
	// route and not a new tab.
	inApp := opts.ShouldEmbedAdmin() || opts.ShouldEmbedAdminInSPA()
	adminConst := `// Where the admin panel is: a route of this same app.
const ADMIN_URL = "/admin/dashboard";`
	if !inApp {
		adminConst = `// Where the admin panel is. It runs as its own app, so this is an origin
// and not a path. NEXT_PUBLIC_ADMIN_URL overrides it.
const ADMIN_URL = process.env.NEXT_PUBLIC_ADMIN_URL || "http://localhost:3001";`
	}
	external := "true"
	if inApp {
		external = "false"
	}
	return fmt.Sprintf(`import {
  BookOpen,
  LayoutDashboard,
  LifeBuoy,
  Rss,
  Shield,
  Sparkles,
  type LucideIcon,
} from "lucide-react";

// The navbar's contents. Edit this file: the menus below name the pages a new
// %s has, which is a starting point and not a sitemap.
//
// An entry with no columns is a plain link. An entry with columns is a mega
// menu: two or three columns of links, and an optional panel on the right for
// whatever you want people to see first.

%s

export interface MegaMenuLink {
  title: string;
  description: string;
  href: string;
  icon: LucideIcon;
  /** Opens in a new tab, and is announced as doing so. */
  external?: boolean;
}

export interface MegaMenuColumn {
  heading: string;
  links: MegaMenuLink[];
}

/** The panel on the right of an open menu. Omit it for a plain two-column menu. */
export interface MegaMenuFeature {
  eyebrow: string;
  title: string;
  description: string;
  href: string;
}

export type NavEntry =
  | { label: string; href: string }
  | { label: string; columns: MegaMenuColumn[]; feature?: MegaMenuFeature };

export function isMegaMenu(
  entry: NavEntry,
): entry is { label: string; columns: MegaMenuColumn[]; feature?: MegaMenuFeature } {
  return "columns" in entry;
}

export const navEntries: NavEntry[] = [
  { label: "Home", href: "/" },
  {
    label: "Product",
    columns: [
      {
        heading: "Explore",
        links: [
          {
            title: "Blog",
            description: "Posts, written and published from the admin panel.",
            href: "/blog",
            icon: Rss,
          },
          {
            title: "Admin panel",
            description: "Tables, forms and filters for every resource.",
            href: ADMIN_URL,
            icon: LayoutDashboard,
            external: %s,
          },
        ],
      },
      {
        heading: "Built in",
        links: [
          {
            title: "Accounts",
            description: "Sign-in, roles and sessions, already wired up.",
            href: ADMIN_URL,
            icon: Shield,
            external: %s,
          },
          {
            title: "Search",
            description: "Across every resource you generate.",
            href: "/blog",
            icon: Sparkles,
          },
        ],
      },
    ],
    feature: {
      eyebrow: "Start here",
      title: "Add your own pages",
      description:
        "This menu is lib/nav-menu.ts. Change the columns, the copy and this panel to say what %s does.",
      href: "/blog",
    },
  },
  {
    label: "Resources",
    columns: [
      {
        heading: "Learn",
        links: [
          {
            title: "Blog",
            description: "Everything published so far.",
            href: "/blog",
            icon: BookOpen,
          },
          {
            title: "Support",
            description: "Replace this with how people reach you.",
            href: "/blog",
            icon: LifeBuoy,
          },
        ],
      },
    ],
  },
];
`, opts.ProjectName, adminConst, external, external, opts.ProjectName)
}

// webMegaMenu emits components/mega-menu.tsx: the desktop menu itself.
func webMegaMenu() string {
	return `"use client";

import Link from "next/link";
import { NavigationMenu } from "@base-ui/react/navigation-menu";
import { ChevronDown, ArrowRight } from "lucide-react";

import { navEntries, isMegaMenu, type MegaMenuLink } from "@/lib/nav-menu";

// One link in an open panel: the icon tile, the title, the one-line
// description. A whole row is the target, not just the title, because a 14px
// line of text is below the 24px minimum a pointer target has to meet.
function PanelLink({ link }: { link: MegaMenuLink }) {
  const Icon = link.icon;
  const body = (
    <>
      <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-accent/20 bg-accent/10">
        <Icon className="h-4 w-4 text-accent" aria-hidden="true" />
      </span>
      <span className="min-w-0">
        <span className="block text-sm font-semibold text-foreground">{link.title}</span>
        <span className="mt-0.5 block text-sm text-text-muted">{link.description}</span>
      </span>
    </>
  );
  const className =
    "flex items-start gap-3 rounded-lg p-3 transition-colors hover:bg-bg-hover focus-visible:bg-bg-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent";

  if (link.external) {
    return (
      <NavigationMenu.Link
        render={
          <a href={link.href} target="_blank" rel="noopener noreferrer" className={className} />
        }
      >
        {body}
        <span className="sr-only"> (opens in a new tab)</span>
      </NavigationMenu.Link>
    );
  }
  return (
    <NavigationMenu.Link render={<Link href={link.href} className={className} />}>
      {body}
    </NavigationMenu.Link>
  );
}

/** The marketing navbar's links, as a mega menu. Hidden below md. */
export function MegaMenu() {
  return (
    <NavigationMenu.Root className="hidden md:flex" delay={100} closeDelay={120}>
      <NavigationMenu.List className="flex items-center gap-1">
        {navEntries.map((entry) =>
          isMegaMenu(entry) ? (
            <NavigationMenu.Item key={entry.label}>
              <NavigationMenu.Trigger className="flex items-center gap-1 rounded-lg px-3 py-2 text-sm text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground data-[popup-open]:bg-bg-hover data-[popup-open]:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent">
                {entry.label}
                <NavigationMenu.Icon className="transition-transform duration-200 data-[popup-open]:rotate-180">
                  <ChevronDown className="h-3.5 w-3.5" aria-hidden="true" />
                </NavigationMenu.Icon>
              </NavigationMenu.Trigger>

              <NavigationMenu.Content className="w-full p-6 transition-[opacity,transform] duration-200 data-[starting-style]:opacity-0 data-[ending-style]:opacity-0">
                <div className="flex gap-8">
                  <div className="grid flex-1 gap-8 sm:grid-cols-2">
                    {entry.columns.map((column) => (
                      <div key={column.heading}>
                        <p className="px-3 pb-2 text-xs font-semibold uppercase tracking-wider text-text-muted">
                          {column.heading}
                        </p>
                        <ul className="space-y-1">
                          {column.links.map((link) => (
                            <li key={link.title}>
                              <PanelLink link={link} />
                            </li>
                          ))}
                        </ul>
                      </div>
                    ))}
                  </div>

                  {entry.feature && (
                    <div className="hidden w-64 shrink-0 rounded-xl border border-border bg-bg-secondary p-5 lg:block">
                      <p className="text-xs font-semibold uppercase tracking-wider text-text-muted">
                        {entry.feature.eyebrow}
                      </p>
                      <p className="mt-3 text-sm font-semibold text-foreground">
                        {entry.feature.title}
                      </p>
                      <p className="mt-1.5 text-sm text-text-muted">{entry.feature.description}</p>
                      <NavigationMenu.Link
                        render={
                          <Link
                            href={entry.feature.href}
                            className="mt-4 inline-flex items-center gap-1.5 text-sm font-medium text-accent hover:text-accent-hover"
                          />
                        }
                      >
                        Take a look
                        <ArrowRight className="h-3.5 w-3.5" aria-hidden="true" />
                      </NavigationMenu.Link>
                    </div>
                  )}
                </div>
              </NavigationMenu.Content>
            </NavigationMenu.Item>
          ) : (
            <NavigationMenu.Item key={entry.label}>
              <NavigationMenu.Link
                render={
                  <Link
                    href={entry.href}
                    className="block rounded-lg px-3 py-2 text-sm text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent"
                  />
                }
              >
                {entry.label}
              </NavigationMenu.Link>
            </NavigationMenu.Item>
          ),
        )}
      </NavigationMenu.List>

      {/* Portalled: the header is sticky and creates a stacking context, so a
          panel rendered inside it is clipped by the header's own height. */}
      <NavigationMenu.Portal>
        <NavigationMenu.Positioner
          sideOffset={10}
          collisionPadding={16}
          className="z-50 box-border max-w-[min(var(--available-width),56rem)] transition-[top,left,right,bottom] duration-200"
        >
          <NavigationMenu.Popup className="relative w-[var(--popup-width)] origin-[var(--transform-origin)] overflow-hidden rounded-xl border border-border bg-background shadow-2xl transition-[opacity,transform,width,height] duration-200 data-[starting-style]:scale-95 data-[starting-style]:opacity-0 data-[ending-style]:scale-95 data-[ending-style]:opacity-0">
            <NavigationMenu.Viewport />
          </NavigationMenu.Popup>
        </NavigationMenu.Positioner>
      </NavigationMenu.Portal>
    </NavigationMenu.Root>
  );
}

/** The same entries in the mobile drawer, where a panel has nowhere to go. */
export function MegaMenuMobile({ onNavigate }: { onNavigate?: () => void }) {
  return (
    <div className="flex flex-col gap-1">
      {navEntries.map((entry) =>
        isMegaMenu(entry) ? (
          <details key={entry.label} className="group">
            <summary className="flex cursor-pointer items-center justify-between rounded-lg px-2 py-2 text-sm text-text-secondary hover:text-foreground">
              {entry.label}
              <ChevronDown
                className="h-4 w-4 transition-transform group-open:rotate-180"
                aria-hidden="true"
              />
            </summary>
            <div className="mt-1 space-y-3 border-l border-border pl-3">
              {entry.columns.map((column) => (
                <div key={column.heading}>
                  <p className="px-2 pb-1 text-xs font-semibold uppercase tracking-wider text-text-muted">
                    {column.heading}
                  </p>
                  {column.links.map((link) =>
                    link.external ? (
                      <a
                        key={link.title}
                        href={link.href}
                        target="_blank"
                        rel="noopener noreferrer"
                        onClick={onNavigate}
                        className="block rounded-lg px-2 py-2 text-sm text-text-secondary hover:text-foreground"
                      >
                        {link.title}
                        <span className="sr-only"> (opens in a new tab)</span>
                      </a>
                    ) : (
                      <Link
                        key={link.title}
                        href={link.href}
                        onClick={onNavigate}
                        className="block rounded-lg px-2 py-2 text-sm text-text-secondary hover:text-foreground"
                      >
                        {link.title}
                      </Link>
                    ),
                  )}
                </div>
              ))}
            </div>
          </details>
        ) : (
          <Link
            key={entry.label}
            href={entry.href}
            onClick={onNavigate}
            className="rounded-lg px-2 py-2 text-sm text-text-secondary hover:text-foreground"
          >
            {entry.label}
          </Link>
        ),
      )}
    </div>
  );
}
`
}
