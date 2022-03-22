package controllers

import (
	"github.com/Meghdut-Mandal/Nimie/models"
	"github.com/Meghdut-Mandal/Nimie/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

// Map of connection id to list of active clients
var clients = make(map[int64][]*websocket.Conn) // connected clients
var broadcast = make(chan models.ChatMessage)   // broadcast channel
// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GetIds  get conversation and user ids from request
func GetIds(r *http.Request) (int64, int64) {
	vars := mux.Vars(r)
	return utils.ParseInt64(vars["conversation_id"]), utils.ParseInt64(vars["user_id"])
}

func HandleChatConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	conversationId, userId := GetIds(r)
	// check if the conversation exists
	conversation := models.GetConversation(conversationId)
	if conversation.ConversationId == 0 {
		http.Error(w, "Conversation does not exist", http.StatusBadRequest)
		return
	}
	// check if the user is in the conversation
	if !(conversation.UserIdB == userId || conversation.UserIdA == userId) {
		http.Error(w, "User is not in the conversation", http.StatusBadRequest)
		return
	}

	// add the connection to the map
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(ws)

	clients[conversationId] = append(clients[conversationId], ws)

	for {
		var msg models.ChatMessage
		// Read in a new message as JSON and map it to a ChatMessage object
		err := ws.ReadJSON(&msg)
		msg.ConversationId = conversationId
		msg.UserId = userId
		models.AddMessage(&msg)
		// print the message to the console
		if err != nil {
			log.Printf("error: %v", err)
			deleteClient(clients[conversationId], ws)
			break
		}
		//println("message sent by ", msg.UserId)
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		//println("Broadcast message ", msg.Message)
		// Send it out to every client that is currently connected
		// Get the clint list having the same Conversation id
		clientList := clients[msg.ConversationId]

		for _, client := range clientList {
			// Send the message
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				// Remove the client from the list
				deleteClient(clientList, client)
			}
		}
	}
}

func deleteClient(clientList []*websocket.Conn, client *websocket.Conn) {
	for i, c := range clientList {
		if c == client {
			clientList = append(clientList[:i], clientList[i+1:]...)
			break
		}
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	requestBody := models.GetConversationMessages{}
	utils.ParseBody(r, &requestBody)
	//check if the conversation id is valid
	if requestBody.MessageId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid conversation id")
		return
	} else if requestBody.ConversationId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid user id")
		return
	}
	messages, err := models.GetMessages(requestBody.MessageId, requestBody.ConversationId)
	// handel error
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, models.ConversationMessages{
		Messages: messages,
		Status:   strconv.Itoa(len(messages)) + " Messages are read.",
	})
}

// GetConversations  conversation messages of a user
func GetConversations(w http.ResponseWriter, r *http.Request) {
	// get the user id from the request
	userid := utils.GetUserId(r)

	// get all Conversations of the user
	conversations := models.GetConversations(userid, 0, 100)

	// if conversation is empty
	if len(conversations) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "No conversation found")
		return
	}

	// respond with the conversation messages
	utils.RespondWithJSON(w, http.StatusOK, conversations)

}
