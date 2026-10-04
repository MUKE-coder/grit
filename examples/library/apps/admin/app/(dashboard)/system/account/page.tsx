"use client";

import { Monitor } from "@/lib/icons";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ProfileForm, CloseAccountCard } from "@/components/account/profile-form";
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
 * One column rather than tabs. Tabs hid six cards behind four labels, so
 * answering "where am I signed in" meant knowing that devices were under
 * Devices and not under Security, and the page this replaced showed all of it
 * at once. Six cards is a scroll, not a navigation problem.
 *
 * PageHeader rather than a hand-written <h1>. It derives the back link for
 * every /system/* route, and carries the refresh, theme and notification
 * controls; writing the heading by hand is how this page ended up the only one
 * in the admin with no way back to the hub.
 */
export default function AccountPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Account"
        subtitle="Your profile, how you sign in, and where you are signed in."
      />

      <div className="mx-auto max-w-4xl space-y-6 px-6 pb-10">
        <ProfileForm />

        <PasswordForm />

        {/* Everything that used to be the Security tab. The id is what every
            existing deep link points at, and scroll-mt keeps the heading clear
            of the sticky header when one lands here. */}
        <section id="security" className="scroll-mt-24 space-y-6">
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
        </section>

        <section
          id="devices"
          className="scroll-mt-24 overflow-hidden rounded-xl border border-border bg-bg-elevated"
        >
          <div className="flex items-start gap-3 border-b border-border px-6 py-4">
            <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent/10 text-accent">
              <Monitor className="h-4 w-4" aria-hidden="true" />
            </span>
            <div className="min-w-0">
              <h2 className="text-base font-semibold text-foreground">Active sessions</h2>
              <p className="mt-1 text-sm leading-relaxed text-text-muted">
                Every device signed in to this account. Changing your password signs the
                others out.
              </p>
            </div>
          </div>

          <div className="p-6">
            <ActiveSessions />
          </div>
        </section>

        {/* Last, because it is the only thing here you cannot undo. */}
        <CloseAccountCard />
      </div>
    </div>
  );
}
