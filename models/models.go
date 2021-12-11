package models

import (
	"Nimie_alpha/config"
	"gorm.io/gorm"
)

type Conversation struct {
	ConversationId int64  `json:"conversation_id" gorm:"primaryKey;autoIncrement;notNull"`
	UserIdA        string `json:"user_id_a"`
	UserIdB        string `json:"user_id_b"`
	CreatedAt      int64  `json:"created_at" gorm:"autoCreateTime"`
}

type User struct {
	UserId     int64  `json:"user_id" gorm:"primaryKey;autoIncrement;notNull"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	PublicKey  string `json:"public_key"`
}

type Message struct {
	MessageId      int64  `json:"message_id" gorm:"primaryKey;autoIncrement;notNull"`
	ConversationId string `json:"conversation_id"`
	CreateTime     int64  `json:"create_time" gorm:"autoCreateTime"`
	UserId         string `json:"user_id"`
	Message        string `json:"message"`
}

type Status struct {
	StatusId   int64  `json:"status_id" gorm:"primaryKey;autoIncrement;notNull"`
	UserId     string `json:"user_id"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	HeaderText string `json:"header_text"`
	LinkId     string `json:"link_id"`
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Conversation{}, &User{}, &Message{}, &Status{})
}

// AddConversation Add conversion to database
func AddConversation(conversation *Conversation) {
	db.Create(conversation)
}

// GetConversation Get conversation from database
func GetConversation(conversationId int64) *Conversation {
	conversation := &Conversation{}
	db.Where("conversation_id = ?", conversationId).First(conversation)
	return conversation
}

// AddNewUser Add user to db
func AddNewUser(publicKey string) *User {
	user := &User{
		PublicKey: publicKey,
	}
	db.Create(user)
	return user
}
