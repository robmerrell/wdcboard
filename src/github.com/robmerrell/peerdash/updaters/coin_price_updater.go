package updaters

import (
	"encoding/json"
	"fmt"
	"github.com/robmerrell/peerdash/models"
	"io/ioutil"
	"net/http"
)

var tickerUrl = "https://btc-e.com/api/2/ppc_%s/ticker"

type CoinPrice struct{}

// Update retrieves PPC buy prices in both USD and BTC and saves
// the prices to the database.
func (c *CoinPrice) Update() error {
	usd, err := getQuoteCurrencyValue("usd")
	if err != nil {
		return err
	}

	btc, err := getQuoteCurrencyValue("btc")
	if err != nil {
		return err
	}

	conn := models.CloneConnection()
	defer conn.Close()
	return models.InsertPrice(conn, usd, btc)
}

// getQuoteCurrencyValue gets the current buy price for a given
// quote currency. At the moment BTC-e only supports USD and BTC
// for PPC trades.
func getQuoteCurrencyValue(quoteCurrency string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf(tickerUrl, quoteCurrency))
	if err != nil {
		return 0.0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var value struct {
		Ticker struct {
			Buy float64 `json:"buy"`
		}
	}
	if err := json.Unmarshal(body, &value); err != nil {
		return 0.0, err
	}

	return value.Ticker.Buy, nil
}
