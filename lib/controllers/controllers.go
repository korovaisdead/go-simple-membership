package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/korovaisdead/go-simple-membership/lib/config"
	"github.com/korovaisdead/go-simple-membership/lib/storage"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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

type AuthResponse struct {
	Token string `json:"token"`
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

	tokenString, err := getToken(string(user.ID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	redisClient := storage.GetRedisClient()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = redisClient.Set(user.ID.Hex(), user.Email, time.Minute*60)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(&AuthResponse{*tokenString})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(response))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func getToken(id string) (*string, error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	claims := &MyCustomClaims{
		id,
		jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 24).Unix()},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.SecuritySecretWorld))
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

type MyCustomClaims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}
