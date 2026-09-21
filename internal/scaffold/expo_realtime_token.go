package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// React Native has no cookie jar, so the Expo realtime client authenticates
// the socket with a token it is given. It was given one by setRealtimeToken,
// which nothing called: the getter returned null, the handshake carried no
// token, and every realtime subscription in an Expo app was refused. Found
// building the WhatsApp blueprint's mobile app. The getter now reads the same
// SecureStore key the API client keeps the access token under, and
// setRealtimeToken stays for an app that stores it elsewhere.
const (
	expoTokenGetterOld = "let tokenGetter: () => string | null | Promise<string | null> = () => null;\n"
	expoTokenGetterNew = `// The access token the API client keeps in SecureStore (lib/api.ts), read on
// every connect so a refreshed token is used, not the one from sign-in.
let tokenGetter: () => string | null | Promise<string | null> = () => SecureStore.getItemAsync("access_token");
`
	// The app's own wrapper, not expo-secure-store: it falls back to
	// localStorage on web, where SecureStore throws. lib/api.ts uses it too.
	expoSecureStoreImport = "import * as SecureStore from \"@/lib/secure-store\";\n\n"
)

// repairExpoRealtimeToken gives an existing Expo app's realtime client the
// token it never had.
func repairExpoRealtimeToken(root string) error {
	path := filepath.Join(root, "apps", "expo", "lib", "realtime.ts")
	if !fileExists(path) {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairTextFile(root, m, path, repairExpoRealtimeTokenSource)
}

func repairExpoRealtimeTokenSource(src string) (string, []string, []string) {
	if strings.Contains(src, expoTokenGetterNew) {
		return src, nil, nil
	}
	if strings.Count(src, expoTokenGetterOld) != 1 {
		if strings.Contains(src, "tokenGetter") && strings.Contains(src, "() => null") {
			return src, nil, []string{"the realtime client is not the one Grit wrote: make its token getter read the access token (SecureStore \"access_token\"), or the socket is refused and nothing live reaches the app"}
		}
		return src, nil, nil
	}
	out := strings.Replace(src, expoTokenGetterOld, expoTokenGetterNew, 1)
	if !strings.Contains(out, `from "@/lib/secure-store"`) {
		out = expoSecureStoreImport + out
	}
	return out, []string{"the realtime socket signs in with the access token, so live updates reach the app"}, nil
}
