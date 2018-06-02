package token

import (
	"btdxcx.com/walker/model"
	"gopkg.in/mgo.v2"
)

const (
	dbName          = "headquarters"
	tokenCollection = "token"
)

// IRepository token repository interface
type IRepository interface {
	Create(*model.Jwtauth) error
	GetAll() ([]*model.Jwtauth, error)
	Close()
}

// Repository db repository
type Repository struct {
	Session *mgo.Session
}

// Create a new account
func (repo *Repository) Create(account *model.Jwtauth) error {
	index := mgo.Index{
		Key:        []string{"client_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	c := repo.collection()

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	return c.Insert(account)
}

// GetAll consignments
func (repo *Repository) GetAll() ([]*model.Jwtauth, error) {
	var jwtauths []*model.Jwtauth
	err := repo.collection().Find(nil).All(&jwtauths)
	return jwtauths, err
}

// Close session
func (repo *Repository) Close() {
	repo.Session.Close()
}

func (repo *Repository) collection() *mgo.Collection {
	return repo.Session.DB(dbName).C(tokenCollection)
}
