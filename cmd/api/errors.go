package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error - method : %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error - method : %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error - method : %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error - method : %s path: %s error: %s", r.Method, r.URL.Path, err)
	writeJsonError(w, http.StatusConflict, err.Error())
}
