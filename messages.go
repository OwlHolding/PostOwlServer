package main

const (
	MessageHello                = "Привет! Я - почтовая сова"
	MessageAddChannel           = "Пришли мне название канала. Например, TelegramTips"
	MessageChannelAlreadyAdded  = "Такой канал уже добавлен ранее"
	MessageChannelNotExists     = "Такой канал не существует"
	MessageChannelOK            = "Отлично! Теперь давай обучим модель"
	MessageRateCycleEnd         = "Обучение завершено, спасибо!"
	MessageRateCycleFormat      = "Пришли 👍 или 👎, пожалуйста"
	MessageRateCycleWait        = "Нужно немного времени, чтобы модель обучилась"
	MessageRateCycleAllPositive = "Ты отметил все посты, как интересные. Это круто, но нужен хотя бы один отрицательный пример. Перешли, пожалуйста, один пост, который тебе неинтересен"
	MessageRateCycleAllNegative = "Ты отметил все посты, как неинтересные. Это круто, но нужен хотя бы один положительный пример. Перешли, пожалуйста, один пост, который тебе интересен"
	MessageDelChannel           = "Пришли мне название канала, который нужно удалить"
	MessageChannelNotListed     = "Канал не найден"
	MessageDelChannelOK         = "Канал успешно удален"
	MessageChangeTime           = "Пришли мне время, в которое присылать посты. Например, 17:00"
	MessageTimeInvalidFormat    = "Ошибка формата времени"
	MessageChangeTimeOK         = "Время успешно изменено"
	MessageUserDisabled         = "Рассылка отключена. Для повторной активации заново установите время"
	MessageInfo                 = "Каналы: %s\nВремя рассылки: %s"
	MessageNoNewPosts           = "Интересных постов сегодня нет 😒"
	MessageUnknownCommand       = "Неизвестная команда"
	MessageError                = "Что-то пошло не так... По этой ошибке уже подняты наши системные администраторы и программисты. Но это не точно"
)