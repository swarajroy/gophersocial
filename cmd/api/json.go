package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	MAX_BYTES = 1_048_578 // limit requests to 1 MB
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return encoder.Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(MAX_BYTES))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		ErrMsg string
	}
	return writeJson(w, status, &envelope{ErrMsg: message})
}
