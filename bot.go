package main

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	fasthttp "github.com/valyala/fasthttp"
)

type ReplyKeyboardMarkup struct {
	Keyboard              [][]tgbotapi.KeyboardButton `json:"keyboard"`
	ResizeKeyboard        bool                        `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool                        `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string                      `json:"input_field_placeholder,omitempty"`
	Selective             bool                        `json:"selective,omitempty"`
	IsPersistent          bool                        `json:"is_persistent,omitempty"`
}

type TgBotUpdateHandler func(int64, string, string)
type TgBotCallbackHandler func(int64, int, string, string)

var BotAPI *tgbotapi.BotAPI
var WebhookPath string
var BotUpdateHandler TgBotUpdateHandler
var BotCallbackHandler TgBotCallbackHandler

func InitBot(config ServerConfig, handler TgBotUpdateHandler,
	callback TgBotCallbackHandler) {
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
	BotCallbackHandler = callback
}

func ProcessRequest(ctx *fasthttp.RequestCtx) {
	defer func() {
		err := recover()
		if err != nil {
			log.Printf("telegram sent bad request %s", err)
		}
	}()

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

	if update.CallbackQuery != nil {
		BotCallbackHandler(update.CallbackQuery.From.ID, update.CallbackQuery.Message.MessageID,
			update.CallbackQuery.Data, update.CallbackQuery.Message.Text)
	} else {
		go BotUpdateHandler(update.Message.From.ID,
			update.Message.Text, update.Message.Chat.UserName)
	}
}

func SendMessage(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	message.ParseMode = "HTML"
	message.DisableWebPagePreview = true
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}

func SendMessageWithKeyboard(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	keyboard := ReplyKeyboardMarkup{
		Keyboard: [][]tgbotapi.KeyboardButton{
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("👍"),
				tgbotapi.NewKeyboardButton("👎"))},
		ResizeKeyboard: true, IsPersistent: true}
	message.ReplyMarkup = keyboard

	message.ParseMode = "HTML"
	message.DisableWebPagePreview = true
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
	message.DisableWebPagePreview = true
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}

func SendMessageWithInlineKeyboard(chatID int64, text string, channel string) {
	message := tgbotapi.NewMessage(chatID, text)
	keyboard := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("👍", "1"+channel),
				tgbotapi.NewInlineKeyboardButtonData("👎", "0"+channel))},
	}
	message.ReplyMarkup = keyboard
	message.ParseMode = "HTML"
	message.DisableWebPagePreview = true
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}

func DisableInlineKeyboard(chatID int64, messageID int) {
	message := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID,
		tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{}))
	_, err := BotAPI.Send(message)
	if err != nil {
		panic(fmt.Errorf("botapi error: %s", err.Error()))
	}
}
