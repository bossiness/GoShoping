package token

import (
	"gopkg.in/mgo.v2/bson"
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
	Delete(string) error
	GetAll() ([]*model.Jwtauth, error)
	Close()
}

// Repository db repository
type Repository struct {
	Session *mgo.Session
}

// Create a new account
func (repo *Repository) Create(jwt *model.Jwtauth) error {
	c := repo.collection()
	c.Remove(bson.M{"client_id": jwt.ClientID})
	return c.Insert(jwt)
}

// Delete a new account
func (repo *Repository) Delete(clientID string) error {
	c := repo.collection()
	return c.Remove(bson.M{"client_id": clientID})
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
