package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// NewFeed constructor
func NewFeed(db *sqlx.DB) *Feed {
	feed := &Feed{}
	feed.db = db
	feed.table = "feeds"
	feed.hasUUID = true

	return feed
}

type FeedRow struct {
	// TODO: ID should be UUID
	UUID              string    `db:"uuid"`
	Message           string    `db:"message"`
	Pokemon           string    `db:"pokemon"`
	CreatedByUserUUID string    `db:"created_by_user_uuid"`
	Lat               float64   `db:"lat"`
	Long              float64   `db:"long"`
	Geocodes          string    `db:"geocodes"` // TODO: can this be json type or map?
	DisplayType       string    `db:"display_type"`
	UpdatedAt         time.Time `db:"updated_at"`
	CreatedAt         time.Time `db:"created_at"`
	DeletedAt         time.Time `db:"deleted_at"`
}

type Feed struct {
	BaseUUID
}

func (f *Feed) feedRowFromSqlResult(tx *sqlx.Tx, sqlResult InsertResultUUID) (*FeedRow, error) {
	// TODO: need to change this to return uuid.
	feedUUID, err := sqlResult.LastInsertUUID()
	if err != nil {
		return nil, err
	}

	return f.GetByUUID(tx, feedUUID)
}

// GetById returns record by id.
func (f *Feed) GetByUUID(tx *sqlx.Tx, uuid string) (*FeedRow, error) {
	feed := &FeedRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE uuid=$1", f.table)
	err := f.db.Get(feed, query, uuid)

	return feed, err
}

// GetByEmail returns record by email.  Needs to really be a bounding box right.
func (f *Feed) GetByLocation(tx *sqlx.Tx, latTop float64, longLeft float64, latBot float64, longRight float64) (*FeedRow, error) {
	feed := &FeedRow{}
	// Query for feeds within this bounding box.
	query := fmt.Sprintf("SELECT * FROM %v WHERE lat > $1 and lat < $2 and long > $3 and long < $4", f.table)
	err := f.db.Get(feed, query, latTop, latBot, longLeft, longRight)

	return feed, err
}

// Signup create a new record of feed.
func (f *Feed) Create(
	tx *sqlx.Tx,
	message string,
	pokemon string,
	createdByUserUUID string,
	lat float64,
	long float64,
	geocodes string,
	displayType string,
) (*FeedRow, error) {
	now := time.Now()

	data := make(map[string]interface{})
	data["message"] = message
	data["pokemon"] = pokemon
	data["updated_at"] = now
	data["created_at"] = now
	data["created_by_user_uuid"] = createdByUserUUID
	data["lat"] = lat
	data["long"] = long
	data["geocodes"] = geocodes
	data["display_type"] = displayType

	sqlResult, err := f.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	// TODO: not sure if this pointer business makes sense.
	return f.feedRowFromSqlResult(tx, *sqlResult)
}
