package ep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/model"
)

type (
	PersonGetter interface {
		PersonGet(context.Context, model.GetPersonQuery) (*model.Person, error)
	}
	PersonCreator interface {
		PersonCreate(context.Context, model.PostPersonQuery) (*model.Person, error)
	}
	PersonSaver interface {
		PersonSave(context.Context, model.PutPersonQuery) (*model.Person, error)
	}
	PersonDeleter interface {
		PersonDelete(context.Context, model.DeletePersonQuery) error
	}
	PersonFinder interface {
		PersonsFind(context.Context, model.FindPersonsQuery) ([]model.Person, error)
	}
	PersonProcessor interface {
		PersonCreator
		PersonGetter
		PersonSaver
		PersonDeleter
		PersonFinder
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
	personID = "person_id"
)

func (p *Persons) SetupRouting(router *mux.Router) {
	const accountsPath = "/accounts"
	const personsPath = "/persons"
	personsWithoutIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s", accountsPath, accountID, personsPath)
	personsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s/{%s:[0-9a-f\\-]+}", accountsPath, accountID, personsPath, personID)

	router.Path(personsWithIDPath).Methods(http.MethodGet).Handler(p.LookupHandler())
	router.Path(personsWithoutIDPath).Methods(http.MethodPost).Handler(p.PostHandler())
	router.Path(personsWithIDPath).Methods(http.MethodPut).Handler(p.PutHandler())
	router.Path(personsWithIDPath).Methods(http.MethodDelete).Handler(p.DeleteHandler())
	router.Path(personsPath).Methods(http.MethodGet).Handler(p.FindHandler())

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

func getPersonMapper(r *http.Request) (q model.GetPersonQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[personID]
	if id == "" {
		err = errors.New(personID + " must not be empty")
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

func postPersonMapper(r *http.Request) (q model.PostPersonQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
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

func putPersonMapper(r *http.Request) (q model.PutPersonQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[personID]
	if id == "" {
		err = errors.New(personID + " must not be empty")
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

func deletePersonMapper(r *http.Request) (q model.DeletePersonQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[personID]
	if id == "" {
		err = errors.New(personID + " must not be empty")
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

const personNameField = "person"

func findPersonMapper(r *http.Request) (q model.FindPersonsQuery, err error) {
	params := queryParams{r: r}
	id, _ := params.vars(accountID)
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	if person, ok := params.vars(personNameField); ok {
		q.PersonFullName = &person
	}
	return
}
