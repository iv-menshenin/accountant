package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/iv-menshenin/accountant/model/generic"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	JWTCore struct {
		private []byte
	}
	Claims struct {
		UserID     uuid.UUID `json:"user_id"`
		Context    []string  `json:"context"`
		Created    int64     `json:"created"`
		Expiration int64     `json:"expiration"`
	}
	Info struct {
		Token  *jwt.Token
		Claims Claims
	}
)

func (c Claims) Valid() error {
	if c.Expiration < time.Now().UTC().Unix() {
		return errors.New("token expired")
	}
	if len(c.Context) == 0 {
		return errors.New("you are not allowed")
	}
	return nil
}

func New(private string) (*JWTCore, error) {
	var err error
	var j JWTCore
	if private != "" {
		j.private, err = base64.URLEncoding.DecodeString(private)
	} else {
		j.private = make([]byte, 0, 32)
		_, err = rand.Read(j.private)
	}
	return &j, err
}

func (c *JWTCore) SignJWT(user generic.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.UUID.String(),
		"context":    user.Context,
		"created":    time.Now().UTC().Unix(),
		"expiration": time.Now().Add(time.Minute * 60).UTC().Unix(),
	})
	return token.SignedString(c.private)
}

func (c *JWTCore) RefreshToken(user generic.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.UUID.String(),
		"context":    user.Context,
		"created":    time.Now().UTC().Unix(),
		"expiration": time.Now().Add(time.Hour * 720).UTC().Unix(),
	})
	return token.SignedString(privateForRefresh(c.private))
}

func (c *JWTCore) ParseJWT(tokenString string) (*Info, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(_ *jwt.Token) (interface{}, error) {
		return c.private, nil
	})
	if err != nil {
		return nil, err
	}
	return &Info{Token: token, Claims: claims}, nil
}

func (c *JWTCore) ParseRefreshToken(tokenString string) (*Info, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(_ *jwt.Token) (interface{}, error) {
		return privateForRefresh(c.private), nil
	})
	if err != nil {
		return nil, err
	}
	return &Info{Token: token, Claims: claims}, nil
}

func privateForRefresh(private []byte) []byte {
	var inverted = make([]byte, len(private))
	for i, b := range private {
		inverted[len(inverted)-i-1] = b
	}
	return inverted
}
