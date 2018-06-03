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

// IntrospectRequest introspect token
type IntrospectRequest struct {
	AccessToken string `json:"access_token,omitempty"`
	ShopID      string `json:"shop_id,omitempty"`
}

// Introspect token
type Introspect struct {
	Username string `json:"username"`
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
