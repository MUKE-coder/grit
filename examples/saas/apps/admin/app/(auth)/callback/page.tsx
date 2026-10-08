"use client";

import { Suspense, useEffect, useRef } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { apiClient } from "@/lib/api-client";
import { useQueryClient } from "@tanstack/react-query";

// OAuth flow (Grit 3.27+):
//   1. /api/auth/oauth/:provider redirects to the provider
//   2. provider hits /api/auth/oauth/:provider/callback (the GO route)
//   3. that handler:
//        a. exchanges code for user info
//        b. finds-or-creates the user
//        c. generates JWT tokens
//        d. SETS HttpOnly grit_access + grit_refresh cookies via Set-Cookie
//        e. redirects to this page with NO query params
//   4. this page just probes /api/auth/me (which works because the cookies
//      from step d arrived on the redirect response) and routes the user
//      to the right area based on their role.
// The tokens never travel as URL params, never sit in browser history,
// never get logged in nginx access logs.
//
// useSearchParams forces a client bailout during prerender unless it's
// wrapped in a Suspense boundary. The inner component does the work;
// the page export just provides the boundary.
export default function AuthCallbackPage() {
  return (
    <Suspense fallback={<CallbackSpinner />}>
      <CallbackInner />
    </Suspense>
  );
}

function CallbackInner() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const queryClient = useQueryClient();
  const processed = useRef(false);

  useEffect(() => {
    if (processed.current) return;
    processed.current = true;

    const error = searchParams.get("error");

    if (error) {
      router.push("/login?error=" + encodeURIComponent(error));
      return;
    }

    // The HttpOnly auth cookies are already in place (set by the Go
    // OAuth handler before it 307-redirected here). Probe /me to grab
    // the user and decide where to route them.
    apiClient
      .get("/api/auth/me")
      .then(({ data }) => {
        const user = data.data;
        queryClient.setQueryData(["me"], user);
        router.push(user.role === "USER" ? "/account" : "/dashboard");
      })
      .catch(() => {
        // /me failed even though we just landed here — either the
        // cookies were dropped (cross-origin same-site weirdness),
        // CORS misconfigured, or the user is disabled. Send them
        // back to login with a generic message.
        router.push("/login?error=Authentication+failed");
      });
  }, [searchParams, router, queryClient]);

  return <CallbackSpinner />;
}

function CallbackSpinner() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="text-center">
        <div className="inline-flex h-10 w-10 animate-spin items-center justify-center rounded-full border-2 border-accent border-t-transparent" />
        <p className="mt-4 text-sm text-text-secondary">Signing you in...</p>
      </div>
    </div>
  );
}
