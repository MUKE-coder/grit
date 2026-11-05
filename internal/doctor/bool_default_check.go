package doctor

import "regexp"

// A bool column with a database default of true cannot store false.
//
// GORM omits a zero-valued field from an INSERT when the column carries a
// default. So on Create, a false bool is left out, the database fills in true,
// and GORM writes that true back into the struct. A form that saves "unticked"
// stores "ticked", and nothing errors anywhere: not the handler, not the
// service, not the database.
//
// It has now been found three times in this project's own models, on
// User.Active, on the SSO flags and on a product variant, and twice more in a
// real application, on marketing consent and on a "sends email" flag (grit#179).
// Marketing consent is the one to think about: the default was true, so a
// person who declined was recorded as having agreed.
//
// The fix is to drop the default and set the value on every create path, which
// is what the framework's own models do, or to make the field a *bool so that
// "not set" and "false" are different values.

// boolDefaultTrueRe is a struct field of type bool whose gorm tag carries
// default:true, in either order and with other tag parts around it.
//
// Not *bool: a pointer has no zero-value problem, because nil and false are
// different things to GORM.
var boolDefaultTrueRe = regexp.MustCompile(
	"(?m)^\\s*(\\w+)\\s+bool\\s+`[^`]*gorm:\"[^\"]*default:(?i:true)[^\"]*\"")

func checkBoolDefaultTrue(p *project) []Finding {
	var out []Finding
	for _, name := range p.resources {
		m := p.models[name]
		if m == nil || m.Source == "" {
			continue
		}
		for _, match := range boolDefaultTrueRe.FindAllStringSubmatch(m.Source, -1) {
			out = append(out, Finding{
				Level:    "error",
				Resource: name + "." + match[1],
				Message: "is a bool with a database default of true, so it can never be " +
					"saved as false: GORM leaves a zero value out of the INSERT, the database " +
					"fills in true, and the row is stored ticked with no error",
				Fix: "drop default:true and set the field on every create path, as the " +
					"framework's own models do, or make it *bool so nil and false differ",
			})
		}
	}
	return out
}
