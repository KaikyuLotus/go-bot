# GoBot

**This is just a small training project and it's not intended for everyday use.**<br>
Please feel free to give me some tips.<br>
This is my first Go project.<br>

## Examples

### Hello World Command
```go
package main

import . "github.com/kaikyudev/gobot"


func main(){
	// Create our bot
	var bot = NewBot("TOKEN")

	// Add a /command handler function
	bot.AddCommandHandler("start", start)

	// Start getting updates, this is a blocking function
	bot.StartPolling(false)
}

// sendMessage example
func start(bot Bot, update Update) {
	msgID := update.Message.MessageID
	chatID := update.Message.Chat.ID
	bot.SendChatAction(update.Message.Chat.ID, Typing)
	                                           // Use SendMessageArgs{} to pass only what you need
	bot.SendMessage(chatID, "*Markdown*", SendMessageArgs{ReplyToMessageID:msgID, ParseMode:Markdown})
}
```

### Getting Updates
```go
package main

import (
	"fmt"
	. "github.com/kaikyudev/gobot"
)


func main(){
	// Create our bot
	var bot = NewBot("TOKEN")

	// Set the update handler function
	bot.SetUpdateHandler(updateHandler)

	// Start getting updates, this is a blocking function
	bot.StartPolling(false)
}

// Here our updates will land
func updateHandler(bot Bot, update Update) {
	fmt.Printf("Update ID: %d\n", update.UpdateID)
}
```


## What is working

* Basic API requests
* Methods: getMe, sendMessage, sendPhoto, sendDocument, sendChatAction [more to come...]
* Commands handling
* Getting Updates

Not much more...