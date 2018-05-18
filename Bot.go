package gobot

import (
	"log"
	"time"
	"strings"
)

/*
I know that Go is not an object-oriented language but i think that this approach is not that bad
So please don't kill me <3
*/

// I need to put return types on every foo in the future
// Bot constructor
func NewBot(token string) *Bot {
	bot := &Bot{Token:token, Running:true}
	bot.getMe()
	return bot
}

func (bot *Bot) SetUpdateHandler(foo UpdateHandlerType) {
	bot.UpdateHandler = foo
}

func (bot *Bot) SetErrorHandler(foo ErrorHandlerType) {
	bot.ErrorHandler = foo
}

func (bot *Bot) AddCommandHandler(command string, foo CommandHandlerType) {
	bot.CommandHandlers = append(bot.CommandHandlers, CommandStruct{command, foo})
}



func (bot *Bot) getMe() GetMeResult {
	result, getMeOk := getMe(bot.Token)
	if getMeOk != RequestOk {
		panic(getMeOk)
	} // For now i'll just panic here if unauthorized
	bot.ThisBot = result
	return result
}

// First parameter is offset and second is timeout
func (bot *Bot) getUpdates(offset int64, timeout bool) (GetUpdateResult, int) {
	return getUpdates(bot.Token, offset, timeout)
}

func (bot *Bot) elaborateUpdate(update Update){
	defer func() {
		// recover from panic if one occured. Set err to nil otherwise.
		if recover() != nil {
			log.Printf("gobot has recovered from a panic caused from the last handler.")
			if bot.ErrorHandler != nil {
				bot.ErrorHandler(bot, update, "error message dummy")
			}
		}
	}()

	bot.Offset = update.UpdateID + 1
	for _, commandStruct := range bot.CommandHandlers {
		if strings.ToLower(update.Message.Text) ==  "/" + strings.ToLower(commandStruct.Command) {
			commandStruct.Function(bot, update)
			return
		}
	}
	// Fix for issue #1
	if bot.UpdateHandler == nil {
		return
	}
	bot.UpdateHandler(bot, update)
}

func (bot *Bot) Stop() {
	bot.Running = false
}

func (bot *Bot) pollingFunction() {
	for bot.Running {
		updates, _ := getUpdates(bot.Token, bot.Offset, true)
		for _, update := range updates.Result {
			bot.elaborateUpdate(update)
		}
	}
	getUpdates(bot.Token, bot.Offset, false) // Clean the last update
	log.Printf("Bot [@%s] has stopped.", bot.ThisBot.Result.Username)
}

func cleanUpdates(bot *Bot){
	bot.Offset = -1
	for true {
		updates, err := bot.getUpdates(bot.Offset, false)
		if err == TimeoutError {
			return
		}
		updateCount := len(updates.Result)
		if updateCount == 0 {
			return
		} else {
			bot.Offset = updates.Result[updateCount - 1].UpdateID + 1
			log.Printf("Cleaned to update #%d", bot.Offset)
		}
	}
}

func (bot *Bot) StartPolling(clean bool) {
	if clean {
		cleanUpdates(bot)
		log.Println("Updates cleaned!")
	}

	go bot.pollingFunction() // I'll take care of the goroutine later
}

func (bot *Bot) Idle() {
	for bot.Running {
		time.Sleep(time.Second * 1)
		// Function to wait while the bot executes
	}
}

func (bot *Bot) SendMessage(chatID int64, text string, args SendMessageArgs) (SendMessageResult, int) {
	return sendMessage(bot.Token, chatID, text, args.ParseMode, args.DisableNotification, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SetChatTitle(chatID int64, title string) {
	setChatTitle(bot.Token, chatID, title)
}

func (bot *Bot) SendChatAction(chatID int64, action string) (BooleanResult, int) {
	return sendChatAction(bot.Token, chatID, action)
}

func (bot *Bot) SendPhoto(chatID int64, fileName string, args SendPhotoArgs) {
	sendPhotoFromFile(bot.Token, chatID, fileName, args.Caption, args.ParseMode, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SendPhotoBytes(chatID int64, photoBytes []byte, args SendPhotoArgs) {
	sendPhotoFromBytes(bot.Token, chatID, photoBytes, args.Caption, args.ParseMode, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SendDocument(chatID int64, fileName string, args SendDocumentArgs) {
	sendDocumentFromFile(bot.Token, chatID, fileName, args.Caption, args.ParseMode, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SendDocumentBytes(chatID int64, fileBytes []byte, args SendDocumentArgs) {
	sendDocumentFromBytes(bot.Token, chatID, fileBytes, args.Caption, args.ParseMode, args.DisableNotification, args.ReplyToMessageID)
}