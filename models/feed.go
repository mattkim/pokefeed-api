package models

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"

	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libuuid"
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
	UUID              string         `db:"uuid"`
	Message           string         `db:"message"`
	Pokemon           string         `db:"pokemon"`
	CreatedByUserUUID string         `db:"created_by_user_uuid"`
	Lat               float64        `db:"lat"`
	Long              float64        `db:"long"`
	Geocodes          types.JSONText `db:"geocodes"` // TODO: can this be json type or map?
	DisplayType       string         `db:"display_type"`
	UpdatedAt         pq.NullTime    `db:"updated_at"`
	CreatedAt         pq.NullTime    `db:"created_at"`
	DeletedAt         pq.NullTime    `db:"deleted_at"`
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

// type GetLatestFeedsStruct struct {
// 	Username         string    `json:"username"`
// 	Message          string    `json:"message"`
// 	Pokemon          string    `json:"pokemon"`
// 	PokemonImageURL  string    `json:"pokemon_image_url"` // Should I fetch this from backend
// 	Lat              float64   `json:"lat"`
// 	Long             float64   `json:"long"`
// 	FormattedAddress string    `json:"formatted_address"`
// 	CreatedAtDate    time.Time `json:"created_at_date"`
// 	UpdatedAtDate    time.Time `json:"updated_at_date"`
// }

// GetLatest returns record by email.  Needs to really be a bounding box right.
func (f *Feed) GetLatest(tx *sqlx.Tx) ([]*FeedRow, error) {
	feeds := []*FeedRow{}

	// feed := *FeedRow{}
	// Query for feeds within this bounding box.
	// query := "SELECT u.username, f.message, f.pokemon, f.lat, f.long, f.geocodes, f.display_type, f.created_at" +
	// 	"FROM feeds as f" +
	// 	"JOIN users as u on u.uuid = f.created_by_user_uuid" + // TODO: users is hard-coded.
	// 	"ORDER BY created_at DESC limit 100"
	var (
		createdByUserUUID string
		message           string
		pokemon           string
		lat               float64
		long              float64
		geocodes          types.JSONText
		displayType       *string
		createdAt         pq.NullTime
	)
	query := "SELECT f.created_by_user_uuid, f.message, f.pokemon, f.lat, f.long, f.geocodes, f.display_type, f.created_at FROM feeds as f ORDER BY f.created_at DESC LIMIT 100"
	rows, err := f.db.Query(query)

	// Info := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	for rows.Next() {
		feed := &FeedRow{}
		if err := rows.Scan(
			&createdByUserUUID,
			&message,
			&pokemon,
			&lat,
			&long,
			&geocodes,
			&displayType,
			&createdAt,
		); err != nil {
			log.Fatal(err)
		}
		feed.CreatedByUserUUID = createdByUserUUID
		feed.Message = message
		feed.Pokemon = pokemon
		feed.Lat = lat
		feed.Long = long
		feed.Geocodes = geocodes

		if displayType != nil {
			feed.DisplayType = *displayType
		}

		feed.CreatedAt = createdAt
		// Info.Println(feed)
		feeds = append(feeds, feed)
	}

	// Info.Println(feeds)

	if err != nil {
		log.Fatal(err)
	}

	return feeds, err
}

// GetByLocation returns record by email.  Needs to really be a bounding box right.
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
	geocodes types.JSONText, // TODO: double check this works?
	displayType string,
) (*FeedRow, error) {
	now := time.Now()

	data := make(map[string]interface{})
	data["uuid"], _ = libuuid.NewUUID()
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
