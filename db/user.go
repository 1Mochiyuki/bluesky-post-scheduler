package db

import (
	"github.com/1Mochiyuki/gosky/db/queries"
	"github.com/jmoiron/sqlx"
)

func InsertNewUser(db *sqlx.DB, user User) (uint8, error) {
	result, err := db.Exec("INSERT OR IGNORE INTO users (handle) VALUES ($1)", user.Handle)
	if err != nil {
		return 0, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint8(userID), nil
}

func UpdateUser(db *sqlx.DB, user User) (uint8, error) {
	return 0, nil
}

func RemoveSavedLogin(db *sqlx.DB, user User) error {
	query, _ := queries.Queries.ReadFile(queries.DELETE_SESSION_FILE)
	_, err := db.Exec(string(query), user.Id)
	return err
}
