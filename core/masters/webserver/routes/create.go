package routes

import (
	"cryptoapi/core/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateResponse struct {
	Status                 string `json:"status"`
	AddressIn              string `json:"address_in"`
	AddressOut             string `json:"address_out"`
	CallbackURL            string `json:"callback_url"`
	MinimumTransactionCoin string `json:"minimum_transaction_coin"`
	Priority               string `json:"priority"`
}

func CreateAPIREQ(c *gin.Context) {
	//* recieve user queries for coin, callback url, and address
	coin := strings.ToLower(c.DefaultQuery("coin", ""))
	callback := c.DefaultQuery("callback", "")
	address := c.DefaultQuery("address", "")
	//* set the url to make a request to
	url := ("https://api.cryptapi.io/" + coin + "/create/?callback=" + callback + "&address=" + address + "&priority=default")
	spaceClient := http.Client{
		Timeout: time.Second * 10,
	}
	//* setup the http request to the url above
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to setup request to CREATE",
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
			"reason": "unable to make request to CREATE",
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
	CoinData := CreateResponse{}
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
	//* check if the api returned a "success" status
	if CoinData.Status != "success" {
		if config.Cfg.Webserver.Debug {
			fmt.Println(string(body))
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unknown error has occured",
		})
		return
	}
	if config.Cfg.Webserver.Debug {
		fmt.Println(string(body))
	}
	c.IndentedJSON(http.StatusOK, gin.H{
		"error":       false,
		"callback":    CoinData.CallbackURL,
		"address_in":  CoinData.AddressIn,
		"address_out": CoinData.AddressOut,
	})
}
