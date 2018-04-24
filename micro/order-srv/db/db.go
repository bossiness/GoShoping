package db

import (
	proto "btdxcx.com/micro/order-srv/proto/order"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Deinit()
	Order
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

// Order is Order
type Order interface {
	CreateOrder(dbname string, customer string) error
	ReadOrders(dbname string, offset int, limst int) (*[]*proto.OrderRecord, error)
	ReadOrder(dbname string, uuid string) (*proto.OrderRecord, error)
	DeleteOrder(dbname string, uuid string) error
}

// CreateOrder create order
func CreateOrder(dbname string, customer string) error {
	return db.CreateOrder(dbname, customer)
}

// ReadOrders read orders
func ReadOrders(dbname string, offset int, limst int) (*[]*proto.OrderRecord, error) {
	return db.ReadOrders(dbname, offset, limst)
}

// ReadOrder read order
func ReadOrder(dbname string, uuid string) (*proto.OrderRecord, error) {
	return db.ReadOrder(dbname, uuid)
}

// DeleteOrder delete order
func DeleteOrder(dbname string, uuid string) error {
	return db.DeleteOrder(dbname, uuid)
}
