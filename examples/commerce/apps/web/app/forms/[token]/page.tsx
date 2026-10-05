import type { Metadata } from "next";

import { apiUrl } from "@/lib/api-core";

import { PublicForm, type ShareInfo } from "./public-form";

export const metadata: Metadata = {
  title: "Submit a form",
  robots: { index: false, follow: false },
};

// A share can be revoked at any moment, so nothing about this page is
// prerendered or revalidated.
export const dynamic = "force-dynamic";

interface PageProps {
  params: Promise<{ token: string }>;
}

interface LoadResult {
  info?: ShareInfo;
  error?: string;
}

// loadShare asks the API about the link. A 404 or a disabled share is a normal
// answer here, not an exception: it renders the "link unavailable" card.
async function loadShare(token: string): Promise<LoadResult> {
  try {
    const res = await fetch(
      apiUrl("/api/public/forms/" + encodeURIComponent(token)),
      { cache: "no-store", headers: { Accept: "application/json" } },
    );
    const body = await res.json().catch(() => null);
    if (!res.ok) {
      return { error: body?.error?.message ?? "Link not found or disabled" };
    }
    if (!body?.data) {
      return { error: "Link not found or disabled" };
    }
    return { info: body.data as ShareInfo };
  } catch {
    // The API being unreachable is the operator's problem, not the visitor's,
    // so they get a sentence rather than a stack.
    return { error: "This form could not be loaded. Try again in a moment." };
  }
}

export default async function PublicFormPage({ params }: PageProps) {
  const { token } = await params;
  const { info, error } = await loadShare(token);

  if (!info) {
    return (
      <main className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
        <div className="w-full max-w-md rounded-2xl border border-slate-200 bg-white p-8 text-center shadow-sm">
          <h1 className="text-xl font-semibold text-slate-900">Link unavailable</h1>
          <p className="mt-2 text-sm text-slate-500">{error ?? "Unknown error"}</p>
        </div>
      </main>
    );
  }

  // Title falls back through three sources:
  //   1. operator-set custom_title (best)
  //   2. operator-set label (legacy — pre-v3.31.50 shares)
  //   3. resource name (worst — bare default)
  const title =
    info.custom_title?.trim() ||
    info.label?.trim() ||
    info.resource_name + " submission";
  const description =
    info.custom_description?.trim() ||
    "Fill out the form below to submit a new " + info.resource_name + ".";

  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
      <div className="w-full max-w-md rounded-2xl border border-slate-200 bg-white p-8 shadow-sm">
        <h1 className="text-2xl font-semibold text-slate-900">{title}</h1>
        <p className="mt-1 text-sm text-slate-500">{description}</p>

        <PublicForm token={token} info={info} />
      </div>
    </main>
  );
}
