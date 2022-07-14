package main

import (
	"github.com/conflux-fans/espace-faucet-go/routers"
	"github.com/gin-gonic/gin"
)

func main() {

	app := gin.Default()
	routers.SetupRoutes(app)
	app.Run() // listen and serve on 0.0.0.0:8080
}

