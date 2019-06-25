package controller

import (
	"../model"
	bookRepository "../repository"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Books struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Books) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book

		books := []model.Book{}
		bookRepo := bookRepository.BookRepository{}

		books = bookRepo.GetBooks(db, book, books)

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(books)
	}
}

func (c Books) GetBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		params := mux.Vars(r)

		bookRepo := bookRepository.BookRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		book = bookRepo.GetBook(db, book, id)

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(book)
	}
}

func (c Books) AddBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		var bookID int

		json.NewDecoder(r.Body).Decode(&book)

		bookRepo := bookRepository.BookRepository{}
		bookID = bookRepo.AddBook(db, book)

		book.Id = bookID
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("location",  "/books/" + strconv.Itoa(bookID))
		json.NewEncoder(w).Encode(book)
	}
}

func (c Books) UpdateBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book model.Book
		json.NewDecoder(r.Body).Decode(&book)

		bookRepo := bookRepository.BookRepository{}
		_ = bookRepo.UpdateBook(db, book)

		w.WriteHeader(204)
	}
}

func (c Books) RemoveBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		bookRepo := bookRepository.BookRepository{}

		id, err := strconv.Atoi(params["id"])
		logFatal(err)

		_ = bookRepo.RemoveBook(db, id)

		w.WriteHeader(204)
	}
}

