import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Testing Your Grit App: Go, Vitest & Playwright',
  description: 'Learn to test every layer of your Grit application — Go unit tests with testify, frontend component tests with Vitest, and end-to-end tests with Playwright.',
}

export default function TestingCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Testing Your Grit App</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">14 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Testing Your Grit App: Go, Vitest & Playwright
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Untested code is broken code — you just don{"'"}t know it yet. In this course, you{"'"}ll learn
            to test every layer of a Grit application: Go unit tests for your API, component tests for
            your React frontend, and end-to-end tests that simulate real users clicking through your app.
            Grit scaffolds test files for you. Your job is to understand them and write more.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: Why Testing Matters ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Why Testing Matters</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every application has bugs. The question is whether you find them before or after your users do.
            Testing is how you catch bugs automatically, before they reach production. There are three main
            types of tests, and they form a pyramid:
          </p>

          <Definition term="Unit Test">
            Tests a single function or method in isolation. Fast to run, easy to write. Example: testing that
            a <Code>HashPassword</Code> function returns a valid bcrypt hash. You write the most of these.
          </Definition>

          <Definition term="Integration Test">
            Tests how multiple parts work together. Example: testing that a POST to <Code>/api/auth/register</Code> creates
            a user in the database AND returns a JWT token. Slower than unit tests, but catches more real-world bugs.
          </Definition>

          <Definition term="End-to-End Test (E2E)">
            Tests the entire application as a real user would — opening a browser, clicking buttons, filling
            forms, and verifying what appears on screen. Slowest to run, but catches the most critical bugs.
            You write the fewest of these.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The testing pyramid: many unit tests at the base, fewer integration tests in the middle, and a
            small number of E2E tests at the top. This gives you fast feedback (unit tests run in milliseconds)
            and high confidence (E2E tests prove the whole system works).
          </p>

          <CodeBlock filename="The Testing Pyramid">
{`        /  E2E  \\           ← Few tests, slow, high confidence
       / Integration \\       ← Some tests, moderate speed
      /   Unit Tests   \\     ← Many tests, fast, isolated
     ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾`}
          </CodeBlock>

          <Note>
            Without tests, every code change is a gamble. You refactor a function and hope nothing broke.
            You add a feature and manually click through 20 pages to verify. With tests, you run one command
            and know in seconds whether everything still works.
          </Note>

          <Challenge number={1} title="Name 3 Things That Can Go Wrong">
            <p>Think about an application with no tests at all. Name 3 specific things that can go wrong
            when a developer pushes new code. For each one, explain how a test would have caught the problem
            before it reached production.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Go Test Basics ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Go Test Basics</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Go has testing built into the language. No external framework needed. Test files end
            with <Code>_test.go</Code>, and test functions start with <Code>Test</Code> followed by a
            capital letter. The <Code>go test</Code> command finds and runs them automatically.
          </p>

          <CodeBlock filename="math.go">
{`package math

func Add(a, b int) int {
    return a + b
}

func IsValidEmail(email string) bool {
    // Must contain @ and at least one dot after @
    atIndex := -1
    for i, ch := range email {
        if ch == '@' {
            atIndex = i
        }
    }
    if atIndex < 1 {
        return false
    }
    domain := email[atIndex+1:]
    for _, ch := range domain {
        if ch == '.' {
            return true
        }
    }
    return false
}`}
          </CodeBlock>

          <CodeBlock filename="math_test.go">
{`package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}

func TestAddNegative(t *testing.T) {
    result := Add(-1, -2)
    if result != -3 {
        t.Errorf("expected -3, got %d", result)
    }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Run all tests in a project with:
          </p>

          <CodeBlock filename="Terminal">
{`go test ./...

# Verbose output (see each test name)
go test -v ./...

# Run a specific test
go test -run TestAdd ./...`}
          </CodeBlock>

          <Definition term="testing.T">
            Go{"'"}s test context, passed to every test function. It provides methods to report failures:
            <Code>t.Errorf()</Code> (fail with message, keep running), <Code>t.Fatal()</Code> (fail and stop
            immediately), <Code>t.Run()</Code> (run a subtest), and <Code>t.Skip()</Code> (skip this test).
          </Definition>

          <Challenge number={2} title="Write an Email Validation Test">
            <p>Write a test function called <Code>TestIsValidEmail</Code> that tests the <Code>IsValidEmail</Code> function.
            Test at least 4 cases: a valid email, an email without @, an email without a dot in the domain,
            and an empty string. Use <Code>t.Errorf</Code> for failures.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Testing with testify ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Testing with testify</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Go{"'"}s built-in testing works, but the assertions are verbose. Writing <Code>if result != expected</Code> for
            every check gets tedious. The <Code>testify</Code> library gives you clean, readable assertions.
          </p>

          <Definition term="testify">
            The most popular Go testing toolkit. It provides two packages: <Code>assert</Code> (reports failure
            but continues the test) and <Code>require</Code> (reports failure and stops immediately). Use
            <Code>assert</Code> when you want to see all failures at once. Use <Code>require</Code> when the
            rest of the test depends on this check passing.
          </Definition>

          <CodeBlock filename="user_test.go">
{`package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
    user := CreateUser("John", "john@test.com")

    // assert: if this fails, keep running the other checks
    assert.Equal(t, "John", user.Name)
    assert.Equal(t, "john@test.com", user.Email)
    assert.NotEmpty(t, user.ID)
    assert.Equal(t, "USER", user.Role)

    // require: if this fails, stop immediately
    require.NotNil(t, user.CreatedAt)
}

func TestCreateUserValidation(t *testing.T) {
    // Empty name should return an error
    user, err := CreateUserSafe("", "john@test.com")
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "name is required")
}`}
          </CodeBlock>

          <Tip>
            Use <Code>require</Code> for setup steps (database connection, creating prerequisite data).
            If setup fails, there{"'"}s no point running the rest of the test. Use <Code>assert</Code> for
            the actual checks so you see all failures at once.
          </Tip>

          <Challenge number={3} title="Test Blog Post Creation">
            <p>Write a test using <Code>assert.Equal</Code> and <Code>assert.NotEmpty</Code> for a blog post
            creation function. Check that the title, slug (auto-generated from title), author_id, and
            created_at fields are set correctly. The slug for {'"'}My First Post{'"'} should be {'"'}my-first-post{'"'}.</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Grit's Scaffolded Tests ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Grit{"'"}s Scaffolded Tests</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you run <Code>grit new myapp</Code>, Grit scaffolds test files automatically. You don{"'"}t
            start from scratch. Here{"'"}s what you get:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">auth_test.go</strong> — 6 tests covering registration (success, validation, duplicate email), login (success, wrong password, unknown email)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">user_test.go</strong> — 4 tests covering auth guard (unauthenticated access blocked), admin list users, user not found (404), get user by ID</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">bench_test.go</strong> — benchmarks for health check, auth login, and auth register endpoints</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            All scaffolded tests use SQLite in-memory databases. This means no Docker, no PostgreSQL, no
            setup — tests run instantly on any machine.
          </p>

          <CodeBlock filename="test_helpers.go">
{`package handler_test

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // Auto-migrate all models
    db.AutoMigrate(&models.User{}, &models.Blog{})
    return db
}

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
    db := setupTestDB(t)
    router := gin.New()
    // Register routes with test database
    return router, db
}`}
          </CodeBlock>

          <Definition term="In-memory Database">
            A database that exists only in RAM, not on disk. Created when the test starts, destroyed when
            the test ends. Extremely fast because there{"'"}s no disk I/O. Each test gets a fresh, empty
            database — no leftover data from previous tests. SQLite supports this with the connection
            string <Code>{'"'}:memory:{'"'}</Code>.
          </Definition>

          <Challenge number={4} title="Find and Run the Scaffolded Tests">
            <p>Navigate to <Code>apps/api/</Code> in your Grit project. Find the test files (they end
            with <Code>_test.go</Code>). Run <Code>go test -v ./...</Code> and observe the output. How
            many tests pass? How many are there total?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Testing Auth Endpoints ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Testing Auth Endpoints</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Authentication is the most critical part of your API. If auth is broken, everything is broken.
            Grit{"'"}s scaffolded <Code>auth_test.go</Code> tests 6 scenarios:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span><strong className="text-foreground">Register success</strong> — valid name, email, password returns 201 + access token</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span><strong className="text-foreground">Register validation</strong> — missing fields return 422 with error details</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">3.</span>
              <span><strong className="text-foreground">Register duplicate</strong> — same email twice returns 409 conflict</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">4.</span>
              <span><strong className="text-foreground">Login success</strong> — correct credentials return 200 + tokens</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">5.</span>
              <span><strong className="text-foreground">Login wrong password</strong> — returns 401 unauthorized</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">6.</span>
              <span><strong className="text-foreground">Login unknown email</strong> — returns 401 unauthorized (same error as wrong password for security)</span>
            </li>
          </ol>

          <CodeBlock filename="auth_test.go">
{`func TestRegisterSuccess(t *testing.T) {
    router, _ := setupTestRouter(t)

    body := map[string]string{
        "name":     "John Doe",
        "email":    "john@test.com",
        "password": "securePass123",
    }
    // POST /api/auth/register with JSON body
    // Assert: status 201
    // Assert: response has "access_token" field
    // Assert: response has "user" object with name and email
}

func TestLoginWrongPassword(t *testing.T) {
    router, _ := setupTestRouter(t)

    // First: register a user
    // Then: attempt login with wrong password
    // Assert: status 401
    // Assert: response has "error" field
}

func TestRegisterDuplicateEmail(t *testing.T) {
    router, _ := setupTestRouter(t)

    // Register once: success (201)
    // Register again with same email: conflict (409)
}`}
          </CodeBlock>

          <Note>
            Notice that login with a wrong password and login with an unknown email both return 401 with
            the same generic message. This is intentional — if the error said {'"'}email not found{'"'} vs
            {'"'}wrong password{'"'}, an attacker could enumerate valid email addresses.
          </Note>

          <Challenge number={5} title="Read the Auth Tests">
            <p>Open <Code>auth_test.go</Code> in your scaffolded project. Read through all 6 test functions.
            For each one, write down: (1) what HTTP method and path it calls, (2) what request body it sends,
            and (3) what status code and response body it asserts.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Testing Generated Resources ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Testing Generated Resources</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After running <Code>grit generate resource Product name:string price:float active:bool</Code>,
            you get CRUD endpoints. Every endpoint needs tests. Here{"'"}s the pattern for testing
            a generated resource:
          </p>

          <CodeBlock filename="product_test.go">
{`func TestCreateProduct(t *testing.T) {
    router, db := setupTestRouter(t)
    token := createTestUser(t, db) // helper: register + get token

    body := map[string]interface{}{
        "name":   "Widget",
        "price":  29.99,
        "active": true,
    }
    // POST /api/products with auth token
    // Assert: status 201
    // Assert: response "data" has name, price, active
    // Assert: product exists in database
}

func TestListProducts(t *testing.T) {
    router, db := setupTestRouter(t)
    token := createTestUser(t, db)

    // Create 5 products
    for i := 0; i < 5; i++ {
        // POST /api/products
    }

    // GET /api/products
    // Assert: status 200
    // Assert: response "data" has 5 items
    // Assert: response "meta" has total, page, page_size
}

func TestGetProduct(t *testing.T) {
    // Create a product, then GET /api/products/:id
    // Assert: 200, correct product returned
}

func TestUpdateProduct(t *testing.T) {
    // Create a product, then PUT /api/products/:id
    // Assert: 200, updated fields match
}

func TestDeleteProduct(t *testing.T) {
    // Create a product, then DELETE /api/products/:id
    // Assert: 200
    // GET /api/products/:id should now return 404
}`}
          </CodeBlock>

          <Tip>
            Create a <Code>createTestUser</Code> helper function that registers a user and returns the
            JWT token. You{"'"}ll use it in almost every test. The scaffolded test files already include
            this helper.
          </Tip>

          <Challenge number={6} title="Write CRUD Tests for a Product">
            <p>Generate a Product resource with <Code>grit generate resource Product name:string price:float
            active:bool category:string</Code>. Then write 5 tests: Create (valid data returns 201),
            List (5 products with pagination), GetByID (correct product), Update (fields change), Delete
            (returns 200, subsequent GET returns 404).</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Benchmarks ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Benchmarks</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Tests tell you if your code is correct. Benchmarks tell you if your code is fast. Go has
            built-in benchmarking — functions named <Code>BenchmarkXxx</Code> that measure how many
            operations per second your code can handle.
          </p>

          <Definition term="Benchmark">
            A performance test that measures how fast code runs. The function receives <Code>*testing.B</Code> and
            runs the code <Code>b.N</Code> times. Go automatically adjusts <Code>b.N</Code> to get a
            reliable measurement. Results show nanoseconds per operation (ns/op).
          </Definition>

          <CodeBlock filename="bench_test.go">
{`package handler_test

import (
    "testing"
    "golang.org/x/crypto/bcrypt"
)

func BenchmarkHashPassword(b *testing.B) {
    for i := 0; i < b.N; i++ {
        bcrypt.GenerateFromPassword([]byte("password123"), 10)
    }
}

func BenchmarkHealthCheck(b *testing.B) {
    router, _ := setupTestRouter(b)
    for i := 0; i < b.N; i++ {
        // GET /api/health
        // Discard response
    }
}

func BenchmarkAuthLogin(b *testing.B) {
    router, db := setupTestRouter(b)
    // Create a user once (outside the loop)
    createTestUser(b, db)

    b.ResetTimer() // Don't count setup time
    for i := 0; i < b.N; i++ {
        // POST /api/auth/login
    }
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Run benchmarks with:
          </p>

          <CodeBlock filename="Terminal">
{`# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkHashPassword ./...

# Run benchmarks with memory allocation stats
go test -bench=. -benchmem ./...`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Typical output looks like:
          </p>

          <CodeBlock filename="Benchmark Output">
{`BenchmarkHealthCheck-8       50000    24350 ns/op    1024 B/op    12 allocs/op
BenchmarkAuthLogin-8          500   2345000 ns/op    8192 B/op    45 allocs/op
BenchmarkHashPassword-8         5 210000000 ns/op      64 B/op     2 allocs/op`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Reading benchmark output: <Code>50000</Code> is how many times the function ran.
            <Code>24350 ns/op</Code> is nanoseconds per operation (lower is faster). Health check
            is fast (~24 microseconds), login is moderate (~2.3 milliseconds because of password
            verification), and bcrypt hashing is intentionally slow (~210 milliseconds).
          </p>

          <Challenge number={7} title="Run the Scaffolded Benchmarks">
            <p>Navigate to your API directory. Run <Code>go test -bench=. -benchmem ./...</Code>. Look at
            the results. Which endpoint is fastest? Which is slowest? Why is <Code>BenchmarkHashPassword</Code>
            so much slower than the others? (Hint: it{"'"}s by design.)</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Frontend Testing with Vitest ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Frontend Testing with Vitest</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your Go API is tested. Now let{"'"}s test the React frontend. Grit uses Vitest — a fast
            test runner built on Vite — combined with React Testing Library for rendering components.
          </p>

          <Definition term="Vitest">
            A fast unit test runner for JavaScript and TypeScript. Like Jest but faster because it uses
            Vite{"'"}s transformation pipeline. Supports TypeScript natively, has hot module reloading for
            tests, and is compatible with Jest{"'"}s API (describe, it, expect).
          </Definition>

          <Definition term="React Testing Library (RTL)">
            A library for testing React components from the user{"'"}s perspective. Instead of testing
            internal state or props, you test what the user sees: rendered text, form fields, buttons.
            The core principle: {'"'}The more your tests resemble the way your software is used, the more
            confidence they give you.{'"'}
          </Definition>

          <CodeBlock filename="__tests__/login-form.test.tsx">
{`import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect } from 'vitest'

describe('LoginForm', () => {
  it('shows email and password fields', () => {
    // render(LoginForm component)
    // screen.getByPlaceholderText('Email') should exist
    // screen.getByPlaceholderText('Password') should exist
  })

  it('shows validation error for empty email', () => {
    // render(LoginForm component)
    // Click submit button without filling fields
    // screen.getByText('Email is required') should appear
  })

  it('disables submit button while loading', () => {
    // render(LoginForm with loading=true)
    // Submit button should be disabled
    // Button text should be 'Signing in...'
  })
})`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit scaffolds test files in <Code>apps/web/__tests__/</Code> and <Code>apps/admin/__tests__/</Code>.
            Run them with:
          </p>

          <CodeBlock filename="Terminal">
{`# Run all frontend tests
pnpm test

# Watch mode (re-runs on file changes)
pnpm test:watch

# Run tests for a specific file
pnpm test login-form`}
          </CodeBlock>

          <Challenge number={8} title="Find and Run Frontend Tests">
            <p>Navigate to <Code>apps/web/</Code> and <Code>apps/admin/</Code>. Find the test files
            in the <Code>__tests__</Code> directories. Run <Code>pnpm test</Code> in each. How many
            component tests does Grit scaffold? What components are tested?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: E2E Testing with Playwright ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">E2E Testing with Playwright</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Unit tests check individual pieces. Integration tests check pieces together. But neither tells
            you if the whole application works from a user{"'"}s perspective. That{"'"}s what Playwright does — it
            opens a real browser, clicks buttons, fills forms, and verifies results.
          </p>

          <Definition term="Playwright">
            A browser automation framework for end-to-end testing. It controls Chromium, Firefox, and
            WebKit browsers programmatically. Tests run headlessly (no visible browser window) in CI,
            or headed (visible) for debugging. Built by Microsoft, it{"'"}s the modern replacement for
            Selenium and Cypress.
          </Definition>

          <CodeBlock filename="e2e/auth.spec.ts">
{`import { test, expect } from '@playwright/test'

test('user can register and see dashboard', async ({ page }) => {
  // 1. Navigate to /register
  // await page.goto('/register')

  // 2. Fill in name, email, password
  // await page.fill('[name="name"]', 'John Doe')
  // await page.fill('[name="email"]', 'john@test.com')
  // await page.fill('[name="password"]', 'securePass123')

  // 3. Click Register button
  // await page.click('button[type="submit"]')

  // 4. Assert: redirected to /dashboard
  // await expect(page).toHaveURL('/dashboard')

  // 5. Assert: welcome message visible
  // await expect(page.locator('h1')).toContainText('Dashboard')
})

test('login with invalid credentials shows error', async ({ page }) => {
  // Navigate to /login
  // Fill wrong credentials
  // Click submit
  // Assert: error message appears
  // Assert: still on /login page
})`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit scaffolds E2E tests in the <Code>e2e/</Code> directory with <Code>auth.spec.ts</Code> and
            <Code>admin.spec.ts</Code>. Run them with:
          </p>

          <CodeBlock filename="Terminal">
{`# Run E2E tests (app must be running)
pnpm test:e2e

# Run with visible browser (for debugging)
pnpm test:e2e -- --headed

# Run a specific test file
pnpm test:e2e auth.spec.ts`}
          </CodeBlock>

          <Tip>
            E2E tests require the application to be running. Start your API and frontend first
            with <Code>grit dev</Code>, then run <Code>pnpm test:e2e</Code> in a separate terminal.
          </Tip>

          <Challenge number={9} title="Find the Playwright Config">
            <p>Find the Playwright configuration file in your project (<Code>playwright.config.ts</Code>).
            What base URL is it configured to use? What browsers does it test against? Start your app
            with <Code>grit dev</Code>, then run <Code>pnpm test:e2e</Code>. Do the tests pass?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Test-Driven Development (TDD) ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Test-Driven Development (TDD)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            So far, we{"'"}ve written tests for existing code. TDD flips this around: you write the test
            first, then write the code to make it pass. It sounds backwards, but it{"'"}s one of the most
            effective ways to write reliable code.
          </p>

          <Definition term="Test-Driven Development (TDD)">
            A development practice where you write a failing test before writing any implementation code.
            The cycle has three steps: <strong className="text-foreground">Red</strong> (write a test that
            fails), <strong className="text-foreground">Green</strong> (write the minimum code to make it
            pass), <strong className="text-foreground">Refactor</strong> (clean up the code while keeping
            tests green). Repeat for each feature.
          </Definition>

          <CodeBlock filename="TDD Example: CalculateDiscount">
{`// STEP 1 — RED: Write the test first
func TestCalculateDiscount(t *testing.T) {
    // 10% discount on $100 = $90
    result := CalculateDiscount(100.0, 10)
    assert.Equal(t, 90.0, result)

    // 25% discount on $200 = $150
    result = CalculateDiscount(200.0, 25)
    assert.Equal(t, 150.0, result)

    // 0% discount = original price
    result = CalculateDiscount(50.0, 0)
    assert.Equal(t, 50.0, result)

    // 100% discount = free
    result = CalculateDiscount(50.0, 100)
    assert.Equal(t, 0.0, result)
}

// STEP 2 — GREEN: Write the minimum code
func CalculateDiscount(price float64, percent int) float64 {
    return price - (price * float64(percent) / 100)
}

// STEP 3 — REFACTOR: Add edge cases
func TestCalculateDiscountEdgeCases(t *testing.T) {
    // Negative discount should return error or original price
    // Discount over 100% should cap at 100%
    // Negative price should return error
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The power of TDD: by writing the test first, you are forced to think about the interface
            before the implementation. What should the function be called? What parameters does it take?
            What does it return? What edge cases exist? The test answers all these questions before you
            write a single line of implementation code.
          </p>

          <Challenge number={10} title="Practice TDD">
            <p>Use TDD to build a <Code>CalculateShipping</Code> function. Write the tests first:
            orders under $50 cost $5.99 shipping, orders $50-$99 cost $2.99, orders $100+ get free
            shipping, and international orders always cost $15.99. Then implement the function to make
            all tests pass.</p>
          </Challenge>

          <Challenge number={11} title="TDD for String Utilities">
            <p>Use TDD to build a <Code>Slugify</Code> function. Write tests first: {'"'}Hello World{'"'} becomes
            {'"'}hello-world{'"'}, {'"'}My First Blog Post!{'"'} becomes {'"'}my-first-blog-post{'"'}, leading/trailing
            spaces are trimmed, multiple spaces become a single dash. Then implement the function.</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You{"'"}ve learned the entire testing stack for a Grit application:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Go unit tests</strong> with <Code>testing.T</Code> and testify assertions</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">In-memory SQLite</strong> for fast, isolated database tests</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Auth endpoint tests</strong> covering registration, login, and error cases</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">CRUD resource tests</strong> for generated resources (create, list, get, update, delete)</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Benchmarks</strong> to measure endpoint performance with <Code>testing.B</Code></span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Vitest</strong> for React component testing with React Testing Library</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Playwright</strong> for end-to-end browser automation tests</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">TDD</strong> — write the test first, then the code</span></li>
          </ul>

          <Challenge number={12} title="Final Challenge: Go API Tests">
            <p>Generate a Task resource (<Code>grit generate resource Task title:string description:text
            status:string priority:int due_date:time</Code>). Write 5 Go API tests:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Create a task with valid data (201)</li>
              <li>List tasks with pagination (200, correct count)</li>
              <li>Get a task by ID (200, correct data)</li>
              <li>Update a task{"'"}s status (200, updated field)</li>
              <li>Delete a task (200, subsequent GET returns 404)</li>
            </ol>
          </Challenge>

          <Challenge number={13} title="Final Challenge: Vitest Component Tests">
            <p>Write 2 Vitest component tests for the Task resource:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Test that the TaskForm renders all fields (title, description, status, priority, due date)</li>
              <li>Test that the TaskList renders task items with correct titles and status badges</li>
            </ol>
          </Challenge>

          <Challenge number={14} title="Final Challenge: Playwright E2E Test">
            <p>Write 1 Playwright end-to-end test for the complete task flow:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Log in as a test user</li>
              <li>Navigate to the tasks page</li>
              <li>Create a new task (fill the form, submit)</li>
              <li>Verify the task appears in the list</li>
              <li>Edit the task{"'"}s status to {'"'}completed{'"'}</li>
              <li>Delete the task and verify it{"'"}s gone</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
