package repoerrs

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyExists    = errors.New("already exists")
	ErrSegmentsNotExist = errors.New("one of the segments does not exist")
	ErrUserNotFound     = errors.New("user not found")
)
