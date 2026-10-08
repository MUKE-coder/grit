package authz

import "testing"

func TestGranted(t *testing.T) {
	cases := []struct {
		name  string
		grant string
		key   string
		want  bool
	}{
		{"exact", "products.create", "products.create", true},
		{"exact mismatch", "products.create", "products.delete", false},
		{"superuser", "*", "anything.at.all", true},
		{"resource wildcard", "products.*", "products.create", true},
		{"resource wildcard is scoped", "products.*", "orders.create", false},
		{"action wildcard", "*.view", "products.view", true},
		{"action wildcard is scoped", "*.view", "products.delete", false},
		{"no partial segment match", "prod.*", "products.view", false},

		// The Shoppleet bug: a matcher that returns true at the first "*"
		// treats this as "products.*" and wrongly grants delete.
		{"middle wildcard checks the tail", "products.*.view", "products.variants.delete", false},
		{"middle wildcard matches its tail", "products.*.view", "products.variants.view", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Granted([]string{tc.grant}, tc.key); got != tc.want {
				t.Errorf("Granted([%q], %q) = %v, want %v", tc.grant, tc.key, got, tc.want)
			}
		})
	}
}

func TestGrantedAnyOf(t *testing.T) {
	grants := []string{"uploads.view", "products.*"}
	if !Granted(grants, "products.delete") {
		t.Error("expected products.* to authorise products.delete")
	}
	if !Granted(grants, "uploads.view") {
		t.Error("expected the exact grant to match")
	}
	if Granted(grants, "users.delete") {
		t.Error("ungranted key must not pass")
	}
	if Granted(nil, "products.view") {
		t.Error("no grants must authorise nothing")
	}
}

func TestKeysAndExpand(t *testing.T) {
	keys := Keys()
	if len(keys) == 0 {
		t.Fatal("catalog produced no keys")
	}

	// Every catalog key must be reachable via the superuser grant.
	if got := len(Expand([]string{"*"})); got != len(keys) {
		t.Errorf("Expand(*) = %d keys, want %d", got, len(keys))
	}
	if !HasAll([]string{"*"}) {
		t.Error("HasAll(*) must be true")
	}

	// A resource wildcard expands only within that resource.
	for _, k := range Expand([]string{"users.*"}) {
		if len(k) < 6 || k[:6] != "users." {
			t.Errorf("users.* leaked %q", k)
		}
	}
	if HasAll([]string{"users.*"}) {
		t.Error("a single resource must not satisfy HasAll")
	}
}
