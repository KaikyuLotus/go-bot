package gobot

import "fmt"


// These function are working placeholders. i'll improve them in future
func assemblyParseMode(mode int) string {
	if mode != None {
		if mode == HTML {
			return "&parse_mode=html"
		} else if mode == Markdown {
			return "&parse_mode=markdown"
		}
	}

	return ""
}

func assemblyReplyToId(messageID int) string {
	if messageID != 0 {
		return fmt.Sprintf("&reply_to_message_id=%d", messageID)
	} else {
		return ""
	}
}

func assemblyBool(paramName string, pbool bool) string {
	if pbool {
		return "&" + paramName + "=true"
	} else {
		return ""
	}
}

func assemblyCaption(caption string) string {
	if caption != "" {
		return "&caption=" + caption
	} else {
		return ""
	}
}


func getMe(botToken string) (GetMeResult, bool) {
	var getMeResult = GetMeResult{}
	response, status := makeRequest(fmt.Sprintf("%sbot%s/getMe", baseUrl, botToken))
	return getMeResult, statusCheck(&getMeResult, response, status)
}

func getUpdates(botToken string, offset int64, timeout bool) (GetUpdateResult, bool){
	var update = GetUpdateResult{}
	response, status := makeTimeoutRequest(fmt.Sprintf("%sbot%s/getUpdates?timeout=%d&offset=%d", baseUrl, botToken, 120, offset), timeout)
	return update, statusCheck(&update, response, status)
}

func sendChatAction(botToken string, chatID int64, action string) (BooleanResult, bool) {
	var booleanResult = BooleanResult{}
	resp, status := makeRequest(fmt.Sprintf("%sbot%s/sendChatAction?chat_id=%d&action=%s", baseUrl, botToken, chatID, action))
	return booleanResult, statusCheck(&booleanResult, resp, status)
}

func sendMessage(botToken string, chatID int64, text string, parseMode int, disableWebPagePreview bool, disableNotification bool, replyToMessageId int) (SendMessageResult, bool) {
	var sendMessageResult = SendMessageResult{}
	// Working placeholder
	url := fmt.Sprintf("%sbot%s/sendMessage?chat_id=%d&text=%s", baseUrl, botToken, chatID, text)
	url += assemblyBool("disable_web_page_preview", disableWebPagePreview)
	url += assemblyBool("disable_notification", disableNotification)
	url += assemblyParseMode(parseMode)
	url += assemblyReplyToId(replyToMessageId)
	resp, status := makeRequest(url)
	return sendMessageResult, statusCheck(&sendMessageResult, resp, status)
}

func setChatTitle(botToken string, chatID int64, title string){
	url := fmt.Sprintf("%sbot%s/setChatTitle?chat_id=%d&title=%s", baseUrl, botToken, chatID, title)
	makeRequest(url)
}

func sendPhotoFromFile(botToken string, chatID int64, fileName string, caption string, parseMode int, disableNotification bool, replyToMessageId int){
	sendPhotoFromBytes(botToken, chatID, readFileBytes(fileName), caption, parseMode, disableNotification, replyToMessageId)
}

func sendPhotoFromBytes(botToken string, chatID int64, fileBytes []byte, caption string, parseMode int, disableNotification bool, replyToMessageId int){
	url := fmt.Sprintf("%sbot%s/sendPhoto?chat_id=%d", baseUrl, botToken, chatID)
	url += assemblyCaption(caption)
	url += assemblyParseMode(parseMode)
	url += assemblyReplyToId(replyToMessageId)
	url += assemblyBool("disable_notification", disableNotification)
	makePost(url, "photo", fileBytes)
}

func sendDocumentFromFile(botToken string, chatID int64, fileName string, caption string, parseMode int, disableNotification bool, replyToMessageId int){
	sendDocumentFromBytes(botToken, chatID, readFileBytes(fileName), caption, parseMode, disableNotification, replyToMessageId)
}

func sendDocumentFromBytes(botToken string, chatID int64, fileBytes []byte, caption string, parseMode int, disableNotification bool, replyToMessageId int){
	url := fmt.Sprintf("%sbot%s/sendDocument?chat_id=%d", baseUrl, botToken, chatID)
	url += assemblyCaption(caption)
	url += assemblyParseMode(parseMode)
	url += assemblyReplyToId(replyToMessageId)
	url += assemblyBool("disable_notification", disableNotification)
	makePost(url, "document", fileBytes)
}