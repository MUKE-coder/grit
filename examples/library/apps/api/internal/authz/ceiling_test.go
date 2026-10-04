package authz

import (
	"reflect"
	"testing"
)

// The ceiling on handing out permissions. Before it, a users.edit or roles.edit
// holder could make themselves ADMIN.
func TestBeyondIsTheGrantCeiling(t *testing.T) {
	cases := []struct {
		name            string
		held, requested []string
		want            []string
	}{
		{"a holder of everything may grant anything", []string{"*"}, []string{"*", "users.edit"}, nil},
		{"nobody else may grant everything", []string{"users.edit"}, []string{"*"}, []string{"*"}},
		{"a held permission may be granted", []string{"users.view", "users.edit"}, []string{"users.edit"}, nil},
		{"one not held may not", []string{"users.view"}, []string{"users.delete"}, []string{"users.delete"}},
		{"a wildcard needs everything it expands to", []string{"users.view"}, []string{"users.*"}, []string{"users.*"}},
		{"a held wildcard covers its permissions", []string{"users.*"}, []string{"users.edit", "users.*"}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Beyond(tc.held, tc.requested); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Beyond(%v, %v) = %v, want %v", tc.held, tc.requested, got, tc.want)
			}
		})
	}
}
