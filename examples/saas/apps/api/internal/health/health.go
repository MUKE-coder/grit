// Package health is the vocabulary /api/health answers in, and the registry a
// subsystem adds itself to.
//
// # Why four states and not a boolean
//
// Every component used to report {"ok": false}, and that one word covered two
// opposite situations: Redis is down, and this deployment has no Redis. A reader
// cannot tell them apart, so a dashboard either cries wolf about a component
// nobody asked for or stays quiet about one that just died. The two that matter
// are separated here:
//
//	ok        wired up, and nothing known to be wrong
//	degraded  wired up, and not working, or working worse than it should
//	off       deliberately not wired up; this is not a problem
//	unknown   wired up, but this probe could not find out
//
// "unknown" is the honest answer when a count has not been taken yet or a probe
// could not run, and it is deliberately not "degraded": not knowing is not the
// same as being broken, and treating it as broken is how an on-call rota learns
// to ignore a page.
//
// # What Overall does with them
//
// Only degraded lowers the overall status. A component that is off was never
// asked for, and one that is unknown has not accused anybody of anything, so
// neither drags a deployment's status down. That rule is in one place, here, so
// a new component cannot quietly invent its own.
//
// # What a probe may do
//
// Nothing slow. A probe reports state its subsystem already holds, or makes one
// bounded call with a deadline the caller sets. It must not block on a
// dependency that is already unwell, because /api/health is exactly what somebody
// reaches for when it is, and a probe that waits on a hung database turns the one
// endpoint that could explain the outage into another symptom of it.
//
// # No probe reports a secret
//
// A report goes out over HTTP, and /api/health is reachable without
// authentication so a load balancer can use it. Report a credential by presence
// or by length, never by value:
//
//	health.OK("configured", map[string]any{"signing_key": true})     // yes
//	health.OK("configured", map[string]any{"signing_key": cfg.Key})  // never
package health

import (
	"encoding/json"
	"sort"
	"sync"
)

// State is what a component is doing, in one word.
type State string

const (
	// StateOK means wired up with nothing known to be wrong. It does not
	// promise a round trip was made: a component that cannot be probed without
	// I/O says so in its detail.
	StateOK State = "ok"
	// StateDegraded means wired up and not working, or working worse than it
	// should. This is the only state that lowers the overall status.
	StateDegraded State = "degraded"
	// StateOff means deliberately not configured. An API with no Redis is not
	// an unhealthy API.
	StateOff State = "off"
	// StateUnknown means wired up, but this probe could not find out.
	StateUnknown State = "unknown"
)

// Report is one component's answer.
//
// Fields carries whatever that component measures: a latency, a queue depth, a
// driver name. It is merged into the JSON object rather than nested under a key,
// so a reader sees {"state":"ok","latency_ms":3} and not {"state":"ok","fields":{...}}.
type Report struct {
	State  State
	Detail string
	Fields map[string]any
}

// MarshalJSON writes state, the legacy ok flag, the detail, and then whatever
// the component measured, flattened into the same object.
//
// ok is still written because clients read it: a desktop heartbeat, a load
// balancer check and an admin page that predates this package. It is exactly
// State == StateOK, so a reader that only knows the boolean sees "off" as false,
// which is what it saw before.
func (r Report) MarshalJSON() ([]byte, error) {
	out := make(map[string]any, len(r.Fields)+3)
	for k, v := range r.Fields {
		out[k] = v
	}
	state := r.State
	if state == "" {
		state = StateUnknown
	}
	out["state"] = state
	out["ok"] = state == StateOK
	if r.Detail != "" {
		out["detail"] = r.Detail
	}
	return json.Marshal(out)
}

// OK reports a working component.
func OK(detail string, fields map[string]any) Report {
	return Report{State: StateOK, Detail: detail, Fields: fields}
}

// Degraded reports a component that is wired up and not working. The detail is
// what an operator reads first, so it says what is wrong rather than naming the
// component again.
func Degraded(detail string, fields map[string]any) Report {
	return Report{State: StateDegraded, Detail: detail, Fields: fields}
}

// Off reports a component nobody configured. The detail says what would turn it
// on, because "off" with no instructions is a dead end.
func Off(detail string) Report {
	return Report{State: StateOff, Detail: detail}
}

// Unknown reports a component that is wired up and could not be measured.
func Unknown(detail string) Report {
	return Report{State: StateUnknown, Detail: detail}
}

// Overall is the status for the whole deployment: degraded when any component
// is, ok otherwise.
//
// Off and unknown never lower it. A deployment with no Redis is healthy, and one
// whose queue depth has not been sampled yet is healthy until something says
// otherwise.
func Overall(components map[string]Report) State {
	for _, c := range components {
		if c.State == StateDegraded {
			return StateDegraded
		}
	}
	return StateOK
}

// A probe is a component's own answer about itself.
type probe struct {
	name string
	fn   func() Report
}

var (
	mu     sync.RWMutex
	probes []probe
)

// Register adds a component to /api/health.
//
// For anything the route file does not already hold a handle to: a plugin, a
// client you wired in main.go, a dependency only your code knows about. Call it
// once while the application is being built, never per request.
//
//	health.Register("billing", func() health.Report {
//		if !gateway.Configured() {
//			return health.Off("no payment gateway; set BILLING_KEY")
//		}
//		return health.OK("", map[string]any{"charges_settled": settled.Load()})
//	})
//
// Registering the same name twice replaces the first, so a reload cannot show
// one component twice.
func Register(name string, fn func() Report) {
	if name == "" || fn == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	for i := range probes {
		if probes[i].name == name {
			probes[i].fn = fn
			return
		}
	}
	probes = append(probes, probe{name: name, fn: fn})
}

// Registered reports the names that have been registered, sorted. For a test
// that wants to know the wiring happened without running the probes.
func Registered() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(probes))
	for _, p := range probes {
		out = append(out, p.name)
	}
	sort.Strings(out)
	return out
}

// Snapshot runs every registered probe and returns what they said.
//
// A probe that panics is reported as unknown rather than taking the endpoint
// down with it: somebody else's broken probe is not a reason to stop answering
// the question "is the database up".
func Snapshot() map[string]Report {
	mu.RLock()
	current := make([]probe, len(probes))
	copy(current, probes)
	mu.RUnlock()

	out := make(map[string]Report, len(current))
	for _, p := range current {
		out[p.name] = run(p)
	}
	return out
}

func run(p probe) (r Report) {
	defer func() {
		if rec := recover(); rec != nil {
			r = Unknown("the probe for this component panicked")
		}
	}()
	return p.fn()
}

// Reset clears the registry. For tests only: an application registers once.
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	probes = nil
}
