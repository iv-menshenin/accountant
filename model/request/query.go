package request

import (
	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/utils/uuid"
)

type (
	GetAccountQuery struct {
		ID uuid.UUID
	}
	PostAccountQuery struct {
		AccountData domain.AccountData
	}
	PutAccountQuery struct {
		ID      uuid.UUID
		Account domain.AccountData
	}
	DeleteAccountQuery struct {
		ID uuid.UUID
	}
	FindAccountsQuery struct {
		Account *string
	}

	PostPersonQuery struct {
		AccountID  uuid.UUID
		PersonData domain.PersonData
	}
	GetPersonQuery struct {
		AccountID uuid.UUID
		PersonID  uuid.UUID
	}
	PutPersonQuery struct {
		AccountID  uuid.UUID
		PersonID   uuid.UUID
		PersonData domain.PersonData
	}
	DeletePersonQuery struct {
		AccountID uuid.UUID
		PersonID  uuid.UUID
	}
	FindPersonsQuery struct {
		AccountID      *uuid.UUID
		PersonFullName *string
	}

	PostObjectQuery struct {
		AccountID  uuid.UUID
		ObjectData domain.ObjectData
	}
	GetObjectQuery struct {
		AccountID uuid.UUID
		ObjectID  uuid.UUID
	}
	PutObjectQuery struct {
		AccountID  uuid.UUID
		ObjectID   uuid.UUID
		ObjectData domain.ObjectData
	}
	DeleteObjectQuery struct {
		AccountID uuid.UUID
		ObjectID  uuid.UUID
	}
	FindObjectsQuery struct {
		AccountID *uuid.UUID
		Address   *string
		Number    *int
	}

	GetTargetQuery struct {
		TargetID uuid.UUID
	}
	PostTargetQuery struct {
		Type   string
		Target domain.TargetData
	}
	PutTargetQuery struct {
		TargetID uuid.UUID
		Target   domain.TargetData
	}
	DeleteTargetQuery struct {
		TargetID uuid.UUID
	}
	FindTargetsQuery struct {
		ShowClosed bool
		Period     *domain.Period
	}

	PostBillQuery struct {
		AccountID uuid.UUID
		Data      domain.BillData
	}
	PutBillQuery struct {
		BillID uuid.UUID
		Data   domain.BillData
	}
	GetBillQuery struct {
		BillID uuid.UUID
	}
	DeleteBillQuery struct {
		BillID uuid.UUID
	}
	FindBillsQuery struct {
		AccountID *uuid.UUID
		TargetID  *uuid.UUID
		Period    *domain.Period
	}

	PostPaymentQuery struct {
		AccountID uuid.UUID
		Data      domain.PaymentData
	}
	GetPaymentQuery struct {
		PaymentID uuid.UUID
	}
	DeletePaymentQuery struct {
		PaymentID uuid.UUID
	}
	FindPaymentsQuery struct {
		AccountID *uuid.UUID
		PersonID  *uuid.UUID
		ObjectID  *uuid.UUID
		TargetID  *uuid.UUID
	}
)
