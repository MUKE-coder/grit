import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { TableOfContents } from "@/components/table-of-contents";
import { PlaygroundChallenge } from "@/components/playground-challenge";
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/prerequisites/golang')

const tocItems = [
  { id: 'go-basics', label: 'Go Basics' },
  { id: 'variables-types', label: 'Variables & Types' },
  { id: 'control-flow', label: 'Control Flow' },
  { id: 'functions-errors', label: 'Functions & Error Handling' },
  { id: 'structs-tags', label: 'Structs & Tags' },
  { id: 'slices-maps', label: 'Slices & Maps' },
  { id: 'pointers', label: 'Pointers' },
  { id: 'methods', label: 'Methods' },
  { id: 'interfaces', label: 'Interfaces' },
  { id: 'goroutines-channels', label: 'Goroutines & Channels' },
  { id: 'packages-structure', label: 'Packages & Project Structure' },
  { id: 'env-variables', label: 'Environment Variables' },
  { id: 'gin-framework', label: 'Gin Framework' },
  { id: 'middleware', label: 'Middleware' },
  { id: 'cors', label: 'CORS' },
  { id: 'handlers', label: 'Handlers' },
  { id: 'services', label: 'Services & The Service Pattern' },
  { id: 'gorm-in-depth', label: 'GORM In Depth' },
  { id: 'migrations-seeding', label: 'Migrations & Seeding' },
  { id: 'jwt-auth', label: 'JWT & Authentication' },
  { id: 'rbac-middleware', label: 'RBAC & Middleware' },
  { id: 'important-packages', label: 'Important Packages' },
  { id: 'putting-it-together', label: 'Putting It Together' },
]

export default function GoForGritPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Prerequisites</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Go for Grit Developers</h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                Everything you need to know about Go to work with Grit&apos;s backend.
                This guide assumes you know another language like JavaScript or Python
                and walks you through Go&apos;s key concepts as they apply to building
                full-stack applications with Grit.
              </p>
            </div>

            <TableOfContents items={tocItems} />

            {/* Try it live callout */}
            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-10 flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-foreground/90">Want to practice as you learn?</p>
                <p className="text-xs text-muted-foreground/70 mt-0.5">Try the code examples in our interactive Go Playground.</p>
              </div>
              <Button size="sm" variant="outline" className="shrink-0 text-xs" asChild>
                <Link href="/playground">Open Playground</Link>
              </Button>
            </div>
            {/* ─────────────────────────────────────────────────── */}
            {/* 1. Go Basics */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="go-basics">1. Go Basics</h2>
              <p>
                Go (often called Golang) is a statically typed, compiled language created at Google.
                It compiles to a single binary with no runtime dependencies, starts up in milliseconds,
                and handles concurrency natively. These qualities make it ideal for building API servers.
              </p>
              <p>
                Every Go file belongs to a <code>package</code>. The special package <code>main</code> is
                the entry point for executables. The <code>func main()</code> function inside
                <code>package main</code> is where your program starts. You import other packages
                using the <code>import</code> keyword.
              </p>
              <p>
                Go uses modules for dependency management. You initialize a module
                with <code>go mod init</code> and run your program with <code>go run .</code>
              </p>
            </div>

            <CodeBlock language="go" filename="main.go" code={`package main

import "fmt"

func main() {
    fmt.Println("Hello, Grit!")
}`} />

            <div className="prose-grit mb-10">
              <p>
                Four things in that file are worth naming, because every Go program has them. The
                <code> package</code> line comes first. The <code>import</code> block lists what the
                file uses. <code>func main()</code> is where execution starts. Anything after
                <code> //</code> is a comment. Here is the same program with a little more in it —
                still no variables, which arrive in the next section.
              </p>
            </div>

            <CodeBlock language="go" filename="anatomy.go" code={`package main

// One import here; several are grouped in parentheses.
import "fmt"

/*
   A block comment. Handy for a paragraph, though most
   Go code uses // for everything.
*/

func main() {
    // Print writes exactly what you give it: no spaces, no newline
    fmt.Print("Starting")
    fmt.Print("...")
    fmt.Println("ready")

    // Println puts a space between arguments and ends the line
    fmt.Println("Grit", "API", 2026)

    // Strings are joined with +
    fmt.Println("Hello, " + "Grit" + "!")

    // A raw string literal keeps line breaks exactly as typed
    fmt.Println(\`usage:
  go run .\`)
}`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                The entry point for every Grit backend is <code>apps/api/cmd/server/main.go</code>.
                This file initializes the database connection, sets up middleware, registers routes,
                and starts the Gin HTTP server. You rarely edit it directly -- the code generator
                handles injecting new routes and models automatically.
              </p>
            </div>

            <PlaygroundChallenge
              title="Your First Program"
              description="Write a program from scratch: the package line, the import, and a main function that prints. No variables yet — just the shape of a Go program."
              challenge={`package main

import "fmt"

func main() {
	// Challenge: print a three-line banner
	// 1. Print "=== Grit API ===" on its own line
	// 2. Print "Starting up" on the next line
	// 3. Print "Ready" on the third line
	//    One fmt.Println per line.
	// 4. Bonus: join two strings with + and print the result,
	//    for example "Hello, " + "Grit"

	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

func main() {
	fmt.Println("=== Grit API ===")
	fmt.Println("Starting up")
	fmt.Println("Ready")

	// Strings join with +
	fmt.Println("Hello, " + "Grit")
}`}
            />

            <PlaygroundChallenge
              title="Print, Println and Raw Strings"
              description="Control where the line breaks fall. Println adds a newline and spaces out its arguments; Print does neither, and a raw string keeps exactly what you typed."
              challenge={`package main

import "fmt"

func main() {
	// Challenge: get the output exactly right
	// 1. Use fmt.Print three times to build "Loading... done" on ONE line,
	//    then end the line with fmt.Println()
	// 2. Use a single fmt.Println with three arguments to print
	//    "Grit API 2026" — it inserts the spaces for you
	// 3. Print this two-line block with ONE raw string literal
	//    (backticks, not quotes):
	//      usage:
	//        go run .

	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

func main() {
	// Print adds nothing: no spaces, no newline
	fmt.Print("Loading")
	fmt.Print("...")
	fmt.Print(" done")
	fmt.Println()

	// Println spaces the arguments and ends the line
	fmt.Println("Grit", "API", 2026)

	// A raw string literal keeps its line breaks and ignores escapes
	fmt.Println(\`usage:
  go run .\`)
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 2. Variables & Types */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="variables-types">2. Variables & Types</h2>
              <p>
                Go is statically typed -- every variable has a fixed type determined at compile time.
                You can declare variables with <code>var</code> (explicit) or <code>:=</code> (short
                assignment, which infers the type). The short form is used inside functions and is by
                far the most common style in Go code.
              </p>
              <p>
                The basic types you will encounter are <code>string</code>, <code>int</code>,
                <code>bool</code>, and <code>float64</code>. Go also has <code>uint</code> (unsigned
                integer, used for database IDs), <code>byte</code>, and <code>rune</code> (for
                Unicode characters). Constants are declared with <code>const</code> and cannot be
                changed after assignment.
              </p>
            </div>

            <CodeBlock language="go" filename="variables.go" code={`package main

import "fmt"

const AppName = "my-saas"

func main() {
    // Explicit declaration
    var name string = "Grit"
    var port int = 8080

    // Short assignment (type inferred)
    host := "localhost"
    debug := true
    price := 29.99

    // Multiple assignment
    width, height := 1920, 1080

    fmt.Println(name, host, port, debug, price, width, height)
    fmt.Println("App:", AppName)
}`} />

            <div className="prose-grit mb-10">
              <p>
                Two things surprise people arriving from JavaScript or Python. Every type has a
                <strong> zero value</strong> — declare a variable without assigning one and it is
                <code>0</code>, <code>&quot;&quot;</code>, <code>false</code> or <code>nil</code>,
                never undefined. And Go never converts implicitly: adding an <code>int</code> to a
                <code>float64</code> is a compile error until you convert one of them yourself.
              </p>
              <p>
                One thing in the next example runs ahead of itself. Turning a string into a number
                can fail, so <code>strconv.Atoi</code> hands back <em>two</em> values: the number and
                an error. Read <code>if err != nil</code> as &quot;if something went wrong&quot; for
                now — section 4 covers errors properly, and this is the only place before then that
                needs them.
              </p>
            </div>

            <CodeBlock language="go" filename="conversion.go" code={`package main

import (
    "fmt"
    "strconv"
)

func main() {
    // Zero values — declared but not assigned
    var count int    // 0
    var name string  // "" (empty, not nil)
    var active bool  // false
    var user *string // nil
    fmt.Printf("%d %q %t %v\n", count, name, active, user)

    // Numeric conversion is always explicit
    total := 10   // int
    price := 2.5  // float64
    // fmt.Println(total * price)        // compile error: mismatched types
    fmt.Println(float64(total) * price)  // 25

    // Integer division truncates — convert BEFORE dividing
    fmt.Println(7 / 2)                   // 3
    fmt.Println(float64(7) / float64(2)) // 3.5

    // Strings are not numbers: strconv returns a value AND an error
    port, err := strconv.Atoi("8080")
    if err != nil {
        fmt.Println("bad port:", err)
        return
    }
    fmt.Println("port + 1 =", port+1)

    // The other direction
    fmt.Println("as string: " + strconv.Itoa(port))
    fmt.Println("as float:  " + strconv.FormatFloat(price, 'f', 2, 64))

    // Something that is not a number gives you an error, not a panic
    if _, err := strconv.Atoi("not-a-port"); err != nil {
        fmt.Println("expected failure:", err)
    }
}`} />


            <div className="prose-grit mb-10">
              <h3 id="format-specifiers">Format Specifiers</h3>
              <p>
                Go&apos;s <code>fmt.Printf</code> and <code>fmt.Sprintf</code> use <strong>format verbs</strong> to
                control how values are printed. You will use these constantly when logging, building strings,
                and debugging. Here are the ones you need to know:
              </p>
            </div>

            <div className="rounded-xl border border-border/40 overflow-hidden mb-8">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/40 bg-accent/10">
                    <th className="text-left py-3 px-4 font-medium text-muted-foreground/80">Specifier</th>
                    <th className="text-left py-3 px-4 font-medium text-muted-foreground/80">Use</th>
                    <th className="text-left py-3 px-4 font-medium text-muted-foreground/80">Example</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/30">
                  <tr>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">%s</code></td>
                    <td className="py-2.5 px-4 text-muted-foreground/80">String</td>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">fmt.Printf(&quot;%s&quot;, &quot;text&quot;)</code></td>
                  </tr>
                  <tr>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">%d</code></td>
                    <td className="py-2.5 px-4 text-muted-foreground/80">Integer</td>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">fmt.Printf(&quot;%d&quot;, 42)</code></td>
                  </tr>
                  <tr>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">%f</code></td>
                    <td className="py-2.5 px-4 text-muted-foreground/80">Float</td>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">fmt.Printf(&quot;%.2f&quot;, 3.14159)</code></td>
                  </tr>
                  <tr>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">%t</code></td>
                    <td className="py-2.5 px-4 text-muted-foreground/80">Boolean</td>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">fmt.Printf(&quot;%t&quot;, true)</code></td>
                  </tr>
                  <tr>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">%v</code></td>
                    <td className="py-2.5 px-4 text-muted-foreground/80">Any value</td>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">fmt.Printf(&quot;%v&quot;, anything)</code></td>
                  </tr>
                  <tr>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">%+v</code></td>
                    <td className="py-2.5 px-4 text-muted-foreground/80">Struct with field names</td>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">fmt.Printf(&quot;%+v&quot;, person)</code></td>
                  </tr>
                  <tr>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">%T</code></td>
                    <td className="py-2.5 px-4 text-muted-foreground/80">Type of value</td>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">fmt.Printf(&quot;%T&quot;, variable)</code></td>
                  </tr>
                  <tr>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">\n</code></td>
                    <td className="py-2.5 px-4 text-muted-foreground/80">Newline</td>
                    <td className="py-2.5 px-4"><code className="text-xs bg-accent/20 px-1.5 py-0.5 rounded">fmt.Printf(&quot;line1\nline2&quot;)</code></td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                You will see <code>:=</code> everywhere in handlers and services. Config values loaded
                from <code>.env</code> are stored in typed struct fields (like <code>Port int</code>,
                <code>JWTSecret string</code>). Constants are used for role names
                (<code>RoleAdmin = &quot;ADMIN&quot;</code>) and error codes.
                Format specifiers are used in error wrapping (<code>fmt.Errorf(&quot;failed to create user: %w&quot;, err)</code>)
                and logging throughout the codebase.
              </p>
            </div>

            <PlaygroundChallenge
              title="Variables & Types"
              description="Declare variables of different types (string, int, float64, bool), convert an int to float64, and print all values with their types."
              challenge={`package main

import "fmt"

func main() {
	// Challenge: Variables & Type Conversion
	// 1. Declare a string variable "name" with any name
	// 2. Declare an int variable "age" with a number
	// 3. Declare a float64 variable "score" with a decimal number
	// 4. Declare a bool variable "passed" set to true
	// 5. Convert age to float64 and add it to score, store in "total"
	// 6. Print each variable with its value and type using %v and %T

	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

func main() {
	age := 25
	name := "Alice"
	score := 95.5
	passed := true

	total := float64(age) + score
	fmt.Printf("Total: %.1f\\n", total)

	fmt.Printf("name: %v (%T)\\n", name, name)
	fmt.Printf("age: %v (%T)\\n", age, age)
	fmt.Printf("score: %v (%T)\\n", score, score)
	fmt.Printf("passed: %v (%T)\\n", passed, passed)
}`}
            />

            <PlaygroundChallenge
              title="Parsing Config Values"
              description="Environment variables always arrive as strings. Convert them to the types you need and handle the failure, which is exactly what config loading does in a real API."
              challenge={`package main

import (
	"fmt"
	"strconv"
)

func main() {
	// These arrive as strings, the way os.Getenv would hand them to you
	rawPort := "8080"
	rawDebug := "true"
	rawRate := "2.5"
	rawBroken := "eight thousand"

	// Challenge: convert each one and print it with its type
	// 1. strconv.Atoi(rawPort)            -> int
	// 2. strconv.ParseBool(rawDebug)      -> bool
	// 3. strconv.ParseFloat(rawRate, 64)  -> float64
	// 4. Each returns (value, error) — check err before using the value
	// 5. Print each with %v and %T so you can see the type you got
	// 6. Convert rawBroken too, and print the error instead of the value

	fmt.Println(rawPort, rawDebug, rawRate, rawBroken)
	fmt.Println("quoted:", strconv.Quote(rawPort))
}`}
              solution={`package main

import (
	"fmt"
	"strconv"
)

func main() {
	rawPort := "8080"
	rawDebug := "true"
	rawRate := "2.5"
	rawBroken := "eight thousand"

	port, err := strconv.Atoi(rawPort)
	if err != nil {
		fmt.Println("bad port:", err)
		return
	}
	fmt.Printf("port  = %v (%T)\n", port, port)

	debug, err := strconv.ParseBool(rawDebug)
	if err != nil {
		fmt.Println("bad debug:", err)
		return
	}
	fmt.Printf("debug = %v (%T)\n", debug, debug)

	rate, err := strconv.ParseFloat(rawRate, 64)
	if err != nil {
		fmt.Println("bad rate:", err)
		return
	}
	fmt.Printf("rate  = %v (%T)\n", rate, rate)

	// The failure path is the point: a bad value is an error, never a panic
	if _, err := strconv.Atoi(rawBroken); err != nil {
		fmt.Println("broken:", err)
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 3. Control Flow */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="control-flow">3. Control Flow</h2>
              <p>
                Go has three keywords for control flow and no more: <code>if</code>,
                <code> for</code> and <code>switch</code>. There is no <code>while</code> — the
                <code> for</code> keyword covers every loop shape — and no ternary operator, so a
                conditional value is written as an ordinary <code>if</code>.
              </p>
              <p>
                Two details catch people out. Parentheses around a condition are not used, but braces
                are always required, even for a single statement. And <code>switch</code> does not
                fall through: each case ends by itself, so there is no <code>break</code> to forget.
              </p>
            </div>

            <CodeBlock language="go" filename="control_flow.go" code={`package main

import "fmt"

func main() {
    // if — no parentheses, braces always required
    port := 8080
    if port < 1024 {
        fmt.Println("privileged port")
    } else if port > 49151 {
        fmt.Println("ephemeral port")
    } else {
        fmt.Println("user port")
    }

    // for — the classic three-part form
    for i := 1; i <= 3; i++ {
        fmt.Println("attempt", i)
    }

    // for as a while loop: one condition, nothing else
    n := 1
    for n < 10 {
        n *= 2
    }
    fmt.Println("n:", n)

    // switch on a value
    env := "staging"
    switch env {
    case "production":
        fmt.Println("be careful")
    case "staging", "qa":
        fmt.Println("safe to experiment")
    default:
        fmt.Println("unknown environment")
    }

    // switch with no value tests conditions instead — often clearer
    // than a chain of else-ifs
    status := 404
    switch {
    case status >= 500:
        fmt.Println("server error")
    case status >= 400:
        fmt.Println("client error")
    default:
        fmt.Println("ok")
    }
}`} />

            <div className="prose-grit mb-10">
              <p>
                <code>break</code> leaves a loop entirely and <code>continue</code> skips to the next
                iteration. An <code>if</code> can also carry a short statement before its condition,
                which is where most Go code puts the variable it is about to test — it keeps that
                variable scoped to the branch that uses it.
              </p>
            </div>

            <CodeBlock language="go" filename="loops.go" code={`package main

import "fmt"

func main() {
    // continue skips the rest of this iteration
    for i := 1; i <= 6; i++ {
        if i%2 != 0 {
            continue // odd numbers are skipped
        }
        fmt.Print(i, " ")
    }
    fmt.Println()

    // break leaves the loop
    total := 0
    for i := 1; ; i++ { // no condition: loops until something breaks it
        total += i
        if total > 20 {
            fmt.Println("stopped at i =", i, "total =", total)
            break
        }
    }

    // A short statement inside if: remainder exists only in these branches
    if remainder := 17 % 5; remainder == 0 {
        fmt.Println("divides evenly")
    } else {
        fmt.Println("remainder is", remainder)
    }

    // Nested loops, and a label to break out of both at once
outer:
    for row := 1; row <= 3; row++ {
        for col := 1; col <= 3; col++ {
            if row*col > 4 {
                fmt.Println("stopping at", row, col)
                break outer
            }
            fmt.Print(row*col, " ")
        }
    }
    fmt.Println()
}`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                The shape you will write most often is the guard clause: an <code>if</code> that
                checks one thing and returns early. Handlers are a stack of them — reject the bad
                request, reject the unauthorised user, then do the work — which keeps the happy path
                at the left margin instead of nested four levels deep.
              </p>
            </div>

            <PlaygroundChallenge
              title="Classify With a Switch"
              description="Loop over a range of numbers and classify each one with a switch. No break statements — Go does not fall through."
              challenge={`package main

import "fmt"

func main() {
	// Challenge: classify the numbers 1 to 15
	// 1. Loop with "for i := 1; i <= 15; i++"
	// 2. Use a switch with no value (switch { case ... }) to print one
	//    line per number:
	//      divisible by 3 and 5 -> "<i> both"
	//      divisible by 3       -> "<i> three"
	//      divisible by 5       -> "<i> five"
	//      otherwise            -> "<i>"
	//    Hint: i%3 == 0 tests divisibility
	// 3. Order matters — check "both" first, or it never matches
	// 4. Bonus: count how many were divisible by neither and print the total

	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

func main() {
	plain := 0

	for i := 1; i <= 15; i++ {
		switch {
		// The most specific case has to come first
		case i%3 == 0 && i%5 == 0:
			fmt.Println(i, "both")
		case i%3 == 0:
			fmt.Println(i, "three")
		case i%5 == 0:
			fmt.Println(i, "five")
		default:
			fmt.Println(i)
			plain++
		}
	}

	fmt.Println("divisible by neither:", plain)
}`}
            />

            <PlaygroundChallenge
              title="break, continue and Guard Clauses"
              description="Skip what you do not want, stop when you have enough, and write the early-return shape that every request handler is built from."
              challenge={`package main

import "fmt"

func main() {
	// Challenge, part one: loop from 1 to 20
	// 1. continue past any number that is not divisible by 4
	// 2. add the rest to a running total
	// 3. break as soon as the total goes above 30, printing where it stopped
	//
	// Challenge, part two: guard clauses
	// 4. Write func describe(age int) string that returns early:
	//      age < 0   -> "invalid"
	//      age < 13  -> "child"
	//      age < 20  -> "teenager"
	//      otherwise -> "adult"
	//    Use four separate returns, not one nested if/else chain.
	// 5. Print describe for -1, 8, 15 and 42

	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

// Guard clauses: handle each rejection and return, so the normal
// answer stays at the bottom instead of nested inside four elses
func describe(age int) string {
	if age < 0 {
		return "invalid"
	}
	if age < 13 {
		return "child"
	}
	if age < 20 {
		return "teenager"
	}
	return "adult"
}

func main() {
	total := 0
	for i := 1; i <= 20; i++ {
		if i%4 != 0 {
			continue // not interested
		}
		total += i
		if total > 30 {
			fmt.Println("stopped at", i, "with total", total)
			break
		}
	}

	// Four plain calls: ranging over a collection arrives with slices
	fmt.Printf("%3d -> %s\n", -1, describe(-1))
	fmt.Printf("%3d -> %s\n", 8, describe(8))
	fmt.Printf("%3d -> %s\n", 15, describe(15))
	fmt.Printf("%3d -> %s\n", 42, describe(42))
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 4. Functions & Error Handling */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="functions-errors">4. Functions & Error Handling</h2>
              <p>
                Go functions can return <strong>multiple values</strong>. This is fundamental
                to Go&apos;s error handling: instead of throwing exceptions, functions return an
                <code>error</code> value as the last return. If the error is <code>nil</code>,
                the operation succeeded. If not, you handle it immediately.
              </p>
              <p>
                The <code>if err != nil</code> pattern appears on nearly every line that calls
                another function. It may look verbose at first, but it makes error flow explicit
                and easy to trace. Use <code>fmt.Errorf(&quot;context: %w&quot;, err)</code> to wrap errors
                with additional context as they bubble up the call stack.
              </p>
            </div>

            <CodeBlock language="go" filename="errors.go" code={`package main

import (
    "errors"
    "fmt"
)

// Functions return (result, error)
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}

func calculateDiscount(price, percent float64) (float64, error) {
    result, err := divide(price * percent, 100)
    if err != nil {
        // Wrap the error with context
        return 0, fmt.Errorf("calculating discount: %w", err)
    }
    return result, nil
}

func main() {
    discount, err := calculateDiscount(100.0, 20.0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Discount:", discount) // 20.0
}`} />

            <div className="prose-grit mb-10">
              <p>
                Wrapping with <code>%w</code> is only half the pattern. The other half is asking what
                an error <em>was</em>, further up the stack. A <strong>sentinel</strong> is a
                package-level error value you compare with <code>errors.Is</code>; a
                package-level error value you compare with <code>errors.Is</code>, and it sees
                through any number of <code>%w</code> wraps. That is what lets a handler map a
                failure raised deep in a service onto the right status code without ever reading the
                message. (Errors that also carry <em>fields</em> need a struct and a method, so they
                wait until after those sections.)
              </p>
            </div>

            <CodeBlock language="go" filename="sentinel_errors.go" code={`package main

import (
    "errors"
    "fmt"
)

// Sentinels: single values, compared by identity rather than by message
var (
    ErrNotFound  = errors.New("record not found")
    ErrForbidden = errors.New("not allowed")
)

func findUser(id int) error {
    if id != 1 {
        // Wrapped, so the caller still finds ErrNotFound underneath
        return fmt.Errorf("findUser %d: %w", id, ErrNotFound)
    }
    return nil
}

func deletePost(role string) error {
    if role != "ADMIN" {
        return fmt.Errorf("deletePost as %s: %w", role, ErrForbidden)
    }
    return nil
}

func main() {
    // errors.Is matches through the wrapping, however deep it goes
    err := findUser(99)
    fmt.Println("error:", err)
    if errors.Is(err, ErrNotFound) {
        fmt.Println("-> respond 404")
    }

    err = deletePost("USER")
    fmt.Println("error:", err)
    if errors.Is(err, ErrForbidden) {
        fmt.Println("-> respond 403")
    }

    // One sentinel never matches another
    fmt.Println("forbidden is not-found?", errors.Is(err, ErrNotFound))

    // The happy path
    fmt.Println("as admin:", deletePost("ADMIN"))
}`} />


            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Every service function in <code>internal/services/</code> returns <code>(result, error)</code>.
                Handlers call services, check for errors, and return the appropriate HTTP response.
                For example, <code>user, err := service.GetUserByID(id)</code> followed by
                an <code>if err != nil</code> block that sends a 404 or 500 JSON response.
              </p>
            </div>

            <PlaygroundChallenge
              title="Functions & Errors"
              description="Write a sqrt function that returns an error for negative numbers. Test it with both positive and negative inputs."
              challenge={`package main

import (
	"fmt"
)

// Challenge: Functions & Error Handling
// 1. Write a function: func sqrt(n float64) (float64, error)
//    - If n is negative, return 0 and an error: "cannot take square root of negative number"
//    - Otherwise return math.Sqrt(n) and nil
//    Hint: use errors.New() to create errors, import "errors" and "math"
// 2. In main, call sqrt(16) and sqrt(-4)
// 3. Handle both cases: print the result on success, print the error on failure
//
// Expected output:
//   sqrt(16) = 4.0
//   Error: cannot take square root of negative number

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"errors"
	"fmt"
	"math"
)

func sqrt(n float64) (float64, error) {
	if n < 0 {
		return 0, errors.New("cannot take square root of negative number")
	}
	return math.Sqrt(n), nil
}

func main() {
	result, err := sqrt(16)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("sqrt(16) = %.1f\\n", result)
	}

	result, err = sqrt(-4)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("sqrt(-4) = %.1f\\n", result)
	}
}`}
            />

            <PlaygroundChallenge
              title="Sentinel Errors"
              description="Define a sentinel error, wrap it with context, then match it with errors.Is — the pattern a handler uses to choose between 409 and 500."
              challenge={`package main

import (
	"errors"
	"fmt"
)

// Challenge: sentinel errors and errors.Is
// 1. Declare a package-level sentinel:
//      var ErrInsufficientStock = errors.New("insufficient stock")
// 2. Write reserve(requested, available int) error that:
//      - wraps ErrInsufficientStock with %w when requested > available,
//        putting both numbers in the message
//      - returns nil otherwise
// 3. In main, call reserve(5, 2) and print the error
// 4. Use errors.Is to detect it and print "-> respond 409"
// 5. Bonus: build an unrelated error and confirm errors.Is says false

func main() {
	fmt.Println(errors.New("start here"))
}`}
              solution={`package main

import (
	"errors"
	"fmt"
)

var ErrInsufficientStock = errors.New("insufficient stock")

func reserve(requested, available int) error {
	if requested > available {
		return fmt.Errorf("reserve %d of %d: %w", requested, available, ErrInsufficientStock)
	}
	return nil
}

func main() {
	err := reserve(5, 2)
	fmt.Println("error:", err)

	// Matches through the wrapping, so nobody has to parse the message
	if errors.Is(err, ErrInsufficientStock) {
		fmt.Println("-> respond 409")
	}

	// Bonus: an unrelated error does not match
	other := fmt.Errorf("db down: %w", errors.New("timeout"))
	fmt.Println("other matches?", errors.Is(other, ErrInsufficientStock))

	// The happy path returns nil
	fmt.Println("reserve(1, 2) =", reserve(1, 2))
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 5. Structs & Tags */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="structs-tags">5. Structs & Tags</h2>
              <p>
                A struct is Go&apos;s way of defining a custom data type -- similar to a class in other
                languages, but without inheritance. Structs group related fields together. Each field
                has a name, a type, and optional <strong>struct tags</strong> (metadata in backtick
                strings after the type).
              </p>
              <p>
                Grit models use three kinds of tags:
              </p>
              <ul>
                <li><strong><code>json:&quot;name&quot;</code></strong> -- controls how the field appears in JSON responses. Use <code>json:&quot;-&quot;</code> to hide a field entirely.</li>
                <li><strong><code>gorm:&quot;...&quot;</code></strong> -- controls the database schema (column type, indexes, constraints, foreign keys).</li>
                <li><strong><code>binding:&quot;required&quot;</code></strong> -- tells Gin to validate incoming request data. If validation fails, Gin returns a 400 error automatically.</li>
              </ul>
            </div>

            <CodeBlock language="go" filename="models/user.go" code={`package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           \`gorm:"primarykey" json:"id"\`
    Name      string         \`gorm:"size:255;not null" json:"name" binding:"required"\`
    Email     string         \`gorm:"size:255;uniqueIndex;not null" json:"email" binding:"required,email"\`
    Password  string         \`gorm:"size:255;not null" json:"-"\`
    Role      string         \`gorm:"size:20;default:USER" json:"role"\`
    Active    bool           \`gorm:"default:true" json:"active"\`
    CreatedAt time.Time      \`json:"created_at"\`
    UpdatedAt time.Time      \`json:"updated_at"\`
    DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\`
}`} />

            <div className="prose-grit mb-10">
              <p>
                That struct is not runnable on its own — it imports GORM. This one is, and it shows
                the tags doing their job: <code>json:&quot;-&quot;</code> keeps the password out of
                every response, and <code>json:&quot;created_at&quot;</code> renames the field on the
                way out. Embedding a struct promotes its fields, which is how a project shares
                <code> ID</code> and timestamps across every model without repeating them.
              </p>
            </div>

            <CodeBlock language="go" filename="struct_json.go" code={`package main

import (
    "encoding/json"
    "fmt"
)

// Embedded into every model — the common fields, declared once
type Base struct {
    ID        uint   \`json:"id"\`
    CreatedAt string \`json:"created_at"\`
}

type User struct {
    Base            // embedded: User gets ID and CreatedAt for free
    Name     string \`json:"name"\`
    Email    string \`json:"email"\`
    Password string \`json:"-"\`                    // never serialised
    Nickname string \`json:"nickname,omitempty"\`   // dropped when empty
}

func main() {
    u := User{
        Base:     Base{ID: 1, CreatedAt: "2026-01-15"},
        Name:     "Ada Lovelace",
        Email:    "ada@example.com",
        Password: "super-secret",
    }

    // Promoted fields are read as if they were declared on User itself
    fmt.Println("id:", u.ID, "created:", u.CreatedAt)

    out, _ := json.MarshalIndent(u, "", "  ")
    fmt.Println(string(out))
    // No "password" key, and no "nickname" because it is empty
}`} />


            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Every model in <code>internal/models/</code> is a struct with these three tag types.
                When you run <code>grit generate resource Product</code>, the CLI creates a struct
                with properly tagged fields, registers it for migration, and generates the matching
                Zod schema and TypeScript type on the frontend.
              </p>
            </div>

            <PlaygroundChallenge
              title="Structs"
              description="Create a Product struct with Name (string), Price (float64), and InStock (bool) fields. Create two products and print them."
              challenge={`package main

import "fmt"

// Challenge: Structs
// 1. Define a Product struct with fields: Name (string), Price (float64), InStock (bool)
// 2. Create two Product instances (e.g. laptop and phone)
// 3. Print each product using %+v to show field names
// 4. Access individual fields: print the name and price of one product
// 5. Check if a product is out of stock using an if statement

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

type Product struct {
	Name    string
	Price   float64
	InStock bool
}

func main() {
	laptop := Product{Name: "Laptop", Price: 999.99, InStock: true}
	phone := Product{Name: "Phone", Price: 699.00, InStock: false}

	fmt.Printf("Product 1: %+v\\n", laptop)
	fmt.Printf("Product 2: %+v\\n", phone)

	fmt.Printf("%s costs $%.2f\\n", laptop.Name, laptop.Price)
	if !phone.InStock {
		fmt.Printf("%s is out of stock!\\n", phone.Name)
	}
}`}
            />

            <PlaygroundChallenge
              title="Struct Tags and Embedding"
              description="Build an Order that embeds a shared Base and uses json tags to rename one field, hide another, and drop an empty one. Then marshal it and read the output."
              challenge={`package main

import (
	"encoding/json"
	"fmt"
)

// Challenge: struct tags and embedding
// 1. Declare a Base struct with:  ID uint  -> json key "id"
// 2. Declare an Order struct that embeds Base and adds:
//      Total        float64 -> json key "total"
//      CustomerID   uint    -> json key "customer_id"
//      InternalNote string  -> never serialised        (json:"-")
//      Coupon       string  -> omitted when empty      (json:"coupon,omitempty")
// 3. Build one Order WITHOUT a coupon, marshal it with json.MarshalIndent
// 4. Print the JSON — there should be no internal note and no coupon key
// 5. Bonus: set a coupon, marshal again, and watch the key appear

func main() {
	var _ = json.Marshal // the import is here ready for step 3
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"encoding/json"
	"fmt"
)

type Base struct {
	ID uint \`json:"id"\`
}

type Order struct {
	Base
	Total        float64 \`json:"total"\`
	CustomerID   uint    \`json:"customer_id"\`
	InternalNote string  \`json:"-"\`
	Coupon       string  \`json:"coupon,omitempty"\`
}

func main() {
	o := Order{
		Base:         Base{ID: 42},
		Total:        99.5,
		CustomerID:   7,
		InternalNote: "flagged for review",
	}

	out, _ := json.MarshalIndent(o, "", "  ")
	fmt.Println(string(out))

	// Bonus: with a coupon set, omitempty stops dropping the key
	o.Coupon = "LAUNCH10"
	out, _ = json.MarshalIndent(o, "", "  ")
	fmt.Println(string(out))
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 6. Slices & Maps */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="slices-maps">6. Slices & Maps</h2>
              <p>
                A <strong>slice</strong> is Go&apos;s dynamic array. Unlike arrays (which have a fixed
                size), slices can grow and shrink. You create them with <code>[]Type{"{}"}</code> or
                <code>make([]Type, length)</code> and add items with <code>append()</code>.
              </p>
              <p>
                A <strong>map</strong> is a key-value data structure (like a JavaScript object or
                Python dictionary). The type <code>map[string]interface{"{}"}</code> (or the modern
                alias <code>map[string]any</code>) can hold any value type -- this is what Gin uses
                for JSON responses.
              </p>
              <p>
                The <code>range</code> keyword iterates over slices and maps, giving you
                both the index/key and value on each iteration.
              </p>
            </div>

            <CodeBlock language="go" filename="collections.go" code={`package main

import "fmt"

func main() {
    // Slices
    names := []string{"Alice", "Bob", "Charlie"}
    names = append(names, "Diana")

    for i, name := range names {
        fmt.Printf("%d: %s\\n", i, name)
    }

    // Maps — every value has the same type here; mixed-type maps need
    // "any", which arrives with interfaces
    user := map[string]string{
        "name":  "Alice",
        "email": "alice@example.com",
    }

    for key, value := range user {
        fmt.Printf("%s = %v\\n", key, value)
    }

    // Access a single value
    fmt.Println("Name:", user["name"])
}`} />

            <div className="prose-grit mb-10">
              <p>
                A slice is a view onto an array: a pointer, a length and a capacity. That is worth
                knowing because it explains the one behaviour that catches everybody — two slices can
                share the same backing array, so writing through one changes the other. Maps have
                their own rule: reading a missing key returns the zero value rather than an error,
                and iteration order is deliberately random, so sort the keys when output has to be
                stable.
              </p>
            </div>

            <CodeBlock language="go" filename="slice_mechanics.go" code={`package main

import (
    "fmt"
    "sort"
)

func main() {
    // len is what is there; cap is how much room before a reallocation
    s := make([]int, 0, 4)
    fmt.Println(len(s), cap(s)) // 0 4

    // Sub-slicing shares the SAME underlying array
    nums := []int{1, 2, 3, 4, 5}
    view := nums[1:3] // [2 3]
    view[0] = 99
    fmt.Println(nums) // [1 99 3 4 5] — nums changed too

    // copy() when you want an independent slice
    safe := make([]int, len(view))
    copy(safe, view)
    safe[0] = 0
    fmt.Println(view, safe) // [99 3] [0 3] — separate now

    // Comma-ok tells "missing" apart from "present but zero"
    stock := map[string]int{"apples": 0}
    n, ok := stock["apples"]
    fmt.Println(n, ok) // 0 true  — present, and genuinely zero
    n, ok = stock["pears"]
    fmt.Println(n, ok) // 0 false — absent

    delete(stock, "apples")
    fmt.Println("size:", len(stock))

    // Map iteration order is random — sort the keys for stable output
    scores := map[string]int{"carol": 9, "alice": 7, "bob": 8}
    keys := make([]string, 0, len(scores))
    for k := range scores {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Printf("%s=%d ", k, scores[k])
    }
    fmt.Println()
}`} />


            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                GORM query results are always slices: <code>var users []models.User</code>. Gin JSON
                responses use <code>gin.H{"{}"}</code> which is just a shortcut for <code>map[string]any</code>.
                For example, <code>c.JSON(200, gin.H{"{"}&quot;data&quot;: users, &quot;message&quot;: &quot;success&quot;{"}"})</code>.
              </p>
            </div>

            <PlaygroundChallenge
              title="Slices & Maps"
              description="Build a word frequency counter: split a sentence into words, count how many times each word appears using a map, and print the results."
              challenge={`package main

import (
	"fmt"
	"strings"
)

// Challenge: Word Frequency Counter
// 1. Write a function: func wordFrequency(sentence string) map[string]int
//    - Split the sentence into words using strings.Fields()
//    - Create a map[string]int to count occurrences
//    - Loop through words, lowercase each with strings.ToLower(), increment count
//    - Return the map
// 2. In main, call it with: "the quick brown fox jumps over the lazy dog the fox"
// 3. Print each word and its count
// 4. Print the total number of unique words using len()

func main() {
	_ = fmt.Sprintf // remove this line when you start
	_ = strings.ToLower
}`}
              solution={`package main

import (
	"fmt"
	"strings"
)

func wordFrequency(sentence string) map[string]int {
	words := strings.Fields(sentence)
	freq := make(map[string]int)
	for _, word := range words {
		freq[strings.ToLower(word)]++
	}
	return freq
}

func main() {
	text := "the quick brown fox jumps over the lazy dog the fox"

	freq := wordFrequency(text)

	for word, count := range freq {
		fmt.Printf("%-10s %d\\n", word, count)
	}

	fmt.Printf("\\nUnique words: %d\\n", len(freq))
}`}
            />

            <PlaygroundChallenge
              title="Grouping With a Map of Slices"
              description="Group records by a key into a map of slices, then sort the keys so the output is identical on every run — the shape of almost every reporting query."
              challenge={`package main

import "fmt"

type Product struct {
	Name     string
	Category string
}

func main() {
	products := []Product{
		{"Laptop", "electronics"},
		{"Desk", "furniture"},
		{"Phone", "electronics"},
		{"Chair", "furniture"},
		{"Cable", "electronics"},
	}

	// Challenge: group the product NAMES by category
	// 1. Build a map[string][]string
	// 2. Loop the products, appending each Name to byCategory[p.Category]
	//    (appending to a missing key works: the zero value is a nil slice)
	// 3. Collect the keys into a []string and sort.Strings them
	// 4. Print one line per category, e.g.
	//      electronics (3): Laptop, Phone, Cable
	//    Hint: strings.Join(names, ", ")

	fmt.Println(products)
}`}
              solution={`package main

import (
	"fmt"
	"sort"
	"strings"
)

type Product struct {
	Name     string
	Category string
}

func main() {
	products := []Product{
		{"Laptop", "electronics"},
		{"Desk", "furniture"},
		{"Phone", "electronics"},
		{"Chair", "furniture"},
		{"Cable", "electronics"},
	}

	// Appending to a missing key is fine — it starts life as a nil slice
	byCategory := map[string][]string{}
	for _, p := range products {
		byCategory[p.Category] = append(byCategory[p.Category], p.Name)
	}

	// Sort the keys, because map iteration order is deliberately random
	categories := make([]string, 0, len(byCategory))
	for c := range byCategory {
		categories = append(categories, c)
	}
	sort.Strings(categories)

	for _, c := range categories {
		names := byCategory[c]
		fmt.Printf("%s (%d): %s\\n", c, len(names), strings.Join(names, ", "))
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 7. Pointers */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="pointers">7. Pointers</h2>
              <p>
                A pointer holds the memory address of a value. Use <code>&amp;</code> to get the
                address of a variable and <code>*</code> to read the value at that address
                (dereference). Pointers let you modify a value in place without copying it, and they
                indicate that a value might be <code>nil</code> (absent).
              </p>
              <p>
                In Go, function arguments are passed by value (copied). If you want a function
                to modify the original value, pass a pointer. This is also why GORM methods take
                pointers to structs: <code>db.Create(&amp;user)</code> writes the new ID back into
                your <code>user</code> variable.
              </p>
            </div>

            <CodeBlock language="go" filename="pointers.go" code={`package main

import "fmt"

func doubleValue(n int) {
    n = n * 2 // Modifies the COPY, not the original
}

func doublePointer(n *int) {
    *n = *n * 2 // Modifies the ORIGINAL via pointer
}

func main() {
    x := 10

    doubleValue(x)
    fmt.Println(x) // Still 10 — the copy was doubled

    doublePointer(&x)
    fmt.Println(x) // Now 20 — modified through pointer

    // Nil pointer: indicates "no value"
    var name *string = nil
    if name == nil {
        fmt.Println("Name is not set")
    }
}`} />

            <div className="prose-grit mb-10">
              <p>
                The place pointers stop being academic is <code>for ... range</code>. The loop
                variable is a <em>copy</em> of the element, so assigning to it changes nothing — a
                bug that produces no error and no output, just a slice that stubbornly refuses to
                update. Reach for the index instead. The same copying rule is why the next section&apos;s
                methods need a pointer receiver whenever they change anything.
              </p>
            </div>

            <CodeBlock language="go" filename="pointer_gotchas.go" code={`package main

import "fmt"

type Product struct {
    Name  string
    Price float64
}

// Takes a copy: this discount goes nowhere
func discountBroken(p Product, pct float64) {
    p.Price = p.Price * (1 - pct/100)
}

// Takes a pointer: changes the original
func discount(p *Product, pct float64) {
    p.Price = p.Price * (1 - pct/100)
}

func main() {
    items := []Product{
        {Name: "Laptop", Price: 1000},
        {Name: "Mouse", Price: 50},
    }

    // WRONG: item is a copy of the element
    for _, item := range items {
        item.Price = 0
    }
    fmt.Println("after range-copy:", items) // unchanged

    // RIGHT: address the element through its index
    for i := range items {
        discount(&items[i], 10)
    }
    fmt.Println("after discount:  ", items)

    // Passing a copy silently does nothing
    discountBroken(items[0], 50)
    fmt.Println("after broken:    ", items)

    // Pointers also let you say "no value" — but check before dereferencing
    var missing *Product
    fmt.Println("missing == nil?", missing == nil)
    if missing != nil {
        fmt.Println(missing.Name) // would panic if reached with nil
    }

    // A pointer into the slice: one element, shared
    first := &items[0]
    first.Price = 1.23
    fmt.Println("via pointer:     ", items)
}`} />


            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                GORM uses pointers for nullable database fields. A regular <code>string</code> defaults
                to <code>&quot;&quot;</code> (empty), but <code>*string</code> can be <code>nil</code> -- meaning the
                database column is NULL. You will see <code>*time.Time</code> for optional timestamps
                like <code>EmailVerifiedAt</code> and <code>gorm.DeletedAt</code> for soft deletes.
                All GORM operations take pointers: <code>db.Create(&amp;user)</code>, <code>db.First(&amp;user, id)</code>.
              </p>
            </div>

            <PlaygroundChallenge
              title="Pointers"
              description="Write a tripleValue function that uses a pointer to modify the original variable, and a swap function that swaps two integers using pointers."
              challenge={`package main

import "fmt"

// Challenge: Pointers
// 1. Write a function: func tripleValue(n *int)
//    - It takes a pointer to int and multiplies the value by 3
//    - Use *n to dereference (read/write the value the pointer points to)
// 2. Write a function: func swap(a, b *int)
//    - Swap the values using: *a, *b = *b, *a
// 3. In main:
//    - Create x := 10, call tripleValue(&x), print x (should be 30)
//    - Create a, b := 5, 15, call swap(&a, &b), print (should be a=15, b=5)

func main() {
	_ = fmt.Sprintf // remove this line when you start
}`}
              solution={`package main

import "fmt"

func tripleValue(n *int) {
	*n = *n * 3
}

func swap(a, b *int) {
	*a, *b = *b, *a
}

func main() {
	x := 10
	fmt.Println("Before triple:", x)

	tripleValue(&x)
	fmt.Println("After triple:", x)

	a, b := 5, 15
	fmt.Printf("Before swap: a=%d, b=%d\\n", a, b)

	swap(&a, &b)
	fmt.Printf("After swap: a=%d, b=%d\\n", a, b)
}`}
            />

            <PlaygroundChallenge
              title="Mutating Through a Pointer"
              description="Fix the classic range-copy bug: a loop that looks like it updates a slice and quietly does nothing. Then write the pointer-receiver method that does work."
              challenge={`package main

import "fmt"

type User struct {
	Name   string
	Visits int
	Active bool
}

func main() {
	users := []User{
		{Name: "Ada", Visits: 0, Active: false},
		{Name: "Grace", Visits: 0, Active: false},
	}

	// Challenge: make these changes actually stick
	// 1. Write a function that takes a POINTER:
	//      func recordVisit(u *User)  -> Visits++ and Active = true
	// 2. Loop with "for i := range users" and call recordVisit(&users[i])
	// 3. Print users and confirm both were updated
	// 4. Now write takesCopy(u User) that does the same thing, call it inside
	//    "for _, u := range users", and show that nothing changes
	// 5. Bonus: write func deactivate(u *User) setting Active = false,
	//    and call it as deactivate(&users[0])

	fmt.Println(users)
}`}
              solution={`package main

import "fmt"

type User struct {
	Name   string
	Visits int
	Active bool
}

// Takes a pointer: operates on the original, not a copy
func recordVisit(u *User) {
	u.Visits++
	u.Active = true
}

// Takes a copy: every change is thrown away when it returns
func takesCopy(u User) {
	u.Visits++
	u.Active = true
}

func deactivate(u *User) {
	u.Active = false
}

func main() {
	users := []User{
		{Name: "Ada"},
		{Name: "Grace"},
	}

	// Indexing gives us the real element to take the address of
	for i := range users {
		recordVisit(&users[i])
	}
	fmt.Println("updated:", users)

	// The range copy: u is a fresh User each iteration, so this is discarded
	fresh := []User{{Name: "Alan"}}
	for _, u := range fresh {
		takesCopy(u)
	}
	fmt.Println("range copy:", fresh) // Visits still 0

	// Bonus: pass the address explicitly
	deactivate(&users[0])
	fmt.Println("deactivated:", users[0])
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 8. Methods */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="methods">8. Methods</h2>
              <p>
                A <strong>method</strong> is a function attached to a type. The difference between
                a function and a method is one thing: <strong>the receiver</strong>. A function
                stands alone, but a method has a receiver parameter before the function name that
                binds it to a specific type.
              </p>
              <p>
                The receiver can be a <strong>value receiver</strong> (<code>func (u User) FullName()</code>)
                or a <strong>pointer receiver</strong> (<code>func (u *User) SetName(name string)</code>).
                Use a pointer receiver when the method needs to modify the struct or when the struct
                is large (to avoid copying). In practice, most methods in Grit use pointer receivers.
              </p>
              <p>
                Methods are how Go achieves object-oriented behavior without classes. Instead of
                <code>class User {"{"}...{"}"}</code>, you define a struct and attach methods to it.
              </p>
            </div>

            <CodeBlock language="go" filename="methods.go" code={`package main

import "fmt"

type User struct {
    FirstName string
    LastName  string
    Email     string
}

// A regular function — takes User as an argument
func getFullName(u User) string {
    return u.FirstName + " " + u.LastName
}

// A method — attached to User with a value receiver
// Use value receiver when you only READ the struct
func (u User) FullName() string {
    return u.FirstName + " " + u.LastName
}

// A method with a pointer receiver
// Use pointer receiver when you MODIFY the struct
func (u *User) SetEmail(email string) {
    u.Email = email // Modifies the original, not a copy
}

func main() {
    user := User{FirstName: "John", LastName: "Doe"}

    // Calling a function — pass the struct as argument
    fmt.Println(getFullName(user)) // "John Doe"

    // Calling a method — use dot notation on the struct
    fmt.Println(user.FullName()) // "John Doe"

    // Pointer receiver method modifies the original
    user.SetEmail("john@example.com")
    fmt.Println(user.Email) // "john@example.com"
}`} />

            <div className="prose-grit mb-10">
              <p>
                Methods are not limited to structs. You can attach them to any type you declare in
                your own package, including one built on <code>string</code> or a slice. That is how
                a bare string becomes a <code>Role</code> that knows what it is allowed to do, and it
                is how <code>String()</code> works: implement that one method and every
                <code> fmt</code> function starts printing your type the way you want.
              </p>
            </div>

            <CodeBlock language="go" filename="named_types.go" code={`package main

import (
    "fmt"
    "strings"
)

// A named type built on string — now it can carry behaviour
type Role string

const (
    RoleAdmin  Role = "ADMIN"
    RoleEditor Role = "EDITOR"
    RoleUser   Role = "USER"
)

// Methods on a named string type
func (r Role) CanPublish() bool {
    return r == RoleAdmin || r == RoleEditor
}

func (r Role) Label() string {
    lower := strings.ToLower(string(r))
    return strings.ToUpper(lower[:1]) + lower[1:]
}

// A named slice type, with a method that reads it
type Cart []float64

func (c Cart) Total() float64 {
    sum := 0.0
    for _, price := range c {
        sum += price
    }
    return sum
}

// Pointer receiver, because this one replaces the slice
func (c *Cart) Add(price float64) {
    *c = append(*c, price)
}

type Money struct {
    Cents int
}

// String() satisfies fmt.Stringer — fmt calls it for you
func (m Money) String() string {
    return fmt.Sprintf("$%d.%02d", m.Cents/100, m.Cents%100)
}

func main() {
    fmt.Println(RoleEditor.CanPublish(), RoleUser.CanPublish()) // true false
    fmt.Println(RoleAdmin.Label())                              // Admin

    cart := Cart{19.99, 5.00}
    cart.Add(3.50)
    fmt.Printf("%d items, total %.2f\\n", len(cart), cart.Total())

    // No .String() call anywhere — fmt finds it
    fmt.Println("price:", Money{Cents: 2599})
}`} />


            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Methods are the foundation of Grit&apos;s architecture. Services are structs with
                a <code>DB *gorm.DB</code> field, and all their operations are methods:
                <code>func (s *ProductService) GetByID(id uint)</code>. Handlers are the same pattern:
                <code>func (h *AuthHandler) Login(c *gin.Context)</code>. GORM hooks are also methods:
                <code>func (u *User) BeforeCreate(tx *gorm.DB) error</code> runs automatically before
                inserting a user into the database.
              </p>
            </div>

            <PlaygroundChallenge
              title="Methods"
              description="Create a Rectangle struct with Width and Height, then add Area() and Perimeter() methods. Use a pointer receiver to add a Scale() method."
              challenge={`package main

import "fmt"

// Challenge: Methods
// 1. Create a Rectangle struct with Width and Height (float64)
// 2. Add an Area() method (value receiver) — returns Width * Height
// 3. Add a Perimeter() method (value receiver) — returns 2 * (Width + Height)
// 4. Add a Scale(factor float64) method (pointer receiver) — multiplies both dimensions
//    Hint: pointer receiver (*Rectangle) so it modifies the original
// 5. In main: create a rect, print area/perimeter, scale it, print new area
//
// Expected output:
//   Rectangle: {Width:10 Height:5}
//   Area: 50.0
//   Perimeter: 30.0
//   After Scale(2): {Width:20 Height:10}
//   New Area: 200.0

func main() {
	_ = fmt.Sprintf // remove this line when you start
}`}
              solution={`package main

import "fmt"

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (r *Rectangle) Scale(factor float64) {
	r.Width *= factor
	r.Height *= factor
}

func main() {
	rect := Rectangle{Width: 10, Height: 5}

	fmt.Printf("Rectangle: %+v\\n", rect)
	fmt.Printf("Area: %.1f\\n", rect.Area())
	fmt.Printf("Perimeter: %.1f\\n", rect.Perimeter())

	rect.Scale(2)
	fmt.Printf("After Scale(2): %+v\\n", rect)
	fmt.Printf("New Area: %.1f\\n", rect.Area())
}`}
            />

            <PlaygroundChallenge
              title="Stringer and Named Types"
              description="Give a named type its own methods, then implement String() so fmt prints it your way without anyone calling a formatter."
              challenge={`package main

import "fmt"

// Challenge: named types and the Stringer interface
// 1. Declare "type Status string" with constants:
//      StatusDraft Status = "DRAFT", StatusLive Status = "LIVE"
// 2. Add a method: func (s Status) IsPublic() bool  -> true only for StatusLive
// 3. Declare "type Temperature float64"
// 4. Give Temperature a String() method returning e.g. "21.5 C"
//    Hint: fmt.Sprintf("%.1f C", float64(t))
// 5. In main, print IsPublic() for both statuses
// 6. Print a Temperature with fmt.Println — String() should be used automatically

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

type Status string

const (
	StatusDraft Status = "DRAFT"
	StatusLive  Status = "LIVE"
)

func (s Status) IsPublic() bool {
	return s == StatusLive
}

type Temperature float64

// Implementing String() is all it takes — fmt looks for it
func (t Temperature) String() string {
	return fmt.Sprintf("%.1f C", float64(t))
}

func main() {
	fmt.Println("draft public?", StatusDraft.IsPublic())
	fmt.Println("live public? ", StatusLive.IsPublic())

	temp := Temperature(21.5)
	fmt.Println("today:", temp) // String() used automatically
	fmt.Printf("also:  %v\\n", temp)
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 9. Interfaces */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="interfaces">9. Interfaces</h2>
              <p>
                An interface defines a set of method signatures. Any type that implements all those
                methods <strong>automatically</strong> satisfies the interface &mdash; there is no
                <code>implements</code> keyword. This is called <strong>implicit implementation</strong> (or
                structural typing), and it is one of Go&apos;s most powerful features.
              </p>
              <p>
                Think of an interface like a <strong>job posting</strong>: it lists the skills required
                (e.g. &quot;must be able to Drive() and Refuel()&quot;), not <em>who</em> you are.
                Anyone who has those skills qualifies &mdash; whether it&apos;s a human or a robot.
              </p>
              <p>
                Interfaces enable polymorphism and are essential for testing. You can swap a real
                database service for a mock that implements the same interface, making unit tests
                fast and isolated.
              </p>
            </div>

            {/* 7.1 Basic Interface */}
            <div className="prose-grit mb-6">
              <h3>7.1 Defining &amp; Implementing an Interface</h3>
              <p>
                Define an interface with the <code>type</code> keyword and a list of method signatures.
                Any type whose methods match is considered an implementation &mdash; no explicit declaration needed.
                This is sometimes called <strong>duck typing</strong>: &quot;if it walks like a duck and
                quacks like a duck, then it&apos;s a duck.&quot;
              </p>
            </div>

            <CodeBlock language="go" filename="interfaces.go" code={`package main

import "fmt"

// ---- THE JOB POSTING (Interface) ----
// This is like a job ad that says:
// "We need someone who can Drive() and Refuel()"
type TruckDriver interface {
    Drive() string
    Refuel() string
}

// ---- CANDIDATE 1: John (Struct) ----
// John never said "I am a TruckDriver"
// He just happens to know how to Drive() and Refuel()
type John struct {
    Name string
    Age  int
}

func (j John) Drive() string {
    return j.Name + " is driving the truck!"
}

func (j John) Refuel() string {
    return j.Name + " is refueling the truck!"
}

// ---- CANDIDATE 2: Robot (Struct) ----
// Robot also never said "I am a TruckDriver"
// But it also knows how to Drive() and Refuel()
type Robot struct {
    Model string
}

func (r Robot) Drive() string {
    return "Robot " + r.Model + " is driving the truck!"
}

func (r Robot) Refuel() string {
    return "Robot " + r.Model + " is refueling the truck!"
}

// ---- THE COMPANY (Function that accepts the interface) ----
// The company doesn't care WHO you are.
// It only cares: "Can you Drive() and Refuel()?"
func HireDriver(d TruckDriver) {
    fmt.Println("Hired!")
    fmt.Println("  ", d.Drive())
    fmt.Println("  ", d.Refuel())
    fmt.Println()
}

func main() {
    // John applies — he can Drive() and Refuel() -> HIRED
    john := John{Name: "John", Age: 35}
    fmt.Println("John applies for the job:")
    HireDriver(john)

    // Robot applies — it can Drive() and Refuel() -> HIRED
    robot := Robot{Model: "TX-500"}
    fmt.Println("Robot applies for the job:")
    HireDriver(robot)
}`} />

            <div className="prose-grit mb-6 mt-8">
              <p>
                Notice that neither <code>John</code> nor <code>Robot</code> declares
                &quot;I am a TruckDriver.&quot; They just <em>have</em> the <code>Drive()</code> and
                <code>Refuel()</code> methods with the right signatures. The compiler checks this
                for you at build time &mdash; if you misspell a method or get the return type wrong,
                you get a compile error, not a runtime crash. The <code>HireDriver</code> function
                doesn&apos;t care about the concrete type &mdash; it only cares that the candidate
                satisfies the <code>TruckDriver</code> contract.
              </p>
              <p>
                Here&apos;s the same pattern in a more real-world context &mdash; a notification
                system where different channels (email, Slack) all satisfy the same <code>Notifier</code> interface:
              </p>
            </div>

            <CodeBlock language="go" filename="notifier.go" code={`package main

import "fmt"

// Notifier — any type that can send a notification
type Notifier interface {
    Send(to string, message string) error
}

// EmailNotifier implements Notifier (implicitly)
type EmailNotifier struct {
    From string
}

func (e *EmailNotifier) Send(to string, message string) error {
    fmt.Printf("Email from %s to %s: %s\\n", e.From, to, message)
    return nil
}

// SlackNotifier also implements Notifier
type SlackNotifier struct {
    Channel string
}

func (s *SlackNotifier) Send(to string, message string) error {
    fmt.Printf("Slack #%s -> %s: %s\\n", s.Channel, to, message)
    return nil
}

// Works with ANY Notifier — email, slack, SMS, webhook...
func alert(n Notifier, user string) {
    n.Send(user, "Your report is ready")
}

func main() {
    email := &EmailNotifier{From: "noreply@app.com"}
    slack := &SlackNotifier{Channel: "alerts"}

    alert(email, "alice@example.com") // Email from noreply@app.com to alice@example.com: Your report is ready
    alert(slack, "alice")             // Slack #alerts -> alice: Your report is ready
}`} />

            {/* 7.2 Why Interfaces Matter */}
            <div className="prose-grit mb-6">
              <h3>7.2 Why Interfaces Matter</h3>
              <p>
                Interfaces solve three real problems:
              </p>
              <ol>
                <li>
                  <strong>Reduce boilerplate</strong> &mdash; Write a function once that works with any
                  type matching the interface. The <code>SaveData</code> function below works with files,
                  network connections, in-memory buffers, and anything else that implements <code>io.Writer</code>.
                </li>
                <li>
                  <strong>Enable testing</strong> &mdash; Swap real services for mocks without changing your
                  business logic. Define a <code>Mailer</code> interface, use the real Resend mailer in
                  production, and a fake one in tests.
                </li>
                <li>
                  <strong>Decouple architecture</strong> &mdash; Your handler layer depends on an interface,
                  not a concrete struct. You can replace the database, switch cloud providers, or refactor
                  internals without touching the code that <em>uses</em> the service.
                </li>
              </ol>
            </div>

            <CodeBlock language="go" filename="why_interfaces.go" code={`package main

import (
    "bytes"
    "fmt"
    "io"
    "os"
)

// SaveData works with ANY io.Writer — files, buffers, HTTP responses, etc.
func SaveData(w io.Writer, data []byte) error {
    _, err := w.Write(data)
    return err
}

func main() {
    data := []byte("Hello, Grit!")

    // Write to a file
    file, _ := os.Create("output.txt")
    SaveData(file, data)
    file.Close()

    // Write to an in-memory buffer (same function!)
    var buf bytes.Buffer
    SaveData(&buf, data)
    fmt.Println(buf.String()) // Hello, Grit!

    // Write to stdout (same function!)
    SaveData(os.Stdout, data) // Hello, Grit!
}`} />

            {/* 7.3 Empty Interface & any */}
            <div className="prose-grit mb-6 mt-8">
              <h3>7.3 The Empty Interface &amp; <code>any</code></h3>
              <p>
                The empty interface <code>interface{'{}'}</code> has zero methods, which means <strong>every
                type satisfies it</strong>. It&apos;s Go&apos;s way of saying &quot;any type at all.&quot;
                Since Go 1.18, you can write <code>any</code> instead &mdash; they are identical.
              </p>
              <p>
                You&apos;ll see empty interfaces in generic data structures, JSON unmarshalling, and
                functions that need to accept truly unpredictable types. But use them sparingly &mdash;
                you lose type safety, so prefer concrete types or named interfaces whenever possible.
              </p>
            </div>

            <CodeBlock language="go" filename="empty_interface.go" code={`package main

import "fmt"

func printAnything(v any) {
    fmt.Printf("Value: %v (type: %T)\\n", v, v)
}

func main() {
    printAnything(42)          // Value: 42 (type: int)
    printAnything("hello")     // Value: hello (type: string)
    printAnything(true)        // Value: true (type: bool)
    printAnything(3.14)        // Value: 3.14 (type: float64)

    // Common in JSON-like data structures
    person := map[string]any{
        "name":  "Alice",
        "age":   30,
        "admin": true,
    }
    fmt.Println(person) // map[admin:true age:30 name:Alice]
}`} />

            {/* 7.4 Type Assertions */}
            <div className="prose-grit mb-6 mt-8">
              <h3>7.4 Type Assertions</h3>
              <p>
                When you have an interface value, you can extract the underlying concrete type
                with a <strong>type assertion</strong>. The syntax is <code>value.(Type)</code>.
                Always use the two-return form <code>val, ok := value.(Type)</code> to avoid
                panics if the assertion fails.
              </p>
            </div>

            <CodeBlock language="go" filename="type_assertions.go" code={`package main

import "fmt"

func describe(v any) {
    // Two-return form — safe, won't panic
    if str, ok := v.(string); ok {
        fmt.Printf("String of length %d: %q\\n", len(str), str)
        return
    }

    if num, ok := v.(int); ok {
        fmt.Printf("Integer: %d (doubled: %d)\\n", num, num*2)
        return
    }

    fmt.Printf("Unknown type: %T\\n", v)
}

func main() {
    describe("hello")  // String of length 5: "hello"
    describe(42)        // Integer: 42 (doubled: 84)
    describe(true)      // Unknown type: bool

    // DANGER: single-return form panics on mismatch!
    // s := someValue.(string) // panics if someValue isn't a string
}`} />

            {/* 7.5 Type Switch */}
            <div className="prose-grit mb-6 mt-8">
              <h3>7.5 Type Switch</h3>
              <p>
                When you need to handle multiple types, a <strong>type switch</strong> is cleaner
                than chaining type assertions. It&apos;s like a regular <code>switch</code> but
                branches on the type of the value using <code>value.(type)</code>.
              </p>
            </div>

            <CodeBlock language="go" filename="type_switch.go" code={`package main

import "fmt"

func process(v any) string {
    switch val := v.(type) {
    case string:
        return fmt.Sprintf("string: %q", val)
    case int:
        return fmt.Sprintf("int: %d", val)
    case bool:
        if val {
            return "bool: yes"
        }
        return "bool: no"
    case []string:
        return fmt.Sprintf("string slice with %d items", len(val))
    default:
        return fmt.Sprintf("unhandled type: %T", val)
    }
}

func main() {
    fmt.Println(process("Go"))              // string: "Go"
    fmt.Println(process(2024))              // int: 2024
    fmt.Println(process(true))              // bool: yes
    fmt.Println(process([]string{"a","b"})) // string slice with 2 items
    fmt.Println(process(3.14))              // unhandled type: float64
}`} />

            {/* 7.6 Common Standard Library Interfaces */}
            <div className="prose-grit mb-6 mt-8">
              <h3>7.6 Common Standard Library Interfaces</h3>
              <p>
                Go&apos;s standard library is built around small, composable interfaces. Learning
                these will make you dramatically more productive:
              </p>
            </div>

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Interface</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Method</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Used For</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-mono text-xs text-foreground/90">fmt.Stringer</td>
                    <td className="px-4 py-2.5 font-mono text-xs">String() string</td>
                    <td className="px-4 py-2.5 text-xs">Custom string representation (like Python&apos;s __str__)</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-mono text-xs text-foreground/90">error</td>
                    <td className="px-4 py-2.5 font-mono text-xs">Error() string</td>
                    <td className="px-4 py-2.5 text-xs">Custom error types with extra context</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-mono text-xs text-foreground/90">io.Reader</td>
                    <td className="px-4 py-2.5 font-mono text-xs">Read(p []byte) (n int, err error)</td>
                    <td className="px-4 py-2.5 text-xs">Reading data &mdash; files, HTTP bodies, buffers</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-mono text-xs text-foreground/90">io.Writer</td>
                    <td className="px-4 py-2.5 font-mono text-xs">Write(p []byte) (n int, err error)</td>
                    <td className="px-4 py-2.5 text-xs">Writing data &mdash; files, HTTP responses, buffers</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-mono text-xs text-foreground/90">io.Closer</td>
                    <td className="px-4 py-2.5 font-mono text-xs">Close() error</td>
                    <td className="px-4 py-2.5 text-xs">Releasing resources &mdash; files, connections, streams</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-mono text-xs text-foreground/90">http.Handler</td>
                    <td className="px-4 py-2.5 font-mono text-xs">ServeHTTP(w, r)</td>
                    <td className="px-4 py-2.5 text-xs">HTTP request handling &mdash; middleware, routers</td>
                  </tr>
                  <tr>
                    <td className="px-4 py-2.5 font-mono text-xs text-foreground/90">sort.Interface</td>
                    <td className="px-4 py-2.5 font-mono text-xs">Len, Less, Swap</td>
                    <td className="px-4 py-2.5 text-xs">Custom sorting for any collection</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <CodeBlock language="go" filename="stringer_error.go" code={`package main

import "fmt"

// Implementing fmt.Stringer — controls how your type prints
type User struct {
    Name string
    Role string
}

func (u User) String() string {
    return fmt.Sprintf("%s (%s)", u.Name, u.Role)
}

// Implementing the error interface — custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

func validateEmail(email string) error {
    if email == "" {
        return &ValidationError{Field: "email", Message: "cannot be empty"}
    }
    return nil
}

func main() {
    user := User{Name: "Alice", Role: "ADMIN"}
    fmt.Println(user) // Alice (ADMIN) — fmt.Stringer in action

    if err := validateEmail(""); err != nil {
        fmt.Println(err) // validation failed on email: cannot be empty
    }
}`} />

            {/* 7.7 Interface Composition */}
            <div className="prose-grit mb-6 mt-8">
              <h3>7.7 Interface Composition</h3>
              <p>
                Go encourages <strong>small, focused interfaces</strong>. You compose larger
                interfaces by embedding smaller ones. This is why the standard library has
                <code>io.Reader</code>, <code>io.Writer</code>, and <code>io.Closer</code> separately,
                then composes them into <code>io.ReadWriter</code>, <code>io.ReadCloser</code>,
                <code>io.WriteCloser</code>, and <code>io.ReadWriteCloser</code>.
              </p>
              <p>
                The Go proverb is: <em>&quot;The bigger the interface, the weaker the abstraction.&quot;</em> Keep
                interfaces small (1-3 methods), and compose when needed.
              </p>
            </div>

            <CodeBlock language="go" filename="composition.go" code={`package main

import "fmt"

// Small, focused interfaces
type Reader interface {
    Read(id string) ([]byte, error)
}

type Writer interface {
    Write(id string, data []byte) error
}

type Deleter interface {
    Delete(id string) error
}

// Compose them into a full storage interface
type Storage interface {
    Reader
    Writer
    Deleter
}

// A function that only needs to read — accepts the smallest interface
func loadConfig(r Reader) ([]byte, error) {
    return r.Read("config.json")
}

// MemoryStore implements all three — so it satisfies Storage
type MemoryStore struct {
    data map[string][]byte
}

func (m *MemoryStore) Read(id string) ([]byte, error) {
    d, ok := m.data[id]
    if !ok {
        return nil, fmt.Errorf("not found: %s", id)
    }
    return d, nil
}

func (m *MemoryStore) Write(id string, data []byte) error {
    m.data[id] = data
    return nil
}

func (m *MemoryStore) Delete(id string) error {
    delete(m.data, id)
    return nil
}

func main() {
    store := &MemoryStore{data: make(map[string][]byte)}
    store.Write("config.json", []byte("port=8080"))

    // Pass the full Storage where only Reader is needed — works fine
    config, _ := loadConfig(store)
    fmt.Println(string(config)) // port=8080
}`} />

            {/* 7.8 Interface Best Practices */}
            <div className="prose-grit mb-6 mt-8">
              <h3>7.8 Best Practices</h3>
            </div>

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Rule</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Why</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">Accept interfaces, return structs</td>
                    <td className="px-4 py-2.5 text-xs">Functions should accept the smallest interface they need, but return concrete types so callers get full functionality</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">Keep interfaces small (1-3 methods)</td>
                    <td className="px-4 py-2.5 text-xs">Small interfaces are easier to implement, mock, and compose. <code>io.Reader</code> has 1 method and is used everywhere</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">Define interfaces where they&apos;re used, not where they&apos;re implemented</td>
                    <td className="px-4 py-2.5 text-xs">The consumer knows what it needs. The implementer doesn&apos;t need to know about every consumer</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 font-medium text-foreground/90">Prefer <code>any</code> over <code>interface{'{}'}</code></td>
                    <td className="px-4 py-2.5 text-xs">Since Go 1.18, <code>any</code> is the idiomatic alias. Use it for readability</td>
                  </tr>
                  <tr>
                    <td className="px-4 py-2.5 font-medium text-foreground/90">Don&apos;t use empty interfaces when you can be specific</td>
                    <td className="px-4 py-2.5 text-xs">Every <code>any</code> you use is type safety you lose. Define a named interface instead</td>
                  </tr>
                </tbody>
              </table>
            </div>

            {/* Quick reference */}
            <div className="prose-grit mb-6">
              <h3>Quick Reference</h3>
            </div>

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Syntax</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">What It Does</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground font-mono text-xs">
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5">type X interface {'{ M() }'}</td>
                    <td className="px-4 py-2.5 font-sans">Define interface X with method M</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5">func (t T) M() {'{ }'}</td>
                    <td className="px-4 py-2.5 font-sans">T implicitly satisfies X (has method M)</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5">val, ok := i.(string)</td>
                    <td className="px-4 py-2.5 font-sans">Type assertion (safe, two-return form)</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5">switch v := i.(type)</td>
                    <td className="px-4 py-2.5 font-sans">Type switch — branch on underlying type</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5">type RW interface {'{ Reader; Writer }'}</td>
                    <td className="px-4 py-2.5 font-sans">Compose interfaces by embedding</td>
                  </tr>
                  <tr>
                    <td className="px-4 py-2.5">any</td>
                    <td className="px-4 py-2.5 font-sans">Alias for interface{'{}'} — accepts any type</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit services follow the interface pattern for testability and flexibility. The mailer
                service, storage service, cache service, and AI service all define interfaces internally.
                For example, you could define a <code>UserService</code> interface with methods
                like <code>GetByID</code>, <code>Create</code>, and <code>Delete</code>, then swap
                in a mock implementation during tests. The <code>routes.Services</code> struct accepts
                these interfaces, making it easy to inject different implementations per environment.
              </p>
            </div>

            <PlaygroundChallenge
              title="Interfaces"
              description="Define a Describable interface with a Describe() string method. Implement it for a Book and a Movie type, then write a function that accepts any Describable."
              challenge={`package main

import "fmt"

// Challenge: Interfaces
// 1. Define a Describable interface with one method: Describe() string
// 2. Create a Book struct with Title and Author (both string)
// 3. Implement Describe() for Book — return "Book: <Title> by <Author>"
// 4. Create a Movie struct with Title and Director (both string)
// 5. Implement Describe() for Movie — return "Movie: <Title> directed by <Director>"
// 6. Write a function: func printDescription(d Describable) that prints d.Describe()
// 7. In main, create a slice of Describable with books and movies, loop and print
//
// Hint: Go interfaces are implicit — no "implements" keyword needed

func main() {
	_ = fmt.Sprintf // remove this line when you start
}`}
              solution={`package main

import "fmt"

type Describable interface {
	Describe() string
}

type Book struct {
	Title  string
	Author string
}

func (b Book) Describe() string {
	return fmt.Sprintf("Book: %s by %s", b.Title, b.Author)
}

type Movie struct {
	Title    string
	Director string
}

func (m Movie) Describe() string {
	return fmt.Sprintf("Movie: %s directed by %s", m.Title, m.Director)
}

func printDescription(d Describable) {
	fmt.Println(d.Describe())
}

func main() {
	items := []Describable{
		Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan"},
		Movie{Title: "The Matrix", Director: "Wachowskis"},
		Book{Title: "Clean Code", Author: "Robert C. Martin"},
	}

	for _, item := range items {
		printDescription(item)
	}
}`}
            />

            <PlaygroundChallenge
              title="Swapping Implementations"
              description="Write one function against an interface and hand it two different implementations. This is the whole argument for interfaces: the caller stops caring which one it got."
              challenge={`package main

import "fmt"

// Challenge: one interface, two implementations
// 1. Declare an interface:
//      type Notifier interface { Send(to, message string) error }
// 2. Implement it twice:
//      type EmailNotifier struct{}  -> prints "email to <to>: <message>"
//      type SMSNotifier struct{}    -> prints "sms to <to>: <message>"
//    Both Send methods return nil.
// 3. Write: func notifyAll(n Notifier, recipients []string, msg string) error
//      - loops the recipients, calls n.Send, returns early on any error
// 4. Call notifyAll twice — once with each implementation — and note that
//    the function never changes
// 5. Bonus: add a FailingNotifier whose Send returns an error, and check
//    that notifyAll stops at the first failure

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"errors"
	"fmt"
)

type Notifier interface {
	Send(to, message string) error
}

type EmailNotifier struct{}

func (EmailNotifier) Send(to, message string) error {
	fmt.Printf("email to %s: %s\\n", to, message)
	return nil
}

type SMSNotifier struct{}

func (SMSNotifier) Send(to, message string) error {
	fmt.Printf("sms to %s: %s\\n", to, message)
	return nil
}

type FailingNotifier struct{}

func (FailingNotifier) Send(to, message string) error {
	return errors.New("transport unavailable")
}

// Written against the interface, so it never needs to change
func notifyAll(n Notifier, recipients []string, msg string) error {
	for _, r := range recipients {
		if err := n.Send(r, msg); err != nil {
			return fmt.Errorf("notify %s: %w", r, err)
		}
	}
	return nil
}

func main() {
	people := []string{"ada@example.com", "grace@example.com"}

	_ = notifyAll(EmailNotifier{}, people, "Build passed")
	_ = notifyAll(SMSNotifier{}, []string{"+15550100"}, "Build passed")

	// Bonus: the first failure stops the loop and is reported with context
	if err := notifyAll(FailingNotifier{}, people, "Build passed"); err != nil {
		fmt.Println("failed:", err)
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 10. Goroutines & Channels */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="goroutines-channels">10. Goroutines & Channels</h2>
              <p>
                Concurrency is one of Go&apos;s most powerful features. Go achieves concurrency through
                two key primitives: <strong>goroutines</strong> and <strong>channels</strong>.
              </p>

              <h3>What is a Goroutine?</h3>
              <p>
                A goroutine is a lightweight thread of execution managed by the Go runtime, not the
                operating system. You start one by putting the <code>go</code> keyword before a
                function call. Goroutines are extremely cheap — you can run <strong>thousands</strong> simultaneously,
                each using only a few kilobytes of memory. This is unlike OS threads which are
                expensive to create and manage.
              </p>
              <p>
                When you call a function normally, it runs <strong>synchronously</strong> — your
                program waits for it to finish before moving to the next line. When you
                prefix it with <code>go</code>, it runs <strong>asynchronously</strong> — execution
                continues immediately while the goroutine runs in the background.
              </p>
            </div>

            <CodeBlock language="go" filename="goroutines-basic.go" code={`package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    fmt.Printf("Hello, %s!\\n", name)
}

func main() {
    // Synchronous — runs and completes before moving on
    sayHello("direct call")

    // Asynchronous — starts a new goroutine
    go sayHello("goroutine")

    // Anonymous goroutine — common pattern
    go func(msg string) {
        fmt.Println(msg)
    }("anonymous goroutine")

    // Without this sleep, main() would exit before
    // the goroutines have a chance to run!
    time.Sleep(100 * time.Millisecond)
    fmt.Println("main done")
}`} />

            <div className="prose-grit mb-10">
              <p>
                <strong>Important:</strong> When the <code>main()</code> function returns, the
                program exits — even if goroutines are still running.
                Using <code>time.Sleep</code> to wait is fragile. In real code, you need proper
                synchronization, which is where <code>sync.WaitGroup</code> and channels come in.
              </p>

              <h3>WaitGroup — Waiting for Goroutines to Finish</h3>
              <p>
                A <code>sync.WaitGroup</code> is a counter that lets you wait for a collection of
                goroutines to finish. Call <code>wg.Add(1)</code> before launching each goroutine,
                call <code>wg.Done()</code> inside the goroutine when it finishes, and
                call <code>wg.Wait()</code> to block until the counter reaches zero.
              </p>
            </div>

            <CodeBlock language="go" filename="waitgroup.go" code={`package main

import (
    "fmt"
    "sync"
    "time"
)

func fetchData(source string, wg *sync.WaitGroup) {
    defer wg.Done() // Signal completion when function returns
    time.Sleep(100 * time.Millisecond) // Simulate work
    fmt.Println("Fetched from:", source)
}

func main() {
    var wg sync.WaitGroup

    sources := []string{"database", "cache", "api"}

    for _, src := range sources {
        wg.Add(1)           // Increment counter
        go fetchData(src, &wg) // Run concurrently
    }

    wg.Wait() // Block until all goroutines call Done()
    fmt.Println("All data fetched!")
    // All 3 goroutines run at the same time (~100ms total, not 300ms)
}`} />

            <div className="prose-grit mb-6 mt-8">
              <h4 className="text-base font-semibold tracking-tight mb-3 text-foreground/80">How the Pieces Connect</h4>
              <p>
                Think of <code>WaitGroup</code> as a simple <strong>counter with a blocking mechanism</strong>:
              </p>
              <ul>
                <li>
                  <code>wg.Add(1)</code> — increments the internal counter. You&apos;re telling the WaitGroup:
                  <em>&quot;one more goroutine is about to start working.&quot;</em> After the loop, the counter is at <strong>3</strong>.
                </li>
                <li>
                  <code>go fetchData(src, &amp;wg)</code> — launches a goroutine. It runs concurrently,
                  meaning it doesn&apos;t block the loop. The loop keeps going and launches all three goroutines almost instantly.
                </li>
                <li>
                  <code>wg.Done()</code> — decrements the counter by 1. Each goroutine calls this
                  (via <code>defer</code>) when it finishes. It&apos;s essentially <code>wg.Add(-1)</code>.
                </li>
                <li>
                  <code>wg.Wait()</code> — blocks the calling goroutine (here, <code>main</code>) until the
                  counter reaches <strong>0</strong>. Once all three goroutines call <code>Done()</code>, the
                  counter hits 0, <code>Wait()</code> unblocks, and <code>main</code> continues.
                </li>
              </ul>
            </div>

            <CodeBlock language="text" filename="Timeline" code={`Time 0ms:
  main:       wg.Add(1), go fetchData("database")  → counter = 1
              wg.Add(1), go fetchData("cache")     → counter = 2
              wg.Add(1), go fetchData("api")       → counter = 3
              wg.Wait()  ← main is now BLOCKED

Time 0-100ms:
  goroutine1: sleeping... (database)
  goroutine2: sleeping... (cache)
  goroutine3: sleeping... (api)

Time ~100ms:
  goroutine1: wg.Done() → counter = 2
  goroutine2: wg.Done() → counter = 1
  goroutine3: wg.Done() → counter = 0  ← Wait() unblocks!

  main:       prints "All data fetched!"`} />

            <div className="prose-grit mb-6 mt-6">
              <h4 className="text-base font-semibold tracking-tight mb-3 text-foreground/80">Two Key Details</h4>
              <p>
                <strong>Why pass <code>&amp;wg</code> (a pointer)?</strong> If you passed <code>wg</code> by value,
                each goroutine would get its own <em>copy</em> of the WaitGroup. Calling <code>Done()</code> on
                a copy wouldn&apos;t decrement the original counter, so <code>Wait()</code> would block forever &mdash;
                a deadlock.
              </p>
              <p>
                <strong>Why <code>defer wg.Done()</code>?</strong> Using <code>defer</code> ensures <code>Done()</code> is
                called even if the function panics. Without it, a panic would leave the counter above 0,
                and <code>Wait()</code> would block forever.
              </p>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">Mental Model</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                <code>Add</code> = &quot;I&apos;m starting work&quot;, <code>Done</code> = &quot;I&apos;m finished&quot;,
                <code>Wait</code> = &quot;hold here until everyone&apos;s finished.&quot; The WaitGroup is just
                the shared scoreboard that makes this coordination possible across goroutines.
              </p>
            </div>

            <div className="prose-grit mb-10">
              <h3>Channels — Communication Between Goroutines</h3>
              <p>
                A <strong>channel is a pipe</strong> that lets goroutines talk to each other:
              </p>
            </div>

            <CodeBlock language="text" filename="How channels work" code={`Goroutine A  ── sends data ──→  [channel]  ──→  receives data ── Goroutine B`} />

            <div className="prose-grit mb-10 mt-6">
              <p>
                You create a channel with <code>make(chan Type)</code>. The two key operations are:
              </p>
              <ul>
                <li>
                  <strong>Send:</strong> <code>{'channel <- value'}</code> — pushes a value into the pipe.
                  The sender <strong>stops and waits</strong> until someone is ready to receive on the other end.
                </li>
                <li>
                  <strong>Receive:</strong> <code>{'value := <-channel'}</code> — pulls a value out of the pipe.
                  The receiver <strong>stops and waits</strong> until someone sends something.
                </li>
              </ul>
              <p>
                This waiting is the magic — it forces goroutines to <strong>synchronize</strong> without
                needing WaitGroups or sleep. By default, channels are <strong>unbuffered</strong> —
                like a phone call where both sides must be on the line at the same time.
              </p>
            </div>

            <CodeBlock language="go" filename="channels.go" code={`package main

import "fmt"

func main() {
    // Create an unbuffered channel of strings
    messages := make(chan string)

    // Launch a goroutine that sends a value
    go func() {
        messages <- "ping" // Send blocks until someone receives
    }()

    // Receive blocks until someone sends
    msg := <-messages
    fmt.Println(msg) // "ping"

    // --- Channel for returning results ---
    results := make(chan int)

    go func() {
        sum := 0
        for i := 1; i <= 100; i++ {
            sum += i
        }
        results <- sum // Send the computed result back
    }()

    total := <-results
    fmt.Println("Sum 1..100 =", total) // 5050
}`} />

            <div className="prose-grit mb-6 mt-8">
              <h4 className="text-base font-semibold tracking-tight mb-3 text-foreground/80">Walking Through the Code</h4>
              <p>
                <strong>Example 1 — Simple message passing:</strong>
              </p>
            </div>

            <CodeBlock language="text" filename="Timeline: message passing" code={`Time 0:
  main:       creates channel "messages"
              launches goroutine
              hits  msg := <-messages  ← BLOCKED (nothing sent yet)

  goroutine:  hits  messages <- "ping" ← finds main is waiting!

  ── handoff happens ──

  main:       msg now equals "ping", prints it
  goroutine:  finishes and exits`} />

            <div className="prose-grit mb-6 mt-6">
              <p>
                The two sides meet at the channel like a <strong>hand-to-hand delivery</strong>. Neither
                side can continue until the other shows up.
              </p>
              <p>
                <strong>Example 2 — Returning a result:</strong>
              </p>
            </div>

            <CodeBlock language="text" filename="Timeline: returning results" code={`  main:       creates channel "results"
              launches goroutine
              hits  total := <-results  ← BLOCKED (waiting for answer)

  goroutine:  calculates sum (1+2+...+100 = 5050)
              hits  results <- 5050     ← delivers the answer

  ── handoff happens ──

  main:       total now equals 5050, prints it`} />

            <div className="prose-grit mb-6 mt-6">
              <p>
                This is the pattern: <strong>send work to a goroutine, get results back through a channel</strong>.
              </p>
            </div>

            {/* Channels vs WaitGroups comparison */}
            <div className="prose-grit mb-6">
              <h4 className="text-base font-semibold tracking-tight mb-3 text-foreground/80">Channels vs WaitGroups</h4>
            </div>

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">WaitGroup</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Channel</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 text-xs">Just says &quot;I&apos;m done&quot;</td>
                    <td className="px-4 py-2.5 text-xs">Sends actual <strong>data</strong> back</td>
                  </tr>
                  <tr className="border-b border-border/20">
                    <td className="px-4 py-2.5 text-xs">Only synchronization</td>
                    <td className="px-4 py-2.5 text-xs">Synchronization <strong>+ communication</strong></td>
                  </tr>
                  <tr>
                    <td className="px-4 py-2.5 text-xs"><code>wg.Done()</code> signals completion</td>
                    <td className="px-4 py-2.5 text-xs">Sending a value signals completion</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">The Go Proverb</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                <em>&quot;Don&apos;t communicate by sharing memory; share memory by communicating.&quot;</em> Channels
                <strong> are</strong> that communication. An unbuffered channel (<code>make(chan string)</code>)
                is like a phone call &mdash; both sides must be on the line at the same time. The sender
                blocks until the receiver is ready, and vice versa.
              </p>
            </div>

            <div className="prose-grit mb-10">
              <h3>Buffered Channels</h3>
              <p>
                By default, channels are unbuffered — a send blocks until a receiver is ready.
                A <strong>buffered channel</strong> has a capacity and can hold values without a
                receiver being ready, up to the buffer size. You create one by passing the
                capacity as the second argument to <code>make</code>.
              </p>
            </div>

            <CodeBlock language="go" filename="buffered-channels.go" code={`package main

import "fmt"

func main() {
    // Buffered channel — can hold up to 2 values
    ch := make(chan string, 2)

    // These sends don't block because the buffer has space
    ch <- "first"
    ch <- "second"

    // Receives pull values out in FIFO order
    fmt.Println(<-ch) // "first"
    fmt.Println(<-ch) // "second"
}`} />

            <div className="prose-grit mb-10">
              <h3>Channel Directions</h3>
              <p>
                When passing channels as function parameters, you can restrict them to
                be <strong>send-only</strong> or <strong>receive-only</strong>. This adds type-safety —
                the compiler prevents you from accidentally reading from a write-only channel or
                vice versa.
              </p>
              <ul>
                <li><code>{'chan<- string'}</code> — send-only channel (can only send strings into it)</li>
                <li><code>{'<-chan string'}</code> — receive-only channel (can only receive strings from it)</li>
                <li><code>chan string</code> — bidirectional channel (can send and receive)</li>
              </ul>
            </div>

            <CodeBlock language="go" filename="channel-directions.go" code={`package main

import "fmt"

// producer can ONLY send to the channel
func producer(ch chan<- string, msg string) {
    ch <- msg
    // <-ch  // This would be a compile error!
}

// consumer can ONLY receive from the channel
func consumer(ch <-chan string) string {
    return <-ch
    // ch <- "x"  // This would be a compile error!
}

func main() {
    ch := make(chan string, 1)
    producer(ch, "hello from producer")
    msg := consumer(ch)
    fmt.Println(msg) // "hello from producer"
}`} />

            <div className="prose-grit mb-10">
              <h3>Ranging Over Channels &amp; Closing</h3>
              <p>
                You can iterate over values received from a channel using <code>for range</code>.
                The loop continues until the channel is <strong>closed</strong>. The
                sender closes a channel with <code>close(ch)</code> to signal that no more values
                will be sent. Closing a channel is important — without it, a <code>range</code> loop
                would block forever waiting for more values.
              </p>
            </div>

            <CodeBlock language="go" filename="range-channels.go" code={`package main

import "fmt"

func generateNumbers(count int, ch chan<- int) {
    for i := 1; i <= count; i++ {
        ch <- i
    }
    close(ch) // Signal: no more values will be sent
}

func main() {
    ch := make(chan int)
    go generateNumbers(5, ch)

    // range automatically stops when channel is closed
    for num := range ch {
        fmt.Printf("Received: %d\\n", num)
    }
    fmt.Println("Channel closed, done!")
    // Output:
    // Received: 1
    // Received: 2
    // Received: 3
    // Received: 4
    // Received: 5
    // Channel closed, done!
}`} />

            <div className="prose-grit mb-10">
              <h3>Select — Waiting on Multiple Channels</h3>
              <p>
                The <code>select</code> statement lets you wait on multiple channel operations
                at once. It&apos;s like a <code>switch</code> statement, but for channels: it blocks
                until one of its cases is ready, then executes that case. If multiple are ready,
                one is chosen at random.
              </p>
              <p>
                <code>select</code> is commonly used for timeouts, cancellation, and multiplexing
                data from multiple sources.
              </p>
            </div>

            <CodeBlock language="go" filename="select.go" code={`package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    // Goroutine 1: slow operation (200ms)
    go func() {
        time.Sleep(200 * time.Millisecond)
        ch1 <- "result from service A"
    }()

    // Goroutine 2: fast operation (100ms)
    go func() {
        time.Sleep(100 * time.Millisecond)
        ch2 <- "result from service B"
    }()

    // Receive results from whichever finishes first
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Got:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Got:", msg2)
        }
    }
    // Output (service B finishes first):
    // Got: result from service B
    // Got: result from service A
}`} />

            <div className="prose-grit mb-10">
              <h3>Practical Pattern: Fan-out / Fan-in</h3>
              <p>
                A common real-world pattern is <strong>fan-out / fan-in</strong>: launch multiple
                goroutines (fan-out), each doing work in parallel, then collect all their results
                through a channel (fan-in). This is the pattern you&apos;ll see in API servers that
                need to fetch data from multiple sources simultaneously.
              </p>
            </div>

            <CodeBlock language="go" filename="fan-out-fan-in.go" code={`package main

import (
    "fmt"
    "time"
)

// Simulates fetching data from different services
func fetch(service string, ch chan<- string) {
    time.Sleep(100 * time.Millisecond) // Simulate network call
    ch <- fmt.Sprintf("data from %s", service)
}

func main() {
    ch := make(chan string)

    // Fan-out: launch 3 goroutines concurrently
    services := []string{"users-api", "orders-api", "payments-api"}
    for _, svc := range services {
        go fetch(svc, ch)
    }

    // Fan-in: collect all results
    for i := 0; i < len(services); i++ {
        result := <-ch
        fmt.Println(result)
    }
    // All 3 fetches run in parallel (~100ms total, not 300ms)
}`} />

            <div className="prose-grit mb-10">
              <h3>Quick Reference</h3>
              <div className="overflow-x-auto">
                <table>
                  <thead>
                    <tr>
                      <th>Operation</th>
                      <th>Syntax</th>
                      <th>Blocks?</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr>
                      <td>Start goroutine</td>
                      <td><code>go func()</code></td>
                      <td>No — runs async</td>
                    </tr>
                    <tr>
                      <td>Create channel</td>
                      <td><code>{'make(chan T)'}</code></td>
                      <td>No</td>
                    </tr>
                    <tr>
                      <td>Create buffered channel</td>
                      <td><code>{'make(chan T, size)'}</code></td>
                      <td>No</td>
                    </tr>
                    <tr>
                      <td>Send to channel</td>
                      <td><code>{'ch <- value'}</code></td>
                      <td>Yes (until receiver ready, or buffer has space)</td>
                    </tr>
                    <tr>
                      <td>Receive from channel</td>
                      <td><code>{'v := <-ch'}</code></td>
                      <td>Yes (until sender sends)</td>
                    </tr>
                    <tr>
                      <td>Close channel</td>
                      <td><code>close(ch)</code></td>
                      <td>No</td>
                    </tr>
                    <tr>
                      <td>Range over channel</td>
                      <td><code>{'for v := range ch'}</code></td>
                      <td>Until channel is closed</td>
                    </tr>
                    <tr>
                      <td>Wait on multiple</td>
                      <td><code>{'select { case ... }'}</code></td>
                      <td>Until one case is ready</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit&apos;s background job system (powered by asynq) uses goroutines under the hood to
                process tasks like sending emails, resizing images, and running cleanup jobs.
                The Gin web server itself handles each HTTP request in its own goroutine — this is
                how it achieves high concurrency without you writing any goroutine code. Pulse&apos;s
                observability tracing also runs in its own goroutines to avoid slowing down your
                API. You generally do not need to write goroutine code directly — asynq, Gin,
                and Pulse manage concurrency for you.
              </p>
            </div>

            <PlaygroundChallenge
              title="Goroutines & Channels"
              description="Create 3 goroutines that each compute a result and send it through a channel. Collect all results in main and print the total."
              challenge={`package main

import "fmt"

// Challenge: Goroutines & Channels
// 1. Write a function: func square(n int, ch chan int)
//    - Compute n*n and send the result to the channel: ch <- n*n
// 2. In main:
//    - Create a channel: ch := make(chan int)
//    - Launch 3 goroutines: go square(3, ch), go square(4, ch), go square(5, ch)
//    - Receive 3 results from the channel in a loop
//    - Sum all results and print the total
//
// Expected: 9 + 16 + 25 = 50 (order of receives may vary)

func main() {
	_ = fmt.Sprintf // remove this line when you start
}`}
              solution={`package main

import "fmt"

func square(n int, ch chan int) {
	ch <- n * n
}

func main() {
	ch := make(chan int)

	go square(3, ch)
	go square(4, ch)
	go square(5, ch)

	total := 0
	for i := 0; i < 3; i++ {
		result := <-ch
		fmt.Printf("Received: %d\\n", result)
		total += result
	}

	fmt.Printf("Sum of squares: %d\\n", total)
}`}
            />

            <PlaygroundChallenge
              title="A Worker Pool"
              description="Fan work out to a fixed number of workers over a channel and collect the results. Sort before printing, because concurrent work never finishes in a predictable order."
              challenge={`package main

import "fmt"

// Challenge: a worker pool
// 1. Make two buffered channels:
//      jobs := make(chan int, 9)     and     results := make(chan string, 9)
// 2. Start 3 workers with "go func(id int) { ... }(w)". Each worker:
//      - ranges over jobs
//      - sends fmt.Sprintf("job %d squared = %d", j, j*j) to results
// 3. Send jobs 1..9 into jobs, then close(jobs)
//    (closing is what lets the workers' range loops end)
// 4. Read exactly 9 values from results into a slice
// 5. sort.Strings the slice and print it, so the output is stable
//    Hint: a sync.WaitGroup is another way to know when the workers are done

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"fmt"
	"sort"
	"sync"
)

func main() {
	const jobCount = 9
	jobs := make(chan int, jobCount)
	results := make(chan string, jobCount)

	var wg sync.WaitGroup
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Ends when jobs is closed and drained
			for j := range jobs {
				results <- fmt.Sprintf("job %d squared = %d", j, j*j)
			}
		}(w)
	}

	for j := 1; j <= jobCount; j++ {
		jobs <- j
	}
	close(jobs) // without this the workers block forever

	wg.Wait()
	close(results)

	out := make([]string, 0, jobCount)
	for r := range results {
		out = append(out, r)
	}

	// Workers finish in whatever order they finish — sort for a stable read
	sort.Strings(out)
	for _, line := range out {
		fmt.Println(line)
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 11. Packages & Project Structure */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="packages-structure">11. Packages & Project Structure</h2>
              <p>
                Go organizes code into packages. Each directory is a package, and the package name
                matches the directory name. A name that starts with an <strong>uppercase letter</strong>
                (like <code>GetUser</code>) is exported (public) -- accessible from other packages.
                A lowercase name (like <code>parseToken</code>) is unexported (private) -- only
                accessible within the same package.
              </p>
              <p>
                The <code>internal/</code> directory is special in Go: packages inside it cannot be
                imported by code outside the parent module. This is a convention enforced by the
                compiler, not just a naming pattern. It keeps your application logic private.
              </p>
            </div>

            <CodeBlock language="bash" filename="project structure" code={`apps/api/
├── cmd/server/
│   └── main.go          # Entry point (package main)
├── cmd/migrate/
│   └── main.go          # Migration CLI (go run cmd/migrate)
├── cmd/seed/
│   └── main.go          # Seeder CLI (go run cmd/seed)
├── internal/
│   ├── config/
│   │   └── config.go    # package config — Config struct, Load()
│   ├── database/
│   │   ├── database.go  # package database — Connect()
│   │   ├── migrate.go   # DropAll() for fresh migrations
│   │   └── seed.go      # Seed() — populate dev data
│   ├── models/
│   │   ├── user.go      # package models — User struct (exported)
│   │   └── upload.go    # package models — Upload struct (exported)
│   ├── handlers/
│   │   ├── auth.go      # package handlers — Login(), Register()
│   │   └── user.go      # package handlers — UserHandler CRUD
│   ├── services/
│   │   └── auth.go      # package services — AuthService (JWT)
│   ├── middleware/
│   │   ├── auth.go      # package middleware — Auth(), RequireRole()
│   │   ├── cors.go      # CORS configuration
│   │   └── logger.go    # Request logging
│   └── routes/
│       └── routes.go    # package routes — Setup() wires everything
└── go.mod               # Module definition`} />

            <div className="prose-grit mb-10">
              <p>
                Go has no <code>public</code> or <code>private</code> keyword. Visibility is decided
                by the first letter: capitalised identifiers are exported from the package,
                lowercase ones are not. That single rule is why service structs expose
                <code> CreateProduct</code> while their helpers stay lowercase — the compiler
                enforces the boundary for you.
              </p>
            </div>

            <CodeBlock language="go" filename="visibility.go" code={`package main

import "fmt"

// Exported — callable from another package as models.Product
type Product struct {
    Name  string  // exported field: appears in JSON, visible everywhere
    Price float64 // exported
    sku   string  // unexported: invisible outside this package
}

// Exported constructor — the usual way to set unexported fields
func NewProduct(name string, price float64, sku string) *Product {
    return &Product{Name: name, Price: price, sku: sku}
}

// Exported method
func (p *Product) SKU() string {
    return p.normalisedSKU()
}

// unexported helper — an implementation detail, free to change
func (p *Product) normalisedSKU() string {
    if p.sku == "" {
        return "UNSET"
    }
    return p.sku
}

// Package-level state and init(), which runs before main()
var registry = map[string]*Product{}

func init() {
    p := NewProduct("Laptop", 999.00, "LAP-1")
    registry[p.SKU()] = p
    fmt.Println("init: registry seeded")
}

func main() {
    fmt.Println("main: registry has", len(registry))

    p := NewProduct("Mouse", 25.00, "")
    fmt.Println(p.Name, p.SKU()) // Mouse UNSET

    // p.sku works here because main is in the same package.
    // From another package it would not compile — that is the whole mechanism.
    fmt.Println("internal sku value:", p.sku)
}`} />

            <PlaygroundChallenge
              title="Exported and Unexported"
              description="Use capitalisation to draw the boundary of a package: an exported constructor and method, with the field and helper behind them kept private."
              challenge={`package main

import "fmt"

// Challenge: visibility by capitalisation
// 1. Declare a struct "Account" with:
//      Owner   string  (exported)
//      balance float64 (unexported — nobody outside may set it directly)
// 2. Write an exported constructor:
//      func NewAccount(owner string, opening float64) *Account
// 3. Write exported methods:
//      func (a *Account) Deposit(amount float64) error  -> reject amounts <= 0
//      func (a *Account) Balance() float64              -> read-only access
// 4. Write an unexported helper: func (a *Account) canWithdraw(n float64) bool
// 5. In main, build an account, deposit twice (once with a bad amount),
//    and print the balance through the method rather than the field

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"errors"
	"fmt"
)

type Account struct {
	Owner   string  // exported
	balance float64 // unexported: only this package can touch it
}

func NewAccount(owner string, opening float64) *Account {
	return &Account{Owner: owner, balance: opening}
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("deposit must be positive")
	}
	a.balance += amount
	return nil
}

func (a *Account) Balance() float64 {
	return a.balance
}

// unexported: an implementation detail, safe to change later
func (a *Account) canWithdraw(n float64) bool {
	return n > 0 && n <= a.balance
}

func main() {
	acct := NewAccount("Ada", 100)

	if err := acct.Deposit(50); err != nil {
		fmt.Println("error:", err)
	}
	if err := acct.Deposit(-5); err != nil {
		fmt.Println("rejected:", err)
	}

	fmt.Printf("%s balance: %.2f\\n", acct.Owner, acct.Balance())
	fmt.Println("can withdraw 120?", acct.canWithdraw(120))
	fmt.Println("can withdraw 200?", acct.canWithdraw(200))
}`}
            />

            <PlaygroundChallenge
              title="init() and Package State"
              description="Package-level variables and init() run before main. Use them to build a lookup table once, which is how registries and default configs get set up."
              challenge={`package main

import "fmt"

// Challenge: package state and init()
// 1. Declare a package-level variable:
//      var statusNames = map[int]string{}
// 2. Write func init() that fills it with:
//      200 "OK", 404 "Not Found", 500 "Internal Server Error"
//    and prints how many entries it added
// 3. Write func describe(code int) string that returns
//      "<code> <name>" if known, or "<code> Unknown" if not
// 4. In main, describe 200, 404 and 418
// 5. Note the ordering: init() output appears BEFORE anything in main

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

// Package-level state — exists before main starts
var statusNames = map[int]string{}

// init() runs automatically, after variable declarations, before main()
func init() {
	statusNames[200] = "OK"
	statusNames[404] = "Not Found"
	statusNames[500] = "Internal Server Error"
	fmt.Println("init: loaded", len(statusNames), "status names")
}

func describe(code int) string {
	name, ok := statusNames[code]
	if !ok {
		return fmt.Sprintf("%d Unknown", code)
	}
	return fmt.Sprintf("%d %s", code, name)
}

func main() {
	fmt.Println("main starts")
	fmt.Println(describe(200))
	fmt.Println(describe(404))
	fmt.Println(describe(418))
}`}
            />


            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit follows Go&apos;s standard project layout exactly. All application code lives
                inside <code>internal/</code>: models, handlers, services, middleware, routes, and
                config. The <code>cmd/</code> directory contains entry points for different commands
                (server, migrate, seed).
                When you import a package, you use the full module path:
                <code>import &quot;my-app/apps/api/internal/models&quot;</code>.
              </p>
            </div>
            {/* ─────────────────────────────────────────────────── */}
            {/* 12. Environment Variables */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="env-variables">12. Environment Variables</h2>
              <p>
                Go reads environment variables with <code>os.Getenv(&quot;KEY&quot;)</code>. For local
                development, you store variables in a <code>.env</code> file and load them
                with the <code>godotenv</code> package. A common pattern is to define a <code>Config</code>
                struct that holds all your settings in one place, loaded once at startup.
              </p>
              <p>
                This pattern keeps configuration centralized, type-safe, and easy to override
                per environment (development, staging, production).
              </p>
            </div>

            <CodeBlock language="go" filename="config/config.go" code={`package config

import (
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

type Config struct {
    Port       int
    DBHost     string
    DBPort     int
    DBName     string
    DBUser     string
    DBPassword string
    JWTSecret  string
    Debug      bool
}

func Load() *Config {
    // Load .env file (ignored in production)
    godotenv.Load()

    port, _ := strconv.Atoi(getEnv("PORT", "8080"))
    dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))

    return &Config{
        Port:       port,
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     dbPort,
        DBName:     getEnv("DB_NAME", "grit_dev"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        JWTSecret:  getEnv("JWT_SECRET", "change-me"),
        Debug:      getEnv("DEBUG", "false") == "true",
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}`} />

            <div className="prose-grit mb-10">
              <p>
                <code>os.Getenv</code> has one weakness: an unset variable and one set to the empty
                string look identical, and both give you <code>&quot;&quot;</code>. That is fine for
                optional values and dangerous for required ones. The fix is two small helpers — one
                that falls back to a default, one that fails loudly — plus typed parsing for
                anything that is not a string.
              </p>
            </div>

            <CodeBlock language="go" filename="env_helpers.go" code={`package main

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

// Optional: fall back when unset or empty
func getEnv(key, fallback string) string {
    if v, ok := os.LookupEnv(key); ok && v != "" {
        return v
    }
    return fallback
}

// Required: fail at start-up rather than at 3am
func mustEnv(key string) (string, error) {
    v, ok := os.LookupEnv(key)
    if !ok || v == "" {
        return "", fmt.Errorf("required env var %s is not set", key)
    }
    return v, nil
}

// Typed, with a fallback when the value is missing or unparseable
func getEnvInt(key string, fallback int) int {
    v, err := strconv.Atoi(os.Getenv(key))
    if err != nil {
        return fallback
    }
    return v
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
    d, err := time.ParseDuration(os.Getenv(key))
    if err != nil {
        return fallback
    }
    return d
}

func main() {
    // Set a couple so the example is self-contained
    os.Setenv("PORT", "9090")
    os.Setenv("TIMEOUT", "45s")
    os.Setenv("EMPTY", "")

    fmt.Println("PORT:      ", getEnvInt("PORT", 8080))       // 9090
    fmt.Println("MAX_CONNS: ", getEnvInt("MAX_CONNS", 25))    // 25 (unset)
    fmt.Println("TIMEOUT:   ", getEnvDuration("TIMEOUT", 30*time.Second))
    fmt.Println("EMPTY:     ", getEnv("EMPTY", "fallback"))   // fallback

    // LookupEnv distinguishes "set to empty" from "not set at all"
    if v, ok := os.LookupEnv("EMPTY"); ok {
        fmt.Printf("EMPTY is set, value = %q\\n", v)
    }

    if _, err := mustEnv("JWT_SECRET"); err != nil {
        fmt.Println("startup error:", err)
    }
}`} />

            <PlaygroundChallenge
              title="Config With Defaults"
              description="Write the two helpers every Go service ends up with: one that falls back to a default, one that refuses to start without a value."
              challenge={`package main

import (
	"fmt"
	"os"
)

// Challenge: environment helpers
// 1. Write func getEnv(key, fallback string) string
//      - use os.LookupEnv so "set but empty" counts as missing
// 2. Write func mustEnv(key string) (string, error)
//      - return an error naming the missing key
// 3. Write func getEnvBool(key string, fallback bool) bool
//      - use strconv.ParseBool, fall back when it errors
// 4. In main: os.Setenv("APP_ENV", "production") and os.Setenv("DEBUG", "false")
// 5. Print getEnv("APP_ENV", "development"), getEnv("REGION", "eu-west-1"),
//    getEnvBool("DEBUG", true), and the error from mustEnv("DATABASE_URL")

func main() {
	fmt.Println(os.Getenv("HOME"))
}`}
              solution={`package main

import (
	"fmt"
	"os"
	"strconv"
)

func getEnv(key, fallback string) string {
	// LookupEnv tells "unset" apart from "set to empty"
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return "", fmt.Errorf("required env var %s is not set", key)
	}
	return v, nil
}

func getEnvBool(key string, fallback bool) bool {
	b, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return b
}

func main() {
	os.Setenv("APP_ENV", "production")
	os.Setenv("DEBUG", "false")

	fmt.Println("env:    ", getEnv("APP_ENV", "development"))
	fmt.Println("region: ", getEnv("REGION", "eu-west-1")) // unset -> fallback
	fmt.Println("debug:  ", getEnvBool("DEBUG", true))     // parsed -> false

	// Required values fail at start-up, where the error is cheap to fix
	if _, err := mustEnv("DATABASE_URL"); err != nil {
		fmt.Println("startup error:", err)
	}
}`}
            />

            <PlaygroundChallenge
              title="Validating Config at Start-up"
              description="Load a whole config struct, collect every problem at once, and report them together — far kinder than failing on one missing variable at a time."
              challenge={`package main

import (
	"fmt"
	"os"
)

type Config struct {
	Port        int
	DatabaseURL string
	JWTSecret   string
	Environment string
}

// Challenge: load and validate config in one pass
// 1. Write func Load() (*Config, []string)
//      - read PORT (default 8080), DATABASE_URL, JWT_SECRET,
//        APP_ENV (default "development")
//      - collect a message into the []string for each problem:
//          * DATABASE_URL missing
//          * JWT_SECRET missing or shorter than 16 characters
//          * PORT set but not a number
// 2. In main, set only APP_ENV and PORT, then call Load()
// 3. If there are problems, print each one and stop; otherwise print the config
// 4. Bonus: set the missing variables and run again to see it succeed

func main() {
	fmt.Println(os.Getenv("PORT"))
}`}
              solution={`package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func Load() (*Config, []string) {
	var problems []string
	cfg := &Config{Port: 8080, Environment: "development"}

	if raw, ok := os.LookupEnv("PORT"); ok && raw != "" {
		p, err := strconv.Atoi(raw)
		if err != nil {
			problems = append(problems, fmt.Sprintf("PORT is not a number: %q", raw))
		} else {
			cfg.Port = p
		}
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		problems = append(problems, "DATABASE_URL is required")
	}

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	switch {
	case cfg.JWTSecret == "":
		problems = append(problems, "JWT_SECRET is required")
	case len(cfg.JWTSecret) < 16:
		problems = append(problems, "JWT_SECRET must be at least 16 characters")
	}

	if v := os.Getenv("APP_ENV"); v != "" {
		cfg.Environment = v
	}
	return cfg, problems
}

func main() {
	os.Setenv("APP_ENV", "production")
	os.Setenv("PORT", "9090")

	cfg, problems := Load()
	if len(problems) > 0 {
		fmt.Println("cannot start:")
		for _, p := range problems {
			fmt.Println("  -", p)
		}
	}

	// Bonus: supply the rest and it loads cleanly
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/app")
	os.Setenv("JWT_SECRET", "a-long-enough-secret")

	cfg, problems = Load()
	fmt.Println("\\nproblems:", len(problems))
	fmt.Printf("%+v\\n", *cfg)
}`}
            />


            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit&apos;s config lives in <code>internal/config/config.go</code>. It loads settings for the
                database, Redis, S3 storage, Resend email, AI keys, Sentinel security, and more -- all from
                the <code>.env</code> file. The <code>Config</code> struct is created once in <code>main.go</code>
                and passed to every service that needs it. A <code>.env.example</code> file is scaffolded
                with every project to document all available variables.
              </p>
            </div>
            {/* ─────────────────────────────────────────────────── */}
            {/* 13. Gin Framework */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="gin-framework">13. Gin Framework</h2>
              <p>
                Gin is Go&apos;s most popular HTTP framework. It provides a fast router, middleware support,
                JSON binding, validation, and route groups. Understanding Gin is essential because
                every handler you write receives a <code>*gin.Context</code> -- the single object that holds
                the request, response, URL parameters, query strings, and more.
              </p>
              <h3>Creating a Server</h3>
              <p>
                You create a Gin engine with <code>gin.New()</code> (bare) or <code>gin.Default()</code>
                (includes logger and recovery middleware). Then you define routes and start the server.
              </p>
            </div>

            <CodeBlock language="go" filename="server.go" code={`package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    // Create a Gin engine (bare, no default middleware)
    r := gin.New()

    // Add middleware globally
    r.Use(gin.Logger())    // Log every request
    r.Use(gin.Recovery())  // Recover from panics

    // Simple route
    r.GET("/api/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
        })
    })

    // Start server on port 8080
    r.Run(":8080")
}`} />

            <div className="prose-grit mb-10">
              <h3>Route Groups & Middleware</h3>
              <p>
                Route groups let you organize related routes under a common prefix and apply
                middleware to all routes in the group at once. This is how Grit separates public
                routes (no auth), protected routes (login required), and admin routes (admin role required).
              </p>
            </div>

            <CodeBlock language="go" filename="route_groups.go" code={`// Public routes — no authentication
auth := r.Group("/api/auth")
{
    auth.POST("/register", authHandler.Register)
    auth.POST("/login", authHandler.Login)
    auth.POST("/refresh", authHandler.Refresh)
}

// Protected routes — requires valid JWT token
protected := r.Group("/api")
protected.Use(middleware.Auth(db, authService))  // Apply auth middleware
{
    protected.GET("/auth/me", authHandler.Me)
    protected.GET("/users/:id", userHandler.GetByID)
}

// Admin routes — requires ADMIN role
admin := r.Group("/api")
admin.Use(middleware.Auth(db, authService))
admin.Use(middleware.RequireRole("ADMIN"))        // Stack middleware
{
    admin.GET("/users", userHandler.List)
    admin.POST("/users", userHandler.Create)
    admin.PUT("/users/:id", userHandler.Update)
    admin.DELETE("/users/:id", userHandler.Delete)
}`} />

            <div className="prose-grit mb-10">
              <h3>The gin.Context Object</h3>
              <p>
                Every handler receives <code>*gin.Context</code>. Here are the methods you will use most:
              </p>
            </div>

            <CodeBlock language="go" filename="gin_context.go" code={`func exampleHandler(c *gin.Context) {
    // ── Reading the request ──────────────────────────────
    id := c.Param("id")                     // URL param: /users/:id
    page := c.Query("page")                 // Query string: ?page=2
    page = c.DefaultQuery("page", "1")      // With default value

    var input CreateUserInput
    err := c.ShouldBindJSON(&input)         // Parse + validate JSON body

    token := c.GetHeader("Authorization")   // Read a header

    // ── Sending responses ────────────────────────────────
    c.JSON(200, gin.H{"data": "hello"})     // Send JSON
    c.JSON(404, gin.H{                      // Send error
        "error": gin.H{
            "code":    "NOT_FOUND",
            "message": "User not found",
        },
    })

    // ── Middleware data ──────────────────────────────────
    c.Set("user_id", uint(42))              // Store data (middleware → handler)
    userID, _ := c.Get("user_id")           // Retrieve data

    // ── Control flow ────────────────────────────────────
    c.Abort()                               // Stop the middleware chain
    c.Next()                                // Continue to next middleware/handler
}`} />

            <div className="prose-grit mb-10">
              <h3>Input Validation with Binding Tags</h3>
              <p>
                Gin uses struct tags to validate incoming JSON. When you call <code>c.ShouldBindJSON(&amp;input)</code>,
                Gin parses the request body, checks the <code>binding</code> tags, and returns an
                error if validation fails. No manual validation code needed.
              </p>
            </div>

            <CodeBlock language="go" filename="validation.go" code={`// Gin validates this struct automatically
type CreateUserInput struct {
    Name     string \`json:"name" binding:"required"\`           // Must be present
    Email    string \`json:"email" binding:"required,email"\`    // Must be valid email
    Password string \`json:"password" binding:"required,min=8"\` // Min 8 characters
    Age      int    \`json:"age" binding:"gte=18,lte=120"\`      // Between 18-120
    Role     string \`json:"role" binding:"oneof=USER EDITOR"\`  // Must be one of these
}

func createUser(c *gin.Context) {
    var input CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        // Gin returns detailed validation errors automatically
        c.JSON(422, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    // input is now validated and safe to use
    fmt.Println(input.Name, input.Email)
}`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                All API routes are defined in <code>internal/routes/routes.go</code>. The <code>Setup()</code>
                function creates a Gin engine, applies global middleware (Logger, Recovery, CORS),
                then organizes routes into groups: public auth, protected, profile, and admin.
                Middleware like <code>Auth()</code> and <code>RequireRole(&quot;ADMIN&quot;)</code> are applied
                per-group. When you generate a new resource, the CLI injects routes into the correct
                group using marker comments.
              </p>
            </div>

            <div className="rounded-xl border border-sky-500/20 bg-sky-500/[0.04] p-5 mb-8">
              <h4 className="text-sm font-semibold text-sky-400 uppercase tracking-wider mb-2">About these challenges</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                The playground compiles against the standard library only, so Gin itself will not
                run there. These challenges use <code>net/http</code> and <code>httptest</code>
                instead, which is what Gin is built on — the router, the handler signature and the
                context are conveniences over exactly this. Everything you practise here is the same
                shape you will write in <code>internal/handlers/</code>, minus the helper methods.
              </p>
            </div>

            <PlaygroundChallenge
              title="Routing and Path Parameters"
              description="Register routes, pull an id out of the path, and return JSON with the right status code. This is what c.Param and c.JSON are doing underneath."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: a tiny router
// 1. Make a mux:  mux := http.NewServeMux()
// 2. Handle "GET /products/{id}" (Go 1.22+ patterns support this)
//      - read the id with r.PathValue("id")
//      - look it up in the products map below
//      - found:   w.WriteHeader(200) and write {"id":1,"name":"Laptop"}
//      - missing: w.WriteHeader(404) and write {"error":"not found"}
//      - always set Content-Type: application/json
// 3. Handle "GET /health" returning 200 and {"status":"ok"}
// 4. Use httptest to call /products/1, /products/99 and /health,
//    printing the status and body of each
//    In Gin this is r.GET("/products/:id", handler) and c.Param("id")

var products = map[string]string{"1": "Laptop", "2": "Mouse"}

func main() {
	fmt.Println(products, http.StatusOK, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

var products = map[string]string{"1": "Laptop", "2": "Mouse"}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func main() {
	mux := http.NewServeMux()

	// Gin: r.GET("/products/:id", ...) with c.Param("id")
	mux.HandleFunc("GET /products/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		name, ok := products[id]
		if !ok {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"id": id, "name": name})
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	for _, path := range []string{"/products/1", "/products/99", "/health"} {
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, httptest.NewRequest("GET", path, nil))
		fmt.Printf("%-14s %d %s", path, rec.Code, rec.Body.String())
	}
}`}
            />

            <PlaygroundChallenge
              title="Binding and Validating JSON"
              description="Decode a request body into a struct, reject what is invalid with 422 and a field-by-field message, and accept what is valid with 201 — the job c.ShouldBindJSON does for you."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type CreateProductRequest struct {
	Name  string  \`json:"name"\`
	Price float64 \`json:"price"\`
}

// Challenge: bind and validate
// 1. Write a handler that decodes the JSON body into CreateProductRequest
//      json.NewDecoder(r.Body).Decode(&req)
// 2. On a decode error, respond 400 with {"error":"invalid JSON"}
// 3. Validate: Name must not be empty, Price must be > 0.
//    Collect problems into map[string]string, e.g. {"name":"is required"}
// 4. If there are problems, respond 422 with {"errors": <that map>}
// 5. Otherwise respond 201 with the created product
// 6. Test it three times with httptest.NewRequest("POST", "/products", strings.NewReader(...)):
//      a valid body, an invalid body, and a body that is not JSON at all
//    In Gin: c.ShouldBindJSON(&req) plus binding:"required" tags

func main() {
	fmt.Println(http.StatusCreated, strings.ToUpper("post"), httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type CreateProductRequest struct {
	Name  string  \`json:"name"\`
	Price float64 \`json:"price"\`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	problems := map[string]string{}
	if strings.TrimSpace(req.Name) == "" {
		problems["name"] = "is required"
	}
	if req.Price <= 0 {
		problems["price"] = "must be greater than zero"
	}
	if len(problems) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": problems})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id": 1, "name": req.Name, "price": req.Price,
	})
}

func main() {
	bodies := []string{
		\`{"name":"Laptop","price":999}\`,
		\`{"name":"","price":0}\`,
		\`not json at all\`,
	}
	for _, b := range bodies {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/products", strings.NewReader(b))
		createProduct(rec, req)
		fmt.Printf("%d %s", rec.Code, rec.Body.String())
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 14. Middleware */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="middleware">14. Middleware</h2>
              <p>
                Middleware is a function that runs <strong>before</strong> (or after) your handler.
                It sits in the request chain and can inspect, modify, or reject requests. Think of it
                as a pipeline: each request passes through a series of middleware functions before
                reaching the handler.
              </p>
              <p>
                In Gin, middleware is a <code>gin.HandlerFunc</code> -- the same type as a handler.
                The difference is that middleware calls <code>c.Next()</code> to pass control to the
                next function in the chain, or <code>c.Abort()</code> to stop the chain entirely
                (e.g., when authentication fails).
              </p>
            </div>

            <CodeBlock language="go" filename="middleware pattern" code={`// A middleware is just a gin.HandlerFunc that calls c.Next()
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next() // ← Run the next handler/middleware

        // This runs AFTER the handler returns
        duration := time.Since(start)
        status := c.Writer.Status()
        log.Printf("%s %s → %d (%v)", c.Request.Method, c.Request.URL.Path, status, duration)
    }
}

// Middleware that blocks requests (c.Abort)
func RequireAPIKey() gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.GetHeader("X-API-Key")
        if key != "valid-key" {
            c.JSON(401, gin.H{"error": "Invalid API key"})
            c.Abort() // ← Stop the chain, handler never runs
            return
        }
        c.Next()
    }
}`} />

            <div className="prose-grit mb-10">
              <h3>The Middleware Chain</h3>
              <p>
                Middleware runs in the order you add it. When a request comes in, it flows through
                each middleware, then the handler, and back out through the middleware in reverse:
              </p>
            </div>

            <CodeBlock language="bash" filename="middleware chain" code={`Request → Logger → CORS → Auth → RequireRole → Handler
                                                         ↓
Response ← Logger ← CORS ← Auth ← RequireRole ← Handler

// If Auth calls c.Abort():
Request → Logger → CORS → Auth ✗ (returns 401, handler never runs)`} />

            <CodeBlock language="go" filename="applying middleware" code={`r := gin.New()

// Global middleware — runs on EVERY request
r.Use(middleware.Logger())
r.Use(gin.Recovery())
r.Use(middleware.CORS(cfg.CORSOrigins))

// Group middleware — runs only on routes in this group
protected := r.Group("/api")
protected.Use(middleware.Auth(db, authService))   // Only protected routes
{
    protected.GET("/users/:id", userHandler.GetByID)
}

// Stacking middleware — multiple on one group
admin := r.Group("/api")
admin.Use(middleware.Auth(db, authService))        // Must be logged in
admin.Use(middleware.RequireRole("ADMIN"))          // AND must be admin
{
    admin.DELETE("/users/:id", userHandler.Delete)
}`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit scaffolds four middleware functions: <code>Logger</code> (request timing),
                <code>CORS</code> (cross-origin access), <code>Auth</code> (JWT validation),
                and <code>RequireRole</code> (role-based access). They are applied in <code>routes.go</code>:
                Logger and CORS are global, Auth is per-group, and RequireRole stacks on top of Auth
                for admin routes.
              </p>
            </div>

            <PlaygroundChallenge
              title="Writing a Middleware"
              description="A middleware is a function that takes a handler and returns a handler. Write one that runs before and after the request, which is what c.Next() splits in Gin."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: your first middleware
// 1. Write: func Logger(next http.Handler) http.Handler
//      - return http.HandlerFunc(func(w, r) { ... })
//      - print "-> METHOD PATH" BEFORE calling next.ServeHTTP(w, r)
//      - print "<- done" AFTER it returns
//    The code before next is Gin's "before c.Next()"; after it is the rest.
// 2. Write a handler that writes "hello" with status 200
// 3. Wrap it:  wrapped := Logger(handler)
// 4. Call it with httptest and print the status and body
// 5. Bonus: write RequestID(next) that sets w.Header().Set("X-Request-ID", "abc123")
//    and wrap with both

func main() {
	fmt.Println(http.StatusOK, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Takes a handler, returns a handler — the whole middleware contract
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("-> %s %s\\n", r.Method, r.URL.Path) // before c.Next()
		next.ServeHTTP(w, r)
		fmt.Println("<- done") // after c.Next()
	})
}

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", "abc123")
		next.ServeHTTP(w, r)
	})
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "hello")
	})

	wrapped := Logger(RequestID(handler))

	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, httptest.NewRequest("GET", "/products", nil))

	fmt.Println("status:", rec.Code)
	fmt.Println("body:  ", rec.Body.String())
	fmt.Println("id:    ", rec.Header().Get("X-Request-ID"))
}`}
            />

            <PlaygroundChallenge
              title="Chaining and Short-Circuiting"
              description="Compose several middleware into one chain, then write one that refuses to call the next handler — which is exactly what c.Abort() does when auth fails."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: chains and aborts
// 1. Write func Chain(h http.Handler, mw ...func(http.Handler) http.Handler) http.Handler
//      - apply the middleware in REVERSE so the first listed runs first
// 2. Write three middleware that print "A start"/"A end", "B start"/"B end",
//    and RequireToken:
//      - if r.Header.Get("Authorization") is empty, respond 401 and RETURN
//        WITHOUT calling next (this is c.Abort())
//      - otherwise call next
// 3. Build the chain: Chain(handler, A, RequireToken, B)
// 4. Call it once WITHOUT the header and once WITH it, printing the status
// 5. Notice the nesting: A start, B start, handler, B end, A end

func main() {
	fmt.Println(http.StatusUnauthorized, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

type middleware func(http.Handler) http.Handler

// Applied in reverse so the first one listed is the outermost
func Chain(h http.Handler, mw ...middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}

func tag(name string) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(name, "start")
			next.ServeHTTP(w, r)
			fmt.Println(name, "end")
		})
	}
}

// Never calls next when the check fails — the equivalent of c.Abort()
func RequireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("handler")
		fmt.Fprint(w, "ok")
	})

	chain := Chain(handler, tag("A"), RequireToken, tag("B"))

	fmt.Println("--- no token ---")
	rec := httptest.NewRecorder()
	chain.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	fmt.Println("status:", rec.Code, rec.Body.String())

	fmt.Println("--- with token ---")
	rec = httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer x")
	chain.ServeHTTP(rec, req)
	fmt.Println("status:", rec.Code, rec.Body.String())
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 15. CORS */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="cors">15. CORS</h2>
              <p>
                <strong>CORS</strong> (Cross-Origin Resource Sharing) is a browser security feature
                that blocks web pages from making requests to a different domain than the one
                that served them. Your Next.js frontend runs on <code>localhost:3000</code> but
                your Go API runs on <code>localhost:8080</code> -- that&apos;s a different origin,
                so the browser blocks the request by default.
              </p>
              <p>
                To fix this, the API must send special headers
                (<code>Access-Control-Allow-Origin</code>) telling the browser which origins
                are allowed. This is handled by CORS middleware.
              </p>
            </div>

            <CodeBlock language="go" filename="middleware/cors.go" code={`package middleware

import (
    "strings"
    "github.com/gin-gonic/gin"
)

// CORS returns middleware that allows cross-origin requests.
func CORS(allowedOrigins string) gin.HandlerFunc {
    origins := strings.Split(allowedOrigins, ",")

    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")

        // Check if the request origin is allowed
        for _, allowed := range origins {
            if strings.TrimSpace(allowed) == origin {
                c.Header("Access-Control-Allow-Origin", origin)
                c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
                c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
                c.Header("Access-Control-Allow-Credentials", "true")
                break
            }
        }

        // Handle preflight requests
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}`} />

            <div className="prose-grit mb-10">
              <p>
                <strong>Preflight requests:</strong> Before making a POST or PUT request, the browser
                sends an OPTIONS request first (called a &quot;preflight&quot;) to check if CORS is allowed.
                The middleware handles this by returning a 204 with the correct headers.
              </p>
            </div>

            <CodeBlock language="bash" filename=".env" code={`# Comma-separated list of allowed frontend origins
CORS_ORIGINS=http://localhost:3000,http://localhost:3001`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit&apos;s CORS middleware reads allowed origins from the <code>CORS_ORIGINS</code> environment
                variable. By default, it allows <code>localhost:3000</code> (web app)
                and <code>localhost:3001</code> (admin panel). In production, update this to your
                actual domain. CORS is applied globally in <code>routes.go</code> so every
                endpoint is accessible from the frontend.
              </p>
            </div>

            <PlaygroundChallenge
              title="An Origin Allowlist"
              description="Decide the Access-Control-Allow-Origin header from a list of permitted origins. Echo the caller's origin when it is allowed, and send nothing at all when it is not."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: allowlist logic
// 1. Write func allowOrigin(origin string, allowed []string) string
//      - return origin when it appears in allowed
//      - return "" otherwise (send NO header rather than a wrong one)
// 2. Write a CORS middleware that:
//      - reads r.Header.Get("Origin")
//      - sets Access-Control-Allow-Origin only when allowOrigin returns non-empty
//      - also sets Access-Control-Allow-Credentials: true in that case
//      - always calls next
// 3. Test with an allowed origin, a disallowed one, and no Origin header,
//    printing the response header each time
//
// Never reflect an arbitrary origin back with credentials enabled — that
// hands any site on the internet an authenticated session.

func main() {
	fmt.Println(http.StatusOK, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func allowOrigin(origin string, allowed []string) string {
	for _, a := range allowed {
		if a == origin {
			return origin
		}
	}
	return "" // not allowed: send no header at all
}

func CORS(allowed []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if o := allowOrigin(r.Header.Get("Origin"), allowed); o != "" {
				w.Header().Set("Access-Control-Allow-Origin", o)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	allowed := []string{"https://app.example.com", "http://localhost:3000"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})
	wrapped := CORS(allowed)(handler)

	for _, origin := range []string{"http://localhost:3000", "https://evil.example", ""} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/products", nil)
		if origin != "" {
			req.Header.Set("Origin", origin)
		}
		wrapped.ServeHTTP(rec, req)
		fmt.Printf("origin=%-28q allow=%q\\n", origin, rec.Header().Get("Access-Control-Allow-Origin"))
	}
}`}
            />

            <PlaygroundChallenge
              title="Handling the Preflight"
              description="Answer the browser's OPTIONS request before it will send the real one. Get this wrong and every non-trivial request fails before your handler is ever called."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: preflight handling
// 1. Extend a CORS middleware so that when r.Method == http.MethodOptions it:
//      - sets Access-Control-Allow-Origin (allowed origins only)
//      - sets Access-Control-Allow-Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
//      - sets Access-Control-Allow-Headers: Content-Type, Authorization
//      - sets Access-Control-Max-Age: 86400
//      - responds 204 and RETURNS WITHOUT calling next
// 2. For every other method, set the origin header and call next as usual
// 3. Test three requests and print status plus the relevant headers:
//      - OPTIONS from an allowed origin       -> 204 with the headers
//      - OPTIONS from a disallowed origin     -> 204, no allow-origin header
//      - GET from an allowed origin           -> 200 from the handler

func main() {
	fmt.Println(http.StatusNoContent, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func allowOrigin(origin string, allowed []string) string {
	for _, a := range allowed {
		if a == origin {
			return origin
		}
	}
	return ""
}

func CORS(allowed []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := allowOrigin(r.Header.Get("Origin"), allowed)
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// The preflight never reaches your handler
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Max-Age", "86400")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func main() {
	allowed := []string{"https://app.example.com"}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "real response")
	})
	wrapped := CORS(allowed)(handler)

	cases := []struct{ method, origin string }{
		{"OPTIONS", "https://app.example.com"},
		{"OPTIONS", "https://evil.example"},
		{"GET", "https://app.example.com"},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(c.method, "/api/products", nil)
		req.Header.Set("Origin", c.origin)
		wrapped.ServeHTTP(rec, req)
		fmt.Printf("%-8s %-26s -> %d allow=%q methods=%q\\n",
			c.method, c.origin, rec.Code,
			rec.Header().Get("Access-Control-Allow-Origin"),
			rec.Header().Get("Access-Control-Allow-Methods"))
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 16. Handlers */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="handlers">16. Handlers</h2>
              <p>
                A handler is the function that runs when an HTTP request matches a route. In Grit,
                handlers follow the <strong>thin handler</strong> pattern: they do four things and
                nothing more:
              </p>
              <ol>
                <li><strong>Parse</strong> the request (URL params, query strings, JSON body)</li>
                <li><strong>Validate</strong> the input (using binding tags)</li>
                <li><strong>Delegate</strong> to a service or database call</li>
                <li><strong>Respond</strong> with the appropriate JSON and status code</li>
              </ol>
              <p>
                Handlers do NOT contain business logic. They don&apos;t hash passwords, calculate totals,
                send emails, or query related data. All of that goes in services. This separation
                makes your code testable and keeps each layer focused on one job.
              </p>
            </div>

            <CodeBlock language="go" filename="handlers/auth.go" code={`package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "myapp/apps/api/internal/models"
    "myapp/apps/api/internal/services"
)

// Handler struct holds dependencies
type AuthHandler struct {
    DB          *gorm.DB
    AuthService *services.AuthService
}

// Request struct — what the client sends
type loginRequest struct {
    Email    string \`json:"email" binding:"required,email"\`
    Password string \`json:"password" binding:"required"\`
}

// Login authenticates a user and returns JWT tokens.
func (h *AuthHandler) Login(c *gin.Context) {
    // 1. Parse & validate
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    // 2. Find user in database
    var user models.User
    if err := h.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "INVALID_CREDENTIALS",
                "message": "Invalid email or password",
            },
        })
        return
    }

    // 3. Check password (delegate to model method)
    if !user.CheckPassword(req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "INVALID_CREDENTIALS",
                "message": "Invalid email or password",
            },
        })
        return
    }

    // 4. Generate tokens (delegate to auth service)
    tokens, err := h.AuthService.GenerateTokenPair(user.ID, user.Email, user.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "TOKEN_ERROR",
                "message": "Failed to generate tokens",
            },
        })
        return
    }

    // 5. Respond
    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{
            "user":   user,
            "tokens": tokens,
        },
        "message": "Logged in successfully",
    })
}`} />

            <div className="prose-grit mb-10">
              <p>
                Strip Gin away and a handler is four steps in a fixed order: read the input, validate
                it, call the service, write the response. The version below is the same shape written
                against <code>net/http</code> so it runs anywhere, including the playground. The part
                worth copying is the last step — one place that turns a service error into a status
                code, so no handler ever invents its own mapping.
              </p>
            </div>

            <CodeBlock language="go" filename="handler_shape.go" code={`package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
)

// Sentinels the service returns; the handler maps them to status codes
var (
    ErrNotFound  = errors.New("not found")
    ErrDuplicate = errors.New("already exists")
)

type Product struct {
    ID    int     \`json:"id"\`
    Name  string  \`json:"name"\`
    Price float64 \`json:"price"\`
}

// The service knows nothing about HTTP
type ProductService struct{ items map[int]Product }

func (s *ProductService) Get(id int) (Product, error) {
    p, ok := s.items[id]
    if !ok {
        return Product{}, fmt.Errorf("get %d: %w", id, ErrNotFound)
    }
    return p, nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(body)
}

// One place that decides which error becomes which status
func writeError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, ErrNotFound):
        writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
    case errors.Is(err, ErrDuplicate):
        writeJSON(w, http.StatusConflict, map[string]string{"error": "already exists"})
    default:
        // Log the real error server-side; never leak it to the client
        fmt.Println("internal:", err)
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
    }
}

func main() {
    svc := &ProductService{items: map[int]Product{
        1: {ID: 1, Name: "Laptop", Price: 999},
    }}

    handler := func(w http.ResponseWriter, r *http.Request) {
        // 1. read input
        id := strings.TrimPrefix(r.URL.Path, "/products/")
        // 2. validate
        if id == "" {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id required"})
            return
        }
        // 3. call the service
        var n int
        fmt.Sscanf(id, "%d", &n)
        p, err := svc.Get(n)
        if err != nil {
            writeError(w, err)
            return
        }
        // 4. write the response
        writeJSON(w, http.StatusOK, p)
    }

    for _, path := range []string{"/products/1", "/products/42"} {
        rec := httptest.NewRecorder()
        handler(rec, httptest.NewRequest("GET", path, nil))
        fmt.Printf("%-14s %d %s", path, rec.Code, rec.Body.String())
    }
}`} />


            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit scaffolds auth handlers (<code>Login</code>, <code>Register</code>, <code>Refresh</code>,
                <code>ForgotPassword</code>, <code>Me</code>) and a user handler (<code>List</code>, <code>Create</code>,
                <code>GetByID</code>, <code>Update</code>, <code>Delete</code>). When you generate a resource,
                the CLI creates a handler with all five CRUD methods plus pagination, search, and sorting --
                all following the same thin pattern.
              </p>
            </div>

            <PlaygroundChallenge
              title="The Four Steps of a Handler"
              description="Read, validate, call the service, respond — in that order, every time. Build one end to end and prove it with httptest."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type CreateUserRequest struct {
	Name  string \`json:"name"\`
	Email string \`json:"email"\`
}

// Challenge: a complete handler
// 1. Decode the body into CreateUserRequest; on failure respond 400
// 2. Validate:
//      - Name must not be blank
//      - Email must contain "@"
//    Collect problems into map[string]string and respond 422 if any
// 3. "Call the service": if Email is "taken@example.com", respond 409
//    with {"error":"email already registered"}
// 4. Otherwise respond 201 with {"id":1,"name":...,"email":...}
// 5. Exercise all four paths with httptest and print status + body
//
// Keep the order. Validating after calling the service is how you end up
// writing a row you were about to reject.

func main() {
	fmt.Println(http.StatusCreated, strings.TrimSpace(" x "), httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

type CreateUserRequest struct {
	Name  string \`json:"name"\`
	Email string \`json:"email"\`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// 1. read
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}

	// 2. validate
	problems := map[string]string{}
	if strings.TrimSpace(req.Name) == "" {
		problems["name"] = "is required"
	}
	if !strings.Contains(req.Email, "@") {
		problems["email"] = "must be a valid email address"
	}
	if len(problems) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": problems})
		return
	}

	// 3. call the service
	if req.Email == "taken@example.com" {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "email already registered"})
		return
	}

	// 4. respond
	writeJSON(w, http.StatusCreated, map[string]any{
		"id": 1, "name": req.Name, "email": req.Email,
	})
}

func main() {
	bodies := []string{
		\`{"name":"Ada","email":"ada@example.com"}\`,
		\`{"name":"","email":"nope"}\`,
		\`{"name":"Bob","email":"taken@example.com"}\`,
		\`{{{\`,
	}
	for _, b := range bodies {
		rec := httptest.NewRecorder()
		createUser(rec, httptest.NewRequest("POST", "/users", strings.NewReader(b)))
		fmt.Printf("%d %s", rec.Code, rec.Body.String())
	}
}`}
            />

            <PlaygroundChallenge
              title="Mapping Errors to Status Codes"
              description="Write the one function that turns a service error into an HTTP status, so no handler has to guess and no internal message leaks to the client."
              challenge={`package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: a single error mapper
// 1. Declare sentinels:
//      ErrNotFound, ErrForbidden, ErrConflict
// 2. Write func writeError(w http.ResponseWriter, err error) that maps:
//      ErrNotFound  -> 404 {"error":"not found"}
//      ErrForbidden -> 403 {"error":"forbidden"}
//      ErrConflict  -> 409 {"error":"already exists"}
//      anything else -> 500 {"error":"internal error"}
//                       and print the real error to the server log only
//    Use errors.Is so wrapped errors still match.
// 3. Build four errors — three wrapped sentinels and one raw database error —
//    pass each to writeError with a fresh httptest.NewRecorder()
// 4. Print status and body, and confirm the database error's text never
//    appears in the response body

func main() {
	fmt.Println(errors.New("x"), http.StatusNotFound, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrConflict  = errors.New("already exists")
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	case errors.Is(err, ErrForbidden):
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
	case errors.Is(err, ErrConflict):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "already exists"})
	default:
		// Logged here, never sent: the client learns nothing about the schema
		fmt.Println("[log] internal:", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}

func main() {
	errs := []error{
		fmt.Errorf("products.Get 42: %w", ErrNotFound),
		fmt.Errorf("products.Delete: %w", ErrForbidden),
		fmt.Errorf("users.Create: %w", ErrConflict),
		errors.New(\`pq: duplicate key value violates unique constraint "users_email_key"\`),
	}
	for _, err := range errs {
		rec := httptest.NewRecorder()
		writeError(rec, err)
		fmt.Printf("%d %s", rec.Code, rec.Body.String())
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 17. Services & The Service Pattern */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="services">17. Services & The Service Pattern</h2>
              <p>
                A <strong>service</strong> is a struct with methods that contain your business logic.
                It sits between the handler (HTTP layer) and the database (data layer). But why not
                just put the logic directly in the handler?
              </p>
              <h3>Why Services Exist</h3>
              <ul>
                <li><strong>Separation of concerns</strong> -- handlers deal with HTTP, services deal with logic. Each layer has one job.</li>
                <li><strong>Testability</strong> -- you can test business logic without spinning up an HTTP server. Just create a service with a test database and call its methods.</li>
                <li><strong>Reusability</strong> -- the same service method can be called from a handler, a background job, a CLI command, or a cron task. If the logic was in the handler, you&apos;d have to duplicate it.</li>
                <li><strong>Maintainability</strong> -- when business rules change, you update one service method instead of hunting through handlers.</li>
              </ul>
            </div>

            <CodeBlock language="go" filename="services/product.go" code={`package services

import (
    "fmt"
    "math"

    "gorm.io/gorm"

    "myapp/apps/api/internal/models"
)

// Service struct — holds the database connection
type ProductService struct {
    DB *gorm.DB
}

// All operations are methods on the service

// List returns paginated products with search and sorting.
func (s *ProductService) List(page, pageSize int, search, sortBy, sortOrder string) ([]models.Product, int64, int, error) {
    query := s.DB.Model(&models.Product{})

    if search != "" {
        query = query.Where("name ILIKE ?", "%"+search+"%")
    }

    var total int64
    query.Count(&total)

    var items []models.Product
    offset := (page - 1) * pageSize
    err := query.Order(sortBy + " " + sortOrder).
        Offset(offset).
        Limit(pageSize).
        Find(&items).Error

    if err != nil {
        return nil, 0, 0, fmt.Errorf("fetching products: %w", err)
    }

    pages := int(math.Ceil(float64(total) / float64(pageSize)))
    return items, total, pages, nil
}

// GetByID returns a single product.
func (s *ProductService) GetByID(id uint) (*models.Product, error) {
    var item models.Product
    if err := s.DB.First(&item, id).Error; err != nil {
        return nil, fmt.Errorf("product not found: %w", err)
    }
    return &item, nil
}

// Create adds a new product.
func (s *ProductService) Create(item *models.Product) error {
    if err := s.DB.Create(item).Error; err != nil {
        return fmt.Errorf("creating product: %w", err)
    }
    return nil
}

// Delete soft-deletes a product.
func (s *ProductService) Delete(id uint) error {
    var item models.Product
    if err := s.DB.First(&item, id).Error; err != nil {
        return fmt.Errorf("product not found: %w", err)
    }
    return s.DB.Delete(&item).Error
}`} />

            <div className="prose-grit mb-10">
              <h3>How Handlers Call Services</h3>
              <p>
                The handler creates or receives a service instance, then calls its methods.
                The handler&apos;s only job is to translate between HTTP and the service layer:
              </p>
            </div>

            <CodeBlock language="go" filename="handlers/product.go (simplified)" code={`type ProductHandler struct {
    Service *services.ProductService
}

func (h *ProductHandler) GetByID(c *gin.Context) {
    id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

    // Delegate to service
    product, err := h.Service.GetByID(uint(id))
    if err != nil {
        c.JSON(404, gin.H{"error": gin.H{"code": "NOT_FOUND", "message": "Product not found"}})
        return
    }

    // Respond
    c.JSON(200, gin.H{"data": product})
}`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Every generated resource gets a service in <code>internal/services/</code> with
                <code>List</code>, <code>GetByID</code>, <code>Create</code>, <code>Update</code>,
                and <code>Delete</code> methods. The auth service (<code>AuthService</code>) handles
                JWT token generation and validation. Background job workers also call services --
                the same <code>ProductService.Create()</code> method can be called from an HTTP handler
                or an async job worker.
              </p>
            </div>

            <PlaygroundChallenge
              title="A Service That Owns the Rules"
              description="Put the business rules in a service that knows nothing about HTTP, then prove it by calling the same service from ordinary code."
              challenge={`package main

import "fmt"

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

// Challenge: a service with real rules
// 1. Declare: type ProductService struct { items map[int]*Product; nextID int }
// 2. Write a constructor NewProductService() *ProductService
// 3. Add methods that return (result, error):
//      Create(name string, price float64) (*Product, error)
//        - reject a blank name or a price <= 0
//        - assign the next id and store it
//      Reserve(id, qty int) error
//        - error if the product does not exist
//        - error if qty > Stock
//        - otherwise decrement Stock
// 4. In main: create two products, reserve some stock, and trigger every
//    error path, printing what happened each time
//
// Notice there is no http anywhere in the service — that is the point.

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrNoStock      = errors.New("insufficient stock")
)

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

type ProductService struct {
	items  map[int]*Product
	nextID int
}

func NewProductService() *ProductService {
	return &ProductService{items: map[int]*Product{}, nextID: 1}
}

func (s *ProductService) Create(name string, price float64) (*Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name is required: %w", ErrInvalidInput)
	}
	if price <= 0 {
		return nil, fmt.Errorf("price must be positive: %w", ErrInvalidInput)
	}
	p := &Product{ID: s.nextID, Name: name, Price: price, Stock: 10}
	s.items[p.ID] = p
	s.nextID++
	return p, nil
}

func (s *ProductService) Reserve(id, qty int) error {
	p, ok := s.items[id]
	if !ok {
		return fmt.Errorf("product %d: %w", id, ErrNotFound)
	}
	if qty > p.Stock {
		return fmt.Errorf("want %d, have %d: %w", qty, p.Stock, ErrNoStock)
	}
	p.Stock -= qty
	return nil
}

func main() {
	svc := NewProductService()

	laptop, _ := svc.Create("Laptop", 999)
	fmt.Printf("created %+v\\n", *laptop)

	if _, err := svc.Create("", 10); err != nil {
		fmt.Println("rejected:", err)
	}
	if _, err := svc.Create("Mouse", -1); err != nil {
		fmt.Println("rejected:", err)
	}

	fmt.Println("reserve 3:", svc.Reserve(laptop.ID, 3))
	fmt.Println("stock now:", svc.items[laptop.ID].Stock)
	fmt.Println("reserve 99:", svc.Reserve(laptop.ID, 99))
	fmt.Println("reserve missing:", svc.Reserve(404, 1))
}`}
            />

            <PlaygroundChallenge
              title="Depending on an Interface"
              description="Have the service depend on a repository interface rather than a concrete database. That is what makes it testable without a running Postgres."
              challenge={`package main

import "fmt"

type User struct {
	ID    int
	Email string
}

// Challenge: dependency injection through an interface
// 1. Declare: type UserRepo interface {
//      FindByEmail(email string) (*User, error)
//      Create(u *User) error
//    }
// 2. Write an in-memory implementation: type MemoryRepo struct { users []*User }
//      - FindByEmail returns a "not found" error when absent
//      - Create appends and assigns an ID
// 3. Write: type UserService struct { repo UserRepo }
//      - Register(email string) (*User, error) that refuses a duplicate email
// 4. In main, build the service with a MemoryRepo, register twice with the
//    same email, and show the second attempt failing
// 5. Bonus: write a FailingRepo whose Create always errors, hand it to the
//    SAME service, and confirm the error surfaces — no database required

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type User struct {
	ID    int
	Email string
}

// The service depends on this, not on a database
type UserRepo interface {
	FindByEmail(email string) (*User, error)
	Create(u *User) error
}

type MemoryRepo struct{ users []*User }

func (r *MemoryRepo) FindByEmail(email string) (*User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (r *MemoryRepo) Create(u *User) error {
	u.ID = len(r.users) + 1
	r.users = append(r.users, u)
	return nil
}

type FailingRepo struct{ MemoryRepo }

func (r *FailingRepo) Create(u *User) error {
	return errors.New("connection refused")
}

type UserService struct{ repo UserRepo }

func (s *UserService) Register(email string) (*User, error) {
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, fmt.Errorf("register %s: already registered", email)
	} else if !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("register %s: %w", email, err)
	}

	u := &User{Email: email}
	if err := s.repo.Create(u); err != nil {
		return nil, fmt.Errorf("register %s: %w", email, err)
	}
	return u, nil
}

func main() {
	svc := &UserService{repo: &MemoryRepo{}}

	u, err := svc.Register("ada@example.com")
	fmt.Println("first: ", u, err)

	_, err = svc.Register("ada@example.com")
	fmt.Println("second:", err)

	// Same service, different repo — no database needed to test the failure
	failing := &UserService{repo: &FailingRepo{}}
	_, err = failing.Register("grace@example.com")
	fmt.Println("failing repo:", err)
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 18. GORM In Depth */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="gorm-in-depth">18. GORM In Depth</h2>
              <p>
                GORM is Go&apos;s most popular ORM. It maps Go structs to database tables and provides
                a chainable API for queries. Let&apos;s cover the key operations you&apos;ll use daily.
              </p>
              <h3>Database Connection</h3>
              <p>
                GORM connects to PostgreSQL using the <code>gorm.io/driver/postgres</code> driver.
                You open a connection once at startup and pass it everywhere via dependency injection.
              </p>
            </div>

            <CodeBlock language="go" filename="database/database.go" code={`package database

import (
    "fmt"
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func Connect(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.New(postgres.Config{
        DSN:                  dsn,
        PreferSimpleProtocol: true,
    }), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }

    // Configure connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)

    log.Println("Database connected successfully")
    return db, nil
}`} />

            <div className="prose-grit mb-10">
              <h3>CRUD Operations</h3>
              <p>
                GORM provides a chainable API for all database operations. Each method returns
                the same <code>*gorm.DB</code>, so you can chain them together.
              </p>
            </div>

            <CodeBlock language="go" filename="gorm_crud.go" code={`// ── CREATE ─────────────────────────────────────────────
product := models.Product{Name: "Widget", Price: 29.99}
db.Create(&product)          // INSERT INTO products ...
fmt.Println(product.ID)      // ID is auto-set after create

// ── READ — single record ──────────────────────────────
var found models.Product
db.First(&found, 42)                          // WHERE id = 42
db.Where("email = ?", "alice@test.com").First(&found)  // WHERE email = ...

// ── READ — multiple records ───────────────────────────
var products []models.Product
db.Find(&products)                            // SELECT * FROM products
db.Where("price > ?", 20.0).Find(&products)  // With condition

// ── READ — pagination and sorting ─────────────────────
db.Order("created_at desc").
    Offset(0).    // Skip 0 records (page 1)
    Limit(20).    // Take 20 records
    Find(&products)

// ── READ — count ──────────────────────────────────────
var total int64
db.Model(&models.Product{}).Count(&total)

// ── READ — search with ILIKE (case-insensitive) ──────
search := "widget"
db.Where("name ILIKE ?", "%"+search+"%").Find(&products)

// ── UPDATE ────────────────────────────────────────────
db.Model(&found).Update("price", 34.99)      // Single field
db.Model(&found).Updates(map[string]any{      // Multiple fields
    "name":  "Super Widget",
    "price": 39.99,
})

// ── DELETE (soft delete) ──────────────────────────────
db.Delete(&found)            // Sets deleted_at, doesn't remove row
// To permanently delete: db.Unscoped().Delete(&found)`} />

            <div className="prose-grit mb-10">
              <h3>Relationships & Preloading</h3>
              <p>
                When a model has relationships (belongs-to, has-many), GORM does NOT load related
                data automatically. You must use <code>Preload()</code> to eagerly load them.
              </p>
            </div>

            <CodeBlock language="go" filename="preloading.go" code={`// Models with relationships
type Category struct {
    ID       uint      \`gorm:"primarykey" json:"id"\`
    Name     string    \`json:"name"\`
    Products []Product \`json:"products"\`  // has many
}

type Product struct {
    ID         uint     \`gorm:"primarykey" json:"id"\`
    Name       string   \`json:"name"\`
    CategoryID uint     \`json:"category_id"\`           // foreign key
    Category   Category \`json:"category"\`               // belongs to
}

// Without Preload — category field will be empty {}
db.First(&product, 1)
fmt.Println(product.Category.Name) // "" (empty!)

// With Preload — category is loaded
db.Preload("Category").First(&product, 1)
fmt.Println(product.Category.Name) // "Electronics"`} />

            <div className="prose-grit mb-10">
              <h3>Hooks (Lifecycle Callbacks)</h3>
              <p>
                GORM hooks are methods on your model that run automatically at specific points
                in the lifecycle. The most common hook is <code>BeforeCreate</code>, used to
                hash passwords before they are stored in the database.
              </p>
            </div>

            <CodeBlock language="go" filename="models/user.go (hooks)" code={`import "golang.org/x/crypto/bcrypt"

// BeforeCreate runs automatically before INSERT
func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.Password != "" {
        hashed, err := bcrypt.GenerateFromPassword(
            []byte(u.Password), bcrypt.DefaultCost,
        )
        if err != nil {
            return err
        }
        u.Password = string(hashed)
    }
    return nil
}

// CheckPassword compares plaintext against stored hash
func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword(
        []byte(u.Password), []byte(password),
    )
    return err == nil
}

// Usage — password is hashed automatically
user := models.User{
    Email:    "alice@example.com",
    Password: "mypassword123",  // Plaintext here
}
db.Create(&user) // BeforeCreate hashes it before INSERT`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Every resource generated by <code>grit generate resource</code> gets a service file
                with these exact GORM operations: <code>Create</code>, <code>List</code> (with
                pagination, search, and sorting), <code>GetByID</code> (with Preload),
                <code>Update</code>, and <code>Delete</code>. The database connection is
                established once in <code>internal/database/database.go</code> and passed to all
                services via dependency injection.
              </p>
            </div>

            <PlaygroundChallenge
              title="Why Preloading Exists"
              description="Reproduce the N+1 query problem with an in-memory store, count the lookups, then fix it the way Preload does — one query for the children instead of one per parent."
              challenge={`package main

import "fmt"

type Author struct {
	ID   int
	Name string
}

type Book struct {
	ID       int
	Title    string
	AuthorID int
}

// Challenge: N+1 and how Preload fixes it
// 1. Write findAuthor(id int) *Author that increments a global queryCount
//    and returns the matching author (this stands in for a database round trip)
// 2. Loop the books calling findAuthor for each — print the count.
//    With 5 books that is 1 + 5 = 6 "queries": the N+1 problem.
// 3. Now do it the Preload way:
//      - collect the distinct author ids
//      - fetch them in ONE call (findAuthorsIn(ids []int) map[int]*Author)
//      - join in memory from the map
// 4. Print the query count for both approaches and the joined output
//
// GORM: db.Find(&books) then db.Preload("Author").Find(&books)

var queryCount int

func main() {
	fmt.Println(queryCount)
}`}
              solution={`package main

import "fmt"

type Author struct {
	ID   int
	Name string
}

type Book struct {
	ID       int
	Title    string
	AuthorID int
}

var (
	authors = map[int]*Author{
		1: {ID: 1, Name: "Ada"},
		2: {ID: 2, Name: "Grace"},
	}
	books = []Book{
		{1, "Go in Practice", 1},
		{2, "Concurrency", 1},
		{3, "Compilers", 2},
		{4, "Databases", 2},
		{5, "Networks", 1},
	}
	queryCount int
)

func findAuthor(id int) *Author {
	queryCount++ // one round trip
	return authors[id]
}

func findAuthorsIn(ids []int) map[int]*Author {
	queryCount++ // still one round trip, however many ids
	out := map[int]*Author{}
	for _, id := range ids {
		if a, ok := authors[id]; ok {
			out[id] = a
		}
	}
	return out
}

func main() {
	// N+1: one query for the books, then one per book
	queryCount = 1
	for _, b := range books {
		_ = findAuthor(b.AuthorID)
	}
	fmt.Println("N+1 queries:", queryCount) // 6

	// Preload: one query for the books, one for all the authors
	queryCount = 1
	seen := map[int]bool{}
	ids := []int{}
	for _, b := range books {
		if !seen[b.AuthorID] {
			seen[b.AuthorID] = true
			ids = append(ids, b.AuthorID)
		}
	}
	byID := findAuthorsIn(ids)
	fmt.Println("preload queries:", queryCount) // 2

	for _, b := range books {
		fmt.Printf("  %-16s by %s\\n", b.Title, byID[b.AuthorID].Name)
	}
}`}
            />

            <PlaygroundChallenge
              title="Hooks and Soft Deletes"
              description="Model two GORM behaviours that surprise people: a BeforeCreate hook that rewrites the record on the way in, and a soft delete that hides rows without removing them."
              challenge={`package main

import "fmt"

type User struct {
	ID       int
	Email    string
	Password string
	Deleted  bool
}

// Challenge: hooks and soft deletes
// 1. Write func (u *User) BeforeCreate() error that:
//      - lowercases and trims the Email
//      - refuses to continue if Password is shorter than 8 (return an error)
//      - replaces Password with "hashed:" + the original
// 2. Write a Store with Create(u *User) error that calls BeforeCreate first
//    and only stores the user if it returns nil
// 3. Add SoftDelete(id int) and List() []User where List SKIPS deleted rows
//    (GORM does this automatically for models with gorm.DeletedAt)
// 4. Create two users (one with a short password), soft delete one,
//    then print List() and the raw slice so the difference is visible

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"errors"
	"fmt"
	"strings"
)

type User struct {
	ID       int
	Email    string
	Password string
	Deleted  bool
}

// GORM calls this automatically; here we call it from Create
func (u *User) BeforeCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	u.Password = "hashed:" + u.Password
	return nil
}

type Store struct {
	rows   []User
	nextID int
}

func (s *Store) Create(u *User) error {
	if err := u.BeforeCreate(); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	s.nextID++
	u.ID = s.nextID
	s.rows = append(s.rows, *u)
	return nil
}

func (s *Store) SoftDelete(id int) {
	for i := range s.rows {
		if s.rows[i].ID == id {
			s.rows[i].Deleted = true // the row stays, it is just hidden
		}
	}
}

// The default query excludes soft-deleted rows
func (s *Store) List() []User {
	out := []User{}
	for _, u := range s.rows {
		if !u.Deleted {
			out = append(out, u)
		}
	}
	return out
}

func main() {
	store := &Store{}

	ada := &User{Email: "  ADA@Example.com ", Password: "correct-horse"}
	fmt.Println("create ada:", store.Create(ada))
	fmt.Printf("stored as: %+v\\n", *ada)

	short := &User{Email: "bob@example.com", Password: "abc"}
	fmt.Println("create bob:", store.Create(short))

	grace := &User{Email: "grace@example.com", Password: "another-long-one"}
	store.Create(grace)

	store.SoftDelete(ada.ID)
	fmt.Println("List() returns:", len(store.List()), "of", len(store.rows), "rows")
	for _, u := range store.List() {
		fmt.Println("  visible:", u.Email)
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 19. Migrations & Seeding */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="migrations-seeding">19. Migrations & Seeding</h2>
              <p>
                <strong>Migrations</strong> create database tables from your Go structs.
                <strong>Seeding</strong> populates tables with initial data for development.
                Both are essential for getting a working database up and running.
              </p>
              <h3>AutoMigrate</h3>
              <p>
                GORM&apos;s <code>AutoMigrate</code> reads your struct fields and creates or updates
                the corresponding database table. It will add new columns but will NOT delete
                removed columns or change existing column types (to prevent data loss).
              </p>
            </div>

            <CodeBlock language="go" filename="models/models.go" code={`package models

import (
    "log"
    "gorm.io/gorm"
)

// Models returns ALL models in migration order.
// Models with no foreign key dependencies come first.
func Models() []interface{} {
    return []interface{}{
        &User{},
        &Upload{},
        &Blog{},
        // grit:models  ← new models are injected here
    }
}

// Migrate creates tables that don't exist yet.
func Migrate(db *gorm.DB) error {
    models := Models()

    for _, model := range models {
        // Skip if table already exists
        if db.Migrator().HasTable(model) {
            log.Printf("  ✓ %T — already exists, skipping", model)
            continue
        }

        if err := db.AutoMigrate(model); err != nil {
            return fmt.Errorf("migrating %T: %w", model, err)
        }
        log.Printf("  ✓ %T — created", model)
    }

    return nil
}`} />

            <div className="prose-grit mb-10">
              <h3>Running Migrations</h3>
              <p>
                Grit provides a dedicated CLI command for migrations with a <code>--fresh</code>
                flag that drops all tables before recreating them (useful during development).
              </p>
            </div>

            <CodeBlock terminal code={`# Run migrations (create missing tables)
grit migrate

# Fresh migration (drop all tables + recreate)
grit migrate --fresh`} />

            <div className="prose-grit mb-10">
              <h3>Seeding</h3>
              <p>
                Seeders create test data for development. A good seeder is <strong>idempotent</strong> --
                it checks if data already exists before creating it, so you can run it multiple times safely.
              </p>
            </div>

            <CodeBlock language="go" filename="database/seed.go" code={`package database

import (
    "log"
    "myapp/apps/api/internal/models"
    "gorm.io/gorm"
)

// Seed populates the database with initial data.
func Seed(db *gorm.DB) error {
    if err := seedAdminUser(db); err != nil {
        return fmt.Errorf("seeding admin: %w", err)
    }
    if err := seedDemoUsers(db); err != nil {
        return fmt.Errorf("seeding users: %w", err)
    }
    // grit:seeders  ← new seeders injected here
    return nil
}

// Idempotent seeder — checks before creating
func seedAdminUser(db *gorm.DB) error {
    var count int64
    db.Model(&models.User{}).Where("email = ?", "admin@example.com").Count(&count)
    if count > 0 {
        log.Println("Admin already exists, skipping...")
        return nil
    }

    admin := models.User{
        FirstName: "Admin",
        LastName:  "User",
        Email:     "admin@example.com",
        Password:  "password",     // Hashed by BeforeCreate hook
        Role:      "ADMIN",
        Active:    true,
    }

    if err := db.Create(&admin).Error; err != nil {
        return fmt.Errorf("creating admin: %w", err)
    }

    log.Println("Created admin: admin@example.com / password")
    return nil
}`} />

            <CodeBlock terminal code={`# Run the seeder
grit seed`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit scaffolds both <code>cmd/migrate/main.go</code> and <code>cmd/seed/main.go</code>
                out of the box. The seed file includes an admin user, demo users with different roles,
                and sample blog posts. When you generate a new resource, the model is automatically
                registered in <code>Models()</code> for migration. You can also use <code>grit migrate</code>
                and <code>grit seed</code> CLI commands.
              </p>
            </div>

            <PlaygroundChallenge
              title="Idempotent Seeding"
              description="A seeder has to be safe to run twice. Write find-or-create so a second run changes nothing, which is the difference between a seeder and a duplicate-row generator."
              challenge={`package main

import "fmt"

type User struct {
	ID    int
	Email string
	Role  string
}

// Challenge: seed without duplicating
// 1. Build a Store with rows []User and nextID int
// 2. Write FirstOrCreate(email, role string) (*User, bool):
//      - if a user with that email exists, return it and false (not created)
//      - otherwise append a new one and return it with true
// 3. Write Seed(s *Store) that seeds three users:
//      admin@example.com ADMIN, editor@example.com EDITOR, user@example.com USER
//    printing "created" or "exists" for each
// 4. Call Seed TWICE and print the row count after each run —
//    it must be 3 both times
// 5. Bonus: make the second run update the role of an existing user
//    without adding a row

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import "fmt"

type User struct {
	ID    int
	Email string
	Role  string
}

type Store struct {
	rows   []User
	nextID int
}

// The seeding primitive: find it, or create it — never blindly insert
func (s *Store) FirstOrCreate(email, role string) (*User, bool) {
	for i := range s.rows {
		if s.rows[i].Email == email {
			return &s.rows[i], false
		}
	}
	s.nextID++
	s.rows = append(s.rows, User{ID: s.nextID, Email: email, Role: role})
	return &s.rows[len(s.rows)-1], true
}

func Seed(s *Store) {
	seeds := []struct{ email, role string }{
		{"admin@example.com", "ADMIN"},
		{"editor@example.com", "EDITOR"},
		{"user@example.com", "USER"},
	}
	for _, sd := range seeds {
		u, created := s.FirstOrCreate(sd.email, sd.role)
		if created {
			fmt.Printf("  created %s (%s)\\n", u.Email, u.Role)
		} else {
			fmt.Printf("  exists  %s (%s)\\n", u.Email, u.Role)
		}
	}
}

func main() {
	store := &Store{}

	fmt.Println("first run:")
	Seed(store)
	fmt.Println("rows:", len(store.rows))

	fmt.Println("second run:")
	Seed(store)
	fmt.Println("rows:", len(store.rows)) // still 3

	// Bonus: update in place, no new row
	if u, created := store.FirstOrCreate("user@example.com", "USER"); !created {
		u.Role = "EDITOR"
	}
	fmt.Println("after promotion:", store.rows[2], "rows:", len(store.rows))
}`}
            />

            <PlaygroundChallenge
              title="Migration Order and Dependencies"
              description="AutoMigrate has to create a table before anything references it. Sort a set of migrations by their dependencies and detect the cycle that would make the order impossible."
              challenge={`package main

import "fmt"

// Challenge: order migrations by dependency
// 1. Given: migrations = map[string][]string where the value lists what a
//    table depends on, e.g. "orders": {"users", "products"}
// 2. Write Resolve(m map[string][]string) ([]string, error) that returns the
//    tables in an order where every dependency comes first
//    Hint: repeatedly take any table whose dependencies are all already placed
// 3. If no progress can be made and tables remain, return an error naming
//    them — that is a dependency cycle
// 4. Print the resolved order for the valid set
// 5. Then add a cycle ("a" needs "b", "b" needs "a") and print the error
//
// Sort the names at each step so the output is the same on every run.

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"errors"
	"fmt"
	"sort"
)

func Resolve(m map[string][]string) ([]string, error) {
	placed := map[string]bool{}
	var order []string

	for len(placed) < len(m) {
		// Collect everything whose dependencies are already satisfied
		var ready []string
		for table, deps := range m {
			if placed[table] {
				continue
			}
			ok := true
			for _, d := range deps {
				if _, known := m[d]; known && !placed[d] {
					ok = false
					break
				}
			}
			if ok {
				ready = append(ready, table)
			}
		}

		if len(ready) == 0 {
			var stuck []string
			for t := range m {
				if !placed[t] {
					stuck = append(stuck, t)
				}
			}
			sort.Strings(stuck)
			return nil, fmt.Errorf("dependency cycle among %v: %w", stuck, errors.New("cannot order migrations"))
		}

		sort.Strings(ready) // stable output
		for _, t := range ready {
			placed[t] = true
			order = append(order, t)
		}
	}
	return order, nil
}

func main() {
	valid := map[string][]string{
		"users":       {},
		"products":    {},
		"orders":      {"users", "products"},
		"order_items": {"orders", "products"},
		"reviews":     {"users", "products"},
	}
	order, err := Resolve(valid)
	fmt.Println("order:", order, "err:", err)

	cyclic := map[string][]string{
		"a": {"b"},
		"b": {"a"},
	}
	_, err = Resolve(cyclic)
	fmt.Println("cycle:", err)
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 20. JWT & Authentication */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="jwt-auth">20. JWT & Authentication</h2>
              <p>
                JWT (JSON Web Token) is how Grit authenticates users. Understanding this flow
                is critical because it connects the frontend, the API, the middleware, and the database.
                Let&apos;s break it down step by step.
              </p>
              <h3>What is a JWT?</h3>
              <p>
                A JWT is a signed string that contains data (called <strong>claims</strong>). The server
                creates a token by encoding claims (user ID, email, role) and signing it with a secret key.
                The client stores this token and sends it with every request. The server validates the
                signature to verify the token hasn&apos;t been tampered with.
              </p>
            </div>

            <CodeBlock language="go" filename="how JWT works" code={`// 1. JWT contains "claims" — data about the user
type Claims struct {
    UserID uint   \`json:"user_id"\`
    Email  string \`json:"email"\`
    Role   string \`json:"role"\`
    jwt.RegisteredClaims          // Expiry, issued-at, etc.
}

// 2. Server creates a token by signing claims with a secret
claims := &Claims{
    UserID: 42,
    Email:  "alice@example.com",
    Role:   "ADMIN",
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    },
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, _ := token.SignedString([]byte("my-secret-key"))
// tokenString = "eyJhbGciOiJIUzI1NiIs..."

// 3. Later, server validates the token
parsed, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
    return []byte("my-secret-key"), nil
})
claims = parsed.Claims.(*Claims)
fmt.Println(claims.UserID) // 42`} />

            <div className="prose-grit mb-10">
              <h3>The Authentication Flow</h3>
              <p>
                Here is the complete flow from registration to authenticated requests. Understanding
                this will make the entire auth system click.
              </p>
            </div>

            <CodeBlock language="bash" filename="authentication flow" code={`┌─────────────────────────────────────────────────────────────────┐
│ 1. REGISTER                                                     │
│                                                                 │
│  Client sends:    POST /api/auth/register                       │
│                   { "email": "alice@test.com",                  │
│                     "password": "mypassword" }                  │
│                                                                 │
│  Server does:     ① Validate input (binding tags)               │
│                   ② Check email doesn't already exist           │
│                   ③ Create user (BeforeCreate hashes password)  │
│                   ④ Generate access token (15min) + refresh     │
│                     token (7 days)                              │
│                   ⑤ Return { user, tokens }                     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 2. LOGIN                                                        │
│                                                                 │
│  Client sends:    POST /api/auth/login                          │
│                   { "email": "alice@test.com",                  │
│                     "password": "mypassword" }                  │
│                                                                 │
│  Server does:     ① Find user by email (db.Where)              │
│                   ② Check password (bcrypt.CompareHashAndPassword)│
│                   ③ Check account is active                     │
│                   ④ Generate new token pair                     │
│                   ⑤ Return { user, tokens }                     │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 3. AUTHENTICATED REQUEST                                        │
│                                                                 │
│  Client sends:    GET /api/users                                │
│                   Authorization: Bearer eyJhbGciOi...           │
│                                                                 │
│  Middleware does:  ① Extract token from "Bearer <token>"        │
│                    ② Validate signature + check expiry          │
│                    ③ Load user from DB by claims.UserID         │
│                    ④ Set user data in context                   │
│                    ⑤ Call c.Next() → handler runs               │
│                                                                 │
│  Handler does:    Read c.Get("user") → return response          │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ 4. TOKEN REFRESH                                                │
│                                                                 │
│  When access token expires (15min), client sends:               │
│  POST /api/auth/refresh { "refresh_token": "eyJ..." }          │
│                                                                 │
│  Server validates the refresh token and returns new tokens.     │
│  Client never needs to log in again until the refresh           │
│  token expires (7 days).                                        │
└─────────────────────────────────────────────────────────────────┘`} />

            <div className="prose-grit mb-10">
              <h3>The Auth Service</h3>
              <p>
                The auth service handles all token operations. It&apos;s a struct with the JWT secret
                and expiry durations, with methods for generating and validating tokens.
              </p>
            </div>

            <CodeBlock language="go" filename="services/auth.go" code={`type AuthService struct {
    Secret        string
    AccessExpiry  time.Duration  // e.g., 15 minutes
    RefreshExpiry time.Duration  // e.g., 7 days
}

type TokenPair struct {
    AccessToken  string \`json:"access_token"\`
    RefreshToken string \`json:"refresh_token"\`
    ExpiresAt    int64  \`json:"expires_at"\`
}

// GenerateTokenPair creates both access and refresh tokens.
func (s *AuthService) GenerateTokenPair(userID uint, email, role string) (*TokenPair, error) {
    // Access token — short-lived, used for API requests
    accessToken, expiresAt, err := s.generateToken(userID, email, role, s.AccessExpiry)
    if err != nil {
        return nil, fmt.Errorf("generating access token: %w", err)
    }

    // Refresh token — long-lived, used only to get new access tokens
    refreshToken, _, err := s.generateToken(userID, email, role, s.RefreshExpiry)
    if err != nil {
        return nil, fmt.Errorf("generating refresh token: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresAt:    expiresAt,
    }, nil
}

// ValidateToken parses and verifies a token string.
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            // Verify the signing method is HMAC (prevent algorithm attacks)
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method")
            }
            return []byte(s.Secret), nil
        },
    )
    if err != nil {
        return nil, fmt.Errorf("parsing token: %w", err)
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                The auth service is created in <code>routes.go</code> with the JWT secret and expiry
                durations from the config. It&apos;s passed to the auth handler and the auth middleware.
                On the frontend, React Query stores the tokens and automatically refreshes
                them when the access token expires. The <code>api-client.ts</code> intercepts 401
                responses and tries a silent refresh before showing the login page.
              </p>
            </div>

            <PlaygroundChallenge
              title="Build a JWT by Hand"
              description="A JWT is three base64 segments joined by dots, the third being an HMAC of the first two. Build and verify one with nothing but the standard library and it stops being magic."
              challenge={`package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Challenge: sign and verify a token yourself
// 1. Write encodeSegment(b []byte) string using base64.RawURLEncoding
// 2. Build the header  {"alg":"HS256","typ":"JWT"}  and a claims payload
//    {"sub":"42","role":"ADMIN","exp":1893456000}  as JSON, then encode both
// 3. signingInput := header + "." + payload
// 4. Sign it: hmac.New(sha256.New, secret), write signingInput, Sum(nil),
//    then encode the result. The token is signingInput + "." + signature
// 5. Write verify(token string, secret []byte) bool that recomputes the
//    signature and compares with hmac.Equal — NOT with ==
// 6. Print the token, verify it, then flip one character in the payload
//    and verify again to watch it fail
//
// Note the payload is only encoded, never encrypted: anyone can read it.

func main() {
	fmt.Println(base64.RawURLEncoding.EncodeToString([]byte("hi")), hmac.Equal(nil, nil), sha256.Size)
}`}
              solution={`package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

func encodeSegment(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

func sign(signingInput string, secret []byte) string {
	m := hmac.New(sha256.New, secret)
	m.Write([]byte(signingInput))
	return encodeSegment(m.Sum(nil))
}

func makeToken(secret []byte) string {
	header := encodeSegment([]byte(\`{"alg":"HS256","typ":"JWT"}\`))
	payload := encodeSegment([]byte(\`{"sub":"42","role":"ADMIN","exp":1893456000}\`))
	signingInput := header + "." + payload
	return signingInput + "." + sign(signingInput, secret)
}

func verify(token string, secret []byte) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	expected := sign(parts[0]+"."+parts[1], secret)
	// Constant-time compare: == would leak timing information
	return hmac.Equal([]byte(expected), []byte(parts[2]))
}

func main() {
	secret := []byte("a-long-secret-from-the-environment")

	token := makeToken(secret)
	fmt.Println("token:", token)
	fmt.Println("valid:", verify(token, secret))

	// The payload is encoded, not encrypted — anyone can read it
	parts := strings.Split(token, ".")
	claims, _ := base64.RawURLEncoding.DecodeString(parts[1])
	fmt.Println("claims:", string(claims))

	// Tamper with the claims and the signature no longer matches
	forged := parts[0] + "." + encodeSegment([]byte(\`{"sub":"42","role":"SUPERADMIN"}\`)) + "." + parts[2]
	fmt.Println("forged valid:", verify(forged, secret))

	// The wrong secret fails too
	fmt.Println("wrong secret:", verify(token, []byte("guess")))
}`}
            />

            <PlaygroundChallenge
              title="Expiry and the Auth Middleware"
              description="Parse a token, reject it when the claims say it has expired, and put the user on the request context so handlers downstream can read it."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: an auth middleware
// 1. Write parseToken(token string, now int64) (userID string, role string, err error)
//    Tokens here are the simple form "userID:role:expiry", e.g. "42:ADMIN:2000"
//      - error if the token does not have three parts
//      - error if expiry <= now  (say "token expired")
// 2. Write Auth(next http.Handler) http.Handler that:
//      - reads the Authorization header, requiring the "Bearer " prefix
//      - responds 401 when it is missing or invalid
//      - on success stores the user id and role on the request context
//        with context.WithValue and calls next with r.WithContext(ctx)
// 3. Write a handler that reads them back from the context and prints them
// 4. Test: no header, a malformed header, an expired token, a valid one
//
// Use your own key type for the context (type ctxKey string) so it cannot
// collide with a key set by another package.

func main() {
	fmt.Println(http.StatusUnauthorized, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type ctxKey string

const (
	ctxUserID ctxKey = "userID"
	ctxRole   ctxKey = "role"
)

func parseToken(token string, now int64) (string, string, error) {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return "", "", errors.New("malformed token")
	}
	exp, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return "", "", errors.New("malformed expiry")
	}
	if exp <= now {
		return "", "", errors.New("token expired")
	}
	return parts[0], parts[1], nil
}

func Auth(now int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			id, role, err := parseToken(strings.TrimPrefix(header, "Bearer "), now)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxUserID, id)
			ctx = context.WithValue(ctx, ctxRole, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func main() {
	const now int64 = 1500

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := r.Context().Value(ctxUserID).(string)
		role, _ := r.Context().Value(ctxRole).(string)
		fmt.Fprintf(w, "user=%s role=%s", id, role)
	})
	wrapped := Auth(now)(handler)

	cases := []struct{ name, header string }{
		{"no header", ""},
		{"malformed", "Bearer nonsense"},
		{"expired", "Bearer 42:ADMIN:1000"},
		{"valid", "Bearer 42:ADMIN:2000"},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/me", nil)
		if c.header != "" {
			req.Header.Set("Authorization", c.header)
		}
		wrapped.ServeHTTP(rec, req)
		fmt.Printf("%-10s %d %s\\n", c.name, rec.Code, strings.TrimSpace(rec.Body.String()))
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 21. RBAC & Middleware */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="rbac-middleware">21. RBAC & Middleware</h2>
              <p>
                <strong>RBAC</strong> (Role-Based Access Control) controls who can do what. Grit uses
                three default roles: <code>ADMIN</code>, <code>EDITOR</code>, and <code>USER</code>.
                This is enforced through two middleware functions that work together.
              </p>
              <h3>Auth Middleware</h3>
              <p>
                The <code>Auth</code> middleware runs on every protected route. It extracts the JWT
                from the Authorization header, validates it, loads the user from the database, and
                stores the user data in the Gin context so handlers can access it.
              </p>
            </div>

            <CodeBlock language="go" filename="middleware/auth.go" code={`func Auth(db *gorm.DB, authService *services.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Get the Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Authorization header required"}})
            c.Abort() // Stop the chain — handler never runs
            return
        }

        // 2. Extract "Bearer <token>"
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Invalid header format"}})
            c.Abort()
            return
        }

        // 3. Validate the token
        claims, err := authService.ValidateToken(parts[1])
        if err != nil {
            c.JSON(401, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Invalid or expired token"}})
            c.Abort()
            return
        }

        // 4. Load user from database
        var user models.User
        if err := db.First(&user, claims.UserID).Error; err != nil {
            c.JSON(401, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "User not found"}})
            c.Abort()
            return
        }

        // 5. Store user data in context for handlers
        c.Set("user", user)
        c.Set("user_id", user.ID)
        c.Set("user_role", user.Role)

        c.Next() // Continue to the handler
    }
}`} />

            <div className="prose-grit mb-10">
              <h3>RequireRole Middleware</h3>
              <p>
                The <code>RequireRole</code> middleware stacks on top of <code>Auth</code>. It reads
                the role that <code>Auth</code> stored in the context and checks if it matches one
                of the allowed roles. If not, it returns a 403 Forbidden.
              </p>
            </div>

            <CodeBlock language="go" filename="middleware/auth.go (RequireRole)" code={`// RequireRole checks if the authenticated user has one of the required roles.
// Uses variadic args — you can pass one or more roles.
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Read the role that Auth middleware stored in context
        userRole, exists := c.Get("user_role")
        if !exists {
            c.JSON(401, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Not authenticated"}})
            c.Abort()
            return
        }

        role := userRole.(string)

        // Check if user's role matches any allowed role
        for _, r := range roles {
            if role == r {
                c.Next() // Role matches — continue
                return
            }
        }

        // No match — forbidden
        c.JSON(403, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "You do not have permission"}})
        c.Abort()
    }
}

// Usage in routes:
// admin.Use(middleware.RequireRole("ADMIN"))           // Only admins
// editor.Use(middleware.RequireRole("ADMIN", "EDITOR")) // Admins + editors`} />

            <div className="prose-grit mb-10">
              <h3>How c.Set / c.Get Passes Data</h3>
              <p>
                The <code>*gin.Context</code> acts as a shared data bag between middleware and handlers
                in the same request. Middleware uses <code>c.Set()</code> to store data, and handlers
                use <code>c.Get()</code> to retrieve it. This is how the user object flows from the
                auth middleware to any handler:
              </p>
            </div>

            <CodeBlock language="go" filename="context data flow" code={`// In Auth middleware:
c.Set("user", user)         // Store the full user struct
c.Set("user_id", user.ID)   // Store just the ID (convenience)
c.Set("user_role", user.Role)

// In any handler on a protected route:
func (h *UserHandler) GetProfile(c *gin.Context) {
    // Get the user object stored by middleware
    userData, _ := c.Get("user")
    user := userData.(models.User) // Type assert from any → User

    // Or get just the ID
    userID, _ := c.Get("user_id")
    id := userID.(uint)

    c.JSON(200, gin.H{"data": user})
}`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit scaffolds three route groups: public (no auth), protected (Auth middleware),
                and admin (Auth + RequireRole). You can add custom role-restricted groups with
                <code>grit generate resource --roles ADMIN,EDITOR</code>. The <code>grit add role MODERATOR</code>
                command adds a new role across the entire codebase (Go constants, Zod schemas, TypeScript types,
                sidebar visibility, form options) in one step.
              </p>
            </div>

            <PlaygroundChallenge
              title="Role Checks as Middleware"
              description="Write RequireRole once and apply it per route. The check belongs in one place, not repeated at the top of every handler where one omission becomes a hole."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: RequireRole
// 1. Assume an earlier middleware has put a role on the context.
//    For this exercise read it from the X-Role header instead.
// 2. Write RequireRole(allowed ...string) func(http.Handler) http.Handler
//      - 401 when no role is present at all
//      - 403 when a role is present but not in allowed
//      - otherwise call next
// 3. Build a mux with:
//      GET  /products         -> open to everyone
//      POST /products         -> RequireRole("ADMIN", "EDITOR")
//      DELETE /products/{id}  -> RequireRole("ADMIN")
// 4. Exercise every combination and print method, role and status
// 5. Note which failures are 401 (who are you?) and which are 403
//    (I know who you are, and no)

func main() {
	fmt.Println(http.StatusForbidden, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func RequireRole(allowed ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := r.Header.Get("X-Role")
			if role == "" {
				http.Error(w, "unauthenticated", http.StatusUnauthorized)
				return
			}
			for _, a := range allowed {
				if a == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "forbidden", http.StatusForbidden)
		})
	}
}

func ok(msg string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, msg)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("GET /products", ok("list"))
	mux.Handle("POST /products", RequireRole("ADMIN", "EDITOR")(ok("created")))
	mux.Handle("DELETE /products/{id}", RequireRole("ADMIN")(ok("deleted")))

	cases := []struct{ method, path, role string }{
		{"GET", "/products", ""},
		{"POST", "/products", ""},
		{"POST", "/products", "USER"},
		{"POST", "/products", "EDITOR"},
		{"DELETE", "/products/1", "EDITOR"},
		{"DELETE", "/products/1", "ADMIN"},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(c.method, c.path, nil)
		if c.role != "" {
			req.Header.Set("X-Role", c.role)
		}
		mux.ServeHTTP(rec, req)
		fmt.Printf("%-7s %-14s role=%-7q %d\\n", c.method, c.path, c.role, rec.Code)
	}
}`}
            />

            <PlaygroundChallenge
              title="Ownership Beats Roles"
              description="Some rules cannot be expressed as a role: a user may edit their own post but not someone else's. Combine a role check with an ownership check and get the order right."
              challenge={`package main

import "fmt"

type Post struct {
	ID       int
	AuthorID int
	Title    string
}

// Challenge: role plus ownership
// 1. Write canEdit(userID int, role string, p Post) (bool, string)
//    returning whether the edit is allowed and a short reason:
//      - role "ADMIN"                     -> true, "admin override"
//      - p.AuthorID == userID             -> true, "owner"
//      - role "EDITOR" and p is a draft   -> true, "editor on draft"
//        (treat a Title starting with "DRAFT:" as a draft)
//      - otherwise                        -> false, "not permitted"
// 2. Check ownership BEFORE the editor rule so an owner is never refused
// 3. Run the table below and print user, role, post and the decision
// 4. Bonus: add a locked post nobody but an admin may edit

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"fmt"
	"strings"
)

type Post struct {
	ID       int
	AuthorID int
	Title    string
	Locked   bool
}

func canEdit(userID int, role string, p Post) (bool, string) {
	// Admin first: the override that beats every other rule
	if role == "ADMIN" {
		return true, "admin override"
	}
	// Locked posts stop here for everyone else
	if p.Locked {
		return false, "post is locked"
	}
	// Ownership before role, so an owner is never turned away
	if p.AuthorID == userID {
		return true, "owner"
	}
	if role == "EDITOR" && strings.HasPrefix(p.Title, "DRAFT:") {
		return true, "editor on draft"
	}
	return false, "not permitted"
}

func main() {
	published := Post{ID: 1, AuthorID: 7, Title: "Shipping Go"}
	draft := Post{ID: 2, AuthorID: 7, Title: "DRAFT: Generics"}
	locked := Post{ID: 3, AuthorID: 7, Title: "Policy", Locked: true}

	cases := []struct {
		userID int
		role   string
		post   Post
	}{
		{7, "USER", published},
		{9, "USER", published},
		{9, "EDITOR", published},
		{9, "EDITOR", draft},
		{9, "ADMIN", locked},
		{7, "USER", locked},
	}

	for _, c := range cases {
		allowed, reason := canEdit(c.userID, c.role, c.post)
		fmt.Printf("user=%d role=%-7s post=%-2d -> %-5t (%s)\\n",
			c.userID, c.role, c.post.ID, allowed, reason)
	}
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 22. Important Packages */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="important-packages">22. Important Packages</h2>
              <p>
                These are the Go packages used in every Grit backend. You don&apos;t need to memorize them --
                they are all pre-configured when you scaffold a project. But knowing what they do helps
                you understand the generated code.
              </p>
            </div>

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-8">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Package</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">What It Does</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/20">
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/gin-gonic/gin</td><td className="px-4 py-2.5 text-muted-foreground">HTTP framework — router, middleware, JSON binding, validation</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">gorm.io/gorm</td><td className="px-4 py-2.5 text-muted-foreground">ORM — maps Go structs to database tables, chainable queries</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">gorm.io/driver/postgres</td><td className="px-4 py-2.5 text-muted-foreground">PostgreSQL driver for GORM</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/golang-jwt/jwt/v5</td><td className="px-4 py-2.5 text-muted-foreground">JWT creation and validation for authentication</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">golang.org/x/crypto/bcrypt</td><td className="px-4 py-2.5 text-muted-foreground">Password hashing (used in User model&apos;s BeforeCreate hook)</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/joho/godotenv</td><td className="px-4 py-2.5 text-muted-foreground">Load .env files into environment variables</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/redis/go-redis/v9</td><td className="px-4 py-2.5 text-muted-foreground">Redis client for caching and session storage</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/hibiken/asynq</td><td className="px-4 py-2.5 text-muted-foreground">Background job queue and cron scheduler (Redis-backed)</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/aws/aws-sdk-go-v2</td><td className="px-4 py-2.5 text-muted-foreground">S3-compatible file storage (AWS S3, Cloudflare R2, MinIO)</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/resend/resend-go/v2</td><td className="px-4 py-2.5 text-muted-foreground">Transactional email service</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/disintegration/imaging</td><td className="px-4 py-2.5 text-muted-foreground">Image resizing and thumbnail generation</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/MUKE-coder/gorm-studio</td><td className="px-4 py-2.5 text-muted-foreground">Visual database browser embedded at /studio</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">github.com/MUKE-coder/sentinel</td><td className="px-4 py-2.5 text-muted-foreground">Security suite — WAF, rate limiting, threat dashboard</td></tr>
                </tbody>
              </table>
            </div>

            <div className="prose-grit mb-10">
              <p>
                The standard library packages you will encounter most often:
              </p>
            </div>

            <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-8">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border/30 bg-accent/20">
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Package</th>
                    <th className="text-left px-4 py-2.5 font-medium text-foreground/80">What It Does</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/20">
                  <tr><td className="px-4 py-2.5 font-mono text-xs">fmt</td><td className="px-4 py-2.5 text-muted-foreground">Formatted printing and string formatting</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">net/http</td><td className="px-4 py-2.5 text-muted-foreground">HTTP status codes (http.StatusOK, http.StatusNotFound, etc.)</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">os</td><td className="px-4 py-2.5 text-muted-foreground">Environment variables, file operations, process exit</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">time</td><td className="px-4 py-2.5 text-muted-foreground">Timestamps, durations, token expiry</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">strings</td><td className="px-4 py-2.5 text-muted-foreground">String manipulation (Split, Contains, ToLower, etc.)</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">strconv</td><td className="px-4 py-2.5 text-muted-foreground">String-to-number conversion (Atoi for page params)</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">log</td><td className="px-4 py-2.5 text-muted-foreground">Logging (log.Println, log.Fatalf)</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">errors</td><td className="px-4 py-2.5 text-muted-foreground">Error creation and wrapping</td></tr>
                  <tr><td className="px-4 py-2.5 font-mono text-xs">math</td><td className="px-4 py-2.5 text-muted-foreground">math.Ceil for pagination page count</td></tr>
                </tbody>
              </table>
            </div>

            <div className="prose-grit mb-10">
              <p>
                Reading the list is one thing; the two programs below actually use them. The first
                covers the text and number packages you reach for on nearly every request, the second
                the time and encoding ones that show up the moment you touch a database or return
                JSON.
              </p>
            </div>

            <CodeBlock language="go" filename="stdlib_text.go" code={`package main

import (
    "fmt"
    "sort"
    "strconv"
    "strings"
)

func main() {
    // strings — the workhorse for anything user-supplied
    raw := "  Laptop, Mouse , Keyboard,  "
    parts := strings.Split(strings.TrimSpace(raw), ",")
    var items []string
    for _, p := range parts {
        if p = strings.TrimSpace(p); p != "" {
            items = append(items, p)
        }
    }
    fmt.Println(items, len(items))

    fmt.Println(strings.ToLower("ADA@Example.com"))
    fmt.Println(strings.Contains("ada@example.com", "@"))
    fmt.Println(strings.HasPrefix("Bearer abc123", "Bearer "))
    fmt.Println(strings.TrimPrefix("Bearer abc123", "Bearer "))
    fmt.Println(strings.Join(items, " | "))
    fmt.Println(strings.ReplaceAll("a/b/c", "/", "-"))

    // strings.Builder — the efficient way to assemble a string in a loop
    var b strings.Builder
    for i, item := range items {
        if i > 0 {
            b.WriteString(", ")
        }
        b.WriteString(item)
    }
    fmt.Println(b.String())

    // strconv — string <-> number, always with an error to check
    n, err := strconv.Atoi("42")
    fmt.Println(n, err)
    fmt.Println(strconv.Itoa(99) + "!")
    fmt.Println(strconv.FormatFloat(3.14159, 'f', 2, 64))
    fmt.Println(strconv.Quote(\`he said "hi"\`))

    // sort — for deterministic output
    prices := []float64{29.99, 5.00, 12.50}
    sort.Float64s(prices)
    fmt.Println(prices)

    sort.Slice(items, func(i, j int) bool { return items[i] < items[j] })
    fmt.Println(items)
}`} />

            <CodeBlock language="go" filename="stdlib_time_json.go" code={`package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "time"
)

type Event struct {
    Name      string        \`json:"name"\`
    StartsAt  time.Time     \`json:"starts_at"\`
    Duration  time.Duration \`json:"duration_ns"\`
    Cancelled bool          \`json:"cancelled"\`
}

func main() {
    // time — parsing uses a reference layout, not format codes
    start, err := time.Parse(time.RFC3339, "2026-03-01T09:30:00Z")
    if err != nil {
        fmt.Println("parse error:", err)
        return
    }
    fmt.Println("start:", start.Format("Mon 2 Jan 2006 15:04"))

    end := start.Add(90 * time.Minute)
    fmt.Println("end:  ", end.Format(time.RFC3339))
    fmt.Println("lasts:", end.Sub(start))
    fmt.Println("after?", end.After(start))

    // Durations are typed, so this reads as what it is
    deadline := 30 * time.Second
    fmt.Println("deadline:", deadline, "in ms:", deadline.Milliseconds())

    // encoding/json — marshal, then unmarshal back
    ev := Event{Name: "Launch", StartsAt: start, Duration: 90 * time.Minute}
    out, _ := json.MarshalIndent(ev, "", "  ")
    fmt.Println(string(out))

    var back Event
    if err := json.Unmarshal(out, &back); err != nil {
        fmt.Println("unmarshal error:", err)
        return
    }
    fmt.Println("round trip:", back.Name, back.StartsAt.Year())

    // Bad JSON gives you a typed error, not a panic
    var broken Event
    err = json.Unmarshal([]byte(\`{"starts_at": 12345}\`), &broken)
    var typeErr *json.UnmarshalTypeError
    if errors.As(err, &typeErr) {
        fmt.Printf("field %q wanted %s\\n", typeErr.Field, typeErr.Type)
    }
}`} />

            <PlaygroundChallenge
              title="Cleaning User Input"
              description="Normalise a messy list of tags with strings and sort — trimming, lowercasing, dropping blanks and removing duplicates, which is most of what validation does in practice."
              challenge={`package main

import "fmt"

// Challenge: normalise a tag list
// 1. Start from: "  Go, backend ,GO,  , API,api , go "
// 2. Split on ",", then for each part:
//      - strings.TrimSpace it
//      - skip it if empty
//      - strings.ToLower it
// 3. Remove duplicates using a map[string]bool as a set
// 4. sort.Strings the result
// 5. Print the count and strings.Join(tags, ", ")
//    Expect: 3 tags -> api, backend, go
// 6. Bonus: reject any tag longer than 20 characters and report it

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"fmt"
	"sort"
	"strings"
)

func normaliseTags(raw string) ([]string, []string) {
	seen := map[string]bool{}
	var tags, rejected []string

	for _, part := range strings.Split(raw, ",") {
		t := strings.ToLower(strings.TrimSpace(part))
		if t == "" {
			continue
		}
		if len(t) > 20 {
			rejected = append(rejected, t)
			continue
		}
		if seen[t] {
			continue // duplicate
		}
		seen[t] = true
		tags = append(tags, t)
	}

	sort.Strings(tags)
	return tags, rejected
}

func main() {
	raw := "  Go, backend ,GO,  , API,api , go , a-very-long-tag-that-is-too-long"

	tags, rejected := normaliseTags(raw)
	fmt.Printf("%d tags -> %s\\n", len(tags), strings.Join(tags, ", "))
	fmt.Println("rejected:", rejected)
}`}
            />

            <PlaygroundChallenge
              title="Times and Durations"
              description="Parse a timestamp, do arithmetic on it, and format it for a response. Getting the reference layout right is the one piece of Go time that trips everybody up."
              challenge={`package main

import "fmt"

// Challenge: time arithmetic
// 1. Parse "2026-03-01T09:30:00Z" with time.Parse(time.RFC3339, ...)
// 2. Build a []time.Time of three sessions: start, +45m, +2h30m
// 3. For each, print the RFC3339 form and a human form using the
//    reference layout "Mon 2 Jan 2006 15:04"
//    (Go formats by example: 01/02 03:04:05PM '06 -0700, not %Y-%m-%d)
// 4. Compute and print the gap between the first and last with Sub()
// 5. Write func isExpired(exp time.Time, now time.Time) bool and test it
//    with a time before and after
// 6. Bonus: parse "45m" with time.ParseDuration and add it as well
//
// Do not use time.Now() here — the playground clock is fixed at 2009,
// which makes any test against "now" confusing.

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"fmt"
	"time"
)

func isExpired(exp, now time.Time) bool {
	return !now.Before(exp)
}

func main() {
	start, err := time.Parse(time.RFC3339, "2026-03-01T09:30:00Z")
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	sessions := []time.Time{
		start,
		start.Add(45 * time.Minute),
		start.Add(2*time.Hour + 30*time.Minute),
	}

	for i, s := range sessions {
		// Go formats by example, using the reference date
		fmt.Printf("%d %s  %s\\n", i+1, s.Format(time.RFC3339), s.Format("Mon 2 Jan 2006 15:04"))
	}

	gap := sessions[len(sessions)-1].Sub(sessions[0])
	fmt.Println("first to last:", gap)

	// Expiry checks compare two explicit times, never a hidden "now"
	fmt.Println("expired at +1h? ", isExpired(sessions[2], start.Add(time.Hour)))
	fmt.Println("expired at +3h? ", isExpired(sessions[2], start.Add(3*time.Hour)))

	// Bonus
	d, _ := time.ParseDuration("45m")
	fmt.Println("start + 45m:", start.Add(d).Format(time.RFC3339))
}`}
            />
            {/* ─────────────────────────────────────────────────── */}
            {/* 23. Putting It Together */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="putting-it-together">23. Putting It Together</h2>
              <p>
                Now you understand all the Go concepts that power a Grit backend. Here is how they
                connect in the request lifecycle. When an HTTP request hits your API, it flows
                through a predictable chain:
              </p>
              <ol>
                <li><strong>main.go</strong> -- loads config, connects to the database, initializes services, starts the server</li>
                <li><strong>routes.go</strong> -- matches the URL to a handler, runs middleware (auth, CORS, logging)</li>
                <li><strong>Auth middleware</strong> -- extracts JWT, validates token, loads user from DB, sets <code>c.Set(&quot;user&quot;, ...)</code></li>
                <li><strong>RequireRole middleware</strong> -- checks <code>c.Get(&quot;user_role&quot;)</code> against allowed roles</li>
                <li><strong>handler</strong> -- parses the request, validates input with struct tags, calls the service</li>
                <li><strong>service</strong> -- contains business logic, uses GORM to query the database, returns (result, error)</li>
                <li><strong>handler</strong> -- checks the error, sends the JSON response with the correct status code</li>
              </ol>
            </div>

            <CodeBlock language="bash" filename="request lifecycle" code={`GET /api/products/42
        │
        ▼
┌─── main.go ───────────────────────────────┐
│ cfg := config.Load()                      │
│ db  := database.Connect(cfg)              │
│ svc := &services.ProductService{DB: db}   │
│ routes.Setup(db, cfg, svc)                │
└───────────────────────────────────────────┘
        │
        ▼
┌─── routes.go ─────────────────────────────┐
│ protected := r.Group("/api")              │
│ protected.Use(middleware.Auth(db, auth))   │
│ protected.GET("/products/:id", h.GetByID) │
└───────────────────────────────────────────┘
        │
        ▼
┌─── middleware/auth.go ────────────────────┐
│ token := c.GetHeader("Authorization")     │
│ claims := authService.ValidateToken(token)│
│ user := db.First(&user, claims.UserID)    │
│ c.Set("user", user)                       │
│ c.Set("user_role", user.Role)             │
│ c.Next()                                  │
└───────────────────────────────────────────┘
        │
        ▼
┌─── handlers/product.go ──────────────────┐
│ id := c.Param("id")                      │
│ product, err := svc.GetByID(id)          │
│ if err != nil { c.JSON(404, ...) }       │
│ c.JSON(200, gin.H{"data": product})      │
└───────────────────────────────────────────┘
        │
        ▼
┌─── services/product.go ──────────────────┐
│ func (s *ProductService) GetByID(id) {   │
│     var product models.Product            │
│     err := s.DB.Preload("Category").      │
│         First(&product, id).Error         │
│     return product, err                   │
│ }                                         │
└───────────────────────────────────────────┘`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                This entire flow is generated for you. When you run <code>grit generate resource Product</code>,
                it creates the model, service, handler, routes, and injects everything into the right
                files. You get a fully working CRUD API with pagination, filtering, authentication,
                and role-based access in seconds. Understanding this flow helps you customize the
                generated code and build features beyond basic CRUD.
              </p>
            </div>

            <div className="prose-grit mb-10">
              <p>
                The diagram above is the Gin version. Here is the same pipeline as a program you can
                actually run: middleware, a handler, a service and a store, wired together and
                exercised with four requests. Strip the framework away and the whole architecture is
                about a hundred lines — which is the point. Gin and GORM save you typing; they do
                not change the shape.
              </p>
            </div>

            <CodeBlock language="go" filename="pipeline.go" code={`package main

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strconv"
    "strings"
)

// ── models ──────────────────────────────────────────────────────────────
type Product struct {
    ID       int     \`json:"id"\`
    Name     string  \`json:"name"\`
    Price    float64 \`json:"price"\`
    Category string  \`json:"category"\`
}

var ErrNotFound = errors.New("not found")

// ── store (stands in for GORM) ──────────────────────────────────────────
type Store struct{ rows map[int]Product }

func (s *Store) First(id int) (Product, error) {
    p, ok := s.rows[id]
    if !ok {
        return Product{}, fmt.Errorf("product %d: %w", id, ErrNotFound)
    }
    return p, nil
}

// ── service: business rules, no HTTP ────────────────────────────────────
type ProductService struct{ store *Store }

func (svc *ProductService) GetByID(id int) (Product, error) {
    if id <= 0 {
        return Product{}, fmt.Errorf("id must be positive: %w", ErrNotFound)
    }
    return svc.store.First(id)
}

// ── transport helpers ───────────────────────────────────────────────────
func writeJSON(w http.ResponseWriter, status int, body any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(body)
}

// ── middleware ──────────────────────────────────────────────────────────
type ctxKey string

const ctxRole ctxKey = "role"

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !strings.HasPrefix(token, "Bearer ") {
            writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
            return
        }
        // A real system decodes claims here; this keeps the shape
        role := strings.TrimPrefix(token, "Bearer ")
        ctx := context.WithValue(r.Context(), ctxRole, role)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func RequireRole(allowed string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if roleFrom(r) != allowed {
                writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// ── handler: read, validate, call the service, respond ──────────────────
type ProductHandler struct{ svc *ProductService }

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.PathValue("id"))
    if err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a number"})
        return
    }

    product, err := h.svc.GetByID(id)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
            return
        }
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{"data": product})
}

func main() {
    store := &Store{rows: map[int]Product{
        42: {ID: 42, Name: "Laptop", Price: 999, Category: "electronics"},
    }}
    h := &ProductHandler{svc: &ProductService{store: store}}

    mux := http.NewServeMux()
    mux.Handle("GET /api/products/{id}", Auth(RequireRole("ADMIN")(http.HandlerFunc(h.GetByID))))

    requests := []struct{ path, token string }{
        {"/api/products/42", ""},          // no token   -> 401
        {"/api/products/42", "USER"},      // wrong role -> 403
        {"/api/products/42", "ADMIN"},     // found      -> 200
        {"/api/products/999", "ADMIN"},    // missing    -> 404
    }
    for _, req := range requests {
        rec := httptest.NewRecorder()
        r := httptest.NewRequest("GET", req.path, nil)
        if req.token != "" {
            r.Header.Set("Authorization", "Bearer "+req.token)
        }
        mux.ServeHTTP(rec, r)
        fmt.Printf("%-20s token=%-6q %d %s", req.path, req.token, rec.Code, rec.Body.String())
    }
}

// Kept at the bottom so the pipeline above reads top to bottom
func roleFrom(r *http.Request) string {
    v, _ := r.Context().Value(ctxRole).(string)
    return v
}`} />

            <PlaygroundChallenge
              title="Wire the Whole Pipeline"
              description="Build the full request path yourself — middleware, handler, service, store — and prove each layer with a request that exercises it."
              challenge={`package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

// Challenge: assemble the pipeline end to end
// 1. Store:   map[int]Order with First(id) (Order, error) returning ErrNotFound
// 2. Service: OrderService.GetByID(id) rejecting id <= 0, otherwise the store
// 3. Handler: read r.PathValue("id"), 400 on a non-number, 404 on ErrNotFound,
//             200 with {"data": order} otherwise
// 4. Middleware: Auth requiring "Bearer " on Authorization, else 401
// 5. Register "GET /api/orders/{id}" with the middleware wrapped around it
// 6. Fire four requests and print the status for each:
//      no token -> 401, bad id -> 400, missing order -> 404, real order -> 200
//
// Keep the layers apart: no http in the service, no business rules in the
// handler. That separation is the entire lesson of this page.

type Order struct {
	ID    int     \`json:"id"\`
	Total float64 \`json:"total"\`
}

func main() {
	fmt.Println(http.StatusOK, httptest.NewRecorder().Code)
}`}
              solution={`package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type Order struct {
	ID    int     \`json:"id"\`
	Total float64 \`json:"total"\`
}

var ErrNotFound = errors.New("not found")

// ── store ───────────────────────────────────────────────────────────────
type Store struct{ rows map[int]Order }

func (s *Store) First(id int) (Order, error) {
	o, ok := s.rows[id]
	if !ok {
		return Order{}, fmt.Errorf("order %d: %w", id, ErrNotFound)
	}
	return o, nil
}

// ── service: rules only, no HTTP ────────────────────────────────────────
type OrderService struct{ store *Store }

func (svc *OrderService) GetByID(id int) (Order, error) {
	if id <= 0 {
		return Order{}, fmt.Errorf("id must be positive: %w", ErrNotFound)
	}
	return svc.store.First(id)
}

// ── transport ───────────────────────────────────────────────────────────
func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

type OrderHandler struct{ svc *OrderService }

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id must be a number"})
		return
	}
	order, err := h.svc.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": order})
}

func main() {
	store := &Store{rows: map[int]Order{7: {ID: 7, Total: 129.99}}}
	h := &OrderHandler{svc: &OrderService{store: store}}

	mux := http.NewServeMux()
	mux.Handle("GET /api/orders/{id}", Auth(http.HandlerFunc(h.GetByID)))

	cases := []struct{ path, token string }{
		{"/api/orders/7", ""},
		{"/api/orders/abc", "x"},
		{"/api/orders/999", "x"},
		{"/api/orders/7", "x"},
	}
	for _, c := range cases {
		rec := httptest.NewRecorder()
		r := httptest.NewRequest("GET", c.path, nil)
		if c.token != "" {
			r.Header.Set("Authorization", "Bearer "+c.token)
		}
		mux.ServeHTTP(rec, r)
		fmt.Printf("%-18s token=%-3q %d %s", c.path, c.token, rec.Code, rec.Body.String())
	}
}`}
            />

            <PlaygroundChallenge
              title="Trace a Request Through the Layers"
              description="Instrument every layer so one request prints the path it took. Seeing the order once is worth more than reading the diagram five times."
              challenge={`package main

import "fmt"

// Challenge: make the flow visible
// 1. Build the same four layers, but have each one append to a []string trace:
//      "cors" -> "logger" -> "auth" -> "handler" -> "service" -> "store"
// 2. Run one successful request and print the trace in order
// 3. Run one that fails auth and print that trace — it should stop at "auth"
//    and never reach the handler
// 4. Print both traces side by side so the short-circuit is obvious
//
// This is what c.Abort() means in practice: the layers after it never run.

func main() {
	fmt.Println("replace me")
}`}
              solution={`package main

import (
	"fmt"
	"strings"
)

type request struct {
	token string
	trace []string
}

func (r *request) mark(layer string) { r.trace = append(r.trace, layer) }

// Each layer either passes the request on, or stops the chain
func cors(r *request, next func(*request) int) int {
	r.mark("cors")
	return next(r)
}

func logger(r *request, next func(*request) int) int {
	r.mark("logger")
	return next(r)
}

func auth(r *request, next func(*request) int) int {
	r.mark("auth")
	if r.token == "" {
		return 401 // abort: nothing below this line runs
	}
	return next(r)
}

func handler(r *request) int {
	r.mark("handler")
	return service(r)
}

func service(r *request) int {
	r.mark("service")
	return store(r)
}

func store(r *request) int {
	r.mark("store")
	return 200
}

func serve(r *request) int {
	return cors(r, func(r *request) int {
		return logger(r, func(r *request) int {
			return auth(r, handler)
		})
	})
}

func main() {
	ok := &request{token: "abc"}
	status := serve(ok)
	fmt.Printf("%d  %s\\n", status, strings.Join(ok.trace, " -> "))

	denied := &request{token: ""}
	status = serve(denied)
	fmt.Printf("%d  %s\\n", status, strings.Join(denied.trace, " -> "))

	fmt.Println("\\nthe denied request never reached:",
		len(ok.trace)-len(denied.trace), "layers")
}`}
            />


            {/* Navigation footer */}
            <div className="mt-16 pt-8 border-t border-border/40 flex items-center justify-between">
              <div />
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/prerequisites/nextjs" className="gap-1.5">
                  Next.js &amp; React for Grit Developers
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>

          </div>
        </div>
      </main>
    </div>
  );
}
