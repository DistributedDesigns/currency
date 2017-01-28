package currency

import "testing"

type conversionPair struct {
	stringCase string
	floatCase  float64
}

var testTable = map[conversionPair]string{
	conversionPair{"100", 100}:       "$100.00",
	conversionPair{"1.01", 1.01}:     "$1.01",
	conversionPair{"1.0001", 1.001}:  "$1.00",
	conversionPair{"0.995", 0.995}:   "$1.00",
	conversionPair{"0", 0}:           "$0.00",
	conversionPair{"34.567", 34.567}: "$34.57",
}

func TestNewFromFloat(t *testing.T) {
	for tc, s := range testTable {
		actual, err := NewFromFloat(tc.floatCase)
		if err != nil {
			t.Errorf("Error parsing %f", tc.floatCase)
		} else if actual.String() != s {
			t.Errorf("Expected %s got %s", s, actual)
		}
	}

	if _, err := NewFromFloat(-50); err == nil {
		t.Error("Should not parse negative floats")
	}
}

func TestNewFromString(t *testing.T) {
	for tc, s := range testTable {
		a, err := NewFromString(tc.stringCase)
		if err != nil {
			t.Errorf("Error parsing %s", tc.stringCase)
		} else if a.String() != s {
			t.Errorf("Expected `%s` got `%s`", s, a)
		}
	}

	unparseableStrings := []string{
		"abcd",
		"-1",
		"",
	}
	for _, s := range unparseableStrings {
		if _, err := NewFromString(s); err == nil {
			t.Errorf("%s should not be parsable", s)
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := map[Currency]float64{
		Currency{45}:   0.45,
		Currency{1000}: 10,
		Currency{101}:  1.01,
		Currency{10}:   0.1,
	}

	for c, expected := range tests {
		actual := c.ToFloat()
		if actual != expected {
			t.Errorf("Expected %f got %f", expected, actual)
		}
	}
}

type testPair struct {
	c1, c2 Currency
}
type expectedPair struct {
	times     uint
	remainder Currency
}

var (
	zeroD    = Currency{0}
	oneD     = Currency{100}
	tenD     = Currency{1000}
	thirtyD  = Currency{3000}
	hundredD = Currency{10000}
)

func TestFitsInto(t *testing.T) {
	tests := map[testPair]expectedPair{
		// $30 fits into $100: 3 times with $10 remainder
		testPair{thirtyD, hundredD}: expectedPair{3, tenD},

		// $100 fits into $30: 0 times with $30 remainder
		testPair{hundredD, thirtyD}: expectedPair{0, thirtyD},

		// $0 fits into $10: 0 times with $0 remainder
		testPair{zeroD, tenD}: expectedPair{0, zeroD},

		// $10 fits into $0: 0 times with $10 remainder
		testPair{tenD, zeroD}: expectedPair{0, tenD},

		// $3.33 fits into $10: 3 times with $0.01 remainder
		testPair{Currency{333}, tenD}: expectedPair{3, Currency{1}},
	}

	for tp, ep := range tests {
		times, remainder := tp.c1.FitsInto(tp.c2)
		if times != ep.times {
			t.Errorf("Expected %d got %d times", ep.times, times)
		} else if remainder != ep.remainder {
			t.Errorf("Expected %s got %s remainder", ep.remainder, remainder)
		}
	}
}

func TestAdd(t *testing.T) {
	tests := map[testPair]Currency{
		// $1.50 + $2.75 = $4.25
		testPair{Currency{150}, Currency{275}}: Currency{425},

		// $0.00 + $1.00 = $1.00
		testPair{zeroD, oneD}: oneD,

		// $1.00 + $0.00 = $1.00
		testPair{oneD, zeroD}: oneD,
	}

	for tp, e := range tests {
		c3 := tp.c1.Add(tp.c2)
		if c3.String() != e.String() {
			t.Errorf("Expected %s got %s", e.String(), c3.String())
		}
	}
}

func TestSub(t *testing.T) {
	tests := map[testPair]Currency{
		// $4.25 - $1.50 = $2.75
		testPair{Currency{425}, Currency{150}}: Currency{275},

		// $1.00 - $0.00 = $1.00
		testPair{oneD, zeroD}: oneD,

		// $1.00 - $1.00 = $0.00
		testPair{oneD, oneD}: zeroD,
	}

	for tp, e := range tests {
		c3, _ := tp.c1.Sub(tp.c2)
		if c3.String() != e.String() {
			t.Errorf("Expected %s got %s", e.String(), c3.String())
		}
	}

	// Make sure we can't get a negative balance
	// $1.00 - $10.00 = << error >>
	_, err := oneD.Sub(tenD)
	if err == nil {
		t.Error("Should not subtract more than balance")
	}
}

type mulPair struct {
	c Currency
	f float64
}

func TestMul(t *testing.T) {
	tests := map[mulPair]Currency{
		// $1.25 * 5 = $6.25
		mulPair{Currency{125}, 5}: Currency{625},

		// $1.00 * 0 = $0.00
		mulPair{oneD, 0}: zeroD,

		// $1.00 * 3.456 = $3.46
		mulPair{oneD, 3.456}: Currency{346},

		// $10.00 * 0.1 = $1.00
		mulPair{tenD, 0.1}: oneD,

		// $10.00 * 0.499999 = $5.00
		mulPair{tenD, 0.499999}: Currency{500},

		// $10.00 * (2/3) = $6.67
		mulPair{tenD, 2.0 / 3.0}: Currency{667},
	}

	for mp, e := range tests {
		c3, _ := mp.c.Mul(mp.f)
		if c3.String() != e.String() {
			t.Errorf("Expected %s got %s", e.String(), c3.String())
		}
	}

	// Negative multiplication should not be allowed
	_, err := oneD.Mul(-1)
	if err == nil {
		t.Error("Should not multiply by negative numbers")
	}
}
