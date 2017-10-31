package controllers

import (
	"bytes"
	"encoding/json"
	tUtils "github.com/korovaisdead/go-simple-membership/utils/testing"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mongoImage string = "mongo"
)

func Test_Register_Handler(t *testing.T) {
	containerID := tUtils.Setup(t)
	defer tUtils.Shutdown(containerID)

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
	if rr.Result().StatusCode != http.StatusCreated {
		t.Errorf("Status is incorrect got: %v", rr.Result().StatusCode)
	}
}
