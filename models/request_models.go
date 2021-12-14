package models

type RegisterUser struct {
	PublicKey string `json:"public_key"`
}

type CreateStatus struct {
	Text string `json:"text"`
}

type InitiateConversation struct {
	Reply    string `json:"reply"`
	StatusId int64  `json:"status_id"`
}

type GetConversationMessages struct {
	MessageId      int64 `json:"message_id"`
	ConversationId int64 `json:"conversation_id"`
}
