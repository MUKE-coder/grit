package plugin

func init() { Register(devicePairingPlugin()) }

// devicePairingPlugin signs a browser in by showing a QR code that an
// already-signed-in phone scans and approves. The WhatsApp Web flow, and how
// Telegram Web, Discord, Steam and most TV apps onboard a second screen.
//
// Design:
//
//   - The browser is anonymous until approval, so the pairing code is the only
//     thing standing between an attacker and somebody else's account. It is 32
//     bytes from crypto/rand, single use, and dead after two minutes.
//
//   - Single use is enforced with a conditional UPDATE, not a read-then-write.
//     Two polls arriving together would otherwise both see "approved, unclaimed"
//     and both walk away with a token pair.
//
//   - Approval is two steps on purpose. The phone fetches what it is about to
//     approve (the browser's user agent, IP and city-less location hint) and
//     shows it before asking. Approving something you cannot see is not consent,
//     and it is the difference between a shoulder-surfed QR being useless and
//     being a silent account takeover.
//
//   - Deny is a first-class endpoint, not an absence. "That wasn't me" has to be
//     one tap, and it burns the code immediately rather than leaving it live for
//     the rest of its TTL.
//
//   - There is no Device model. A paired browser gets an ordinary session row,
//     so it shows up in Active Sessions and is revoked by the machinery that is
//     already there. A parallel device table would be a second list of the same
//     thing, drifting.
func devicePairingPlugin() Plugin {
	return Plugin{
		Name:    "device-pairing",
		Version: "1.0.0",
		Summary: "Sign in a browser by scanning a QR code from a signed-in phone",
		Description: `Adds the WhatsApp Web flow: a browser shows a QR code, an already-signed-in
device scans it and approves, and the browser is signed in.

  • POST /pair/start issues a single-use code and a ready-to-render QR PNG
  • The browser polls and collects a real token pair on approval
  • The approving device sees the user agent and IP it is about to trust
    BEFORE approving, and can deny in one tap
  • A paired browser becomes an ordinary session, so it appears in Active
    Sessions and is revoked from there
  • Codes are 32 bytes of crypto/rand, single use, and expire in two minutes
  • Rate limited per IP, because the poll endpoint is anonymous

Adds a /link page to the web app and a Link a device screen to the admin.`,

		NextSteps: []string{
			"Run the migration:  grit migrate",
			"Rebuild the API:    cd apps/api && go build ./...",
			"Open http://localhost:3000/link to see the QR",
			"Approve it from the admin at System -> Link a device",
			"Paired browsers appear under Account -> Security as active sessions",
		},

		GoDeps: []Dependency{
			// Already present for TOTP enrolment in most projects; declared so
			// an --api project that never enabled 2FA still resolves.
			{Name: "github.com/skip2/go-qrcode", Version: "v0.0.0-20200617195104-da1b6568686e"},
		},

		Files:      devicePairingFiles,
		Injections: devicePairingInjections,
	}
}

// apiPath resolves a path inside the Go API for either project shape.
func apiPath(ctx Context, rel string) string {
	if ctx.Architecture == "single" {
		return rel
	}
	return "apps/api/" + rel
}

func hasAdmin(ctx Context) bool {
	return ctx.Architecture == "triple" || ctx.Architecture == "full"
}

func hasWeb(ctx Context) bool {
	return ctx.Architecture == "double" || ctx.Architecture == "triple" || ctx.Architecture == "full"
}

func devicePairingFiles(ctx Context) map[string]string {
	files := map[string]string{
		apiPath(ctx, "internal/models/pairing_request.go"):       pairingModelGo(ctx),
		apiPath(ctx, "internal/handlers/device_pairing.go"):      pairingHandlerGo(ctx),
		apiPath(ctx, "internal/handlers/device_pairing_test.go"): pairingHandlerTestGo(ctx),
	}

	// The QR lives on the device being paired, which is the browser.
	if hasWeb(ctx) {
		files["apps/web/app/link/page.tsx"] = pairingWebPage()
		files["apps/web/hooks/use-pairing.ts"] = pairingWebHook()
	}

	// The approving surface. A phone is the usual one, but somebody has to be
	// able to approve from a desktop admin too, so this takes a typed code as
	// well as a scanned one.
	if hasAdmin(ctx) {
		files[adminFile(ctx, "hooks/use-device-pairing.ts")] = adminSource(ctx, pairingAdminHook())

		page, route := adminPageFile(ctx, "system/link-device", "system-link-device")
		files[page] = adminSource(ctx, pairingAdminPage())
		if route != "" {
			files[route] = tanStackRouteShim("system/link-device", "@/pages/system-link-device")
		}
	}

	return files
}

func devicePairingInjections(ctx Context) []Injection {
	injections := []Injection{
		{
			File:   apiPath(ctx, "internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&PairingRequest{},",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code:   "\tpairingHandler := handlers.NewPairingHandler(db, authService)",
		},
		{
			// The anonymous half. grit:routes:custom sits at function level with
			// v1 in scope, and v1 carries no auth middleware: the browser has no
			// account yet and no API key, so neither the protected nor the
			// public group can serve it.
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:custom",
			Code: "\t// Device pairing. Anonymous by necessity: the browser being paired\n" +
				"\t// has no session yet. Rate limited inside the handler.\n" +
				"\tv1.POST(\"/pair/start\", pairingHandler.Start)\n" +
				"\tv1.GET(\"/pair/:code\", pairingHandler.Status)",
		},
		{
			// The approving half, from a device that is already signed in.
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: "\t\tprotected.GET(\"/pair/:code/request\", pairingHandler.Describe)\n" +
				"\t\tprotected.POST(\"/pair/:code/approve\", pairingHandler.Approve)\n" +
				"\t\tprotected.POST(\"/pair/:code/deny\", pairingHandler.Deny)",
		},
	}

	if hasAdmin(ctx) {
		admin := adminDir(ctx)
		injections = append(injections,
			// getIcon falls back to FileText for a key it does not know, so
			// without these two the nav entry renders a document icon and
			// nothing anywhere complains.
			Injection{
				File:   admin + "/lib/icons.ts",
				Marker: "// grit:icons:import",
				Code:   "  QrCode,",
			},
			Injection{
				File:   admin + "/lib/icons.ts",
				Marker: "// grit:icons:map",
				Code:   "  QrCode,",
			},
			Injection{
				File:   admin + "/components/chrome/CollapsibleSidebar.tsx",
				Marker: "// grit:nav:system",
				// adminOnly is false: linking your own device is a personal action, not
				// an administrative one, and an EDITOR signed into the admin has the
				// same need for it as an ADMIN.
				Code: `  { href: "/system/link-device", label: "Link a device", iconKey: "QrCode", adminOnly: false },`,
			},
		)
	}

	return injections
}
