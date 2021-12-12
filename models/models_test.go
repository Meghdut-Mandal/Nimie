package models

import "testing"

// test add conversation method
func TestConverstionsCrud(t *testing.T) {
	// add conversation
	conversation := &Conversation{
		UserIdA: 134,
		UserIdB: 31434,
	}
	AddConversation(conversation)
	readData := GetConversation(conversation.ConversationId)
	// assert if conversation and readData are equal
	if readData.ConversationId != conversation.ConversationId {
		println("ConversationId is not equal")
	}
}
