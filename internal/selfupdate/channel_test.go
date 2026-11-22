package selfupdate

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// withChannelAt points ChannelURL at a test server for one test.
func withChannelAt(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	was := channelURL
	channelURL = srv.URL
	t.Cleanup(func() { channelURL = was })
}

func TestResolveChannelReadsTheStableVersion(t *testing.T) {
	withChannelAt(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"stable":"3.527.0","edge":"3.566.0"}`))
	})

	got, err := ResolveChannel("stable")
	if err != nil {
		t.Fatalf("ResolveChannel: %v", err)
	}
	if got != "3.527.0" {
		t.Errorf("got %q, want 3.527.0", got)
	}
}

// edge is the default and must not touch the network.
//
// A channel file that is unreachable has to be unable to break `grit update`,
// which is the command people run when something is already wrong.
func TestEdgeNeedsNoChannelFile(t *testing.T) {
	withChannelAt(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("edge asked for the channel file")
		http.Error(w, "should not be called", http.StatusInternalServerError)
	})

	for _, name := range []string{"", "edge", "EDGE", "latest"} {
		got, err := ResolveChannel(name)
		if err != nil {
			t.Errorf("ResolveChannel(%q): %v", name, err)
		}
		if got != "" {
			t.Errorf("ResolveChannel(%q) = %q, want the empty string, which "+
				"means the latest release", name, got)
		}
	}
}

func TestAnUnknownChannelIsRefusedByName(t *testing.T) {
	_, err := ResolveChannel("nightly")
	if err == nil {
		t.Fatal("an unknown channel was accepted")
	}
	if !strings.Contains(err.Error(), "nightly") {
		t.Errorf("the error does not name the channel asked for: %v", err)
	}
}

// A channel file that says nothing is not a version.
//
// Silently installing the empty string means installing the latest, which is
// the opposite of what somebody asking for stable wants.
func TestAnEmptyChannelFileIsAnError(t *testing.T) {
	for _, body := range []string{`{}`, `{"stable":""}`, `{"stable":"  "}`} {
		withChannelAt(t, func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(body))
		})
		if _, err := ResolveChannel("stable"); err == nil {
			t.Errorf("%s was accepted as a stable channel", body)
		}
	}
}

func TestAChannelFileThatIsNotJSONSaysSo(t *testing.T) {
	withChannelAt(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("<html>404</html>"))
	})
	_, err := ResolveChannel("stable")
	if err == nil {
		t.Fatal("HTML was accepted as a channel file")
	}
	if !strings.Contains(err.Error(), "JSON") {
		t.Errorf("the error does not say what was wrong: %v", err)
	}
}

func TestANon200IsAnError(t *testing.T) {
	withChannelAt(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusServiceUnavailable)
	})
	if _, err := ResolveChannel("stable"); err == nil {
		t.Fatal("a 503 was accepted as a channel file")
	}
}

// withChannelBytes serves a fixed body as the channel file.
func withChannelBytes(t *testing.T, body []byte) {
	t.Helper()
	withChannelAt(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(body)
	})
}
