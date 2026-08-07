'use client'

import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { Check, Eye, EyeOff, X } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  useFormField,
} from '@/components/ui/form'

/*
 * A registration card with live password requirements.
 *
 * No positive tabindex. The source numbered its fields tabIndex={1} through
 * tabIndex={6}, which does not order the form: a positive tabindex pulls an
 * element in front of every element in the natural order, across the whole
 * page. Drop this form into a page with a header and the first Tab jumps into
 * the middle of it. Source order is the tab order, and this form is already in
 * the right order.
 *
 * The requirements list is the field's description, not a live region. It
 * updates on every keystroke, and a live region would read the whole list out
 * on every character typed, over the top of the typing. As a description it is
 * read when the field takes focus, which is when the rules matter, and each
 * row still carries its met or not met state for anyone going back to check.
 *
 * It takes its id from useFormField rather than setting aria-describedby on
 * the input. FormControl already points aria-describedby at the description
 * and, once validation fails, at the error message too. Writing the attribute
 * by hand replaces that pair with one id, so the error silently stops being
 * announced with the field. FormDescription itself renders a <p>, which cannot
 * legally contain this list.
 *
 * The show/hide toggle carries aria-pressed and stays out of the tab order of
 * the input itself. Toggling it does not move focus, so a password manager and
 * a keyboard user both keep their place.
 *
 * autocomplete is new-password on both password fields, which is what makes a
 * manager offer to generate and save one rather than trying to fill it. name,
 * email and tel are equally load-bearing on the fields above.
 *
 * The schema checks the same rules the checklist shows. A checklist that
 * disagrees with the validation is worse than no checklist.
 */

const RULES = [
  { id: 'length', label: 'At least 8 characters', test: (v: string) => v.length >= 8 },
  { id: 'lower', label: 'A lowercase letter', test: (v: string) => /[a-z]/.test(v) },
  { id: 'upper', label: 'An uppercase letter', test: (v: string) => /[A-Z]/.test(v) },
  { id: 'number', label: 'A number', test: (v: string) => /\d/.test(v) },
] as const

const registerSchema = z
  .object({
    name: z.string().min(2, { message: 'Enter your full name' }),
    email: z.string().email({ message: 'Enter a valid email address' }),
    phone: z.string().min(7, { message: 'Enter a phone number we can reach you on' }),
    password: z
      .string()
      .min(8, { message: 'Password must be at least 8 characters' })
      .regex(/[a-z]/, { message: 'Password needs a lowercase letter' })
      .regex(/[A-Z]/, { message: 'Password needs an uppercase letter' })
      .regex(/\d/, { message: 'Password needs a number' }),
    confirmPassword: z.string(),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Both passwords must match',
    path: ['confirmPassword'],
  })

export type RegisterValues = z.infer<typeof registerSchema>

/*
 * The relative wrapper lives here rather than inside FormControl. FormControl
 * is a Slot: it puts the generated id and the aria attributes on its direct
 * child. Wrap the input in a div and the id lands on the div, so the label's
 * htmlFor points at a container and the input has no id at all.
 */
function PasswordField({
  field,
  label,
  autoComplete,
}: {
  field: React.ComponentProps<typeof Input>
  label: string
  autoComplete: string
}) {
  const [shown, setShown] = useState(false)

  return (
    <div className="relative">
      <FormControl>
        <Input
          {...field}
          type={shown ? 'text' : 'password'}
          autoComplete={autoComplete}
          className="h-11 pr-11"
        />
      </FormControl>
      <button
        type="button"
        onClick={() => setShown((current) => !current)}
        aria-pressed={shown}
        className="absolute inset-y-0 right-0 inline-flex w-11 items-center justify-center text-gray-500 hover:text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-400 dark:hover:text-white"
      >
        {shown ? (
          <EyeOff aria-hidden="true" className="size-4" />
        ) : (
          <Eye aria-hidden="true" className="size-4" />
        )}
        <span className="sr-only">Show {label.toLowerCase()}</span>
      </button>
    </div>
  )
}

/** Renders as the field's description, so FormControl already links it. */
function PasswordRules({ value }: { value: string }) {
  const { formDescriptionId } = useFormField()

  return (
    <ul id={formDescriptionId} role="list" className="mt-3 grid gap-1.5 sm:grid-cols-2">
      {RULES.map((rule) => {
        const met = rule.test(value ?? '')
        return (
          <li
            key={rule.id}
            className={`flex items-center gap-2 text-sm ${
              met ? 'text-emerald-700 dark:text-emerald-400' : 'text-gray-500 dark:text-gray-400'
            }`}
          >
            {met ? (
              <Check aria-hidden="true" className="size-4 shrink-0" />
            ) : (
              <X aria-hidden="true" className="size-4 shrink-0" />
            )}
            {rule.label}
            {/* The tick is a colour and a glyph. This is the state in words. */}
            <span className="sr-only">{met ? ': met' : ': not met yet'}</span>
          </li>
        )
      })}
    </ul>
  )
}

export default function RegisterCardWithPasswordRules({
  onSubmit = async () => {},
  signInHref = '#',
}: {
  onSubmit?: (values: RegisterValues) => Promise<void> | void
  signInHref?: string
}) {
  const [submitting, setSubmitting] = useState(false)

  const form = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { name: '', email: '', phone: '', password: '', confirmPassword: '' },
    mode: 'onTouched',
  })

  const password = form.watch('password')

  async function handleSubmit(values: RegisterValues) {
    setSubmitting(true)
    try {
      await onSubmit(values)
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 dark:bg-gray-950">
      <div className="w-full max-w-lg rounded-2xl border border-gray-200 bg-white p-8 shadow-sm dark:border-white/10 dark:bg-gray-900">
        <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          Create your account
        </h1>
        <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
          Already have one?{' '}
          <a
            href={signInHref}
            className="font-medium text-indigo-600 underline-offset-4 hover:underline dark:text-indigo-400"
          >
            Sign in
          </a>
          .
        </p>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(handleSubmit)} className="mt-8 space-y-5">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Full name</FormLabel>
                  <FormControl>
                    <Input {...field} autoComplete="name" className="h-11" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="grid gap-5 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input {...field} type="email" autoComplete="email" className="h-11" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="phone"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Phone</FormLabel>
                    <FormControl>
                      <Input {...field} type="tel" autoComplete="tel" className="h-11" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <PasswordField field={field} label="password" autoComplete="new-password" />
                  <PasswordRules value={password ?? ''} />
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Confirm password</FormLabel>
                  <PasswordField field={field} label="confirmation" autoComplete="new-password" />
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" disabled={submitting} className="h-11 w-full font-medium">
              {submitting ? 'Creating account...' : 'Create account'}
            </Button>
          </form>
        </Form>
      </div>
    </div>
  )
}
