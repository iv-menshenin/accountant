package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/iv-menshenin/accountant/model/domain"
	"github.com/iv-menshenin/accountant/test/randomizer"
)

const (
	randomAccountsCount = 100
	serverName          = "acc.devaliada.ru"
	serverPort          = "8080"
	serverProto         = "http"
)

var rnd = randomizer.New()

func main() {
	rand.Seed(time.Now().UnixNano())
	var server = fmt.Sprintf("%s://%s:%s", serverProto, serverName, serverPort)
	var token = getAuthData(server)

	for i := 0; i < randomAccountsCount; i++ {
		account := createNewAccount(server, token, i)
		createNewPersons(server, token, account)
		createNewObjects(server, token, account)
	}
}

var counter = 0

func randomAccount() domain.AccountData {
	counter++
	agrDate := time.Now().Add(time.Duration(rand.Intn(1500)) * time.Hour * -24)
	return domain.AccountData{
		Account:       fmt.Sprintf("A-0000%d", counter),
		CadNumber:     fmt.Sprintf("%d:%d:000%d:%d", 10+rand.Intn(90), 10+rand.Intn(90), 100+rand.Intn(900), 100+rand.Intn(900)),
		AgreementNum:  fmt.Sprintf("%d", 100+rand.Intn(1900)),
		AgreementDate: &agrDate,
		PurchaseKind:  "",
		PurchaseDate:  agrDate,
		Comment:       "",
	}
}

func randomPerson() domain.PersonData {
	agrDate := time.Now().Add(time.Duration(rand.Intn(600)) * time.Hour * -24 * 30)
	return domain.PersonData{
		Name:     rnd.RandomName(),
		Surname:  rnd.RandomSurname(),
		PatName:  "",
		DOB:      &agrDate,
		IsMember: rand.Intn(10) > 3,
		Phone:    fmt.Sprintf("(%d) %d-%d-%d", 100+rand.Intn(900), 100+rand.Intn(900), 10+rand.Intn(90), 10+rand.Intn(90)),
		EMail:    "test@yandex.ru",
	}
}

func randomObject() domain.ObjectData {
	return domain.ObjectData{
		PostalCode: "",
		City:       "",
		Village:    "",
		Street:     rnd.RandomStreet(),
		Number:     counter,
		Area:       float64(rand.Intn(5)) + 399,
	}
}
