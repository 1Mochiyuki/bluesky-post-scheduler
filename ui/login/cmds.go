package login

import (
	"fmt"

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

type Credentials struct {
	Handle  string
	AppPass string
}

func (c Credentials) String() string {
	return fmt.Sprintf("Handle: %s Pass: %s", c.Handle, c.AppPass)
}

type AnyUserExistsMsg struct {
	Err     error
	Results []Credentials
}

var LOGIN_CACHCE = map[string]string{}

func AnyCredentialsExist() tea.Msg {
	var rowCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&rowCount)
	if err != nil {
		l.Error().Err(err).Msg("error checking any creds exist")
		return AnyUserExistsMsg{
			Err: err,
		}
	}

	var results []Credentials
	var handle, pass string

	if rowCount == 0 {
		l.Warn().Msg("db is empty")
		return AnyUserExistsMsg{
			Err: err,
		}
	}
	if rowCount == len(LOGIN_CACHCE) {
		l.Debug().Msg("pulling from cache")
		for handle, pass := range LOGIN_CACHCE {
			results = append(results, Credentials{Handle: handle, AppPass: pass})
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
		err = rows.Scan(&handle, &pass)
		if err != nil {
			l.Fatal().Err(err).Msg("error retrieving details from db")
			return AnyUserExistsMsg{
				Err: err,
			}
		}
		LOGIN_CACHCE[handle] = pass
		results = append(results, Credentials{Handle: handle, AppPass: pass})
	}
	if err = rows.Err(); err != nil {
		l.Fatal().Err(err).Msg("an error occurred getting next row")
	}
	l.Debug().Msg("pulled from db, added to cache")
	return AnyUserExistsMsg{
		Results: results,
	}
}
