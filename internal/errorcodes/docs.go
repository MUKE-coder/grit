package errorcodes

import (
	"fmt"
	"strings"
)

// DocsTypeScript renders the catalogue as the data module the docs site reads.
//
// The page is generated rather than written for the same reason the project's own
// files are: a table of error codes maintained by hand is a table that lists what
// the API used to return. The entries keep the catalogue's declaration order, so
// the page reads in the order somebody meets these: the shared envelope first,
// then sign-in, then the rest.
func DocsTypeScript() string {
	areas, grouped := ByArea()

	var b strings.Builder
	b.WriteString(`// Generated from internal/errorcodes by tools/errorcodes. DO NOT EDIT.
//
// Run: go run ./tools/errorcodes
//
// The same catalogue generates internal/respond/codes.go and
// packages/shared/types/errors.ts in every project, which is what makes this
// table the documentation of what the API actually returns rather than a list
// somebody kept up by hand.

export interface ErrorCodeRow {
  code: string
  status: number
  category: string
  /** What happened, from the caller's side. */
  meaning: string
  /** What the caller should do about it. */
  client: string
}

export interface ErrorCodeArea {
  /** The key the API uses. */
  area: string
  /** The heading for a reader. */
  label: string
  codes: ErrorCodeRow[]
}

export const errorCodeAreas: ErrorCodeArea[] = [
`)

	for _, area := range areas {
		b.WriteString("  {\n")
		b.WriteString(fmt.Sprintf("    area: %s,\n", tsQuote(area)))
		b.WriteString(fmt.Sprintf("    label: %s,\n", tsQuote(AreaLabel(area))))
		b.WriteString("    codes: [\n")
		for _, entry := range grouped[area] {
			b.WriteString("      {\n")
			b.WriteString(fmt.Sprintf("        code: %s,\n", tsQuote(entry.Code)))
			b.WriteString(fmt.Sprintf("        status: %d,\n", entry.Status))
			b.WriteString(fmt.Sprintf("        category: %s,\n", tsQuote(string(entry.Category))))
			b.WriteString(fmt.Sprintf("        meaning: %s,\n", tsQuote(entry.Meaning)))
			b.WriteString(fmt.Sprintf("        client: %s,\n", tsQuote(entry.Client)))
			b.WriteString("      },\n")
		}
		b.WriteString("    ],\n")
		b.WriteString("  },\n")
	}

	b.WriteString("]\n\n")
	b.WriteString(fmt.Sprintf("/** How many codes the API documents. Shown on the page, so it cannot be stale. */\nexport const errorCodeCount = %d\n", len(All())))
	b.WriteString("\n/** Every row, flattened, for searching and for a test that checks coverage. */\nexport const errorCodes: ErrorCodeRow[] = errorCodeAreas.flatMap((area) => area.codes)\n")
	return b.String()
}

// tsQuote writes a TypeScript single-quoted string.
func tsQuote(value string) string {
	escaped := strings.ReplaceAll(value, `\`, `\\`)
	escaped = strings.ReplaceAll(escaped, "'", `\'`)
	return "'" + escaped + "'"
}
