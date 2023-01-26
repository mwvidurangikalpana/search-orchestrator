package models

type Document struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type SearchRequest struct {
	Query string `json:"query"`
}
