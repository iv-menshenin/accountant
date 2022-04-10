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
	TargetCollection struct {
		storage  *mongo.Collection
		mapError func(error) error
	}
	targetRecord struct {
		ID       mid.UUID     `bson:"_id"`
		Data     model.Target `bson:"data"`
		Created  time.Time    `bson:"created"`
		Updated  time.Time    `bson:"updated"`
		Deleted  *time.Time   `bson:"deleted"`
		OwnerCtx mid.UUID     `bson:"ownerCtx"`
	}
)

func (t *TargetCollection) Create(ctx context.Context, target model.Target) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var record = mapTargetToRecord(ctx, target)
		_, err := t.storage.InsertOne(ctx, record, options.InsertOne())
		return t.mapError(err)
	}
}

func mapTargetToRecord(ctx context.Context, target model.Target) targetRecord {
	return targetRecord{
		ID:       mid.UUID(target.TargetID),
		Data:     target,
		Created:  time.Now(),
		Updated:  time.Now(),
		OwnerCtx: mid.UUID(getOwnerCtx(ctx)),
	}
}

func (t *TargetCollection) Lookup(ctx context.Context, targetID uuid.UUID) (*model.Target, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = targetIdFilter(targetID)
		result := t.storage.FindOne(ctx, filter, options.FindOne().SetShowRecordID(true))
		if err := result.Err(); err != nil {
			return nil, t.mapError(err)
		}
		var tar targetRecord
		if err := result.Decode(&tar); err != nil {
			return nil, t.mapError(err)
		}
		return mapRecordToTarget(tar), nil
	}
}

func targetIdFilter(id uuid.UUID) interface{} {
	return bson.M{"_id": bson.M{"$eq": mid.UUID(id)}}
}

func mapRecordToTarget(rec targetRecord) *model.Target {
	return &rec.Data
}

func (t *TargetCollection) Delete(ctx context.Context, targetID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var filter = targetIdFilter(targetID)
		_, err := t.storage.UpdateOne(ctx, filter, bson.M{"$set": bson.D{{Key: "deleted", Value: time.Now()}}}, options.Update())
		return t.mapError(err)
	}
}

func (t *TargetCollection) FindByPeriod(ctx context.Context, period model.Period) (targets []model.Target, eut error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = targetPeriodFilter(period)
		cur, err := t.storage.Find(ctx, filter, options.Find().SetShowRecordID(true))
		if err != nil {
			return nil, t.mapError(err)
		}
		defer func() {
			if e := cur.Close(ctx); e != nil && eut == nil {
				eut = e
			}
		}()
		for cur.Next(ctx) {
			var record targetRecord
			if err = cur.Decode(&record); err != nil {
				return nil, err
			}
			targets = append(targets, *mapRecordToTarget(record))
		}
		return targets, t.mapError(cur.Err())
	}
}

func targetPeriodFilter(period model.Period) interface{} {
	var filter = bson.D{
		bson.E{Key: "deleted", Value: nil},
		bson.E{Key: "data.period", Value: period},
	}
	return filter
}

func (s *Storage) NewTargetCollection(mapError func(error) error) *TargetCollection {
	return &TargetCollection{
		storage:  s.mongo.Targets(),
		mapError: mapError,
	}
}
