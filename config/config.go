package config

import (
	"database/sql"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache   bool
	InfoLog    *log.Logger
	Production bool
	Session    *scs.SessionManager
	DB         *sql.DB
}
