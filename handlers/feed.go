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
	"github.com/pokefeed/pokefeed-api/libuuid"
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
	c := models.NewComment(db)
	fi := models.NewFacebookInfo(db)

	feedItems, _ := f.GetByLocation(nil, lat, long, latRadius, longRadius)

	results := []*structs.FeedItemStruct{}

	for _, feedItem := range feedItems {
		user, _ := u.GetByUUID(nil, feedItem.CreatedByUserUUID)
		facebookInfo, _ := fi.GetByUserUUID(nil, feedItem.CreatedByUserUUID)

		feedTags, _ := ft.GetByFeedUUID(nil, feedItem.UUID)
		comments, _ := c.GetByFeedUUID(nil, feedItem.UUID)

		commentsResult := []*structs.CommentStruct{}

		for _, comment := range comments {
			// TODO: include these in the join query man....
			user2, _ := u.GetByUUID(nil, comment.CreatedByUserUUID)
			facebookInfo2, _ := fi.GetByUserUUID(nil, comment.CreatedByUserUUID)

			commentResult := mappers.CommentMapperDBToJSON(comment, user2, facebookInfo2)
			commentsResult = append(commentsResult, &(commentResult))
		}

		feedTagsResult := mappers.FeedTagMapperArrayDBToJSON(feedTags)
		feedItemResult := mappers.FeedItemMapperDBToJSON(
			feedItem,
			user,
			facebookInfo,
			feedTagsResult,
			commentsResult,
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
	c := models.NewComment(db)
	fi := models.NewFacebookInfo(db)

	feedItems, _ := f.GetLatest(nil)

	results := []*structs.FeedItemStruct{}

	for _, feedItem := range feedItems {
		user, _ := u.GetByUUID(nil, feedItem.CreatedByUserUUID)
		facebookInfo, _ := fi.GetByUserUUID(nil, feedItem.CreatedByUserUUID)
		feedTags, _ := ft.GetByFeedUUID(nil, feedItem.UUID)
		comments, _ := c.GetByFeedUUID(nil, feedItem.UUID)

		commentsResult := []*structs.CommentStruct{}

		for _, comment := range comments {
			// TODO: include these in the join query man....
			user2, _ := u.GetByUUID(nil, comment.CreatedByUserUUID)
			facebookInfo2, _ := fi.GetByUserUUID(nil, comment.CreatedByUserUUID)
			commentResult := mappers.CommentMapperDBToJSON(comment, user2, facebookInfo2)
			commentsResult = append(commentsResult, &(commentResult))
		}

		feedTagsResult := mappers.FeedTagMapperArrayDBToJSON(feedTags)
		feedItemResult := mappers.FeedItemMapperDBToJSON(
			feedItem,
			user,
			facebookInfo,
			feedTagsResult,
			commentsResult,
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

	if len(t.UUID) > 0 && !libuuid.ValidateUUIDv4(t.UUID) {
		response := structs.ResultStruct{
			Result: "PostFeed failed.  Given UUID is not valid.",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	db := context.Get(r, "db").(*sqlx.DB)

	feedItem, err2 := models.NewFeedItem(db).Create(
		nil,
		t.UUID,
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

// PostComment method
func PostComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	decoder := json.NewDecoder(r.Body)
	var t structs.CommentStruct
	err := decoder.Decode(&t)

	if err != nil {
		libhttp.HandleBadRequest(w, err)
		return
	}

	db := context.Get(r, "db").(*sqlx.DB)

	comment, err2 := models.NewComment(db).Create(
		nil,
		t.FeedItemUUID,
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

	if comment == nil {
		libhttp.HandleBadRequest(w, errors.New("Comment does not exist."))
		return
	}

	// Don't return anything to client.  Optimistically displaying.
	response := structs.ResultStruct{
		Result: "PostComment successful",
	}
	json.NewEncoder(w).Encode(response)
}
