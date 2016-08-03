package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libhttp"
	"github.com/pokefeed/pokefeed-api/models"
)

type PostFeedStruct struct {
	CreatedByUserUUID string  `json:"created_by_user_uuid"`
	Username          string  `json:"username"`
	Message           string  `json:"message"`
	Pokemon           string  `json:"pokemon"`
	Lat               float64 `json:"lat"`
	Long              float64 `json:"long"`
	Geocodes          string  `json:"geocodes"` // TODO: can we honor the jsonness here.
	DisplayType       string  `json:"display_type"`
}

type ResultStruct struct {
	Result string `json:"result"`
}

type GetFeedResultStruct struct {
	UUID              string    `json:"uuid"`
	Username          string    `json:"username"`
	CreatedByUserUUID string    `json:"created_by_user_uuid"`
	Message           string    `json:"message"`
	Pokemon           string    `json:"pokemon"`
	PokemonImageURL   string    `json:"pokemon_image_url"`
	Lat               float64   `json:"lat"`
	Long              float64   `json:"long"`
	FormattedAddress  string    `json:"formatted_address"`
	CreatedAtDate     time.Time `json:"created_at_date"`
	UpdatedAtDate     time.Time `json:"updated_at_date"`
	DeletedAtDate     time.Time `json:"deleted_at_date"`
}

func GetFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	result := GetFeedResultStruct{
		UUID:              "a551ebe9-8b11-466f-ad25-797073b05b8b",
		Username:          "ilovepokemon23",
		CreatedByUserUUID: "b89d86f1-5502-4f17-8e68-6945206f2b3c",
		Message:           "Everyone get over here and catch this guy!",
		Pokemon:           "Charmander",
		// Create a map between pokemon name and image url, and then return as base64 image.
		PokemonImageURL:  "http://static.giantbomb.com/uploads/scale_small/0/6087/2438704-1202149925_t.png",
		Lat:              37.7752315,
		Long:             -122.4197165,
		FormattedAddress: "11 Oak St, San Francisco, CA 94102, USA",
		// Use google apis to reverse encode the address
		// http://maps.googleapis.com/maps/api/geocode/json?latlng=37.7752315,-122.4197165&sensor=true
		CreatedAtDate: time.Date(2016, 7, 17, 20, 34, 58, 651387237, time.UTC),
		// TODO: this should return UTC but it does not seem to.
		UpdatedAtDate: time.Now(),
		DeletedAtDate: time.Now(),
	}
	response := []GetFeedResultStruct{result, result, result, result, result}
	json.NewEncoder(w).Encode(response)
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
		t.Pokemon,
		t.CreatedByUserUUID,
		t.Lat,
		t.Long,
		t.Geocodes,
		t.DisplayType,
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
