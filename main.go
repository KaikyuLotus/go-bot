package main

import "fmt"

var botToken = "TOKEN"

func main(){
	var bot = NewBot(botToken)
	bot.setUpdateHandler(updateHandler)
	bot.addCommandHandler("go", start)
	bot.startPolling(false)
}

func updateHandler(bot Bot, update Update) {
	fmt.Printf("Update ID: %d\n", update.UpdateID)
}

func start(bot Bot, update Update) {
	bot.sendChatAction(update.Message.Chat.ID, Typing)
	bot.sendMessage(update.Message.Chat.ID, "Works!")
}