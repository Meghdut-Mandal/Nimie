package models

import (
	"github.com/Meghdut-Mandal/Nimie/config"
	"github.com/Meghdut-Mandal/Nimie/utils"
	"gorm.io/gorm"
	"time"
)

type Conversation struct {
	ConversationId int64 `json:"conversation_id" gorm:"primaryKey;autoIncrement;notNull"`
	UserIdA        int64 `json:"user_id_a"`
	UserIdB        int64 `json:"user_id_b"`
	CreatedAt      int64 `json:"created_at" gorm:"autoCreateTime"`
	StatusId       int64 `json:"status_id"`
}

type User struct {
	UserId     int64  `json:"user_id" gorm:"primaryKey;autoIncrement;notNull"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	PublicKey  string `json:"public_key"`
}

type ChatMessage struct {
	MessageId      int64  `json:"message_id" gorm:"primaryKey;autoIncrement;notNull"`
	ConversationId int64  `json:"conversation_id"`
	CreateTime     int64  `json:"create_time" gorm:"autoCreateTime"`
	UserId         int64  `json:"user_id"`
	Message        string `json:"message"`
}

type Status struct {
	StatusId   int64  `json:"status_id" gorm:"primaryKey;autoIncrement;notNull"`
	UserId     int64  `json:"user_id"`
	CreateTime int64  `json:"create_time" gorm:"autoCreateTime"`
	HeaderText string `json:"header_text"`
	LinkId     string `json:"link_id"`
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetSqlDB()
	db.AutoMigrate(&Conversation{}, &User{}, &ChatMessage{}, &Status{})
}

func NewConversation(statusId int64, reply string, userIdB int64) (int64, string, error) {

	// Read status from database
	status := Status{}
	db.Where("status_id = ?", statusId).First(&status)
	// check if the status is valid
	if status.StatusId == 0 {
		return -1, "", utils.NewError("Status not found")
	}

	userIdA := status.UserId
	// Read userA from database
	userA := User{}
	db.Where("user_id = ?", userIdA).First(&userA)
	// check if the user is valid
	if userA.UserId == 0 {
		return -1, "", utils.NewError("User not found")
	}
	// read userB from database
	userB := User{}
	db.Where("user_id = ?", userIdB).First(&userB)
	// check if the user is valid
	if userB.UserId == 0 {
		return -1, "", utils.NewError("User not found")
	}

	conversation := Conversation{
		UserIdA:  userIdA,
		UserIdB:  userIdB,
		StatusId: statusId,
	}
	err := db.Create(&conversation).Error
	if err != nil {
		return -1, "", err
	}
	chatMessage := ChatMessage{
		Message:        reply,
		ConversationId: conversation.ConversationId,
		UserId:         userIdB,
	}
	AddMessage(&chatMessage)
	return conversation.ConversationId, userA.PublicKey, nil
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

// AddStatus Create a new status
func AddStatus(text *string, userId int64) *Status {
	status := &Status{
		UserId:     userId,
		HeaderText: *text,
		LinkId:     utils.Base62Encode(time.Now().UnixNano() - utils.UnixTime2021),
	}
	db.Create(status)
	return status
}

// GetAllStatus  to read all status from database given a offset and limit
func GetBulkStatus(offset int, limit int) []Status {
	var statuses []Status
	db.Order("create_time desc").Offset(offset).Limit(limit).Find(&statuses)
	return statuses
}

// AddMessage Add messages to db
func AddMessage(message *ChatMessage) {
	db.Create(message)
}

// RemoveStatus remove status from the database
func RemoveStatus(statusId int64, userId int64) string {
	status := &Status{}
	db.Where("status_id = ?", statusId).First(status)

	if status.StatusId == 0 {
		return "Status not found"
	} else if status.UserId != userId {
		return "You are not the owner of this status"
	}

	if status.UserId == userId {
		db.Delete(status)
		return "Status deleted"
	}
	return "You are not allowed to delete this status"
}

func GetMessages(messageId int64, conversationId int64) ([]ChatMessage, error) {
	// read conversation from db
	conversation := GetConversation(conversationId)
	// check if the conversation is valid
	if conversation.ConversationId == 0 {
		return []ChatMessage{}, utils.NewError("Conversation not found")
	}
	var messages []ChatMessage
	db.Where("conversation_id = ? and message_id < ?", conversationId, messageId).Limit(25).Find(&messages)
	return messages, nil
}

// GetStatusFromLink get status form unique link
func GetStatusFromLink(linkId string) *Status {
	status := &Status{}
	db.Where("link_id = ?", linkId).First(status)
	return status
}

// GetConversations Get conversation from database of a user
func GetConversations(userId int64) []Conversation {
	var conversations []Conversation
	db.Select("conversation_id").Where("user_id_a = ? or user_id_b = ?", userId, userId).Find(&conversations)
	return conversations
}
