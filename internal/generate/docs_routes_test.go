package generate

import (
	"strings"
	"testing"
)

// removeDocsRoutes deletes fluent call chains rather than brace-delimited
// blocks, so the interesting cases are all about where a chain ends and
// whether a neighbouring resource gets caught in the sweep.

const docsRoutesFixture = `package routes

func Setup() {
	docs.Route("POST /api/v1/auth/login").
		Summary("Sign in").
		RequestBody(handlers.LoginRequest{}).
		Response(200, handlers.AuthResponse{}, "Signed in")
	docs.Route("GET /api/v1/products").
		Summary("List products").
		Response(200, []models.Product{}, "A page of products")
	docs.Route("POST /api/v1/products").
		Summary("Create a product").
		RequestBody(handlers.CreateProductRequest{}).
		Response(201, models.Product{}, "Created")
	docs.Route("GET /api/v1/products/:id").
		Summary("Get one product").
		Response(200, models.Product{}, "The product")
	docs.Route("DELETE /api/v1/products/:id").
		Summary("Delete a product").
		Response(204, nil, "Deleted")
	docs.Route("GET /api/v1/orders").
		Summary("List orders").
		Response(200, []models.Order{}, "A page of orders")
	// grit:docs:routes:end
	log.Println("ready")
}
`

func TestRemoveDocsRoutes(t *testing.T) {
	t.Run("removes every chain for the resource", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", docsRoutesFixture)

		if err := removeDocsRoutes(f, "/api/v1/products"); err != nil {
			t.Fatalf("removeDocsRoutes: %v", err)
		}

		got := readFile(t, f)
		if strings.Contains(got, "/api/v1/products") {
			t.Errorf("product routes survived:\n%s", got)
		}
		// Every line of a chain must go, not just the docs.Route( line.
		for _, orphan := range []string{"List products", "Create a product", "CreateProductRequest"} {
			if strings.Contains(got, orphan) {
				t.Errorf("orphaned chain line %q left behind:\n%s", orphan, got)
			}
		}
	})

	t.Run("leaves other resources and core routes alone", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", docsRoutesFixture)

		if err := removeDocsRoutes(f, "/api/v1/products"); err != nil {
			t.Fatalf("removeDocsRoutes: %v", err)
		}

		got := readFile(t, f)
		for _, keep := range []string{
			`docs.Route("POST /api/v1/auth/login")`,
			`docs.Route("GET /api/v1/orders")`,
			"handlers.LoginRequest{}",
			"// grit:docs:routes:end",
			`log.Println("ready")`,
		} {
			if !strings.Contains(got, keep) {
				t.Errorf("removed too much — %q is gone:\n%s", keep, got)
			}
		}
	})

	t.Run("a plural that prefixes another resource is not over-matched", func(t *testing.T) {
		// "order" must not take "orders" with it: the match includes the
		// closing quote precisely so this cannot happen.
		f := writeTempFile(t, "routes.go", `package routes

	docs.Route("GET /api/v1/order").
		Summary("List order").
		Response(200, []models.Order{}, "orders")
	docs.Route("GET /api/v1/orders").
		Summary("List orders").
		Response(200, []models.Orders{}, "many")
`)

		if err := removeDocsRoutes(f, "/api/v1/order"); err != nil {
			t.Fatalf("removeDocsRoutes: %v", err)
		}

		got := readFile(t, f)
		if strings.Contains(got, `docs.Route("GET /api/v1/order")`) {
			t.Error("the targeted route survived")
		}
		if !strings.Contains(got, `docs.Route("GET /api/v1/orders")`) {
			t.Errorf("the neighbouring plural was swept up too:\n%s", got)
		}
	})

	t.Run("reports when there is nothing to remove", func(t *testing.T) {
		f := writeTempFile(t, "routes.go", docsRoutesFixture)

		if err := removeDocsRoutes(f, "/api/v1/widgets"); err == nil {
			t.Error("expected an error for a resource with no docs routes")
		}
	})
}
