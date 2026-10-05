/**
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

/** A credential as the server sends it: the id is base64url, not bytes. */
interface CredentialJSON {
  id: string;
  [key: string]: unknown;
}

/** The server's creation options, as JSON. */
export interface CreationOptionsJSON {
  challenge: string;
  user: { id: string; [key: string]: unknown };
  excludeCredentials?: CredentialJSON[];
  [key: string]: unknown;
}

/** The server's request options, as JSON. */
export interface RequestOptionsJSON {
  challenge: string;
  allowCredentials?: CredentialJSON[];
  [key: string]: unknown;
}

/** Turn the server's creation options into what navigator.credentials wants. */
export function toCreationOptions(o: CreationOptionsJSON): PublicKeyCredentialCreationOptions {
  return {
    ...o,
    challenge: fromBase64url(o.challenge),
    user: { ...o.user, id: fromBase64url(o.user.id) },
    excludeCredentials: (o.excludeCredentials ?? []).map((c) => ({
      ...c,
      id: fromBase64url(c.id),
    })),
  } as unknown as PublicKeyCredentialCreationOptions;
}

/** And the request options, for signing in. */
export function toRequestOptions(o: RequestOptionsJSON): PublicKeyCredentialRequestOptions {
  return {
    ...o,
    challenge: fromBase64url(o.challenge),
    allowCredentials: (o.allowCredentials ?? []).map((c) => ({
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
