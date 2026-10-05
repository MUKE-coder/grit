"use client";

import { useState } from "react";
import { useForm } from "react-hook-form";
import { Check, Eye, EyeOff, Loader2, Lock } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";
import { useMe } from "@/hooks/use-auth";
import { useChangePassword } from "@/hooks/use-profile";
import { buttonClasses } from "@/components/ui/button";

/**
 * Changing a password, with the rules shown while you type.
 *
 * The checklist is the same four rules the API enforces (internal/password),
 * evaluated here so the answer is immediate. The server is still the
 * authority: this only decides what the list looks like, never whether the
 * save is allowed, so a checklist that drifted would cost a refused save and
 * not a weak password.
 */
const RULES = [
  { id: "length", label: "At least 8 characters" },
  { id: "variety", label: "Letters and something else: a number, a symbol or a space" },
  { id: "not-common", label: "Not a password everyone tries first" },
  { id: "not-personal", label: "Nothing from your name or email" },
] as const;

// The few hundred that attackers try first, in the order they try them. The
// long lists belong in a service; most of the value is in the first page.
const COMMON = new Set([
  "123456", "password", "12345678", "qwerty", "123456789", "12345", "1234", "111111",
  "1234567", "dragon", "123123", "baseball", "abc123", "football", "monkey", "letmein",
  "shadow", "master", "666666", "qwertyuiop", "123321", "mustang", "1234567890",
  "superman", "1qaz2wsx", "7777777", "121212", "000000", "qazwsx", "123qwe", "killer",
  "trustno1", "zxcvbnm", "asdfgh", "iloveyou", "starwars", "112233", "computer",
  "zxcvbn", "555555", "11111111", "131313", "freedom", "777777", "pass", "159753",
  "aaaaaa", "princess", "welcome", "admin", "letmein1", "password1", "password123",
  "passw0rd", "p@ssw0rd", "qwerty123", "iloveyou1", "welcome1", "admin123", "root",
  "changeme", "secret123", "hello", "test", "1111", "0000", "sunshine", "whatever",
]);

function personalPieces(raw: string): string[] {
  const lower = raw.toLowerCase().trim();
  if (!lower) return [];
  const at = lower.indexOf("@");
  const source = at > 0 ? lower.slice(0, at) + " " + lower.slice(at + 1).split(".")[0] : lower;
  return source.split(/[^a-z0-9]+/).filter(Boolean);
}

export function passwordFailures(candidate: string, about: string[]): string[] {
  const failed: string[] = [];
  if ([...candidate].length < 8) failed.push("length");

  const letters = /\p{L}/u.test(candidate);
  const others = [...candidate].some((c) => !/\p{L}/u.test(c));
  if (!letters || !others) failed.push("variety");

  if (COMMON.has(candidate.toLowerCase().trim())) failed.push("not-common");

  const lower = candidate.toLowerCase();
  const personal = about
    .flatMap(personalPieces)
    .some((piece) => piece.length >= 4 && lower.includes(piece));
  if (personal) failed.push("not-personal");

  return failed;
}

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
  const met = RULES.filter((rule) => !failed.includes(rule.id));
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

          {/* Four segments, one per rule met. A bar with no explanation tells
              somebody they are wrong without telling them how to be right, so
              the list below it is the part that matters. */}
          <div className="flex gap-1.5" aria-hidden="true">
            {RULES.map((rule, i) => (
              <span
                key={rule.id}
                className={
                  "h-1 flex-1 rounded-full transition-colors " +
                  (next.length > 0 && met.length > i ? "bg-success" : "bg-border")
                }
              />
            ))}
          </div>

          <ul id="account-password-rules" className="space-y-1.5 text-sm">
            {RULES.map((rule) => {
              const ok = next.length > 0 && !failed.includes(rule.id);
              return (
                <li key={rule.id} className="flex items-start gap-2">
                  <span
                    className={
                      "mt-0.5 inline-flex h-4 w-4 shrink-0 items-center justify-center rounded-full border " +
                      (ok ? "border-success bg-success/15 text-success" : "border-border text-transparent")
                    }
                  >
                    <Check className="h-3 w-3" aria-hidden="true" />
                  </span>
                  <span className={ok ? "text-text-secondary" : "text-text-muted"}>
                    {rule.label}
                    <span className="sr-only">{ok ? ": met" : ": not met yet"}</span>
                  </span>
                </li>
              );
            })}
          </ul>

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
