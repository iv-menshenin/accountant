package transport

import (
	"context"
	"log"

	"github.com/iv-menshenin/accountant/transport/internal/http"
)

type (
	Transport interface {
		ListenAndServe(chan<- error)
		Shutdown(context.Context) error
	}
)

func New(logger *log.Logger) Transport {
	return http.New(logger)
}
