package routes

import (
	"cryptoapi/core/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CoinPriceResponse struct {
	Coin                   string `json:"coin"`
	Logo                   string `json:"logo"`
	Ticker                 string `json:"ticker"`
	MinimumTransaction     int    `json:"minimum_transaction"`
	MinimumTransactionCoin string `json:"minimum_transaction_coin"`
	MinimumFee             int    `json:"minimum_fee"`
	MinimumFeeCoin         string `json:"minimum_fee_coin"`
	FeePercent             string `json:"fee_percent"`
	Status                 string `json:"status"`
	Prices                 struct {
		USD string `json:"USD"`
		EUR string `json:"EUR"`
		GBP string `json:"GBP"`
		CAD string `json:"CAD"`
		JPY string `json:"JPY"`
		AED string `json:"AED"`
		DKK string `json:"DKK"`
		BRL string `json:"BRL"`
		CNY string `json:"CNY"`
		HKD string `json:"HKD"`
		INR string `json:"INR"`
		MXN string `json:"MXN"`
		UGX string `json:"UGX"`
		PLN string `json:"PLN"`
		PHP string `json:"PHP"`
		CZK string `json:"CZK"`
		HUF string `json:"HUF"`
		BGN string `json:"BGN"`
		RON string `json:"RON"`
		LKR string `json:"LKR"`
		TRY string `json:"TRY"`
	} `json:"prices"`
	PricesUpdated time.Time `json:"prices_updated"`
}

func CoinPriceAPIREQ(c *gin.Context) {
	//* recieve user queries for coin and usd amount
	coin := strings.ToLower(c.DefaultQuery("coin", ""))
	usd_amount := c.DefaultQuery("usd_amount", "")
	//* convert the provided usd amount from string to integer
	int_usd_amount, err := strconv.ParseFloat(usd_amount, 64)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to convert usd_amount to FLOAT",
		})
		return
	}
	//* set the url to make a request to
	url := ("https://api.cryptapi.io/" + coin + "/info/")
	spaceClient := http.Client{
		Timeout: time.Second * 10,
	}
	//* setup the http request to the url above
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to setup request to COIN_PRICE",
		})
		return
	}
	//* set a custom user agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36")
	//* make the http request to the url above
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to make request to COIN_PRICE",
		})
		return
	}
	//* if the request returns a string body, make sure it is closed after use
	if res.Body != nil {
		defer res.Body.Close()
	}
	//* read the string body that the request returns
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to read STRING BODY response",
		})
		return
	}
	//* parse the string body that the request returns into json
	CoinData := CoinPriceResponse{}
	jsonErr := json.Unmarshal(body, &CoinData)
	if jsonErr != nil {
		if config.Cfg.Webserver.Debug {
			fmt.Println(string(body))
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to read JSON BODY response",
		})
		return
	}
	int_coin_price, err := strconv.ParseFloat(CoinData.Prices.USD, 64)
	if err != nil {
		if config.Cfg.Webserver.Debug {
			fmt.Println(string(body))
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to convert coin_price to FLOAT",
		})
		return
	}
	//* calculate the coin amount which is equal to the usd amount given
	coin_amount := int_usd_amount / int_coin_price
	if config.Cfg.Webserver.Debug {
		fmt.Println(string(body))
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"error":       false,
		"coin_amount": coin_amount,
	})
}
