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
	// No admin panel, nothing to write. The same guard writeAdminSecurityFiles
	// carries, for the same reason: every file below is an admin screen or a
	// component one of them renders.
	//
	// Without it, `grit new x --api` grew an apps/admin/ holding seven files and
	// nothing to build them with: no package.json, no layout, no tsconfig. The
	// home page offers that command as "the Go API alone", and it produced a
	// frontend directory that never compiles, which reads as a scaffold that
	// stopped halfway.
	//
	// The account, trash and deleted-account API endpoints are unaffected and
	// still exist in every architecture.
	if !opts.HasAdminPanel() {
		return nil
	}

	files := map[string]string{
		adminComponent(root, opts, "account", "password-form.tsx"):  adminAccountPasswordFormTSX(),
		adminComponent(root, opts, "account", "avatar-cropper.tsx"): adminAvatarCropperTSX(),
		adminComponent(root, opts, "account", "sign-in-links.tsx"):  adminSignInLinksCardTSX(),
	}
	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	// The page at /account, because the dashboard layout confines a plain USER
	// to /profile and /account: under /system it bounced exactly the people
	// whose account it is.
	for path, content := range adminPageFiles(root, opts, "account", "account",
		strings.ReplaceAll(adminAccountPageTSX(), "{{MODULE}}", opts.Module())) {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	// Closed accounts, beside the bin and separate from it.
	for path, content := range adminPageFiles(root, opts, "system/deleted-accounts", "system-deleted-accounts", adminDeletedAccountsPage()) {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	// The bin, which is a system page like the others.
	for path, content := range adminPageFiles(root, opts, "system/trash", "system-trash", adminTrashPage()) {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	// And a redirect where it used to be, for the System Hub card and anything
	// else that still points at the old path.
	for path, content := range adminPageFiles(root, opts, "system/account", "system-account",
		adminSystemAccountRedirect()) {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

// adminAccountPageTSX emits the page itself.
// adminAccountPageTSX is the merged account page, served at /account.
func adminAccountPageTSX() string { return tmpl("admin/app/account/account-page.tsx") }

// adminSystemAccountRedirect keeps /system/account working for anything
// that linked to it, including the System Hub card.
func adminSystemAccountRedirect() string {
	return `"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

// Merged into /account. Kept as a redirect rather than deleted: this path was
// the user menu's target for years, it is where a USER lands after signing in,
// and it is in release notes and bookmarks.
export default function SystemAccountRedirect() {
	const router = useRouter();
	useEffect(() => {
		// With the fragment, or the dashboard's "Set it up" lands at the top of a
		// long page instead of the two-factor card it was pointing at. A redirect
		// that drops the anchor is a redirect that half worked.
		router.replace("/account" + window.location.hash);
	}, [router]);

	return (
		<p className="p-8 text-sm text-text-secondary">
			Moved to <a className="underline" href="/account">your account</a>.
		</p>
	);
}
`
}

func adminAvatarCropperTSX() string { return tmpl("admin/components/account/avatar-cropper.tsx") }

// adminAccountPasswordFormTSX emits components/account/password-form.tsx.
func adminAccountPasswordFormTSX() string {
	return `"use client";

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
`
}
