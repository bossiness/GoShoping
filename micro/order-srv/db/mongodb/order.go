package mongodb

import (
	"gopkg.in/mgo.v2/bson"
	proto "btdxcx.com/micro/order-srv/proto/order"
)

const (
	ordersCollectionName = "orders"
)

// Order DB
type Order struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	UUID string        `bson:"uuid,omitempty"`
	// Items              []OrderItem       `bson:"items,omitempty"`
	ItemsTotal         int64              `bson:"itemsTotal,omitempty"`
	Adjustments        []Order_Adjustment `bson:"adjustments,omitempty"`
	AdjustmentsTotal   int64              `bson:"adjustmentsTotal,omitempty"`
	Total              int64              `bson:"total,omitempty"`
	ShippingAddress    string             `bson:"shippingAddress,omitempty"`
	BillingAddress     string             `bson:"billingAddress,omitempty"`
	Shipment           Order_Shipment     `bson:"shipment,omitempty"`
	Payment            Order_Payment      `bson:"payment,omitempty"`
	Customer           string             `bson:"customer,omitempty"`
	State              string             `bson:"state,omitempty"`
	CheckoutState      string             `bson:"checkoutState,omitempty"`
	CheckoutCompleteAt int64              `bson:"checkoutCompleteAt,omitempty"`
	CreatedAt          int64              `bson:"created_at,omitempty"`
	UpdatedAt          int64              `bson:"updated_at,omitempty"`
}

// OrderItem DB
type OrderItem struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	OrderID string        `bson:"order_id,omitempty"`
}

type Address struct {
}

type Customer struct {
}

type Order_Adjustment struct {
}

type Order_Shipment struct {
}

type Order_Payment struct {
}




// CreateOrder create order
func (m *Mongo) CreateOrder(dbname string, customer string) error {
	return nil
}

// ReadOrders read orders
func (m *Mongo) ReadOrders(dbname string, offset int, limst int) (*[]*proto.OrderRecord, error) {
	return nil, nil
}

// ReadOrder read order
func (m *Mongo) ReadOrder(dbname string, uuid string) (*proto.OrderRecord, error) {
	return nil, nil
}

// DeleteOrder delete order
func (m *Mongo) DeleteOrder(dbname string, uuid string) error {
	return nil
}