"use client";

import { useState } from "react";
import { useDeleteAccount } from "@/hooks/use-profile";
import { AlertTriangle, Loader2, X } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";

interface DeleteAccountDialogProps {
  open: boolean;
  onClose: () => void;
}

export function DeleteAccountDialog({ open, onClose }: DeleteAccountDialogProps) {
  const [confirmation, setConfirmation] = useState("");
  const deleteAccount = useDeleteAccount();

  if (!open) return null;

  const handleDelete = () => {
    if (confirmation !== "DELETE") return;
    deleteAccount.mutate();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-md rounded-xl border border-border bg-bg-secondary p-6 shadow-2xl">
        <button
          onClick={onClose}
          className="absolute right-4 top-4 text-text-muted hover:text-foreground transition-colors"
        >
          <X className="h-4 w-4" />
        </button>

        <div className="flex items-center gap-3 text-danger">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-danger/10">
            <AlertTriangle className="h-5 w-5" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground">Delete Account</h3>
            <p className="text-xs text-text-muted">This action is permanent</p>
          </div>
        </div>

        <div className="mt-4 space-y-3">
          <p className="text-sm text-text-secondary">
            This will permanently delete your account, including all your data, settings, and
            access. This action <strong className="text-foreground">cannot be undone</strong>.
          </p>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-text-secondary">
              Type <span className="font-mono text-danger">DELETE</span> to confirm
            </label>
            <input
              type="text"
              value={confirmation}
              onChange={(e) => setConfirmation(e.target.value)}
              placeholder="DELETE"
              className={inputClasses({ invalid: true })}
            />
          </div>
        </div>

        <div className="mt-6 flex justify-end gap-3">
          <button
            onClick={onClose}
            className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleDelete}
            disabled={confirmation !== "DELETE" || deleteAccount.isPending}
            className="flex items-center gap-2 rounded-lg bg-danger px-4 py-2 text-sm font-medium text-white hover:bg-danger/90 disabled:opacity-50 transition-colors"
          >
            {deleteAccount.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <AlertTriangle className="h-4 w-4" />
            )}
            Delete my account
          </button>
        </div>
      </div>
    </div>
  );
}
