package gobot

import (
	"strconv"
	"encoding/json"
	"io"
	"log"
	"fmt"
)

var baseUrl = "https://api.telegram.org/"

// ToDO: Convert every return value to its proper interface
func toApiResult(resp io.Reader, outStruct interface{}) {
	json.NewDecoder(resp).Decode(outStruct)
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
	response, err := MakeRequest(baseUrl+"bot"+botToken+"/getMe", nil, make(map[string]string))
	if err != nil {
		return getMeResult, err
	}
	toApiResult(response, &getMeResult)
	return getMeResult, nil
}

func getUpdates(botToken string, offset int64, timeout bool) (GetUpdateResult, *RequestsError) {
	var updates = GetUpdateResult{}
	kwargs := make(map[string]string)
	if timeout {
		kwargs["timeout"] = "120"
	}
	kwargs["offset"] = strconv.Itoa(int(offset))
	response, err := MakeRequest(baseUrl+"bot"+botToken+"/getUpdates", nil, kwargs)
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
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/sendChatAction", nil, kwargs)
	return booleanResult, err // statusCheck(&booleanResult, resp, status)
}

func sendMessage(botToken string, chatID int64, text string, parseMode int, disableWebPagePreview bool,
	disableNotification bool, replyToMessageId int) (SendMessageResult, *RequestsError) {
	var sendMessageResult = SendMessageResult{}
	// Working placeholder
	kwargs := make(map[string]string)
	kwargs["disable_notification"] = strconv.FormatBool(disableNotification)
	kwargs["disable_web_page_preview"] = strconv.FormatBool(disableWebPagePreview)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["text"] = text
	kwargs["parse_mode"] = getParseMode(parseMode)
	if replyToMessageId != 0 {
		kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	}
	result, err := MakeRequest(baseUrl+"bot"+botToken+"/sendMessage", nil, kwargs)
	if err != nil {
		return sendMessageResult, err
	}
	toApiResult(result, &sendMessageResult)
	return sendMessageResult, nil
}

func setChatTitle(botToken string, chatID int64, title string) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["title"] = title
	MakeRequest(baseUrl+"bot"+botToken+"/setChatTitle", nil, kwargs)
}

func sendSticker(botToken string, chatID int64, fileID string, replyToMessageId int, disableNotification bool) (SendStickerResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["sticker"] = fileID
	if disableNotification {
		kwargs["disable_notification"] = "true"
	}
	if replyToMessageId != 0 {
		kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	}
	MakeRequest(baseUrl+"bot"+botToken+"/sendSticker", nil, kwargs)
	return SendStickerResult{}, nil
}

func forwardMessage(botToken string, chatID int64, fromChatID int64, messageID int, disableNotification bool) (ForwardMessageResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["from_chat_id"] = strconv.Itoa(int(fromChatID))
	kwargs["message_id"] = strconv.Itoa(int(messageID))
	if disableNotification {
		kwargs["disable_notification"] = "true"
	}
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/forwardMessage", nil, kwargs)
	return ForwardMessageResult{}, err
}

// This is a bit of a mess... But the final user won't even see this
func sendPhotoBytes(botToken string, chatID int64, fileName string, fileBytes []byte, caption string, parseMode int,
	disableNotification bool, replyToMessageId int) (SendPhotoResult, *RequestsError) {
	// url := fmt.Sprintf("%sbot%s/sendPhoto?chat_id=%d", baseUrl, botToken, chatID)
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	if fileBytes != nil {
		kwargs["filename"] = fileName
		kwargs["filetype"] = "photo"
	} else {
		kwargs["photo"] = fileName
	}
	kwargs["caption"] = caption
	kwargs["parse_mode"] = getParseMode(parseMode)
	if disableNotification {
		kwargs["disable_notification"] = "true"
	}
	if replyToMessageId != 0 {
		kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	}
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/sendPhoto", fileBytes, kwargs)
	return SendPhotoResult{}, err
}

func sendPhotoByID(botToken string, chatID int64, fileID string, caption string, parseMode int,
	disableNotification bool, replyToMessageId int) (SendPhotoResult, *RequestsError) {
	sendPhotoResult := SendPhotoResult{}
	_, err := sendPhotoBytes(botToken, chatID, fileID, nil, caption, parseMode, disableNotification, replyToMessageId)
	if err != nil {
		log.Println(fmt.Sprintf("Error while sending by ID (%s): %s", fileID, err.Cause))
	}
	return sendPhotoResult, err
}

func sendDocumentBytes(botToken string, chatID int64, fileName string, fileBytes []byte, caption string,
	parseMode int, disableNotification bool, replyToMessageId int) (SendDocumentResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	if fileBytes != nil {
		kwargs["filename"] = fileName
		kwargs["filetype"] = "document"
	} else {
		kwargs["document"] = fileName
	}
	kwargs["caption"] = caption
	kwargs["parse_mode"] = getParseMode(parseMode)
	if disableNotification {
		kwargs["disable_notification"] = "true"
	}
	if replyToMessageId != 0 {
		kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	}
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/sendDocument", fileBytes, kwargs)
	return SendDocumentResult{}, err
}

func sendDocumentByID(botToken string, chatID int64, fileID string, caption string, parseMode int,
	disableNotification bool, replyToMessageId int) (SendDocumentResult, *RequestsError) {
	return sendDocumentBytes(botToken, chatID, fileID, nil, caption, parseMode, disableNotification, replyToMessageId)
}

func sendAudioBytes(botToken string, chatID int64, fileBytes []byte, fileName string, caption string, parseMode int,
	duration int, performer string, title string, disableNotification bool, replyToMessageId int) (SendAudioResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	if fileBytes != nil {
		kwargs["filename"] = fileName
		kwargs["filetype"] = "audio"
	} else {
		kwargs["audio"] = fileName
	}
	kwargs["caption"] = caption
	kwargs["parse_mode"] = getParseMode(parseMode)
	kwargs["performer"] = performer
	kwargs["title"] = title
	if disableNotification {
		kwargs["disable_notification"] = "true"
	}
	if replyToMessageId != 0 {
		kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	}
	if duration != 0 {
		kwargs["duration"] = strconv.Itoa(duration)
	}
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/sendAudio", fileBytes, kwargs)

	return SendAudioResult{}, err
}

func sendAudioByID(botToken string, chatID int64, fileID string, caption string, parseMode int,
	duration int, performer string, title string, disableNotification bool, replyToMessageId int) (SendAudioResult, *RequestsError) {
	sendAudioResult := SendAudioResult{}
	_, err := sendAudioBytes(botToken, chatID, nil, fileID, caption, parseMode, duration, performer, title, disableNotification, replyToMessageId)
	return sendAudioResult, err
}

func sendVoiceBytes(botToken string, chatID int64, fileBytes []byte, fileName string, caption string, parseMode int,
	duration int, disableNotification bool, replyToMessageId int) (SendVoiceResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	if fileBytes != nil {
		kwargs["filename"] = fileName
		kwargs["filetype"] = "voice"
	} else {
		kwargs["voice"] = fileName
	}
	kwargs["caption"] = caption
	kwargs["parse_mode"] = getParseMode(parseMode)
	if disableNotification {
		kwargs["disable_notification"] = "true"
	}
	if replyToMessageId != 0 {
		kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	}
	if duration != 0 {
		kwargs["duration"] = strconv.Itoa(duration)
	}
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/sendVoice", fileBytes, kwargs)

	return SendVoiceResult{}, err
}

func sendVoiceByID(botToken string, chatID int64, fileID string, caption string, parseMode int,
	duration int, disableNotification bool, replyToMessageId int) (SendVoiceResult, *RequestsError) {
	sendVoiceResult := SendVoiceResult{}
	_, err := sendVoiceBytes(botToken, chatID, nil, fileID, caption, parseMode, duration, disableNotification, replyToMessageId)
	return sendVoiceResult, err
}
