package scaffold

import (
	"strings"
	"testing"
)

// A new Expo API client refreshes once for every request waiting on it.
func TestExpoAPIClientSharesOneRefresh(t *testing.T) {
	src := expoAPIClient()
	for _, want := range []string{expoRefreshMethod, expoRefreshBlockNew, "this.refreshing ??="} {
		if !strings.Contains(src, want) {
			t.Errorf("the Expo API client is missing %q", strings.SplitN(want, "\n", 2)[0])
		}
	}
	if strings.Contains(src, expoRefreshBlockOld) {
		t.Error("the Expo API client still refreshes once per failed request")
	}
	if n := strings.Count(src, "/auth/refresh`"); n != 1 {
		t.Errorf("the refresh endpoint is called from %d places, want 1", n)
	}
}

func TestExpoRefreshRepair(t *testing.T) {
	now := expoAPIClient()
	before := strings.Replace(strings.Replace(now, expoRefreshBlockNew, expoRefreshBlockOld, 1), expoRefreshMethod, "", 1)
	if before == now {
		t.Fatal("could not rebuild the old client")
	}
	got, fixed, warn := repairExpoRefreshSource(before)
	if got != now || len(fixed) != 1 || len(warn) != 0 {
		t.Fatalf("repair did not produce the current client (fixed %v, warned %v)", fixed, warn)
	}
	if again, fixed, _ := repairExpoRefreshSource(got); again != got || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
	custom := "fetch(API_URL + \"/auth/refresh\")\n"
	if out, _, warn := repairExpoRefreshSource(custom); out != custom || len(warn) != 1 {
		t.Errorf("a client Grit did not write was changed, or not warned about: %v", warn)
	}
}
