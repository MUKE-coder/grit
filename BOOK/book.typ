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
