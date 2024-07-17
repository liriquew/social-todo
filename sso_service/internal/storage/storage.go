package storage

import "fmt"

var (
	ErrNotFound  = fmt.Errorf("user not found")
	ErrUserExist = fmt.Errorf("user exist")
)
