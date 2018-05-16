package main

type CommandStruct struct {
	Command string
	Function CommandHandlerType
}

type UpdateHandlerType func(Bot, Update)

type CommandHandlerType func(Bot, Update)

type Bot struct {
	ThisBot GetMeResult
	Token string
	UpdateHandler UpdateHandlerType
	CommandHandlers []CommandStruct
	Offset int64
}

// Bot constructor
func NewBot(token string) *Bot {
	bot := &Bot{Token: token}
	bot.getMe()
	return bot
}

func (bot *Bot) setUpdateHandler(foo UpdateHandlerType) {
	bot.UpdateHandler = foo
}

func (bot *Bot) addCommandHandler(command string, foo CommandHandlerType){
	bot.CommandHandlers = append(bot.CommandHandlers, CommandStruct{command, foo})
}

func (bot *Bot) getMe() GetMeResult {
	result, getMeOk := getMe()
	if !getMeOk { panic(getMeOk) }
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

	return getUpdates(offset, true)
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

func (bot *Bot) startPolling(clean bool){
	for true {
		updates, _ := getUpdates(bot.Offset, true)
		for _, update := range updates.Result {
			bot.elaborateUpdate(update)
		}
	}
}

func (bot *Bot) sendMessage(chatID int64, text string) (SendMessageResult, bool) {
	return sendMessage(chatID, text)
}

func (bot *Bot) sendChatAction(chatID int64, action string) (BooleanResult, bool) {
	return sendChatAction(chatID, action)
}


