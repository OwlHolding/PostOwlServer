package main

import (
	"time"
)

func InitScheduler(config ServerConfig) {
	go func() {
		time.Sleep(time.Second * time.Duration((60 - time.Now().Second())))
		go SchedulerTick()
	}()
}

func SchedulerTick() {
	for {
		users := DatabaseForScheduler(int16(time.Now().Hour()*60 + time.Now().Minute()))
		for _, user := range users {
			go SendPosts(user)
		}
		time.Sleep(time.Minute)
	}
}
