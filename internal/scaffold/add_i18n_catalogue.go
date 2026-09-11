package scaffold

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Catalogues are the one kind of i18n file a project is meant to edit, so grit
// add i18n never overwrites one. That also meant a project scaffolded before a
// key existed never got it: the admin asked for nav.systemHub, found nothing,
// and showed English in the middle of a French page. mergeCatalogue adds the
// keys a catalogue lacks and changes nothing else: every value already there,
// translated or not, stays, in the order it was in.

// isCatalogue reports a translation file: messages/<locale>.json in a Next.js
// app, or internal/i18n/locales/<locale>.json in the API.
func isCatalogue(path string) bool {
	dir := filepath.Base(filepath.Dir(path))
	return filepath.Ext(path) == ".json" && (dir == "messages" || dir == "locales")
}

// catalogueEntry is one key of a JSON object, kept in its original order.
type catalogueEntry struct {
	key   string
	raw   json.RawMessage // a string, number or array, byte for byte
	child []catalogueEntry
	isObj bool
}

// mergeCatalogue adds to the catalogue at path every key in template it lacks.
// Returns whether it changed the file. A catalogue that is not valid JSON is
// left alone, with a warning, rather than failing the command.
func mergeCatalogue(path, template string) (bool, error) {
	current, err := os.ReadFile(path)
	if err != nil {
		return false, nil
	}
	have, err := parseCatalogue(current)
	if err != nil {
		fmt.Printf("  ⚠ %s is not valid JSON, so its new keys were not added: %v\n", path, err)
		return false, nil
	}
	want, err := parseCatalogue([]byte(template))
	if err != nil {
		return false, fmt.Errorf("the %s template is not valid JSON: %w", filepath.Base(path), err)
	}
	merged, changed := mergeCatalogueEntries(have, want)
	if !changed {
		return false, nil
	}
	var b bytes.Buffer
	writeCatalogue(&b, merged, "")
	b.WriteString("\n")
	if err := writeEdit(path, b.String()); err != nil {
		return false, fmt.Errorf("writing %s: %w", path, err)
	}
	return true, nil
}

func parseCatalogue(data []byte) ([]catalogueEntry, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil, fmt.Errorf("not a JSON object")
	}
	var out []catalogueEntry
	for dec.More() {
		keyTok, err := dec.Token()
		if err != nil {
			return nil, err
		}
		key, ok := keyTok.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected %v where a key belongs", keyTok)
		}
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return nil, err
		}
		entry := catalogueEntry{key: key}
		if trimmed := bytes.TrimSpace(raw); len(trimmed) > 0 && trimmed[0] == '{' {
			child, err := parseCatalogue(trimmed)
			if err != nil {
				return nil, err
			}
			entry.child, entry.isObj = child, true
		} else {
			entry.raw = raw
		}
		out = append(out, entry)
	}
	return out, nil
}

// mergeCatalogueEntries appends each key of want that have lacks, recursing
// into objects both have. Existing values are never replaced.
func mergeCatalogueEntries(have, want []catalogueEntry) ([]catalogueEntry, bool) {
	changed := false
	for _, w := range want {
		i := -1
		for j := range have {
			if have[j].key == w.key {
				i = j
				break
			}
		}
		if i < 0 {
			have = append(have, w)
			changed = true
			continue
		}
		if have[i].isObj && w.isObj {
			if merged, c := mergeCatalogueEntries(have[i].child, w.child); c {
				have[i].child = merged
				changed = true
			}
		}
	}
	return have, changed
}

// writeCatalogue writes entries with two-space indentation, the shape the
// templates use.
func writeCatalogue(b *bytes.Buffer, entries []catalogueEntry, indent string) {
	b.WriteString("{\n")
	for i, e := range entries {
		b.WriteString(indent + "  ")
		b.Write(catalogueKey(e.key))
		b.WriteString(": ")
		if e.isObj {
			writeCatalogue(b, e.child, indent+"  ")
		} else {
			var compact bytes.Buffer
			if json.Compact(&compact, e.raw) == nil {
				b.Write(compact.Bytes())
			} else {
				b.Write(bytes.TrimSpace(e.raw))
			}
		}
		if i < len(entries)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}
	b.WriteString(indent + "}")
}

// catalogueKey quotes a key without HTML-escaping it.
func catalogueKey(key string) []byte {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(key)
	return bytes.TrimRight(b.Bytes(), "\n")
}
