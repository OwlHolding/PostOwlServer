package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const (
	StateIdle           = 0
	StateWaitAddChannel = 1
	StateRateCycle      = 2
	StateWaitDelChannel = 3
	StateWaitChangeTime = 4
	StateAddDiffPost    = 5
)

type DialogEmpty struct{}

func (dialog *DialogEmpty) Encode() []byte     { return make([]byte, 0) }
func (dialog *DialogEmpty) Decode(data []byte) {}

type DialogPostRate struct {
	Channel string
	Posts   []string
	Labels  []int8
	Iter    rune
}

func (dialog *DialogPostRate) Encode() []byte {
	bytedata, err := json.Marshal(dialog)
	if err != nil {
		log.Fatal(err)
	}
	return bytedata
}

func (dialog *DialogPostRate) Decode(bytedata []byte) {
	err := json.Unmarshal(bytedata, dialog)
	if err != nil {
		log.Fatal(err)
	}
}

var AdminChatID int64

func InitStateMachine(config ServerConfig) {
	AdminChatID = config.AdminChatID
}

func StateMachine(chatID int64, text string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Print(err)
			SendMessage(chatID, MessageError)
		}
	}()

	user := User{ID: chatID, Channels: "&", Time: -1}
	if !user.Get() {
		user.Create()
		ApiRegUser(chatID, user.Location)
	}

	userstate := UserState{ID: chatID, Data: &DialogEmpty{}}
	userstate.Get()

	if text == "/start" {
		SendMessage(chatID, MessageHello)
		return
	}

	if text == "/addchannel" {
		userstate := UserState{ID: chatID, State: StateWaitAddChannel, Data: &DialogEmpty{}}
		userstate.Set()
		SendMessage(chatID, MessageAddChannel)
		return
	}

	if text == "/delchannel" {
		userstate := UserState{ID: chatID, State: StateWaitDelChannel, Data: &DialogEmpty{}}
		userstate.Set()
		SendMessage(chatID, MessageDelChannel)
		return
	}

	if text == "/changetime" {
		userstate := UserState{ID: chatID, State: StateWaitChangeTime, Data: &DialogEmpty{}}
		userstate.Set()
		SendMessage(chatID, MessageChangeTime)
		return
	}

	if text == "/info" {
		channels := strings.ReplaceAll(user.Channels, "&", " ")

		time := "не установлено"
		if user.Time != -1 {
			time = fmt.Sprint(user.Time/60) + ":" + fmt.Sprint(user.Time%60)
		}

		SendMessage(chatID, fmt.Sprintf(MessageInfo, channels, time))
		return
	}

	if text == "/disable" {
		user.Time = -1
		user.Update()
		SendMessage(chatID, MessageUserDisabled)
		return
	}

	if userstate.State == StateWaitAddChannel {
		if strings.Contains(user.Channels, "&"+text+"&") {
			SendMessage(chatID, MessageChannelAlreadyAdded)
			return
		}

		posts := ApiRegChannel(chatID, user.Location, text)

		if len(posts) == 0 {
			SendMessage(chatID, MessageChannelNotExists)
			return
		}

		userstate.State = StateRateCycle
		userstate.Data = &DialogPostRate{Channel: text, Posts: posts}
		userstate.Set()
		SendMessage(chatID, MessageChannelOK)
		SendMessageWithKeyboard(chatID, posts[0])
		return
	}

	if userstate.State == StateRateCycle {
		userstate.Data = &DialogPostRate{}
		userstate.Get()
		data := userstate.Data.(*DialogPostRate)

		if text == "👍" || text == "👎" {
			label := 1
			if text == "👎" {
				label = 0
			}
			data.Labels = append(data.Labels, int8(label))
			data.Iter++
			if len(data.Labels) == len(data.Posts) {
				sum := 0
				for _, label := range data.Labels {
					sum += int(label)
				}
				if sum == 0 {
					SendMessage(chatID, MessageRateCycleAllNegative)
					data.Labels = append(data.Labels, 1)
					userstate.Data = data
					userstate.State = StateAddDiffPost
					userstate.Set()
					return
				}
				if sum == len(data.Labels) {
					SendMessage(chatID, MessageRateCycleAllPositive)
					data.Labels = append(data.Labels, 0)
					userstate.Data = data
					userstate.State = StateAddDiffPost
					userstate.Set()
					return
				}

				SendMessage(chatID, MessageRateCycleWait)
				ApiTrainChannel(
					chatID, user.Location, data.Channel, data.Posts,
					data.Labels)
				userstate.State = StateIdle
				userstate.Data = &DialogEmpty{}
				userstate.Set()
				user.Channels += data.Channel + "&"
				user.Update()
				SendMessageRemoveKeyboard(chatID, MessageRateCycleEnd)
			} else {
				userstate.Data = data
				userstate.Set()
				SendMessageWithKeyboard(chatID, data.Posts[data.Iter])
			}
			return
		}
		SendMessage(chatID, MessageRateCycleFormat)
		return
	}

	if userstate.State == StateWaitDelChannel {
		if !strings.Contains(user.Channels, "&"+text+"&") {
			SendMessage(chatID, MessageChannelNotListed)
			return
		}
		user.Channels = strings.ReplaceAll(user.Channels, text+"&", "")
		user.Update()
		userstate.State = StateIdle
		userstate.Set()
		SendMessage(chatID, MessageDelChannelOK)
		return
	}

	if userstate.State == StateWaitChangeTime {
		time := strings.Split(text, ":")
		if len(time) != 2 {
			SendMessage(chatID, MessageTimeInvalidFormat)
			return
		}
		hours, err := strconv.Atoi(time[0])
		if err != nil {
			SendMessage(chatID, MessageTimeInvalidFormat)
			return
		}

		minutes, err := strconv.Atoi(time[1])
		if err != nil {
			SendMessage(chatID, MessageTimeInvalidFormat)
			return
		}

		if hours > 23 || hours < 0 || minutes > 59 || minutes < 0 {
			SendMessage(chatID, MessageTimeInvalidFormat)
			return
		}

		user.Time = int16(hours)*60 + int16(minutes)
		user.Update()
		userstate.State = StateIdle
		userstate.Set()
		SendMessage(chatID, MessageChangeTimeOK)
		return
	}

	if userstate.State == StateAddDiffPost {
		userstate.Data = &DialogPostRate{}
		userstate.Get()
		data := userstate.Data.(*DialogPostRate)
		data.Posts = append(data.Posts, text)
		SendMessage(chatID, MessageRateCycleWait)
		ApiTrainChannel(
			chatID, user.Location, data.Channel, data.Posts,
			data.Labels)
		userstate.State = StateIdle
		userstate.Data = &DialogEmpty{}
		userstate.Set()
		user.Channels += data.Channel + "&"
		user.Update()
		SendMessageRemoveKeyboard(chatID, MessageRateCycleEnd)
		userstate.State = StateIdle
		userstate.Data = &DialogEmpty{}
		userstate.Set()
		return
	}

	SendMessage(chatID, MessageUnknownCommand)
}

func SendPosts(chatID int64) {
	defer func() {
		err := recover()
		if err != nil {
			log.Print(err)
		}
	}()

	user := User{ID: chatID}
	user.Get()
	channels := strings.Split(user.Channels, "&")
	channels = channels[1 : len(channels)-1]

	avalposts := false
	for _, channel := range channels {
		posts := ApiPredict(user.ID, user.Location, channel, user.Time)
		for _, post := range posts {
			if post != "" {
				avalposts = true
				SendMessage(user.ID, post+"\n\n<b>"+channel+"</b>")
			}
		}
	}
	if !avalposts {
		SendMessage(user.ID, MessageNoNewPosts)
	}
}
