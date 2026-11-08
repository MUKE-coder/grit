package scaffold

// /system is the platform's own console: jobs, backups, audit, health, roles.
// A tenant user must never reach it, whatever role they hold inside a tenant.
//
// The check lives in its own component, rendered by the dashboard layout, and
// not in AdminLayout. Plugins inject into admin-layout.tsx, which makes the
// upgrade treat the file as edited and leave it alone: a project that installed
// a plugin would never receive a guard written there. The dashboard layouts are
// framework-owned and no plugin touches them.
func adminSystemGuardComponent() string {
	return `"use client";

import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";

import { useMe } from "@/hooks/use-auth";
import { usePermissions } from "@/hooks/use-permissions";

/**
 * Keeps a tenant user out of the platform console.
 *
 * The API refuses every /system request without one of these permissions, so
 * this is not the security boundary; it is what stops the shell from rendering
 * for somebody who is about to be refused. A tenant role (owner included) is
 * not a platform role, and only staff hold these grants.
 */
export function SystemGuard() {
  const pathname = usePathname();
  const router = useRouter();
  const { data: user, isLoading: userLoading } = useMe();
  const { can, isLoading: permsLoading } = usePermissions();

  useEffect(() => {
    if (userLoading || permsLoading || !user) return;
    if (!pathname.startsWith("/system")) return;
    const staff =
      user.role === "ADMIN" ||
      ["system.view", "jobs.view", "backups.view", "audit.view", "roles.view"].some((key) => can(key));
    if (!staff) router.replace("/account");
  }, [pathname, router, user, userLoading, permsLoading, can]);

  return null;
}
`
}
