package server

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Response struct {
	Href  *Link   `json:"self"`
	Links []*Link `json:"_links"`
}

type Link struct {
	Name   string `json:"name,omitempty"`
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
	Rel    string `json:"rel,omitempty"`
}

func Respond(responsePayload interface{}, selfHref string, okStatusCode int, w http.ResponseWriter) {
	bytes, marshalErr := json.Marshal(responsePayload)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if marshalErr != nil {
		httpErr := &HttpError{Href: &Link{Href: selfHref}, ApiError: &ApiError{marshalErr.Error(), http.StatusInternalServerError}}
		RespondWithError(httpErr, w)
	} else {
		w.WriteHeader(okStatusCode)
		_, err := w.Write(bytes)
		if err != nil {
			logrus.WithError(err).Error("unable to write")
		}
	}
}

func RespondWithError(err *HttpError, w http.ResponseWriter) {
	if err != nil {
		logrus.WithError(err).WithField("p", err.Href.Href).Error("")
		bytes, jsonErr := json.Marshal(err)
		if jsonErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(err.StatusCode)
			_, err := w.Write(bytes)
			if err != nil {
				logrus.WithError(err).Error("unable to write")
			}
		}
	}
}
