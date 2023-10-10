package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJson(writer http.ResponseWriter, request *http.Request, data interface{}) error {
	var maxBytes = 1048576

	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))

	dec := json.NewDecoder(request.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("multiple json values are not allowed at body")
	}

	return nil
}
