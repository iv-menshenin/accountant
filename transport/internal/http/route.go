package http

import (
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/transport/internal/http/ep"
	"net/http"
)

type (
	RequestProcessor interface {
		ep.AccountCreator
		ep.AccountGetter
		ep.AccountSaver
	}
)

func makeRouter(rp RequestProcessor) http.Handler {
	router := mux.NewRouter()

	accounts := ep.NewAccountsEP(rp)
	accounts.RegisterRoute(router)

	return router
}
