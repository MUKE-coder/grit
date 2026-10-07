package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// The password meter: how strong, said out loud, next to what is still missing.
//
// Before this, the only feedback anywhere was a five-segment bar on the admin
// account page that filled one segment per rule met. That reads as a strength
// meter and is not one. "a" fails only two of the five rules, so it lit three
// segments: a one-character password shown as sixty percent of the way there.
// Every other place a password is chosen (register, reset, both profile pages,
// the desktop client) showed nothing at all beyond the placeholder text.
//
// So there are two separate questions, and the UI now answers both:
//
//   - Is this allowed? The checklist, one line per rule, with a check when it
//     is satisfied. The server decides this; the list mirrors internal/password
//     so the answer is immediate, and the worst a drift can cost is a refused
//     save rather than a weak password.
//   - How hard is it to guess? The ladder, Very weak through Very strong. This
//     is advisory and deliberately client-only: there is no server rule to
//     mirror, and a Strength() in the Go package that no handler called would
//     be dead weight in every generated project.
//
// Both live in one component so that the five places that need it cannot drift
// from each other, and the component is written against no palette at all:
// layout utilities, the five ladder colours, and currentColor for everything
// else. That is what lets the same file serve the admin (Tailwind tokens), the
// web app (different Tailwind tokens), the auth pages (--auth-* variables) and
// the desktop client (inline styles) without a variant per frontend.

// passwordRulesTS emits lib/password-rules.ts: the rules, the failures and the
// strength estimate, with no JSX, so a form can use the logic without the
// component and the component has one source for both.
func passwordRulesTS() string {
	return `// The password rules, mirrored from the API's internal/password package, and
// the strength estimate, which exists only here.
//
// Keep RULES in step with password.Rules() on the server. The server is the
// authority: this decides what the checklist looks like, never whether a save
// is allowed, so a drift costs a confusing refusal and not a weak password.

export interface PasswordRule {
  id: string;
  label: string;
}

export const RULES: PasswordRule[] = [
  { id: "length", label: "At least 8 characters" },
  { id: "max-length", label: "At most 72 characters" },
  { id: "variety", label: "Letters and something else: a number, a symbol or a space" },
  { id: "not-common", label: "Not a password everyone tries first" },
  { id: "not-personal", label: "Nothing from your name or email" },
];

export const MIN_LENGTH = 8;

// bcrypt's ceiling, not ours: golang.org/x/crypto/bcrypt refuses anything over
// 72 bytes outright, so the API cannot hash a longer one. Bytes, because that
// is what bcrypt counts: an accented letter is two and an emoji is four, so a
// 60-character password can be over the limit.
export const MAX_BYTES = 72;

// The few hundred that attackers try first. The long lists belong in a
// service; most of the value is in the first page of this one.
const COMMON = new Set([
  "123456", "password", "12345678", "qwerty", "123456789", "12345", "1234", "111111",
  "1234567", "dragon", "123123", "baseball", "abc123", "football", "monkey", "letmein",
  "shadow", "master", "666666", "qwertyuiop", "123321", "mustang", "1234567890",
  "superman", "1qaz2wsx", "7777777", "121212", "000000", "qazwsx", "123qwe", "killer",
  "trustno1", "zxcvbnm", "asdfgh", "iloveyou", "starwars", "112233", "computer",
  "zxcvbn", "555555", "11111111", "131313", "freedom", "777777", "pass", "159753",
  "aaaaaa", "princess", "welcome", "admin", "letmein1", "password1", "password123",
  "passw0rd", "p@ssw0rd", "qwerty123", "iloveyou1", "welcome1", "admin123", "root",
  "changeme", "secret123", "hello", "test", "1111", "0000", "sunshine", "whatever",
]);

// Bytes of UTF-8, because that is the unit bcrypt's 72 is measured in.
//
// Counted by code point rather than with TextEncoder so this works in a test
// runner, a server render and a browser alike, with no feature check.
function byteLength(candidate: string): number {
  let bytes = 0;
  for (const char of candidate) {
    const code = char.codePointAt(0) ?? 0;
    if (code <= 0x7f) bytes += 1;
    else if (code <= 0x7ff) bytes += 2;
    else if (code <= 0xffff) bytes += 3;
    else bytes += 4;
  }
  return bytes;
}

function personalPieces(raw: string): string[] {
  const lower = raw.toLowerCase().trim();
  if (!lower) return [];
  const at = lower.indexOf("@");
  // The local part, and the domain without its suffix: somebody at acme.com
  // should not be using "acme" either.
  const source = at > 0 ? lower.slice(0, at) + " " + lower.slice(at + 1).split(".")[0] : lower;
  return source.split(/[^a-z0-9]+/).filter(Boolean);
}

/**
 * passwordFailures returns the ids of the rules this password does not pass,
 * in the order RULES lists them so the checklist can line them up.
 *
 * about is what the account already tells us: an email, a first name, a last
 * name. Passing nothing is fine where there is nothing to compare against yet.
 */
export function passwordFailures(candidate: string, about: string[] = []): string[] {
  const failed: string[] = [];

  if ([...candidate].length < MIN_LENGTH) failed.push("length");
  if (byteLength(candidate) > MAX_BYTES) failed.push("max-length");

  const letters = /\p{L}/u.test(candidate);
  const others = [...candidate].some((c) => !/\p{L}/u.test(c));
  if (!letters || !others) failed.push("variety");

  if (COMMON.has(candidate.toLowerCase().trim())) failed.push("not-common");

  const lower = candidate.toLowerCase();
  // Four characters or more only: a surname like "Ng" turning up inside an
  // otherwise good password is a coincidence, and refusing it would be a rule
  // nobody could satisfy.
  const personal = about
    .flatMap(personalPieces)
    .some((piece) => piece.length >= 4 && lower.includes(piece));
  if (personal) failed.push("not-personal");

  return failed;
}

export interface StrengthLevel {
  label: string;
  color: string;
}

// The ladder, weakest first. Red through green, with the two middle rungs
// distinguishable from each other without colour as well, because they sit at
// different heights on the bar.
export const STRENGTH_LEVELS: StrengthLevel[] = [
  { label: "Very weak", color: "#ff6b6b" },
  { label: "Weak", color: "#ff9f43" },
  { label: "Fair", color: "#fdcb6e" },
  { label: "Strong", color: "#55c57a" },
  { label: "Very strong", color: "#00b894" },
];

export interface PasswordStrengthResult extends StrengthLevel {
  /** 0 for Very weak through 4 for Very strong. */
  score: number;
  /** The guess-entropy estimate, in bits, the score came from. */
  bits: number;
}

// How many distinct characters an attacker has to consider, given which kinds
// this password actually uses. Rough on purpose: the point is the shape of the
// answer, not a number to four significant figures.
function poolSize(candidate: string): number {
  const chars = [...candidate];
  let pool = 0;
  if (/\p{Ll}/u.test(candidate)) pool += 26;
  if (/\p{Lu}/u.test(candidate)) pool += 26;
  if (/\p{Nd}/u.test(candidate)) pool += 10;
  // A space or a punctuation mark: anything printable that is not a letter or
  // a digit. Written as a negation rather than a character class so this file
  // needs no backtick, which the generator cannot carry.
  if (chars.some((c) => (c.codePointAt(0) ?? 0) <= 127 && !/[\p{L}\p{N}]/u.test(c))) pool += 33;
  // Outside ASCII: accents, other scripts, emoji. Counted conservatively,
  // because the real alphabet there is enormous and crediting all of it would
  // call one accented letter a strong password.
  if (chars.some((c) => (c.codePointAt(0) ?? 0) > 127)) pool += 100;
  return Math.max(pool, 2);
}

// Length, with the parts an attacker gets for free discounted.
//
// A character that repeats the one before it, or continues a run up or down by
// one, costs a guesser almost nothing: "aaaaaaaaaaaaaaaaaaaa" and
// "abcdefghijklmnopqrst" are twenty characters and about as hard to guess as
// four. Counting them at a quarter each is what stops the bar calling either
// of them very strong.
function effectiveLength(candidate: string): number {
  let total = 0;
  let previous = -1;
  for (const char of [...candidate]) {
    const code = char.codePointAt(0) ?? 0;
    const continues =
      previous >= 0 && (code === previous || code === previous + 1 || code === previous - 1);
    total += continues ? 0.25 : 1;
    previous = code;
  }
  return total;
}

/**
 * passwordStrength estimates how hard this password is to guess and puts it on
 * the five-rung ladder.
 *
 * Advisory, and separate from passwordFailures on purpose: the failures say
 * whether the save is allowed, this says whether it is a good idea. The two
 * disagree in both directions, which is the useful part. "Passw0rd!" passes
 * every rule and is Very weak. A 90-character passphrase is Very strong and
 * still cannot be saved, because bcrypt will not hash it.
 */
export function passwordStrength(candidate: string, about: string[] = []): PasswordStrengthResult {
  if (!candidate) return { score: 0, bits: 0, ...STRENGTH_LEVELS[0] };

  let bits = Math.round(effectiveLength(candidate) * Math.log2(poolSize(candidate)));

  const failed = passwordFailures(candidate, about);
  // Counting characters cannot see either of these. A password on the list is
  // guessed in seconds however long it is, and one built out of the person's
  // own name is the second thing tried.
  if (failed.includes("not-common")) bits = Math.min(bits, 8);
  if (failed.includes("not-personal")) bits = Math.min(bits, 20);

  // The usual bands. 28 bits is where a password stops being instant, 60 is
  // where it stops being cheap, and 80 is past the point where guessing it is
  // the attacker's best move.
  let score = 0;
  if (bits >= 80) score = 4;
  else if (bits >= 60) score = 3;
  else if (bits >= 36) score = 2;
  else if (bits >= 28) score = 1;

  // Nothing the form will refuse gets to look encouraging.
  //
  // "lamp7!" is six characters of three different kinds, which the arithmetic
  // above puts at 37 bits and therefore Fair. It is also below the minimum
  // length, so the save is blocked, and a bar reading Fair next to a cross
  // reading "At least 8 characters" tells somebody two different things at
  // once. Capped at Weak: still typing, not yet allowed.
  //
  // max-length is deliberately not in this list. A 90-character passphrase
  // genuinely is very strong and the checklist is the right place to explain
  // that bcrypt will not hash it; calling it weak would be a lie told to
  // somebody doing the best possible thing.
  if (failed.includes("length") || failed.includes("variety")) {
    score = Math.min(score, 1);
  }

  return { score, bits, ...STRENGTH_LEVELS[score] };
}
`
}

// passwordStrengthTSXRelative is the same component for an app with no @/
// alias.
//
// `grit new-desktop` produces a standalone Wails app whose routes import by
// relative path. One rewritten import is cheaper than a second copy of the
// component, and far cheaper than leaving that app's sign-up as the one screen
// in the product with no meter on it.
func passwordStrengthTSXRelative() string {
	return strings.ReplaceAll(
		passwordStrengthTSX(),
		`from "@/lib/password-rules"`,
		`from "../lib/password-rules"`,
	)
}

// passwordStrengthTSX emits components/password-strength.tsx: the ladder, the
// label and the checklist.
//
// Written against no palette. Layout utilities exist in every frontend; the
// five ladder colours are literals because they have to mean the same thing
// everywhere; and the track, the text and the check marks are currentColor at
// an opacity, so the component takes the colour of whatever it is dropped
// into. That is the whole reason one file can serve four apps.
func passwordStrengthTSX() string {
	return `"use client";

import { RULES, passwordFailures, passwordStrength, STRENGTH_LEVELS } from "@/lib/password-rules";

interface PasswordStrengthProps {
  /** What the person has typed so far. */
  value: string;
  /** Their email and name, so "nothing from your name" can be checked. */
  about?: string[];
  /**
   * Set this as the password input's aria-describedby, so a screen reader
   * reaching the field hears what is still missing.
   */
  id?: string;
  className?: string;
}

/**
 * The strength ladder and the rule checklist, under a new-password field.
 *
 * Renders nothing until there is something to say, so an untouched sign-up
 * form is not a wall of grey crosses.
 */
export function PasswordStrength({
  value,
  about = [],
  id = "password-strength",
  className = "",
}: PasswordStrengthProps) {
  const failed = passwordFailures(value, about);
  const strength = passwordStrength(value, about);

  if (!value) return null;

  return (
    <div id={id} className={"space-y-2.5 " + className}>
      <div className="flex items-center gap-3">
        {/* One rung per level, filled up to the score. aria-hidden because the
            sentence beside it says the same thing in words. */}
        <div className="flex flex-1 gap-1.5" aria-hidden="true">
          {STRENGTH_LEVELS.map((level, i) => (
            <span
              key={level.label}
              className="h-1 flex-1 rounded-full transition-colors"
              style={
                i <= strength.score
                  ? { background: strength.color }
                  : { background: "currentColor", opacity: 0.15 }
              }
            />
          ))}
        </div>
        {/* polite, so the level is announced once the typing pauses rather
            than on every keystroke. */}
        <p
          role="status"
          aria-live="polite"
          className="shrink-0 text-xs font-medium tabular-nums"
          style={{ color: strength.color }}
        >
          {strength.label}
        </p>
      </div>

      <ul className="space-y-1.5 text-sm">
        {RULES.map((rule) => {
          const ok = !failed.includes(rule.id);
          return (
            <li key={rule.id} className="flex items-start gap-2">
              <span
                className="mt-0.5 inline-flex h-4 w-4 shrink-0 items-center justify-center rounded-full border"
                style={
                  ok
                    ? { borderColor: "#00b894", background: "rgba(0, 184, 148, 0.15)", color: "#00b894" }
                    : { borderColor: "currentColor", opacity: 0.3, color: "transparent" }
                }
              >
                <svg
                  className="h-3 w-3"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth={3}
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  aria-hidden="true"
                >
                  <path d="M20 6 9 17l-5-5" />
                </svg>
              </span>
              <span style={ok ? undefined : { opacity: 0.65 }}>
                {rule.label}
                <span className="sr-only">{ok ? ": met" : ": not met yet"}</span>
              </span>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
`
}

// writePasswordStrengthFiles delivers the rules module and the meter to every
// frontend this project has.
//
// One writer, two callers: `grit new` and `grit upgrade`. A framework-owned
// frontend file needs both, and the ones that only made the first reached new
// projects and no existing one, quietly, because the missing file is not a
// compile error until something imports it.
func writePasswordStrengthFiles(root string, opts Options) error {
	// Where lib/ and components/ live for one frontend, and whether it speaks
	// Vite rather than Next: a Vite app cannot have "use client" at the top of
	// a module, which is what nextToTanStack takes off.
	type target struct {
		lib       string
		component string
		tanstack  bool
	}
	var targets []target

	// The admin panel, wherever it is: its own app, a route group inside the web
	// app, or a section of the SPA. adminLib and adminComponent know which.
	if opts.HasAdminPanel() {
		targets = append(targets, target{
			lib:       adminLib(root, opts, "password-rules.ts"),
			component: adminComponent(root, opts, "password-strength.tsx"),
			tanstack:  opts.UseTanStack(),
		})
	}

	// The web app, and a single project's frontend, which is the same app at a
	// different path: the project root rather than apps/web.
	if opts.ShouldIncludeWeb() || opts.ShouldIncludeSingleSPA() {
		base := webAppRoot(root, opts)
		vite := opts.UseTanStack()
		if vite {
			base = filepath.Join(base, "src")
		}
		targets = append(targets, target{
			lib:       filepath.Join(base, "lib", "password-rules.ts"),
			component: filepath.Join(base, "components", "password-strength.tsx"),
			tanstack:  vite,
		})
	}

	// The desktop client, which is a Vite app in a Wails window.
	if opts.ShouldIncludeDesktop() {
		base := filepath.Join(root, "apps", "desktop", "frontend", "src")
		targets = append(targets, target{
			lib:       filepath.Join(base, "lib", "password-rules.ts"),
			component: filepath.Join(base, "components", "password-strength.tsx"),
			tanstack:  true,
		})
	}

	for _, t := range targets {
		// An admin embedded in the web app shares that app's lib and components,
		// so two targets can name the same file. Harmless: both write the same
		// bytes, and writeFile is manifest-guarded either way.
		if err := writeFile(t.lib, passwordRulesTS()); err != nil {
			return fmt.Errorf("writing %s: %w", t.lib, err)
		}
		body := passwordStrengthTSX()
		if t.tanstack {
			body = nextToTanStack(body)
		}
		if err := writeFile(t.component, body); err != nil {
			return fmt.Errorf("writing %s: %w", t.component, err)
		}
	}
	return nil
}
