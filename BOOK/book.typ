#import "template.typ": *

#show: book.with(
  title: "Go for JavaScript Developers",
  subtitle: "From React and Next.js to Go, and building real applications with Grit",
  author: "Muke JohnBaptist",
  version: "Written against Grit v3.403.0 and Go 1.25",
)

// Part I: the language, taught against what the reader already knows.
#include "chapters/01-why-go.typ"
#include "chapters/02-hello-go.typ"
#include "chapters/03-types.typ"
#include "chapters/04-functions-and-errors.typ"
#include "chapters/05-structs.typ"
#include "chapters/06-slices.typ"
#include "chapters/07-maps.typ"
#include "chapters/08-pointers.typ"
#include "chapters/09-packages.typ"

// Part II: the ideas that separate Go from what the reader knows.
#include "chapters/10-interfaces.typ"
#include "chapters/11-error-design.typ"
#include "chapters/12-defer-panic-recover.typ"
#include "chapters/13-goroutines.typ"
#include "chapters/14-channels.typ"
#include "chapters/15-sync.typ"
#include "chapters/16-context.typ"
#include "chapters/17-generics.typ"
#include "chapters/18-testing.typ"
#include "chapters/19-json.typ"

// Part III: building real things with the standard library.
#include "chapters/20-io.typ"
#include "chapters/21-http-server.typ"
#include "chapters/22-database.typ"
#include "chapters/23-url-shortener.typ"
