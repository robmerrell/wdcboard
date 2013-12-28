package models

import (
	"github.com/robmerrell/wdcboard/config"
	"labix.org/v2/mgo/bson"
	. "launchpad.net/gocheck"
	"testing"
	"time"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type priceSuite struct{}

var _ = Suite(&priceSuite{})

func (s *priceSuite) SetUpTest(c *C) {
	config.LoadConfig("test")
	ConnectToDB(config.String("database.host"), config.String("database.db"))
	DropCollections()
}

func (s *priceSuite) TestInserting(c *C) {
	conn := CloneConnection()
	defer conn.Close()

	p := &Price{
		UsdPerBtc: 100.0,
		Cryptsy: &ExchangePrice{
			Btc: 0.3456,
		},
	}
	p.Insert(conn)

	var newp Price
	conn.DB.C(priceCollection).Find(bson.M{"usdperbtc": 100.0}).One(&newp)

	c.Check(newp.Cryptsy.Btc, Equals, 0.3456)
}

func (s *priceSuite) TestSettingPercentChange(c *C) {
	conn := CloneConnection()
	defer conn.Close()

	yesterdayDuration, _ := time.ParseDuration("-25h")

	p1 := &Price{
		UsdPerBtc: 100.0,
		Cryptsy: &ExchangePrice{
			Btc: 0.3456,
		},
		GeneratedAt: time.Now().UTC().Add(yesterdayDuration),
	}
	p1.Insert(conn)

	p2 := &Price{
		UsdPerBtc: 100.0,
		Cryptsy: &ExchangePrice{
			Btc: 0.55,
		},
		GeneratedAt: time.Now().UTC(),
	}
	p2.SetPercentChange(conn)

	c.Check(p2.Cryptsy.PercentChange, Equals, "59.14")
}

func (s *priceSuite) TestPercentChange(c *C) {
	c.Check(percentChange(1, 2), Equals, "100.00")
	c.Check(percentChange(1, 5), Equals, "400.00")
	c.Check(percentChange(3, 1.63), Equals, "-45.67")
	c.Check(percentChange(0.456, 0.457), Equals, "0.22")
}
