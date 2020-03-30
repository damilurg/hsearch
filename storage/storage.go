package storage

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/comov/hsearch/configs"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

type (
	// Connector - the interface to the storage.
	Connector struct {
		DB              *sql.DB
		skipTime        time.Duration
		freshOffersTime time.Duration
	}
)

// regexContain - sqlite3 does not contain a certain type of error, so we have
//  to search through the text to understand what error was caused.
var regexContain = regexp.MustCompile(`UNIQUE constraint failed*`)

// New - creates a connection to the base and returns the interface to work
//  with the storage.
func New(cnf *configs.Config) (*Connector, error) {
	db, err := sql.Open("sqlite3", "hsearch.db?cache=shared")
	if err != nil {
		return nil, err
	}

	// sqlite3 is a simple file storage that is very fragile for multi-threaded
	//  operation. Therefore, we set limits that do not apply to PG or MySql.
	db.SetMaxOpenConns(1)

	return &Connector{
		DB:              db,
		skipTime:        cnf.SkipTime,
		freshOffersTime: cnf.FreshOffers,
	}, nil
}

// Migrate - Applies the changes recorded in the migration files to the
//  database.
func (c *Connector) Migrate(path string) error {
	err := goose.SetDialect("sqlite3")
	if err != nil {
		return err
	}
	err = goose.Run("up", c.DB, path)
	if err == goose.ErrNoNextVersion {
		return nil
	}

	return err
}
