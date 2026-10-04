// safeHref returns a stored value as a link target only when following it is
// safe: an http or https URL, a mailto: or tel: address, or a path on this site.
// Anything else (javascript:, data:, vbscript:, a protocol-relative //host)
// comes back as "#".
//
// React refuses javascript: URLs, but not the others, and the values that reach
// these links are written by people outside the team: public form submissions,
// notification payloads.
const SAFE_HREF = /^(?:https?:\/\/|mailto:|tel:|\/(?![/\\]))/i;

export function safeHref(value: string | null | undefined): string {
  const href = typeof value === "string" ? value.trim() : "";
  return SAFE_HREF.test(href) ? href : "#";
}
