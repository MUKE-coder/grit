package scaffold

// Biome is a generated project's linter and formatter. It replaced ESLint and
// Prettier: one tool, one config file, one dependency.
//
// What it replaced was not working. The Next.js apps ran `next lint`, which
// Next.js 16 removed, and the Vite apps ran `eslint .` with no ESLint
// installed, so `pnpm lint` failed on every new project. Prettier was
// configured (single quotes, no semicolons) against templates written in the
// opposite style, and nothing ever checked it.

// biomeVersion is pinned exactly: a new Biome release can add a recommended
// rule, and a lint that starts failing on its own is one people switch off.
const biomeVersion = "2.5.14"

// biomeDevDependency is the devDependencies line, as the templates write it.
const biomeDevDependency = `"@biomejs/biome": "` + biomeVersion + `"`

// biomeLintScript checks lint rules, not formatting. The templates are not
// Biome-formatted (hundreds of files would change), and every
// `grit generate resource` writes more, so a format check would fail on a new
// project and again after each resource. Formatting is `pnpm format`, on
// demand.
const biomeLintScript = "biome check --formatter-enabled=false ."

// biomeFormatScript rewrites files in place.
const biomeFormatScript = "biome format --write ."

// biomeConfigFile is the config's name. .jsonc because Biome rejects comments
// in biome.json, and every rule switched off below says why.
const biomeConfigFile = "biome.jsonc"

// biomeConfig is the project's biome.jsonc.
//
// The formatter follows the code the templates already contain: double quotes,
// semicolons, two spaces, 100 columns. The Prettier config it replaces asked
// for single quotes and no semicolons, which no template followed; matching it
// would have changed 258 of a new project's 272 files on the first format
// instead of 175.
//
// The rules are Biome's recommended set, minus the ones switched off below,
// each with its reason. Several of those fire on the generated admin panel
// itself and can be switched back on as the templates are fixed.
//
// inFrontend is a single app, whose config sits in frontend/ while its
// .gitignore sits a level up.
func biomeConfig(inFrontend bool) string {
	vcsRoot := ""
	if inFrontend {
		vcsRoot = `, "root": ".."`
	}
	return `{
  "$schema": "https://biomejs.dev/schemas/` + biomeVersion + `/schema.json",
  "vcs": { "enabled": true, "clientKind": "git", "useIgnoreFile": true` + vcsRoot + ` },
  "files": {
    "ignoreUnknown": true,
    "includes": [
      "**",
      "!**/.next",
      "!**/dist",
      "!**/build",
      "!**/out",
      "!**/.turbo",
      "!**/coverage",
      "!**/.source",
      // Written by tools, not people.
      "!**/routeTree.gen.ts",
      "!**/next-env.d.ts",
      "!**/wailsjs",
      // Not linted yet: the Expo app and the docs site have no lint script.
      "!apps/expo",
      "!apps/docs"
    ]
  },
  "formatter": {
    "enabled": true,
    "indentStyle": "space",
    "indentWidth": 2,
    "lineWidth": 100,
    "lineEnding": "lf"
  },
  "javascript": {
    "formatter": {
      "quoteStyle": "double",
      "jsxQuoteStyle": "double",
      "semicolons": "always",
      "trailingCommas": "es5",
      "arrowParentheses": "always",
      "bracketSpacing": true
    }
  },
  // Tailwind 4's @theme, @apply and @custom-variant are errors to a plain CSS parser.
  "css": { "parser": { "tailwindDirectives": true } },
  // index.html carries one minified inline script: the theme, applied before first paint.
  "html": { "linter": { "enabled": false } },
  // Sorting imports would rewrite nearly every generated file.
  "assist": { "enabled": false },
  "linter": {
    "enabled": true,
    "rules": {
      "preset": "recommended"` + biomeRuleOverrides + `
    }
  }
}
`
}

// biomeRuleOverrides switches off the recommended rules that fire on what Grit
// generates, each with its reason in the config itself.
//
// Two kinds of reason. Some rules misread a pattern the templates use on
// purpose. The rest ("the admin panel ... switch it on when") are real findings
// in the generated admin that need template changes, and are the list to work
// down.
const biomeRuleOverrides = `,
      "a11y": {
        // The admin panel has 94 buttons with no type attribute; switch it on when they have one.
        "useButtonType": "off",
        // The admin panel has 35 labels not tied to their control; switch it on when they are.
        "noLabelWithoutControl": "off",
        // Modal backdrops close on click; Escape closes them from the keyboard.
        "noStaticElementInteractions": "off",
        "useKeyWithClickEvents": "off",
        // Inline icons sit beside a text label that already names the action.
        "noSvgWithoutTitle": "off",
        // The sign-in and dialog forms focus their first field on purpose.
        "noAutofocus": "off",
        // The checkbox and radio fields are styled role="checkbox" elements; switch it on when they are inputs.
        "useSemanticElements": "off",
        // Listbox options take focus through aria-activedescendant, not tabIndex.
        "useFocusableInteractive": "off",
        // A combobox's list is a ul with role="listbox", as the ARIA pattern describes.
        "noNoninteractiveElementToInteractiveRole": "off"
      },
      "suspicious": {
        // Skeleton rows and fixed-length lists that never reorder are keyed by position.
        "noArrayIndexKey": "off",
        // A handler may return nothing or a promise, which is what void | Promise says.
        "noConfusingVoidType": "off",
        // Turbo's env list covers build tasks; tests and Playwright read their own variables.
        "noUndeclaredEnvVars": "off",
        // Tailwind 3's @tailwind, in the desktop app.
        "noUnknownAtRules": { "level": "error", "options": { "ignore": ["tailwind"] } }
      },
      "performance": {
        // Images come from whatever storage host the project uses, which next/image needs listed in advance.
        "noImgElement": "off"
      },
      "complexity": {
        // The print stylesheet has to beat every utility class.
        "noImportantStyles": "off",
        // Seven a && a.b checks in the admin panel. They are correct; the rewrite is style.
        "useOptionalChain": "off",
        // An escaped - in a character class cannot be misread as a range.
        "noUselessEscapeInRegex": "off",
        // A case beside default names a value the switch knows about.
        "noUselessSwitchCase": "off"
      },
      "style": {
        // The templates are Go raw strings, which cannot hold a backtick, so strings are joined with +.
        "useTemplate": "off",
        // Used after checks TypeScript cannot follow.
        "noNonNullAssertion": "off",
        // Math.pow is what the templates use, and it reads the same.
        "useExponentiationOperator": "off"
      }`
