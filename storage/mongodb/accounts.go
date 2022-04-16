package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mid "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
)

type (
	AccountsCollection struct {
		storage  *mongo.Collection
		mapError func(error) error
	}
	accountRecord struct {
		ID       mid.UUID          `bson:"_id"`
		Data     model.AccountData `bson:"data"`
		Persons  []model.Person    `bson:"persons"`
		Objects  []model.Object    `bson:"objects"`
		Created  time.Time         `bson:"created"`
		Updated  time.Time         `bson:"updated"`
		OwnerCtx mid.UUID          `bson:"ownerCtx"`
	}
)

func (a *AccountsCollection) Create(ctx context.Context, account model.Account) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var record = mapAccountToRecord(ctx, account)
		_, err := a.storage.InsertOne(ctx, record, options.InsertOne())
		return a.mapError(err)
	}
}

func (a *AccountsCollection) Lookup(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = accountIdFilter(id)
		result := a.storage.FindOne(ctx, filter, options.FindOne().SetShowRecordID(true))
		if err := result.Err(); err != nil {
			return nil, a.mapError(err)
		}
		var acc accountRecord
		if err := result.Decode(&acc); err != nil {
			return nil, a.mapError(err)
		}
		return mapRecordToAccount(acc), nil
	}
}

func (a *AccountsCollection) Replace(ctx context.Context, id uuid.UUID, account model.Account) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var filter = accountIdFilter(id)
		var record = mapAccountToRecord(ctx, account)
		var document = updateDocument(record)
		_, err := a.storage.UpdateOne(ctx, filter, document, options.Update())
		return a.mapError(err)
	}
}

func mapAccountToRecord(ctx context.Context, account model.Account) accountRecord {
	return accountRecord{
		ID:       mid.UUID(account.AccountID),
		Data:     account.AccountData,
		Persons:  account.Persons,
		Objects:  account.Objects,
		Created:  time.Now(),
		Updated:  time.Now(),
		OwnerCtx: mid.UUID(getOwnerCtx(ctx)),
	}
}

func updateDocument(record accountRecord) interface{} {
	return bson.M{"$set": bson.D{
		{"updated", record.Updated},
		{"data", record.Data},
		{"persons", record.Persons},
		{"objects", record.Objects},
	}}
}

func (a *AccountsCollection) Delete(ctx context.Context, id uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var filter = accountIdFilter(id)
		_, err := a.storage.DeleteOne(ctx, filter, options.Delete())
		return a.mapError(err)
	}
}

func (a *AccountsCollection) Find(ctx context.Context, option model.FindAccountOption) (accounts []model.Account, eut error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = accountFilter(option)
		cur, err := a.storage.Find(ctx, filter, options.Find().SetShowRecordID(true))
		if err != nil {
			return nil, a.mapError(err)
		}
		defer func() {
			if e := cur.Close(ctx); e != nil && eut == nil {
				eut = e
			}
		}()
		for cur.Next(ctx) {
			var record accountRecord
			if err = cur.Decode(&record); err != nil {
				return nil, err
			}
			accounts = append(accounts, *mapRecordToAccount(record))
		}
		return accounts, a.mapError(cur.Err())
	}
}

func accountIdFilter(id uuid.UUID) interface{} {
	return bson.M{"_id": bson.M{"$eq": mid.UUID(id)}}
}

func accountFilter(options model.FindAccountOption) interface{} {
	var filter = bson.D{}
	if options.Account != nil {
		filter = append(filter, bson.E{Key: "data.account", Value: *options.Account})
	}
	if options.Address != nil {
		filter = append(filter, bson.E{Key: "objects.street", Value: *options.Address})
	}
	if options.Number != nil {
		filter = append(filter, bson.E{Key: "objects.number", Value: *options.Number})
	}
	return filter
}

func mapRecordToAccount(record accountRecord) *model.Account {
	return &model.Account{
		AccountID:   uuid.UUID(record.ID),
		Persons:     record.Persons,
		Objects:     record.Objects,
		AccountData: record.Data,
	}
}

func (s *Storage) NewAccountsCollection(mapError func(error) error) *AccountsCollection {
	return &AccountsCollection{
		storage:  s.mongo.Accounts(),
		mapError: mapError,
	}
}
