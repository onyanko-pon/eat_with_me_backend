package entity

import "time"

type Event struct {
	ID             uint64    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	OrganizeUserID uint64    `json:"organize_user_id"`
	StateDatetime  time.Time `json:"start_datetime"`
	EndDatetime    time.Time `json:"end_datetime"`
	OrganizeUser   User      `json:"organize_user"`
	JoinUsers      []User    `json:"join_users"`
}
