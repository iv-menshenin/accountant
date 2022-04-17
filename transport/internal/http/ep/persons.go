package ep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/model/request"
)

type (
	PersonProcessor interface {
		request.PersonCreator
		request.PersonGetter
		request.PersonSaver
		request.PersonDeleter
		request.PersonFinder
	}
	Persons struct {
		processor PersonProcessor
	}
)

func NewPersonsEP(pp PersonProcessor) *Persons {
	return &Persons{
		processor: pp,
	}
}

const (
	pathSegmentPersons    = "/persons"
	parameterNamePersonID = "person_id"
	parameterNamePerson   = "person"
)

func (p *Persons) SetupRouting(router *mux.Router) {
	personsWithoutIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s", pathSegmentAccounts, parameterNameAccountID, pathSegmentPersons)
	personsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s/{%s:[0-9a-f\\-]+}", pathSegmentAccounts, parameterNameAccountID, pathSegmentPersons, parameterNamePersonID)

	router.Path(personsWithIDPath).Methods(http.MethodGet).Handler(p.LookupHandler())
	router.Path(personsWithoutIDPath).Methods(http.MethodPost).Handler(p.PostHandler())
	router.Path(personsWithIDPath).Methods(http.MethodPut).Handler(p.PutHandler())
	router.Path(personsWithIDPath).Methods(http.MethodDelete).Handler(p.DeleteHandler())
	router.Path(pathSegmentPersons).Methods(http.MethodGet).Handler(p.FindHandler())

}

func (p *Persons) LookupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getPersonMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := p.processor.PersonGet(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func getPersonMapper(r *http.Request) (q request.GetPersonQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[parameterNamePersonID]
	if id == "" {
		err = errors.New(parameterNamePersonID + " must not be empty")
		return
	}
	err = q.PersonID.FromString(id)
	return
}

func (p *Persons) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postPersonMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		person, err := p.processor.PersonCreate(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, person)
	}
}

func postPersonMapper(r *http.Request) (q request.PostPersonQuery, err error) {
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
	err = decoder.Decode(&q.PersonData)
	return
}

func (p *Persons) PutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := putPersonMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := p.processor.PersonSave(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func putPersonMapper(r *http.Request) (q request.PutPersonQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[parameterNamePersonID]
	if id == "" {
		err = errors.New(parameterNamePersonID + " must not be empty")
		return
	}
	err = q.PersonID.FromString(id)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q.PersonData)
	return
}

func (p *Persons) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := deletePersonMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		err = p.processor.PersonDelete(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeNoContent(w)
	}
}

func deletePersonMapper(r *http.Request) (q request.DeletePersonQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[parameterNamePersonID]
	if id == "" {
		err = errors.New(parameterNamePersonID + " must not be empty")
		return
	}
	err = q.PersonID.FromString(id)
	return
}

func (p *Persons) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := findPersonMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		accounts, err := p.processor.PersonsFind(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, accounts)
	}
}

func findPersonMapper(r *http.Request) (q request.FindPersonsQuery, err error) {
	params := queryParams{r: r}
	id, _ := params.vars(parameterNameAccountID)
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	if person, ok := params.vars(parameterNamePerson); ok {
		q.PersonFullName = &person
	}
	return
}
