package cmds

import (
	"github.com/codegangsta/martini"
	"github.com/hoisie/mustache"
	"github.com/robmerrell/wdcboard/lib"
	"github.com/robmerrell/wdcboard/models"
	"log"
	"net/http"
	"strconv"
)

var ServerDoc = `
Starts the WDCBoard webserver.
`

func webError(err error, res http.ResponseWriter) {
	log.Println(err)
	http.Error(res, "There was an error, try again later", 500)
}

func ServeAction() error {
	m := martini.Classic()
	m.Use(martini.Static("resources/public"))

	mainView, err := mustache.ParseFile("resources/views/main.html.mustache")
	if err != nil {
		panic(err)
	}

	m.Get("/", func(res http.ResponseWriter) string {
		conn := models.CloneConnection()
		defer conn.Close()

		// get the latest pricing data
		price, err := models.GetLatestPrice(conn)
		if err != nil {
			webError(err, res)
			return ""
		}

		// get data for the graph

		// get the forum posts

		// get reddit posts

		// get the mining information
		network, err := models.GetLatestNetworkSnapshot(conn)
		if err != nil {
			webError(err, res)
			return ""
		}

		// generate the HTML

		return mainView.Render(generateTplVars(price, network))
	})

	// returns basic information about the state of the service. If any hardcoded checks fail
	// the message is returned with a 500 status. We can then use pingdom or another service
	// to alert when data integrity may be off.
	m.Get("/health", func() string {
		return "ok"
	})

	m.Run()

	return nil
}

func generateTplVars(price *models.Price, network *models.Network) map[string]string {
	// apply the necessary style for the percent change box
	changeStyle := "percent-change-stat-up"
	if price.Cryptsy.PercentChange != "" && string(price.Cryptsy.PercentChange[0]) == "-" {
		changeStyle = "percent-change-stat-down"
	}

	// marketcap
	minedNum, _ := strconv.Atoi(network.Mined)
	marketCap := float64(minedNum) * price.Cryptsy.Usd

	// coins left to be mined
	remainingCoins := 265420800 - minedNum

	vars := map[string]string{
		"usd":         lib.RenderFloat("", price.Cryptsy.Usd),
		"btc":         strconv.FormatFloat(price.Cryptsy.Btc, 'f', 8, 64),
		"marketCap":   lib.RenderInteger("", int(marketCap)),
		"change":      price.Cryptsy.PercentChange,
		"changeStyle": changeStyle,

		"hashRate":   lib.RenderFloatFromString("", network.HashRate),
		"difficulty": lib.RenderFloatFromString("", network.Difficulty),
		"mined":      lib.RenderIntegerFromString("", network.Mined),
		"remaining":  lib.RenderInteger("", remainingCoins),
	}

	return vars
}
