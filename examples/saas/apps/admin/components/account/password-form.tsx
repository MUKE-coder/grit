"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { Eye, EyeOff, Loader2, Lock } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";
import { passwordFailures } from "@/lib/password-rules";
import { PasswordStrength } from "@/components/password-strength";
import { useMe } from "@/hooks/use-auth";
import { useChangePassword } from "@/hooks/use-profile";
import { buttonClasses } from "@/components/ui/button";

/**
 * Changing a password, with the strength meter and the rules shown while you
 * type.
 *
 * Both come from lib/password-rules and components/password-strength, which the
 * register and reset screens use as well. This page used to carry its own copy
 * of the rules, and the copy checked four of the five: "At most 72 characters"
 * was in its list and in none of its code, so a password bcrypt cannot hash got
 * a green check here and a 422 from the API.
 *
 * The server is still the authority. The meter decides what the person sees,
 * never whether the save is allowed, so the worst a drift can cost is a
 * refusal that was not predicted rather than a password that should not exist.
 */
interface Values {
  current_password: string;
  password: string;
}

export function PasswordForm() {
  const { data: user } = useMe();
  const changePassword = useChangePassword();
  const [showCurrent, setShowCurrent] = useState(false);
  const [showNext, setShowNext] = useState(false);

  const { register, handleSubmit, watch, reset } = useForm<Values>({
    defaultValues: { current_password: "", password: "" },
  });

  const next = watch("password") ?? "";
  const current = watch("current_password") ?? "";
  const about = [user?.email ?? "", user?.first_name ?? "", user?.last_name ?? ""].filter(Boolean);
  const failed = passwordFailures(next, about);
  const ready = next.length > 0 && failed.length === 0 && current.length > 0;

  function onSubmit(values: Values) {
    changePassword.mutate(
      { current_password: values.current_password, password: values.password },
      { onSuccess: () => reset({ current_password: "", password: "" }) },
    );
  }

  return (
    <section className="overflow-hidden rounded-xl border border-border bg-bg-elevated">
      <div className="flex items-start gap-3 border-b border-border px-6 py-4">
        <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent/10 text-accent">
          <Lock className="h-4 w-4" aria-hidden="true" />
        </span>
        <div className="min-w-0">
          <h2 className="text-base font-semibold text-foreground">Password</h2>
          <p className="mt-1 text-sm leading-relaxed text-text-muted">
            Changing it signs you out everywhere else.
          </p>
        </div>
      </div>

      <div className="p-6">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
          <div className="space-y-1.5">
            <label htmlFor="account-current-password" className="block text-sm font-medium text-text-secondary">
              Current password
            </label>
            <div className="relative">
              <input
                id="account-current-password"
                type={showCurrent ? "text" : "password"}
                autoComplete="current-password"
                {...register("current_password")}
                className={inputClasses({ className: "pr-11" })}
              />
              <button
                type="button"
                onClick={() => setShowCurrent((on) => !on)}
                aria-label={showCurrent ? "Hide the password" : "Show the password"}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 text-text-muted transition-colors hover:text-foreground"
              >
                {showCurrent ? <EyeOff className="h-4 w-4" aria-hidden="true" /> : <Eye className="h-4 w-4" aria-hidden="true" />}
              </button>
            </div>
          </div>

          <div className="space-y-1.5">
            <label htmlFor="account-new-password" className="block text-sm font-medium text-text-secondary">
              New password
            </label>
            <div className="relative">
              <input
                id="account-new-password"
                type={showNext ? "text" : "password"}
                autoComplete="new-password"
                aria-describedby="account-password-rules"
                {...register("password")}
                className={inputClasses({ className: "pr-11" })}
              />
              <button
                type="button"
                onClick={() => setShowNext((on) => !on)}
                aria-label={showNext ? "Hide the password" : "Show the password"}
                className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 text-text-muted transition-colors hover:text-foreground"
              >
                {showNext ? <EyeOff className="h-4 w-4" aria-hidden="true" /> : <Eye className="h-4 w-4" aria-hidden="true" />}
              </button>
            </div>
          </div>

          <PasswordStrength value={next} about={about} id="account-password-rules" />

          {changePassword.isError && (
            <p role="alert" className="text-sm text-danger">
              {(changePassword.error as { response?: { data?: { error?: { message?: string } } } } | null)
                ?.response?.data?.error?.message ?? "That did not save. Try again."}
            </p>
          )}
          {changePassword.isSuccess && (
            <p role="status" className="text-sm text-success">
              Password changed. Every other device has been signed out.
            </p>
          )}

          <button
            type="submit"
            disabled={!ready || changePassword.isPending}
            className={buttonClasses({ variant: "primary" })}
          >
            {changePassword.isPending && <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />}
            Change password
          </button>
        </form>
      </div>
    </section>
  );
}
