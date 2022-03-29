package models

import (
	"errors"
	"github.com/Meghdut-Mandal/Nimie/config"
	"github.com/Meghdut-Mandal/Nimie/utils"
	"gorm.io/gorm"
	"time"
)

type Conversation struct {
	ConversationId int64 `json:"conversation_id" gorm:"primaryKey;autoIncrement;notNull"`
	UserIdA        int64 `json:"user_id_a"`
	UserIdB        int64 `json:"user_id_b"`
	CreatedAt      int64 `json:"created_at" gorm:"autoUpdateTime:milli"`
	StatusId       int64 `json:"status_id"`
}

type User struct {
	UserId     int64  `json:"user_id" gorm:"primaryKey;autoIncrement;notNull"`
	CreateTime int64  `json:"create_time" gorm:"autoUpdateTime:milli"`
	PublicKey  []byte `json:"public_key"`
}

type ChatMessage struct {
	MessageId      int64  `json:"message_id" gorm:"primaryKey;autoIncrement;notNull"`
	ConversationId int64  `json:"conversation_id"`
	CreateTime     int64  `json:"create_time" gorm:"autoUpdateTime:milli"`
	UserId         int64  `json:"user_id"`
	Message        []byte `json:"message"`
	IsSeen         bool   `json:"is_seen"`
	MessageType    string `json:"message_type"`
}

type Status struct {
	StatusId   int64  `json:"status_id" gorm:"primaryKey;autoIncrement;notNull"`
	UserId     int64  `json:"user_id"`
	CreateTime int64  `json:"create_time" gorm:"autoUpdateTime:milli"`
	HeaderText string `json:"header_text"`
	LinkId     string `json:"link_id"`
}

type KeyExchangeRequest struct {
	ConversationId int64  `json:"conversation_id"`
	AESKey         []byte `json:"aes_key"`
}

var db *gorm.DB

func init() {
	config.Connect()
	db = config.GetSqlDB()
	db.AutoMigrate(&Conversation{}, &User{}, &ChatMessage{}, &Status{}, &KeyExchangeRequest{})
}

// AddKeyExchangeRequest add new KeyExchangeRequest request
func AddKeyExchangeRequest(conversationId int64, aesKey []byte) error {
	// verify if conversation exists
	var conversation Conversation
	db.Where("conversation_id = ?", conversationId).First(&conversation)
	if conversation.ConversationId == 0 {
		return errors.New("conversation not found")
	}

	db.Create(&KeyExchangeRequest{
		ConversationId: conversationId,
		AESKey:         aesKey,
	})
	return nil
}

// GetKeyExchangeRequest Get the AES key for the key exchange request
func GetKeyExchangeRequest(conversationId int64) ([]byte, error) {
	var keyExchangeRequest KeyExchangeRequest
	db.Where("conversation_id = ?", conversationId).First(&keyExchangeRequest)
	if keyExchangeRequest.ConversationId == 0 {
		return nil, errors.New("key exchange request not found")
	}
	return keyExchangeRequest.AESKey, nil
}

// DeleteKeyExchangeRequest delete the key exchange request
func DeleteKeyExchangeRequest(conversationId int64) error {
	var keyExchangeRequest KeyExchangeRequest
	db.Where("conversation_id = ?", conversationId).First(&keyExchangeRequest)
	if keyExchangeRequest.ConversationId == 0 {
		return errors.New("key exchange request not found")
	}
	db.Delete(&keyExchangeRequest).Where("conversation_id = ?", conversationId)
	return nil
}

// NewConversation Here user_id_b is the sender and user_id_a is the receiver
func NewConversation(statusId int64, reply []byte, userIdB int64) (int64, []byte, error) {

	// Read status from database
	status := Status{}
	db.Where("status_id = ?", statusId).First(&status)
	// check if the status is valid
	if status.StatusId == 0 {
		return -1, nil, utils.NewError("Status not found")
	}

	userIdA := status.UserId
	// Read userA from database
	userA := User{}
	db.Where("user_id = ?", userIdA).First(&userA)
	// check if the user is valid
	if userA.UserId == 0 {
		return -1, nil, utils.NewError("User not found")
	}
	// read userB from database
	userB := User{}
	db.Where("user_id = ?", userIdB).First(&userB)
	// check if the user is valid
	if userB.UserId == 0 {
		return -1, nil, utils.NewError("User not found")
	}

	conversation := Conversation{
		UserIdA:  userIdA,
		UserIdB:  userIdB,
		StatusId: statusId,
	}
	err := db.Create(&conversation).Error
	if err != nil {
		return -1, nil, err
	}
	chatMessage := ChatMessage{
		Message:        reply,
		ConversationId: conversation.ConversationId,
		UserId:         userIdB,
		IsSeen:         false,
		MessageType:    "text",
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
func AddNewUser(publicKey []byte) *User {
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
func AddMessage(message *ChatMessage) ChatMessage {
	db.Create(message)
	return *message
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

func GetUserPublicKey(userId int64) []byte {
	user := &User{}
	db.Where("user_id = ?", userId).First(user)
	return user.PublicKey
}

func GetMessages(messageId int64, conversationId int64) ([]ChatMessage, error) {
	// read conversation from db
	conversation := GetConversation(conversationId)
	// check if the conversation is valid
	if conversation.ConversationId == 0 {
		return []ChatMessage{}, utils.NewError("Conversation not found")
	}
	var messages []ChatMessage
	db.Where("conversation_id = ? and message_id > ?", conversationId, messageId).Find(&messages)
	return messages, nil
}

// GetStatusFromLink get status form unique link
func GetStatusFromLink(linkId string) *Status {
	status := &Status{}
	db.Where("link_id = ?", linkId).First(status)
	return status
}

// GetConversations Get conversation from database of a user
func GetConversations(userId int64, offset int, limit int) []Conversation {
	var conversations []Conversation
	db.Where("user_id_a = ? or user_id_b = ?", userId, userId).
		Offset(offset).Limit(limit).
		Find(&conversations)
	return conversations
}

func GetLastMessage(conversationId int64) []byte {
	message := &ChatMessage{}
	db.Where("conversation_id = ?", conversationId).Order("create_time desc").First(message)
	return message.Message
}
