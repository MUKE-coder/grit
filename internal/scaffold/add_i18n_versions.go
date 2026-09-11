package scaffold

import "regexp"

// nextIntl3 matches a next-intl 3 dependency in package.json. next-intl 3
// supports Next up to 15; the scaffold is on Next 16, whose production builds
// failed with it (see injectJSONDep).
var nextIntl3 = regexp.MustCompile(`"next-intl":\s*"[~^]?3\.[^"]*"`)
