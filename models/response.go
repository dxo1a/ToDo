package models

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Response struct {
	Message string `json:"message"`
}
