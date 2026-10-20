#import "../template.typ": *

= Build your own JSON parser

JSON is the smallest grammar worth parsing properly. It has five kinds of
value, no ambiguity, and a specification you can read in an afternoon. Once you
have written a parser for it, every other parser you meet, a config format, a
query language, a template engine, is the same shape with more rules.

We are building a lexer, a recursive descent parser, errors that name a line
and a column, and the defences a parser needs when the input comes from
strangers.

== The grammar, in full

```
value   := object | array | string | number | true | false | null
object  := '{' (string ':' value (',' string ':' value)*)? '}'
array   := '[' (value (',' value)*)? ']'
```

That is the whole language. Strings and numbers have their own rules, which is
where most of the work turns out to be.

#bridge[
  `JSON.parse` in the browser is roughly this, plus a reviver function and some
  extremely optimised string handling.

  The one behavioural difference worth knowing before we start: JavaScript
  numbers are float64, so `JSON.parse('{"id": 9007199254740993}')` gives you
  9007199254740992. The parser is not wrong, the type is. We will hit the same
  decision later in this chapter, and in Go we get to choose differently.
]

== Tokens

```go
type Kind int

const (
    EOF Kind = iota
    LBrace; RBrace; LBracket; RBracket
    Comma; Colon
    String; Number; True; False; Null
)

type Token struct {
    Kind  Kind
    Text  string   // the decoded string, or the raw number
    Line  int
    Col   int
}
```

The position goes on the token, not on the parser, because the parser reports
errors about tokens it has already read. Retrofitting positions later is
miserable, and every parser author does it once.

== The lexer

```go
type Lexer struct {
    src       []byte
    pos       int
    line, col int
}

func (l *Lexer) Next() (Token, error) {
    l.skipSpace()
    if l.pos >= len(l.src) {
        return Token{Kind: EOF, Line: l.line, Col: l.col}, nil
    }

    start := l.mark()
    switch c := l.src[l.pos]; {
    case c == '{':
        l.advance()
        return start.token(LBrace, "{"), nil
    case c == '"':
        return l.lexString(start)
    case c == '-' || (c >= '0' && c <= '9'):
        return l.lexNumber(start)
    case c == 't' || c == 'f' || c == 'n':
        return l.lexLiteral(start)
    default:
        return Token{}, l.errorf(start,
            "unexpected character %q", c)
    }
}
```

The dull part is `advance`, and it is the part to get right first:

```go
func (l *Lexer) advance() byte {
    c := l.src[l.pos]
    l.pos++
    if c == '\n' {
        l.line++
        l.col = 1
    } else {
        l.col++
    }
    return c
}
```

Every byte goes through it, so the line and column are always correct, and no
later code has to remember to update them.

== Strings, where the work is

```go
func (l *Lexer) lexString(start mark) (Token, error) {
    l.advance()                   // the opening quote
    var sb strings.Builder

    for {
        if l.pos >= len(l.src) {
            return Token{}, l.errorf(start, "unterminated string")
        }
        c := l.advance()
        switch {
        case c == '"':
            return start.token(String, sb.String()), nil
        case c == '\\':
            if err := l.escape(&sb); err != nil {
                return Token{}, err
            }
        case c < 0x20:
            return Token{}, l.errorf(start,
                "control character %#02x in string", c)
        default:
            sb.WriteByte(c)
        }
    }
}
```

Three correctness points that separate a parser from a toy.

*Raw control characters are illegal* inside a JSON string. A literal newline in
a string is a malformed document, and accepting it means your parser disagrees
with everybody else's.

*The escapes are a fixed list*: `" \ / b f n r t` and `uXXXX`. Anything else is
an error, not a passthrough.

*Surrogate pairs are real.* `😀` is one emoji expressed as two
escapes, and handling each in isolation produces two broken runes:

```go
if utf16.IsSurrogate(r) {
    r2, err := l.readEscapeU()   // must be the low half
    if err != nil {
        return err
    }
    sb.WriteRune(utf16.DecodeRune(r, r2))
    return nil
}
```

#gotcha[
  `utf16.DecodeRune` returns `U+FFFD` for a pair that does not combine, which
  means a lone high surrogate becomes a replacement character rather than an
  error.

  Decide deliberately which you want. Being strict is right for a validator and
  wrong for a parser that has to read whatever a legacy system emits, and
  either way it should be a decision with a comment and not an accident.
]

== Numbers, and the decision JavaScript cannot make

```go
func (l *Lexer) lexNumber(start mark) (Token, error) {
    begin := l.pos
    l.accept("-")
    l.acceptDigits()
    if l.accept(".") {
        l.acceptDigits()
    }
    if l.accept("eE") {
        l.accept("+-")
        l.acceptDigits()
    }
    raw := string(l.src[begin:l.pos])
    if !validNumber(raw) {
        return Token{}, l.errorf(start, "invalid number %q", raw)
    }
    return start.token(Number, raw), nil
}
```

Note what it does *not* do: it does not convert. The token carries the raw
text, and the decision about what Go type to produce belongs to the caller.

```go
type Options struct {
    // UseNumber keeps numbers as their original text rather than
    // converting to float64.
    UseNumber bool
}
```

#gotcha[
  This is not a detail. A 64-bit ID, a currency amount in minor units, a
  nanosecond timestamp: all of them exceed the 53 bits a float64 can hold
  exactly, and converting silently changes the value.

  `9007199254740993` becomes `9007199254740992`, and no error is produced
  anywhere. If your parser is reading somebody's API responses, keeping the raw
  text and converting on demand is the only safe default. This is the same
  lesson as chapter forty's money, arrived at from the parsing side.
]

== The parser

Recursive descent means one function per production, and the structure reads
like the grammar:

```go
func (p *Parser) parseValue() (any, error) {
    switch p.tok.Kind {
    case LBrace:
        return p.parseObject()
    case LBracket:
        return p.parseArray()
    case String:
        v := p.tok.Text
        p.next()
        return v, nil
    case Number:
        return p.number()
    case True, False:
        v := p.tok.Kind == True
        p.next()
        return v, nil
    case Null:
        p.next()
        return nil, nil
    default:
        return nil, p.errorf("expected a value, found %s",
            p.tok.Kind)
    }
}
```

An object is a loop over pairs:

```go
func (p *Parser) parseObject() (map[string]any, error) {
    p.next()                       // past '{'
    out := map[string]any{}

    if p.tok.Kind == RBrace {
        p.next()
        return out, nil            // {}
    }
    for {
        if p.tok.Kind != String {
            return nil, p.errorf("expected a key, found %s",
                p.tok.Kind)
        }
        key := p.tok.Text
        p.next()

        if err := p.expect(Colon); err != nil {
            return nil, err
        }
        val, err := p.parseValue()
        if err != nil {
            return nil, err
        }
        out[key] = val

        if p.tok.Kind == Comma {
            p.next()
            continue
        }
        return out, p.expect(RBrace)
    }
}
```

Arrays are the same shape with no keys, and that is the parser.

== Two defences against hostile input

A parser that will be fed data from outside needs both of these, and neither is
obvious until you have been bitten.

*A depth limit.* `[[[[[[...]]]]]]` with ten thousand levels recurses ten
thousand times, and a deep enough document is a stack overflow, which in Go is
a crash you cannot recover from.

```go
func (p *Parser) enter() error {
    p.depth++
    if p.depth > p.maxDepth {       // 10000 is generous
        return p.errorf("nesting too deep")
    }
    return nil
}
```

*A size limit, upstream.* Not the parser's job, but the parser is where the
cost appears: a hundred megabyte document becomes several hundred megabytes of
`map[string]any`. Cap the reader with `io.LimitReader` before parsing, and
remember chapter forty two's decompression bomb is the same family of attack.

#gotcha[
  Duplicate keys are where the specification shrugs: `{"a":1,"a":2}` is legal
  and the result is "implementation defined".

  Most parsers keep the last. Some keep the first. A few error. The attack this
  enables is real: a document parsed by two systems that disagree is a way to
  show one value to a validator and another to a consumer.

  Pick last-wins to match everybody else, and offer a strict mode that
  errors. Then write it down, because the next person will assume whichever one
  you did not choose.
]

== Errors people can act on

```
parse error at line 7, column 24: expected ',' or '}', found string
   7 | { "name": "Ada" "role": "engineer" }
     |                 ^
```

Three pieces: the position, what was expected, and what was found. The caret
line costs about fifteen lines of code and is the difference between a parser
people debug with and one they swear at.

== Testing it properly

*Table tests* for the cases you thought of. Then two techniques that find the
ones you did not.

*Differential fuzzing* against the standard library:

```go
func FuzzAgainstStdlib(f *testing.F) {
    f.Add([]byte(`{"a":[1,2,{"b":null}]}`))

    f.Fuzz(func(t *testing.T, data []byte) {
        var want any
        wantErr := json.Unmarshal(data, &want)

        got, gotErr := Parse(data)

        if (wantErr == nil) != (gotErr == nil) {
            t.Fatalf("stdlib err=%v, ours err=%v", wantErr, gotErr)
        }
        if wantErr == nil && !reflect.DeepEqual(want, got) {
            t.Fatalf("mismatch:\n want %#v\n got  %#v", want, got)
        }
    })
}
```

```bash
go test -fuzz FuzzAgainstStdlib -fuzztime 60s
```

Sixty seconds of that finds the cases a human table never contains: a lone
surrogate, a number like `-0`, `1e999`, a string of nothing but escapes.

*The JSON test suite.* There is a well-known corpus of files named `y_` (must
accept), `n_` (must reject) and `i_` (implementation defined). Running against
it is how you find out your parser is not as correct as it feels.

#bridge[
  Fuzzing is built into `go test`, which is unusual and underused. In most
  languages it is a separate tool with its own configuration, so people do not
  bother.

  Here it is three lines and a flag, and for anything that parses untrusted
  input it finds more real bugs per minute than any other technique available
  to you.
]

== How fast is it

Benchmark against `encoding/json` and report the number honestly. A clean
hand-written parser usually lands in the same order of magnitude for parsing
into `any`, and loses badly to the standard library's struct decoding, which
has years of optimisation in it.

That is the right outcome. You did not build this to replace the standard
library. You built it so that the next time you read a stack trace through
`encoding/json`, you know what each frame is doing.

#project(title: "Project 46: a parser that passes the suite")[
  In `part6-from-scratch/50-json`, build it completely.

  *One: lexer and parser*, producing `any`, with positions on every token.

  *Two: errors* with line, column, expected and found, plus the caret line.

  *Three: the number decision.* Implement `UseNumber`, and write a test with
  `9007199254740993` proving the default loses it and the option keeps it.

  *Four: surrogates.* A test with an emoji written as a surrogate pair, one
  with a lone high surrogate, and your documented behaviour for each.

  *Five: defences.* A depth limit with a test that would otherwise overflow the
  stack, and a documented size limit.

  *Six: the fuzzer.* Run the differential fuzz for five minutes. Every
  difference it finds is a bug in your parser until proven otherwise. Record
  the first three.

  *Seven: the suite.* Run the standard conformance corpus and record your pass
  rate for `y_`, `n_` and `i_` separately. Getting `n_` high is harder and more
  interesting than `y_`.

  *Eight: the benchmark.* Yours against `encoding/json` on a 1 MB document, in
  `RESULTS.md`, with one sentence on what you would optimise first.
]

#tryit[
  Feed your parser a 50 MB JSON file and watch the memory with
  `/usr/bin/time -v` or the Task Manager.

  The ratio between the file size and the resident memory is the real cost of
  `map[string]any`, and it is the argument for streaming, which is the next
  thing to build: a decoder that yields tokens rather than constructing the
  whole tree.
]

#assignment(name: "JSON Parser - Assignment One")[
  Build the parser, then prove it against things you did not write.

  *One.* The full parser with positions, errors, the depth limit and the number
  option.

  *Two.* `CONFORMANCE.md`: your pass rate on the standard corpus, broken into
  the three categories, with three specific files you fail and why.

  *Three.* Five minutes of differential fuzzing, with the first three
  differences it found and the fix for each.

  *Four.* The number test, with both values printed.

  *Five.* One paragraph comparing your error messages with
  `encoding/json`'s for the same malformed document. Paste both. If yours is
  not better, say so and say what you would change.
]
