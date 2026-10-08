"use client";

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
