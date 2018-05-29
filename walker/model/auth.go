package model

import (
	"gopkg.in/mgo.v2/bson"
)

// AuthRequest Auth Request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code,omitempty"`
	ShopID   string
	Scopes   []string
	Metadata map[string]string
	Type     string
}

// Token Auth Response
type Token struct {
	// jwt
	AccessToken string `json:"access_token"`
	// jwt
	RefreshToken string            `json:"refresh_token"`
	ExpiresAt    int64             `json:"expires_at"`
	Scopes       []string          `json:"scopes,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// Account db
type Account struct {
	ID           bson.ObjectId     `bson:"_id,omitempty"`
	Type         string            `bson:"type"`
	ClientID     string            `bson:"client_id"`
	ClientSecret string            `bson:"client_secret"`
	Metadata     map[string]string `bson:"metadata"`
	CreatedAt    int64             `bson:"created_at"`
	UpdatedAt    int64             `bson:"updated_at"`
}

// Jwtauth db
type Jwtauth struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	ClientID  string        `bson:"client_id"`
	Scopes    []string      `bson:"scopes"`
	Access    string        `bson:"access_token"`
	Refresh   string        `bson:"refresh_token"`
	ExpiresAt int64         `bson:"expires_at"`
	Cipher    string        `bson:"cipher"`
}

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
