package scaffold

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// adoptNextTSConfig takes back a tsconfig.json that Next.js rewrote.
//
// next build and next dev rewrite tsconfig.json when it lacks what they want
// (Next 16 sets "jsx" to "react-jsx" and includes .next/dev/types) and
// reformat the whole file while they are at it. After the first build the file
// no longer matched what Grit had recorded, so every upgrade reported it as
// edited by you and never updated it. The templates carry those values now, so
// Next leaves a new project's file alone. A file Next already rewrote is
// adopted when it means exactly what the template means, key for key, whatever
// its formatting; one with comments or any real difference stays yours.
//
// Called before an app's files are written, with the templates about to be
// written, while upgrade's manifest recording is active.
func adoptNextTSConfig(appRoot string, templates map[string]string) int {
	path := filepath.Join(appRoot, "tsconfig.json")
	template, ok := templates[path]
	if !ok || !fileExists(path) || manifest.IsUnchanged(path) {
		return 0
	}
	current, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	var have, want interface{}
	if json.Unmarshal(current, &have) != nil || json.Unmarshal([]byte(template), &want) != nil {
		return 0
	}
	if !reflect.DeepEqual(have, want) {
		return 0
	}
	manifest.Note(path, string(current))
	guardAdopt(path, string(current))
	return 1
}
