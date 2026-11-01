package setup

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MUKE-coder/grit/v3/internal/plugin"
	"github.com/MUKE-coder/grit/v3/internal/project"
)

// pluginContextAt describes the project to a plugin. setup cannot call the
// command package's helper, and the new single layout puts the Go module in
// api/, which the command helper predates; using scaffolding's own APIRoot is
// the one answer that stays right for every shape setup can create.
func pluginContextAt(root string, o Options) (plugin.Context, error) {
	scaffoldOpts := o.ScaffoldOptions()

	ctx := plugin.Context{
		Root:         root,
		Architecture: o.Architecture,
		Frontend:     o.Frontend,
		APIRoot:      scaffoldOpts.APIRoot(root),
		Module:       scaffoldOpts.Module(),
	}

	if info, err := project.DetectProjectFrom(root); err == nil && info.Module != "" {
		ctx.Module = info.Module
	}

	if data, err := os.ReadFile(filepath.Join(root, "grit.json")); err == nil {
		var meta struct {
			Architecture string `json:"architecture"`
			Frontend     string `json:"frontend"`
		}
		if json.Unmarshal(data, &meta) == nil {
			if meta.Frontend != "" {
				ctx.Frontend = meta.Frontend
			}
			// grit.json writes the shape that was built, not the shorthand
			// that was asked for: full and desktop are both triple on disk,
			// with their extra apps recorded separately.
			if meta.Architecture != "" {
				ctx.Architecture = meta.Architecture
			}
		}
	}

	if o.Architecture == "full" || o.Architecture == "desktop" {
		ctx.Architecture = "triple"
	}

	return ctx, nil
}

// InstallPlugins installs the chosen built-ins through the same call
// `grit plugin add` makes. The project was created a second ago, so there is
// no dirty tree to protect and no lockfile entry that could collide.
func InstallPlugins(ctx context.Context, root string, o Options) ([]plugin.Plugin, error) {
	pctx, err := pluginContextAt(root, o)
	if err != nil {
		return nil, err
	}

	installed := make([]plugin.Plugin, 0, len(o.Plugins))
	for _, name := range o.Plugins {
		p, err := plugin.Get(name)
		if err != nil {
			return installed, err
		}
		if _, err := plugin.Install(pctx, p); err != nil {
			return installed, fmt.Errorf("installing %s: %w\n\nRetry with: grit plugin add %s", name, err, name)
		}
		installed = append(installed, p)
	}
	return installed, nil
}
