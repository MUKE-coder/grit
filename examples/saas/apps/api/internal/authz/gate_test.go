package authz

import (
	"context"
	"testing"
)

type fakeRecord struct {
	Status string
}

func reset(t *testing.T) {
	t.Helper()
	ResetForTest()
	t.Cleanup(ResetForTest)
}

func actorCtx(userID string, admin bool) context.Context {
	return WithActor(context.Background(), Actor{UserID: userID, Admin: admin})
}

// The property the whole design rests on. A gate narrows a decision the route's
// permission has already made, so an ability nobody wrote a rule for behaves
// exactly as the project did before it had policies. Failing closed instead
// would mean that delivering this file to an existing project refused every
// write in it, which is not a safe default, it is an outage.
func TestAnUndefinedAbilityAllows(t *testing.T) {
	reset(t)

	if d := Inspect(actorCtx("u1", false), "widgets.update", &fakeRecord{}); !d.Allowed() {
		t.Errorf("an ability with no rule denied: %q", d.Reason())
	}
}

func TestARuleDecidesAndSaysWhy(t *testing.T) {
	reset(t)

	Define("widgets.update", func(actor Actor, record any) Decision {
		item, ok := record.(*fakeRecord)
		if !ok {
			return Allow()
		}
		if item.Status == "draft" {
			return Allow()
		}
		return Deny("published widgets are edited through a revision")
	})

	ctx := actorCtx("u1", false)

	if d := Inspect(ctx, "widgets.update", &fakeRecord{Status: "draft"}); !d.Allowed() {
		t.Errorf("a draft was refused: %q", d.Reason())
	}

	d := Inspect(ctx, "widgets.update", &fakeRecord{Status: "published"})
	if d.Allowed() {
		t.Fatal("a published row was allowed")
	}
	// The reason is the point: a refusal that cannot explain itself produces a
	// support ticket rather than an answer.
	if d.Reason() != "published widgets are edited through a revision" {
		t.Errorf("reason = %q", d.Reason())
	}
}

// A denial with no reason still has to say something, because the string is
// rendered to a person.
func TestADenialAlwaysHasWords(t *testing.T) {
	reset(t)

	Define("widgets.delete", func(Actor, any) Decision { return Deny("") })

	d := Inspect(actorCtx("u1", false), "widgets.delete", &fakeRecord{})
	if d.Allowed() || d.Reason() == "" {
		t.Errorf("Decision = %+v, want a refusal with words", d)
	}
}

// Before is where a blanket exemption belongs, so it is written once instead of
// at the top of every rule, where the fifteenth copy is the one that forgets it.
func TestBeforeAnswersOutright(t *testing.T) {
	reset(t)

	Define("widgets.update", func(Actor, any) Decision { return Deny("never") })
	Before(func(actor Actor, ability string, record any) (Decision, bool) {
		if actor.Admin {
			return Allow(), true
		}
		return Decision{}, false
	})

	if d := Inspect(actorCtx("admin", true), "widgets.update", &fakeRecord{}); !d.Allowed() {
		t.Error("the Before hook did not exempt the admin")
	}
	if d := Inspect(actorCtx("u1", false), "widgets.update", &fakeRecord{}); d.Allowed() {
		t.Error("the Before hook exempted somebody it should not have")
	}
}

// Two rules for one ability is a merge gone wrong, and the half that loses
// would be invisible.
func TestDefiningAnAbilityTwicePanics(t *testing.T) {
	reset(t)

	Define("widgets.update", func(Actor, any) Decision { return Allow() })

	defer func() {
		if recover() == nil {
			t.Error("defining the same ability twice was accepted")
		}
	}()
	Define("widgets.update", func(Actor, any) Decision { return Deny("other") })
}

// A rule that reaches a record it was not written for is a wiring mistake, not
// a refusal: the generated stub returns Allow, and this is the test that says
// why that is the right shape.
func TestARuleSeesTheActor(t *testing.T) {
	reset(t)

	Define("widgets.delete", func(actor Actor, record any) Decision {
		if actor.UserID == "owner" {
			return Allow()
		}
		return Deny("only the owner may delete this")
	})

	if d := Inspect(actorCtx("owner", false), "widgets.delete", &fakeRecord{}); !d.Allowed() {
		t.Error("the owner was refused")
	}
	if d := Inspect(actorCtx("someone-else", false), "widgets.delete", &fakeRecord{}); d.Allowed() {
		t.Error("a stranger was allowed")
	}
}

func TestDeniedErrorCarriesTheReasonToTheHandler(t *testing.T) {
	err := Denied("widgets.update", Deny("it is live"))

	var denied *DeniedError
	if !asDenied(err, &denied) {
		t.Fatalf("Denied returned %T, which a handler cannot tell from any other error", err)
	}
	if denied.Ability != "widgets.update" || denied.Reason != "it is live" {
		t.Errorf("DeniedError = %+v", denied)
	}
}

// errors.As, written out so this file does not import errors for one call.
func asDenied(err error, target **DeniedError) bool {
	d, ok := err.(*DeniedError)
	if ok {
		*target = d
	}
	return ok
}

func TestAbilitiesListsWhatIsDefined(t *testing.T) {
	reset(t)

	Define("widgets.update", func(Actor, any) Decision { return Allow() })
	Define("posts.read", func(Actor, any) Decision { return Allow() })

	got := Abilities()
	if len(got) != 2 || got[0] != "posts.read" || got[1] != "widgets.update" {
		t.Errorf("Abilities() = %v, want them sorted", got)
	}
}
