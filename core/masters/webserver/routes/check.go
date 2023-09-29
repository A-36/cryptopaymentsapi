package routes

import (
	"cryptoapi/core/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CheckResponse struct {
	Status              string `json:"status"`
	CallbackURL         string `json:"callback_url"`
	AddressIn           string `json:"address_in"`
	AddressOut          string `json:"address_out"`
	NotifyPending       bool   `json:"notify_pending"`
	NotifyConfirmations int    `json:"notify_confirmations"`
	Priority            string `json:"priority"`
	Callbacks           []struct {
		UUID               string        `json:"uuid"`
		LastUpdate         string        `json:"last_update"`
		Result             string        `json:"result"`
		Confirmations      int           `json:"confirmations"`
		FeePercent         string        `json:"fee_percent"`
		Fee                int           `json:"fee"`
		Value              int           `json:"value"`
		ValueCoin          string        `json:"value_coin"`
		ValueForwarded     int           `json:"value_forwarded"`
		ValueForwardedCoin string        `json:"value_forwarded_coin"`
		TxidIn             string        `json:"txid_in"`
		TxidOut            string        `json:"txid_out"`
		Coin               string        `json:"coin"`
		Logs               []interface{} `json:"logs"`
	} `json:"callbacks"`
}

func CheckAPIREQ(c *gin.Context) {
	//* recieve user queries for coin, callback url, coin amount, and confirmations
	coin := c.DefaultQuery("coin", "")
	callback := c.DefaultQuery("callback", "")
	coin_amount := c.DefaultQuery("coin_amount", "")
	confirmations := c.DefaultQuery("confirmations", "")
	//* convert the provided confirmations from string to integer
	int_confirmations, err := strconv.Atoi(confirmations)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to convert confirmations to INT",
		})
		return
	}
	//* convert the provided coin amount from string to float
	int_coin_amount, err := strconv.ParseFloat(coin_amount, 64)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to convert coin_amount to FLOAT",
		})
		return
	}
	//* check if the provided confirmations is more than one
	if int_confirmations < 1 {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "invalid confirmations provided (minimum: 1)",
		})
		return
	}
	//* set the url to make a request to
	url := ("https://api.cryptapi.io/" + coin + "/logs/?callback=" + callback)
	spaceClient := http.Client{
		Timeout: time.Second * 10,
	}
	//* setup the http request to the url above
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":  true,
			"reason": "unable to setup request to CHECK",
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
			"reason": "unable to make request to CHECK",
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
	CoinData := CheckResponse{}
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
	//* set and create variables to later calculate recieved and confirmed amount of money
	var total_recieved float64
	var money_recieved float64
	total_recieved = 0
	money_recieved = 0
	//* range through the callbacks
	for _, callbacks := range CoinData.Callbacks {
		//* parse the amount of crypto recieved in the callback
		each_recieved, err := strconv.ParseFloat(callbacks.ValueCoin, 64)
		if err != nil {
			if config.Cfg.Webserver.Debug {
				fmt.Println(string(body))
			}
			c.IndentedJSON(http.StatusOK, gin.H{
				"error":  true,
				"reason": "unable to get coin amount paid FLOAT",
			})
			return
		}
		//* if it has been confirmed add it to "total_recieved"
		if callbacks.Confirmations >= int_confirmations {
			total_recieved += each_recieved
			//* if it has not been confirmed add it to "money_recieved"
		} else {
			money_recieved += each_recieved
		}
	}
	//* if the money has been confirmed and is equal to or more than the amount requested in the query return "confirmed" and "recieved"
	if float64(total_recieved) >= float64(int_coin_amount) {
		if config.Cfg.Webserver.Debug {
			fmt.Println(string(body))
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":     false,
			"status":    "confirmed",
			"received":  true,
			"confirmed": true,
		})
		return
	}
	//* if the money has not been confirmed but is equal to or more than the amount requested in the query return "recieved"
	if float64(money_recieved) >= float64(int_coin_amount) {
		if config.Cfg.Webserver.Debug {
			fmt.Println(string(body))
		}
		c.IndentedJSON(http.StatusOK, gin.H{
			"error":     false,
			"status":    "unconfirmed",
			"received":  true,
			"confirmed": false,
		})
		return
	}
	if config.Cfg.Webserver.Debug {
		fmt.Println(string(body))
	}
	//* if the money has not been confirmed or recieved return "unpaid"
	c.IndentedJSON(http.StatusOK, gin.H{
		"error":     false,
		"status":    "unpaid",
		"received":  false,
		"confirmed": false,
	})
}
