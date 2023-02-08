package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/valyala/fastjson"
)

var MlServers []string

func InitApi(config ServerConfig) {
	MlServers = config.MlServers
}

func ApiRegUser(id int64, location int16) {
	req := MlServers[location] + "/register/" + fmt.Sprint(id)

	resp, err := http.Post(req, "", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml error: %d", resp.StatusCode))
	}
}

type ApiRegChannelRequest struct {
	Channel string `json:"channel"`
}

func ApiRegChannel(id int64, location int16, channel string) []string {
	req := MlServers[location] + "/regchannel/" + fmt.Sprint(id)
	reqdata, err := json.Marshal(ApiRegChannelRequest{Channel: channel})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(req, "application/json", bytes.NewBuffer(reqdata))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.Status == "Channel Not Exists" {
		return make([]string, 0)
	}

	if resp.StatusCode != 201 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}

	byte_body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(errors.New("can't read body of response"))
	}

	var parser fastjson.Parser
	body, err := parser.Parse(string(byte_body))
	if err != nil {
		panic(err)
	}

	rawposts := body.GetArray("posts")
	var posts []string

	for _, post := range rawposts {
		posts = append(posts, string(post.GetStringBytes()))
	}

	return posts
}

type ApiTrainChannelRequest struct {
	Posts   []string `json:"posts"`
	Labels  []int8   `json:"labels"`
	Channel string   `json:"channel"`
}

func ApiTrainChannel(
	id int64, location int16, channel string, posts []string, labels []int8) {

	req := MlServers[location] + "/train/" + fmt.Sprint(id)
	reqdata, err := json.Marshal(
		ApiTrainChannelRequest{Posts: posts, Labels: labels, Channel: channel})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(req, "application/json", bytes.NewBuffer(reqdata))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 202 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}
}

type ApiPredictReq struct {
	Channel string `json:"channel"`
	Time    int16  `json:"time"`
	Count   int16  `json:"count"`
}

func ApiPredict(
	id int64, location int16, channel string, time int16, count int16) []string {

	req := MlServers[location] + "/predict/" + fmt.Sprint(id)
	reqdata, err := json.Marshal(ApiPredictReq{Channel: channel, Time: time, Count: count})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(req, "application/json", bytes.NewBuffer(reqdata))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}

	byte_body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(errors.New("can't read body of response"))
	}

	var parser fastjson.Parser
	body, err := parser.Parse(string(byte_body))
	if err != nil {
		panic(err)
	}

	rawposts := body.GetArray("posts")
	var posts []string

	for _, post := range rawposts {
		posts = append(posts, string(post.GetStringBytes()))
	}

	return posts
}
