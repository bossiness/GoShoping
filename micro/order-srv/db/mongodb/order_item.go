package mongodb

import (
	proto "btdxcx.com/micro/order-srv/proto/order"
)

// CreateOrderItem create order item
func (m *Mongo) CreateOrderItem(dbname string, order string, item *proto.OrderRecord_Item) (string, error) {
	// c := m.session.DB(dbname).C(orderItemsCollectionName)
	return "", nil
}

// UpdateOrderItem update order item
func (m *Mongo) UpdateOrderItem(dbname string, order string, id string, item *proto.OrderRecord_Item) error {
	return nil
}

// DeleteOrderItem delete order item
func (m *Mongo) DeleteOrderItem(dbname string, order string, id string) error {
	return nil
}
