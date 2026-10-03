package mcp

import "fmt"

// Mode decides which tools exist on a server.
//
// # Absence, not a permission check
//
// A read-only server does not contain the tools that can change a project. They
// are not registered, so there is no code path that reaches one: no token, no
// scope, no misconfigured client and no prompt-injected README can call what is
// not in the map. tools/list and tools/call read the same map, so a tool that is
// not listed is also not callable, from one decision rather than two that have
// to agree.
//
// The alternative, which is what most servers do, is to register everything and
// check a flag inside each handler. That works until somebody adds a tool and
// forgets the check, and the failure is silent and total: the tool simply works
// for everyone. Making the mode decide what is built removes the class of
// mistake instead of asking every future contributor to remember.
//
// # Why read is the default
//
// The agent pointed at this server is usually reading a checkout to answer a
// question. Running a generator is a different act with a different blast
// radius, and it should take saying so: `--mode write`. An editor that
// registers the server once and forgets is registering the safe one.
type Mode string

const (
	// ModeRead answers questions about a project and cannot change it.
	ModeRead Mode = "read"
	// ModeWrite also exposes the generators, which write files.
	ModeWrite Mode = "write"
)

// ParseMode reads a mode name, defaulting to read when empty.
func ParseMode(s string) (Mode, error) {
	switch s {
	case "", string(ModeRead):
		return ModeRead, nil
	case string(ModeWrite):
		return ModeWrite, nil
	default:
		return "", fmt.Errorf("unknown mode %q: use read or write", s)
	}
}

// allows reports whether a tool with this mutates flag belongs on a server in
// this mode.
func (m Mode) allows(mutates bool) bool {
	return !mutates || m == ModeWrite
}
