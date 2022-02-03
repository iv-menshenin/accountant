package http

import (
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/transport/internal/http/ep"
	"github.com/iv-menshenin/accountant/transport/internal/http/static"
	"net/http"
)

type (
	RequestProcessor interface {
		ep.AccountProcessor
		ep.PersonProcessor
		ep.ObjectProcessor
	}
)

func makeRouter(rp RequestProcessor, auth AuthCore) http.Handler {
	router := mux.NewRouter()
	wwwSubRouter := router.PathPrefix("/www").Subrouter()

	stat := static.New()
	stat.SetupRouting(wwwSubRouter)

	apiSubRouter := router.PathPrefix("/api").Subrouter()
	if auth != nil {
		apiSubRouter.Use(auth.Middleware())
	}

	accounts := ep.NewAccountsEP(rp)
	accounts.SetupRouting(apiSubRouter)

	persons := ep.NewPersonsEP(rp)
	persons.SetupRouting(apiSubRouter)

	objects := ep.NewObjectsEP(rp)
	objects.SetupRouting(apiSubRouter)

	authSubRouter := router.PathPrefix("/auth").Subrouter()

	authP := ep.NewAuthEP(auth)
	authP.SetupRouting(authSubRouter)

	return router
}
