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
const jwtKey string = "hB2sPfLoqJKIRGE_WF8OERaZBchR1S1urvKCWUEMGQ7"
const expHrs time.Duration = (12 * time.Hour)

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
		expiry := time.Now().Add(expHrs)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"UserId":   user.UserId,
			"UserName": user.UserName,
			"RegisteredClaims": jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiry),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		})

		signedToken, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			log.Println(err.Error())
		}

		cookie := http.Cookie{
			Name:    "_jwt",
			Value:   signedToken,
			Expires: expiry,
		}
		http.SetCookie(res, &cookie)

		var auth = AuthResponse{
			UserId:   user.UserId,
			UserName: user.UserName,
			Expires:  expiry,
		}

		response := JsonResponse{
			Data:   auth,
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
		expiry := time.Now().Add(expHrs)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"UserId":   user.UserId,
			"UserName": user.UserName,
			"RegisteredClaims": jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expiry),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		})

		signedToken, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			log.Println(err.Error())
		}

		cookie := http.Cookie{
			Name:    "_jwt",
			Value:   signedToken,
			Expires: expiry,
		}
		http.SetCookie(res, &cookie)
		var auth = AuthResponse{
			UserId:   user.UserId,
			UserName: user.UserName,
			Expires:  expiry,
		}

		response := JsonResponse{
			Data:   auth,
			Status: 200,
		}
		response.write(res)
	}
}

func postHandler(res http.ResponseWriter, req *http.Request) {
	auth := req.Header.Get("Authorization")
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(strings.Replace(auth, "Bearer ", "", 1), claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
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
		userId := int(claims["UserId"].(float64)) //not sure why it's float64

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
				UserName:  claims["UserName"].(string),
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
	id := req.PostFormValue("id")
	auth := req.Header.Get("Authorization")

	//log.Println("content id", id)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(strings.Replace(auth, "Bearer ", "", 1), claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	})
	log.Println("token", token)
	userId := int(claims["UserId"].(float64)) //not sure why it's float64

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
			if content.UserId != userId {
				response := JsonResponse{
					Errors: []string{"Unauthorized"},
					Status: 500,
				}
				response.write(res)
				return
			} else {
				// happy path
				content.Access = true
				response := JsonResponse{
					Data:   content,
					Status: 200,
				}
				response.write(res)
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

func publishHandler(res http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	id := query["id"]
	status := query["status"]
	//auth := req.Header.Get("Authorization")
	auth, err := req.Cookie("_jwt")
	if err != nil {
		response := JsonResponse{
			Errors: []string{"Cookie is invalid: " + err.Error()},
			Status: 500,
		}
		response.write(res)
		return

	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(auth.Value, claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		response := JsonResponse{
			Errors: []string{"User id is invalid: " + err.Error()},
			Status: 500,
		}
		response.write(res)
		return

	}

	//log.Println("token", token)
	_ = token
	userId := int(claims["UserId"].(float64)) //not sure why it's float64

	Id, err := strconv.Atoi(id[0])
	Status, err := strconv.Atoi(status[0])
	if err != nil {
		response := JsonResponse{
			Errors: []string{"Content id is invalid: " + err.Error()},
			Status: 500,
		}
		response.write(res)
		return
	}
	content := GetContent(Id)

	if content.UserId != userId {
		response := JsonResponse{
			Errors: []string{"Unauthorized"},
			Status: 500,
		}
		response.write(res)
		return
	}

	// happy path
	content.Status = Status
	_, err = PostContent(&content)
	if err != nil {
		response := JsonResponse{
			Errors: []string{err.Error()},
			Status: 500,
		}
		response.write(res)
	} else {
		/*
			response := JsonResponse{
				Data:   contents,
				Status: 200,
			}
			response.write(res) */
		log.Println("REDIRECT to /")
		http.Redirect(res, req, "/", http.StatusFound)
	}

}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/join", joinHandler)
	http.HandleFunc("/content", contentHandler)
	http.HandleFunc("/publish", publishHandler)
	http.HandleFunc("/post", postHandler)

	fmt.Printf("listening on http://localhost:%d/\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
