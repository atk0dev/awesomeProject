package repository

import (
	"../model"
	"database/sql"
)

type BookRepository struct{}

func (b BookRepository) GetBooks(db *sql.DB, book model.Book, books []model.Book) []model.Book {
	rows, err := db.Query("select id, title, author, year from books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(books, book)
	}

	return books
}

func (b BookRepository) GetBook(db *sql.DB, book model.Book, id int) model.Book {
	rows := db.QueryRow("select id, title, author, year from books where id=$1", id)

	err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	return book
}

func (b BookRepository) AddBook(db *sql.DB, book model.Book) int {
	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) RETURNING id;",
		book.Title, book.Author, book.Year).Scan(&book.Id)

	logFatal(err)

	return book.Id
}

func (b BookRepository) UpdateBook(db *sql.DB, book model.Book) int {
	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id",
		&book.Title, &book.Author, &book.Year, &book.Id)

	logFatal(err)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	return int(rowsUpdated)
}

func (b BookRepository) RemoveBook(db *sql.DB, id int) int {
	result, err := db.Exec("delete from books where id = $1", id)
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	return int(rowsDeleted)
}
