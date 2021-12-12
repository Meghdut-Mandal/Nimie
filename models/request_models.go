package models

type RegisterUser struct {
	PublicKey string `json:"public_key"`
}

type CreateStatus struct {
	Text string `json:"text"`
}
