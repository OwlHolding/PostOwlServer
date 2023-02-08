package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	redis "github.com/redis/go-redis/v9"
)

type UserData interface {
	Encode() []byte
	Decode([]byte)
}

type UserState struct {
	ID    int64
	State int16
	Data  UserData
}

type RedisRecord struct {
	State int16
	Data  []byte
}

var MainCtx context.Context
var RedisClient *redis.Client

func InitRedis(config ServerConfig) {
	MainCtx = context.Background()
	RedisClient = redis.NewClient(&redis.Options{Addr: config.RedisUrl})
	_, err := RedisClient.Ping(MainCtx).Result()
	if err != nil {
		log.Fatal(err)
	}
}

func (state *UserState) Get() bool {
	value, err := RedisClient.Get(MainCtx, fmt.Sprint(state.ID)).Result()
	if err != nil {
		if err == redis.Nil {
			return false
		} else {
			log.Fatal(err)
		}
	}

	var record RedisRecord
	err = json.Unmarshal([]byte(value), &record)
	if err != nil {
		return false
	}

	state.State = record.State
	state.Data.Decode(record.Data)

	return true
}

func (state *UserState) Set() {
	key := fmt.Sprint(state.ID)
	record := RedisRecord{State: state.State, Data: state.Data.Encode()}
	value, err := json.Marshal(record)
	if err != nil {
		log.Fatal(err)
	}

	err = RedisClient.Set(MainCtx, key, value, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}
