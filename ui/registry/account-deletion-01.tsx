"use client";

// account-deletion-01 — Account deletion confirmation (GitHub / Linear-inspired)
import { useState } from "react";
export function AccountDeletion({
  workspaceName = "Acme Corp",
  onDelete,
}: {
  workspaceName?: string; onDelete?: () => void;
}) {
  const [input, setInput] = useState("");
  const confirmed = input === workspaceName;
  return (
    <div className="border border-danger/40 bg-danger/5 rounded-xl p-5 flex flex-col gap-4">
      <div className="flex items-start gap-3">
        <svg className="w-5 h-5 text-danger flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <div>
          <h3 className="text-sm font-semibold text-danger">Delete workspace</h3>
          <p className="text-xs text-text-secondary mt-0.5">This action cannot be undone. All projects, data, and team members will be permanently removed.</p>
        </div>
      </div>
      <div className="flex flex-col gap-1.5">
        <label className="text-xs text-text-secondary">
          Type <strong className="text-foreground font-mono">{workspaceName}</strong> to confirm
        </label>
        <input
          value={input} onChange={(e) => setInput(e.target.value)}
          placeholder={workspaceName}
          className="bg-background border border-border rounded-lg px-3 py-2 text-sm text-foreground outline-none focus:border-danger transition-colors"
        />
      </div>
      <button
        onClick={onDelete} disabled={!confirmed}
        className="w-full py-2.5 rounded-lg bg-danger hover:bg-danger/90 text-white font-semibold text-sm transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
      >
        Delete {workspaceName}
      </button>
    </div>
  );
}

export default AccountDeletion;
