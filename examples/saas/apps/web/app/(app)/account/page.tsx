"use client";

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
