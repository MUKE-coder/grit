"use client";

import { useMutation, type UseMutationOptions } from "@tanstack/react-query";
import { toast } from "sonner";
import { getApiErrorMessage } from "@/lib/api-core";

interface ToastedOptions<TData, TError, TVariables, TContext>
  extends UseMutationOptions<TData, TError, TVariables, TContext> {
  /** Success message. When a function, receives the mutation result. */
  successMessage?: string | ((data: TData) => string);
  /** Error message. When a function, receives the thrown error. */
  errorMessage?: string | ((err: TError) => string);
  /** Skip the success toast (e.g. when navigating away on success). */
  silentSuccess?: boolean;
}

/**
 * Drop-in replacement for useMutation. Times every mutation and emits a
 * toast on settle:
 *
 *   ✓ Ticket opened — 142ms
 *   ✗ Couldn't reach the API — 1203ms
 *
 * The timing is intentionally surfaced. Users get tactile feedback that
 * the request ran (vs hung) and developers spot regressions early.
 */
export function useToastedMutation<TData = unknown, TError = unknown, TVariables = void, TContext = unknown>(
  options: ToastedOptions<TData, TError, TVariables, TContext>
) {
  const { successMessage, errorMessage, silentSuccess, onSuccess, onError, ...rest } = options;

  // Use spread + cast to stay agnostic between react-query v4 (3-arg
  // callbacks) and v5 (4-arg callbacks with onMutateResult + context).
  // The wrapper only cares about data, vars, and the per-mutation context
  // where we stamped __startedAt — everything else passes through.
  return useMutation<TData, TError, TVariables, TContext>({
    ...rest,
    onSuccess: ((...args: unknown[]) => {
      const [data, , ctx] = args as [TData, TVariables, TContext];
      const elapsed = readElapsed(ctx as MutationContext | undefined);
      if (!silentSuccess) {
        const msg = typeof successMessage === "function" ? successMessage(data) : (successMessage || "Done");
        toast.success(msg + (elapsed != null ? " — " + elapsed + "ms" : ""));
      }
      passThrough(onSuccess, args);
    }) as never,
    onError: ((...args: unknown[]) => {
      const [err, , ctx] = args as [TError, TVariables, TContext];
      const elapsed = readElapsed(ctx as MutationContext | undefined);
      const fallback = pickErrorMessage(err);
      const msg = typeof errorMessage === "function" ? errorMessage(err) : (errorMessage || fallback);
      toast.error(msg + (elapsed != null ? " — " + elapsed + "ms" : ""));
      passThrough(onError, args);
    }) as never,
    onMutate: ((...args: unknown[]) => {
      const userCtxPromise = passThrough(rest.onMutate, args);
      return Promise.resolve(userCtxPromise).then((userCtx) =>
        ({ ...(userCtx as object), __startedAt: performance.now() } as TContext)
      );
    }) as never,
  });
}

type MutationContext = { __startedAt?: number } | undefined;

// Calls the caller's own callback with every argument react-query passed,
// however many that is in the installed version.
function passThrough(callback: unknown, args: unknown[]): unknown {
  return typeof callback === "function" ? (callback as (...a: unknown[]) => unknown)(...args) : undefined;
}

function readElapsed(ctx: MutationContext): number | null {
  const t = ctx?.__startedAt;
  if (typeof t !== "number") return null;
  return Math.round(performance.now() - t);
}

// The message this hook shows when the caller gave no errorMessage. It was a
// private copy of the same cast 21 other files had written out; now there is
// one, in lib/api-core.ts, and this names the default it wants.
function pickErrorMessage(err: unknown): string {
  return getApiErrorMessage(err, "Something went wrong");
}

// Re-export sonner's toast so pages can drop ad-hoc toasts (e.g. on a
// copy-to-clipboard) without a second import line.
export { toast };
