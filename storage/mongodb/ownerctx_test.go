package mongodb

import (
	"context"
	"reflect"
	"testing"

	"github.com/iv-menshenin/accountant/utils/uuid"
)

func TestSetOwnerCtx(t *testing.T) {
	var (
		ctx       = context.Background()
		needID    = uuid.NewUUID()
		changedID = uuid.NewUUID()
	)
	ctx = SetOwnerCtx(ctx, needID)
	owner := getOwnerCtx(ctx)
	if !reflect.DeepEqual(owner, needID) {
		t.Errorf("need: %s, got: %s", needID, owner)
	}
	ctx = SetOwnerCtx(ctx, changedID)
	owner = getOwnerCtx(ctx)
	if !reflect.DeepEqual(owner, changedID) {
		t.Errorf("need: %s, got: %s", changedID, owner)
	}
}
