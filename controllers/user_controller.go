package controllers

import (
	"Nimie/models"
	"Nimie/utils"
	"net/http"
)

// RegisterUser controller for adding user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	requestBody := &models.RegisterUser{}
	utils.ParseBody(r, requestBody)
	_, err := utils.PublicKeyFrom64(requestBody.PublicKey)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Public key is invalid")
		return
	}

	user := models.AddNewUser(requestBody.PublicKey)
	utils.RespondWithJSON(w, http.StatusOK, models.UserCreated{
		UserId:    user.UserId,
		CreatedAt: user.CreateTime,
		Message:   "User created successfully",
	})

}
