package store

import (
	"github.com/iv-menshenin/accountant/model"
	"github.com/iv-menshenin/accountant/model/uuid"
	"strings"
)

type (
	FindPersonOption struct {
		AccountID      *uuid.UUID
		PersonFullName *string
	}
)

func (q FindPersonOption) FillFromQuery(query model.FindPersonsQuery) {
	if query.PersonFullName != nil {
		q.PersonFullName = query.PersonFullName
	}
	if query.AccountID != nil {
		q.AccountID = query.AccountID
	}
}

func checkPersonFullName(person model.Person, pattern string) bool {
	return strings.Contains(person.Name, pattern) ||
		strings.Contains(person.Surname, pattern) ||
		strings.Contains(person.PatName, pattern)
}
