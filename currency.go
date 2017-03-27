// Package currency implements basic math operations that avoids
// rounding problems common to floating point arithmetic. This is
// accomplished by only storing cents
//
// To use it as a struct:
//  type Foo struct {
//    Balance Currency
//  }
//
// Its zero value is $0.00
//
package currency

import (
	"errors"
	"fmt"
	"strconv"
)

// Currency : is a dollar and cent amount
type Currency struct{ cents uint }

// NewFromFloat : Parses a float into a new Currency.
func NewFromFloat(f float64) (Currency, error) {
	// Currencies are strictly non negative
	if f < 0 {
		return Currency{}, errors.New("Currency must be positive")
	}

	// Shifting by 0.5 will make the uint cast do rounding
	cents := uint((f * 100) + 0.5)

	return Currency{cents}, nil
}

// NewFromString : Parses a string into a new Currency.
func NewFromString(s string) (Currency, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return Currency{}, err
	}

	return NewFromFloat(f)
}

func (c Currency) String() string {
	dollars := c.cents / 100
	cents := c.cents % 100
	return fmt.Sprintf("$%d.%02d", dollars, cents)
}

// ToFloat : Converts $1.23 -> 1.23
func (c Currency) ToFloat() float64 {
	return float64(c.cents) / 100
}

// FitsInto : Finds whole number of divisions of two currencies, with maximum whole
// multiple. If either argument is zero then (0, $0.00) is returned
func (c *Currency) FitsInto(total Currency) (uint, Currency) {
	// Return default values if c is empty or total is smaller than divisor
	if c.cents == 0 || c.cents > total.cents {
		return 0, Currency{}
	}

	times := total.cents / c.cents
	wholeMultiple := total.cents - (total.cents % c.cents)

	return times, Currency{wholeMultiple}
}

// Add : Increase the value of a Currency
func (c *Currency) Add(c2 Currency) {
	c.cents += c2.cents
}

// Sub : Decrease the value of a currency
func (c *Currency) Sub(c2 Currency) error {
	if c2.cents > c.cents {
		return errors.New("Cannot create negative currency")
	}

	c.cents -= c2.cents

	return nil
}

// Mul : Scales the value of a currency
func (c *Currency) Mul(f float64) error {
	if f < 0 {
		return errors.New("Cannot multiply by negative numbers")
	}

	scaledCents := float64(c.cents) * f

	// Do the shift / cast rounding trick again
	c.cents = uint(scaledCents + 0.5)

	return nil
}
