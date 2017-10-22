package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	var router *mux.Router
	router = mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/register", RegisterHandler).Methods(http.MethodPost)
	router.HandleFunc("/authenticate", RegisterHandler).Methods(http.MethodPost)

	if err := http.ListenAndServe(":8080", handlers.LoggingHandler(os.Stdout, router)); err != nil {
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

	registerModel := RegisterModel{}
	if err = json.Unmarshal(body, &registerModel); err != nil {
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

	if registerModel.Password == registerModel.PasswordConfirmation {
		w.WriteHeader(http.StatusBadRequest)
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
