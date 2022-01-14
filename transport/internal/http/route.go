package http

import (
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/transport/internal/http/ep"
	"net/http"
)

type (
	RequestProcessor interface {
		ep.AccountGetter
		ep.AccountSaver
	}
)

func makeRouter(rp RequestProcessor) http.Handler {
	router := mux.NewRouter()

	accounts := ep.NewAccountsEP(rp)

	router.Path(accounts.LookupPathPattern()).
		Methods(http.MethodGet).
		Handler(accounts.LookupHandler())

	return router
}
