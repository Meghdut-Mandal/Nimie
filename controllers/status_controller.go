package controllers

import (
	"Nimie/models"
	"Nimie/utils"
	"github.com/gorilla/mux"
	"net/http"
)

// add status controller

func CreateStatus(w http.ResponseWriter, r *http.Request) {
	// get CreateStatus struct from request body
	requestBody := &models.CreateStatus{}
	utils.ParseBody(r, requestBody)
	userId := utils.GetUserId(r)

	if requestBody.Text == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "text is required")
		return
	}
	status := models.AddStatus(&requestBody.Text, userId)
	utils.RespondWithJSON(w, http.StatusOK, models.StatusCreated{
		UniqueId: status.StatusId,
	})
}

// DeleteStatus `Delete status controller
func DeleteStatus(w http.ResponseWriter, r *http.Request) {
	// get DeleteStatus struct from request body
	userId := utils.GetUserId(r)
	vars := mux.Vars(r)
	statusId := utils.ParseInt64(vars["status_id"])

	if statusId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "statusId is required")
		return
	}
	status := models.RemoveStatus(statusId, userId)
	utils.RespondWithJSONMessage(w, http.StatusOK, status)
}

func ReplyStatus(w http.ResponseWriter, r *http.Request) {
	// get InitiateConversation struct from request body
	requestBody := &models.InitiateConversation{}
	utils.ParseBody(r, requestBody)
	userIdB := utils.GetUserId(r)

	// check if userIdA is valid
	if requestBody.StatusId == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "statusId is required")
		return
	} else if requestBody.Reply == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "reply is required")
		return
	}

	ConversationId, publicKey, err := models.NewConversation(requestBody.StatusId, requestBody.Reply, userIdB)
	// handle error response
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.ConversationCreated{
		ConversationID: ConversationId,
		PublicKey:      publicKey,
	})
}
