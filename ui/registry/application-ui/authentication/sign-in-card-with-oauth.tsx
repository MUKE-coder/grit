'use client'

import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'

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
 * A sign-in card: OAuth above, email and password below.
 *
 * This is one of the few blocks in the library that depends on the shadcn
 * primitives rather than being self-contained markup, and the reason is
 * validation. FormField hands the input an id, ties the label to it, points
 * aria-describedby at the error and flips aria-invalid when it fires. Rolling
 * that by hand for every field is how you end up with a form that looks fine
 * and reports nothing to a screen reader. `registryDependencies` lists what it
 * needs, so installing this pulls button, input and form in with it.
 *
 * The provider logos ARE the real marks, unlike the placeholder tiles in the
 * Integrations blocks. That is not an inconsistency: Google's sign-in branding
 * guidelines require their mark on the button, and a "Continue with Google"
 * button carrying a grey square is a button people do not trust enough to
 * press. Marketing pages showing off logos are a different question.
 *
 * The autocomplete attributes are load-bearing. `email` and `current-password`
 * are what let a password manager fill this form; without them it either does
 * nothing or fills the wrong field, and users blame your site rather than their
 * manager. `current-password` also distinguishes this from a registration form,
 * where the value is `new-password` and managers offer to generate one instead.
 *
 * `onSubmit` is a prop that takes the values. The block ships with a resolved
 * promise so the preview works, but nothing here talks to a server: wire it to
 * your own action and handle the failure case, which is the part demos skip.
 */

const signInSchema = z.object({
  email: z.string().email({ message: 'Enter a valid email address' }),
  password: z.string().min(8, { message: 'Password must be at least 8 characters' }),
})

export type SignInValues = z.infer<typeof signInSchema>

function GoogleMark() {
  return (
    <svg viewBox="0 0 256 262" aria-hidden="true" className="size-4">
      <path
        fill="#4285f4"
        d="M255.878 133.451c0-10.734-.871-18.567-2.756-26.69H130.55v48.448h71.947c-1.45 12.04-9.283 30.172-26.69 42.356l-.244 1.622l38.755 30.023l2.685.268c24.659-22.774 38.875-56.282 38.875-96.027"
      />
      <path
        fill="#34a853"
        d="M130.55 261.1c35.248 0 64.839-11.605 86.453-31.622l-41.196-31.913c-11.024 7.688-25.82 13.055-45.257 13.055c-34.523 0-63.824-22.773-74.269-54.25l-1.531.13l-40.298 31.187l-.527 1.465C35.393 231.798 79.49 261.1 130.55 261.1"
      />
      <path
        fill="#fbbc05"
        d="M56.281 156.37c-2.756-8.123-4.351-16.827-4.351-25.82c0-8.994 1.595-17.697 4.206-25.82l-.073-1.73L15.26 71.312l-1.335.635C5.077 89.644 0 109.517 0 130.55s5.077 40.905 13.925 58.602z"
      />
      <path
        fill="#eb4335"
        d="M130.55 50.479c24.514 0 41.05 10.589 50.479 19.438l36.844-35.974C195.245 12.91 165.798 0 130.55 0C79.49 0 35.393 29.301 13.925 71.947l42.211 32.783c10.59-31.477 39.891-54.251 74.414-54.251"
      />
    </svg>
  )
}

function GitHubMark() {
  return (
    <svg viewBox="0 0 256 256" aria-hidden="true" className="size-4 fill-gray-900 dark:fill-white">
      <path d="M128 0C57.307 0 0 57.307 0 128c0 56.562 36.665 104.535 87.535 121.469 6.405 1.19 8.753-2.781 8.753-6.17 0-3.063-.122-11.175-.174-21.935-35.606 7.738-43.122-17.178-43.122-17.178-5.83-14.814-14.222-18.757-14.222-18.757-11.64-7.96.877-7.803.877-7.803 12.88.907 19.67 13.237 19.67 13.237 11.437 19.59 29.992 13.928 37.292 10.646 1.14-8.292 4.48-13.933 8.147-17.14-28.426-3.237-58.337-14.213-58.337-63.287 0-13.977 5-25.412 13.237-34.374-1.32-3.24-5.73-16.274 1.256-33.932 0 0 10.8-3.457 35.4 13.146 10.263-2.847 21.285-4.27 32.243-4.32 10.957.05 21.982 1.474 32.25 4.32 24.59-16.603 35.38-13.146 35.38-13.146 7 17.658 2.58 30.692 1.26 33.932 8.25 8.96 13.22 20.397 13.22 34.374 0 49.2-29.96 60-58.44 63.167 4.6 3.97 8.72 11.767 8.72 23.77 0 17.17-.155 30.983-.155 35.195 0 3.407 2.328 7.4 8.8 6.158C219.34 232.508 256 184.546 256 128 256 57.307 198.693 0 128 0z" />
    </svg>
  )
}

export default function SignInCardWithOauth({
  title = 'Sign in',
  subtitle = 'Welcome back. Sign in to continue.',
  forgotHref = '#',
  registerHref = '#',
  onSubmit,
}: {
  title?: string
  subtitle?: string
  forgotHref?: string
  registerHref?: string
  onSubmit?: (values: SignInValues) => Promise<void> | void
}) {
  const [submitting, setSubmitting] = useState(false)
  const [failure, setFailure] = useState<string | null>(null)

  const form = useForm<SignInValues>({
    resolver: zodResolver(signInSchema),
    defaultValues: { email: '', password: '' },
  })

  async function handleSubmit(values: SignInValues) {
    setSubmitting(true)
    setFailure(null)
    try {
      await onSubmit?.(values)
    } catch {
      setFailure('We could not sign you in. Check your details and try again.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-16 dark:bg-gray-950">
      <div className="w-full max-w-md">
        <div className="rounded-2xl border border-gray-200 bg-white p-8 shadow-sm dark:border-white/10 dark:bg-gray-900">
          <h1 className="text-xl font-semibold text-gray-900 dark:text-white">{title}</h1>
          <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">{subtitle}</p>

          <div className="mt-6 grid grid-cols-2 gap-3">
            <Button type="button" variant="outline">
              <GoogleMark />
              <span>Google</span>
            </Button>
            <Button type="button" variant="outline">
              <GitHubMark />
              <span>GitHub</span>
            </Button>
          </div>

          <div className="my-6 flex items-center gap-3">
            <span aria-hidden="true" className="h-px flex-1 bg-gray-200 dark:bg-white/10" />
            <span className="text-xs text-gray-500 dark:text-gray-400">or</span>
            <span aria-hidden="true" className="h-px flex-1 bg-gray-200 dark:bg-white/10" />
          </div>

          {/* A failed sign-in has to be announced, not just displayed. Without
              role="alert" someone using a screen reader submits the form, hears
              nothing, and has no idea why the page did not move. */}
          {failure && (
            <p
              role="alert"
              className="mb-4 rounded-md bg-red-50 px-3 py-2 text-sm text-red-700 dark:bg-red-500/10 dark:text-red-400"
            >
              {failure}
            </p>
          )}

          <Form {...form}>
            <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-5">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input
                        type="email"
                        autoComplete="email"
                        placeholder="you@example.com"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <div className="flex items-center justify-between">
                      <FormLabel>Password</FormLabel>
                      <a
                        href={forgotHref}
                        className="text-sm font-medium text-indigo-600 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
                      >
                        Forgot password?
                      </a>
                    </div>
                    <FormControl>
                      {/* current-password, not new-password: this tells a
                          password manager to fill rather than to offer to
                          generate one. */}
                      <Input type="password" autoComplete="current-password" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" className="w-full" disabled={submitting}>
                {submitting ? 'Signing in...' : 'Sign in'}
              </Button>
            </form>
          </Form>
        </div>

        <p className="mt-6 text-center text-sm text-gray-600 dark:text-gray-400">
          Don&rsquo;t have an account?{' '}
          <a
            href={registerHref}
            className="font-medium text-indigo-600 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
          >
            Create one
          </a>
        </p>
      </div>
    </section>
  )
}
