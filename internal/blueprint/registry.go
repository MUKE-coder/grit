// Package blueprint lists the complete applications built with Grit.
//
// A blueprint is not something the CLI installs. It is a whole repository you
// clone and work from, which is the point: a plugin adds a feature to your
// project, a blueprint IS a project. `grit blueprint list` exists so that
// somebody deciding how to start can see what already exists without leaving
// the terminal to find out.
//
// The list is compiled in, like the plugin registry, and so it is as current
// as the binary. `grit update` is the refresh.
package blueprint

import "sort"

// Access says how a blueprint is obtained.
type Access string

const (
	// Open is a repository anyone can clone.
	Open Access = "open"
	// Paid needs a membership from the blueprints site, which grants an
	// invitation to the private repository.
	Paid Access = "paid"
	// Planned is announced but not built. Listed because "is there a Linear
	// blueprint?" is a question worth answering with "not yet" rather than
	// with silence.
	Planned Access = "planned"
)

// Surface is a client a blueprint ships.
type Surface string

const (
	Web     Surface = "web"
	Admin   Surface = "admin"
	Mobile  Surface = "mobile"
	Desktop Surface = "desktop"
)

// Blueprint is one complete application.
type Blueprint struct {
	Name string
	// Title is what it is called, as a person would say it.
	Title string
	// Summary is one line, shown by `grit blueprint list`.
	Summary string
	// Teaches is what somebody would open this one to learn, rather than what
	// it is a clone of. This is the field that makes the list useful: every
	// blueprint "shows how to build an app", and that sentence helps nobody
	// choose between them.
	Teaches []string
	// Architecture is the grit architecture it was scaffolded with.
	Architecture string
	Surfaces     []Surface
	Access       Access
	// Repo is the GitHub repository. Empty while Planned.
	Repo string
}

// Site is where blueprints are bought and the repository invitation is sent.
const Site = "https://gritframework.dev/blueprints"

var registry = []Blueprint{
	{
		Name:         "whatsapp",
		Title:        "WhatsApp",
		Summary:      "Realtime chat across web, mobile and desktop, on one Go API",
		Architecture: "triple",
		Surfaces:     []Surface{Web, Admin, Mobile, Desktop},
		Access:       Open,
		Repo:         "https://github.com/MUKE-coder/grit-whatsapp",
		Teaches: []string{
			"One conversation model reaching four clients without four implementations",
			"Message delivery and read state over the domain event bus",
			"Attachments through the same storage the rest of the app uses",
			"An Expo app and a Wails desktop app against the same shared types",
		},
	},
	{
		Name:         "instagram",
		Title:        "Instagram",
		Summary:      "Photo and video feed with an Expo app",
		Architecture: "triple",
		Surfaces:     []Surface{Web, Admin, Mobile},
		Access:       Paid,
		Repo:         "https://github.com/MUKE-coder/grit-instagram",
		Teaches: []string{
			"Image and video processing on upload, and what to do off the request",
			"A feed that paginates without drifting as new posts arrive",
			"Mobile capture and upload from Expo",
		},
	},
	{
		Name:         "storefront",
		Title:        "Kaya Goods",
		Summary:      "A shop: products, variants, cart and Stripe checkout",
		Architecture: "triple",
		Surfaces:     []Surface{Web, Admin},
		Access:       Paid,
		Repo:         "https://github.com/MUKE-coder/grit-storefront",
		Teaches: []string{
			"Product variants, and the option matrix that generates them",
			"A cart that survives sign-in, and a checkout that survives a retry",
			"Stripe webhooks that are idempotent when Stripe redelivers",
		},
	},
	{
		Name:         "stripe-dashboard",
		Title:        "Stripe dashboard",
		Summary:      "Payments, subscriptions and revenue reporting",
		Architecture: "triple",
		Surfaces:     []Surface{Admin},
		Access:       Planned,
		Teaches: []string{
			"Reporting over money without recomputing it on every page view",
			"Subscription lifecycle: trials, dunning, cancellation",
		},
	},
	{
		Name:         "uber",
		Title:        "Uber",
		Summary:      "Live tracking on a map, with PostGIS",
		Architecture: "triple",
		Surfaces:     []Surface{Web, Admin, Mobile},
		Access:       Planned,
		Teaches: []string{
			"Geospatial queries in Postgres, and what they cost",
			"A position stream that does not melt the database",
		},
	},
	{
		Name:         "linear",
		Title:        "Linear",
		Summary:      "An issue tracker built for the keyboard",
		Architecture: "triple",
		Surfaces:     []Surface{Web, Admin},
		Access:       Planned,
		Teaches: []string{
			"Keyboard-first navigation across a whole application",
			"Optimistic updates that stay correct when the write is refused",
		},
	},
}

// All returns every blueprint, available ones first and planned ones last,
// each group sorted by name. Somebody scanning this list wants to know what
// they can have today; what is coming is worth saying, but not first.
func All() []Blueprint {
	out := make([]Blueprint, len(registry))
	copy(out, registry)
	sort.SliceStable(out, func(i, j int) bool {
		if (out[i].Access == Planned) != (out[j].Access == Planned) {
			return out[j].Access == Planned
		}
		return out[i].Name < out[j].Name
	})
	return out
}

// Get returns one blueprint by name.
func Get(name string) (Blueprint, bool) {
	for _, b := range registry {
		if b.Name == name {
			return b, true
		}
	}
	return Blueprint{}, false
}

// Names returns every blueprint name, for error messages and completion.
func Names() []string {
	out := make([]string, 0, len(registry))
	for _, b := range registry {
		out = append(out, b.Name)
	}
	sort.Strings(out)
	return out
}

// SurfaceList renders surfaces for display.
func (b Blueprint) SurfaceList() []string {
	out := make([]string, 0, len(b.Surfaces))
	for _, s := range b.Surfaces {
		out = append(out, string(s))
	}
	return out
}
