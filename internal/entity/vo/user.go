package vo

type UserVO struct {
	UserId   uint64 `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
