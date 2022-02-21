package mongodb

import (
	"context"

	model "github.com/iv-menshenin/accountant/model/uuid"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

type ownerCtx struct{}

func SetOwnerCtx(ctx context.Context, owner model.UUID) context.Context {
	return context.WithValue(ctx, ownerCtx{}, uuid.UUID(owner))
}

func getOwnerCtx(ctx context.Context) uuid.UUID {
	if v := ctx.Value(ownerCtx{}); v != nil {
		return v.(uuid.UUID)
	}
	return uuid.UUID{}
}
