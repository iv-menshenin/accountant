package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/iv-menshenin/accountant/model"
)

func createNewAccount(server, token string, _ int) model.Account {
	var acc = randomAccount()
	var r = bytes.NewBufferString("")
	encoder := json.NewEncoder(r)
	if err := encoder.Encode(acc); err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/accounts", server), r)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("unexpected http code: %d", resp.StatusCode))
	}
	var account struct {
		Meta ResponseMeta  `json:"meta"`
		Data model.Account `json:"data,omitempty"`
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&account); err != nil {
		panic(err)
	}
	if err = resp.Body.Close(); err != nil {
		panic(err)
	}
	return account.Data
}

func createNewPersons(server, token string, account model.Account) {
	var personData = randomPerson()
	var r = bytes.NewBufferString("")
	encoder := json.NewEncoder(r)
	if err := encoder.Encode(personData); err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/accounts/%s/persons", server, account.AccountID), r)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("unexpected http code: %d", resp.StatusCode))
	}
	var person struct {
		Meta ResponseMeta `json:"meta"`
		Data model.Person `json:"data,omitempty"`
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&person); err != nil {
		panic(err)
	}
	if err = resp.Body.Close(); err != nil {
		panic(err)
	}
	if rand.Intn(5) == 0 {
		createNewPersons(server, token, account)
	}
}

func createNewObjects(server, token string, account model.Account) {
	var objectData = randomObject()
	var r = bytes.NewBufferString("")
	encoder := json.NewEncoder(r)
	if err := encoder.Encode(objectData); err != nil {
		panic(err)
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/accounts/%s/objects", server, account.AccountID), r)
	if err != nil {
		panic(err)
	}
	req.Header.Set("X-Auth-Token", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("unexpected http code: %d", resp.StatusCode))
	}
	var person struct {
		Meta ResponseMeta `json:"meta"`
		Data model.Object `json:"data,omitempty"`
	}
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	if err = decoder.Decode(&person); err != nil {
		panic(err)
	}
	if err = resp.Body.Close(); err != nil {
		panic(err)
	}
	if rand.Intn(5) == 0 {
		createNewObjects(server, token, account)
	}
}
