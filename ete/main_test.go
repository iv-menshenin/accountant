package ete_test

import (
	"bytes"
	"encoding/json"
	"github.com/iv-menshenin/accountant/business"
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/store"
	"github.com/iv-menshenin/accountant/transport"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_ete(t *testing.T) {
	var (
		listeningError    = make(chan error)
		logger            = log.Default()
		accountCollection = store.NewAccountMemoryCollection()

		appHnd         = business.New(accountCollection)
		queryTransport = transport.NewHTTPServer(logger, appHnd)
	)
	go queryTransport.ListenAndServe(listeningError)

	var someDate = time.Date(2016, 12, 12, 0, 0, 0, 0, time.UTC)
	var buf = bytes.NewBufferString("")
	var account = model.AccountData{
		Comment:       "test comment",
		AgreementNum:  "NNN-0000001",
		AgreementDate: &someDate,
	}
	enc := json.NewEncoder(buf)
	if err := enc.Encode(account); err != nil {
		t.Error(err)
	}
	req, err := http.NewRequest(http.MethodPost, "/accounts", buf)
	if err != nil {
		t.Error(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		t.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	select {
	case err = <-listeningError:
		t.Error(err)
	}
}
