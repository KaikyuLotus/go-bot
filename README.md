# GoBot

**This is just a small training project and it's not intended for everyday use, not yet.**<br>
Please feel free to give me some tips.<br>
This is my first Go project.<br>

## Examples

### Hello World Command
```go
package main

import . "github.com/kaikyudev/gobot"


func main() {
    // instantiate our bot with out custom costructor
    var bot = NewBot("TOKEN")

    // Add our funcion in response to /start
    bot.addCommandHandler("start", start)

    // Start getting updates
    bot.startPolling(true)

    // Wait while our bot runs
    bot.Idle()
}

func start(bot *Bot, update Update) {
    msgID := update.Message.MessageID
    chatID := update.Message.Chat.ID

    // Typing...
    bot.SendChatAction(update.Message.Chat.ID, Typing)

    // With SendMessageArgs{} you can pass extra args for sendMessage method
    bot.SendMessage(chatID, "*Markdown*", SendMessageArgs{ReplyToMessageID:msgID, ParseMode:Markdown})
} // That's it!
```

### Getting Updates
```go
package main

import (
    . "github.com/kaikyudev/gobot"
    "log"
)


func main() {
    var bot = NewBot("TOKEN")

    // Add our funcion in response to any update
    bot.setUpdateHandler(updateHandler)

    bot.startPolling(true)
    bot.Idle()
}

func updateHandler(bot *Bot, update Update) {
    // Log every update
    log.Printf("[%s] %s", update.Message.From.Username, update.Message.Text)
}
```


## What is working

* Basic API requests
* Methods: getMe, sendMessage, sendChatAction, sendPhoto, sendDocument [more to come...]
* Getting updates
* Commands handling

Not much more...