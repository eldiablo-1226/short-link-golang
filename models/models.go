package models

type InsertShortLink struct {
	Url string `json:"url"`
	Tag string `json:"tag"`
}

type Result struct {
	Url string `json:"url"`
}
