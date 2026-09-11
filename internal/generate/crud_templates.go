package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// The generated handler and service, built from one description of the
// resource.
//
// The handler used to do all of it: 27 direct GORM calls across list, export,
// read, write, patch, delete and bulk, while the service beside it was called
// by nothing and had drifted from what the handler actually did. The service
// owns every query now. The handler reads the request, calls the service and
// writes the answer; a job, a command or a test calls the same methods with
// no request at all, so a rule enforced there is enforced everywhere.
//
// Both files are rendered from crudParts, worked out once, so they cannot
// disagree about a resource again.

// crudParts is everything about a resource the handler and the service need.
type crudParts struct {
	// Handler: the request.
	createFields    string
	updateFields    string
	createAssign    string
	updateMap       string
	itemsBuild      string // Create: req.Items onto item.Items
	itemsFromUpdate string // Update: req.Items as *[]models.Child
	linksFromCreate string
	linksFromUpdate string
	linksArg        string
	itemsArg        string
	fkFilters       string
	exportCols      string
	pdfFields       string
	pdfSections     string
	optionalID      string
	identExpr       string
	handlerImports  string
	datatypesImport string

	// Service: the queries.
	preloads     string
	searchCols   string
	sortCols     string
	filterCols   string
	writable     string // body of the writable-columns map, one tab in
	linksType    string
	linksParam   string
	linksApply   string // inside a transaction on tx
	itemsParam   string
	itemsReplace string // inside a transaction on tx
	patchM2M     string
	single       bool // one statement: RETURNING, no transaction
	hasFiles     bool
	owned        bool
	ownerCol     string
	ownerField   string
	ownerName    string
}

// crud works out a resource's parts.
func (g *Generator) crud(names Names) crudParts {
	var p crudParts
	patchAllowed := ""
	var preloads []string

	ownerName := ""
	if owner := g.Definition.OwnerField(); owner != nil {
		ownerName = owner.Name
		p.owned = true
		p.ownerName = owner.Name
		p.ownerCol = toSnakeCase(owner.Name) + "_id"
		p.ownerField = toPascalCase(owner.Name) + "ID"
	}

	var linkFields, linkCreate, linkUpdate string
	hasLinks := false

	for _, f := range g.Definition.Fields {
		if f.IsSlug() {
			continue
		}

		if f.IsBelongsTo() {
			baseName := strings.TrimSuffix(f.Name, "_id")
			fkGoName := toPascalCase(baseName) + "ID"
			fkJson := toSnakeCase(baseName) + "_id"
			preloads = append(preloads, toPascalCase(baseName))

			// The owner is stamped from the signed-in caller, in the service.
			// Accepting it from the body would let a caller create rows that
			// belong to somebody else.
			if ownerName != "" && f.Name == ownerName {
				continue
			}

			// A self-reference is optional: the root of a tree has no parent.
			selfRef := f.RelatedModelName() == toPascalCase(g.Definition.Name)
			if selfRef {
				p.createFields += fmt.Sprintf("\t\t%s string `json:\"%s\"`\n", fkGoName, fkJson)
				p.createAssign += fmt.Sprintf("\t\t%s: optional%sID(req.%s),\n", fkGoName, names.Pascal, fkGoName)
			} else {
				p.createFields += fmt.Sprintf("\t\t%s string `json:\"%s\" binding:\"required\"`\n", fkGoName, fkJson)
				p.createAssign += fmt.Sprintf("\t\t%s: req.%s,\n", fkGoName, fkGoName)
			}
			p.updateFields += fmt.Sprintf("\t\t%s *string `json:\"%s\"`\n", fkGoName, fkJson)
			if selfRef {
				p.updateMap += fmt.Sprintf("\tif req.%s != nil {\n\t\tupdates[\"%s\"] = optional%sID(*req.%s)\n\t}\n", fkGoName, fkJson, names.Pascal, fkGoName)
			} else {
				p.updateMap += fmt.Sprintf("\tif req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}\n", fkGoName, fkJson, fkGoName)
			}
			patchAllowed += fmt.Sprintf("\t\t\"%s\": true,\n", fkJson)
			continue
		}

		if f.IsManyToMany() {
			hasLinks = true
			relModel := f.RelatedModelName()
			assocName := toPascalCase(f.Name)
			idsName := toPascalCase(f.Name) + "IDs"
			idsJson := strings.TrimSuffix(toSnakeCase(f.Name), "s") + "_ids"
			preloads = append(preloads, assocName)

			p.createFields += fmt.Sprintf("\t\t%s []string `json:\"%s\"`\n", idsName, idsJson)
			p.updateFields += fmt.Sprintf("\t\t%s *[]string `json:\"%s\"`\n", idsName, idsJson)

			linkFields += fmt.Sprintf("\t%s *[]string\n", idsName)
			linkCreate += fmt.Sprintf("\tif len(req.%s) > 0 {\n\t\tlinks.%s = &req.%s\n\t}\n", idsName, idsName, idsName)
			linkUpdate += fmt.Sprintf("\t\t%s: req.%s,\n", idsName, idsName)

			related := "related" + assocName
			p.linksApply += fmt.Sprintf(`		if links.%[1]s != nil {
			var %[2]s []models.%[3]s
			if len(*links.%[1]s) > 0 {
				if err := tx.Find(&%[2]s, *links.%[1]s).Error; err != nil {
					return err
				}
				// An id that matches nothing is refused, not dropped: dropping it
				// emptied the set when the only id given was mistyped.
				wanted := map[string]bool{}
				for _, id := range *links.%[1]s {
					wanted[id] = true
				}
				if len(%[2]s) != len(wanted) {
					return respond.Rule("%[6]s: %%d of the %%d ids given do not exist", len(wanted)-len(%[2]s), len(wanted))
				}
			}
			if err := tx.Model(item).Association(%[4]q).Replace(%[2]s); err != nil {
				return fmt.Errorf("setting %[5]s: %%w", err)
			}
		}
`, idsName, related, relModel, assocName, toSnakeCase(f.Name), idsJson)
			continue
		}

		goName := toPascalCase(f.Name)
		goType := f.GoType()
		jsonTag := toSnakeCase(f.Name)
		bindingTag := ""
		if f.Required {
			bindingTag = ` binding:"required"`
		}
		p.createFields += fmt.Sprintf("\t\t%s %s `json:\"%s\"%s`\n", goName, goType, jsonTag, bindingTag)
		p.createAssign += fmt.Sprintf("\t\t%s: req.%s,\n", goName, goName)
		patchAllowed += fmt.Sprintf("\t\t\"%s\": true,\n", jsonTag)

		// Update uses pointers, so "sent" and "not sent" can be told apart.
		switch goType {
		case "bool":
			p.updateFields += fmt.Sprintf("\t\t%s *%s `json:\"%s\"`\n", goName, goType, jsonTag)
			p.updateMap += fmt.Sprintf("\tif req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}\n", goName, jsonTag, goName)
		case "string":
			p.updateFields += fmt.Sprintf("\t\t%s %s `json:\"%s\"`\n", goName, goType, jsonTag)
			p.updateMap += fmt.Sprintf("\tif req.%s != \"\" {\n\t\tupdates[\"%s\"] = req.%s\n\t}\n", goName, jsonTag, goName)
		case "*time.Time":
			p.updateFields += fmt.Sprintf("\t\t%s %s `json:\"%s\"`\n", goName, goType, jsonTag)
			p.updateMap += fmt.Sprintf("\tif req.%s != nil {\n\t\tupdates[\"%s\"] = req.%s\n\t}\n", goName, jsonTag, goName)
		default:
			p.updateFields += fmt.Sprintf("\t\t%s *%s `json:\"%s\"`\n", goName, goType, jsonTag)
			p.updateMap += fmt.Sprintf("\tif req.%s != nil {\n\t\tupdates[\"%s\"] = *req.%s\n\t}\n", goName, jsonTag, goName)
		}
	}

	if hasLinks {
		p.linksType = fmt.Sprintf(`// %sLinks carries the many-to-many ids a write sets. A nil list leaves that
// set as it is; a list, empty or not, replaces it.
type %sLinks struct {
%s}
`, names.Pascal, names.Pascal, linkFields)
		p.linksParam = fmt.Sprintf(", links %sLinks", names.Pascal)
		p.linksArg = ", links"
		p.linksFromCreate = fmt.Sprintf("\n\tlinks := services.%sLinks{}\n%s", names.Pascal, linkCreate)
		p.linksFromUpdate = fmt.Sprintf("\tlinks := services.%sLinks{\n%s\t}\n", names.Pascal, linkUpdate)
	}

	// Inline line items (--items): the request carries the lines, a create
	// cascades them, and an update replaces them in the same transaction as
	// the row. The parent key is set here, never sent per line.
	if g.Definition.Items != nil {
		child := BuildNames(g.Definition.Items).Pascal
		reqFields, assign := "", ""
		parent := toPascalCase(g.Definition.Name)
		for _, cf := range g.Definition.Items.Fields {
			if cf.IsSlug() || cf.IsManyToMany() {
				continue
			}
			// The key back to the parent is set on save. Every other belongs_to
			// is what the line is about (the product, the account), so it is
			// required on each line.
			if cf.IsBelongsTo() {
				if cf.RelatedModelName() == parent {
					continue
				}
				base := strings.TrimSuffix(cf.Name, "_id")
				fkName := toPascalCase(base) + "ID"
				fkJSON := toSnakeCase(base) + "_id"
				reqFields += fmt.Sprintf("\t\t\t%s string `json:\"%s\" binding:\"required\"`\n", fkName, fkJSON)
				assign += fmt.Sprintf("\t\t\t\t%s: it.%s,\n", fkName, fkName)
				continue
			}
			gName := toPascalCase(cf.Name)
			reqFields += fmt.Sprintf("\t\t\t%s %s `json:\"%s\"`\n", gName, cf.GoType(), toSnakeCase(cf.Name))
			assign += fmt.Sprintf("\t\t\t\t%s: it.%s,\n", gName, gName)
		}
		itemsReq := fmt.Sprintf("\t\tItems []struct {\n%s\t\t} `json:\"items\"`\n", reqFields)
		p.createFields += itemsReq
		p.updateFields += itemsReq
		p.itemsBuild = fmt.Sprintf("\n\tif len(req.Items) > 0 {\n\t\titems := make([]models.%s, 0, len(req.Items))\n\t\tfor _, it := range req.Items {\n\t\t\titems = append(items, models.%s{\n%s\t\t\t})\n\t\t}\n\t\titem.Items = items\n\t}\n", child, child, assign)
		p.itemsFromUpdate = fmt.Sprintf(`	// Nil leaves the lines as they are; a list, empty or not, replaces them.
	var items *[]models.%s
	if req.Items != nil {
		rows := make([]models.%s, 0, len(req.Items))
		for _, it := range req.Items {
			rows = append(rows, models.%s{
%s			})
		}
		items = &rows
	}
`, child, child, child, assign)
		p.itemsParam = fmt.Sprintf(", items *[]models.%s", child)
		p.itemsArg = ", items"
		p.itemsReplace = fmt.Sprintf(`		if items != nil {
			if err := tx.Where("%s = ?", item.ID).Delete(&models.%s{}).Error; err != nil {
				return err
			}
			if len(*items) > 0 {
				rows := *items
				for i := range rows {
					rows[i].%sID = item.ID
				}
				if err := tx.Create(&rows).Error; err != nil {
					return err
				}
			}
		}
`, names.Snake+"_id", child, names.Pascal)
		preloads = append(preloads, "Items")
	}

	for _, pl := range preloads {
		p.preloads += fmt.Sprintf(".Preload(%q)", pl)
	}
	p.writable = dedentOneTab(patchAllowed)

	// One statement needs neither the reload nor the transaction GORM wraps
	// writes in: RETURNING brings back what the database filled in. Not when
	// the model has relations to preload, children, join rows or a sequence
	// hook writing alongside.
	hasAuto := false
	for _, f := range g.Definition.Fields {
		if f.Auto {
			hasAuto = true
			break
		}
	}
	p.single = p.preloads == "" && !hasLinks && g.Definition.Items == nil && !hasAuto

	// Sortable: never a relation, and a money field sorts by its amount
	// column, the one that exists.
	p.sortCols = `"id": true, "created_at": true`
	for _, f := range g.Definition.Fields {
		if f.IsRelationship() {
			continue
		}
		if FieldType(f.Type) == FieldMoney {
			p.sortCols += fmt.Sprintf(`, "%s_amount": true`, toSnakeCase(f.Name))
			continue
		}
		if f.GoType() == "string" || f.GoType() == "int" || f.GoType() == "uint" {
			p.sortCols += fmt.Sprintf(`, "%s": true`, toSnakeCase(f.Name))
		}
	}

	// Filterable: wider than sortable (a bool or a foreign key is worth
	// filtering by), but ciphertext cannot be matched.
	p.filterCols = `"id": true`
	for _, f := range g.Definition.Fields {
		if f.IsManyToMany() || f.Encrypted {
			continue
		}
		col := toSnakeCase(f.Name)
		if f.IsBelongsTo() {
			col = f.FKColumnName()
		}
		if FieldType(f.Type) == FieldMoney {
			p.filterCols += fmt.Sprintf(`, "%s_amount": true, "%s_currency": true`, col, col)
			continue
		}
		p.filterCols += fmt.Sprintf(`, "%s": true`, col)
	}
	p.searchCols = g.buildHandlerSearchCols()

	// PATCH replaces a many-to-many set when its key is in the body, the same
	// as POST and PUT. touchedRelations is declared even with no relations,
	// because the guard after it reads it either way.
	p.patchM2M = "\ttouchedRelations := false\n"
	for _, f := range g.Definition.Fields {
		if !f.IsManyToMany() {
			continue
		}
		rel := MakeNames(f.RelatedModelName())
		key := strings.TrimSuffix(toSnakeCase(f.Name), "s") + "_ids"
		p.patchM2M += fmt.Sprintf(`	// %[1]s: replace the whole set. An empty list is a real instruction, so
	// presence in the body is what counts, not whether it has anything in it.
	if rawIDs, ok := body[%[1]q]; ok {
		list, isList := rawIDs.([]interface{})
		if !isList {
			return nil, nil, respond.Rule(%[2]q)
		}
		ids := make([]string, 0, len(list))
		wanted := map[string]bool{}
		for _, raw := range list {
			id, ok := raw.(string)
			if !ok {
				return nil, nil, respond.Rule(%[2]q)
			}
			ids = append(ids, id)
			wanted[id] = true
		}
		var related []models.%[3]s
		if len(ids) > 0 {
			if err := db.Find(&related, ids).Error; err != nil {
				return nil, nil, err
			}
			// An id that matches nothing is refused, not dropped: dropping it
			// emptied the set when the only id given was mistyped.
			if len(related) != len(wanted) {
				return nil, nil, respond.Rule("%[1]s: %%d of the %%d ids given do not exist", len(wanted)-len(related), len(wanted))
			}
		}
		if err := db.Model(item).Association(%[4]q).Replace(related); err != nil {
			return nil, nil, fmt.Errorf("updating %[5]s: %%w", err)
		}
		touchedRelations = true
	}
`, key, key+" must be a list of ids", rel.Pascal, toPascalCase(f.Name), toSnakeCase(f.Name))
	}

	// FK filters: GET /<plural>?category_id=... returns one parent's children.
	for _, f := range g.Definition.Fields {
		if f.IsBelongsTo() {
			fk := f.FKColumnName()
			p.fkFilters += fmt.Sprintf(".With(%q, c.Query(%q))", fk, fk)
		}
	}

	p.exportCols = "\t\t\t{Header: \"ID\", Field: \"ID\"},\n"
	for _, f := range g.Definition.Fields {
		if f.IsRelationship() {
			continue
		}
		header := strings.ReplaceAll(toPascalCase(f.Name), "ID", " ID")
		format := ""
		switch {
		case f.GoType() == "time.Time" || f.GoType() == "*time.Time":
			format = "date:2006-01-02"
		case f.GoType() == "bool":
			format = "bool"
		}
		if format != "" {
			p.exportCols += fmt.Sprintf("\t\t\t{Header: %q, Field: %q, Format: %q},\n", header, toPascalCase(f.Name), format)
		} else {
			p.exportCols += fmt.Sprintf("\t\t\t{Header: %q, Field: %q},\n", header, toPascalCase(f.Name))
		}
	}
	p.exportCols += "\t\t\t{Header: \"Created At\", Field: \"CreatedAt\", Format: \"date:2006-01-02\"},"

	// Imports the handler's request structs need, parent and child alike: the
	// inline Items struct lives in the handler file too.
	needsDatatypes, needsJSONTime, needsMoney, needsCrypto := false, false, false, false
	fields := g.Definition.Fields
	if g.Definition.Items != nil {
		fields = append(append([]Field{}, fields...), g.Definition.Items.Fields...)
	}
	for _, f := range fields {
		if f.NeedsJSONTimeImport() {
			needsJSONTime = true
		}
		if FieldType(f.Type) == FieldMoney {
			needsMoney = true
		}
		if f.Encrypted {
			needsCrypto = true
		}
		if f.NeedsDatatypesImport() {
			needsDatatypes = true
		}
	}
	for _, f := range g.Definition.Fields {
		if f.IsFileField() {
			p.hasFiles = true
		}
	}
	if needsDatatypes {
		p.datatypesImport = "\n\t\"gorm.io/datatypes\""
	}
	if needsJSONTime {
		p.handlerImports += "\n\t\"" + g.Module + "/internal/jsontime\""
	}
	if needsMoney {
		p.handlerImports += "\n\t\"" + g.Module + "/internal/money\""
	}
	if needsCrypto {
		p.handlerImports += "\n\t\"" + g.Module + "/internal/crypto\""
	}
	// The request structs name files.FileRef, which the service does not see.
	if strings.Contains(p.createFields+p.updateFields, "files.") {
		p.handlerImports += "\n\t\"" + g.Module + "/internal/files\""
	}
	if p.hasFiles {
		p.handlerImports += "\n\t\"" + g.Module + "/internal/storage\""
	}
	if g.Definition.WorkflowField() != nil {
		p.handlerImports += "\n\t\"" + g.Module + "/internal/workflow\""
	}

	p.identExpr = pickIdentifierExpr(g.Definition.Fields)

	// The PDF: the resource's own fields as a detail grid, line items as a
	// table. Files, media and many-to-many have no sensible text rendering.
	for _, f := range g.Definition.Fields {
		if f.IsManyToMany() || f.IsSlug() {
			continue
		}
		switch f.FormFieldType() {
		case "image", "images", "video", "videos", "file", "files", "richtext":
			continue
		}
		label := strings.Join(splitPascal(toPascalCase(f.Name)), " ")
		if f.IsBelongsTo() {
			assoc := toPascalCase(strings.TrimSuffix(f.Name, "_id"))
			p.pdfFields += fmt.Sprintf("\t\t\t{Label: %q, Value: pdf.Display(item.%s)},\n", label, assoc)
			continue
		}
		p.pdfFields += fmt.Sprintf("\t\t\t{Label: %q, Value: pdf.Value(item.%s)},\n", label, toPascalCase(f.Name))
	}
	p.pdfFields += fmt.Sprintf("\t\t\t{Label: %q, Value: pdf.Value(item.CreatedAt)},\n", "Created")
	if g.Definition.Items != nil {
		childNames := BuildNames(g.Definition.Items)
		headers, cells, aligns := "", "", ""
		for _, cf := range g.Definition.Items.Fields {
			if cf.IsBelongsTo() || cf.IsManyToMany() {
				continue
			}
			switch cf.FormFieldType() {
			case "image", "images", "video", "videos", "file", "files", "richtext":
				continue
			}
			headers += fmt.Sprintf("%q, ", strings.Join(splitPascal(toPascalCase(cf.Name)), " "))
			cells += fmt.Sprintf("pdf.Value(row.%s), ", toPascalCase(cf.Name))
			if cf.FormFieldType() == "number" {
				aligns += `"R", `
			} else {
				aligns += `"L", `
			}
		}
		p.pdfSections = fmt.Sprintf(`
	itemRows := make([][]string, 0, len(item.Items))
	for _, row := range item.Items {
		itemRows = append(itemRows, []string{%s})
	}
	if len(itemRows) > 0 {
		rec.Sections = append(rec.Sections, pdf.Section{
			Title:   %q,
			Headers: []string{%s},
			Aligns:  []string{%s},
			Rows:    itemRows,
		})
	}
`, strings.TrimSuffix(cells, ", "), strings.Join(splitPascal(childNames.PluralPascal), " "),
			strings.TrimSuffix(headers, ", "), strings.TrimSuffix(aligns, ", "))
	}

	// optionalID only for a resource with a self-reference: the only nullable
	// foreign key the generator writes, and an unused function does not compile.
	for _, f := range g.Definition.Fields {
		if f.IsBelongsTo() && f.RelatedModelName() == toPascalCase(g.Definition.Name) {
			p.optionalID = strings.ReplaceAll(OPTIONAL_ID_HELPER_SRC, "{{Pascal}}", names.Pascal)
			break
		}
	}
	return p
}

// serviceWriteHelper is the single-statement write, as two methods on the
// service. They do what database.Write and database.SupportsReturning do, and
// they are here rather than there because the database package imports
// services for its seeders: a service that imported database back would be an
// import cycle, and the project would not build. Methods, so that two
// resources in one package cannot collide.
const serviceWriteHelper = `
// write is a session for a single-statement write: no wrapping transaction,
// and RETURNING where the dialect has it. Safe only because every write it is
// used for is exactly one statement, which the generator knew.
func (s *{{Pascal}}Service) write(db *gorm.DB) *gorm.DB {
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if s.returning(db) {
		tx = tx.Clauses(clause.Returning{})
	}
	return tx
}

// returning reports whether the dialect hands back the written row. MySQL
// does not, and does not say so: the clause is dropped and the defaults come
// back empty, so there the row is read again.
func (s *{{Pascal}}Service) returning(db *gorm.DB) bool {
	switch db.Dialector.Name() {
	case "postgres", "sqlite":
		return true
	}
	return false
}
`

// writeGoService writes services/<resource>.go: every query the resource has.
func (g *Generator) writeGoService(names Names) error {
	path := filepath.Join(g.APIRoot(), "internal", "services", names.Snake+".go")
	return writeFileWithDirs(path, g.serviceSource(names))
}

// serviceSource renders the service.
func (g *Generator) serviceSource(names Names) string {
	p := g.crud(names)

	imports := "\t\"" + g.Module + "/internal/concurrency\""
	if p.owned {
		imports = "\t\"" + g.Module + "/internal/authz\"\n" + imports
	}
	if p.hasFiles {
		imports += "\n\t\"" + g.Module + "/internal/files\""
	}
	imports += "\n\t\"" + g.Module + "/internal/models\"\n\t\"" + g.Module + "/internal/paginate\"\n\t\"" + g.Module + "/internal/respond\""
	if p.hasFiles {
		imports += "\n\t\"" + g.Module + "/internal/storage\""
	}

	storageField, storageNote := "", ""
	if p.hasFiles {
		storageField = "\n\t// Storage removes files a write replaced and claims the ones it keeps.\n\tStorage *storage.Storage"
	}

	// --owned-by: scope, guard and stamp, from the actor on the context.
	ownerScope, ownerExportScope, ownerBulkScope, ownerGuard, ownerStamp := "", "", "", "", ""
	if p.owned {
		ownerScope = "\n\t// --owned-by " + p.ownerName + ": a caller sees only their own rows. Without this" +
			"\n\t// the list hands over every row, and the ids below stop being worth" +
			"\n\t// protecting. ADMIN, and the system, are exempt." +
			"\n\tquery = authz.ScopeOwned(ctx, query, \"" + p.ownerCol + "\")\n"
		ownerExportScope = "\n\t// --owned-by " + p.ownerName + ": an export is the list without pages, scoped the same way." +
			"\n\tquery = authz.ScopeOwned(ctx, query, \"" + p.ownerCol + "\")\n"
		ownerBulkScope = "\t// --owned-by " + p.ownerName + ": only the caller's own rows are acted on, and" +
			"\n\t// somebody else's id drops out as if it did not exist." +
			"\n\tscope = authz.ScopeOwned(ctx, scope, \"" + p.ownerCol + "\")\n"
		ownerGuard = "\t// --owned-by " + p.ownerName + ": somebody else's row is not found rather than" +
			"\n\t// forbidden, so a wrong guess cannot be told from a right one." +
			"\n\tif !authz.Owns(ctx, &item) {" +
			"\n\t\treturn nil, gorm.ErrRecordNotFound" +
			"\n\t}\n"
		ownerStamp = "\t// --owned-by " + p.ownerName + ": the owner is whoever is signed in, never the body." +
			"\n\titem." + p.ownerField + " = authz.UserIDFrom(ctx)\n"
	}

	reload := "\tif err := db" + p.preloads + ".First(item, \"id = ?\", item.ID).Error; err != nil {\n\t\treturn %s\n\t}\n"
	createReload := fmt.Sprintf(reload, "err")
	updateReload := fmt.Sprintf(reload, "nil, err")
	createCall, updateCall := "db.Create(item)", "db.Model(item).Scopes(pre.Scope).Updates(updates)"
	gormClause, writeHelper := "", ""
	if p.single {
		// RETURNING where the dialect has it; the reload stays for the ones
		// that do not.
		createCall = "s.write(db).Create(item)"
		updateCall = "s.write(db).Model(item).Scopes(pre.Scope).Updates(updates)"
		createReload = "\tif !s.returning(db) {\n\t" + strings.ReplaceAll(createReload, "\n\t", "\n\t\t") + "}\n"
		updateReload = "\tif !s.returning(db) {\n\t" + strings.ReplaceAll(updateReload, "\n\t", "\n\t\t") + "}\n"
		gormClause = "\n\t\"gorm.io/gorm/clause\""
		// Not replaced by the template's own pass: a Replacer does not rescan
		// what it has just put in.
		writeHelper = strings.ReplaceAll(serviceWriteHelper, "{{Pascal}}", names.Pascal)
	}

	createClaim, updateSnapshot, updateCleanup := "", "", ""
	if p.hasFiles {
		createClaim = "\tif s.Storage != nil {\n\t\tfiles.ClaimRefs(ctx, db, item)\n\t}\n"
		updateSnapshot = "\n\toldItem := *item // the files before the write, to find the ones it replaced"
		updateCleanup = "\tif s.Storage != nil {\n\t\tfiles.CleanupRemoved(ctx, s.Storage, &oldItem, item)\n\t\tfiles.ClaimRefs(ctx, db, item)\n\t}\n"
	}

	createBody := "\tif err := " + createCall + ".Error; err != nil {\n\t\treturn err\n\t}\n" + createReload
	if p.linksApply != "" {
		createBody = `	// The row and its links land together or not at all.
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
` + p.linksApply + `		return nil
	}); err != nil {
		return err
	}
` + createReload
	}

	updateBody := "\twritten := " + updateCall + `
	if err := written.Error; err != nil {
		return nil, err
	}
	// pre named a version this record has moved past: someone else saved
	// first. A conflict rather than overwriting their change.
	if pre.Missed(written) {
		return nil, s.conflict(ctx, item.ID)
	}
` + updateReload
	if p.linksApply != "" || p.itemsReplace != "" {
		updateBody = `	// The row, its links and its lines change together or not at all. The
	// lines used to be deleted and recreated outside any transaction with the
	// errors dropped, so a failed insert lost every line.
	missed := false
	if err := db.Transaction(func(tx *gorm.DB) error {
		written := tx.Model(item).Scopes(pre.Scope).Updates(updates)
		if written.Error != nil {
			return written.Error
		}
		if pre.Missed(written) {
			missed = true
			return nil
		}
` + p.linksApply + p.itemsReplace + `		return nil
	}); err != nil {
		return nil, err
	}
	if missed {
		return nil, s.conflict(ctx, item.ID)
	}
` + updateReload
	}

	r := strings.NewReplacer(
		"{{IMPORTS}}", imports,
		"{{GORM_CLAUSE}}", gormClause,
		"{{WRITE_HELPER}}", writeHelper,
		"{{STORAGE_FIELD}}", storageField,
		"{{STORAGE_NOTE}}", storageNote,
		"{{SEARCH_COLS}}", p.searchCols,
		"{{SORT_COLS}}", p.sortCols,
		"{{FILTER_COLS}}", p.filterCols,
		"{{WRITABLE}}", p.writable,
		"{{LINKS_TYPE}}", p.linksType,
		"{{LINKS_PARAM}}", p.linksParam,
		"{{ITEMS_PARAM}}", p.itemsParam,
		"{{PRELOADS}}", p.preloads,
		"{{OWNER_SCOPE}}", ownerScope,
		"{{OWNER_EXPORT_SCOPE}}", ownerExportScope,
		"{{OWNER_BULK_SCOPE}}", ownerBulkScope,
		"{{OWNER_GUARD}}", ownerGuard,
		"{{OWNER_STAMP}}", ownerStamp,
		"{{CREATE_BODY}}", createBody,
		"{{CREATE_CLAIM}}", createClaim,
		"{{UPDATE_SNAPSHOT}}", updateSnapshot,
		"{{UPDATE_BODY}}", updateBody,
		"{{UPDATE_CLEANUP}}", updateCleanup,
		// Final text: a Replacer does not rescan what it puts in.
		"{{PUBLIC_METHODS}}", g.publicServiceMethods(names),
		"{{PATCH_M2M}}", p.patchM2M,
		"{{Pascal}}", names.Pascal,
		"{{camel}}", names.Camel,
		"{{lower}}", names.Lower,
		"{{plural}}", names.Plural,
	)
	return r.Replace(`package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"{{GORM_CLAUSE}}

{{IMPORTS}}
)

// {{Pascal}}Service owns every database read and write for {{plural}}.
//
// The handler reads the request, calls one of these and writes the answer. A
// job, a command or a test calls the same methods with no request at all, so a
// rule enforced here is enforced everywhere, not only on the HTTP route.
//
// Each method takes the context it runs in. It carries the organization the
// multitenant plugin resolved, the actor ownership is scoped by (see
// authz.WithActor, and authz.AsSystem for a job), and the cancellation that
// fires when a client goes away.
type {{Pascal}}Service struct {
	DB *gorm.DB{{STORAGE_FIELD}}
}

// {{camel}}ListConfig is what a client may search, sort and filter {{plural}}
// by. Whitelisted, because each name ends up in SQL.
var {{camel}}ListConfig = paginate.Config{
	Searchable: []string{{{SEARCH_COLS}}},
	Sortable:   map[string]bool{{{SORT_COLS}}},
	Filterable: map[string]bool{{{FILTER_COLS}}},
}

// writable{{Pascal}} is every column Patch and Bulk may write. id, the
// timestamps and the version are the framework's, and are dropped.
var writable{{Pascal}} = map[string]bool{
{{WRITABLE}}}

{{LINKS_TYPE}}
// db binds the database to ctx, so whatever a middleware put there reaches
// GORM's callbacks: the multitenant plugin scopes by it.
func (s *{{Pascal}}Service) db(ctx context.Context) *gorm.DB {
	return s.DB.WithContext(ctx)
}
{{WRITE_HELPER}}
// List returns one page of {{plural}}.
//
//	archived "true" or "1"   only archived rows
//	archived "all"           both
//	anything else            only live rows
func (s *{{Pascal}}Service) List(ctx context.Context, p paginate.Params, archived string) (paginate.Result[models.{{Pascal}}], error) {
	query := s.db(ctx).Model(&models.{{Pascal}}{}){{PRELOADS}}

	// Archived rows are excluded by default. Anything else means an operator
	// archives twelve rows, sees the count go down, and finds them again the
	// next time somebody sorts by a different column.
	switch archived {
	case "true", "1":
		query = query.Where("archived_at IS NOT NULL")
	case "all":
		// no filter
	default:
		query = query.Where("archived_at IS NULL")
	}
{{OWNER_SCOPE}}
	return paginate.List[models.{{Pascal}}](query, p, {{camel}}ListConfig)
}

// Export hands every matching {{lower}} to each, a batch at a time, so a large
// table is never in memory at once. search matches the columns List searches.
//
// No ORDER BY of its own: FindInBatches pages by primary key, which is
// creation order for the time-ordered ids Grit issues, and a sort in front of
// that key repeated rows from the second batch on.
func (s *{{Pascal}}Service) Export(ctx context.Context, search string, each func(rows []models.{{Pascal}}) error) error {
	query := s.db(ctx).Model(&models.{{Pascal}}{}){{PRELOADS}}
	if search != "" {
		clause := ""
		args := []any{}
		wild := "%" + search + "%"
		for i, col := range {{camel}}ListConfig.Searchable {
			if i > 0 {
				clause += " OR "
			}
			clause += "LOWER(" + col + ") LIKE LOWER(?)"
			args = append(args, wild)
		}
		if clause != "" {
			query = query.Where(clause, args...)
		}
	}
{{OWNER_EXPORT_SCOPE}}
	// FindInBatches fills rows and pages by primary key. The tx it hands the
	// callback is a fresh session with no query on it, so each batch is read
	// from rows: re-reading it through tx.Scan found nothing, and every export
	// was an empty file with a 200.
	var rows []models.{{Pascal}}
	return query.FindInBatches(&rows, 1000, func(tx *gorm.DB, batch int) error {
		return each(rows)
	}).Error
}

// GetByID returns one {{lower}} with the relations shown beside it.
func (s *{{Pascal}}Service) GetByID(ctx context.Context, id string) (*models.{{Pascal}}, error) {
	var item models.{{Pascal}}
	if err := s.db(ctx){{PRELOADS}}.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
{{OWNER_GUARD}}	return &item, nil
}

// load reads the row a write is about to change, without its relations.
func (s *{{Pascal}}Service) load(ctx context.Context, id string) (*models.{{Pascal}}, error) {
	var item models.{{Pascal}}
	if err := s.db(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
{{OWNER_GUARD}}	return &item, nil
}

// conflict is the answer to a write whose precondition failed: the version
// the row is at now.
func (s *{{Pascal}}Service) conflict(ctx context.Context, id string) error {
	var current models.{{Pascal}}
	if err := s.db(ctx).Select("version").First(&current, "id = ?", id).Error; err != nil {
		return err
	}
	return &concurrency.ErrConflict{Current: current.Version}
}

// Create saves a new {{lower}} and fills item in as it was stored.
func (s *{{Pascal}}Service) Create(ctx context.Context, item *models.{{Pascal}}{{LINKS_PARAM}}) error {
	db := s.db(ctx)
{{OWNER_STAMP}}{{CREATE_BODY}}{{CREATE_CLAIM}}	return nil
}

// Update writes updates to one {{lower}}. With a precondition it lands only
// if the row is still at that version, and otherwise returns an
// *concurrency.ErrConflict naming the version it is at.
func (s *{{Pascal}}Service) Update(ctx context.Context, id string, updates map[string]interface{}{{LINKS_PARAM}}{{ITEMS_PARAM}}, pre *concurrency.Precondition) (*models.{{Pascal}}, error) {
	item, err := s.load(ctx, id)
	if err != nil {
		return nil, err
	}
	db := s.db(ctx){{UPDATE_SNAPSHOT}}
{{UPDATE_BODY}}{{UPDATE_CLEANUP}}	return item, nil
}

// Patch writes only the columns body names, leaving every other one as it
// is. Keys that are not writable columns are dropped. It returns the row and
// the columns it wrote.
func (s *{{Pascal}}Service) Patch(ctx context.Context, id string, body map[string]interface{}, pre *concurrency.Precondition) (*models.{{Pascal}}, map[string]interface{}, error) {
	item, err := s.load(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	db := s.db(ctx)

	updates := map[string]interface{}{}
	for k, v := range body {
		if writable{{Pascal}}[k] {
			updates[k] = v
		}
	}
{{PATCH_M2M}}	if len(updates) == 0 && !touchedRelations {
		return nil, nil, respond.Rule("No writable fields in request body")
	}

	if len(updates) > 0 {
		written := db.Model(item).Scopes(pre.Scope).Updates(updates)
		if err := written.Error; err != nil {
			return nil, nil, err
		}
		if pre.Missed(written) {
			return nil, nil, s.conflict(ctx, item.ID)
		}
	}
	if err := db{{PRELOADS}}.First(item, "id = ?", item.ID).Error; err != nil {
		return nil, nil, err
	}
	return item, updates, nil
}

// Delete soft-deletes one {{lower}} and returns it as it was.
func (s *{{Pascal}}Service) Delete(ctx context.Context, id string) (*models.{{Pascal}}, error) {
	item, err := s.load(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.db(ctx).Delete(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

// {{Pascal}}BulkResult is what a bulk action did.
type {{Pascal}}BulkResult struct {
	// IDs are the rows acted on: those requested that exist, that the caller
	// may touch, and that the action applies to.
	IDs []string
	// Updates are the columns a patch wrote, after the whitelist.
	Updates map[string]interface{}
}

// Bulk applies one action to many {{plural}} in a single transaction: all of
// it lands or none of it does. action is delete, archive, restore or patch.
func (s *{{Pascal}}Service) Bulk(ctx context.Context, action string, ids []string, patch map[string]interface{}) ({{Pascal}}BulkResult, error) {
	db := s.db(ctx)
	var result {{Pascal}}BulkResult

	// Unarchived rows for archive, archived for restore: without it a mixed
	// selection reports "12 archived" having changed three.
	scope := db.Model(&models.{{Pascal}}{}).Where("id IN ?", ids)
	if action == "restore" {
		scope = scope.Where("archived_at IS NOT NULL")
	} else if action == "archive" {
		scope = scope.Where("archived_at IS NULL")
	}
{{OWNER_BULK_SCOPE}}	var items []models.{{Pascal}}
	if err := scope.Find(&items).Error; err != nil {
		return result, fmt.Errorf("loading {{plural}}: %w", err)
	}
	for _, item := range items {
		result.IDs = append(result.IDs, item.ID)
	}
	if len(result.IDs) == 0 {
		return result, nil
	}

	if action == "patch" {
		// The same whitelist as Patch. Framework-owned columns are dropped
		// rather than refused, so a client sending the whole row is not wrong.
		result.Updates = map[string]interface{}{}
		for k, v := range patch {
			if writable{{Pascal}}[k] {
				result.Updates[k] = v
			}
		}
		if len(result.Updates) == 0 {
			return {{Pascal}}BulkResult{}, respond.Rule("No writable fields in patch")
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		switch action {
		case "delete":
			return tx.Where("id IN ?", result.IDs).Delete(&models.{{Pascal}}{}).Error
		case "archive":
			return tx.Model(&models.{{Pascal}}{}).Where("id IN ?", result.IDs).
				Update("archived_at", time.Now()).Error
		case "restore":
			return tx.Model(&models.{{Pascal}}{}).Where("id IN ?", result.IDs).
				Update("archived_at", nil).Error
		case "patch":
			return tx.Model(&models.{{Pascal}}{}).Where("id IN ?", result.IDs).
				Updates(result.Updates).Error
		}
		return respond.Rule("unknown bulk action %q", action)
	})
	if err != nil {
		return {{Pascal}}BulkResult{}, err
	}
	return result, nil
}
{{PUBLIC_METHODS}}`)
}

// writeGoHandler writes handlers/<resource>.go: the HTTP edge of the resource,
// with no query of its own.
func (g *Generator) writeGoHandler(names Names) error {
	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+".go")
	return writeFileWithDirs(path, g.handlerSource(names))
}

// handlerSource renders the handler.
func (g *Generator) handlerSource(names Names) string {
	p := g.crud(names)
	ar := g.auditReadSnippets(names)

	storageField, storageAssign := "", ""
	if p.hasFiles {
		storageField = "\n\tStorage *storage.Storage"
		storageAssign = ", Storage: h.Storage"
	}

	r := strings.NewReplacer(
		"{{AUDIT_IMPORT}}", ar.Import,
		"{{AUDIT_READ_LIST}}", ar.List,
		"{{AUDIT_READ_ONE}}", ar.One,
		"{{AUDIT_EXPORT_DECL}}", ar.ExportDecl,
		"{{AUDIT_EXPORT_COUNT}}", ar.ExportCount,
		"{{AUDIT_EXPORT_MARK}}", ar.ExportMark,
		"{{AUDIT_XLSX_MARK}}", ar.XLSXMark,
		"{{DATATYPES_IMPORT}}", p.datatypesImport,
		"{{EXTRA_IMPORTS}}", p.handlerImports,
		"{{STORAGE_FIELD}}", storageField,
		"{{STORAGE_ASSIGN}}", storageAssign,
		"{{OPTIONAL_ID_HELPER}}", p.optionalID,
		"{{FK_FILTERS}}", p.fkFilters,
		"{{EXPORT_COLS}}", p.exportCols,
		"{{UPPER_LABEL}}", strings.ToUpper(strings.Join(splitPascal(names.Pascal), " ")),
		"{{PDF_SUBTITLE}}", "pdf.Value("+p.identExpr+")",
		"{{PDF_FIELDS}}", p.pdfFields,
		"{{PDF_SECTIONS}}", p.pdfSections,
		"{{CREATE_FIELDS_TOP}}", dedentOneTab(p.createFields),
		"{{UPDATE_FIELDS_TOP}}", dedentOneTab(p.updateFields),
		"{{CREATE_ASSIGN}}", p.createAssign,
		"{{UPDATE_MAP}}", p.updateMap,
		"{{ITEMS_BUILD}}", p.itemsBuild,
		"{{ITEMS_FROM_UPDATE}}", p.itemsFromUpdate,
		"{{LINKS_FROM_CREATE}}", p.linksFromCreate,
		"{{LINKS_FROM_UPDATE}}", p.linksFromUpdate,
		"{{LINKS_ARG}}", p.linksArg,
		"{{ITEMS_ARG}}", p.itemsArg,
		"{{IDENT_EXPR}}", p.identExpr,
		"{{kebab}}", names.Kebab,
		"{{MODULE}}", g.Module,
		"{{Pascal}}", names.Pascal,
		"{{lower}}", names.Lower,
		"{{plural}}", names.Plural,
		"{{Plural}}", names.PluralPascal,
	)

	content := r.Replace(`package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"{{DATATYPES_IMPORT}}
	"gorm.io/gorm"

	{{AUDIT_IMPORT}}"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/concurrency"
	"{{MODULE}}/internal/events"
	"{{MODULE}}/internal/export"{{EXTRA_IMPORTS}}
	"{{MODULE}}/internal/models"
	"{{MODULE}}/internal/paginate"
	"{{MODULE}}/internal/pdf"
	"{{MODULE}}/internal/respond"
	"{{MODULE}}/internal/services"
)

// {{Pascal}}Handler serves the {{lower}} endpoints. It reads the request, asks
// services.{{Pascal}}Service, and writes the answer; it runs no query of its
// own, so everything a route does is also available to a job or a test.
type {{Pascal}}Handler struct {
	DB *gorm.DB{{STORAGE_FIELD}}
}

// service is the {{lower}} service over this handler's database.
func (h *{{Pascal}}Handler) service() *services.{{Pascal}}Service {
	return &services.{{Pascal}}Service{DB: h.DB{{STORAGE_ASSIGN}}}
}

// ctx is the request's context with the caller on it, which the service
// scopes owned rows by. The organization the multitenant plugin resolved and
// the cancellation when the client goes away travel with it.
func (h *{{Pascal}}Handler) ctx(c *gin.Context) context.Context {
	return authz.WithActor(c.Request.Context(), authz.ActorOf(c))
}

// fail answers an error from the service: a version conflict with the version
// the record is at, a missing row with 404, a broken rule with 422 and its
// message, and anything else as an opaque 500 that is logged.
func (h *{{Pascal}}Handler) fail(c *gin.Context, err error, fallback string) {
	var conflict *concurrency.ErrConflict
	switch {
	case errors.As(err, &conflict):
		concurrency.WriteConflict(c, conflict.Current)
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "{{Pascal}} not found",
			},
		})
	default:
		respond.WriteError(c, err, fallback)
	}
}

{{OPTIONAL_ID_HELPER}}

// List returns a paginated list of {{plural}}.
//
//	?archived=true   only archived rows
//	?archived=all    both
//	(default)        only live rows
func (h *{{Pascal}}Handler) List(c *gin.Context) {
	res, err := h.service().List(h.ctx(c), paginate.Bind(c){{FK_FILTERS}}, c.Query("archived"))
	if err != nil {
		h.fail(c, err, "Failed to fetch {{plural}}")
		return
	}

{{AUDIT_READ_LIST}}	c.JSON(http.StatusOK, res)
}

// Export streams the full filtered list as CSV (default) or XLSX: the same
// search as List, without pages.
//
//	GET /api/{{plural}}/export?format=csv
//	GET /api/{{plural}}/export?format=xlsx&search=foo
func (h *{{Pascal}}Handler) Export(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")
	search := c.Query("search")

	opts := export.Options{
		Sheet: "{{Plural}}",
		Columns: []export.Column{
{{EXPORT_COLS}}
		},
	}

	if format == "xlsx" {
		// excelize has no streaming writer, so the sheet is built in memory.
		var all []models.{{Pascal}}
		if err := h.service().Export(h.ctx(c), search, func(rows []models.{{Pascal}}) error {
			all = append(all, rows...)
			return nil
		}); err != nil {
			h.fail(c, err, "Failed to export {{plural}}")
			return
		}
{{AUDIT_XLSX_MARK}}		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", ` + "`" + `attachment; filename="{{plural}}.xlsx"` + "`" + `)
		if err := export.XLSX(c.Writer, all, opts); err != nil {
			log.Printf("export {{plural}} as xlsx: %v", err)
		}
		return
	}

	// CSV streams: the header with the first batch, then each batch as it is
	// read.
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", ` + "`" + `attachment; filename="{{plural}}.csv"` + "`" + `)

{{AUDIT_EXPORT_DECL}}	headerWritten := false
	err := h.service().Export(h.ctx(c), search, func(rows []models.{{Pascal}}) error {
{{AUDIT_EXPORT_COUNT}}		if !headerWritten {
			headerWritten = true
			return export.CSV(c.Writer, rows, opts)
		}
		return export.CSVRows(c.Writer, rows, opts)
	})
	if err == nil && !headerWritten {
		// Nothing matched: still a CSV, with its header row.
		err = export.CSV(c.Writer, []models.{{Pascal}}{}, opts)
	}
	if err != nil {
		if !c.Writer.Written() {
			c.Writer.Header().Del("Content-Type")
			c.Writer.Header().Del("Content-Disposition")
			h.fail(c, err, "Failed to export {{plural}}")
			return
		}
		// Rows are already on the wire, so the status cannot change. Logged,
		// so a truncated file has an explanation somewhere.
		log.Printf("export {{plural}}: %v", err)
	}
{{AUDIT_EXPORT_MARK}}}

// GetByID returns a single {{lower}} by ID.
func (h *{{Pascal}}Handler) GetByID(c *gin.Context) {
	item, err := h.service().GetByID(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to load {{lower}}")
		return
	}
{{AUDIT_READ_ONE}}
	// The version to send back as If-Match when saving.
	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}

// PDF streams this {{lower}} as a print-ready PDF: a repeating header and
// footer with page numbers, the record's fields as a detail grid, and any
// line items as a table. Edit the pdf.Record below to restyle it; the
// renderer itself lives in internal/pdf/record.go.
func (h *{{Pascal}}Handler) PDF(c *gin.Context) {
	id := c.Param("id")

	item, err := h.service().GetByID(h.ctx(c), id)
	if err != nil {
		h.fail(c, err, "Failed to load {{lower}}")
		return
	}
{{AUDIT_READ_ONE}}
	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "{{Pascal}}"
	}

	rec := pdf.Record{
		Title:      "{{UPPER_LABEL}}",
		Subtitle:   {{PDF_SUBTITLE}},
		Brand:      appName,
		FooterNote: appName + " · generated " + time.Now().Format("2 Jan 2006 15:04"),
		Fields: []pdf.Field{
{{PDF_FIELDS}}		},
	}
{{PDF_SECTIONS}}
	out, err := pdf.RenderRecord(rec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "PDF_ERROR",
				"message": "could not render the PDF",
			},
		})
		return
	}

	filename := "{{kebab}}-" + id + ".pdf"
	c.Header("Content-Disposition", "inline; filename=\""+filename+"\"")
	c.Data(http.StatusOK, "application/pdf", out)
}

// Create{{Pascal}}Request is the JSON body accepted by POST /{{plural}}.
//
// Named rather than anonymous so the API reference can document it: gindocs
// builds a request schema by reflecting over a real type. routes.go passes
// this type to docs.Route(...).RequestBody().
type Create{{Pascal}}Request struct {
{{CREATE_FIELDS_TOP}}}

// Update{{Pascal}}Request is the JSON body accepted by PUT /{{plural}}/:id.
// Every field is optional: only what the client sends is applied.
type Update{{Pascal}}Request struct {
{{UPDATE_FIELDS_TOP}}}

// Create adds a new {{lower}}.
func (h *{{Pascal}}Handler) Create(c *gin.Context) {
	var req Create{{Pascal}}Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	item := models.{{Pascal}}{
{{CREATE_ASSIGN}}	}
{{ITEMS_BUILD}}{{LINKS_FROM_CREATE}}
	if err := h.service().Create(h.ctx(c), &item{{LINKS_ARG}}); err != nil {
		h.fail(c, err, "Failed to create {{lower}}")
		return
	}

	events.Emitted(c, "{{plural}}", "{{Pascal}}", "created", item.ID, {{IDENT_EXPR}}, "", nil, item)

	c.JSON(http.StatusCreated, gin.H{
		"data":    item,
		"message": "{{Pascal}} created successfully",
	})
}

// Update modifies an existing {{lower}}. A client that sends the version it
// read as If-Match gets 409 instead of overwriting a newer save.
func (h *{{Pascal}}Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req Update{{Pascal}}Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	updates := map[string]interface{}{}
{{UPDATE_MAP}}{{LINKS_FROM_UPDATE}}{{ITEMS_FROM_UPDATE}}
	item, err := h.service().Update(h.ctx(c), id, updates{{LINKS_ARG}}{{ITEMS_ARG}}, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to update {{lower}}")
		return
	}

	events.Emitted(c, "{{plural}}", "{{Pascal}}", "updated", item.ID, {{IDENT_EXPR}}, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "{{Pascal}} updated successfully",
	})
}

// Patch applies a partial update to a {{lower}}. Used by the admin's grouped
// update view: each form group's Save button sends only the fields it owns, so
// editing "Address" does not rewrite "Pricing".
func (h *{{Pascal}}Handler) Patch(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	item, updates, err := h.service().Patch(h.ctx(c), c.Param("id"), body, concurrency.FromRequest(c))
	if err != nil {
		h.fail(c, err, "Failed to patch {{lower}}")
		return
	}

	events.Emitted(c, "{{plural}}", "{{Pascal}}", "updated", item.ID, {{IDENT_EXPR}}, services.DiffSummary(updates), nil, item)

	c.Header("ETag", concurrency.Tag(item.Version))
	c.JSON(http.StatusOK, gin.H{
		"data":    item,
		"message": "{{Pascal}} updated successfully",
	})
}

// Delete soft-deletes a {{lower}}.
func (h *{{Pascal}}Handler) Delete(c *gin.Context) {
	item, err := h.service().Delete(h.ctx(c), c.Param("id"))
	if err != nil {
		h.fail(c, err, "Failed to delete {{lower}}")
		return
	}

	events.Emitted(c, "{{plural}}", "{{Pascal}}", "deleted", item.ID, {{IDENT_EXPR}}, "", item, nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "{{Pascal}} deleted successfully",
	})
}

// Bulk{{Pascal}}Request is one operation applied to a set of rows.
//
// One request rather than one per row. The admin used to fire N parallel
// DELETEs, which means N transactions, N audit entries, and a half-applied
// result when the eleventh fails.
type Bulk{{Pascal}}Request struct {
	// delete removes, archive puts away, restore brings back, patch writes the
	// same field values to every selected row.
	Action string ` + "`" + `json:"action" binding:"required,oneof=delete archive restore patch"` + "`" + `
	// Capped: an unbounded IN clause is a way to lock a table by accident.
	IDs []string ` + "`" + `json:"ids" binding:"required,min=1,max=500"` + "`" + `
	// Only read when action is "patch". Whitelisted the same way Patch is.
	Patch map[string]interface{} ` + "`" + `json:"patch"` + "`" + `
}

// Bulk applies one action to many {{plural}} in a single transaction.
func (h *{{Pascal}}Handler) Bulk(c *gin.Context) {
	var req Bulk{{Pascal}}Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	result, err := h.service().Bulk(h.ctx(c), req.Action, req.IDs, req.Patch)
	if err != nil {
		h.fail(c, err, "Failed to "+req.Action+" {{plural}}")
		return
	}
	if len(result.IDs) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"data":    gin.H{"affected": 0, "requested": len(req.IDs)},
			"message": "Nothing to do",
		})
		return
	}
	ids := result.IDs

	// One audit entry naming the action and the count, not N entries that bury
	// everything else somebody did today. A local map, not a package-level
	// helper: every resource has its own handler file in package handlers.
	past := map[string]string{
		"delete":  "deleted",
		"archive": "archived",
		"restore": "restored",
		"patch":   "updated",
	}[req.Action]

	noun := "{{plural}}"
	if len(ids) == 1 {
		noun = "{{lower}}"
	}

	summary := req.Action + " " + strconv.Itoa(len(ids)) + " " + noun
	if req.Action == "patch" {
		summary += ": " + services.DiffSummary(result.Updates)
	}
	// resourceID holds ONE id, not all of them: it is a lookup key, and the
	// count lives in the summary, where it can be read.
	events.Emitted(c, "{{plural}}", "{{Pascal}}", "bulk", ids[0], summary, summary, nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"affected": len(ids), "requested": len(req.IDs)},
		"message": strconv.Itoa(len(ids)) + " " + noun + " " + past,
	})
}
`)

	// The transition endpoint, only for a resource whose status field is a
	// state machine.
	return content + g.workflowHandlerMethod(names)
}
