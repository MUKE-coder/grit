package scaffold

import (
	"fmt"
	"strings"
)

// The customer area: the third kind of page a web app serves.
//
// The web app had two layouts and needed three. The public site is one thing,
// the admin panel another, and the pages a signed-in customer sees are a third:
// their account, their profile, whatever the product gives them once they have
// logged in. Those had no home at all. The navbar's user menu linked to /account
// and the middleware protected /account, and the page did not exist: signing in
// and clicking your own name was a 404.
//
// It lives in app/(app)/, a route group, so it has a layout of its own without
// changing a single URL. Written by `grit add web-auth`, because a customer area
// with no way to log in is furniture.
//
// The shape is a dashboard, not a page with a back link: a sidebar that names the
// product and lists the sections, a header that says who is signed in, and the
// content on a tinted ground. Sections are added by adding to one array.

// webAccountLayout is the shell around every signed-in customer page.
func webAccountLayout(opts Options) string {
	initial := "A"
	if opts.ProjectName != "" {
		initial = strings.ToUpper(string([]rune(opts.ProjectName)[0]))
	}
	return fmt.Sprintf(`"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutDashboard, LogOut, User as UserIcon } from "lucide-react";
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
              %s
            </span>
            <span className="text-lg font-bold tracking-tight text-foreground">%s</span>
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
`, initial, opts.ProjectName)
}

// webAccountOverviewPage is what a customer lands on after signing in.
func webAccountOverviewPage() string {
	return `"use client";

import Link from "next/link";
import { Mail, ShieldCheck, CalendarDays } from "lucide-react";
import { useMe } from "@/hooks/use-auth";

// The account overview. Deliberately small: it reads the session and shows what
// the account is, with a way into everything else. Add your product's own cards
// here (orders, subscriptions, usage) rather than starting a new shell.
export default function AccountOverviewPage() {
  const { data: user, isLoading } = useMe();

  if (isLoading) {
    return <p className="text-sm text-text-secondary">Loading your account...</p>;
  }

  const fullName = [user?.first_name, user?.last_name].filter(Boolean).join(" ") || "--";
  const joined = user?.created_at
    ? new Date(user.created_at).toLocaleDateString(undefined, {
        year: "numeric",
        month: "long",
        day: "numeric",
      })
    : "--";

  const cards = [
    { icon: Mail, label: "Email", value: user?.email ?? "--" },
    { icon: ShieldCheck, label: "Role", value: user?.role ?? "USER" },
    { icon: CalendarDays, label: "Member since", value: joined },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <ShieldCheck className="h-5 w-5" />
        </span>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">Overview</h1>
          <p className="text-sm text-text-secondary">Your account at a glance</p>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-3">
        {cards.map((card) => {
          const Icon = card.icon;
          return (
            <div
              key={card.label}
              className="rounded-xl border border-border bg-background p-5"
            >
              <div className="flex items-center gap-2 text-text-secondary">
                <Icon className="h-4 w-4" />
                <p className="text-xs font-semibold uppercase tracking-wide">{card.label}</p>
              </div>
              <p className="mt-3 truncate text-sm text-foreground">{card.value}</p>
            </div>
          );
        })}
      </div>

      <div className="rounded-xl border border-border bg-background">
        <div className="border-b border-border px-6 py-4">
          <h2 className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
            Profile
          </h2>
        </div>
        <dl className="divide-y divide-border/60">
          {[
            ["Full name", fullName],
            ["Email", user?.email ?? "--"],
            ["Job title", user?.job_title || "--"],
          ].map(([label, value]) => (
            <div key={label} className="flex items-center justify-between gap-4 px-6 py-4">
              <dt className="text-xs font-semibold uppercase tracking-wide text-text-secondary">
                {label}
              </dt>
              <dd className="truncate text-sm text-foreground">{value}</dd>
            </div>
          ))}
        </dl>
        <div className="border-t border-border px-6 py-4">
          <Link
            href="/account/profile"
            className="inline-flex items-center rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-accent-hover"
          >
            Edit profile
          </Link>
        </div>
      </div>
    </div>
  );
}
`
}

// webAccountProfilePage lets a customer edit their own record.
func webAccountProfilePage() string {
	return `"use client";

import { useEffect, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { UserCog } from "lucide-react";
import { api } from "@/lib/api";
import { useMe } from "@/hooks/use-auth";

// One endpoint does all of this: PUT /api/profile takes the name fields, the
// email, and a password when the customer wants to change it. An empty password
// is left out of the payload rather than sent as an empty string, which the API
// would hash.
export default function AccountProfilePage() {
  const { data: user } = useMe();
  const queryClient = useQueryClient();

  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    email: "",
    job_title: "",
    password: "",
  });
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    if (!user) return;
    setForm({
      first_name: user.first_name ?? "",
      last_name: user.last_name ?? "",
      email: user.email ?? "",
      job_title: user.job_title ?? "",
      password: "",
    });
  }, [user]);

  const save = useMutation({
    mutationFn: async () => {
      const payload: Record<string, string> = {
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        job_title: form.job_title,
      };
      if (form.password) payload.password = form.password;
      const { data } = await api.put("/api/profile", payload);
      return data;
    },
    onSuccess: () => {
      setSaved(true);
      setForm((f) => ({ ...f, password: "" }));
      queryClient.invalidateQueries({ queryKey: ["me"] });
      window.setTimeout(() => setSaved(false), 4000);
    },
  });

  const field = (key: keyof typeof form) => ({
    value: form[key],
    onChange: (e: React.ChangeEvent<HTMLInputElement>) =>
      setForm((f) => ({ ...f, [key]: e.target.value })),
    className:
      "w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none transition-colors focus:border-accent",
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <UserCog className="h-5 w-5" />
        </span>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">Profile</h1>
          <p className="text-sm text-text-secondary">
            This is what the rest of the app knows about you
          </p>
        </div>
      </div>

      <form
        className="space-y-5 rounded-xl border border-border bg-background p-6"
        onSubmit={(e) => {
          e.preventDefault();
          save.mutate();
        }}
      >
        <div className="grid gap-5 sm:grid-cols-2">
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">First name</span>
            <input {...field("first_name")} />
          </label>
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">Last name</span>
            <input {...field("last_name")} />
          </label>
        </div>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Email</span>
          <input type="email" {...field("email")} />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Job title</span>
          <input {...field("job_title")} />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">New password</span>
          <input type="password" autoComplete="new-password" {...field("password")} />
          <span className="block text-xs text-text-secondary">
            Leave this empty to keep the password you have.
          </span>
        </label>

        <div className="flex items-center gap-3 pt-1">
          <button
            type="submit"
            disabled={save.isPending}
            className="inline-flex items-center rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-accent-hover disabled:opacity-60"
          >
            {save.isPending ? "Saving..." : "Save changes"}
          </button>
          {saved ? <span className="text-sm text-text-secondary">Saved.</span> : null}
          {save.isError ? (
            <span className="text-sm text-danger">
              That did not save. Check the fields and try again.
            </span>
          ) : null}
        </div>
      </form>
    </div>
  );
}
`
}
