// Package money is an amount and the currency it is denominated in.
//
// The reason this exists rather than a float64: 0.1 + 0.2 is not 0.3 in binary
// floating point, and after tax, a discount, a split payment and a refund, the
// ledger drifts. Finance finds it before you do. Money is stored and moved as
// an integer count of minor units, which is exactly representable and exactly
// what every payment provider's API expects.
//
//	money.New(1999, "USD")   // $19.99
//	money.New(50000, "UGX")  // USh 50,000 -- UGX has no minor unit
//
// Currency travels with the amount, always. A bare number crossing a function
// boundary is how a USD total ends up added to a UGX one.
package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var (
	ErrCurrencyMismatch = errors.New("cannot combine amounts in different currencies")
	ErrBadCurrency      = errors.New("currency must be a three-letter ISO 4217 code")
)

// exponents lists the currencies whose minor unit is not 1/100.
//
// Only the exceptions are listed; everything absent is assumed to have two
// decimal places, which is true of most of ISO 4217. Getting this wrong is not
// cosmetic: dividing a UGX amount by 100 for display turns 50,000 shillings
// into 500, and a hardcoded amount/100 anywhere in a frontend is a bug waiting
// for its first international customer.
var exponents = map[string]int{
	"BIF": 0, "CLP": 0, "DJF": 0, "GNF": 0, "ISK": 0, "JPY": 0, "KMF": 0,
	"KRW": 0, "PYG": 0, "RWF": 0, "UGX": 0, "UYI": 0, "VND": 0, "VUV": 0,
	"XAF": 0, "XOF": 0, "XPF": 0,
	"BHD": 3, "IQD": 3, "JOD": 3, "KWD": 3, "LYD": 3, "OMR": 3, "TND": 3,
}

// Exponent returns how many decimal places a currency has.
func Exponent(currency string) int {
	if e, ok := exponents[strings.ToUpper(currency)]; ok {
		return e
	}
	return 2
}

// Money is an amount in minor units, plus its currency.
//
// Stored by GORM as two columns via an embedded struct, so the amount stays an
// integer in the database and the currency is queryable. A single column
// holding "19.99 USD" is not something you can SUM.
type Money struct {
	Amount   int64  `json:"amount" gorm:"not null;default:0"`
	Currency string `json:"currency" gorm:"size:3;not null;default:'USD'"`
}

// New builds an amount from minor units.
func New(amount int64, currency string) Money {
	return Money{Amount: amount, Currency: strings.ToUpper(strings.TrimSpace(currency))}
}

// FromMajor converts a human-facing figure, 19.99, into minor units.
//
// Rounds half away from zero at the currency's own exponent. This is the one
// place a float is allowed near money, because a person typed it into a form,
// and it is converted immediately rather than carried.
func FromMajor(major float64, currency string) Money {
	scale := math.Pow(10, float64(Exponent(currency)))
	return New(int64(math.Round(major*scale)), currency)
}

// Major returns the amount as a decimal figure, for display only.
//
// Never feed this back into arithmetic: converting to float and back is how the
// cent you were protecting gets lost.
func (m Money) Major() float64 {
	return float64(m.Amount) / math.Pow(10, float64(Exponent(m.Currency)))
}

// String renders the amount with the right number of decimals.
func (m Money) String() string {
	e := Exponent(m.Currency)
	cur := m.Currency
	if cur == "" {
		cur = "USD"
	}
	return strconv.FormatFloat(m.Major(), 'f', e, 64) + " " + cur
}

func (m Money) IsZero() bool { return m.Amount == 0 }

// Add returns the sum, refusing to mix currencies.
func (m Money) Add(other Money) (Money, error) {
	if err := m.sameCurrency(other); err != nil {
		return Money{}, err
	}
	return New(m.Amount+other.Amount, m.currencyOrDefault()), nil
}

// Sub returns the difference, refusing to mix currencies.
func (m Money) Sub(other Money) (Money, error) {
	if err := m.sameCurrency(other); err != nil {
		return Money{}, err
	}
	return New(m.Amount-other.Amount, m.currencyOrDefault()), nil
}

// MulInt scales by a whole number, which is what a line total is: a unit price
// times a quantity. Exact, with no rounding decision to make.
func (m Money) MulInt(n int64) Money {
	return New(m.Amount*n, m.currencyOrDefault())
}

// MulFloat scales by a rate, for tax and percentage discounts.
//
// Rounds half away from zero, once, here. Percentage maths has to round
// somewhere, and doing it at one named point beats letting it happen wherever
// the value is eventually printed.
func (m Money) MulFloat(rate float64) Money {
	return New(int64(math.Round(float64(m.Amount)*rate)), m.currencyOrDefault())
}

// Convert changes currency at an exchange rate.
//
// The rate is major units of the target currency per major unit of this one,
// the way rates are quoted, and it is a decimal string ("1.0837", "148.2")
// rather than a float64. A float cannot hold most rates exactly, 1.005 is
// 1.00499999999999989 in binary, and an FX rate is multiplied by the largest
// amounts in the system, which is where that error stops being invisible. The
// string is parsed as an exact fraction, multiplied once and rounded once, half
// away from zero, the same rule MulFloat uses.
//
// The minor-unit exponents are handled here. USD has two decimals and JPY none,
// so 10.00 USD at 148.2 is 1482 JPY. Multiplying the amount by the rate and
// relabelling the currency gives 148200: a hundred times too much, silently.
func (m Money) Convert(to, rate string) (Money, error) {
	to = strings.ToUpper(strings.TrimSpace(to))
	if len(to) != 3 {
		return Money{}, ErrBadCurrency
	}
	r, ok := new(big.Rat).SetString(strings.TrimSpace(rate))
	if !ok || r.Sign() <= 0 {
		return Money{}, fmt.Errorf("exchange rate must be a positive decimal, got %q", rate)
	}

	v := new(big.Rat).SetInt64(m.Amount)
	v.Mul(v, r)
	shift := Exponent(to) - Exponent(m.currencyOrDefault())
	if shift != 0 {
		exp := shift
		if exp < 0 {
			exp = -exp
		}
		scale := new(big.Rat).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(exp)), nil))
		if shift > 0 {
			v.Mul(v, scale)
		} else {
			v.Quo(v, scale)
		}
	}

	n, err := roundHalfAway(v)
	if err != nil {
		return Money{}, err
	}
	return New(n, to), nil
}

// roundHalfAway rounds an exact fraction to the nearest whole number, with
// halves going away from zero.
func roundHalfAway(v *big.Rat) (int64, error) {
	num := new(big.Int).Set(v.Num())
	neg := num.Sign() < 0
	num.Abs(num)
	q, rem := new(big.Int).QuoRem(num, v.Denom(), new(big.Int))
	// rem/den is at least a half exactly when 2*rem >= den.
	if new(big.Int).Lsh(rem, 1).Cmp(v.Denom()) >= 0 {
		q.Add(q, big.NewInt(1))
	}
	if neg {
		q.Neg(q)
	}
	if !q.IsInt64() {
		return 0, fmt.Errorf("converted amount does not fit in minor units")
	}
	return q.Int64(), nil
}

// Allocate splits an amount into n parts without losing or inventing a minor
// unit.
//
// Splitting 10.00 three ways gives 3.34, 3.33, 3.33, not three lots of 3.33
// with a cent evaporated. The remainder goes to the earliest parts, which is
// the convention every accounting system expects.
func (m Money) Allocate(n int) []Money {
	if n <= 0 {
		return nil
	}
	base := m.Amount / int64(n)
	rem := m.Amount % int64(n)
	out := make([]Money, n)
	for i := 0; i < n; i++ {
		amt := base
		if int64(i) < rem {
			amt++
		}
		out[i] = New(amt, m.currencyOrDefault())
	}
	return out
}

func (m Money) currencyOrDefault() string {
	if m.Currency == "" {
		return "USD"
	}
	return m.Currency
}

func (m Money) sameCurrency(other Money) error {
	a, b := m.currencyOrDefault(), other.currencyOrDefault()
	if a != b {
		return fmt.Errorf("%w: %s and %s", ErrCurrencyMismatch, a, b)
	}
	return nil
}

// Validate checks the currency looks like ISO 4217.
func (m Money) Validate() error {
	c := m.Currency
	if len(c) != 3 {
		return ErrBadCurrency
	}
	for _, r := range c {
		if r < 'A' || r > 'Z' {
			return ErrBadCurrency
		}
	}
	return nil
}

// UnmarshalJSON accepts the object form, and a bare number.
//
// A bare number is read as MAJOR units: 19.99 means $19.99, and 2500 means
// $2,500.00, not $25.00. There is no way to tell those two intents apart from
// the wire, so read this before pointing a client at it -- if your caller
// already speaks in cents, as Stripe's API does, sending the bare number will
// silently multiply every price by a hundred. Send the object form and the
// question does not arise.
//
// Major units win the coin toss because the bare-number path exists for
// hand-written and legacy callers, and a human writing a price by hand writes
// 19.99. Output is always the object form: a bare number on the wire, with the
// currency left to be inferred, is the exact thing this type exists to stop.
func (m *Money) UnmarshalJSON(b []byte) error {
	trimmed := strings.TrimSpace(string(b))
	if trimmed == "null" || trimmed == "" {
		return nil
	}
	if trimmed[0] == '{' {
		var alias struct {
			Amount   int64  `json:"amount"`
			Currency string `json:"currency"`
		}
		if err := json.Unmarshal(b, &alias); err != nil {
			return err
		}
		m.Amount = alias.Amount
		if alias.Currency != "" {
			m.Currency = strings.ToUpper(alias.Currency)
		}
		return nil
	}
	major, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return fmt.Errorf("reading money: expected {amount, currency} or a number, got %s", trimmed)
	}
	*m = FromMajor(major, m.currencyOrDefault())
	return nil
}

// Deliberately no driver.Valuer or sql.Scanner.
//
// Implementing them is the obvious thing and it silently breaks the type.
// GORM treats anything with Value/Scan as a scalar, so it ignores the
// embedded tag and collapses the field to a single INTEGER column: the amount
// survives and the currency is thrown away, which is the exact failure this
// package exists to prevent. Verified by looking at the created table, not by
// reading the tag.
//
// Without them, GORM expands the struct into <field>_amount and
// <field>_currency, which is what you want anyway: the amount stays an integer
// the database can SUM and the currency is queryable.
