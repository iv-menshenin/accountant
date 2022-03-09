package mongodb

import (
	"log"

	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/storage/internal/mongodb"
)

type (
	Storage struct {
		logger *log.Logger
		mongo  *mongodb.Database
	}
)

func NewStorage(config *config.ConfigStorage, logger *log.Logger) (*Storage, error) {
	db, err := mongodb.New(config, logger)
	if err != nil {
		return nil, err
	}
	return &Storage{
		logger: logger,
		mongo:  db,
	}, nil
}

func (s *Storage) Close() error {
	return s.mongo.Close()
}
