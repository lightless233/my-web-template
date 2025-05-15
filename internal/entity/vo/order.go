package vo

type OrderVO struct {
	OrderId string   `json:"order_id"`
	UserId  uint64   `json:"user_id"`
	Total   uint64   `json:"total"`
	Items   []string `json:"items"`
}
