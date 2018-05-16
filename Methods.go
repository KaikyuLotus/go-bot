package main

import (
	"fmt"
)

const (
	Typing = "typing"
	UploadPhoto = "upload_photo"
	UploadVideo = "upload_video"
)

func getMe() (GetMeResult, bool) {
	var getMeResult = GetMeResult{}
	response, status := makeRequest(fmt.Sprintf("%sbot%s/getMe", baseUrl, botToken))
	return getMeResult, statusCheck(&getMeResult, response, status)
}

func getUpdates(offset int64, timeout bool) (GetUpdateResult, bool){
	var update = GetUpdateResult{}
	response, status := makeTimeoutRequest(fmt.Sprintf("%sbot%s/getUpdates?timeout=%d&offset=%d", baseUrl, botToken, 120, offset), timeout)
	return update, statusCheck(&update, response, status)
}

func sendMessage(chatID int64, text string) (SendMessageResult, bool) {
	var sendMessageResult = SendMessageResult{}
	resp, status := makeRequest(fmt.Sprintf("%sbot%s/sendMessage?chat_id=%d&text=%s", baseUrl, botToken, chatID, text))
	return sendMessageResult, statusCheck(&sendMessageResult, resp, status)
}

func sendChatAction(chatID int64, action string) (BooleanResult, bool) {
	var booleanResult = BooleanResult{}
	resp, status := makeRequest(fmt.Sprintf("%sbot%s/sendChatAction?chat_id=%d&action=%s", baseUrl, botToken, chatID, action))
	return booleanResult, statusCheck(&booleanResult, resp, status)
}