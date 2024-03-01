package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/parkpcm/standard/secret"
	log "github.com/sirupsen/logrus"
	"os"
)

var ErrConnectingToDatabase = errors.New("unexpected error connecting to the database")
var ErrMissingConfigurationVariable = errors.New("missing configuration variable")

type secretData struct {
	Host     string `json:"instance"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	Private  string `json:"private"`
}

type Connection struct {
	host       string
	username   string
	password   string
	database   string
	private    string
	credLoaded bool
	CredsPath  string
}

func (c *Connection) fromSecret(s secretData) {
	c.host = s.Host
	c.username = s.Username
	c.password = s.Password
	c.database = s.Database
	c.private = s.Private
}

// Validate checks if the configuration is valid by ensuring that the required fields are not empty.
// It returns an error if any of the required fields is missing.
func (c *Connection) Validate() error {

	if len(c.username) == 0 {
		return fmt.Errorf("%w missing username", ErrMissingConfigurationVariable)
	}

	if len(c.private) == 0 {
		return fmt.Errorf("%w missing password", ErrMissingConfigurationVariable)
	}

	if len(c.database) == 0 {
		return fmt.Errorf("%w missing database", ErrMissingConfigurationVariable)
	}

	return nil
}

// loadCredentials loads the credentials for the connection from either the secret path specified in the environment variable "SECRET_PATH" or from the specified `CredsPath`.
func (c *Connection) loadCredentials(ctx context.Context) error {

	if !c.credLoaded {

		var data []byte
		var err error

		if len(os.Getenv("SECRET_PATH")) > 0 {
			if data, err = secret.Get(ctx, os.Getenv("SECRET_PATH")); err != nil {
				return fmt.Errorf("%w error getting credentials from SECRET_PATH", err)
			}
		} else if len(c.CredsPath) > 0 {
			if data, err = secret.GetFromVolume(c.CredsPath); err != nil {
				return fmt.Errorf("%w error getting credentials from %s", err, c.CredsPath)
			}
		} else {
			return fmt.Errorf("unable to find a credentials source")
		}

		s := secretData{}

		if err = json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("%w error reading data into credentails", err)
		}

		c.credLoaded = true

		c.fromSecret(s)

		if err = c.Validate(); err != nil {
			return err
		}

	}

	return nil

}

// ConnectByIP establishes a database connection using the provided IP address.
// It returns a pointer to sql.DB and any error encountered during the connection process.
func (c *Connection) ConnectByIP(ctx context.Context) (*sql.DB, error) {

	if err := c.loadCredentials(ctx); err != nil {
		return nil, err
	}

	if len(c.private) == 0 {
		return nil, fmt.Errorf("%w missing private IP", ErrMissingConfigurationVariable)
	}

	var link *sql.DB
	var err error

	sqlS := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?autocommit=true&parseTime=true", c.username, c.password, c.private, c.database)

	if link, err = sql.Open("mysql", sqlS); err != nil {
		log.Errorf("ConnectionIP: %v", err)
		return nil, fmt.Errorf("%w [%v]", ErrConnectingToDatabase, err)
	}

	if _, err := link.Exec("SET time_zone = 'Europe/London'"); err != nil {
		log.Errorf("ConnectionIP: %v", err)
		return nil, fmt.Errorf("%w [%v]", ErrConnectingToDatabase, err)
	}

	return link, err

}

// ConnectBySocket establishes a database connection using the provided socket directory and host.
// It returns a pointer to sql.DB and any error encountered during the connection process.
func (c *Connection) ConnectBySocket(ctx context.Context) (*sql.DB, error) {

	if err := c.loadCredentials(ctx); err != nil {
		return nil, err
	}

	if len(c.host) == 0 {
		return nil, fmt.Errorf("%w missing database host", ErrMissingConfigurationVariable)
	}

	var link *sql.DB
	var err error

	var dbURI string
	socketDir := "/cloudsql"

	if len(os.Getenv("SQLPATH")) > 0 {
		socketDir = os.Getenv("SQLPATH")
	}
	dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true&autocommit=true", c.username, c.password, socketDir, c.host, c.database)

	if link, err = sql.Open("mysql", dbURI); err != nil {
		log.Errorf("ConnectionBySocket: %v", err)
		return nil, fmt.Errorf("%w [%v]", ErrConnectingToDatabase, err)
	}

	if _, err := link.Exec("SET time_zone = 'Europe/London'"); err != nil {
		log.Errorf("ConnectionBySocket: %v", err)
		return nil, fmt.Errorf("%w [%v]", ErrConnectingToDatabase, err)
	}

	return link, err

}
