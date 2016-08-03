// Package handlers provides request handlers.
package handlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/pokefeed/pokefeed-api/models"
)

func getCurrentUser(w http.ResponseWriter, r *http.Request) *models.UserRow {
	sessionStore := context.Get(r, "sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "pokefeed-api-session")
	return session.Values["user"].(*models.UserRow)
}

func getUUIDFromPath(w http.ResponseWriter, r *http.Request) (string, error) {
	userUUIDString := mux.Vars(r)["uuid"]
	if userUUIDString == "" {
		// TODO: why do I have to return empty string
		return "", errors.New("user uuid cannot be empty.")
	}

	return userUUIDString, nil
}
