package login

import (
	"fmt"

	"github.com/1Mochiyuki/gosky/config/logger"
	"github.com/1Mochiyuki/gosky/db"
	"github.com/bluesky-social/indigo/xrpc"
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

// func SaveCredentials(handle, appPass string) tea.Cmd {
// 	saveErr := db.CreateSavedLogin(db.DB, handle, appPass, nil)
// 	if saveErr != nil {
// 		l.Error().Err(saveErr).Msg("there was an error saving the credentials")
// 		return func() tea.Msg {
// 			return SaveCredentialsCmd{
// 				handle:  handle,
// 				appPass: appPass,
// 				err:     saveErr,
// 			}
// 		}
// 	}
// 	l.Info().Msgf("saved %s to db", handle)
//
// 	return func() tea.Msg {
// 		return SaveCredentialsCmd{
// 			handle:  handle,
// 			appPass: appPass,
// 			err:     nil,
// 		}
// 	}
// }

type Credentials struct {
	handle  string
	appPass string
	auth    *xrpc.AuthInfo
}

func (c Credentials) String() string {
	return fmt.Sprintf("Handle: %s Pass: %s", c.handle, c.appPass)
}

type AnyUserExistsMsg struct {
	err     error
	results []Credentials
}

func AnyCredentialsExist() tea.Msg {
	var rowCount int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&rowCount)
	if err != nil {
		l.Error().Err(err).Msg("error checking any creds exist")
		return AnyUserExistsMsg{
			err: err,
		}
	}
	if rowCount == 0 {
		l.Warn().Msg("db is empty")
		return AnyUserExistsMsg{
			err: err,
		}
	}
	rows, err := db.DB.Query("SELECT * FROM users")
	if err != nil {
		l.Error().Err(err).Msg("error selecting all users from table")
		return AnyUserExistsMsg{
			err: err,
		}
	}
	defer rows.Close()
	var results []Credentials
	var handle, pass string
	for rows.Next() {
		err = rows.Scan(&handle, &pass)
		if err != nil {
			l.Fatal().Err(err).Msg("error retrieving details from db")
			return AnyUserExistsMsg{
				err: err,
			}
		}
		results = append(results, Credentials{handle: handle, appPass: pass})
	}
	if err = rows.Err(); err != nil {
		l.Fatal().Err(err).Msg("an error occurred getting next row")
	}
	return AnyUserExistsMsg{
		results: results,
	}
}
