all:
	@go build -o bin/postowlserver main.go config.go bot.go redis.go database.go api.go dialog.go messages.go scheduler.go