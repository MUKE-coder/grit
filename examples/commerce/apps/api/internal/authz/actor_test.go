package authz

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
