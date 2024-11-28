package db

import (
	"context"
	"net/http"
	"time"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
)

type User struct {
	Handle string `db:"user_handle"`
	Session
	Id uint8 `db:"user_id"`
}

type Session struct {
	TimeUpdated time.Time `db:"last_updated"`
	AccessJWT   string    `db:"access_jwt"`
	RefreshJWT  string    `db:"refresh_jwt"`
	Handle      string    `db:"session_user_handle"`
	Did         string    `db:"did"`
	Id          int       `db:"session_id"`
	UserId      uint8     `db:"user_id"`
}

func (s Session) ToXRPCAuthInfo() *xrpc.AuthInfo {
	return &xrpc.AuthInfo{
		RefreshJwt: s.RefreshJWT,
		AccessJwt:  s.AccessJWT,
		Did:        s.Did,
		Handle:     s.Handle,
	}
}

func (s *Session) Renew(server string) error {
	if server == "" {
		server = "https://bsky.social"
	}
	client := &xrpc.Client{
		Client: new(http.Client),
		Auth:   s.ToXRPCAuthInfo(),
		Host:   server,
	}
	output, refreshErr := atproto.ServerRefreshSession(context.Background(), client)
	if refreshErr != nil {
		l.Debug().Err(refreshErr).Msg("an error occured refreshing the previous session")
		return refreshErr
	}
	s.AccessJWT = output.AccessJwt
	s.RefreshJWT = output.RefreshJwt

	updateDbErr := UpdateSession(DB, *s, s.UserId)
	if updateDbErr != nil {
		l.Debug().Err(updateDbErr).Msg("an error occurred updating the db with a new session")
		return updateDbErr
	}

	return nil
}
