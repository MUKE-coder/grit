package prompt

import (
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// The picker and the --db flag have to offer the same engines.
//
// The same arrangement the theme picker has, and for the same reason: when the
// flag learns about an engine and the picker does not, nothing fails. The
// option is simply absent, so for everybody who does not read the help text it
// has not shipped.
func TestThePickerOffersEveryDatabase(t *testing.T) {
	options := databaseOptions()
	if len(options) != len(scaffold.DBProviders) {
		t.Fatalf("the picker offers %d engines, the flag accepts %d", len(options), len(scaffold.DBProviders))
	}

	offered := map[string]string{}
	for _, o := range options {
		offered[o.Value] = o.Key
	}
	for name, meaning := range scaffold.DBProviders {
		label, ok := offered[name]
		if !ok {
			t.Errorf("the picker does not offer %q, so it is reachable only with --db", name)
			continue
		}
		// The label carries the description the flag documents, so the two
		// cannot say different things about the same engine.
		if !strings.Contains(label, meaning) {
			t.Errorf("the option for %q reads %q, which does not say %q", name, label, meaning)
		}
	}

	// Postgres stays first: it is the default, and a picker that opens on
	// something else quietly changes what `grit new` does for anybody who
	// presses enter.
	if options[0].Value != "postgres" {
		t.Errorf("the first option is %q, not postgres, which is still the default", options[0].Value)
	}

	// The order is fixed. A picker whose options move between runs is one
	// people stop reading.
	for i, want := range scaffold.DBProviderOrder {
		if options[i].Value != want {
			t.Errorf("option %d is %q, want %q", i, options[i].Value, want)
		}
	}
}

// Every engine the order names is an engine the validator accepts.
func TestEveryOfferedDatabaseIsValid(t *testing.T) {
	for _, name := range scaffold.DBProviderOrder {
		opts := scaffold.Options{DBProvider: name}
		if err := opts.ValidateDBProvider(); err != nil {
			t.Errorf("the picker offers %q and the validator refuses it: %v", name, err)
		}
		if opts.DBProvider != name {
			t.Errorf("%q normalised to %q, so the picker's value is not the one written to .env", name, opts.DBProvider)
		}
	}
}
