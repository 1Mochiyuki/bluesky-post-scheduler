package db

import (
	"time"

	"github.com/1Mochiyuki/gosky/db/queries"
	"github.com/jmoiron/sqlx"
)

func InsertNewSession(db *sqlx.DB, session Session, userId uint8) error {
	_, err := db.Exec(
		"INSERT OR IGNORE INTO sessions (access_jwt, refresh_jwt, session_user_handle, did, user_id, last_updated) VALUES ($1, $2, $3, $4, $5, $6)",
		session.AccessJWT, session.RefreshJWT, session.Handle, session.Did, userId, time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSession(db *sqlx.DB, session Session, userId uint8) error {
	return nil
}

func RemoveSavedSession(db *sqlx.DB, session Session) error {
	query, _ := queries.Queries.ReadFile(queries.DELETE_SESSION_FILE)
	_, err := db.Exec(string(query), session.Id)
	return err
}
