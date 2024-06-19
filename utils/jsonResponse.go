package utils

import (
	"encoding/json"
	"net/http"
)

type JSONResponse interface {
	GetStatus() bool
	GetMessage() string
}

type SuccessResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (sr *SuccessResponse) GetStatus() bool {
	return sr.Status
}

func (sr *SuccessResponse) GetMessage() string {
	return sr.Message
}

type ErrorResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func (er *ErrorResponse) GetStatus() bool {
	return er.Status
}

func (er *ErrorResponse) GetMessage() string {
	return er.Message
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, response JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
	return
}
