package model

import "time"

type Note struct {
	Id string `json:"id"`
	Message string `json:"message"`
	CreatedTimestamp time.Time `json:"created_timestamp"`
}

type CreateNoteRequest struct {
	Message string `json:"message"`
}

type Response struct {
	StatusCode int
	Message    string
}

type Message struct {
	Content string `json:"content"`
}
