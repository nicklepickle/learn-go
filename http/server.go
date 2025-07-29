package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

func postHandler(res http.ResponseWriter, req *http.Request) {
	auth := req.Header.Get("Authorization")
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(strings.Replace(auth, "Bearer ", "", 1), claims, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	})
	if err != nil {
		response := JsonResponse{
			Errors: []string{"User id is invalid: " + err.Error()},
			Status: 500,
		}
		response.write(res)
		return
	} else {
		//log.Println("token", token)
		userId := int(claims["id"].(float64)) //not sure why it's float64

		id := req.PostFormValue("id")
		title := req.PostFormValue("title")
		body := req.PostFormValue("body")
		//status := req.PostFormValue("status")

		content := &Content{}

		log.Println("user id = ", userId)
		if id == "0" || id == "" {
			// new content
			content = &Content{
				ContentId: 0,
				UserId:    userId,
				UserName:  claims["user"].(string),
				Title:     title,
				Body:      body,
				Created:   time.Now(),
				Updated:   time.Now(),
				Status:    1,
			}
		} else {
			// edit content
			userContent := GetUserContent(userId)
			contentId, err := strconv.Atoi(id)
			if err != nil {
				response := JsonResponse{
					Errors: []string{err.Error()},
					Status: 500,
				}
				response.write(res)
			} else {
				for _, c := range userContent {
					if c.ContentId == contentId {
						// does the content belong to this user?
						if c.UserId != userId {
							response := JsonResponse{
								Errors: []string{"Unauthorized"},
								Status: 500,
							}
							response.write(res)
							return
						} else {
							content = &c
							content.Title = title
							content.Body = body
							content.Updated = time.Now()
						}

						break
					}
				}
			}
		}

		contents, err := PostContent(content)
		if err != nil {
			response := JsonResponse{
				Errors: []string{err.Error()},
				Status: 500,
			}
			response.write(res)
		} else {
			response := JsonResponse{
				Data:   contents,
				Status: 200,
			}
			response.write(res)
		}
	}

}

func contentHandler(res http.ResponseWriter, req *http.Request) {
	auth := req.Header.Get("Authorization")
	id := req.PostFormValue("id")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(strings.Replace(auth, "Bearer ", "", 1), claims, func(token *jwt.Token) (any, error) {
		return []byte(key), nil
	})
	log.Println("token", token)
	userId := int(claims["id"].(float64)) //not sure why it's float64

	if err != nil {
		response := JsonResponse{
			Errors: []string{"User id is invalid: " + err.Error()},
			Status: 500,
		}
		response.write(res)
	} else if id != "" {
		Id, err := strconv.Atoi(id)
		if err != nil {
			response := JsonResponse{
				Errors: []string{"Content id is invalid: " + err.Error()},
				Status: 500,
			}
			response.write(res)
			return
		} else {
			content := GetContent(Id)
			if content.UserId == userId {
				response := JsonResponse{
					Data:   content,
					Status: 200,
				}
				response.write(res)
			} else {
				response := JsonResponse{
					Errors: []string{"Unauthorized"},
					Status: 500,
				}
				response.write(res)
				return
			}
		}
	} else {

		//log.Println("user id = ", userId)
		userContent := GetUserContent(userId)
		response := JsonResponse{
			Data:   userContent,
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
	http.HandleFunc("/post", postHandler)

	fmt.Printf("listening on http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
