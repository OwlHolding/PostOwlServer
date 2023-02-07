package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	config := LoadConfigFromEnv("POSTOWLCONFIG")

	log.Print("Server started")

	err := fasthttp.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile,
		nil)
	log.Fatal(err)
}
