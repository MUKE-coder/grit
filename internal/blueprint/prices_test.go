package blueprint

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// The website's copy of the prices. The CLI prints them, the checkout charges
// them, and neither can import the other.
const plansFile = "docs/config/blueprint-plans.ts"

// TestBlueprintPricesMatchTheWebsite compares the two.
//
// `grit blueprint info instagram` says $10. The checkout route takes its
// amount from blueprint-plans.ts. Change one and the other keeps saying the
// old number, with nothing anywhere complaining: the first report is a
// customer with a screenshot of a price that was never charged.
func TestBlueprintPricesMatchTheWebsite(t *testing.T) {
	page, ok := readRepoFile(t, plansFile)
	if !ok {
		t.Skipf("%s not found; this test only runs inside the Grit repository", plansFile)
	}

	for _, c := range []struct {
		constant string
		goValue  int
	}{
		{"SINGLE_USD", PriceSingleUSD},
		{"MEMBERSHIP_MONTHLY_USD", PriceMembershipMonthlyUSD},
		{"MEMBERSHIP_YEARLY_USD", PriceMembershipYearlyUSD},
		{"LIFETIME_USD", PriceLifetimeUSD},
	} {
		want, found := readConst(page, c.constant)
		if !found {
			t.Errorf("%s no longer exports %s, so the CLI's prices are comparing "+
				"against nothing. Fix the reader here, or the two drift silently.",
				plansFile, c.constant)
			continue
		}
		if want != c.goValue {
			t.Errorf("%s is $%d in %s and $%d in registry.go",
				c.constant, want, plansFile, c.goValue)
		}
	}
}

// TestPaidBlueprintsHaveAPageToBuyThemOn: a paid blueprint whose page does not
// exist is a dead end at the moment somebody wants to pay.
func TestPaidBlueprintsHaveAPageToBuyThemOn(t *testing.T) {
	page, ok := readRepoFile(t, blueprintsPage)
	if !ok {
		t.Skipf("%s not found; this test only runs inside the Grit repository", blueprintsPage)
	}

	for _, b := range All() {
		if b.Access == Planned {
			continue
		}
		want := Site + "/" + b.Name
		if got := b.Page(); got != want {
			t.Errorf("%s: Page() is %s, expected %s", b.Name, got, want)
		}
		// The page renders from the website's list, so being in that list is
		// what makes the URL resolve. Being absent is a 404 on a page the CLI
		// prints, which is the failure this catches.
		if !strings.Contains(page, fmt.Sprintf("name: '%s'", b.Name)) {
			t.Errorf("%s prints %s, but %s does not list it, so that URL 404s",
				"grit blueprint info "+b.Name, b.Page(), blueprintsPage)
		}
	}
}

// TestEveryAvailableBlueprintSaysHowToRunIt: a detail page that charges $10
// and then does not say what to type is worse than no page.
func TestEveryAvailableBlueprintSaysHowToRunIt(t *testing.T) {
	page, ok := readRepoFile(t, blueprintsPage)
	if !ok {
		t.Skipf("%s not found; this test only runs inside the Grit repository", blueprintsPage)
	}

	for _, b := range All() {
		if b.Access == Planned {
			continue
		}
		block, found := entryBlock(page, b.Name)
		if !found {
			t.Errorf("could not find the %s entry in %s", b.Name, blueprintsPage)
			continue
		}
		for _, field := range []string{"install: [", "requires: [", "signIn: ["} {
			if !strings.Contains(block, field) {
				t.Errorf("%s has no %s in %s, so its page cannot say how to run it",
					b.Name, strings.TrimSuffix(field, " ["), blueprintsPage)
			}
		}
		if strings.Contains(block, "install: []") {
			t.Errorf("%s is available but its install steps are empty", b.Name)
		}
	}
}

var constPattern = `export const %s = ([0-9]+)`

func readConst(src, name string) (int, bool) {
	m := regexp.MustCompile(fmt.Sprintf(constPattern, regexp.QuoteMeta(name))).FindStringSubmatch(src)
	if len(m) != 2 {
		return 0, false
	}
	var out int
	if _, err := fmt.Sscanf(m[1], "%d", &out); err != nil {
		return 0, false
	}
	return out, true
}

// entryBlock returns the text of one blueprint's object, from its name field
// to the start of the next one.
func entryBlock(src, name string) (string, bool) {
	start := strings.Index(src, fmt.Sprintf("name: '%s'", name))
	if start < 0 {
		return "", false
	}
	rest := src[start+len(name):]
	if next := strings.Index(rest, "\n  {\n"); next >= 0 {
		return rest[:next], true
	}
	return rest, true
}
