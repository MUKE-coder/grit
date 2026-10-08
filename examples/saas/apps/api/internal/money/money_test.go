package money

import (
	"encoding/json"
	"errors"
	"testing"
)

// USD has two decimals, JPY none and KWD three. Relabelling a USD amount as yen
// is off by a factor of a hundred; Convert moves the decimal point for you.
func TestConvertHandlesMinorUnitExponents(t *testing.T) {
	cases := []struct {
		from       Money
		to, rate   string
		wantAmount int64
	}{
		{New(1000, "USD"), "JPY", "148.2", 1482},     // 10.00 USD -> 1482 JPY
		{New(1482, "JPY"), "USD", "0.0067476", 1000}, // 1482 JPY -> 10.00 USD
		{New(1000, "USD"), "KWD", "0.3075", 3075},    // 10.00 USD -> 3.075 KWD
		{New(1000, "USD"), "EUR", "0.9227", 923},     // 10.00 USD -> 9.23 EUR
	}
	for _, c := range cases {
		got, err := c.from.Convert(c.to, c.rate)
		if err != nil {
			t.Fatalf("%s -> %s: %v", c.from, c.to, err)
		}
		if got.Amount != c.wantAmount || got.Currency != c.to {
			t.Errorf("%s at %s -> %s: got %d %s, want %d", c.from, c.rate, c.to, got.Amount, got.Currency, c.wantAmount)
		}
	}
}

// 1.005 is not representable in binary, so a float rate turns 100.5 cents into
// 100.49999999999999 and rounds it down. The rate is a string for this reason.
func TestConvertRoundsTheExactValue(t *testing.T) {
	got, err := New(100, "USD").Convert("EUR", "1.005")
	if err != nil {
		t.Fatal(err)
	}
	if got.Amount != 101 {
		t.Errorf("1.00 USD at 1.005 is 100.5 cents, which rounds to 101: got %d", got.Amount)
	}
	neg, err := New(-1, "USD").Convert("EUR", "2.5")
	if err != nil {
		t.Fatal(err)
	}
	if neg.Amount != -3 {
		t.Errorf("halves round away from zero: -2.5 should be -3, got %d", neg.Amount)
	}
}

func TestConvertRefusesNonsense(t *testing.T) {
	for _, rate := range []string{"", "abc", "0", "-1.2"} {
		if _, err := New(100, "USD").Convert("EUR", rate); err == nil {
			t.Errorf("rate %q was accepted", rate)
		}
	}
	if _, err := New(100, "USD").Convert("EU", "1.1"); !errors.Is(err, ErrBadCurrency) {
		t.Errorf("a two-letter currency was accepted: %v", err)
	}
}

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
	if string(out) != `{"amount":1999,"currency":"USD"}` {
		t.Errorf("unexpected wire format: %s", out)
	}

	var m Money
	if err := json.Unmarshal([]byte(`{"amount":2500,"currency":"UGX"}`), &m); err != nil {
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
