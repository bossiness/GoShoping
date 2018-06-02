package model

import (
	"gopkg.in/mgo.v2/bson"
)

// AdminUser db
type AdminUser struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Username     string        `bson:"username,omitempty"`
	FirstName    string        `bson:"first_name,omitempty"`
	LastName     string        `bson:"last_name,omitempty"`
	Phone        string        `bson:"phone,omitempty"`
	Email        string        `bson:"email,omitempty"`
	Portrait     string        `bson:"portrait,omitempty"`
	RegisteredAt int64         `bson:"registered_at,omitempty"`
	WeXinID      string        `bson:"wexin_id,omitempty"`
	AccessAt     int64         `bson:"access_at,omitempty"`
	Enable       bool          `bson:"enable,omitempty"`
	Role         []string      `bson:"role,omitempty"`
	ShopID       string        `bson:"shop_id,omitempty"`
}
