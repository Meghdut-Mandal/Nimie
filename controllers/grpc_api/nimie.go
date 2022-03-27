package grpc_api

import (
	"context"
	"github.com/Meghdut-Mandal/Nimie/models"
	"github.com/Meghdut-Mandal/Nimie/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type NimieApiServerImpl struct {
}

var chatClientConnections = make(map[int64]*NimieApi_ChatConnectServer)
var conversationCache = make(map[int64]*models.Conversation)

const SimpleMsgType = 1
const PingPongType = 4

func (*NimieApiServerImpl) RegisterUser(_ context.Context, request *RegisterUserRequest) (*RegisterUserResponse, error) {
	_, err := utils.PublicKeyFrom(request.GetPubicKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Public key is invalid")
	}

	user := models.AddNewUser(request.GetPubicKey())

	return &RegisterUserResponse{
		UserId:    user.UserId,
		CreatedAt: user.CreateTime,
	}, nil
}
func (*NimieApiServerImpl) CreateStatus(_ context.Context, r *CreateStatusRequest) (*CreateStatusResponse, error) {

	if r.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "Status text is empty")
	}

	statusCreated := models.AddStatus(&r.Text, r.UserId)

	return &CreateStatusResponse{
		StatusId:   statusCreated.StatusId,
		CreateTime: statusCreated.CreateTime,
		LinkId:     statusCreated.LinkId,
	}, nil
}

func (*NimieApiServerImpl) GetBulkStatus(_ context.Context, in *GetBulkStatusRequest) (*GetBulkStatusResponse, error) {
	statuses := models.GetBulkStatus(int(in.GetOffset()), int(in.GetLimit()))

	bulkStatus := make([]*ApiStatus, len(statuses))
	for i, statusObj := range statuses {

		// get the public key of the user
		publicKey := models.GetUserPublicKey(statusObj.UserId)

		bulkStatus[i] = &ApiStatus{
			StatusId:   statusObj.StatusId,
			CreateTime: statusObj.CreateTime,
			LinkId:     statusObj.LinkId,
			Text:       statusObj.HeaderText,
			PublicKey:  publicKey,
		}
	}

	return &GetBulkStatusResponse{
		BulkStatus: bulkStatus,
	}, nil
}

func (*NimieApiServerImpl) DeleteStatus(context.Context, *DeleteStatusRequest) (*GenericResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteStatus not implemented")
}
func (*NimieApiServerImpl) ReplyStatus(_ context.Context, request *InitiateConversationRequest) (*InitiateConversationResponse, error) {

	if request.Reply == nil {
		return nil, status.Error(codes.InvalidArgument, "Status text is empty")
	}

	ConversationId, publicKey, err := models.NewConversation(request.StatusId, request.Reply, request.UserId)

	return &InitiateConversationResponse{
		ConversationId: ConversationId,
		PublicKey:      publicKey,
	}, err
}

func (*NimieApiServerImpl) GetConversationList(_ context.Context, request *ConversationListRequest) (*ConversationListResponse, error) {
	conversations := models.GetConversations(request.GetUserId(), int(request.GetOffset()), int(request.GetLimit()))

	bulkApiConversation := make([]*ApiConversation, len(conversations))
	for i, conversation := range conversations {

		otherUserId := int64(0)
		// get the public key of the other user
		if conversation.UserIdA == request.GetUserId() {
			otherUserId = conversation.UserIdB
		} else {
			otherUserId = conversation.UserIdA
		}

		otherUserPublicKey := models.GetUserPublicKey(otherUserId)

		// Get last message of the conversation
		lastMessage := models.GetLastMessage(conversation.ConversationId)

		bulkApiConversation[i] = &ApiConversation{
			ConversationId: conversation.ConversationId,
			StatusId:       conversation.StatusId,
			CreateTime:     conversation.CreatedAt,
			OtherPublicKey: otherUserPublicKey,
			LastReply:      lastMessage,
		}
	}

	return &ConversationListResponse{
		Conversations: bulkApiConversation,
		Status:        "Conversations fetched successfully",
	}, nil
}

func (*NimieApiServerImpl) ChatConnect(stream NimieApi_ChatConnectServer) error {
	userId := int64(0)
	// handle client messages
	for {
		rr, err := stream.Recv() // Recv is a blocking method which returns on client data
		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			println("Client has closed connection")
			chatClientConnections[userId] = nil
			break
		}

		// any other error means the transport between the server and client is unavailable
		if err != nil {
			println("Unable to read from client", "error", err)
			chatClientConnections[userId] = nil
			return err
		}

		msg := rr.GetMessage()
		println("Received message from client", "message", msg)

		userId = msg.UserId

		// handle the message
		chatClientConnections[userId] = &stream
		println("Client connected with", "conversationId", msg.ConversationId)
		println("Client id is ", "clientId", userId)

		// find the conversation from the cache
		if conversationCache[msg.ConversationId] == nil {
			println("Conversation not found in cache, fetching from database")
			conversationCache[msg.ConversationId] = models.GetConversation(msg.ConversationId)
			println("Conversation fetched from database and added to cache ", conversationCache[msg.ConversationId])
		}

		conversation := conversationCache[msg.ConversationId]

		// find the other user's id
		otherUserId := int64(0)
		if conversationCache[msg.ConversationId].UserIdA != userId {
			otherUserId = conversation.UserIdA
		} else {
			otherUserId = conversation.UserIdB
		}

		// find the other user's conversation stream
		println("Other user id", otherUserId)
		otherClientConnection := chatClientConnections[otherUserId]

		if rr.MessageType == SimpleMsgType {
			// convert the received message to a Message object
			dbMessage := models.ChatMessage{
				ConversationId: msg.ConversationId,
				UserId:         userId,
				Message:        msg.Message,
				MessageType:    msg.ContentType,
				IsSeen:         false,
			}
			savedMessage := models.AddMessage(&dbMessage)

			go func() {
				err := stream.Send(&ChatServerResponse{
					Messages: &ApiTextMessage{
						ConversationId: savedMessage.ConversationId,
						UserId:         0,
						Message:        savedMessage.Message,
						ContentType:    savedMessage.MessageType,
						IsSeen:         savedMessage.IsSeen,
						CreateTime:     savedMessage.CreateTime,
						MessageId:      savedMessage.MessageId,
					},
					MessageType: 1,
				})
				if err != nil {
					println("Unable to send message to client", "error", err)
				}
			}()

			// send the message to the other user
			if otherClientConnection != nil {
				println("Sending message to other user with id ", otherUserId)
				go func() {
					err := (*otherClientConnection).Send(&ChatServerResponse{
						Messages: &ApiTextMessage{
							ConversationId: savedMessage.ConversationId,
							UserId:         1,
							Message:        savedMessage.Message,
							ContentType:    savedMessage.MessageType,
							IsSeen:         savedMessage.IsSeen,
							CreateTime:     savedMessage.CreateTime,
							MessageId:      savedMessage.MessageId,
						},
						MessageType: SimpleMsgType,
					})
					if err != nil {
						println("Unable to send message to client", "error", err)
					}
				}()
			}
		} else if rr.MessageType == PingPongType {
			// send a pong message to the current user
			go func() {
				err := stream.Send(&ChatServerResponse{
					Messages:    nil,
					MessageType: PingPongType,
				})
				if err != nil {
					println("Unable to send message to client", "error", err)
				}
			}()

		} else {
			println("Unknown message type")
		}

	}

	chatClientConnections[userId] = nil
	return nil
}

func (*NimieApiServerImpl) GetConversationMessages(request *GetConversationMessagesRequest, stream NimieApi_GetConversationMessagesServer) error {

	messages, err := models.GetMessages(request.LastMessageId, request.ConversationId)
	if err != nil {
		return err
	}

	// send the messages to the client
	for _, message := range messages {

		maskedId := int64(0)

		if message.UserId != request.UserId {
			maskedId = 1
		} else {
			maskedId = 0
		}

		println("Sending message to client", "message", message.MessageId, request.UserId)

		err := stream.Send(&ChatServerResponse{
			Messages: &ApiTextMessage{
				ConversationId: message.ConversationId,
				UserId:         maskedId,
				Message:        message.Message,
				ContentType:    message.MessageType,
				IsSeen:         message.IsSeen,
				CreateTime:     message.CreateTime,
				MessageId:      message.MessageId,
			},
			MessageType: SimpleMsgType,
		})
		if err != nil {
			println("Unable to send message to client", "error", err)
		}
	}
	stream.Context().Done()
	return nil
}

/*
Look for UnimplementedNimie piServer methods in the pb.go file
*/
