package storage

import "errors"

var (
	NotFoundErr    = errors.New("no such element in storage")
	DuplicateErr   = errors.New("such element already exists")
	OpenDbErr      = errors.New("can't open database by provided connection string")
	UnSpecifiedErr = errors.New("unspecified error while working with db")
)
