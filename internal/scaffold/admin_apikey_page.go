package scaffold

// The API-keys screen.
//
// One rule drives the whole design: the raw key exists for exactly one render.
// The server stores a hash, so there is no "show key" affordance to add later
// and no support path that recovers it. The reveal panel therefore has to be
// unmissable and refuses to close until the key has been copied or downloaded
// — the same contract as the 2FA backup codes, for the same reason.
//
// Revoked keys stay in the list rather than disappearing. Revocation is a
// timestamp, not a delete, because "which key did this?" is a question people
// ask about keys that were turned off months ago.

func adminAPIKeysPage() string {
	return `"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { PageHeader } from "@/components/chrome/PageHeader";
import { ConfirmModal } from "@/components/ui/confirm-modal";
import { Pill, type PillTone } from "@/components/ui/pill";
import { Note } from "@/components/ui/note";
import { EmptyState } from "@/components/ui/empty-state";
import { Segmented } from "@/components/ui/segmented";
import { RadioCard } from "@/components/ui/radio-card";
import { useAPIKeys, useCreateAPIKey, useRevokeAPIKey } from "@/hooks/use-api-keys";
import { toast } from "sonner";
import {
  KeyRound,
  Plus,
  Copy,
  Check,
  Loader2,
  Trash2,
  AlertTriangle,
  Download,
  Globe,
  Server,
} from "@/lib/icons";
import { Button } from "@/components/ui/button";
import { inputClasses } from "@/components/ui/input";

interface APIKey {
  id: string;
  name: string;
  kind: "publishable" | "secret";
  prefix: string;
  /** Present only for publishable keys, which were never secret. */
  token?: string | null;
  scopes: string[] | null;
  endpoints: string[] | null;
  origins: string[] | null;
  rate_limit: number;
  last_used_at?: string | null;
  expires_at?: string | null;
  revoked_at?: string | null;
  created_at: string;
}

/** Splits a textarea into a list, dropping blanks and stray whitespace. */
function lines(raw: string): string[] {
  return raw
    .split(/[\n,]/)
    .map((l) => l.trim())
    .filter(Boolean);
}

function relative(iso?: string | null): string {
  if (!iso) return "never";
  const diff = Date.now() - new Date(iso).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 1) return "just now";
  if (mins < 60) return mins + "m ago";
  const hours = Math.floor(mins / 60);
  if (hours < 24) return hours + "h ago";
  return Math.floor(hours / 24) + "d ago";
}

/** When a key stops working, in words rather than a timestamp nobody subtracts. */
function expiry(iso?: string | null): string {
  if (!iso) return "never";
  const days = Math.ceil((new Date(iso).getTime() - Date.now()) / 86400000);
  if (days < 0) return "expired";
  if (days === 0) return "today";
  if (days === 1) return "tomorrow";
  if (days < 60) return "in " + days + " days";
  return "in " + Math.round(days / 30) + " months";
}

/**
 * The one and only render of a secret key.
 *
 * It is a dialog rather than a panel on the page because of where the eye is
 * when it appears. The Create button is at the bottom of a tall form, so an
 * in-page panel rendered above that form arrives outside the viewport: you
 * could create a secret key, see a form reset, scroll away, and never learn
 * the key existed. There is no second chance, because the server stores a
 * hash.
 *
 * Closing it is still allowed. A dialog you cannot leave is a keyboard trap,
 * and the people it traps are the ones using a keyboard. So the exits work
 * and, until the key has been copied or downloaded, they ask once.
 */
function RevealModal({ token, onClose }: { token: string; onClose: () => void }) {
  const [saved, setSaved] = useState(false);
  const [leaving, setLeaving] = useState(false);
  const panel = useRef<HTMLDivElement>(null);
  const copyButton = useRef<HTMLButtonElement>(null);

  // Focus starts on Copy, which is what somebody is here to do.
  useEffect(() => {
    copyButton.current?.focus();
  }, []);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        e.preventDefault();
        if (saved) onClose();
        else setLeaving(true);
        return;
      }
      if (e.key !== "Tab") return;
      // Tab stays inside the dialog, which is what aria-modal promises.
      const stops = panel.current?.querySelectorAll<HTMLElement>(
        "button:not([disabled]), [href], [tabindex]:not([tabindex='-1'])",
      );
      if (!stops || stops.length === 0) return;
      const first = stops[0];
      const last = stops[stops.length - 1];
      if (e.shiftKey && document.activeElement === first) {
        e.preventDefault();
        last.focus();
      } else if (!e.shiftKey && document.activeElement === last) {
        e.preventDefault();
        first.focus();
      }
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [saved, onClose]);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(token);
      toast.success("Secret key copied.");
    } catch {
      // Insecure origin, or a document that is not focused. Counted as saved
      // anyway: refusing would leave somebody holding a key they can read on
      // screen and a dialog that will not believe them. Download still works.
      toast.error("Could not reach the clipboard. Select the key and copy it, or download it.");
    }
    setSaved(true);
  };

  const download = () => {
    const blob = new Blob([token + "\n"], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "api-key.txt";
    a.click();
    URL.revokeObjectURL(url);
    setSaved(true);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="fixed inset-0 bg-black/60"
        onClick={() => (saved ? onClose() : setLeaving(true))}
      />
      <div
        ref={panel}
        role="dialog"
        aria-modal="true"
        aria-labelledby="reveal-title"
        aria-describedby="reveal-body"
        className="relative z-10 w-full max-w-lg rounded-2xl border border-border bg-bg-elevated shadow-2xl motion-safe:animate-in motion-safe:fade-in motion-safe:zoom-in-95"
      >
        <div className="flex items-start gap-3 border-b border-border px-5 py-4">
          <span
            className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-warning-soft text-warning"
            aria-hidden="true"
          >
            <AlertTriangle className="h-4 w-4" />
          </span>
          <div className="min-w-0">
            <h2 id="reveal-title" className="text-[17px] font-semibold text-foreground">
              Copy your key now
            </h2>
            <p id="reveal-body" className="mt-1 text-[13px] text-text-secondary">
              This is the only time it will be shown. The server keeps a hash, so nobody can
              recover it later, not even you. If you lose it, revoke this key and make another.
            </p>
          </div>
        </div>

        <div className="px-5 py-4">
          <code className="block break-all rounded-lg border border-border bg-bg-tertiary px-4 py-3 font-mono text-[13px] text-foreground">
            {token}
          </code>
        </div>

        {leaving ? (
          <div className="flex flex-wrap items-center justify-between gap-3 border-t border-danger-line bg-danger-soft px-5 py-4">
            <p className="text-[13px] text-foreground">
              Close without saving it? The key is gone for good.
            </p>
            <div className="flex items-center gap-2">
              <Button variant="ghost" size="sm" onClick={() => setLeaving(false)}>
                Keep it open
              </Button>
              <Button variant="danger" size="sm" onClick={onClose}>
                Close anyway
              </Button>
            </div>
          </div>
        ) : (
          <div className="flex flex-wrap items-center gap-2 border-t border-border px-5 py-4">
            <Button ref={copyButton} variant="outline" size="sm" onClick={copy}>
              {saved ? <Check className="h-3.5 w-3.5" /> : <Copy className="h-3.5 w-3.5" />}
              Copy
            </Button>
            <Button variant="outline" size="sm" onClick={download}>
              <Download className="h-3.5 w-3.5" />
              Download
            </Button>
            <Button
              size="sm"
              className="ml-auto"
              onClick={() => (saved ? onClose() : setLeaving(true))}
            >
              {saved ? "I have saved it" : "Close"}
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}

type Shown = "active" | "inactive" | "all";

export default function APIKeysPage() {
  const [name, setName] = useState("");
  const [kind, setKind] = useState<"publishable" | "secret">("publishable");
  const [endpoints, setEndpoints] = useState("");
  const [origins, setOrigins] = useState("");
  const [rateLimit, setRateLimit] = useState("");
  const [expiresIn, setExpiresIn] = useState("");
  const [token, setToken] = useState<string | null>(null);
  const [pendingRevoke, setPendingRevoke] = useState<APIKey | null>(null);
  const [shown, setShown] = useState<Shown>("active");

  // The requests, their query keys and their invalidation are in
  // hooks/use-api-keys.ts. What is left here is the form and the reveal.
  const { data, isLoading } = useAPIKeys();
  const create = useCreateAPIKey();
  const revoke = useRevokeAPIKey();

  const submit = () => {
    const body: Record<string, unknown> = { name: name.trim(), kind };
    if (lines(endpoints).length) body.endpoints = lines(endpoints);
    // Origins are a browser-only control and the field is only shown for a
    // publishable key, so a value typed before switching type is not sent.
    if (kind === "publishable" && lines(origins).length) body.origins = lines(origins);
    const rpm = parseInt(rateLimit, 10);
    if (!Number.isNaN(rpm) && rpm > 0) body.rate_limit = rpm;
    const days = parseInt(expiresIn, 10);
    if (!Number.isNaN(days) && days > 0) body.expires_in_days = days;
    create.mutate(body, {
      onSuccess: (d) => {
        // Only a secret key gets the one-time reveal. A publishable key is
        // readable from the table forever, so putting it behind a "save this
        // now or lose it" dialog would teach the wrong lesson about what it is.
        if (kind === "secret") setToken(d.token);
        else toast.success("Publishable key created. It stays readable below.");
        setName("");
        setEndpoints("");
        setOrigins("");
        setRateLimit("");
        setExpiresIn("");
      },
    });
  };

  const keys = (data ?? []) as APIKey[];

  const dead = (k: APIKey) =>
    !!k.revoked_at || (!!k.expires_at && new Date(k.expires_at) < new Date());

  // "Revoked" was the obvious name for the second filter and it would have
  // been a lie: an expired key is just as dead and is not revoked, so it
  // would have belonged to neither tab and been unreachable from this page.
  const counts = useMemo(
    () => ({
      active: keys.filter((k) => !dead(k)).length,
      inactive: keys.filter((k) => dead(k)).length,
      all: keys.length,
    }),
    [keys],
  );

  const visible = useMemo(
    () =>
      keys.filter((k) =>
        shown === "all" ? true : shown === "active" ? !dead(k) : dead(k),
      ),
    [keys, shown],
  );

  return (
    <div>
      <PageHeader
        title="API keys"
        subtitle="Long-lived credentials for scripts and integrations, instead of a login."
      />

      {token && <RevealModal token={token} onClose={() => setToken(null)} />}

      {/* Create a key.
          The type is first because it decides the rest of the form: what the
          key may reach, whether an origin list means anything, and whether
          this is the last time anybody sees it. */}
      <section className="mb-8 rounded-2xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-5 py-4">
          <h2 className="text-[15px] font-semibold text-foreground">Create a key</h2>
          <p className="mt-0.5 text-[13px] text-text-secondary">
            Name it after the thing that will use it, so it can be revoked without guesswork.
          </p>
        </header>

        <div className="space-y-5 px-5 py-5">
          <fieldset>
            <legend className="mb-2 text-xs font-medium uppercase tracking-wider text-text-muted">
              Type
            </legend>
            <div className="grid gap-3 sm:grid-cols-2">
              <RadioCard
                name="key-kind"
                value="publishable"
                checked={kind === "publishable"}
                onChange={(v) => setKind(v as "publishable")}
                title="Publishable"
                meta="grit_pk_..."
                icon={<Globe className="h-4 w-4" />}
                description="Safe in a browser or a mobile app. Reaches public endpoints only, never a protected one, and stays readable here because it is not a secret."
              />
              <RadioCard
                name="key-kind"
                value="secret"
                checked={kind === "secret"}
                onChange={(v) => setKind(v as "secret")}
                title="Secret"
                meta="grit_sk_..."
                icon={<Server className="h-4 w-4" />}
                description="Server to server. Acts as you, with your roles and permissions. Shown once when it is created and never again."
              />
            </div>
          </fieldset>

          <div className="grid gap-4 sm:grid-cols-[minmax(0,1fr)_10rem_10rem]">
            <div>
              <label htmlFor="key-name" className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-text-muted">
                Name
              </label>
              <input
                id="key-name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="nightly export"
                className={inputClasses()}
              />
            </div>
            <div>
              <label htmlFor="key-expiry" className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-text-muted">
                Expires in (days)
              </label>
              <input
                id="key-expiry"
                value={expiresIn}
                onChange={(e) => setExpiresIn(e.target.value)}
                placeholder="never"
                inputMode="numeric"
                className={inputClasses()}
              />
            </div>
            <div>
              <label htmlFor="key-rpm" className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-text-muted">
                Limit (req/min)
              </label>
              <input
                id="key-rpm"
                value={rateLimit}
                onChange={(e) => setRateLimit(e.target.value)}
                placeholder="unlimited"
                inputMode="numeric"
                className={inputClasses()}
              />
            </div>
          </div>

          {/* Restrictions. Two axes: what you may call, and where from. The
              second only exists for a publishable key, so for a secret key it
              is not shown rather than shown and ignored. */}
          <div className={"grid gap-4 " + (kind === "publishable" ? "sm:grid-cols-2" : "")}>
            <div>
              <label htmlFor="key-endpoints" className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-text-muted">
                Allowed endpoints
              </label>
              <textarea
                id="key-endpoints"
                value={endpoints}
                onChange={(e) => setEndpoints(e.target.value)}
                rows={3}
                placeholder={"GET /api/v1/public/products\nGET /api/v1/public/products/*"}
                className={inputClasses({ multiline: true }) + " font-mono text-xs"}
              />
              <p className="mt-1.5 text-[13px] text-text-muted">
                One per line, method and path. A trailing * matches a prefix. Empty means every
                route this type already allows.
              </p>
            </div>
            {kind === "publishable" && (
              <div>
                <label htmlFor="key-origins" className="mb-1.5 block text-xs font-medium uppercase tracking-wider text-text-muted">
                  Browser origins
                </label>
                <textarea
                  id="key-origins"
                  value={origins}
                  onChange={(e) => setOrigins(e.target.value)}
                  rows={3}
                  placeholder={"https://shopfront.com\nhttp://localhost:3000"}
                  className={inputClasses({ multiline: true }) + " font-mono text-xs"}
                />
                <p className="mt-1.5 text-[13px] text-text-muted">
                  Leave it empty for a mobile app: native clients send no Origin header, so an
                  allowlist would reject every request they make. It also stops nothing that is
                  not a browser.
                </p>
              </div>
            )}
          </div>

          <div>
            <Button onClick={submit} disabled={!name.trim()} loading={create.isPending}>
              {!create.isPending && <Plus className="h-3.5 w-3.5" />}
              Create {kind === "publishable" ? "publishable" : "secret"} key
            </Button>
            <Note className="mt-2.5">
              {kind === "publishable"
                ? "A publishable key is readable from the table below whenever you need it again."
                : "A secret key is shown once, in a dialog, and the server keeps only a hash of it."}
            </Note>
          </div>
        </div>
      </section>

      {/* The keys.
          Revocation is a timestamp rather than a delete, because "which key
          did this?" is a question people ask about keys that were turned off
          months ago. So dead keys stay, behind a filter that opens on the
          ones that still work. */}
      <section className="overflow-hidden rounded-2xl border border-border bg-bg-elevated">
        <header className="flex flex-wrap items-center justify-between gap-3 border-b border-border px-5 py-4">
          <h2 className="text-[15px] font-semibold text-foreground">Keys</h2>
          <Segmented
            label="Which keys to show"
            value={shown}
            onChange={setShown}
            options={[
              { value: "active", label: "Active", count: counts.active },
              { value: "inactive", label: "Inactive", count: counts.inactive },
              { value: "all", label: "All", count: counts.all },
            ]}
          />
        </header>

        {isLoading ? (
          <div className="px-5 py-12 text-center">
            <Loader2 className="mx-auto h-4 w-4 animate-spin text-text-muted" aria-hidden="true" />
          </div>
        ) : visible.length === 0 ? (
          <EmptyState
            icon={<KeyRound className="h-5 w-5" />}
            title={keys.length === 0 ? "No API keys yet" : "No " + shown + " keys"}
            description={
              keys.length === 0
                ? "Create one above for a script, a cron job, or an integration, and it appears here with its last use."
                : "Every key you have is in one of the other tabs."
            }
          />
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left text-xs font-medium uppercase tracking-wider text-text-muted">
                  <th className="px-5 py-3 font-medium">Name</th>
                  <th className="px-4 py-3 font-medium">Key</th>
                  <th className="px-4 py-3 font-medium">Last used</th>
                  <th className="px-4 py-3 font-medium">Expires</th>
                  <th className="px-4 py-3 font-medium">Status</th>
                  <th className="px-5 py-3" />
                </tr>
              </thead>
              <tbody>
                {visible.map((k) => {
                  const expired = !!k.expires_at && new Date(k.expires_at) < new Date();
                  const off = !!k.revoked_at || expired;
                  const status: { label: string; tone: PillTone } = k.revoked_at
                    ? { label: "Revoked", tone: "bad" }
                    : expired
                      ? { label: "Expired", tone: "warn" }
                      : { label: "Active", tone: "good" };
                  return (
                    <tr key={k.id} className="border-b border-border/60 last:border-0">
                      <td className="px-5 py-3 align-top">
                        <span className={off ? "text-text-muted line-through" : "text-foreground"}>
                          {k.name}
                        </span>
                        <div className="mt-1.5 flex flex-wrap items-center gap-1.5">
                          <Pill tone={k.kind === "publishable" ? "info" : "warn"}>{k.kind}</Pill>
                          {k.rate_limit > 0 && <Pill>{k.rate_limit}/min</Pill>}
                          {k.endpoints && k.endpoints.length > 0 && (
                            <Pill title={k.endpoints.join("\n")}>
                              {k.endpoints.length} endpoint{k.endpoints.length === 1 ? "" : "s"}
                            </Pill>
                          )}
                          {k.origins && k.origins.length > 0 && (
                            <Pill title={k.origins.join("\n")}>
                              {k.origins.length} origin{k.origins.length === 1 ? "" : "s"}
                            </Pill>
                          )}
                        </div>
                      </td>
                      <td className="px-4 py-3 align-top">
                        {/*
                          A publishable key is shown in full, with a copy button.
                          It is already in every copy of your app, so hiding it here
                          would protect nothing and cost you the ability to read it
                          when setting up a new environment.

                          A secret key shows only its prefix, because only the hash
                          was ever stored. The prefix is how the server finds the row
                          to compare against, so it is not a secret and it is the one
                          part that tells two secret keys apart.
                        */}
                        {k.kind === "publishable" && k.token ? (
                          <div className="flex items-center gap-2">
                            <code className="max-w-64 truncate font-mono text-xs text-foreground">
                              {k.token}
                            </code>
                            <button
                              type="button"
                              onClick={async () => {
                                try {
                                  await navigator.clipboard.writeText(k.token as string);
                                  toast.success("Publishable key copied.");
                                } catch {
                                  toast.error("Could not copy. Select it and copy by hand.");
                                }
                              }}
                              className="rounded p-1 text-text-muted transition-grit hover:bg-bg-hover hover:text-foreground"
                              aria-label={"Copy " + k.name}
                            >
                              <Copy className="h-3.5 w-3.5" />
                            </button>
                          </div>
                        ) : (
                          <code className="font-mono text-xs text-text-muted">
                            grit_sk_{k.prefix}_&hellip;
                          </code>
                        )}
                      </td>
                      <td className="px-4 py-3 align-top text-[13px] text-text-muted">
                        {relative(k.last_used_at)}
                      </td>
                      <td className="px-4 py-3 align-top text-[13px] text-text-muted">
                        {expiry(k.expires_at)}
                      </td>
                      <td className="px-4 py-3 align-top">
                        <Pill tone={status.tone} dot>
                          {status.label}
                        </Pill>
                      </td>
                      <td className="px-5 py-3 align-top text-right">
                        {!off && (
                          <button
                            type="button"
                            onClick={() => setPendingRevoke(k)}
                            className="inline-flex items-center gap-1.5 rounded-md px-2 py-1 text-[13px] text-text-muted transition-grit hover:bg-danger-soft hover:text-danger"
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                            Revoke
                          </button>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <Note className="mt-4">
        Send a key as <code className="font-mono">X-API-Key: grit_&hellip;</code> or{" "}
        <code className="font-mono">Authorization: Bearer grit_&hellip;</code>. Revoked keys stay
        listed under Inactive, because &ldquo;which key did this?&rdquo; is a question people ask
        about keys that were turned off months ago.
      </Note>

      <ConfirmModal
        open={!!pendingRevoke}
        title="Revoke this key?"
        description={
          "Anything using \"" + (pendingRevoke?.name ?? "") +
          "\" will start getting 401s immediately. This cannot be undone."
        }
        confirmLabel="Revoke"
        variant="danger"
        onCancel={() => setPendingRevoke(null)}
        onConfirm={() => {
          if (pendingRevoke) revoke.mutate(pendingRevoke.id);
          setPendingRevoke(null);
        }}
      />
    </div>
  );
}
`
}
