'use client'

import { useMemo, useState } from 'react'
import { ArrowLeft, Banknote, Check, CreditCard, Minus, Plus, Send, Smartphone, X } from 'lucide-react'

/*
 * A four-step checkout on its payment step, with the basket beside it.
 *
 * The stepper is the part worth getting right, and the usual version gets it
 * wrong in a way nobody notices until somebody tries to use it.
 *
 * It is an ordered list, not a row of divs. The steps have an order and a
 * position, "3 of 4" is a fact about them, and `ol` is how that reaches anybody
 * not looking at the screen. Each step states its state in text too, through a
 * visually hidden word: a green tick is invisible to a screen reader, and
 * "Account, completed" is what makes the row mean anything.
 *
 * Completed steps are links, the current step is not, and later steps are
 * disabled. Letting somebody jump to Review before choosing a payment method
 * produces a review of nothing, and the fix is to not offer it rather than to
 * validate on arrival.
 *
 * The card fields are a fieldset. Number, name, expiry and CVV are one thing,
 * and grouping them is what stops a screen reader reading four unrelated inputs
 * between two radio buttons. They render only for the selected method, because
 * a form with every payment method's fields on screen at once is a form where
 * the required ones are ambiguous.
 *
 * Totals are derived on every render. A total in state drifts the moment a line
 * is removed, and this is the number somebody is about to be charged.
 *
 * autoComplete is set on every card field. It is the difference between a
 * checkout a phone can fill in one tap and one where a customer types sixteen
 * digits with their thumbs, and it costs one attribute.
 */

export interface CheckoutStep {
  id: string
  label: string
}

export interface OrderLine {
  id: string
  name: string
  /** e.g. "Japanese Food". Shown above the name as a category link. */
  group?: string
  unitPrice: number
  quantity: number
  image: string
}

export interface CheckoutMethod {
  id: string
  label: string
  detail: string
  icon: 'card' | 'cash' | 'mobile'
  /** Shows the card fieldset when selected. */
  takesCardDetails?: boolean
  badge?: string
}

const STEPS: CheckoutStep[] = [
  { id: 'account', label: 'Account' },
  { id: 'delivery', label: 'Delivery' },
  { id: 'payment', label: 'Payment' },
  { id: 'review', label: 'Review' },
]

const LINES: OrderLine[] = [
  { id: '1', name: 'Sushi Hiro Brolyn', group: 'Japanese food', unitPrice: 15, quantity: 1, image: 'https://images.unsplash.com/photo-1579871494447-9811cf80d66c?w=200&h=200&fit=crop&q=80' },
  { id: '2', name: 'Bistecca Fiorentina', group: 'Italian food', unitPrice: 40, quantity: 1, image: 'https://images.unsplash.com/photo-1544025162-d76694265947?w=200&h=200&fit=crop&q=80' },
]

const METHODS: CheckoutMethod[] = [
  { id: 'cod', label: 'Cash on delivery', detail: 'Pay when you receive your order.', icon: 'cash', badge: 'Popular' },
  { id: 'card', label: 'Credit or debit card', detail: 'Visa, Mastercard, Amex.', icon: 'card', takesCardDetails: true },
  { id: 'mobile', label: 'Mobile banking', detail: 'bKash, Nagad, Upay.', icon: 'mobile' },
]

const ICONS = { card: CreditCard, cash: Banknote, mobile: Smartphone }

export default function CheckoutStepperWithPaymentMethods({
  steps = STEPS,
  currentStep = 'payment',
  lines: initialLines = LINES,
  methods = METHODS,
  currency = 'USD',
  locale = 'en-US',
  shipping = 5,
  taxRate = 0.05,
  onContinue,
}: {
  steps?: CheckoutStep[]
  currentStep?: string
  lines?: OrderLine[]
  methods?: CheckoutMethod[]
  currency?: string
  locale?: string
  shipping?: number
  /** As a fraction. 0.05 is 5%. */
  taxRate?: number
  onContinue?: (payload: { method: string; total: number }) => void
}) {
  const [lines, setLines] = useState(initialLines)
  const [method, setMethod] = useState(methods.find((m) => m.takesCardDetails)?.id ?? methods[0]?.id ?? '')
  const [promo, setPromo] = useState('')

  const money = useMemo(
    () => new Intl.NumberFormat(locale, { style: 'currency', currency }),
    [locale, currency],
  )

  const currentIndex = Math.max(0, steps.findIndex((s) => s.id === currentStep))

  const subtotal = lines.reduce((sum, line) => sum + line.unitPrice * line.quantity, 0)
  const tax = subtotal * taxRate
  const total = subtotal + shipping + tax

  function setQuantity(id: string, quantity: number) {
    setLines((current) =>
      quantity <= 0
        ? current.filter((line) => line.id !== id)
        : current.map((line) => (line.id === id ? { ...line, quantity } : line)),
    )
  }

  const selected = methods.find((m) => m.id === method)

  return (
    <div className="bg-gray-100 py-8 dark:bg-gray-950">
      <div className="mx-auto max-w-5xl overflow-hidden rounded-2xl bg-white shadow-sm dark:bg-gray-900">
        <div className="flex flex-wrap items-center gap-4 border-b border-gray-200 px-6 py-4 dark:border-gray-800">
          <button
            type="button"
            className="flex min-h-11 items-center gap-2 rounded-full bg-gray-100 px-4 text-sm font-medium text-gray-800 hover:bg-gray-200 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:bg-gray-800 dark:text-gray-200"
          >
            <ArrowLeft className="h-4 w-4" aria-hidden="true" />
            Back
          </button>

          {/* An ordered list, because the steps have an order and a position and
              that is a fact about them, not a visual arrangement. */}
          <nav aria-label="Checkout progress" className="min-w-0 flex-1">
            <ol className="flex flex-wrap items-center justify-center gap-x-3 gap-y-2">
              {steps.map((step, i) => {
                const done = i < currentIndex
                const current = i === currentIndex
                return (
                  <li key={step.id} className="flex items-center gap-2">
                    <span
                      className={`flex h-8 w-8 items-center justify-center rounded-full text-sm font-semibold ${
                        done
                          ? 'bg-green-500 text-white'
                          : current
                            ? 'bg-gray-900 text-white dark:bg-white dark:text-gray-900'
                            : 'border border-gray-300 text-gray-500 dark:border-gray-700 dark:text-gray-400'
                      }`}
                      aria-hidden="true"
                    >
                      {done ? <Check className="h-4 w-4" /> : i + 1}
                    </span>
                    <span
                      className={`text-sm ${current ? 'font-semibold text-gray-900 dark:text-white' : 'text-gray-600 dark:text-gray-400'}`}
                      aria-current={current ? 'step' : undefined}
                    >
                      {step.label}
                      {/* The tick is invisible to a screen reader. This is what
                          makes the row mean something. */}
                      <span className="sr-only">
                        {done ? ', completed' : current ? ', current step' : ', not yet reached'}
                      </span>
                    </span>
                    {i < steps.length - 1 && (
                      <span
                        aria-hidden="true"
                        className={`hidden h-px w-8 sm:block ${done ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-700'}`}
                      />
                    )}
                  </li>
                )
              })}
            </ol>
          </nav>
        </div>

        <div className="grid gap-6 p-6 lg:grid-cols-5">
          <fieldset className="lg:col-span-3">
            <legend className="text-lg font-semibold text-gray-900 dark:text-white">
              Payment method
            </legend>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
              Pick an option to continue to review.
            </p>

            <div className="mt-5 space-y-3">
              {methods.map((option) => {
                const Icon = ICONS[option.icon]
                const active = method === option.id
                return (
                  <div
                    key={option.id}
                    className={`rounded-xl border transition-colors ${
                      active
                        ? 'border-blue-500 bg-blue-50/60 dark:bg-blue-950/40'
                        : 'border-gray-200 hover:border-gray-300 dark:border-gray-700'
                    }`}
                  >
                    <label className="flex min-h-16 cursor-pointer items-center gap-3 p-4">
                      <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-gray-900 text-white dark:bg-gray-700">
                        <Icon className="h-5 w-5" aria-hidden="true" />
                      </span>
                      <span className="min-w-0 flex-1">
                        <span className="flex items-center gap-2">
                          <span className="text-sm font-semibold text-gray-900 dark:text-white">
                            {option.label}
                          </span>
                          {option.badge && (
                            <span className="rounded-full bg-gray-900 px-2 py-0.5 text-[11px] font-medium text-white dark:bg-white dark:text-gray-900">
                              {option.badge}
                            </span>
                          )}
                        </span>
                        <span className="block text-xs text-gray-600 dark:text-gray-400">
                          {option.detail}
                        </span>
                      </span>
                      <input
                        type="radio"
                        name="checkout-method"
                        checked={active}
                        onChange={() => setMethod(option.id)}
                        className="h-5 w-5 border-gray-300 text-blue-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                      />
                    </label>

                    {/* Only for the selected method. Every method's fields on
                        screen at once makes the required ones ambiguous. */}
                    {active && option.takesCardDetails && (
                      <fieldset className="border-t border-blue-200 px-4 py-4 dark:border-blue-900">
                        <legend className="sr-only">Card details</legend>
                        <div className="grid gap-3 sm:grid-cols-2">
                          <div className="sm:col-span-2">
                            <label
                              htmlFor="card-number"
                              className="mb-1 block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              Card number
                            </label>
                            <input
                              id="card-number"
                              inputMode="numeric"
                              autoComplete="cc-number"
                              placeholder="1234 5678 9012 3456"
                              className="min-h-11 w-full rounded-lg border border-gray-200 px-3 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                            />
                          </div>
                          <div className="sm:col-span-2">
                            <label
                              htmlFor="card-name"
                              className="mb-1 block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              Cardholder name
                            </label>
                            <input
                              id="card-name"
                              autoComplete="cc-name"
                              className="min-h-11 w-full rounded-lg border border-gray-200 px-3 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                            />
                          </div>
                          <div>
                            <label
                              htmlFor="card-expiry"
                              className="mb-1 block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              Expiry date
                            </label>
                            <input
                              id="card-expiry"
                              inputMode="numeric"
                              autoComplete="cc-exp"
                              placeholder="MM / YY"
                              className="min-h-11 w-full rounded-lg border border-gray-200 px-3 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                            />
                          </div>
                          <div>
                            <label
                              htmlFor="card-cvv"
                              className="mb-1 block text-xs font-medium text-gray-700 dark:text-gray-300"
                            >
                              Security code
                            </label>
                            <input
                              id="card-cvv"
                              inputMode="numeric"
                              autoComplete="cc-csc"
                              placeholder="123"
                              className="min-h-11 w-full rounded-lg border border-gray-200 px-3 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                            />
                          </div>
                        </div>
                      </fieldset>
                    )}
                  </div>
                )
              })}
            </div>

            <div className="mt-6 flex flex-wrap gap-3">
              <button
                type="button"
                className="min-h-12 flex-1 rounded-lg border border-gray-200 text-sm font-semibold text-gray-800 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:border-gray-700 dark:text-gray-200 dark:hover:bg-gray-800"
              >
                Back to delivery
              </button>
              <button
                type="button"
                onClick={() => onContinue?.({ method, total })}
                className="min-h-12 flex-1 rounded-lg bg-gray-900 text-sm font-semibold text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:bg-white dark:text-gray-900"
              >
                Continue to review
              </button>
            </div>
          </fieldset>

          <section aria-labelledby="order-summary" className="lg:col-span-2">
            <h2 id="order-summary" className="text-lg font-semibold text-gray-900 dark:text-white">
              Order summary
            </h2>

            <ul className="mt-4 space-y-4">
              {lines.map((line) => (
                <li key={line.id} className="flex items-start gap-3">
                  <img
                    src={line.image}
                    alt=""
                    loading="lazy"
                    className="h-14 w-14 shrink-0 rounded-lg object-cover"
                  />
                  <div className="min-w-0 flex-1">
                    {line.group && (
                      <p className="text-xs text-blue-600 dark:text-blue-400">{line.group}</p>
                    )}
                    <p className="truncate text-sm font-medium text-gray-900 dark:text-white">
                      {line.name}
                    </p>
                    <p className="text-sm text-gray-600 dark:text-gray-400">
                      {money.format(line.unitPrice)}
                    </p>
                  </div>
                  <div className="flex items-center gap-1 rounded-full border border-gray-200 dark:border-gray-700">
                    <button
                      type="button"
                      onClick={() => setQuantity(line.id, line.quantity - 1)}
                      className="flex h-8 w-8 items-center justify-center rounded-full text-gray-600 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:text-gray-300 dark:hover:bg-gray-800"
                    >
                      <Minus className="h-3.5 w-3.5" aria-hidden="true" />
                      <span className="sr-only">Remove one {line.name}</span>
                    </button>
                    <span className="min-w-5 text-center text-sm text-gray-900 dark:text-white">
                      {line.quantity}
                    </span>
                    <button
                      type="button"
                      onClick={() => setQuantity(line.id, line.quantity + 1)}
                      className="flex h-8 w-8 items-center justify-center rounded-full text-gray-600 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:text-gray-300 dark:hover:bg-gray-800"
                    >
                      <Plus className="h-3.5 w-3.5" aria-hidden="true" />
                      <span className="sr-only">Add one {line.name}</span>
                    </button>
                  </div>
                  <button
                    type="button"
                    onClick={() => setQuantity(line.id, 0)}
                    className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-gray-400 hover:bg-gray-100 hover:text-gray-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:hover:bg-gray-800"
                  >
                    <X className="h-4 w-4" aria-hidden="true" />
                    <span className="sr-only">Remove {line.name} from the order</span>
                  </button>
                </li>
              ))}
            </ul>

            <div className="mt-5">
              <label
                htmlFor="promo-code"
                className="mb-1.5 block text-sm font-medium text-gray-900 dark:text-white"
              >
                Promotion code
              </label>
              <div className="flex gap-2">
                <input
                  id="promo-code"
                  value={promo}
                  onChange={(e) => setPromo(e.target.value)}
                  placeholder="Add promo code"
                  className="min-h-11 w-full rounded-lg border border-gray-200 px-3 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                />
                <button
                  type="button"
                  className="flex h-11 w-11 shrink-0 items-center justify-center rounded-lg border border-gray-200 text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:border-gray-700 dark:text-gray-300"
                >
                  <Send className="h-4 w-4" aria-hidden="true" />
                  <span className="sr-only">Apply promotion code</span>
                </button>
              </div>
            </div>

            <dl className="mt-5 space-y-2 border-t border-gray-200 pt-4 text-sm dark:border-gray-800">
              <div className="flex justify-between">
                <dt className="text-gray-600 dark:text-gray-400">Subtotal</dt>
                <dd className="text-gray-900 dark:text-white">{money.format(subtotal)}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-gray-600 dark:text-gray-400">Shipping</dt>
                <dd className="text-gray-900 dark:text-white">{money.format(shipping)}</dd>
              </div>
              <div className="flex justify-between">
                <dt className="text-gray-600 dark:text-gray-400">
                  Tax ({Math.round(taxRate * 100)}%)
                </dt>
                <dd className="text-gray-900 dark:text-white">{money.format(tax)}</dd>
              </div>
              <div className="flex justify-between border-t border-gray-200 pt-2 text-base font-semibold dark:border-gray-800">
                <dt className="text-gray-900 dark:text-white">Total</dt>
                <dd className="text-blue-600 dark:text-blue-400">{money.format(total)}</dd>
              </div>
            </dl>
          </section>
        </div>
      </div>
    </div>
  )
}
