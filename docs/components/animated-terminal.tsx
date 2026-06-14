'use client'

import { useState, useEffect, useCallback } from 'react'

interface TerminalLine {
  type: 'command' | 'output' | 'success' | 'step'
  text: string
  delay?: number
}

const lines: TerminalLine[] = [
  { type: 'command', text: 'curl -fsSL https://gritframework.dev/install.sh | sh' },
  { type: 'step', text: 'Downloading grit CLI...', delay: 400 },
  { type: 'success', text: 'Installed grit v3.26.0', delay: 300 },
  { type: 'command', text: 'grit new myapp', delay: 600 },
  { type: 'step', text: 'Creating directory structure...', delay: 250 },
  { type: 'step', text: 'Scaffolding Go API (Gin + GORM)...', delay: 250 },
  { type: 'step', text: 'Adding batteries (cache, storage, mail, jobs, AI)...', delay: 250 },
  { type: 'step', text: 'Setting up Next.js web app...', delay: 250 },
  { type: 'step', text: 'Creating admin panel with data tables and form builders...', delay: 250 },
  { type: 'step', text: 'Writing shared schemas and types...', delay: 250 },
  { type: 'success', text: 'Project "myapp" created successfully!', delay: 300 },
  { type: 'command', text: 'cd myapp && grit start server', delay: 600 },
  { type: 'step', text: 'Starting Go API on :8080...', delay: 300 },
  { type: 'success', text: 'API running at http://localhost:8080', delay: 300 },
  { type: 'command', text: 'grit start client', delay: 600 },
  { type: 'step', text: 'Starting Next.js web on :3000...', delay: 200 },
  { type: 'step', text: 'Starting admin panel on :3001...', delay: 200 },
  { type: 'success', text: 'All services running. Happy building!', delay: 300 },
]

const TYPING_SPEED = 30
const PAUSE_AFTER_COMMAND = 200

export function AnimatedTerminal() {
  const [visibleLines, setVisibleLines] = useState<{ type: string; text: string; typed: string }[]>([])
  const [currentLineIndex, setCurrentLineIndex] = useState(0)
  const [currentCharIndex, setCurrentCharIndex] = useState(0)
  const [isTyping, setIsTyping] = useState(true)
  const [isPaused, setIsPaused] = useState(false)
  const [hasStarted, setHasStarted] = useState(false)

  // Start animation after mount
  useEffect(() => {
    const timer = setTimeout(() => setHasStarted(true), 500)
    return () => clearTimeout(timer)
  }, [])

  const advanceToNextLine = useCallback(() => {
    const nextIndex = currentLineIndex + 1
    if (nextIndex >= lines.length) {
      // Loop: reset after a pause
      setTimeout(() => {
        setVisibleLines([])
        setCurrentLineIndex(0)
        setCurrentCharIndex(0)
        setIsTyping(true)
        setIsPaused(false)
      }, 3000)
      return
    }

    const nextLine = lines[nextIndex]
    const delay = nextLine.delay || 0

    setIsPaused(true)
    setTimeout(() => {
      setCurrentLineIndex(nextIndex)
      setCurrentCharIndex(0)
      setIsTyping(nextLine.type === 'command')
      setIsPaused(false)

      if (nextLine.type !== 'command') {
        // Output lines appear instantly
        setVisibleLines(prev => [...prev, { type: nextLine.type, text: nextLine.text, typed: nextLine.text }])
      }
    }, delay)
  }, [currentLineIndex])

  useEffect(() => {
    if (!hasStarted || isPaused) return

    const line = lines[currentLineIndex]
    if (!line) return

    if (line.type === 'command' && isTyping) {
      if (currentCharIndex <= line.text.length) {
        const timer = setTimeout(() => {
          const typed = line.text.slice(0, currentCharIndex)

          // Update or add current typing line
          setVisibleLines(prev => {
            const updated = [...prev]
            const existingIndex = updated.findIndex(
              (l, i) => i === prev.length - 1 && l.type === 'command' && prev.length > 0
            )

            if (prev.length > 0 && prev[prev.length - 1].type === 'command' && prev[prev.length - 1].text === line.text) {
              updated[updated.length - 1] = { type: 'command', text: line.text, typed }
            } else if (currentCharIndex === 0) {
              updated.push({ type: 'command', text: line.text, typed: '' })
            }

            return updated
          })

          setCurrentCharIndex(prev => prev + 1)
        }, TYPING_SPEED)
        return () => clearTimeout(timer)
      } else {
        // Done typing this command
        const timer = setTimeout(() => {
          setIsTyping(false)
          advanceToNextLine()
        }, PAUSE_AFTER_COMMAND)
        return () => clearTimeout(timer)
      }
    } else if (line.type !== 'command') {
      // Non-command line was already added, advance
      advanceToNextLine()
    }
  }, [hasStarted, currentLineIndex, currentCharIndex, isTyping, isPaused, advanceToNextLine])

  return (
    <div className="rounded-xl border border-border/60 bg-card/80 backdrop-blur-sm overflow-hidden glow-purple">
      {/* Terminal header */}
      <div className="flex items-center gap-2 px-4 py-3 border-b border-border/40 bg-accent/30">
        <div className="flex items-center gap-1.5">
          <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
          <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
          <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
        </div>
        <span className="ml-2 text-[11px] font-mono text-muted-foreground/50">terminal</span>
      </div>
      {/* Terminal body */}
      <div className="p-6 font-mono text-sm space-y-1.5 min-h-[340px]">
        {visibleLines.map((line, i) => (
          <div key={`${i}-${line.text}`} className="flex items-start gap-2">
            {line.type === 'command' ? (
              <>
                <span className="text-primary/50 select-none shrink-0">$ </span>
                <span className="text-foreground/80">
                  {line.typed}
                  {i === visibleLines.length - 1 && isTyping && line.text === lines[currentLineIndex]?.text && (
                    <span className="inline-block w-[7px] h-[14px] bg-primary/70 ml-[1px] animate-pulse align-middle" />
                  )}
                </span>
              </>
            ) : line.type === 'step' ? (
              <span className="text-muted-foreground/40 text-xs pl-4 flex items-center gap-2">
                <span className="text-primary/60">+</span>
                <span>{line.text}</span>
              </span>
            ) : (
              <span className="text-xs pl-4 flex items-center gap-2 text-primary/80">
                <span>&#10003;</span>
                <span>{line.text}</span>
              </span>
            )}
          </div>
        ))}
        {/* Initial cursor */}
        {visibleLines.length === 0 && hasStarted && (
          <div className="flex items-center gap-2">
            <span className="text-primary/50 select-none">$ </span>
            <span className="inline-block w-[7px] h-[14px] bg-primary/70 animate-pulse" />
          </div>
        )}
      </div>
    </div>
  )
}
