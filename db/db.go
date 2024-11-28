package db

import (
	"fmt"

	"github.com/1Mochiyuki/gosky/app"
	"github.com/1Mochiyuki/gosky/config/logger"
	"github.com/1Mochiyuki/gosky/db/queries"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB *sqlx.DB
	l  = logger.Get()
)

var (
	SESSION_CACHE = make(map[string]Session)
	USER_CACHE    = make(map[uint8]string)
)

func SessionFromAuthInfo(info *xrpc.AuthInfo) Session {
	return Session{
		AccessJWT:  info.AccessJwt,
		RefreshJWT: info.RefreshJwt,
		Handle:     info.Handle,
		Did:        info.Did,
	}
}

func InitDB() error {
	appHome, homeErr := app.AppHome()
	if homeErr != nil {
		return homeErr
	}

	db, err := sqlx.Open("sqlite3", fmt.Sprintf("%s/app.db", appHome))
	if err != nil {
		return err
	}
	DB = db
	db.MustExec("PRAGMA foreign_keys = ON;")
	user_schema, readErr := queries.Queries.ReadFile(queries.USER_SCHEMA_FILE)
	if readErr != nil {
		return readErr
	}
	db.MustExec(string(user_schema))
	session_schema, readErr := queries.Queries.ReadFile(queries.SESSION_SCHEMA_FILE)
	if readErr != nil {
		return readErr
	}
	db.MustExec(string(session_schema))
	l.Info().Msg("db connection successful")

	return nil
}
