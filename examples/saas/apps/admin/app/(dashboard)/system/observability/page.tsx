"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

// Merged into the operations page: everything this showed, and the latency,
// throughput and slow routes it never managed to.
export default function ObservabilityRedirect() {
	const router = useRouter();
	useEffect(() => {
		router.replace("/system/performance");
	}, [router]);

	return (
		<p className="p-8 text-sm text-text-secondary">
			Observability has moved to <a className="underline" href="/system/performance">Operations</a>.
		</p>
	);
}
