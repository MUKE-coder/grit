package scaffold

import (
	"strings"
	"testing"
)

// The web app's client refreshes an expired session once, shared by every
// request that failed with it, and never for the auth endpoints themselves.
func TestWebAPIClientRefreshesAnExpiredSession(t *testing.T) {
	src := webAPIClient()
	for _, want := range []string{
		`api.interceptors.response.use(`,
		`error.response?.status !== 401`,
		`config._retried = true;`,
		`refreshing ??= api`,
		`.post("/api/auth/refresh")`,
		`/\/auth\/(login|register|refresh|logout)/`,
		`return api.request(config);`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("lib/api.ts lacks %q", want)
		}
	}
	if strings.Count(src, `.post("/api/auth/refresh")`) != 1 {
		t.Error("more than one refresh call: failed requests must share one")
	}
}
