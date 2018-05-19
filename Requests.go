package gobot

import (
	"net/http"
	"net/url"
	"fmt"
	"time"
	"net"
	"bytes"
	"mime/multipart"
)


var client = &http.Client{}

/*func statusCheck(result interface{}, resp *http.Response, status int) int {
	var old = &result
	var apiError = ApiError{}

	if status != ResponseError {
		if status == RequestNotOk {
			json.NewDecoder(resp.Body).Decode(&apiError)
			fmt.Println("Excecution failed")
			fmt.Printf("Telegram has returned error '%s' with status code '%d'.", apiError.Description, apiError.ErrorCode)
			return RequestNotOk
		} else {
			json.NewDecoder(resp.Body).Decode(result)
		}

		// Defer the closing of the body
		defer resp.Body.Close()

		if result == old {
			fmt.Println("Excecution failed, network error?")
			return RequestNotOk
		}

		return RequestOk
	} else {
		return ResponseError
	}
}*/

func makeRequest(urll string, kwargs map[string]string) (*http.Response, *RequestsError) {
	// Build the request
	u, _ := url.Parse(urll)
	q := u.Query()
	for key, value := range kwargs {
		if value == "" {
			continue
		}
		q.Add(key, value)
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, &RequestsError{Enum:RequestNotOk, Args:kwargs, Url:urll, Cause:fmt.Sprintf("Error in the URL '%s'", u.String())}
	}

	client.Timeout = time.Duration(time.Second * 122)

	resp, err := client.Do(req)
	if err != nil {
		switch err := err.(type) {
		case net.Error:
			if err.Timeout() {
				return nil, &RequestsError{Enum:TimeoutError, Args:kwargs, Url:urll, Cause:fmt.Sprintf("Timeout exceeded (%s)", err)}
			}
		case *url.Error:
			if err, ok := err.Err.(net.Error); ok && err.Timeout() {
				return nil, &RequestsError{Enum:TimeoutError, Args:kwargs, Url:urll, Cause:fmt.Sprintf("Timeout exceeded (%s)", err)}
			}
		}
		return nil, &RequestsError{Enum:ResponseError, Args:kwargs, Url:urll, Cause:fmt.Sprintf("Request failed due to an unknown error: %s", err.Error())}
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 401 {
			return resp, &RequestsError{Enum:Unauthorized, Args:kwargs, Url:urll, Cause:fmt.Sprintf("Unauthorized (%d)", resp.StatusCode), Response:resp}
		} else {
			return resp, &RequestsError{Enum:StatusNot200, Args:kwargs, Url:urll, Cause:fmt.Sprintf("Status Code not 200: %d", resp.StatusCode), Response:resp}
		}

	}

	/* body, _ := ioutil.ReadAll(resp.Body)
	result := string(body)
	println(result) */

	return resp, nil
}

// ToDo: Improve POST function with custom args and fix error return
func makePost(url string, fileType string, content []byte) (*http.Response, int) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fileType, "file.jpg")
	part.Write(content)
	writer.Close()

	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Printf("Error in the URL '%s'", url)
		fmt.Println(err)
		return nil, WrongRequest
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(r)
	if err != nil {
		println("Error while requesting...")
		fmt.Println(err)
		return nil, ResponseError
	}

	if resp.StatusCode != 200 {
		fmt.Println("Status code not 200.")
		return resp, RequestNotOk
	}

	return resp, RequestOk
}