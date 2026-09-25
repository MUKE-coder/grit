package scaffold

import (
	"fmt"
	"strings"
)

// The account screen: one page for who you are, how you sign in, and where you
// are signed in.
//
// These pieces existed and were scattered. The password form sat on a page
// called "profile" next to a job title and a bio; two-factor and passkeys sat
// on a second page; sessions on a third. Somebody trying to lock their account
// down had to know all three existed. This is one screen, reachable from the
// System hub under Security & Access, and the old routes send you here.

// writeAdminAccountFiles writes the account page and the two cards it owns.
func writeAdminAccountFiles(root string, opts Options) error {
	files := map[string]string{
		adminComponent(root, opts, "account", "password-form.tsx"): adminAccountPasswordFormTSX(),
		adminComponent(root, opts, "account", "profile-form.tsx"):  adminAccountProfileFormTSX(),
		adminComponent(root, opts, "account", "sign-in-links.tsx"): adminSignInLinksCardTSX(),
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	for path, content := range adminPageFiles(root, opts, "system/account", "system-account",
		strings.ReplaceAll(adminAccountPageTSX(), "{{MODULE}}", opts.Module())) {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// adminAccountPasswordFormTSX emits components/account/password-form.tsx.
func adminAccountPasswordFormTSX() string {
	return `"use client";

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
    <section className="rounded-xl border border-border bg-bg-elevated">
      <div className="grid gap-6 p-6 sm:grid-cols-[minmax(0,14rem)_minmax(0,1fr)]">
        <div>
          <h2 className="flex items-center gap-2.5 text-base font-semibold text-foreground">
            <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent/10 text-accent">
              <Lock className="h-4 w-4" aria-hidden="true" />
            </span>
            Password
          </h2>
          <p className="mt-2 text-sm leading-relaxed text-text-muted">
            Changing it signs you out everywhere else.
          </p>
        </div>

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
`
}

// adminAccountProfileFormTSX emits components/account/profile-form.tsx.
//
// Everything the old /profile page did, so that page can redirect here and
// nothing is lost: the avatar, the names, the job title and bio, the email
// (which asks for the password, because it changes where a reset goes), and
// closing the account.
func adminAccountProfileFormTSX() string {
	return `"use client";

import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { Loader2, Trash2, Upload, User } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";
import { useMe } from "@/hooks/use-auth";
import { useUpdateProfile } from "@/hooks/use-profile";
import { DeleteAccountDialog } from "@/components/profile/delete-account-dialog";
import { buttonClasses } from "@/components/ui/button";
import { uploadFile } from "@/lib/api-client";

interface Values {
  first_name: string;
  last_name: string;
  email: string;
  job_title: string;
  bio: string;
  current_password: string;
}

export function ProfileForm() {
  const { data: user } = useMe();
  const updateProfile = useUpdateProfile();
  const [uploading, setUploading] = useState(false);
  const [uploadError, setUploadError] = useState("");
  const [confirmDelete, setConfirmDelete] = useState(false);
  const avatarInput = useRef<HTMLInputElement>(null);

  const { register, handleSubmit, reset, watch } = useForm<Values>({
    defaultValues: { first_name: "", last_name: "", email: "", job_title: "", bio: "", current_password: "" },
  });

  useEffect(() => {
    if (!user) return;
    const extra = user as { job_title?: string; bio?: string };
    reset({
      first_name: user.first_name ?? "",
      last_name: user.last_name ?? "",
      email: user.email ?? "",
      job_title: extra.job_title ?? "",
      bio: extra.bio ?? "",
      current_password: "",
    });
  }, [user, reset]);

  const emailChanged = (watch("email") ?? "") !== (user?.email ?? "");

  async function onAvatarPicked(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (!file) return;
    setUploading(true);
    setUploadError("");
    try {
      const result = await uploadFile(file);
      const url = (result.data as Record<string, unknown>)?.url as string | undefined;
      if (url) {
        updateProfile.mutate({ avatar: url });
      } else {
        setUploadError("That upload came back without a URL. Try again.");
      }
    } catch {
      // Said out loud rather than swallowed: a picture that silently does not
      // change is a bug report nobody can describe.
      setUploadError("That picture did not upload. Try again.");
    } finally {
      setUploading(false);
    }
  }

  function onSubmit(values: Values) {
    const payload: Record<string, string> = {
      first_name: values.first_name,
      last_name: values.last_name,
      job_title: values.job_title,
      bio: values.bio,
    };
    if (emailChanged) {
      payload.email = values.email;
      payload.current_password = values.current_password;
    }
    updateProfile.mutate(payload);
  }

  const field = inputClasses();

  return (
    <div className="space-y-6">
      <section className="rounded-xl border border-border bg-bg-elevated">
        <div className="grid gap-6 p-6 sm:grid-cols-[minmax(0,14rem)_minmax(0,1fr)]">
          <div>
            <h2 className="flex items-center gap-2.5 text-base font-semibold text-foreground">
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent/10 text-accent">
                <User className="h-4 w-4" aria-hidden="true" />
              </span>
              Profile
            </h2>
            <p className="mt-2 text-sm leading-relaxed text-text-muted">
              Your name and picture are what other people in this app see.
            </p>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            <div className="flex items-center gap-4">
              <span className="inline-flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-full bg-bg-tertiary text-lg font-semibold text-text-secondary">
                {user?.avatar ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img src={user.avatar} alt="" className="h-full w-full object-cover" />
                ) : (
                  (user?.first_name?.[0] ?? "") + (user?.last_name?.[0] ?? "")
                )}
              </span>
              <div>
                <input
                  ref={avatarInput}
                  id="account-avatar"
                  type="file"
                  accept="image/*"
                  className="sr-only"
                  onChange={onAvatarPicked}
                />
                <button
                  type="button"
                  onClick={() => avatarInput.current?.click()}
                  disabled={uploading}
                  className={buttonClasses({ variant: "outline", size: "sm" })}
                >
                  {uploading ? (
                    <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
                  ) : (
                    <Upload className="h-4 w-4" aria-hidden="true" />
                  )}
                  Change picture
                </button>
                {uploadError && (
                  <p role="alert" className="mt-1.5 text-xs text-danger">
                    {uploadError}
                  </p>
                )}
              </div>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <label htmlFor="account-first-name" className="block text-sm font-medium text-text-secondary">
                  First name
                </label>
                <input id="account-first-name" className={field} {...register("first_name")} />
              </div>
              <div className="space-y-1.5">
                <label htmlFor="account-last-name" className="block text-sm font-medium text-text-secondary">
                  Last name
                </label>
                <input id="account-last-name" className={field} {...register("last_name")} />
              </div>
            </div>

            <div className="space-y-1.5">
              <label htmlFor="account-job-title" className="block text-sm font-medium text-text-secondary">
                Job title
              </label>
              <input id="account-job-title" className={field} {...register("job_title")} />
            </div>

            <div className="space-y-1.5">
              <label htmlFor="account-bio" className="block text-sm font-medium text-text-secondary">
                Bio
              </label>
              <textarea id="account-bio" rows={3} className={field} {...register("bio")} />
            </div>

            <div className="space-y-1.5">
              <label htmlFor="account-email" className="block text-sm font-medium text-text-secondary">
                Email
              </label>
              <input id="account-email" type="email" autoComplete="email" className={field} {...register("email")} />
            </div>

            {emailChanged && (
              <div className="space-y-1.5 rounded-lg border border-warning/40 bg-warning/5 p-4">
                <label htmlFor="account-email-password" className="block text-sm font-medium text-text-secondary">
                  Current password
                </label>
                <input
                  id="account-email-password"
                  type="password"
                  autoComplete="current-password"
                  className={field}
                  {...register("current_password")}
                />
                <p className="text-xs text-text-muted">
                  Changing your email changes where a password reset goes, so it asks who you are first.
                </p>
              </div>
            )}

            {updateProfile.isError && (
              <p role="alert" className="text-sm text-danger">
                {(updateProfile.error as { response?: { data?: { error?: { message?: string } } } } | null)
                  ?.response?.data?.error?.message ?? "That did not save. Try again."}
              </p>
            )}
            {updateProfile.isSuccess && (
              <p role="status" className="text-sm text-success">
                Saved.
              </p>
            )}

            <button type="submit" disabled={updateProfile.isPending} className={buttonClasses({ variant: "primary" })}>
              {updateProfile.isPending && <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />}
              Save changes
            </button>
          </form>
        </div>
      </section>

      <section className="rounded-xl border border-danger/40 bg-danger/5">
        <div className="grid gap-6 p-6 sm:grid-cols-[minmax(0,14rem)_minmax(0,1fr)]">
          <div>
            <h2 className="flex items-center gap-2.5 text-base font-semibold text-foreground">
              <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-danger/10 text-danger">
                <Trash2 className="h-4 w-4" aria-hidden="true" />
              </span>
              Close this account
            </h2>
            <p className="mt-2 text-sm leading-relaxed text-text-muted">
              Your account and the data attached to it. This cannot be undone.
            </p>
          </div>
          <div>
            <button
              type="button"
              onClick={() => setConfirmDelete(true)}
              className={buttonClasses({ variant: "danger" })}
            >
              Close my account
            </button>
          </div>
        </div>
      </section>

      <DeleteAccountDialog open={confirmDelete} onClose={() => setConfirmDelete(false)} />
    </div>
  );
}
`
}

// adminAccountPageTSX emits the page itself.
func adminAccountPageTSX() string {
	return `"use client";

import { Suspense } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { KeyRound, Monitor, ShieldCheck, User } from "@/lib/icons";
import { ProfileForm } from "@/components/account/profile-form";
import { PasswordForm } from "@/components/account/password-form";
import { TwoFactorCard } from "@/components/profile/two-factor-card";
import { PasskeysCard } from "@/components/security/passkeys";
import { SignInLinksCard } from "@/components/account/sign-in-links";
import { ActiveSessions } from "@/components/profile/active-sessions";

/**
 * The account screen.
 *
 * Deliberately not /system/security, which is the operator's threat dashboard:
 * blocked addresses, rate-limit pressure, recent attacks. That page is about
 * other people. This one is about you, and putting "change your password" next
 * to a list of intrusion attempts helps nobody.
 *
 * The tabs are links rather than an ARIA tablist on purpose. A link is
 * keyboard-operable, shareable, and survives a refresh, and it needs no
 * roving-focus code to be correct: the pattern that fails a keyboard user is
 * always the hand-rolled one.
 */
const TABS = [
  { id: "profile", label: "Profile", icon: User },
  { id: "password", label: "Password", icon: KeyRound },
  { id: "security", label: "Security", icon: ShieldCheck },
  { id: "devices", label: "Devices", icon: Monitor },
] as const;

type TabID = (typeof TABS)[number]["id"];

function AccountTabs() {
  const params = useSearchParams();
  const requested = params.get("tab");
  const active: TabID = (TABS.find((t) => t.id === requested)?.id ?? "profile") as TabID;

  return (
    <div className="mx-auto max-w-4xl space-y-6 p-6">
      <header>
        <h1 className="text-2xl font-bold tracking-tight text-foreground">Account</h1>
        <p className="mt-1 text-sm text-text-muted">
          Your profile, how you sign in, and where you are signed in.
        </p>
      </header>

      <nav
        aria-label="Account sections"
        className="-mb-px flex gap-1 overflow-x-auto border-b border-border"
      >
        {TABS.map((tab) => {
          const Icon = tab.icon;
          const on = tab.id === active;
          return (
            <Link
              key={tab.id}
              href={"/system/account?tab=" + tab.id}
              aria-current={on ? "page" : undefined}
              className={
                // Three signals, not one: the underline, the weight and the
                // colour. An underline alone is a two-pixel line somebody has
                // to go looking for, and it is the only signal a person with
                // low vision loses first.
                "inline-flex min-h-[44px] items-center gap-2 whitespace-nowrap border-b-2 px-4 text-sm transition-colors " +
                (on
                  ? "border-accent font-semibold text-accent"
                  : "border-transparent font-medium text-text-muted hover:border-border hover:text-foreground")
              }
            >
              <Icon className="h-4 w-4" aria-hidden="true" />
              {tab.label}
            </Link>
          );
        })}
      </nav>

      {active === "profile" && <ProfileForm />}

      {active === "password" && <PasswordForm />}

      {active === "security" && (
        <div className="space-y-6">
          {/* Two-factor first: it is the single largest improvement available
              on this page, and an account without it is one leaked password
              away from gone. */}
          <TwoFactorCard />
          {/* Then passkeys, which are the thing that makes the password matter
              less. The card hides itself where the browser has no
              authenticator. */}
          <PasskeysCard />
          {/* And what has been asked for by email. The request form answers the
              same to every address so that it cannot enumerate accounts, which
              also means it can never warn anybody: this is the only place a
              request somebody did not make becomes visible. */}
          <SignInLinksCard />
        </div>
      )}

      {active === "devices" && <ActiveSessions />}
    </div>
  );
}

export default function AccountPage() {
  return (
    <Suspense
      fallback={<div className="mx-auto max-w-4xl p-6 text-sm text-text-muted">One moment...</div>}
    >
      <AccountTabs />
    </Suspense>
  );
}
`
}
