package auth

import (
	"fmt"
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"testing"
)

func TestJWTCore_InputKey(t *testing.T) {
	jwtC, err := New("YmFzZTY0IGRhdGE=")
	if err != nil {
		t.Fatal(err)
	}
	var user = model.User{
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
	jwtC, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	var user = model.User{
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