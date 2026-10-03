package authz

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// # Policies: the rules that ownership and permissions cannot express
//
// This project already answers two authorization questions. The route's
// middleware answers "may this user touch widgets at all", from the permission
// catalogue. The --owned-by scoping answers "is this their row", from the
// owner column. Between them they cover most of what an application needs.
//
// What neither can say is anything conditional on the record itself. A post may
// be edited while it is a draft and not after it is published. An order may be
// approved by anyone except the person who raised it. A document may be read by
// its author, its reviewers, and nobody else. Those are rules about this row and
// this user together, and before policies they were hand-written ifs at the top
// of a handler, which is where they stop being findable and start being
// forgotten in the second handler that needs them.
//
// # Absence is not denial, and why that is right here
//
// Inspect returns Allow when no rule is defined for an ability. That is
// deliberate, and it is the opposite of what a security layer usually wants.
//
// A gate is a *narrowing*. The route's permission has already run and already
// said yes; this decides whether this particular record is an exception. An
// undefined rule means "no exception", which is exactly how the project behaved
// before policies existed. Failing closed instead would mean that delivering
// this file to an existing project denied every write in it, which is not a
// safe default, it is an outage.
//
// So a rule is a restriction you add, never a permission you forget to grant.
// If you want a deny-by-default ability, write the rule that denies.
//
// # Where it runs
//
// A generated service checks three abilities, around the row it has just
// loaded: <resource>.read in GetByID, <resource>.update before an update or a
// patch, and <resource>.delete before a delete. A denial comes back to the
// handler as a DeniedError and is answered with 403 and the rule's own reason,
// which is the part ownership checks cannot do: a 404 that will not say why is
// right for hiding a row's existence, and wrong for telling somebody their post
// is already published.

// Decision is a rule's answer, and the reason for it.
//
// The reason is the point. A permission check can only say no, so a user who
// cannot do something is told that they cannot do it, and files a ticket. A
// rule that knows why can say "this post was published on Tuesday and published
// posts are edited through a revision", which is an answer rather than a wall.
type Decision struct {
	allowed bool
	reason  string
}

// Allow permits the action.
func Allow() Decision { return Decision{allowed: true} }

// Deny refuses it, and says why in words the person reading them can act on.
//
// The reason reaches the user, so write it for them: not "policy violation" but
// "this order was raised by you, and an approver cannot be the raiser".
func Deny(reason string) Decision { return Decision{reason: reason} }

// Allowed reports whether the action may proceed.
func (d Decision) Allowed() bool { return d.allowed }

// Reason is why it may not. Empty when it may.
func (d Decision) Reason() string {
	if d.reason == "" && !d.allowed {
		return "You cannot do that to this record"
	}
	return d.reason
}

// Rule decides one ability over one record.
//
// record is the row as the service loaded it, so a rule can read any column.
// It is any rather than a generic, because the registry holds rules for every
// resource at once and a map cannot be keyed by type; the rule asserts to the
// model it is written for and returns Allow for anything else, which is what
// the generated stub does.
type Rule func(actor Actor, record any) Decision

// BeforeFunc runs ahead of every rule. Returning handled=true answers the
// question outright and the rule is never consulted.
//
// This is where a blanket exemption belongs, so it is written once instead of
// at the top of every rule, where the fifteenth copy is the one that forgets it.
type BeforeFunc func(actor Actor, ability string, record any) (decision Decision, handled bool)

var (
	mu     sync.RWMutex
	rules  = map[string]Rule{}
	before []BeforeFunc
)

// Define registers the rule for an ability, such as "posts.update".
//
// Called from internal/policies during startup. Defining the same ability twice
// panics rather than silently keeping one of them: two rules for one ability is
// a merge gone wrong, and the half that loses would be invisible.
func Define(ability string, rule Rule) {
	ability = strings.TrimSpace(ability)
	if ability == "" || rule == nil {
		panic("authz: Define needs an ability and a rule")
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := rules[ability]; exists {
		panic(fmt.Sprintf("authz: %q already has a rule; one ability, one rule", ability))
	}
	rules[ability] = rule
}

// Before registers a check that runs ahead of every rule.
func Before(fn BeforeFunc) {
	if fn == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	before = append(before, fn)
}

// Inspect asks whether the actor on ctx may do ability to record.
//
// Allow when no rule is defined. See the note at the top of this file: a gate
// narrows a decision the route's permission has already made, so an undefined
// rule means no narrowing.
func Inspect(ctx context.Context, ability string, record any) Decision {
	actor, _ := ActorFrom(ctx)

	mu.RLock()
	hooks := make([]BeforeFunc, len(before))
	copy(hooks, before)
	rule, defined := rules[ability]
	mu.RUnlock()

	for _, hook := range hooks {
		if decision, handled := hook(actor, ability, record); handled {
			return decision
		}
	}
	if !defined {
		return Allow()
	}
	return rule(actor, record)
}

// Allows is Inspect without the reason, for a caller that only branches.
func Allows(ctx context.Context, ability string, record any) bool {
	return Inspect(ctx, ability, record).Allowed()
}

// Abilities lists every ability that has a rule, sorted. For `grit doctor` and
// for a test that wants to assert the set did not change by accident.
func Abilities() []string {
	mu.RLock()
	defer mu.RUnlock()

	out := make([]string, 0, len(rules))
	for ability := range rules {
		out = append(out, ability)
	}
	sort.Strings(out)
	return out
}

// DeniedError is a refused action travelling from a service back to a handler.
//
// A typed error rather than a bare one, because the handler has to tell this
// apart from "not found" and from a database failure, and because the reason
// has to survive the trip: it is the thing the user reads.
type DeniedError struct {
	Ability string
	Reason  string
}

func (e *DeniedError) Error() string {
	return fmt.Sprintf("authz: %s denied: %s", e.Ability, e.Reason)
}

// Denied builds the error for a decision that refused.
func Denied(ability string, d Decision) error {
	return &DeniedError{Ability: ability, Reason: d.Reason()}
}

// ResetForTest clears every rule. Only for tests, which register their own and
// would otherwise inherit whatever the last one defined.
func ResetForTest() {
	mu.Lock()
	defer mu.Unlock()
	rules = map[string]Rule{}
	before = nil
}
