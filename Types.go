package gobot

import "net/http"

// enums
const (
	None = iota
	Markdown
	HTML
)

const (
	Typing = "typing"
	UploadPhoto = "upload_photo"
	UploadVideo = "upload_video"
	UploadDocument = "upload_document"
)

const (
	ResponseError = iota
	RequestNotOk
	RequestOk
	WrongRequest
	TimeoutError
	StatusNot200
	Unauthorized
)


type RequestsError struct {
	Enum	int
	Url  	string
	Cause 	string
	Args 	map[string]string
	Response *http.Response
}


type CommandStruct struct {
	Command	 string
	Function CommandHandlerType
}

type UpdateHandlerType func(*Bot, Update)

type CommandHandlerType func(*Bot, Update)

type PanicHandlerType func(*Bot, Update, interface{})

// Thanks to https://mholt.github.io/json-to-go/
type Thumb   struct {
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
}

type Sticker struct {
	Width   int    	`json:"width"`
	Height  int    	`json:"height"`
	Emoji   string 	`json:"emoji"`
	SetName string 	`json:"set_name"`
	Thumb Thumb		`json:"thumb"`
	FileID   string `json:"file_id"`
	FileSize int    `json:"file_size"`
}

type From      struct {
	ID           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Bot      struct {
	ID           	int    `json:"id"`
	IsBot        	bool   `json:"is_bot"`
	FirstName    	string `json:"first_name"`
	Username     	string `json:"username"`
	token 			string
	authorized		bool
	UpdateHandler 	UpdateHandlerType
	CommandHandlers []CommandStruct
	ErrorHandler 	PanicHandlerType
	Offset 			int64
	Running 		bool
}

type Chat struct {
	ID        int64		`json:"id"`
	FirstName string 	`json:"first_name"`
	Username  string 	`json:"username"`
	Type      string 	`json:"type"`
}

type Message  struct {
	MessageID int	`json:"message_id"`
	Text string		`json:"text"`
	From From 		`json:"from"`
	Chat Chat 		`json:"chat"`
	Date    int 	`json:"date"`
	Sticker Sticker `json:"sticker"`
	Args []string
}

type Update struct {
	UpdateID 	int64 	`json:"update_id"`
	Message		Message `json:"message"`
}

type GetUpdateResult struct {
	Ok     bool		`json:"ok"`
	Result []Update	`json:"result"`
}

type SendMessageResult struct {
	Ok     bool `json:"ok"`
	Message Message `json:"result"`
}

type GetMeResult struct {
	Ok		bool	`json:"ok"`
	Bot		Bot	`json:"result"`
}

type BooleanResult struct {
	Ok     bool `json:"ok"`
	Result bool `json:"result"`
}

type SendStickerResult struct {

}

type ApiError struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type SendMessageArgs struct {
	ReplyToMessageID 		int
	ParseMode 				int
	DisableWebPagePreview	bool
	DisableNotification 	bool
}

type SendPhotoArgs struct {
	Caption				string
	ParseMode			int
	ReplyToMessageID	int
	DisableNotification bool
}

type SendDocumentArgs struct {
	Caption string
	ParseMode int
	ReplyToMessageID int
	DisableNotification bool
}

type SendStickerArgs struct {
	ReplyToMessageID int
	DisableNotification bool
}