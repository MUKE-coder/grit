#import "../template.typ": *

= Project: a task tracker, end to end

Part IV explained the pieces. This chapter uses them, in order, to build
something a small team could actually run: projects, tasks with a status and an
assignee, comments, a dashboard, and a web page that is not the admin panel.

Everything here is a command you type and a file you open. Nothing is skipped
for space, and the parts that are genuinely tedious are called tedious rather
than hidden.

== What we are building

A tracker with three models and one rule that matters:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`Project`], [A name, a colour, an owner. Archived rather than deleted],
  [`Task`], [Belongs to a project, has a status, an assignee and a due date],
  [`Comment`], [Belongs to a task, written by a user],
)

The rule: *a task cannot be closed while it has open subtasks*, which is the
kind of thing that is trivial to state, easy to enforce in one place, and
impossible to enforce if it lives in a form.

== Step one: the project

```bash
grit new tracker
cd tracker
grit start
```

Three applications come up: the API on `:8080`, the web app on `:3000` and the
admin on `:3001`. Log in with the seeded admin account printed in the output,
and confirm the dashboard renders before generating anything. A broken starting
point is worth five minutes to rule out.

#gotcha[
  If the admin loads and the lists are empty while the API logs show nothing,
  the usual cause is that migrations have not run. `grit migrate` once, then
  reload.

  The second most common cause is a port already in use from an earlier
  project, which shows up as a login that posts and does nothing because
  `:8080` belongs to somebody else. `grit doctor` reports both.
]

== Step two: the resources

Generate the parent first, because the child refers to it:

```bash
grit generate resource Project \
  --fields "name:string, colour:color, archived:bool" \
  --owned-by user
```

Then the task, with the relationship and a real status:

```bash
grit generate resource Task --fields "\
title:string, \
notes:richtext, \
status:select:todo=To do|doing=In progress|done=Done, \
priority:select:low=Low|normal=Normal|high=High, \
due:date, \
project:belongs_to:Project, \
assignee:belongs_to:User"
```

And the comment:

```bash
grit generate resource Comment \
  --fields "body:text, task:belongs_to:Task"
```

Then once:

```bash
grit migrate
```

Reload the admin. Three new sections in the sidebar, each with a working list
and form. The task form has two dropdowns with labels rather than raw values, a
date picker, a rich text editor, and two searchable relationship pickers.

#bridge[
  Count what you did not do. No table component, no form library wiring, no
  `<select>` with hand-written options, no fetch hook, no TypeScript interface,
  no Zod schema, no migration file, no route registration.

  That is chapter twenty nine's claim arriving. The reason it is worth
  distrusting until you have done it is that most tools make the same claim and
  stop at the list page.
]

== Step three: the rule, in the only place it can live

Open `apps/api/internal/services/task.go` and add a method. This is
chapter thirty two, applied:

```go
var ErrHasOpenSubtasks = errors.New(
    "task has open subtasks")

func (s *TaskService) Close(ctx context.Context, id uint) (
    *models.Task, error) {

    task, err := s.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    var open int64
    if err := s.db(ctx).Model(&models.Task{}).
        Where("parent_id = ? AND status != ?", id, "done").
        Count(&open).Error; err != nil {
        return nil, fmt.Errorf("counting subtasks: %w", err)
    }
    if open > 0 {
        return nil, ErrHasOpenSubtasks
    }

    task.Status = "done"
    if err := s.db(ctx).Save(task).Error; err != nil {
        return nil, fmt.Errorf("closing task %d: %w", id, err)
    }
    return task, nil
}
```

The subtask relationship needs a column, which is one command rather than a
migration file:

```bash
grit generate field Task parent:belongs_to:Task
grit migrate
```

A self-referencing `belongs_to` is legal and useful: it is how a task gets
subtasks without a second model.

#gotcha[
  `s.db(ctx)` rather than `s.DB`. The first is scoped by the actor in the
  context, the second is not.

  On a `--owned-by user` resource that one difference decides whether the
  subtask count includes other people's tasks. It will look correct in your own
  testing, because you are the only user, which is exactly why this is the bug
  that reaches production.
]

== Step four: the endpoint

```go
func (h *TaskHandler) Close(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        respond.BadRequest(c, "invalid id")
        return
    }

    task, err := h.service().Close(h.ctx(c), uint(id))
    if err != nil {
        respond.Fail(c, err)
        return
    }
    respond.OK(c, task)
}
```

Register it next to the generated routes:

```go
tasks.POST("/:id/close", h.Task.Close)
```

And teach the boundary what the new error means, so the frontend gets a 409
rather than a 500:

```go
case errors.Is(err, services.ErrHasOpenSubtasks):
    respond.Conflict(c, err)
```

Check it is wired:

```bash
grit routes
```

That prints every registered route with its method, handler and middleware
group, which is the fastest way to catch a route registered on the wrong group
and therefore unauthenticated.

== Step five: data worth looking at

An empty application tells you nothing about whether it works:

```bash
grit seed Project --count 8
grit seed Task --count 400
grit seed Comment --count 1200
```

Four hundred tasks is where the list page stops being a demo. Sort by due date,
filter by status, search a title, and page through. Pagination is server side,
so the browser is holding twenty rows no matter how many exist.

#tryit[
  Seed a hundred thousand tasks and try the same screen. It stays quick,
  because the query is doing the work rather than the client.

  Then sort by a column the service did not whitelist and watch nothing happen.
  That is chapter thirty two's whitelist, and it is also the reason the first
  sort you add by hand appears to be ignored.
]

== Step six: the admin screen, customised

The generated list is correct and plain. Two changes make it usable, and both
live in `apps/admin/app/_dashboard/resources/task/index.tsx`.

A status badge rather than a word:

```tsx
{
  key: "status",
  label: "Status",
  sortable: true,
  render: (row) => {
    const tone = {
      todo:  "bg-slate-500/15 text-slate-300",
      doing: "bg-amber-500/15 text-amber-400",
      done:  "bg-emerald-500/15 text-emerald-400",
    }[row.status]
    return <span className={`rounded px-2 py-0.5 ${tone}`}>
      {row.status}</span>
  },
}
```

And an overdue due date, because a date that is late should look late:

```tsx
{
  key: "due",
  label: "Due",
  sortable: true,
  render: (row) => {
    const late = row.status !== "done" &&
      new Date(row.due) < new Date()
    return <span className={late ? "text-red-400" : ""}>
      {formatDate(row.due)}</span>
  },
}
```

Re-run the generator for `Task` afterwards and confirm these survive. They are
in the file the generator writes once.

== Step seven: a page the admin cannot be

An operator console is not a product. The team wants a board, so write one as
an ordinary Next.js page using the generated hooks:

```tsx
export default function BoardPage() {
  const { data, isLoading } = useTasks({ page_size: 200 })
  if (isLoading) return <BoardSkeleton />

  const columns = ["todo", "doing", "done"] as const
  return (
    <div className="grid gap-4 md:grid-cols-3">
      {columns.map((status) => (
        <Column
          key={status}
          title={LABELS[status]}
          tasks={data.data.filter((t) => t.status === status)}
        />
      ))}
    </div>
  )
}
```

Moving a card calls the same API the admin does:

```tsx
const { mutate } = useUpdateTask()
mutate({ id: task.id, status: "doing" })
```

Nothing here is framework-specific. It is React against typed hooks, and the
types came from the Go model through chapter thirty six's shared package, so
renaming `due` in the model makes this page fail to compile rather than fail in
front of a user.

#gotcha[
  Resist putting the "cannot close with open subtasks" check in this page. It
  is one line of JavaScript and it is wrong, because the board is not the only
  caller.

  Call the endpoint, show the 409's message, and let the rule live once. The
  test is always chapter twenty eight's: would it still have to be true if a
  nightly job did it?
]

== Step eight: a form for people without accounts

Bug reports should not require an admin login. One command writes the page:

```bash
grit expose form Task --to apps/web/app/report/page.tsx
```

That emits a Next.js page with one input per field, submitting through the
authenticated create hook. For a genuinely public form, create a form share in
the admin under System, then:

```bash
grit expose form Task --to apps/web/app/report/page.tsx \
  --public-share --token 9CkLh7gJZ...
```

which posts to the public submission endpoint instead, with no account
required.

#gotcha[
  A public write endpoint is a spam target the hour it is discovered. The
  share token is revocable in the admin, and the form share is rate limited,
  but the decision about what a stranger may create is yours.

  The safer shape for a bug reporter: a public form that creates a `Report`,
  and a human who turns a report into a `Task`. One extra model, and no
  stranger ever writes to the table the team works from.
]

== Step nine: the things that make it a product

Four small commands, each closing something a client will ask about:

```bash
grit add role MANAGER
grit backup
grit doctor
grit routes
```

A role that the sidebar, the forms and the seeded permissions all know about. A
backup you have actually taken. A health check across the project. A route list
you have read, looking for anything registered outside the authenticated group.

== What it cost

Count honestly, including the parts nobody counts:

#table(
  columns: (1fr, auto),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Three resources, migrated, in the admin], [four commands],
  [The business rule and its endpoint], [about 40 lines],
  [Two customised columns], [about 25 lines],
  [The board page], [about 80 lines],
  [Everything else: auth, lists, forms, validation, types, hooks], [generated],
)

Compare that with chapter forty eight, where one resource with no admin panel
and no frontend took about six hundred lines. That ratio is the entire argument
of this book, and now you have both numbers from your own machine rather than
from a marketing page.

#project(title: "Project 34: finish the tracker")[
  Take the tracker from working to usable. Four additions, each of which uses a
  different part of Part IV.

  *One: activity.* Generate a `Activity` resource with `--append-only` and
  write a row whenever a task changes status. Append-only means no edits and no
  deletes, which is what an audit trail is.

  *Two: my work.* Add a page listing the signed-in user's open tasks across
  projects, ordered by due date. Use the generated hook with a filter rather
  than fetching everything and filtering in the browser, and explain in a
  comment why that distinction matters at ten thousand rows.

  *Three: a dashboard widget.* Add a stat card showing overdue tasks, and a
  chart of completions per week. Both read from the API, not from the browser's
  copy.

  *Four: the rule, tested.* A Go test that creates a task with one open subtask
  and asserts `Close` returns `ErrHasOpenSubtasks`. Then a second test that
  closes the subtask first and asserts success. Run them in CI.

  Then delete something on purpose: `grit remove resource Comment`, confirm the
  project still builds, and read what the command removed. A generator that
  cannot take its own work back is a generator you cannot experiment with.
]

#tryit[
  Open the tracker in a browser as two different users in two windows: an admin
  and an ordinary member. Walk the same screens.

  The differences you see are chapter thirty five working. The ones you expected
  and do not see are permissions you have not configured yet, and it is much
  better to find them here than in a client demo.
]

#assignment(name: "Task Tracker - Assignment One")[
  Build the tracker in your own repository, under
  `part5-building/37-tracker`, and push it as a working project rather than a
  set of notes.

  *One.* The three resources, the subtask relationship, the `Close` rule, its
  endpoint and its 409. The rule must be enforced from a Go test that never
  touches HTTP.

  *Two.* The board page, with the status badge and the overdue styling in the
  admin resource file.

  *Three.* Seed four hundred tasks and record, in `NOTES.md`, the time to first
  paint on the list page and what changes at a hundred thousand rows.

  *Four.* `grit routes` output committed as `ROUTES.txt`, with one line naming
  any route you would not want public.

  *Five.* The honest comparison. Count the lines you wrote yourself, put them
  next to the chapter forty eight number, and write a paragraph on which parts
  of the generated code you would keep if you had to maintain this for five
  years.
]
