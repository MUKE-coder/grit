package doctor

import (
	"os"
	"strings"
	"testing"
)

// Every check the binary runs is listed on the page that claims to list them.
//
// The page said "17 checks" while the binary ran 25, and twelve of them were
// undocumented, including two added the week before. A test on either side
// alone cannot see that: the registry is complete by construction and the page
// is valid TSX whatever it omits. So this reads both.
//
// It reads the page as text rather than parsing TSX. The only thing it needs is
// whether the check's name appears in it, and the name is quoted there exactly
// as it is here.
func TestEveryCheckIsOnTheDoctorPage(t *testing.T) {
	const page = "../../docs/app/docs/security/doctor/page.tsx"
	body, err := os.ReadFile(page)
	if err != nil {
		// The docs site is in this repository. If it has moved, this test has
		// to be pointed at the new path rather than quietly passing.
		t.Fatalf("reading %s: %v", page, err)
	}
	text := string(body)

	var missing []string
	for _, c := range checks {
		if !strings.Contains(text, "'"+c.name+"'") {
			missing = append(missing, c.name)
		}
	}
	if len(missing) > 0 {
		t.Errorf("%d of %d checks are not on the doctor page: %s\n\n"+
			"    Add each one to the checks array in %s, with the level it reports\n"+
			"    and one sentence on the mistake it finds. A check nobody can read\n"+
			"    about is a check whose warning gets ignored.",
			len(missing), len(checks), strings.Join(missing, ", "), page)
	}

	// And the other direction: a name on the page that the binary no longer
	// runs describes a check that cannot fire.
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "name: '") {
			continue
		}
		name := strings.TrimSuffix(strings.TrimPrefix(line, "name: '"), "',")
		found := false
		for _, c := range checks {
			if c.name == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("the doctor page lists %q, which the binary does not run", name)
		}
	}
}

// The count in the page's sample output is the count the binary prints.
func TestTheDoctorPageQuotesTheRightCheckCount(t *testing.T) {
	const page = "../../docs/app/docs/security/doctor/page.tsx"
	body, err := os.ReadFile(page)
	if err != nil {
		t.Fatalf("reading %s: %v", page, err)
	}
	want := itoa(len(checks)) + " checks over"
	if !strings.Contains(string(body), want) {
		t.Errorf("the page's sample output does not say %q; it is a transcript of a "+
			"command, so a stale number there is a page that disagrees with the terminal "+
			"beside it", want)
	}
}

// itoa, so this file needs no import for one number.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
