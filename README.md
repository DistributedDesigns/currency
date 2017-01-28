Currency
=====

[![Build Status](https://travis-ci.org/DistributedDesigns/currency.svg?branch=master)](https://travis-ci.org/DistributedDesigns/currency)

Nice currency calculations.

#### Examples
```go
// Currency looks like
type Currency struct { cents uint }

// You can make a new Currency by declaring the number of cents directly
c1 := Currency{1234} // $123.40

// You can also parse strings and floats
c2, _ := NewFromString("500") // $500.00
c3, _ := NewFromFloat(2.4855) // $2.49

// Conversions that don't work return an error
_, err := NewFromString("abcd") // err != nil

// Basic operations are available
payment := NewFromString("10")
balance := NewFromString("100")

taxes := payment.Mul(0.07) // $0.70
balance = balance.Add(payment.Sub(taxes)) // $109.30

// There are convenience mod and remainder functions
stockPrice := NewFromFloat(3.33)
account := NewFromFloat(10.00)
times, remainder := stockPrice.FitsInto(account) // times: 3, remainder: $0.01
```
