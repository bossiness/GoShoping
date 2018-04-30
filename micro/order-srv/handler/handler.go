package handler

import (
	"context"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"btdxcx.com/micro/order-srv/db"
	"github.com/micro/go-micro/errors"

	proto "btdxcx.com/micro/order-srv/proto/order"
)

const (
	svrName = "com.btdxcx.micro.srv.order"
)

// Handler order
type Handler struct{}

// CreateCart is a single request handler called via client.CreateCart or the generated client code
func (h *Handler) CreateCart(ctx context.Context, req *proto.CreateCartRequest, rsp *proto.OrderResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	orderID, err := db.CreateOrder(shopID, req.Customer)
	if err != nil {
		return errors.InternalServerError(svrName + ".CreateCart", err.Error())
	}
	record, err1 := db.ReadOrder(shopID, orderID)
	if err1 != nil {
		return errors.InternalServerError(svrName + ".CreateCart", err1.Error())
	}

	rsp.Record = record
	return nil
}

// ReadOrders is a single request handler called via client.ReadOrders or the generated client code
func (h *Handler) ReadOrders(ctx context.Context, req *proto.ReadOrdersRequest, rsp *proto.ReadOrdersResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	orders, err := db.ReadOrders(shopID, req.State, req.CheckoutState, int(req.Offset), int(req.Limit))
	if err != nil {
		return errors.NotFound(svrName + ".ReadOrders", err.Error())
	}

	rsp.Limit = req.Limit
	rsp.Offset = req.Offset
	rsp.Total = int32(len(*orders))
	rsp.Records = *orders

	return nil
}

// ReadOrder is a single request handler called via client.ReadOrder or the generated client code
func (h *Handler) ReadOrder(ctx context.Context, req *proto.ReadOrderRequest, rsp *proto.OrderResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	product, err := db.ReadOrder(shopID, req.Uuid)
	if err != nil {
		return errors.NotFound(svrName + ".ReadOrder", err.Error())
	}
	rsp.Record = product
	return nil
}

// DeleteOrder is a single request handler called via client.DeleteOrder or the generated client code
func (h *Handler) DeleteOrder(ctx context.Context, req *proto.DeleteOrderRequest, rsp *proto.Response) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.DeleteOrder(shopID, req.Uuid); err != nil {
		return errors.NotFound(svrName + ".DeleteOrder", err.Error())
	}
	return nil
}

// ReadCustomerOrders is a single request handler called via client.ReadCustomerOrders or the generated client code
func (h *Handler) ReadCustomerOrders(ctx context.Context, req *proto.ReadCustomerOrdersRequest, rsp *proto.ReadCustomerOrdersResponse) error {
	return nil
}

// CreateCartItem is a single request handler called via client.CreateCartItem or the generated client code
func (h *Handler) CreateCartItem(context.Context, *proto.CreateCartItemRequest, *proto.CartItemResponse) error {
	return nil
}

// UpdateCartItem is a single request handler called via client.UpdateCartItem or the generated client code
func (h *Handler) UpdateCartItem(context.Context, *proto.UpdateCartItemRequest, *proto.CartItemResponse) error {
	return nil
}

// DeleteCartItem is a single request handler called via client.DeleteCartItem or the generated client code
func (h *Handler) DeleteCartItem(context.Context, *proto.DeleteCartItemRequest, *proto.Response) error {
	return nil
}

// CheckoutAddressing is a single request handler called via client.CheckoutAddressing or the generated client code
func (h *Handler) CheckoutAddressing(context.Context, *proto.CheckoutAddressingRequest, *proto.Response) error {
	return nil
}

// CheckoutSelectShipping is a single request handler called via client.CheckoutSelectShipping or the generated client code
func (h *Handler) CheckoutSelectShipping(context.Context, *proto.CheckoutSelectShippingRequest, *proto.Response) error {
	return nil
}

// CheckoutSelectPayment is a single request handler called via client.CheckoutSelectPayment or the generated client code
func (h *Handler) CheckoutSelectPayment(context.Context, *proto.CheckoutSelectPaymentRequest, *proto.Response) error {
	return nil
}

// CheckoutComplete is a single request handler called via client.CheckoutComplete or the generated client code
func (h *Handler) CheckoutComplete(context.Context, *proto.CheckoutCompleteRequest, *proto.Response) error {
	return nil
}

// CancelOrder is a single request handler called via client.CancelOrder or the generated client code
func (h *Handler) CancelOrder(context.Context, *proto.CancelOrderRequest, *proto.Response) error {
	return nil
}

// ShipmentShip is a single request handler called via client.ShipmentShip or the generated client code
func (h *Handler) ShipmentShip(context.Context, *proto.ShipmentShipRequest, *proto.Response) error {
	return nil
}

// ShipmentComplete is a single request handler called via client.ShipmentComplete or the generated client code
func (h *Handler) ShipmentComplete(context.Context, *proto.ShipmentCompleteRequest, *proto.Response) error {
	return nil
}
