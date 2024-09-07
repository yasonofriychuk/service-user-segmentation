package service

import "fmt"

var (
	ErrSegmentAlreadyExists = fmt.Errorf("segment already exists")
	ErrCannotCreateSegment  = fmt.Errorf("cannot create segment")
	ErrSegmentNotFound      = fmt.Errorf("segment not found")
	ErrUserNotFound         = fmt.Errorf("user not found")
	ErrUserNoData           = fmt.Errorf("this user has no data")
)
