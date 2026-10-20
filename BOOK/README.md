# Go for JavaScript Developers

A textbook that teaches Go to someone who already ships React and Next.js, then
teaches Grit and builds real applications with it. Target: **500 pages**.

## Building

```bash
python build.py            # compile book.pdf, print progress against 500 pages
python build.py --pages 8  # also render the first 8 pages as PNGs to look at
```

Requires [Typst](https://typst.app) on `PATH` (`typst --version`). No LaTeX.

## Why Typst and not Markdown

The book leans on four repeated teaching devices and a JavaScript/Go comparison
block. Those need real layout: a coloured rule, a label, controlled breaking
across pages. Markdown plus a converter gives none of that without fighting it,
and the devices are the pedagogy rather than decoration.

## The shape of a chapter

Every chapter uses the same instruments, defined in `template.typ`. Use them;
do not style anything inline, or the book stops looking like one book.

| Device | Use it for |
|---|---|
| `#bridge[...]` | The core move: name the JavaScript thing first, then where it breaks |
| `#gotcha[...]` | A trap that will cost the reader an hour |
| `#project[...]` | A small thing to build, complete and runnable |
| `#tryit[...]` | A short exercise, left to the reader |
| `#assignment(name: "...")[...]` | The chapter's assignment, with the exact commit message |
| `#compare(js, go)` | The same task in both languages, stacked |

### Assignments

Every chapter ends with exactly one `#assignment`, and chapter 0 explains the
system to the reader: one public GitHub repository for the whole book, one
commit per assignment, scored on the Gopher League leaderboard in the Committed
app.

The `name:` argument is the commit message, and it is not decoration. The
leaderboard matches on it, so it must follow `<Topic> - Assignment <Number>`
exactly, in words rather than digits, and it must never change once a chapter
has shipped. Chapter 3 carries two, which is why the numbering exists.

The website's Go guide at `/docs/prerequisites/golang` has 46 browser exercises
over the same ground. Where one lines up with a chapter, say so at the end of
that chapter's assignment. They are for practising a concept; an assignment is
for building something and pushing it, and only the second one earns points.

Two rules learned the hard way:

- **Code must fit 75 columns**, enforced by `build.py`. Measured, not
  calculated: the font-metric estimate said 80 and a 72 column line wrapped on
  the page. Body code is Consolas 8.2pt and callout code 7.6pt, scaled so one
  budget holds everywhere.
- **`#compare` stacks, it does not sit side by side.** Two columns was the first
  attempt; half the measure is about thirty characters and indented Go does not
  survive it.

## House style

- **No em dashes.** Project-wide rule. Use a colon, a comma or a full stop.
- **Typst bold is `*one star*`, not Markdown's `**two**`.** `build.py` checks.
- Say what a thing costs, not only what it does. A chapter that only lists
  advantages is marketing.
- Every project must be finished and runnable, never a toy that stops before
  the interesting part.
- Grit facts must come from the codebase, not from memory. The framework half
  of this book is checkable, and a wrong claim there is worse than a vague one.

## Progress

`build.py` prints pages written, pages remaining and the current density. At the
time of writing it runs about **8.5 pages per chapter**, so the plan below is
sized at roughly 50 chapters.

## Plan

### Before Part I

0. Before you start: a repository, and a leaderboard **(written)**

Numbered 0 so Part I still opens at chapter 1, which is what every
cross-reference in the book says. `template.typ` prints the chapter number from
its own counter and subtracts one.

### Part I: Go for the JavaScript developer

1. Why Go, when you already ship JavaScript **(written)**
2. Installing Go and your first program **(written)**
3. Variables, types and the compiler that argues with you
4. Functions, multiple returns, and the end of try/catch
5. Structs and methods, against objects and classes
6. Slices: the array that is not an array
7. Maps, and why iteration order is random on purpose
8. Pointers, without the fear
9. Packages, modules and imports, against npm and ESM

### Part II: the ideas that make Go different

10. Interfaces, and why Go's point the other way
11. Designing with errors: what a package should return
12. `defer`, `panic` and `recover`
13. Goroutines, against the event loop
14. Channels, and the select statement
15. `sync`: mutexes, wait groups, and when to use them instead
16. Context: cancellation, deadlines, and request scope
17. Generics, briefly and honestly
18. Testing: the standard library runner, table tests, and coverage
19. JSON, struct tags, and the capital letter rule that bites

### Part III: building real things in plain Go

20. Files, streams and the `io` interfaces **(written)**
21. `net/http`: a server with no framework **(written)**
22. Talking to a database with `database/sql` **(written)**
23. Project: a URL shortener **(written)**
24. Configuration, environment and secrets **(written)**
25. Building, cross compiling and shipping one binary **(written)**

The JSON CLI that was planned here is now the project in chapter 19, and
routing and middleware are covered inside chapter 21 rather than split out.

### Part IV: Grit fundamentals

26. The twenty six problems, and what a framework is for **(written)**
27. `grit new`: what lands, and the shape of a project **(written)**
28. The architecture: a modular monolith, layer by layer **(written)**
29. The resource, and why it is the unit of generation **(written)**
30. Field types: all thirty one, and what each one decides **(written)**
31. Models, migrations and seeding at scale **(written)**
32. Services: where the rules live, and why not the handler **(written)**
33. Handlers and the HTTP contract **(written)**
34. The generated admin panel **(written)**
35. Authentication, roles and per-record ownership **(written)**
36. The shared package: one definition, four outputs **(written)**

### Part V: building with Grit

37. Project: a task tracker, end to end **(written)**
38. Relationships, and the forms they generate **(written)**
39. Project: a publication, from draft to public page **(written)**
40. Money, dates and the types that prevent bugs **(written)**
41. Project: an invoicing application **(written)**
42. File uploads and the image pipeline **(written)**
43. Background jobs, cron and transactional email **(written)**
44. The event bus, realtime, and webhooks **(written)**
45. Observability: Studio, Pulse and Sentinel **(written)**
46. Testing a Grit application **(written)**
47. Upgrading, deploying, and handing it over **(written)**

### Part VI: build it from scratch

Eleven builds in the spirit of the build-your-own-x list. They deepen Go rather
than teach Grit, and they depend only on Parts I to III, so a reader who wants
them early can jump here straight after chapter 25.

**Written from first principles.** Two third-party books on the same theme sit
in the author's reference folder and are deliberately not drawn from. The idea
of building a database from scratch is nobody's property; the prose and code of
somebody else's book is. Nothing in this part is adapted from them.

48. Build your own CRUD API from scratch **(written)**
49. Build your own CLI framework: flags, subcommands and help **(written)**
50. Build your own JSON parser: a lexer and a recursive descent parser **(written)**
51. Build your own key-value store: a write-ahead log and crash recovery **(written)**
52. Build your own on-disk index: pages and a B-tree **(written)**
53. Build your own Redis: RESP over raw TCP **(written)**
54. Build your own HTTP server, from a TCP socket **(written)**
55. Build your own Git **(written)**
56. Build your own reverse proxy and load balancer **(written)**
57. Build your own job queue **(written)**
58. Build your own test runner **(written)**

Chapter 48 comes first because it needs only Parts I to III and it is the one
build every reader will recognise from work. It also sets up Part IV: having
written six hundred lines for one resource by hand, `grit generate resource` is
a measured claim rather than a demo.

File names carry the plan number, so `51-kv-store.typ` is chapter 51. A chapter
written out of order keeps its planned number rather than its writing order.

### Back matter

Unnumbered, via the `backmatter` state in `template.typ`: the chapter counter
still steps and only the number is hidden, so nothing later depends on where it
was set.

- Appendix: the JavaScript to Go phrasebook **(written)**
- Afterword: what to do on Monday **(written)**

## Status

**503 pages, 61 chapters. The 500-page target is met and every chapter in the
plan is written.** What remains is revision rather than drafting: tightening
Part I now that the later parts exist, checking every Grit claim against the
current codebase, and a pass for repetition across chapters that were written
weeks apart.
