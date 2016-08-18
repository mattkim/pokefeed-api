package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libuuid"
)

func NewFacebookInfo(db *sqlx.DB) *FacebookInfo {
	facebookInfo := &FacebookInfo{}
	facebookInfo.db = db
	facebookInfo.table = "facebook_info"
	facebookInfo.hasUUID = true

	return facebookInfo
}

type FacebookInfoRow struct {
	UUID         string      `db:"uuid"`
	FacebookID   string      `db:"facebook_id"`
	FacebookName string      `db:"facebook_name"`
	UserUUID     string      `db:"user_uuid"`
	Email        string      `db:"email"`
	UpdatedAt    pq.NullTime `db:"updated_at"`
	CreatedAt    pq.NullTime `db:"created_at"`
	DeletedAt    pq.NullTime `db:"deleted_at"`
}

type FacebookInfo struct {
	BaseUUID
}

// GetById returns record by id.
func (f *FacebookInfo) GetByUUID(tx *sqlx.Tx, uuid string) (*FacebookInfoRow, error) {
	FacebookInfo := &FacebookInfoRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE uuid=$1", f.table)
	err := f.db.Get(FacebookInfo, query, uuid)

	return FacebookInfo, err
}

func (f *FacebookInfo) facebookInfoRowFromSqlResult(tx *sqlx.Tx, sqlResult InsertResultUUID) (*FacebookInfoRow, error) {
	facebookInfoUUID, err := sqlResult.LastInsertUUID()
	if err != nil {
		return nil, err
	}

	return f.GetByUUID(tx, facebookInfoUUID)
}

// CreateFacebookUser create a new record of user the facebook way
func (f *FacebookInfo) CreateFacebookUser(tx *sqlx.Tx, userUUID, facebookID, facebookName string) (*FacebookInfoRow, error) {
	if userUUID == "" {
		return nil, errors.New("userUUID cannot be blank.")
	}
	if facebookID == "" {
		return nil, errors.New("FacebookID cannot be blank.")
	}
	if facebookName == "" {
		return nil, errors.New("FacebookName cannot be blank.")
	}

	now := time.Now().UTC()

	data := make(map[string]interface{})
	// TODO: make username and password nullable.
	// Can we have a nullable unique index?
	data["uuid"], _ = libuuid.NewUUID()
	data["user_uuid"] = userUUID
	data["facebook_id"] = facebookID
	data["facebook_name"] = facebookName
	data["updated_at"] = now
	data["created_at"] = now

	sqlResult, err := f.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return f.facebookInfoRowFromSqlResult(tx, *sqlResult)
}

func (f *FacebookInfo) GetByFacebookID(tx *sqlx.Tx, facebook_id string) (*FacebookInfoRow, error) {
	facebookInfo := &FacebookInfoRow{}
	query := fmt.Sprintf(
		`SELECT * FROM %v as f
		WHERE f.facebook_id=$1`,
		f.table,
	)
	err := f.db.Get(facebookInfo, query, facebook_id)

	return facebookInfo, err
}

func (f *FacebookInfo) GetByUserUUID(tx *sqlx.Tx, user_uuid string) (*FacebookInfoRow, error) {
	facebookInfo := &FacebookInfoRow{}
	query := fmt.Sprintf(
		`SELECT * FROM %v as f
		WHERE f.user_uuid=$1`,
		f.table,
	)
	err := f.db.Get(facebookInfo, query, user_uuid)

	return facebookInfo, err
}
