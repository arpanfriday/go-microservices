package main

import (
	"log-service/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read json into var
	var requestPayload jsonPayload
	_ = app.readJson(w, r, &requestPayload)

	// insert the data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	if err := app.Models.LogEntry.Insert(event); err != nil {
		app.errorJson(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
	}
	app.writeJson(w, http.StatusAccepted, resp)
}
