"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Eye, EyeOff, Fingerprint } from "@/lib/icons";
import { useLogin, useMe, usePasskeyLogin, useVerifyTOTP } from "@/hooks/use-auth";
import { passkeysSupported } from "@/lib/webauthn";
import { apiClient } from "@/lib/api-client";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { LoginSchema, type LoginInput } from "@repo/shared/schemas";
import { AuthShell } from "@/components/auth/AuthShell";

const inputBase =
  "w-full rounded-[var(--auth-radius)] border bg-[var(--auth-card)] px-4 py-3 text-[var(--auth-fg)] placeholder:text-[var(--auth-muted)] focus:outline-none focus:ring-2 transition-colors";
const inputOk = inputBase + " border-[var(--auth-border)] focus:border-[var(--auth-primary)] focus:ring-[var(--auth-primary)]/30";
const inputErr = inputBase + " border-red-400 focus:border-red-500 focus:ring-red-400/30";

export default function LoginPage() {
  const [showPassword, setShowPassword] = useState(false);
  // Holding a pending token means the password was right and 2FA is owed.
  const [pendingToken, setPendingToken] = useState<string | null>(null);
  // "app" or "email": which second factor this account uses, as the
  // sign-in response reported it.
  const [factorMethod, setFactorMethod] = useState("app");
  const [code, setCode] = useState("");
  const [useBackup, setUseBackup] = useState(false);
  // Asked of the browser, not of the server. Before anybody has identified
  // themselves the server does not know whether this person has a passkey, and
  // asking it would be an oracle for which addresses have accounts. The
  // browser knows whether this device can produce one, which is the question
  // that decides whether the button can do anything.
  const [passkeyReady, setPasskeyReady] = useState(false);
  const [passkeyError, setPasskeyError] = useState("");
  const [trustDevice, setTrustDevice] = useState(false);
  // Signing in by link instead of by password.
  const [linkState, setLinkState] = useState<"idle" | "sending" | "sent">("idle");
  const [linkError, setLinkError] = useState("");
  const { mutate: login, isPending, error: serverError } = useLogin();
  const { mutate: verifyTOTP, isPending: verifying, error: verifyError } = useVerifyTOTP();
  const { data: existingUser, isLoading: meLoading } = useMe();
  const router = useRouter();
  const { register, handleSubmit, getValues, formState: { errors } } = useForm<LoginInput>({
    resolver: zodResolver(LoginSchema),
  });

  // The address already typed above is the one we send to. Asking for it a
  // second time on a second screen is the step people abandon.
  const sendMagicLink = async () => {
    const email = (getValues("email") || "").trim();
    if (!email.includes("@")) {
      setLinkError("Type your email address above first.");
      return;
    }
    setLinkError("");
    setLinkState("sending");
    try {
      await apiClient.post("/api/auth/magic-link", { email });
      setLinkState("sent");
    } catch (err: unknown) {
      setLinkState("idle");
      setLinkError(
        (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message ??
          "Could not send the link. Try your password instead.",
      );
    }
  };

  useEffect(() => {
    let cancelled = false;
    passkeysSupported().then((ok) => {
      if (!cancelled) setPasskeyReady(ok);
    });
    return () => {
      cancelled = true;
    };
  }, []);

  const { mutate: signInWithPasskey, isPending: passkeyPending } = usePasskeyLogin();

  const onPasskey = () => {
    setPasskeyError("");
    signInWithPasskey(undefined, {
      onError: (err: unknown) => {
        // Closing the sheet or tapping Cancel is a NotAllowedError, and it is
        // not a failure: telling somebody "passkey sign-in failed" because they
        // changed their mind is how a button stops being trusted.
        if ((err as { name?: string })?.name === "NotAllowedError") return;
        setPasskeyError(
          (err as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error
            ?.message ?? "That passkey was not accepted. Use your password instead.",
        );
      },
    });
  };

  // v3.31.15: if the session cookie is still valid, don't show the
  // login form — bounce straight to the dashboard.
  useEffect(() => {
    if (!meLoading && existingUser) {
      router.replace(existingUser.role === "USER" ? "/account" : "/dashboard");
    }
  }, [meLoading, existingUser, router]);

  const onSubmit = (data: LoginInput) =>
    login(data, {
      onSuccess: (res) => {
        const d = res.data as { totp_required?: boolean; pending_token?: string; method?: string };
        if (d?.totp_required && d.pending_token) {
          setPendingToken(d.pending_token);
          // Which factor the account uses. Empty means an authenticator: that
          // is what every account meant before codes by email existed.
          setFactorMethod(d.method ?? "app");
        }
      },
    });

  const onVerify = (e: React.FormEvent) => {
    e.preventDefault();
    if (!pendingToken) return;
    verifyTOTP({
      pending_token: pendingToken,
      code,
      trust_device: trustDevice,
      backup: useBackup,
    });
  };

  const errText = (e: unknown) =>
    (e as { response?: { data?: { error?: { message?: string } } } })?.response?.data?.error?.message;
  const message = errText(serverError);
  const verifyMessage = errText(verifyError);

  // Step two: the password was accepted and the account has 2FA on. Rendered
  // instead of the credentials form rather than beneath it, so there is one
  // obvious thing to do.
  if (pendingToken) {
    return (
      <AuthShell
        mode="login"
        title="Two-factor authentication"
        subtitle={
          useBackup
            ? "Enter one of your backup codes"
            : factorMethod === "email"
              ? "We emailed you a 6-digit code. It expires in a few minutes."
              : "Enter the 6-digit code from your authenticator app"
        }
        errorMessage={verifyMessage}
      >
        <form onSubmit={onVerify} className="space-y-5">
          <div className="space-y-2">
            <label htmlFor="totp-code" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
              {useBackup ? "Backup code" : "Authentication code"}
            </label>
            <input
              id="totp-code"
              value={code}
              onChange={(e) => setCode(e.target.value)}
              className={inputOk + " text-center tracking-[0.4em] text-lg"}
              placeholder={useBackup ? "XXXXXXXX" : "000000"}
              inputMode={useBackup ? "text" : "numeric"}
              autoComplete="one-time-code"
              maxLength={useBackup ? 8 : 6}
              autoFocus
            />
          </div>

          {!useBackup && (
            <label className="flex items-center gap-2 cursor-pointer text-sm" style={{ color: "var(--auth-muted)" }}>
              <input
                type="checkbox"
                checked={trustDevice}
                onChange={(e) => setTrustDevice(e.target.checked)}
                className="h-4 w-4 rounded border-[var(--auth-border)]"
              />
              Trust this device for 30 days
            </label>
          )}

          <button
            type="submit"
            disabled={verifying || code.trim().length < (useBackup ? 8 : 6)}
            className="w-full rounded-[var(--auth-radius)] py-3 font-medium text-white transition-opacity disabled:opacity-50"
            style={{ background: "var(--auth-primary)" }}
          >
            {verifying ? "Verifying…" : "Verify and sign in"}
          </button>

          <div className="flex items-center justify-between text-sm">
            <button
              type="button"
              onClick={() => { setUseBackup(!useBackup); setCode(""); }}
              style={{ color: "var(--auth-primary)" }}
            >
              {useBackup ? "Use your authenticator app" : "Use a backup code"}
            </button>
            <button
              type="button"
              onClick={() => { setPendingToken(null); setCode(""); setUseBackup(false); }}
              style={{ color: "var(--auth-muted)" }}
            >
              Back
            </button>
          </div>
        </form>
      </AuthShell>
    );
  }

  return (
    <AuthShell
      mode="login"
      title="Welcome back"
      subtitle="Sign in to your account"
      errorMessage={message}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
        <div className="space-y-2">
          <label htmlFor="email" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Email
          </label>
          <input
            id="email"
            type="email"
            {...register("email")}
            className={errors.email ? inputErr : inputOk}
            placeholder="you@example.com"
            autoFocus
          />
          {errors.email && <p className="text-sm text-red-500">{errors.email.message}</p>}
        </div>

        <div className="space-y-2">
          <label htmlFor="password" className="block text-sm font-medium" style={{ color: "var(--auth-muted)" }}>
            Password
          </label>
          <div className="relative">
            <input
              id="password"
              type={showPassword ? "text" : "password"}
              autoComplete="current-password"
              {...register("password")}
              className={(errors.password ? inputErr : inputOk) + " pr-12"}
              placeholder="Enter your password"
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 -translate-y-1/2"
              style={{ color: "var(--auth-muted)" }}
              aria-label={showPassword ? "Hide password" : "Show password"}
            >
              {showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
            </button>
          </div>
          {errors.password && <p className="text-sm text-red-500">{errors.password.message}</p>}
        </div>

        <div className="flex items-center justify-between text-sm">
          <label className="flex items-center gap-2 cursor-pointer" style={{ color: "var(--auth-muted)" }}>
            <input type="checkbox" className="h-4 w-4 rounded border-[var(--auth-border)]" />
            Remember me
          </label>
          <Link href="/forgot-password" style={{ color: "var(--auth-primary)" }}>
            Forgot password?
          </Link>
        </div>

        <button
          type="submit"
          disabled={isPending}
          className="w-full rounded-[var(--auth-radius)] py-3 font-medium disabled:opacity-50 transition-colors"
          style={{ background: "var(--auth-primary)", color: "var(--auth-primary-fg)" }}
        >
          {isPending ? "Signing in..." : "Sign In"}
        </button>

        {/* A password is a cognitive test, and WCAG 3.3.8 asks for a way past
            one. A passkey is the better way past it than an emailed link: it
            never leaves the device, there is nothing to phish, and it is one
            touch. Shown whenever this browser can produce one, which is the
            only question that can be answered before somebody says who they
            are. */}
        {passkeyReady && (
          <>
            <div className="flex items-center gap-3" aria-hidden="true">
              <span className="h-px flex-1" style={{ background: "var(--auth-border)" }} />
              <span className="text-xs" style={{ color: "var(--auth-muted)" }}>
                or
              </span>
              <span className="h-px flex-1" style={{ background: "var(--auth-border)" }} />
            </div>

            <button
              type="button"
              onClick={onPasskey}
              disabled={passkeyPending}
              className="flex w-full items-center justify-center gap-2 rounded-[var(--auth-radius)] border py-3 font-medium transition-colors disabled:opacity-50"
              style={{ borderColor: "var(--auth-border)", color: "var(--auth-text)" }}
            >
              <Fingerprint className="h-5 w-5" aria-hidden="true" />
              {passkeyPending ? "Waiting for your passkey..." : "Sign in with a passkey"}
            </button>

            {passkeyError && (
              <p role="alert" className="text-center text-sm text-red-500">
                {passkeyError}
              </p>
            )}
          </>
        )}

        {/* And the emailed link, for the people who never set a password
            because they signed up with Google. */}
        <div className="space-y-2 text-center text-sm">
          {linkState === "sent" ? (
            <p role="status" style={{ color: "var(--auth-muted)" }}>
              If that address has an account, a sign-in link is on its way. It works once and expires in 15 minutes.
            </p>
          ) : (
            <button
              type="button"
              onClick={sendMagicLink}
              disabled={linkState === "sending"}
              className="underline disabled:opacity-50"
              style={{ color: "var(--auth-primary)" }}
            >
              {linkState === "sending" ? "Sending a link..." : "Email me a sign-in link instead"}
            </button>
          )}
          {linkError && <p className="text-red-500">{linkError}</p>}
        </div>
      </form>
    </AuthShell>
  );
}
