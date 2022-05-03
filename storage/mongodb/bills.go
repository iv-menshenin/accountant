package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mid "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	BillsCollection struct {
		storage  *mongo.Collection
		mapError func(error) error
	}
	billRecord struct {
		ID       mid.UUID    `bson:"_id"`
		Data     domain.Bill `bson:"data"`
		Created  time.Time   `bson:"created"`
		Updated  time.Time   `bson:"updated"`
		Deleted  *time.Time  `bson:"deleted"`
		OwnerCtx mid.UUID    `bson:"ownerCtx"`
	}
)

func (b *BillsCollection) Create(ctx context.Context, bill domain.Bill) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var record = mapBillToRecord(ctx, bill)
		_, err := b.storage.InsertOne(ctx, record, options.InsertOne())
		return b.mapError(err)
	}
}

func mapBillToRecord(ctx context.Context, bill domain.Bill) billRecord {
	return billRecord{
		ID:       mid.UUID(bill.BillID),
		Data:     bill,
		Created:  time.Now(),
		Updated:  time.Now(),
		OwnerCtx: mid.UUID(getOwnerCtx(ctx)),
	}
}

func (b *BillsCollection) Lookup(ctx context.Context, billID uuid.UUID) (*domain.Bill, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = billIdFilter(billID)
		result := b.storage.FindOne(ctx, filter, options.FindOne().SetShowRecordID(true))
		if err := result.Err(); err != nil {
			return nil, b.mapError(err)
		}
		var bill billRecord
		if err := result.Decode(&bill); err != nil {
			return nil, b.mapError(err)
		}
		return mapRecordToBill(bill), nil
	}
}

func (b *BillsCollection) Replace(ctx context.Context, bill domain.Bill) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var filter = billIdFilter(bill.BillID)
		var record = mapBillToRecord(ctx, bill)
		var document = updateBillDocument(record)
		_, err := b.storage.UpdateOne(ctx, filter, document, options.Update())
		return b.mapError(err)
	}
}

func updateBillDocument(record billRecord) interface{} {
	return bson.M{"$set": bson.D{
		{Key: "updated", Value: record.Updated},
		{Key: "data", Value: record.Data},
	}}
}

func billIdFilter(id uuid.UUID) interface{} {
	return bson.M{"_id": bson.M{"$eq": mid.UUID(id)}}
}

func mapRecordToBill(rec billRecord) *domain.Bill {
	return &rec.Data
}

func (b *BillsCollection) FindBy(ctx context.Context, accountID, targetID *uuid.UUID, period *domain.Period) (bills []domain.Bill, eut error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = billFilter(accountID, targetID, period)
		cur, err := b.storage.Find(ctx, filter, options.Find().SetShowRecordID(true))
		if err != nil {
			return nil, b.mapError(err)
		}
		defer func() {
			if e := cur.Close(ctx); e != nil && eut == nil {
				eut = e
			}
		}()
		for cur.Next(ctx) {
			var record billRecord
			if err = cur.Decode(&record); err != nil {
				return nil, err
			}
			bills = append(bills, *mapRecordToBill(record))
		}
		return bills, b.mapError(cur.Err())
	}
}

func billFilter(accountID *uuid.UUID, targetID *uuid.UUID, period *domain.Period) interface{} {
	var filter = bson.D{
		bson.E{Key: "deleted", Value: nil},
	}
	if accountID != nil {
		filter = append(filter, bson.E{Key: "data.account_id", Value: *accountID})
	}
	if targetID != nil {
		filter = append(filter, bson.E{Key: "data.target.target_id", Value: *targetID})
	}
	if period != nil {
		filter = append(filter, bson.E{Key: "data.period", Value: *period})
	}
	return filter
}

func (b *BillsCollection) Delete(ctx context.Context, billID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var filter = billIdFilter(billID)
		_, err := b.storage.UpdateOne(ctx, filter, bson.M{"$set": bson.D{{Key: "deleted", Value: time.Now()}}}, options.Update())
		return b.mapError(err)
	}
}

func (s *Storage) NewBillsCollection(mapError func(error) error) *BillsCollection {
	return &BillsCollection{
		storage:  s.mongo.Bills(),
		mapError: mapError,
	}
}
