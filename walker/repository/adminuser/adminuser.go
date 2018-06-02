package adminuser

import (
	"gopkg.in/mgo.v2"
	"btdxcx.com/walker/model"
)

const (
	dbName = "headquarters"
	adminCollection = "adminuser"
)

// IRepository adminuser repository interface
type IRepository interface {
	Create(*model.AdminUser) error
	GetAll() ([]*model.AdminUser, error)
	Close()
}

// Repository db repository
type Repository struct {
	Session *mgo.Session
}

// Create a new account
func (repo *Repository) Create(m *model.AdminUser) error {
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	c := repo.collection()

	if err := c.EnsureIndex(index); err != nil {
		return err
	}
	return c.Insert(m)
}

// Close session
func (repo *Repository) Close()  {
	repo.Session.Close()
}

func (repo *Repository) collection() *mgo.Collection {
	return repo.Session.DB(dbName).C(adminCollection)
}
