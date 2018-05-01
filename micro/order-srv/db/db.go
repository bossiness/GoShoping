package db

import (
	proto "btdxcx.com/micro/order-srv/proto/order"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Deinit()
	Order
	OrderItem
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
	CreateOrder(dbname string, customer string) (string, error)
	ReadOrders(dbname string, state string, checkoutState string, offset int, limit int) (*[]*proto.OrderRecord, error)
	ReadOrder(dbname string, uuid string) (*proto.OrderRecord, error)
	DeleteOrder(dbname string, uuid string) error
	ReadCustomerOrders(dbname string, customer string) (*[]*proto.OrderRecord, error)
}

// CreateOrder create order
func CreateOrder(dbname string, customer string) (string, error) {
	return db.CreateOrder(dbname, customer)
}

// ReadOrders read orders
func ReadOrders(dbname string, state string, checkoutState string, offset int, limit int) (*[]*proto.OrderRecord, error) {
	return db.ReadOrders(dbname, state, checkoutState, offset, limit)
}

// ReadOrder read order
func ReadOrder(dbname string, uuid string) (*proto.OrderRecord, error) {
	return db.ReadOrder(dbname, uuid)
}

// DeleteOrder delete order
func DeleteOrder(dbname string, uuid string) error {
	return db.DeleteOrder(dbname, uuid)
}

// ReadCustomerOrders delete order
func ReadCustomerOrders(dbname string, customer string) (*[]*proto.OrderRecord, error) {
	return db.ReadCustomerOrders(dbname, customer)
}

// OrderItem is Order
type OrderItem interface {
	CreateOrderItem(dbname string, order string, item *proto.OrderRecord_Item) (string, error)
	UpdateOrderItem(dbname string, id string, item *proto.OrderRecord_Item) error
	DeleteOrderItem(dbname string, id string) error
}

// CreateOrderItem create order item
func CreateOrderItem(dbname string, order string, item *proto.OrderRecord_Item) (string, error) {
	return db.CreateOrderItem(dbname, order, item)
}

// UpdateOrderItem update order item
func UpdateOrderItem(dbname string, id string, item *proto.OrderRecord_Item) error {
	return db.UpdateOrderItem(dbname, id, item)
}

// DeleteOrderItem delete order item
func DeleteOrderItem(dbname string, id string) error {
	return db.DeleteOrderItem(dbname, id)
}
