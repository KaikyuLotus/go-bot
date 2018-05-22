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
	disableNotification bool, replyToMessageId int, replyMarkup InlineKeyboardMarkup) (SendMessageResult, *RequestsError) {
	var sendMessageResult = SendMessageResult{}
	kwargs := make(map[string]string)
	kwargs["disable_notification"] = strconv.FormatBool(disableNotification)
	kwargs["disable_web_page_preview"] = strconv.FormatBool(disableWebPagePreview)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["text"] = text
	kwargs["parse_mode"] = getParseMode(parseMode)
	if replyToMessageId != 0 {
		kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageId)
	}

	if replyMarkup.InlineKeyboard != nil {
		b, e := json.Marshal(replyMarkup)
		if e != nil {
			fmt.Println(e)
			return sendMessageResult, nil
		}

		kwargs["reply_markup"] = string(b)
		fmt.Println(kwargs["reply_markup"])
	}

	result, err := MakeRequest(baseUrl+"bot"+botToken+"/sendMessage", nil, kwargs)
	if err != nil {
		return sendMessageResult, err
	}
	toApiResult(result, &sendMessageResult)
	return sendMessageResult, nil
}

func answerCallbackQuery(botToken string, callbackQueryID string, text string, showAlert bool, url string, cacheTime int) (BooleanResult, *RequestsError) {
	bres := BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["callback_query_id"] = callbackQueryID
	if text != "" {
		kwargs["text"] = text
	}
	if showAlert {
		kwargs["show_alert"] = "true"
	}
	if url != "" {
		kwargs["url"] = url
	}
	if cacheTime != 0 {
		kwargs["cache_time"] = strconv.Itoa(int(cacheTime))
	}
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/answerCallbackQuery", nil, kwargs)
	toApiResult(res, &bres)
	return bres, err
}

func setChatTitle(botToken string, chatID int64, title string) (BooleanResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["title"] = title
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/setChatTitle", nil, kwargs)
	return BooleanResult{}, err
}

func setChatDescription(botToken string, chatID int64, description string) (BooleanResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["description"] = description
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/setChatDescription", nil, kwargs)
	return BooleanResult{}, err
}

func pinChatMessage(botToken string, chatID int64, messageID int, disableNotification bool) (BooleanResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["message_id"] = strconv.Itoa(int(messageID))
	if disableNotification {
		kwargs["disable_notification"] = "true"
	}
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/pinChatMessage", nil, kwargs)
	return BooleanResult{}, err
}

func unpinChatMessage(botToken string, chatID int64) (BooleanResult, *RequestsError) {
	bres := BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/unpinChatMessage", nil, kwargs)
	toApiResult(res, &bres)
	return bres, err
}

func kickChatMember(botToken string, chatID int64, userID int, untilDate int64) (BooleanResult, *RequestsError) {
	bres := BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["user_id"] = strconv.Itoa(int(userID))
	kwargs["until_date"] = strconv.Itoa(int(untilDate))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/kickChatMember", nil, kwargs)
	toApiResult(res, &bres)
	return bres, err
}

func unbanChatMember(botToken string, chatID int64, userID int) (BooleanResult, *RequestsError) {
	bres := BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["user_id"] = strconv.Itoa(int(userID))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/unbanChatMember", nil, kwargs)
	toApiResult(res, &bres)
	return bres, err
}

func getUserProfilePhotos(botToken string, userID int, offset int, limit int) (GetUserProfilePhotosResult, *RequestsError) {
	sres := GetUserProfilePhotosResult{}
	kwargs := make(map[string]string)
	kwargs["user_id"] = strconv.Itoa(int(userID))
	if limit > 0 {
		kwargs["limit"] = strconv.Itoa(limit)
	}
	if offset > 0 {
		kwargs["offset"] = strconv.Itoa(offset)
	}
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/getUserProfilePhotos", nil, kwargs)
	toApiResult(res, &sres)
	return sres, err
}

func getFile(botToken string, fileID string) (GetFileResult, *RequestsError) {
	fres := GetFileResult{}
	kwargs := make(map[string]string)
	kwargs["file_id"] = fileID
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/getFile", nil, kwargs)
	toApiResult(res, &fres)
	return fres, err
}

func leaveChat(botToken string, chatID int64) (BooleanResult, *RequestsError) {
	bres := BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/leaveChat", nil, kwargs)
	toApiResult(res, &bres)
	return bres, err
}

func exportChatInviteLink(botToken string, chatID int64) (StringResult, *RequestsError) {
	sres := StringResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/exportChatInviteLink", nil, kwargs)
	toApiResult(res, &sres)
	return sres, err
}

func setChatPhoto(botToken string, chatID int64, photoBytes []byte) (BooleanResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["filename"] = "chatPhoto.jpg"
	kwargs["filetype"] = "photo"
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/setChatPhoto", photoBytes, kwargs)
	return BooleanResult{}, err
}

func deleteChatPhoto(botToken string, chatID int64) (BooleanResult, *RequestsError) {
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	_, err := MakeRequest(baseUrl+"bot"+botToken+"/deleteChatPhoto", nil, kwargs)
	return BooleanResult{}, err
}

func getChat(botToken string, chatID int64) (GetChatResult, *RequestsError) {
	gres := GetChatResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/getChat", nil, kwargs)
	toApiResult(res, &gres)
	return gres, err
}

func getChatAdministrators(botToken string, chatID int64) (GetChatAdministratorsResult, *RequestsError) {
	gres := GetChatAdministratorsResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/getChatAdministrators", nil, kwargs)
	toApiResult(res, &gres)
	return gres, err
}

func getChatMembersCount(botToken string, chatID int64) (IntegerResult, *RequestsError) {
	ires := IntegerResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/getChatMembersCount", nil, kwargs)
	toApiResult(res, &ires)
	return ires, err
}

func getChatMember(botToken string, chatID int64, userID int)(GetChatMemberResult, *RequestsError){
	cres := GetChatMemberResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["user_id"] = strconv.Itoa(userID)
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/getChatMember", nil, kwargs)
	toApiResult(res, &cres)
	return cres, err
}

func setChatStickerSet(botToken string, chatID int64, stickerSetName string) (BooleanResult, *RequestsError) {
	bres := BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["sticker_set_name"] = stickerSetName
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/setChatStickerSet", nil, kwargs)
	toApiResult(res, &bres)
	return bres, err
}

func deleteChatStickerSet(botToken string, chatID int64) (BooleanResult, *RequestsError) {
	bres := BooleanResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/deleteChatStickerSet", nil, kwargs)
	toApiResult(res, &bres)
	return bres, err
}

func sendContact(botToken string, chatID int64, phoneNumber string, firstName string, lastName string, disableNotification bool,
	replyToMessageID int) (SendMessageResult, *RequestsError){
	sres := SendMessageResult{}
	kwargs := make(map[string]string)
	kwargs["chat_id"] = strconv.Itoa(int(chatID))
	kwargs["phone_number"] = phoneNumber
	kwargs["first_name"] = firstName
	kwargs["last_name"] = lastName
	if disableNotification {
		kwargs["disable_notification"] = "true"
	}
	if replyToMessageID != 0 {
		kwargs["reply_to_message_id"] = strconv.Itoa(replyToMessageID)
	}
	res, err := MakeRequest(baseUrl+"bot"+botToken+"/sendContact", nil, kwargs)
	toApiResult(res, &sres)
	return sres, err
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
