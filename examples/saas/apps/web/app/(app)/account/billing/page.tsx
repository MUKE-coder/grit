"use client";

import { useQuery } from "@tanstack/react-query";
import { CreditCard, Receipt } from "lucide-react";
import { api } from "@/lib/api";
import { formatMoney, type Invoice, type Subscription } from "@repo/shared/types";
import { API_ROUTES } from "@repo/shared/constants";

/**
 * Billing, from the customer's side.
 *
 * Every number on this page comes from the same endpoints the operator console
 * reads. There is no portal-only API and no second copy of the queries.
 *
 * The paths come from API_ROUTES in packages/shared, which `grit generate
 * resource` writes and keeps current. Worth using rather than typing the path:
 * a resource whose name is two words is mounted at /api/usage_records, not
 * usage-records, and a hand-typed guess is a 404 the page renders as "nothing
 * here yet". The version is left off on purpose, because the api client pins
 * it in an interceptor.
 *
 * What makes it the customer's billing and not everybody's is one thing, and it
 * is on the server: both resources were generated with `--owned-by user`, so
 * the service puts authz.ScopeOwned on the list and authz.Owns on the read. A
 * customer asking for the whole table gets their own rows back.
 *
 * That is worth knowing before you add a page here. The confinement is not in
 * this file, and it is not a filter in the request below. Writing
 * `?user_id=<me>` on the client would look identical and protect nobody,
 * because anybody can change it. The test that proves the real thing is
 * apps/api/internal/handlers/portal_isolation_test.go.
 */

function StatusPill({ status }: { status: Subscription["status"] | Invoice["status"] }) {
  const tone: Record<string, string> = {
    active: "bg-success/15 text-success",
    paid: "bg-success/15 text-success",
    trialing: "bg-info/15 text-info",
    open: "bg-warning/15 text-warning",
    past_due: "bg-danger/15 text-danger",
    canceled: "bg-bg-hover text-text-secondary",
    draft: "bg-bg-hover text-text-secondary",
    void: "bg-bg-hover text-text-secondary",
  };
  return (
    <span className={"rounded-full px-2.5 py-0.5 text-xs font-medium " + (tone[status] ?? "bg-bg-hover text-text-secondary")}>
      {status.replace("_", " ")}
    </span>
  );
}

function longDate(value: string | null | undefined) {
  if (!value) return "--";
  return new Date(value).toLocaleDateString(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

export default function BillingPage() {
  // page_size 1: a customer has one subscription. Asking for one row rather
  // than the whole table is the difference between a page that stays fast when
  // somebody has eighty invoices and one that does not.
  const subscription = useQuery({
    queryKey: ["portal", "subscription"],
    queryFn: async () => {
      const { data } = await api.get(`${API_ROUTES.SUBSCRIPTIONS.LIST}?page_size=1&sort=created_at&order=desc`);
      return (data?.data?.[0] ?? null) as Subscription | null;
    },
  });

  const invoices = useQuery({
    queryKey: ["portal", "invoices"],
    queryFn: async () => {
      const { data } = await api.get(`${API_ROUTES.INVOICES.LIST}?page_size=12&sort=issued_on&order=desc`);
      return (data?.data ?? []) as Invoice[];
    },
  });

  const current = subscription.data;
  const plan = current?.plan;

  return (
    <div className="space-y-8">
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <CreditCard className="h-5 w-5" />
        </span>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">Billing</h1>
          <p className="text-sm text-text-secondary">Your plan and what you have been charged</p>
        </div>
      </div>

      <section className="rounded-xl border border-border bg-background p-6">
        {subscription.isLoading ? (
          <p className="text-sm text-text-secondary">Loading your plan...</p>
        ) : !current ? (
          <div className="space-y-1">
            <p className="text-sm font-medium text-foreground">No subscription yet</p>
            <p className="text-sm text-text-secondary">
              Once you are on a plan it shows here, with the date it renews.
            </p>
          </div>
        ) : (
          <div className="space-y-5">
            <div className="flex flex-wrap items-start justify-between gap-4">
              <div>
                <div className="flex items-center gap-3">
                  <h2 className="text-lg font-semibold text-foreground">{plan?.name ?? "Your plan"}</h2>
                  <StatusPill status={current.status} />
                </div>
                {plan?.blurb && <p className="mt-1 text-sm text-text-secondary">{plan.blurb}</p>}
              </div>
              {plan && (
                <p className="text-right text-lg font-semibold text-foreground tabular-nums">
                  {formatMoney(plan.price)}
                  <span className="ml-1 text-sm font-normal text-text-secondary">
                    / {plan.interval === "yearly" ? "year" : "month"}
                  </span>
                </p>
              )}
            </div>

            <dl className="grid gap-4 border-t border-border pt-5 sm:grid-cols-3">
              <div>
                <dt className="text-xs font-medium uppercase tracking-wider text-text-muted">Seats</dt>
                <dd className="mt-1 text-sm text-foreground tabular-nums">
                  {current.seats}
                  {plan?.seats ? <span className="text-text-secondary"> of {plan.seats}</span> : null}
                </dd>
              </div>
              <div>
                <dt className="text-xs font-medium uppercase tracking-wider text-text-muted">
                  {current.status === "canceled" ? "Ended" : "Renews"}
                </dt>
                <dd className="mt-1 text-sm text-foreground">
                  {longDate(current.status === "canceled" ? current.canceled_on : current.current_period_end)}
                </dd>
              </div>
              <div>
                <dt className="text-xs font-medium uppercase tracking-wider text-text-muted">Started</dt>
                <dd className="mt-1 text-sm text-foreground">{longDate(current.created_at)}</dd>
              </div>
            </dl>
          </div>
        )}
      </section>

      <section className="space-y-4">
        <div className="flex items-center gap-2">
          <Receipt className="h-4 w-4 text-text-secondary" />
          <h2 className="text-sm font-semibold uppercase tracking-wider text-text-secondary">Invoices</h2>
        </div>

        {invoices.isLoading ? (
          <p className="text-sm text-text-secondary">Loading your invoices...</p>
        ) : invoices.data && invoices.data.length > 0 ? (
          <div className="overflow-hidden rounded-xl border border-border">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border bg-bg-secondary text-left">
                  <th className="px-4 py-3 font-medium text-text-secondary">Invoice</th>
                  <th className="px-4 py-3 font-medium text-text-secondary">Issued</th>
                  <th className="px-4 py-3 font-medium text-text-secondary">Status</th>
                  <th className="px-4 py-3 text-right font-medium text-text-secondary">Amount</th>
                </tr>
              </thead>
              <tbody>
                {invoices.data.map((invoice) => (
                  <tr key={invoice.id} className="border-b border-border last:border-0">
                    <td className="px-4 py-3 font-medium text-foreground">{invoice.number}</td>
                    <td className="px-4 py-3 text-text-secondary">{longDate(invoice.issued_on)}</td>
                    <td className="px-4 py-3">
                      <StatusPill status={invoice.status} />
                    </td>
                    <td className="px-4 py-3 text-right text-foreground tabular-nums">
                      {formatMoney(invoice.amount)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="rounded-xl border border-border bg-background p-6 text-sm text-text-secondary">
            Nothing billed yet. Invoices appear here as they are issued.
          </p>
        )}
      </section>
    </div>
  );
}
