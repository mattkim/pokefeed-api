package application

import (
	"net/http"

	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"github.com/pokefeed/pokefeed-api/handlers"
	"github.com/pokefeed/pokefeed-api/middlewares"
)

// New is the constructor for Application struct.
func New(config *viper.Viper) (*Application, error) {
	dsn := config.Get("dsn").(string)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	cookieStoreSecret := config.Get("cookie_secret").(string)

	app := &Application{}
	app.config = config
	app.dsn = dsn
	app.db = db
	app.sessionStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, err
}

// Application is the application object that runs HTTP server.
type Application struct {
	config       *viper.Viper
	dsn          string
	db           *sqlx.DB
	sessionStore sessions.Store
}

func (app *Application) MiddlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetDB(app.db))
	middle.Use(middlewares.SetSessionStore(app.sessionStore))

	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()

	router.Handle("/", MustLogin(http.HandlerFunc(handlers.GetHome))).Methods("GET")

	// New stuff
	router.HandleFunc("/allfeedtags", handlers.Options).Methods("OPTIONS")
	router.HandleFunc("/latestfeeds", handlers.Options).Methods("OPTIONS")
	router.HandleFunc("/getfeeds", handlers.Options).Methods("OPTIONS")
	router.HandleFunc("/postfeed", handlers.Options).Methods("OPTIONS")
	router.HandleFunc("/postcomment", handlers.Options).Methods("OPTIONS")
	router.HandleFunc("/signup", handlers.Options).Methods("OPTIONS")
	router.HandleFunc("/login", handlers.Options).Methods("OPTIONS")
	router.HandleFunc("/create_facebook_user", handlers.Options).Methods("OPTIONS")

	router.HandleFunc("/allfeedtags", handlers.GetAllFeedTags).Methods("GET")
	router.HandleFunc("/latestfeeds", handlers.GetLatestFeeds).Methods("GET")
	router.HandleFunc("/getfeeds", handlers.GetFeeds).Methods("GET")
	router.HandleFunc("/postfeed", handlers.PostFeed).Methods("POST")
	router.HandleFunc("/postcomment", handlers.PostComment).Methods("POST")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST")
	router.HandleFunc("/create_facebook_user", handlers.CreateFacebookUser).Methods("POST")
	router.HandleFunc("/get_facebook_user", handlers.GetFacebookUser).Methods("GET")

	router.Handle("/users/{id:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersUUID))).Methods("POST", "PUT", "DELETE")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}
