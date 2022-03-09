package main

import (
	"context"
	"flag"
	"github.com/iv-menshenin/accountant/storage/mongodb"
	"log"
	"os"
	"time"

	"github.com/iv-menshenin/accountant/auth"
	"github.com/iv-menshenin/accountant/business"
	"github.com/iv-menshenin/accountant/config"
	"github.com/iv-menshenin/accountant/logger"
	"github.com/iv-menshenin/accountant/storage"
	"github.com/iv-menshenin/accountant/transport"
	"github.com/iv-menshenin/appctl"
)

func main() {
	flag.Parse()
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
	authCore, err := auth.New("")
	if err != nil {
		return err
	}
	var logWriter = log.Default()

	mongoStorage, err := mongodb.NewStorage(config.New("db"), logWriter)
	if err != nil {
		return err
	}
	var (
		listeningError = make(chan error)
		appLogger      = logger.NewFromLogger(logWriter, logger.LogLevelDebug)

		accountCollection = mongoStorage.NewAccountCollection(storage.MapMongodbErrors)
		personsCollection = mongoStorage.NewPersonCollection(accountCollection, storage.MapMongodbErrors)
		objectsCollection = mongoStorage.NewObjectCollection(accountCollection, storage.MapMongodbErrors)

		appHnd         = business.New(appLogger, accountCollection, personsCollection, objectsCollection)
		queryTransport = transport.NewHTTPServer(config.New("http"), logWriter, appHnd, authCore)
	)

	go queryTransport.ListenAndServe(listeningError)
	select {
	case err = <-listeningError:
		return err
	case <-halt:
		return queryTransport.Shutdown(ctx)
	}
}
