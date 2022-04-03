package ep

import (
	"net/http"

	"github.com/gorilla/mux"
)

type queryParams struct {
	r *http.Request
}

func (q queryParams) vars(name string) (string, bool) {
	if s, ok := mux.Vars(q.r)[name]; ok {
		return s, ok
	}
	if s, ok := q.r.URL.Query()[name]; ok {
		if len(s) == 0 {
			return "", ok
		}
		return s[0], ok
	}
	return "", false
}
