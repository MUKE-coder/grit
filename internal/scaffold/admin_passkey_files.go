package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeAdminPasskeyFiles writes the passkey card and its WebAuthn plumbing.
func writeAdminPasskeyFiles(root string, opts Options) error {
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{
		filepath.Join(adminRoot, "lib", "webauthn.ts"):                 adminWebauthnLibTS(),
		filepath.Join(adminRoot, "components", "security", "passkeys.tsx"): adminPasskeysCardTSX(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", opts.Module())
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func adminWebauthnLibTS() string {
	return `/**
 * The base64url plumbing WebAuthn needs, and nothing else.
 *
 * The browser's credential APIs take and return ArrayBuffers. JSON does not
 * carry those, so the server sends base64url strings and every field has to be
 * converted on the way in and out. Getting one of them wrong produces a
 * DOMException that names no field, which is why this lives in one place
 * instead of being inlined at three call sites.
 *
 * base64url, not base64: the alphabet uses - and _ and drops the padding.
 * Feeding a standard-base64 decoder a base64url string mostly works and then
 * fails on the inputs containing + or /, which is roughly one challenge in
 * thirty and looks like a flaky authenticator.
 */

export function fromBase64url(value: string): Uint8Array {
  const padded = value.replace(/-/g, "+").replace(/_/g, "/");
  const binary = atob(padded + "=".repeat((4 - (padded.length % 4)) % 4));
  const out = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) out[i] = binary.charCodeAt(i);
  return out;
}

export function toBase64url(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  let binary = "";
  for (let i = 0; i < bytes.length; i++) binary += String.fromCharCode(bytes[i]);
  return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

/** Is a platform authenticator available at all? */
export async function passkeysSupported(): Promise<boolean> {
  if (typeof window === "undefined" || !window.PublicKeyCredential) return false;
  try {
    return await window.PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
  } catch {
    return false;
  }
}

/** Turn the server's creation options into what navigator.credentials wants. */
export function toCreationOptions(o: Record<string, any>): PublicKeyCredentialCreationOptions {
  return {
    ...o,
    challenge: fromBase64url(o.challenge),
    user: { ...o.user, id: fromBase64url(o.user.id) },
    excludeCredentials: (o.excludeCredentials ?? []).map((c: any) => ({
      ...c,
      id: fromBase64url(c.id),
    })),
  } as unknown as PublicKeyCredentialCreationOptions;
}

/** And the request options, for signing in. */
export function toRequestOptions(o: Record<string, any>): PublicKeyCredentialRequestOptions {
  return {
    ...o,
    challenge: fromBase64url(o.challenge),
    allowCredentials: (o.allowCredentials ?? []).map((c: any) => ({
      ...c,
      id: fromBase64url(c.id),
    })),
  } as unknown as PublicKeyCredentialRequestOptions;
}

/** The registration answer, in the shape the server parses. */
export function encodeAttestation(cred: PublicKeyCredential) {
  const r = cred.response as AuthenticatorAttestationResponse;
  return {
    id: cred.id,
    rawId: toBase64url(cred.rawId),
    type: cred.type,
    clientExtensionResults: cred.getClientExtensionResults(),
    response: {
      clientDataJSON: toBase64url(r.clientDataJSON),
      attestationObject: toBase64url(r.attestationObject),
    },
  };
}

/** The sign-in answer. userHandle is what tells the server who this is. */
export function encodeAssertion(cred: PublicKeyCredential) {
  const r = cred.response as AuthenticatorAssertionResponse;
  return {
    id: cred.id,
    rawId: toBase64url(cred.rawId),
    type: cred.type,
    clientExtensionResults: cred.getClientExtensionResults(),
    response: {
      clientDataJSON: toBase64url(r.clientDataJSON),
      authenticatorData: toBase64url(r.authenticatorData),
      signature: toBase64url(r.signature),
      userHandle: r.userHandle ? toBase64url(r.userHandle) : null,
    },
  };
}

/** A readable default name, so the list is not four rows of "Passkey". */
export function guessDeviceName(): string {
  if (typeof navigator === "undefined") return "Passkey";
  const ua = navigator.userAgent;
  if (/iPhone|iPad/.test(ua)) return "iPhone or iPad";
  if (/Android/.test(ua)) return "Android device";
  if (/Mac OS X/.test(ua)) return "Mac";
  if (/Windows/.test(ua)) return "Windows PC";
  if (/Linux/.test(ua)) return "Linux machine";
  return "Passkey";
}
`
}

func adminPasskeysCardTSX() string {
	return `"use client";

import { useEffect, useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { KeyRound, Trash2, ShieldCheck } from "lucide-react";
import {
  passkeysSupported,
  toCreationOptions,
  encodeAttestation,
  guessDeviceName,
} from "@/lib/webauthn";

interface Passkey {
  id: string;
  name: string;
  synced: boolean;
  created_at: string;
  last_used_at?: string;
}

/**
 * Passkeys on the account.
 *
 * A passkey is a key pair the authenticator holds; the private half never
 * leaves the device and the server only ever stores the public one. That is
 * why this is worth having over a password: there is nothing here for a breach
 * to leak and nothing for a phishing page to collect.
 *
 * The card hides itself when the browser has no platform authenticator, rather
 * than offering a button that opens a dialog and fails. Support is a runtime
 * fact, not a configuration one.
 */
export function PasskeysCard() {
  const qc = useQueryClient();
  const [supported, setSupported] = useState<boolean | null>(null);
  const [busy, setBusy] = useState(false);

  useEffect(() => {
    passkeysSupported().then(setSupported);
  }, []);

  const { data: keys = [], isLoading } = useQuery<Passkey[]>({
    queryKey: ["passkeys"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/auth/passkeys");
      return data.data ?? [];
    },
  });

  const remove = useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete("/api/auth/passkeys/" + id);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["passkeys"] }),
  });

  async function addPasskey() {
    setBusy(true);
    try {
      const { data: begin } = await apiClient.post("/api/auth/passkeys/register/begin", {});
      const options = toCreationOptions(begin.data.options.publicKey ?? begin.data.options);

      const cred = (await navigator.credentials.create({
        publicKey: options,
      })) as PublicKeyCredential | null;
      if (!cred) throw new Error("No passkey was created");

      await apiClient.post(
        "/api/auth/passkeys/register/finish?session=" +
          encodeURIComponent(begin.data.session_id) +
          "&name=" +
          encodeURIComponent(guessDeviceName()),
        encodeAttestation(cred),
      );
      await qc.invalidateQueries({ queryKey: ["passkeys"] });
      toast.success("Passkey added");
    } catch (err) {
      // NotAllowedError is the user closing the dialog, which is not a failure
      // worth shouting about.
      const e = err as { name?: string; response?: { data?: { error?: { message?: string } } } };
      if (e?.name === "NotAllowedError") return;
      toast.error(e?.response?.data?.error?.message ?? "Could not add that passkey");
    } finally {
      setBusy(false);
    }
  }

  if (supported === false) return null;

  return (
    <section className="rounded-xl border border-border bg-bg-elevated p-6">
      <div className="flex items-start justify-between gap-4">
        <div className="flex items-start gap-3">
          <span className="mt-0.5 flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-accent">
            <KeyRound className="h-4 w-4" aria-hidden="true" />
          </span>
          <div>
            <h2 className="text-sm font-semibold text-foreground">Passkeys</h2>
            <p className="mt-0.5 text-xs text-text-muted">
              Sign in with your fingerprint, face or device PIN. Nothing to remember, and
              nothing a phishing page can collect.
            </p>
          </div>
        </div>
        {keys.length > 0 && (
          <span className="inline-flex shrink-0 items-center gap-1.5 rounded-full bg-success/10 px-2.5 py-1 text-xs font-medium text-success">
            <ShieldCheck className="h-3.5 w-3.5" aria-hidden="true" />
            {keys.length} registered
          </span>
        )}
      </div>

      <div className="mt-5 space-y-2">
        {isLoading ? (
          <p className="text-sm text-text-muted">Loading...</p>
        ) : keys.length === 0 ? (
          <p className="text-sm text-text-secondary">
            No passkeys yet. Add one and you can sign in on this device without a password.
          </p>
        ) : (
          <ul className="space-y-2">
            {keys.map((k) => (
              <li
                key={k.id}
                className="flex items-center justify-between gap-3 rounded-lg border border-border bg-bg-secondary px-3 py-2.5"
              >
                <div className="min-w-0">
                  <div className="truncate text-sm font-medium text-foreground">{k.name}</div>
                  <div className="text-xs text-text-muted">
                    Added {new Date(k.created_at).toLocaleDateString()}
                    {k.last_used_at
                      ? " · last used " + new Date(k.last_used_at).toLocaleDateString()
                      : " · never used"}
                    {k.synced ? " · synced across your devices" : ""}
                  </div>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  aria-label={"Remove " + k.name}
                  disabled={remove.isPending}
                  onClick={() => remove.mutate(k.id)}
                >
                  <Trash2 className="h-3.5 w-3.5" aria-hidden="true" />
                </Button>
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="mt-4">
        <Button onClick={addPasskey} disabled={busy || supported === null}>
          {busy ? "Waiting for your device..." : "Add a passkey"}
        </Button>
      </div>
    </section>
  );
}
`
}
