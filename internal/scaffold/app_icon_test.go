package scaffold

import (
	"bytes"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A generated mobile app shipped Grit's brand as its own: icon.png,
// adaptive-icon.png, splash.png and favicon.png were four copies of the Grit
// logo, so an app called cm2 had Grit's G on the home screen, on the splash,
// in the browser tab and in whatever a store listing would show. The auth
// screens rendered the same file, and the sign-up screen told the user to
// "Get started with Grit in seconds".
//
// The rule was already written down for the web app and the admin panel. The
// mobile app is the surface where breaking it is most visible, because an app
// icon is what somebody sees before they open anything.
//
// Found by running a generated Expo app and looking at its login screen.

func TestTheAppIconIsDrawnFromTheTheme(t *testing.T) {
	// Two themes with different accents produce different icons. If they did
	// not, the icon is not coming from the theme.
	one, err := appIconPNG("atlas")
	if err != nil {
		t.Fatal(err)
	}
	two, err := appIconPNG("emerald")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(one, two) {
		t.Error("two themes produced the same icon, so it is not drawn from the theme")
	}

	img, err := png.Decode(bytes.NewReader(one))
	if err != nil {
		t.Fatalf("the icon is not a readable PNG: %v", err)
	}
	if b := img.Bounds(); b.Dx() != iconSize || b.Dy() != iconSize {
		t.Errorf("the icon is %dx%d, not %dx%d, which is what the stores ask for",
			b.Dx(), b.Dy(), iconSize, iconSize)
	}

	// The accent has to actually appear in it.
	accent := parseHexColour(themeColor("atlas", "accent", "#6c5ce7"))
	centre := img.At(iconSize/2, iconSize/2)
	r, g, b, _ := centre.RGBA()
	if uint8(r>>8) != accent.R || uint8(g>>8) != accent.G || uint8(b>>8) != accent.B {
		t.Errorf("the middle of the icon is not the theme's accent: got %v, want %v",
			centre, accent)
	}

	// And the corners are transparent, so the platform's own mask has
	// something to work with.
	if _, _, _, a := img.At(1, 1).RGBA(); a != 0 {
		t.Error("the icon's corner is not transparent, so a platform mask will clip a square")
	}
}

// The Expo scaffold writes its four assets, and they are this project's mark.
func TestTheExpoScaffoldDrawsItsOwnIcons(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Theme: "emerald"}
	if err := writeExpoFiles(root, opts); err != nil {
		t.Fatalf("writeExpoFiles: %v", err)
	}

	want, err := appIconPNG(opts.Theme)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"icon.png", "splash.png", "adaptive-icon.png", "favicon.png"} {
		got := readTestBytes(t, filepath.Join(root, "apps", "expo", "assets", name))
		if bytes.Equal(got, gritLogoPNG) {
			t.Errorf("%s is the Grit logo, so this app ships somebody else's brand", name)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s is not the icon drawn for this project's theme", name)
		}
	}

	// And a note beside them saying they are placeholders.
	readme := string(readTestBytes(t, filepath.Join(root, "apps", "expo", "assets", "README.md")))
	if !strings.Contains(readme, "placeholders") {
		t.Error("nothing beside the icons says they are placeholders to replace")
	}
}

// readTestBytes reads a generated file or fails the test.
func readTestBytes(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return data
}

// And the auth screens draw this project's initial rather than an image of
// somebody else's mark.
func TestTheExpoAuthScreensShowTheProject(t *testing.T) {
	for name, src := range map[string]string{
		"login":    expoLoginScreen(),
		"register": expoRegisterScreen(),
	} {
		if strings.Contains(src, `require("../../assets/icon.png")`) {
			t.Errorf("the %s screen renders assets/icon.png, which is a placeholder, not "+
				"an identity", name)
		}
		if !strings.Contains(src, "APP_INITIAL") {
			t.Errorf("the %s screen does not show the project's initial", name)
		}
		if !strings.Contains(src, `import { APP_INITIAL } from "@/lib/brand"`) {
			t.Errorf("the %s screen does not import APP_INITIAL", name)
		}
	}

	if strings.Contains(expoRegisterScreen(), "with Grit in seconds") {
		t.Error("the sign-up screen still advertises the framework to the app's own users")
	}
}

// And nothing else in the app talks about the framework to the app's users.
//
// The home screen's quick-start card read "Your Grit mobile app is connected
// to the API", which is a sentence about the tooling shown to somebody using
// the product.
func TestNoExpoScreenNamesTheFrameworkToItsUsers(t *testing.T) {
	screens := map[string]string{
		"home":     expoHomeScreen(),
		"explore":  expoExploreScreen(),
		"profile":  expoProfileScreen(),
		"settings": expoSettingsScreen(),
		"login":    expoLoginScreen(),
		"register": expoRegisterScreen(),
	}
	for name, src := range screens {
		// A comment or an import path may legitimately carry the word, and a
		// JSX block comment spans several lines, so the state has to be
		// tracked rather than tested per line.
		inBlock := false
		for _, line := range strings.Split(src, "\n") {
			trimmed := strings.TrimSpace(line)
			opened := strings.Contains(trimmed, "{/*") || strings.Contains(trimmed, "/*")
			closed := strings.Contains(trimmed, "*/")
			wasInBlock := inBlock
			if opened && !closed {
				inBlock = true
			} else if closed {
				inBlock = false
			}
			if wasInBlock || opened {
				continue
			}
			if !strings.Contains(trimmed, "Grit") {
				continue
			}
			if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "import ") ||
				strings.Contains(trimmed, "@/lib/") || strings.Contains(trimmed, "apps/expo") {
				continue
			}
			t.Errorf("the %s screen shows the framework's name to the app's own users: %s",
				name, trimmed)
		}
	}
}

// lib/brand.ts reads the name rather than having it spliced in, so renaming
// the project in app.json renames it on screen.
func TestTheBrandModuleReadsTheNameFromAppJSON(t *testing.T) {
	src := expoBrand()
	if !strings.Contains(src, "Constants.expoConfig?.name") {
		t.Error("lib/brand.ts does not read the app's name from app.json")
	}
	if !strings.Contains(src, "export const APP_INITIAL") {
		t.Error("lib/brand.ts exports no APP_INITIAL")
	}
}

// The splash and adaptive-icon backgrounds follow the theme too.
//
// They were Grit's near-black on every project. The default theme, atlas, has
// a white background, so a new app launched by flashing a colour from a
// palette it does not use.
func TestTheSplashBackgroundFollowsTheTheme(t *testing.T) {
	for _, theme := range []string{"atlas", "emerald", "mono"} {
		src := expoAppJSON(Options{ProjectName: "app", Theme: theme})
		want := themeColor(theme, "bg-primary", "")
		if want == "" {
			t.Fatalf("the %s theme declares no bg-primary", theme)
		}
		if strings.Count(src, `"backgroundColor": "`+want+`"`) != 2 {
			t.Errorf("app.json does not use the %s theme's background (%s) for both the "+
				"splash and the adaptive icon", theme, want)
		}
	}

	// And Grit's own near-black is gone from every one of them.
	for _, theme := range []string{"atlas", "emerald", "sky"} {
		src := expoAppJSON(Options{ProjectName: "app", Theme: theme})
		if themeColor(theme, "bg-primary", "") != "#0a0a0f" && strings.Contains(src, "#0a0a0f") {
			t.Errorf("the %s theme's app.json still carries Grit's background", theme)
		}
	}
}
