package main

import (
	"encoding/json"
	"fmt"
	//"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	config *Configuration
)

func main() {
	c, err := loadConfig()
	if err != nil {
		panic(err)
	}
	config = c

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

	user, err := loadUserByEmail(registerModel.Email)
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

	if err = saveUser(registerModel); err != nil {
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

	user, err := loadUserByEmail(loginModel.Email)
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

func getSession() (*mgo.Session, error) {
	di := &mgo.DialInfo{
		Addrs:    []string{config.Db.Host + config.Db.Port},
		Database: config.Db.Database,
	}
	session, err := mgo.DialWithInfo(di)
	if err != nil {
		return nil, err
	}
	return session, nil
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

//User represents the user model inside database
type User struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Email     string        `bson:"email" json:"email"`
	Firstname string        `bson:"firstname" json:"firstname"`
	Lastname  string        `bson:"lastname" json:"lastname"`
	Password  string        `bson:"password" json:"password"`
	Phone     string        `bson:"phone" json:"phone"`
	Salt      string        `bson:"salt" json:"salt"`
}

func getUsers() (*[]User, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var users []User
	if err = session.DB("Auth").C("Users").Find(nil).All(&users); err != nil {
		return nil, err
	}

	return &users, nil
}

func loadUserByEmail(email string) (*User, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var user User
	if err = session.DB("Auth").C("Users").Find(bson.M{"email": email}).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func saveUser(model *RegisterModel) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()

	salt := getRandomString()
	hash, err := bcrypt.GenerateFromPassword([]byte(model.Password+salt), config.Security.BcryptCost)
	if err != nil {
		return err
	}

	user := User{
		ID:        bson.NewObjectId(),
		Firstname: model.Firstname,
		Lastname:  model.Lastname,
		Email:     model.Email,
		Phone:     model.Phone,
		Password:  string(hash),
		Salt:      salt,
	}

	return session.DB("Auth").C("Users").Insert(user)
}

func getRandomString() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, config.Security.SaltLength)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

type Configuration struct {
	Web struct {
		Port string `json:"port"`
	} `json:"web"`
	Db struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
	} `json:"db"`
	Security struct {
		SaltLength int `json:"saltLength"`
		BcryptCost int `json:"bcryptCost"`
	} `json:"security"`
}

func loadConfig() (*Configuration, error) {
	file, err := os.Open("config.local.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Configuration{}
	if err = decoder.Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}
