package scaffold

// Roles & permissions for the Wails desktop client.
//
// Desktop knew only the coarse user.role string: no usePermissions, no call to
// /auth/permissions, and a hardcoded USER/EDITOR/ADMIN <select> on the users
// screen — so a role defined in the admin could never be assigned from the
// desktop app. These screens use the same endpoints and the same wildcard
// semantics as the web admin and the mobile app.
//
// Note the paths: apiClient's baseURL already ends in /api, so requests are
// "/roles", not "/api/roles" and not "/admin/roles".
//
// This is a UX layer. Every route is enforced server-side.

// desktopClientUsePermissions emits hooks/use-permissions.ts.
func desktopClientUsePermissions() string {
	return `import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

interface MyPermissions {
  grants: string[];
  permissions: string[];
  is_super: boolean;
}

export function usePermissions() {
  const { data, isLoading } = useQuery<MyPermissions>({
    queryKey: ["my-permissions"],
    staleTime: 60_000,
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: MyPermissions }>("/auth/permissions");
      return data.data;
    },
  });

  const granted = new Set(data?.permissions ?? []);
  const isSuper = data?.is_super ?? false;

  /**
   * can("users.delete") — exact permission
   * can("users.*")      — any permission on that resource
   *
   * False while loading, so gated UI stays hidden rather than flashing in.
   */
  function can(permission: string): boolean {
    if (isSuper) return true;
    if (permission.endsWith(".*")) {
      const prefix = permission.slice(0, -1);
      for (const p of granted) if (p.startsWith(prefix)) return true;
      return false;
    }
    return granted.has(permission);
  }

  return { can, isSuper, isLoading, permissions: data?.permissions ?? [] };
}
`
}

// desktopClientSystemRolesPage emits routes/app/system.roles.tsx.
func desktopClientSystemRolesPage() string {
	return `import { useState, useEffect } from "react";
import { createFileRoute } from "@tanstack/react-router";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Shield, Plus, Trash2 } from "lucide-react";
import { PageHeader } from "@/components/layout/page-header";
import { SectionCard, EmptyState } from "@/components/system-ui";
import { Button } from "@/components/ui/button";
import { usePermissions } from "@/hooks/use-permissions";
import { apiClient } from "@/lib/api-client";

export const Route = createFileRoute("/app/system/roles")({
  component: SystemRolesPage,
});

interface Feature { key: string; name: string; actions: string[] }
interface Group { key: string; name: string; features: Feature[] }
interface Module { key: string; name: string; groups: Group[] }
interface Catalog { keys: string[]; modules: Module[] }

interface Role {
  id: string;
  name: string;
  description: string;
  grants: string[];
  expanded: string[];
  is_system: boolean;
  user_count: number;
}

const inputCls =
  "w-full rounded-lg border border-border bg-surface px-3 py-2 text-[13px] text-foreground outline-none focus:border-accent";

function SystemRolesPage() {
  const queryClient = useQueryClient();
  const { can, isLoading: permsLoading } = usePermissions();
  const [editing, setEditing] = useState<Role | "new" | null>(null);

  const { data: roles } = useQuery<Role[]>({
    queryKey: ["system", "roles"],
    enabled: can("roles.view"),
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: Role[] }>("/roles");
      return data.data ?? [];
    },
  });

  if (permsLoading) return null;

  if (!can("roles.view")) {
    return (
      <div>
        <PageHeader title="Roles & Permissions" description="Control exactly what each role can see and do" />
        <div className="mt-6">
          <EmptyState
            icon={Shield}
            title="Not permitted"
            hint="You do not have permission to view roles."
          />
        </div>
      </div>
    );
  }

  return (
    <div>
      <PageHeader
        title="Roles & Permissions"
        description="Control exactly what each role can see and do"
        actions={
          can("roles.create") ? (
            <Button onClick={() => setEditing("new")}>
              <Plus className="mr-1.5 h-3.5 w-3.5" />
              New role
            </Button>
          ) : undefined
        }
      />

      <div className="mt-6 grid gap-4 lg:grid-cols-2">
        <SectionCard title="Roles" description="A role granted a whole resource keeps actions added to it later">
          <div className="divide-y divide-border">
            {(roles ?? []).map((role) => (
              <button
                key={role.id}
                onClick={() => setEditing(role)}
                className="flex w-full items-center justify-between px-1 py-3 text-left hover:bg-surface-hover"
              >
                <div>
                  <div className="flex items-center gap-2">
                    <span className="text-[13px] font-medium text-foreground">{role.name}</span>
                    {role.is_system && (
                      <span className="rounded-full bg-surface-hover px-2 py-0.5 text-[11px] text-foreground-muted">
                        Built-in
                      </span>
                    )}
                  </div>
                  <p className="mt-0.5 text-[12px] text-foreground-muted">{role.description}</p>
                </div>
                <div className="text-right">
                  <div className="text-[12px] text-foreground-muted">
                    {role.user_count} {role.user_count === 1 ? "user" : "users"}
                  </div>
                  <div className="text-[12px] text-accent">
                    {role.grants?.includes("*") ? "all permissions" : (role.expanded?.length ?? 0) + " granted"}
                  </div>
                </div>
              </button>
            ))}
            {roles?.length === 0 && (
              <p className="py-6 text-center text-[12px] text-foreground-muted">No roles yet.</p>
            )}
          </div>
        </SectionCard>

        {editing ? (
          <RoleEditor
            role={editing === "new" ? null : editing}
            onDone={async () => {
              setEditing(null);
              await queryClient.invalidateQueries({ queryKey: ["system", "roles"] });
            }}
            onCancel={() => setEditing(null)}
          />
        ) : (
          <SectionCard title="Permissions" description="Select a role to edit what it can do">
            <p className="py-6 text-center text-[12px] text-foreground-muted">No role selected.</p>
          </SectionCard>
        )}
      </div>
    </div>
  );
}

function RoleEditor({
  role,
  onDone,
  onCancel,
}: {
  role: Role | null;
  onDone: () => void;
  onCancel: () => void;
}) {
  const { can } = usePermissions();
  const [name, setName] = useState(role?.name ?? "");
  const [description, setDescription] = useState(role?.description ?? "");
  const [selected, setSelected] = useState<Set<string>>(new Set(role?.expanded ?? []));
  // Permission modules are collapsed by default; expand one at a time.
  const [openModules, setOpenModules] = useState<Set<string>>(new Set());
  const toggleModule = (key: string) =>
    setOpenModules((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");

  // Reseed when a different role is picked — the editor is reused, not remounted.
  useEffect(() => {
    setName(role?.name ?? "");
    setDescription(role?.description ?? "");
    setSelected(new Set(role?.expanded ?? []));
    setError("");
  }, [role]);

  const { data: catalog } = useQuery<Catalog>({
    queryKey: ["permission-catalog"],
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: Catalog }>("/permissions");
      return data.data;
    },
  });

  function toggle(key: string) {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(key)) next.delete(key);
      else next.add(key);
      return next;
    });
  }

  function toggleFeature(feature: Feature) {
    const keys = feature.actions.map((a) => feature.key + "." + a);
    const allOn = keys.every((k) => selected.has(k));
    setSelected((prev) => {
      const next = new Set(prev);
      for (const k of keys) {
        if (allOn) next.delete(k);
        else next.add(k);
      }
      return next;
    });
  }

  async function save() {
    setError("");
    if (!name.trim()) {
      setError("Name is required");
      return;
    }
    setSaving(true);
    try {
      const body = { name: name.trim(), description: description.trim(), grants: Array.from(selected) };
      if (role) await apiClient.put("/roles/" + role.id, body);
      else await apiClient.post("/roles", body);
      onDone();
    } catch (e: any) {
      setError(e?.response?.data?.error?.message ?? "Failed to save role");
    } finally {
      setSaving(false);
    }
  }

  async function remove() {
    if (!role) return;
    setSaving(true);
    try {
      await apiClient.delete("/roles/" + role.id);
      onDone();
    } catch (e: any) {
      setError(e?.response?.data?.error?.message ?? "Failed to delete role");
    } finally {
      setSaving(false);
    }
  }

  return (
    <SectionCard
      title={role ? role.name : "New role"}
      description={selected.size + " of " + (catalog?.keys?.length ?? 0) + " permissions granted"}
    >
      {error && <p className="mb-3 text-[12px] text-danger">{error}</p>}

      <label className="mb-1 block text-[12px] text-foreground-muted">Name</label>
      <input
        className={inputCls}
        value={name}
        onChange={(e) => setName(e.target.value)}
        disabled={!!role?.is_system}
        placeholder="EDITOR"
      />

      <label className="mb-1 mt-3 block text-[12px] text-foreground-muted">Description</label>
      <input
        className={inputCls}
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        placeholder="What this role is for"
      />

      {role?.is_system && (
        <p className="mt-2 text-[12px] text-foreground-muted">
          Built-in role — the name is fixed, but its permissions can still be changed.
        </p>
      )}

      <div className="mt-5 max-h-[420px] overflow-y-auto pr-1">
        {catalog?.modules?.map((mod) => {
          const modKeys = mod.groups.flatMap((g) =>
            g.features.flatMap((f) => f.actions.map((a) => f.key + "." + a))
          );
          const modOn = modKeys.filter((k) => selected.has(k)).length;
          const isOpen = openModules.has(mod.key);
          return (
          <div key={mod.key} className="mb-3">
            {/* Collapsed by default — click to expand one module at a time. */}
            <button
              onClick={() => toggleModule(mod.key)}
              className="mb-2 flex w-full items-center justify-between text-left"
            >
              <span className="flex items-center gap-1.5">
                <span className={"text-foreground-muted transition-transform " + (isOpen ? "rotate-90" : "")}>&rsaquo;</span>
                <span className="text-[11px] font-semibold uppercase tracking-wide text-foreground-muted">{mod.name}</span>
              </span>
              <span className="text-[11px] text-accent">{modOn}/{modKeys.length}</span>
            </button>
            {isOpen && mod.groups.map((group) =>
              group.features.map((feature) => {
                const keys = feature.actions.map((a) => feature.key + "." + a);
                const onCount = keys.filter((k) => selected.has(k)).length;
                return (
                  <div key={feature.key} className="mb-2 rounded-lg border border-border p-3">
                    <button
                      onClick={() => toggleFeature(feature)}
                      className="mb-2 flex w-full items-center justify-between text-left"
                    >
                      <span className="text-[13px] font-medium text-foreground">{feature.name}</span>
                      <span className="text-[12px] text-accent">
                        {onCount}/{keys.length}
                      </span>
                    </button>
                    <div className="flex flex-wrap gap-1.5">
                      {feature.actions.map((action) => {
                        const key = feature.key + "." + action;
                        const on = selected.has(key);
                        return (
                          <button
                            key={key}
                            onClick={() => toggle(key)}
                            className={
                              on
                                ? "rounded-full bg-accent px-2.5 py-1 text-[12px] capitalize text-white"
                                : "rounded-full bg-surface-hover px-2.5 py-1 text-[12px] capitalize text-foreground-muted"
                            }
                          >
                            {action}
                          </button>
                        );
                      })}
                    </div>
                  </div>
                );
              })
            )}
          </div>
          );
        })}
      </div>

      <div className="mt-4 flex items-center gap-2">
        <Button onClick={save} disabled={saving}>
          {saving ? "Saving…" : "Save role"}
        </Button>
        <Button variant="ghost" onClick={onCancel}>
          Cancel
        </Button>
        {role && !role.is_system && can("roles.delete") && (
          <Button variant="ghost" onClick={remove} className="ml-auto text-danger">
            <Trash2 className="mr-1.5 h-3.5 w-3.5" />
            Delete
          </Button>
        )}
      </div>
    </SectionCard>
  );
}
`
}
