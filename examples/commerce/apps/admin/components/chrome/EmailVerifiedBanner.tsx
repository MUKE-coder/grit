"use client";

import { useEffect, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { useMe } from "@/hooks/use-auth";
import { apiClient } from "@/lib/api-client";
import { AlertTriangle, Check, Loader2, X } from "@/lib/icons";

const DISMISSED_KEY = "grit.verify-email.dismissed";

/**
 * Shown to a signed-in user whose address is still unconfirmed.
 *
 * Renders nothing at all when there is nothing to say, including while the
 * user query is still loading, so it never flashes in and out on every page
 * load for people who verified months ago.
 *
 * Dismissal is per browser session, not permanent. An unconfirmed address
 * means a password reset has nowhere to go, so the banner earns its place;
 * what it does not deserve is to be unclosable while somebody is trying to
 * read the page behind it.
 */
export function EmailVerifiedBanner() {
  const { data: user, isLoading } = useMe();
  const [sent, setSent] = useState(false);
  const [devLink, setDevLink] = useState<string | null>(null);
  const [dismissed, setDismissed] = useState(false);

  // Read after mount: sessionStorage does not exist while this renders on the
  // server, and reading it during render would mismatch the hydration.
  useEffect(() => {
    try {
      setDismissed(window.sessionStorage.getItem(DISMISSED_KEY) === "1");
    } catch {
      // Private windows and blocked site data. Showing the banner is the
      // safe answer.
    }
  }, []);

  const dismiss = () => {
    setDismissed(true);
    try {
      window.sessionStorage.setItem(DISMISSED_KEY, "1");
    } catch {
      // Closed for this view either way.
    }
  };

  const resend = useMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post("/api/auth/verify-email/send", {});
      return (data as { verify_url?: string }) ?? {};
    },
    onSuccess: (data) => {
      setSent(true);
      // Development only: the server hands the link back rather than posting
      // it, so verifying does not mean digging through storage/mail.
      if (data.verify_url) setDevLink(data.verify_url);
    },
  });

  if (isLoading || !user || user.email_verified_at || dismissed) return null;

  return (
    <div className="mx-6 mt-4 rounded-lg border border-warning/40 bg-warning/[0.07] px-4 py-3">
      <div className="flex flex-wrap items-center gap-3">
        <AlertTriangle className="h-4 w-4 shrink-0 text-warning" />
        <p className="min-w-0 flex-1 text-sm text-foreground">
          Confirm your email address.{" "}
          <span className="text-text-secondary">
            We sent a link to {user.email} when you signed up.
          </span>
        </p>

        {sent && !devLink ? (
          <span className="inline-flex items-center gap-1.5 text-sm text-success">
            <Check className="h-3.5 w-3.5" />
            Sent, check your inbox
          </span>
        ) : !sent ? (
          <button
            type="button"
            onClick={() => resend.mutate()}
            disabled={resend.isPending}
            className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-3 py-1.5 text-sm hover:bg-bg-hover disabled:opacity-50"
          >
            {resend.isPending && <Loader2 className="h-3.5 w-3.5 animate-spin" />}
            Resend link
          </button>
        ) : null}

        <button
          type="button"
          onClick={dismiss}
          aria-label="Dismiss until next sign-in"
          className="-m-1.5 shrink-0 rounded-md p-1.5 text-text-muted transition-colors hover:bg-bg-hover hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent"
        >
          <X className="h-4 w-4" aria-hidden="true" />
        </button>
      </div>

      {devLink && (
        <p className="mt-2 pl-7 text-sm">
          <span className="text-text-muted">Development: </span>
          <a href={devLink} className="font-medium text-accent underline underline-offset-2">
            open the verification link
          </a>
        </p>
      )}
    </div>
  );
}
