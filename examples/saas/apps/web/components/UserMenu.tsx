"use client";

// v3.31.42 -- UserMenu shows different chrome in the navbar
// depending on whether the visitor is signed in:
//   - signed in : avatar dropdown with Account + Sign out
//   - signed out: Log in + Sign up buttons
//   - loading   : a small placeholder so the navbar doesn't shift
//                 layout once useMe() resolves
//
// The signed-in / signed-out decision uses the same useMe() that
// every other auth-aware page reads, so it picks up login + logout
// immediately via React Query's cache.

import { useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useMe, useLogout } from "@/hooks/use-auth";
import { ChevronDown, LogOut, User as UserIcon } from "lucide-react";

export function UserMenu() {
  const { data: user, isLoading } = useMe();
  const logout = useLogout();
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);

  // Close on outside click. Attached only while the dropdown is open
  // so we don't pay the cost on every page.
  useEffect(() => {
    if (!open) return;
    function handle(e: MouseEvent) {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    }
    document.addEventListener("mousedown", handle);
    return () => document.removeEventListener("mousedown", handle);
  }, [open]);

  if (isLoading) {
    // Reserve roughly the same width as the signed-out CTA pair so
    // the navbar doesn't jump on auth resolve.
    return <div className="h-8 w-32 rounded-md bg-bg-elevated/40 animate-pulse" />;
  }

  if (!user) {
    return (
      <div className="flex items-center gap-2">
        <Link
          href="/login"
          className="rounded-lg border border-border bg-transparent px-3 py-1.5 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
        >
          Log in
        </Link>
        <Link
          href="/register"
          className="rounded-lg bg-accent px-3 py-1.5 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
        >
          Sign up
        </Link>
      </div>
    );
  }

  const firstName = (user as { first_name?: string }).first_name || "";
  const lastName = (user as { last_name?: string }).last_name || "";
  const email = (user as { email?: string }).email || "";
  const fullName =
    [firstName, lastName].filter(Boolean).join(" ") || email || "Account";
  const initials =
    ((firstName[0] || "") + (lastName[0] || "")).toUpperCase() ||
    email.slice(0, 2).toUpperCase() ||
    "U";

  return (
    <div ref={rootRef} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        className="flex items-center gap-2 rounded-lg border border-border bg-bg-elevated/60 px-2 py-1.5 text-sm text-foreground hover:bg-bg-hover transition-colors"
      >
        <span className="inline-flex h-7 w-7 items-center justify-center rounded-full bg-accent/15 text-accent text-xs font-semibold">
          {initials}
        </span>
        <span className="hidden sm:inline max-w-[100px] truncate">{fullName}</span>
        <ChevronDown className="h-3.5 w-3.5 text-text-muted" />
      </button>

      {open && (
        <div className="absolute right-0 top-full z-50 mt-1 w-56 rounded-lg border border-border bg-bg-elevated shadow-lg">
          <div className="border-b border-border px-3 py-2.5">
            <p className="truncate text-sm font-medium text-foreground">{fullName}</p>
            {email && (
              <p className="truncate text-xs text-text-muted">{email}</p>
            )}
          </div>
          <div className="p-1">
            <Link
              href="/account"
              onClick={() => setOpen(false)}
              className="flex items-center gap-2 rounded-md px-2.5 py-2 text-sm text-foreground hover:bg-bg-hover transition-colors"
            >
              <UserIcon className="h-3.5 w-3.5 text-text-muted" />
              Account
            </Link>
            <button
              type="button"
              onClick={() => {
                setOpen(false);
                logout.mutate();
              }}
              className="flex w-full items-center gap-2 rounded-md px-2.5 py-2 text-left text-sm text-foreground hover:bg-bg-hover transition-colors"
            >
              <LogOut className="h-3.5 w-3.5 text-text-muted" />
              Sign out
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
