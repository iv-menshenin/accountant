package ep

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/iv-menshenin/accountant/model/request"
)

type (
	ObjectProcessor interface {
		request.ObjectCreator
		request.ObjectGetter
		request.ObjectSaver
		request.ObjectDeleter
		request.ObjectFinder
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
	pathSegmentObjects    = "/objects"
	parameterNameObjectID = "object_id"
	parameterNameNumber   = "number"
	parameterNameAddress  = "address"
)

func (o *Objects) SetupRouting(router *mux.Router) {
	objectsWithoutIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s", pathSegmentAccounts, parameterNameAccountID, pathSegmentObjects)
	objectsWithIDPath := fmt.Sprintf("%s/{%s:[0-9a-f\\-]+}%s/{%s:[0-9a-f\\-]+}", pathSegmentAccounts, parameterNameAccountID, pathSegmentObjects, parameterNameObjectID)

	router.Path(objectsWithIDPath).Methods(http.MethodGet).Handler(o.LookupHandler())
	router.Path(objectsWithoutIDPath).Methods(http.MethodPost).Handler(o.PostHandler())
	router.Path(objectsWithIDPath).Methods(http.MethodPut).Handler(o.PutHandler())
	router.Path(objectsWithIDPath).Methods(http.MethodDelete).Handler(o.DeleteHandler())
	router.Path(pathSegmentObjects).Methods(http.MethodGet).Handler(o.FindHandler())

}

func (o *Objects) LookupHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := getObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		object, err := o.processor.ObjectGet(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, object)
	}
}

func getObjectMapper(r *http.Request) (q request.GetObjectQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[parameterNameObjectID]
	if id == "" {
		err = errors.New(parameterNameObjectID + " must not be empty")
		return
	}
	err = q.ObjectID.FromString(id)
	return
}

func (o *Objects) PostHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := postObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		object, err := o.processor.ObjectCreate(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, object)
	}
}

func postObjectMapper(r *http.Request) (q request.PostObjectQuery, err error) {
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
	err = decoder.Decode(&q.ObjectData)
	return
}

func (o *Objects) PutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := putObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		object, err := o.processor.ObjectSave(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, object)
	}
}

func putObjectMapper(r *http.Request) (q request.PutObjectQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[parameterNameObjectID]
	if id == "" {
		err = errors.New(parameterNameObjectID + " must not be empty")
		return
	}
	err = q.ObjectID.FromString(id)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&q.ObjectData)
	return
}

func (o *Objects) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := deleteObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		err = o.processor.ObjectDelete(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeNoContent(w)
	}
}

func deleteObjectMapper(r *http.Request) (q request.DeleteObjectQuery, err error) {
	id := mux.Vars(r)[parameterNameAccountID]
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	id = mux.Vars(r)[parameterNameObjectID]
	if id == "" {
		err = errors.New(parameterNameObjectID + " must not be empty")
		return
	}
	err = q.ObjectID.FromString(id)
	return
}

func (o *Objects) FindHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, err := findObjectMapper(r)
		if err != nil {
			writeQueryError(w, err)
			return
		}
		objects, err := o.processor.ObjectsFind(r.Context(), q)
		if err != nil {
			writeError(w, err)
			return
		}
		writeData(w, objects)
	}
}

func findObjectMapper(r *http.Request) (q request.FindObjectsQuery, err error) {
	params := queryParams{r: r}
	id, _ := params.vars(parameterNameAccountID)
	if id == "" {
		err = errors.New(parameterNameAccountID + " must not be empty")
		return
	}
	if err = q.AccountID.FromString(id); err != nil {
		return
	}
	if address, ok := params.vars(parameterNameAddress); ok {
		q.Address = &address
	}
	if numStr, ok := params.vars(parameterNameNumber); ok {
		var num int
		if num, err = strconv.Atoi(numStr); err != nil {
			return
		}
		q.Number = &num
	}
	return
}
