package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mid "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	PaymentCollection struct {
		storage  *mongo.Collection
		mapError func(error) error
	}
	paymentRecord struct {
		ID       mid.UUID      `bson:"_id"`
		Data     model.Payment `bson:"data"`
		Created  time.Time     `bson:"created"`
		Updated  time.Time     `bson:"updated"`
		Deleted  *time.Time    `bson:"deleted"`
		OwnerCtx mid.UUID      `bson:"ownerCtx"`
	}
)

func (p *PaymentCollection) Create(ctx context.Context, payment model.Payment) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var record = mapPaymentToRecord(ctx, payment)
		_, err := p.storage.InsertOne(ctx, record, options.InsertOne())
		return p.mapError(err)
	}
}

func mapPaymentToRecord(ctx context.Context, payment model.Payment) paymentRecord {
	return paymentRecord{
		ID:       mid.UUID(payment.PaymentID),
		Data:     payment,
		Created:  time.Now(),
		Updated:  time.Now(),
		OwnerCtx: mid.UUID(getOwnerCtx(ctx)),
	}
}

func (p *PaymentCollection) Delete(ctx context.Context, paymentID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var filter = paymentIdFilter(paymentID)
		_, err := p.storage.UpdateOne(ctx, filter, bson.M{"$set": bson.D{{Key: "deleted", Value: time.Now()}}}, options.Update())
		return p.mapError(err)
	}
}

func paymentIdFilter(id uuid.UUID) interface{} {
	return bson.M{"_id": bson.M{"$eq": mid.UUID(id)}}
}

func (p *PaymentCollection) FindByAccount(ctx context.Context, accountID uuid.UUID) (payments []model.Payment, eut error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = paymentFilter(&accountID, nil)
		cur, err := p.storage.Find(ctx, filter, options.Find().SetShowRecordID(true))
		if err != nil {
			return nil, p.mapError(err)
		}
		defer func() {
			if e := cur.Close(ctx); e != nil && eut == nil {
				eut = e
			}
		}()
		for cur.Next(ctx) {
			var record paymentRecord
			if err = cur.Decode(&record); err != nil {
				return nil, err
			}
			payments = append(payments, *mapRecordToPayment(record))
		}
		return payments, p.mapError(cur.Err())
	}
}

func (p *PaymentCollection) FindByIDs(ctx context.Context, uuids []uuid.UUID) (payments []model.Payment, eut error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = paymentFilter(nil, uuids)
		cur, err := p.storage.Find(ctx, filter, options.Find().SetShowRecordID(true))
		if err != nil {
			return nil, p.mapError(err)
		}
		defer func() {
			if e := cur.Close(ctx); e != nil && eut == nil {
				eut = e
			}
		}()
		for cur.Next(ctx) {
			var record paymentRecord
			if err = cur.Decode(&record); err != nil {
				return nil, err
			}
			payments = append(payments, *mapRecordToPayment(record))
		}
		return payments, p.mapError(cur.Err())
	}
}

func mapRecordToPayment(rec paymentRecord) *model.Payment {
	return &rec.Data
}

func paymentFilter(accountID *uuid.UUID, uuids []uuid.UUID) interface{} {
	var filter = bson.D{
		bson.E{Key: "deleted", Value: nil},
	}
	if accountID != nil {
		filter = append(filter, bson.E{Key: "data.account_id", Value: *accountID})
	}
	if len(uuids) > 0 {
		ids := make([]mid.UUID, len(uuids))
		for i := range uuids {
			ids[i] = mid.UUID(uuids[i])
		}
		filter = append(filter, bson.E{Key: "_id", Value: bson.M{"$in": ids}})
	}
	return filter
}

func (s *Storage) NewPaymentCollection(mapError func(error) error) *PaymentCollection {
	return &PaymentCollection{
		storage:  s.mongo.Payments(),
		mapError: mapError,
	}
}
