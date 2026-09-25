package scaffold

// The Sign-in links card on the Account screen's Security tab.
//
// A sign-in link is a bearer credential sitting in a mailbox, and the person
// whose mailbox it is has no other way to notice somebody keeps asking for
// one: the request form answers the same to every address by design, which is
// what stops it enumerating accounts, and which also means it can never warn
// anybody. This card is where they find out. It shows when a link was asked
// for and from where, and never the token.

func adminSignInLinksCardTSX() string {
	return `"use client";

import { useQuery } from "@tanstack/react-query";
import { Link as LinkIcon, ShieldCheck } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";

interface LinkActivity {
  created_at: string;
  used_at?: string | null;
  expires_at: string;
  ip_address?: string;
  user_agent?: string;
}

function when(iso: string): string {
  const then = new Date(iso).getTime();
  const mins = Math.round((Date.now() - then) / 60000);
  if (mins < 1) return "just now";
  if (mins < 60) return mins + " minute" + (mins === 1 ? "" : "s") + " ago";
  const hours = Math.round(mins / 60);
  if (hours < 24) return hours + " hour" + (hours === 1 ? "" : "s") + " ago";
  return new Date(iso).toLocaleDateString();
}

function outcome(row: LinkActivity): string {
  if (row.used_at) return "used";
  if (new Date(row.expires_at).getTime() < Date.now()) return "expired unused";
  return "still valid";
}

/**
 * Sign-in links, on the Security tab.
 *
 * Read-only on purpose. There is no button here to send yourself a link,
 * because you are already signed in, and no toggle to turn the feature off,
 * because for an account created through Google or GitHub a link is the only
 * way in that does not depend on that provider still working.
 */
export function SignInLinksCard() {
  const { data, isLoading } = useQuery({
    queryKey: ["magic-link-activity"],
    queryFn: async () => {
      const res = await apiClient.get("/api/auth/magic-link/recent");
      return ((res.data as { data?: LinkActivity[] })?.data ?? []) as LinkActivity[];
    },
  });

  const rows = data ?? [];

  return (
    <section className="rounded-xl border border-border bg-bg-elevated p-6">
      <div className="flex items-start gap-3">
        <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent/10 text-accent">
          <LinkIcon className="h-4 w-4" aria-hidden="true" />
        </span>
        <div className="min-w-0 flex-1">
          <h2 className="text-base font-semibold text-foreground">Sign-in links</h2>
          <p className="mt-1 text-sm leading-relaxed text-text-muted">
            You can sign in from the login page without a password, using a link emailed to this
            address. A link lasts 15 minutes, works once, and never gets past two-factor.
          </p>

          <h4 className="mt-5 text-sm font-medium text-foreground">Recently requested</h4>
          {isLoading ? (
            <p className="mt-2 text-sm text-text-secondary">Loading...</p>
          ) : rows.length === 0 ? (
            <p className="mt-2 text-sm text-text-secondary">
              Nobody has asked for a sign-in link for this account.
            </p>
          ) : (
            <ul className="mt-2 space-y-2">
              {rows.map((row) => (
                <li
                  key={row.created_at + (row.ip_address ?? "")}
                  className="flex flex-wrap items-baseline gap-x-2 gap-y-0.5 text-sm"
                >
                  <span className="text-foreground">{when(row.created_at)}</span>
                  <span className="text-text-secondary">from {row.ip_address || "an unknown address"}</span>
                  <span className="text-text-muted">&middot; {outcome(row)}</span>
                </li>
              ))}
            </ul>
          )}

          <p className="mt-4 flex items-start gap-2 text-sm text-text-secondary">
            <ShieldCheck className="mt-0.5 h-4 w-4 shrink-0 text-success" aria-hidden="true" />
            <span>
              If you see a request you did not make, somebody knows your email address, which is not
              a secret. They cannot sign in without the mailbox itself. Turning on two-factor here
              means even the mailbox is not enough.
            </span>
          </p>
        </div>
      </div>
    </section>
  );
}
`
}
