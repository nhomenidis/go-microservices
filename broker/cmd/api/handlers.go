package main

import (
	"bytes"
	"common"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(writer http.ResponseWriter, request *http.Request) {
	payload := common.JsonResponse{
		Error:   false,
		Message: "Call the Broker",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write(out)
}

func (app *Config) HandleSubmission(writer http.ResponseWriter, request *http.Request) {
	var requestPayload RequestPayload

	err := common.ReadJson(writer, request, &requestPayload)
	if err != nil {
		common.ErrorJSON(writer, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(writer, requestPayload.Auth)
	default:
		common.ErrorJSON(writer, errors.New("invalid action"), http.StatusBadRequest)
	}

	out, _ := json.MarshalIndent(requestPayload, "", "\t")
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusAccepted)
	writer.Write(out)
}

func (app *Config) authenticate(writer http.ResponseWriter, authPayload AuthPayload) {
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	request, err := http.NewRequest(
		"POST",
		"http://authentication-service/authenticate",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		common.ErrorJSON(writer, err, http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		common.ErrorJSON(writer, err, http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		common.ErrorJSON(writer, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	} else if response.StatusCode != http.StatusAccepted {
		common.ErrorJSON(writer, errors.New("authentication service error"), http.StatusInternalServerError)
		return
	}

	var authResponse common.JsonResponse

	err = json.NewDecoder(response.Body).Decode(&authResponse)

	if authResponse.Error {
		common.ErrorJSON(writer, errors.New(authResponse.Message), http.StatusUnauthorized)
		return
	}

	var payload common.JsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = authResponse.Data

	common.WriteJSON(writer, http.StatusAccepted, payload)
}
