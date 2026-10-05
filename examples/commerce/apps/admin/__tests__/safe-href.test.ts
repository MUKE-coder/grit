import { describe, it, expect } from "vitest";
import { safeHref } from "@/lib/safe-href";

describe("safeHref", () => {
  it("keeps web, mail and phone links and paths on this site", () => {
    for (const href of ["https://example.com/a", "http://example.com", "mailto:a@example.com", "tel:+256700000000", "/system/jobs"]) {
      expect(safeHref(href)).toBe(href);
    }
  });

  it("refuses every other scheme and protocol-relative hosts", () => {
    for (const href of ["javascript:alert(1)", " JavaScript:alert(1)", "java\tscript:alert(1)", "data:text/html,<script>alert(1)</script>", "vbscript:msgbox(1)", "//evil.example", "/\\evil.example", "", null, undefined]) {
      expect(safeHref(href)).toBe("#");
    }
  });
});
