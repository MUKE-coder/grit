package plugin

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
)

func init() { Register(pushPlugin()) }

//go:embed push/*.tmpl
var pushTemplates embed.FS

// pushPlugin sends push notifications to the Expo app.
//
// Design:
//
//   - One provider. Expo's push service takes one token format for iOS and
//     Android and talks to APNs and FCM itself, so the API holds no Apple
//     certificate or Firebase key and makes one kind of request.
//
//   - A token is unique and belongs to whoever registered it last. A phone that
//     signs in as someone else moves to them; otherwise it would keep getting the
//     previous person's messages.
//
//   - Dead tokens clean themselves up. Expo answers DeviceNotRegistered for an
//     app that was deleted or a user who turned notifications off, and that
//     token is removed at once instead of being tried forever.
//
//   - Sending never blocks a request. Push.Go sends in the background; Send is
//     there for jobs and tests that want the result.
//
// Found missing while building the WhatsApp blueprint, whose plan listed push
// as already shipped.
func pushPlugin() Plugin {
	return Plugin{
		Name:    "push",
		Version: "1.0.0",
		Summary: "Push notifications to the Expo app, through Expo's push service",
		Description: `Sends push notifications to your Expo app on iOS and Android.

  • POST /push/tokens registers the signed-in device; /push/tokens/remove
    unregisters it on sign-out; POST /push/test sends yourself a test
  • services.Push sends to every device of a list of users, 100 per request,
    in the background with Push.Go
  • Tokens Expo reports as no longer registered are deleted automatically
  • A phone that signs in as someone else stops getting the last user's pushes
  • apps/expo/lib/push.ts asks permission, registers, unregisters, and tells
    you which notification was tapped

No certificates or Firebase keys in the API: Expo relays to APNs and FCM.`,

		NextSteps: []string{
			"Run the migration:   grit migrate",
			"Install the package: pnpm install",
			"Register after sign-in in the Expo app:  import { registerForPush } from \"@/lib/push\"",
			"Send from Go:        services.NewPush(db).Go([]string{userID}, services.PushMessage{Title: \"Hi\", Body: \"...\"})",
			"Try it:              POST /api/v1/push/test from a signed-in device with a real token",
			"Real pushes need a physical device (not a simulator) and, for store builds, an EAS project id",
		},

		NodeDeps: []Dependency{
			// The version Expo SDK 54 bundles (expo/bundledNativeModules.json).
			{Name: "expo-notifications", Version: "~0.32.17", Workspace: "apps/expo"},
		},

		Files:      pushFiles,
		Injections: pushInjections,
	}
}

func pushTemplate(ctx Context, name string) string {
	raw, err := pushTemplates.ReadFile("push/" + name)
	if err != nil {
		// The templates are compiled in; a missing one is a build mistake.
		panic("push plugin template missing: " + name)
	}
	return strings.ReplaceAll(string(raw), "{{MODULE}}", ctx.Module)
}

func hasExpo(ctx Context) bool {
	info, err := os.Stat(filepath.Join(ctx.Root, "apps", "expo"))
	return err == nil && info.IsDir()
}

func pushFiles(ctx Context) map[string]string {
	files := map[string]string{
		apiPath(ctx, "internal/models/push_token.go"):  pushTemplate(ctx, "push_token.go.tmpl"),
		apiPath(ctx, "internal/services/push.go"):      pushTemplate(ctx, "push_service.go.tmpl"),
		apiPath(ctx, "internal/services/push_test.go"): pushTemplate(ctx, "push_service_test.go.tmpl"),
		apiPath(ctx, "internal/handlers/push.go"):      pushTemplate(ctx, "push_handler.go.tmpl"),
	}
	if hasExpo(ctx) {
		files["apps/expo/lib/push.ts"] = pushTemplate(ctx, "expo_push.ts.tmpl")
	}
	return files
}

func pushInjections(ctx Context) []Injection {
	return []Injection{
		{
			File:   apiPath(ctx, "internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&PushToken{},",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code:   "\tpushHandler := handlers.NewPushHandler(db)",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: "\t\t// Push notifications: this device's token, for the signed-in user.\n" +
				"\t\tprotected.POST(\"/push/tokens\", pushHandler.Register)\n" +
				"\t\tprotected.POST(\"/push/tokens/remove\", pushHandler.Unregister)\n" +
				"\t\tprotected.POST(\"/push/test\", pushHandler.Test)",
		},
	}
}
