package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 208 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}
}

func ApiRegChannel(id int64, location int16, channel string) []string {
	req := MlServers[location] + "/regchannel/" + fmt.Sprint(id) + "/" + fmt.Sprint(channel)

	resp, err := http.Post(req, "", nil)
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
	Posts  []string `json:"posts"`
	Labels []int8   `json:"labels"`
}

func ApiTrainChannel(
	id int64, location int16, channel string, posts []string, labels []int8) {

	req := MlServers[location] + "/train/" + fmt.Sprint(id) + "/" + fmt.Sprint(channel)
	reqdata, err := json.Marshal(
		ApiTrainChannelRequest{Posts: posts, Labels: labels})
	if err != nil {
		panic(fmt.Errorf("apitrainchannel error: %s", err.Error()))
	}

	resp, err := http.Post(req, "application/json", bytes.NewBuffer(reqdata))
	if err != nil {
		panic(fmt.Errorf("apitrainchannel error: %s", err.Error()))
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		panic(fmt.Errorf("ml server error: %d", resp.StatusCode))
	}
}

type ApiPredictReq struct {
	Time int16 `json:"time"`
}

func ApiPredict(id int64, location int16, channel string, time int16) []string {

	req := MlServers[location] + "/predict/" + fmt.Sprint(id) + "/" + fmt.Sprint(channel)
	reqdata, err := json.Marshal(ApiPredictReq{Time: time})
	if err != nil {
		panic(fmt.Errorf("apipredict error: %s", err.Error()))
	}

	resp, err := http.Post(req, "application/json", bytes.NewBuffer(reqdata))
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

	return posts
}
