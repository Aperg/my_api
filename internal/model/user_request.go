package model

import (
	"database/sql"
	"time"
)

// UserRequest is a request for equipment
type UserRequest struct {
	ID_user   uint64       `db:"id_user"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
	DoneAt    sql.NullTime `db:"done_at"`
}
