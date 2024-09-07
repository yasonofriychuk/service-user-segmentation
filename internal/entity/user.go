package entity

import "time"

type User struct {
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
