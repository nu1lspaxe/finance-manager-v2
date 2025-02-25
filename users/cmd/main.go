package main

import (
	"os"
	"users/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	certFile := os.Getenv("CERT_FILE")
	keyFile := os.Getenv("KEY_FILE")

	server, err := server.NewServer(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	if err := server.Run(certFile, keyFile); err != nil {
		panic(err)
	}
}
