#import "../template.typ": *

= Build your own CLI framework

Every tool you write grows a command line eventually, and most people reach
straight for Cobra without ever knowing what it does. A command framework is
about three hundred lines. Writing one makes every CLI you use afterwards
legible, and it is the fastest way to learn what separates a tool people enjoy
from one they tolerate.

We are building: nested subcommands, per-command flags, generated help, "did
you mean", correct exit codes, and cancellation on Ctrl-C.

== What a good CLI does, before any code

The design decisions matter more than the parser:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [Exit codes], [0 success, 1 failure, 2 usage error. Scripts read these],
  [Streams], [Results to stdout, diagnostics to stderr. Always],
  [Help], [Generated from the commands, never written twice],
  [Errors], [Say what to do next, not just what went wrong],
  [Cancellation], [Ctrl-C stops the work, it does not kill it mid-write],
)

#bridge[
  The stream rule is the one JavaScript developers break most, because
  `console.log` goes to stdout and it is right there.

  The consequence shows up the first time somebody writes
  `mytool list > out.json` and finds a progress message in the middle of their
  JSON. Results on stdout, everything else on stderr, and the tool becomes
  something people can build on.
]

== The command

```go
type Command struct {
    Name  string
    Short string                       // one line, for the index
    Long  string                       // paragraphs, for the help page
    Usage string                       // "add <file> [more...]"

    Flags *flag.FlagSet
    Run   func(ctx context.Context, args []string) error

    subs map[string]*Command
}

func (c *Command) Add(sub *Command) {
    if c.subs == nil {
        c.subs = map[string]*Command{}
    }
    c.subs[sub.Name] = sub
}
```

Everything else in this chapter is built from that struct. Note what it does
not have: no global state, no `init()` registration, no singleton root. A
command tree is a value you can build in a test.

== Dispatch

Walking the tree is a loop, not recursion, and it stops at the first argument
that is not a known subcommand:

```go
func (c *Command) find(args []string) (*Command, []string) {
    cmd := c
    for len(args) > 0 {
        sub, ok := cmd.subs[args[0]]
        if !ok {
            return cmd, args
        }
        cmd, args = sub, args[1:]
    }
    return cmd, args
}
```

So `tool db migrate --dry-run status.json` finds `migrate` and hands it
`["--dry-run", "status.json"]`.

Then parse that command's own flags and run it:

```go
func (c *Command) Execute(ctx context.Context, argv []string) error {
    cmd, rest := c.find(argv)

    if cmd.Flags != nil {
        cmd.Flags.SetOutput(io.Discard) // we print our own errors
        if err := cmd.Flags.Parse(rest); err != nil {
            return usageError{cmd: cmd, err: err}
        }
        rest = cmd.Flags.Args()
    }

    if cmd.Run == nil {
        return usageError{cmd: cmd, err: errNoRun}
    }
    return cmd.Run(ctx, rest)
}
```

#gotcha[
  `flag.Parse` stops at the first non-flag argument. `tool add file.txt -v`
  leaves `-v` in the arguments rather than setting the flag, and the user sees
  nothing happen.

  Go's own tools accept this; GNU tools do not. Decide which you are, say so in
  the help, and if you want flags anywhere, sort the arguments before parsing
  rather than writing your own parser.
]

== Two kinds of failure

This is the design decision that most improves a CLI, and almost nobody makes
it deliberately.

```go
type usageError struct {
    cmd *Command
    err error
}
```

*A usage error* means the command was typed wrong: an unknown flag, a missing
argument. Print the error, then the usage, and exit 2. The flag list is exactly
what the user needs.

*Every other error* means the command was fine and the work failed: a file was
missing, the server refused, the database was down. Print the error and exit 1,
*without the usage*.

#gotcha[
  Printing the flag list under "connection refused" is the most common CLI
  mistake there is. It buries the one sentence that matters under forty lines
  of options, and it implies the user typed something wrong when they did not.

  Grit's own CLI sets `SilenceUsage` the moment it knows the flags parsed and
  the project was found, for exactly this reason. If you take one thing from
  this chapter into a tool you already maintain, take this one.
]

```go
func Main(root *Command, argv []string) int {
    ctx := withSignals(context.Background())

    err := root.Execute(ctx, argv)
    switch {
    case err == nil:
        return 0
    case errors.Is(err, context.Canceled):
        fmt.Fprintln(os.Stderr, "cancelled")
        return 130                       // 128 + SIGINT
    case errors.As(err, new(usageError)):
        var ue usageError
        errors.As(err, &ue)
        fmt.Fprintf(os.Stderr, "error: %v\n\n", ue.err)
        ue.cmd.PrintUsage(os.Stderr)
        return 2
    default:
        fmt.Fprintf(os.Stderr, "error: %v\n", err)
        return 1
    }
}
```

`130` is the conventional code for "killed by Ctrl-C", and scripts do check it.

== Help, generated

Help that is written by hand goes stale on the first new flag. Generate it:

```go
func (c *Command) PrintUsage(w io.Writer) {
    if c.Long != "" {
        fmt.Fprintf(w, "%s\n\n", strings.TrimSpace(c.Long))
    }
    fmt.Fprintf(w, "Usage:\n  %s %s\n", c.path(), c.Usage)

    if len(c.subs) > 0 {
        fmt.Fprintf(w, "\nCommands:\n")
        for _, name := range sortedKeys(c.subs) {
            fmt.Fprintf(w, "  %-12s %s\n",
                name, c.subs[name].Short)
        }
    }
    if c.Flags != nil {
        fmt.Fprintf(w, "\nFlags:\n")
        c.Flags.SetOutput(w)
        c.Flags.PrintDefaults()
    }
    fmt.Fprintf(w, "\nRun '%s <command> --help' for more.\n", c.path())
}
```

The `io.Writer` parameter is not decoration. It is what makes the help testable
without capturing os.Stdout, which is chapter twenty's lesson applied to a
command line.

== Did you mean

A typo should not produce a wall of help. It should produce one suggestion:

```
tool: unknown command "migarte"

Did you mean?
    migrate
```

Edit distance, bounded so nonsense suggests nothing:

```go
func suggest(want string, have []string) []string {
    limit := len(want)/3 + 1
    var out []string
    for _, candidate := range have {
        if distance(want, candidate) <= limit {
            out = append(out, candidate)
        }
    }
    return out
}
```

`distance` is Levenshtein, about twenty lines with two rows rather than a full
matrix. Writing it once is worth it: the same function powers "did you mean" in
every tool you will ever write, and in your API's error messages for an unknown
field.

== Cancellation that is not a kill

```go
func withSignals(parent context.Context) context.Context {
    ctx, cancel := signal.NotifyContext(parent,
        os.Interrupt, syscall.SIGTERM)

    go func() {
        <-ctx.Done()
        cancel()
        // A second Ctrl-C means business.
        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        <-c
        os.Exit(130)
    }()
    return ctx
}
```

The first interrupt cancels the context, so the command can finish writing the
file it is in the middle of and exit cleanly. The second one exits immediately,
because a user pressing Ctrl-C twice has stopped asking politely.

Every long-running command then has one obligation, from chapter sixteen:
*check `ctx.Done()`*. A command that ignores it is a command Ctrl-C cannot
stop.

== Output humans read and machines parse

The same command often has two audiences:

```go
if *asJSON {
    return json.NewEncoder(out).Encode(results)
}
for _, r := range results {
    fmt.Fprintf(out, "%-20s %s\n", r.Name, r.Status)
}
```

Two rules make this work. *The JSON goes to stdout and nothing else does*, or
piping it breaks. And *colour is disabled when stdout is not a terminal*, or
your pipeline fills with escape codes:

```go
func useColour(w io.Writer) bool {
    if os.Getenv("NO_COLOR") != "" {
        return false
    }
    f, ok := w.(*os.File)
    if !ok {
        return false
    }
    info, err := f.Stat()
    return err == nil && info.Mode()&os.ModeCharDevice != 0
}
```

`NO_COLOR` is an agreed convention. Honour it, and the people who need it will
never have to ask you to.

== Testing a CLI

Because nothing is global, a test builds a tree and runs it:

```go
func TestMigrateRejectsUnknownFlag(t *testing.T) {
    var out bytes.Buffer
    root := newRoot(&out)

    code := Main(root, []string{"db", "migrate", "--nope"})

    if code != 2 {
        t.Fatalf("exit = %d, want 2 for a usage error", code)
    }
    if !strings.Contains(out.String(), "Usage:") {
        t.Error("a usage error should print the usage")
    }
}
```

#bridge[
  Compare this with testing a CLI built on package-level variables and `init()`
  registration, which is the shape most tutorials teach.

  There, every test shares one root, flags keep values between tests, and the
  only way to assert on output is to swap `os.Stdout`. The design in this
  chapter is not cleverer; it just passes its dependencies in, which is the
  same argument chapter ten made about interfaces.
]

#project(title: "Project 45: build the framework, then a real tool with it")[
  In `part6-from-scratch/49-cli`, build the framework above as a package, and
  then a tool that uses it.

  *One: the framework.* `Command`, `Add`, `find`, `Execute`, `PrintUsage`,
  `Main`, suggestions, signals. No global state.

  *Two: the tool.* A `notes` CLI with real subcommands:

  ```
  notes add "buy milk" --tag shopping
  notes list --tag shopping --json
  notes done 3
  notes export > notes.json
  ```

  Store them in a JSON file under the user's config directory, found with
  `os.UserConfigDir`.

  *Three: the exit codes.* Verify with a script:

  ```bash
  notes list --nope;      echo "usage:   $?"   # 2
  notes done 999;         echo "failure: $?"   # 1
  notes list > /dev/null; echo "ok:      $?"   # 0
  ```

  *Four: the pipe.* `notes list --json | jq '.[0].text'` must work, which means
  no progress output on stdout and no colour codes when piped.

  *Five: cancellation.* Add `notes import <big-file>` that reports progress and
  honours `ctx`. Interrupt it halfway and confirm the file it was writing is
  not left half written. Then interrupt twice and confirm it exits immediately.

  *Six: the typo.* `notes lst` suggests `list`. `notes xyzzy` suggests nothing
  and exits 2.

  *Seven: tests.* Every exit code, the help output, the suggestion, and one
  test that runs the whole tree against a temporary config directory.
]

#tryit[
  Run `--help` on three tools you use daily and read them as a designer rather
  than a user.

  Which one tells you what to do next when you get it wrong? Which one buries
  the useful line? The difference is almost never the parser; it is the three
  decisions at the top of this chapter.
]

#assignment(name: "CLI Framework - Assignment One")[
  Build the framework and the `notes` tool, and make it something you would
  actually install.

  *One.* The framework as its own package with no global state, and tests that
  build a command tree per test.

  *Two.* The tool, with at least four subcommands, one of them nested two
  levels deep.

  *Three.* `CLI.md` recording the exit code for every failure your tool can
  produce, and which stream each message goes to.

  *Four.* Proof it composes: one shell pipeline using your tool with `jq`,
  `grep` and a redirect, pasted with its output.

  *Five.* Cross compile it for three platforms, as in chapter twenty five, and
  install it on your own machine. Use it for a week before you decide it is
  finished, and note the first thing that annoyed you.
]
