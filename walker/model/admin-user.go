package model

import (
	"gopkg.in/mgo.v2/bson"
)

// AdminUser db
type AdminUser struct {
	ID           bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Username     string        `json:"username" bson:"username,omitempty"`
	FirstName    string        `json:"first_name" bson:"first_name,omitempty"`
	LastName     string        `json:"last_name" bson:"last_name,omitempty"`
	Phone        string        `json:"phone" bson:"phone,omitempty"`
	Email        string        `json:"email" bson:"email,omitempty"`
	Portrait     string        `json:"portrait" bson:"portrait,omitempty"`
	RegisteredAt int64         `json:"registered_at" bson:"registered_at,omitempty"`
	WeXinID      string        `json:"wexin_id" bson:"wexin_id,omitempty"`
	AccessAt     int64         `json:"access_at" bson:"access_at,omitempty"`
	Enable       bool          `json:"enable" bson:"enable,omitempty"`
	Role         []string      `json:"role" bson:"role,omitempty"`
	ShopID       string        `json:"shop_id" bson:"shop_id,omitempty"`
}

// AdminUsersPage adminuser page
type AdminUsersPage struct {
	Offset  int32        `json:"offset"`
	Limit   int32        `json:"limit"`
	Total   int32        `json:"total"`
	Records []*AdminUser `json:"records,omitempty"`
}

// AdminUsersRecord adminuser record
type AdminUsersRecord struct {
	Record *AdminUser `json:"record,omitempty"`
}
