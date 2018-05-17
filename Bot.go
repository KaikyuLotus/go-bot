package gobot

// I need to put return types on every foo in the future

// Bot constructor
func NewBot(token string) *Bot {
	bot := &Bot{Token: token}
	bot.getMe()
	return bot
}

func (bot *Bot) SetUpdateHandler(foo UpdateHandlerType) {
	bot.UpdateHandler = foo
}

func (bot *Bot) AddCommandHandler(command string, foo CommandHandlerType){
	bot.CommandHandlers = append(bot.CommandHandlers, CommandStruct{command, foo})
}

func (bot *Bot) getMe() GetMeResult {
	result, getMeOk := getMe(bot.Token)
	if !getMeOk { panic(getMeOk) } // For now i'll just panic here if unauthorized
	bot.ThisBot = result
	return result
}

// First parameter is offset and second is timeout
func (bot *Bot) getUpdates(values ...int64) (GetUpdateResult, bool) {
	updates := GetUpdateResult{}
	argCount := len(values)
	if argCount > 1 {
		return updates, false
	}

	var offset int64 = 0

	if argCount != 0 {
		if argCount == 1 {
			offset = values[0]
		}
	}

	return getUpdates(bot.Token, offset, true)
}

func (bot *Bot) elaborateUpdate(update Update){
	bot.Offset = update.UpdateID + 1
	for _, commandStruct := range bot.CommandHandlers {
		if update.Message.Text ==  "/" + commandStruct.Command {
			commandStruct.Function(*bot, update)
			return
		}
	}
	bot.UpdateHandler(*bot, update)
}

func (bot *Bot) StartPolling(clean bool){
	for true {
		updates, _ := getUpdates(bot.Token, bot.Offset, true)
		for _, update := range updates.Result {
			bot.elaborateUpdate(update)
		}
	}
}

func (bot *Bot) SendMessage(chatID int64, text string, args SendMessageArgs) (SendMessageResult, bool) {
	return sendMessage(bot.Token, chatID, text, args.ParseMode, args.DisableNotification, args.DisableNotification, args.ReplyToMessageID)
}

func (bot *Bot) SetChatTitle(chatID int64, title string){
	setChatTitle(bot.Token, chatID, title)
}

func (bot *Bot) SendChatAction(chatID int64, action string) (BooleanResult, bool) {
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