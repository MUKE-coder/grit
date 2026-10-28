package plugin

import (
	_ "embed"
	"encoding/json"
	"sort"
)

// communityJSON is the index of plugins other people wrote.
//
// Embedded rather than fetched, for the same reason the built-in registry is
// compiled in: `grit plugin list` works on a plane, and a list that sometimes
// takes four seconds and sometimes fails is worse than one that is as current
// as the binary. `grit update` is the refresh.
//
//go:embed community.json
var communityJSON []byte

// CommunityPlugin is one entry in the index.
type CommunityPlugin struct {
	Name    string `json:"name"`
	Author  string `json:"author"`
	Repo    string `json:"repo"`
	Summary string `json:"summary"`
	// Added is the date the entry was accepted, as YYYY-MM-DD. It is shown
	// because an index without dates cannot be told from an abandoned one.
	Added string `json:"added"`
	// Architectures it supports, empty meaning all of them.
	Architectures []string `json:"architectures,omitempty"`
}

// SubmitURL opens a submission against the index.
const SubmitURL = "https://github.com/MUKE-coder/grit/issues/new?template=plugin-submission.yml&labels=plugin-submission"

// MarketplaceURL is the browsable version of this list.
const MarketplaceURL = "https://gritframework.dev/docs/plugins/marketplace"

var communityCache []CommunityPlugin

// Community returns the community plugin index, sorted by name.
func Community() []CommunityPlugin {
	if communityCache != nil {
		return communityCache
	}
	var doc struct {
		Plugins []CommunityPlugin `json:"plugins"`
	}
	// A malformed index is a build-time mistake in this repository, not a
	// user's problem, and an empty list is the honest answer either way.
	// TestCommunityIndexParses makes sure it never ships malformed.
	_ = json.Unmarshal(communityJSON, &doc)

	sort.Slice(doc.Plugins, func(i, j int) bool { return doc.Plugins[i].Name < doc.Plugins[j].Name })
	communityCache = doc.Plugins
	return communityCache
}

// CommunityJSON returns the raw index, for the script that syncs it into the
// docs site.
func CommunityJSON() []byte { return communityJSON }
