package mongodb

import (
	"errors"
	"time"

	proto "btdxcx.com/micro/order-srv/proto/order"
	"gopkg.in/mgo.v2/bson"

	customerDB "btdxcx.com/micro/member-srv/db/mongodb"
)

const (
	ordersCollectionName     = "orders"
	customersCollectionName  = "customers"
)

// Order DB
type Order struct {
	ID bson.ObjectId `bson:"_id,omitempty"`
	// UUID string        `bson:"uuid,omitempty"`
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
func (m *Mongo) CreateOrder(dbname string, customer string) (string, error) {
	c := m.session.DB(dbname).C(ordersCollectionName)
	cc := m.session.DB(dbname).C(customersCollectionName)

	resultCustomer := customerDB.Customer{}
	if err := cc.Find(bson.M{"username": customer}).One(&resultCustomer); err != nil {
		if err.Error() == "not found" {
			return "", errors.New("not found customer")
		}
		return "", err
	}

	doc := &Order{
		ID:               bson.NewObjectId(),
		Customer:         customer,
		Total:            0,
		AdjustmentsTotal: 0,
		ItemsTotal:       0,
		State:            "cart",
		CheckoutState:    "cart",
		CreatedAt:        time.Now().Unix(),
	}
	if err := c.Insert(doc); err != nil {
		return "", err
	}

	return doc.ID.Hex(), nil
}

// ReadOrders read orders
func (m *Mongo) ReadOrders(dbname string, state string, checkoutState string, offset int, limit int) (*[]*proto.OrderRecord, error) {
	c := m.session.DB(dbname).C(ordersCollectionName)

	results := []Order{}
	query := bson.M{}
	if len(state) != 0 && len(checkoutState) != 0 {
		query = bson.M{"state": state, "checkoutState": checkoutState}
	} else if len(state) != 0 {
		query = bson.M{"state": state}
	} else if len(checkoutState) != 0 {
		query = bson.M{"checkoutState": checkoutState}
	} else {
		query = nil
	}
	if err := c.Find(query).Skip(offset).Limit(limit).All(&results); err != nil {
		return nil, err
	}

	records := []*proto.OrderRecord{}
	for _, order := range results {
		records = append(records, m.mapOrder(dbname, &order))
	}

	return nil, nil
}

func (m *Mongo) mapOrder(dbname string, order *Order) *proto.OrderRecord {
	adjustmentsTotal, adjustments := m.readAdjustments(dbname, order.Adjustments)
	return &proto.OrderRecord{
		Uuid:               order.ID.Hex(),
		Items:              m.readOrderItems(dbname, order.ID),
		ItemsTotal:         order.ItemsTotal,
		Adjustments:        adjustments,
		AdjustmentsTotal:   adjustmentsTotal,
		Total:              order.Total,
		ShippingAddress:    m.readShippingAddress(dbname, order.ShippingAddress),
		BillingAddress:     m.readBillingAddress(dbname, order.BillingAddress),
		Shipment:           m.mapShipment(&order.Shipment),
		Payment:            m.mapPayment(&order.Payment),
		Customer:           m.readCustomer(dbname, order.Customer),
		State:              order.State,
		CheckoutState:      order.CheckoutState,
		CheckoutCompleteAt: order.CheckoutCompleteAt,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
	}
}

func (m *Mongo) readCustomer(dbname string, username string) *proto.Customer {
	cc := m.session.DB(dbname).C(customersCollectionName)

	rc := customerDB.Customer{}
	if err := cc.Find(bson.M{"username": username}).One(&rc); err != nil {
		return nil
	}

	return &proto.Customer{
		Username:  rc.Username,
		FirstName: rc.FirstName,
		LastName:  rc.LastName,
		Phone:     rc.Phone,
		Email:     rc.Email,
		Portrait:  rc.Portrait,
	}
}

func (m *Mongo) readOrderItems(dbname string, orderID bson.ObjectId) []*proto.OrderRecord_Item {
	c := m.session.DB(dbname).C(orderItemsCollectionName)

	items := []OrderItem{}
	recordItems := []*proto.OrderRecord_Item{}
	if err := c.Find(bson.M{"order_id": orderID}).All(&items); err != nil {
		return recordItems
	}
	for _, item := range items {
		adjustmentsTotal, adjustments := m.readAdjustments(dbname, item.Adjustments)
		total := (item.Quantity * item.UnitPrice) + adjustmentsTotal
		oi := proto.OrderRecord_Item{
			Uuid:             item.OrderID.Hex(),
			Quantity:         item.Quantity,
			UnitPrice:        item.UnitPrice,
			Total:            total,
			Adjustments:      adjustments,
			AdjustmentsTotal: adjustmentsTotal,
			Variant:          item.Variant,
		}
		recordItems = append(recordItems, &oi)
	}
	return recordItems
}

func (m *Mongo) readAdjustments(dbname string, adjustments []string) (int64, []*proto.OrderRecord_Adjustment) {
	return 0, []*proto.OrderRecord_Adjustment{}
}

func (m *Mongo) readShippingAddress(dbname string, address string) *proto.Address {
	return nil
}

func (m *Mongo) readBillingAddress(dbname string, address string) *proto.Address {
	return nil
}

func (m *Mongo) mapShipment(os *OrderShipment) *proto.OrderRecord_Shipment {
	return nil
}

func (m *Mongo) mapPayment(os *OrderPayment) *proto.OrderRecord_Payment {
	return nil
}

// ReadOrder read order
func (m *Mongo) ReadOrder(dbname string, uuid string) (*proto.OrderRecord, error) {
	c := m.session.DB(dbname).C(ordersCollectionName)

	if !bson.IsObjectIdHex(uuid) {
		return nil, errors.New("unexpected ID")
	}

	result := Order{}
	id := bson.ObjectIdHex(uuid)
	if err := c.FindId(id).One(&result); err != nil {
		return nil, err
	}
	return m.mapOrder(dbname, &result), nil
}

// DeleteOrder delete order
func (m *Mongo) DeleteOrder(dbname string, uuid string) error {
	c := m.session.DB(dbname).C(ordersCollectionName)

	if !bson.IsObjectIdHex(uuid) {
		return errors.New("unexpected ID")
	}

	result := Order{}
	if err := c.FindId(bson.ObjectIdHex(uuid)).One(&result); err != nil {
		return err
	}

	ic := m.session.DB(dbname).C(orderItemsCollectionName)
	ic.RemoveAll(bson.M{"order_id": result.ID})

	return c.RemoveId(bson.ObjectIdHex(uuid))
}

// ReadCustomerOrders delete order
func (m *Mongo) ReadCustomerOrders(dbname string, customer string, state string, checkoutState string) (*[]*proto.OrderRecord, error) {
	c := m.session.DB(dbname).C(ordersCollectionName)

	results := []Order{}
	query := bson.M{"customer": customer}
	if len(state) != 0 && len(checkoutState) != 0 {
		query = bson.M{"customer": customer, "state": state, "checkoutState": checkoutState}
	} else if len(state) != 0 {
		query = bson.M{"customer": customer, "state": state}
	} else if len(checkoutState) != 0 {
		query = bson.M{"customer": customer, "checkoutState": checkoutState}
	} 
	if err := c.Find(query).All(&results); err != nil {
		return nil, err
	}

	records := []*proto.OrderRecord{}
	for _, order := range results {
		records = append(records, m.mapOrder(dbname, &order))
	}

	return &records, nil
}
