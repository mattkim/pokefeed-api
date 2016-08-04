package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libhttp"
	"github.com/pokefeed/pokefeed-api/models"
)

type PostFeedStruct struct {
	CreatedByUserUUID string         `json:"created_by_user_uuid"`
	Username          string         `json:"username"`
	Message           string         `json:"message"`
	Pokemon           string         `json:"pokemon"`
	Lat               float64        `json:"lat"`
	Long              float64        `json:"long"`
	Geocodes          types.JSONText `json:"geocodes"` // TODO: can we honor the jsonness here.
	DisplayType       string         `json:"display_type"`
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

type GetLatestFeedsStruct struct {
	Username         string    `json:"username"`
	Message          string    `json:"message"`
	Pokemon          string    `json:"pokemon"`
	PokemonImageURL  string    `json:"pokemon_image_url"` // Should I fetch this from backend
	Lat              float64   `json:"lat"`
	Long             float64   `json:"long"`
	FormattedAddress string    `json:"formatted_address"`
	CreatedAt        time.Time `json:"created_at"`
}

func GetLatestFeeds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

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

		// Find the right geocode here
		var goodg map[string]interface{}
		// var f []interface{}
		var gs []interface{}
		// json.Unmarshal(feed.Geocodes, &f)
		json.Unmarshal(feed.Geocodes, &gs)
		if gs == nil {
			// TODO: this is super weird, but sometimes we cannot unmarshal geocodes even though it exists.
			Info.Println("**** skipping")
			Info.Println(fmt.Sprintf("%+v\n", feed))
			Info.Println(feed)
		} else {
			// gs := f.([]interface{})
			goodg = gs[0].(map[string]interface{}) // default to the first one.
			for _, g := range gs {
				gn := g.(map[string]interface{})
				gnTypes := gn["types"].([]interface{})
				for _, t := range gnTypes {
					if t == feed.DisplayType {
						goodg = gn
					}
				}
			}

			Info.Println(goodg)
			// Fetch the formatted address and lat long here.
			formattedAddress := goodg["formatted_address"].(string)
			lat := goodg["geometry"].(map[string]interface{})["location"].(map[string]interface{})["lat"].(float64)
			long := goodg["geometry"].(map[string]interface{})["location"].(map[string]interface{})["lng"].(float64)

			result.Username = user.Username
			result.Message = feed.Message
			result.Pokemon = feed.Pokemon
			// TODO: fetch url from map.
			result.CreatedAt = feed.CreatedAt.Time
			result.Lat = lat
			result.Long = long
			result.FormattedAddress = formattedAddress
			results = append(results, result)
		}
	}
	if err != nil {
		libhttp.HandleBadRequest(w, err)
		return
	}

	json.NewEncoder(w).Encode(results)
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
