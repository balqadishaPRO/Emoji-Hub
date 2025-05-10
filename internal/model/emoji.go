package model

type Emoji struct {
	ID       string   `db:"id"        json:"id"`
	Name     string   `db:"name"      json:"name"`
	Category string   `db:"category"  json:"category"`
	Group    string   `db:"group"     json:"group"`
	HtmlCode []string `db:"html_code" json:"htmlCode"`
	Unicode  []string `db:"unicode"   json:"unicode"`
}

type EmojiDetail struct {
	Emoji
	Mood string `json:"mood"`
}
