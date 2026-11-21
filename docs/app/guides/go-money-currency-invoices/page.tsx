import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Money in Go: minor units, ISO 4217 and gap-free invoice numbers',
  description:
    'Why float is wrong, what zero-decimal and three-decimal currencies break, allocating a remainder without losing a cent, and sequences with no gaps.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-money-currency-invoices' },
  openGraph: {
    title: 'Money in Go: minor units, ISO 4217 and gap-free invoice numbers',
    description: 'Storing, splitting and numbering money so an accountant agrees with you.',
    url: 'https://gritframework.dev/guides/go-money-currency-invoices',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-money-currency-invoices"
      lead={
        <>
          <p className="mb-4">
            Money looks like a number and behaves like a different kind of object. It has a
            currency, it cannot be divided freely, it has to round in the direction somebody
            agreed to, and the total of the parts must equal the whole, exactly, every time.
          </p>
          <p>
            Getting this wrong produces errors of a cent, which is worse than a large error: it
            survives testing, reaches an accountant, and takes a day to trace.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A money representation that cannot silently lose precision.</p>
        <p>Allocation and tax that add up exactly, including the awkward thirds.</p>
        <p>Invoice numbers with no gaps, which is a legal requirement in several countries.</p>
      </Outcome>

      <Section id="float" title="Why not float">
        <CodeBlock
          code={`total := 0.1 + 0.2
fmt.Println(total)          // 0.30000000000000004
fmt.Println(total == 0.3)   // false

// At a scale that matters:
var sum float64
for i := 0; i < 100000; i++ {
	sum += 0.01
}
fmt.Println(sum)            // 1000.0000000001588`}
        />
        <p>
          Binary floating point cannot represent 0.1 exactly, for the same reason decimal cannot
          represent a third. The error is tiny and it accumulates, and the places it surfaces are
          the worst ones: a reconciliation that is off by a cent, a balance that will not reach
          zero, a total that disagrees with the sum of its lines.
        </p>
      </Section>

      <Section id="minor" title="Store minor units, as an integer">
        <CodeBlock
          code={`type Money struct {
	// The smallest unit of the currency: cents for USD, pence for GBP,
	// and whole yen for JPY. int64 reaches about 92 quadrillion cents,
	// which is enough.
	Amount   int64  \`json:"amount"\`
	Currency string \`json:"currency"\`   // ISO 4217: "USD", "JPY", "KWD"
}

// 1234 USD cents
price := Money{Amount: 1234, Currency: "USD"}`}
        />
        <Checklist
          items={[
            <>
              <strong>Exact.</strong> Integers add and subtract with no error, ever.
            </>,
            <>
              <strong>Unambiguous in JSON.</strong> <code>1234</code> is 1234 cents, in a way
              that <code>12.34</code> is not, since a JSON number is a double on the other side.
            </>,
            <>
              <strong>Agrees with the payment providers.</strong> Stripe, Adyen and most others
              take and return minor units, so there is no conversion at the boundary.
            </>,
          ]}
        />
        <Warn title="Not every currency has two decimal places">
          <p>
            Hardcoding division by 100 is the most common money bug after float. Japanese yen and
            Korean won have none, so 1000 is one thousand yen rather than ten. Kuwaiti dinar,
            Bahraini dinar and Omani rial have three, so 1000 is one dinar. Chilean peso has
            none. The exponent is per currency, from ISO 4217, and it is a lookup rather than a
            constant.
          </p>
        </Warn>
        <CodeBlock
          code={`// The exponents that are not 2. Everything else is 2.
var exponents = map[string]int{
	"JPY": 0, "KRW": 0, "CLP": 0, "VND": 0, "ISK": 0, "UGX": 0,
	"BHD": 3, "KWD": 3, "OMR": 3, "JOD": 3, "TND": 3,
}

func Exponent(currency string) int {
	if e, ok := exponents[currency]; ok {
		return e
	}
	return 2
}`}
        />
        <H3>Arithmetic that refuses to mix currencies</H3>
        <CodeBlock
          code={`// Adding dollars to euros should be a compile-time or a runtime
// refusal, never a silent sum. This is the whole reason for a type
// rather than a bare int64.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("cannot add %s to %s", other.Currency, m.Currency)
	}
	return Money{Amount: m.Amount + other.Amount, Currency: m.Currency}, nil
}

// Multiplying by a quantity is fine: money times a scalar is money.
// Multiplying money by money is not a thing, and there should be no
// method for it.
func (m Money) Times(n int64) Money {
	return Money{Amount: m.Amount * n, Currency: m.Currency}
}`}
        />
      </Section>

      <Section id="allocation" title="Splitting, without losing a cent">
        <p>
          Divide 100 cents three ways. Each gets 33, and one cent is unaccounted for. Rounding
          each to 33.33 and storing it loses the cent; rounding up loses it in the other
          direction. The requirement is that the parts sum exactly to the whole.
        </p>
        <CodeBlock
          code={`// Largest remainder: give everyone the floor, then distribute the
// leftover units one at a time to whoever was rounded down hardest.
// The sum is exact by construction.
func (m Money) Allocate(ratios []int64) ([]Money, error) {
	var total int64
	for _, r := range ratios {
		if r < 0 {
			return nil, errors.New("negative ratio")
		}
		total += r
	}
	if total == 0 {
		return nil, errors.New("ratios sum to zero")
	}

	out := make([]Money, len(ratios))
	remainder := m.Amount
	for i, r := range ratios {
		share := m.Amount * r / total        // integer division: floors
		out[i] = Money{Amount: share, Currency: m.Currency}
		remainder -= share
	}

	// Hand out the remaining units. One each, in order, which is
	// deterministic and therefore testable. Some domains require
	// largest-remainder-first instead; whichever you choose, write it
	// down, because an accountant will eventually ask.
	for i := 0; remainder > 0; i, remainder = i+1, remainder-1 {
		out[i%len(out)].Amount++
	}
	return out, nil
}`}
        />
        <CodeBlock
          code={`// 100 cents, three ways → [34, 33, 33]. Sums to 100.
// 100 cents, 1:1:2    → [25, 25, 50].`}
        />
        <p>
          The same function covers splitting a payment across invoices, a discount across line
          items, and a shipping charge across a basket. It is worth having once, tested, rather
          than open-coded in three places with three different rounding behaviours.
        </p>
      </Section>

      <Section id="tax" title="Tax, and the order of operations">
        <p>
          Whether tax is computed per line or on the total changes the answer by a cent or two,
          and which is correct depends on the jurisdiction rather than on preference.
        </p>
        <CodeBlock
          code={`// Per line, then summed. Three lines at 1.00 with 7.5% tax:
//   round(7.5) = 8 each → 24 total
line := 100
tax := roundHalfUp(line * 75, 1000)     // 8
total := 3 * (line + tax)               // 324

// On the total:
//   300 × 7.5% = 22.5 → 23 → 323
subtotal := 3 * 100
tax = roundHalfUp(subtotal*75, 1000)    // 23
total = subtotal + tax                  // 323

// A cent apart, and both are "right" somewhere. Decide, write it in a
// comment with the reason, and test it.`}
        />
        <CodeBlock
          code={`// Round half up, in integers, with no float anywhere in the path.
func roundHalfUp(numerator, denominator int64) int64 {
	if numerator >= 0 {
		return (numerator + denominator/2) / denominator
	}
	return -((-numerator + denominator/2) / denominator)
}`}
        />
        <Note title="Store the rate you used, not just the result">
          <p>
            Tax rates change. An invoice from 2024 must still show the rate that applied in 2024,
            and recomputing it from today&apos;s configuration produces a figure that disagrees
            with the document the customer has. Store the rate and the computed amount on the
            line, and treat both as immutable once issued.
          </p>
        </Note>
      </Section>

      <Section id="gorm" title="In the database, and in GORM">
        <CodeBlock
          language="sql"
          code={`-- bigint for the amount, char(3) for the currency. Not numeric: an
-- integer is exact, smaller and faster, and numeric invites somebody
-- to insert a fractional value.
total_amount   bigint      NOT NULL,
total_currency char(3)     NOT NULL,

CONSTRAINT total_non_negative CHECK (total_amount >= 0)`}
        />
        <CodeBlock
          code={`// A composite type stored as two columns, so the currency can never
// be separated from the amount.
type Money struct {
	Amount   int64
	Currency string
}

type Invoice struct {
	ID    uint
	Total Money \`gorm:"embedded;embeddedPrefix:total_"\`
}
// → total_amount, total_currency`}
        />
        <Warn title="A BeforeUpdate hook is overwritten by Save">
          <p>
            A GORM hook that calls <code>SetColumn</code> to normalise an amount is discarded
            when the write is <code>Save(&amp;struct)</code>, and honoured when it is{' '}
            <code>Updates(map)</code>. Money fields are a common place for such hooks, and the
            resulting bug is that the normalisation works on some write paths and not others,
            with nothing reporting it.
          </p>
        </Warn>
      </Section>

      <Section id="display" title="Formatting, at the edge only">
        <p>
          Convert to a decimal at the moment of display and never before. An amount that becomes
          a float on its way through the system has already lost the property you chose integers
          for.
        </p>
        <CodeBlock
          language="ts"
          code={`// The frontend receives minor units and the currency, and formats with
// the platform's own rules, which know about placement, separators and
// the exponent.
export function formatMoney(amount: number, currency: string, locale = 'en-GB') {
  const exponent = new Intl.NumberFormat(locale, { style: 'currency', currency })
    .resolvedOptions().maximumFractionDigits
  return new Intl.NumberFormat(locale, { style: 'currency', currency })
    .format(amount / 10 ** exponent)
}

formatMoney(123456, 'USD')          // $1,234.56
formatMoney(123456, 'JPY', 'ja-JP') // ￥123,456   no division by 100
formatMoney(123456, 'KWD', 'ar-KW') // 123.456 KWD, three places`}
        />
        <p>
          Asking <code>Intl</code> for the exponent rather than hardcoding it is what makes the
          yen and dinar cases work without a special branch.
        </p>
        <p>
          In a form, group the digits as the user types: <code>50000000</code> is unreadable and{' '}
          <code>50,000,000</code> is not, and a missing zero in a payment field is an expensive
          typo.
        </p>
      </Section>

      <Section id="sequences" title="Invoice numbers with no gaps">
        <p>
          In several jurisdictions, including much of the EU, invoice numbers must be sequential
          with no gaps. A gap has to be explainable, and &quot;the transaction rolled back&quot;
          is not a satisfactory explanation to a tax authority.
        </p>
        <CodeBlock
          code={`// A Postgres sequence will NOT do. Sequences are deliberately
// non-transactional: nextval() does not roll back, so a failed insert
// consumes a number permanently. That is correct behaviour for
// surrogate keys and wrong for invoice numbers.

// A row with a lock does the job: the increment is part of your
// transaction, so a rollback returns the number.
func (s *Sequences) Next(ctx context.Context, tx *gorm.DB, orgID uint, kind string) (int64, error) {
	var seq models.Sequence
	err := tx.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).      // SELECT ... FOR UPDATE
		Where("organisation_id = ? AND kind = ?", orgID, kind).
		First(&seq).Error
	if err != nil {
		return 0, fmt.Errorf("lock sequence %s: %w", kind, err)
	}

	seq.Current++
	if err := tx.Save(&seq).Error; err != nil {
		return 0, fmt.Errorf("advance sequence %s: %w", kind, err)
	}
	return seq.Current, nil
}`}
        />
        <Checklist
          items={[
            <>
              <strong>It serialises.</strong> Two concurrent invoice creations for the same
              organisation queue behind each other. That is the cost of the guarantee, and
              invoices are not a high-throughput operation.
            </>,
            <>
              <strong>Lock per tenant and per kind,</strong> so one organisation&apos;s invoicing
              does not block another&apos;s, and invoices do not block credit notes.
            </>,
            <>
              <strong>Hold the lock for as little as possible.</strong> Take the number last, not
              first: a transaction that locks the sequence and then calls a payment provider
              holds it for the length of an HTTP request.
            </>,
            <>
              <strong>Assign the number at issue, not at draft.</strong> A deleted draft that
              consumed a number is the gap you were avoiding.
            </>,
          ]}
        />
        <Note title="Reset per period, if the format says so">
          <p>
            <code>INV-2026-0001</code> means the counter resets each year, so the sequence row is
            keyed by period as well as by tenant and kind. Decide the format before the first
            invoice is issued: changing it later means either a visible discontinuity or
            renumbering documents customers already have.
          </p>
        </Note>
      </Section>

      <Section id="multi" title="More than one currency">
        <p>Two rules, and the second is the one people skip.</p>
        <Checklist
          items={[
            <>
              <strong>Never sum across currencies in the database.</strong> A dashboard adding
              USD and EUR amounts because both are <code>bigint</code> produces a number with no
              meaning. Group by currency, and if a single total is required, convert explicitly
              at a stated rate.
            </>,
            <>
              <strong>Store the rate with the transaction.</strong> A payment of 100 EUR recorded
              against a USD account needs the amount, the currency, the rate used and when it was
              taken. Recomputing later with today&apos;s rate gives a different figure, and the
              difference is a real gain or loss that accounting wants to see.
            </>,
          ]}
        />
        <CodeBlock
          code={`type Payment struct {
	Amount     Money     // what was actually paid: 10000 EUR
	BaseAmount Money     // in the account's currency: 10850 USD
	Rate       int64     // 108500, scaled by 10^5, not a float
	RateAt     time.Time
}`}
        />
      </Section>

      <Section id="with-grit" title="As a field type">
        <CodeBlock
          terminal
          code={`grit generate resource Invoice \\
  --fields "number:string:unique,total:money,issued_at:datetime"

grit generate sequence invoice --format "INV-{YYYY}-{0000}" --reset yearly`}
        />
        <p>
          A <code>money</code> field writes the two columns, the Go type with currency-aware
          arithmetic, the non-negative constraint, the TypeScript type, the Zod schema, a form
          input that groups digits as you type, and a right-aligned, correctly formatted table
          column.
        </p>
        <p>
          The sequence generator writes the locked-row implementation above. One thing to know
          about it: the auto-numbering hook calls the sequence package directly rather than going
          through a service, because models cannot import services without a cycle, and the
          workaround for that cycle is worse than the direct call.
        </p>
        <p>
          <Link href="/docs/concepts/money" className="text-accent hover:underline">
            The money reference
          </Link>{' '}
          covers the field type, and{' '}
          <Link href="/docs/backend/invoices" className="text-accent hover:underline">
            the invoices reference
          </Link>{' '}
          covers numbering and documents.
        </p>
      </Section>

      <Faq
        slug="go-money-currency-invoices"
        items={[
          {
            q: 'Integer minor units or a decimal library?',
            a: (
              <>
                Integers for stored amounts: exact, small, and the same representation payment
                providers use. A decimal type such as shopspring/decimal is worth it for
                intermediate calculations with fractional rates, interest accrual or
                multi-step percentages, where you want more precision before the final round.
                Convert back to integer minor units before storing.
              </>
            ),
            text:
              'Integers for stored amounts: exact, compact and the same representation payment providers use. A decimal library is worth it for intermediate calculations with fractional rates or interest, where more precision is wanted before the final round. Convert back to integer minor units before storing.',
          },
          {
            q: 'Why not use numeric in Postgres?',
            a: (
              <>
                <code>numeric</code> is exact and is a reasonable choice. Integers win on three
                counts: they are smaller and faster, they match what payment APIs exchange, and
                they make a fractional cent impossible to insert rather than merely discouraged.
                The main argument for <code>numeric</code> is if non-application tooling reads
                the column and expects a decimal.
              </>
            ),
            text:
              'numeric is exact and reasonable. Integers are smaller and faster, match what payment APIs exchange, and make a fractional cent impossible rather than merely discouraged. numeric is preferable mainly when other tooling reads the column expecting a decimal.',
          },
          {
            q: 'How do I handle a currency with no minor unit?',
            a: (
              <>
                Store the amount in whatever the smallest unit is, which for yen is the yen
                itself, and look the exponent up by currency instead of dividing by 100. The code
                does not special-case anything; only the display and parsing consult the
                exponent. Hardcoded hundreds are the bug.
              </>
            ),
            text:
              'Store the amount in the smallest unit, which for yen is the yen itself, and look up the exponent per currency rather than dividing by 100. Only display and parsing consult the exponent, and hardcoded hundreds are the bug.',
          },
          {
            q: 'Do invoice numbers really have to be gap-free?',
            a: (
              <>
                In much of the EU and several other jurisdictions, yes, and a gap must be
                explainable. Elsewhere it is often only a convention. Check the rule for where
                you invoice before choosing, because the gap-free version serialises invoice
                creation and the non-gap-free version does not, and that is a decision worth
                making deliberately rather than discovering during an audit.
              </>
            ),
            text:
              'In much of the EU and several other jurisdictions, yes, and a gap must be explainable. Elsewhere it is often only a convention. Check the rule for where you invoice, since the gap-free version serialises invoice creation and the alternative does not.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/build-saas-mvp-weekend',
            label: 'Build a SaaS MVP in a weekend',
            blurb: 'Where the billing model is chosen, which decides what you must count.',
          },
          {
            href: '/guides/go-saga-pattern',
            label: 'Sagas and compensating transactions',
            blurb: 'Refunds, and why a compensation is a second transaction rather than a rollback.',
          },
          {
            href: '/docs/concepts/money',
            label: 'The money field type',
            blurb: 'What it generates across the stack.',
          },
          {
            href: '/docs/backend/invoices',
            label: 'Invoices',
            blurb: 'Numbering, documents and the sequence implementation.',
          },
        ]}
      />
    </GuidePage>
  )
}
