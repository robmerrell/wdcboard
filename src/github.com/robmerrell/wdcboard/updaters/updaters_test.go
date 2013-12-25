package updaters

import (
	"fmt"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

// --------------------------------
// Tests for retrieving coin prices
// --------------------------------
type coinPriceSuite struct {
	usdServer *httptest.Server
	btcServer *httptest.Server
	badServer *httptest.Server
}

var _ = Suite(&coinPriceSuite{})

func (s *coinPriceSuite) SetUpSuite(c *C) {
	s.usdServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		usd := "{\"ticker\":{\"high\":4.37,\"low\":4.01,\"avg\":4.19,\"vol\":283714.84819,\"vol_cur\":67726.42618,\"last\":4.21,\"buy\":4.213,\"sell\":4.21,\"updated\":1387037358,\"server_time\":1387037358}}"
		fmt.Fprintln(w, usd)
	}))

	s.btcServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		btc := "{\"ticker\":{\"high\":0.00507,\"low\":0.00477,\"avg\":0.00492,\"vol\":380.68146,\"vol_cur\":77679.12094,\"last\":0.00492,\"buy\":0.00493,\"sell\":0.00492,\"updated\":1387037763,\"server_time\":1387037763}}"
		fmt.Fprintln(w, btc)
	}))

	s.badServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "An error occured, I'm not returning valid JSON")
	}))
}

func (s *coinPriceSuite) TearDownSuite(c *C) {
	s.usdServer.Close()
	s.btcServer.Close()
	s.badServer.Close()
}

func replaceCoinPriceUrl(newUrl string, testFunc func()) {
	oldTicker := tickerUrl
	tickerUrl = newUrl + "/%s"

	testFunc()

	tickerUrl = oldTicker
}

func (s *coinPriceSuite) TestTradePrices(c *C) {
	replaceCoinPriceUrl(s.usdServer.URL, func() {
		value, _ := getQuoteCurrencyValue("usd")
		c.Check(value, Equals, 4.213)
	})

	replaceCoinPriceUrl(s.btcServer.URL, func() {
		value, _ := getQuoteCurrencyValue("btc")
		c.Check(value, Equals, 0.00493)
	})

	replaceCoinPriceUrl(s.badServer.URL, func() {
		value, err := getQuoteCurrencyValue("usd")
		c.Check(value, Equals, 0.0)
		c.Assert(err, NotNil)
	})
}

func (s *coinPriceSuite) TestSavingPrices(c *C) {

}
