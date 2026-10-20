#import "../template.typ": *

= Before you start: a repository, and a leaderboard

Every chapter in this book ends with an assignment. You write it, you commit it
with an exact message, and you push it. This chapter sets that up, and it takes
about ten minutes.

Do it now rather than at chapter four, because the habit is the point and a
habit started late is usually not started.

== Why the assignments are pushed rather than just done

Three reasons, and the third is the one that matters.

*You will finish them.* Work that leaves your machine gets finished. Work that
stays in a scratch file gets abandoned at the interesting part, which is
usually one step before the thing it was teaching.

*You will have something to show.* By the end of this book your repository
holds fifty-odd small Go programs, each one solving a named problem, each with
a readable commit history. That is a more convincing portfolio than a single
large project, because it shows range and it shows persistence.

*You are not doing it alone.* The commits are scored on a leaderboard shared
with everybody else working through this book, which is what the rest of this
chapter is about.

== Step one: make the repository

```bash
mkdir go-journey && cd go-journey
git init
```

One repository for the whole book. Not one per chapter: the history is the
record, and a single timeline is what makes the progression visible.

Give it a structure that will still make sense at chapter fifty:

```
go-journey/
├── README.md
├── part1-language/
│   ├── 01-why-go/
│   ├── 02-hello/
│   └── ...
├── part2-ideas/
├── part3-standard-library/
├── part4-grit/
├── part5-building/
└── part6-from-scratch/
```

Each assignment gets a directory with its own `main.go`, and where it needs
one, its own `go.mod`:

```bash
mkdir -p part1-language/02-hello && cd part1-language/02-hello
go mod init go-journey/part1/hello
```

#gotcha[
  A separate `go.mod` per assignment keeps them independent, so an experiment
  in chapter twenty cannot break chapter three. The alternative, one module at
  the root with a package per chapter, also works and means one `go.mod`, but
  every `go build ./...` then compiles everything.

  For a learning repository, independent modules are the better trade. Use a
  Go workspace if you want `go build ./...` to still work across all of them:

  ```bash
  go work init
  go work use ./part1-language/02-hello
  ```
]

Write a `README.md` on day one, because you will not write it later:

```markdown
# go-journey

Working through *Go for JavaScript Developers*, one assignment per chapter.

| Part | What it covers |
|---|---|
| 1 | The Go language, against the JavaScript I already knew |
| 2 | Interfaces, errors, goroutines, channels, context |
| 3 | Files, HTTP, databases, shipping a binary |
| 4 | The Grit framework |
| 5 | Building real applications |
| 6 | Building things from scratch: a parser, a database, a server |
```

Then push it:

```bash
git add .
git commit -m "Setup - Repository"
git branch -M main
git remote add origin git@github.com:YOUR-NAME/go-journey.git
git push -u origin main
```

Make the repository *public*. The leaderboard reads public activity, and a
public learning repository has never cost anybody a job.

== Step two: the Committed app

Install Committed from the Play Store:

```
play.google.com/store/apps/details?id=com.desishub.committed
```

Sign up *with your GitHub account*, which is how it sees your commits. There is
nothing to configure after that.

== Step three: join the Gopher League

Open the app and join with the code:

```
5TCLJ27Z
```

Or follow the link, which does the same thing:

```
https://committed-api.gritframework.dev/join/5TCLJ27Z
```

The group's description says what it is:

#align(center)[
  #block(width: 92%, inset: 10pt, radius: 3pt,
    fill: rgb("#f6f7f9"), stroke: 0.5pt + rgb("#e3e6ea"))[
    #set text(size: 10pt, style: "italic")
    Gopher League is a friendly leaderboard for developers learning Go
    together. Every commit, merged PR, review and closed issue earns points.
    Ship Go, climb the table, and get crowned weekly champion.
  ]
]

Commits score. So do merged pull requests, reviews and closed issues, which is
worth knowing: reviewing somebody else's assignment is worth points, and it is
also one of the better ways to learn.

== Step four: the commit convention

Every assignment in this book carries an exact commit message, printed in its
own box:

```bash
git commit -m "Variables - Assignment One"
```

*Use the message exactly as printed.* It is not decoration. The leaderboard
matches on it, and a commit called "did the variables thing" earns nothing,
which is a sad way to lose a week's points.

The pattern is always the same:

```
<Topic> - Assignment <Number>
```

So the first few are `Hello Go - Assignment One`, `Variables - Assignment One`,
`Variables - Assignment Two`, and so on.

#gotcha[
  Commit *per assignment*, not once per chapter and not once per day.

  Partly because each one scores. Mostly because a commit that says exactly
  what it contains is the single most useful habit in this book, and a history
  of fifty precise commits is worth more to you in six months than fifty files
  in one commit called "work".
]

== The loop, for every chapter from here

```bash
# 1. a directory for the assignment
mkdir -p part1-language/03-types && cd part1-language/03-types
go mod init go-journey/part1/types

# 2. write it, and run it until it works
go run .

# 3. commit with the exact message from the book
git add .
git commit -m "Variables - Assignment One"
git push
```

Four commands, maybe fifteen minutes per chapter. The leaderboard updates on
its own.

== There is more practice on the website

The Grit documentation has a Go guide with forty six interactive exercises,
which run in the browser with a solution to check against:

```
gritframework.dev/docs/prerequisites/golang
```

They cover the same ground as Parts I to III and are shorter than this book's
assignments, so they work well as a warm-up before a chapter or as revision
after it. Where one lines up with a chapter, that chapter's assignment says so.

The division worth knowing: *the website exercises are for practising a
concept, and this book's assignments are for building something and pushing
it*. Do both if you have time, and if you only do one, do the assignment,
because that is the one that leaves your machine.

#assignment(name: "Setup - Repository")[
  Set up everything above, and prove it works.

  *One.* Create the `go-journey` repository, public, with the directory
  structure and the README.

  *Two.* Install Committed, sign up with GitHub, and join Gopher League with
  the code `5TCLJ27Z`.

  *Three.* Add a file `part1-language/00-setup/notes.md` answering the question
  from the end of chapter one in two sentences: what is the specific thing in
  your current stack that made you pick this book up? You will want that answer
  again at chapter twenty six.

  *Four.* Commit and push, then open the app and confirm your commit appears on
  the leaderboard. If it does not, the usual cause is that the repository is
  private or the email on the commit does not match your GitHub account:

  ```bash
  git config user.email "the-email-on-your-github-account"
  ```
]

#tryit[
  Open the leaderboard and find your position. It will be last, which is the
  correct place to start and the only time the number matters.

  Then look at who is near the top and how often they commit. Almost always the
  answer is "a little, most days" rather than "a lot, once". That pattern is
  the actual lesson of the leaderboard, and it applies to learning Go whether
  or not anybody is counting.
]
