package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/iv-menshenin/appctl"

	"github.com/iv-menshenin/accountant/auth"
	"github.com/iv-menshenin/accountant/business"
	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/logger"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/storage/mongodb"
	"github.com/iv-menshenin/accountant/transport"
)

func main() {
	var app = appctl.Application{
		MainFunc:              mainFunc,
		Resources:             nil,
		TerminationTimeout:    time.Second,
		InitializationTimeout: time.Millisecond * 500,
	}
	if err := app.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func mainFunc(ctx context.Context, halt <-chan struct{}) (err error) {
	var logWriter = log.Default()

	mongoStorage, err := mongodb.NewStorage(config.New("db"), logWriter)
	if err != nil {
		return err
	}
	var (
		listeningError = make(chan error)
		appLogger      = logger.NewFromLogger(logWriter, logger.LogLevelDebug)

		accountCollection  = mongoStorage.NewAccountsCollection(storage.MapMongodbErrors)
		personsCollection  = mongoStorage.NewPersonsCollection(accountCollection, storage.MapMongodbErrors)
		objectsCollection  = mongoStorage.NewObjectsCollection(accountCollection, storage.MapMongodbErrors)
		targetsCollection  = mongoStorage.NewTargetsCollection(storage.MapMongodbErrors)
		billsCollection    = mongoStorage.NewBillsCollection(storage.MapMongodbErrors)
		paymentsCollection = mongoStorage.NewPaymentsCollection(storage.MapMongodbErrors)
		usersCollection    = mongoStorage.NewUsersCollection(storage.MapMongodbErrors)

		appHnd = business.New(
			appLogger,
			accountCollection,
			personsCollection,
			objectsCollection,
			targetsCollection,
			billsCollection,
			paymentsCollection,
			usersCollection,
		)
	)

	authCore, err := auth.New(usersCollection, "")
	if err != nil {
		return err
	}
	queryTransport := transport.NewHTTPServer(config.New("http"), logWriter, appHnd, authCore)

	go queryTransport.ListenAndServe(listeningError)
	select {
	case err = <-listeningError:
		return err
	case <-halt:
		return queryTransport.Shutdown(ctx)
	}
}
