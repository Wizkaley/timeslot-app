package models

type ServiceError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
