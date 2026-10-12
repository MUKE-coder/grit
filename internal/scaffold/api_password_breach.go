package scaffold

// Checking a chosen password against a public breach corpus.
//
// The rules in password.go judge a password on its shape: long enough, mixed
// enough, not one of a couple of hundred obvious ones. Shape is a weak proxy
// for whether a password is actually safe. "Summer2024!" passes every one of
// them and is in every credential-stuffing list in the world, because the rules
// measure the thing that is easy to measure rather than the thing that matters.
//
// What matters is whether this exact password already appears in a breach,
// because that is the list an attacker is working from. Credential stuffing is
// not guessing: it is replaying known passwords. A password that is already on
// that list is compromised on the day it is chosen, however strong it looks.
//
// The check is off by default, because it is a call to a third party on a path
// that must not depend on one, and it fails open for the same reason. It is
// deliberately not applied at sign-in: an existing user whose password turns up
// in tomorrow's breach needs to be told, not locked out mid-session.

// apiPasswordBreachGo emits internal/password/breach.go.
func apiPasswordBreachGo() string {
	return `package password

import (
	"bufio"
	"context"
	"crypto/sha1" // #nosec G505 -- the range API this calls is defined over SHA-1; it is a lookup key here, never a password hash
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// BreachAPI is the range endpoint. Only the first five characters of the hash
// are ever sent to it.
const BreachAPI = "https://api.pwnedpasswords.com/range/"

// breachTimeout bounds the call. A password form that hangs for ten seconds
// because somebody else's service is slow is worse than not checking at all,
// and this runs while a person waits on a submit button.
const breachTimeout = 2 * time.Second

var (
	breachMu      sync.RWMutex
	breachOn      bool
	breachBaseURL = BreachAPI
	breachClient  = &http.Client{Timeout: breachTimeout}
)

// SetBreachCheck turns the check on or off for the process. Off is the default.
func SetBreachCheck(enabled bool) {
	breachMu.Lock()
	defer breachMu.Unlock()
	breachOn = enabled
	if enabled {
		log.Println("password: breach checking is on; chosen passwords are checked against the public corpus")
	}
}

// BreachCheckEnabled reports whether the check is on, so a caller can tell the
// difference between "checked and clean" and "not checked".
func BreachCheckEnabled() bool {
	breachMu.RLock()
	defer breachMu.RUnlock()
	return breachOn
}

// Breached reports whether candidate appears in the public breach corpus.
//
// It answers false whenever it cannot get a definite yes: disabled, no network,
// a slow or broken endpoint, an unreadable response. That is a deliberate
// choice about which failure is worse. Failing closed would mean an outage at
// somebody else's service stops every registration and every password change
// here, which is a self-inflicted incident caused by an advisory check. Failing
// open means that during such an outage passwords are judged by the local rules
// alone, which is exactly where this feature started.
func Breached(ctx context.Context, candidate string) bool {
	breachMu.RLock()
	on, base, client := breachOn, breachBaseURL, breachClient
	breachMu.RUnlock()
	if !on || candidate == "" {
		return false
	}

	// k-anonymity: the password is hashed here, and only the first five
	// characters of that hash leave the process. The endpoint answers with
	// every suffix sharing that prefix, several hundred of them, and the
	// comparison happens locally. The service never learns which one was
	// asked about, and never sees the password.
	sum := sha1.Sum([]byte(candidate)) // #nosec G401 -- lookup key for the range API, not a password hash
	full := strings.ToUpper(hex.EncodeToString(sum[:]))
	prefix, suffix := full[:5], full[5:]

	ctx, cancel := context.WithTimeout(ctx, breachTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+prefix, nil)
	if err != nil {
		log.Printf("password: breach check request for prefix %s: %v", prefix, err)
		return false
	}
	// Pads the response to a fixed size, so the number of results cannot be
	// inferred from the length of an observed response.
	req.Header.Set("Add-Padding", "true")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("password: breach check unavailable, allowing the password: %v", err)
		return false
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("password: breach check answered %d, allowing the password", resp.StatusCode)
		return false
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		tail, count, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(tail), suffix) {
			continue
		}
		// A padded response includes entries with a count of zero, which are
		// filler rather than real matches.
		if strings.TrimSpace(count) == "0" {
			return false
		}
		return true
	}
	if err := scanner.Err(); err != nil {
		log.Printf("password: reading the breach response, allowing the password: %v", err)
	}
	return false
}

// BreachMessage is what to tell somebody whose password is on the list. It says
// what happened and what to do, and does not say how many times it appeared:
// that number is interesting to an engineer and alarming to everybody else.
const BreachMessage = "This password has appeared in a public data breach, so it is already on the lists attackers use. Please choose a different one."

// setBreachEndpoint points the check at another URL. Tests only.
func setBreachEndpoint(url string) func() {
	breachMu.Lock()
	previous := breachBaseURL
	breachBaseURL = url
	breachMu.Unlock()
	return func() {
		breachMu.Lock()
		breachBaseURL = previous
		breachMu.Unlock()
	}
}
`
}

// apiPasswordBreachTestGo emits internal/password/breach_test.go.
//
// The tests run against a local server rather than the real endpoint, so they
// are fast, offline and deterministic. The one thing they must prove is the
// fail-open behaviour, because that is the property the production path
// depends on and the one that is invisible until the day it matters.
func apiPasswordBreachTestGo() string {
	return `package password

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The suffix of SHA-1("password"), whose full hash is
// 5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8.
const passwordSuffix = "1E4C9B93F3F0682250B6CF8331B7EE68FD8"

func TestBreachedIsOffByDefault(t *testing.T) {
	if BreachCheckEnabled() {
		t.Fatal("breach checking should be off until it is switched on")
	}
	if Breached(context.Background(), "password") {
		t.Fatal("a disabled check must not call out, and must answer false")
	}
}

func TestBreachedFindsAKnownPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Add-Padding"); got != "true" {
			t.Errorf("Add-Padding header = %q, want \"true\"", got)
		}
		// Only the prefix may be sent, never the whole hash.
		if len(r.URL.Path) != len("/")+5 {
			t.Errorf("path %q sent more than a five character prefix", r.URL.Path)
		}
		_, _ = w.Write([]byte("0000000000000000000000000000000000A:12\n" + passwordSuffix + ":37359195\n"))
	}))
	defer server.Close()
	defer setBreachEndpoint(server.URL + "/")()

	SetBreachCheck(true)
	defer SetBreachCheck(false)

	if !Breached(context.Background(), "password") {
		t.Fatal("a password present in the corpus should be reported as breached")
	}
	if Breached(context.Background(), "a-passphrase-that-is-not-in-the-list") {
		t.Fatal("a password absent from the response should not be reported as breached")
	}
}

func TestBreachedIgnoresPaddingEntries(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// A padded response carries filler entries with a count of zero.
		_, _ = w.Write([]byte(passwordSuffix + ":0\n"))
	}))
	defer server.Close()
	defer setBreachEndpoint(server.URL + "/")()

	SetBreachCheck(true)
	defer SetBreachCheck(false)

	if Breached(context.Background(), "password") {
		t.Fatal("a zero count is padding, not a match")
	}
}

// The property the production path depends on: when the service is unavailable,
// a password change still works rather than every one of them failing.
func TestBreachedFailsOpen(t *testing.T) {
	cases := map[string]http.HandlerFunc{
		"server error": func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusInternalServerError) },
		"rate limited": func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusTooManyRequests) },
		"nonsense":     func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("not a range response")) },
	}
	for name, handler := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(handler)
			defer server.Close()
			defer setBreachEndpoint(server.URL + "/")()

			SetBreachCheck(true)
			defer SetBreachCheck(false)

			if Breached(context.Background(), "password") {
				t.Fatal("an unusable answer must not block the password")
			}
		})
	}
}

func TestBreachedFailsOpenWhenNothingIsListening(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := server.URL + "/"
	server.Close() // nothing is listening now
	defer setBreachEndpoint(url)()

	SetBreachCheck(true)
	defer SetBreachCheck(false)

	if Breached(context.Background(), "password") {
		t.Fatal("an unreachable endpoint must not block the password")
	}
}
`
}
