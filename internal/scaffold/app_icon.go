package scaffold

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"path/filepath"
	"strconv"
	"strings"
)

// The app icon a generated mobile app ships with.
//
// It used to be four copies of the Grit logo, so an app called cm2 had Grit's
// G on the home screen, on the splash, in the browser tab and in whatever a
// store listing would show. The rule is already written down for the web app
// and the admin panel, which draw the project's own initial in its own accent
// colour, and the mobile app is the surface where getting it wrong is most
// visible: an app icon is what somebody sees before they open anything.
//
// Drawn rather than embedded, because the colour belongs to the project's
// theme. No letter: rendering one needs a font in the binary, and a plain mark
// in the right colour is honestly a placeholder, which is what this is. The
// README says to replace it, and the auth screens draw the initial themselves
// in a view, where text costs nothing.

// iconSize is the square Expo asks for. 1024 is what the stores want and what
// `expo prebuild` downsamples from.
const iconSize = 1024

// appIconPNG draws the placeholder mark for a project on this theme.
func appIconPNG(theme string) ([]byte, error) {
	accent := parseHexColour(themeColor(theme, "accent", "#6c5ce7"))
	// A darker shade of the same hue behind it, so the mark reads as a mark
	// rather than as a missing image.
	back := shade(accent, 0.45)

	img := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	// The rounded square Apple and Android both mask into their own shape.
	// Generous radius: a tighter one looks clipped once the platform mask is
	// applied on top of it.
	outer := float64(iconSize) * 0.22
	// Small enough that the inner shape stays a rounded square. Above half its
	// side it degenerates into a circle, which is a different mark.
	inner := float64(iconSize) * 0.09
	pad := float64(iconSize) * 0.26

	for y := 0; y < iconSize; y++ {
		for x := 0; x < iconSize; x++ {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			if !insideRounded(fx, fy, 0, 0, iconSize, iconSize, outer) {
				img.Set(x, y, color.RGBA{0, 0, 0, 0})
				continue
			}
			if insideRounded(fx, fy, pad, pad, float64(iconSize)-pad, float64(iconSize)-pad, inner) {
				img.Set(x, y, accent)
				continue
			}
			img.Set(x, y, back)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encoding the app icon: %w", err)
	}
	return buf.Bytes(), nil
}

// insideRounded reports whether a point is inside a rounded rectangle.
func insideRounded(x, y, x0, y0, x1, y1, r float64) bool {
	if x < x0 || x > x1 || y < y0 || y > y1 {
		return false
	}
	// Only the four corner squares need the distance test.
	cx, cy := x, y
	switch {
	case x < x0+r:
		cx = x0 + r
	case x > x1-r:
		cx = x1 - r
	default:
		return true
	}
	switch {
	case y < y0+r:
		cy = y0 + r
	case y > y1-r:
		cy = y1 - r
	default:
		return true
	}
	dx, dy := x-cx, y-cy
	return math.Hypot(dx, dy) <= r
}

// shade darkens a colour towards black by f.
func shade(c color.RGBA, f float64) color.RGBA {
	k := 1 - f
	return color.RGBA{
		R: uint8(float64(c.R) * k),
		G: uint8(float64(c.G) * k),
		B: uint8(float64(c.B) * k),
		A: 255,
	}
}

// parseHexColour reads #rgb or #rrggbb, falling back to the Grit accent so a
// malformed theme value cannot stop a scaffold.
func parseHexColour(s string) color.RGBA {
	h := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(h) == 3 {
		h = string([]byte{h[0], h[0], h[1], h[1], h[2], h[2]})
	}
	if len(h) != 6 {
		return color.RGBA{0x6c, 0x5c, 0xe7, 255}
	}
	v, err := strconv.ParseUint(h, 16, 32)
	if err != nil {
		return color.RGBA{0x6c, 0x5c, 0xe7, 255}
	}
	return color.RGBA{uint8(v >> 16), uint8(v >> 8), uint8(v), 255}
}

// writeAppIcons writes the mobile app's icon, splash, adaptive icon and
// favicon, all of them this project's mark rather than the framework's.
func writeAppIcons(dir, theme string, names ...string) error {
	data, err := appIconPNG(theme)
	if err != nil {
		return err
	}
	for _, name := range names {
		if err := writeBytes(filepath.Join(dir, name), data); err != nil {
			return fmt.Errorf("writing app icon %s: %w", name, err)
		}
	}
	return nil
}
