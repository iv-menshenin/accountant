package memory

import (
	"strings"

	"github.com/iv-menshenin/accountant/model"
)

func checkPersonFullName(person model.Person, pattern string) bool {
	return strings.Contains(person.Name, pattern) ||
		strings.Contains(person.Surname, pattern) ||
		strings.Contains(person.PatName, pattern)
}
