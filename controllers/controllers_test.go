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

	if rr.Code != http.StatusCreated {
		t.Fatalf("Status is incorrect got: %v", rr.Code)
	}
}

func Test_Can_Not_Register_Duplicate(t *testing.T) {
	handler := http.HandlerFunc(RegisterHandler)

	request := createRegisterRequest("test1@test.com", "test", "test", "password", "password", t)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, request)

	//Add the new one with the same email
	rr = httptest.NewRecorder()
	request = createRegisterRequest("test1@test.com", "test", "test", "password", "password", t)
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusConflict {
		t.Fatalf("Status is incorrect got: %v", rr.Code)
	}
}

func Test_Can_Not_Create_User_With_Empty_Email(t *testing.T) {
	handler := http.HandlerFunc(RegisterHandler)

	request := createRegisterRequest("", "test", "test", "password", "password", t)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Status is incorrect got: %v", rr.Code)
	}
}

func Test_Can_Not_Create_User_With_Empty_Password(t *testing.T) {
	handler := http.HandlerFunc(RegisterHandler)

	request := createRegisterRequest("test1@gmail.com", "test", "test", "", "", t)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Status is incorrect got: %v", rr.Code)
	}
}

func Test_Can_Not_Create_User_With_Not_Equal_Passwords(t *testing.T) {
	handler := http.HandlerFunc(RegisterHandler)

	request := createRegisterRequest("test12@gmail.com", "test", "test", "c", "not equal", t)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Status is incorrect got: %v", rr.Code)
	}
}

func Test_Can_Not_Create_User_With_Wrong_Body(t *testing.T) {
	handler := http.HandlerFunc(RegisterHandler)

	request, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte("")))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Status is incorrect got: %v", rr.Code)
	}
}

func Test_User_Can_Login(t *testing.T) {
	request := createRegisterRequest("test_of_the_login@test.com", "test", "test", "password", "password", t)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, request)

	request = createLoginModel("test_of_the_login@test.com", "password", t)
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, request)

	if rr.Code != http.StatusOK {
		t.Fatalf("Login failed. Status: %v", rr.Code)
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

func createLoginModel(email, password string, t *testing.T) *http.Request {
	loginModel := &LoginModel{
		Email:    email,
		Password: password,
	}

	sJson, err := json.Marshal(loginModel)
	if err != nil {
		t.Fail()
	}

	request, err := http.NewRequest("POST", "/login", bytes.NewBuffer(sJson))
	if err != nil {
		t.Fatal(err)
	}

	return request
}
