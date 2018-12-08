package gobot

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

/*
I know that Go is not an object-oriented language but i think that this approach is not that bad
So please don't kill me <3
*/

// I need to put return types on every foo in the future
// Bot constructor
func NewBot(token string) (*Bot, *RequestsError) {
	bot := &Bot{
		token:   token,
		Running: true,
	}

	_, err := bot.getMe()
	if err != nil {
		return &Bot{}, err
	}
	bot.authorized = true
	return bot, nil
}

func (bot *Bot) SetUpdateHandler(foo UpdateHandlerType) *Bot {
	bot.UpdateHandler = foo
	return bot
}

func (bot *Bot) SetCallbackQueryHandler(foo CallbackHandlerType) *Bot {
	bot.CallbackQueryHandler = foo
	return bot
}

func (bot *Bot) SetPanicHandler(foo PanicHandlerType) *Bot {
	bot.ErrorHandler = foo
	return bot
}

func (bot *Bot) AddCommandHandler(command string, foo CommandHandlerType) {
	bot.CommandHandlers = append(bot.CommandHandlers, CommandStruct{strings.ToLower(command), foo})
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

func (bot *Bot) pushUpdateHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Thanks!"))
}

func (bot *Bot) webhookUpdateHandler(rw http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)

	var update Update
	json.Unmarshal(body, &update)

	bot.elaborateUpdate(update)

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Thanks!"))
}

func (bot *Bot) elaborateUpdate(update Update) {
	var uType string
	if update.Message.MessageID != 0 {
		uType = "message"
	} else if update.CallbackQuery.ID != "" {
		uType = "callback_query"
	}

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

	if uType == "message" {
		// Thanks to https://gist.github.com/sk22/cc02d95cd2d24c882835c1dddb33e1da#file-telegramcmd-regex
		regex, _ := regexp.Compile("(?i)^/([^@\\s]+)@?(?:(\\S+)|)\\s?([\\s\\S]*)$")
		for _, commandStruct := range bot.CommandHandlers {
			res := regex.FindStringSubmatch(update.Message.Text)
			// Check if the regex has no matches
			if len(res) == 0 {
				break // this is not a /command
			}

			if res[2] != "" { // is it /command@username ?
				if strings.ToLower(res[2]) != strings.ToLower(bot.Username) {
					break // This isn't out bot's username
				}
			}
			// Is it our command?
			if strings.ToLower(res[1]) == commandStruct.Command {
				update.Message.Args = strings.Split(res[3], " ") // Split args
				commandStruct.Function(bot, update)
				return // we're done with this update
			}
		}
		// Fix for issue #1
		if bot.UpdateHandler == nil {
			return
		}
		bot.UpdateHandler(bot, update)
		return
	}

	if uType == "callback_query" {
		query := update.CallbackQuery
		query.Message.Args = strings.Split(query.Message.Text, " ")

		if bot.CallbackQueryHandler == nil {
			return
		}
		bot.CallbackQueryHandler(bot, query)
		return
	}

	log.Printf("Unknown update type...")
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

func (bot *Bot) StartWebhook(url string, port int, certificateFileName string,
	serverKeyFileName string, clean bool) {

	if clean {
		cleanUpdates(bot)
		log.Println("Updates cleaned!")
	}

	bot.DeleteWebhook()
	bot.SetWebhoook(
		url, port, "update",
		SetWebhookArgs{Certificate: certificateFileName},
	)

	router := http.NewServeMux()
	router.Handle("/push", http.HandlerFunc(bot.pushUpdateHandler))
	router.Handle("/update", http.HandlerFunc(bot.webhookUpdateHandler))

	bot.server = &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   router,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	go bot.server.ListenAndServeTLS(certificateFileName, serverKeyFileName)
}

func (bot *Bot) StartPolling(clean bool) *Bot {
	if clean {
		cleanUpdates(bot)
		log.Println("Updates cleaned!")
	}

	go bot.pollingFunction() // I'll take care of the goroutine later
	return bot
}

func (bot *Bot) Idle() {
	for bot.Running {
		time.Sleep(time.Second * 1)
		// Function to wait while the bot executes
	}
}

// Wrappers for RAW functions, maybe i'll join them...

func (bot *Bot) EditMessageText(text string, args EditMessageArgs) (BooleanResult, *RequestsError) {
	return editMessageText(bot.token, text, args.ChatID, args.MessageID, args.InlineMessageID, args.ParseMode, args.DisableWebPagePreview, args.ReplyMarkup)
}

func (bot *Bot) SendMessage(chatID int64, text string, args SendMessageArgs) (SendMessageResult, *RequestsError) {
	return sendMessage(bot.token, chatID, text, args.ParseMode, args.DisableNotification, args.DisableNotification, args.ReplyToMessageID, args.ReplyMarkup)
}

func (bot *Bot) AnswerCallbackQuery(callbackQueryID string, args AnswerCallbackQueryArgs) (BooleanResult, *RequestsError) {
	return answerCallbackQuery(bot.token, callbackQueryID, args.Text, args.ShowAlert, args.Url, args.CacheTime)
}

func (bot *Bot) SetChatTitle(chatID int64, title string) (BooleanResult, *RequestsError) {
	return setChatTitle(bot.token, chatID, title)
}

func (bot *Bot) SetChatDescription(chatID int64, description string) (BooleanResult, *RequestsError) {
	return setChatDescription(bot.token, chatID, description)
}

func (bot *Bot) PinChatMessage(chatID int64, messageID int, disableNotification bool) (BooleanResult, *RequestsError) {
	return pinChatMessage(bot.token, chatID, messageID, disableNotification)
}

func (bot *Bot) UnpinChatMessage(chatID int64) (BooleanResult, *RequestsError) {
	return unpinChatMessage(bot.token, chatID)
}

func (bot *Bot) KickChatMember(chatID int64, userID int, untilDate int64) (BooleanResult, *RequestsError) {
	return kickChatMember(bot.token, chatID, userID, untilDate)
}

func (bot *Bot) UnbanChatMember(chatID int64, userID int) (BooleanResult, *RequestsError) {
	return unbanChatMember(bot.token, chatID, userID)
}

func (bot *Bot) ExportChatInviteLink(chatID int64) (StringResult, *RequestsError) { // (BooleanResult, *RequestsError) {
	return exportChatInviteLink(bot.token, chatID)
}

func (bot *Bot) GetUserProfilePhotos(userID int, args GetUserProfilePhotosArgs) (GetUserProfilePhotosResult, *RequestsError) {
	return getUserProfilePhotos(bot.token, userID, args.Offset, args.Limit)
}

func (bot *Bot) GetFile(fileID string) (GetFileResult, *RequestsError) {
	return getFile(bot.token, fileID)
}

func (bot *Bot) LeaveChat(chatID int64) (BooleanResult, *RequestsError) {
	return leaveChat(bot.token, chatID)
}

func (bot *Bot) GetChatAdministrators(chatID int64) (GetChatAdministratorsResult, *RequestsError) {
	return getChatAdministrators(bot.token, chatID)
}

func (bot *Bot) GetChatMembersCount(chatID int64) (IntegerResult, *RequestsError) {
	return getChatMembersCount(bot.token, chatID)
}

func (bot *Bot) getChatMember(chatID int64, userID int) (GetChatMemberResult, *RequestsError) {
	return getChatMember(bot.token, chatID, userID)
}

func (bot *Bot) SetChatStickerSet(chatID int64, stickerSetName string) (BooleanResult, *RequestsError) {
	return setChatStickerSet(bot.token, chatID, stickerSetName)
}

func (bot *Bot) deleteChatStickerSet(chatID int64) (BooleanResult, *RequestsError) {
	return deleteChatStickerSet(bot.token, chatID)
}

func (bot *Bot) SendContact(chatID int64, phoneNumber string, firstName string, contact SendContactArgs) (SendMessageResult, *RequestsError) {
	return sendContact(bot.token, chatID, phoneNumber, firstName, contact.LastName, contact.DisableNotification, contact.ReplyToMessageID)
}

func (bot *Bot) SetChatPhoto(chatID int64, photoBytes []byte) (BooleanResult, *RequestsError) {
	return setChatPhoto(bot.token, chatID, photoBytes)
}

func (bot *Bot) DeleteChatPhoto(chatID int64) (BooleanResult, *RequestsError) {
	return deleteChatPhoto(bot.token, chatID)
}

func (bot *Bot) GetChat(chatID int64) (GetChatResult, *RequestsError) {
	return getChat(bot.token, chatID)
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

func (bot *Bot) SetWebhoook(url string, port int, path string, args SetWebhookArgs) (BooleanResult, *RequestsError) {
	return setWebhook(bot.token, url, port, path, args.Certificate, args.MaxConnections, args.AllowedUpdates)
}

func (bot *Bot) DeleteWebhook() (BooleanResult, *RequestsError) {
	return deleteWebhook(bot.token)
}

func (bot *Bot) GetWebhookInfo() (GetWebhookInfoResult, *RequestsError) {
	return getWebhookInfo(bot.token)
}
