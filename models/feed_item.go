package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libuuid"
)

// NewFeed constructor
func NewFeedItem(db *sqlx.DB) *FeedItem {
	feedItem := &FeedItem{}
	feedItem.db = db
	feedItem.table = "feed_items"
	feedItem.hasUUID = true

	return feedItem
}

type FeedItemRow struct {
	UUID              string      `db:"uuid"`
	Message           string      `db:"message"`
	CreatedByUserUUID string      `db:"created_by_user_uuid"`
	Lat               float64     `db:"lat"`
	Long              float64     `db:"long"`
	FormattedAddress  string      `db:"formatted_address"`
	UpdatedAt         pq.NullTime `db:"updated_at"`
	CreatedAt         pq.NullTime `db:"created_at"`
	DeletedAt         pq.NullTime `db:"deleted_at"`
}

type FeedItem struct {
	BaseUUID
}

func (f *FeedItem) feedRowFromSqlResult(tx *sqlx.Tx, sqlResult InsertResultUUID) (*FeedItemRow, error) {
	feedItemUUID, err := sqlResult.LastInsertUUID()
	if err != nil {
		return nil, err
	}

	return f.GetByUUID(tx, feedItemUUID)
}

// GetByUUID returns record by uuid.
func (f *FeedItem) GetByUUID(tx *sqlx.Tx, uuid string) (*FeedItemRow, error) {
	feedItem := &FeedItemRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE uuid=$1", f.table)
	err := f.db.Get(feedItem, query, uuid)

	return feedItem, err
}

func (f *FeedItem) GetByLocation(
	tx *sqlx.Tx,
	lat float64,
	long float64,
	latRadius float64,
	longRadius float64,
) ([]FeedItemRow, error) {
	feedItems := []FeedItemRow{}

	// TODO: how to do paging.
	query := `SELECT *
	FROM feed_items as f
	WHERE f.lat > $1
	AND f.lat < $2
	AND f.long > $3
	AND f.long < $4
	ORDER BY f.created_at DESC
	LIMIT 50`

	err := f.db.Select(
		&feedItems,
		query,
		lat-latRadius,
		lat+latRadius,
		long-longRadius,
		long+longRadius,
	)

	if err != nil {
		log.Fatal(err)
	}

	return feedItems, err
}

// GetLatest returns record by email.
func (f *FeedItem) GetLatest(tx *sqlx.Tx) ([]FeedItemRow, error) {
	feedItems := []FeedItemRow{}
	query := "SELECT * FROM feed_items as f ORDER BY f.created_at DESC LIMIT 50"
	err := f.db.Select(&feedItems, query)

	if err != nil {
		log.Fatal(err)
	}

	return feedItems, err
}

// Create create a new record of feedItem.
func (f *FeedItem) Create(
	tx *sqlx.Tx,
	uuid string,
	message string,
	createdByUserUUID string,
	lat float64,
	long float64,
	formattedAddress string,
) (*FeedItemRow, error) {
	now := time.Now().UTC()

	if len(uuid) > 0 {
		if !libuuid.ValidateUUIDv4(uuid) {
			return nil, errors.New("UUID v4 validation fails")
		}
	} else {
		// Assign a new uuid if it doesn't already exist
		newUUID, _ := libuuid.NewUUID()
		uuid = newUUID
	}

	data := make(map[string]interface{})
	data["uuid"] = uuid
	data["message"] = message
	data["updated_at"] = now
	data["created_at"] = now
	data["created_by_user_uuid"] = createdByUserUUID
	data["lat"] = lat
	data["long"] = long
	data["formatted_address"] = formattedAddress

	sqlResult, err := f.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	// TODO: not sure if this pointer business makes sense.
	return f.feedRowFromSqlResult(tx, *sqlResult)
}
