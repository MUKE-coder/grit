#import "../template.typ": *

= Money, dates, and the types that prevent bugs

There is a class of bug you cannot test your way out of, because the test and
the code share the same wrong assumption. A total that is a cent off. A birthday
that is a day early for half your users. A currency conversion that is a hundred
times too large and looks plausible.

The defence is not care. It is a type that cannot hold the wrong thing. This
chapter is the two the framework ships, and the general move behind both.

== Why a float is the wrong type for money

```go
fmt.Println(0.1 + 0.2)        // 0.30000000000000004
fmt.Println(1.1 * 3)          // 3.3000000000000003
```

This is not a Go problem. It is binary floating point, and JavaScript does the
same thing, which is why you have seen `toFixed(2)` sprinkled through every
invoice page you have ever worked on.

#bridge[
  `toFixed(2)` is a *display* fix applied to a *storage* problem. The error is
  already in the number; rounding it at the end hides one cent per operation
  and accumulates in the direction nobody notices until a month-end total is
  off by four pounds and somebody has to explain it.

  The fix is the same in both languages: store integers of the smallest unit,
  and keep the decimal point as a rendering concern.
]

== Money: an integer and a currency

```go
type Money struct {
    Amount   int64  `json:"amount"`
    Currency string `json:"currency" gorm:"size:3"`
}
```

`Amount` is in *minor units*: 2800 is twenty eight dollars. The currency rides
with it, because an amount without one is not a quantity of anything.

A `money` field embeds this as two columns with the field name as prefix:

```bash
grit generate resource Invoice --fields "total:money"
```

gives `total_amount` and `total_currency`.

#gotcha[
  Two columns rather than one string, for a reason you feel at the first
  report: a column holding `"19.99 USD"` is not something you can `SUM`, index
  or compare in SQL.

  The amount stays an integer the database understands, and the currency is
  queryable, so "total revenue by currency this quarter" is one query rather
  than a loop in application code.
]

== The operations, and which ones can fail

```go
money.New(2800, "USD")          // from minor units
money.FromMajor(28.00, "USD")   // from a figure a person typed

m.Major()                       // 28.00, for display only
m.String()                      // "$28.00"

sum, err := a.Add(b)            // errors on mixed currencies
diff, err := a.Sub(b)
line := unit.MulInt(qty)        // a line total
tax := line.MulFloat(0.2)       // a percentage
parts := total.Allocate(3)      // a split that loses nothing
eur, err := usd.Convert("EUR", "0.9213")
```

Three of these are worth stopping on.

*`FromMajor` is the one place a float is allowed near money*, because a person
typed `28.00` into a form. It is converted immediately rather than carried, and
`Major()` goes the other way for display. Feeding `Major()` back into
arithmetic is how the cent you were protecting gets lost.

*`Add` returns an error*, and the error is a currency mismatch.

#gotcha[
  You will be tempted to ignore that error, because your application has one
  currency. Do not.

  The day it has two, an ignored error is not a crash and not a validation
  message. It is a total that is silently wrong, computed from amounts that
  were never comparable, and it will be wrong in reports that have already been
  sent.

  This is the clearest case in the book for chapter four's rule: the error
  exists because the operation genuinely cannot be defined, and `_` is a
  decision to produce nonsense.
]

*`Allocate` is the one people write by hand and get wrong.* Splitting 10.00
three ways gives 3.34, 3.33, 3.33. Three lots of 3.33 loses a cent, and in
accounting a lost cent is not a rounding detail, it is a reconciliation that
fails. The remainder goes to the earliest parts, which is the convention
accounting systems expect.

== Not every currency has two decimal places

```go
var exponents = map[string]int{
    "BIF": 0, "CLP": 0, "DJF": 0, "GNF": 0, "ISK": 0, "JPY": 0,
    "KMF": 0, "KRW": 0, "PYG": 0, "RWF": 0, "UGX": 0, "VND": 0,
    "XAF": 0, "XOF": 0, "XPF": 0,
    "BHD": 3, "IQD": 3, "JOD": 3, "KWD": 3, "LYD": 3, "OMR": 3,
    "TND": 3,
}
```

Yen and Ugandan shillings have no minor unit. Kuwaiti dinars have three. Every
other currency is assumed to have two, which is true of most of ISO 4217.

#gotcha[
  A hardcoded `amount / 100` anywhere in your frontend is a bug waiting for its
  first international customer. Dividing a UGX amount by a hundred turns fifty
  thousand shillings into five hundred, on the invoice, in the email, in the
  report.

  Format through the shared helper, which knows the exponent. This is a
  one-line habit that saves an embarrassing afternoon.
]

Conversion handles it too, which is the subtle part:

```go
usd := money.New(1000, "USD")        // 10.00
jpy, _ := usd.Convert("JPY", "148.2") // 1482 JPY
```

Multiplying the amount by the rate and relabelling the currency would give
148200: a hundred times too much, silently, because USD has two decimals and
JPY has none.

The rate is a *decimal string* rather than a float, and that is deliberate:
1.005 is 1.00499999999999989 in binary, and an exchange rate multiplies the
largest amounts in the system, which is exactly where that error stops being
invisible. The string is parsed as an exact fraction, multiplied once, and
rounded once.

== Dates: the type that exists because of one error message

A date field used to be a `*time.Time`, and `time.Time` accepts RFC 3339 and
nothing else. The admin's date picker sends `"2001-08-06"`. So every date field
failed to save with this:

```
parsing time "2001-08-06" as "2006-01-02T15:04:05Z07:00":
cannot parse "" as "T"
```

The fix was not to make the browser send a timezone it does not have. It was
two wire types that accept what a browser actually sends:

```go
type Date struct{ time.Time }      // "2001-08-06"
type DateTime struct{ time.Time }  // and the browser's variants
```

The accepted layouts, widest first:

```go
time.RFC3339Nano
time.RFC3339
"2006-01-02T15:04:05"   // datetime-local, with seconds
"2006-01-02T15:04"      // datetime-local, without
"2006-01-02"            // a date picker
```

The middle two are what `<input type="datetime-local">` produces: no zone, and
often no seconds. Neither is RFC 3339, and both are what you get from a real
browser.

#bridge[
  The JavaScript habit that causes the worst version of this is building a
  `Date` from a date string and sending that.

  `new Date("2001-08-06")` is parsed as UTC midnight, so a user west of
  Greenwich sees the fifth of August. A date of birth that is a day early for
  everybody in the Americas is a bug that survives for years, because it is
  correct in the office where it was written.

  A date is not a timestamp. Send the string.
]

== date against datetime, and how to choose

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`date`], [A day that is the same day everywhere: a birthday, an invoice
    date, a public holiday, a due date],
  [`datetime`], [A moment on the clock: a login, an appointment, a reading
    from a sensor],
)

The test: *would two people in different countries disagree about what day it
is?* A birthday is the sixth of August in Kampala and in Lima. A meeting at
14:00 is not the same moment in both.

Store timestamps in UTC and render them local. Store dates as dates, with no
zone at all, because a zone on a birthday is a question nobody asked.

#gotcha[
  Date sorting failed silently once, and the way it failed is worth knowing
  because it will happen to you with some other column.

  The sortable-column check tested the Go type for `string`, `int` or `uint`. A
  date is `*jsontime.Date`, so it matched none of them and was never added to
  the whitelist. The pagination layer *drops* a sort it does not recognise
  rather than rejecting it, so the request succeeded, the rows came back in
  whatever order the database felt like, and an invoice list could not be
  sorted by its issue date.

  Nothing errored. Nothing logged. The general lesson: *a whitelist that
  silently ignores what it does not know needs a test that asserts the
  whitelist's contents*, not just a test that the endpoint returns 200.
]

== The same move, in smaller places

Money and time are the big two. The field type list has several more of the
same idea, where choosing the right one closes a bug rather than saving typing:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`email`, `url`, `tel`], [Validated on write, so the column holds something
    sendable],
  [`percent`], [Stored as a number with a known meaning, not "0.2 or 20?"],
  [`slug`], [Lowercased, dashed, unique. A URL is a key],
  [`select`], [A closed set, in the database and in the type],
  [`uint`], [A count that cannot be negative, enforced by the type],
)

#bridge[
  This is the thing functional programmers mean by "make illegal states
  unrepresentable", arrived at from the opposite direction: not from theory,
  but from reading support tickets.

  Every type in that table exists because somebody shipped a string, and then
  shipped a validation, and then shipped a second validation in a second place
  that disagreed with the first.
]

#project(title: "Project 37: an invoice that adds up")[
  Build the arithmetic properly, then try to break it.

  ```bash
  grit new billing && cd billing
  grit generate resource Invoice \
    --fields "number:string, issued:date, due:date, \
client:belongs_to:Client, status:select:draft=Draft|sent=Sent|paid=Paid" \
    --items "InvoiceItem:description:string,qty:int,unit:money"
  grit generate resource Client --fields "name:string, email:email"
  grit migrate
  ```

  *One: line totals.* Add a `Total()` method on the invoice that multiplies
  each line with `MulInt` and sums with `Add`, handling the error. Write a test
  with three lines whose total would be wrong in floats, for example three
  lines of 0.1 each.

  *Two: tax.* Add twenty percent with `MulFloat`, rounded once. Assert the
  exact minor-unit value in a test rather than comparing formatted strings.

  *Three: the split.* A payment plan over three instalments using `Allocate`.
  Assert that the instalments sum back to exactly the total. This is the test
  that catches a hand-rolled division.

  *Four: the mismatch.* Add a line in a different currency and watch `Add`
  return an error. Then ignore the error on purpose, print the total, and look
  at the number you get. That number is what an ignored error buys you.

  *Five: the dates.* Save an invoice through the admin with a due date, then
  read the raw row in `/studio`. Confirm it is a date and not a timestamp at
  midnight in somebody's zone. Then sort the invoice list by `issued` and
  confirm it actually sorts.

  *Six: the exponents.* Create one invoice in JPY and one in KWD. Check both
  render correctly in the admin, in the list and in the detail view. If either
  is wrong by a factor of a hundred, you have found the hardcoded division.
]

#tryit[
  In a scratch file, add `0.1 + 0.2` in Go, in JavaScript and in your database
  with a `NUMERIC` column, and print all three to twenty decimal places.

  The third one is right, which is the argument for `NUMERIC` in the database.
  The first two are the reason the application should not be doing the
  arithmetic in a float on the way there.
]

#assignment(name: "Money and Time - Assignment One")[
  In `part5-building/40-money`, build the invoicing slice and write the tests
  that prove the arithmetic.

  *One.* The invoice with line items and a `Total()` that handles the currency
  error rather than discarding it.

  *Two.* Four tests: a float-unsafe total, tax rounded once, an `Allocate`
  that sums back exactly, and a mixed-currency addition that fails.

  *Three.* `MONEY.md` recording what 10.00 USD converts to in JPY at 148.2, and
  what you would have got by multiplying and relabelling. Show both numbers.

  *Four.* A date field saved from the admin, with the stored row pasted in, and
  one paragraph on why it is a date and not a timestamp.

  *Five.* Find one place in code you have already written where money is a
  float or a date is a timestamp. Describe the specific bug it would produce
  and which user would see it first. If you genuinely cannot find one, say so:
  that is a good sign and a rarer one than you would think.
]
