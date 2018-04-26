package mongodb

import (
	proto "btdxcx.com/micro/order-srv/proto/order"
	"gopkg.in/mgo.v2/bson"
)

const (
	ordersCollectionName = "orders"
	customersCollectionName = "customers"
	orderItemsCollectionName = "orderItems"
)

// Order DB
type Order struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	UUID string        `bson:"uuid,omitempty"`
	// Items              []OrderItem       `bson:"items,omitempty"`
	ItemsTotal         int64         `bson:"itemsTotal,omitempty"`
	Adjustments        []string      `bson:"adjustments,omitempty"`
	AdjustmentsTotal   int64         `bson:"adjustmentsTotal,omitempty"`
	Total              int64         `bson:"total,omitempty"`
	ShippingAddress    string        `bson:"shippingAddress,omitempty"`
	BillingAddress     string        `bson:"billingAddress,omitempty"`
	Shipment           OrderShipment `bson:"shipment,omitempty"`
	Payment            OrderPayment  `bson:"payment,omitempty"`
	Customer           string        `bson:"customer,omitempty"`
	State              string        `bson:"state,omitempty"`
	CheckoutState      string        `bson:"checkoutState,omitempty"`
	CheckoutCompleteAt int64         `bson:"checkoutCompleteAt,omitempty"`
	CreatedAt          int64         `bson:"created_at,omitempty"`
	UpdatedAt          int64         `bson:"updated_at,omitempty"`
}

// OrderItem DB
type OrderItem struct {
	ID               bson.ObjectId `bson:"_id,omitempty"`
	OrderID          string        `bson:"order_id,omitempty"`
	Quantity         int64         `bson:"quantity,omitempty"`
	UnitPrice        int64         `bson:"unitPrice,omitempty"`
	Total            int64         `bson:"total,omitempty"`
	Adjustments      []string      `bson:"adjustments,omitempty"`
	AdjustmentsTotal int64         `bson:"adjustmentsTotal,omitempty"`
	Variant          string        `bson:"variant,omitempty"`
	CreatedAt        int64         `bson:"created_at,omitempty"`
	UpdatedAt        int64         `bson:"updated_at,omitempty"`
}

// Address DB
type Address struct {
	ID        string `bson:"_id,omitempty"`
	FirstName string `bson:"firstName,omitempty"`
	LastName  string `bson:"lastName,omitempty"`
	City      string `bson:"city,omitempty"`
	Postcode  string `bson:"postcode,omitempty"`
	Street    string `bson:"street,omitempty"`
	Country   string `bson:"country,omitempty"`
}

// Customer DB
type Customer struct {
	ID        string `bson:"_id,omitempty"`
	Username  string `bson:"username,omitempty"`
	FirstName string `bson:"firstName,omitempty"`
	LastName  string `bson:"lastName,omitempty"`
	Phone     string `bson:"phone,omitempty"`
	Email     string `bson:"email,omitempty"`
	Portrait  string `bson:"portrait,omitempty"`
	Role      string `bson:"role,omitempty"`
}

// OrderAdjustment DB
type OrderAdjustment struct {
	ID     string `bson:"_id,omitempty"`
	Type   string `bson:"type,omitempty"`
	Label  string `bson:"label,omitempty"`
	Amount int64  `bson:"amount,omitempty"`
}

// OrderShipment DB
type OrderShipment struct {
	State  string `bson:"state,omitempty"`
	Method string `bson:"method,omitempty"`
}

// OrderPayment DB
type OrderPayment struct {
	State  string `bson:"state,omitempty"`
	Amount int64  `bson:"amount,omitempty"`
	Method string `bson:"method,omitempty"`
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
