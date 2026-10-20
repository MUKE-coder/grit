#import "../template.typ": *

= Relationships, and the forms they generate

A model on its own is a spreadsheet. The moment a task belongs to a project and
carries tags, you have the thing that makes applications hard: two tables that
must agree, a form that has to offer the right options, a query that must not
run once per row, and a JSON shape the frontend can rely on.

Grit has four relationship shapes. This chapter is what each one writes, and
where each one bites.

== The four shapes, and the question each answers

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`belongs_to`], [This row points at one other row. A task has a project],
  [`one_to_one`], [The same, but only one row may point at each target. A user
    has one profile],
  [`many_to_many`], [Both sides have many. A task has tags, a tag has tasks],
  [`--items`], [A child that exists only inside its parent. An invoice has
    line items],
)

The fourth is not a field type, and the difference is worth naming early: an
`--items` child has no life of its own. You do not browse line items; you edit
an invoice. Everything else is a relation between two things that each exist.

== belongs_to, and everything it writes

```bash
grit generate resource Task \
  --fields "title:string, project:belongs_to:Project"
```

The related model is the third part, and it is optional when the field name
already names it: `project:belongs_to` finds `Project`. Spell it out when the
field is called something else, which is the usual case for people:
`author:belongs_to:User`.

In the model you get two things, not one:

```go
type Task struct {
    ID        uint     `json:"id" gorm:"primaryKey"`
    Title     string   `json:"title"`
    ProjectID uint     `json:"project_id" gorm:"index"`
    Project   *Project `gorm:"foreignKey:ProjectID"
                        json:"project,omitempty"`
}
```

`ProjectID` is the column the database knows about, indexed because every
listing filters on it. `Project` is the association GORM fills in when asked,
and it is a *pointer with `omitempty`*.

#gotcha[
  That pointer is not a style choice, and the version before it shipped a
  convincing bug.

  As a value, a relation that was never preloaded marshalled as a complete
  zero-value object: `"project": {"id": 0, "name": "", "colour": ""}`. To a
  client that reads as a real project whose every field happens to be blank,
  and the generated TypeScript agreed, because it said the object was always
  there.

  As a pointer with `omitempty`, an absent relation is absent. The frontend can
  tell "not loaded" from "loaded and empty", which are different facts and used
  to be the same one.
]

The service preloads it, so a list of tasks arrives with its projects:

```go
db.Preload("Project").Find(&tasks)
```

That is chapter twenty two's N+1 solved by default: one query for tasks, one
for the projects they reference, rather than one per row. You get it without
asking, which is also why it is worth knowing it is there when a list of
twelve columns suddenly reads from five tables.

== The admin dropdown, and the label it shows

The form field for a `belongs_to` is a searchable dropdown, and the thing that
makes it usable is what it displays. It is not the ID, and it is not a field
called whatever came first. The generator looks for a label column on the
related model, in this order:

```
name, title, label, subject, number, code, reference, email
```

The first one that exists wins. A `Project` shows its name, an `Invoice` shows
its number, and a `User` with none of those falls back to its email, which is
the reason the fallback list ends there rather than at `code`.

#bridge[
  Every admin panel you have built by hand has had this bug: a select showing
  `[object Object]`, or a list of IDs, or a dropdown that loads six thousand
  rows to find one.

  The generated one searches server side and shows a human label. It is a small
  thing that is tedious to get right once per relation, and there are usually
  twenty relations.
]

== one_to_one is belongs_to plus one index

```bash
grit generate resource Profile \
  --fields "bio:text, user:one_to_one:User"
```

The generated code differs from `belongs_to` in exactly one place: the foreign
key gets a unique index rather than an ordinary one.

That single line is what makes the name true. Without it, "one profile per
user" is a comment, because the database would happily accept a second row
pointing at the same user, and the second one usually arrives from a retried
request rather than from a user doing something strange.

== A relation pointing at itself is a tree

```bash
grit generate field Category parent:belongs_to:Category
```

Electronics above Cameras above Lenses, in one table. Two details differ from
an ordinary relation, and both were found the hard way.

*The association must be a pointer.* Go rejects a struct that contains itself
by value: `invalid recursive type`. A generated project with a self-reference
did not compile at all before this was handled.

*The foreign key must be nullable.* GORM creates a real foreign key constraint,
and a root category has no parent, so the column has to hold something the
constraint accepts. SQL has exactly one such value, and it is `NULL`.

#gotcha[
  An earlier version used the empty string for "no parent", matching what every
  other `belongs_to` field does. It passed every test and failed on the first
  real project:

  ```
  ERROR: insert or update on table "categories" violates
  foreign key constraint "fk_categories_children"
  ```

  The tests ran on SQLite, which does not enforce foreign keys unless you ask
  it to. Postgres does. The generated tests now switch foreign keys on for
  exactly this reason, and the general lesson is larger than trees: *a test
  database that is more permissive than production is a test database that
  hides this class of bug entirely*.
]

For a real hierarchy, the flag does more than the field:

```bash
grit generate resource Category --fields "name:string" --tree
```

That gives the self-referential parent plus a materialized path, a depth, a
sibling order, the tree queries and a move endpoint. The difference matters
when you need "everything under Electronics" as one query rather than as
recursion in application code.

== many_to_many, and the table you did not name

```bash
grit generate resource Task \
  --fields "title:string, tags:many_to_many:Tag"
```

The related model is required here. There is nothing to infer from, because
`tags` is a plural of the field, not of the model.

The model gets a slice and no foreign key, because the keys live in a third
table:

```go
Tags []Tag `gorm:"many2many:tasks_tags" json:"tags"`
```

The join table is named from the owning model and the field: `tasks_tags`.
GORM creates and maintains it, and you do not get a model for it.

#gotcha[
  The moment you need a column *on the relationship itself*, a quantity, a
  role, a date added, `many_to_many` is the wrong shape and no flag fixes it.

  "A user belongs to many projects" is many-to-many. "A user belongs to many
  projects *as owner or as member*" is a `Membership` model with two
  `belongs_to` fields and a `role`, generated like any other resource.

  Reaching for the join table as a place to put data is the single most common
  modelling mistake in this chapter, and it is expensive because the migration
  that fixes it has to move rows.
]

In the admin this renders as a multi-select picker, and in JSON it is an array,
which means a create or update sends IDs:

```json
{ "title": "Fix the leak", "tag_ids": [3, 7] }
```

== Line items: a child inside its parent

```bash
grit generate resource Invoice \
  --fields "number:string, client:belongs_to:Client" \
  --items "InvoiceItem:description:string,qty:int,unit:money"
```

The child is generated fully, and the parent's form gets a line-items table:
add a row, remove a row, totals that update. That is a component people build
by hand for a week.

The part worth understanding is the update, because it is not what you would
write first:

```go
// Inside the same transaction as the parent update
tx.Where("invoice_id = ?", item.ID).
    Delete(&models.InvoiceItem{})
tx.Create(&rows)
```

The children are deleted and recreated rather than diffed.

#gotcha[
  Delete-and-recreate is the right call for line items and the wrong call for
  anything with its own identity. New IDs on every save means anything pointing
  at a line item now points at nothing.

  So the test for `--items` is: *would anybody ever link to this row?* An
  invoice line, no. A comment on a task, yes, and a comment is therefore a
  resource with a `belongs_to`, not an item.

  Note also that the delete and the insert are inside one transaction. Without
  that, a failure between them leaves an invoice with no lines, which is a
  document that silently became worth nothing.
]

== What relationships cost on a list page

Three consequences that surprise people, all of them from chapter thirty two's
whitelists:

*You cannot sort by a relation.* The sortable list contains columns, and
`project.name` is not a column of `tasks`. Sorting by a related field needs a
join the generated query does not make. Sort by `project_id` if order matters,
or denormalise the name onto the row if it matters a lot.

*You filter by the key, not the label.* `?project_id=4` works. `?project=Apollo`
does not, unless you add it.

*Search does not cross the relation.* Searching tasks searches the task's own
searchable columns. Finding "every task whose project is named Apollo" is a
query you write.

#bridge[
  These are not framework limitations so much as the honest shape of SQL
  showing through. An ORM that silently joined four tables to satisfy a sort
  parameter would be faster to demonstrate and slower in production, and you
  would find out at the worst moment.

  The framework's choice is to make the fast thing automatic and the expensive
  thing explicit.
]

== Choosing well

Two questions settle most modelling arguments.

*Is the other side a real thing, or a word?* A status is a word: `select`, not
a relation. A supplier is a thing: a relation. The giveaway is whether anybody
would ever want to edit it, list it, or attach anything else to it.

*Does the link itself carry data?* If yes, it is a model, not a
`many_to_many`.

#project(title: "Project 35: model something genuinely awkward")[
  Build a small school timetable, which has every shape in this chapter and one
  real trap.

  ```bash
  grit new timetable
  cd timetable
  ```

  *One: the obvious ones.*

  ```bash
  grit generate resource Teacher --fields "name:string, email:email"
  grit generate resource Room --fields "code:string, seats:int"
  grit generate resource Course \
    --fields "title:string, teacher:belongs_to:Teacher"
  ```

  Open the Course form and confirm the teacher dropdown shows names, not IDs.

  *Two: the trap.* Model "a student attends many courses" as a
  `many_to_many`. Then try to answer: when did they enrol, and did they pass?

  You cannot, so replace it with an `Enrolment` resource that belongs to both
  and carries `enrolled_on:date` and `grade:select`. Write the two commands and
  the migration, and note in `NOTES.md` what you would have had to do if the
  first version had been live for a year.

  *Three: the tree.* Give `Course` a `prerequisite:belongs_to:Course`. Create a
  root with no prerequisite and confirm the column stores `NULL`. Then run the
  same thing against Postgres rather than SQLite:

  ```bash
  docker compose up -d db
  ```

  *Four: line items.* Generate a `Timetable` with
  `--items "Slot:day:select:mon=Mon|tue=Tue,period:int,room:belongs_to:Room"`.
  Add three slots, save, and watch the IDs change. Then write one sentence on
  why that is acceptable here and would not be for enrolments.

  *Five: the cost.* Load the Course list with 500 rows and look at the API log.
  Count the queries. Then remove the `Preload` from the service, reload, and
  count again. That difference is the N+1, measured.
]

#tryit[
  Take the data model of something you have already built, and write down each
  relationship as one of the four shapes in this chapter.

  The ones you cannot classify are the interesting ones. In my experience they
  are almost always a join table that grew a column, which is the modelling
  mistake this chapter exists to name.
]

#assignment(name: "Relationships - Assignment One")[
  In `part5-building/38-relationships`, build a small library system with every
  shape, and prove the behaviour rather than describing it.

  *One.* `Author`, `Book` with `author:belongs_to:Author`, `Member`, and a
  `Loan` that belongs to both a book and a member and carries `due:date`.
  Explain in one line why `Loan` is a model and not a `many_to_many`.

  *Two.* Add `categories:many_to_many:Category` to `Book`. Find the join table
  in `/studio`, record its name, and paste one row into `NOTES.md`.

  *Three.* A self-referential `Category` with a parent. Create a root, confirm
  the stored value is `NULL` and not an empty string, and say what would break
  on Postgres if it were.

  *Four.* The JSON evidence. Fetch one book with its author preloaded, and one
  endpoint where it is not. Paste both responses: one has an `author` object,
  the other has no `author` key at all. That difference is the pointer and the
  `omitempty`, and it is the one thing in this chapter a frontend developer
  will thank you for.

  *Five.* Try to sort the book list by author name through the API. Record what
  happens, then make it work, and write down what you had to do.
]
