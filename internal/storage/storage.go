package storage

import (
	"errors"
)

var (
	ErrNotFound = errors.New("bookmark not found")
	ErrExists   = errors.New("bookmark for this url already exists")
)
