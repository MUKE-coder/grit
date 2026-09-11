package scaffold

// adminI18nLib emits lib/i18n.tsx: translation for the admin, with no
// dependency.
//
// The base admin cannot import next-intl, because only projects that ran
// grit add i18n have it installed. Before v3.222.0 that meant the admin was
// English-only whatever --i18n did: the catalogues were written, the provider
// was wrapped around the layout, and not one component read from them. This
// is the seam that lets every component translate without depending on the
// library: t(key, fallback) returns the English the admin always showed until
// an I18nProvider supplies a catalogue, which is what grit add i18n wires.
func adminI18nLib() string {
	return `"use client";

import { createContext, useCallback, useContext, useMemo, type ReactNode } from "react";
import type { ResourceDefinition } from "@/lib/resource";

/*
 * Translation for the admin, with no dependency.
 *
 * Every visible string goes through t(key, fallback). Without an I18nProvider
 * (a project that never ran grit add i18n) t returns the fallback, which is the
 * English the admin has always shown, so this costs nothing until it is used.
 * grit add i18n wraps the layout in an I18nProvider fed from the same next-intl
 * catalogue as the rest of the app, and from then on any key present in
 * messages/<locale>.json replaces its fallback.
 *
 * Keys are dotted paths into the catalogue: "table.export", or
 * "resources.purchase-requests.fields.total" for a field label, where the
 * middle part is the definition's slug. A missing key
 * falls back, so a half-translated catalogue shows English where it has nothing
 * better, never a raw key.
 */

type Messages = Record<string, unknown>;

const MessagesContext = createContext<Messages | null>(null);

export function I18nProvider({ messages, children }: { messages: Messages; children: ReactNode }) {
  return <MessagesContext.Provider value={messages}>{children}</MessagesContext.Provider>;
}

function lookup(messages: Messages | null, key: string): string | undefined {
  if (!messages) return undefined;
  let node: unknown = messages;
  for (const part of key.split(".")) {
    if (!node || typeof node !== "object") return undefined;
    node = (node as Record<string, unknown>)[part];
  }
  return typeof node === "string" ? node : undefined;
}

export type Translate = (key: string, fallback: string, vars?: Record<string, string | number>) => string;

/** t(key, fallback, vars): the catalogue's text for key, or fallback, with {name} placeholders filled. */
export function useT(): Translate {
  const messages = useContext(MessagesContext);
  return useCallback<Translate>(
    (key, fallback, vars) => {
      const text = lookup(messages, key) ?? fallback;
      if (!vars) return text;
      return text.replace(/\{(\w+)\}/g, (whole, name: string) => (name in vars ? String(vars[name]) : whole));
    },
    [messages],
  );
}

/**
 * A resource definition with its labels translated: its own name, every column
 * and every form field, keyed under resources.<slug>. Pages call it where the
 * definition comes in, so the table, the forms and the headings below all read
 * translated labels without knowing translation exists.
 *
 *   "resources": {
 *     "purchase_requests": {
 *       "singular": "Demande d'achat",
 *       "plural": "Demandes d'achat",
 *       "fields": { "total": "Montant" }
 *     }
 *   }
 */
export function useLocalizedResource(resource: ResourceDefinition): ResourceDefinition {
  const t = useT();
  return useMemo(() => {
    const base = "resources." + resource.slug;
    const singular = resource.label?.singular ?? resource.name;
    const plural = resource.label?.plural ?? resource.slug;
    return {
      ...resource,
      label: { singular: t(base + ".singular", singular), plural: t(base + ".plural", plural) },
      table: {
        ...resource.table,
        columns: resource.table.columns.map((c) => ({ ...c, label: t(base + ".fields." + c.key, c.label) })),
      },
      form: resource.form && {
        ...resource.form,
        fields: resource.form.fields.map((f) => ({ ...f, label: t(base + ".fields." + f.key, f.label) })),
      },
    };
  }, [resource, t]);
}
`
}
