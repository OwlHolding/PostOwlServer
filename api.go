package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fastjson"
)

var MlServers []string

func InitApi(config ServerConfig) {
	MlServers = config.MlServers
}

func ApiRegUser(id int64, location int16) {
	req := MlServers[location] + "/register/" + fmt.Sprint(id)

	client := http.Client{Timeout: time.Minute}
	resp, err := client.Post(req, "", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}
}

func ApiRegChannel(id int64, location int16, channel string) []string {
	req := MlServers[location] + "/regchannel/" + fmt.Sprint(id) + "/" + channel

	client := http.Client{Timeout: time.Minute}
	resp, err := client.Post(req, "", nil)
	if err != nil {
		panic(fmt.Errorf("apiregchannel error: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		return make([]string, 0)
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}

	byte_body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("apiregchannel error: %s", err.Error()))
	}

	var parser fastjson.Parser
	body, err := parser.Parse(string(byte_body))
	if err != nil {
		panic(fmt.Errorf("apiregchannel error: %s", err.Error()))
	}

	rawposts := body.GetArray("posts")
	var posts []string
	var post string

	for _, rawpost := range rawposts {
		post = string(rawpost.GetStringBytes())
		if post != "" {
			posts = append(posts, post)
		}
	}

	return posts
}

type ApiTrainChannelRequest struct {
	Posts    []string `json:"posts"`
	Labels   []int8   `json:"labels"`
	Finetune bool     `json:"finetune"`
}

func ApiTrainChannel(
	id int64, location int16, channel string, posts []string, labels []int8,
	finetune bool) {

	req := MlServers[location] + "/train/" + fmt.Sprint(id) + "/" + channel
	reqdata, err := json.Marshal(
		ApiTrainChannelRequest{Posts: posts, Labels: labels, Finetune: finetune})
	if err != nil {
		panic(fmt.Errorf("apitrainchannel error: %s", err.Error()))
	}

	client := http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Post(req, "application/json", bytes.NewBuffer(reqdata))
	if err != nil {
		panic(fmt.Errorf("apitrainchannel error: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}
}

func ApiTrainChannelSafe(
	id int64, location int16, channel string, posts []string, labels []int8,
	finetune bool) {

	defer func() {
		err := recover()
		if err != nil {
			log.Print(err)
		}
	}()

	ApiTrainChannel(id, location, channel, posts, labels, finetune)
}

type ApiPredictReq struct {
	Time int16 `json:"time"`
}

func ApiPredict(id int64, location int16, channel string,
	sendtime int16) ([]string, string) {

	req := MlServers[location] + "/predict/" + fmt.Sprint(id) + "/" + channel
	reqdata, err := json.Marshal(ApiPredictReq{Time: sendtime})
	if err != nil {
		panic(fmt.Errorf("apipredict error: %s", err.Error()))
	}

	client := http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Post(req, "application/json", bytes.NewBuffer(reqdata))
	if err != nil {
		panic(fmt.Errorf("apipredict error: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}

	byte_body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("apipredict error: %s", err.Error()))
	}

	var parser fastjson.Parser
	body, err := parser.Parse(string(byte_body))
	if err != nil {
		panic(fmt.Errorf("apipredict error: %s", err.Error()))
	}

	rawposts := body.GetArray("posts")
	var posts []string

	for _, post := range rawposts {
		posts = append(posts, string(post.GetStringBytes()))
	}

	markup := string(body.Get("markup").GetStringBytes())

	return posts, markup
}
