package main

import (
	"encoding/json"
	"fmt"
	//"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	c "github.com/korovaisdead/go-simple-membership/config"
	"github.com/korovaisdead/go-simple-membership/storage"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var (
	config *c.Configuration
)

func main() {
	config, err := c.GetConfig()
	if err != nil {
		panic(err)
	}

	var router *mux.Router
	router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/register", RegisterHandler).Methods(http.MethodPost)
	router.HandleFunc("/authenticate", LoginHandler).Methods(http.MethodPost)

	if err := http.ListenAndServe(config.Web.Port, handlers.LoggingHandler(os.Stdout, router)); err != nil {
		panic(err)
	}

	select {}
}

//RegisterHandler is the function to handle the registration requests
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	registerModel := &RegisterModel{}
	if err = json.Unmarshal(body, registerModel); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(registerModel.Email)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(registerModel.Password)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if registerModel.Password != registerModel.PasswordConfirmation {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := storage.LoadUserByEmail(registerModel.Email)
	if err != nil {
		if err != mgo.ErrNotFound {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
	}

	if user != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err = storage.SaveUser(registerModel.Firstname, registerModel.Lastname, registerModel.Email, registerModel.Phone, registerModel.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

//LoginHandler represens the handler of the login action
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loginModel := &LoginModel{}
	if err = json.Unmarshal(body, loginModel); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(loginModel.Email)) == 0 || len(strings.TrimSpace(loginModel.Password)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := storage.LoadUserByEmail(loginModel.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginModel.Password+user.Salt)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//RegisterModel structure is represent the entoty to add
type RegisterModel struct {
	Firstname            string `json:"firstname"`
	Lastname             string `json:"lastname"`
	Email                string `json:"email"`
	Phone                string `json:"phone"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"confirmation"`
}

//LoginModel represents the model for the login action
type LoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
