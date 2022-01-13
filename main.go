package main

import (
	"context"
	"flag"
	"github.com/iv-menshenin/accountant/business"
	"log"
	"os"
	"time"

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
	var (
		logger         = log.Default()
		appHnd         = business.New()
		queryTransport = transport.New(logger, appHnd)
		listeningError = make(chan error)
	)
	queryTransport.ListenAndServe(listeningError)
	select {
	case err = <-listeningError:
		return err
	case <-halt:
		return queryTransport.Shutdown(ctx)
	}
}
