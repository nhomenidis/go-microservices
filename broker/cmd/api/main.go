package main

import (
	"bytes"
	"common"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct{}

func (app *Config) logItem(writer http.ResponseWriter, logPayload LogPayload) {
	jsonData, _ := json.MarshalIndent(logPayload, "", "\t")

	logServiceUrl := "http://logger-service:8085/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		common.ErrorJSON(writer, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	// TODO: remove this
	log.Println("We are just about to call the log service")
	response, err := client.Do(request)
	if err != nil {
		common.ErrorJSON(writer, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		common.ErrorJSON(writer, err)
		return
	}

	var payload common.JsonResponse
	payload.Error = false
	payload.Message = "Log has been written"

	common.WriteJSON(writer, http.StatusAccepted, payload)
}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s\n", webPort)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
