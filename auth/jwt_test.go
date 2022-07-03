package auth

import (
	"context"
	"fmt"
	"github.com/iv-menshenin/accountant/model/domain"
	"testing"

	"github.com/iv-menshenin/accountant/model/generic"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func TestJWTCore_InputKey(t *testing.T) {
	jwtC, err := New(mockRepo{}, "YmFzZTY0IGRhdGE=")
	if err != nil {
		t.Fatal(err)
	}
	var user = generic.User{
		UUID:     uuid.NewUUID(),
		UserName: "devalio",
		Context:  []string{"admin", "read", "write"},
	}
	JWT, err := jwtC.SignJWT(user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("got JWT token:", JWT)
	info, err := jwtC.ParseJWT(JWT)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Claims.UserID.Equal(user.UUID) {
		t.Fatalf("wrong user_id, want: %s, got: %s", user.UUID, info.Claims.UserID)
	}
}

func TestJWTCore_RandomKey(t *testing.T) {
	jwtC, err := New(mockRepo{}, "")
	if err != nil {
		t.Fatal(err)
	}
	var user = generic.User{
		UUID:     uuid.NewUUID(),
		UserName: "devalio",
		Context:  []string{"green plum", "read", "write"},
	}
	JWT, err := jwtC.SignJWT(user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("got JWT token:", JWT)
	info, err := jwtC.ParseJWT(JWT)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Claims.UserID.Equal(user.UUID) {
		t.Fatalf("wrong user_id, want: %s, got: %s", user.UUID, info.Claims.UserID)
	}
}

type mockRepo struct{}

func (mockRepo) FindByLogin(ctx context.Context, login string) (domain.UserInfo, error) {
	return domain.UserInfo{
		ID:          uuid.NewUUID(),
		Name:        "devalio",
		Permissions: []domain.Permission{"green plum", "read", "write"},
	}, nil
}

func (mockRepo) Lookup(ctx context.Context, ID uuid.UUID) (domain.UserInfo, error) {
	return domain.UserInfo{
		ID:          uuid.NewUUID(),
		Name:        "devalio",
		Permissions: []domain.Permission{"green plum", "read", "write"},
	}, nil
}
