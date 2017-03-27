Currency
=====

[![Build Status](https://travis-ci.org/DistributedDesigns/currency.svg?branch=master)](https://travis-ci.org/DistributedDesigns/currency)

Nice currency calculations.

#### Examples
```go
// Currency looks like
type Currency struct { cents uint }

// You can make new Currencies from strings and floats with normal rounding
c2, _ := NewFromString("500") // $500.00
c3, _ := NewFromFloat(2.4855) // $2.49

// Conversions that don't work return an error
_, err := NewFromString("abcd") // err != nil
_, err := NewFromFloat(-50)     // err != nil

// Basic operations are available
payment, _ := NewFromString("10.0")
bill, _ := NewFromFloat(25)
balance, _ := NewFromString("100")
balance.Add(payment) // balance: $110.00
balance.Sub(bill)    // balance:  $85.00
balance.Mul(0.1)     // balance:   $8.50

// It's easy to determine the number of times one currency fits into another
stockPrice := NewFromFloat(3.33)
balance := NewFromFloat(10.00)
times, maxStockPurchase := stockPrice.FitsInto(balance) // times: 3, maxStockPurchase: $9.99
balance.Sub(maxStockPurchase) // balance: $0.01
```
