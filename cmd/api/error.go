package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s path: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem")
}
func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request error: %s path: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found error: %s path: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusNotFound, "Not found")
}
func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflct error: %s path: %s", r.Method, r.URL.Path, err)
	writeJSONError(w, http.StatusConflict, err.Error())
}
