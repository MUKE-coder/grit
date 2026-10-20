#import "../template.typ": *

= Why Go, when you already ship JavaScript

You already build and ship software. You know React, you know Next.js, and you
can take an idea to production without asking anybody how. So the first
question this book owes you is not "what is Go", it is "why would you spend a
month learning a second language when the one you have works".

Here is the honest answer, and then the honest list of what you give up.

== What actually pushed you here

Most JavaScript developers do not come to Go out of curiosity. They come
because something in their stack got painful in a specific way. Usually one of
these four.

*The deployment got heavy.* Your Next.js app needs Node on the server, a
`node_modules` directory with a few hundred megabytes in it, and a build step
that has to run somewhere. You want to copy one file onto a server and run it.
A compiled Go program is one file. It has no runtime to install, no package
directory, and no build step on the target machine. You build it on your laptop
and the output is the thing that runs.

*The work stopped being request-and-respond.* An API that answers a question in
forty milliseconds is a perfect job for Node. A process that holds ten thousand
open connections, or runs a job queue, or chews through a large file while
staying responsive, is where the single-threaded event loop stops being an
elegance and starts being a constraint you work around.

*You wanted the errors earlier.* In JavaScript, a renamed field is found by a
user. TypeScript helps enormously and you should keep using it, but it is a
layer on top of a language that does not have types, and it disappears at
runtime. In Go the types are the language. A program that compiles has already
had a whole class of mistake removed from it.

*Somebody said the word "backend" and meant it.* Payment handling, a scheduler,
an audit trail, something that must still be correct at 3am when nobody is
watching. Go was built by people who operate systems like that for a living,
and it shows in what the language refuses to let you do.

#bridge[
  None of this makes Go better than JavaScript. It makes it different in ways
  that matter for one kind of work. You are not replacing React. In the second
  half of this book you will have a Go backend and a React frontend in the same
  project, talking to each other with types generated from one definition. That
  is the shape most people actually want.
]

== The shape of Go, in one page

Before any syntax, here is the character of the language, because it explains
every decision you are about to meet.

Go is *small on purpose*. The entire specification is about one hundred pages.
There are twenty five keywords. JavaScript has roughly double that and a
standard library that has grown for thirty years. You can hold all of Go in
your head, which is a design goal rather than an accident, and it is why Go
teams onboard people in days rather than weeks.

Go is *explicit on purpose*. There is no implicit type coercion, no `this` that
changes meaning depending on how a function was called, no prototype chain to
walk. When you read a line of Go, what happens is what is written. This is
sometimes tedious and it is never surprising, and after a few weeks you will
notice you have stopped debugging a category of problem entirely.

Go is *opinionated about formatting*, which sounds trivial and is not. There is
one correct way to format Go code, a tool called `gofmt` produces it, and every
Go project on earth uses it. No Prettier config, no debate about semicolons, no
`.editorconfig`. The argument is over because the language ended it.

Go *compiles fast*. Not "fast for a compiled language", fast in absolute terms:
a medium project builds in a couple of seconds. This matters more than it
sounds, because it means the compile step never becomes a reason to avoid
checking your work.

== What you give up, stated plainly

A book that only lists advantages is marketing. Here is the other side, and you
should weigh it before chapter two.

*You will write more lines.* Go has no `map`, `filter` or `reduce` built into
the language the way JavaScript does. Where you write one chained expression,
you will often write a four-line loop. Go's position is that the loop is easier
to read six months later. You may disagree, and you will still write the loop.

*Error handling is verbose.* You are about to meet `if err != nil` and you will
type it hundreds of times. There is no `try`/`catch` wrapping half a function.
Every single operation that can fail says so, at the point it can fail. This is
the single biggest adjustment coming from JavaScript, and chapter four is
entirely about making peace with it.

*The ecosystem is smaller.* npm has over two million packages. Go has far
fewer. For web application work this matters less than you would think, because
Go's standard library covers an enormous amount of what you would reach for a
package to do in Node. But if you need an obscure SDK, check before you commit.

*No null safety.* Go has `nil`, and a nil pointer dereference will crash your
program at runtime exactly like reading a property of `undefined` crashes
yours. Go's type system is simpler than TypeScript's, and in this one specific
area TypeScript is genuinely ahead.

#gotcha[
  If someone tells you Go is simply better than JavaScript, stop listening to
  them. The two languages were built for different jobs by people with
  different problems. Go's concurrency is better than the event loop for
  sustained work, and the event loop is a beautiful fit for a UI. Use both.
]

== The one example that explains everything

Here is the same program in both languages: read three URLs, and print how long
each took. It is worth studying before you understand any Go syntax, because
the difference in shape is the whole argument.

#compare(
  "async function timeAll(urls) {
  const results = await Promise.all(
    urls.map(async (url) => {
      const start = Date.now();
      await fetch(url);
      return { url, ms: Date.now() - start };
    })
  );
  for (const r of results) {
    console.log(r.url, r.ms);
  }
}",
  "func timeAll(urls []string) {
    var wg sync.WaitGroup
    results := make([]result, len(urls))

    for i, url := range urls {
        wg.Add(1)
        go func() {
            defer wg.Done()
            start := time.Now()
            http.Get(url)
            results[i] = result{url, time.Since(start)}
        }()
    }

    wg.Wait()
    for _, r := range results {
        fmt.Println(r.url, r.ms)
    }
}",
)

The JavaScript is shorter and, if you know JavaScript, clearer. Do not let
anybody tell you otherwise.

But look at what the Go version is doing. The word `go` in front of that
function call starts a *goroutine*: an independent line of execution that the
Go runtime schedules onto a real operating system thread. Three URLs here, but
the same code with fifty thousand URLs would start fifty thousand goroutines,
and Go would handle it, because a goroutine costs about two kilobytes of memory
rather than the megabyte an OS thread costs.

The JavaScript version is also concurrent, and for three HTTP requests it is
just as fast, because all three are waiting on the network and the event loop
is excellent at waiting. The difference appears when the work is not waiting.
If each of those tasks had to parse a large document or resize an image, the
JavaScript version would do them one at a time on one core, and the Go version
would use every core you have.

#bridge(title: "The mental model shift")[
  In JavaScript, concurrency is something you *await*, and parallelism is
  something you get by starting more processes. In Go, concurrency is something
  you *start*, with one keyword, and parallelism comes free because the runtime
  spreads goroutines across cores. Chapter eleven is where this clicks properly.
  For now, notice only that `go doSomething()` is the entire syntax.
]

== What this book does

This book has two halves and they have different jobs.

*Part one through three teach you Go*, assuming you know JavaScript and nothing
else. Every concept is introduced by naming the JavaScript thing it resembles,
then showing exactly where the resemblance breaks. You will build small things
constantly: a CLI tool, a JSON API with no framework at all, a URL shortener
backed by a real database. By the end you can read and write Go without this
book open.

*Part four and five teach you Grit*, which is a framework that generates a Go
API, a React admin panel and the typed client between them from one description
of your data. The reason it comes second is that a framework you cannot read is
a liability. Once you know Go, Grit's generated code is just Go you did not
have to type, and you can change any of it.

#project(title: "What you will have built by the end")[
  A command line tool that formats JSON. An HTTP API written with nothing but
  the standard library. A URL shortener with a SQLite database. A task tracker,
  a blog with a working admin panel, and an invoicing application with
  customers, line items and money handled correctly. All of them small, all of
  them finished, none of them a toy that stops at the interesting part.
]

== How to read it

Type the code. Do not copy and paste it. This is advice everybody gives and
almost nobody follows, so here is the specific reason it matters for Go: the
compiler's error messages are a large part of how you learn the language, and
you only see them if you make mistakes. Typing produces mistakes. Pasting does
not.

Run everything. Every code block that is a complete program is marked, and you
should run it, break it deliberately, and run it again. Go builds in seconds,
so the cost of an experiment is almost nothing.

Do the "Your turn" boxes. They are short, and they are the difference between
recognising Go and writing it.

#tryit[
  Before the next chapter, answer this for your own situation in two sentences:
  what is the specific thing in your current stack that made you pick this book
  up? Write it down. In chapter twenty three there is a list of twenty six
  problems that every project runs into, and you will want to know which of
  them is yours.
]

#assignment(name: "Why Go - Assignment One")[
  Before any Go, write down why you are here, because chapter twenty six will
  ask you to check the answer.

  Create `part1-language/01-why-go/NOTES.md` and answer three things in a
  paragraph each.

  *One.* Which of the four pressures in this chapter is yours: the deployment,
  the long-running work, the late errors, or the word "backend"? Name the
  actual system it happened on.

  *Two.* What did you try first, and why did it not settle the matter? There is
  almost always a library or a rewrite that nearly worked.

  *Three.* What would have to be true in six months for this book to have been
  worth the time? Be specific enough that you could be wrong.]
