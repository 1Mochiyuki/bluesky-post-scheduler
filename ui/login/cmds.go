package login

import (
	"github.com/1Mochiyuki/gosky/config/logger"
	"github.com/1Mochiyuki/gosky/db"
	tea "github.com/charmbracelet/bubbletea"
)

type SuccessPing bool

func PingDb() tea.Msg {
	connErr := db.DB.Ping()
	return SuccessPing(connErr != nil)
}

var l = logger.Get()

type SaveCredentialsCmd struct {
	err     error
	handle  string
	appPass string
}

type AnyUserExistsMsg struct {
	Err     error
	Results []db.User
}

func AnyCredentialsExist() tea.Msg {
	var rowCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&rowCount)
	if err != nil {
		l.Error().Err(err).Msg("error checking any creds exist")
		return AnyUserExistsMsg{
			Err: err,
		}
	}

	var results []db.User
	var handle string
	var userId uint8

	if rowCount == 0 {
		l.Warn().Msg("db is empty")
		return AnyUserExistsMsg{
			Err: err,
		}
	}
	if rowCount == len(db.USER_CACHE) {
		l.Debug().Msg("pulling from cache")
		for userId, handle := range db.USER_CACHE {
			results = append(results, db.User{Id: userId, Handle: handle})
		}
		return AnyUserExistsMsg{
			Results: results,
			Err:     nil,
		}
	}
	rows, err := db.DB.Query("SELECT * FROM users")
	if err != nil {
		l.Error().Err(err).Msg("error selecting all users from table")
		return AnyUserExistsMsg{
			Err: err,
		}
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&userId, &handle)
		if err != nil {
			l.Fatal().Err(err).Msg("error retrieving details from db")
			return AnyUserExistsMsg{
				Err: err,
			}
		}
		db.USER_CACHE[userId] = handle
		results = append(results, db.User{Id: userId, Handle: handle})
	}
	if err = rows.Err(); err != nil {
		l.Fatal().Err(err).Msg("an error occurred getting next row")
	}
	l.Debug().Msg("pulled from db, added to cache")
	return AnyUserExistsMsg{
		Results: results,
	}
}
