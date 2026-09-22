package scaffold

import (
	"regexp"
	"strings"
)

// perResource renames a template's package-level helpers after its resource.
//
// grit add variants writes a service, a seeder and their tests per resource,
// all into shared packages, so a second resource with variants wrote a second
// shortID and a second variantTestDB beside the first and the API stopped
// compiling. The docs promised a second run adds only that resource's tables.
// Found building the storefront blueprint.
func perResource(src, pascal string, helpers ...string) string {
	camel := strings.ToLower(pascal[:1]) + pascal[1:]
	for _, helper := range helpers {
		named := camel + strings.ToUpper(helper[:1]) + helper[1:]
		src = regexp.MustCompile(`\b`+regexp.QuoteMeta(helper)+`\b`).ReplaceAllString(src, named)
	}
	return src
}
