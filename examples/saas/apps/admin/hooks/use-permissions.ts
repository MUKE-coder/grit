"use client";

import { useCallback } from "react";
import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

interface MyPermissions {
	grants: string[];
	permissions: string[];
	is_super: boolean;
}

interface PermissionState {
	granted: Set<string>;
	isSuper: boolean;
	permissions: string[];
}

const noPermissions: string[] = [];
const noGrants = new Set<string>();

// Module scope on purpose: React Query reruns select only when the response or
// the function itself changes, so the Set is built once per fetch instead of
// once per render of every component that asks can().
function toPermissionState(data: MyPermissions): PermissionState {
	return {
		granted: new Set(data.permissions ?? []),
		isSuper: data.is_super ?? false,
		permissions: data.permissions ?? noPermissions,
	};
}

export function usePermissions() {
	const { data, isLoading } = useQuery({
		queryKey: ["my-permissions"],
		staleTime: 60 * 1000,
		queryFn: async () => {
			const res = await apiClient.get("/api/auth/permissions");
			return res.data.data as MyPermissions;
		},
		select: toPermissionState,
	});

	const granted = data?.granted ?? noGrants;
	const isSuper = data?.isSuper ?? false;

	/**
	 * can("users.delete")  — exact permission
	 * can("users.*")       — any permission on that resource
	 *
	 * Returns false while loading. Nav items and action buttons therefore stay
	 * hidden until permissions are known, rather than flashing into view and
	 * disappearing — a flash of forbidden UI looks broken and leaks the shape of
	 * the admin to users who can't use it.
	 */
	const can = useCallback(
		(permission: string): boolean => {
			if (isSuper) return true;
			if (permission.endsWith(".*")) {
				const prefix = permission.slice(0, -1); // "users."
				for (const p of granted) {
					if (p.startsWith(prefix)) return true;
				}
				return false;
			}
			return granted.has(permission);
		},
		[granted, isSuper],
	);

	return { can, isSuper, isLoading, permissions: data?.permissions ?? noPermissions };
}
