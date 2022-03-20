package grpc_api

import (
	"context"
	"github.com/Meghdut-Mandal/Nimie/models"
	"github.com/Meghdut-Mandal/Nimie/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NimieApiServerImpl struct {
}

func (*NimieApiServerImpl) RegisterUser(_ context.Context, request *RegisterUserRequest) (*RegisterUserResponse, error) {
	_, err := utils.PublicKeyFrom64(request.GetPubicKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Public key is invalid")
	}

	user := models.AddNewUser(request.GetPubicKey())

	return &RegisterUserResponse{
		UserId:    user.UserId,
		CreatedAt: user.CreateTime,
		Jwt:       "User created successfully",
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
		bulkStatus[i] = &ApiStatus{
			StatusId:   statusObj.StatusId,
			CreateTime: statusObj.CreateTime,
			LinkId:     statusObj.LinkId,
			Text:       statusObj.HeaderText,
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

	if request.Reply == "" {
		return nil, status.Error(codes.InvalidArgument, "Status text is empty")
	}

	ConversationId, publicKey, err := models.NewConversation(request.StatusId, request.Reply, request.UserId)

	return &InitiateConversationResponse{
		ConversationId: ConversationId,
		PublicKey:      publicKey,
	}, err
}
func (*NimieApiServerImpl) GetConversationMessages(context.Context, *GetConversationMessagesRequest) (*GetConversationMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversationMessages not implemented")
}
func (*NimieApiServerImpl) GetConversations(context.Context, *GetConversationsRequest) (*GetConversationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversations not implemented")
}
