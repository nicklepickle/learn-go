package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt" // import bcrypt
)

const bcryptCost int = 8
const minPassLen int = 8
const port int = 8080

/*
// test

	func apiHandler(res http.ResponseWriter, req *http.Request) {
		json := fmt.Sprintf(`{"method":"%s","path":"%s"}`, req.Method, req.URL)
		fmt.Println(json)
		res.Write([]byte(json))
	}
*/
type JsonResponse struct {
	Status int
	Data   any
	Errors []string
}

func writeResponse(res http.ResponseWriter, data any, status int, errors []string) {
	response := JsonResponse{
		Status: status,
		Data:   data,
		Errors: errors,
	}
	bytes, err := json.Marshal(response)
	if err != nil {
		log.Println(err.Error())
		return
	}
	res.WriteHeader(status)
	res.Header().Add("Content-Type", "application/json")

	res.Write(bytes)
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	username := req.PostFormValue("user")
	pw := req.PostFormValue("password")

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
			//http.Redirect(res, req, "/?ok=true", http.StatusSeeOther)

			writeResponse(res, user, 200, nil)
			return
		} else {
			log.Println(err.Error())
		}
	}
	//http.Redirect(res, req, "/?ok=false", http.StatusSeeOther)
	errors := []string{"Log in failed"}
	writeResponse(res, nil, 401, errors)
}

func joinHandler(res http.ResponseWriter, req *http.Request) {
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
		errors = append(errors, fmt.Sprintf("Email addess %s is invalid %s", email, err.Error()))
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw1), bcryptCost)
	if err != nil {
		errors = append(errors, fmt.Sprintf("Password is invalid %s", err.Error()))
	}

	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Println(e)
		}
		http.Redirect(res, req, "/?ok=false", http.StatusSeeOther)
		return
	}

	//fmt.Printf("POST user=%s email=%s hash=%s\n", username, email, hash)

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
		log.Println(err.Error())
		// something else went wrong
		http.Redirect(res, req, "/?ok=false", http.StatusSeeOther)
		return
	}

	// happy path
	http.Redirect(res, req, "/?ok=true", http.StatusSeeOther)
}

func main() {
	ReadUsers()

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	//http.HandleFunc("/api", apiHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/join", joinHandler)

	fmt.Printf("listening on http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
