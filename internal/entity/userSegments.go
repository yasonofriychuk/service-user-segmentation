package entity

type UserSegments struct {
	UserID      string `db:"user_id"`
	SegmentSlug string `db:"segment_slug"`
}
