package gobot

import (
	"net/http"
	"bytes"
	"mime/multipart"
	"net/url"
	"strings"
	"fmt"
	"time"
	"net"
	"io"
)

var client = &http.Client{}

/*
Maybe this function is a mess...
But hey it works!
I'm going to improve it later
Please if you have suggestions feel free to tell me!
*/
func MakeRequest(urll string, contentBytes []byte, kwargs map[string]string) (io.Reader, *RequestsError) {
	// Initialize some vars...
	var body = &bytes.Buffer{}
	var writer = multipart.NewWriter(body)
	var req *http.Request
	var err error
	var isPost = false
	var fileName = ""
	var fileType = ""
	var mode = "GET"

	// Check if we've passed some bytes
	if contentBytes != nil {
		isPost = true
		mode = "POST"
	}

	// Build the request
	u, _ := url.Parse(urll)
	q := u.Query()
	for key, value := range kwargs {
		if value == "" {
			continue
		}

		if isPost {
			if strings.ToLower(key) == "filename" {
				fileName = value
				continue
			} else if strings.ToLower(key) == "filetype" {
				fileType = value
				continue
			}
		}

		q.Add(key, value)
	}

	u.RawQuery = q.Encode()

	if isPost {
		// We NEED these two values when POSTing
		if fileName == "" {
			return nil, &RequestsError{Enum: ArgsError, Args: kwargs, Url: urll, Cause: "Missing 'fileName' from POST args"}
		}

		if fileType == "" {
			return nil, &RequestsError{Enum: ArgsError, Args: kwargs, Url: urll, Cause: "Missing 'fileType' from POST args"}
		}

		part, _ := writer.CreateFormFile(fileType, fileName)
		part.Write(contentBytes)
		writer.Close()
	}

	req, err = http.NewRequest(mode, u.String(), body)
	if err != nil {
		return nil, &RequestsError{Enum: RequestNotOk, Args: kwargs, Url: u.String(), Cause: fmt.Sprintf("Error in the URL '%s'", u.String())}
	}

	// Set a timeout, 120 is a lot tbh but we need it fro GetUpdates
	client.Timeout = time.Duration(time.Second * 122)

	if isPost {
		req.Header.Add("Content-Type", writer.FormDataContentType())
	}

	httpResp, err := client.Do(req)
	if err != nil {
		switch err := err.(type) {
		case net.Error:
			if err.Timeout() {
				return nil, &RequestsError{Enum: TimeoutError, Args: kwargs, Url: urll, Cause: fmt.Sprintf("Timeout exceeded (%s)", err)}
			}
		case *url.Error:
			if err, ok := err.Err.(net.Error); ok && err.Timeout() {
				return nil, &RequestsError{Enum: TimeoutError, Args: kwargs, Url: urll, Cause: fmt.Sprintf("Timeout exceeded (%s)", err)}
			}
		}
		return nil, &RequestsError{Enum: ResponseError, Args: kwargs, Url: urll, Cause: fmt.Sprintf("Request failed due to an unknown error: %s", err.Error())}
	}

	if httpResp.StatusCode != 200 {
		// Add some other cases with a switch
		if httpResp.StatusCode == 401 {
			return httpResp.Body, &RequestsError{Enum: Unauthorized, Args: kwargs, Url: urll, Cause: fmt.Sprintf("Unauthorized (%d)", httpResp.StatusCode), Response: httpResp.Body}
		} else if httpResp.StatusCode == 400 {
			return httpResp.Body, &RequestsError{Enum: BadRequest, Args: kwargs, Url: urll, Cause: fmt.Sprintf("BadRequest (%d)", httpResp.StatusCode), Response: httpResp.Body}
		} else {
			return httpResp.Body, &RequestsError{Enum: StatusNot200, Args: kwargs, Url: urll, Cause: fmt.Sprintf("Status Code not 200: %d", httpResp.StatusCode), Response: httpResp.Body}
		}
	}

	// Everything is ok!
	return httpResp.Body, nil
}
