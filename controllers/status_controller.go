package controllers

import (
	"Nimie_alpha/models"
	"Nimie_alpha/utils"
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
