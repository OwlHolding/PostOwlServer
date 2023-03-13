package main

const (
	MessageHello                = "Привет!\nЯ – Почтовая сова. Моя главная задача – экономить твоё время. Я обучаюсь на твоих предпочтениях и каждый день присылаю сводку самых полезных постов за сутки. Для начала нам нужно добавить хотя бы один канал с помощью команды /addchannel, и установить время рассылки командой /changetime"
	MessageAddChannel           = "Пришли мне название канала или ссылку на него. Например, <code>forbesrussia</code>"
	MessageChannelAlreadyAdded  = "Такой канал уже добавлен ранее"
	MessageChannelNotExists     = `К сожалению, я не могу найти такой канал. Возможно, он не существуют, или наши коллеги из <a href="rapidapi.com">RapidAPI</a> еще не научились с ним работать😭`
	MessageChannelOK            = "Отлично! Теперь давай обучим модель"
	MessageChannelOverflow      = "К сожалению, пока нельзя добавить более %d каналов😭"
	MessageRateCycleEnd         = "Обучение завершено, спасибо🥰"
	MessageRateCycleFormat      = "Пришли 👍 или 👎, пожалуйста"
	MessageRateCycleWait        = "Нужно немного времени, чтобы модель обучилась"
	MessageRateCycleAllPositive = "Нужен хотя бы один отрицательный пример. Перешли, пожалуйста, один пост, который тебе неинтересен"
	MessageRateCycleAllNegative = "Нужен хотя бы один положительный пример. Перешли, пожалуйста, один пост, который тебе интересен"
	MessageDelChannel           = "Пришли мне название канала или ссылку на канал, который нужно удалить"
	MessageChannelNotListed     = "Такого канал нет в списке рассылки"
	MessageDelChannelOK         = "Канал <code>%s</code> успешно удален"
	MessageChangeTime           = "Пришли мне время рассылки, которое тебе удобно. Например, 17:00"
	MessageTimeInvalidFormat    = "Пришли время в формате мм:сс, пожалуйста"
	MessageChangeTimeOK         = "Спасибо! Теперь я буду прислать рассылку в %s"
	MessageUserDisabled         = "Рассылка успешно отключена. Для повторной активации заново установи время с помощью команды /changetime"
	MessageInfo                 = "Время рассылки: %s\nСписок рассылки:\n%s"
	MessageNoNewPosts           = "Интересных постов сегодня нет 😒"
	MessageMarkPlease           = "Нужно дообучить модель"
	MessageMarkCycleEnd         = "Дообучение завершено"
	MessageUnknownCommand       = "Неизвестная команда"
	MessageCancel               = "Команда отменена"
	MessageError                = "Что-то пошло не так... По этой ошибке уже подняты наши системные администраторы и программисты. Но это не точно🙃"
)
