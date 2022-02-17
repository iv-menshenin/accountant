package transport

import (
	"context"
	"log"

	"github.com/iv-menshenin/accountant/configstorage"
	ht "github.com/iv-menshenin/accountant/transport/internal/http"
)

type (
	Transport interface {
		ListenAndServe(chan<- error)
		Shutdown(context.Context) error
	}
)

func NewHTTPServer(logger *log.Logger, rp ht.RequestProcessor, auth ht.AuthCore) Transport {
	return ht.New(configstorage.New(""), logger, rp, auth)
}
