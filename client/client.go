package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/1Mochiyuki/gosky/errs"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
)

type BskyAgent struct {
	ctx     context.Context
	client  *xrpc.Client
	handle  string
	apiPass string
}

func NewAgent(ctx context.Context, server, handle, apiPass string) BskyAgent {
	if server == "" {
		server = "https://bsky.social"
	}
	if !strings.Contains(handle, ".") {
		handle = fmt.Sprintf("%s.bsky.social", handle)
	}
	return BskyAgent{
		ctx: ctx,
		client: &xrpc.Client{
			Client: new(http.Client),
			Host:   server,
		},
		handle:  handle,
		apiPass: apiPass,
	}
}

var ErrConnection = errors.New("unable to connect")

func (c *BskyAgent) Connect() error {
	input := &atproto.ServerCreateSession_Input{Identifier: c.handle, Password: c.apiPass}
	session, err := atproto.ServerCreateSession(c.ctx, c.client, input)
	if err != nil {
		return errs.NewIncorrectCredentialsError()
	}

	c.client.Auth = &xrpc.AuthInfo{
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
		Handle:     session.Handle,
		Did:        session.Did,
	}

	return nil
}

func (c *BskyAgent) CreatePost(post bsky.FeedPost) (string, string, error) {
	input := atproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Record:     &util.LexiconTypeDecoder{Val: &post},
		Repo:       c.handle,
	}

	response, err := atproto.RepoCreateRecord(c.ctx, c.client, &input)
	if err != nil {
		return "", "", fmt.Errorf("unable to post: %v", err)
	}
	return response.Cid, response.Uri, nil
}
