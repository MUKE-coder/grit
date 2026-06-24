package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// injectAll injects code into all existing files that have markers.
func (g *Generator) injectAll(names Names) error {
	apiRoot := g.APIRoot()
	sharedRoot := filepath.Join(g.Root, "packages", "shared")
	adminRoot := g.AdminRoot()

	// 1. Inject model into AutoMigrate
	modelFile := filepath.Join(apiRoot, "internal", "models", "user.go")
	if fileExists(modelFile) {
		if err := injectBefore(modelFile, "// grit:models",
			fmt.Sprintf("\t\t&%s{},", names.Pascal)); err != nil {
			return fmt.Errorf("injecting model: %w", err)
		}
		fmt.Println("  ✓ Injected model into AutoMigrate")
	}

	// 2. Inject model into GORM Studio (inline, before the marker)
	routesFile := filepath.Join(apiRoot, "internal", "routes", "routes.go")
	if fileExists(routesFile) {
		if err := injectInline(routesFile, "/* grit:studio */",
			fmt.Sprintf("&models.%s{}, ", names.Pascal)); err != nil {
			return fmt.Errorf("injecting studio model: %w", err)
		}
		fmt.Println("  ✓ Injected model into GORM Studio")
	}

	// 2b. Register the model with the sync registry so desktop clients
	// can push/pull it via /api/sync. Tolerant of older projects that
	// don't have the // grit:sync marker yet.
	if fileExists(routesFile) {
		register := fmt.Sprintf("\tsyncRegistry.Register(\"%s\", &models.%s{})", names.PluralSnake, names.Pascal)
		if err := injectBefore(routesFile, "// grit:sync", register); err == nil {
			fmt.Println("  ✓ Registered model with sync registry")
		}
	}

	// 3. Inject handler init
	if fileExists(routesFile) {
		// v3.31.33 -- if the resource has file fields, wire Storage so
		// the Create/Update flows can do immediate S3 cleanup on
		// replace and mark uploads as claimed.
		hasFileFields := false
		for _, f := range g.Definition.Fields {
			if f.IsFileField() {
				hasFileFields = true
				break
			}
		}
		extraField := ""
		if hasFileFields {
			extraField = "\n\t\tStorage: svc.Storage,"
		}
		handlerInit := fmt.Sprintf(`	%sHandler := &handlers.%sHandler{
		DB: db,%s
	}`, names.Camel, names.Pascal, extraField)
		if err := injectBefore(routesFile, "// grit:handlers", handlerInit); err != nil {
			return fmt.Errorf("injecting handler: %w", err)
		}
		fmt.Println("  ✓ Injected handler initialization")
	}

	// 4. Inject routes (role-restricted or default split)
	if fileExists(routesFile) {
		if len(g.Roles) > 0 {
			// Role-restricted: inject all routes into // grit:routes:custom as a group
			roleArgs := make([]string, len(g.Roles))
			for i, r := range g.Roles {
				roleArgs[i] = fmt.Sprintf("%q", r)
			}
			rolesStr := strings.Join(roleArgs, ", ")
			customRoutes := fmt.Sprintf(`	// %s routes (restricted to %s)
	%sGroup := protected.Group("/%s")
	%sGroup.Use(middleware.RequireRole(%s))
	{
		%sGroup.GET("", %sHandler.List)
		%sGroup.GET("/export", %sHandler.Export)
		%sGroup.GET("/:id", %sHandler.GetByID)
		%sGroup.POST("", %sHandler.Create)
		%sGroup.PUT("/:id", %sHandler.Update)
		%sGroup.PATCH("/:id", %sHandler.Patch)
		%sGroup.DELETE("/:id", %sHandler.Delete)
	}`,
				names.PluralPascal, strings.Join(g.Roles, ", "),
				names.Camel, names.Plural,
				names.Camel, rolesStr,
				names.Camel, names.Camel,
				names.Camel, names.Camel,
				names.Camel, names.Camel,
				names.Camel, names.Camel,
				names.Camel, names.Camel,
				names.Camel, names.Camel,
				names.Camel, names.Camel)
			if err := injectBefore(routesFile, "// grit:routes:custom", customRoutes); err != nil {
				return fmt.Errorf("injecting role-restricted routes: %w", err)
			}
			fmt.Printf("  ✓ Injected role-restricted routes (%s)\n", strings.Join(g.Roles, ", "))
		} else {
			// Default: CRUD in protected, DELETE in admin
			protectedRoutes := fmt.Sprintf(`		protected.GET("/%s", %sHandler.List)
		protected.GET("/%s/export", %sHandler.Export)
		protected.GET("/%s/:id", %sHandler.GetByID)
		protected.POST("/%s", %sHandler.Create)
		protected.PUT("/%s/:id", %sHandler.Update)
		protected.PATCH("/%s/:id", %sHandler.Patch)`,
				names.Plural, names.Camel,
				names.Plural, names.Camel,
				names.Plural, names.Camel,
				names.Plural, names.Camel,
				names.Plural, names.Camel,
				names.Plural, names.Camel)
			if err := injectBefore(routesFile, "// grit:routes:protected", protectedRoutes); err != nil {
				return fmt.Errorf("injecting protected routes: %w", err)
			}
			fmt.Println("  ✓ Injected protected routes")

			adminRoutes := fmt.Sprintf(`		admin.DELETE("/%s/:id", %sHandler.Delete)`,
				names.Plural, names.Camel)
			if err := injectBefore(routesFile, "// grit:routes:admin", adminRoutes); err != nil {
				return fmt.Errorf("injecting admin routes: %w", err)
			}
			fmt.Println("  ✓ Injected admin routes")
		}
	}

	// 5b. v3.31.20 — inject a dispatch case into the form-share submit
	// service so public submissions can create records of this resource.
	dispatchFile := filepath.Join(apiRoot, "internal", "services", "form_share_dispatch.go")
	if fileExists(dispatchFile) {
		labelExpr := pickLabelExpr(g.Definition.Fields)
		dispatchCase := fmt.Sprintf(`	case %q:
		item := &models.%s{}
		body, _ := json.Marshal(fields)
		if err := json.Unmarshal(body, item); err != nil {
			return nil, fmt.Errorf("decoding %s body: %%w", err)
		}
		if err := db.Create(item).Error; err != nil {
			return nil, fmt.Errorf("creating %s: %%w", err)
		}
		return &SharedResourceSubmission{ID: item.ID, Label: %s}, nil
`,
			names.Pascal,
			names.Pascal,
			names.Pascal,
			names.Pascal,
			labelExpr,
		)
		if err := injectBefore(dispatchFile, "// grit:form-share:dispatch", dispatchCase); err != nil {
			return fmt.Errorf("injecting form-share dispatch: %w", err)
		}
		// v3.31.43: inject a matching case into PublicFields so the
		// public form renders the right inputs. The reflection helper
		// PublicFields(...) takes a model pointer and walks its struct
		// tags -- no per-resource field list to maintain here.
		fieldsCase := fmt.Sprintf(`	case %q:
		return reflectPublicFields(&models.%s{})
`,
			names.Pascal,
			names.Pascal,
		)
		if err := injectBefore(dispatchFile, "// grit:form-share:fields", fieldsCase); err != nil {
			// Pre-v3.31.43 projects don't have the fields marker yet.
			// Surface a warning so the operator knows to add it but
			// don't fail the whole generate -- the dispatch case
			// landed fine.
			fmt.Println("  ⚠ form-share:fields marker missing; public form will fall back to no-fields. Add `// grit:form-share:fields` to services/form_share_dispatch.go inside PublicFields().")
		} else {
			fmt.Println("  ✓ Injected form-share fields case")
		}
		// Make sure the imports the case needs are present.
		if err := ensureDispatchImports(dispatchFile, g.Module); err != nil {
			return fmt.Errorf("updating dispatch imports: %w", err)
		}
		fmt.Println("  ✓ Injected form-share dispatch case")
	}

	// 6. Inject schema export
	schemaIndex := filepath.Join(sharedRoot, "schemas", "index.ts")
	if fileExists(schemaIndex) {
		schemaExport := fmt.Sprintf(`export {
  Create%sSchema,
  Update%sSchema,
  type Create%sInput,
  type Update%sInput,
} from "./%s";`, names.Pascal, names.Pascal, names.Pascal, names.Pascal, names.Kebab)
		if err := injectBefore(schemaIndex, "// grit:schemas", schemaExport); err != nil {
			return fmt.Errorf("injecting schema export: %w", err)
		}
		fmt.Println("  ✓ Injected schema export")
	}

	// 7. Inject type export
	typesIndex := filepath.Join(sharedRoot, "types", "index.ts")
	if fileExists(typesIndex) {
		typeExport := fmt.Sprintf(`export type { %s } from "./%s";`, names.Pascal, names.Kebab)
		if err := injectBefore(typesIndex, "// grit:types", typeExport); err != nil {
			return fmt.Errorf("injecting type export: %w", err)
		}
		fmt.Println("  ✓ Injected type export")
	}

	// 8. Inject API routes constant
	constantsIndex := filepath.Join(sharedRoot, "constants", "index.ts")
	if fileExists(constantsIndex) {
		upper := strings.ToUpper(names.Plural)
		routeConst := fmt.Sprintf(`  %s: {
    LIST: "/api/%s",
    GET: (id: number) => `+"`"+`/api/%s/${id}`+"`"+`,
    CREATE: "/api/%s",
    UPDATE: (id: number) => `+"`"+`/api/%s/${id}`+"`"+`,
    DELETE: (id: number) => `+"`"+`/api/%s/${id}`+"`"+`,
  },`, upper, names.Plural, names.Plural, names.Plural, names.Plural, names.Plural)
		if err := injectBefore(constantsIndex, "// grit:api-routes", routeConst); err != nil {
			return fmt.Errorf("injecting API routes: %w", err)
		}
		fmt.Println("  ✓ Injected API route constants")
	}

	// 9. Inject resource import into resource registry
	registryFile := filepath.Join(adminRoot, "resources", "index.ts")
	// TanStack admin: src/resources/index.ts
	if !fileExists(registryFile) {
		registryFile = filepath.Join(adminRoot, "src", "resources", "index.ts")
	}
	if fileExists(registryFile) {
		resourceImport := fmt.Sprintf(`import { %sResource } from "./%s";`,
			names.Camel, names.PluralKebab)
		if err := injectBefore(registryFile, "// grit:resources", resourceImport); err != nil {
			return fmt.Errorf("injecting resource import: %w", err)
		}
		fmt.Println("  ✓ Injected resource import into registry")
	}

	// 10. Inject resource into registry array
	if fileExists(registryFile) {
		resourceEntry := fmt.Sprintf(`  %sResource,`, names.Camel)
		if err := injectBefore(registryFile, "// grit:resource-list", resourceEntry); err != nil {
			return fmt.Errorf("injecting resource entry: %w", err)
		}
		fmt.Println("  ✓ Injected resource into registry list")
	}

	return nil
}

// injectInline inserts code directly before a marker on the same line.
func injectInline(filePath, marker, code string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	content := string(data)
	idx := strings.Index(content, marker)
	if idx == -1 {
		return fmt.Errorf("marker %q not found in %s", marker, filePath)
	}

	newContent := content[:idx] + code + content[idx:]
	return os.WriteFile(filePath, []byte(newContent), 0644)
}

// injectBefore finds a marker in a file and inserts code on the line before it.
func injectBefore(filePath, marker, code string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", filePath, err)
	}

	content := string(data)
	idx := strings.Index(content, marker)
	if idx == -1 {
		return fmt.Errorf("marker %q not found in %s", marker, filePath)
	}

	// Find the start of the line containing the marker
	lineStart := idx
	for lineStart > 0 && content[lineStart-1] != '\n' {
		lineStart--
	}

	// Get the indentation of the marker line (for reference, not used for injected code)
	_ = content[lineStart:idx]

	// Insert the code before the marker line
	newContent := content[:lineStart] + code + "\n" + content[lineStart:]

	return os.WriteFile(filePath, []byte(newContent), 0644)
}

// guessLucideIcon returns a Lucide icon name based on the resource name.
func guessLucideIcon(name string) string {
	lower := strings.ToLower(name)
	icons := map[string]string{
		"post":         "FileText",
		"article":      "Newspaper",
		"blog":         "Newspaper",
		"comment":      "MessageSquare",
		"category":     "FolderTree",
		"tag":          "Tag",
		"product":      "Package",
		"order":        "ShoppingCart",
		"invoice":      "Receipt",
		"payment":      "CreditCard",
		"customer":     "UserCircle",
		"user":         "Users",
		"project":      "Briefcase",
		"task":         "CheckSquare",
		"event":        "Calendar",
		"file":         "File",
		"image":        "Image",
		"media":        "Image",
		"message":      "Mail",
		"notification": "Bell",
		"setting":      "Settings",
		"role":         "Shield",
		"permission":   "Lock",
		"team":         "UsersRound",
		"company":      "Building2",
		"organization": "Building2",
		"report":       "BarChart3",
		"analytic":     "TrendingUp",
		"log":          "ScrollText",
		"page":         "FileText",
		"document":     "FileText",
		"review":       "Star",
		"subscription": "CreditCard",
		"plan":         "Gem",
		"coupon":       "Ticket",
		"discount":     "Percent",
		"shipping":     "Truck",
		"address":      "MapPin",
		"location":     "MapPin",
		"contact":      "Contact",
		"lead":         "Target",
		"deal":         "Handshake",
		"pipeline":     "GitBranch",
		"workflow":     "Workflow",
		"template":     "LayoutTemplate",
		"email":        "Mail",
		"campaign":     "Megaphone",
		"survey":       "ClipboardList",
		"form":         "FormInput",
		"question":     "HelpCircle",
		"answer":       "MessageCircle",
		"ticket":       "Ticket",
		"issue":        "AlertCircle",
		"bug":          "Bug",
		"feature":      "Sparkles",
		"release":      "Rocket",
		"version":      "GitCommit",
		"deploy":       "CloudUpload",
	}

	for key, icon := range icons {
		if strings.Contains(lower, key) {
			return icon
		}
	}
	return "Database"
}
