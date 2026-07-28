package scaffold

import "strings"

// adminSSOPage emits apps/admin/app/(dashboard)/system/sso/page.tsx — the
// screen where an operator adds a customer's identity provider.
//
// The shape follows the form-shares page: one card listing connections, a
// ResponsiveSheet for create/edit. The client secret is write-only — the API
// never returns it, so editing shows a blank field that means "keep what's
// stored" rather than pretending to show the current value.
func adminSSOPage() string {
	src := `"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ResponsiveSheet } from "@/components/ui/ResponsiveSheet";
import { apiClient } from "@/lib/api-client";
import { toast } from "sonner";
import { ShieldCheck, Plus, Trash2, Loader2, Check, AlertTriangle, Copy } from "@/lib/icons";

interface Connection {
  id: string;
  slug: string;
  name: string;
  protocol: string;
  metadata_url: string;
  domains: string;
  issuer_url: string;
  client_id: string;
  has_secret: boolean;
  enabled: boolean;
  jit_provisioning: boolean;
  groups_claim: string;
  group_mappings: string;
  last_used_at: string | null;
}

const BLANK = {
  protocol: "oidc",
  metadata_url: "",
  metadata_xml: "",
  email_attribute: "",
  groups_attribute: "",
  allow_idp_initiated: true,
  slug: "",
  name: "",
  domains: "",
  issuer_url: "",
  client_id: "",
  client_secret: "",
  groups_claim: "groups",
  group_mappings: "",
  jit_provisioning: true,
  enabled: true,
};

const inputCls =
  "w-full rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent";

function Field({ label, hint, children }: { label: string; hint?: string; children: React.ReactNode }) {
  return (
    <div className="space-y-1.5">
      <label className="block text-sm font-medium text-foreground">{label}</label>
      {children}
      {hint && <p className="text-xs text-text-muted">{hint}</p>}
    </div>
  );
}

export default function SSOPage() {
  const qc = useQueryClient();
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<Connection | null>(null);
  const [form, setForm] = useState({ ...BLANK });

  const listQ = useQuery({
    queryKey: ["sso-connections"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/sso/connections");
      return { rows: (data.data ?? []) as Connection[], live: data.meta?.live ?? 0 };
    },
  });

  const saveM = useMutation({
    mutationFn: async () => {
      if (editing) {
        const { data } = await apiClient.put("/api/sso/connections/" + editing.id, form);
        return data;
      }
      const { data } = await apiClient.post("/api/sso/connections", form);
      return data;
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["sso-connections"] });
      setOpen(false);
      setEditing(null);
      setForm({ ...BLANK });
      toast.success(editing ? "Connection updated" : "Connection created");
    },
    onError: (e: unknown) => {
      const err = e as { response?: { data?: { error?: { message?: string } } } };
      toast.error(err?.response?.data?.error?.message || "Could not save the connection");
    },
  });

  const deleteM = useMutation({
    mutationFn: async (id: string) => apiClient.delete("/api/sso/connections/" + id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["sso-connections"] });
      toast.success("Connection deleted");
    },
  });

  const rows = listQ.data?.rows ?? [];

  const startEdit = (c: Connection) => {
    setEditing(c);
    setForm({
      protocol: c.protocol || "oidc",
      metadata_url: c.metadata_url || "",
      metadata_xml: "",
      email_attribute: "",
      groups_attribute: "",
      allow_idp_initiated: true,
      slug: c.slug,
      name: c.name,
      domains: c.domains,
      issuer_url: c.issuer_url,
      client_id: c.client_id,
      client_secret: "",
      groups_claim: c.groups_claim || "groups",
      group_mappings: c.group_mappings || "",
      jit_provisioning: c.jit_provisioning,
      enabled: c.enabled,
    });
    setOpen(true);
  };

  // SAML needs two URLs on the IdP side (where to POST the assertion, and
  // where to fetch our metadata); OIDC needs one redirect URI.
  const copyCallback = (c: Connection) => {
    const api = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    if ((c.protocol || "oidc") === "saml") {
      navigator.clipboard?.writeText(
        [
          "ACS URL: " + api + "/api/auth/saml/" + c.slug + "/acs",
          "SP metadata: " + api + "/api/auth/saml/" + c.slug + "/metadata",
        ].join("\n")
      );
      toast.success("ACS + metadata URLs copied — give these to the IdP admin");
      return;
    }
    navigator.clipboard?.writeText(api + "/api/auth/sso/" + c.slug + "/callback");
    toast.success("Callback URL copied — paste it into the provider");
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Single sign-on"
        subtitle="Let a customer's team sign in with their own identity provider. One connection per organisation, routed by email domain."
        actions={
          <button
            type="button"
            onClick={() => { setEditing(null); setForm({ ...BLANK }); setOpen(true); }}
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover"
          >
            <Plus className="h-4 w-4" /> New connection
          </button>
        }
      />

      <div className="rounded-xl border border-border bg-bg-elevated p-6">
        {listQ.isLoading ? (
          <p className="flex items-center gap-2 text-sm text-text-secondary">
            <Loader2 className="h-4 w-4 animate-spin" /> Loading&hellip;
          </p>
        ) : rows.length === 0 ? (
          <div className="py-8 text-center">
            <ShieldCheck className="mx-auto mb-3 h-8 w-8 text-text-muted" />
            <p className="text-sm font-medium text-foreground">No SSO connections yet</p>
            <p className="mt-1 text-sm text-text-secondary">
              Add one and anybody whose email matches its domains signs in through that provider
              instead of a password.
            </p>
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border text-left text-xs uppercase tracking-wider text-text-muted">
                <th className="pb-3">Connection</th>
                <th className="pb-3">Domains</th>
                <th className="pb-3">Status</th>
                <th className="pb-3 text-right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((c) => (
                <tr key={c.id} className="border-b border-border/50">
                  <td className="py-3">
                    <span className="block font-medium text-foreground">{c.name}</span>
                    <span className="block text-xs text-text-muted">{c.issuer_url}</span>
                  </td>
                  <td className="py-3 text-text-secondary">
                    {c.domains || <span className="text-text-muted">&mdash;</span>}
                    <span className="ml-2 rounded bg-bg-tertiary px-1.5 py-0.5 text-[10px] uppercase tracking-wide text-text-muted">
                      {c.protocol || "oidc"}
                    </span>
                  </td>
                  <td className="py-3">
                    {c.enabled ? (
                      <span className="inline-flex items-center gap-1 text-xs text-success">
                        <Check className="h-3.5 w-3.5" /> Enabled
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 text-xs text-text-muted">
                        <AlertTriangle className="h-3.5 w-3.5" /> Disabled
                      </span>
                    )}
                    {!c.has_secret && <span className="ml-2 text-xs text-danger">no secret</span>}
                  </td>
                  <td className="py-3 text-right">
                    <button onClick={() => copyCallback(c)} className="mr-3 text-xs text-text-secondary hover:text-accent">
                      <Copy className="mr-1 inline h-3 w-3" />IdP URLs
                    </button>
                    <button onClick={() => startEdit(c)} className="mr-3 text-xs text-text-secondary hover:text-accent">
                      Edit
                    </button>
                    <button
                      onClick={() => {
                        if (window.confirm("Delete " + c.name + "? Their users fall back to password login.")) {
                          deleteM.mutate(c.id);
                        }
                      }}
                      className="text-xs text-text-secondary hover:text-danger"
                    >
                      <Trash2 className="mr-1 inline h-3 w-3" />Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <ResponsiveSheet
        open={open}
        onClose={() => setOpen(false)}
        title={editing ? "Edit connection" : "New SSO connection"}
        description="Anything with an OpenID Connect discovery document works — Okta, Entra ID, Auth0, Keycloak, Google Workspace."
        size="lg"
        footer={
          <div className="flex items-center justify-end gap-3">
            <button onClick={() => setOpen(false)} className="rounded-lg border border-border px-4 py-2 text-sm text-text-secondary hover:bg-bg-hover">
              Cancel
            </button>
            <button
              onClick={() => saveM.mutate()}
              disabled={saveM.isPending}
              className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white hover:bg-accent-hover disabled:opacity-50"
            >
              {saveM.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
              {editing ? "Save" : "Create"}
            </button>
          </div>
        }
      >
        <div className="space-y-4">
          <Field label="Name" hint="Shown to users on the login page.">
            <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} placeholder="Acme Corp" className={inputCls} />
          </Field>
          <Field label="Slug" hint="Used in the callback URL. Lowercase, no spaces. Cannot change after creation.">
            <input value={form.slug} onChange={(e) => setForm({ ...form, slug: e.target.value })} disabled={!!editing} placeholder="acme" className={inputCls} />
          </Field>
          <Field label="Protocol" hint="OIDC works with every modern IdP. Choose SAML only if the customer's provider requires it.">
            <select
              value={form.protocol}
              onChange={(e) => setForm({ ...form, protocol: e.target.value })}
              disabled={!!editing}
              className={inputCls}
            >
              <option value="oidc">OpenID Connect</option>
              <option value="saml">SAML 2.0</option>
            </select>
          </Field>

          <Field label="Email domains" hint="Comma-separated. Anyone with an address at these domains is sent to this provider.">
            <input value={form.domains} onChange={(e) => setForm({ ...form, domains: e.target.value })} placeholder="acme.com, acme.co.uk" className={inputCls} />
          </Field>
          {form.protocol === "saml" ? (
            <>
              <Field label="IdP metadata URL" hint="Where the provider publishes its metadata. Leave blank if pasting the XML below.">
                <input value={form.metadata_url} onChange={(e) => setForm({ ...form, metadata_url: e.target.value })} placeholder="https://acme.okta.com/app/xxx/sso/saml/metadata" className={inputCls} />
              </Field>
              <Field label="IdP metadata XML" hint="Paste the document instead of fetching it. A pasted document wins — it is pinned, where a URL can start serving something different tomorrow.">
                <textarea rows={4} value={form.metadata_xml} onChange={(e) => setForm({ ...form, metadata_xml: e.target.value })} placeholder="<EntityDescriptor ...>" className={inputCls} />
              </Field>
              <Field label="Email attribute" hint="Blank tries the usual names (email, mail, the Microsoft claim URI).">
                <input value={form.email_attribute} onChange={(e) => setForm({ ...form, email_attribute: e.target.value })} className={inputCls} />
              </Field>
              <Field label="Groups attribute" hint="Blank tries groups, memberOf and the Microsoft groups claim.">
                <input value={form.groups_attribute} onChange={(e) => setForm({ ...form, groups_attribute: e.target.value })} className={inputCls} />
              </Field>
              <label className="flex items-start gap-2 text-sm">
                <input type="checkbox" checked={form.allow_idp_initiated} onChange={(e) => setForm({ ...form, allow_idp_initiated: e.target.checked })} className="mt-0.5" />
                <span>
                  <span className="font-medium text-foreground">Allow IdP-initiated sign-in</span>
                  <span className="block text-xs text-text-muted">
                    Lets users start from their provider&rsquo;s app tile rather than this login page. Most enterprise users sign in that way. Assertions are still signature-checked and time-bounded.
                  </span>
                </span>
              </label>
            </>
          ) : (
          <>
          <Field label="Issuer URL" hint="Discovery is fetched from /.well-known/openid-configuration beneath it.">
            <input value={form.issuer_url} onChange={(e) => setForm({ ...form, issuer_url: e.target.value })} placeholder="https://acme.okta.com" className={inputCls} />
          </Field>
          <Field label="Client ID">
            <input value={form.client_id} onChange={(e) => setForm({ ...form, client_id: e.target.value })} className={inputCls} />
          </Field>
          <Field
            label="Client secret"
            hint={editing ? "Leave blank to keep the stored secret." : "Stored encrypted at rest and never shown again."}
          >
            <input type="password" value={form.client_secret} onChange={(e) => setForm({ ...form, client_secret: e.target.value })} className={inputCls} />
          </Field>
          </>
          )}
          <Field label="Groups claim" hint="The claim carrying the user's groups: groups for Okta/Keycloak, roles for Entra ID.">
            <input value={form.groups_claim} onChange={(e) => setForm({ ...form, groups_claim: e.target.value })} className={inputCls} />
          </Field>
          <Field
            label="Group to role mapping"
            hint="JSON mapping IdP groups to role names. Roles are re-applied on every login, so removing someone from a group revokes it here too."
          >
            <textarea
              rows={3}
              value={form.group_mappings}
              onChange={(e) => setForm({ ...form, group_mappings: e.target.value })}
              placeholder={JSON.stringify({ "it-admins": "ADMIN", engineering: "EDITOR" }, null, 2)}
              className={inputCls}
            />
          </Field>

          <label className="flex items-start gap-2 text-sm">
            <input
              type="checkbox"
              checked={form.jit_provisioning}
              onChange={(e) => setForm({ ...form, jit_provisioning: e.target.checked })}
              className="mt-0.5"
            />
            <span>
              <span className="font-medium text-foreground">Create accounts on first login</span>
              <span className="block text-xs text-text-muted">
                Off means the user must already exist here before they can sign in.
              </span>
            </span>
          </label>
          <label className="flex items-center gap-2 text-sm">
            <input type="checkbox" checked={form.enabled} onChange={(e) => setForm({ ...form, enabled: e.target.checked })} />
            <span className="font-medium text-foreground">Enabled</span>
          </label>
        </div>
      </ResponsiveSheet>
    </div>
  );
}
`
	return strings.ReplaceAll(src, "~", "`")
}
