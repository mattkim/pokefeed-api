package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/constants"
	"github.com/pokefeed/pokefeed-api/libhttp"
	"github.com/pokefeed/pokefeed-api/models"
)

type PostFeedStruct struct {
	CreatedByUserUUID string  `json:"created_by_user_uuid"`
	Message           string  `json:"message"`
	PokemonName       string  `json:"pokemon_name"`
	Lat               float64 `json:"lat"`
	Long              float64 `json:"long"`
	FormattedAddress  string  `json:"formatted_address"`
}

type ResultStruct struct {
	Result string `json:"result"`
}

type GetFeedResultStruct struct {
	UUID              string `json:"uuid"`
	Username          string `json:"username"`
	CreatedByUserUUID string `json:"created_by_user_uuid"`
	Message           string `json:"message"`
	Pokemon           string `json:"pokemon"`
	// PokemonImageURL   string    `json:"pokemon_image_url"`
	Lat              float64   `json:"lat"`
	Long             float64   `json:"long"`
	FormattedAddress string    `json:"formatted_address"`
	CreatedAtDate    time.Time `json:"created_at_date"`
	UpdatedAtDate    time.Time `json:"updated_at_date"`
	DeletedAtDate    time.Time `json:"deleted_at_date"`
}

type GetLatestFeedsStruct struct {
	Username         string    `json:"username"`
	Message          string    `json:"message"`
	PokemonName      string    `json:"pokemon_name"`
	Lat              float64   `json:"lat"`
	Long             float64   `json:"long"`
	FormattedAddress string    `json:"formatted_address"`
	CreatedAt        time.Time `json:"created_at"`
}

type GetFeedsStruct struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

type GetAllPokemonStruct struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	ImageURL    string `json:"image_url"`
}

func GetLatestFeeds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	db := context.Get(r, "db").(*sqlx.DB)
	f := models.NewFeed(db)
	u := models.NewUser(db)

	feeds, err := f.GetLatest(nil)

	results := []*GetLatestFeedsStruct{}

	for _, feed := range feeds {
		user, err := u.GetByUUID(nil, feed.CreatedByUserUUID)

		if err != nil {
			libhttp.HandleBadRequest(w, err)
			return
		}
		result := &GetLatestFeedsStruct{}

		result.Username = user.Username
		result.Message = feed.Message
		result.PokemonName = feed.PokemonName
		result.CreatedAt = feed.CreatedAt.Time
		result.Lat = feed.Lat
		result.Long = feed.Long
		result.FormattedAddress = feed.FormattedAddress
		results = append(results, result)
	}
	if err != nil {
		libhttp.HandleBadRequest(w, err)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func GetFeeds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	lat, _ := strconv.ParseFloat(r.URL.Query()["lat"][0], 64)
	long, _ := strconv.ParseFloat(r.URL.Query()["long"][0], 64)
	latRadius, _ := strconv.ParseFloat(r.URL.Query()["latRadius"][0], 64)
	longRadius, _ := strconv.ParseFloat(r.URL.Query()["longRadius"][0], 64)

	db := context.Get(r, "db").(*sqlx.DB)
	f := models.NewFeed(db)
	u := models.NewUser(db)

	Info.Println(lat)
	Info.Println(long)

	feeds, err := f.GetFeeds(nil, lat, long, latRadius, longRadius)

	Info.Println(feeds)

	results := []*GetLatestFeedsStruct{}

	for _, feed := range feeds {
		user, err := u.GetByUUID(nil, feed.CreatedByUserUUID)

		if err != nil {
			libhttp.HandleBadRequest(w, err)
			return
		}
		result := &GetLatestFeedsStruct{}

		result.Username = user.Username
		result.Message = feed.Message
		result.PokemonName = feed.PokemonName
		result.CreatedAt = feed.CreatedAt.Time
		result.Lat = feed.Lat
		result.Long = feed.Long
		result.FormattedAddress = feed.FormattedAddress
		results = append(results, result)
	}
	if err != nil {
		libhttp.HandleBadRequest(w, err)
		return
	}

	json.NewEncoder(w).Encode(results)
}

func OptionsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func PostFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	decoder := json.NewDecoder(r.Body)
	var t PostFeedStruct
	err := decoder.Decode(&t)

	if err != nil {
		libhttp.HandleBadRequest(w, err)
		return
	}

	db := context.Get(r, "db").(*sqlx.DB)

	feed, err2 := models.NewFeed(db).Create(
		nil,
		t.Message,
		t.PokemonName,
		t.CreatedByUserUUID,
		t.Lat,
		t.Long,
		t.FormattedAddress,
	)

	if err2 != nil {
		if ae, ok := err2.(*pq.Error); ok {
			libhttp.HandlePostgresError(w, *ae)
			return
		}
		libhttp.HandleBadRequest(w, err2)
		return
	}

	if feed == nil {
		libhttp.HandleBadRequest(w, errors.New("Feed does not exist."))
		return
	}

	// Don't return anything to client.  Optimistically displaying.
	response := ResultStruct{
		Result: "PostFeed successful",
	}
	json.NewEncoder(w).Encode(response)
}

func OptionsAllPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func GetAllPokemon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	pokemon := constants.GetMap()

	results := []*GetAllPokemonStruct{}

	// Because pokemon returns a map... we need to grab it out.
	// TODO: how to just return the map anyways.
	// TODO: I guess I should just use IDs...
	for _, p := range pokemon {
		result := &GetAllPokemonStruct{}
		result.ID = p.ID
		result.Name = p.Name
		result.DisplayName = p.DisplayName
		result.ImageURL = p.ImageURL
		results = append(results, result)
	}

	json.NewEncoder(w).Encode(results)
}

func GetAllPokemonDB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	db := context.Get(r, "db").(*sqlx.DB)

	pokemon, err2 := models.NewPokemon(db).GetAll(
		nil,
	)

	if err2 != nil {
		if ae, ok := err2.(*pq.Error); ok {
			libhttp.HandlePostgresError(w, *ae)
			return
		}
		libhttp.HandleBadRequest(w, err2)
		return
	}

	if pokemon == nil {
		libhttp.HandleBadRequest(w, errors.New("pokemon does not exist."))
		return
	}

	results := []*GetAllPokemonStruct{}

	for _, p := range pokemon {
		result := &GetAllPokemonStruct{}
		result.ID = p.ID
		result.Name = p.Name
		// result.ImageURL = p.ImageURL
		results = append(results, result)
	}

	json.NewEncoder(w).Encode(results)
}
