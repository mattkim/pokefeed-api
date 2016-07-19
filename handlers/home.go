package handlers

import (
	"github.com/pokefeed/pokefeed-api/models"
	"github.com/pokefeed/pokefeed-api/libhttp"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := context.Get(r, "sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "pokefeed-api-session")
	currentUser, ok := session.Values["user"].(*models.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 302)
		return
	}

	data := struct {
		CurrentUser *models.UserRow
	}{
		currentUser,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
