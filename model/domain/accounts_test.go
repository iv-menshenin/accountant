package domain

import (
	"encoding/json"
	"fmt"

	"github.com/iv-menshenin/accountant/utils/uuid"
)

func ExampleNestedPerson() {
	var acc_id, pers_id uuid.UUID
	if err := acc_id.FromString("e7435d15-3d31-4aeb-91dd-08f787dab2d8"); err != nil {
		panic(err)
	}
	if err := pers_id.FromString("f2b8962f-51f1-4f02-9eb1-53311ab94cba"); err != nil {
		panic(err)
	}
	var np = NestedPerson{
		Person: Person{
			PersonID: pers_id,
			PersonData: PersonData{
				Name:     "Igor",
				Surname:  "Menshenin",
				PatName:  "Vladimirovitch",
				DOB:      nil,
				IsMember: true,
				Phone:    "(555) 55-55-55",
				EMail:    "devalio@yandex.ru",
			},
		},
		AccountID: acc_id,
	}
	data, err := json.Marshal(np)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	// Output: {"person_id":"f2b8962f-51f1-4f02-9eb1-53311ab94cba","name":"Igor","surname":"Menshenin","pat_name":"Vladimirovitch","is_member":true,"phone":"(555) 55-55-55","email":"devalio@yandex.ru","accountID":"e7435d15-3d31-4aeb-91dd-08f787dab2d8"}
}
