package models

type StatusCreated struct {
	UniqueId int64 `json:"unique_id"`
}
type UserCreated struct {
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at"`
	UserId    int64  `json:"user_id"`
}

type ConversationCreated struct {
	ConversationID int64  `json:"conversation_id"`
	PublicKey      string `json:"public_key"`
}

type ConversationMessages struct {
	Messages []ChatMessage `json:"messages"`
	Status   string        `json:"status"`
}
