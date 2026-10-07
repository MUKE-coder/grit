package scaffold

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every field where somebody chooses a password shows the meter.
//
// Written as a sweep rather than a list of the nine screens, because the list
// is what went wrong. The admin's account page had a checklist for months and
// the other eight screens had a placeholder reading "At least 8 characters" and
// nothing else, and nobody noticed, because there was no one place that said
// how many of them there were meant to be. A tenth screen added later is
// caught by this without anybody remembering to add it here.
//
// autoComplete="new-password" is the marker, and it is the right one: the
// browser asks for it on exactly the fields where a new password is set, and a
// screen that forgets it has a worse bug than a missing meter.
func TestEveryNewPasswordFieldShowsTheMeter(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}

	// The two places that are the meter itself, and the Expo app, which is
	// React Native: it has no DOM, so a component built from div and span
	// cannot render there. Its screens are a separate job.
	skipFile := map[string]bool{
		"password_strength.go": true,
		"expo_files.go":        true,
	}

	checked := 0
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") || skipFile[name] {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		src := strings.ReplaceAll(string(raw), "\r\n", "\n")

		starts := tmplHeader.FindAllStringSubmatchIndex(src, -1)
		for i, loc := range starts {
			end := len(src)
			if i+1 < len(starts) {
				end = starts[i+1][0]
			}
			body := src[loc[0]:end]
			fn := src[loc[2]:loc[3]]

			if !strings.Contains(body, `autoComplete="new-password"`) {
				continue
			}
			checked++
			if !strings.Contains(body, "<PasswordStrength") {
				t.Errorf("%s: %s has a new-password field and no <PasswordStrength>, so it asks for a password and says nothing about it",
					name, fn)
			}
			if !strings.Contains(body, "password-strength") {
				t.Errorf("%s: %s does not import the meter", name, fn)
			}
		}
	}
	if checked == 0 {
		t.Fatal("no new-password field was examined, so this test proves nothing")
	}
}

// Every password field says which password it is.
//
// This is the test that should have existed first. The sweep above finds the
// screens by their autoComplete="new-password", and five of the admin's six
// auth styles never set autoComplete at all: the marker used to find the
// screens was missing from exactly the screens missing everything else, so the
// sweep found one sign-up page of six and passed.
//
// The attribute matters on its own account. Without current-password a manager
// does not fill the login form; without new-password it neither offers to
// generate a password nor offers to save the one that was typed, which pushes
// people towards a password they can retype from memory. That is the opposite
// of what the meter beside it is asking for.
func TestEveryPasswordFieldDeclaresWhichPasswordItIs(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}

	checked := 0
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		src := strings.ReplaceAll(string(raw), "\r\n", "\n")

		for _, loc := range passwordTypeAttr.FindAllStringIndex(src, -1) {
			// The enclosing element: back to its <input, forward to its />.
			open := strings.LastIndex(src[:loc[0]], "<input")
			close := strings.Index(src[loc[1]:], "/>")
			if open < 0 || close < 0 {
				continue
			}
			tag := src[open : loc[1]+close]
			// A CSS selector, not an element: the e2e specs are full of
			// page.fill('input[type="password"]', ...). The backwards search for
			// "<input" then lands on some earlier element and the span runs
			// through half the file, so the giveaway is another tag inside it.
			if strings.Contains(tag[len("<input"):], "<") {
				continue
			}
			checked++
			if !strings.Contains(tag, "autoComplete=") {
				line := strings.Count(src[:open], "\n") + 1
				t.Errorf("%s:%d: a password field with no autoComplete, so a password manager can neither fill it nor save what is typed into it:\n  %s",
					name, line, strings.Join(strings.Fields(tag), " "))
			}
		}
	}
	if checked == 0 {
		t.Fatal("no password field was examined, so this test proves nothing")
	}
}

// A password input, whether the type is a literal or a show/hide toggle.
var passwordTypeAttr = regexp.MustCompile(`type=(?:"password"|\{[^}\n]*"password"[^}\n]*\})`)

// The meter has to reach every frontend, in every architecture.
//
// A framework-owned frontend file needs the new-project list and the upgrade
// writer both, and the ones that only made the first reached new projects and
// no existing one. writePasswordStrengthFiles is the single writer, so what
// this checks is that it knows about every shape: a missing branch here is an
// app whose register page imports a module that is not there, which is a blank
// screen rather than a build error in a Vite app.
func TestTheMeterReachesEveryFrontend(t *testing.T) {
	cases := []struct {
		name string
		opts Options
		want []string
	}{
		{
			name: "a triple with Next",
			opts: Options{ProjectName: "t", Architecture: ArchTriple, Frontend: FrontendNext},
			want: []string{
				"apps/admin/lib/password-rules.ts",
				"apps/admin/components/password-strength.tsx",
				"apps/web/lib/password-rules.ts",
				"apps/web/components/password-strength.tsx",
			},
		},
		{
			name: "a triple with Vite",
			opts: Options{ProjectName: "t", Architecture: ArchTriple, Frontend: FrontendTanStack},
			want: []string{
				"apps/admin/src/lib/password-rules.ts",
				"apps/admin/src/components/password-strength.tsx",
				"apps/web/src/lib/password-rules.ts",
				"apps/web/src/components/password-strength.tsx",
			},
		},
		{
			// A double has no admin app: the panel is a route group inside the
			// web app, and shares its lib and components.
			name: "a double",
			opts: Options{ProjectName: "t", Architecture: ArchDouble, Frontend: FrontendNext},
			want: []string{
				"apps/web/lib/password-rules.ts",
				"apps/web/components/password-strength.tsx",
			},
		},
		{
			// A single's frontend is the project root, not apps/web.
			name: "a single with Next",
			opts: Options{ProjectName: "t", Architecture: ArchSingle, Frontend: FrontendNext},
			want: []string{
				"lib/password-rules.ts",
				"components/password-strength.tsx",
			},
		},
		{
			name: "a single with Vite",
			opts: Options{ProjectName: "t", Architecture: ArchSingle, Frontend: FrontendTanStack},
			want: []string{
				"src/lib/password-rules.ts",
				"src/components/password-strength.tsx",
			},
		},
		{
			name: "the desktop client",
			opts: Options{ProjectName: "t", Architecture: ArchTriple, Frontend: FrontendNext, IncludeDesktop: true},
			want: []string{
				"apps/desktop/frontend/src/lib/password-rules.ts",
				"apps/desktop/frontend/src/components/password-strength.tsx",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := writePasswordStrengthFiles(root, tc.opts); err != nil {
				t.Fatalf("writePasswordStrengthFiles: %v", err)
			}
			for _, rel := range tc.want {
				path := filepath.Join(root, filepath.FromSlash(rel))
				body, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("%s was not written: %v", rel, err)
					continue
				}
				if len(body) == 0 {
					t.Errorf("%s is empty", rel)
				}
			}
		})
	}
}

// A Vite app cannot have "use client" at the top of a module.
//
// Next ignores the directive where it is unnecessary; Vite's React plugin does
// not know it, and esbuild leaves it as an expression statement at module
// scope, which is a lint error in every project that runs one and noise in the
// rest. nextToTanStack exists to take it off, and the thing to check is that
// the writer actually calls it for the Vite shapes.
func TestTheMeterDropsUseClientForVite(t *testing.T) {
	if !strings.Contains(passwordStrengthTSX(), `"use client"`) {
		t.Fatal("the Next.js component has no \"use client\", so this test proves nothing")
	}

	root := t.TempDir()
	opts := Options{ProjectName: "t", Architecture: ArchSingle, Frontend: FrontendTanStack}
	if err := writePasswordStrengthFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(root, "src", "components", "password-strength.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), `"use client"`) {
		t.Error(`the Vite component still carries "use client"`)
	}
	if !strings.Contains(string(body), "export function PasswordStrength") {
		t.Error("the Vite component does not export PasswordStrength")
	}
}

// The ladder reads weakest to strongest, and the strength is separate from the
// rules.
//
// Both halves matter. A bar that fills one segment per rule met is not a
// strength meter, and the account page shipped one for months: "a" fails two
// rules of five, so it lit three segments and showed a single character as more
// than half way there. And a meter with no words on it tells somebody they are
// wrong without telling them what would be right.
func TestTheLadderIsOrderedAndNamed(t *testing.T) {
	src := passwordRulesTS()

	ladder := []string{"Very weak", "Weak", "Fair", "Strong", "Very strong"}
	at := -1
	for _, rung := range ladder {
		next := strings.Index(src, `"`+rung+`"`)
		if next < 0 {
			t.Fatalf("the ladder has no %q rung", rung)
		}
		if next <= at {
			t.Errorf("%q comes before the rung below it, so the ladder is not weakest first", rung)
		}
		at = next
	}

	// The score comes from an entropy estimate, not from counting rules. If it
	// ever goes back to counting rules met, these are the names that disappear.
	for _, needed := range []string{"effectiveLength", "poolSize", "bits"} {
		if !strings.Contains(src, needed) {
			t.Errorf("passwordStrength does not use %s, so it is not measuring guessability", needed)
		}
	}

	// A repeat or a run costs a guesser almost nothing, and the discount for it
	// is the only reason twenty a's does not score as a very strong password.
	if !strings.Contains(src, "0.25") {
		t.Error("nothing discounts a repeated or sequential character, so aaaaaaaaaaaaaaaaaaaa scores as very strong")
	}
}
