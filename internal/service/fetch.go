package service

import (
	"encoding/json"
	"net/http"
	"time"
)

const emojihubURL = "https://emojihub.yurace.pro/api/all"

var client = &http.Client{Timeout: 10 * time.Second}

type EmojiRaw struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Group    string   `json:"group"`
	HtmlCode []string `json:"htmlCode"`
	Unicode  []string `json:"unicode"`
}

func FetchAll() ([]EmojiRaw, error) {
	resp, err := client.Get(emojihubURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []EmojiRaw
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
