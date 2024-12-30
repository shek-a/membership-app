package models

type Response struct {
	StatusCode int
	Body       any
}

type ErrorMessage struct {
	Error string `json:"error" binding:"required"`
}
