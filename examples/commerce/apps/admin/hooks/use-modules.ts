"use client";

import { useQuery } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export function useModules() {
	const { data, isLoading } = useQuery({
		queryKey: ["system-modules"],
		// Module flags come from the server's env and only change on restart.
		staleTime: 5 * 60 * 1000,
		queryFn: async () => {
			const res = await apiClient.get("/api/system/modules");
			return res.data.data as Record<string, boolean>;
		},
	});

	/**
	 * Fails OPEN, unlike can(): an unknown module, or one whose flags haven't
	 * loaded yet, is treated as enabled. Hiding navigation because a request is
	 * slow would look like the feature vanished — and unlike permissions, a
	 * visible link to a disabled module is a 404, not a privilege leak.
	 */
	function moduleEnabled(name: string): boolean {
		if (isLoading || !data) return true;
		return data[name] !== false;
	}

	return { moduleEnabled, modules: data ?? {}, isLoading };
}
