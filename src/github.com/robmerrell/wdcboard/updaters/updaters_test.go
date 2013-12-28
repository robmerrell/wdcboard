package updaters

import (
	"fmt"
	"github.com/robmerrell/wdcboard/config"
	"github.com/robmerrell/wdcboard/models"
	"labix.org/v2/mgo/bson"
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
	coinBaseServer *httptest.Server
	cryptsyServer  *httptest.Server

	usdServer *httptest.Server
	badServer *httptest.Server
}

var _ = Suite(&coinPriceSuite{})

func (s *coinPriceSuite) SetUpSuite(c *C) {
	s.coinBaseServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := "{\"btc_to_usd\":\"676.58046\"}"
		fmt.Fprintln(w, res)
	}))

	s.cryptsyServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := "{\"success\": 1, \"return\": {\"markets\": {\"WDC\": {\"recenttrades\": [{\"id\": \"9496223\", \"time\": \"2013-12-25 16:27:42\", \"price\": \"0.00053275\", \"quantity\": \"32.67455654\", \"total\": \"0.01740737\"}]}}}}"
		fmt.Fprintln(w, res)
	}))

	s.badServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "An error occured, I'm not returning valid JSON")
	}))
}

func (s *coinPriceSuite) SetUpTest(c *C) {
	config.LoadConfig("test")
	models.ConnectToDB(config.String("database.host"), config.String("database.db"))
	models.DropCollections()
}

func (s *coinPriceSuite) TearDownSuite(c *C) {
	s.coinBaseServer.Close()
	s.cryptsyServer.Close()
	s.badServer.Close()
}

func replaceQuoteUrl(newUrl string, val *string, testFunc func()) {
	oldUrl := *val
	*val = newUrl

	testFunc()

	val = &oldUrl
}

func (s *coinPriceSuite) TestTradePrices(c *C) {
	replaceQuoteUrl(s.coinBaseServer.URL, &coinbaseUrl, func() {
		value, _ := coinbaseQuote()
		c.Check(value, Equals, 676.58046)
	})

	replaceQuoteUrl(s.badServer.URL, &coinbaseUrl, func() {
		value, err := coinbaseQuote()
		c.Check(value, Equals, 0.0)
		c.Assert(err, NotNil)
	})

	replaceQuoteUrl(s.cryptsyServer.URL, &cryptsyUrl, func() {
		value, _ := cryptsyQuote()
		c.Check(value, Equals, 0.00053275)
	})

	replaceQuoteUrl(s.badServer.URL, &cryptsyUrl, func() {
		value, err := cryptsyQuote()
		c.Check(value, Equals, 0.0)
		c.Assert(err, NotNil)
	})
}

func (s *coinPriceSuite) TestSavingPrices(c *C) {
	replaceQuoteUrl(s.coinBaseServer.URL, &coinbaseUrl, func() {
		replaceQuoteUrl(s.cryptsyServer.URL, &cryptsyUrl, func() {
			conn := models.CloneConnection()
			defer conn.Close()

			coinPrice := &CoinPrice{}
			coinPrice.Update()

			var saved models.Price
			conn.DB.C("prices").Find(bson.M{}).One(&saved)

			c.Check(saved.UsdPerBtc, Equals, 676.58046)
			c.Check(saved.Cryptsy.Btc, Equals, 0.00053275)
		})
	})
}
