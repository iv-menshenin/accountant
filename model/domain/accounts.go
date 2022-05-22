package domain

import (
	"time"

	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	// Account представляет собой верхний уровень иерархии - лицевой счет
	Account struct {
		AccountID   uuid.UUID `json:"account_id"`
		Persons     []Person  `json:"persons"`
		Objects     []Object  `json:"objects"`
		AccountData `json:",inline"`
	}
	AccountData struct {
		// Account это номер лицевого счета. Т.к. не является первичным ключем, то может изменяться без проблем с консистентностью
		Account string `bson:"account" json:"account"`
		// CadNumber это кадастровый номер [АА:ВВ:CCCCСCC:КК]
		//  АА — кадастровый округ
		//  ВВ — кадастровый район
		//  CCCCCCС — кадастровый квартал состоит из 6 или 7 цифр
		//  КК — номер объекта недвижимости
		CadNumber string `bson:"cad_number" json:"cad_number"`
		// AgreementNum это номер договора (купли/продажи или аренды)
		AgreementNum string `bson:"agreement" json:"agreement"`
		// AgreementDate это дата договора, если есть
		AgreementDate *time.Time `bson:"agreement_date,omitempty" json:"agreement_date,omitempty"`
		// PurchaseKind вид приобретения
		PurchaseKind string `bson:"purchase_kind" json:"purchase_kind"`
		// PurchaseDate дата приобретения
		PurchaseDate time.Time `bson:"purchase_date" json:"purchase_date"`
		// Comment просто текстовый комментарий для заметок
		Comment string `bson:"comment,omitempty" json:"comment,omitempty"`
	}
	// Person представляет собой физическое лицо, закрепленное за лицевым счетом
	Person struct {
		PersonID   uuid.UUID `bson:"person_id" json:"person_id"`
		PersonData `bson:",inline" json:",inline"`
	}
	PersonData struct {
		Name    string `bson:"name" json:"name"`
		Surname string `bson:"surname" json:"surname"`
		PatName string `bson:"pat_name" json:"pat_name"`
		// DOB дата рождения
		DOB *time.Time `bson:"dob,omitempty" json:"dob,omitempty"`
		// IsMember это признак, является ли членом общества
		IsMember bool `bson:"is_member" json:"is_member"`
		// Phone это номер телефона
		Phone string `bson:"phone,omitempty" json:"phone,omitempty"`
		// EMail это адрес электронной почты
		EMail string `bson:"email,omitempty" json:"email,omitempty"`
	}
	// Object представляет собой территорию, которая закреплена за лицевым счетом (дачный участок)
	Object struct {
		ObjectID   uuid.UUID `bson:"object_id" json:"object_id"`
		ObjectData `bson:",inline" json:",inline"`
	}
	ObjectData struct {
		PostalCode string `bson:"postal_code" json:"postal_code"`
		City       string `bson:"city" json:"city"`
		Village    string `bson:"village" json:"village"`
		Street     string `bson:"street" json:"street"`
		Number     int    `bson:"number" json:"number"`
		// Area это площадь территории
		Area float64 `bson:"area,omitempty" json:"area,omitempty"`
		// FIXME CadNumber кадастровый номер
		CadNumber string `bson:"cad_number" json:"cad_number"`
	}
)
