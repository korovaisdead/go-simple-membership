package storage

import (
	c "github.com/korovaisdead/go-simple-membership/config"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"time"
)

func GetUsers() (*[]User, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	var users []User
	config := c.Get()
	if err = session.DB(config.Db.Database).C("Users").Find(nil).All(&users); err != nil {
		return nil, err
	}

	return &users, nil
}

func LoadUserByEmail(email string) (*User, error) {
	session, err := getSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	config := c.Get()
	var user User
	if err = session.DB(config.Db.Database).C("Users").Find(bson.M{"email": email}).One(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func SaveUser(firstname, lastname, email, phone, password string) error {
	session, err := getSession()
	if err != nil {
		return err
	}
	defer session.Close()

	config := c.Get()
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
	config := c.Get()

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, config.Security.SaltLength)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
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
