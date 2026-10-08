"use client";

import { useEffect } from "react";
import { buttonClasses } from "@/components/ui/button";

export default function ErrorPage({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="w-full max-w-md text-center">
        <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-full bg-danger/10 border border-danger/20">
          <svg className="h-8 w-8 text-danger" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <h2 className="mb-2 text-2xl font-bold text-foreground">Something went wrong</h2>
        <p className="mb-6 text-text-secondary">
          An unexpected error occurred. You can try again or go back.
        </p>
        {error.digest && (
          <p className="mb-4 text-xs text-text-muted font-mono">Error ID: {error.digest}</p>
        )}
        <div className="flex gap-3 justify-center">
          <button
            onClick={() => window.history.back()}
            className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-4 py-2.5 text-sm font-medium text-foreground hover:bg-bg-tertiary transition-colors"
          >
            Go Back
          </button>
          <button
            onClick={reset}
            className={buttonClasses()}
          >
            Try Again
          </button>
        </div>
      </div>
    </div>
  );
}
