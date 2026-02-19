import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";

export default function GoForGritPage() {
  return (
    <div className="min-h-screen bg-background">
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 1. Go Basics */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="go-basics">Go Basics</h2>
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

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">main.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

import "fmt"

func main() {
    fmt.Println("Hello, Grit!")
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                The entry point for every Grit backend is <code>apps/api/cmd/server/main.go</code>.
                This file initializes the database connection, sets up middleware, registers routes,
                and starts the Gin HTTP server. You rarely edit it directly -- the code generator
                handles injecting new routes and models automatically.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 2. Variables & Types */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="variables-types">Variables & Types</h2>
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

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">variables.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

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
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                You will see <code>:=</code> everywhere in handlers and services. Config values loaded
                from <code>.env</code> are stored in typed struct fields (like <code>Port int</code>,
                <code>JWTSecret string</code>). Constants are used for role names
                (<code>RoleAdmin = &quot;ADMIN&quot;</code>) and error codes.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 3. Structs & Tags */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="structs-tags">Structs & Tags</h2>
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

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">models/user.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package models

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
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Every model in <code>internal/models/</code> is a struct with these three tag types.
                When you run <code>grit generate resource Product</code>, the CLI creates a struct
                with properly tagged fields, registers it for migration, and generates the matching
                Zod schema and TypeScript type on the frontend.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 4. Functions & Error Handling */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="functions-errors">Functions & Error Handling</h2>
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

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">errors.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

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
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Every service function in <code>internal/services/</code> returns <code>(result, error)</code>.
                Handlers call services, check for errors, and return the appropriate HTTP response.
                For example, <code>user, err := service.GetUserByID(id)</code> followed by
                an <code>if err != nil</code> block that sends a 404 or 500 JSON response.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 5. Slices & Maps */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="slices-maps">Slices & Maps</h2>
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

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">collections.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

import "fmt"

func main() {
    // Slices
    names := []string{"Alice", "Bob", "Charlie"}
    names = append(names, "Diana")

    for i, name := range names {
        fmt.Printf("%d: %s\n", i, name)
    }

    // Maps
    user := map[string]any{
        "id":    1,
        "name":  "Alice",
        "email": "alice@example.com",
    }

    for key, value := range user {
        fmt.Printf("%s = %v\n", key, value)
    }

    // Access a single value
    fmt.Println("Name:", user["name"])
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                GORM query results are always slices: <code>var users []models.User</code>. Gin JSON
                responses use <code>gin.H{"{}"}</code> which is just a shortcut for <code>map[string]any</code>.
                For example, <code>c.JSON(200, gin.H{"{"}&quot;data&quot;: users, &quot;message&quot;: &quot;success&quot;{"}"})</code>.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 6. Interfaces */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="interfaces">Interfaces</h2>
              <p>
                An interface defines a set of method signatures. Any type that implements all those
                methods <strong>automatically</strong> satisfies the interface -- there is no
                <code>implements</code> keyword. This is called implicit implementation, and it is
                one of Go&apos;s most powerful features.
              </p>
              <p>
                Interfaces enable polymorphism and are essential for testing. You can swap a real
                database service for a mock that implements the same interface, making unit tests
                fast and isolated.
              </p>
            </div>

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">interfaces.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

import "fmt"

// Define an interface
type Notifier interface {
    Send(to string, message string) error
}

// EmailNotifier implements Notifier (implicitly)
type EmailNotifier struct {
    From string
}

func (e *EmailNotifier) Send(to string, message string) error {
    fmt.Printf("Email from %s to %s: %s\n", e.From, to, message)
    return nil
}

// SlackNotifier also implements Notifier
type SlackNotifier struct {
    Channel string
}

func (s *SlackNotifier) Send(to string, message string) error {
    fmt.Printf("Slack #%s -> %s: %s\n", s.Channel, to, message)
    return nil
}

// Works with ANY Notifier
func alert(n Notifier, user string) {
    n.Send(user, "Your report is ready")
}

func main() {
    email := &EmailNotifier{From: "noreply@app.com"}
    slack := &SlackNotifier{Channel: "alerts"}

    alert(email, "alice@example.com")
    alert(slack, "alice")
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit services can implement interfaces for testability. For example, you could
                define a <code>UserService</code> interface with methods like <code>GetByID</code>,
                <code>Create</code>, and <code>Delete</code>, then swap in a mock implementation
                during tests. The built-in mailer and storage services also follow this pattern.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 7. Pointers */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="pointers">Pointers</h2>
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

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">pointers.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

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
    fmt.Println(x) // Still 10 -- the copy was doubled

    doublePointer(&x)
    fmt.Println(x) // Now 20 -- modified through pointer

    // Nil pointer: indicates "no value"
    var name *string = nil
    if name == nil {
        fmt.Println("Name is not set")
    }
}`}
              </pre>
            </div>

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

            {/* ─────────────────────────────────────────────────── */}
            {/* 8. Goroutines & Channels */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="goroutines-channels">Goroutines & Channels</h2>
              <p>
                A goroutine is a lightweight thread managed by the Go runtime. You start one by
                putting <code>go</code> before a function call. Goroutines are extremely cheap --
                you can run thousands simultaneously, unlike OS threads.
              </p>
              <p>
                Channels are the way goroutines communicate. A channel is a typed pipe: one
                goroutine sends a value in, another receives it out. Use <code>sync.WaitGroup</code>
                when you need to wait for multiple goroutines to finish before proceeding.
              </p>
            </div>

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">goroutines.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

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
        wg.Add(1)
        go fetchData(src, &wg) // Run concurrently
    }

    wg.Wait() // Wait for all goroutines to finish
    fmt.Println("All data fetched!")

    // Channel example
    ch := make(chan string)

    go func() {
        ch <- "hello from goroutine" // Send
    }()

    msg := <-ch // Receive
    fmt.Println(msg)
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit&apos;s background job system (powered by asynq) uses goroutines under the hood to
                process tasks like sending emails, resizing images, and running cleanup jobs.
                The Gin web server itself handles each HTTP request in its own goroutine. You
                generally do not need to write goroutine code directly -- asynq and Gin manage
                concurrency for you.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 9. Packages & Project Structure */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="packages-structure">Packages & Project Structure</h2>
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

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">project structure</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`apps/api/
├── cmd/server/
│   └── main.go          # Entry point (package main)
├── internal/
│   ├── config/
│   │   └── config.go    # package config -- Config struct, Load()
│   ├── database/
│   │   └── database.go  # package database -- Connect()
│   ├── models/
│   │   ├── user.go      # package models -- User struct (exported)
│   │   └── upload.go    # package models -- Upload struct (exported)
│   ├── handlers/
│   │   └── auth.go      # package handlers -- Login(), Register()
│   ├── services/
│   │   └── auth.go      # package services -- AuthService
│   ├── middleware/
│   │   └── auth.go      # package middleware -- RequireAuth()
│   └── routes/
│       └── routes.go    # package routes -- Setup()
└── go.mod               # Module definition`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit follows Go&apos;s standard project layout exactly. All application code lives
                inside <code>internal/</code>: models, handlers, services, middleware, routes, and
                config. The <code>cmd/server/</code> directory contains only the entry point.
                When you import a package, you use the full module path:
                <code>import &quot;my-app/apps/api/internal/models&quot;</code>.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 10. HTTP with Gin */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="http-gin">HTTP with Gin</h2>
              <p>
                Gin is Go&apos;s most popular HTTP framework. It provides a fast router, middleware support,
                JSON binding, validation, and route groups. Every handler function receives a
                <code>*gin.Context</code> which holds the request, response, URL parameters, and more.
              </p>
              <p>
                Key Gin methods you will use constantly:
              </p>
              <ul>
                <li><code>c.JSON(status, data)</code> -- send a JSON response</li>
                <li><code>c.ShouldBindJSON(&amp;input)</code> -- parse and validate the request body</li>
                <li><code>c.Param(&quot;id&quot;)</code> -- get a URL parameter like <code>/users/:id</code></li>
                <li><code>c.Query(&quot;page&quot;)</code> -- get a query string parameter like <code>?page=2</code></li>
                <li><code>c.Set() / c.Get()</code> -- pass data through middleware</li>
              </ul>
            </div>

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">server.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type CreateUserInput struct {
    Name  string \`json:"name" binding:"required"\`
    Email string \`json:"email" binding:"required,email"\`
}

func main() {
    r := gin.Default()

    // Middleware
    r.Use(gin.Logger())

    // Route group
    api := r.Group("/api")
    {
        api.GET("/users", getUsers)
        api.POST("/users", createUser)
        api.GET("/users/:id", getUserByID)
    }

    r.Run(":8080")
}

func getUsers(c *gin.Context) {
    page := c.DefaultQuery("page", "1")
    c.JSON(http.StatusOK, gin.H{
        "data":    []string{"Alice", "Bob"},
        "message": "Users fetched",
        "page":    page,
    })
}

func createUser(c *gin.Context) {
    var input CreateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{
        "data":    input,
        "message": "User created",
    })
}

func getUserByID(c *gin.Context) {
    id := c.Param("id")
    c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": id}})
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                All API routes are defined in <code>internal/routes/routes.go</code>. The <code>Setup()</code>
                function creates route groups for public routes, authenticated routes, and admin-only
                routes. Handlers are thin functions that validate input, call a service, and return
                JSON. Middleware like <code>RequireAuth()</code> and <code>RequireRole(&quot;ADMIN&quot;)</code>
                protect route groups.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 11. GORM Basics */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="gorm-basics">GORM Basics</h2>
              <p>
                GORM is Go&apos;s most popular ORM. It maps Go structs to database tables and provides
                a chainable API for queries. You connect to a database, define models as structs,
                run <code>AutoMigrate</code> to create tables, and then use methods like
                <code>Create</code>, <code>Find</code>, <code>Where</code>, <code>Save</code>,
                and <code>Delete</code> for CRUD operations.
              </p>
              <p>
                GORM supports preloading related models (<code>Preload</code>), soft deletes,
                hooks (lifecycle callbacks like <code>BeforeCreate</code>), transactions, and
                pagination with <code>Offset().Limit()</code>.
              </p>
            </div>

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">gorm_crud.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package main

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Product struct {
    ID    uint    \`gorm:"primarykey" json:"id"\`
    Name  string  \`gorm:"size:255;not null" json:"name"\`
    Price float64 \`gorm:"not null" json:"price"\`
}

func main() {
    dsn := "host=localhost user=postgres password=postgres dbname=myapp port=5432"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("failed to connect to database")
    }

    // Create tables from structs
    db.AutoMigrate(&Product{})

    // Create
    product := Product{Name: "Widget", Price: 29.99}
    db.Create(&product) // product.ID is now set

    // Read one
    var found Product
    db.First(&found, product.ID) // Find by primary key

    // Read many with conditions
    var expensive []Product
    db.Where("price > ?", 20.0).Find(&expensive)

    // Update
    db.Model(&found).Update("price", 34.99)

    // Delete (soft delete if DeletedAt field exists)
    db.Delete(&found)

    // Pagination
    var page []Product
    db.Offset(0).Limit(20).Find(&page) // First 20 records

    // Preload relationships
    // db.Preload("Category").Find(&products)
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Every resource generated by <code>grit generate resource</code> gets a service file
                in <code>internal/services/</code> with these exact GORM operations: <code>Create</code>,
                <code>GetAll</code> (with pagination and filtering), <code>GetByID</code> (with
                Preload), <code>Update</code>, and <code>Delete</code>. The database connection is
                established in <code>internal/database/database.go</code> and passed to all services.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 12. Environment Variables */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="env-variables">Environment Variables</h2>
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

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">config/config.go</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`package config

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
}`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit&apos;s config lives in <code>internal/config/config.go</code>. It loads settings for the
                database, Redis, S3 storage, Resend email, AI keys, and more -- all from
                the <code>.env</code> file. The <code>Config</code> struct is created once in <code>main.go</code>
                and passed to every service that needs it. A <code>.env.example</code> file is scaffolded
                with every project to document all available variables.
              </p>
            </div>

            {/* ─────────────────────────────────────────────────── */}
            {/* 13. Putting It Together */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="putting-it-together">Putting It Together</h2>
              <p>
                Now you understand all the Go concepts that power a Grit backend. Here is how they
                connect in the request lifecycle. When an HTTP request hits your API, it flows
                through a predictable chain:
              </p>
              <ol>
                <li><strong>main.go</strong> -- loads config, connects to the database, initializes services, starts the server</li>
                <li><strong>routes.go</strong> -- matches the URL to a handler, runs middleware (auth, CORS, logging)</li>
                <li><strong>middleware</strong> -- validates JWT tokens, checks roles, adds request context</li>
                <li><strong>handler</strong> -- parses the request, validates input with struct tags, calls the service</li>
                <li><strong>service</strong> -- contains business logic, uses GORM to query the database, returns (result, error)</li>
                <li><strong>handler</strong> -- checks the error, sends the JSON response with the correct status code</li>
              </ol>
            </div>

            <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
              <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                <span className="text-[11px] font-mono text-muted-foreground/40">request lifecycle</span>
              </div>
              <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">
{`GET /api/products/42
        │
        ▼
┌─── main.go ───────────────────────────────┐
│ cfg := config.Load()                      │
│ db  := database.Connect(cfg)              │
│ svc := services.NewProductService(db)     │
│ routes.Setup(router, db, cfg, svc)        │
└───────────────────────────────────────────┘
        │
        ▼
┌─── routes.go ─────────────────────────────┐
│ api := router.Group("/api")               │
│ api.Use(middleware.RequireAuth())          │
│ api.GET("/products/:id", handlers.GetProduct) │
└───────────────────────────────────────────┘
        │
        ▼
┌─── middleware/auth.go ────────────────────┐
│ token := c.GetHeader("Authorization")     │
│ claims, err := jwt.Parse(token)           │
│ c.Set("userID", claims.UserID)            │
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
│     err := s.db.Preload("Category").      │
│         First(&product, id).Error         │
│     return product, err                   │
│ }                                         │
└───────────────────────────────────────────┘`}
              </pre>
            </div>

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                This entire flow is generated for you. When you run <code>grit generate resource Product</code>,
                it creates the model, service, handler, routes, and injects everything into the right
                files. You get a fully working CRUD API with pagination, filtering, and validation
                in seconds. Understanding this flow helps you customize the generated code and build
                features beyond basic CRUD.
              </p>
            </div>

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
