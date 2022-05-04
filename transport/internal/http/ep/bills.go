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
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	BillProcessor interface {
		request.BillCreator
		request.BillSaver
		request.BillGetter
		request.BillDeleter
		request.BillFinder
	}
	Bills struct {
		processor BillProcessor
	}
)

func NewBillsEP(ap BillProcessor) *Bills {
	return &Bills{
		processor: ap,
	}
}

const (
	pathSegmentBills    = "/bills"
	parameterNameBillID = "bill_id"
	parameterNamePeriod = "period"
)

func (b *Bills) SetupRouting(router *mux.Router) {
	billsWithAccountPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s", pathSegmentAccounts, parameterNameAccountID, pathSegmentBills)
	billsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}", pathSegmentBills, parameterNameBillID)

	router.Path(billsWithAccountPath).Methods(http.MethodPost).Handler(b.PostHandler())
	router.Path(billsWithIDPath).Methods(http.MethodGet).Handler(b.LookupHandler())
	router.Path(billsWithIDPath).Methods(http.MethodPut).Handler(b.PutHandler())
	router.Path(billsWithIDPath).Methods(http.MethodDelete).Handler(b.DeleteHandler())
	router.Path(pathSegmentBills).Methods(http.MethodGet).Handler(b.FindHandler())

}

func (b *Bills) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postBillsMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		bill, err := b.processor.BillCreate(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, bill)
	}
}

func postBillsMapper(r *http.Request) (q request.PostBillQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q.Data)
	return
}

func (b *Bills) LookupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getBillMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		bill, err := b.processor.BillGet(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, bill)
	}
}

func getBillMapper(r *http.Request) (q request.GetBillQuery, err error) {
	id := mux.Vars(r)[parameterNameBillID]
	if id == "" {
		err = errors.New(parameterNameBillID + " must not be empty")
		return
	}
	err = q.BillID.FromString(id)
	return
}

func (b *Bills) PutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := putBillMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		bill, err := b.processor.BillSave(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, bill)
	}
}

func putBillMapper(r *http.Request) (q request.PutBillQuery, err error) {
	id := mux.Vars(r)[parameterNameBillID]
	if id == "" {
		err = errors.New(parameterNameBillID + " must not be empty")
		return
	}
	if err = q.BillID.FromString(id); err != nil {
		return q, err
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q.Data)
	return
}

func (b *Bills) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := deleteBillMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		err = b.processor.BillDelete(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeNoContent(w)
	}
}

func deleteBillMapper(r *http.Request) (q request.DeleteBillQuery, err error) {
	id := mux.Vars(r)[parameterNameBillID]
	if id == "" {
		err = errors.New(parameterNameBillID + " must not be empty")
		return
	}
	err = q.BillID.FromString(id)
	return
}

func (b *Bills) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := findBillMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		bills, err := b.processor.BillsFind(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, bills)
	}
}

func findBillMapper(r *http.Request) (q request.FindBillsQuery, err error) {
	params := queryParams{r: r}
	if target, ok := params.vars(parameterNameTargetID); ok {
		var targetID uuid.UUID
		if err = targetID.FromString(target); err != nil {
			return
		}
		q.TargetID = &targetID
	}
	if account, ok := params.vars(parameterNameAccountID); ok {
		var accountID uuid.UUID
		if err = accountID.FromString(account); err != nil {
			return
		}
		q.AccountID = &accountID
	}
	if period, ok := params.vars(parameterNamePeriod); ok {
		split := strings.Split(period, ".")
		if len(split) != 2 {
			return q, fmt.Errorf("parameter %s unknown format, expected MM.YYYY", parameterNamePeriod)
		}
		var m, y int
		if m, err = strconv.Atoi(split[0]); err != nil {
			return
		}
		if y, err = strconv.Atoi(split[1]); err != nil {
			return
		}
		q.Period = &domain.Period{
			Month: m,
			Year:  y,
		}
	}
	return
}
