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
func NewBot(token string) (*Bot, *RequestsError) {
	bot := &Bot{token: token, Running: true}
	_, err := bot.getMe()
	if err != nil {
		return &Bot{}, err
	}
	bot.authorized = true
	return bot, nil
}

func (bot *Bot) SetUpdateHandler(foo UpdateHandlerType) {
	bot.UpdateHandler = foo
}

func (bot *Bot) SetPanicHandler(foo PanicHandlerType) {
	bot.ErrorHandler = foo
}

func (bot *Bot) AddCommandHandler(command string, foo CommandHandlerType) {
	bot.CommandHandlers = append(bot.CommandHandlers, CommandStruct{command, foo})
}

func (bot *Bot) getMe() (GetMeResult, *RequestsError) {
	result, err := getMe(bot.token)
	if err != nil {
		return GetMeResult{}, err
	}
	bot.ID = result.Bot.ID
	bot.Username = result.Bot.Username
	bot.FirstName = result.Bot.FirstName
	bot.IsBot = result.Bot.IsBot

	return result, nil
}

// First parameter is offset and second is timeout
func (bot *Bot) getUpdates(offset int64, timeout bool) (GetUpdateResult, *RequestsError) {
	updates, err := getUpdates(bot.token, offset, timeout)
	if err != nil {
		return GetUpdateResult{}, err
	}
	return updates, nil
}

func (bot *Bot) elaborateUpdate(update Update) {
	if bot.ErrorHandler != nil {
		defer func() {
			// recover from panic if one occured. Set err to nil otherwise.
			rec := recover()
			if rec != nil {
				log.Printf("gobot has recovered from a panic caused from the last handler.")
				if bot.ErrorHandler != nil {
					bot.ErrorHandler(bot, update, rec)
				}
			}
		}()
	}
	bot.Offset = update.UpdateID + 1
	update.Message.Args = strings.Split(update.Message.Text, " ")
	for _, commandStruct := range bot.CommandHandlers {
		if strings.HasPrefix(strings.ToLower(update.Message.Text), "/"+strings.ToLower(commandStruct.Command)) {
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
		updates, _ := getUpdates(bot.token, bot.Offset, true)
		for _, update := range updates.Result {
			bot.elaborateUpdate(update)
		}
	}
	getUpdates(bot.token, bot.Offset, false) // Clean the last update
	log.Printf("Bot [@%s] has stopped.", bot.Username)
}

func cleanUpdates(bot *Bot) {
	bot.Offset = -1
	for true {
		updates, err := bot.getUpdates(bot.Offset, false)
		if err != nil {
			if err.Enum == TimeoutError {
				return
			}
		}
		updateCount := len(updates.Result)
		if updateCount == 0 {
			return
		} else {
			bot.Offset = updates.Result[updateCount-1].UpdateID + 1
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

// Wrappers for RAW functions, maybe i'll join them...
func (bot *Bot) SendMessage(chatID int64, text string, args SendMessageArgs) (SendMessageResult, *RequestsError) {
	return sendMessage(bot.token, chatID, text, args.ParseMode, args.DisableNotification, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SetChatTitle(chatID int64, title string) {
	setChatTitle(bot.token, chatID, title)
}

func (bot *Bot) SendChatAction(chatID int64, action string) (BooleanResult, *RequestsError) {
	return sendChatAction(bot.token, chatID, action)
}

func (bot *Bot) SendSticker(chatID int64, sticker string, args SendStickerArgs) (SendStickerResult, *RequestsError) {
	return sendSticker(bot.token, chatID, sticker, args.ReplyToMessageID, args.DisableNotification)
}

func (bot *Bot) ForwardMessage(chatID int64, fromChatID int64, messageID int, disableNotification bool) (ForwardMessageResult, *RequestsError) {
	return forwardMessage(bot.token, chatID, fromChatID, messageID, disableNotification)
}

// SendPhoto
func (bot *Bot) SendPhotoBytes(chatID int64, fileName string, photoBytes []byte, args SendPhotoArgs) (SendPhotoResult, *RequestsError) {
	return sendPhotoBytes(bot.token, chatID, fileName, photoBytes, args.Caption, args.ParseMode, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SendPhoto(chatID int64, fileID string, args SendPhotoArgs) (SendPhotoResult, *RequestsError) {
	return sendPhotoByID(bot.token, chatID, fileID, args.Caption, args.ParseMode, args.DisableNotification, args.ReplyToMessageID)
}

// SendDocument
func (bot *Bot) SendDocumentBytes(chatID int64, fileName string, fileBytes []byte, args SendDocumentArgs) (SendDocumentResult, *RequestsError) {
	return sendDocumentBytes(bot.token, chatID, fileName, fileBytes, args.Caption, args.ParseMode, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SendDocument(chatID int64, fileID string, args SendDocumentArgs) (SendDocumentResult, *RequestsError) {
	return sendDocumentByID(bot.token, chatID, fileID, args.Caption, args.ParseMode, args.DisableNotification, args.ReplyToMessageID)
}

// SendAudio
func (bot *Bot) SendAudioBytes(chatID int64, fileName string, fileBytes []byte, args SendAudioArgs) (SendAudioResult, *RequestsError) {
	return sendAudioBytes(bot.token, chatID, fileBytes, fileName, args.Caption, args.ParseMode, args.Duration,
		args.Performer, args.Title, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SendAudio(chatID int64, fileID string, args SendAudioArgs) (SendAudioResult, *RequestsError) {
	return sendAudioByID(bot.token, chatID, fileID, args.Caption, args.ParseMode, args.Duration,
		args.Performer, args.Title, args.DisableNotification, args.ReplyToMessageID)
}

// SendVoice
func (bot *Bot) SendVoiceBytes(chatID int64, fileName string, fileBytes []byte, args SendVoiceArgs) (SendVoiceResult, *RequestsError) {
	return sendVoiceBytes(bot.token, chatID, fileBytes, fileName, args.Caption, args.ParseMode, args.Duration, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SendVoice(chatID int64, fileID string, args SendVoiceArgs) (SendVoiceResult, *RequestsError) {
	return sendVoiceByID(bot.token, chatID, fileID, args.Caption, args.ParseMode, args.Duration, args.DisableNotification, args.ReplyToMessageID)
}