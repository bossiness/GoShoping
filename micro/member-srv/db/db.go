package db

import (
	proto "btdxcx.com/micro/member-srv/proto/member"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Deinit()
	Customer
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

// Deinit 析构
func Deinit() {
	db.Deinit()
}

// Customer is Customer
type Customer interface {
	CreateCustomer(string, *proto.Customer) error
	ReadCustomers(string, int, int) (*[]*proto.Customer, error)
	ReadCustomer(string, string) (*proto.Customer ,error)
	UpdateCustomer(string, string, *proto.Customer) error
	DeleteCustomer(string, string) error
}
