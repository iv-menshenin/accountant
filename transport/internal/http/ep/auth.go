package ep

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/model/generic"
)

type (
	AuthProcessor interface {
		Auth(context.Context, generic.AuthQuery) (generic.AuthData, error)
		Refresh(context.Context, generic.RefreshTokenQuery) (generic.AuthData, error)
	}
	Auth struct {
		processor AuthProcessor
	}
)

func NewAuthEP(ap AuthProcessor) *Auth {
	return &Auth{
		processor: ap,
	}
}

func (a *Auth) SetupRouting(router *mux.Router) {
	const loginPath = "/login"
	const refreshPath = "/refresh"

	router.Path(loginPath).Methods(http.MethodPost).Handler(a.LoginHandler())
	router.Path(refreshPath).Methods(http.MethodPost).Handler(a.RefreshHandler())

}
func (a *Auth) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var q generic.AuthQuery
		j := json.NewDecoder(r.Body)
		j.DisallowUnknownFields()
		if err := j.Decode(&q); err != nil {
			writeQueryError(w, err)
			return
		}
		info, err := a.processor.Auth(r.Context(), q)
		if err != nil {
			writeUnauthorizedError(w, err)
			return
		}
		writeData(w, info)
	}
}

func (a *Auth) RefreshHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var q generic.RefreshTokenQuery
		j := json.NewDecoder(r.Body)
		j.DisallowUnknownFields()
		if err := j.Decode(&q); err != nil {
			writeQueryError(w, err)
			return
		}
		info, err := a.processor.Refresh(r.Context(), q)
		if err != nil {
			writeUnauthorizedError(w, err)
			return
		}
		writeData(w, info)
	}
}
