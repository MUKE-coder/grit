// Client-side rules for the formatted field types. The API applies the same
// rules in internal/fieldtypes and is the authority; these let a form say what
// is wrong before it asks.

import { COUNTRY_CODES } from "./countries";

const COUNTRIES = new Set<string>(COUNTRY_CODES);

/** What the JSON editor holds while its text does not parse. */
export class InvalidJSON {
  readonly text: string;
  readonly message: string;
  constructor(text: string, message: string) {
    this.text = text;
    this.message = message;
  }
}

/** A pasted address reduced to its host: https://www.example.com/a becomes www.example.com. */
export function toDomain(raw: string): string {
  let s = raw.trim().toLowerCase();
  const scheme = s.indexOf("://");
  if (scheme >= 0) s = s.slice(scheme + 3);
  s = s.split(/[/?#]/)[0];
  const at = s.lastIndexOf("@");
  if (at >= 0) s = s.slice(at + 1);
  s = s.split(":")[0];
  return s.replace(/\.$/, "");
}

/** #abc and ABCDEF become #aabbcc and #abcdef; anything else comes back as it was. */
export function normalizeColor(raw: string): string {
  const m = /^#?([0-9a-f]{3}|[0-9a-f]{6})$/i.exec(raw.trim());
  if (!m) return raw.trim();
  let hex = m[1].toLowerCase();
  if (hex.length === 3) hex = hex.split("").map((c) => c + c).join("");
  return "#" + hex;
}

/** "14:30" in the viewer's own clock, such as 2:30 PM. */
export function formatTimeOfDay(value: string): string {
  const m = /^(\d{1,2}):(\d{2})/.exec(value);
  if (!m) return value;
  const d = new Date(2000, 0, 1, Number(m[1]), Number(m[2]));
  return d.toLocaleTimeString([], { hour: "numeric", minute: "2-digit" });
}

/** The message for a value that breaks its type's rule, or true. */
export function validateFormat(type: string, value: unknown, max = 5): true | string {
  if (value === undefined || value === null || value === "") return true;
  const s = typeof value === "string" ? value.trim() : "";
  switch (type) {
    case "email":
      return /^[^\s@<>()]+@[^\s@<>()]+\.[^\s@<>()]+$/.test(s) || "Enter a valid email address";
    case "url":
      try {
        const u = new URL(s);
        return ((u.protocol === "http:" || u.protocol === "https:") && !!u.hostname) || "Enter an http or https address";
      } catch {
        return "Enter a web address, such as https://example.com";
      }
    case "domain":
      return /^([a-z0-9¡-￿]([a-z0-9¡-￿-]{0,61}[a-z0-9¡-￿])?\.)+[a-z¡-￿-]{2,63}$/i.test(s) || "Enter a domain, such as example.com";
    case "tel":
      return /^\+[1-9][0-9]{6,14}$/.test(s) || "Enter a valid phone number";
    case "country":
      return COUNTRIES.has(s.toUpperCase()) || "Choose a country";
    case "color":
      return /^#[0-9a-f]{6}$/i.test(s) || "Enter a hex colour, such as #6c5ce7";
    case "time":
      return /^([01][0-9]|2[0-3]):[0-5][0-9]/.test(s) || "Enter a time, such as 14:30";
    case "percent": {
      const n = typeof value === "number" ? value : Number(s);
      return (Number.isFinite(n) && n >= 0 && n <= 100) || "Enter a percentage from 0 to 100";
    }
    case "rating": {
      const n = typeof value === "number" ? value : Number(s);
      return (Number.isInteger(n) && n >= 0 && n <= max) || "Choose from 1 to " + max + " stars";
    }
    case "json":
      return value instanceof InvalidJSON ? value.message : true;
    default:
      return true;
  }
}
