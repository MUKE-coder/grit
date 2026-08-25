// Package uploadpkg holds the client-side upload and image optimisation
// library.
//
// This directory is two things at once, deliberately. It is a real npm package
// that any React app can install, and it is the source the scaffolder embeds
// into every generated project. Keeping one copy is the point: the alternative
// is a published package and a Go string template that agree only until
// somebody edits one of them.
//
// Go ignores the npm files and npm ignores this one, via the "files" list in
// package.json.
package uploadpkg

import "embed"

// Sources are the TypeScript files written into a generated project.
//
// package.json is not embedded: the published package and a workspace copy
// need different names and entry points, so the scaffolder writes its own.
//
// The .tsx pattern is separate because src/*.ts does not match it. Missing it
// would drop the Dropzone from every generated project silently, since a file
// that is not embedded simply is not written.
//
//go:embed src/*.ts src/*.tsx styles.css README.md
var Sources embed.FS
