package scaffold

import (
	"strings"
	"testing"
)

// The seeder wrote NEXT_PUBLIC_ names into every frontend, and only one kind
// of frontend reads them.
//
//	Next   reads NEXT_PUBLIC_*   and got the right file
//	Vite   reads VITE_*          and got a file nothing read
//	Expo   reads EXPO_PUBLIC_*   and got no file at all
//
// A TanStack project's apps/admin/.env.local therefore named a publishable key
// the app could not see and an API address it ignored, and an Expo project had
// neither. Each symptom points somewhere else: a public endpoint answering
// INVALID_API_KEY looks like a key problem, and a phone calling a port nothing
// answers looks like the phone.
//
// Found on the Expo tier of the end-to-end sweep, and verified on three
// generated projects: Vite got VITE_API_URL, Expo got EXPO_PUBLIC_API_PORT,
// Next was unchanged.

func TestTheSeederAsksWhatEachFrontendIs(t *testing.T) {
	src := apiAPIKeySeederGo()

	for _, want := range []string{
		`prefix := clientEnvPrefix(dir)`,
		`func clientEnvPrefix(dir string) string {`,
		`return "VITE_"`,
		`return "EXPO_PUBLIC_"`,
		`return "NEXT_PUBLIC_"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the seeder does not ask which frontend it is writing for: missing %q", want)
		}
	}

	// The prefix has to reach the names it writes, not just be computed.
	if strings.Contains(src, `"NEXT_PUBLIC_API_KEY=" + publishable`) {
		t.Error("the key line still hardcodes NEXT_PUBLIC_, so a Vite or Expo app gets a " +
			"key under a name it cannot read")
	}
	if !strings.Contains(src, `prefix + "API_KEY=" + publishable`) {
		t.Error("the key line does not use the prefix")
	}
}

// An Expo app is in the list at all, which it was not.
func TestTheSeederWritesAnEnvFileForExpo(t *testing.T) {
	src := apiAPIKeySeederGo()
	if !strings.Contains(src, `filepath.Join("..", "expo", ".env.local")`) {
		t.Error("apps/expo is not among the frontends the seeder writes to, so an Expo " +
			"project gets no publishable key and no API address")
	}
}

// And it gets a port rather than a URL, because a phone cannot reach
// localhost: that address is the phone.
func TestAnExpoAppIsGivenThePortNotAURL(t *testing.T) {
	src := apiAPIKeySeederGo()

	if !strings.Contains(src, `prefix + "API_PORT=" + apiPort()`) {
		t.Error("the Expo file does not name the API's port. An absolute localhost URL " +
			"would break the LAN derivation the app uses on a real device.")
	}
	if !strings.Contains(src, "func apiPort() string {") {
		t.Error("there is no apiPort helper for the Expo file to use")
	}
}

// The Expo client reads that port instead of a literal.
func TestTheExpoClientPortFollowsTheProject(t *testing.T) {
	src := expoAPIClient()

	if strings.Contains(src, "const API_PORT = 8080;") {
		t.Error("the Expo API client still hardcodes port 8080, so moving APP_PORT leaves " +
			"the phone calling a port nothing answers on all three resolution paths")
	}
	if !strings.Contains(src, "process.env.EXPO_PUBLIC_API_PORT") {
		t.Error("the Expo API client does not read EXPO_PUBLIC_API_PORT")
	}
	// The LAN derivation is the reason the port is configured separately, so
	// it has to still be there.
	if !strings.Contains(src, "hostUri") {
		t.Error("the Expo client no longer derives the host from Metro, which is what " +
			"makes it work on a real device")
	}
}
