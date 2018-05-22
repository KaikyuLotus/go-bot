# GoBot

**This is just a small training project and it's not intended for everyday use, not yet.**<br>
Please feel free to give me some tips.<br>
This is my first Go project.<br>

## Examples

# These examples are outdated, please be patient.

### Hello World Command
```go
package main

import . "github.com/kaikyudev/gobot"


func main() {
    // instantiate our bot with out custom costructor
    var bot = NewBot("TOKEN")

    // Add our funcion in response to /start
    bot.AddCommandHandler("start", start)

    // Start getting updates, pass true to clean pending updates, otherwise pass false
    bot.StartPolling(true)

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
    bot.SetUpdateHandler(updateHandler)

    bot.StartPolling(true)
    bot.Idle()
}

func updateHandler(bot *Bot, update Update) {
    // Log every update
    log.Printf("[%s] %s", update.Message.From.Username, update.Message.Text)
}
```

### Handling Panics
```go
package main

import (
    . "github.com/kaikyudev/gobot"
    "log"
)


func main() {
    var bot = NewBot("TOKEN")

    // Add our funcion in response to every panic
    bot.SetErrorHandler(errorHandler)

    bot.AddCommandHandler("panic", panicFoo)

    bot.StartPolling(true)
    bot.Idle()
}

func errorHandler(bot *Bot, update Update, errorMessage string) {
    // Log every panic
    log.Printf("Update #%d has caused a Panic with error message %s", update.UpdateID, error)
}

func panicFoo(bot *Bot, update Update){
    log.Print("Starting panic...")
    panic(0)
}
```

## What is working

* Basic API requests
* Methods: GetMe, GetUpdates, SendMessage, SendChatAction, SendPhoto, SendDocument [more to come...]
* Getting updates
* Commands handling
* Errors handling

Not much more...