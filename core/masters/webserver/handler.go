package webserver

import (
	"cryptoapi/core/config"
	"cryptoapi/core/masters/webserver/routes"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func Start() {
	if config.Cfg.Webserver.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("assets/html", true)))
	group := router.Group("api")
	{
		group.GET("create", routes.CreateAPIREQ)
		group.GET("check", routes.CheckAPIREQ)
		group.GET("coin_price", routes.CoinPriceAPIREQ)
	}
	router.Run(config.Cfg.Webserver.Host)
}
