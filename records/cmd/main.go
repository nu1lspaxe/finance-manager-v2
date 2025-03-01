package main

import (
	"fmt"
	"os"
	"records/server"

	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()
}

func main() {
	app, err := server.NewServer()
	if err != nil {
		fmt.Printf("failed to start application: %v", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		fmt.Printf("failed to run application: %v", err)
		os.Exit(1)
	}
}
