package memory

import "errors"

var (
	ErrNotFound  = errors.New("entity not found")
	ErrDuplicate = errors.New("entity duplicated by UUID")
)
