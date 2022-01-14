package ep

import (
	"encoding/json"
	"github.com/iv-menshenin/accountant/model"
	"log"
	"net/http"
)

type (
	ResponseMeta struct {
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	}
	ErrorResponse struct {
		Meta ResponseMeta `json:"meta"`
	}
	DataResponse struct {
		Meta ResponseMeta `json:"meta"`
		Data interface{}  `json:"data,omitempty"`
	}
)

func writeQueryError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write(ErrorResponse{
		Meta: ResponseMeta{
			Status:  "QueryError",
			Message: e.Error(),
		},
	}.data())
	if err != nil {
		log.Println(err)
	}
}

func writeError(w http.ResponseWriter, e error) {
	switch e.(type) {
	case model.Unauthorized:
		writeUnauthorizedError(w, e)
	case model.Forbidden:
		writeDataAccessError(w, e)
	case model.NotFound:
		writeNotFoundError(w, e)
	default:
		writeInternalError(w, e)
	}
}

func writeInternalError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write(ErrorResponse{
		Meta: ResponseMeta{
			Status:  "InternalError",
			Message: e.Error(),
		},
	}.data())
	if err != nil {
		log.Println(err)
	}
}

func writeDataAccessError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusForbidden)
	_, err := w.Write(ErrorResponse{
		Meta: ResponseMeta{
			Status:  "Forbidden",
			Message: e.Error(),
		},
	}.data())
	if err != nil {
		log.Println(err)
	}
}

func writeUnauthorizedError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write(ErrorResponse{
		Meta: ResponseMeta{
			Status:  "Unauthorized",
			Message: e.Error(),
		},
	}.data())
	if err != nil {
		log.Println(err)
	}
}

func writeNotFoundError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write(ErrorResponse{
		Meta: ResponseMeta{
			Status:  "NotFound",
			Message: e.Error(),
		},
	}.data())
	if err != nil {
		log.Println(err)
	}
}

func writeData(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(DataResponse{
		Meta: ResponseMeta{
			Status: "Ok",
		},
		Data: data,
	}.data())
	if err != nil {
		log.Println(err)
	}
}

func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (r ErrorResponse) data() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}
	return b
}

func (r DataResponse) data() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}
	return b
}
