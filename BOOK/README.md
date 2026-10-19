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
| `#compare(js, go)` | The same task in both languages, stacked |

Two rules learned the hard way:

- **Code must fit 62 characters.** The text block is 5.35in and code is Consolas
  8.8pt. Anything wider wraps mid-token and reads as nonsense.
- **`#compare` stacks, it does not sit side by side.** Two columns was the first
  attempt; half the measure is about thirty characters and indented Go does not
  survive it.

## House style

- **No em dashes.** Project-wide rule. Use a colon, a comma or a full stop.
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
11. Errors are values, not exceptions
12. `defer`, `panic` and `recover`
13. Goroutines, against the event loop
14. Channels, and the select statement
15. `sync`: mutexes, wait groups, and when to use them instead
16. Context: cancellation, deadlines, and request scope
17. Generics, briefly and honestly
18. Testing: the standard library runner, table tests, and coverage
19. JSON, struct tags, and the capital letter rule that bites

### Part III: building real things in plain Go

20. Reading files, streams and the `io` interfaces
21. Project: a JSON formatting CLI
22. `net/http`: a server with no framework
23. Routing, middleware and handlers
24. Project: a JSON API, standard library only
25. Talking to a database with `database/sql`
26. Project: a URL shortener with SQLite
27. Configuration, environment and secrets
28. Building, cross compiling and shipping one binary

### Part IV: Grit fundamentals

29. The twenty six problems, and what a framework is for
30. `grit new`: what lands, and the shape of a project
31. The architecture: a modular monolith, layer by layer
32. The resource, and why it is the unit of generation
33. Field types: all thirty one, and what each one decides
34. Models, migrations and seeding at scale
35. Services: where the rules live, and why not the handler
36. Handlers and the HTTP contract
37. The generated admin panel
38. Authentication, roles and per-record ownership
39. The shared package: one definition, four outputs

### Part V: building with Grit

40. Project: a task tracker, end to end
41. Relationships: `belongs_to`, `many_to_many` and the forms they generate
42. Project: a blog with a working admin
43. Money, dates and the types that prevent bugs
44. Project: an invoicing application
45. File uploads and the image pipeline
46. Background jobs, cron and transactional email
47. Realtime, the event bus and webhooks
48. Observability: Studio, Pulse and Sentinel
49. Testing a Grit application
50. `grit upgrade`, deployment, and where to go next
