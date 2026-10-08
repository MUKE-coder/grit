"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

// Merged into /account. Kept as a redirect rather than deleted: this path was
// the user menu's target for years, it is where a USER lands after signing in,
// and it is in release notes and bookmarks.
export default function SystemAccountRedirect() {
	const router = useRouter();
	useEffect(() => {
		// With the fragment, or the dashboard's "Set it up" lands at the top of a
		// long page instead of the two-factor card it was pointing at. A redirect
		// that drops the anchor is a redirect that half worked.
		router.replace("/account" + window.location.hash);
	}, [router]);

	return (
		<p className="p-8 text-sm text-text-secondary">
			Moved to <a className="underline" href="/account">your account</a>.
		</p>
	);
}
