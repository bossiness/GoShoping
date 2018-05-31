package account

import (
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
	GetAll() ([]*model.Account, error)
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

	if err := repo.collection().EnsureIndex(index); err != nil {
		return err
	}
	return repo.collection().Insert(account)
}

// GetAll consignments
func (repo *Repository) GetAll() ([]*model.Account, error) {
	var accounts []*model.Account
	err := repo.collection().Find(nil).All(&accounts)
	return accounts, err
}

// Close session
func (repo *Repository) Close()  {
	repo.Session.Close()
}

func (repo *Repository) collection() *mgo.Collection {
	return repo.Session.DB(dbName).C(accountCollection)
}