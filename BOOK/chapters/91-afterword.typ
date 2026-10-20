#import "../template.typ": *

= Afterword: what to do on Monday

You have read a book, which is the easy part. What follows is the shortest path
from here to being the person on a team who can be handed a backend.

== The three things that actually move you

*Finish the repository.* If you did the assignments, you have fifty-odd
programs and a commit history that shows a year of progress compressed into
however long this took. If you skipped them, the single highest-value thing you
can do this week is go back to chapter zero, make the repository, and work
forward doing only the assignments. Reading teaches recognition; writing
teaches recall, and only one of those survives an interview or a deadline.

*Ship one thing to one real user.* Not a portfolio piece, not a clone. Pick the
application from Part V closest to something somebody you know does by hand,
build it, and get them to use it for a month. Everything that matters and is
not in this book, the questions users ask, the data they paste in, the ways
they misunderstand your forms, arrives in that month.

*Read code that is better than yours.* The Go standard library is unusually
readable, and it is right there: `go doc -src net/http.Server`. Read `io`,
`errors`, `sync.Once`. They are short and they are the work of people who
thought hard about exactly the things this book has been teaching you to think
about.

== The uncomfortable parts, named

Four things that will be harder than this book made them look, said plainly so
that you recognise them rather than concluding you are bad at this.

*Concurrency is easy to start and hard to be sure about.* `go` is one keyword
and correct concurrent code is a skill. Use the race detector, prefer channels
and simple ownership, and be suspicious of any design you cannot explain in two
sentences.

*The framework is not the hard part.* Modelling is. Deciding that an enrolment
is a thing rather than a join table, that a sent invoice is immutable, that a
date is not a timestamp: those decisions determine whether the project is
pleasant in a year, and no generator makes them for you.

*Production is a different skill.* Writing the feature is a day. Operating it,
backing it up, knowing when it is slow, deploying it without dropping requests,
is the other half of the job and the half that is usually paid for.

*You will be wrong about performance.* Every time. Measure before optimising,
and prefer the clearer code until a profile says otherwise. Chapter forty five
exists for this.

== What to learn next, in order

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [SQL, properly], [Indexes, query plans, transactions and isolation levels.
    It will do more for your applications than any language feature],
  [One more Part VI build], [Pick the one you were most curious about. They
    are the fastest way to deepen the language],
  [Profiling], [`pprof`, flame graphs, allocation profiles. Three evenings,
    and you can answer questions most teams guess at],
  [Linux basics], [Processes, file descriptors, signals, systemd. Your binary
    runs there],
  [A second domain], [Payments, scheduling, inventory. Depth in a domain
    outlasts depth in a framework],
)

== On using a framework at all

This book argued for one for two hundred pages and then had you write a CRUD
API by hand to see what it costs. Hold both.

A framework is worth it when it removes work you have already done and would do
the same way again. It is a trap when it removes decisions you should be
making, or when you cannot read what it produced. The reason Part VI exists is
so that neither is true for you: the generated service is code you could have
written, which means you can also judge when to leave it.

If you take one sentence from the whole book, take this one: *the thing that
makes software expensive is rarely the code you wrote, it is the code nobody
understands*. Generated or hand-written, read it.

== Thank you, and one request

If you got here, you did a serious amount of work. Go is a language people
describe as boring, and that is the compliment it is meant to be: it gets out
of the way of the part that is actually hard, which is understanding the
problem you were asked to solve.

The request: when you find an error in this book, say so. A wrong claim about
the framework is checkable against the codebase, and a wrong claim about Go is
checkable against the compiler. Both are worth fixing for the next person, who
is in exactly the position you were in on page one.

Now go and build something somebody needs.
