package mongodb

import (
	"context"
	"github.com/iv-menshenin/accountant/utils/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/iv-menshenin/accountant/model/domain"
	mid "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

/*
GetUserInfo(ctx context.Context, user, pass string) (UserInfo, error)
GetUserInfoByRefresh(uuid.UUID) (UserInfo, error)
*/

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

func (c *UsersCollection) NewLogin(ctx context.Context, info domain.UserInfo, identity domain.UserIdentity) (domain.UserInfo, error) {
	_, err := c.storage.InsertOne(
		ctx,
		newRecord(info, identity),
		options.InsertOne(),
	)
	if err != nil {
		return domain.UserInfo{}, c.mapError(err)
	}
	return info, nil
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

func (c *UsersCollection) FindByLogin(ctx context.Context, login string) (domain.UserInfo, error) {
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
		return domain.UserInfo{}, c.mapError(err)
	}
	user := mapDbToUserInfo(record)
	return user, nil
}

func mapDbToUserInfo(record userRecord) domain.UserInfo {
	return domain.UserInfo{
		ID:          uuid.UUID(record.ID),
		Name:        record.Data.Name,
		Surname:     record.Data.Surname,
		Permissions: record.Data.Permissions,
	}
}

func (c *UsersCollection) UserUpdate(ctx context.Context, info domain.UserInfo) (domain.UserInfo, error) {
	ur, err := c.storage.UpdateOne(
		ctx,
		bson.M{"_id": info.ID},
		bson.M{"$set": bson.M{
			"data": info,
		}},
		options.Update(),
	)
	if err != nil {
		return info, c.mapError(err)
	}
	if ur.MatchedCount == 0 {
		return info, c.mapError(mongo.ErrNoDocuments)
	}
	return info, nil
}

func (c *UsersCollection) UserDelete(ctx context.Context, ID uuid.UUID) error {
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

func (s *Storage) NewUsersCollection(mapError func(error) error) *UsersCollection {
	users := s.mongo.Users()
	return &UsersCollection{
		storage:   users.Collection,
		logger:    users.Logger,
		mapErrorF: mapError,
	}
}
