package mongodb

import (
	"gopkg.in/mgo.v2"
	"errors"
	"time"

	proto "btdxcx.com/micro/order-srv/proto/order"
	productdb "btdxcx.com/micro/product-srv/db/mongodb"
	"gopkg.in/mgo.v2/bson"
)

const (
	variantsCollectionName = "variants"
	orderItemsCollectionName = "orderItems"
)

// OrderItem DB
type OrderItem struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	OrderID     bson.ObjectId `bson:"order_id,omitempty"`
	Quantity    int64         `bson:"quantity,omitempty"`
	UnitPrice   int64         `bson:"unitPrice,omitempty"`
	Adjustments []string      `bson:"adjustments,omitempty"`
	Variant     string        `bson:"variant,omitempty"`
	CreatedAt   int64         `bson:"created_at,omitempty"`
	UpdatedAt   int64         `bson:"updated_at,omitempty"`
}

// CreateOrderItem create order item
func (m *Mongo) CreateOrderItem(dbname string, order string, item *proto.OrderRecord_Item) (string, error) {
	c := m.session.DB(dbname).C(orderItemsCollectionName)

	vc := m.session.DB(dbname).C(variantsCollectionName)

	variant := new(productdb.Variant)
	if err := vc.Find(bson.M{"sku": item.Variant}).One(&variant); err != nil {
		return "", errors.New("sku not found")
	}

	index := mgo.Index{
		Key:        []string{"order_id", "variant"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return "", err
	}

	unitPrice := int64(variant.Pricings.Current)
	doc := &OrderItem{
		ID:        bson.NewObjectId(),
		OrderID:   bson.ObjectIdHex(order),
		Quantity:  item.Quantity,
		UnitPrice: unitPrice,
		Variant:   item.Variant,
		CreatedAt: time.Now().Unix(),
	}
	if err := c.Insert(doc); err != nil {
		return "", err
	}

	return doc.ID.Hex(), nil
}

// UpdateOrderItem update order item
func (m *Mongo) UpdateOrderItem(dbname string, id string, item *proto.OrderRecord_Item) error {
	c := m.session.DB(dbname).C(orderItemsCollectionName)
	selector := bson.ObjectIdHex(id)
	updataData := bson.M{"$set": bson.M{
		"quantity":   item.Quantity,
		// "unitPrice":  item.UnitPrice,
		"updated_at": time.Now()}}

	return c.UpdateId(selector, updataData)
}

// DeleteOrderItem delete order item
func (m *Mongo) DeleteOrderItem(dbname string, id string) error {
	c := m.session.DB(dbname).C(orderItemsCollectionName)
	return c.RemoveId(bson.ObjectIdHex(id))
}
