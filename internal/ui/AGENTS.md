# internal/ui

## Context

The CLI's design system: the palette, the glyphs, the line components and the
stepper. Everything `grit` prints in a terminal comes from here, or should.

Small on purpose. The whole package is a palette of seven colour roles, a
handful of glyphs with ASCII fallbacks, and about twenty functions that each
return or print one line.

## Source of truth

- **The palette**: `ui.go`. Seven roles, each a `pair(light, dark, ansiLight,
  ansiDark)`. `GRIT_THEME=light|dark` overrides detection.
- **The glyphs**: `ui.go`, with their ASCII fallbacks and `hasUTF8()`.
- **The line components**: `components.go`.
- **Progressive output**: `stepper.go`.
- **The huh prompt theme**: `prompt_theme.go`.

## Boundary rules

1. **Seven roles, and no more.** `cBrand`, `cMuted`, `cFaint`, `cSuccess`,
   `cWarning`, `cError`, and `Text`, which carries no colour on purpose. A new
   colour is a new role, and a new role needs a reason that is not "this line
   should stand out".

2. **Every colour is a pair with an ANSI fallback.** A terminal with 16 colours
   is not a broken terminal, and a hex-only style renders as nothing there.
   Use `pair(...)`, never a bare `lipgloss.Color`.

3. **Every glyph has an ASCII fallback.** A Windows console without a UTF-8 code
   page prints a box. `hasUTF8()` decides; see `GlyphOK` and `GlyphMark` for the
   shape.

4. **Honour `NO_COLOR` and `CI`.** `Interactive()` is the gate: no cursor
   movement, no erase, no spinner when it is false. A stepper that rewrites
   lines in CI produces a log nobody can read.

5. **Do not emit SGR by hand**, with one exception that is already written:
   `URL` does, because lipgloss v1.1.0 emits the underline pair per character.
   If you find another case, comment it the way that one is.

6. **The stepper erases exactly what it wrote.** `\033[1A\033[2K\r` assumes one
   line. A component that prints two and is then erased leaves half of itself on
   screen.

7. **This package prints; it does not decide.** No business logic, no file
   system, no network. A component takes strings and returns or prints a line.

## Validation workflow

```bash
go test ./internal/ui/
```

Then look at it, because a colour test asserts a code and not a contrast:

1. A **dark** terminal and a **light** one. The light-background case is the one
   that breaks, because a colour chosen on black is often invisible on white.
2. `NO_COLOR=1 grit ...`: still legible, still aligned.
3. `CI=1 grit ...`: no cursor movement, one line per step.
4. A console **without** UTF-8 (`chcp 437` on Windows): no boxes.
5. A narrow terminal, 80 columns. Check the padding helpers still line up.

The accessibility floor is WCAG AA on the background the text actually lands
on: 4.5:1 for body text. `#7c7c96` measures 4.9:1 on `#0a0a0f` and passes;
`#606078` measured 3.1:1 and did not.
