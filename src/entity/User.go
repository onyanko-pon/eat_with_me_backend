package entity

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	ImageURL string `json:"imageURL"`
}
