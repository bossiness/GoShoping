package repository

import (
	"gopkg.in/mgo.v2"
)

// CreateSession creates the main session to our mongodb instance
func CreateSession(dbURL string) (*mgo.Session, error) {
	session, err := mgo.Dial(dbURL)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	return session, nil
}