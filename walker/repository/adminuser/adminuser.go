package adminuser

import (
	"time"
	"gopkg.in/mgo.v2/bson"
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
	GetList(int, int) (*[]model.AdminUser, int, error)
	Get(string) (*model.AdminUser, error)
	Update(string, *model.AdminUser) error
	Delete(string) error
	GetProfile(string) (*model.AdminUser, error)
	Close()
}

// Repository db repository
type Repository struct {
	Session *mgo.Session
}

// Create a new adminuser
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
	m.RegisteredAt = time.Now().Unix()
	return c.Insert(m)
}

// GetList get adminuser
func (repo *Repository) GetList(offset int, limit int) (*[]model.AdminUser, int, error) {
	query := repo.collection().Find(nil)
	total, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	result := &[]model.AdminUser{}
	if err = query.Skip(offset).Limit(limit).All(result); err != nil {
		return nil, total, err
	}
	return result, total, nil
}

// Get get adminuser
func (repo *Repository) Get(id string) (*model.AdminUser, error) {
	result := &model.AdminUser{}
	err := repo.collection().FindId(bson.ObjectIdHex(id)).One(result)
	return result, err
}

// Update get adminuser
func (repo *Repository) Update(id string, m *model.AdminUser) error {
	updataData := bson.M{"$set": bson.M{
		"first_name": m.FirstName,
		"last_name": m.LastName,
		"phone": m.Phone,
		"email": m.Email,
		"portrait": m.Portrait,
		"wexin_id": m.WeXinID,
	}}
	return repo.collection().UpdateId(bson.ObjectIdHex(id), updataData)
}

// Delete get adminuser
func (repo *Repository) Delete(id string) error {
	return repo.collection().RemoveId(bson.ObjectIdHex(id))
}

// GetProfile get adminuser
func (repo *Repository) GetProfile(username string) (*model.AdminUser, error) {
	result := &model.AdminUser{}
	err := repo.collection().Find(bson.M{"username": username}).One(result)
	return result, err
}

// Close session
func (repo *Repository) Close()  {
	repo.Session.Close()
}

func (repo *Repository) collection() *mgo.Collection {
	return repo.Session.DB(dbName).C(adminCollection)
}
