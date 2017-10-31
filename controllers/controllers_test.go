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

func TestMain(m *testing.M) {
	containerID := tUtils.Setup()
	m.Run()
	tUtils.Shutdown(containerID)
}

func Test_Can_Register(t *testing.T) {
	request := createRegisterRequest("test@test.com", "test", "test", "password", "password", t)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)

	handler.ServeHTTP(rr, request)

	if rr.Result().StatusCode != http.StatusCreated {
		t.Errorf("Status is incorrect got: %v", rr.Result().StatusCode)
	}
}

func Test_Can_Not_Register_Duplicate(t *testing.T) {
	handler := http.HandlerFunc(RegisterHandler)

	request := createRegisterRequest("test@test.com", "test", "test", "password", "password", t)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, request)

	//Add the new one with the same email
	rr = httptest.NewRecorder()
	request = createRegisterRequest("test@test.com", "test", "test", "password", "password", t)
	handler.ServeHTTP(rr, request)

	if rr.Result().StatusCode != http.StatusConflict {
		t.Errorf("Status is incorrect got: %v", rr.Result().StatusCode)
	}
}

func createRegisterRequest(email, firstname, lastname, password, confirmation string, t *testing.T) *http.Request {
	registerModel := &RegisterModel{
		Email:                email,
		Firstname:            firstname,
		Lastname:             lastname,
		Password:             password,
		PasswordConfirmation: confirmation,
	}
	sJson, err := json.Marshal(registerModel)
	if err != nil {
		t.Fail()
	}

	request, err := http.NewRequest("POST", "/register", bytes.NewBuffer(sJson))
	if err != nil {
		t.Fatal(err)
	}

	return request
}
