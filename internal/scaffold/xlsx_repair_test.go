package scaffold

import (
	"strings"
	"testing"
)

func TestVulnerableXLSXMessage(t *testing.T) {
	for src, vulnerable := range map[string]bool{
		`{"dependencies": {"xlsx": "^0.18.5"}}`:                    true,
		`{"dependencies": {"xlsx": "0.19.3"}}`:                     true,
		`{"dependencies": {"xlsx": "~0.20.1"}}`:                    true,
		`{"dependencies": {"xlsx": "^0.20.2"}}`:                    false,
		`{"dependencies": {"react": "19.2.7"}}`:                    false,
		`{"dependencies": {"xlsx": "` + xlsxPatchedTarball + `"}}`: false,
	} {
		if got := vulnerableXLSXMessage(src) != ""; got != vulnerable {
			t.Errorf("%s: vulnerable = %v, want %v", src, got, vulnerable)
		}
	}
}

// No template takes xlsx from npm any more.
func TestFreshTemplatesUsePatchedXLSX(t *testing.T) {
	double := Options{ProjectName: "demo", Architecture: ArchDouble, Frontend: FrontendNext}
	triple := Options{ProjectName: "demo", Architecture: ArchTriple, Frontend: FrontendNext}
	tanstack := Options{ProjectName: "demo", Architecture: ArchTriple, Frontend: FrontendTanStack}
	single := Options{ProjectName: "demo", Architecture: ArchSingle}
	for name, src := range map[string]string{
		"web app with the admin panel": webAdminDependencies(double),
		"Next admin app":               adminPackageJSON(triple),
		"TanStack admin app":           adminTanStackPackageJSON(tanstack),
		"single frontend":              singleFrontendPackageJSON(single),
	} {
		if !strings.Contains(src, xlsxPatchedTarball) || vulnerableXLSXMessage(src) != "" {
			t.Errorf("%s still takes xlsx from npm", name)
		}
	}
}
