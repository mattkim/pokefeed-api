package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libhttp"
	"github.com/pokefeed/pokefeed-api/models"
)

// PostSignupStruct struct
type PostSignupStruct struct {
	Email         string `json:"email"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	PasswordAgain string `json:"password_again"`
}

// PostLoginStruct struct
type PostLoginStruct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GetSignup method
func GetSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/users/users-external.html.tmpl", "templates/users/signup.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, nil)
}

// OptionsSignup method
func OptionsSignup(w http.ResponseWriter, r *http.Request) {
	// TODO: move these header setters into a lib
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// PostSignup method
func PostSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	db := context.Get(r, "db").(*sqlx.DB)

	// TODO: move into a util method
	decoder := json.NewDecoder(r.Body)
	var t PostSignupStruct
	err := decoder.Decode(&t)

	if err != nil {
		libhttp.HandleBadRequest(w, err)
		return
	}

	user, err2 := models.NewUser(db).Signup(nil, t.Email, t.Username, t.Password, t.PasswordAgain)
	if err2 != nil {
		if ae, ok := err2.(*pq.Error); ok {
			libhttp.HandlePostgresError(w, *ae)
			return
		}
		libhttp.HandleBadRequest(w, err2)
		return
	}

	if user == nil {
		libhttp.HandleBadRequest(w, errors.New("User does not exist."))
		return
	}

	json.NewEncoder(w).Encode(*user)
}

// GetLoginWithoutSession method
func GetLoginWithoutSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.ParseFiles("templates/users/users-external.html.tmpl", "templates/users/login.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, nil)
}

// GetLogin get login page.
func GetLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := context.Get(r, "sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "pokefeed-api-session")

	currentUserInterface := session.Values["user"]
	if currentUserInterface != nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	GetLoginWithoutSession(w, r)
}

// OptionsLogin performs login.
func OptionsLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// PostLogin performs login.
func PostLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	db := context.Get(r, "db").(*sqlx.DB)

	decoder := json.NewDecoder(r.Body)
	var t PostLoginStruct
	err := decoder.Decode(&t)

	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	u := models.NewUser(db)

	user, err2 := u.GetUserByEmailAndPassword(nil, t.Email, t.Password)
	if err2 != nil {
		libhttp.HandleBadRequest(w, err2)
		return
	}

	if user == nil {
		libhttp.HandleBadRequest(w, errors.New("User does not exist."))
		return
	}

	json.NewEncoder(w).Encode(*user)
}

// GetLogout method
func GetLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	sessionStore := context.Get(r, "sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "pokefeed-api-session")

	delete(session.Values, "user")
	session.Save(r, w)

	http.Redirect(w, r, "/login", 302)
}

// PostPutDeleteUsersID method
func PostPutDeleteUsersUUID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	method := r.FormValue("_method")
	if method == "" || strings.ToLower(method) == "post" || strings.ToLower(method) == "put" {
		PutUsersUUID(w, r)
	} else if strings.ToLower(method) == "delete" {
		DeleteUsersUUID(w, r)
	}
}

// PutUsersID method
func PutUsersUUID(w http.ResponseWriter, r *http.Request) {
	userUUID, err := getUUIDFromPath(w, r)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	db := context.Get(r, "db").(*sqlx.DB)

	sessionStore := context.Get(r, "sessionStore").(sessions.Store)

	session, _ := sessionStore.Get(r, "pokefeed-api-session")

	currentUser := session.Values["user"].(*models.UserRow)

	if currentUser.UUID != userUUID {
		err := errors.New("Modifying other user is not allowed.")
		libhttp.HandleErrorJson(w, err)
		return
	}

	email := r.FormValue("Email")
	password := r.FormValue("Password")
	passwordAgain := r.FormValue("PasswordAgain")

	u := models.NewUser(db)

	currentUser, err = u.UpdateEmailAndPasswordByUUID(nil, currentUser.UUID, email, password, passwordAgain)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	// Update currentUser stored in session.
	session.Values["user"] = currentUser
	err = session.Save(r, w)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	http.Redirect(w, r, "/", 302)
}

// DeleteUsersID method
func DeleteUsersUUID(w http.ResponseWriter, r *http.Request) {
	err := errors.New("DELETE method is not implemented.")
	libhttp.HandleErrorJson(w, err)
	return
}
