package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/1Mochiyuki/gosky/app"
	"github.com/1Mochiyuki/gosky/config/logger"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Handle  string
	AppPass string
	Id      int
}

var (
	DB *sql.DB
	l  = logger.Get()
)

func InitDB() error {
	appHome, homeErr := app.AppHome()
	if homeErr != nil {
		return homeErr
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/app.db", appHome))
	if err != nil {
		return err
	}
	if connErr := db.Ping(); connErr != nil {
		l.Fatal().Err(connErr).Msg("there was an error pinging the database")
		return connErr
	}
	DB = db
	l.Info().Msg("db connection successful")

	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS users (handle TEXT PRIMARY KEY NOT NULL UNIQUE, app_pass TEXT NOT NULL)")
	if err != nil {
		return err
	}
	return nil
}

func CreateSavedLogin(db *sql.DB, handle, appPass string) error {
	hashedAppPass, err := bcrypt.GenerateFromPassword([]byte(appPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT OR IGNORE INTO users (handle, app_pass) VALUES (?, ?)", handle, hashedAppPass)
	return err
}

func RemoveSavedLogin(db *sql.DB, handle string) error {
	_, err := db.Exec("DELETE FROM users WHERE handle = ?", handle)
	return err
}

func VerifyLogin(db *sql.DB, handle, appPass string) (bool, error) {
	var hashedAppPass string
	err := db.QueryRow("SELECT app_pass FROM users WHERE handle= ?", handle).Scan(&hashedAppPass)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			l.Warn().Err(err).Msg("Couldnt find user")
			return false, err

		}
		l.Warn().Err(err).Msg("unhandled error")

		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedAppPass), []byte(appPass))
	if err != nil {
		return false, err
	}

	return true, nil
}
