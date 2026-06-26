import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Most customer-facing apps ship with the same catalog
        endpoints: list categories, list products, view one product,
        filter products by category, show related products on a
        detail page. The previous lesson covered{' '}
        <em>why</em> some endpoints have to leave the{' '}
        <code>protected</code> route group; this one is the
        cheatsheet for <em>what</em> to wire up, using Category +
        Product as the worked example.
      </p>
      <p>
        For each endpoint below you get the same four pieces:
      </p>
      <ul>
        <li>
          <strong>What the generator already gave you</strong> --
          handler method, service method, route entry. Often the
          answer is &ldquo;already there, just move the route.&rdquo;
        </li>
        <li>
          <strong>What to add by hand</strong> -- the service +
          handler code when the endpoint is new, or the query-param
          plumbing when the existing handler needs to learn a new
          trick.
        </li>
        <li>
          <strong>Route registration</strong> -- which group the
          route belongs in.
        </li>
        <li>
          <strong>React Query hook</strong> -- the web-side hook
          customer pages call. The generator wrote a starter hook
          file; you extend it.
        </li>
      </ul>

      <h2>The starting point</h2>
      <p>
        Generate the two resources with the relationship:
      </p>
      <CodeBlock
        language="bash"
        code={`grit generate resource Category --fields=name:string!,image:file
grit generate resource Product --fields=name:string!,price:int,thumbnail:file,description:text,images:files,category:belongs_to:Category
grit migrate`}
      />
      <p>
        After the two generates, the routes file has every CRUD
        operation behind <code>protected</code> (the auth-required
        group). Open{' '}
        <code>apps/api/internal/routes/routes.go</code> and find
        the eight routes the generator created for you:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go (generated)"
        code={`protected.GET("/categories",       categoryHandler.List)
protected.GET("/categories/:id",   categoryHandler.GetByID)
protected.POST("/categories",      categoryHandler.Create)
protected.PUT("/categories/:id",   categoryHandler.Update)
protected.PATCH("/categories/:id", categoryHandler.Patch)

protected.GET("/products",         productHandler.List)
protected.GET("/products/:id",     productHandler.GetByID)
protected.POST("/products",        productHandler.Create)
protected.PUT("/products/:id",     productHandler.Update)
protected.PATCH("/products/:id",   productHandler.Patch)`}
      />
      <p>
        All ten lines need to move. The reads belong in a new
        public group; the writes belong in the admin group (which
        already exists in the scaffold for any resource the
        operator marked admin-only). Five reads -- four endpoints
        you&apos;ll expose publicly + four endpoints you&apos;ll
        add -- are the cheatsheet below.
      </p>

      <TipBox tone="info">
        Every &ldquo;already there&rdquo; section in this lesson
        refers to code the generator emitted on your first{' '}
        <code>grit generate resource</code>. If you regenerated
        with different fields the line numbers shift but the
        structure stays.
      </TipBox>

      {/* ─────────────────────────────────────────────────── */}
      <h2>1. List categories</h2>
      <p>
        <strong>Public</strong>. Anonymous visitors browse the
        category navigation, so this is the most-hit endpoint on a
        typical catalog.
      </p>

      <h3>What you already have</h3>
      <p>
        <code>categoryHandler.List</code> is generated and works as
        is -- pagination, search, sort, filter, all wired through
        the shared <code>paginate.List</code> helper.
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/category.go (already generated)"
        code={`func (h *CategoryHandler) List(c *gin.Context) {
    query := h.DB.Model(&models.Category{})

    res, err := paginate.List[models.Category](
        query,
        paginate.Bind(c),
        paginate.Config{
            Searchable: []string{"name", "slug"},
            Sortable:   map[string]bool{"id": true, "created_at": true, "name": true},
        },
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{"code": "INTERNAL_ERROR", "message": "Failed to fetch categories"},
        })
        return
    }
    c.JSON(http.StatusOK, res)
}`}
      />

      <h3>What to change</h3>
      <p>
        Move the route out of <code>protected</code> into a new{' '}
        <code>public</code> group:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go"
        code={`// PUBLIC: catalog. Anyone can browse.
public := r.Group("/api")
{
    public.GET("/categories",     categoryHandler.List)
    public.GET("/categories/:id", categoryHandler.GetByID)
    public.GET("/products",       productHandler.List)
    public.GET("/products/:id",   productHandler.GetByID)
    // ... add the four new endpoints below as you build them
}`}
      />

      <h3>React Query hook</h3>
      <p>
        The generator wrote{' '}
        <code>apps/web/hooks/use-categories.ts</code> with{' '}
        <code>useCategories(...)</code> already in it. No edit
        needed -- it just starts working once the route is public:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/(shop)/page.tsx"
        code={`"use client";
import { useCategories } from "@/hooks/use-categories";

export default function HomePage() {
  const { data, isLoading } = useCategories({ pageSize: 100 });
  if (isLoading) return <p>Loading…</p>;
  return (
    <nav className="grid grid-cols-2 sm:grid-cols-4 gap-3">
      {data?.data.map((c) => (
        <a key={c.id} href={"/category/" + c.slug} className="rounded-lg border p-4">
          {c.name}
        </a>
      ))}
    </nav>
  );
}`}
      />

      {/* ─────────────────────────────────────────────────── */}
      <h2>2. List products</h2>
      <p>
        <strong>Public</strong>. The product grid + the search bar.
        Both consumers (anonymous list page + admin) call this; the
        admin&apos;s axios just adds the auth cookie that the public
        version doesn&apos;t need.
      </p>

      <h3>What you already have</h3>
      <p>
        <code>productHandler.List</code> -- generated. Preloads the
        Category relation so the response includes nested category
        data; supports search, sort, pagination via the same{' '}
        <code>paginate.List</code> helper.
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/product.go (already generated)"
        code={`func (h *ProductHandler) List(c *gin.Context) {
    query := h.DB.Model(&models.Product{}).Preload("Category")

    res, err := paginate.List[models.Product](
        query,
        paginate.Bind(c),
        paginate.Config{
            Searchable: []string{"name", "slug", "description"},
            Sortable:   map[string]bool{"id": true, "created_at": true, "name": true, "price": true},
        },
    )
    if err != nil { /* 500 */ return }
    c.JSON(http.StatusOK, res)
}`}
      />

      <h3>What to change</h3>
      <p>
        Just move the route. The handler already does everything
        you need.
      </p>

      <h3>React Query hook</h3>
      <p>
        <code>useProducts(...)</code> in{' '}
        <code>apps/web/hooks/use-products.ts</code> works as is.
        Same shape as <code>useCategories</code>; the response{' '}
        <code>data.data[i].category</code> is the nested object
        thanks to the <code>Preload</code> above.
      </p>

      {/* ─────────────────────────────────────────────────── */}
      <h2>3. Get single product</h2>
      <p>
        <strong>Public</strong>. The product detail page.
      </p>

      <h3>What you already have</h3>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/product.go (already generated)"
        code={`func (h *ProductHandler) GetByID(c *gin.Context) {
    id := c.Param("id")
    var item models.Product
    if err := h.DB.Preload("Category").First(&item, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{"code": "NOT_FOUND", "message": "Product not found"},
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": item})
}`}
      />

      <h3>Bonus: lookup by slug too</h3>
      <p>
        Customer URLs read better as{' '}
        <code>/p/red-running-shoes</code> than{' '}
        <code>/p/01HXP...</code>. Add a second handler method that
        accepts either:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/product.go (add)"
        code={`// GetByIDOrSlug -- v3.31.49 cheatsheet. Tries ID first (UUIDs are
// 36 chars; slugs aren't), then slug. Lets customer URLs use the
// slug form without a second endpoint.
func (h *ProductHandler) GetByIDOrSlug(c *gin.Context) {
    key := c.Param("id")
    var item models.Product
    q := h.DB.Preload("Category")
    if len(key) == 36 {
        q = q.Where("id = ?", key)
    } else {
        q = q.Where("slug = ?", key)
    }
    if err := q.First(&item).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{"code": "NOT_FOUND", "message": "Product not found"},
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": item})
}`}
      />
      <p>
        Then route it:
      </p>
      <CodeBlock
        language="go"
        code={`public.GET("/products/:id", productHandler.GetByIDOrSlug)`}
      />

      <h3>React Query hook</h3>
      <p>
        <code>useGetProduct(id)</code> already exists. It treats the
        param as a generic identifier, so passing either an ID or a
        slug just works once the handler is the slug-aware one:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/p/[slug]/page.tsx"
        code={`"use client";
import { use } from "react";
import { useGetProduct } from "@/hooks/use-products";

export default function ProductDetailPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = use(params);
  const { data, isLoading, error } = useGetProduct(slug);
  if (isLoading) return <p>Loading…</p>;
  if (error || !data)  return <p>Product not found.</p>;
  return (
    <article>
      <h1>{data.name}</h1>
      <p className="text-lg">UGX {data.price.toLocaleString()}</p>
      <p className="text-text-muted">{data.category?.name}</p>
      <p>{data.description}</p>
    </article>
  );
}`}
      />

      {/* ─────────────────────────────────────────────────── */}
      <h2>4. Get single category</h2>
      <p>
        <strong>Public</strong>. The category landing page header
        (&ldquo;Phones &amp; Tablets&rdquo; with a hero image).
      </p>

      <h3>What you already have</h3>
      <p>
        <code>categoryHandler.GetByID</code> works as is. Add the
        same slug fallback you did for products if you want
        readable category URLs (<code>/category/phones</code> vs{' '}
        <code>/category/01HXP...</code>) -- pattern is identical.
      </p>

      {/* ─────────────────────────────────────────────────── */}
      <h2>5. List products by category</h2>
      <p>
        <strong>Public</strong>. The category landing page&apos;s
        product grid. Two ways to ship this -- pick one.
      </p>

      <h3>Option A: query param on the existing list endpoint (recommended)</h3>
      <p>
        Don&apos;t add a new endpoint. Teach{' '}
        <code>productHandler.List</code> a new filter:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/product.go (extend List)"
        code={`func (h *ProductHandler) List(c *gin.Context) {
    query := h.DB.Model(&models.Product{}).Preload("Category")

    // v3.31.49 cheatsheet -- filter by category. Accepts both
    // ?category_id=<uuid>  and  ?category=<slug> so the customer
    // app can pass whichever it has.
    if catID := c.Query("category_id"); catID != "" {
        query = query.Where("category_id = ?", catID)
    } else if slug := c.Query("category"); slug != "" {
        query = query.Joins("JOIN categories ON categories.id = products.category_id").
            Where("categories.slug = ?", slug)
    }

    res, err := paginate.List[models.Product](
        query,
        paginate.Bind(c),
        paginate.Config{
            Searchable: []string{"name", "slug", "description"},
            Sortable:   map[string]bool{"id": true, "created_at": true, "name": true, "price": true},
        },
    )
    if err != nil { /* 500 */ return }
    c.JSON(http.StatusOK, res)
}`}
      />

      <h3>React Query hook -- extend useProducts</h3>
      <CodeBlock
        language="ts"
        filename="apps/web/hooks/use-products.ts (extend)"
        code={`interface UseProductsParams {
  page?: number;
  pageSize?: number;
  search?: string;
  sortBy?: string;
  sortOrder?: string;
  // v3.31.49 cheatsheet -- either filter works; the API picks
  // whichever it receives. categoryId is faster (indexed PK);
  // categorySlug is friendlier in URLs.
  categoryId?: string;
  categorySlug?: string;
}

export function useProducts({
  page = 1, pageSize = 20, search = "",
  sortBy = "created_at", sortOrder = "desc",
  categoryId, categorySlug,
}: UseProductsParams = {}) {
  return useQuery<ProductsResponse>({
    queryKey: ["products", { page, pageSize, search, sortBy, sortOrder, categoryId, categorySlug }],
    queryFn: async () => {
      const params = new URLSearchParams({
        page: String(page), page_size: String(pageSize),
        sort_by: sortBy, sort_order: sortOrder,
      });
      if (search)       params.set("search", search);
      if (categoryId)   params.set("category_id", categoryId);
      if (categorySlug) params.set("category", categorySlug);
      const { data } = await apiClient.get(\`/api/products?\${params}\`);
      return data;
    },
  });
}`}
      />
      <p>
        Customer page:
      </p>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/category/[slug]/page.tsx"
        code={`"use client";
import { use } from "react";
import { useProducts } from "@/hooks/use-products";

export default function CategoryPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = use(params);
  const { data, isLoading } = useProducts({ categorySlug: slug, pageSize: 24 });
  if (isLoading) return <p>Loading…</p>;
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {data?.data.map((p) => <ProductCard key={p.id} product={p} />)}
    </div>
  );
}`}
      />

      <h3>Option B: nested route (alternative)</h3>
      <p>
        If you prefer{' '}
        <code>GET /api/categories/:slug/products</code> over a
        query param, add a new handler method:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/category.go (add)"
        code={`// Products returns the paginated product list for one category.
// Nested-route flavour of the same data the
// GET /products?category=<slug> form returns.
func (h *CategoryHandler) Products(c *gin.Context) {
    slug := c.Param("slug")

    var cat models.Category
    if err := h.DB.Where("slug = ?", slug).First(&cat).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{"code": "NOT_FOUND", "message": "Category not found"},
        })
        return
    }

    query := h.DB.Model(&models.Product{}).
        Where("category_id = ?", cat.ID).
        Preload("Category")

    res, err := paginate.List[models.Product](
        query,
        paginate.Bind(c),
        paginate.Config{
            Searchable: []string{"name", "slug", "description"},
            Sortable:   map[string]bool{"id": true, "created_at": true, "name": true, "price": true},
        },
    )
    if err != nil { /* 500 */ return }
    c.JSON(http.StatusOK, res)
}`}
      />
      <CodeBlock
        language="go"
        code={`public.GET("/categories/:slug/products", categoryHandler.Products)`}
      />

      <TipBox tone="info">
        Pick one or the other -- don&apos;t ship both. The two
        return identical JSON; running both just doubles the
        surface area you have to keep in sync. Option A wins for
        99% of catalogs because it reuses every search / sort /
        pagination feature the existing List handler already has.
      </TipBox>

      {/* ─────────────────────────────────────────────────── */}
      <h2>6. Related products on the detail page</h2>
      <p>
        <strong>Public</strong>. The &ldquo;You might also
        like&rdquo; carousel under a product detail page. The
        cheap good-enough heuristic: top-N products in the same
        category, excluding the one being viewed.
      </p>

      <h3>Add a service method</h3>
      <p>
        Business logic stays in the service so the handler can
        stay thin. The generator left an empty surface in{' '}
        <code>services/product.go</code> -- add this method below
        the existing ones:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/product.go (add)"
        code={`// Related returns up to N products in the same category as the
// passed product, ordered by created_at desc, excluding the
// product itself. The simplest heuristic that still feels
// relevant -- swap in a smarter scoring later (tags, popularity,
// purchased-together) without touching the handler.
func (s *ProductService) Related(productID string, limit int) ([]models.Product, error) {
    if limit <= 0 || limit > 50 {
        limit = 8
    }

    // Look up the source product's category. One query, no preload
    // needed -- we only want the category_id column.
    var source models.Product
    if err := s.DB.Select("id", "category_id").First(&source, "id = ?", productID).Error; err != nil {
        return nil, fmt.Errorf("source product not found: %w", err)
    }

    var items []models.Product
    err := s.DB.
        Where("category_id = ? AND id <> ?", source.CategoryID, source.ID).
        Order("created_at DESC").
        Limit(limit).
        Preload("Category").
        Find(&items).Error
    if err != nil {
        return nil, fmt.Errorf("fetching related products: %w", err)
    }
    return items, nil
}`}
      />

      <h3>Add a handler method</h3>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/product.go (add)"
        code={`// Related is the public &ldquo;you might also like&rdquo; carousel.
// Reads a ?limit= query param (default 8, capped at 50 in the service).
// Returns a flat array under data{} for symmetry with the rest of the
// catalog endpoints.
func (h *ProductHandler) Related(c *gin.Context) {
    id := c.Param("id")
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "8"))

    svc := &services.ProductService{DB: h.DB}
    items, err := svc.Related(id, limit)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error": gin.H{"code": "NOT_FOUND", "message": err.Error()},
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": items})
}`}
      />
      <p>
        Don&apos;t forget the new import in the handler file --{' '}
        <code>strconv</code> and{' '}
        <code>{`"<your-module>/internal/services"`}</code> (the
        latter is already imported by every generated handler, so
        nothing to add there).
      </p>

      <h3>Route it</h3>
      <CodeBlock
        language="go"
        code={`public.GET("/products/:id/related", productHandler.Related)`}
      />

      <h3>React Query hook</h3>
      <p>
        Add a sibling hook in{' '}
        <code>apps/web/hooks/use-products.ts</code>:
      </p>
      <CodeBlock
        language="ts"
        filename="apps/web/hooks/use-products.ts (add)"
        code={`export function useRelatedProducts(id: string, limit = 8) {
  return useQuery<{ data: Product[] }>({
    queryKey: ["products", id, "related", limit],
    queryFn: async () => {
      const { data } = await apiClient.get(
        \`/api/products/\${id}/related?limit=\${limit}\`,
      );
      return data;
    },
    enabled: !!id,
  });
}`}
      />

      <h3>Use it on the product detail page</h3>
      <CodeBlock
        language="tsx"
        filename="apps/web/app/p/[slug]/page.tsx (extend)"
        code={`const { data: product } = useGetProduct(slug);
const { data: related } = useRelatedProducts(product?.id ?? "");

// ... below the product detail ...
{related?.data && related.data.length > 0 && (
  <section className="mt-12">
    <h2 className="mb-4 text-xl font-semibold">You might also like</h2>
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {related.data.map((p) => <ProductCard key={p.id} product={p} />)}
    </div>
  </section>
)}`}
      />

      {/* ─────────────────────────────────────────────────── */}
      <h2>The full public route table</h2>
      <p>
        After all the moves + the two new endpoints, your{' '}
        <code>routes.go</code> public block should look like this:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go (final)"
        code={`// PUBLIC: catalog. Anyone can browse.
public := r.Group("/api")
{
    public.GET("/categories",                 categoryHandler.List)
    public.GET("/categories/:id",             categoryHandler.GetByID)
    public.GET("/products",                   productHandler.List)
    public.GET("/products/:id",               productHandler.GetByIDOrSlug)
    public.GET("/products/:id/related",       productHandler.Related)
}

protected := r.Group("/api")
protected.Use(middleware.Auth(db, authService))
{
    // Customer-account routes go here later (orders, profile, ...).
    // grit:routes:protected
}

admin := r.Group("/api")
admin.Use(middleware.Auth(db, authService))
admin.Use(middleware.RequireRole("ADMIN"))
{
    // Writes for the catalog -- staff only.
    admin.POST("/categories",         categoryHandler.Create)
    admin.PUT("/categories/:id",      categoryHandler.Update)
    admin.PATCH("/categories/:id",    categoryHandler.Patch)
    admin.DELETE("/categories/:id",   categoryHandler.Delete)

    admin.POST("/products",           productHandler.Create)
    admin.PUT("/products/:id",        productHandler.Update)
    admin.PATCH("/products/:id",      productHandler.Patch)
    admin.DELETE("/products/:id",     productHandler.Delete)
    // grit:routes:admin
}`}
      />

      <TipBox tone="warning">
        Make sure each route is registered in exactly one group.
        If you forget to delete the old <code>protected.GET</code>{' '}
        line after copying it to <code>public.GET</code>, Gin
        panics at startup with{' '}
        <code>handlers are already registered for path</code>.
      </TipBox>

      <KnowledgeCheck
        question="You moved `GET /products` from `protected` to `public` and the customer page loads. But hitting `POST /products` from the customer browser console still creates a product. What did you forget?"
        choices={[
          {
            label: 'Move the POST too -- writes were already on the public group',
            feedback:
              "Wrong direction. The generator put POST in `protected`, not `public`. The default protects writes; you need to make sure they're protected ENOUGH.",
          },
          {
            label: 'You left POST in `protected`, which only requires auth (any logged-in user) -- move it to `admin`, which adds the ADMIN role check',
            correct: true,
            feedback:
              "Right. `protected` only checks that SOMEONE is signed in. A regular customer with an account would pass that check. The `admin` group adds `middleware.RequireRole(\"ADMIN\")` on top -- that's what blocks non-staff from POSTing.",
          },
          {
            label: 'Add CSRF middleware to the route',
            feedback:
              "Wrong -- the AutoCSRF middleware is already mounted globally. The issue is authorisation, not the CSRF check.",
          },
          {
            label: "Restart the API; Gin caches routes",
            feedback:
              "Wrong -- Gin doesn't cache routes the way you'd need for this symptom. The handler genuinely ran; the route just lives in the wrong group.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Ship the full catalog in your jumia / ecom project:
            </p>
            <ol>
              <li>
                Move the four{' '}
                <code>GET /categories</code>,{' '}
                <code>GET /categories/:id</code>,{' '}
                <code>GET /products</code>,{' '}
                <code>GET /products/:id</code> routes from{' '}
                <code>protected</code> to a new <code>public</code>{' '}
                group.
              </li>
              <li>
                Move all <code>POST/PUT/PATCH/DELETE</code> for
                both resources from <code>protected</code> to{' '}
                <code>admin</code>.
              </li>
              <li>
                Extend <code>productHandler.List</code> to honour{' '}
                <code>?category_id=</code> and{' '}
                <code>?category=</code> filters.
              </li>
              <li>
                Add <code>GetByIDOrSlug</code> on{' '}
                <code>productHandler</code> and route it.
              </li>
              <li>
                Add{' '}
                <code>productService.Related</code> +{' '}
                <code>productHandler.Related</code>, route at{' '}
                <code>GET /products/:id/related</code>.
              </li>
              <li>
                Build customer pages:{' '}
                <code>/</code> (categories),{' '}
                <code>/category/[slug]</code> (products in category),{' '}
                <code>/p/[slug]</code> (product detail + related).
              </li>
            </ol>
            <p>
              Test as an anonymous visitor (incognito) -- you
              should see every page. Then try{' '}
              <code>POST /api/products</code> from the same
              incognito session -- it should return 401.
            </p>
          </>
        }
        hint={
          <>
            For step 6 you don&apos;t have to write{' '}
            <code>ProductCard</code> from scratch -- copy the card
            shape from the v3.31.31 Grit UI ecommerce registry (the{' '}
            <em>ecom-product-card</em> component). Drop it under{' '}
            <code>apps/web/components/product-card.tsx</code> and
            import it on both grid pages.
          </>
        }
        solution={
          <>
            <p>
              Three things commonly go wrong:
            </p>
            <ul>
              <li>
                <strong>Duplicate route panic.</strong> You left a
                copy of the route in <code>protected</code>. Search
                <code> routes.go</code> for each path -- each
                method+path appears exactly once.
              </li>
              <li>
                <strong>401 on the public page.</strong> The route
                still lives in <code>protected</code>. Restart{' '}
                <code>grit dev</code> after editing the routes
                file -- the API doesn&apos;t HMR.
              </li>
              <li>
                <strong>category is null on the product card.</strong>{' '}
                The new <code>List</code> filter accidentally
                replaced the <code>Preload(&quot;Category&quot;)</code>{' '}
                call. Make sure both the original{' '}
                <code>Preload</code> AND the new filters survive
                the merge.
              </li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Catalog reads are now public, writes stay locked behind{' '}
        <code>RequireRole(&quot;ADMIN&quot;)</code>, and the
        customer pages call typed React Query hooks. Next chapter
        switches to <em>going public</em> in the other direction
        -- letting an anonymous visitor submit a public form for{' '}
        <em>any</em> resource via a token-gated link, no auth
        required. See{' '}
        <code>concepts/generators/public-form-sharing</code>.
      </p>
    </>
  )
}
