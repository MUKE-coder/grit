package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// xlsxPatchedTarball is SheetJS's own build of xlsx with the fixes. SheetJS
// stopped publishing to npm at 0.18.5, which has two high advisories (prototype
// pollution, and a regular expression denial of service), so npm offers no
// patched version at all and pnpm audit fails on every project that has it.
const xlsxPatchedTarball = "https://cdn.sheetjs.com/xlsx-0.20.3/xlsx-0.20.3.tgz"

var xlsxNpmRangeRe = regexp.MustCompile(`"xlsx":\s*"[\^~]?0\.(\d+)\.(\d+)"`)

// warnVulnerableXLSX points out every package.json that still takes xlsx from
// npm, below the first fixed release. package.json is the developer's, so
// upgrade names the command rather than editing the dependency.
func warnVulnerableXLSX(root string) {
	candidates, _ := filepath.Glob(filepath.Join(root, "apps", "*", "package.json"))
	candidates = append(candidates, filepath.Join(root, "frontend", "package.json"))
	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if msg := vulnerableXLSXMessage(string(data)); msg != "" {
			rel, _ := filepath.Rel(root, filepath.Dir(path))
			fmt.Printf("  ⚠ %s: %s. In %s run:\n      pnpm add xlsx@%s\n", filepath.ToSlash(rel), msg, filepath.ToSlash(rel), xlsxPatchedTarball)
		}
	}
}

// vulnerableXLSXMessage reports why a package.json's xlsx needs replacing, or
// "" when it does not.
func vulnerableXLSXMessage(packageJSON string) string {
	m := xlsxNpmRangeRe.FindStringSubmatch(packageJSON)
	if m == nil {
		return ""
	}
	minor, _ := strconv.Atoi(m[1])
	patch, _ := strconv.Atoi(m[2])
	if minor > 20 || (minor == 20 && patch >= 2) {
		return ""
	}
	return fmt.Sprintf("xlsx 0.%d.%d from npm has two high advisories and no patched npm release, so pnpm audit fails", minor, patch)
}
