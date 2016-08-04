package models

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pokefeed/pokefeed-api/libuuid"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(db *sqlx.DB) *User {
	user := &User{}
	user.db = db
	user.table = "users"
	user.hasUUID = true

	return user
}

type UserRow struct {
	UUID      string      `db:"uuid"`
	Username  string      `db:"username"`
	Email     string      `db:"email"`
	Password  string      `db:"password"`
	UpdatedAt pq.NullTime `db:"updated_at"`
	CreatedAt pq.NullTime `db:"created_at"`
	DeletedAt pq.NullTime `db:"deleted_at"`
}

type User struct {
	BaseUUID
}

func (u *User) userRowFromSqlResult(tx *sqlx.Tx, sqlResult InsertResultUUID) (*UserRow, error) {
	userUUID, err := sqlResult.LastInsertUUID()
	if err != nil {
		return nil, err
	}

	return u.GetByUUID(tx, userUUID)
}

// GetById returns record by id.
func (u *User) GetByUUID(tx *sqlx.Tx, uuid string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE uuid=$1", u.table)
	err := u.db.Get(user, query, uuid)

	return user, err
}

// GetByEmail returns record by email.
func (u *User) GetByEmail(tx *sqlx.Tx, email string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE email=$1", u.table)
	err := u.db.Get(user, query, email)

	return user, err
}

// GetByEmail returns record by email but checks password first.
func (u *User) GetUserByEmailAndPassword(tx *sqlx.Tx, email, password string) (*UserRow, error) {
	user, err := u.GetByEmail(tx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, err
}

// Signup create a new record of user.
func (u *User) Signup(tx *sqlx.Tx, email, username, password, passwordAgain string) (*UserRow, error) {
	if email == "" {
		return nil, errors.New("Email cannot be blank.")
	}
	match, _ := regexp.MatchString(".+@.+\\..+", email)
	if !match {
		return nil, errors.New("Not a valid email format.")
	}
	if username == "" {
		return nil, errors.New("Username cannot be blank.")
	}
	if password == "" {
		return nil, errors.New("Password cannot be blank.")
	}
	if password != passwordAgain {
		return nil, errors.New("Passwords do not match.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return nil, err
	}

	// TODO: move into base_UUID
	// TODO: also make sure this is utc
	now := time.Now()

	data := make(map[string]interface{})
	// TODO: ignoring potential error here.
	data["uuid"], _ = libuuid.NewUUID()
	data["email"] = email
	data["username"] = username
	data["password"] = hashedPassword
	data["updated_at"] = now
	data["created_at"] = now

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.userRowFromSqlResult(tx, *sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (u *User) UpdateEmailAndPasswordByUUID(tx *sqlx.Tx, userUUID string, email, password, passwordAgain string) (*UserRow, error) {
	data := make(map[string]interface{})

	if email != "" {
		data["email"] = email
	}

	if password != "" && passwordAgain != "" && password == passwordAgain {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			return nil, err
		}

		data["password"] = hashedPassword
	}

	if len(data) > 0 {
		_, err := u.UpdateByUUID(tx, data, userUUID)
		if err != nil {
			return nil, err
		}
	}

	return u.GetByUUID(tx, userUUID)
}
