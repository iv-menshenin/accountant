package storage

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/iv-menshenin/accountant/storage/internal/memory"
)

var (
	ErrNotFound  = errors.New("entity not found")
	ErrDuplicate = errors.New("duplicate entity")
)

func MapMemoryErrors(err error) error {
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

func MapMongodbErrors(err error) error {
	if err == nil {
		return nil
	}
	if err == mongo.ErrNoDocuments {
		return ErrNotFound
	}
	if err == mongo.ErrInvalidIndexValue {
		return ErrDuplicate
	}
	return err
}
