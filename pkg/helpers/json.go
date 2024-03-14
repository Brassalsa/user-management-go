package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

type DataResponse struct {
	Data interface{} `json:"data"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Println("Failed to marshal JSON ", DataResponse{
			Data: payload,
		})
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

type errorResponse struct {
	Error string `json:"error"`
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Error code 5XX :", msg)
	}

	RespondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
