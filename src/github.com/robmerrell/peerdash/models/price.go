package models

import (
	"time"
)

type Price struct {
	Usd       float64   "usd"
	Btc       float64   "btc"
	CreatedAt time.Time "createdAt"
}

var priceCollection = "prices"

// InsertPrice create new PPC price points
func InsertPrice(conn *MgoConnection, usd, btc float64) error {
	price := &Price{
		Usd:       usd,
		Btc:       btc,
		CreatedAt: time.Now().UTC(),
	}

	err := conn.DB.C(priceCollection).Insert(price)
	return err
}
