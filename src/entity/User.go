package entity

type User struct {
	ID                  uint64 `json:"id"`
	Username            string `json:"username"`
	ImageURL            string `json:"imageURL"`
	TwitterScreenName   string `json:"twitter_screen_name"`
	TwitterUsername     string `json:"twitter_username"`
	TwitterUserID       uint64 `json:"twitter_user_id"`
	AppleUserIdentifier string `json:"-"`
}
