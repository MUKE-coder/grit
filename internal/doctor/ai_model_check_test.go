package doctor

import (
	"strings"
	"testing"
)

func TestAIModelIDFindings(t *testing.T) {
	cases := []struct {
		name  string
		env   map[string]string
		level string // "" means nothing reported
		says  string
	}{
		{
			name: "the id Grit itself shipped",
			env: map[string]string{
				"AI_GATEWAY_MODEL":   "anthropic/claude-sonnet-4-6",
				"AI_GATEWAY_API_KEY": "vck_live_xxx",
			},
			level: "error",
			says:  "anthropic/claude-sonnet-4.6",
		},
		{
			name:  "the same id with no key is not failing yet",
			env:   map[string]string{"AI_GATEWAY_MODEL": "anthropic/claude-sonnet-4-6"},
			level: "warning",
			says:  "anthropic/claude-sonnet-4.6",
		},
		{
			name:  "no provider in front of it",
			env:   map[string]string{"AI_GATEWAY_MODEL": "claude-sonnet-5.5"},
			level: "warning",
			says:  "provider/model",
		},
		{
			name:  "the current default",
			env:   map[string]string{"AI_GATEWAY_MODEL": "anthropic/claude-sonnet-5.5"},
			level: "",
		},
		{
			name:  "a dotted Anthropic minor",
			env:   map[string]string{"AI_GATEWAY_MODEL": "anthropic/claude-sonnet-4.6"},
			level: "",
		},
		{
			// Anthropic's older naming. A date on the end is not a dashed minor,
			// and calling it one would be a check that is wrong about a real id.
			name:  "a dated Anthropic name",
			env:   map[string]string{"AI_GATEWAY_MODEL": "anthropic/claude-3-5-sonnet-20241022"},
			level: "",
		},
		{
			// Another provider's dashes are its own business: only Anthropic's
			// catalogue was read closely enough to say.
			name:  "another provider with dashes",
			env:   map[string]string{"AI_GATEWAY_MODEL": "meta/llama-4-scout"},
			level: "",
		},
		{
			name:  "unset, so the compiled-in default applies",
			env:   map[string]string{},
			level: "",
		},
		{
			name:  "set to whitespace",
			env:   map[string]string{"AI_GATEWAY_MODEL": "   "},
			level: "",
		},
		{
			name:  "a batch variant",
			env:   map[string]string{"AI_GATEWAY_MODEL": "anthropic/claude-sonnet-5.5-batch"},
			level: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := checkAIModelID(&project{env: c.env})
			if c.level == "" {
				if len(got) != 0 {
					t.Fatalf("expected nothing, got %+v", got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("expected one finding, got %d: %+v", len(got), got)
			}
			if got[0].Level != c.level {
				t.Errorf("level: want %q, got %q", c.level, got[0].Level)
			}
			if c.says != "" && !strings.Contains(got[0].Fix+got[0].Message, c.says) {
				t.Errorf("neither the message nor the fix mentions %q:\n  %s\n  %s",
					c.says, got[0].Message, got[0].Fix)
			}
		})
	}
}

// TestGritsOwnDefaultPassesItsOwnCheck reads both sides: the default the
// scaffold writes into a project's .env.example, and the check that judges it.
// Shipping a default this check would flag is the exact mistake grit#186 was.
func TestGritsOwnDefaultPassesItsOwnCheck(t *testing.T) {
	for _, model := range scaffoldedAIModels(t) {
		got := checkAIModelID(&project{env: map[string]string{
			"AI_GATEWAY_MODEL":   model,
			"AI_GATEWAY_API_KEY": "vck_live_xxx",
		}})
		if len(got) != 0 {
			t.Errorf("the scaffold writes AI_GATEWAY_MODEL=%s, which this check rejects: %s",
				model, got[0].Message)
		}
	}
}
