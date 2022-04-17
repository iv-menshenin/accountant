package mongodb

import (
	"context"

	"github.com/iv-menshenin/accountant/utils/uuid"
)

type ownerCtx struct{}

func SetOwnerCtx(ctx context.Context, owner uuid.UUID) context.Context {
	return context.WithValue(ctx, ownerCtx{}, owner)
}

func getOwnerCtx(ctx context.Context) uuid.UUID {
	if v := ctx.Value(ownerCtx{}); v != nil {
		return v.(uuid.UUID)
	}
	return uuid.UUID{}
}
