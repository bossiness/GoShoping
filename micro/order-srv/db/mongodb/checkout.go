package mongodb

import (
	proto "btdxcx.com/micro/order-srv/proto/order"
)

// CheckoutNew new chekout
func (m *Mongo) CheckoutNew(dbname string, orderID string) error {
	return nil
}

// CheckoutAddressing chekout addressing
func (m *Mongo) CheckoutAddressing(dbname string, orderID string, shipping *proto.Address, billing *proto.Address) error {
	return nil
}

// CheckoutSelectShipping chekout shipping
func (m *Mongo) CheckoutSelectShipping(dbname string, orderID string, method string) error {
	return nil
}

// CheckoutSelectPayment chekout payment
func (m *Mongo) CheckoutSelectPayment(dbname string, orderID string, method string) error {
	return nil
}

// CheckoutComplete chekout complete
func (m *Mongo) CheckoutComplete(dbname string, orderID string) error {
	return nil
}
