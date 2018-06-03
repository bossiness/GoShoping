package account

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"btdxcx.com/walker/model"
)

const (
	dbName = "headquarters"
	accountCollection = "account"
)

// IRepository auth repository interface
type IRepository interface {
	Create(*model.Account) error
	Get(string) (*model.Account, error)
	Close()
}

// Repository db repository
type Repository struct {
	Session *mgo.Session
}

// Create a new account
func (repo *Repository) Create(account *model.Account) error {
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

// Get consignments
func (repo *Repository) Get(clientID string) (*model.Account, error) {
	account := new(model.Account)
	query := bson.M{
		"client_id": clientID,
	}
	err := repo.collection().Find(query).One(account)
	return account, err
}

// Close session
func (repo *Repository) Close()  {
	repo.Session.Close()
}

func (repo *Repository) collection() *mgo.Collection {
	return repo.Session.DB(dbName).C(accountCollection)
}