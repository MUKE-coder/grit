#import "template.typ": *

#show: book.with(
  title: "Go for JavaScript Developers",
  subtitle: "From React and Next.js to Go, and building real applications with Grit",
  author: "Muke JohnBaptist",
  version: "Written against Grit v3.403.0 and Go 1.25",
)

// Chapter 0: the repository and the leaderboard, before any Go.
#include "chapters/00-before-you-start.typ"

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
#include "chapters/24-config.typ"
#include "chapters/25-shipping.typ"

// Part IV: Grit, now that the reader can read what it generates.
#include "chapters/26-the-problems.typ"
#include "chapters/27-grit-new.typ"
#include "chapters/28-architecture.typ"
#include "chapters/29-the-resource.typ"
#include "chapters/30-field-types.typ"
#include "chapters/31-migrations.typ"
#include "chapters/32-services.typ"
#include "chapters/33-handlers.typ"
#include "chapters/34-admin-panel.typ"
#include "chapters/35-auth.typ"
#include "chapters/36-shared-package.typ"

// Part V: building with Grit.
#include "chapters/37-task-tracker.typ"
#include "chapters/38-relationships.typ"
#include "chapters/39-publication.typ"
#include "chapters/40-money-and-time.typ"
#include "chapters/41-invoicing.typ"
#include "chapters/42-uploads.typ"
#include "chapters/43-jobs-and-mail.typ"
#include "chapters/44-events-realtime-webhooks.typ"
#include "chapters/45-observability.typ"
#include "chapters/46-testing-grit.typ"
#include "chapters/47-upgrade-and-deploy.typ"

// Part VI: builds from first principles. Written out of order because the
// chapter format needed proving; the rest of Part IV slots in above.
#include "chapters/48-crud-api.typ"
#include "chapters/49-cli-framework.typ"
#include "chapters/50-json-parser.typ"
#include "chapters/51-kv-store.typ"
#include "chapters/52-btree-index.typ"
#include "chapters/53-redis-resp.typ"
#include "chapters/54-http-from-tcp.typ"
#include "chapters/55-build-git.typ"
#include "chapters/56-reverse-proxy.typ"
#include "chapters/57-job-queue.typ"
#include "chapters/58-test-runner.typ"

// Back matter: unnumbered, per the backmatter state in template.typ.
#include "chapters/90-phrasebook.typ"
#include "chapters/91-afterword.typ"
