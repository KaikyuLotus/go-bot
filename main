package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
)

const (
	TelegramError = iota
	ResponseError
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
var botToken = "495913948:AAHrir-w5aaymPPpPyrp7pYiZqL81DlBA80"

func main(){
	result, err := sendMessage(487353090, "Hello World")

	if (ApiError{}) != err {
		fmt.Println("Excecution failed")
		fmt.Printf("Telegram has returned error '%s' with status code '%d'.", err.Description, err.ErrorCode)
		return
	}

	if result == (SendMessageResult{}) {
		fmt.Println("Excecution failed, network error?")
		return
	}

	fmt.Println(result.Result.Text + " sent.")
	fmt.Println("Execution ok!")
}

func getMe() (GetMeResult, ApiError) {
	// Always return an empty error
	var apiError = ApiError{}
	var result = GetMeResult{}

	var resp, status = makeRequest(baseUrl + "bot" + botToken + "/" + "getMe")
	statusCheck(&result, &apiError, resp, status)

	return result, apiError
}

func sendMessage(chatID int, text string) (SendMessageResult, ApiError) {
	// Always return an empty error
	var apiError = ApiError{}
	var result = SendMessageResult{}

	var resp, status = makeRequest(baseUrl + "bot" + botToken + "/" + "sendMessage?chat_id=487353090&text=" + text)
	statusCheck(&result, &apiError, resp, status)

	return result, apiError
}

func statusCheck(result interface{}, apiError *ApiError, resp *http.Response, status int){
	if status != ResponseError {
		if status == TelegramError {
			// Compile the error
			json.NewDecoder(resp.Body).Decode(&apiError)
		} else {
			// Use json.Decode for reading streams of JSON data
			// body, _ := ioutil.ReadAll(resp.Body)
			// result := string(body)
			// fmt.Println(result)
			json.NewDecoder(resp.Body).Decode(&result)
		}
	}

	// Defer the closing of the body
	defer resp.Body.Close()
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
