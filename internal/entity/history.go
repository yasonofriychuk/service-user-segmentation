package entity

import "time"

type History struct {
	UserID      string    `db:"user_id"`
	SegmentSlug string    `db:"segment_slug"`
	Type        string    `db:"type"`
	CreatedAt   time.Time `db:"created_at"`
}
