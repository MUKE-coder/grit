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
    if (activeID && !organizations.some((o) => o.id === activeID)) {
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
  if (isLoading || organizations.length < 2) return null;

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
