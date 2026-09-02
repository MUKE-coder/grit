package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeMoneyFiles writes internal/money: the type commerce code needs and
// float64 cannot be.
func writeMoneyFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	files := map[string]string{
		filepath.Join(apiRoot, "internal", "money", "money.go"):      moneyGo(),
		filepath.Join(apiRoot, "internal", "money", "money_test.go"): moneyTestGo(),
	}
	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func moneyGo() string {
	return `// Package money is an amount and the currency it is denominated in.
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
	Amount   int64  ` + "`" + `json:"amount" gorm:"not null;default:0"` + "`" + `
	Currency string ` + "`" + `json:"currency" gorm:"size:3;not null;default:'USD'"` + "`" + `
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
			Amount   int64  ` + "`" + `json:"amount"` + "`" + `
			Currency string ` + "`" + `json:"currency"` + "`" + `
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
`
}

func moneyTestGo() string {
	return `package money

import (
	"encoding/json"
	"testing"
)

// The reason the package exists. In float64 this sum is 0.30000000000000004,
// and after enough of them the ledger stops balancing.
func TestAdditionIsExact(t *testing.T) {
	a, b := FromMajor(0.1, "USD"), FromMajor(0.2, "USD")
	sum, err := a.Add(b)
	if err != nil {
		t.Fatal(err)
	}
	if sum.Amount != 30 {
		t.Errorf("expected 30 cents, got %d", sum.Amount)
	}
	if sum.String() != "0.30 USD" {
		t.Errorf("expected 0.30 USD, got %s", sum.String())
	}
}

// Zero-decimal currencies are the ones a hardcoded /100 destroys.
func TestZeroDecimalCurrency(t *testing.T) {
	m := New(50000, "UGX")
	if m.Major() != 50000 {
		t.Errorf("UGX has no minor unit: expected 50000, got %v", m.Major())
	}
	if m.String() != "50000 UGX" {
		t.Errorf("expected 50000 UGX, got %s", m.String())
	}
	if got := FromMajor(50000, "UGX"); got.Amount != 50000 {
		t.Errorf("round trip lost the amount: %d", got.Amount)
	}
}

func TestThreeDecimalCurrency(t *testing.T) {
	m := FromMajor(1.5, "KWD")
	if m.Amount != 1500 {
		t.Errorf("KWD has three decimals: expected 1500 fils, got %d", m.Amount)
	}
}

// Mixing currencies must be refused, not silently added.
func TestCurrencyMismatchIsRefused(t *testing.T) {
	usd, ugx := New(100, "USD"), New(100, "UGX")
	if _, err := usd.Add(ugx); err == nil {
		t.Error("adding USD to UGX must fail")
	}
	if _, err := usd.Sub(ugx); err == nil {
		t.Error("subtracting UGX from USD must fail")
	}
}

// A line total is a unit price times a quantity: exact, no rounding.
func TestMulIntIsExact(t *testing.T) {
	if got := New(1999, "USD").MulInt(3); got.Amount != 5997 {
		t.Errorf("expected 5997, got %d", got.Amount)
	}
}

func TestMulFloatRoundsOnce(t *testing.T) {
	// 19.99 at 7.5% tax is 1.49925, which must land on 150 rather than 149.
	if got := New(1999, "USD").MulFloat(0.075); got.Amount != 150 {
		t.Errorf("expected 150, got %d", got.Amount)
	}
}

// Splitting must not lose or invent a minor unit.
func TestAllocateKeepsEveryUnit(t *testing.T) {
	parts := New(1000, "USD").Allocate(3)
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}
	var total int64
	for _, p := range parts {
		total += p.Amount
	}
	if total != 1000 {
		t.Errorf("allocation lost money: %d of 1000", total)
	}
	if parts[0].Amount != 334 || parts[1].Amount != 333 || parts[2].Amount != 333 {
		t.Errorf("remainder should go to the earliest part: %v", parts)
	}
}

// The wire format carries the currency. A bare number is what this type exists
// to stop, so it is accepted on input for older clients but never emitted.
func TestJSONShape(t *testing.T) {
	out, err := json.Marshal(New(1999, "USD"))
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != ` + "`" + `{"amount":1999,"currency":"USD"}` + "`" + ` {
		t.Errorf("unexpected wire format: %s", out)
	}

	var m Money
	if err := json.Unmarshal([]byte(` + "`" + `{"amount":2500,"currency":"UGX"}` + "`" + `), &m); err != nil {
		t.Fatal(err)
	}
	if m.Amount != 2500 || m.Currency != "UGX" {
		t.Errorf("round trip changed the value: %+v", m)
	}
}

func TestBareNumberIsAcceptedAsMajorUnits(t *testing.T) {
	m := Money{Currency: "USD"}
	if err := json.Unmarshal([]byte("19.99"), &m); err != nil {
		t.Fatal(err)
	}
	if m.Amount != 1999 {
		t.Errorf("a bare 19.99 should read as 1999 cents, got %d", m.Amount)
	}
}

func TestValidateRejectsNonISOCodes(t *testing.T) {
	for _, bad := range []string{"", "US", "usd1", "DOLLAR"} {
		if err := (Money{Currency: bad}).Validate(); err == nil {
			t.Errorf("%q should not validate as a currency", bad)
		}
	}
	if err := (Money{Currency: "USD"}).Validate(); err != nil {
		t.Errorf("USD should validate: %v", err)
	}
}
`
}
