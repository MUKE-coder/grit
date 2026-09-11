package scaffold

// The actor: who a query runs for, carried on the context.
//
// Generated services own every query now, and a service has no gin request:
// the same method runs from a handler, a job, a command or a test. Ownership
// scoping used to read the caller straight off the gin context
// (authz.ScopeToOwner(c, ...)), which only a handler has. These helpers read
// it from context.Context instead. The handler puts the caller there with
// WithActor; a job states that it acts as the application with AsSystem; and a
// context with no actor at all sees no owned rows, so forgetting it fails
// closed rather than open.

// APIAuthzActorGo is exported so the generator can add the file to a project
// from before it existed: generated services call authz.ScopeOwned.
func APIAuthzActorGo() string { return authzActorGo() }

func authzActorGo() string {
	return `package authz

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Actor is who a query runs for: the signed-in user, and whether ownership
// scoping should let them through, as it does an ADMIN.
type Actor struct {
	UserID string
	Admin  bool
}

type actorKey struct{}

// ActorOf is the caller of a request, as the auth middleware identified them.
func ActorOf(c *gin.Context) Actor {
	return Actor{UserID: CurrentUserID(c), Admin: IsAdmin(c)}
}

// WithActor returns ctx carrying a.
func WithActor(ctx context.Context, a Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, a)
}

// ActorFrom returns the actor on ctx, if there is one.
func ActorFrom(ctx context.Context) (Actor, bool) {
	a, ok := ctx.Value(actorKey{}).(Actor)
	return a, ok
}

// AsSystem marks ctx as the application acting on its own behalf, for a job or
// a command. Ownership scoping lets it through, as it does an ADMIN.
func AsSystem(ctx context.Context) context.Context {
	return WithActor(ctx, Actor{Admin: true})
}

// UserIDFrom is the acting user's id, or "".
func UserIDFrom(ctx context.Context) string {
	a, _ := ActorFrom(ctx)
	return a.UserID
}

// ScopeOwned narrows a query to the rows the actor on ctx owns. An ADMIN or
// the system gets the query back untouched. A context with no actor, or an
// actor with no user, matches nothing: a scope that opened up when it could
// not tell who was asking would read as protection and give none.
//
// column is a fixed string from the generator, never user input.
func ScopeOwned(ctx context.Context, q *gorm.DB, column string) *gorm.DB {
	a, ok := ActorFrom(ctx)
	switch {
	case ok && a.Admin:
		return q
	case !ok || a.UserID == "":
		return q.Where("1 = 0")
	}
	return q.Where(column+" = ?", a.UserID)
}

// Owns reports whether the actor on ctx may see row: an ADMIN, the system, or
// the row's owner.
func Owns(ctx context.Context, row Ownable) bool {
	a, ok := ActorFrom(ctx)
	if !ok {
		return false
	}
	return a.Admin || (a.UserID != "" && row.GetOwnerID() == a.UserID)
}
`
}

func authzActorTestGo() string {
	return `package authz

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type actorTestRow struct{ owner string }

func (r actorTestRow) GetOwnerID() string { return r.owner }

// A service has no request, only a context. Ownership has to hold there too,
// and fail closed when nobody said who is asking.
func TestOwnershipFromTheContext(t *testing.T) {
	bg := context.Background()
	if Owns(bg, actorTestRow{"u1"}) {
		t.Error("a context with no actor saw an owned row")
	}

	alice := WithActor(bg, Actor{UserID: "u1"})
	if !Owns(alice, actorTestRow{"u1"}) {
		t.Error("an owner could not see their own row")
	}
	if Owns(alice, actorTestRow{"u2"}) {
		t.Error("a user saw somebody else's row")
	}
	if !Owns(AsSystem(bg), actorTestRow{"u2"}) {
		t.Error("the system could not see a row")
	}
	if UserIDFrom(alice) != "u1" || UserIDFrom(bg) != "" {
		t.Error("UserIDFrom read the wrong user")
	}

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "u9")
	c.Set("user_role", "ADMIN")
	if a := ActorOf(c); a.UserID != "u9" || !a.Admin {
		t.Errorf("ActorOf read %+v from the request", a)
	}
}
`
}
