package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/utils/uuid"
	mid "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type (
	UsersCollection Collection
	userRecord      struct {
		ID       mid.UUID            `bson:"_id"`
		Data     domain.UserInfo     `bson:"data"`
		Identity domain.UserIdentity `bson:"identity"`
		Created  time.Time           `bson:"created"`
		Updated  time.Time           `bson:"updated"`
		Deleted  *time.Time          `bson:"deleted"`
	}
)

func (c *UsersCollection) mapError(err error) error {
	return (*Collection)(c).mapError(err)
}

func (c *UsersCollection) Create(ctx context.Context, info domain.UserInfo, identity domain.UserIdentity) (*domain.UserInfo, error) {
	_, err := c.storage.InsertOne(
		ctx,
		newRecord(info, identity),
		options.InsertOne(),
	)
	if err != nil {
		return nil, c.mapError(err)
	}
	return &info, nil
}

func newRecord(info domain.UserInfo, identity domain.UserIdentity) userRecord {
	return userRecord{
		ID:       mid.UUID(info.ID),
		Data:     info,
		Identity: identity,
		Created:  time.Now(),
		Updated:  time.Now(),
	}
}

func (c *UsersCollection) FindByLogin(ctx context.Context, login string) (*domain.UserInfo, error) {
	sr := c.storage.FindOne(
		ctx,
		bson.M{
			"identity.login": login,
			"deleted":        nil,
		},
		options.FindOne(),
	)
	var record userRecord
	err := sr.Decode(&record)
	if err != nil {
		return nil, c.mapError(err)
	}
	user := mapDbToUserInfo(record)
	return &user, nil
}

func (c *UsersCollection) Lookup(ctx context.Context, ID uuid.UUID) (*domain.UserInfo, error) {
	sr := c.storage.FindOne(
		ctx,
		bson.M{"_id": mid.UUID(ID)},
		options.FindOne(),
	)
	var record userRecord
	err := sr.Decode(&record)
	if err != nil {
		return nil, c.mapError(err)
	}
	user := mapDbToUserInfo(record)
	return &user, nil
}

func mapDbToUserInfo(record userRecord) domain.UserInfo {
	return record.Data
}

func (c *UsersCollection) Update(ctx context.Context, info domain.UserInfo) (*domain.UserInfo, error) {
	ur, err := c.storage.UpdateOne(
		ctx,
		bson.M{"_id": info.ID},
		bson.M{"$set": bson.M{
			"data": info,
		}},
		options.Update(),
	)
	if err != nil {
		return nil, c.mapError(err)
	}
	if ur.MatchedCount == 0 {
		return nil, c.mapError(mongo.ErrNoDocuments)
	}
	return &info, nil
}

func (c *UsersCollection) Delete(ctx context.Context, ID uuid.UUID) error {
	ur, err := c.storage.UpdateOne(
		ctx,
		bson.M{"_id": ID},
		bson.M{"$set": bson.M{
			"deleted": time.Now(),
		}},
		options.Update(),
	)
	if err != nil {
		return c.mapError(err)
	}
	if ur.MatchedCount == 0 {
		return c.mapError(mongo.ErrNoDocuments)
	}
	return nil
}

func (c *UsersCollection) Find(ctx context.Context, searchPattern string) ([]domain.UserInfo, error) {
	cur, err := c.storage.Find(
		ctx,
		bson.M{
			"$or": bson.A{
				bson.M{"data.name": bson.M{"$regex": searchPattern}},
				bson.M{"data.surname": bson.M{"$regex": searchPattern}},
				bson.M{"identity.login": bson.M{"$regex": searchPattern}},
				bson.M{"data.e_mail": bson.M{"$regex": searchPattern}},
			},
		},
		options.Find().SetShowRecordID(true),
	)
	if err != nil {
		return nil, c.mapError(err)
	}
	var data []userRecord
	if err := cur.All(ctx, &data); err != nil {
		return nil, c.mapError(err)
	}
	var result []domain.UserInfo
	for _, user := range data {
		result = append(result, mapDbToUserInfo(user))
	}
	return result, nil
}

func (s *Storage) NewUsersCollection(mapError func(error) error) *UsersCollection {
	users := s.mongo.Users()
	return &UsersCollection{
		storage:   users.Collection,
		logger:    users.Logger,
		mapErrorF: mapError,
	}
}
