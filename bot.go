package main

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	fasthttp "github.com/valyala/fasthttp"
)

type TgBotUpdateHandler func(int64, string)

var BotAPI *tgbotapi.BotAPI
var WebhookPath string
var BotUpdateHandler TgBotUpdateHandler

func InitBot(config ServerConfig, handler TgBotUpdateHandler) {
	WebhookPath = "/" + config.Token

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatal(err)
	}

	webhook, _ := tgbotapi.NewWebhookWithCert(config.Url+":"+config.Port+WebhookPath,
		tgbotapi.FilePath(config.CertFile))
	webhook.MaxConnections = config.MaxBotConns

	_, err = bot.Request(webhook)
	if err != nil {
		log.Fatal(err)
	}

	BotAPI = bot
	BotUpdateHandler = handler
}

func ProcessRequest(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != WebhookPath {
		ctx.Error("", fasthttp.StatusForbidden)
		return
	}

	var update tgbotapi.Update
	err := json.Unmarshal(ctx.PostBody(), &update)
	if err != nil {
		ctx.Error("", fasthttp.StatusBadRequest)
		return
	}
	BotUpdateHandler(update.Message.From.ID, update.Message.Text)
}

func SendMessage(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	message.ParseMode = "HTML"
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}

func SendMessageWithKeyboard(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	keyboard := tgbotapi.ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("👍"),
				tgbotapi.NewKeyboardButton("👎"))},
		ResizeKeyboard: true, OneTimeKeyboard: true,
	}
	message.ReplyMarkup = keyboard
	message.ParseMode = "HTML"
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}

func SendMessageRemoveKeyboard(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	keyboard := tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
	message.ReplyMarkup = keyboard
	message.ParseMode = "HTML"
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}
