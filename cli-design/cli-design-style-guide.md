# Grit CLI style guide

How every `grit` command should look in the terminal. The goal is output that reads well on dark and light terminals without the user configuring anything.

The Go examples assume [lipgloss](https://github.com/charmbracelet/lipgloss) v1 and [huh](https://github.com/charmbracelet/huh), which the current prompts appear to use. Check field names against the versions in `go.mod`.

## 1. Principles

1. **Body text has no color.** Print it with the terminal's default foreground. That is the only color guaranteed to be readable on every theme.
2. **Every color is a pair.** Each role has a dark-terminal value and a light-terminal value. Never hard-code a single hex.
3. **Color never carries meaning alone.** Each status has its own glyph (`✓ ! ✗ +`), so output still reads with `NO_COLOR` set or when piped to a file.
4. **Brand color means "you can act on this".** Commands to type, the selected option, the spinner. Emphasis for names and totals is **bold**, not color.
5. **Report results, not activity.** Finished steps are a check and a noun. Only the running step animates.

## 2. Color roles

| Role | Dark terminal | Light terminal | ANSI-16 fallback (dark / light) | Use for |
|---|---|---|---|---|
| `brand` | `#6AA7FF` | `#0B57C9` | 12 bright blue / 4 blue | Wordmark, `❯` pointer, `?` prompt mark, spinner, progress fill, commands, note bar `│` |
| `text` | terminal default | terminal default | default | All body text |
| `muted` | `#939CAE` | `#586170` | 8 bright black / 8 bright black | Descriptions, comments, timings, labels, key hints, log timestamps |
| `faint` | `#4A5160` | `#C4C9D2` | 8 bright black / 7 white | Rules, pending dots, progress track. Never for words |
| `success` | `#3DD68C` | `#0A7343` | 10 bright green / 2 green | `✓`, `+`, the final result word ("Created", "Updated") |
| `warning` | `#F0B847` | `#855300` | 11 bright yellow / 3 yellow | The `!` glyph only |
| `error` | `#FF7B7B` | `#BE1F2B` | 9 bright red / 1 red | The `✗` glyph only |

Alternative brand, **Ember**: `#FF9255` dark, `#AD3F08` light. Swap only the `brand` pair; nothing else changes.

Contrast: every text role is at least 5.1:1 on black, `#0D1117`, `#1E1E1E`, One Dark `#282C34`, Solarized Dark, white, `#FAFAFA`, `#EEEEEE` and Solarized Light. `faint` is below that on purpose, which is why it is never used for words.

### Allowed combinations

| Element | Style |
|---|---|
| Command to type | `brand` |
| Selected option name | `brand` + bold |
| Unselected option name | `text` |
| Description next to an option or command | `muted` (selected row: `text`) |
| Result headline word | `success` + bold |
| Warning or error title | `text` + bold, after a colored glyph |
| Names, counts, new version | `text` + bold |
| URL | `text` + underline |
| Timing, version, path detail | `muted` |

Do not use: colored backgrounds, brand-colored body sentences, more than one colored word in a sentence, bright ANSI green/cyan/magenta directly.

## 3. Glyphs

| Meaning | Glyph | ASCII fallback | Color |
|---|---|---|---|
| Done | `✓` | `ok` | success |
| Added | `+` | `+` | success |
| Warning | `!` | `!` | warning |
| Failed | `✗` | `x` | error |
| Running | `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` | `-\|/` | brand |
| Pending | `·` | `.` | faint |
| Prompt | `?` | `?` | brand + bold |
| Pointer | `❯` | `>` | brand + bold |
| Note bar | `│` | `\|` | brand |
| Brand mark | `░▒▓█` | `grit` | brand |

Use the ASCII column when the output is not UTF-8 (`LANG`/`LC_ALL` without `UTF-8`, or legacy Windows console).

## 4. Identity

**Wordmark.** Shown only by `grit new`, `grit version` and bare `grit`. All in `brand`.

```
  ░▒▓████████  ███████   ██  ████████
  ░▒▓██        ██    ██  ██     ██
  ░▒▓██  ████  ███████   ██     ██
  ░▒▓██    ██  ██  ██    ██     ██
  ░▒▓████████  ██    ██  ██     ██

  Describe your data. Get the whole app.  v3.387.0
```

**One-line header.** Every other command starts with this instead of the wordmark:

```
  ░▒▓█ grit migrate  v3.387.0 · sqlite
```

Mark in `brand`, command name bold, the rest `muted`.

**Progress bar.** The grain is the leading edge: `████████████▓▒░·········  58%`. Fill in `brand`, track in `faint`, percent in `muted`.

## 5. Layout

- Two-space left gutter on every line. Nothing starts at column 1 except the shell prompt.
- One blank line between groups. No `=====` or `-----` rules. One `faint` line is allowed above the URL table.
- Second columns line up: pad command names and labels to the longest in the group plus 3 spaces.
- Keep lines under 80 columns. Wrap long text with a hanging indent of 4 spaces.
- Timings are right-hand, `muted`, formatted `8.4s` or `120ms`. Omit them under 100ms.
- Lists longer than 8 rows are cut with a `muted` count: `… 35 more, see them all with --verbose`.

## 6. Components

**Step list**

```
  ✓ Go dependencies                               8.4s
  ⠹ Installing frontend deps  ████████████▓▒░·········  58%
  · Admin panel
```

Done steps are nouns ("Go API"), not "Scaffolding Go API...". Pending step text is `muted`.

**Result and next steps**

```
  Created my-app  Single · Atlas · 14.2s

  Next steps

    cd my-app
    docker compose up -d    Postgres, Redis, MinIO, Mailhog
    grit migrate            create database tables
```

"Created" is `success` + bold, the name is bold, commands are `brand`, comments are `muted` with no `#`.

**Answered prompt** collapses to one line: `✓ Architecture  Single · Go API + embedded React SPA` (label and detail `muted`).

**Warning.** Title, what still works, fix.

```
  ! MinIO is not answering  http://localhost:9000
    Files are kept on local disk at storage/app. Uploads still work.
    Fix  docker compose up -d minio
```

**Error.** Same shape with `✗` and a `Try` line. Exit non-zero.

```
  ✗ Port 8080 is already in use
    The API could not start because another process holds the port.
    Try  grit start --port 8081
```

**Note.** For information that is not a problem, such as the baseline migration:

```
  │ Baseline run
  │ This run built the schema, so it cannot be rolled back.
  │ To start over: grit migrate --fresh
```

**Runtime logs.** Time only, a fixed-width `muted` source, then the message. Drop the date and banner lines.

```
  14:43:11  db       connected  sqlite, pool of 1
  14:43:11  migrate  43 models registered
```

## 7. Implementation

```go
package ui

import "github.com/charmbracelet/lipgloss"

func pair(light, dark, ansiLight, ansiDark string) lipgloss.CompleteAdaptiveColor {
	return lipgloss.CompleteAdaptiveColor{
		Light: lipgloss.CompleteColor{TrueColor: light, ANSI256: light, ANSI: ansiLight},
		Dark:  lipgloss.CompleteColor{TrueColor: dark, ANSI256: dark, ANSI: ansiDark},
	}
}

var (
	cBrand   = pair("#0B57C9", "#6AA7FF", "4", "12")
	cMuted   = pair("#586170", "#939CAE", "8", "8")
	cFaint   = pair("#C4C9D2", "#4A5160", "7", "8")
	cSuccess = pair("#0A7343", "#3DD68C", "2", "10")
	cWarning = pair("#855300", "#F0B847", "3", "11")
	cError   = pair("#BE1F2B", "#FF7B7B", "1", "9")

	Text    = lipgloss.NewStyle() // no color on purpose
	Bold    = lipgloss.NewStyle().Bold(true)
	Brand   = lipgloss.NewStyle().Foreground(cBrand)
	Muted   = lipgloss.NewStyle().Foreground(cMuted)
	Faint   = lipgloss.NewStyle().Foreground(cFaint)
	Success = lipgloss.NewStyle().Foreground(cSuccess)
	Warning = lipgloss.NewStyle().Foreground(cWarning)
	Error   = lipgloss.NewStyle().Foreground(cError)
	Link    = lipgloss.NewStyle().Underline(true)
)

const gutter = "  "

func Header(cmd, meta string) string {
	return gutter + Brand.Render("░▒▓█") + " " + Bold.Render("grit "+cmd) + "  " + Muted.Render(meta)
}
func Done(label, detail string) string {
	return gutter + Success.Render("✓") + " " + label + "  " + Muted.Render(detail)
}
func Warn(title, detail string) string {
	return gutter + Warning.Render("!") + " " + Bold.Render(title) + "  " + Muted.Render(detail)
}
func Fail(title string) string {
	return gutter + Error.Render("✗") + " " + Bold.Render(title)
}
```

Route every print through these helpers so no command builds its own colors.

**Prompts (huh)**

```go
func Theme() *huh.Theme {
	t := huh.ThemeBase()
	t.Focused.Base = lipgloss.NewStyle().PaddingLeft(2) // no left border
	t.Focused.Title = Bold
	t.Focused.Description = Muted
	t.Focused.SelectSelector = Brand.Bold(true).SetString("❯ ")
	t.Focused.SelectedOption = Brand.Bold(true)
	t.Focused.UnselectedOption = Text
	t.Help.ShortKey, t.Help.ShortDesc, t.Help.ShortSeparator = Muted, Muted, Faint
	t.Blurred = t.Focused
	return t
}
```

**Light or dark detection.** lipgloss asks the terminal for its background color and picks the right half of each pair. When detection is impossible it assumes dark. Give users an override:

```go
switch os.Getenv("GRIT_THEME") {
case "light":
	lipgloss.SetHasDarkBackground(false)
case "dark":
	lipgloss.SetHasDarkBackground(true)
}
```

On lipgloss v2, replace the pairs with `lipgloss.LightDark(isDark)` and detect once at startup.

**No color.** lipgloss already drops color when `NO_COLOR` is set or stdout is not a TTY. Also stop the spinner and progress bar in that case and print one plain line per finished step, so CI logs stay clean.

## 8. Checklist for a new command

- [ ] Starts with the one-line header, not the wordmark
- [ ] Body text uses no color
- [ ] Every colored word comes from a role in section 2
- [ ] Each status line has a glyph
- [ ] Ends with one result line: glyph, bold headline, muted timing
- [ ] Readable with `NO_COLOR=1` and with `GRIT_THEME=light`
