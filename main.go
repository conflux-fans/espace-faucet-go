package main

import (
	"fmt"
	"github.com/conflux-fans/espace-faucet-go/routers"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func initConfig() {
	viper.SetConfigName("config")             // name of config file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/eSpace-faucet/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.eSpace-faucet") // call multiple times to add many search paths
	viper.AddConfigPath(".")                  // optionally look for config in the working directory
	err := viper.ReadInConfig()               // Find and read the config file
	if err != nil {                           // Handle errors reading the config file
		log.Fatalln(fmt.Errorf("fatal error config file: %w", err))
	}
}

func init() {
	initConfig()
}


func main() {

	app := gin.Default()
	routers.SetupRoutes(app)
	app.Run() // listen and serve on 0.0.0.0:8080
}

