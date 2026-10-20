#import "../template.typ": *

= The resource, and why it is the unit of generation

This is the idea the framework rests on. If you disagree with it, the rest of
part four is wrong, so it gets a chapter of its own.

== The claim

*A resource is the only thing that spans the whole stack, so it is the thing
worth generating.* Not a file. Not a layer. The invoice, from the database
column to the form field.

Consider what "add invoices to this application" actually requires:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Database], [a table, columns, indexes, a foreign key],
  [Go], [a model, a service with the rules, a handler, routes],
  [Contract], [a validation schema, a TypeScript type],
  [Frontend], [a data hook, a list page, a form, a detail view],
)

Thirteen or so files. A scaffolder that generates *a model* has automated the
easy part, because the model was never the expensive bit. The expense is that
one concept, "invoice", has to be spelled correctly in thirteen places, and
kept correct as it changes.

#bridge[
  This is D3 from chapter twenty six, and you have felt it. The Go struct, the
  Zod schema, the TypeScript type and the mobile client are four descriptions
  of one thing, kept in step by memory. Nothing fails to compile when they
  drift. You find out in the browser, from a user, as a blank column or a
  wrong total.
]

== One command

```bash
grit generate resource Invoice \
  --fields "number:string, total:money, due:date, client:belongs_to"
```

Four field types. One of them is a relationship. That is the entire input.

What it writes, by name:

```
apps/api/internal/models/invoice.go
apps/api/internal/services/invoice.go
apps/api/internal/services/invoice_import.go
apps/api/internal/handlers/invoice.go
apps/api/internal/handlers/invoice_import.go
apps/api/internal/handlers/invoice_test.go
apps/api/internal/policies/invoice.go
apps/api/internal/factory/invoice.go
packages/shared/schemas/invoice.ts
packages/shared/types/invoice.ts
apps/admin/app/_dashboard/resources/invoice/index.tsx
```

Plus the routes registered, the model added to the migration list, and the
resource added to the admin's navigation.

#gotcha[
  Note `index.tsx` in the resource folder. The generator writes it once and
  then *never regenerates it*, because that is where your custom cells and
  forms go. The distinction between "regenerated" and "written once" runs
  through the whole system, and chapter thirty one is about where that line
  sits.
]

== Why one description is the point

The thirteen files are not the achievement. The achievement is that they come
from *one description*, so they cannot disagree.

Add a field later:

```bash
grit generate field Invoice paid_at:datetime
```

That reaches the model, the migration, the Zod schema, the TypeScript type, the
admin table and the admin form in one step. There is no version of this where
the backend has `paid_at` and the frontend does not, because there is no second
place to update.

That is the real answer to D3. Not "less typing": *no drift*.

== A field type is a decision made once

Thirty one field types ship. The important ones are not the obvious ones.

```
string  text  richtext  int  uint  float  bool  toggle  select  radio
check  datetime  date  money  slug  belongs_to  one_to_one  many_to_many
string_array  file  files  email  url  domain  tel  country  color
percent  rating  time  json
```

Take `money`. Writing `total:money` decides a column type, a Go type, a Zod
rule, an input mask and how the number is formatted in a table cell. Say
`float` instead and you have shipped a rounding bug that will surface in a
total months later, found by an accountant.

```go
// What money generates: an integer, in minor units.
Total money.Amount    // stored as cents, never a float64
```

This is the lesson from chapter nineteen, where you implemented `MarshalJSON`
for exactly this. The difference is that here it is a *type* rather than
advice, so the correct decision is the default rather than something you must
remember.

Take `belongs_to`. One word, and you get:

- a `ClientID` foreign key column, with a constraint
- a `Client` field on the struct, preloaded on read
- a nested type in the generated TypeScript
- a searchable dropdown in the admin form that shows the client's *name*
  rather than a UUID

#bridge[
  In a hand-written stack each of those is a separate decision someone makes
  slightly differently. The dropdown showing a UUID instead of a name is the
  canonical example: it is nobody's bug, everybody has shipped it, and it
  happens because the frontend had no way to know which column is the human
  label. Here the generator knows, because `belongs_to` told it.
]

== It is still your code afterwards

Open the generated service. It is a struct with a `*gorm.DB` and methods, with
comments explaining *why* rather than what:

```go
// ContactService owns every database read and write for contacts.
//
// The handler reads the request, calls one of these and writes the
// answer. A job, a command or a test calls the same methods with no
// request at all, so a rule enforced here is enforced everywhere.
type ContactService struct {
    DB *gorm.DB
}
```

One dependency, and it is GORM. You can edit any of it. The generated code is
the starting point, not a sealed artifact.

#gotcha[
  Re-running `grit generate resource` on an existing resource does *not*
  silently overwrite your edits. Every generated file is hashed into a manifest
  when written, and a file whose hash has moved is reported rather than
  replaced.

  This behaviour exists because the opposite happened: an early version
  rewrote the model, service and handler on a re-run, printing a green tick for
  each, and everything hand-written in them was gone. No prompt, no backup, and
  the documentation had promised that the service was where your business logic
  goes. A newcomer lost a `BeforeSave` rule and a service method in their first
  week.

  The lesson generalises past this framework: a generator that overwrites
  without asking is a generator you cannot trust with a file you have edited.
]

== What it does not generate

Worth saying plainly, because the demo makes it look like magic.

It generates the *shape*. It does not generate the rule that an invoice cannot
be billed twice, the report the client asks for every month, or the pricing
logic that is the reason they hired you. That is still yours, it is still the
hard part, and it is still where the value is.

What has gone is the two weeks before you got to any of it.

#project(title: "Project 26: generate a resource and read every file")[
  In the project from chapter twenty seven:

  ```bash
  grit generate resource Book \
    --fields "title:string, isbn:string, price:money, published:date"
  grit migrate
  grit seed Book --count 500
  grit start
  ```

  Then open every file it named, in this order, and answer one question about
  each.

  *`models/book.go`.* Find the `money` field. What Go type is it, and what does
  the GORM tag say the column is? Why is that not a `float64`?

  *`services/book.go`.* Find `paginate.Config`. Which columns can a client sort
  by, and which cannot? Now try `?sort=price` and `?sort=secret` against the
  running API.

  *`handlers/book.go`.* Count the lines that are about HTTP. Find the one call
  into the service per endpoint.

  *`packages/shared/schemas/book.ts`.* Find the Zod rule for `price` and
  compare it to the Go validation. They came from the same word in your
  command.

  *`packages/shared/types/book.ts`.* Compare the TypeScript type to the Go
  struct field by field.

  *The admin panel at `:3001`.* Find the Books page. Sort by title, search,
  open the form. Notice the price input: it is a money input, not a number
  input, because of one word in your command.

  Now add a field and watch the whole chain move:

  ```bash
  grit generate field Book pages:int
  grit migrate
  ```

  Reload the admin. The column is in the table and the field is in the form,
  and you edited nothing. Check `types/book.ts`: it is there too.

  Finally, prove the manifest. Add a method to `services/book.go` by hand, then
  re-run the original `grit generate resource Book` command. It should refuse
  to overwrite rather than quietly eating your method, and tell you why.
]

#tryit[
  Generate a resource with a deliberately wrong field type, `price:float`
  instead of `price:money`, and look at what changes in the admin form and the
  Zod schema. Then compute `0.1 + 0.2` in the Go REPL of your choice, or just
  `fmt.Println(0.1 + 0.2)`, and decide for yourself whether a `money` type is
  pedantry or the single most valuable field type in the list.
]

#assignment(name: "Resources - Assignment One")[
  Generate the domain you described in chapter twenty six.

  *One.* Three resources, with their fields chosen rather than defaulted:

  ```bash
  grit generate resource Site --fields "name:string, code:slug"
  ```

  *Two.* Before running the command, write down what you expect it to create.
  Then run it with the diff open and check. The list you did not predict is the
  list worth reading.

  *Three.* Edit one generated service by hand, re-run the generator for that
  resource, and show what happened to your edit. Record it in `NOTES.md`,
  because trusting a generator you have not tested this on is how edits get
  lost.

  *Four.* Commit the whole diff as one commit, so the history shows what one
  command produced.]
