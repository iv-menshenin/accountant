package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/generic"
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

func (c *JWTCore) Auth(ctx context.Context, q generic.AuthQuery) (generic.AuthData, error) {
	userInfo, err := c.repository.FindByLogin(ctx, q.Login)
	if err != nil {
		return generic.AuthData{}, err
	}
	user := mapUserInfoToParticipant(*userInfo)

	token, err := c.SignJWT(user)
	if err != nil {
		return generic.AuthData{}, err
	}
	refresh, err := c.RefreshToken(user)
	if err != nil {
		return generic.AuthData{}, err
	}
	return generic.AuthData{
		JWT:     token,
		UserID:  user.UUID,
		Context: user.Context,
		Refresh: refresh,
	}, nil
}

func (c *JWTCore) Refresh(ctx context.Context, r generic.RefreshTokenQuery) (generic.AuthData, error) {
	info, err := c.ParseRefreshToken(r.Token)
	if err != nil {
		return generic.AuthData{}, err
	}
	if err := info.Claims.Valid(); err != nil {
		return generic.AuthData{}, err
	}

	userInfo, err := c.repository.Lookup(ctx, info.Claims.UserID)
	if err != nil {
		return generic.AuthData{}, err
	}
	user := mapUserInfoToParticipant(*userInfo)

	token, err := c.SignJWT(user)
	if err != nil {
		return generic.AuthData{}, err
	}
	refresh, err := c.RefreshToken(user)
	if err != nil {
		return generic.AuthData{}, err
	}
	return generic.AuthData{
		JWT:     token,
		UserID:  user.UUID,
		Context: user.Context,
		Refresh: refresh,
	}, nil
}

func mapUserInfoToParticipant(i domain.UserInfo) generic.User {
	var u = generic.User{
		UUID:     i.ID,
		UserName: i.Name,
		Context:  make([]string, len(i.Permissions)),
	}
	for _, perm := range i.Permissions {
		u.Context = append(u.Context, fmt.Sprintf("permission:%s", perm))
	}
	return u
}
