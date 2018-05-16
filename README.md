# GoBot

**This is just a small training project and it's not intended for everyday use.**<br>
Please feel free to give me some tips.<br>
This is my first Go project.<br>

## Examples

### Hello World Command
```go
package main

func main() {
	// instantiate our bot with out custom costructor
	var bot = NewBot("TOKEN")

	// Add our funcion in response to /helloWorld
	bot.addCommandHandler("helloWorld", helloWorld)

	// Start getting updates
	bot.startPolling(true)
}

func helloWorld(bot Bot, update Update) {
	bot.sendChatAction(update.Message.Chat.ID, Typing)
	bot.sendMessage(update.Message.Chat.ID, "Works!")

} // That's it!
```

### Getting Updates
```go
package main

func main() {
	// instantiate our bot with out custom costructor
	var bot = NewBot("TOKEN")

	// Add our funcion in response to any update
	bot.setUpdateHandler(updateHandler)

	// Start getting updates
	bot.startPolling(true)
}

func updateHandler(bot Bot, update Update) {
	bot.sendChatAction(update.Message.Chat.ID, Typing)
	bot.sendMessage(update.Message.Chat.ID, "Works!")
}
```


## What is working

* Basic API requests
* Methods: getMe, sendMessage, sendChatAction, getUpdate [more to come...]
* Commands handling

Not much more...