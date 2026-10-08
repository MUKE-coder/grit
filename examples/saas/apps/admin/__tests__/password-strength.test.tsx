import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { PasswordStrength } from "@/components/password-strength";
import { passwordFailures, passwordStrength } from "@/lib/password-rules";

/**
 * The password meter: the ladder, and the checklist beside it.
 *
 * Two separate questions, and the tests are separate too. passwordFailures
 * says whether the save is allowed and mirrors the API's internal/password
 * package. passwordStrength says how hard the password is to guess, and is
 * advisory. They disagree in both directions, which is the useful part.
 */
describe("the strength ladder", () => {
  it("runs from weakest to strongest", () => {
    const cases: Array<[string, string]> = [
      ["ab1", "Very weak"],
      // Twenty characters and about as hard to guess as four: every one after
      // the first continues a run, so it counts for a quarter.
      ["aaaaaaaaaaaaaaaaaaaa", "Very weak"],
      ["abcdefghijklmnopqrst", "Very weak"],
      // On the list attackers try first, however it is capitalised.
      ["password1", "Very weak"],
      ["Passw0rd", "Very weak"],
      // Three kinds of character, and too short to save: capped at Weak so the
      // bar does not contradict the cross beside it.
      ["lamp7!", "Weak"],
      ["aB3!xQ7z", "Fair"],
      ["Tr0ub4dor&3", "Strong"],
      ["correct horse battery staple", "Very strong"],
    ];
    for (const [password, want] of cases) {
      expect(`${password} -> ${passwordStrength(password).label}`).toBe(
        `${password} -> ${want}`,
      );
    }
  });

  it("does not call a password strong while the form refuses it", () => {
    for (const weak of ["lamp7!", "abcdefghij", "1234567890"]) {
      expect(passwordStrength(weak).score).toBeLessThanOrEqual(1);
    }
  });

  it("counts bytes, because that is what bcrypt counts", () => {
    // 71 accented letters is 142 bytes, well over the 72 bcrypt accepts, while
    // being 72 characters. Counting characters here would put back the 500 the
    // rule exists to prevent, for anybody not typing in ASCII.
    expect(passwordFailures("é".repeat(71) + "1")).toContain("max-length");
    // 73 bytes is over, 72 is the limit itself and allowed.
    expect(passwordFailures("a1" + "b".repeat(71))).toContain("max-length");
    expect(passwordFailures("a1" + "b".repeat(70))).not.toContain("max-length");
  });

  it("keeps the two questions apart for a passphrase bcrypt cannot hash", () => {
    const long =
      "correct horse battery staple correct horse battery staple correct horse 9!";
    // The checklist refuses it...
    expect(passwordFailures(long)).toContain("max-length");
    // ...and the meter still says what it is, rather than calling the
    // strongest kind of password weak.
    expect(passwordStrength(long).label).toBe("Very strong");
  });
});

describe("the checklist", () => {
  it("renders nothing until there is something to say", () => {
    const { container } = render(<PasswordStrength value="" />);
    expect(container.firstChild).toBeNull();
  });

  it("names the level and marks the rules that are met", () => {
    render(<PasswordStrength value="correct horse battery staple" />);
    expect(screen.getByRole("status").textContent).toBe("Very strong");
    expect(screen.getByText(/At least 8 characters/).textContent).toContain(": met");
    expect(screen.getByText(/At most 72 characters/).textContent).toContain(": met");
  });

  it("tells a screen reader which rules are not met yet", () => {
    render(<PasswordStrength value="abcdefghij" />);
    expect(screen.getByRole("status").textContent).toBe("Very weak");
    expect(
      screen.getByText(/a number, a symbol or a space/).textContent,
    ).toContain(": not met yet");
  });

  it("sees the person's own name inside their password", () => {
    render(
      <PasswordStrength value="okello-was-here-99" about={["Ada", "Okello"]} />,
    );
    expect(
      screen.getByText(/Nothing from your name or email/).textContent,
    ).toContain(": not met yet");
  });
});
