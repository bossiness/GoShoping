package db

import (
	proto "btdxcx.com/micro/shop-srv/proto/shop"
	dproto "btdxcx.com/micro/shop-srv/proto/shop/details"
)

// DB is 数据库接口
type DB interface {
	Init() error
	ShopKey
	ShopDetails
}

// ShopKey is ShopKey数据接口
type ShopKey interface {
	ReadKey(string) (*proto.ShopTagKeys, error)
	CreateKey(string, *proto.ShopTagKeys) error
	DeleteKey(string) error
}

// ShopDetails is ShopDetails数据接口
type ShopDetails interface {
	CreateDetails(*dproto.CreateRequest) (*dproto.CreateResponse, error)
	ReadDetails(string) (*dproto.ReadResponse, error)
	DeleteDetails(string) error
	UpdateDetails(*dproto.UpdateRequest) error
	ListDetails(*dproto.ListRequest) (*dproto.ListResponse, error)
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


// CreateDetails create 
func CreateDetails(req *dproto.CreateRequest) (*dproto.CreateResponse, error) {
	return db.CreateDetails(req)
}

// ReadDetails read 
func ReadDetails(id string) (*dproto.ReadResponse, error) {
	return db.ReadDetails(id)
}

// DeleteDetails delete 
func DeleteDetails(id string) error {
	return db.DeleteDetails(id)
}

// UpdateDetails update 
func UpdateDetails(req *dproto.UpdateRequest) error {
	return db.UpdateDetails(req)
}

// ListDetails list
func ListDetails(req *dproto.ListRequest) (*dproto.ListResponse, error) {
	return db.ListDetails(req)
}