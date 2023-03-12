package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	config := LoadConfigFromEnv("POSTOWLCONFIG")

	InitBot(config, StateMachine)
	InitRedis(config)
	InitDatabase(config)
	InitApi(config)
	InitStateMachine(config)
	InitScheduler(config, SendPosts)

	log.Print("Server started")

	err := fasthttp.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile,
		ProcessRequest)
	log.Fatal(err)
}
