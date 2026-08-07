'use client'

import { useEffect, useId, useRef, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { AlertCircle } from 'lucide-react'

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
 * A create-or-edit form for one record: text, a pair of prices, a description
 * with a counter, a category and a published switch.
 *
 * Number fields do not coerce empty to zero. The source wrote
 * `parseFloat(e.target.value) || 0`, and every part of that is a trap: clearing
 * the field gives NaN, NaN is falsy, so the value silently becomes 0. A price
 * you deleted reads as free, and "required" can never fire because the field
 * is never empty. One of its fields did the same with `|| 5`, so clearing a
 * rating jumped it to five. Empty becomes undefined here and the schema
 * rejects it.
 *
 * Submitting an invalid form moves focus to a summary that lists what is
 * wrong, with each entry a button that focuses its field. Field-level messages
 * alone leave someone who pressed Save with no idea anything happened,
 * especially when the first bad field is scrolled off the top. This needs
 * shouldFocusError turned off, or react-hook-form races the summary for focus
 * and the winner changes between submits.
 *
 * The character counter is tied to the textarea with aria-describedby and is
 * only a live region in the last fifty characters. Announcing a number on
 * every keystroke is unusable; saying nothing until the limit truncates your
 * sentence is worse.
 *
 * The published control is a checkbox with role="switch", not a Radix Switch.
 * It is a real form control that submits, requires no primitive and no
 * JavaScript to reflect its state, and role="switch" gets it announced as on
 * or off rather than checked or unchecked.
 *
 * The category is a native select. A listbox built out of divs has to
 * reimplement type-ahead, arrow keys, Escape, and the way phones present
 * options natively, and usually reimplements none of it.
 */

const CATEGORIES = ['Bags', 'Audio', 'Watches', 'Footwear', 'Kitchen'] as const

const DESCRIPTION_MAX = 500

/* z.coerce would turn '' into 0, the same bug the source had in its onChange
   handler, so the fields hand up undefined instead. Note this is Zod 4: the
   old required_error and invalid_type_error options are one `error` now. */
const price = z
  .number({ error: 'Enter a price' })
  .min(0.01, { message: 'Price must be more than zero' })

const schema = z
  .object({
    name: z.string().min(1, { message: 'Give the product a name' }).max(100, {
      message: 'Name must be 100 characters or fewer',
    }),
    price,
    salePrice: z
      .number({ error: 'Enter a number or leave this empty' })
      .min(0, { message: 'Sale price cannot be negative' })
      .optional(),
    description: z
      .string()
      .min(10, { message: 'Describe the product in at least 10 characters' })
      .max(DESCRIPTION_MAX, { message: `Description must be ${DESCRIPTION_MAX} characters or fewer` }),
    category: z.enum(CATEGORIES, { error: 'Pick a category' }),
    published: z.boolean(),
  })
  .refine((data) => data.salePrice === undefined || data.salePrice < data.price, {
    message: 'Sale price must be lower than the price',
    path: ['salePrice'],
  })

export type ResourceValues = z.infer<typeof schema>

/** Empty means empty, not zero. */
function toNumber(value: string) {
  return value === '' ? undefined : Number(value)
}

function CharacterCount({ value, max }: { value: string; max: number }) {
  const { formDescriptionId } = useFormField()
  const used = value.length
  const remaining = max - used
  const close = remaining <= 50

  return (
    <p
      id={formDescriptionId}
      /* Silent until the limit is in sight, then it starts speaking. */
      aria-live={close ? 'polite' : 'off'}
      className={`text-sm ${
        remaining < 0
          ? 'text-red-600 dark:text-red-400'
          : close
            ? 'text-amber-700 dark:text-amber-400'
            : 'text-gray-500 dark:text-gray-400'
      }`}
    >
      {remaining < 0
        ? `${Math.abs(remaining)} characters over the limit`
        : `${used} of ${max} characters`}
    </p>
  )
}

export default function ResourceFormWithErrorSummary({
  onSubmit = async () => {},
}: {
  onSubmit?: (values: ResourceValues) => Promise<void> | void
}) {
  const [submitting, setSubmitting] = useState(false)
  const [showSummary, setShowSummary] = useState(false)
  /* Counts rejections rather than tracking a boolean, so the focus effect
     fires again on a second failed submit when the summary was already up. */
  const [rejections, setRejections] = useState(0)
  const summary = useRef<HTMLDivElement>(null)
  const baseId = useId()

  const form = useForm<ResourceValues>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: '',
      price: undefined,
      salePrice: undefined,
      description: '',
      category: undefined,
      published: false,
    },
    mode: 'onTouched',
    /* react-hook-form focuses the first invalid field on a rejected submit by
       default, which fights the summary for focus and wins intermittently.
       Focus belongs in one place, and the summary is the better landing spot:
       it says how many problems there are before dropping you into one. */
    shouldFocusError: false,
  })

  const description = form.watch('description') ?? ''
  const errors = form.formState.errors

  async function handleValid(values: ResourceValues) {
    setShowSummary(false)
    setSubmitting(true)
    try {
      await onSubmit(values)
    } finally {
      setSubmitting(false)
    }
  }

  /* react-hook-form calls this instead of the success handler when the schema
     rejects, which is the only reliable moment to raise the summary. */
  function handleInvalid() {
    setShowSummary(true)
    setRejections((count) => count + 1)
  }

  /* Focused from an effect, not from the handler. A requestAnimationFrame
     inside handleInvalid runs before React has committed the summary, so
     there is nothing to focus yet. */
  useEffect(() => {
    if (rejections > 0) summary.current?.focus()
  }, [rejections])

  const errorEntries = Object.entries(errors) as [keyof ResourceValues, { message?: string }][]

  const fieldClass =
    'h-11 w-full rounded-md border border-gray-300 bg-white px-3 text-sm text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:bg-gray-900 dark:text-white'

  return (
    <div className="bg-gray-50 py-12 dark:bg-gray-950">
      <div className="mx-auto max-w-2xl px-4">
        <div className="rounded-2xl border border-gray-200 bg-white p-8 dark:border-white/10 dark:bg-gray-900">
          <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
            New product
          </h1>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            Everything except the sale price is required.
          </p>

          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(handleValid, handleInvalid)}
              noValidate
              className="mt-8 space-y-6"
            >
              {showSummary && errorEntries.length > 0 && (
                <div
                  ref={summary}
                  tabIndex={-1}
                  role="alert"
                  className="rounded-lg border border-red-300 bg-red-50 p-4 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600 dark:border-red-500/40 dark:bg-red-500/10"
                >
                  <p className="flex items-center gap-2 font-medium text-red-800 dark:text-red-300">
                    <AlertCircle aria-hidden="true" className="size-4" />
                    {errorEntries.length === 1
                      ? 'There is one problem to fix'
                      : `There are ${errorEntries.length} problems to fix`}
                  </p>
                  <ul role="list" className="mt-2 list-disc space-y-1 pl-9 text-sm">
                    {errorEntries.map(([name, error]) => (
                      <li key={name}>
                        {/* Focuses the field rather than jumping the page, so
                            the caret lands where the fix has to happen. */}
                        <button
                          type="button"
                          onClick={() => form.setFocus(name)}
                          className="text-red-800 underline underline-offset-2 hover:no-underline dark:text-red-300"
                        >
                          {error?.message}
                        </button>
                      </li>
                    ))}
                  </ul>
                </div>
              )}

              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input {...field} className="h-11" />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="grid gap-6 sm:grid-cols-2">
                <FormField
                  control={form.control}
                  name="price"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Price</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          value={field.value ?? ''}
                          onChange={(event) => field.onChange(toNumber(event.target.value))}
                          type="number"
                          step="0.01"
                          min="0"
                          inputMode="decimal"
                          className="h-11"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="salePrice"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Sale price (optional)</FormLabel>
                      <FormControl>
                        <Input
                          {...field}
                          value={field.value ?? ''}
                          onChange={(event) => field.onChange(toNumber(event.target.value))}
                          type="number"
                          step="0.01"
                          min="0"
                          inputMode="decimal"
                          className="h-11"
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <textarea
                        {...field}
                        rows={5}
                        className="w-full rounded-md border border-gray-300 bg-white p-3 text-sm text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:bg-gray-900 dark:text-white"
                      />
                    </FormControl>
                    <CharacterCount value={description} max={DESCRIPTION_MAX} />
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="category"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Category</FormLabel>
                    <FormControl>
                      <select {...field} value={field.value ?? ''} className={fieldClass}>
                        <option value="" disabled>
                          Choose a category
                        </option>
                        {CATEGORIES.map((category) => (
                          <option key={category} value={category}>
                            {category}
                          </option>
                        ))}
                      </select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="published"
                render={({ field }) => (
                  <FormItem>
                    <div className="flex items-center justify-between rounded-lg border border-gray-200 p-4 dark:border-white/10">
                      <div>
                        <FormLabel>Published</FormLabel>
                        <p
                          id={`${baseId}-published-hint`}
                          className="mt-1 text-sm text-gray-500 dark:text-gray-400"
                        >
                          Visible in the storefront as soon as you save.
                        </p>
                      </div>

                      <FormControl>
                        {/* A checkbox with role="switch": announced as on or
                            off, submits with the form, needs no primitive. */}
                        <input
                          type="checkbox"
                          role="switch"
                          name={field.name}
                          ref={field.ref}
                          checked={field.value}
                          onBlur={field.onBlur}
                          onChange={(event) => field.onChange(event.target.checked)}
                          aria-describedby={`${baseId}-published-hint`}
                          className="size-6 shrink-0 accent-indigo-600"
                        />
                      </FormControl>
                    </div>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="flex justify-end gap-3 border-t border-gray-200 pt-6 dark:border-white/10">
                <Button type="button" variant="outline" className="h-11">
                  Cancel
                </Button>
                <Button type="submit" disabled={submitting} className="h-11 font-medium">
                  {submitting ? 'Saving...' : 'Save product'}
                </Button>
              </div>
            </form>
          </Form>
        </div>
      </div>
    </div>
  )
}
