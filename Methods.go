package gobot

import (
	"fmt"
	"strconv"
	"net/http"
	"encoding/json"
)

var baseUrl = "https://api.telegram.org/"

func toApiResult(resp *http.Response, outStruct interface{}) {
	json.NewDecoder(resp.Body).Decode(outStruct)
}

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

// Returns an empty struct if an error happens and the error, otherwise returns the result and nil
func getMe(botToken string) (GetMeResult, *RequestsError) {
	var getMeResult = GetMeResult{}
	response, err := makeRequest(baseUrl + "bot" + botToken +"/getMe", make(map[string]string))
	if err != nil {
		return getMeResult, err
	}
	toApiResult(response, &getMeResult)
	return getMeResult, nil
}

func getUpdates(botToken string, offset int64, timeout bool) (GetUpdateResult, *RequestsError){
	var updates = GetUpdateResult{}
	kwargs := make(map[string]string)
	if timeout {
		kwargs["timeout"] = "120"
	}
	kwargs["offset"] = strconv.Itoa(int(offset))
	response, err := makeRequest(baseUrl + "bot" + botToken + "/getUpdates", kwargs)
	if err != nil {
		return updates, err
	}
	toApiResult(response, &updates)
	return updates, nil // statusCheck(&update, response, status)
}

func sendChatAction(botToken string, chatID int64, action string) (BooleanResult, *RequestsError) {
	var booleanResult = BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["action"] = action
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	_, err := makeRequest(baseUrl + "bot" + botToken + "/sendChatAction", kwargs)
	return booleanResult, err// statusCheck(&booleanResult, resp, status)
}

func sendMessage(botToken string, chatID int64, text string, parseMode int, disableWebPagePreview bool, disableNotification bool, replyToMessageId int) (SendMessageResult, *RequestsError) {
	var sendMessageResult = SendMessageResult{}
	// Working placeholder
	kwargs := make(map[string]string)
	kwargs["disable_notification"] = strconv.FormatBool(disableNotification)
	kwargs["disable_web_page_preview"] = strconv.FormatBool(disableWebPagePreview)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["text"] = text
	kwargs["parse_mode"] = getParseMode(parseMode)
	kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	result, err := makeRequest(baseUrl + "bot" + botToken + "/sendMessage", kwargs)
	if err != nil {
		return sendMessageResult, err
	}
	toApiResult(result, &sendMessageResult)
	return sendMessageResult, nil
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

func sendSticker(botToken string, chatID int64, fileID string, replyToMessageId int, disableNotification bool) (SendStickerResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["sticker"] = fileID
	kwargs["disable_notification"] = strconv.FormatBool(disableNotification)
	kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	makeRequest(baseUrl + "bot" + botToken + "/sendSticker", kwargs)
	return SendStickerResult{}, nil

}