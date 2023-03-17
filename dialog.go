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
	StateAddMarkCycle   = 6
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

type DialogAddMark struct {
	Posts    []string
	Channels []string
}

func (dialog *DialogAddMark) Encode() []byte {
	bytedata, err := json.Marshal(dialog)
	if err != nil {
		log.Fatal(err)
	}
	return bytedata
}

func (dialog *DialogAddMark) Decode(bytedata []byte) {
	err := json.Unmarshal(bytedata, dialog)
	if err != nil {
		log.Fatal(err)
	}
}

func FormatTime(hours int16, minutes int16) string {
	time := ""
	if hours < 10 {
		time += "0"
	}
	time += fmt.Sprint(hours) + ":"
	if minutes < 10 {
		time += "0"
	}
	time += fmt.Sprint(minutes)
	return time
}

var AdminChatIDs []int64
var ChansPerUser int
var LocationsCount int16
var BanList []int64
var WhiteList []int64

func InitStateMachine(config ServerConfig) {
	AdminChatIDs = config.AdminChatIDs
	ChansPerUser = config.ChansPerUser
	LocationsCount = int16(len(config.MlServers))
	BanList = config.BanList
	WhiteList = config.WhiteList
}

func StateMachine(chatID int64, text string, username string) {
	defer func() {
		err := recover()
		if err != nil {
			log.Print(err)
			SendMessage(chatID, MessageError)
			if len(AdminChatIDs) != 0 {
				for _, AdminChatID := range AdminChatIDs {
					if AdminChatID != 0 {
						SendMessage(AdminChatID,
							fmt.Sprintf(`Error: "%s"; Username: @%s`, err, username))
					}
				}
			}
		}
	}()

	if len(WhiteList) > 0 {
		ban := true
		for _, id := range WhiteList {
			if id == chatID {
				ban = false
				break
			}
		}
		if ban {
			SendMessage(chatID, MessageBanned)
			return
		}
	} else {
		for _, id := range BanList {
			if id == chatID {
				SendMessage(chatID, MessageBanned)
				return
			}
		}
	}

	user := User{ID: chatID, Channels: "&", Time: -1}
	if !user.Get() {
		var min int64 = 9223372036854775807
		var minloc int16 = 0
		var iter int16 = 0
		for iter < LocationsCount {
			entry := DatabaseCountLocation(iter)
			if entry < min {
				min = entry
				minloc = iter
			}
			iter++
		}
		user.Location = minloc
		ApiRegUser(chatID, user.Location)
		user.Create()
		if len(AdminChatIDs) != 0 {
			for _, AdminChatID := range AdminChatIDs {
				if AdminChatID != 0 {
					SendMessage(AdminChatID, fmt.Sprintf("New user registered: @%s %d", username, chatID))
				}
			}
		}
	}

	userstate := UserState{ID: chatID, Data: &DialogEmpty{}}
	userstate.Get()

	if text == "/start" {
		userstate := UserState{ID: chatID, State: StateIdle, Data: &DialogEmpty{}}
		userstate.Set()
		SendMessageRemoveKeyboard(chatID, MessageHello)
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
		channelslist := strings.Split(user.Channels, "&")
		channelslist = channelslist[1 : len(channelslist)-1]
		channels := ""
		for _, channel := range channelslist {
			channels += "\t<code>" + channel + "</code>\n"
		}
		if len(channels) > 2 {
			channels = channels[:len(channels)-1]
		}

		time := "не установлено"
		if user.Time != -1 {
			time = FormatTime(user.Time/60, user.Time%60)
		}

		userstate := UserState{ID: chatID, State: StateIdle, Data: &DialogEmpty{}}
		userstate.Set()
		SendMessage(chatID, fmt.Sprintf(MessageInfo, time, channels))
		return
	}

	if text == "/disable" {
		user.Time = -1
		user.Update()
		userstate := UserState{ID: chatID, State: StateIdle, Data: &DialogEmpty{}}
		userstate.Set()
		SendMessage(chatID, MessageUserDisabled)
		return
	}

	if text == "/cancel" {
		userstate := UserState{ID: chatID, State: StateIdle, Data: &DialogEmpty{}}
		userstate.Set()
		SendMessageRemoveKeyboard(chatID, MessageCancel)
		return
	}

	if userstate.State == StateWaitAddChannel {
		text = strings.ReplaceAll(text, "https://t.me/", "")
		text = strings.ReplaceAll(text, "@", "")

		if strings.Contains(user.Channels, "&"+text+"&") {
			SendMessage(chatID, MessageChannelAlreadyAdded)
			return
		}

		if strings.Count(user.Channels, "&")-1 >= ChansPerUser {
			SendMessage(chatID, fmt.Sprintf(MessageChannelOverflow, ChansPerUser))
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
			var label int8 = 1
			if text == "👎" {
				label = 0
			}
			data.Labels = append(data.Labels, label)
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

				SendMessageRemoveKeyboard(chatID, MessageRateCycleWait)
				ApiTrainChannel(
					chatID, user.Location, data.Channel, data.Posts,
					data.Labels, false)
				userstate.State = StateIdle
				userstate.Data = &DialogEmpty{}
				userstate.Set()
				user.Channels += data.Channel + "&"
				user.Update()
				SendMessage(chatID, MessageRateCycleEnd)
			} else {
				userstate.Data = data
				userstate.Set()
				SendMessageWithKeyboard(chatID, data.Posts[data.Iter])
			}
		} else {
			SendMessage(chatID, MessageRateCycleFormat)
		}
		return
	}

	if userstate.State == StateWaitDelChannel {
		text = strings.ReplaceAll(text, "https://t.me/", "")
		text = strings.ReplaceAll(text, "@", "")

		if !strings.Contains(user.Channels, "&"+text+"&") {
			SendMessage(chatID, MessageChannelNotListed)
			return
		}
		user.Channels = strings.ReplaceAll(user.Channels, text+"&", "")
		user.Update()
		userstate.State = StateIdle
		userstate.Set()
		SendMessage(chatID, fmt.Sprintf(MessageDelChannelOK, text))
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
		SendMessage(chatID,
			fmt.Sprintf(MessageChangeTimeOK, FormatTime(int16(hours), int16(minutes))))
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
			data.Labels, false)
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

	if userstate.State == StateAddMarkCycle {
		userstate.Data = &DialogAddMark{}
		userstate.Get()
		data := userstate.Data.(*DialogAddMark)

		if text == "👍" || text == "👎" {
			var label int8 = 1
			if text == "👎" {
				label = 0
			}

			go ApiTrainChannelSafe(user.ID, user.Location,
				data.Channels[len(data.Channels)-1],
				data.Posts[len(data.Posts)-1:len(data.Posts)], []int8{int8(label)}, true)

			data.Channels = data.Channels[:len(data.Channels)-1]
			data.Posts = data.Posts[:len(data.Posts)-1]

			if len(data.Posts) > 0 {
				SendMessageWithKeyboard(chatID,
					data.Posts[len(data.Posts)-1]+"\n"+fmt.Sprintf(`<a href="t.me/%s">%s</a>`,
						data.Channels[len(data.Channels)-1], data.Channels[len(data.Channels)-1]))
				userstate.Data = data
			} else {
				SendMessageRemoveKeyboard(chatID, MessageMarkCycleEnd)
				userstate.State = StateIdle
				userstate.Data = &DialogEmpty{}
			}
			userstate.Set()

		} else {
			SendMessage(chatID, MessageRateCycleFormat)
		}
		return
	}

	SendMessage(chatID, MessageUnknownCommand)
}

func SendPosts(time int16) {
	defer func() {
		err := recover()
		if err != nil {
			log.Print(err)
			if len(AdminChatIDs) != 0 {
				for _, AdminChatID := range AdminChatIDs {
					if AdminChatID != 0 {
						SendMessage(AdminChatID,
							fmt.Sprintf(`Scheduler error: "%s"`, err))
					}
				}
			}
		}
	}()

	users := DatabaseForScheduler(time)
	for _, chatID := range users {
		user := User{ID: chatID}
		user.Get()
		if user.Time == 0 {
			user.Time = 1
		}
		channels := strings.Split(user.Channels, "&")
		channels = channels[1 : len(channels)-1]

		userstate := UserState{ID: chatID, State: StateAddMarkCycle}
		data := DialogAddMark{}

		avalposts := false
		for _, channel := range channels {
			posts, markup := ApiPredict(user.ID, user.Location, channel, user.Time)
			for _, post := range posts {
				if post != "" {
					avalposts = true
					SendMessage(chatID,
						post+"\n"+fmt.Sprintf(`<a href="t.me/%s">%s</a>`, channel, channel))
				}
			}
			if markup != "" {
				data.Posts = append(data.Posts, markup)
				data.Channels = append(data.Channels, channel)
			}
		}

		if !avalposts {
			SendMessage(user.ID, MessageNoNewPosts)
		} else {
			SendMessage(chatID, MessageMarkPlease)
			SendMessageWithKeyboard(chatID,
				data.Posts[len(data.Posts)-1]+"\n"+fmt.Sprintf(`<a href="t.me/%s">%s</a>`,
					data.Channels[len(data.Channels)-1], data.Channels[len(data.Channels)-1]))
			userstate.Data = &data
			userstate.Set()
		}
	}
}
