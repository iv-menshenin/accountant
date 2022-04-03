package ep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iv-menshenin/accountant/model"
	"net/http"
	"strconv"
)

type (
	ObjectGetter interface {
		ObjectGet(context.Context, model.GetObjectQuery) (*model.Object, error)
	}
	ObjectCreator interface {
		ObjectCreate(context.Context, model.PostObjectQuery) (*model.Object, error)
	}
	ObjectSaver interface {
		ObjectSave(context.Context, model.PutObjectQuery) (*model.Object, error)
	}
	ObjectDeleter interface {
		ObjectDelete(context.Context, model.DeleteObjectQuery) error
	}
	ObjectFinder interface {
		ObjectsFind(context.Context, model.FindObjectsQuery) ([]model.Object, error)
	}
	ObjectProcessor interface {
		ObjectCreator
		ObjectGetter
		ObjectSaver
		ObjectDeleter
		ObjectFinder
	}
	Objects struct {
		processor ObjectProcessor
	}
)

func NewObjectsEP(pp ObjectProcessor) *Objects {
	return &Objects{
		processor: pp,
	}
}

const (
	objectID = "object_id"
)

func (o *Objects) SetupRouting(router *mux.Router) {
	const accountsPath = "/accounts"
	const objectsPath = "/objects"
	objectsWithoutIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s", accountsPath, accountID, objectsPath)
	objectsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s/{%s:[0-9a-f\\-]+}", accountsPath, accountID, objectsPath, objectID)

	router.Path(objectsWithIDPath).Methods(http.MethodGet).Handler(o.LookupHandler())
	router.Path(objectsWithoutIDPath).Methods(http.MethodPost).Handler(o.PostHandler())
	router.Path(objectsWithIDPath).Methods(http.MethodPut).Handler(o.PutHandler())
	router.Path(objectsWithIDPath).Methods(http.MethodDelete).Handler(o.DeleteHandler())
	router.Path(objectsPath).Methods(http.MethodGet).Handler(o.FindHandler())

}

func (p *Objects) LookupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := p.processor.ObjectGet(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func getObjectMapper(r *http.Request) (q model.GetObjectQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[objectID]
	if id == "" {
		err = errors.New(objectID + " must not be empty")
		return
	}
	err = q.ObjectID.FromString(id)
	return
}

func (p *Objects) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := p.processor.ObjectCreate(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func postObjectMapper(r *http.Request) (q model.PostObjectQuery, err error) {
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
	err = decoder.Decode(&q.ObjectData)
	return
}

func (p *Objects) PutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := putObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		account, err := p.processor.ObjectSave(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, account)
	}
}

func putObjectMapper(r *http.Request) (q model.PutObjectQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[objectID]
	if id == "" {
		err = errors.New(objectID + " must not be empty")
		return
	}
	err = q.ObjectID.FromString(id)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q.ObjectData)
	return
}

func (p *Objects) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := deleteObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		err = p.processor.ObjectDelete(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeNoContent(w)
	}
}

func deleteObjectMapper(r *http.Request) (q model.DeleteObjectQuery, err error) {
	id := mux.Vars(r)[accountID]
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[objectID]
	if id == "" {
		err = errors.New(objectID + " must not be empty")
		return
	}
	err = q.ObjectID.FromString(id)
	return
}

func (p *Objects) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := findObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		accounts, err := p.processor.ObjectsFind(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, accounts)
	}
}

func findObjectMapper(r *http.Request) (q model.FindObjectsQuery, err error) {
	params := queryParams{r: r}
	id, _ := params.vars(accountID)
	if id == "" {
		err = errors.New(accountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	if address, ok := params.vars(addressField); ok {
		q.Address = &address
	}
	if numStr, ok := params.vars(objectNumField); ok {
		var num int
		if num, err = strconv.Atoi(numStr); err != nil {
			return
		}
		q.Number = &num
	}
	return
}
