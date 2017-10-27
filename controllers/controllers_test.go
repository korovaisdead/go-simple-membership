package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/korovaisdead/go-simple-membership/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Register_Handler(t *testing.T) {
	_, err := config.BuildConfig("test")

	registerModel := &RegisterModel{
		Email:                "test@test.com",
		Firstname:            "test",
		Lastname:             "test",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	sJson, err := json.Marshal(registerModel)
	if err != nil {
		t.Fail()
	}

	request, err := http.NewRequest("POST", "/register", bytes.NewBuffer(sJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)

	handler.ServeHTTP(rr, request)
}
