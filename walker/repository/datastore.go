package repository

import (
	"gopkg.in/mgo.v2"
	"os"
	"github.com/micro/go-log"
)

var (
	session *mgo.Session
)

const (
	defaultHost = "localhost:27017"
)

// SingleSession single instance
func SingleSession() *mgo.Session {

	if session != nil {
		return session
	}

	// Database host from the environment variables
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := createSession(host)

	// Mgo creates a 'master' session, we need to end that session
	// before the main function closes.
	// defer session.Close()

	if err != nil {
		log.Fatal("create session: ", err)
	}

	return session
}

// CreateSession creates the main session to our mongodb instance
func createSession(dbURL string) (*mgo.Session, error) {
	var err error
	session, err = mgo.Dial(dbURL)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	return session, nil
}