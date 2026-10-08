"use client";

import { useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function AccountSecurityPage() {
  const router = useRouter();

  useEffect(() => {
    router.replace("/account#security");
  }, [router]);

  return (
    <div className="p-6 text-sm text-text-secondary">
      This page moved.{" "}
      <Link href="/account#security" className="text-accent underline">
        Open your account
      </Link>
      .
    </div>
  );
}
