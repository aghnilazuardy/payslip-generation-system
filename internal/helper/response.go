package helper

import (
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"errors"`
	Data    interface{} `json:"data"`
}

type EmptyObj struct{}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, message string, data, errors interface{}) Response {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	status := "error"
	if statusCode >= 200 && statusCode < 300 {
		status = "success"
	}

	res := Response{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Error:   errors,
		Data:    data,
	}
	return res
}
