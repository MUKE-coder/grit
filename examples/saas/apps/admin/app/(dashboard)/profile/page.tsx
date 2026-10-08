"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

// Merged into /account. Kept as a redirect rather than deleted: this path was
// the user menu's target for years, it is where a USER lands after signing in,
// and it is in release notes and bookmarks.
export default function ProfileRedirect() {
	const router = useRouter();
	useEffect(() => {
		router.replace("/account");
	}, [router]);

	return (
		<p className="p-8 text-sm text-text-secondary">
			Moved to <a className="underline" href="/account">your account</a>.
		</p>
	);
}
