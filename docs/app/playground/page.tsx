'use client'

import { useState, useCallback, useEffect, Suspense } from 'react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { Play, Loader2, RotateCcw, Share2, Check, BookOpen } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import CodeMirror from '@uiw/react-codemirror'
import { go } from '@codemirror/lang-go'
import { vscodeDark } from '@uiw/codemirror-theme-vscode'

const DEFAULT_CODE = `package main

import "fmt"

func main() {
\tfmt.Println("Hello from Grit Playground!")

\t// Try editing this code and press Run (or Ctrl+Enter)
\tnames := []string{"Go", "React", "Grit"}
\tfor i, name := range names {
\t\tfmt.Printf("%d: %s\\n", i, name)
\t}
}
`

const EXAMPLES = [
  {
    label: 'Hello World',
    code: `package main

import "fmt"

func main() {
\tfmt.Println("Hello, World!")
}
`,
  },
  {
    label: 'Structs & Methods',
    code: `package main

import "fmt"

type User struct {
\tFirstName string
\tLastName  string
\tEmail     string
\tRole      string
}

func (u User) FullName() string {
\treturn u.FirstName + " " + u.LastName
}

func (u *User) Promote() {
\tu.Role = "ADMIN"
}

func main() {
\tuser := User{
\t\tFirstName: "Alice",
\t\tLastName:  "Johnson",
\t\tEmail:     "alice@example.com",
\t\tRole:      "USER",
\t}

\tfmt.Println("Name:", user.FullName())
\tfmt.Println("Role:", user.Role)

\tuser.Promote()
\tfmt.Println("After Promote:", user.Role)
}
`,
  },
  {
    label: 'Error Handling',
    code: `package main

import (
\t"errors"
\t"fmt"
)

func divide(a, b float64) (float64, error) {
\tif b == 0 {
\t\treturn 0, errors.New("cannot divide by zero")
\t}
\treturn a / b, nil
}

func main() {
\tresult, err := divide(10, 3)
\tif err != nil {
\t\tfmt.Println("Error:", err)
\t\treturn
\t}
\tfmt.Printf("10 / 3 = %.2f\\n", result)

\t_, err = divide(10, 0)
\tif err != nil {
\t\tfmt.Println("Error:", err)
\t}
}
`,
  },
  {
    label: 'Interfaces',
    code: `package main

import "fmt"

type Shape interface {
\tArea() float64
\tName() string
}

type Circle struct {
\tRadius float64
}

func (c Circle) Area() float64 {
\treturn 3.14159 * c.Radius * c.Radius
}

func (c Circle) Name() string {
\treturn "Circle"
}

type Rectangle struct {
\tWidth, Height float64
}

func (r Rectangle) Area() float64 {
\treturn r.Width * r.Height
}

func (r Rectangle) Name() string {
\treturn "Rectangle"
}

func printShape(s Shape) {
\tfmt.Printf("%s: area = %.2f\\n", s.Name(), s.Area())
}

func main() {
\tshapes := []Shape{
\t\tCircle{Radius: 5},
\t\tRectangle{Width: 4, Height: 6},
\t\tCircle{Radius: 3},
\t}

\tfor _, s := range shapes {
\t\tprintShape(s)
\t}
}
`,
  },
  {
    label: 'Goroutines',
    code: `package main

import (
\t"fmt"
\t"sync"
\t"time"
)

func worker(id int, wg *sync.WaitGroup) {
\tdefer wg.Done()
\tfmt.Printf("Worker %d starting\\n", id)
\ttime.Sleep(time.Millisecond * 100)
\tfmt.Printf("Worker %d done\\n", id)
}

func main() {
\tvar wg sync.WaitGroup

\tfor i := 1; i <= 5; i++ {
\t\twg.Add(1)
\t\tgo worker(i, &wg)
\t}

\twg.Wait()
\tfmt.Println("All workers finished!")
}
`,
  },
  {
    label: 'Maps & Slices',
    code: `package main

import "fmt"

func main() {
\t// Slices
\tfruits := []string{"apple", "banana", "cherry"}
\tfruits = append(fruits, "date", "elderberry")

\tfmt.Println("Fruits:", fruits)
\tfmt.Println("First 3:", fruits[:3])
\tfmt.Println("Length:", len(fruits))

\t// Maps
\tscores := map[string]int{
\t\t"Alice": 95,
\t\t"Bob":   87,
\t\t"Carol": 92,
\t}

\tscores["Diana"] = 98

\tfor name, score := range scores {
\t\tfmt.Printf("%s scored %d\\n", name, score)
\t}

\t// Check if key exists
\tif val, ok := scores["Eve"]; ok {
\t\tfmt.Println("Eve:", val)
\t} else {
\t\tfmt.Println("Eve not found")
\t}
}
`,
  },
]

interface CompileEvent {
  Message: string
  Kind: string
  Delay: number
}

interface CompileResponse {
  Errors: string
  Events: CompileEvent[] | null
}

export default function PlaygroundPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background" />}>
      <PlaygroundInner />
    </Suspense>
  )
}

function PlaygroundInner() {
  const searchParams = useSearchParams()
  const [code, setCode] = useState(DEFAULT_CODE)
  const [output, setOutput] = useState('')
  const [running, setRunning] = useState(false)
  const [copied, setCopied] = useState(false)
  const [duration, setDuration] = useState<number | null>(null)
  const [challengeTitle, setChallengeTitle] = useState<string | null>(null)

  // Load code from URL query params (?code=base64&title=name)
  useEffect(() => {
    const codeParam = searchParams.get('code')
    const titleParam = searchParams.get('title')
    if (codeParam) {
      try {
        setCode(atob(codeParam))
      } catch {
        // Invalid base64, ignore
      }
    }
    if (titleParam) {
      setChallengeTitle(titleParam)
    }
  }, [searchParams])

  const runCode = useCallback(async () => {
    setRunning(true)
    setOutput('')
    setDuration(null)

    const start = performance.now()

    try {
      const res = await fetch('/api/compile', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        body: `version=2&body=${encodeURIComponent(code)}&withVet=true`,
      })

      const data: CompileResponse = await res.json()
      const elapsed = Math.round(performance.now() - start)
      setDuration(elapsed)

      if (data.Errors) {
        setOutput(data.Errors)
      } else if (data.Events && data.Events.length > 0) {
        setOutput(data.Events.map((e) => e.Message).join(''))
      } else {
        setOutput('(no output)')
      }
    } catch (err) {
      setOutput(`Error: Could not reach the Go Playground service.\n${err}`)
    } finally {
      setRunning(false)
    }
  }, [code])

  const resetCode = useCallback(() => {
    setCode(DEFAULT_CODE)
    setOutput('')
    setDuration(null)
  }, [])

  const shareCode = useCallback(async () => {
    try {
      const res = await fetch('/api/share', {
        method: 'POST',
        body: code,
      })
      const data = await res.json()
      await navigator.clipboard.writeText(data.url)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      // Fallback: copy code to clipboard
      await navigator.clipboard.writeText(code)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }, [code])

  // Keyboard shortcut: Ctrl+Enter / Cmd+Enter to run
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'Enter') {
        e.preventDefault()
        runCode()
      }
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [runCode])

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <SiteHeader />

      {/* Toolbar */}
      <div className="border-b border-border/40 bg-background/80 backdrop-blur-sm">
        <div className="container max-w-screen-2xl flex items-center justify-between px-4 py-2">
          <div className="flex items-center gap-3">
            <h1 className="text-sm font-semibold text-foreground/90">Go Playground</h1>
            {challengeTitle ? (
              <span className="hidden sm:inline-flex items-center rounded-full bg-emerald-500/10 border border-emerald-500/20 px-2.5 py-0.5 text-[11px] font-medium text-emerald-400">
                Challenge: {challengeTitle}
              </span>
            ) : (
              <span className="text-xs text-muted-foreground/50 hidden sm:inline">
                Powered by go.dev
              </span>
            )}
          </div>

          <div className="flex items-center gap-2">
            {/* Examples dropdown */}
            <select
              className="h-8 rounded-md border border-border/40 bg-background px-2 text-xs text-muted-foreground focus:outline-none focus:ring-1 focus:ring-primary/50"
              onChange={(e) => {
                const example = EXAMPLES.find((ex) => ex.label === e.target.value)
                if (example) {
                  setCode(example.code)
                  setOutput('')
                  setDuration(null)
                }
              }}
              defaultValue=""
            >
              <option value="" disabled>
                Examples
              </option>
              {EXAMPLES.map((ex) => (
                <option key={ex.label} value={ex.label}>
                  {ex.label}
                </option>
              ))}
            </select>

            <Button
              variant="ghost"
              size="sm"
              className="h-8 px-2 text-xs text-muted-foreground"
              onClick={resetCode}
            >
              <RotateCcw className="h-3.5 w-3.5 mr-1" />
              Reset
            </Button>

            <Button
              variant="ghost"
              size="sm"
              className="h-8 px-2 text-xs text-muted-foreground"
              onClick={shareCode}
            >
              {copied ? (
                <Check className="h-3.5 w-3.5 mr-1 text-green-500" />
              ) : (
                <Share2 className="h-3.5 w-3.5 mr-1" />
              )}
              {copied ? 'Copied!' : 'Share'}
            </Button>

            <Button
              variant="ghost"
              size="sm"
              className="h-8 px-2 text-xs text-muted-foreground"
              asChild
            >
              <Link href="/docs/prerequisites/golang">
                <BookOpen className="h-3.5 w-3.5 mr-1" />
                Go Guide
              </Link>
            </Button>

            <Button
              size="sm"
              className="h-8 px-3 text-xs font-medium bg-green-600 hover:bg-green-700 text-white"
              onClick={runCode}
              disabled={running}
            >
              {running ? (
                <Loader2 className="h-3.5 w-3.5 mr-1 animate-spin" />
              ) : (
                <Play className="h-3.5 w-3.5 mr-1" />
              )}
              {running ? 'Running...' : 'Run'}
              <span className="hidden sm:inline ml-1 text-white/60">
                {navigator.platform?.includes('Mac') ? '⌘↵' : 'Ctrl+↵'}
              </span>
            </Button>
          </div>
        </div>
      </div>

      {/* Editor + Output */}
      <div className="flex-1 flex flex-col lg:flex-row">
        {/* Editor panel */}
        <div className="flex-1 flex flex-col min-h-0 border-b lg:border-b-0 lg:border-r border-border/40">
          <div className="flex items-center justify-between px-4 py-1.5 border-b border-border/20 bg-accent/10">
            <span className="text-xs text-muted-foreground/60 font-mono">main.go</span>
          </div>
          <div className="flex-1 min-h-[300px] lg:min-h-0 overflow-auto">
            <CodeMirror
              value={code}
              onChange={setCode}
              extensions={[go()]}
              theme={vscodeDark}
              basicSetup={{
                lineNumbers: true,
                highlightActiveLineGutter: true,
                highlightActiveLine: true,
                bracketMatching: true,
                closeBrackets: true,
                autocompletion: true,
                foldGutter: true,
                indentOnInput: true,
              }}
              style={{ height: '100%', fontSize: '13px' }}
            />
          </div>
        </div>

        {/* Output panel */}
        <div className="flex-1 lg:max-w-[50%] flex flex-col min-h-0">
          <div className="flex items-center justify-between px-4 py-1.5 border-b border-border/20 bg-accent/10">
            <span className="text-xs text-muted-foreground/60 font-mono">Output</span>
            {duration !== null && (
              <span className="text-xs text-muted-foreground/40 font-mono">{duration}ms</span>
            )}
          </div>
          <div className="flex-1 min-h-[200px] lg:min-h-0 overflow-auto p-4 bg-[#1e1e1e]">
            {running ? (
              <div className="flex items-center gap-2 text-muted-foreground/60 text-sm">
                <Loader2 className="h-4 w-4 animate-spin" />
                Compiling and running...
              </div>
            ) : output ? (
              <pre className="text-sm font-mono text-foreground/90 whitespace-pre-wrap leading-relaxed">
                {output}
              </pre>
            ) : (
              <p className="text-sm text-muted-foreground/40">
                Press <strong>Run</strong> or <kbd className="px-1.5 py-0.5 rounded bg-accent/20 border border-border/30 text-xs font-mono">Ctrl+Enter</kbd> to execute your code.
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
