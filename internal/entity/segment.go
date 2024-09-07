package entity

import "time"

type Segment struct {
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
}

const (
	OperationTypeAdd           = "add"
	OperationTypeDelete        = "delete"
	OperationTypeAutoAdd       = "auto_add"
	OperationTypeSegmentDelete = "delete_segment"
)
