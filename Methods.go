package gobot

import (
	"fmt"
	"strconv"
)

var baseUrl = "https://api.telegram.org/"

func getParseMode(mode int) string {
	if mode != None {
		if mode == HTML {
			return "html"
		} else if mode == Markdown {
			return "markdown"
		}
	}
	return ""
}

func getMe(botToken string) (GetMeResult, int) {
	var getMeResult = GetMeResult{}
	response, status := makeRequest(baseUrl + "bot" + botToken +"/getMe", make(map[string]string))
	return getMeResult, statusCheck(&getMeResult, response, status)
}

func getUpdates(botToken string, offset int64, timeout bool) (GetUpdateResult, int){
	var update = GetUpdateResult{}
	kwargs := make(map[string]string)
	if timeout {
		kwargs["timeout"] = "120"
	}
	kwargs["offset"] = strconv.Itoa(int(offset))
	response, status := makeRequest(baseUrl + "bot" + botToken + "/getUpdates", kwargs)
	return update, statusCheck(&update, response, status)
}

func sendChatAction(botToken string, chatID int64, action string) (BooleanResult, int) {
	var booleanResult = BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["action"] = action
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	resp, status := makeRequest(baseUrl + "bot" + botToken + "/sendChatAction", kwargs)
	return booleanResult, statusCheck(&booleanResult, resp, status)
}

func sendMessage(botToken string, chatID int64, text string, parseMode int, disableWebPagePreview bool, disableNotification bool, replyToMessageId int) (SendMessageResult, int) {
	var sendMessageResult = SendMessageResult{}
	// Working placeholder
	kwargs := make(map[string]string)
	kwargs["disable_notification"] = strconv.FormatBool(disableNotification)
	kwargs["disable_web_page_preview"] = strconv.FormatBool(disableWebPagePreview)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["text"] = text
	kwargs["parse_mode"] = getParseMode(parseMode)
	kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	resp, status := makeRequest(baseUrl + "bot" + botToken + "/sendMessage", kwargs)
	return sendMessageResult, statusCheck(&sendMessageResult, resp, status)
}

func setChatTitle(botToken string, chatID int64, title string){
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["title"] = title
	makeRequest(baseUrl + "bot" + botToken + "/setChatTitle", kwargs)
}

func sendPhotoFromFile(botToken string, chatID int64, fileName string, caption string, parseMode int, disableNotification bool, replyToMessageId int){
	sendPhotoFromBytes(botToken, chatID, ReadFileBytes(fileName), caption, parseMode, disableNotification, replyToMessageId)
}

// ToDo: Fix POST url composition
func sendPhotoFromBytes(botToken string, chatID int64, fileBytes []byte, caption string, parseMode int, disableNotification bool, replyToMessageId int){
	url := fmt.Sprintf("%sbot%s/sendPhoto?chat_id=%d", baseUrl, botToken, chatID)
	makePost(url, "photo", fileBytes) // later
}

func sendDocumentFromFile(botToken string, chatID int64, fileName string, caption string, parseMode int, disableNotification bool, replyToMessageId int){
	sendDocumentFromBytes(botToken, chatID, ReadFileBytes(fileName), caption, parseMode, disableNotification, replyToMessageId)
}

func sendDocumentFromBytes(botToken string, chatID int64, fileBytes []byte, caption string, parseMode int, disableNotification bool, replyToMessageId int){
	url := fmt.Sprintf("%sbot%s/sendDocument?chat_id=%d", baseUrl, botToken, chatID)
	makePost(url, "document", fileBytes) // later
}