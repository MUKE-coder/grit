package scaffold

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// A project with the Expo app gets the PostCSS override and the two accepted
// image-size advisories; one without gets neither.
func TestPnpmWorkspaceExpoAudit(t *testing.T) {
	var ws struct {
		Overrides   map[string]string `yaml:"overrides"`
		AuditConfig struct {
			IgnoreGhsas []string `yaml:"ignoreGhsas"`
		} `yaml:"auditConfig"`
	}
	if err := yaml.Unmarshal([]byte(pnpmWorkspace(true, true)), &ws); err != nil {
		t.Fatalf("pnpm-workspace.yaml does not parse: %v", err)
	}
	if ws.Overrides["@expo/metro-config>postcss"] == "" {
		t.Error("no PostCSS override for Expo's bundler")
	}
	if len(ws.AuditConfig.IgnoreGhsas) != 2 {
		t.Errorf("accepted advisories: %v, want the two image-size ones", ws.AuditConfig.IgnoreGhsas)
	}
	if plain := pnpmWorkspace(true, false); strings.Contains(plain, "overrides:") || strings.Contains(plain, "auditConfig:") {
		t.Error("a project without Expo got Expo's audit settings")
	}
}
