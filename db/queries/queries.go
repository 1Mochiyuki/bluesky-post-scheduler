package queries

import (
	"embed"
)

//go:embed *.sql
var Queries embed.FS

const (
	USER_SCHEMA_FILE    = "user_schema.sql"
	SESSION_SCHEMA_FILE = "session_schema.sql"
	DELETE_USER_FILE    = "delete_user.sql"
	DELETE_SESSION_FILE = "delete_session.sql"
)
