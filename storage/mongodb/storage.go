package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"

	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/storage/internal/mongodb"
)

type (
	Storage struct {
		logger *log.Logger
		mongo  *mongodb.Database
	}
	Collection struct {
		storage   *mongo.Collection
		logger    *log.Logger
		mapErrorF func(error) error
	}
)

func (c *Collection) mapError(err error) error {
	if c.mapErrorF != nil {
		err = c.mapErrorF(err)
	}
	if err != nil {
		c.logger.Printf("DB ERROR: %v\n", err)
	}
	return err
}

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

func NewTestStorage(config *config.ConfigStorage, logger *log.Logger) (*Storage, error) {
	db, err := mongodb.New(config, logger)
	if err != nil {
		return nil, err
	}

	db.Accounts().Collection.DeleteMany(context.Background(), bson.D{})
	db.Bills().Collection.DeleteMany(context.Background(), bson.D{})
	db.Payments().Collection.DeleteMany(context.Background(), bson.D{})
	db.Targets().Collection.DeleteMany(context.Background(), bson.D{})
	db.Users().Collection.DeleteMany(context.Background(), bson.D{})

	return &Storage{
		logger: logger,
		mongo:  db,
	}, nil
}

func (s *Storage) Close() error {
	return s.mongo.Close()
}
