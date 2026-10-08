package authz

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireRolesAllowsListed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_role", "admin"); c.Next() })
	r.GET("/admin", RequireRoles("admin"), func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRequireRolesBlocksOthers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) { c.Set("user_role", "user"); c.Next() })
	r.GET("/admin", RequireRoles("admin"), func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin", nil)
	r.ServeHTTP(w, req)
	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
