package generate

import (
	"regexp"
	"strings"
)

// Reading a generated model back, for the shapes that have no other signal.
//
// ── Why these live here and not where they are used ─────────────────────
//
// `grit doctor` suggests a `retention.Register` line fitted to a resource's
// shape, and a record with lines or a tree needs options in it or the database
// refuses the nightly sweep for ever. The only place that shape is recorded is
// the model's own source, so the check reads it back.
//
// It did that with its own copy of the pattern, which is the arrangement this
// repository keeps getting wrong: the writer is in this package and the reader
// was in another, nothing made them agree, and the day somebody reformats the
// has-many or rewrites the tree banner the suggestion quietly loses its
// options. It still compiles. The sweep then fails every night on a foreign
// key, which is the failure retention exists to prevent.
//
// So the reader lives beside the writer, one function each, and
// TestTheModelReadersMatchWhatTheGeneratorWrites runs the real generator and
// asserts both. One definition cannot disagree with itself.

// itemsHasManyRe matches the has-many `--items` writes on the parent model:
//
//	Items []InvoiceLine `gorm:"foreignKey:InvoiceID" json:"items"`
//
// Written by writeGoModel, a few lines above treeFields. The capture is the
// child model's name.
var itemsHasManyRe = regexp.MustCompile("Items\\s+\\[\\](\\w+)\\s+`gorm:\"foreignKey:\\w+\"")

// treeBanner opens the block treeFields writes. Matched on the words alone:
// the rule of dashes after them is decoration, and matching it would fail on
// a reflow that changed nothing.
const treeBanner = "--- tree (grit --tree)"

// ItemsChildModel is the model holding this record's lines, or "" when the
// resource was not generated with --items.
func ItemsChildModel(modelSource string) string {
	if m := itemsHasManyRe.FindStringSubmatch(modelSource); m != nil {
		return m[1]
	}
	return ""
}

// IsTreeModel reports a resource generated with --tree, whose rows point at a
// parent row of their own.
func IsTreeModel(modelSource string) bool {
	return strings.Contains(modelSource, treeBanner)
}
