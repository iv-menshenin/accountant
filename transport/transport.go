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

func NewHTTPServer(logger *log.Logger, rp http.RequestProcessor) Transport {
	return http.New(logger, rp)
}
