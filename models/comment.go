package models

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libuuid"
)

// NewComment constructor
func NewComment(db *sqlx.DB) *Comment {
	comment := &Comment{}
	comment.db = db
	comment.table = "comments"
	comment.hasUUID = true

	return comment
}

type CommentRow struct {
	UUID              string      `db:"uuid"`
	FeedItemUUID      string      `db:"feed_item_uuid"`
	Message           string      `db:"message"`
	CreatedByUserUUID string      `db:"created_by_user_uuid"`
	Lat               float64     `db:"lat"`
	Long              float64     `db:"long"`
	FormattedAddress  string      `db:"formatted_address"`
	UpdatedAt         pq.NullTime `db:"updated_at"`
	CreatedAt         pq.NullTime `db:"created_at"`
	DeletedAt         pq.NullTime `db:"deleted_at"`
}

type Comment struct {
	BaseUUID
}

func (c *Comment) commentRowFromSqlResult(tx *sqlx.Tx, sqlResult InsertResultUUID) (*CommentRow, error) {
	commentUUID, err := sqlResult.LastInsertUUID()
	if err != nil {
		return nil, err
	}

	return c.GetByUUID(tx, commentUUID)
}

// GetByUUID returns record by uuid.
func (c *Comment) GetByUUID(tx *sqlx.Tx, uuid string) (*CommentRow, error) {
	comment := &CommentRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE uuid=$1", c.table)
	err := c.db.Get(comment, query, uuid)

	return comment, err
}

// GetByFeedUUID comments
func (c *Comment) GetByFeedUUID(tx *sqlx.Tx, feedUUID string) ([]CommentRow, error) {
	comments := []CommentRow{}
	query := `SELECT c.*
		FROM comments as c
		WHERE c.feed_item_uuid=$1`
	err := c.db.Select(&comments, query, feedUUID)
	return comments, err
}

// Create create a new record of comment.
func (c *Comment) Create(
	tx *sqlx.Tx,
	feedItemUUID string,
	message string,
	createdByUserUUID string,
	lat float64,
	long float64,
	formattedAddress string,
) (*CommentRow, error) {
	now := time.Now().UTC()
	uuid, _ := libuuid.NewUUID()
	data := make(map[string]interface{})
	data["uuid"] = uuid
	data["feed_item_uuid"] = feedItemUUID
	data["message"] = message
	data["updated_at"] = now
	data["created_at"] = now
	data["created_by_user_uuid"] = createdByUserUUID
	data["lat"] = lat
	data["long"] = long
	data["formatted_address"] = formattedAddress

	sqlResult, err := c.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	// TODO: not sure if this pointer business makes sense.
	return c.commentRowFromSqlResult(tx, *sqlResult)
}
