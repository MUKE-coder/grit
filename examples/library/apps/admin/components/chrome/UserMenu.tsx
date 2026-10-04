"use client";

import Link from "next/link";
import { useEffect, useRef, useState } from "react";
import { useMe, useLogout } from "@/hooks/use-auth";
import { Activity, User as UserIcon, LogOut, ShieldCheck } from "@/lib/icons";

/**
 * Avatar button + dropdown. Click outside to close. Click items to
 * navigate (Link) or perform an action (Logout). The dropdown links
 * default to routes Grit ships — change the hrefs to fit your app.
 */
export function UserMenu() {
  const { data: user } = useMe();
  const { mutate: logout } = useLogout();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, [open]);

  if (!user) return null;

  const initials = ((user.first_name?.[0] || "") + (user.last_name?.[0] || "")).toUpperCase() || "U";
  const fullName = [user.first_name, user.last_name].filter(Boolean).join(" ") || "User";

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-label="Open user menu"
        className="h-9 w-9 overflow-hidden rounded-full ring-2 ring-accent/40 hover:ring-accent transition-colors bg-bg-elevated"
      >
        {user.avatar ? (
          <img src={user.avatar} alt={fullName} width={36} height={36} decoding="async" className="h-full w-full object-cover" />
        ) : (
          <span className="flex h-full w-full items-center justify-center text-sm font-semibold text-foreground">
            {initials}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-50 mt-2 w-64 overflow-hidden rounded-xl border border-border bg-bg-elevated shadow-xl">
          <div className="border-b border-border px-4 py-3">
            <p className="text-sm font-semibold text-foreground truncate">{fullName}</p>
            <p className="text-xs text-text-muted truncate">{user.email}</p>
          </div>
          <nav className="py-1 text-sm">
            <Link
              href="/system/activity"
              onClick={() => setOpen(false)}
              className="flex items-center gap-3 px-4 py-2.5 text-text-secondary hover:bg-bg-hover hover:text-foreground"
            >
              <Activity className="h-4 w-4" />
              User Activity
            </Link>
            <Link
              href="/profile"
              onClick={() => setOpen(false)}
              className="flex items-center gap-3 px-4 py-2.5 text-text-secondary hover:bg-bg-hover hover:text-foreground"
            >
              <UserIcon className="h-4 w-4" />
              Profile
            </Link>
            <Link
              href="/system/account#security"
              onClick={() => setOpen(false)}
              className="flex items-center gap-3 px-4 py-2.5 text-text-secondary hover:bg-bg-hover hover:text-foreground"
            >
              <ShieldCheck className="h-4 w-4" />
              Security
            </Link>
            <button
              type="button"
              onClick={() => { setOpen(false); logout(); }}
              className="flex w-full items-center gap-3 px-4 py-2.5 text-left text-text-secondary hover:bg-bg-hover hover:text-foreground"
            >
              <LogOut className="h-4 w-4" />
              Log out
            </button>
          </nav>
        </div>
      )}
    </div>
  );
}
