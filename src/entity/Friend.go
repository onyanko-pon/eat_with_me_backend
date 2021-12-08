package entity

type Friend struct {
	User   User   `json:"user"`
	Status string `json:"status"`
}
