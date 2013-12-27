package models

import (
	"github.com/robmerrell/wdcboard/config"
	. "launchpad.net/gocheck"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type priceSuite struct{}

var _ = Suite(&priceSuite{})

func (s *priceSuite) SetUpTest(c *C) {
	config.LoadConfig("test")
	ConnectToDB(config.String("database.host"), config.String("database.db"))
	DropCollection(priceCollection)
}

func (s *priceSuite) TestPercentChange(c *C) {
	c.Check(percentChange(1, 2), Equals, "100.00")
	c.Check(percentChange(1, 5), Equals, "400.00")
	c.Check(percentChange(3, 1.63), Equals, "-45.67")
	c.Check(percentChange(0.456, 0.457), Equals, "0.22")
}

func (s *priceSuite) TestInserting(c *C) {

}
