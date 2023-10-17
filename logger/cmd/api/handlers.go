package main

import (
	"common"
	"log"
	"logger/data"
	"net/http"
)

type JsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(writer http.ResponseWriter, request *http.Request) {
	var requestPayload JsonPayload

	// TODO: remove this
	log.Println("We are ready to write a log")

	_ = common.ReadJson(writer, request, &requestPayload)

	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.MongoService.Insert(event)
	if err != nil {
		common.ErrorJSON(writer, err)
		return
	}

	response := common.JsonResponse{
		Error:   false,
		Message: "The log was written",
	}

	common.WriteJSON(writer, http.StatusAccepted, response)
}
