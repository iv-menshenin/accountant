package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
)

type contextAuth struct{}

const httpHeader = "X-Auth-Token"

func (c *JWTCore) Middleware() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearer := r.Header.Get(httpHeader)
			if bearer == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			info, err := c.ParseJWT(bearer)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			for _, c := range info.Claims.Context {
				if strings.EqualFold(c, "ban") {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextAuth{}, info)))
		})
	}
}

func (c *JWTCore) Auth(context.Context, model.AuthQuery) (model.AuthData, error) {
	var user = model.User{
		UUID:     uuid.NilUUID(),
		UserName: "test",
		Context:  []string{"test", "admin"},
	}

	// dummy

	token, err := c.SignJWT(user)
	if err != nil {
		return model.AuthData{}, err
	}
	return model.AuthData{
		JWT:     token,
		UserID:  user.UUID,
		Context: user.Context,
	}, nil
}

func (c *JWTCore) Refresh(context.Context, model.RefreshTokenQuery) (model.AuthData, error) {
	var user = model.User{
		UUID:     uuid.NilUUID(),
		UserName: "test",
		Context:  []string{"test", "admin"},
	}

	// dummy

	token, err := c.SignJWT(user)
	if err != nil {
		return model.AuthData{}, err
	}
	return model.AuthData{
		JWT:     token,
		UserID:  user.UUID,
		Context: user.Context,
	}, nil
}
