package transport

import (
	"context"
	"log"

	"github.com/iv-menshenin/accountant/transport/internal/http"
	"github.com/iv-menshenin/accountant/utils"
)

type (
	Transport interface {
		ListenAndServe(chan<- error)
		Shutdown(context.Context) error
	}
)

func NewHTTPServer(config utils.Config, logger *log.Logger, rp http.RequestProcessor, auth http.AuthCore) Transport {
	return http.New(config, logger, rp, auth)
}
