package models

type StatusCreated struct {
	UniqueId int64 `json:"unique_id"`
}
type UserCreated struct {
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at"`
	UserId    int64  `json:"user_id"`
}
