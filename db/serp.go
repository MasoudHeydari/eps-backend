package db

import (
	"time"
)

type SERP struct {
	URL         string
	Title       string
	Description string
	Phones      []string
	Emails      []string
	Keywords    []string
	IsRead      bool
	CreatedAt   time.Time
}

type SearchQuery struct {
	Id         int       `json:"id"`
	Query      string    `json:"query"`
	Language   string    `json:"language"`
	Location   string    `json:"location"`
	IsCanceled bool      `json:"is_canceled"`
	CreatedAt  time.Time `json:"created_at"`
}
