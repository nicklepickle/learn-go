package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const port int = 8080
const key string = "hB2sPfLoqJKIRGE_WF8OERaZBchR1S1urvKCWUEMGQ7"

type JsonResponse struct {
	Status int
	Data   any
	Errors []string
}

func (r *JsonResponse) write(w http.ResponseWriter) error {
	bytes, err := json.Marshal(r)
	if err != nil {
		return err
	}
	w.WriteHeader(r.Status)
	w.Header().Add("Content-Type", "application/json")
	w.Write(bytes)
	return nil
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	user, errors := AuthenticateUser(req)

	if len(errors) > 0 {
		response := JsonResponse{
			Errors: errors,
			Status: 401,
		}
		response.write(res)
	} else {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":   user.UserId,
			"user": user.UserName,
		})

		signedToken, err := token.SignedString([]byte(key))
		if err != nil {
			log.Println(err.Error())
		}

		cookie := http.Cookie{
			Name:  "_jwt",
			Value: signedToken,
		}
		http.SetCookie(res, &cookie)

		response := JsonResponse{
			Data:   signedToken,
			Status: 200,
		}
		response.write(res)
	}
}

func joinHandler(res http.ResponseWriter, req *http.Request) {
	user, errors := RegisterUser(req)

	if len(errors) > 0 {
		response := JsonResponse{
			Errors: errors,
			Status: 401,
		}
		response.write(res)
	} else {
		response := JsonResponse{
			Data:   user,
			Status: 200,
		}
		response.write(res)
	}
}

func contentHandler(res http.ResponseWriter, req *http.Request) {
	auth := req.Header.Get("Authorization")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(strings.Replace(auth, "Bearer ", "", 1), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("token", token)
		//for c := range claims {
		//	fmt.Printf("c: %v\n", c)
		//}
		fmt.Println("user", claims["user"])

		response := JsonResponse{
			Data:   claims,
			Status: 200,
		}
		response.write(res)
	}

}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/content", contentHandler)

	fmt.Printf("listening on http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
