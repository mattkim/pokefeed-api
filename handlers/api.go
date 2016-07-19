package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pokefeed/pokefeed-api/libhttp"
)

type PostFeedStruct struct {
	CreatedByUserUUID string  `json:"created_by_user_uuid"`
	Message           string  `json:"message"`
	Pokemon           string  `json:"pokemon"`
	Lat               float64 `json:"lat"`
	Long              float64 `json:"long"`
}

type ResultStruct struct {
	ID     int32  `json:"id"`
	Result string `json:"result"`
}

type GetFeedResultStruct struct {
	UUID              string    `json:"uuid"`
	CreatedByUserUUID string    `json:"created_by_user_uuid"`
	Message           string    `json:"message"`
	Pokemon           string    `json:"pokemon"`
	Lat               float64   `json:"lat"`
	Long              float64   `json:"long"`
	CreatedAtDate     time.Time `json:"created_at_date"`
	UpdatedAtDate     time.Time `json:"updated_at_date"`
	DeletedAtDate     time.Time `json:"deleted_at_date"`
}

func GetFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Note adding an ID is recommended by reactjs
	result := GetFeedResultStruct{
		UUID:              "a551ebe9-8b11-466f-ad25-797073b05b8b",
		CreatedByUserUUID: "b89d86f1-5502-4f17-8e68-6945206f2b3c",
		Message:           "Message",
		Pokemon:           "Charmander",
		Lat:               51.5032510,
		Long:              51.5032510,
		CreatedAtDate:     time.Now(),
		UpdatedAtDate:     time.Now(),
		DeletedAtDate:     time.Now(),
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
		libhttp.HandleErrorJson(w, err)
		return
	}

	response := ResultStruct{
		ID:     1,
		Result: fmt.Sprintf("PostFeed: %v", t),
	}
	json.NewEncoder(w).Encode(response)
}
