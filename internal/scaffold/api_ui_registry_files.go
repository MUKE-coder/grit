package scaffold

// apiUIComponentModelGo returns the UIComponent GORM model.
func apiUIComponentModelGo() string {
	return `package models

import "time"

// UIComponent stores a shadcn-compatible UI component in the database.
// Components are served as a public registry at GET /r.json and GET /r/:name.json.
type UIComponent struct {
	ID           uint      ` + "`" + `gorm:"primarykey" json:"id"` + "`" + `
	Name         string    ` + "`" + `gorm:"size:100;uniqueIndex" json:"name"` + "`" + `        // slug: "hero-01"
	DisplayName  string    ` + "`" + `gorm:"size:200" json:"display_name"` + "`" + `            // "Hero Section"
	Description  string    ` + "`" + `gorm:"type:text" json:"description"` + "`" + `
	Category     string    ` + "`" + `gorm:"size:50;index" json:"category"` + "`" + `           // "marketing","ecommerce","layout"
	Tags         string    ` + "`" + `gorm:"type:text" json:"tags"` + "`" + `                  // JSON: ["hero","landing"]
	Files        string    ` + "`" + `gorm:"type:text" json:"files"` + "`" + `                 // JSON: shadcn files array
	Dependencies string    ` + "`" + `gorm:"type:text" json:"dependencies"` + "`" + `          // JSON: ["lucide-react"]
	RegistryDeps string    ` + "`" + `gorm:"type:text" json:"registry_deps"` + "`" + `         // JSON: other grit-ui component deps
	PreviewCode  string    ` + "`" + `gorm:"type:text" json:"preview_code"` + "`" + `          // short JSX snippet for browser
	IsPublic     bool      ` + "`" + `gorm:"default:true" json:"is_public"` + "`" + `
	CreatedAt    time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt    time.Time ` + "`" + `json:"updated_at"` + "`" + `
}
`
}

// apiUIComponentSeedGo returns the seed file that populates the DB with starter components.
func apiUIComponentSeedGo() string {
	return `package models

import (
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

// SeedUIComponents inserts the 4 starter Grit UI components if the table is empty.
// Safe to call on every startup — skips silently if any component already exists.
func SeedUIComponents(db *gorm.DB) {
	var count int64
	db.Model(&UIComponent{}).Count(&count)
	if count > 0 {
		return
	}

	// filesJSON builds the shadcn-compatible files array for a single TSX component.
	filesJSON := func(name, tsx string) string {
		data, _ := json.Marshal([]map[string]string{{
			"path":    "components/grit-ui/" + name + ".tsx",
			"type":    "registry:component",
			"target":  "components/grit-ui/" + name + ".tsx",
			"content": tsx,
		}})
		return string(data)
	}

	components := []UIComponent{
		{
			Name:         "hero-01",
			DisplayName:  "Hero Section",
			Description:  "Bold hero section with gradient text, announcement badge, and dual CTA buttons",
			Category:     "marketing",
			Tags:         ` + "`" + `["hero","landing","marketing"]` + "`" + `,
			Files:        filesJSON("hero-01", hero01TSX),
			Dependencies: ` + "`" + `[]` + "`" + `,
			RegistryDeps: ` + "`" + `[]` + "`" + `,
			PreviewCode:  "<Hero01 />",
			IsPublic:     true,
		},
		{
			Name:         "pricing-card-01",
			DisplayName:  "Pricing Card",
			Description:  "Pricing tier card with feature list, popular badge, and gradient CTA button",
			Category:     "marketing",
			Tags:         ` + "`" + `["pricing","card","saas"]` + "`" + `,
			Files:        filesJSON("pricing-card-01", pricingCard01TSX),
			Dependencies: ` + "`" + `["lucide-react"]` + "`" + `,
			RegistryDeps: ` + "`" + `[]` + "`" + `,
			PreviewCode:  ` + "`" + `<PricingCard01 name="Pro" price={49} period="month" description="For growing teams" features={["Unlimited projects"]} ctaLabel="Get Started" ctaHref="#" />` + "`" + `,
			IsPublic:     true,
		},
		{
			Name:         "product-card-01",
			DisplayName:  "Product Card",
			Description:  "E-commerce product card with image, star rating, hover overlay, and add-to-cart button",
			Category:     "ecommerce",
			Tags:         ` + "`" + `["product","card","ecommerce","shop"]` + "`" + `,
			Files:        filesJSON("product-card-01", productCard01TSX),
			Dependencies: ` + "`" + `["lucide-react"]` + "`" + `,
			RegistryDeps: ` + "`" + `[]` + "`" + `,
			PreviewCode:  ` + "`" + `<ProductCard01 name="Wireless Headphones" price={99} rating={4.5} category="Electronics" />` + "`" + `,
			IsPublic:     true,
		},
		{
			Name:         "feature-grid-01",
			DisplayName:  "Feature Grid",
			Description:  "Feature showcase section with icon cards in a responsive 2-col/3-col grid",
			Category:     "marketing",
			Tags:         ` + "`" + `["features","grid","cards","marketing"]` + "`" + `,
			Files:        filesJSON("feature-grid-01", featureGrid01TSX),
			Dependencies: ` + "`" + `["lucide-react"]` + "`" + `,
			RegistryDeps: ` + "`" + `[]` + "`" + `,
			PreviewCode:  ` + "`" + `<FeatureGrid01 title="Why Grit?" subtitle="Build faster." features={[]} />` + "`" + `,
			IsPublic:     true,
		},
	}

	if err := db.Create(&components).Error; err != nil {
		log.Printf("Warning: failed to seed UI components: %v", err)
	} else {
		log.Printf("Seeded %d Grit UI components", len(components))
	}
}

// ── Component TSX source ──────────────────────────────────────────────────────

const hero01TSX = ` + "`" + `import Link from "next/link";

export function Hero01() {
  return (
    <section className="relative overflow-hidden bg-[#0a0a0f] py-24 sm:py-32">
      <div className="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(ellipse_80%_50%_at_50%_-20%,rgba(108,92,231,0.3),transparent)]" />
      <div className="mx-auto max-w-5xl px-6 text-center">
        <span className="inline-flex items-center gap-2 rounded-full border border-[#2a2a3a] bg-[#111118] px-4 py-1.5 text-sm text-[#9090a8]">
          <span className="inline-block h-2 w-2 rounded-full bg-[#6c5ce7]" />
          Now in beta &mdash; built with Go + React
        </span>
        <h1 className="mt-8 text-5xl font-bold tracking-tight text-[#e8e8f0] sm:text-7xl">
          Ship faster with{" "}
          <span className="bg-gradient-to-r from-[#6c5ce7] to-[#74b9ff] bg-clip-text text-transparent">
            Grit
          </span>
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-lg text-[#9090a8]">
          A full-stack meta-framework that fuses Go with Next.js. One command to scaffold
          a complete production-ready app.
        </p>
        <div className="mt-10 flex flex-col items-center gap-4 sm:flex-row sm:justify-center">
          <Link
            href="/register"
            className="rounded-lg bg-[#6c5ce7] px-6 py-3 text-sm font-semibold text-white shadow-lg transition hover:bg-[#7c6cf7]"
          >
            Get started free
          </Link>
          <Link
            href="/docs"
            className="rounded-lg border border-[#2a2a3a] bg-[#111118] px-6 py-3 text-sm font-semibold text-[#e8e8f0] transition hover:bg-[#1a1a24]"
          >
            View docs &rarr;
          </Link>
        </div>
      </div>
    </section>
  );
}
` + "`" + `

const pricingCard01TSX = ` + "`" + `import { Check } from "lucide-react";

interface PricingCard01Props {
  name: string;
  price: number;
  period: string;
  description: string;
  features: string[];
  isPopular?: boolean;
  ctaLabel: string;
  ctaHref: string;
}

export function PricingCard01({
  name, price, period, description, features, isPopular, ctaLabel, ctaHref,
}: PricingCard01Props) {
  const popular = isPopular ?? false;
  return (
    <div className={"relative rounded-2xl border bg-[#111118] p-8 " + (popular ? "border-[#6c5ce7] shadow-lg shadow-[#6c5ce7]/20" : "border-[#2a2a3a]")}>
      {popular && (
        <span className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-[#6c5ce7] px-4 py-1 text-xs font-semibold text-white">
          Most Popular
        </span>
      )}
      <h3 className="text-lg font-semibold text-[#e8e8f0]">{name}</h3>
      <p className="mt-1 text-sm text-[#9090a8]">{description}</p>
      <div className="mt-6 flex items-baseline gap-1">
        <span className="text-4xl font-bold text-[#e8e8f0]">${price}</span>
        <span className="text-[#9090a8]">/{period}</span>
      </div>
      <ul className="mt-6 space-y-3">
        {features.map((f) => (
          <li key={f} className="flex items-center gap-2 text-sm text-[#9090a8]">
            <Check className="h-4 w-4 text-[#00b894]" />
            {f}
          </li>
        ))}
      </ul>
      <a
        href={ctaHref}
        className={"mt-8 block rounded-lg px-4 py-3 text-center text-sm font-semibold transition " + (popular ? "bg-[#6c5ce7] text-white hover:bg-[#7c6cf7]" : "border border-[#2a2a3a] text-[#e8e8f0] hover:bg-[#1a1a24]")}
      >
        {ctaLabel}
      </a>
    </div>
  );
}
` + "`" + `

const productCard01TSX = ` + "`" + `import { Star, ShoppingCart } from "lucide-react";

interface ProductCard01Props {
  name: string;
  price: number;
  originalPrice?: number;
  rating: number;
  category: string;
  imageUrl?: string;
  href?: string;
}

export function ProductCard01({ name, price, originalPrice, rating, category, imageUrl, href = "#" }: ProductCard01Props) {
  return (
    <a href={href} className="group block overflow-hidden rounded-2xl border border-[#2a2a3a] bg-[#111118] transition hover:border-[#6c5ce7]/50">
      <div className="relative aspect-square overflow-hidden bg-[#1a1a24]">
        {imageUrl ? (
          <img src={imageUrl} alt={name} className="h-full w-full object-cover transition group-hover:scale-105" />
        ) : (
          <div className="flex h-full items-center justify-center">
            <ShoppingCart className="h-16 w-16 text-[#606078]" />
          </div>
        )}
        <div className="absolute inset-0 bg-black/0 transition group-hover:bg-black/20" />
      </div>
      <div className="p-4">
        <p className="text-xs text-[#6c5ce7]">{category}</p>
        <h3 className="mt-1 font-semibold text-[#e8e8f0] line-clamp-2">{name}</h3>
        <div className="mt-2 flex items-center gap-1">
          {Array.from({ length: 5 }).map((_, i) => (
            <Star key={i} className={"h-3 w-3 " + (i < Math.floor(rating) ? "fill-[#fdcb6e] text-[#fdcb6e]" : "text-[#606078]")} />
          ))}
          <span className="ml-1 text-xs text-[#9090a8]">{rating}</span>
        </div>
        <div className="mt-3 flex items-center justify-between">
          <div className="flex items-baseline gap-2">
            <span className="text-lg font-bold text-[#e8e8f0]">${price}</span>
            {originalPrice && <span className="text-sm text-[#9090a8] line-through">${originalPrice}</span>}
          </div>
          <button className="rounded-lg bg-[#6c5ce7] px-3 py-1.5 text-xs font-semibold text-white transition hover:bg-[#7c6cf7]">
            Add to cart
          </button>
        </div>
      </div>
    </a>
  );
}
` + "`" + `

const featureGrid01TSX = ` + "`" + `import type { LucideIcon } from "lucide-react";

interface Feature {
  icon: LucideIcon;
  title: string;
  description: string;
}

interface FeatureGrid01Props {
  title: string;
  subtitle: string;
  features: Feature[];
}

export function FeatureGrid01({ title, subtitle, features }: FeatureGrid01Props) {
  return (
    <section className="bg-[#0a0a0f] py-24">
      <div className="mx-auto max-w-6xl px-6">
        <div className="text-center">
          <h2 className="text-4xl font-bold text-[#e8e8f0]">{title}</h2>
          <p className="mt-4 text-lg text-[#9090a8]">{subtitle}</p>
        </div>
        <div className="mt-16 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => (
            <div key={feature.title} className="rounded-2xl border border-[#2a2a3a] bg-[#111118] p-6 transition hover:border-[#6c5ce7]/50">
              <div className="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-[#6c5ce7]/10">
                <feature.icon className="h-5 w-5 text-[#6c5ce7]" />
              </div>
              <h3 className="mt-4 font-semibold text-[#e8e8f0]">{feature.title}</h3>
              <p className="mt-2 text-sm text-[#9090a8]">{feature.description}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
` + "`" + `
`
}

// apiUIRegistryHandlerGo returns the UIRegistry handler for public + admin endpoints.
func apiUIRegistryHandlerGo() string {
	return `package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/models"
)

// UIRegistryHandler serves the shadcn-compatible component registry.
type UIRegistryHandler struct {
	DB     *gorm.DB
	APIURL string // e.g. "http://localhost:8080"
}

// NewUIRegistryHandler constructs the handler.
func NewUIRegistryHandler(db *gorm.DB, apiURL string) *UIRegistryHandler {
	return &UIRegistryHandler{DB: db, APIURL: apiURL}
}

// ── Public registry endpoints ────────────────────────────────────────────────

// GetRegistry serves GET /r.json — shadcn-compatible root registry.
func (h *UIRegistryHandler) GetRegistry(c *gin.Context) {
	var components []models.UIComponent
	if err := h.DB.Where("is_public = ?", true).
		Order("category, name").
		Find(&components).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch registry"})
		return
	}

	items := make([]gin.H, 0, len(components))
	for _, comp := range components {
		items = append(items, gin.H{
			"name":        comp.Name,
			"type":        "registry:block",
			"title":       comp.DisplayName,
			"description": comp.Description,
			"category":    comp.Category,
			"registryUrl": h.APIURL + "/r/" + comp.Name + ".json",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"$schema":  "https://ui.shadcn.com/schema/registry.json",
		"name":     "grit-ui",
		"homepage": "https://gritframework.dev/components",
		"items":    items,
	})
}

// GetComponent serves GET /r/:name.json — full shadcn component JSON for npx shadcn add.
func (h *UIRegistryHandler) GetComponent(c *gin.Context) {
	name := c.Param("name")
	// Strip .json suffix if present
	if len(name) > 5 && name[len(name)-5:] == ".json" {
		name = name[:len(name)-5]
	}

	var comp models.UIComponent
	if err := h.DB.Where("name = ? AND is_public = ?", name, true).First(&comp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "component not found"})
		return
	}

	var files []gin.H
	_ = json.Unmarshal([]byte(comp.Files), &files)

	var dependencies []string
	_ = json.Unmarshal([]byte(comp.Dependencies), &dependencies)

	var registryDeps []string
	_ = json.Unmarshal([]byte(comp.RegistryDeps), &registryDeps)

	c.JSON(http.StatusOK, gin.H{
		"$schema":              "https://ui.shadcn.com/schema/registry-item.json",
		"name":                 comp.Name,
		"type":                 "registry:block",
		"title":                comp.DisplayName,
		"description":          comp.Description,
		"files":                files,
		"dependencies":         dependencies,
		"registryDependencies": registryDeps,
	})
}

// ── Authenticated list endpoints ──────────────────────────────────────────────

// ListComponents serves GET /api/ui-components.
func (h *UIRegistryHandler) ListComponents(c *gin.Context) {
	category := c.Query("category")

	type componentSummary struct {
		ID          uint   ` + "`" + `json:"id"` + "`" + `
		Name        string ` + "`" + `json:"name"` + "`" + `
		DisplayName string ` + "`" + `json:"display_name"` + "`" + `
		Description string ` + "`" + `json:"description"` + "`" + `
		Category    string ` + "`" + `json:"category"` + "`" + `
		Tags        string ` + "`" + `json:"tags"` + "`" + `
		PreviewCode string ` + "`" + `json:"preview_code"` + "`" + `
	}

	query := h.DB.Model(&models.UIComponent{}).Where("is_public = ?", true)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var components []componentSummary
	if err := query.Order("category, name").Find(&components).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch components"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": components})
}

// GetComponentDetail serves GET /api/ui-components/:name.
func (h *UIRegistryHandler) GetComponentDetail(c *gin.Context) {
	name := c.Param("name")

	var comp models.UIComponent
	if err := h.DB.Where("name = ? AND is_public = ?", name, true).First(&comp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "component not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": comp})
}

// ── Admin endpoints ───────────────────────────────────────────────────────────

type uiComponentInput struct {
	Name         string ` + "`" + `json:"name" binding:"required"` + "`" + `
	DisplayName  string ` + "`" + `json:"display_name" binding:"required"` + "`" + `
	Description  string ` + "`" + `json:"description"` + "`" + `
	Category     string ` + "`" + `json:"category" binding:"required"` + "`" + `
	Tags         string ` + "`" + `json:"tags"` + "`" + `
	Files        string ` + "`" + `json:"files" binding:"required"` + "`" + `
	Dependencies string ` + "`" + `json:"dependencies"` + "`" + `
	RegistryDeps string ` + "`" + `json:"registry_deps"` + "`" + `
	PreviewCode  string ` + "`" + `json:"preview_code"` + "`" + `
	IsPublic     *bool  ` + "`" + `json:"is_public"` + "`" + `
}

// CreateComponent serves POST /api/admin/ui-components.
func (h *UIRegistryHandler) CreateComponent(c *gin.Context) {
	var req uiComponentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	comp := models.UIComponent{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Description:  req.Description,
		Category:     req.Category,
		Tags:         req.Tags,
		Files:        req.Files,
		Dependencies: req.Dependencies,
		RegistryDeps: req.RegistryDeps,
		PreviewCode:  req.PreviewCode,
		IsPublic:     isPublic,
	}

	if err := h.DB.Create(&comp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create component"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": comp, "message": "Component created"})
}

// UpdateComponent serves PUT /api/admin/ui-components/:name.
func (h *UIRegistryHandler) UpdateComponent(c *gin.Context) {
	name := c.Param("name")

	var comp models.UIComponent
	if err := h.DB.Where("name = ?", name).First(&comp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "component not found"})
		return
	}

	var req uiComponentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"display_name":  req.DisplayName,
		"description":   req.Description,
		"category":      req.Category,
		"tags":          req.Tags,
		"files":         req.Files,
		"dependencies":  req.Dependencies,
		"registry_deps": req.RegistryDeps,
		"preview_code":  req.PreviewCode,
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if err := h.DB.Model(&comp).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update component"})
		return
	}

	h.DB.Where("name = ?", name).First(&comp)
	c.JSON(http.StatusOK, gin.H{"data": comp, "message": "Component updated"})
}

// DeleteComponent serves DELETE /api/admin/ui-components/:name.
func (h *UIRegistryHandler) DeleteComponent(c *gin.Context) {
	name := c.Param("name")

	var comp models.UIComponent
	if err := h.DB.Where("name = ?", name).First(&comp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "component not found"})
		return
	}

	if err := h.DB.Delete(&comp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete component"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Component deleted"})
}
`
}
