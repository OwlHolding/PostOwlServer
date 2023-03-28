package main

import (
	"log"

	"github.com/valyala/fasthttp"
)

func main() {
	config := LoadConfigFromEnv("POSTOWLCONFIG")

	InitDatabase(config)
	InitRedis(config)
	InitApi(config)
	InitStateMachine(config)
	InitBot(config, StateMachine, RatePost)
	InitScheduler(config, SendPosts)

	log.Print("Server started")

	err := fasthttp.ListenAndServeTLS(":"+config.Port, config.CertFile, config.KeyFile,
		ProcessRequest)
	log.Fatal(err)
}
