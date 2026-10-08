"use client";

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
