package scaffold

import (
	"strings"
	"testing"
)

// grit#147: the admin's 401 interceptor had two paths that replay a request
// and marked only one of them.
//
// A request that queues behind an in-flight refresh is replayed too. Without
// the mark, a second 401 on that replay reads as a first attempt and starts
// the refresh again: a refresh storm rather than one clean sign-out, which is
// what several tabs on a dead session produce.
//
// The desktop client already had it right, setting the flag before the branch,
// so this is the two clients agreeing rather than a rule invented here.
func TestEveryRefreshInterceptorMarksBeforeItBranches(t *testing.T) {
	for _, client := range []struct {
		name   string
		source string
	}{
		{"admin", adminAPIClient()},
		{"desktop", desktopClientApiClientTS()},
	} {
		mark := strings.Index(client.source, "_retry = true")
		branch := strings.Index(client.source, "if (isRefreshing)")
		if mark < 0 {
			t.Errorf("%s: the interceptor never sets _retry, so a retry is indistinguishable from a first attempt", client.name)
			continue
		}
		if branch < 0 {
			t.Errorf("%s: no isRefreshing branch; this test is reading the wrong thing", client.name)
			continue
		}
		if mark > branch {
			t.Errorf("%s: _retry is set after the isRefreshing branch, so a request that "+
				"queues behind a refresh is replayed unmarked and a second 401 on it "+
				"starts the refresh again (grit#147)", client.name)
		}
	}
}
