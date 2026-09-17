package scaffold

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// ─── L17: a generated write does not read its row back ───────────────────────
//
// A generated service wrote a row and then read it again with its relations
// preloaded: a create was BEGIN, INSERT, COMMIT, the row, and one query per
// relation, and an update and a patch were the same after loading the row they
// change. A resource with no relations already wrote with RETURNING in one
// statement, except in Patch, which still read the row back.
//
// Now a resource whose relations are all belongs_to writes the same way, and
// then reads only the rows it points at. MySQL has no RETURNING, so there the
// row is still read back with its relations.
//
// The text below is shared by the generator (internal/generate renders it into
// every new service) and by grit upgrade, which rewrites a service still in the
// shape the generator wrote.

// ServiceWriteHelper is the single-statement write, as two methods on a
// generated service. They do what database.Write and database.SupportsReturning
// do, and they are in the service rather than there because the database
// package imports services for its seeders: a service that imported database
// back would be an import cycle, and the project would not build. Methods, so
// that two resources in one package cannot collide. {{Pascal}} is the resource.
const ServiceWriteHelper = `
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

// ServiceRelation is one belongs_to relation a generated service fills in
// after a write.
type ServiceRelation struct {
	// Field is the association on the model, e.g. Group.
	Field string
	// Model is the related model, e.g. Group.
	Model string
	// FK is the foreign key field, e.g. GroupID.
	FK string
	// Nullable is a *string foreign key: a self-reference, where a root has none.
	Nullable bool
}

// ServiceRelationsMethod renders the relations method of a generated service.
func ServiceRelationsMethod(pascal string, rels []ServiceRelation) string {
	var b strings.Builder
	b.WriteString("\n// relations reads the rows item points at, after a write whose RETURNING\n")
	b.WriteString("// brought the row itself back. Reading the row again with its relations\n")
	b.WriteString("// preloaded cost one more query on every create, update and patch.\n")
	fmt.Fprintf(&b, "func (s *%sService) relations(db *gorm.DB, item *models.%s) error {\n", pascal, pascal)
	for _, r := range rels {
		fmt.Fprintf(&b, "\titem.%s = nil\n", r.Field)
		key := "item." + r.FK
		if r.Nullable {
			fmt.Fprintf(&b, "\tif item.%s != nil && *item.%s != \"\" {\n", r.FK, r.FK)
			key = "*item." + r.FK
		} else {
			fmt.Fprintf(&b, "\tif item.%s != \"\" {\n", r.FK)
		}
		fmt.Fprintf(&b, "\t\tvar related []models.%s\n", r.Model)
		fmt.Fprintf(&b, "\t\tif err := db.Where(\"id = ?\", %s).Limit(1).Find(&related).Error; err != nil {\n", key)
		b.WriteString("\t\t\treturn err\n\t\t}\n")
		b.WriteString("\t\tif len(related) == 1 {\n")
		fmt.Fprintf(&b, "\t\t\titem.%s = &related[0]\n", r.Field)
		b.WriteString("\t\t}\n\t}\n")
	}
	b.WriteString("\treturn nil\n}\n")
	return b.String()
}

// ServiceRelationsReload is what follows a write in a service with belongs_to
// relations. preloads is the Preload chain; ret is what the method returns on
// an error, e.g. "nil, err".
func ServiceRelationsReload(preloads, ret string) string {
	return "\t// RETURNING filled the row in, so only the rows it points at are read.\n" +
		"\tif s.returning(db) {\n" +
		"\t\tif err := s.relations(db, item); err != nil {\n" +
		"\t\t\treturn " + ret + "\n" +
		"\t\t}\n" +
		"\t} else if err := db" + preloads + ".First(item, \"id = ?\", item.ID).Error; err != nil {\n" +
		"\t\treturn " + ret + "\n" +
		"\t}\n"
}

// ServiceSingleReload is what follows a write in a service with no relations:
// the row is read again only where the dialect has no RETURNING.
func ServiceSingleReload(ret string) string {
	return "\tif !s.returning(db) {\n" +
		"\t\tif err := db.First(item, \"id = ?\", item.ID).Error; err != nil {\n" +
		"\t\t\treturn " + ret + "\n" +
		"\t\t}\n" +
		"\t}\n"
}

func plainReload(preloads, ret string) string {
	return "\tif err := db" + preloads + ".First(item, \"id = ?\", item.ID).Error; err != nil {\n\t\treturn " + ret + "\n\t}\n"
}

const (
	patchWriteOld = "\t\twritten := db.Model(item).Scopes(pre.Scope).Updates(updates)\n"
	patchWriteNew = "\t\twritten := s.write(db).Model(item).Scopes(pre.Scope).Updates(updates)\n"
)

var (
	serviceLoadRe    = regexp.MustCompile(`func \(s \*(\w+)Service\) load\(ctx context\.Context, id string\) \(\*models\.(\w+), error\) \{`)
	createPreloadsRe = regexp.MustCompile(`\n\tif err := db\.Create\(item\)\.Error; err != nil \{\n\t\treturn err\n\t\}\n\tif err := db((?:\.Preload\("\w+"\))+)\.First\(item, "id = \?", item\.ID\)\.Error; err != nil \{\n\t\treturn err\n\t\}\n`)
	modelAssocRe     = regexp.MustCompile("(?m)^\\t(\\w+)\\s+\\*(\\w+)\\s+`gorm:\"foreignKey:(\\w+)\"")
	preloadNameRe    = regexp.MustCompile(`\.Preload\("(\w+)"\)`)
)

// RepairGeneratedServiceWrites brings a generated service up to L17. model is
// the source of its model file. It reports whether it changed anything, and
// leaves a service alone when it is not in a shape the generator wrote: links
// or line items written in a transaction, a sequence or a tree hook, or code
// somebody has changed.
func RepairGeneratedServiceWrites(src, model string) (string, bool) {
	m := serviceLoadRe.FindStringSubmatch(src)
	if m == nil || m[1] != m[2] {
		return src, false
	}
	pascal := m[1]
	if strings.Contains(model, "sequence.") || strings.Contains(model, "resolveTreePath") {
		return src, false
	}

	// Already writing with RETURNING: a resource with no relations. Only Patch
	// still read its row back.
	if strings.Contains(src, "func (s *"+pascal+"Service) write(db *gorm.DB) *gorm.DB {") {
		if strings.Contains(src, patchWriteNew) || strings.Count(src, patchWriteOld) != 1 || strings.Count(src, plainReload("", "nil, nil, err")) != 1 {
			return src, false
		}
		out := strings.Replace(src, patchWriteOld, patchWriteNew, 1)
		out = strings.Replace(out, plainReload("", "nil, nil, err"), ServiceSingleReload("nil, nil, err"), 1)
		return out, true
	}

	cm := createPreloadsRe.FindStringSubmatch(src)
	if cm == nil {
		return src, false
	}
	preloads := cm[1]
	assocs := map[string]ServiceRelation{}
	for _, a := range modelAssocRe.FindAllStringSubmatch(model, -1) {
		nullable := regexp.MustCompile(`(?m)^\t` + a[3] + `\s+\*string\s`).MatchString(model)
		assocs[a[1]] = ServiceRelation{Field: a[1], Model: a[2], FK: a[3], Nullable: nullable}
	}
	var rels []ServiceRelation
	for _, name := range preloadNameRe.FindAllStringSubmatch(preloads, -1) {
		r, ok := assocs[name[1]]
		if !ok {
			return src, false // a has-many or many-to-many: not a single row to point at
		}
		rels = append(rels, r)
	}

	updateWriteOld := "\n\twritten := db.Model(item).Scopes(pre.Scope).Updates(updates)\n"
	updateWriteNew := "\n\twritten := s.write(db).Model(item).Scopes(pre.Scope).Updates(updates)\n"
	dbFunc := "func (s *" + pascal + "Service) db(ctx context.Context) *gorm.DB {\n\treturn s.DB.WithContext(ctx)\n}\n"
	const gormImport = "\t\"gorm.io/gorm\"\n"
	for _, anchor := range []struct {
		text string
		want int
	}{
		{updateWriteOld, 1}, {patchWriteOld, 1}, {dbFunc, 1}, {gormImport, 1},
		{plainReload(preloads, "nil, err"), 1}, {plainReload(preloads, "nil, nil, err"), 1},
	} {
		if strings.Count(src, anchor.text) != anchor.want {
			return src, false
		}
	}

	out := strings.Replace(src, cm[0],
		"\n\tif err := s.write(db).Create(item).Error; err != nil {\n\t\treturn err\n\t}\n"+ServiceRelationsReload(preloads, "err"), 1)
	out = strings.Replace(out, updateWriteOld, updateWriteNew, 1)
	out = strings.Replace(out, plainReload(preloads, "nil, err"), ServiceRelationsReload(preloads, "nil, err"), 1)
	out = strings.Replace(out, patchWriteOld, patchWriteNew, 1)
	out = strings.Replace(out, plainReload(preloads, "nil, nil, err"), ServiceRelationsReload(preloads, "nil, nil, err"), 1)
	helpers := strings.ReplaceAll(ServiceWriteHelper, "{{Pascal}}", pascal) + ServiceRelationsMethod(pascal, rels)
	out = strings.Replace(out, dbFunc, dbFunc+helpers, 1)
	out = strings.Replace(out, gormImport, gormImport+"\t\"gorm.io/gorm/clause\"\n", 1)
	return out, true
}

// repairGeneratedServiceWrites applies L17 to every generated service in an
// existing project.
func repairGeneratedServiceWrites(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	services, err := goSources(filepath.Join(apiRoot, "internal", "services"))
	if err != nil {
		return err
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, path := range services {
		model := filepath.Join(apiRoot, "internal", "models", filepath.Base(path))
		if strings.HasSuffix(path, "_test.go") || !fileExists(model) {
			continue
		}
		modelSrc, err := readNormalized(model)
		if err != nil {
			return err
		}
		if err := repairSourceFile(root, m, path, func(src string) (string, []string, []string) {
			out, changed := RepairGeneratedServiceWrites(src, modelSrc)
			if !changed {
				return src, nil, nil
			}
			return out, []string{"a create, update or patch writes with RETURNING and reads only the rows it points at, instead of reading the row back"}, nil
		}); err != nil {
			return err
		}
	}
	return nil
}
