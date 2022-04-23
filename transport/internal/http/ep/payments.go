package ep

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/model/request"
	"github.com/iv-menshenin/accountant/utils/uuid"
	"net/http"
)

type (
	PaymentProcessor interface {
		request.PaymentCreator
		request.PaymentGetter
		request.PaymentDeleter
		request.PaymentFinder
	}
	Payments struct {
		processor PaymentProcessor
	}
)

func NewPaymentsEP(pp PaymentProcessor) *Payments {
	return &Payments{
		processor: pp,
	}
}

const (
	pathSegmentPayments    = "/payments"
	parameterNamePaymentID = "payment_id"
)

func (p *Payments) SetupRouting(router *mux.Router) {
	paymentsWithAccountPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s", pathSegmentAccounts, parameterNameAccountID, pathSegmentPayments)
	paymentsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}", pathSegmentPayments, parameterNamePaymentID)

	router.Path(paymentsWithAccountPath).Methods(http.MethodPost).Handler(p.PostHandler())
	router.Path(paymentsWithIDPath).Methods(http.MethodGet).Handler(p.LookupHandler())
	router.Path(paymentsWithIDPath).Methods(http.MethodDelete).Handler(p.DeleteHandler())
	router.Path(pathSegmentPayments).Methods(http.MethodGet).Handler(p.FindHandler())
}

func (p *Payments) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postPaymentsMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := p.processor.PaymentCreate(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func postPaymentsMapper(r *http.Request) (q request.PostPaymentQuery, err error) {
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

func (p *Payments) LookupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getPaymentMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := p.processor.PaymentGet(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func getPaymentMapper(r *http.Request) (q request.GetPaymentQuery, err error) {
	id := mux.Vars(r)[parameterNamePaymentID]
	if id == "" {
		err = errors.New(parameterNamePaymentID + " must not be empty")
		return
	}
	err = q.PaymentID.FromString(id)
	return
}

func (p *Payments) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := deletePaymentMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		err = p.processor.PaymentDelete(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeNoContent(w)
	}
}

func deletePaymentMapper(r *http.Request) (q request.DeletePaymentQuery, err error) {
	id := mux.Vars(r)[parameterNamePaymentID]
	if id == "" {
		err = errors.New(parameterNamePaymentID + " must not be empty")
		return
	}
	err = q.PaymentID.FromString(id)
	return
}

func (p *Payments) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := findPaymentsMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		accounts, err := p.processor.PaymentsFind(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, accounts)
	}
}

func findPaymentsMapper(r *http.Request) (q request.FindPaymentsQuery, err error) {
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
	if person, ok := params.vars(parameterNamePersonID); ok {
		var personID uuid.UUID
		if err = personID.FromString(person); err != nil {
			return
		}
		q.PersonID = &personID
	}
	if object, ok := params.vars(parameterNameObjectID); ok {
		var objectID uuid.UUID
		if err = objectID.FromString(object); err != nil {
			return
		}
		q.ObjectID = &objectID
	}
	return
}
