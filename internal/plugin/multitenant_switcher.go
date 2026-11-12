package plugin

// The organization switcher (grit#141).
//
// The plugin resolved the active organization from X-Organization-ID and
// nothing in a generated client ever set it. A user who belongs to one
// organization is fine, because the middleware activates it automatically; a
// user who belongs to two had no way to choose, which is the case the plugin
// exists for.
//
// The header itself is the framework's: lib/active-org.ts and the api-core
// interceptor ship in every project and are inert until something sets a
// value. That is deliberate. api-client.ts is rewritten by every upgrade, so
// injecting the header into it would make the file read as the developer's own
// edit: the injection would survive and the file would stop receiving the
// fixes that matter in it, which are CSRF, token refresh and the 401 retry.
//
// So what is left here is the part that is genuinely about organizations: the
// list, the choice, and clearing the cache when the choice changes.

// mtOrganizationsHook emits hooks/use-organizations.ts.
func mtOrganizationsHook() string {
	return `"use client";

import { useCallback, useEffect, useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";

import { apiClient } from "@/lib/api-client";
import { getActiveOrg, setActiveOrg, subscribeActiveOrg } from "@/lib/active-org";

export interface Organization {
  id: string;
  name: string;
  slug: string;
  description?: string;
  owner_id: string;
  active: boolean;
  owner_label?: string;
  admin_label?: string;
  member_label?: string;
  logo_url?: string;
  website?: string;
  legal_name?: string;
  tax_id?: string;
  tax_country?: string;
  address_line1?: string;
  address_line2?: string;
  city?: string;
  region?: string;
  postal_code?: string;
  country?: string;
}

/**
 * The organizations this user belongs to.
 *
 * GET /organizations reads memberships and looks the organizations up by id,
 * so it can only ever return too few, never another tenant's.
 */
export function useMyOrganizations() {
  return useQuery<Organization[]>({
    queryKey: ["organizations", "mine"],
    queryFn: async () => {
      // "/api/..." and not "/organizations": the client's versionedPath
      // rewrite only touches a path that starts with /api/, so a bare path
      // goes to the host root and 404s. It puts the version in from there.
      const { data } = await apiClient.get("/api/organizations");
      return data.data ?? [];
    },
    // The list changes when somebody is added to an organization, which is
    // rare and never while they are looking at this dropdown.
    staleTime: 5 * 60 * 1000,
  });
}

/**
 * The active organization, and how to change it.
 *
 * Changing it clears the query cache, which is not a nicety: every list on
 * screen was fetched for the previous organization, and React Query would go
 * on showing those rows under the new organization's name until each one
 * refetched. Showing another tenant's data under this tenant's name is worse
 * than showing a spinner.
 */
export function useActiveOrg() {
  const queryClient = useQueryClient();
  const [activeID, setActiveID] = useState<string | null>(() => getActiveOrg());

  // Another tab, or another component, may have changed it.
  useEffect(() => subscribeActiveOrg(setActiveID), []);

  const switchTo = useCallback(
    (orgID: string | null) => {
      if (orgID === getActiveOrg()) return;
      setActiveOrg(orgID);
      // Not invalidate: remove. An invalidated query keeps its old data while
      // it refetches, and that data belongs to the organization being left.
      queryClient.clear();
    },
    [queryClient],
  );

  return { activeID, switchTo };
}

/**
 * The organizations plus the active one, with a stored id that no longer
 * applies cleared.
 *
 * A membership can be revoked while the id sits in this browser's storage.
 * Sending it then means every request carries an organization the API will
 * refuse, and the reader sees an app that is broken rather than one that has
 * moved them back to where they belong.
 */
export function useOrganizationChoice() {
  const { data: organizations = [], isLoading } = useMyOrganizations();
  const { activeID, switchTo } = useActiveOrg();

  useEffect(() => {
    if (isLoading || organizations.length === 0) return;
    if (!activeID || !organizations.some((o) => o.id === activeID)) {
      // A first organization is the natural active choice. Persisting it makes
      // later requests explicit and avoids a client that displays one tenant
      // while its API requests have no tenant selected.
      if (!activeID) {
        switchTo(organizations[0].id);
        return;
      }
      switchTo(null);
    }
  }, [isLoading, organizations, activeID, switchTo]);

  const active = organizations.find((o) => o.id === activeID) ?? null;
  return { organizations, active, activeID, switchTo, isLoading };
}
`
}

// mtOrganizationSwitcher emits components/chrome/OrganizationSwitcher.tsx.
//
// It renders nothing for a reader with fewer than two organizations, which is
// most readers and every single-tenant project. A control that offers one
// choice is furniture.
func mtOrganizationSwitcher() string {
	return `"use client";

import { useState } from "react";
import { Building2, Check, ChevronDown } from "lucide-react";

import { useOrganizationChoice } from "@/hooks/use-organizations";

/**
 * Which organization the reader is working in.
 *
 * Hidden when there is no choice to make: a reader in one organization has it
 * activated by the API automatically, and a control with a single option is
 * furniture that implies a decision nobody has.
 */
export function OrganizationSwitcher() {
  const { organizations, active, switchTo, isLoading } = useOrganizationChoice();
  const [open, setOpen] = useState(false);

  // Nothing at all when there is no choice to make: a reader in one
  // organization has it activated by the API automatically, and a control with
  // a single option is furniture that implies a decision nobody has. The strip
  // is inside this component for the same reason, so there is no empty bar
  // taking up space above every page in a single-organization project.
  if (isLoading) return null;
  if (organizations.length < 2) {
    const only = organizations[0];
    return <div className="flex items-center gap-2 px-4 pt-4 text-sm md:px-8">
      {only?.logo_url ? <img src={only.logo_url} alt="" className="h-5 w-5 shrink-0 rounded object-contain" /> : null}
      <span className="truncate text-text-secondary">{only?.name ?? "No organization yet"}</span>
    </div>;
  }

  // No stored choice yet: the API activates the first membership, so that is
  // what the reader is actually looking at and what the label must say.
  const shown = active ?? organizations[0];

  return (
    <div className="relative flex items-center gap-3 px-4 pt-4 md:px-8">
      <button
        type="button"
        onClick={() => setOpen((v) => !v)}
        aria-haspopup="listbox"
        aria-expanded={open}
        className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-elevated px-3 py-1.5 text-sm text-text-primary transition-colors hover:bg-bg-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
      >
        <Building2 className="h-4 w-4 text-text-secondary" aria-hidden="true" />
        {shown.logo_url ? (
          <img src={shown.logo_url} alt="" className="h-5 w-5 shrink-0 rounded object-contain" />
        ) : null}
        <span className="max-w-[12rem] truncate">{shown.name}</span>
        <ChevronDown className="h-4 w-4 text-text-secondary" aria-hidden="true" />
      </button>

      {open && (
        <>
          {/* Click-away. Before the dropdown in the DOM so it cannot cover
              the options. */}
          <div
            className="fixed inset-0 z-40"
            aria-hidden="true"
            onClick={() => setOpen(false)}
          />
          <ul
            role="listbox"
            aria-label="Organization"
            className="absolute left-4 top-full z-50 mt-1 max-h-72 w-64 overflow-auto rounded-lg border border-border bg-bg-elevated py-1 shadow-lg md:left-8"
          >
            {organizations.map((org) => {
              const current = org.id === shown.id;
              return (
                <li key={org.id}>
                  <button
                    type="button"
                    role="option"
                    aria-selected={current}
                    onClick={() => {
                      switchTo(org.id);
                      setOpen(false);
                    }}
                    className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-text-primary transition-colors hover:bg-bg-hover focus-visible:outline-none focus-visible:bg-bg-hover"
                  >
                    <Check
                      className={
                        "h-4 w-4 shrink-0 " + (current ? "text-accent" : "opacity-0")
                      }
                      aria-hidden="true"
                    />
                    {org.logo_url ? (
                      <img src={org.logo_url} alt="" className="h-5 w-5 shrink-0 rounded object-contain" />
                    ) : null}
                    <span className="truncate">{org.name}</span>
                  </button>
                </li>
              );
            })}
          </ul>
        </>
      )}
    </div>
  );
}
`
}

func mtWebOrganizationPage() string {
	return `"use client";

import { FormEvent, useEffect, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { adminLoginUrl, orgClient } from "@/lib/organization-client";
import { useOrganizationChoice } from "@/hooks/use-organizations";

type Member = { user_id: string; email: string; first_name: string; last_name: string; role: string; created_at: string };
type Invitation = { id: string; email: string; role: string; expires_at: string; created_at: string };
type ArchivedOrganization = { id: string; name: string };

type Profile = {
  name: string; slug: string; description: string; website: string; legal_name: string;
  tax_id: string; tax_country: string; address_line1: string; address_line2: string;
  city: string; region: string; postal_code: string; country: string;
};

const EMPTY_PROFILE: Profile = {
  name: "", slug: "", description: "", website: "", legal_name: "", tax_id: "",
  tax_country: "", address_line1: "", address_line2: "", city: "", region: "",
  postal_code: "", country: "",
};

export default function OrganizationSettingsPage() {
  const { active, organizations, isLoading, switchTo } = useOrganizationChoice();
  const queryClient = useQueryClient();
  const [members, setMembers] = useState<Member[]>([]);
  const [invitations, setInvitations] = useState<Invitation[]>([]);
  const [archived, setArchived] = useState<ArchivedOrganization[]>([]);
  const [email, setEmail] = useState("");
  const [organizationName, setOrganizationName] = useState("");
  const [role, setRole] = useState("member");
  const [labels, setLabels] = useState({ owner_label: "Owner", admin_label: "Admin", member_label: "Member" });
  const [profile, setProfile] = useState<Profile>(EMPTY_PROFILE);
  const [logoURL, setLogoURL] = useState("");
  const [logoBusy, setLogoBusy] = useState(false);
  const [slugState, setSlugState] = useState({ checked: false, available: false, message: "" });
  const [notice, setNotice] = useState("");
  const [busy, setBusy] = useState(false);
  const [me, setMe] = useState<{ id: string; email: string } | null>(null);
  const [canManageInvitations, setCanManageInvitations] = useState(false);
  // The web app has no sign-in page of its own: the session is the one the
  // admin panel sets on the same origin. Rendering the creation form to a
  // signed-out visitor only produced a 401 from the API, so the page asks them
  // to sign in first and points at the login that exists.
  const [signedIn, setSignedIn] = useState<boolean | null>(null);
  const isOwner = Boolean(active && me && active.owner_id === me.id);

  // Who is looking at the page decides which controls exist. The API enforces
  // every one of them; this only keeps a member from being shown buttons that
  // answer 403, which is what made the page say "could not load" for a member.
  // The session cookies are HttpOnly, so the browser cannot answer "is somebody
  // signed in" from document.cookie. Ask the API once; the same answer names the
  // account that owns the controls.
  useEffect(() => {
    orgClient.get("/api/auth/me")
      .then(({ data }) => { setMe(data.data ?? null); setSignedIn(Boolean(data?.data)); })
      .catch(() => setSignedIn(false));
  }, []);

  async function refresh() {
    if (!active) return;
    const [memberResponse, invitationResponse] = await Promise.all([
      orgClient.get("/api/organizations/" + active.id + "/members"),
      orgClient.get("/api/organizations/" + active.id + "/invitations").catch((error: any) => {
        // A member may read the membership but not the invitations. That is a
        // permission, not a failure to load the page.
        if (error?.response?.status === 403) return null;
        throw error;
      }),
    ]);
    setMembers(memberResponse.data.data ?? []);
    setInvitations(invitationResponse?.data?.data ?? []);
    setCanManageInvitations(Boolean(invitationResponse));
    setLabels({ owner_label: active.owner_label || "Owner", admin_label: active.admin_label || "Admin", member_label: active.member_label || "Member" });
  }

  useEffect(() => { refresh().catch(() => setNotice("Could not load organization settings.")); }, [active?.id, me?.id]);
  useEffect(() => {
    if (!active) { setProfile(EMPTY_PROFILE); setLogoURL(""); return; }
    setProfile({
      name: active.name || "", slug: active.slug || "", description: active.description || "",
      website: active.website || "", legal_name: active.legal_name || "", tax_id: active.tax_id || "",
      tax_country: active.tax_country || "", address_line1: active.address_line1 || "",
      address_line2: active.address_line2 || "", city: active.city || "",
      region: active.region || "", postal_code: active.postal_code || "", country: active.country || "",
    });
    setLogoURL(active.logo_url || "");
    setSlugState({ checked: false, available: false, message: "" });
  }, [active?.id]);
  useEffect(() => {
    if (active) return;
    orgClient.get("/api/organizations/archived").then(({ data }) => setArchived(data.data ?? [])).catch(() => {});
  }, [active?.id]);

  async function invite(event: FormEvent) {
    event.preventDefault();
    if (!active) return;
    setBusy(true); setNotice("");
    try {
      await orgClient.post("/api/organizations/" + active.id + "/invitations", { email, role });
      setEmail(""); setNotice("Invitation sent."); await refresh();
    } catch (error: any) {
      setNotice(error?.response?.data?.error?.message || "Could not send invitation. Check the email or pending invitations.");
    }
    finally { setBusy(false); }
  }

  async function createOrganization(event: FormEvent) {
    event.preventDefault();
    setBusy(true); setNotice("");
    try {
      const { data } = await orgClient.post("/api/organizations", { name: organizationName });
      const organization = data.data;
      switchTo(organization.id);
      setOrganizationName(""); setNotice("Organization created.");
    } catch { setNotice("Could not create organization."); }
    finally { setBusy(false); }
  }

  async function saveProfile(event: FormEvent) {
    event.preventDefault();
    if (!active) return;
    setBusy(true); setNotice("");
    try {
      const { data } = await orgClient.patch("/api/organizations/" + active.id, profile);
      setProfile({
        name: data.data?.name || "", slug: data.data?.slug || "", description: data.data?.description || "",
        website: data.data?.website || "", legal_name: data.data?.legal_name || "", tax_id: data.data?.tax_id || "",
        tax_country: data.data?.tax_country || "", address_line1: data.data?.address_line1 || "",
        address_line2: data.data?.address_line2 || "", city: data.data?.city || "",
        region: data.data?.region || "", postal_code: data.data?.postal_code || "", country: data.data?.country || "",
      });
      setLogoURL(data.data?.logo_url || logoURL);
      await queryClient.invalidateQueries({ queryKey: ["organizations", "mine"] });
      setNotice("Organization details saved.");
    } catch (error: any) {
      const code = error?.response?.data?.error?.code;
      setNotice(
        code === "SLUG_TAKEN" ? "That slug is already in use."
          : code === "TAX_ID_TAKEN" ? "That tax id already belongs to another organization."
          : error?.response?.data?.error?.message || "Could not save organization details."
      );
    }
    finally { setBusy(false); }
  }

  async function checkSlug() {
    if (!active) return;
    const slug = profile.slug.trim().toLowerCase();
    try {
      const { data } = await orgClient.get("/api/organizations/slug-availability", { params: { slug, exclude: active.id } });
      const available = Boolean(data.data?.available);
      setSlugState({ checked: true, available, message: data.data?.reason || (available ? "Available" : "Already in use") });
    } catch {
      setSlugState({ checked: true, available: false, message: "Could not check this slug" });
    }
  }

  async function uploadLogo(file: File | undefined) {
    if (!active || !file) return;
    setLogoBusy(true); setNotice("");
    const body = new FormData();
    body.append("file", file);
    try {
      const { data } = await orgClient.post("/api/organizations/" + active.id + "/logo", body);
      setLogoURL(data.data?.logo_url || "");
      await queryClient.invalidateQueries({ queryKey: ["organizations", "mine"] });
      setNotice("Logo updated.");
    } catch (error: any) {
      setNotice(error?.response?.data?.error?.message || "Could not upload the logo.");
    } finally { setLogoBusy(false); }
  }

  async function removeLogo() {
    if (!active) return;
    try {
      await orgClient.delete("/api/organizations/" + active.id + "/logo");
      setLogoURL("");
      await queryClient.invalidateQueries({ queryKey: ["organizations", "mine"] });
      setNotice("Logo removed.");
    } catch { setNotice("Could not remove the logo."); }
  }

  async function saveLabels(event: FormEvent) {
    event.preventDefault();
    if (!active) return;
    setBusy(true); setNotice("");
    try {
      await orgClient.patch("/api/organizations/" + active.id + "/role-labels", labels);
      await queryClient.invalidateQueries({ queryKey: ["organizations", "mine"] });
      setNotice("Role labels saved.");
    } catch { setNotice("Could not save role labels."); }
    finally { setBusy(false); }
  }

  async function revoke(invitationID: string) {
    if (!active) return;
    try { await orgClient.delete("/api/organizations/" + active.id + "/invitations/" + invitationID); await refresh(); }
    catch { setNotice("Could not revoke invitation."); }
  }

  async function resend(invitationID: string) {
    if (!active) return;
    try { await orgClient.post("/api/organizations/" + active.id + "/invitations/" + invitationID + "/resend"); setNotice("Invitation resent."); await refresh(); }
    catch { setNotice("Could not resend invitation."); }
  }

  async function setMemberRole(member: Member, nextRole: string) {
    if (!active) return;
    try { await orgClient.post("/api/organizations/" + active.id + "/members", { user_id: member.user_id, role: nextRole }); setNotice("Member role updated."); await refresh(); }
    catch { setNotice("Could not update member role."); }
  }

  async function removeMember(member: Member) {
    if (!active || !window.confirm("Remove " + member.email + " from this organization?")) return;
    try { await orgClient.delete("/api/organizations/" + active.id + "/members/" + member.user_id); await refresh(); }
    catch { setNotice("Could not remove member."); }
  }

  async function transferOwner(member: Member) {
    if (!active || !window.confirm("Transfer organization ownership to " + member.email + "?")) return;
    try { await orgClient.post("/api/organizations/" + active.id + "/transfer-owner", { user_id: member.user_id }); setNotice("Ownership transferred."); await refresh(); }
    catch { setNotice("Could not transfer ownership."); }
  }

  async function archiveOrganization() {
    if (!active || !window.confirm("Archive this organization? Its data will be retained and access suspended.")) return;
    try { await orgClient.post("/api/organizations/" + active.id + "/archive", {}); switchTo(null); setNotice("Organization archived."); }
    catch { setNotice("Could not archive organization."); }
  }

  async function restoreOrganization(org: ArchivedOrganization) {
    try { await orgClient.post("/api/organizations/" + org.id + "/restore", {}); setNotice("Organization restored."); window.location.reload(); }
    catch { setNotice("Could not restore organization."); }
  }

  if (signedIn === null) return <main className="mx-auto max-w-4xl p-8">Checking your session…</main>;
  if (!signedIn) return <main className="mx-auto max-w-xl space-y-4 p-6 md:p-12"><p className="text-sm text-text-secondary">Organization settings</p><h1 className="text-3xl font-semibold">Sign in to manage your organization</h1><p className="text-text-secondary">The account area uses the session from the admin panel. Sign in there, then come back to this page.</p><div className="flex flex-wrap gap-3"><a href={adminLoginUrl} className="rounded-lg bg-accent px-4 py-2 text-accent-fg">Sign in</a><a href="/invitations/accept" className="rounded-lg border border-border px-4 py-2">Have an invitation? Open it</a></div></main>;
  if (isLoading) return <main className="mx-auto max-w-4xl p-8">Loading organizations…</main>;
  if (!active) return <main className="mx-auto max-w-4xl space-y-5 p-8"><h1 className="text-2xl font-semibold">Set up your organization</h1><p className="text-text-secondary">Create an organization or accept an invitation to continue.</p>{notice && <p role="status">{notice}</p>}<form onSubmit={createOrganization} className="flex gap-3"><label className="sr-only" htmlFor="organization-name">Organization name</label><input id="organization-name" required maxLength={160} value={organizationName} onChange={e => setOrganizationName(e.target.value)} placeholder="Organization name" className="min-w-0 flex-1 rounded-lg border border-border bg-bg-tertiary px-3 py-2" /><button disabled={busy} className="rounded-lg bg-accent px-4 py-2 text-accent-fg">Create organization</button></form><a className="underline" href="/invitations/accept">Have an invitation? Accept it</a>{archived.length > 0 && <section className="rounded-xl border border-border p-5"><h2 className="text-lg font-semibold">Archived organizations</h2>{archived.map(org => <div key={org.id} className="flex items-center justify-between py-2"><span>{org.name}</span><button type="button" onClick={() => restoreOrganization(org)} className="rounded border border-border px-3 py-1">Restore</button></div>)}</section>}</main>;

  return <main className="mx-auto max-w-5xl space-y-8 p-6 md:p-10">
    <header className="space-y-3">
      <a href="/" className="text-sm text-text-secondary underline-offset-4 hover:underline">Back to the site</a>
      <h1 className="text-3xl font-semibold">Organization settings</h1>
      {organizations.length > 1 && <label className="flex max-w-sm flex-col gap-1 text-sm" htmlFor="active-organization">Active organization<select id="active-organization" value={active.id} onChange={(event) => switchTo(event.target.value)} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2">{organizations.map((org) => <option key={org.id} value={org.id}>{org.name}</option>)}</select></label>}
    </header>
    {notice && <p role="status" className="rounded-lg border border-border p-3 text-sm">{notice}</p>}

    <section className="rounded-xl border border-border p-5">
      <h2 className="text-xl font-semibold">Organization details</h2>
      <p className="mb-4 mt-1 text-sm text-text-secondary">Name and slug identify the organization. The rest is optional and used on documents, invoices and support pages.</p>
      {!isOwner && <p className="mb-4 rounded-lg border border-border p-3 text-sm text-text-secondary">Only the organization owner can change these details.</p>}
      <div className="mb-5 flex flex-wrap items-center gap-4">
        {logoURL ? <img src={logoURL} alt="Organization logo" className="h-16 w-16 rounded-lg border border-border object-contain" /> : <div className="flex h-16 w-16 items-center justify-center rounded-lg border border-dashed border-border text-xs text-text-secondary">No logo</div>}
        {isOwner && <div className="flex flex-wrap items-center gap-2">
          <label htmlFor="organization-logo" className="cursor-pointer rounded-lg border border-border px-3 py-2 text-sm">{logoBusy ? "Uploading…" : "Upload logo"}</label>
          <input id="organization-logo" type="file" accept="image/png,image/jpeg,image/webp,image/gif" className="sr-only" disabled={logoBusy} onChange={event => { uploadLogo(event.target.files?.[0]); event.target.value = ""; }} />
          {logoURL && <button type="button" onClick={removeLogo} className="rounded-lg border border-border px-3 py-2 text-sm">Remove logo</button>}
        </div>}
        <p className="text-xs text-text-secondary">PNG, JPEG, WebP or GIF, up to 2 MB.</p>
      </div>
      <form onSubmit={saveProfile} className="grid gap-4 md:grid-cols-2">
        <label className="grid gap-1 text-sm md:col-span-2" htmlFor="organization-display-name">Name<input id="organization-display-name" required disabled={!isOwner} maxLength={160} value={profile.name} onChange={e => setProfile({...profile, name: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="organization-slug">Slug<input id="organization-slug" required maxLength={160} value={profile.slug} onChange={e => { setProfile({...profile, slug: e.target.value}); setSlugState({ checked: false, available: false, message: "" }); }} onBlur={checkSlug} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" />{slugState.checked && <span className={slugState.available ? "text-xs text-green-600" : "text-xs text-red-600"}>{slugState.message}</span>}</label>
        <label className="grid gap-1 text-sm" htmlFor="organization-website">Website<input id="organization-website" type="url" placeholder="https://example.com" value={profile.website} onChange={e => setProfile({...profile, website: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm md:col-span-2" htmlFor="organization-description">Description<textarea id="organization-description" rows={3} maxLength={500} value={profile.description} onChange={e => setProfile({...profile, description: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>

        <label className="grid gap-1 text-sm" htmlFor="organization-legal-name">Legal name<input id="organization-legal-name" maxLength={200} value={profile.legal_name} onChange={e => setProfile({...profile, legal_name: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="organization-tax-id">Tax ID<input id="organization-tax-id" maxLength={64} value={profile.tax_id} onChange={e => setProfile({...profile, tax_id: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="organization-tax-country">Tax country<input id="organization-tax-country" maxLength={2} placeholder="ES" value={profile.tax_country} onChange={e => setProfile({...profile, tax_country: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="organization-country">Country<input id="organization-country" maxLength={2} placeholder="ES" value={profile.country} onChange={e => setProfile({...profile, country: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>

        <label className="grid gap-1 text-sm md:col-span-2" htmlFor="organization-address-1">Address line 1<input id="organization-address-1" maxLength={200} value={profile.address_line1} onChange={e => setProfile({...profile, address_line1: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm md:col-span-2" htmlFor="organization-address-2">Address line 2<input id="organization-address-2" maxLength={200} value={profile.address_line2} onChange={e => setProfile({...profile, address_line2: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="organization-city">City<input id="organization-city" maxLength={120} value={profile.city} onChange={e => setProfile({...profile, city: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="organization-region">Region<input id="organization-region" maxLength={120} value={profile.region} onChange={e => setProfile({...profile, region: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="organization-postal-code">Postal code<input id="organization-postal-code" maxLength={32} value={profile.postal_code} onChange={e => setProfile({...profile, postal_code: e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>

        {isOwner && <button disabled={busy} className="w-fit rounded-lg bg-accent px-4 py-2 text-accent-fg disabled:opacity-50 md:col-span-2">{busy ? "Saving…" : "Save details"}</button>}
      </form>
    </section>

    {canManageInvitations && <section className="rounded-xl border border-border p-5">
      <h2 className="text-xl font-semibold">Invite a teammate</h2>
      <p className="mb-4 mt-1 text-sm text-text-secondary">The invitation expires in seven days. The recipient must verify this email before joining.</p>
      <form onSubmit={invite} className="flex flex-col gap-3 sm:flex-row">
        <label className="sr-only" htmlFor="invite-email">Email address</label>
        <input id="invite-email" type="email" required value={email} onChange={e => setEmail(e.target.value)} placeholder="name@example.com" className="min-w-0 flex-1 rounded-lg border border-border bg-bg-tertiary px-3 py-2" />
        <label className="sr-only" htmlFor="invite-role">Role</label>
        <select id="invite-role" value={role} onChange={e => setRole(e.target.value)} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2"><option value="member">{active.member_label || "Member"}</option><option value="admin">{active.admin_label || "Admin"}</option></select>
        <button disabled={busy} className="rounded-lg bg-accent px-4 py-2 text-accent-fg disabled:opacity-50">Send invite</button>
      </form>
    </section>}

    <section className="rounded-xl border border-border p-5">
      <h2 className="mb-4 text-xl font-semibold">Members</h2>
      <ul className="divide-y divide-border">{members.map(member => <li key={member.user_id} className="flex flex-wrap items-center justify-between gap-3 py-3"><span><strong>{member.first_name} {member.last_name}</strong><span className="ml-2 text-sm text-text-secondary">{member.email}</span></span><span className="flex flex-wrap items-center gap-2"><span className="rounded-full border border-border px-3 py-1 text-sm">{member.role === "owner" ? labels.owner_label : member.role === "admin" ? labels.admin_label : labels.member_label}</span>{isOwner && member.role !== "owner" && <><select aria-label={"Role for " + member.email} value={member.role === "admin" ? "admin" : "member"} onChange={e => setMemberRole(member, e.target.value)} className="rounded border border-border bg-bg-tertiary px-2 py-1 text-sm"><option value="member">{labels.member_label}</option><option value="admin">{labels.admin_label}</option></select><button type="button" onClick={() => transferOwner(member)} className="rounded border border-border px-2 py-1 text-sm">Make owner</button><button type="button" onClick={() => removeMember(member)} className="rounded border border-border px-2 py-1 text-sm">Remove</button></>}</span></li>)}</ul>
      {members.length === 0 && <p className="text-sm text-text-secondary">No members found.</p>}
    </section>

    {canManageInvitations && <section className="rounded-xl border border-border p-5">
      <h2 className="mb-4 text-xl font-semibold">Pending invitations</h2>
      <ul className="divide-y divide-border">{invitations.map(invitation => <li key={invitation.id} className="flex flex-wrap items-center justify-between gap-3 py-3"><span>{invitation.email}<span className="ml-2 text-sm text-text-secondary">{invitation.role}</span></span><span className="flex gap-2"><button type="button" onClick={() => resend(invitation.id)} className="rounded border border-border px-3 py-1 text-sm">Resend</button><button type="button" onClick={() => revoke(invitation.id)} className="rounded border border-border px-3 py-1 text-sm">Revoke</button></span></li>)}</ul>
      {invitations.length === 0 && <p className="text-sm text-text-secondary">No pending invitations.</p>}
    </section>}

    {isOwner && <section className="rounded-xl border border-border p-5">
      <h2 className="mb-2 text-xl font-semibold">Role names</h2><p className="mb-4 text-sm text-text-secondary">These labels affect only what people see. Permission keys stay unchanged.</p>
      <form onSubmit={saveLabels} className="grid gap-3 sm:grid-cols-3">
        <label className="grid gap-1 text-sm" htmlFor="label-owner">Owner<input id="label-owner" required maxLength={40} value={labels.owner_label} onChange={e => setLabels({...labels, owner_label:e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="label-admin">Admin<input id="label-admin" required maxLength={40} value={labels.admin_label} onChange={e => setLabels({...labels, admin_label:e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <label className="grid gap-1 text-sm" htmlFor="label-member">Member<input id="label-member" required maxLength={40} value={labels.member_label} onChange={e => setLabels({...labels, member_label:e.target.value})} className="rounded-lg border border-border bg-bg-tertiary px-3 py-2" /></label>
        <button disabled={busy} className="w-fit rounded-lg border border-border px-4 py-2 sm:col-span-3">Save labels</button>
      </form>
    </section>}

    {isOwner && <section className="rounded-xl border border-red-500/30 p-5"><h2 className="text-xl font-semibold">Archive organization</h2><p className="my-2 text-sm text-text-secondary">Archiving suspends access but keeps organization and tenant-owned data.</p><button type="button" onClick={archiveOrganization} className="rounded-lg border border-red-500 px-4 py-2 text-red-600">Archive organization</button></section>}
  </main>;
}
`
}

func mtWebInvitationPage() string {
	return `"use client";

import { FormEvent, useEffect, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { orgClient } from "@/lib/organization-client";
import { setActiveOrg } from "@/lib/active-org";

type Preview = { organization_name: string; email: string; role: string; status: string };
type Me = { id: string; email: string; email_verified_at?: string | null };
type Mode =
  | "loading" | "invalid" | "expired" | "accepted" | "signin" | "register"
  | "verify" | "ready" | "mismatch" | "totp" | "error";

export default function AcceptOrganizationInvitationPage() {
  const queryClient = useQueryClient();
  const [token, setToken] = useState("");
  const [preview, setPreview] = useState<Preview | null>(null);
  const [me, setMe] = useState<Me | null>(null);
  const [mode, setMode] = useState<Mode>("loading");
  const [message, setMessage] = useState("");
  const [notice, setNotice] = useState("");
  const [password, setPassword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [busy, setBusy] = useState(false);
  // The second step of a two-factor sign-in. /api/auth/login hands back a
  // pending token rather than a session, and this page used to throw it away
  // and tell the person to finish in an app they have no reason to have.
  const [pendingToken, setPendingToken] = useState("");
  const [code, setCode] = useState("");
  const [useBackupCode, setUseBackupCode] = useState(false);

  useEffect(() => {
    const value = new URLSearchParams(window.location.search).get("token") || "";
    setToken(value);
    if (!value) {
      setMode("invalid");
      setMessage("This invitation link is missing its token.");
      return;
    }
    void load(value);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function describe(error: any, fallback: string) {
    return error?.response?.data?.error?.message || fallback;
  }

  async function currentUser(): Promise<Me | null> {
    try {
      const { data } = await orgClient.get("/api/auth/me");
      return data.data as Me;
    } catch {
      return null;
    }
  }

  function decide(found: Preview, user: Me | null) {
    if (found.status === "accepted") {
      setMode("accepted");
      setMessage("This invitation has already been accepted.");
      return;
    }
    if (found.status === "revoked") {
      setMode("invalid");
      setMessage("This invitation was revoked.");
      return;
    }
    if (found.status === "expired") {
      setMode("expired");
      setMessage("This invitation has expired. Ask the organization for a new one.");
      return;
    }
    if (!user) {
      setMode("signin");
      setMessage("");
      return;
    }
    if (user.email.toLowerCase() !== found.email.toLowerCase()) {
      setMode("mismatch");
      setMessage("You are signed in as " + user.email + ", but this invitation is for " + found.email + ".");
      return;
    }
    if (!user.email_verified_at) {
      setMode("verify");
      setMessage("Verify " + found.email + " before joining.");
      return;
    }
    setMode("ready");
    setMessage("");
  }

  async function load(value: string) {
    setMode("loading");
    setMessage("");
    setNotice("");
    try {
      // A public lookup, and also the safe request that issues the CSRF cookie
      // this page's POSTs need. Without it, a fresh session's Accept answered
      // 403 CSRF_INVALID.
      const { data } = await orgClient.get("/api/organization-invitations/preview", { params: { token: value } });
      const found = data.data as Preview;
      setPreview(found);
      // Session cookies are HttpOnly, so the only way to know is to ask.
      const user = await currentUser();
      setMe(user);
      decide(found, user);
    } catch (error: any) {
      setMode("invalid");
      setMessage(describe(error, "This invitation is invalid."));
    }
  }

  async function signIn(event: FormEvent) {
    event.preventDefault();
    if (!preview) return;
    setBusy(true);
    setMessage("");
    try {
      const { data } = await orgClient.post("/api/auth/login", { email: preview.email, password });
      if (data.data?.totp_required) {
        setPendingToken(String(data.data.pending_token ?? ""));
        setCode("");
        setUseBackupCode(false);
        setMode("totp");
        setMessage("");
        return;
      }
      const user = await currentUser();
      setMe(user);
      setPassword("");
      decide(preview, user);
    } catch (error: any) {
      setMessage(describe(error, "Could not sign in with that password."));
    } finally {
      setBusy(false);
    }
  }

  // Finish the two-factor sign-in, which is the same call and the same
  // contract the admin's own TOTP screen makes: a pending token, a code, and
  // whether to trust this device. trust_device is false here on purpose: a
  // one-off visit to accept an invitation is not the moment to mark a browser
  // as trusted for every later sign-in.
  async function verifyTotp(event: FormEvent) {
    event.preventDefault();
    if (!preview || !pendingToken) return;
    setBusy(true);
    setMessage("");
    try {
      const path = useBackupCode
        ? "/api/auth/totp/backup-codes/verify"
        : "/api/auth/totp/verify";
      await orgClient.post(path, {
        pending_token: pendingToken,
        code: code.trim(),
        trust_device: false,
      });
      // From here it is exactly the password path: the session cookie is set,
      // so read who we are and decide what this invitation means for them.
      const user = await currentUser();
      setMe(user);
      setPassword("");
      setCode("");
      setPendingToken("");
      decide(preview, user);
    } catch (error: any) {
      setMessage(describe(error, useBackupCode
        ? "That backup code was not accepted."
        : "That code was not accepted. Codes last 30 seconds, so try the current one."));
    } finally {
      setBusy(false);
    }
  }

  async function registerAccount(event: FormEvent) {
    event.preventDefault();
    if (!preview) return;
    setBusy(true);
    setMessage("");
    try {
      await orgClient.post("/api/auth/register", {
        first_name: firstName,
        last_name: lastName,
        email: preview.email,
        password,
      });
      const user = await currentUser();
      setMe(user);
      setPassword("");
      setMode("verify");
      setMessage("Account created. Check " + preview.email + " for the verification link, then come back to this page.");
    } catch (error: any) {
      setMessage(describe(error, "Could not create the account."));
    } finally {
      setBusy(false);
    }
  }

  async function accept() {
    if (!preview) return;
    setBusy(true);
    setMessage("");
    try {
      const { data } = await orgClient.post("/api/organization-invitations/accept", { token });
      const organizationID = data.data?.organization_id;
      if (organizationID) setActiveOrg(organizationID);
      queryClient.clear();
      setNotice(data.message || "Invitation accepted.");
      window.location.assign("/account/organization");
    } catch (error: any) {
      const code = error?.response?.data?.error?.code;
      if (error?.response?.status === 401) {
        setMode("signin");
        setMessage("Sign in to accept this invitation.");
      } else if (code === "EMAIL_NOT_VERIFIED") {
        setMode("verify");
        setMessage("Verify " + preview.email + " before joining.");
      } else if (code === "EMAIL_MISMATCH") {
        setMode("mismatch");
        setMessage("This invitation belongs to " + preview.email + ".");
      } else {
        setMode("error");
        setMessage(describe(error, "This invitation could not be accepted."));
      }
    } finally {
      setBusy(false);
    }
  }

  async function switchAccount() {
    try {
      await orgClient.post("/api/auth/logout");
    } catch {}
    setMe(null);
    setPassword("");
    setMessage("");
    setMode("signin");
  }

  async function resendVerification() {
    setBusy(true);
    setNotice("");
    try {
      await orgClient.post("/api/auth/verify-email/send");
      setNotice("Verification email sent. Check " + (preview?.email ?? "your inbox") + ".");
    } catch (error: any) {
      setNotice(describe(error, "Could not send the verification email."));
    } finally {
      setBusy(false);
    }
  }

  const heading = preview?.organization_name ? "Join " + preview.organization_name : "Join your team";
  const input = "rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-foreground";
  const primary = "rounded-lg bg-accent px-4 py-2 text-accent-fg disabled:opacity-50";
  const secondary = "rounded-lg border border-border px-4 py-2 text-foreground disabled:opacity-50";

  return (
    <main className="mx-auto max-w-xl p-6 md:p-12">
      <section className="space-y-5 rounded-xl border border-border p-6">
        <div>
          <p className="text-sm text-text-secondary">Organization invitation</p>
          <h1 className="mt-1 text-3xl font-semibold text-foreground">{heading}</h1>
          {preview && <p className="mt-2 text-sm text-text-secondary">{preview.email} · {preview.role}</p>}
        </div>

        {message && <p role="alert" className="rounded-lg border border-border p-3 text-sm text-text-secondary">{message}</p>}
        {notice && <p role="status" className="rounded-lg border border-border p-3 text-sm text-text-secondary">{notice}</p>}

        {mode === "loading" && <p className="text-sm text-text-secondary">Loading invitation…</p>}

        {mode === "ready" && (
          <button type="button" onClick={accept} disabled={busy} className={primary}>
            {busy ? "Joining…" : "Accept invitation"}
          </button>
        )}

        {mode === "signin" && preview && (
          <form onSubmit={signIn} className="space-y-3">
            <p className="text-sm text-text-secondary">
              Sign in with the invited address. No account yet? Create one below.
            </p>
            <label className="grid gap-1 text-sm text-foreground" htmlFor="invitation-email">Email
              <input id="invitation-email" type="email" readOnly value={preview.email} className={input} />
            </label>
            <label className="grid gap-1 text-sm text-foreground" htmlFor="invitation-password">Password
              <input id="invitation-password" type="password" required autoComplete="current-password" value={password} onChange={(event) => setPassword(event.target.value)} className={input} />
            </label>
            <div className="flex flex-wrap gap-3">
              <button type="submit" disabled={busy} className={primary}>{busy ? "Signing in…" : "Sign in and join"}</button>
              <button type="button" onClick={() => { setMode("register"); setMessage(""); }} className={secondary}>Create an account</button>
            </div>
          </form>
        )}

        {mode === "register" && preview && (
          <form onSubmit={registerAccount} className="space-y-3">
            <p className="text-sm text-text-secondary">
              The account is created for the invited address. You will verify it by email before joining.
            </p>
            <label className="grid gap-1 text-sm text-foreground" htmlFor="invitation-first-name">First name
              <input id="invitation-first-name" required value={firstName} onChange={(event) => setFirstName(event.target.value)} className={input} />
            </label>
            <label className="grid gap-1 text-sm text-foreground" htmlFor="invitation-last-name">Last name
              <input id="invitation-last-name" required value={lastName} onChange={(event) => setLastName(event.target.value)} className={input} />
            </label>
            <label className="grid gap-1 text-sm text-foreground" htmlFor="invitation-register-email">Email
              <input id="invitation-register-email" type="email" readOnly value={preview.email} className={input} />
            </label>
            <label className="grid gap-1 text-sm text-foreground" htmlFor="invitation-register-password">Password
              <input id="invitation-register-password" type="password" required autoComplete="new-password" value={password} onChange={(event) => setPassword(event.target.value)} className={input} />
            </label>
            <div className="flex flex-wrap gap-3">
              <button type="submit" disabled={busy} className={primary}>{busy ? "Creating…" : "Create account"}</button>
              <button type="button" onClick={() => { setMode("signin"); setMessage(""); }} className={secondary}>I already have an account</button>
            </div>
          </form>
        )}

        {mode === "verify" && (
          <div className="space-y-3">
            <p className="text-sm text-text-secondary">
              We sent a verification link to {preview?.email ?? "your address"}. Open it, then return here.
            </p>
            <div className="flex flex-wrap gap-3">
              <button type="button" onClick={() => void load(token)} className={primary}>I verified my email</button>
              <button type="button" onClick={resendVerification} disabled={busy} className={secondary}>{busy ? "Sending…" : "Resend the link"}</button>
            </div>
          </div>
        )}

        {mode === "mismatch" && (
          <div className="space-y-3">
            <p className="text-sm text-text-secondary">
              Sign in with {preview?.email ?? "the invited address"} to accept, or ask the organization for a new invitation.
            </p>
            <button type="button" onClick={switchAccount} className={secondary}>Sign in with another account</button>
          </div>
        )}

        {mode === "totp" && (
          <form onSubmit={verifyTotp} className="space-y-3">
            <p className="text-sm text-text-secondary">
              {useBackupCode
                ? "Enter one of the backup codes you saved when you turned on two-factor sign-in."
                : "Enter the 6-digit code from your authenticator app."}
            </p>
            <input
              value={code}
              onChange={(event) => setCode(event.target.value)}
              placeholder={useBackupCode ? "backup code" : "123456"}
              inputMode={useBackupCode ? "text" : "numeric"}
              autoComplete="one-time-code"
              autoFocus
              className={input}
            />
            <button type="submit" disabled={busy || code.trim() === ""} className={primary}>
              {busy ? "Checking…" : "Verify and accept"}
            </button>
            <button
              type="button"
              onClick={() => { setUseBackupCode(!useBackupCode); setCode(""); setMessage(""); }}
              className="text-sm text-text-secondary underline-offset-4 hover:underline"
            >
              {useBackupCode ? "Use my authenticator app instead" : "I have lost my authenticator"}
            </button>
          </form>
        )}

        {(mode === "invalid" || mode === "expired" || mode === "error") && (
          <a href="/" className="inline-block text-sm text-text-secondary underline-offset-4 hover:underline">Back to the site</a>
        )}

        {mode === "accepted" && (
          <div className="space-y-3">
            <a href="/account/organization" className="inline-block text-sm underline-offset-4 hover:underline">Go to your organizations</a>
            {me && <button type="button" onClick={switchAccount} className={secondary}>Sign in with another account</button>}
          </div>
        )}
      </section>
    </main>
  );
}
`
}

func mtExpoSwitcher() string {
	return `import { useEffect, useState } from "react";
import { Image, Pressable, ScrollView, Text } from "react-native";
import * as SecureStore from "@/lib/secure-store";
import { api } from "@/lib/api";
import { queryClient } from "@/lib/query-client";

type Organization = { id: string; name: string; logo_url?: string };

export function OrganizationSwitcher() {
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [activeID, setActiveID] = useState<string | null>(null);
  useEffect(() => {
    let mounted = true;
    Promise.all([api.get("/organizations"), SecureStore.getItemAsync("active_org_id")]).then(async ([result, selected]) => {
      if (!mounted) return;
      const available: Organization[] = result.data?.data ?? [];
      const active = selected && available.some(org => org.id === selected) ? selected : available[0]?.id ?? null;
      setOrganizations(available);
      setActiveID(active);
      if (active && active !== selected) await SecureStore.setItemAsync("active_org_id", active);
    }).catch(() => {});
    return () => { mounted = false; };
  }, []);
  if (organizations.length < 2) return null;
  return <ScrollView horizontal showsHorizontalScrollIndicator={false} className="max-h-14 border-b border-slate-200 bg-white px-3 py-2 dark:border-slate-800 dark:bg-slate-950">{organizations.map(org => <Pressable key={org.id} onPress={async () => { if (org.id === activeID) return; await SecureStore.setItemAsync("active_org_id", org.id); setActiveID(org.id); queryClient.clear(); }} className={"mr-2 flex-row items-center rounded-full px-4 py-2 " + (org.id === activeID ? "bg-blue-600" : "bg-slate-100 dark:bg-slate-800")}>{org.logo_url ? <Image source={{ uri: org.logo_url }} className="mr-1.5 h-5 w-5 rounded-full" /> : null}<Text className={org.id === activeID ? "text-white" : "text-slate-900 dark:text-white"}>{org.name}</Text></Pressable>)}</ScrollView>;
}
`
}

func mtExpoInvitationPage() string {
	return `import { useState } from "react";
import { ActivityIndicator, Pressable, Text, View } from "react-native";
import { useLocalSearchParams, useRouter } from "expo-router";
import * as SecureStore from "@/lib/secure-store";
import { api } from "@/lib/api";
import { queryClient } from "@/lib/query-client";

export default function AcceptOrganizationInvitationPage() {
  const { token } = useLocalSearchParams<{ token?: string }>();
  const router = useRouter();
  const [message, setMessage] = useState("Review the invitation and accept to join the organization.");
  const [busy, setBusy] = useState(false);
  async function accept() {
    if (!token) { setMessage("This invitation link is missing its token."); return; }
    setBusy(true);
    try {
      const result = await api.post("/organization-invitations/accept", { token });
      const organizationID = result.data?.organization_id;
      if (organizationID) await SecureStore.setItemAsync("active_org_id", organizationID);
      queryClient.clear();
      setMessage("Invitation accepted."); router.replace("/(tabs)");
    }
    catch (error: any) { setMessage(error?.message || "Sign in with the invited and verified email, then try again."); }
    finally { setBusy(false); }
  }
  return <View className="flex-1 justify-center gap-5 bg-white p-6 dark:bg-slate-950"><Text className="text-3xl font-bold text-slate-950 dark:text-white">Join your team</Text><Text className="text-base text-slate-600 dark:text-slate-300">{message}</Text><Pressable disabled={busy} onPress={accept} className="items-center rounded-xl bg-blue-600 p-4">{busy ? <ActivityIndicator color="white" /> : <Text className="font-semibold text-white">Accept invitation</Text>}</Pressable></View>;
}
`
}

func mtWebOrganizationHook() string {
	return `"use client";

import { useCallback, useEffect, useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api";
import { getActiveOrg, setActiveOrg, subscribeActiveOrg } from "@/lib/active-org";

export type Organization = { id:string; name:string; slug:string; description?:string; owner_id:string; owner_label?:string; admin_label?:string; member_label?:string; logo_url?:string; website?:string; legal_name?:string; tax_id?:string; tax_country?:string; address_line1?:string; address_line2?:string; city?:string; region?:string; postal_code?:string; country?:string };

export function useOrganizationChoice() {
  const queryClient = useQueryClient();
  const { data: organizations = [], isLoading, isError } = useQuery<Organization[]>({
    queryKey:["organizations","mine"],
    queryFn:async () => { const { data } = await apiClient.get("/api/organizations"); return data.data ?? []; },
    staleTime:5 * 60 * 1000,
  });
  const [activeID, setActiveID] = useState<string | null>(() => getActiveOrg());
  useEffect(() => subscribeActiveOrg(setActiveID), []);
  const switchTo = useCallback((id:string | null) => {
    if (id === getActiveOrg()) return;
    setActiveOrg(id);
    queryClient.clear();
  }, [queryClient]);
  useEffect(() => {
    if (isLoading || organizations.length === 0) return;
    if (!activeID) { switchTo(organizations[0].id); return; }
    if (!organizations.some(org => org.id === activeID)) switchTo(null);
  }, [activeID, isLoading, organizations, switchTo]);
  const active = organizations.find(org => org.id === activeID) ?? null;
  return { organizations, active, activeID, switchTo, isLoading, isError };
}
`
}

func mtWebOrganizationSwitcher() string {
	return `"use client";

import { useState } from "react";
import { useOrganizationChoice } from "@/hooks/use-organizations";

export function OrganizationSwitcher() {
  const { organizations, active, switchTo, isLoading, isError } = useOrganizationChoice();
  const [open, setOpen] = useState(false);
  if (isLoading || isError) return null;
  if (organizations.length === 0) return <a href="/account/organization" className="mr-3 text-xs text-text-muted">Create or join an organization</a>;
  if (organizations.length === 1) return <a href="/account/organization" className="mr-3 flex max-w-40 items-center gap-1.5 truncate text-xs text-text-muted">{organizations[0].logo_url ? <img src={organizations[0].logo_url} alt="" className="h-4 w-4 rounded object-contain" /> : null}{organizations[0].name}</a>;
  const shown = active ?? organizations[0];
  return <div className="relative mr-3">
    <button type="button" onClick={() => setOpen(value => !value)} aria-expanded={open} className="flex max-w-40 items-center gap-1.5 truncate rounded-lg border border-border bg-bg-elevated px-2 py-1.5 text-xs text-foreground">{shown.logo_url ? <img src={shown.logo_url} alt="" className="h-4 w-4 rounded object-contain" /> : null}{shown.name} <span aria-hidden="true">▾</span></button>
    {open && <div role="listbox" aria-label="Organization" className="absolute right-0 top-full z-50 mt-1 max-h-64 w-56 overflow-auto rounded-lg border border-border bg-bg-elevated p-1 shadow-lg">{organizations.map(org => <button key={org.id} type="button" role="option" aria-selected={org.id === shown.id} onClick={() => { switchTo(org.id); setOpen(false); }} className="flex w-full items-center gap-2 truncate rounded px-3 py-2 text-left text-sm text-foreground hover:bg-bg-hover">{org.logo_url ? <img src={org.logo_url} alt="" className="h-4 w-4 rounded object-contain" /> : null}{org.name}</button>)}</div>}
  </div>;
}
`
}

func mtTanStackActiveOrg() string {
	return `import { api } from "@/lib/api";

// Which organization this browser is working in, and the interceptor that
// sends it. The TanStack web app has no framework active-org module (only the
// Next app does), so the plugin supplies the one the switcher needs. The
// interceptor is registered on import and the switcher is mounted at the root,
// so every request made through the shared client carries the choice.
const KEY = "grit.active_org";

type Listener = (id: string | null) => void;
const listeners = new Set<Listener>();

export function getActiveOrg(): string | null {
  try { return localStorage.getItem(KEY); } catch { return null; }
}

export function setActiveOrg(id: string | null): void {
  try {
    if (id) localStorage.setItem(KEY, id);
    else localStorage.removeItem(KEY);
  } catch {}
  listeners.forEach((listener) => listener(id));
}

export function subscribeActiveOrg(listener: Listener): () => void {
  listeners.add(listener);
  const onStorage = (event: StorageEvent) => {
    if (event.key === KEY) listener(event.newValue);
  };
  window.addEventListener("storage", onStorage);
  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", onStorage);
  };
}

api.interceptors.request.use((config) => {
  const id = getActiveOrg();
  if (!id || !config.headers) return config;
  const headers = config.headers as { set?: (name: string, value: string) => void } & Record<string, unknown>;
  if (typeof headers.set === "function") headers.set("X-Organization-ID", id);
  else headers["X-Organization-ID"] = id;
  return config;
});
`
}

func mtTanStackWebOrganizationHook() string {
	return `import { useCallback, useEffect, useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { getActiveOrg, setActiveOrg, subscribeActiveOrg } from "@/lib/active-org";

export type Organization = { id:string; name:string; slug:string; description?:string; owner_id:string; owner_label?:string; admin_label?:string; member_label?:string; logo_url?:string; website?:string; legal_name?:string; tax_id?:string; tax_country?:string; address_line1?:string; address_line2?:string; city?:string; region?:string; postal_code?:string; country?:string };

export function useOrganizationChoice() {
  const queryClient = useQueryClient();
  const { data: organizations = [], isLoading, isError } = useQuery<Organization[]>({
    queryKey:["organizations","mine"],
    queryFn:async () => { const { data } = await api.get("/api/organizations"); return data.data ?? []; },
    staleTime:5 * 60 * 1000,
  });
  const [activeID, setActiveID] = useState<string | null>(() => getActiveOrg());
  useEffect(() => subscribeActiveOrg(setActiveID), []);
  const switchTo = useCallback((id:string | null) => {
    if (id === getActiveOrg()) return;
    setActiveOrg(id);
    queryClient.clear();
  }, [queryClient]);
  useEffect(() => {
    if (isLoading || organizations.length === 0) return;
    if (!activeID) { switchTo(organizations[0].id); return; }
    if (!organizations.some(org => org.id === activeID)) switchTo(null);
  }, [activeID, isLoading, organizations, switchTo]);
  const active = organizations.find(org => org.id === activeID) ?? null;
  return { organizations, active, activeID, switchTo, isLoading, isError };
}
`
}

func mtOrganizationClient(imported, adminLogin string) string {
	return `import { ` + imported + ` } from "@/lib/api";

// The client the organization screens call. The Next web app exports
// apiClient from lib/api; the Vite one exports api. Keeping the difference in
// one generated module lets the account page itself be identical in both.
export const orgClient = ` + imported + `;

// Where this stack's sign-in page is. The web app has none of its own; the
// admin panel sets the session cookie the account area uses.
export const adminLoginUrl = ` + adminLogin + `;
`
}
