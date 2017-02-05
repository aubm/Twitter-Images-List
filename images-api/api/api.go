package api

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

type apiError struct {
	Error       bool   `json:"error"`
	Description string `json:"error_description"`
}

func newError(description string) apiError {
	return apiError{
		Error:       true,
		Description: description,
	}
}

var serverError = newError("An error occured, please try again later")
var invalidJSONError = newError("The request payload must be a valid JSON")

func writeJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func ok(w http.ResponseWriter) {
	writeJSON(w, map[string]interface{}{
		"error": false,
	}, 200)
}

func httpError(w http.ResponseWriter, code int, error apiError) {
	writeJSON(w, error, code)
}

type ContextProviderInterface interface {
	New(r *http.Request) context.Context
}

type ContextProvider struct{}

func (p ContextProvider) New(r *http.Request) context.Context {
	return appengine.NewContext(r)
}
