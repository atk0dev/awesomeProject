package repository

import (
	"../model"
	"database/sql"
)

type UserRepository struct{}

func (b UserRepository) AddUser(db *sql.DB, user model.User) int {
	err := db.QueryRow("insert into users (email, password) values($1, $2) RETURNING id;",
		user.Email, user.Password).Scan(&user.Id)
	logFatal(err)
	return user.Id
}

func (b UserRepository) GetUser(db *sql.DB, user model.User, email string) model.User {
	rows := db.QueryRow("select id, email, password from users where email=$1", email)
	err := rows.Scan(&user.Id, &user.Email, &user.Password)
	logFatal(err)
	return user
}
