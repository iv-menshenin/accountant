package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	mid "go.mongodb.org/mongo-driver/x/mongo/driver/uuid"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	TargetsCollection Collection
	targetRecord      struct {
		ID       mid.UUID      `bson:"_id"`
		Data     domain.Target `bson:"data"`
		Created  time.Time     `bson:"created"`
		Updated  time.Time     `bson:"updated"`
		Deleted  *time.Time    `bson:"deleted"`
		OwnerCtx mid.UUID      `bson:"ownerCtx"`
	}
)

func (t *TargetsCollection) mapError(err error) error {
	return (*Collection)(t).mapError(err)
}

func (t *TargetsCollection) Create(ctx context.Context, target domain.Target) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var record = mapTargetToRecord(ctx, target)
		_, err := t.storage.InsertOne(ctx, record, options.InsertOne())
		return (*Collection)(t).mapError(err)
	}
}

func mapTargetToRecord(ctx context.Context, target domain.Target) targetRecord {
	return targetRecord{
		ID:       mid.UUID(target.TargetID),
		Data:     target,
		Created:  time.Now(),
		Updated:  time.Now(),
		OwnerCtx: mid.UUID(getOwnerCtx(ctx)),
	}
}

func (t *TargetsCollection) Update(ctx context.Context, targetID uuid.UUID, target domain.TargetData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		update := bson.M{"$set": bson.D{
			{Key: "data.period", Value: target.Period},
			{Key: "data.closed", Value: target.Closed},
			{Key: "data.cost", Value: target.Cost},
			{Key: "data.comment", Value: target.Comment},
			{Key: "updated", Value: time.Now()},
		}}
		_, err := t.storage.UpdateOne(ctx, bson.D{{Key: "_id", Value: mid.UUID(targetID)}}, update, options.Update())
		return t.mapError(err)
	}
}

func (t *TargetsCollection) Lookup(ctx context.Context, targetID uuid.UUID) (*domain.Target, error) {
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

func mapRecordToTarget(rec targetRecord) *domain.Target {
	return &rec.Data
}

func (t *TargetsCollection) Delete(ctx context.Context, targetID uuid.UUID) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		var filter = targetIdFilter(targetID)
		_, err := t.storage.UpdateOne(ctx, filter, bson.M{"$set": bson.D{{Key: "deleted", Value: time.Now()}}}, options.Update())
		return t.mapError(err)
	}
}

func (t *TargetsCollection) FindByPeriod(ctx context.Context, option storage.FindTargetOption) (targets []domain.Target, eut error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		var filter = targetPeriodFilter(option)
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

func targetPeriodFilter(option storage.FindTargetOption) interface{} {
	var filter = bson.D{
		bson.E{Key: "deleted", Value: nil},
	}
	if option.Year > 0 {
		filter = append(filter, bson.E{Key: "data.period.year", Value: option.Year})
	}
	if option.Month > 0 {
		filter = append(filter, bson.E{Key: "data.period.month", Value: option.Month})
	}
	if !option.ShowClosed {
		// filter = append(filter, bson.E{Key: "data.closed", Value: nil})
	}
	return filter
}

func (s *Storage) NewTargetsCollection(mapError func(error) error) *TargetsCollection {
	targets := s.mongo.Targets()
	return &TargetsCollection{
		storage:   targets.Collection,
		logger:    targets.Logger,
		mapErrorF: mapError,
	}
}
