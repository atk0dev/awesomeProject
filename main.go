package main

import (
	"./controller"
	"./driver"
	"./model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

type WithCORS struct {
	r *mux.Router
}

func (s *WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	if origin := req.Header.Get("Origin"); origin != "" {
		res.Header().Set("Access-Control-Allow-Origin", origin)
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	if req.Method == "OPTIONS" {
		return
	}

	s.r.ServeHTTP(res, req)
}

func main() {

	db := driver.ConnectDB()
	router := mux.NewRouter()

	booksController := controller.Books{}
	router.HandleFunc("/books", TokenVerifyMiddleWare(booksController.GetBooks(db))).Methods("GET")
	router.HandleFunc("/books/{id}", TokenVerifyMiddleWare(booksController.GetBook(db))).Methods("GET")
	router.HandleFunc("/books", TokenVerifyMiddleWare(booksController.AddBook(db))).Methods("POST")
	router.HandleFunc("/books", TokenVerifyMiddleWare(booksController.UpdateBook(db))).Methods("PUT")
	router.HandleFunc("/books/{id}", TokenVerifyMiddleWare(booksController.RemoveBook(db))).Methods("DELETE")

	usersController := controller.Users{}
	router.HandleFunc("/signup", usersController.Signup(db)).Methods("POST")
	router.HandleFunc("/login", usersController.Login(db)).Methods("POST")

	port := 8000
	log.Println("Starting to listen on port ", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), &WithCORS{router}))
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var errorObject model.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte("secret"), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				controller.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				controller.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Message = "Invalid token."
			controller.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
