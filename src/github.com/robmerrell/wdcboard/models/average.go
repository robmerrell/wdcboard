package models

import (
	"time"
)

type ExchangeAverage struct {
	Btc float64 "btc"
	Usd float64 "usd"
}

type Average struct {
	Cryptsy   *ExchangeAverage "cryptsy"
	TimeBlock time.Time        "timeBlock"
}

var averageCollection = "averages"

func GenerateAverage(conn *MgoConnection, startTime, endTime time.Time) (*Average, error) {
	prices, err := GetPricesBetweenDates(conn, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// find the average of the prices
	avg := &Average{TimeBlock: startTime, Cryptsy: &ExchangeAverage{}}
	btcSum := 0.0
	usdSum := 0.0
	for _, price := range prices {
		usdSum += price.Cryptsy.Usd
		btcSum += price.Cryptsy.Btc
	}
	avg.Cryptsy.Btc = btcSum / float64(len(prices))
	avg.Cryptsy.Usd = usdSum / float64(len(prices))

	err = avg.Insert(conn)
	return avg, err
}

// Insert saves a new WDC average price point to the database.
func (a *Average) Insert(conn *MgoConnection) error {
	return conn.DB.C(averageCollection).Insert(a)
}
