package main

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/MUKE-coder/grit/v3/internal/plugin"
	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// --saas, and what it actually does.
//
// Nothing here is new machinery. The signed-in customer area is what
// `grit add web-auth` writes, organizations are the multitenant plugin and
// charging is the stripe plugin, each already installed by a command of its
// own and each already tested. What --saas adds is the question, asked once,
// at the point where somebody knows the answer.
//
// The reason it is worth asking at all: a resource generated in a project
// with customers needs two screens, not one. The operator's console has
// always been generated; the customer's half was a React Query hook and a
// blank page. grit.json carries the bit, and internal/generate reads it.

// installCustomerArea writes the signed-in half of the web app.
//
// First of the three, and the order is load-bearing: the organization screen
// and the billing screen are pages inside app/(app), and each plugin looks
// for that directory to decide where to put its page.
func installCustomerArea(root string, opts scaffold.Options) error {
	// An API-only or mobile-only project has no browser surface to sign into.
	// Not an error: --saas on an --api project is a reasonable thing to type,
	// and the bit in grit.json still means something to a later `grit add`.
	if !opts.ShouldIncludeWeb() {
		color.Yellow("  ! --saas with no web app: grit.json records it, and the customer")
		color.Yellow("    screens arrive when this project gets one.")
		return nil
	}

	if err := scaffold.AddWebAuth(root, false); err != nil {
		return fmt.Errorf("adding the customer sign-in area: %w", err)
	}
	fmt.Println("  ✓ Customer area at apps/web/app/(app), with sign-in, guards and the account pages")
	return nil
}

// installPluginAt runs one plugin's install, the way `grit plugin add` does.
//
// installMultitenantAt is the same thing with a message of its own, written
// before there was a second caller. This is the general form.
func installPluginAt(root, name string) error {
	p, err := plugin.Get(name)
	if err != nil {
		return err
	}
	ctx, err := pluginContextAt(root)
	if err != nil {
		return err
	}
	record, err := plugin.Install(ctx, p)
	if err != nil {
		return err
	}
	fmt.Printf("  ✓ %s added (%d files, %d patches)\n", p.Summary, len(record.Files), len(record.Injections))
	return nil
}
