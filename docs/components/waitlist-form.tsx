'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Loader2, CheckCircle2, Mail } from 'lucide-react'

const GOOGLE_SCRIPT_URL = 'https://script.google.com/macros/s/AKfycbwyfriICvROwcecMLoK61XIadB13FsYxf_AglVDX4jP79XKzpgj0sm2YIRJHgOmUeoS/exec'

export function WaitlistForm() {
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [status, setStatus] = useState<'idle' | 'loading' | 'success' | 'error'>('idle')
  const [errorMsg, setErrorMsg] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!name.trim() || !email.trim()) return

    setStatus('loading')
    setErrorMsg('')

    try {
      const formData = new URLSearchParams()
      formData.append('name', name.trim())
      formData.append('email', email.trim())

      await fetch(GOOGLE_SCRIPT_URL, {
        method: 'POST',
        mode: 'no-cors',
        body: formData,
      })

      setStatus('success')
      setName('')
      setEmail('')
    } catch {
      setStatus('error')
      setErrorMsg('Something went wrong. Please try again.')
    }
  }

  if (status === 'success') {
    return (
      <div className="rounded-xl border border-primary/20 bg-primary/5 p-8 text-center">
        <CheckCircle2 className="h-8 w-8 text-primary mx-auto mb-3" />
        <h3 className="text-lg font-semibold mb-1">You&apos;re on the list!</h3>
        <p className="text-sm text-muted-foreground">
          We&apos;ll notify you when the course launches. Check your inbox.
        </p>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div className="flex flex-col sm:flex-row gap-3">
        <input
          type="text"
          placeholder="Your name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
          className="flex-1 h-11 rounded-lg border border-border/60 bg-background px-4 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary/40 transition-colors"
        />
        <input
          type="email"
          placeholder="you@example.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          className="flex-1 h-11 rounded-lg border border-border/60 bg-background px-4 text-sm text-foreground placeholder:text-muted-foreground/50 focus:outline-none focus:ring-2 focus:ring-primary/30 focus:border-primary/40 transition-colors"
        />
        <Button
          type="submit"
          disabled={status === 'loading'}
          className="h-11 px-6 glow-purple-sm"
        >
          {status === 'loading' ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <>
              <Mail className="h-4 w-4 mr-2" />
              Join Waitlist
            </>
          )}
        </Button>
      </div>
      {status === 'error' && (
        <p className="text-xs text-red-400">{errorMsg}</p>
      )}
      <p className="text-[11px] text-muted-foreground/50">
        No spam. We&apos;ll only email you when the course is ready.
      </p>
    </form>
  )
}
