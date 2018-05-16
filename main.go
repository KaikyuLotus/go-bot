package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
)

const (
	ResponseError = iota
	RequestNotOk
	RequestOk
)

// {"ok":true,"result":{"id":495913948,"is_bot":true,"first_name":"Icy","username":"Icyppbot"}}
// https://mholt.github.io/json-to-go/
type SendMessageResult struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

type GetMeResult struct {
	Ok     bool `json:"ok"`
	Result struct {
		ID        int    `json:"id"`
		IsBot     bool   `json:"is_bot"`
		FirstName string `json:"first_name"`
		Username  string `json:"username"`
	} `json:"result"`
}

type ApiError struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

var baseUrl = "https://api.telegram.org/"
var botToken = "495913948:AAE-dymFTf_Sj5pxpR4KNf1GwMgcjvMYxwI"

func main(){
	bot, getMeOk := getMe()
	if getMeOk {
		println(bot.Result.Username)
	}

	sentMessage, ok := sendMessage(487353090, "Hello World")
	if ok{
		fmt.Println(sentMessage.Result.Text + " sent.")
	}
}

func getMe() (GetMeResult, bool) {
	response, status := makeRequest(fmt.Sprintf("%sbot%s/getMe", baseUrl, botToken))
	result, ok := statusCheck(GetMeResult{}, response, status)
	getMeResult, convOk := result.(GetMeResult)
	if !convOk {
		println("Casting failed")
	}
	return getMeResult, ok
}

func sendMessage(chatID int, text string) (SendMessageResult, bool) {
	resp, status := makeRequest(fmt.Sprintf("%sbot%s/sendMessage?chat_id=%d&text=%s", baseUrl, botToken, chatID, text))
	result, ok := statusCheck(SendMessageResult{}, resp, status)
	sendMessageResult, _ := result.(SendMessageResult)
	return sendMessageResult, ok
}

func statusCheck(result interface{}, resp *http.Response, status int) (interface{}, bool) {
	var old = result
	var apiError = ApiError{}

	if status != ResponseError {
		if status == RequestNotOk {
			json.NewDecoder(resp.Body).Decode(&apiError)
		} else {
			json.NewDecoder(resp.Body).Decode(&result)
		}
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	if (ApiError{}) != apiError {
		fmt.Println("Excecution failed")
		fmt.Printf("Telegram has returned error '%s' with status code '%d'.", apiError.Description, apiError.ErrorCode)
		return result, false
	}

	if result == old {
		fmt.Println("Excecution failed, network error?")
		return result, false
	}

	return result, true
}

func makeRequest(url string) (*http.Response, int)  {
	// safePhone := url.QueryEscape(phone)
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, ResponseError
	}

	client := &http.Client{}
	client.Timeout = time.Duration(120 * time.Second)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, ResponseError
	}

	if resp.StatusCode != 200 {
		fmt.Println("Status code not 200!")
		// body, _ := ioutil.ReadAll(resp.Body)
		// result := string(body)
		return resp, RequestNotOk
	}

	return resp, RequestOk
}
