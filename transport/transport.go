package transport

import (
	"context"
	"log"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/transport/internal/http"
)

type (
	Transport interface {
		ListenAndServe(chan<- error)
		Shutdown(context.Context) error
	}
)

func NewHTTPServer(config model.Config, logger *log.Logger, rp http.RequestProcessor, auth http.AuthCore) Transport {
	return http.New(config, logger, rp, auth)
}
