package gobot

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
)


// Bot types
type Bot struct {
	ThisBot GetMeResult
	Token string
	UpdateHandler UpdateHandlerType
	CommandHandlers []CommandStruct
	Offset int64
}

type CommandStruct struct {
	Command string
	Function CommandHandlerType
}

type UpdateHandlerType func(Bot, Update)

type CommandHandlerType func(Bot, Update)


// Thanks to https://mholt.github.io/json-to-go/
// JSON structs
type Update struct {
	UpdateID int64 `json:"update_id"`
	Message struct {
		MessageID int `json:"message_id"`
		From struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID    int64  `json:"id"`
			Title string `json:"title"`
			Type  string `json:"type"`
		} `json:"chat"`
		Date int `json:"date"`
		ReplyToMessage struct {
			MessageID int `json:"message_id"`
			From struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID    int64  `json:"id"`
				Title string `json:"title"`
				Type  string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"reply_to_message"`
		Text string `json:"text"`
	} `json:"message"`
}

type GetUpdateResult struct {
	Ok     bool `json:"ok"`
	Result []Update `json:"result"`
}

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

type BooleanResult struct {
	Ok     bool `json:"ok"`
	Result bool `json:"result"`
}

type ApiError struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type SendMessageArgs struct {
	ReplyToMessageID int
	ParseMode int
	DisableWebPagePreview bool
	DisableNotification bool
}

type SendPhotoArgs struct {
	Caption string
	ParseMode int
	ReplyToMessageID int
	DisableNotification bool
}

type SendDocumentArgs struct {
	Caption string
	ParseMode int
	ReplyToMessageID int
	DisableNotification bool
}