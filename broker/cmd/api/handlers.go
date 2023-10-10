package main

import (
	"encoding/json"
	"net/http"
)

func (app *Config) Broker(writer http.ResponseWriter, request *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Call the Broker",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write(out)
}
