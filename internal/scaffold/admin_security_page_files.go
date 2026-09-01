package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeAdminSecurityFiles writes the account security page: one screen for
// everything that protects a login.
//
// The pieces already existed and were scattered. Two-factor and sessions lived
// on a page called "profile", next to a bio and an avatar, which is a page
// somebody opens to change their job title. Security decisions deserve their
// own screen, and the moment there is more than one of them, they need
// somewhere to sit together.
func writeAdminSecurityFiles(root string, opts Options) error {
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{
		filepath.Join(adminRoot, "app", "(dashboard)", "account", "security", "page.tsx"): adminSecurityPageTSX(),
		filepath.Join(adminRoot, "components", "security", "recovery-contacts.tsx"):       adminRecoveryContactsTSX(),
		filepath.Join(adminRoot, "hooks", "use-security.ts"):                              adminUseSecurityTS(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", opts.Module())
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func adminUseSecurityTS() string {
	return `import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

/**
 * What the server says about this account and this deployment.
 *
 * Both halves matter. The account half is state; the deployment half is
 * capability, and it is why sms_provider_configured is here: phone recovery
 * only exists if somebody wired a provider, and the page leaves the card out
 * rather than rendering a control that cannot work.
 */
export interface SecurityOverview {
  email: string;
  email_verified: boolean;
  has_password: boolean;
  provider: string;
  /** Masked, e.g. "b****p@example.com". The server never returns it in full. */
  recovery_email: string;
  recovery_email_verified: boolean;
  recovery_phone: string;
  recovery_phone_verified: boolean;
  sms_provider_configured: boolean;
}

export function useSecurityOverview() {
  return useQuery<SecurityOverview>({
    queryKey: ["security-overview"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/auth/security");
      return data.data;
    },
  });
}

type Kind = "email" | "phone";

/** Starts adding a recovery contact. Sends a code; stores nothing yet. */
export function useSetRecoveryContact(kind: Kind) {
  return useMutation({
    mutationFn: async (input: { password: string; value: string }) => {
      const body =
        kind === "email"
          ? { password: input.password, email: input.value }
          : { password: input.password, phone: input.value };
      const { data } = await apiClient.post("/api/auth/recovery/" + kind, body);
      return data;
    },
  });
}

export function useVerifyRecoveryContact(kind: Kind) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (code: string) => {
      const { data } = await apiClient.post("/api/auth/recovery/" + kind + "/verify", { code });
      return data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["security-overview"] }),
  });
}

export function useClearRecoveryContact(kind: Kind) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (password: string) => {
      // A body on DELETE, because removing a recovery contact is as sensitive
      // as adding one and takes the same password.
      const { data } = await apiClient.delete("/api/auth/recovery/" + kind, { data: { password } });
      return data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["security-overview"] }),
  });
}
`
}

func adminRecoveryContactsTSX() string {
	return `"use client";

import { useState } from "react";
import { Mail, Phone, ShieldCheck, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import {
  useSetRecoveryContact,
  useVerifyRecoveryContact,
  useClearRecoveryContact,
  type SecurityOverview,
} from "@/hooks/use-security";

/**
 * One recovery contact: add, confirm with a code, or remove.
 *
 * The password field is not friction for its own sake. A recovery address is a
 * second way into the account, so somebody holding a live session on a borrowed
 * laptop could otherwise attach their own address and keep the account forever.
 * The password is the thing they do not have, and the server checks it on both
 * add and remove.
 */
export function RecoveryContactCard({
  kind,
  overview,
}: {
  kind: "email" | "phone";
  overview: SecurityOverview;
}) {
  const isEmail = kind === "email";
  const current = isEmail ? overview.recovery_email : overview.recovery_phone;
  const verified = isEmail ? overview.recovery_email_verified : overview.recovery_phone_verified;

  const [value, setValue] = useState("");
  const [password, setPassword] = useState("");
  const [code, setCode] = useState("");
  const [pending, setPending] = useState<string | null>(null);

  const set = useSetRecoveryContact(kind);
  const verify = useVerifyRecoveryContact(kind);
  const clear = useClearRecoveryContact(kind);

  const label = isEmail ? "Recovery email" : "Recovery phone";
  const Icon = isEmail ? Mail : Phone;

  function apiMessage(err: unknown, fallback: string) {
    const e = err as { response?: { data?: { error?: { message?: string } } } };
    return e?.response?.data?.error?.message || fallback;
  }

  async function onSend() {
    try {
      const res = await set.mutateAsync({ password, value });
      setPending(res?.data?.sent_to ?? value);
      setPassword("");
      toast.success("Code sent. Enter it below to confirm.");
    } catch (err) {
      toast.error(apiMessage(err, "Could not send the code"));
    }
  }

  async function onVerify() {
    try {
      await verify.mutateAsync(code);
      setPending(null);
      setValue("");
      setCode("");
      toast.success(label + " confirmed");
    } catch (err) {
      toast.error(apiMessage(err, "That code is not valid"));
    }
  }

  async function onRemove() {
    const pw = window.prompt("Enter your password to remove this recovery contact");
    if (!pw) return;
    try {
      await clear.mutateAsync(pw);
      toast.success(label + " removed");
    } catch (err) {
      toast.error(apiMessage(err, "Could not remove it"));
    }
  }

  return (
    <section className="rounded-xl border border-border bg-bg-elevated p-6">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3">
          <span className="mt-0.5 flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-accent">
            <Icon className="h-4 w-4" aria-hidden="true" />
          </span>
          <div>
            <h2 className="text-sm font-semibold text-foreground">{label}</h2>
            <p className="mt-0.5 text-xs text-text-muted">
              {isEmail
                ? "Where we can reach you if you lose access to your sign-in address."
                : "A number we can text if you lose access to your email."}
            </p>
          </div>
        </div>

        {verified && current && (
          <span className="inline-flex shrink-0 items-center gap-1.5 rounded-full bg-success/10 px-2.5 py-1 text-xs font-medium text-success">
            <ShieldCheck className="h-3.5 w-3.5" aria-hidden="true" />
            Confirmed
          </span>
        )}
      </div>

      <div className="mt-5">
        {verified && current ? (
          <div className="flex items-center justify-between gap-3 rounded-lg border border-border bg-bg-secondary px-3 py-2.5">
            {/* Masked by the server, not here. Whoever is reading this screen
                might be the problem, and the full address tells them where to
                go next. */}
            <span className="font-mono text-sm text-foreground">{current}</span>
            <Button variant="ghost" size="sm" onClick={onRemove} disabled={clear.isPending}>
              <Trash2 className="mr-1.5 h-3.5 w-3.5" aria-hidden="true" />
              Remove
            </Button>
          </div>
        ) : pending ? (
          <div className="space-y-3">
            <p className="text-sm text-text-secondary">
              We sent a six-digit code to <span className="font-mono">{pending}</span>. It
              expires in 15 minutes.
            </p>
            <div className="flex gap-2">
              <Input
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder="123456"
                inputMode="numeric"
                maxLength={6}
                className="w-32 font-mono"
                aria-label="Verification code"
              />
              <Button onClick={onVerify} disabled={code.length < 6 || verify.isPending}>
                {verify.isPending ? "Confirming..." : "Confirm"}
              </Button>
              <Button variant="ghost" onClick={() => setPending(null)}>
                Cancel
              </Button>
            </div>
          </div>
        ) : (
          <div className="grid gap-3 sm:grid-cols-2">
            <div>
              <label className="mb-1 block text-xs font-medium text-text-muted">
                {isEmail ? "Recovery address" : "Phone number"}
              </label>
              <Input
                value={value}
                onChange={(e) => setValue(e.target.value)}
                placeholder={isEmail ? "you@example.com" : "+256700000000"}
                type={isEmail ? "email" : "tel"}
              />
            </div>
            <div>
              <label className="mb-1 block text-xs font-medium text-text-muted">
                Your password
              </label>
              <Input
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                type="password"
                placeholder="Confirm it is you"
                autoComplete="current-password"
              />
            </div>
            <div className="sm:col-span-2">
              <Button onClick={onSend} disabled={!value || !password || set.isPending}>
                {set.isPending ? "Sending..." : "Send code"}
              </Button>
            </div>
          </div>
        )}
      </div>
    </section>
  );
}
`
}

func adminSecurityPageTSX() string {
	return `"use client";

import Link from "next/link";
import { AlertCircle, ArrowLeft, KeyRound, Monitor, ShieldCheck } from "lucide-react";
import { TwoFactorCard } from "@/components/profile/two-factor-card";
import { ActiveSessions } from "@/components/profile/active-sessions";
import { RecoveryContactCard } from "@/components/security/recovery-contacts";
import { PasskeysCard } from "@/components/security/passkeys";
import { useSecurityOverview } from "@/hooks/use-security";
import { SkeletonCards } from "@/components/ui/Skeleton";

/**
 * Everything that protects this account, on one screen.
 *
 * Deliberately not /system/security, which is the operator's threat dashboard:
 * blocked addresses, recent attacks, the health of the perimeter. That page is
 * about other people. This one is about you, and merging them would put a
 * "change your password" box next to a list of intrusion attempts.
 */
export default function AccountSecurityPage() {
  const { data: overview, isLoading } = useSecurityOverview();

  if (isLoading || !overview) {
    return (
      <div className="p-6">
        <SkeletonCards count={4} />
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-3xl space-y-6 p-6">
      <div>
        <Link
          href="/profile"
          className="mb-2 inline-flex items-center gap-1.5 text-xs text-text-muted transition-colors hover:text-foreground"
        >
          <ArrowLeft className="h-3.5 w-3.5" aria-hidden="true" />
          Back to profile
        </Link>
        <h1 className="text-2xl font-bold tracking-tight text-foreground">Security</h1>
        <p className="mt-1 text-sm text-text-muted">
          How you sign in, and how you get back in if something goes wrong.
        </p>
      </div>

      {/* Ordered by what actually protects the account. Two-factor first
          because it is the single largest improvement available here, and an
          account without it is one leaked password away from gone. */}
      <TwoFactorCard />

      {/* Passkeys before recovery, because a passkey is the thing that makes
          the password matter less, and recovery is what you need when it does
          not. The card hides itself when the browser has no authenticator. */}
      <PasskeysCard />

      <RecoveryContactCard kind="email" overview={overview} />

      {/* Only when the deployment can actually send a text. A disabled control
          with no explanation is worse than no control. */}
      {overview.sms_provider_configured && (
        <RecoveryContactCard kind="phone" overview={overview} />
      )}

      {!overview.sms_provider_configured && (
        <div className="flex gap-3 rounded-xl border border-border bg-bg-secondary/50 p-4">
          <AlertCircle className="mt-0.5 h-4 w-4 shrink-0 text-text-muted" aria-hidden="true" />
          <p className="text-xs leading-relaxed text-text-muted">
            <span className="font-medium text-text-secondary">Phone recovery is not available.</span>{" "}
            This deployment has no SMS provider configured. Register one in{" "}
            <code className="rounded bg-bg-hover px-1 py-0.5 font-mono">internal/sms</code> and this
            option appears.
          </p>
        </div>
      )}

      {!overview.has_password && (
        <div className="flex gap-3 rounded-xl border border-warning/30 bg-warning/5 p-4">
          <KeyRound className="mt-0.5 h-4 w-4 shrink-0 text-warning" aria-hidden="true" />
          <p className="text-xs leading-relaxed text-text-secondary">
            This account signs in with {overview.provider}. Set a password on your profile before
            adding recovery contacts, since confirming them requires one.
          </p>
        </div>
      )}

      {/* ActiveSessions renders the list and nothing else: on the profile page
          it sits under that page's own heading. Given one here so it reads as
          a card like the rest of this screen rather than a loose list. */}
      <section className="rounded-xl border border-border bg-bg-elevated p-6">
        <div className="mb-4 flex items-start gap-3">
          <span className="mt-0.5 flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-accent">
            <Monitor className="h-4 w-4" aria-hidden="true" />
          </span>
          <div>
            <h2 className="text-sm font-semibold text-foreground">Active sessions</h2>
            <p className="mt-0.5 text-xs text-text-muted">
              Everywhere you are currently signed in. Sign out anything you do not recognise.
            </p>
          </div>
        </div>
        <ActiveSessions />
      </section>

      <div className="flex gap-3 rounded-xl border border-border bg-bg-secondary/50 p-4">
        <ShieldCheck className="mt-0.5 h-4 w-4 shrink-0 text-text-muted" aria-hidden="true" />
        <p className="text-xs leading-relaxed text-text-muted">
          Change your password on the{" "}
          <Link href="/profile" className="text-accent hover:underline">
            profile page
          </Link>
          . Changing it signs out every other device.
        </p>
      </div>
    </div>
  );
}
`
}
