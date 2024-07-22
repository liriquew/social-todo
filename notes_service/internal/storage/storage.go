package storage

import "fmt"

var (
	ErrNotFound      = fmt.Errorf("note not found")
	ErrAlreadyExists = fmt.Errorf("note with that title already exists")
)
