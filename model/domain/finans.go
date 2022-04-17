package domain

import (
	"time"

	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	// Payment представляет внесенную оплату
	Payment struct {
		PaymentID   uuid.UUID  `bson:"payment_id" json:"payment_id"`
		AccountID   uuid.UUID  `bson:"account_id" json:"account_id"`
		PersonID    *uuid.UUID `bson:"person_id" json:"person_id"`
		ObjectID    *uuid.UUID `bson:"object_id" json:"object_id"`
		Period      Period     `bson:"period" json:"period"`
		Target      TargetHead `bson:"target" json:"target"`
		Payment     float64    `bson:"payment" json:"payment"`
		PaymentDate *time.Time `bson:"payment_date" json:"payment_date"`
		Receipt     string     `bson:"receipt" json:"receipt"`
	}
	TargetHead struct {
		TargetID uuid.UUID `bson:"target_id" json:"target_id"`
		Type     string    `bson:"type" json:"type"`
	}
	Period struct {
		Month int `bson:"month" json:"month"`
		Year  int `bson:"year" json:"year"`
	}
	// Target содержит описание целевых взносов
	Target struct {
		TargetHead `bson:",inline" json:",inline"`
		TargetData `bson:",inline" json:",inline"`
	}
	// TargetData содержит изменяемые данные структуры Target
	TargetData struct {
		Period  Period     `bson:"period" json:"period"`
		Closed  *time.Time `bson:"closed" json:"closed"`
		Cost    float64    `bson:"cost" json:"cost"`
		Comment string     `bson:"comment" json:"comment"`
	}
	// Bill описывает начисления (счет на оплату)
	Bill struct {
		BillID    uuid.UUID   `bson:"bill_id" json:"bill_id"`
		AccountID uuid.UUID   `bson:"account_id" json:"account_id"`
		PersonID  *uuid.UUID  `bson:"person_id" json:"person_id"`
		ObjectID  *uuid.UUID  `bson:"object_id" json:"object_id"`
		Period    Period      `bson:"period" json:"period"`
		Target    TargetHead  `bson:"target" json:"target"`
		Bill      float64     `bson:"bill" json:"bill"`
		Payments  []uuid.UUID `bson:"payment_linked" json:"payment_linked"`
	}
)
