package db

import (
	proto "btdxcx.com/micro/jwtauth-srv/proto/auth"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Auth
}

// Auth is Auth数据接口
type Auth interface {
	Read(string, string) (*proto.Record, error)
	Create(string, *proto.Record) error
	Update(string, *proto.Record) error
	DeleteAccessToken(string, string) error
	DeleteRefreshToken(string, string) error
	ReadFormRefreshToken(string, string) (*proto.Record, error)
}

var (
	db DB
)

// Register db Imp
func Register(backend DB) {
	db = backend
}

// Init 数据库初始化
func Init() error {
	return db.Init()
}

// Read Account
func Read(dbName string, token string) (*proto.Record, error) {
	return db.Read(dbName, token)
}

// Create Account
func Create(dbName string, record *proto.Record) error {
	return db.Create(dbName, record)
}

// Update Account
func Update(dbName string, record *proto.Record) error {
	return db.Update(dbName, record)
}

// DeleteAccessToken Account
func DeleteAccessToken(dbName string, token string) error {
	return db.DeleteAccessToken(dbName, token)
}

// DeleteRefreshToken Account
func DeleteRefreshToken(dbName string, token string) error {
	return db.DeleteRefreshToken(dbName, token)
}

// ReadFormRefreshToken read token
func ReadFormRefreshToken(dbName string, refreshToken string) (*proto.Record, error) {
	return db.ReadFormRefreshToken(dbName, refreshToken)
}