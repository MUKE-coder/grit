package generate

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// ModelShape hashes what the generator reads out of a model, as opposed to the
// bytes of the file it reads it from.
//
// It exists because the two are not the same question, and `grit doctor` was
// asking the wrong one. A generated file is behind its model when the model
// says something different about the columns; the manifest recorded the model's
// bytes, so any edit at all counted, including the edits that change nothing a
// generator would write. The case that forced it (grit#307) is a gorm tag
// pinning a column's type:
//
//	FxRate float64 `gorm:"type:numeric(18,8)" json:"fx_rate"`
//
// That tag is the only place the instruction can live, because `grit migrate`
// runs AutoMigrate first and would alter the column back to match the Go type.
// It is therefore a correct and expected hand edit, and it made the model
// permanently drifted, with a printed remedy that would have reverted it.
//
// The shape is the definition the model rebuilds into, which is exactly the
// input every downstream generator works from: RebuildDefinitionFromModel is
// what `grit generate resource <Name>` uses when the --fields are not given
// again. So this cannot fall behind what the generators read by being a second
// opinion about it, which a hand-maintained list of "tag keys that matter"
// would have been. A tag the rebuild keeps (type:text, foreignKey, uniqueIndex)
// changes the shape and is still drift; one it discards does not.
//
// Empty when the shape cannot be computed: a file that does not parse, a model
// with no declarable fields, a path that is not a model. An empty shape is
// recorded as absent and the comparison falls back to the bytes, which is the
// behaviour this replaces.
func ModelShape(modelPath string) string {
	def := shapeDefinition(modelPath)
	if def == nil {
		return ""
	}
	// Field order is source order and every map Go marshals is key-sorted, so
	// one model gives one shape on every machine.
	encoded, err := json.Marshal(def)
	if err != nil {
		return ""
	}
	return manifest.Hash(string(encoded))
}

// shapeDefinition rebuilds the definition a model file describes, or nil.
//
// modelForRebuild is the mode that keeps every column somebody declared,
// including the slug and the foreign keys a seeder has to leave out. That is
// the right mode here for the same reason it is right there: this stands in
// for the whole model, not for a subset of it.
func shapeDefinition(modelPath string) *ResourceDefinition {
	name := modelNameFromPath(modelPath)
	if name == "" {
		return nil
	}
	apiRoot := apiRootFromModelPath(modelPath)
	if apiRoot == "" {
		return nil
	}
	def, err := definitionFromModelFileMode(apiRoot, name, modelForRebuild)
	if err != nil || def == nil || len(def.Fields) == 0 {
		return nil
	}
	return def
}

// modelNameFromPath turns .../internal/models/blog_post.go into BlogPost.
func modelNameFromPath(modelPath string) string {
	base := strings.TrimSuffix(filepath.Base(modelPath), ".go")
	if base == "" || base == filepath.Base(modelPath) {
		return "" // not a .go file
	}
	return MakeNames(base).Pascal
}

// apiRootFromModelPath walks up from internal/models/<x>.go to the directory
// the API module lives in, which is what definitionFromModelFileMode joins
// its own paths onto.
func apiRootFromModelPath(modelPath string) string {
	dir := filepath.Dir(modelPath) // .../internal/models
	if filepath.Base(dir) != "models" {
		return ""
	}
	internal := filepath.Dir(dir) // .../internal
	if filepath.Base(internal) != "internal" {
		return ""
	}
	return filepath.Dir(internal)
}
