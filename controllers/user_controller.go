package controllers

import (
	"Nimie_alpha/models"
	"Nimie_alpha/utils"
	"net/http"
)

// RegisterUser controller for adding user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	requestBody := &models.RegisterUser{}
	utils.ParseBody(r, requestBody)
	if requestBody.PublicKey == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Public key is required")
		return
	}
	user := models.AddNewUser(requestBody.PublicKey)
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":    "User registered successfully",
		"user_id":    user.UserId,
		"created_at": user.CreateTime,
	})

}
