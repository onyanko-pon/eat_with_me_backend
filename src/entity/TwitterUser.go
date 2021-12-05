package entity

type TwitterUser struct {
	ID                   uint64 `json:"id"`
	Name                 string `json:"name"`
	ScreenName           string `json:"screen_name"`
	ProfileImageUrlHttps string `json:"profile_image_url_https"`
}
