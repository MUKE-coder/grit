package plugin

// pairingWebHook drives the browser being paired: ask for a code, poll, stop.
func pairingWebHook() string {
	return `"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { api } from "@/lib/api";

export type PairingState =
  | { status: "idle" }
  | { status: "waiting"; code: string; qr: string; expiresAt: string }
  | { status: "approved" }
  | { status: "denied" }
  | { status: "expired" }
  | { status: "error"; message: string };

// A second is short enough to feel immediate on the browser and long enough
// that a phone taking twenty seconds to approve costs twenty requests, not two
// hundred. The code dies after two minutes, so this bounds the whole exchange
// at around 120 polls.
const POLL_MS = 1000;

export function usePairing() {
  const [state, setState] = useState<PairingState>({ status: "idle" });
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null);
  const stopped = useRef(false);

  const clear = useCallback(() => {
    if (timer.current) {
      clearTimeout(timer.current);
      timer.current = null;
    }
  }, []);

  const poll = useCallback(
    async (code: string) => {
      if (stopped.current) return;
      try {
        const { data } = await api.get("/api/pair/" + code);
        if (data?.data?.status === "approved") {
          // The server set the auth cookies alongside this response, so the
          // app is signed in the moment we get here.
          setState({ status: "approved" });
          return;
        }
        timer.current = setTimeout(() => poll(code), POLL_MS);
      } catch (err: any) {
        const code = err?.response?.data?.error?.code;
        if (code === "DENIED") setState({ status: "denied" });
        else if (code === "EXPIRED" || code === "NOT_FOUND") setState({ status: "expired" });
        else setState({ status: "error", message: "Could not reach the server" });
      }
    },
    [],
  );

  const start = useCallback(async () => {
    clear();
    stopped.current = false;
    setState({ status: "idle" });
    try {
      const { data } = await api.post("/api/pair/start");
      const d = data.data;
      setState({ status: "waiting", code: d.code, qr: d.qr_png, expiresAt: d.expires_at });
      timer.current = setTimeout(() => poll(d.code), POLL_MS);
    } catch (err: any) {
      setState({
        status: "error",
        message:
          err?.response?.data?.error?.message ?? "Could not start pairing",
      });
    }
  }, [clear, poll]);

  useEffect(() => {
    return () => {
      stopped.current = true;
      clear();
    };
  }, [clear]);

  return { state, start, restart: start };
}
`
}

// pairingWebPage is the screen the QR appears on.
func pairingWebPage() string {
	return `"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { usePairing } from "@/hooks/use-pairing";

/**
 * Sign in by scanning a QR from a device that is already signed in.
 *
 * The code is shown as text underneath as well as in the QR. A camera is not
 * always available, or working, or the thing the user reaches for, and a flow
 * with no fallback strands them.
 */
export default function LinkDevicePage() {
  const router = useRouter();
  const { state, start, restart } = usePairing();

  useEffect(() => {
    void start();
    // Once, on mount. Restarting is explicit, via the button.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (state.status === "approved") {
      router.replace("/dashboard");
    }
  }, [state.status, router]);

  return (
    <main className="flex min-h-screen items-center justify-center px-4 py-16">
      <div className="w-full max-w-md rounded-2xl border border-border bg-card p-8 text-center">
        <h1 className="text-2xl font-semibold text-foreground">Link this device</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          Open the app on a device you are already signed in on, and scan this code.
        </p>

        {state.status === "waiting" && (
          <>
            <div className="mt-8 flex justify-center">
              <Image
                src={state.qr}
                alt="Pairing QR code"
                width={240}
                height={240}
                unoptimized
                className="rounded-xl border border-border bg-white p-3"
              />
            </div>
            <p className="mt-6 text-xs uppercase tracking-wide text-muted-foreground">
              Or enter this code
            </p>
            <code className="mt-1 block break-all rounded-lg bg-muted px-3 py-2 font-mono text-xs text-foreground">
              {state.code}
            </code>
            <p className="mt-6 flex items-center justify-center gap-2 text-sm text-muted-foreground">
              <span className="h-2 w-2 animate-pulse rounded-full bg-primary" />
              Waiting for approval
            </p>
          </>
        )}

        {state.status === "approved" && (
          <p className="mt-8 text-sm text-emerald-500">Approved. Signing you in.</p>
        )}

        {state.status === "denied" && (
          <div className="mt-8">
            <p className="text-sm text-red-500">
              The request was refused on the other device.
            </p>
            <button
              onClick={() => void restart()}
              className="mt-4 rounded-lg border border-border px-4 py-2 text-sm text-foreground hover:bg-muted"
            >
              Try again
            </button>
          </div>
        )}

        {state.status === "expired" && (
          <div className="mt-8">
            <p className="text-sm text-muted-foreground">
              This code expired. Codes are short-lived on purpose.
            </p>
            <button
              onClick={() => void restart()}
              className="mt-4 rounded-lg border border-border px-4 py-2 text-sm text-foreground hover:bg-muted"
            >
              Show a new code
            </button>
          </div>
        )}

        {state.status === "error" && (
          <div className="mt-8">
            <p className="text-sm text-red-500">{state.message}</p>
            <button
              onClick={() => void restart()}
              className="mt-4 rounded-lg border border-border px-4 py-2 text-sm text-foreground hover:bg-muted"
            >
              Try again
            </button>
          </div>
        )}
      </div>
    </main>
  );
}
`
}

// pairingAdminHook is the approving side: look up a code, approve or deny it.
func pairingAdminHook() string {
	return `"use client";

import { useState, useCallback } from "react";
import { apiClient } from "@/lib/api-client";

export interface PairingRequestInfo {
  user_agent: string;
  ip: string;
  created_at: string;
  expires_at: string;
}

/**
 * Approving a device is two steps, deliberately.
 *
 * describe() fetches what the code is asking for so the screen can show it,
 * and only then does approve() send anything. Approving an opaque code is not
 * consent: the user has to be able to see that it is their own browser and not
 * somebody who read the QR over their shoulder.
 */
export function useDevicePairing() {
  const [info, setInfo] = useState<PairingRequestInfo | null>(null);
  const [code, setCode] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [done, setDone] = useState<"approved" | "denied" | null>(null);

  const fail = (err: any, fallback: string) =>
    setError(err?.response?.data?.error?.message ?? fallback);

  const describe = useCallback(async (raw: string) => {
    // A scanned QR carries the whole callback URL; a typed one is just the
    // code. Accept either rather than making the user edit what they scanned.
    const value = raw.trim().split("/").pop() ?? "";
    setBusy(true);
    setError(null);
    setInfo(null);
    setDone(null);
    try {
      const { data } = await apiClient.get("/api/pair/" + value + "/request");
      setCode(value);
      setInfo(data.data);
    } catch (err: any) {
      fail(err, "That code is not valid any more");
    } finally {
      setBusy(false);
    }
  }, []);

  const act = useCallback(
    async (what: "approve" | "deny") => {
      setBusy(true);
      setError(null);
      try {
        await apiClient.post("/api/pair/" + code + "/" + what);
        setDone(what === "approve" ? "approved" : "denied");
        setInfo(null);
      } catch (err: any) {
        fail(err, "That code is not valid any more");
      } finally {
        setBusy(false);
      }
    },
    [code],
  );

  const reset = useCallback(() => {
    setInfo(null);
    setCode("");
    setError(null);
    setDone(null);
  }, []);

  return {
    info,
    busy,
    error,
    done,
    describe,
    approve: () => act("approve"),
    deny: () => act("deny"),
    reset,
  };
}
`
}

// pairingAdminPage is the screen an approver uses.
func pairingAdminPage() string {
	return `"use client";

import { useState } from "react";
import { useDevicePairing } from "@/hooks/use-device-pairing";

/**
 * Approve a device that is asking to be linked.
 *
 * The confirmation step shows the browser and address the request came from,
 * because that is the only thing standing between a QR somebody photographed
 * across a room and their being signed into this account.
 */
export default function LinkDevicePage() {
  const { info, busy, error, done, describe, approve, deny, reset } = useDevicePairing();
  const [input, setInput] = useState("");

  return (
    <div className="mx-auto max-w-xl p-6">
      <h1 className="text-2xl font-semibold text-foreground">Link a device</h1>
      <p className="mt-2 text-sm text-muted-foreground">
        Enter or scan the code shown on the device asking to be linked. It will be
        signed in as you.
      </p>

      {!info && !done && (
        <form
          className="mt-6 flex gap-2"
          onSubmit={(e) => {
            e.preventDefault();
            void describe(input);
          }}
        >
          <input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Paste the code"
            autoComplete="off"
            spellCheck={false}
            className="flex-1 rounded-lg border border-border bg-background px-3 py-2 font-mono text-sm text-foreground"
          />
          <button
            type="submit"
            disabled={busy || input.trim() === ""}
            className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-50"
          >
            Continue
          </button>
        </form>
      )}

      {info && (
        <div className="mt-6 rounded-xl border border-amber-500/30 bg-amber-500/5 p-5">
          <p className="text-sm font-medium text-foreground">
            Approve this device?
          </p>
          <dl className="mt-4 space-y-2 text-sm">
            <div className="flex justify-between gap-4">
              <dt className="text-muted-foreground">Browser</dt>
              <dd className="text-right font-mono text-xs text-foreground">
                {info.user_agent || "unknown"}
              </dd>
            </div>
            <div className="flex justify-between gap-4">
              <dt className="text-muted-foreground">Address</dt>
              <dd className="font-mono text-xs text-foreground">{info.ip}</dd>
            </div>
          </dl>
          <p className="mt-4 text-xs text-muted-foreground">
            If you did not just open this on another device, refuse it.
          </p>
          <div className="mt-5 flex gap-2">
            <button
              onClick={() => void approve()}
              disabled={busy}
              className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground disabled:opacity-50"
            >
              Approve
            </button>
            <button
              onClick={() => void deny()}
              disabled={busy}
              className="rounded-lg border border-red-500/40 px-4 py-2 text-sm font-medium text-red-500 disabled:opacity-50"
            >
              That wasn&apos;t me
            </button>
          </div>
        </div>
      )}

      {done === "approved" && (
        <div className="mt-6 rounded-xl border border-emerald-500/30 bg-emerald-500/5 p-5">
          <p className="text-sm text-emerald-500">
            Device linked. It appears under Account -&gt; Security, where you can
            sign it out again.
          </p>
          <button
            onClick={() => {
              reset();
              setInput("");
            }}
            className="mt-4 text-sm text-muted-foreground underline"
          >
            Link another
          </button>
        </div>
      )}

      {done === "denied" && (
        <div className="mt-6 rounded-xl border border-border p-5">
          <p className="text-sm text-muted-foreground">
            Refused. The code is dead and cannot be used.
          </p>
          <button
            onClick={() => {
              reset();
              setInput("");
            }}
            className="mt-4 text-sm text-muted-foreground underline"
          >
            Start over
          </button>
        </div>
      )}

      {error && <p className="mt-4 text-sm text-red-500">{error}</p>}
    </div>
  );
}
`
}
