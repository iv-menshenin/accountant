package ep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/model/request"
)

type (
	TargetProcessor interface {
		request.TargetCreator
		request.TargetGetter
		request.TargetDeleter
		request.TargetFinder
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
	pathSegmentTargets      = "/targets"
	parameterNameTargetID   = "target_id"
	parameterNameType       = "type"
	parameterNameShowClosed = "show_closed"
	parameterNameYear       = "year"
	parameterNameMonth      = "month"
)

func (t *Targets) SetupRouting(router *mux.Router) {
	targetsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}", pathSegmentTargets, parameterNameTargetID)

	router.Path(targetsWithIDPath).Methods(http.MethodGet).Handler(t.LookupHandler())
	router.Path(pathSegmentTargets).Methods(http.MethodPost).Handler(t.PostHandler())
	router.Path(targetsWithIDPath).Methods(http.MethodDelete).Handler(t.DeleteHandler())
	router.Path(pathSegmentTargets).Methods(http.MethodGet).Handler(t.FindHandler())

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

func getTargetMapper(r *http.Request) (q request.GetTargetQuery, err error) {
	id := mux.Vars(r)[parameterNameTargetID]
	if id == "" {
		err = errors.New(parameterNameTargetID + " must not be empty")
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

func postTargetMapper(r *http.Request) (q request.PostTargetQuery, err error) {
	if rs, ok := r.URL.Query()[parameterNameType]; ok {
		q.Type = strings.Join(rs, ",")
	}
	if q.Type == "" {
		return q, fmt.Errorf("%s parameter must not be empty", parameterNameType)
	}
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

func deleteTargetMapper(r *http.Request) (q request.DeleteTargetQuery, err error) {
	id := mux.Vars(r)[parameterNameTargetID]
	if id == "" {
		err = errors.New(parameterNameTargetID + " must not be empty")
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

func findTargetMapper(r *http.Request) (q request.FindTargetsQuery, err error) {
	params := queryParams{r: r}
	if sc, ok := params.vars(parameterNameShowClosed); ok && sc != "false" && sc != "0" {
		q.ShowClosed = true
	}
	var period domain.Period
	if y, ok := params.vars(parameterNameYear); ok {
		period.Year, err = strconv.Atoi(y)
		if err != nil {
			err = fmt.Errorf("%s must be integer: %w", parameterNameYear, err)
			return
		}
	}
	if y, ok := params.vars(parameterNameMonth); ok {
		period.Month, err = strconv.Atoi(y)
		if err != nil {
			err = fmt.Errorf("%s must be integer: %w", parameterNameMonth, err)
			return
		}
	}
	if period.Month != 0 || period.Year != 0 {
		q.Period = &period
	}
	return
}
