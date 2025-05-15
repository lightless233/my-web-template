package dto

type UserDTO struct {
	UserId   uint64 `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	State    uint8  `json:"state"`
}
