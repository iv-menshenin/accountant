package http

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/transport/internal/http/ep"
	"github.com/iv-menshenin/accountant/transport/internal/http/static"
)

type (
	RequestProcessor interface {
		ep.AccountProcessor
		ep.PersonProcessor
		ep.ObjectProcessor
		ep.TargetProcessor
		ep.BillProcessor
	}
)

const (
	PathAuth = "/auth"
	PathAPI  = "/api"
	PathWWW  = "/www"
)

func makeRouter(rp RequestProcessor, auth AuthCore) http.Handler {
	router := mux.NewRouter()
	wwwSubRouter := router.PathPrefix(PathWWW).Subrouter()

	stat := static.New()
	stat.SetupRouting(wwwSubRouter)

	apiSubRouter := router.PathPrefix(PathAPI).Subrouter()
	if auth != nil {
		apiSubRouter.Use(auth.Middleware())
	}

	accounts := ep.NewAccountsEP(rp)
	accounts.SetupRouting(apiSubRouter)

	persons := ep.NewPersonsEP(rp)
	persons.SetupRouting(apiSubRouter)

	objects := ep.NewObjectsEP(rp)
	objects.SetupRouting(apiSubRouter)

	targets := ep.NewTargetsEP(rp)
	targets.SetupRouting(apiSubRouter)

	bills := ep.NewBillsEP(rp)
	bills.SetupRouting(apiSubRouter)

	authSubRouter := router.PathPrefix(PathAuth).Subrouter()

	authP := ep.NewAuthEP(auth)
	authP.SetupRouting(authSubRouter)

	return router
}
