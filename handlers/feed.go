package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libhttp"
	"github.com/pokefeed/pokefeed-api/mappers"
	"github.com/pokefeed/pokefeed-api/models"
	"github.com/pokefeed/pokefeed-api/structs"
)

// GetFeeds by location
func GetFeeds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	lat, _ := strconv.ParseFloat(r.URL.Query()["lat"][0], 64)
	long, _ := strconv.ParseFloat(r.URL.Query()["long"][0], 64)
	// TODO: change json param to underscores?
	latRadius, _ := strconv.ParseFloat(r.URL.Query()["latRadius"][0], 64)
	longRadius, _ := strconv.ParseFloat(r.URL.Query()["longRadius"][0], 64)

	db := context.Get(r, "db").(*sqlx.DB)
	f := models.NewFeedItem(db)
	u := models.NewUser(db)
	ft := models.NewFeedTag(db)

	feedItems, _ := f.GetByLocation(nil, lat, long, latRadius, longRadius)

	results := []*structs.FeedItemStruct{}

	for _, feedItem := range feedItems {
		user, _ := u.GetByUUID(nil, feedItem.CreatedByUserUUID)
		feedTags, _ := ft.GetByFeedUUID(nil, feedItem.UUID)

		// TODO is pointer to struct necessary?
		feedTagsResult := mappers.FeedTagMapperArrayDBToJSON(feedTags)
		feedItemResult := mappers.FeedItemMapperDBToJSON(
			feedItem,
			user,
			feedTagsResult,
		)

		results = append(results, &feedItemResult)
	}

	json.NewEncoder(w).Encode(results)
}

// GetLatestFeeds order by date
func GetLatestFeeds(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	db := context.Get(r, "db").(*sqlx.DB)
	f := models.NewFeedItem(db)
	u := models.NewUser(db)
	ft := models.NewFeedTag(db)

	feedItems, _ := f.GetLatest(nil)

	results := []*structs.FeedItemStruct{}

	for _, feedItem := range feedItems {
		user, _ := u.GetByUUID(nil, feedItem.CreatedByUserUUID)
		feedTags, _ := ft.GetByFeedUUID(nil, feedItem.UUID)

		// TODO is pointer to struct necessary?
		feedTagsResult := mappers.FeedTagMapperArrayDBToJSON(feedTags)
		feedItemResult := mappers.FeedItemMapperDBToJSON(
			feedItem,
			user,
			feedTagsResult,
		)

		results = append(results, &feedItemResult)
	}

	json.NewEncoder(w).Encode(results)
}

// PostFeed method
func PostFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	decoder := json.NewDecoder(r.Body)
	var t structs.PostFeedStruct
	err := decoder.Decode(&t)

	if err != nil {
		libhttp.HandleBadRequest(w, err)
		return
	}

	if len(t.FeedTags) > 5 {
		response := structs.ResultStruct{
			Result: "PostFeed failed.  Feed Tags is greater than 5.",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	db := context.Get(r, "db").(*sqlx.DB)

	feedItem, err2 := models.NewFeedItem(db).Create(
		nil,
		t.Message,
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

	if feedItem == nil {
		libhttp.HandleBadRequest(w, errors.New("FeedItem does not exist."))
		return
	}

	for _, feedTag := range t.FeedTags {
		models.NewFeedItemFeedTag(db).Create(
			nil,
			feedItem.UUID,
			feedTag.UUID,
		)
	}

	// Don't return anything to client.  Optimistically displaying.
	response := structs.ResultStruct{
		Result: "PostFeed successful",
	}
	json.NewEncoder(w).Encode(response)
}

func GetAllFeedTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	db := context.Get(r, "db").(*sqlx.DB)

	feedTags, err := models.NewFeedTag(db).GetAll(nil)

	if err != nil {
		libhttp.HandleBadRequest(w, err)
		return
	}

	feedTagsResult := mappers.FeedTagMapperArrayDBToJSON(feedTags)

	json.NewEncoder(w).Encode(feedTagsResult)
}
