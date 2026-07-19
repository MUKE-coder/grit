package scaffold

import "strings"

// adminUseRoles emits hooks/use-roles.ts — React Query hooks for the roles API.
//
// Note what is NOT here: any wildcard-matching logic. The API serves each role's
// grants already expanded, so the client renders checkboxes straight from
// `expanded` and never reimplements the matcher. Keeping a second copy of those
// rules in TypeScript is how the system this was modelled on ended up with a
// client that disagreed with its server about who could do what.
func adminUseRoles() string {
	src := `"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export type PermAction = "create" | "view" | "edit" | "delete";

export interface PermFeature {
	key: string;
	name: string;
	actions: PermAction[];
}
export interface PermGroup {
	key: string;
	name: string;
	features: PermFeature[];
}
export interface PermModule {
	key: string;
	name: string;
	groups: PermGroup[];
}

export interface Role {
	id: string;
	name: string;
	description: string;
	/** As authored — may contain wildcards like "products.*" or "*". */
	grants: string[];
	/** Wildcards resolved by the server. Render from this. */
	expanded: string[];
	is_system: boolean;
	user_count: number;
}

export function usePermissionCatalog() {
	return useQuery({
		queryKey: ["permission-catalog"],
		// The catalog only changes when code changes, so it can be cached hard.
		staleTime: 5 * 60 * 1000,
		queryFn: async () => {
			const res = await apiClient.get("/api/permissions");
			return res.data.data as { modules: PermModule[]; keys: string[] };
		},
	});
}

export function useRoles() {
	return useQuery({
		queryKey: ["roles"],
		queryFn: async () => {
			const res = await apiClient.get("/api/roles");
			return res.data.data as Role[];
		},
	});
}

export function useCreateRole() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: async (input: { name: string; description: string; grants: string[] }) => {
			const res = await apiClient.post("/api/roles", input);
			return res.data.data as Role;
		},
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["roles"] });
			// A role change can alter the current user's own permissions, which
			// drive nav visibility — refetch those too.
			qc.invalidateQueries({ queryKey: ["my-permissions"] });
		},
	});
}

export function useUpdateRole() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: async (input: { id: string; name: string; description: string; grants: string[] }) => {
			const { id, ...body } = input;
			const res = await apiClient.put("/api/roles/" + id, body);
			return res.data.data as Role;
		},
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["roles"] });
			qc.invalidateQueries({ queryKey: ["my-permissions"] });
		},
	});
}

export function useDeleteRole() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: async (id: string) => {
			await apiClient.delete("/api/roles/" + id);
		},
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["roles"] });
			qc.invalidateQueries({ queryKey: ["my-permissions"] });
		},
	});
}

export function useAssignUserRoles() {
	const qc = useQueryClient();
	return useMutation({
		mutationFn: async (input: { userId: string; roleIds: string[] }) => {
			await apiClient.put("/api/users/" + input.userId + "/roles", { role_ids: input.roleIds });
		},
		onSuccess: () => {
			qc.invalidateQueries({ queryKey: ["users"] });
			qc.invalidateQueries({ queryKey: ["roles"] });
			qc.invalidateQueries({ queryKey: ["my-permissions"] });
		},
	});
}

/** Every concrete key a feature contributes, e.g. ["users.view","users.edit"]. */
export function featureKeys(f: PermFeature): string[] {
	return f.actions.map((a) => f.key + "." + a);
}

export function groupKeys(g: PermGroup): string[] {
	return g.features.flatMap(featureKeys);
}

export function moduleKeys(m: PermModule): string[] {
	return m.groups.flatMap(groupKeys);
}

/**
 * Collapse a selection back to the shortest equivalent grant list.
 *
 * When every action of a resource is selected we store "<resource>.*" rather
 * than the four leaves. That is what lets a role keep working when a new action
 * is added to the catalog later — storing expanded leaves is precisely why the
 * reference implementation's roles silently stopped inheriting new permissions.
 */
export function collapseGrants(selected: Set<string>, modules: PermModule[]): string[] {
	const out: string[] = [];
	const covered = new Set<string>();

	for (const m of modules) {
		for (const g of m.groups) {
			for (const f of g.features) {
				const keys = featureKeys(f);
				if (keys.length > 0 && keys.every((k) => selected.has(k))) {
					out.push(f.key + ".*");
					keys.forEach((k) => covered.add(k));
				}
			}
		}
	}
	for (const k of selected) {
		if (!covered.has(k)) out.push(k);
	}
	return out.sort();
}
`
	return strings.ReplaceAll(src, "~", "`")
}
