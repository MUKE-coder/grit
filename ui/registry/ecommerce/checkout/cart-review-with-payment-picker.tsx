'use client'

import { useMemo, useState } from 'react'
import { Gift, MapPin, Minus, Plus, Tag, Trash2 } from 'lucide-react'

/*
 * The one-page checkout: review the basket, pick a card, apply a code, pay.
 *
 * Everything on this screen is arithmetic somebody is about to be charged for,
 * so the decisions are about not lying rather than about layout.
 *
 * Totals are derived on every render, never accumulated. The tempting version
 * keeps a `total` in state and adjusts it as quantities change, and it drifts
 * the first time a line is removed while a discount is applied. Here the
 * subtotal is a fold over the lines and everything below it is computed from
 * that, so the sum on screen is always the sum of what is on screen.
 *
 * Money is formatted once, by Intl, in the store's currency. Any component that
 * writes "SAR " + n.toFixed(2) has quietly decided where the symbol goes and
 * what the separator is for every locale it will ever be shown in.
 *
 * The quantity steppers are buttons that name their line. Six identical plus
 * signs is a screen reader reading "plus, plus, plus" down a basket, and the
 * customer has no idea which jacket they just added.
 *
 * The credit toggle states the shortfall rather than hiding it. Applying store
 * credit that does not cover the order is the moment a customer needs to be
 * told they still owe something, and it is exactly the moment most checkouts
 * say nothing.
 *
 * A discount code that is accepted says so in a live region, because the change
 * it makes is a number further down the page that nobody is looking at.
 */

export interface CartLine {
  id: string
  name: string
  brand?: string
  /** e.g. "XL / Black". Rendered as-is under the name. */
  variant?: string
  unitPrice: number
  quantity: number
  image: string
}

export interface PaymentMethod {
  id: string
  label: string
  /** e.g. "Visa" or "Wrap your items". Shown under the label. */
  detail?: string
}

const LINES: CartLine[] = [
  { id: '1', name: 'Dri-Fit training jacket, summer kit', brand: 'Nike', variant: 'XL / Black', unitPrice: 40, quantity: 2, image: 'https://images.unsplash.com/photo-1551028719-00167b16eac5?w=200&h=200&fit=crop&q=80' },
  { id: '2', name: 'Varsity bomber jacket', brand: 'Nike', variant: 'XL / Black', unitPrice: 40, quantity: 2, image: 'https://images.unsplash.com/photo-1521223890158-f9f7c3d5d504?w=200&h=200&fit=crop&q=80' },
  { id: '3', name: 'Fleece crew sweatshirt', brand: 'Nike', variant: 'XL / Grey', unitPrice: 40, quantity: 2, image: 'https://images.unsplash.com/photo-1556821840-3a63f95609a7?w=200&h=200&fit=crop&q=80' },
]

const METHODS: PaymentMethod[] = [
  { id: 'visa', label: 'Visa ending 0912', detail: 'Expires 09 / 28' },
  { id: 'mastercard', label: 'Mastercard ending 0912', detail: 'Expires 04 / 27' },
  { id: 'later', label: 'Pay in instalments', detail: 'Four payments, no interest' },
]

export default function CartReviewWithPaymentPicker({
  currency = 'SAR',
  locale = 'en',
  lines: initialLines = LINES,
  methods = METHODS,
  availableCredit = 4000,
  giftWrapPrice = 20,
  vatRate = 0.15,
  onCheckout,
}: {
  currency?: string
  locale?: string
  lines?: CartLine[]
  methods?: PaymentMethod[]
  availableCredit?: number
  giftWrapPrice?: number
  /** As a fraction. 0.15 is 15%. */
  vatRate?: number
  onCheckout?: (summary: { lines: CartLine[]; total: number }) => void
}) {
  const [lines, setLines] = useState(initialLines)
  const [method, setMethod] = useState(methods[0]?.id ?? '')
  const [useCredit, setUseCredit] = useState(false)
  const [giftWrap, setGiftWrap] = useState(false)
  const [code, setCode] = useState('')
  const [appliedCode, setAppliedCode] = useState<string | null>(null)

  const money = useMemo(
    () => new Intl.NumberFormat(locale, { style: 'currency', currency, maximumFractionDigits: 2 }),
    [locale, currency],
  )

  /* Derived, never accumulated. A total held in state drifts the first time a
     line is removed while a discount is applied. */
  const subtotal = lines.reduce((sum, line) => sum + line.unitPrice * line.quantity, 0)
  const discount = appliedCode ? subtotal * 0.2 : 0
  const wrap = giftWrap ? giftWrapPrice : 0
  const vat = (subtotal - discount + wrap) * vatRate
  const total = subtotal - discount + wrap + vat
  const creditApplied = useCredit ? Math.min(availableCredit, total) : 0
  const due = total - creditApplied

  function setQuantity(id: string, quantity: number) {
    setLines((current) =>
      quantity <= 0
        ? current.filter((line) => line.id !== id)
        : current.map((line) => (line.id === id ? { ...line, quantity } : line)),
    )
  }

  return (
    <div className="bg-gray-50 py-8 dark:bg-gray-950">
      <div className="mx-auto max-w-6xl px-4 sm:px-6 lg:px-8">
        <h1 className="text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">
          Checkout
        </h1>

        <div className="mt-6 grid gap-6 lg:grid-cols-3">
          <div className="space-y-6 lg:col-span-2">
            <section
              aria-labelledby="address-heading"
              className="rounded-xl bg-white p-6 text-center dark:bg-gray-900"
            >
              <h2 id="address-heading" className="sr-only">
                Delivery address
              </h2>
              <span className="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 dark:bg-gray-800">
                <MapPin className="h-5 w-5 text-gray-500" aria-hidden="true" />
              </span>
              <p className="mt-3 font-medium text-gray-900 dark:text-white">No address saved</p>
              <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
                Add one so we can track the delivery.
              </p>
              <button
                type="button"
                className="mt-4 min-h-11 rounded-lg bg-teal-500 px-5 text-sm font-semibold text-white hover:bg-teal-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal-600"
              >
                Add new location
              </button>
            </section>

            <section aria-labelledby="cart-heading" className="rounded-xl bg-white p-6 dark:bg-gray-900">
              <div className="flex items-baseline justify-between">
                <h2 id="cart-heading" className="text-base font-semibold text-gray-900 dark:text-white">
                  Cart{' '}
                  <span className="font-normal text-gray-500 dark:text-gray-400">
                    {lines.length} item{lines.length === 1 ? '' : 's'}
                  </span>
                </h2>
                <button
                  type="button"
                  onClick={() => setLines([])}
                  className="flex min-h-9 items-center gap-1.5 text-sm text-red-600 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
                >
                  <Trash2 className="h-4 w-4" aria-hidden="true" />
                  Remove all
                </button>
              </div>

              {lines.length === 0 ? (
                <p className="py-8 text-center text-sm text-gray-600 dark:text-gray-400">
                  Your cart is empty.
                </p>
              ) : (
                <ul className="mt-4 divide-y divide-gray-100 dark:divide-gray-800">
                  {lines.map((line) => (
                    <li key={line.id} className="flex flex-wrap items-center gap-4 py-4">
                      <img
                        src={line.image}
                        alt=""
                        loading="lazy"
                        className="h-16 w-16 shrink-0 rounded-lg object-cover"
                      />

                      <div className="min-w-0 flex-1">
                        {line.brand && (
                          <p className="text-xs font-semibold text-gray-900 dark:text-white">
                            {line.brand}
                          </p>
                        )}
                        <p className="truncate text-sm text-gray-800 dark:text-gray-200">{line.name}</p>
                        {line.variant && (
                          <p className="text-xs text-gray-500 dark:text-gray-400">{line.variant}</p>
                        )}
                        <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                          {money.format(line.unitPrice)} per item
                        </p>
                      </div>

                      {/* Each stepper names its line. Six identical plus signs
                          is a screen reader saying "plus, plus, plus". */}
                      <div className="flex items-center gap-1 rounded-full border border-gray-200 dark:border-gray-700">
                        <button
                          type="button"
                          onClick={() => setQuantity(line.id, line.quantity - 1)}
                          className="flex h-9 w-9 items-center justify-center rounded-full text-gray-600 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal-600 dark:text-gray-300 dark:hover:bg-gray-800"
                        >
                          <Minus className="h-4 w-4" aria-hidden="true" />
                          <span className="sr-only">Remove one {line.name}</span>
                        </button>
                        <span className="min-w-6 text-center text-sm font-medium text-gray-900 dark:text-white">
                          {line.quantity}
                        </span>
                        <button
                          type="button"
                          onClick={() => setQuantity(line.id, line.quantity + 1)}
                          className="flex h-9 w-9 items-center justify-center rounded-full text-gray-600 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal-600 dark:text-gray-300 dark:hover:bg-gray-800"
                        >
                          <Plus className="h-4 w-4" aria-hidden="true" />
                          <span className="sr-only">Add one {line.name}</span>
                        </button>
                      </div>

                      <p className="w-24 text-right text-sm font-semibold text-gray-900 dark:text-white">
                        {money.format(line.unitPrice * line.quantity)}
                      </p>
                    </li>
                  ))}
                </ul>
              )}
            </section>
          </div>

          <div className="space-y-6">
            <fieldset className="rounded-xl bg-white p-6 dark:bg-gray-900">
              <legend className="text-base font-semibold text-gray-900 dark:text-white">
                Choose how to pay
              </legend>
              <div className="mt-4 space-y-2">
                {methods.map((option) => (
                  <label
                    key={option.id}
                    className={`flex min-h-14 cursor-pointer items-center gap-3 rounded-lg border p-3 transition-colors ${
                      method === option.id
                        ? 'border-teal-500 bg-teal-50 dark:bg-teal-950'
                        : 'border-gray-200 hover:border-gray-300 dark:border-gray-700'
                    }`}
                  >
                    <input
                      type="radio"
                      name="payment-method"
                      checked={method === option.id}
                      onChange={() => setMethod(option.id)}
                      className="h-4 w-4 border-gray-300 text-teal-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal-600"
                    />
                    <span className="min-w-0">
                      <span className="block text-sm font-medium text-gray-900 dark:text-white">
                        {option.label}
                      </span>
                      {option.detail && (
                        <span className="block text-xs text-gray-500 dark:text-gray-400">
                          {option.detail}
                        </span>
                      )}
                    </span>
                  </label>
                ))}
              </div>
            </fieldset>

            <section aria-labelledby="summary-heading" className="rounded-xl bg-white p-6 dark:bg-gray-900">
              <h2 id="summary-heading" className="sr-only">
                Order summary
              </h2>

              <label className="flex items-start gap-3">
                <input
                  type="checkbox"
                  checked={useCredit}
                  onChange={(e) => setUseCredit(e.target.checked)}
                  className="mt-1 h-4 w-4 rounded border-gray-300 text-teal-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal-600"
                />
                <span className="min-w-0 flex-1">
                  <span className="block text-sm font-medium text-gray-900 dark:text-white">
                    Use credit for this purchase
                  </span>
                  <span className="block text-xs text-gray-500 dark:text-gray-400">
                    Available balance {money.format(availableCredit)}
                  </span>
                  {/* Said out loud, because a credit that does not cover the
                      order is exactly when a customer needs telling. */}
                  {useCredit && due > 0 && (
                    <span className="mt-1 block text-xs text-amber-700 dark:text-amber-400">
                      Credit does not cover the order. {money.format(due)} still to pay by your
                      chosen method.
                    </span>
                  )}
                </span>
              </label>

              <label className="mt-4 flex items-start gap-3 border-t border-gray-100 pt-4 dark:border-gray-800">
                <input
                  type="checkbox"
                  checked={giftWrap}
                  onChange={(e) => setGiftWrap(e.target.checked)}
                  className="mt-1 h-4 w-4 rounded border-gray-300 text-teal-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal-600"
                />
                <span>
                  <span className="flex items-center gap-1.5 text-sm font-medium text-gray-900 dark:text-white">
                    <Gift className="h-4 w-4" aria-hidden="true" />
                    Make it a gift
                  </span>
                  <span className="block text-xs text-gray-500 dark:text-gray-400">
                    Wrapped in one box, {money.format(giftWrapPrice)}
                  </span>
                </span>
              </label>

              <div className="mt-4 border-t border-gray-100 pt-4 dark:border-gray-800">
                <label
                  htmlFor="discount-code"
                  className="mb-1.5 flex items-center gap-1.5 text-sm font-medium text-gray-900 dark:text-white"
                >
                  <Tag className="h-4 w-4" aria-hidden="true" />
                  Discount code
                </label>
                <div className="flex gap-2">
                  <input
                    id="discount-code"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    placeholder="Enter a code"
                    className="min-h-11 w-full rounded-lg border border-gray-200 px-3 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                  />
                  <button
                    type="button"
                    onClick={() => setAppliedCode(code.trim() ? code.trim().toUpperCase() : null)}
                    className="min-h-11 shrink-0 rounded-lg bg-gray-900 px-4 text-sm font-semibold text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:bg-white dark:text-gray-900"
                  >
                    Apply
                  </button>
                </div>
                {/* Live: the effect of applying a code is a number further down
                    the page that nobody is looking at. */}
                <p aria-live="polite" className="mt-1.5 min-h-5 text-xs text-teal-700 dark:text-teal-400">
                  {appliedCode ? `${appliedCode} applied, 20% off this order.` : ''}
                </p>
              </div>

              <dl className="mt-4 space-y-2 border-t border-gray-100 pt-4 text-sm dark:border-gray-800">
                <div className="flex justify-between">
                  <dt className="text-gray-600 dark:text-gray-400">
                    Subtotal ({lines.length} item{lines.length === 1 ? '' : 's'})
                  </dt>
                  <dd className="text-gray-900 dark:text-white">{money.format(subtotal)}</dd>
                </div>
                {discount > 0 && (
                  <div className="flex justify-between">
                    <dt className="text-gray-600 dark:text-gray-400">Discount</dt>
                    <dd className="text-teal-700 dark:text-teal-400">-{money.format(discount)}</dd>
                  </div>
                )}
                {wrap > 0 && (
                  <div className="flex justify-between">
                    <dt className="text-gray-600 dark:text-gray-400">Gift wrap</dt>
                    <dd className="text-gray-900 dark:text-white">{money.format(wrap)}</dd>
                  </div>
                )}
                <div className="flex justify-between">
                  <dt className="text-gray-600 dark:text-gray-400">Shipping</dt>
                  <dd className="text-teal-700 dark:text-teal-400">Free</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-gray-600 dark:text-gray-400">VAT</dt>
                  <dd className="text-gray-900 dark:text-white">{money.format(vat)}</dd>
                </div>
                {creditApplied > 0 && (
                  <div className="flex justify-between">
                    <dt className="text-gray-600 dark:text-gray-400">Credit applied</dt>
                    <dd className="text-teal-700 dark:text-teal-400">-{money.format(creditApplied)}</dd>
                  </div>
                )}
                <div className="flex justify-between border-t border-gray-100 pt-2 text-base font-semibold dark:border-gray-800">
                  <dt className="text-gray-900 dark:text-white">Total</dt>
                  <dd className="text-gray-900 dark:text-white">{money.format(due)}</dd>
                </div>
              </dl>

              <button
                type="button"
                disabled={lines.length === 0}
                onClick={() => onCheckout?.({ lines, total: due })}
                className="mt-5 min-h-12 w-full rounded-lg bg-teal-500 text-sm font-semibold text-white hover:bg-teal-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-teal-600 disabled:cursor-not-allowed disabled:opacity-50"
              >
                Checkout
              </button>
            </section>
          </div>
        </div>
      </div>
    </div>
  )
}
