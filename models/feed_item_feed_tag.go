package models

import (
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

func NewFeedItemFeedTag(db *sqlx.DB) *FeedItemFeedTag {
	feedItemFeedTag := &FeedItemFeedTag{}
	feedItemFeedTag.db = db
	feedItemFeedTag.table = "feed_items_feed_tags"
	feedItemFeedTag.hasUUID = false

	return feedItemFeedTag
}

type FeedItemFeedTagRow struct {
	FeedItemUUID string      `db:"feed_item_uuid"`
	FeedTagUUID  string      `db:"feed_tag_uuid"`
	UpdatedAt    pq.NullTime `db:"updated_at"`
	CreatedAt    pq.NullTime `db:"created_at"`
	DeletedAt    pq.NullTime `db:"deleted_at"`
}

type FeedItemFeedTag struct {
	BaseUUID
}

func (ft *FeedItemFeedTag) Create(
	tx *sqlx.Tx,
	feedItemUUID string,
	feedTagUUID string,
) error {
	now := time.Now().UTC()
	data2 := make(map[string]interface{})
	data2["feed_item_uuid"] = feedItemUUID
	data2["feed_tag_uuid"] = feedTagUUID
	data2["updated_at"] = now
	data2["created_at"] = now

	_, err2 := ft.InsertIntoTable(tx, data2)
	if err2 != nil {
		return err2
	}

	return nil
}
