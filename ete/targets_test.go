package ete_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"

	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

func testTargets(t *testing.T, logData fmt.Stringer, actor httpActor) {
	var uuid uuid.UUID
	want := model.TargetData{
		Period: model.Period{
			Month: 12,
			Year:  2009,
		},
		Cost:    1200000000,
		Comment: "проверка куку",
	}
	t.Run("create", func(t *testing.T) {
		got, err := actor.createTarget("testTarget", want)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("can not create target: %s", err)
		}
		if !reflect.DeepEqual(want, got.TargetData) {
			t.Log(logData.String())
			t.Fatalf("matching error. want: %v, got: %v", want, got.TargetData)
		}
		uuid = got.TargetID
		found, err := actor.findTarget(true, 12, 2009)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("unexpected error: %s", err)
		}
		if len(found) != 1 {
			t.Fatalf("unexpected error, want nil, got: %v", got)
		}
		if !reflect.DeepEqual(want, found[0].TargetData) {
			t.Log(logData.String())
			t.Fatalf("matching error. want: %v, got: %v", want, found[0].TargetData)
		}
	})
	t.Run("lookup", func(t *testing.T) {
		got, err := actor.getTarget(uuid)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("can not lookup target: %s", err)
		}
		if !reflect.DeepEqual(want, got.TargetData) {
			t.Log(logData.String())
			t.Fatalf("matching error. want: %v, got: %v", want, got.TargetData)
		}
	})
	t.Run("delete", func(t *testing.T) {
		err := actor.deleteTarget(uuid)
		if err != nil {
			t.Log(logData.String())
			t.Fatalf("can not lookup target: %s", err)
		}
		got, err := actor.findTarget(true, 12, 2009)
		if err != nil && err != ErrNotFound {
			t.Log(logData.String())
			t.Fatalf("unexpected error: %s", err)
		}
		if got != nil {
			t.Fatalf("unexpected error, want nil, got: %v", got)
		}
	})
}

func (a *httpActor) createTarget(targetType string, data model.TargetData) (result *model.Target, err error) {
	var buf = bytes.NewBufferString("")
	enc := json.NewEncoder(buf)
	if err = enc.Encode(data); err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, a.targetsURL, buf)
	if err != nil {
		return
	}
	req.URL.RawQuery = url.Values{"type": []string{targetType}}.Encode()
	var respData TargetDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) getTarget(uuid uuid.UUID) (result *model.Target, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.targetsURL+"/"+uuid.String(), nil)
	if err != nil {
		return
	}
	var respData TargetDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) findTarget(showClosed bool, periodMonth, periodYear int) (result []model.Target, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, a.targetsURL, nil)
	if err != nil {
		return
	}
	req.URL.RawQuery = url.Values{
		"show_closed": []string{strconv.FormatBool(showClosed)},
		"year":        []string{strconv.Itoa(periodYear)},
		"month":       []string{strconv.Itoa(periodMonth)},
	}.Encode()
	var respData TargetsDataResponse
	if err = a.exec(req, &respData); err != nil {
		return nil, err
	}
	result = respData.Data
	return
}

func (a *httpActor) deleteTarget(uuid uuid.UUID) error {
	req, err := http.NewRequest(http.MethodDelete, a.targetsURL+"/"+uuid.String(), nil)
	if err != nil {
		return err
	}
	var respData TargetDataResponse
	if err = a.exec(req, &respData); err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}
	return nil
}
