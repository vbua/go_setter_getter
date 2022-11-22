package httpserver

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Error   interface{} `json:"error,omitempty"`
	Message interface{} `json:"result,omitempty"`
}

func sendResponse(error bool, message interface{}, code int, w http.ResponseWriter) {
	resp := Response{}
	if error {
		resp.Error = message
		w.WriteHeader(code)
	} else {
		resp.Message = message
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(respBytes)
}
