package store

import (
	"errors"

	"github.com/iv-menshenin/accountant/store/internal/memory"
)

var (
	ErrNotFound  = errors.New("entity not found")
	ErrDuplicate = errors.New("duplicate entity")
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if err == memory.ErrNotFound {
		return ErrNotFound
	}
	if err == memory.ErrDuplicate {
		return ErrDuplicate
	}
	return err
}
