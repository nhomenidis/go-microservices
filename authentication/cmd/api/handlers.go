package main

import (
	"common"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(writer http.ResponseWriter, request *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := common.ReadJson(writer, request, &requestPayload)
	if err != nil {
		common.ErrorJSON(writer, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		common.ErrorJSON(writer, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		common.ErrorJSON(writer, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := common.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	common.WriteJSON(writer, http.StatusAccepted, payload)
}
