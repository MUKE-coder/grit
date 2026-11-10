package prompt

import (
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// The picker offers every provider the flag accepts.
//
// The same shape as TestThePickerOffersEveryDatabase, and for the same reason:
// the picker is what most people see, so a provider reachable only with
// --storage may as well not exist for them. That is exactly how this started.
// Every driver worked for a long time; the only way to find out there were
// five was to scaffold with minio and then read a comment in .env.
func TestThePickerOffersEveryStorageDriver(t *testing.T) {
	options := storageOptions()

	if len(options) != len(scaffold.StorageDrivers) {
		t.Fatalf("the picker offers %d providers, the flag accepts %d",
			len(options), len(scaffold.StorageDrivers))
	}

	// Each option's value is the name .env will hold, and its label carries
	// the same sentence the flag's own map does, so the two cannot describe
	// the same provider differently.
	for name, meaning := range scaffold.StorageDrivers {
		var found bool
		for _, opt := range options {
			if opt.Value != name {
				continue
			}
			found = true
			if !strings.Contains(opt.Key, meaning) {
				t.Errorf("the %s option reads %q, which does not carry %q",
					name, opt.Key, meaning)
			}
		}
		if !found {
			t.Errorf("%q is accepted by --storage and is not in the picker", name)
		}
	}
}

// The order is the one StorageDriverOrder declares. A map has none, and a
// picker whose options move between runs is one people stop reading.
func TestTheStoragePickerKeepsItsOrder(t *testing.T) {
	options := storageOptions()
	if len(options) != len(scaffold.StorageDriverOrder) {
		t.Fatalf("the picker has %d options and the order lists %d",
			len(options), len(scaffold.StorageDriverOrder))
	}
	for i, want := range scaffold.StorageDriverOrder {
		if options[i].Value != want {
			t.Errorf("option %d is %q, want %q", i, options[i].Value, want)
		}
	}
}

// What the picker returns is what Validate accepts unchanged.
//
// A picker that offered "Local disk" and set a value Validate then rewrote,
// or refused, would be a silent disagreement between the two ways in.
func TestEveryPickedStorageDriverSurvivesValidation(t *testing.T) {
	for _, name := range scaffold.StorageDriverOrder {
		opts := scaffold.Options{StorageDriver: name}
		if err := opts.ValidateStorageDriver(); err != nil {
			t.Errorf("the picker can return %q and the flag refuses it: %v", name, err)
			continue
		}
		if opts.StorageDriver != name {
			t.Errorf("%q normalised to %q, so the picker's value is not the one written to .env",
				name, opts.StorageDriver)
		}
	}
}
