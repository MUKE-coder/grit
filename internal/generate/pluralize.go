package generate

import "strings"

// Pluralize returns the plural form of a word.
// Handles common English pluralization rules.
func Pluralize(word string) string {
	if word == "" {
		return word
	}

	lower := strings.ToLower(word)

	// Irregular plurals
	irregulars := map[string]string{
		"person":  "people",
		"child":   "children",
		"mouse":   "mice",
		"goose":   "geese",
		"man":     "men",
		"woman":   "women",
		"tooth":   "teeth",
		"foot":    "feet",
		"ox":      "oxen",
		"datum":   "data",
		"medium":  "media",
		"index":   "indices",
		"matrix":  "matrices",
		"vertex":  "vertices",
		"crisis":  "crises",
		"axis":    "axes",
		"analysis":"analyses",
	}

	if plural, ok := irregulars[lower]; ok {
		if word[0] >= 'A' && word[0] <= 'Z' {
			return strings.ToUpper(plural[:1]) + plural[1:]
		}
		return plural
	}

	// Words ending in -y preceded by a consonant → -ies
	if strings.HasSuffix(lower, "y") && len(lower) > 1 {
		prev := lower[len(lower)-2]
		if prev != 'a' && prev != 'e' && prev != 'i' && prev != 'o' && prev != 'u' {
			return word[:len(word)-1] + "ies"
		}
	}

	// Words ending in -s, -ss, -sh, -ch, -x, -z → -es
	if strings.HasSuffix(lower, "s") ||
		strings.HasSuffix(lower, "ss") ||
		strings.HasSuffix(lower, "sh") ||
		strings.HasSuffix(lower, "ch") ||
		strings.HasSuffix(lower, "x") ||
		strings.HasSuffix(lower, "z") {
		return word + "es"
	}

	// Words ending in -f or -fe → -ves
	if strings.HasSuffix(lower, "fe") {
		return word[:len(word)-2] + "ves"
	}
	if strings.HasSuffix(lower, "f") {
		return word[:len(word)-1] + "ves"
	}

	// Default: add -s
	return word + "s"
}
