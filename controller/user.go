package controller

import (
	"../model"
	userRepository "../repository"
	"database/sql"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type Users struct{}

func (c Users) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var error model.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing."
			RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing."
			RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)

		userRepo := userRepository.UserRepository{}
		userID := userRepo.AddUser(db, user)

		if userID <= 0 {
			error.Message = "Server error."
			RespondWithError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = ""

		w.Header().Set("Content-Type", "application/json")
		responseJSON(w, user)
	}
}

func (c Users) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		var jwt model.JWT
		var error model.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing."
			RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing."
			RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		password := user.Password
		email := user.Email

		userRepo := userRepository.UserRepository{}
		user = userRepo.GetUser(db, user, email)
		hashedPassword := user.Password

		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err != nil {
			error.Message = "Invalid User credentials"
			RespondWithError(w, http.StatusUnauthorized, error)
			return
		}

		token, err := GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		jwt.Token = token

		responseJSON(w, jwt)
	}
}

func GenerateToken(user model.User) (string, error) {
	var err error
	secret := "secret"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "myAmazingServer",
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}


func RespondWithError(w http.ResponseWriter, status int, error model.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func responseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

