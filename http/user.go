package main

import (
	"encoding/json"
	"os"
	"time"
)

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

func (u *User) save() error {
	Users[u.UserName] = *u
	err := WriteUsers()
	return err
}

func ReadUsers() error {
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

func WriteUsers() error {
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
