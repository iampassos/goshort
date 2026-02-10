package domain

var Urls = make(map[string]Url)

type Url struct {
	Id       int    `json:"id"`
	ShortUrl string `json:"short_url"`
	LongUrl  string `json:"long_url"`
	UserId   int    `json:"user_id"`
}
