package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost int = 8
const minPassLen int = 8

var Users map[string]User = make(map[string]User)

type User struct {
	UserId   int
	UserName string
	Email    string
	Hash     string
	Created  time.Time
	Seen     time.Time
	Status   int
}

type AuthResponse struct {
	UserId   int
	UserName string
	Expires  time.Time
}

func (u *User) save() error {
	Users[u.UserName] = *u
	err := writeUsers()
	return err
}

func AuthenticateUser(req *http.Request) (*User, []string) {
	username := req.PostFormValue("user")
	pw := req.PostFormValue("password")
	if len(Users) == 0 {
		err := readUsers()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	user, exists := Users[username]
	if exists && user.Status == 1 {
		//fmt.Printf("USERS user=%s password=%s\n", user.UserName, user.Hash)
		err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(pw))
		if err == nil {
			user.Seen = time.Now()
			err = user.save()
			if err != nil {
				log.Println(err.Error())
			}
			return &user, nil
		}
	}

	return nil, []string{"Log in failed"}
}

func RegisterUser(req *http.Request) (*User, []string) {
	username := req.PostFormValue("user")
	email := req.PostFormValue("email")
	pw1 := req.PostFormValue("password")
	pw2 := req.PostFormValue("confirm")

	// validate
	errors := []string{}

	match, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !match {
		errors = append(errors, "User name contains non-alphanumeric characters")
	}

	_, exists := Users[username]
	if exists {
		errors = append(errors, fmt.Sprintf("User name %s is already taken", username))
	}

	if len(pw1) < minPassLen {
		errors = append(errors, fmt.Sprintf("Password is less than %d characters", minPassLen))
	}

	if pw1 != pw2 {
		errors = append(errors, "Password and confirmation do not match")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Email addess %s is invalid", email))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw1), bcryptCost)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Password is invalid %s", err.Error()))
	}

	user := User{
		UserId:   len(Users) + 1, //  let db do this
		UserName: username,
		Email:    email,
		Hash:     string(hash),
		Created:  time.Now(),
		Seen:     time.Now(),
		Status:   1,
	}
	err = user.save()
	if err != nil {
		errors = append(errors, fmt.Sprintf("Could not save user %s", err.Error()))
	}

	if len(errors) > 0 {
		return nil, errors
	}
	return &user, nil
}

func readUsers() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	bytes, err := os.ReadFile(path + "/users.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &Users)
	//fmt.Println("ReadUsers: " + path + "/users.json = " + string(bytes))
	if err != nil {
		return err
	}
	return nil
}

func writeUsers() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(Users)
	if err != nil {
		return err
	}
	err = os.WriteFile(path+"/users.json", bytes, 0777)
	if err != nil {
		return err
	}
	return nil
}
