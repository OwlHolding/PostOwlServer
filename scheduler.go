package main

import (
	"strings"
	"time"
)

func InitScheduler(config ServerConfig) {
	go Scheduler()
}

func Scheduler() {
	for {
		ids := DatabaseForScheduler(int16(time.Now().Hour()*60 + time.Now().Minute()))

		for _, id := range ids {
			user := User{ID: id}
			user.Get()
			channels := strings.Split(user.Channels, "&")
			channels = channels[1 : len(channels)-1]

			for _, channel := range channels {
				posts := ApiPredict(user.ID, user.Location, channel, user.Time)
				SendMessage(user.ID, posts[0])
			}
		}
		time.Sleep(time.Minute)
	}
}
