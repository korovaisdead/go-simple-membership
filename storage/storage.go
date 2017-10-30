package storage

import (
	c "github.com/korovaisdead/go-simple-membership/config"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"time"
)

func getSession() (*mgo.Session, error) {
	config, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

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

func GetUsers() (*[]User, error) {
	config, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var users []User
	if err = session.DB(config.Db.Database).C("Users").Find(nil).All(&users); err != nil {
		return nil, err
	}

	return &users, nil
}

func LoadUserByEmail(email string) (*User, error) {
	config, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var user User
	if err = session.DB(config.Db.Database).C("Users").Find(bson.M{"email": email}).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func SaveUser(firstname, lastname, email, phone, password string) error {
	config, err := c.GetConfig()
	if err != nil {
		return err
	}

	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()

	salt := getRandomString()
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), config.Security.BcryptCost)
	if err != nil {
		return err
	}

	user := User{
		ID:        bson.NewObjectId(),
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Phone:     phone,
		Password:  string(hash),
		Salt:      salt,
	}

	return session.DB(config.Db.Database).C("Users").Insert(user)
}

func getRandomString() string {
	config, err := c.GetConfig()
	if err != nil {
		panic(err)
	}

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, config.Security.SaltLength)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
