"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Activity, CreditCard, LayoutDashboard, LogOut, User as UserIcon } from "lucide-react";
import { ProtectedWebRoute } from "@/components/ProtectedWebRoute";
import { useMe, useLogout } from "@/hooks/use-auth";

// The signed-in customer area. Anything under app/(app) renders inside this: a
// sidebar, a header that says who is signed in, and none of the marketing chrome.
//
// The group name is in parentheses, so it is not part of the URL. These pages are
// at /account and /account/profile, which is what middleware.ts protects.
//
// Add a section by adding a line here and a page under app/(app)/. Orders,
// addresses, invoices, whatever the product has: the shell does not change.
const NAV = [
  { href: "/account", label: "Overview", icon: LayoutDashboard },
  { href: "/account/billing", label: "Billing", icon: CreditCard },
  { href: "/account/usage", label: "Usage", icon: Activity },
  { href: "/account/profile", label: "Profile", icon: UserIcon },
];

function initials(first?: string | null, last?: string | null, email?: string | null) {
  const a = (first ?? "").trim();
  const b = (last ?? "").trim();
  if (a || b) return ((a[0] ?? "") + (b[0] ?? "")).toUpperCase();
  return (email ?? "?").slice(0, 2).toUpperCase();
}

export default function AccountLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname() ?? "";
  const { data: user } = useMe();
  const logout = useLogout();

  const isActive = (href: string) =>
    href === "/account" ? pathname === "/account" : pathname.startsWith(href);

  const navLinks = NAV.map((item) => {
    const Icon = item.icon;
    const active = isActive(item.href);
    return (
      <Link
        key={item.href}
        href={item.href}
        className={
          "flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm transition-colors " +
          (active
            ? "bg-accent/10 font-medium text-accent"
            : "text-text-secondary hover:bg-bg-hover hover:text-foreground")
        }
      >
        <Icon className="h-4.5 w-4.5 shrink-0" />
        {item.label}
      </Link>
    );
  });

  return (
    <ProtectedWebRoute>
      <div className="flex min-h-screen bg-bg-secondary">
        {/* Sidebar */}
        <aside className="hidden w-64 shrink-0 flex-col border-r border-border bg-background md:flex">
          <Link href="/" className="flex h-20 items-center gap-3 px-6">
            <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-sm font-bold text-accent">
              S
            </span>
            <span className="text-lg font-bold tracking-tight text-foreground">saas</span>
          </Link>

          <nav className="flex flex-1 flex-col gap-1 px-4">{navLinks}</nav>

          <div className="border-t border-border p-4">
            <button
              type="button"
              onClick={() => logout.mutate()}
              className="flex w-full items-center gap-3 rounded-lg px-3 py-2.5 text-sm text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground"
            >
              <LogOut className="h-4.5 w-4.5 shrink-0" />
              Sign out
            </button>
          </div>
        </aside>

        <div className="flex min-h-screen flex-1 flex-col">
          {/* Header */}
          <header className="flex h-20 items-center justify-between border-b border-border bg-background px-6">
            <div className="min-w-0">
              <p className="truncate text-sm font-semibold text-foreground">
                Welcome, {user?.first_name || user?.email || "there"}
              </p>
              <p className="truncate text-xs text-text-secondary">{user?.email}</p>
            </div>
            <div className="flex items-center gap-3">
              <Link
                href="/"
                className="hidden text-sm text-text-secondary transition-colors hover:text-foreground sm:block"
              >
                Back to site
              </Link>
              <span className="flex h-10 w-10 items-center justify-center rounded-full bg-accent/10 text-sm font-semibold text-accent">
                {initials(user?.first_name, user?.last_name, user?.email)}
              </span>
            </div>
          </header>

          {/* On a phone the sidebar is gone, so the sections ride under the header. */}
          <nav className="flex gap-2 overflow-x-auto border-b border-border bg-background px-4 py-2 md:hidden">
            {navLinks}
            <button
              type="button"
              onClick={() => logout.mutate()}
              className="flex items-center gap-2 rounded-lg px-3 py-2.5 text-sm text-text-secondary"
            >
              <LogOut className="h-4.5 w-4.5" />
              Sign out
            </button>
          </nav>

          <main className="flex-1 px-6 py-8 lg:px-10">
            <div className="mx-auto max-w-4xl">{children}</div>
          </main>
        </div>
      </div>
    </ProtectedWebRoute>
  );
}
