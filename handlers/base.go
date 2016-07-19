// Package handlers provides request handlers.
package handlers

import (
	"errors"
	"github.com/pokefeed/pokefeed-api/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

func getCurrentUser(w http.ResponseWriter, r *http.Request) *models.UserRow {
	sessionStore := context.Get(r, "sessionStore").(sessions.Store)
	session, _ := sessionStore.Get(r, "pokefeed-api-session")
	return session.Values["user"].(*models.UserRow)
}

func getIdFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	userIdString := mux.Vars(r)["id"]
	if userIdString == "" {
		return -1, errors.New("user id cannot be empty.")
	}

	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return -1, err
	}

	return userId, nil
}
