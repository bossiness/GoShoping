package db

import (
	proto "btdxcx.com/micro/shop-srv/proto/shop"
)

// DB is 数据库接口
type DB interface {
	Init() error
	ShopKey
}

// ShopKey is ShopKey数据接口
type ShopKey interface {
	ReadKey(string) (*proto.ShopTagKeys, error)
	CreateKey(string, *proto.ShopTagKeys) error
	DeleteKey(string) error
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

// ReadKey form uuid
func ReadKey(id string) (*proto.ShopTagKeys, error) {
	return db.ReadKey(id)
}

// CreateKey form uuid
func CreateKey(id string, shopKey *proto.ShopTagKeys) error {
	return db.CreateKey(id, shopKey)
}

// DeleteKey form uuid
func DeleteKey(id string) error {
	return db.DeleteKey(id)
}
