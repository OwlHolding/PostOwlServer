package main

import (
	"time"
)

type SchedulerHander func(int16)

var SchedulerTickHandler SchedulerHander

func InitScheduler(config ServerConfig, handler SchedulerHander) {
	SchedulerTickHandler = handler

	go func() {
		time.Sleep(time.Second * time.Duration((60 - time.Now().Second())))
		go SchedulerTick()
	}()
}

func SchedulerTick() {
	for {
		go SchedulerTickHandler(int16(time.Now().Hour()*60 + time.Now().Minute()))
		time.Sleep(time.Minute)
	}
}
