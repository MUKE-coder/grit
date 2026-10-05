package authz

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
