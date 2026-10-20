#import "../template.typ": *

= The twenty six problems, and what a framework is for

You can now write Go. You have built a CLI, an HTTP service, a database layer
and a deployable binary with no framework at all. So the fair question before
part four is: *why would you add one?*

The answer is not "to write less code". It is that every project runs into the
same list of problems, the list is longer than people admit, and most of it has
nothing to do with the product you were hired to build.

== Two lists, and they belong to different people

The problems split in two, and the split matters because the two halves are
paid for by different people at different times.

*The developer's twelve* cost you the first two weeks. You feel them
immediately and you can name them from memory.

*The client's fourteen* are questions asked before anybody mentions your actual
feature. They decide whether there is a project at all, and they are always
promised for later.

== The developer's twelve

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*D1*], [The copy-paste start. Auth, storage and email lifted from the last repository, bugs included],
  [*D2*], [Wiring the monorepo: workspaces, ports, env files, Docker, before anything runs],
  [*D3*], [The same type written twice, or four times: Go, TypeScript, the validation schema, the mobile client],
  [*D4*], [The same UI for every table. Twenty tables, twenty slightly different list pages],
  [*D5*], [Testing performance is its own project, so it waits for a quiet week that does not come],
  [*D6*], [No window into the database without installing one more tool],
  [*D7*], [The API nobody wrote down. Hand-written OpenAPI is accurate for one afternoon],
  [*D8*], [Import and export are always a ticket, and always a week rather than an afternoon],
  [*D9*], [Nothing to develop against: empty tables, and mail you cannot read],
  [*D10*], [Changing the schema is the scary part: the migration that half ran, with no way back],
  [*D11*], [The frontend is somebody else's job, and there is nobody else],
  [*D12*], [Generating it with a model costs money every time, and none of it compounds],
)

Most of these you have met in this book. You wrote the error boundary in
chapter eleven, the middleware in twenty one, the config validation in twenty
four. None of it was hard. All of it was time, and none of it was the product.

#bridge[
  D11 deserves a sentence, because it is the reason many readers picked this
  book up. You are fluent in React and now also in Go, which makes you the
  rare person who can do both. The cost is that you are now the person who does
  both, for every project, forever. A generated admin panel is not about
  disliking frontend work; it is about not rebuilding the same table for the
  twentieth time.
]

== The client's fourteen

These arrive as questions, usually in the first meeting, before anybody has
described the actual product.

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*B1*], [Can different people see different things?],
  [*B2*], [What happens if we get hacked?],
  [*B3*], [Does it work without internet?],
  [*B4*], [How will we know when it breaks?],
  [*B5*], [Is any of this backed up?],
  [*B6*], [What if it gets busy?],
  [*B7*], [Who changed this?],
  [*B8*], [Our users photograph things on their phones],
  [*B9*], [Can two of our clients use one system?],
  [*B10*], [Does the screen update by itself?],
  [*B11*], [Can we switch this on for one customer only?],
  [*B12*], [What if two people edit the same record?],
  [*B13*], [Does it talk to our other systems?],
  [*B14*], [Someone has asked us to delete their data],
)

The honest answer to every one of them is "yes, in about three weeks", and the
client hears that as a no.

Notice what they have in common: each is real work with *no demo attached*. You
cannot show a stakeholder an audit log. So each gets promised for phase two,
and offline, audit and observability are the three that never arrive.

== Why these survive a decade of excellent libraries

Go's ecosystem is superb at solving one problem each. A router, a queue, a
migration tool, an image decoder, all excellent. So why is the list still here?

*Libraries solve one problem each. Nobody hands you the shape.* A router
solves routing. Nothing solves "what does an application made of these look
like", so the shape is hand-carved on every project, slightly differently.

*A template rots the day you clone it.* A boilerplate repository freezes the
bugs of the day you started. The fix you make on project five never reaches
projects one to four, and you cannot tell which repository has which. This is
D1 with extra steps.

*The client's fourteen are always phase two.* Expensive, invisible, always cut.
Which is a precise definition of what belongs in a framework.

#gotcha[
  D12 is the newest and it is worth stating carefully, because the obvious
  version of the argument is wrong.

  An agent will write this scaffold for you. A Grit project measures about
  74,500 lines, which is roughly three quarters of a million output tokens: at
  current frontier pricing, tens of dollars. That is cheap for two weeks of
  work and anybody who tells you the cost alone is the objection is
  overselling.

  The objection is the other three properties. You pay it *again* on the next
  project. The model writes something *different* each time, so the auth it
  gives you in March and in June are different code with different bugs. And
  because they are different, a fix you make in one project has nowhere to go.

  That is D1, the copy-paste start, with a bill attached.
]

== The question a framework has to answer

So the question is not *which library*. It is:

#align(center)[
  #block(width: 92%, inset: 10pt, radius: 3pt, fill: rgb("#f6f7f9"), stroke: 0.5pt + rgb("#e3e6ea"))[
    #set text(size: 11pt, style: "italic")
    Can the shape of an application be generated, owned outright by the team
    that generated it, and still improve after they own it?
  ]
]

Three clauses, and they fight each other.

*Generated* is easy. Plenty of tools generate code.

*Generated and owned* is harder, and is where most scaffolders stop: they hand
you a repository and walk away.

*Generated, owned, and still improving* is the actual problem, because the
moment code is yours, nobody can fix it but you. That third clause is what
`grit upgrade` exists for, and chapter 29 is about how the line between "yours"
and "the framework's" is drawn.

== What this means for the rest of the book

Part four explains Grit's answer to that question: what `grit new` produces,
how the architecture is laid out, why the *resource* is the unit of generation,
and how the generated code stays plain Go that you can read with everything you
learned in parts one to three.

Part five builds real applications with it.

Two promises about how this is written. Every claim about what Grit does comes
from the codebase rather than from marketing, and where it does not close a
problem, the book says so: B6, "what if it gets busy", is the honest example,
because a framework can hand you a cache, read replicas and a job queue, and it
cannot tell you which of them your bottleneck needs.

#bridge[
  You could also solve all twenty six yourself. You now know enough to. The
  question is whether you want to solve them on this project, and then again on
  the next one, and whether the version you write on project five can reach
  projects one to four.

  If your application is genuinely unusual, a handful of endpoints with an
  unusual data model and no admin surface, then the honest advice is to keep
  writing what you wrote in part three. A framework earns its place on
  CRUD-shaped products with an admin panel, several entities and a client who
  asks the fourteen questions. Most business software is exactly that, which is
  why these frameworks exist.
]

#tryit[
  Before the next chapter, go through both tables and mark each problem as one
  of three things: *solved* on your current project, *promised* for later, or
  *not yet asked about*. Be honest about the middle column.

  Then count the "promised" ones. That number is what a framework is offering
  to buy from you, and whether the trade is worth it depends almost entirely on
  how large it is.
]
