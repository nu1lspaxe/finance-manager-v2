package main

import (
	"fmt"
	"os"
	"records_bank/server"
)

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
