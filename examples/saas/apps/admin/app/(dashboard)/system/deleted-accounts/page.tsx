"use client";

import { useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient } from "@/lib/api-client";
import { getApiErrorMessage } from "@/lib/api-core";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { AlertTriangle, RotateCcw, Trash2 } from "@/lib/icons";

/*
 * Accounts somebody closed.
 *
 * Closing an account is a soft delete, so the row stayed on disk and left every
 * list at once: "I deleted my account by mistake" had no answer short of a SQL
 * console, and neither did "who closed their account last month".
 *
 * Deliberately not the same page as the bin. The bin lists rows of the
 * resources an app defines; this is the people who can sign in to it, which is
 * a different question asked by a different person at a different time. They
 * also come from different places: users are kept out of the sync registry the
 * bin reads, because syncing a row that carries its own role is how an account
 * makes itself an administrator.
 */

interface DeletedUser {
  id: string;
  first_name?: string;
  last_name?: string;
  email: string;
  role?: string;
  deleted_at?: string;
}

export default function DeletedAccountsPage() {
  const queryClient = useQueryClient();
  const [confirm, setConfirm] = useState<DeletedUser | null>(null);
  const [error, setError] = useState<string | null>(null);

  const { data, isPending } = useQuery<{ data: DeletedUser[] }>({
    queryKey: ["admin", "deleted-accounts"],
    queryFn: async () => (await apiClient.get<{ data: DeletedUser[] }>("/api/admin/deleted-accounts")).data,
  });

  function refresh() {
    void queryClient.invalidateQueries({ queryKey: ["admin", "deleted-accounts"] });
    // The account is going back into the Users list, which is now stale.
    void queryClient.invalidateQueries({ queryKey: ["resource"] });
  }

  const restore = useMutation({
    mutationFn: (user: DeletedUser) => apiClient.post("/api/admin/deleted-accounts/" + user.id + "/restore"),
    onSuccess: refresh,
    onError: (e) => setError(getApiErrorMessage(e, "Could not restore that account")),
  });

  const purge = useMutation({
    mutationFn: (user: DeletedUser) => apiClient.delete("/api/admin/deleted-accounts/" + user.id),
    onSuccess: refresh,
    onError: (e) => setError(getApiErrorMessage(e, "Could not delete that account")),
  });

  const users = data?.data ?? [];

  return (
    <div>
      <PageHeader
        title="Deleted accounts"
        subtitle="Accounts that were closed. Restore one, or remove it for good."
      />

      {error && (
        <div role="alert" className="mb-4 flex items-start gap-2 rounded-lg border border-danger/30 bg-danger/5 px-4 py-3 text-sm text-danger">
          <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" aria-hidden="true" />
          <span>{error}</span>
        </div>
      )}

      <div className="rounded-xl border border-border bg-bg-elevated">
        {isPending ? (
          <p className="py-16 text-center text-sm text-text-muted">Loading...</p>
        ) : users.length === 0 ? (
          <p className="py-16 text-center text-sm text-text-muted">
            No account has been closed. One that is closed lands here rather than disappearing.
          </p>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border text-left text-xs text-text-muted">
                <th className="px-5 py-3 font-medium">Who</th>
                <th className="px-5 py-3 font-medium">Role</th>
                <th className="px-5 py-3 font-medium">Closed</th>
                <th className="px-5 py-3 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {users.map((user) => (
                <tr key={user.id} className="border-b border-border/60 last:border-0">
                  <td className="px-5 py-3">
                    <span className="font-medium text-foreground">
                      {[user.first_name, user.last_name].filter(Boolean).join(" ") || "No name"}
                    </span>
                    <span className="ml-2 text-xs text-text-muted">{user.email}</span>
                  </td>
                  <td className="px-5 py-3 text-text-secondary">{user.role ?? "—"}</td>
                  <td className="px-5 py-3 text-text-secondary">{when(user.deleted_at)}</td>
                  <td className="px-5 py-3">
                    <div className="flex items-center justify-end gap-1">
                      <button
                        type="button"
                        onClick={() => restore.mutate(user)}
                        disabled={restore.isPending}
                        className="inline-flex items-center gap-1.5 rounded-lg border border-border px-2.5 py-1 text-xs text-text-secondary transition-colors hover:bg-bg-hover hover:text-foreground"
                      >
                        <RotateCcw className="h-3 w-3" aria-hidden="true" />
                        Restore
                      </button>
                      <button
                        type="button"
                        onClick={() => setConfirm(user)}
                        className="inline-flex items-center gap-1.5 rounded-lg border border-transparent px-2.5 py-1 text-xs text-danger transition-colors hover:bg-danger/10"
                      >
                        <Trash2 className="h-3 w-3" aria-hidden="true" />
                        Delete forever
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <p className="mt-4 text-xs leading-relaxed text-text-muted">
        Restoring an account lets it sign in again with the password it had. Deleting it for good removes
        the row and nothing else: erasing what a person left behind across every table is a different
        operation, with different law attached to it, and it lives on the GDPR page.
      </p>

      <ConfirmModal
        open={confirm !== null}
        title="Delete this account for good?"
        description={
          (confirm?.email ?? "") +
          " cannot be restored after this. It removes the account row; use the GDPR page to erase what they left behind elsewhere."
        }
        confirmLabel="Delete forever"
        variant="danger"
        onConfirm={() => {
          if (confirm) purge.mutate(confirm);
          setConfirm(null);
        }}
        onCancel={() => setConfirm(null)}
      />
    </div>
  );
}

function when(iso?: string): string {
  if (!iso) return "—";
  const d = new Date(iso);
  return Number.isNaN(d.getTime())
    ? "—"
    : d.toLocaleDateString(undefined, { day: "numeric", month: "short", year: "numeric" });
}