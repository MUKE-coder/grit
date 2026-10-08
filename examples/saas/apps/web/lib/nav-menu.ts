import {
  BookOpen,
  LayoutDashboard,
  LifeBuoy,
  Rss,
  Shield,
  Sparkles,
  type LucideIcon,
} from "lucide-react";

// The navbar's contents. Edit this file: the menus below name the pages a new
// saas has, which is a starting point and not a sitemap.
//
// An entry with no columns is a plain link. An entry with columns is a mega
// menu: two or three columns of links, and an optional panel on the right for
// whatever you want people to see first.

// Where the admin panel is. It runs as its own app, so this is an origin
// and not a path. NEXT_PUBLIC_ADMIN_URL overrides it.
const ADMIN_URL = process.env.NEXT_PUBLIC_ADMIN_URL || "http://localhost:3001";

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
            external: true,
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
            external: true,
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
        "This menu is lib/nav-menu.ts. Change the columns, the copy and this panel to say what saas does.",
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
