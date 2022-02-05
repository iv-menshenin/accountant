package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iv-menshenin/accountant/model"
)

type (
	ResponseMeta struct {
		Status  string `json:"status"`
		Message string `json:"message,omitempty"`
	}
	AuthResponse struct {
		Meta ResponseMeta   `json:"meta"`
		Data model.AuthData `json:"data,omitempty"`
	}
)

func getAuthData(server string) string {
	var login = model.AuthQuery{
		Login:    "test",
		Password: "test",
	}
	var r = bytes.NewBufferString("")
	encoder := json.NewEncoder(r)
	if err := encoder.Encode(login); err != nil {
		panic(err)
	}
	resp, err := http.Post(fmt.Sprintf("%s/auth/login", server), "application/json", r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	var token AuthResponse
	if err := decoder.Decode(&token); err != nil {
		panic(err)
	}
	return token.Data.JWT
}
