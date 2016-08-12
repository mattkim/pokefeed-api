package models

import (
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

func NewFeedTag(db *sqlx.DB) *FeedTag {
	feedTag := &FeedTag{}
	feedTag.db = db
	feedTag.table = "feed_tags"
	feedTag.hasUUID = true

	return feedTag
}

type FeedTagRow struct {
	UUID        string      `db:"uuid"`
	Type        string      `db:"type"`
	Name        string      `db:"name"`
	DisplayName string      `db:"display_name"`
	ImageURL    string      `db:"image_url"`
	UpdatedAt   pq.NullTime `db:"updated_at"`
	CreatedAt   pq.NullTime `db:"created_at"`
	DeletedAt   pq.NullTime `db:"deleted_at"`
}

type FeedTag struct {
	BaseUUID
}

// TODO: test this select join works
func (ft *FeedTag) GetByFeedUUID(tx *sqlx.Tx, feedUUID string) ([]FeedTagRow, error) {
	feedTags := []FeedTagRow{}
	query := `SELECT ft.*
		FROM feed_tags as ft
		JOIN feed_items_feed_tags as fift on fift.feed_tag_uuid = ft.uuid
		WHERE fift.feed_item_uuid=$1`
	err := ft.db.Select(&feedTags, query, feedUUID)
	return feedTags, err
}

func (f *FeedTag) GetAll(
	tx *sqlx.Tx,
) ([]FeedTagRow, error) {
	feedTags := []FeedTagRow{}

	query := `SELECT * FROM feed_tags ORDER BY type, name`
	err := f.db.Select(&feedTags, query)

	if err != nil {
		log.Fatal(err)
	}

	return feedTags, err
}

// func (ft *FeedTag) Create(
// 	tx *sqlx.Tx,
// 	feedItemUUID string,
// 	feedTagUUID string,
// ) error {
// 	now := time.Now().UTC()
// 	data2 := make(map[string]interface{})
// 	data2["feed_item_uuid"] = feedItemUUID
// 	data2["feed_tag_uuid"] = feedTagUUID
// 	data2["updated_at"] = now
// 	data2["created_at"] = now
//
// 	_, err2 := ft.InsertIntoTable(tx, data2)
// 	if err2 != nil {
// 		return err2
// 	}
//
// 	return nil
// }
