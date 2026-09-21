package scaffold

import (
	"strings"
	"testing"
)

// The Expo realtime client finds the access token by itself, and the web one
// is untouched: a browser sends its cookie and has no SecureStore.
func TestExpoRealtimeClientReadsTheAccessToken(t *testing.T) {
	native := realtimeClientTS(true)
	if !strings.HasPrefix(native, expoSecureStoreImport) {
		t.Error("the Expo client does not import expo-secure-store first")
	}
	if !strings.Contains(native, expoTokenGetterNew) {
		t.Error("the Expo client's token getter does not read SecureStore")
	}
	if strings.Contains(native, expoTokenGetterOld) {
		t.Error("the Expo client still has the getter that returns null")
	}
	web := realtimeClientTS(false)
	if strings.Contains(web, "SecureStore") || strings.Contains(web, "tokenGetter") {
		t.Error("the browser client picked up the native token code")
	}
}

func TestExpoRealtimeTokenRepair(t *testing.T) {
	now := realtimeClientTS(true)
	before := strings.Replace(strings.TrimPrefix(now, expoSecureStoreImport), expoTokenGetterNew, expoTokenGetterOld, 1)
	got, fixed, warn := repairExpoRealtimeTokenSource(before)
	if got != now || len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("repair did not produce the current client (fixed %v, warned %v)", fixed, warn)
	}
	if again, fixed, _ := repairExpoRealtimeTokenSource(got); again != got || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
	custom := "let tokenGetter = () => null; // mine\n"
	if out, _, warn := repairExpoRealtimeTokenSource(custom); out != custom || len(warn) != 1 {
		t.Errorf("a client Grit did not write was changed, or not warned about: %v", warn)
	}
}
