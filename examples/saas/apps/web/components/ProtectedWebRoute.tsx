"use client";

// v3.31.22 — apps/web/components/ProtectedWebRoute.tsx
//
// Client-side route guard. Wrap a page's contents to enforce
// authentication. Pairs with apps/web/middleware.ts — middleware
// handles the SSR cookie check; this component handles the
// post-hydration "is the cookie actually valid?" probe.
//
// Use this when you need:
//   - Role/permission checks beyond cookie existence (the cookie
//     doesn't carry role; useMe() returns the full user payload).
//   - Per-page protection without editing middleware.ts.
//
// For simple "is the visitor signed in?" pages, middleware.ts is
// faster — it bounces unauthenticated requests before the page
// ever loads.

import { useEffect, type ReactNode } from "react";
import { useRouter, usePathname } from "next/navigation";
import { useMe } from "@/hooks/use-auth";

interface ProtectedWebRouteProps {
  children: ReactNode;
  // Optional role gate. When set, the user's role must match (or
  // include, for arrays). Visitors with a session but the wrong
  // role get bounced to / (the landing page).
  roles?: string | string[];
  // Where to send unauthenticated visitors. Defaults to /login with
  // a ?next= param so they return to the original page after login.
  loginPath?: string;
  // Custom loading view while useMe() is mid-flight.
  fallback?: ReactNode;
}

export function ProtectedWebRoute({
  children,
  roles,
  loginPath = "/login",
  fallback,
}: ProtectedWebRouteProps) {
  const { data: user, isLoading, isError } = useMe();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (isLoading) return;
    if (isError || user === null) {
      const next = encodeURIComponent(pathname ?? "/");
      router.replace(loginPath + "?next=" + next);
      return;
    }
    if (roles) {
      const allowed = Array.isArray(roles) ? roles : [roles];
      // user.role is the source of truth — falls back to a friendlier
      // redirect (to the landing page) for authenticated-but-wrong-role
      // visitors, since a /login bounce wouldn't help them.
      const userWithRole = user as { role?: string };
      if (!userWithRole.role || !allowed.includes(userWithRole.role)) {
        router.replace("/");
      }
    }
  }, [isLoading, isError, user, roles, router, pathname, loginPath]);

  if (isLoading || !user) {
    return (
      fallback ?? (
        <div className="flex min-h-screen items-center justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-slate-300 border-t-slate-900" />
        </div>
      )
    );
  }

  return <>{children}</>;
}
