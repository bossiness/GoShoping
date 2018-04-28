package db

import (
	proto "btdxcx.com/micro/member-srv/proto/member"
)

// DB is 数据库接口
type DB interface {
	Init() error
	Deinit()
	Customer
	AdminUser
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

// Customer is Customer db interface
type Customer interface {
	CreateCustomer(string, *proto.CustomerRecord) error
	ReadCustomers(string, int, int) (*[]*proto.CustomerRecord, error)
	ReadCustomer(string, string) (*proto.CustomerRecord ,error)
	ReadCustomerFromName(string, string) (*proto.CustomerRecord ,error)
	UpdateCustomer(string, string, *proto.CustomerRecord) error
	DeleteCustomer(string, string) error
}

// CreateCustomer create Customer
func CreateCustomer(dbname string, record *proto.CustomerRecord) error {
	return db.CreateCustomer(dbname, record)
}

// ReadCustomers read Customers
func ReadCustomers(dbname string, offset int, limit int) (*[]*proto.CustomerRecord, error) {
	return db.ReadCustomers(dbname, offset, limit)
}

// ReadCustomer read a Customer
func ReadCustomer(dbname string, id string) (*proto.CustomerRecord ,error) {
	return db.ReadCustomer(dbname, id)
}

// ReadCustomerFromName read a Customer
func ReadCustomerFromName(dbname string, name string) (*proto.CustomerRecord ,error) {
	return db.ReadCustomerFromName(dbname, name)
}

// UpdateCustomer update a Customer
func UpdateCustomer(dbname string, id string, record *proto.CustomerRecord) error {
	return db.UpdateCustomer(dbname, id, record)
}

// DeleteCustomer delete a Customer
func DeleteCustomer(dbname string, id string) error {
	return db.DeleteCustomer(dbname, id)
}

// AdminUser is AdminUser db interface
type AdminUser interface {
	CreateAdminUser(string, *proto.AdminUserRecord) error
	ReadAdminUsers(string, int, int) (*[]*proto.AdminUserRecord, error)
	ReadAdminUser(string, string) (*proto.AdminUserRecord ,error)
	ReadAdminUserFromName(string, string) (*proto.AdminUserRecord ,error)
	UpdateAdminUser(string, string, *proto.AdminUserRecord) error
	DeleteAdminUser(string, string) error
}

// CreateAdminUser create AdminUser
func CreateAdminUser(dbname string, record *proto.AdminUserRecord) error {
	return db.CreateAdminUser(dbname, record)
}

// ReadAdminUsers read AdminUsers
func ReadAdminUsers(dbname string, offset int, limit int) (*[]*proto.AdminUserRecord, error) {
	return db.ReadAdminUsers(dbname, offset, limit)
}

// ReadAdminUser read a AdminUser
func ReadAdminUser(dbname string, id string) (*proto.AdminUserRecord ,error) {
	return db.ReadAdminUser(dbname, id)
}

// ReadAdminUserFromName read a AdminUser
func ReadAdminUserFromName(dbname string, name string) (*proto.AdminUserRecord ,error) {
	return db.ReadAdminUserFromName(dbname, name)
}

// UpdateAdminUser update a AdminUser
func UpdateAdminUser(dbname string, id string, record *proto.AdminUserRecord) error {
	return db.UpdateAdminUser(dbname, id, record)
}

// DeleteAdminUser delete a AdminUser
func DeleteAdminUser(dbname string, id string) error {
	return db.DeleteAdminUser(dbname, id)
}