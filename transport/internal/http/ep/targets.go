package ep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/model"
)

type (
	TargetGetter interface {
		TargetGet(context.Context, model.GetTargetQuery) (*model.Target, error)
	}
	TargetCreator interface {
		TargetCreate(context.Context, model.PostTargetQuery) (*model.Target, error)
	}
	TargetDeleter interface {
		TargetDelete(context.Context, model.DeleteTargetQuery) error
	}
	TargetFinder interface {
		TargetsFind(context.Context, model.FindTargetQuery) ([]model.Target, error)
	}
	TargetProcessor interface {
		TargetCreator
		TargetGetter
		TargetDeleter
		TargetFinder
	}
	Targets struct {
		processor TargetProcessor
	}
)

func NewTargetsEP(ap TargetProcessor) *Targets {
	return &Targets{
		processor: ap,
	}
}

const (
	targetID = "target_id"
)

func (t *Targets) SetupRouting(router *mux.Router) {
	const targetsPath = "/targets"
	targetsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}", targetsPath, targetID)

	router.Path(targetsWithIDPath).Methods(http.MethodGet).Handler(t.LookupHandler())
	router.Path(targetsPath).Methods(http.MethodPost).Handler(t.PostHandler())
	router.Path(targetsWithIDPath).Methods(http.MethodDelete).Handler(t.DeleteHandler())
	router.Path(targetsPath).Methods(http.MethodGet).Handler(t.FindHandler())

}

func (t *Targets) LookupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getTargetMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := t.processor.TargetGet(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func getTargetMapper(r *http.Request) (q model.GetTargetQuery, err error) {
	id := mux.Vars(r)[targetID]
	if id == "" {
		err = errors.New(targetID + " must not be empty")
		return
	}
	err = q.TargetID.FromString(id)
	return
}

func (t *Targets) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postTargetMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := t.processor.TargetCreate(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func postTargetMapper(r *http.Request) (q model.PostTargetQuery, err error) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q.Target)
	return
}

func (t *Targets) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := deleteTargetMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		err = t.processor.TargetDelete(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeNoContent(w)
	}
}

func deleteTargetMapper(r *http.Request) (q model.DeleteTargetQuery, err error) {
	id := mux.Vars(r)[targetID]
	if id == "" {
		err = errors.New(targetID + " must not be empty")
		return
	}
	err = q.TargetID.FromString(id)
	return
}

func (t *Targets) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := findTargetMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		accounts, err := t.processor.TargetsFind(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, accounts)
	}
}

const (
	showClosed  = "show_closed"
	periodYear  = "year"
	periodMonth = "month"
)

func findTargetMapper(r *http.Request) (q model.FindTargetQuery, err error) {
	params := queryParams{r: r}
	if sc, ok := params.vars(showClosed); ok && sc != "false" && sc != "0" {
		q.ShowClosed = true
	}
	var period model.Period
	if y, ok := params.vars(periodYear); ok {
		period.Year, err = strconv.Atoi(y)
		if err != nil {
			err = fmt.Errorf("%s must be integer: %w", periodYear, err)
			return
		}
	}
	if y, ok := params.vars(periodMonth); ok {
		period.Month, err = strconv.Atoi(y)
		if err != nil {
			err = fmt.Errorf("%s must be integer: %w", periodMonth, err)
			return
		}
	}
	if period.Month != 0 || period.Year != 0 {
		q.Period = &period
	}
	return
}
