'use client'

import { useEffect, useRef, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { LoaderCircle, MailCheck } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'

/*
 * Request a reset link, then the state that comes after it.
 *
 * The form is replaced rather than annotated. The source left the form in
 * place and put a green line above it, so the page looked the same after
 * submitting as before, and nothing stopped you sending three more. Swapping
 * to a confirmation panel says what happened and names the address it went to,
 * which is the one thing people actually need to check.
 *
 * Focus moves to the confirmation heading. When a view replaces itself, a live
 * region announces the new text but leaves focus back where the form used to
 * be, so the next Tab goes somewhere unrelated. Moving focus to the heading
 * puts the reader at the top of what just appeared. The heading takes
 * tabIndex={-1} for that and nothing else.
 *
 * Resend is aria-disabled during its cooldown, not disabled. A disabled button
 * loses focus the instant it is pressed, dropping a keyboard user back to the
 * top of the document; aria-disabled keeps the button focusable and the click
 * handler ignores the press.
 *
 * The countdown is not a live region. A number changing every second would be
 * announced every second. The status region speaks twice: when the mail is
 * sent, and once when resending becomes available again.
 *
 * autocomplete is "email". The source set autocomplete="off" on an email
 * field, which stops a password manager filling the one thing it definitely
 * knows. It also carried autoFocus, which drops a screen reader user into the
 * input past the heading that explains the page.
 */

const schema = z.object({
  email: z.string().email({ message: 'Enter a valid email address' }),
})

export type ForgotPasswordValues = z.infer<typeof schema>

const COOLDOWN_SECONDS = 30

export default function ForgotPasswordWithSentState({
  onSubmit = async () => {},
  signInHref = '#',
}: {
  onSubmit?: (values: ForgotPasswordValues) => Promise<void> | void
  signInHref?: string
}) {
  const [sentTo, setSentTo] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [cooldown, setCooldown] = useState(0)
  const [announcement, setAnnouncement] = useState('')
  const heading = useRef<HTMLHeadingElement>(null)

  const form = useForm<ForgotPasswordValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: '' },
    mode: 'onTouched',
  })

  /* The dep is the boolean, not the number, so the interval is created once
     when the cooldown starts and cleared when it ends rather than being torn
     down and rebuilt every second. */
  const counting = cooldown > 0
  useEffect(() => {
    if (!counting) return
    const timer = setInterval(() => setCooldown((current) => Math.max(0, current - 1)), 1000)
    return () => clearInterval(timer)
  }, [counting])

  /* Announced from an effect rather than from inside the setCooldown updater.
     Calling setState from within another updater is a render side effect, and
     React runs updaters twice in development to surface exactly that. */
  const wasCounting = useRef(false)
  useEffect(() => {
    if (wasCounting.current && cooldown === 0) {
      setAnnouncement('You can send the email again now.')
    }
    wasCounting.current = counting
  }, [counting, cooldown])

  /* Runs after the confirmation renders, so the heading exists to receive
     focus. */
  useEffect(() => {
    if (sentTo) heading.current?.focus()
  }, [sentTo])

  async function send(values: ForgotPasswordValues) {
    setSubmitting(true)
    try {
      await onSubmit(values)
      setSentTo(values.email)
      setCooldown(COOLDOWN_SECONDS)
      setAnnouncement(`Reset link sent to ${values.email}.`)
    } finally {
      setSubmitting(false)
    }
  }

  function resend() {
    if (cooldown > 0 || !sentTo) return
    void send({ email: sentTo })
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 dark:bg-gray-950">
      <div className="w-full max-w-md rounded-2xl border border-gray-200 bg-white p-8 shadow-sm dark:border-white/10 dark:bg-gray-900">
        <p role="status" aria-live="polite" className="sr-only">
          {announcement}
        </p>

        {sentTo ? (
          <div className="text-center">
            <span
              aria-hidden="true"
              className="mx-auto flex size-12 items-center justify-center rounded-full bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400"
            >
              <MailCheck className="size-6" />
            </span>

            <h1
              ref={heading}
              tabIndex={-1}
              className="mt-4 text-2xl font-bold tracking-tight text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-indigo-600 dark:text-white"
            >
              Check your email
            </h1>

            <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
              We sent a reset link to{' '}
              <span className="font-medium text-gray-900 dark:text-white">{sentTo}</span>. It
              expires in an hour.
            </p>

            <div className="mt-8 space-y-3">
              <Button
                type="button"
                variant="outline"
                onClick={resend}
                /* aria-disabled, not disabled. See the note above. */
                aria-disabled={cooldown > 0 || submitting}
                className={`h-11 w-full ${cooldown > 0 || submitting ? 'opacity-50' : ''}`}
              >
                {submitting && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
                {cooldown > 0 ? `Resend in ${cooldown}s` : 'Resend email'}
              </Button>

              <Button
                type="button"
                variant="ghost"
                onClick={() => {
                  setSentTo(null)
                  setCooldown(0)
                  setAnnouncement('')
                }}
                className="h-11 w-full"
              >
                Use a different address
              </Button>
            </div>

            <p className="mt-6 text-sm text-gray-600 dark:text-gray-400">
              Or return to{' '}
              <a
                href={signInHref}
                className="font-medium text-indigo-600 underline-offset-4 hover:underline dark:text-indigo-400"
              >
                sign in
              </a>
              .
            </p>
          </div>
        ) : (
          <>
            <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
              Forgot your password?
            </h1>
            <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
              Give us the address on the account and we will send a link to reset it.
            </p>

            <Form {...form}>
              <form onSubmit={form.handleSubmit(send)} className="mt-8 space-y-6">
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Email address</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          type="email"
                          autoComplete="email"
                          placeholder="you@example.com"
                          className="h-11"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <Button type="submit" disabled={submitting} className="h-11 w-full font-medium">
                  {submitting && <LoaderCircle aria-hidden="true" className="size-4 animate-spin" />}
                  Email a reset link
                </Button>
              </form>
            </Form>

            <p className="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
              Or return to{' '}
              <a
                href={signInHref}
                className="font-medium text-indigo-600 underline-offset-4 hover:underline dark:text-indigo-400"
              >
                sign in
              </a>
              .
            </p>
          </>
        )}
      </div>
    </div>
  )
}
