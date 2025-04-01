package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("Internal server error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem")
}
func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Bad request error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}
func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("Not found error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusNotFound, "Not found")
}
func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Conflict error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusConflict, err.Error())
}
func (app *application) forbiddenError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Forbidden error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusForbidden, err.Error())
}
func (app *application) unauthorisedError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Unauthorised error", "method", r.Method, "path", r.URL.Path, "error", err)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
func (app *application) unauthorisedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("Unauthorised error", "method", r.Method, "path", r.URL.Path, "error", err)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
