import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { TableOfContents } from "@/components/table-of-contents";
import { PlaygroundChallenge } from "@/components/playground-challenge";

const tocItems = [
  { id: 'go-basics', label: 'Go Basics' },
  { id: 'variables-types', label: 'Variables & Types' },
  { id: 'structs-tags', label: 'Structs & Tags' },
  { id: 'functions-errors', label: 'Functions & Error Handling' },
  { id: 'methods', label: 'Methods' },
  { id: 'slices-maps', label: 'Slices & Maps' },
  { id: 'interfaces', label: 'Interfaces' },
  { id: 'pointers', label: 'Pointers' },
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
              title="Go Basics"
              description="Print your name and the current year using fmt.Printf with format verbs (%s and %d)."
              challenge={`package main

import "fmt"

func main() {
	// Challenge: Print your name and the current year
	// 1. Create a string variable "name" with your name
	// 2. Create an int variable "year" with 2026
	// 3. Use fmt.Printf to print: "Hi, I'm <name> and it's <year>!"
	//    Hint: use %s for strings and %d for integers
	// 4. Bonus: Print the type of each variable using %T

}`}
              solution={`package main

import "fmt"

func main() {
	name := "Alice"
	year := 2026

	fmt.Printf("Hi, I'm %s and it's %d!\\n", name, year)

	// Bonus: Print the type of each variable using %T
	fmt.Printf("name is %T, year is %T\\n", name, year)
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 3. Structs & Tags */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="structs-tags">3. Structs & Tags</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 5. Methods */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="methods">5. Methods</h2>
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

    // Maps
    user := map[string]any{
        "id":    1,
        "name":  "Alice",
        "email": "alice@example.com",
    }

    for key, value := range user {
        fmt.Printf("%s = %v\\n", key, value)
    }

    // Access a single value
    fmt.Println("Name:", user["name"])
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 7. Interfaces */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="interfaces">7. Interfaces</h2>
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

            <CodeBlock language="go" filename="interfaces.go" code={`package main

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

// Works with ANY Notifier
func alert(n Notifier, user string) {
    n.Send(user, "Your report is ready")
}

func main() {
    email := &EmailNotifier{From: "noreply@app.com"}
    slack := &SlackNotifier{Channel: "alerts"}

    alert(email, "alice@example.com")
    alert(slack, "alice")
}`} />

            <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
              <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
              <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                Grit services can implement interfaces for testability. For example, you could
                define a <code>UserService</code> interface with methods like <code>GetByID</code>,
                <code>Create</code>, and <code>Delete</code>, then swap in a mock implementation
                during tests. The built-in mailer and storage services also follow this pattern.
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 8. Pointers */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="pointers">8. Pointers</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 9. Goroutines & Channels */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="goroutines-channels">9. Goroutines & Channels</h2>
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

            <CodeBlock language="go" filename="goroutines.go" code={`package main

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
}`} />

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

            {/* ─────────────────────────────────────────────────── */}
            {/* 10. Packages & Project Structure */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="packages-structure">10. Packages & Project Structure</h2>
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
            {/* 11. Environment Variables */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="env-variables">11. Environment Variables</h2>
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
            {/* 12. Gin Framework */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="gin-framework">12. Gin Framework</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 13. Middleware */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="middleware">13. Middleware</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 14. CORS */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="cors">14. CORS</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 15. Handlers */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="handlers">15. Handlers</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 16. Services & The Service Pattern */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="services">16. Services & The Service Pattern</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 17. GORM In Depth */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="gorm-in-depth">17. GORM In Depth</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 18. Migrations & Seeding */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="migrations-seeding">18. Migrations & Seeding</h2>
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
go run cmd/migrate/main.go

# Fresh migration (drop all tables + recreate)
go run cmd/migrate/main.go --fresh`} />

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
go run cmd/seed/main.go`} />

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

            {/* ─────────────────────────────────────────────────── */}
            {/* 19. JWT & Authentication */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="jwt-auth">19. JWT & Authentication</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 20. RBAC & Middleware */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="rbac-middleware">20. RBAC & Middleware</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 21. Important Packages */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="important-packages">21. Important Packages</h2>
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

            {/* ─────────────────────────────────────────────────── */}
            {/* 22. Putting It Together */}
            {/* ─────────────────────────────────────────────────── */}
            <div className="prose-grit mb-10">
              <h2 id="putting-it-together">22. Putting It Together</h2>
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
