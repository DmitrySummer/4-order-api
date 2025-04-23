package res

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, statusCode int, message interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	var errMessage string
	switch v := message.(type) {
	case error:
		errMessage = v.Error()
	case string:
		errMessage = v
	default:
		errMessage = "internal server error"
	}

	json.NewEncoder(w).Encode(Response{
		Success: false,
		Error:   errMessage,
	})
}

func Json(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-type", "aplication/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
