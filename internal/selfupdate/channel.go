package selfupdate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// The stable channel.
//
// Grit ships every merged change, several times a day, and that is right for
// the project and wrong for a team with a change-control process. A handful
// of releases have fixed a regression from the days before, so the newest
// release is occasionally the worst one to be on (grit#245).
//
// stable is the newest release that is at least a week old and had every
// scheduled canary run in that week come back green. scripts/promote-stable.py
// decides it, docs/config/stable.ts records it, and the site serves it here.
//
// There are no backports. A fix lands on edge and reaches stable when it has
// aged; if you need it sooner, take the edge release, which is what `grit
// update` does by default.

// ChannelURL is where the channel file lives.
//
// The site rather than a moving git tag. A moving tag has to be force-pushed,
// and a force-pushed tag is a change nobody can audit afterwards; this is a
// file in the repository, and the commit that moved it is in the history.
const ChannelURL = "https://gritframework.dev/stable.json"

// channelURL is the address actually used, so a test can point it at a
// local server. The exported constant stays the documented answer.
var channelURL = ChannelURL

// Channel is what that file says.
type Channel struct {
	Stable     string `json:"stable"`
	Edge       string `json:"edge"`
	Promoted   string `json:"promoted"`
	Released   string `json:"released"`
	CanaryRuns int    `json:"canaryRuns"`
	MinAgeDays int    `json:"minAgeDays"`
	Rule       string `json:"rule"`
}

// FetchChannel reads the channel file.
//
// Short timeout and no retry: this is one step of a command somebody is
// waiting on, and a channel that cannot be read should say so rather than
// hang. The caller decides whether to fall back.
func FetchChannel() (*Channel, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, channelURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "grit-cli")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("reading the stable channel: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reading the stable channel: %s returned %s",
			channelURL, resp.Status)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return nil, fmt.Errorf("reading the stable channel: %w", err)
	}
	var ch Channel
	if err := json.Unmarshal(body, &ch); err != nil {
		return nil, fmt.Errorf("the stable channel is not JSON: %w", err)
	}
	if strings.TrimSpace(ch.Stable) == "" {
		return nil, fmt.Errorf("the stable channel names no version")
	}
	return &ch, nil
}

// ResolveChannel turns a channel name into the version to install.
//
// "edge" and "" both mean the latest release, which is what the GitHub API
// answers and what this command has always done. Only "stable" needs the
// channel file, so an unreachable site cannot break the default.
func ResolveChannel(name string) (version string, err error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "edge", "latest":
		return "", nil
	case "stable":
		ch, err := FetchChannel()
		if err != nil {
			return "", err
		}
		return ch.Stable, nil
	default:
		return "", fmt.Errorf("unknown channel %q: use stable or edge", name)
	}
}
