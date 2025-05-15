package constant

const (
	UserStatusActive   = 1
	UserStatusDisabled = 2
)

var UserStatusMap = map[int]string{
	UserStatusActive:   "Active",
	UserStatusDisabled: "Disabled",
}

func GetUserStatusName(status int) string {
	return UserStatusMap[status]
}
