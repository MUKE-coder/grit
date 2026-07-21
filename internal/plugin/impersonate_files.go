package plugin

import "strings"

// impersonateHandlerGo emits internal/handlers/impersonate.go.
func impersonateHandlerGo(ctx Context) string {
	src := `package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/respond"
	"{{MODULE}}/internal/services"
)

// ImpersonateHandler swaps the current session to another user and back.
//
// The swap is entirely server-side: auth is in HttpOnly cookies the browser
// can't read, so Start re-issues the auth cookies for the target user and
// stashes the admin's own token in grit_impersonator (also HttpOnly). Stop
// reads that cookie and restores the admin. The admin never touches a raw token.
type ImpersonateHandler struct {
	DB   *gorm.DB
	Auth *services.AuthService
}

func NewImpersonateHandler(db *gorm.DB, auth *services.AuthService) *ImpersonateHandler {
	return &ImpersonateHandler{DB: db, Auth: auth}
}

// Start begins impersonating a user. Mounted on the ADMIN group, so only an
// admin reaches it.
func (h *ImpersonateHandler) Start(c *gin.Context) {
	adminID, _ := c.Get("user_id")
	adminIDStr, _ := adminID.(string)

	// The admin's current access token, stashed so Stop can restore the session.
	adminToken, err := c.Cookie("grit_access")
	if err != nil || adminToken == "" {
		respond.BadRequest(c, "No active session to return to")
		return
	}

	targetID := c.Param("id")
	if targetID == adminIDStr {
		respond.BadRequest(c, "You are already yourself")
		return
	}

	var target models.User
	if err := h.DB.Where("id = ?", targetID).First(&target).Error; err != nil {
		respond.NotFound(c, "User not found")
		return
	}

	pair, err := h.Auth.GenerateTokenPair(target.ID, target.Email, target.Role)
	if err != nil {
		respond.Internal(c, err)
		return
	}

	// Swap the session to the target, stash the admin's token, and set a
	// readable flag the UI can see (HttpOnly cookies are invisible to JS).
	h.Auth.SetAuthCookies(c, pair)
	setImpersonatorCookie(c, adminToken)
	setImpersonatingFlag(c, strings.TrimSpace(target.FirstName+" "+target.LastName), target.Email)

	services.LogActivity(h.DB, c, services.ActivityArgs{
		UserID:       adminIDStr,
		Action:       "user.impersonate.start",
		Severity:     "warn",
		Summary:      "Started impersonating " + target.Email,
		ResourceType: "user",
		ResourceID:   target.ID,
	})

	respond.OK(c, gin.H{"user": target}, "Impersonation started")
}

// Stop returns to the original admin. Mounted on the PROTECTED group, because
// the caller is currently the impersonated user, who may not be an admin.
func (h *ImpersonateHandler) Stop(c *gin.Context) {
	adminToken, err := c.Cookie("grit_impersonator")
	if err != nil || adminToken == "" {
		respond.BadRequest(c, "Not impersonating anyone")
		return
	}

	claims, err := h.Auth.ValidateToken(adminToken)
	if err != nil {
		clearImpersonatorCookies(c)
		respond.Unauthorized(c, "Impersonation session expired; sign in again")
		return
	}

	pair, err := h.Auth.GenerateTokenPair(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		respond.Internal(c, err)
		return
	}

	h.Auth.SetAuthCookies(c, pair)
	clearImpersonatorCookies(c)

	services.LogActivity(h.DB, c, services.ActivityArgs{
		UserID:   claims.UserID,
		Action:   "user.impersonate.stop",
		Severity: "info",
		Summary:  "Stopped impersonating",
	})

	var admin models.User
	if err := h.DB.Where("id = ?", claims.UserID).First(&admin).Error; err != nil {
		respond.OK(c, gin.H{}, "Returned to your account")
		return
	}
	respond.OK(c, gin.H{"user": admin}, "Returned to your account")
}

func impersonateHTTPS(c *gin.Context) bool {
	return c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
}

// One hour is enough to poke around and return; the admin re-authenticates
// otherwise, which is the safe failure.
func setImpersonatorCookie(c *gin.Context, adminToken string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("grit_impersonator", adminToken, 3600, "/", "", impersonateHTTPS(c), true)
}

// grit_impersonating is NOT HttpOnly — the UI reads it to show the banner. It
// carries only display text ("Name|email"), never a token.
func setImpersonatingFlag(c *gin.Context, name, email string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("grit_impersonating", url.QueryEscape(name+"|"+email), 3600, "/", "", impersonateHTTPS(c), false)
}

func clearImpersonatorCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("grit_impersonator", "", -1, "/", "", impersonateHTTPS(c), true)
	c.SetCookie("grit_impersonating", "", -1, "/", "", impersonateHTTPS(c), false)
}
`
	return strings.ReplaceAll(src, "{{MODULE}}", ctx.Module)
}

// impersonateHook emits apps/admin/hooks/use-impersonate.ts.
func impersonateHook() string {
	return `"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

// Start impersonating a user. On success the whole app reloads as that user,
// so every cached query and screen reflects the new identity.
export function useImpersonate() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (userId: string) =>
      (await apiClient.post("/api/admin/impersonate/" + userId)).data,
    onSuccess: () => {
      qc.clear();
      window.location.href = "/dashboard";
    },
  });
}

// Return to your own account.
export function useStopImpersonate() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async () => (await apiClient.post("/api/auth/impersonate/stop")).data,
    onSuccess: () => {
      qc.clear();
      window.location.href = "/dashboard";
    },
  });
}

// Reads the non-HttpOnly grit_impersonating flag cookie. Returns who you're
// impersonating, or null. The auth token itself lives in an HttpOnly cookie
// this can't see — this is only the banner's display text.
export function readImpersonating(): { name: string; email: string } | null {
  if (typeof document === "undefined") return null;
  const m = document.cookie.match(/(?:^|; )grit_impersonating=([^;]+)/);
  if (!m) return null;
  const [name, email] = decodeURIComponent(m[1]).split("|");
  return { name: name || "", email: email || "" };
}
`
}

// impersonateBanner emits apps/admin/components/impersonation-banner.tsx.
func impersonateBanner() string {
	return `"use client";

import { useEffect, useState } from "react";
import { readImpersonating, useStopImpersonate } from "@/hooks/use-impersonate";

// A persistent amber strip shown whenever you're impersonating someone, with a
// one-click return. Rendered in the dashboard layout above the page content.
export function ImpersonationBanner() {
  const [who, setWho] = useState<{ name: string; email: string } | null>(null);
  const stop = useStopImpersonate();

  // Cookie can only be read on the client, after mount.
  useEffect(() => {
    setWho(readImpersonating());
  }, []);

  if (!who) return null;

  return (
    <div className="flex items-center justify-between gap-3 border-b border-warning/30 bg-warning/10 px-4 py-2.5 text-sm md:px-8">
      <span className="text-warning">
        You are impersonating{" "}
        <strong className="font-semibold">{who.name || who.email}</strong>
        {who.name && who.email ? " (" + who.email + ")" : ""}. Actions are performed as this user.
      </span>
      <button
        onClick={() => stop.mutate()}
        disabled={stop.isPending}
        className="shrink-0 rounded-lg border border-warning/40 bg-bg-elevated px-3 py-1.5 text-xs font-medium text-warning hover:bg-warning/10 disabled:opacity-50"
      >
        {stop.isPending ? "Returning…" : "Return to your account"}
      </button>
    </div>
  );
}
`
}

// impersonatePage emits apps/admin/app/(dashboard)/system/impersonate/page.tsx.
func impersonatePage() string {
	return `"use client";

import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { usePermissions } from "@/hooks/use-permissions";
import { useImpersonate } from "@/hooks/use-impersonate";
import { useMe } from "@/hooks/use-auth";
import { PageHeader } from "@/components/chrome/PageHeader";

interface UserRow {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  role: string;
}

export default function ImpersonatePage() {
  const { can, isLoading: permsLoading } = usePermissions();
  const { data: me } = useMe();
  const impersonate = useImpersonate();

  const { data: users, isLoading } = useQuery<UserRow[]>({
    queryKey: ["impersonate-users"],
    enabled: can("users.edit"),
    queryFn: async () => {
      const res = await apiClient.get("/api/users?page_size=200");
      return (res.data?.data ?? []) as UserRow[];
    },
  });

  if (permsLoading) return null;

  if (!can("users.edit")) {
    return (
      <div>
        <PageHeader title="Impersonate" subtitle="Sign in as another user" />
        <p className="mt-6 text-sm text-text-muted">
          You do not have permission to impersonate users.
        </p>
      </div>
    );
  }

  return (
    <div>
      <PageHeader title="Impersonate" subtitle="Sign in as another user to reproduce an issue or check their access. Every impersonation is logged." />
      <div className="mt-6 overflow-hidden rounded-xl border border-border bg-bg-secondary">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border text-left text-xs uppercase tracking-wide text-text-muted">
              <th className="px-4 py-3 font-medium">Name</th>
              <th className="px-4 py-3 font-medium">Email</th>
              <th className="px-4 py-3 font-medium">Role</th>
              <th className="px-4 py-3" />
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr><td colSpan={4} className="px-4 py-8 text-center text-text-muted">Loading…</td></tr>
            ) : (
              (users ?? []).map((u) => {
                const isSelf = me?.id === u.id;
                return (
                  <tr key={u.id} className="border-b border-border/50 last:border-b-0">
                    <td className="px-4 py-3 text-foreground">{u.first_name} {u.last_name}</td>
                    <td className="px-4 py-3 text-text-secondary">{u.email}</td>
                    <td className="px-4 py-3 text-text-secondary">{u.role}</td>
                    <td className="px-4 py-3 text-right">
                      <button
                        disabled={isSelf || impersonate.isPending}
                        onClick={() => impersonate.mutate(u.id)}
                        className="rounded-lg border border-border bg-bg-tertiary px-3 py-1.5 text-xs font-medium text-foreground hover:border-accent/40 disabled:opacity-40"
                      >
                        {isSelf ? "You" : "Impersonate"}
                      </button>
                    </td>
                  </tr>
                );
              })
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
`
}
