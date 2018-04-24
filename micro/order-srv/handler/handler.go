package handler

import (
	"context"

	proto "btdxcx.com/micro/order-srv/proto/order"
)

// Handler order
type Handler struct{}

// CreateCart is a single request handler called via client.CreateCart or the generated client code
func (h *Handler) CreateCart(context.Context, *proto.CreateCartRequest, *proto.OrderResponse) error {
	return nil
}

// ReadOrders is a single request handler called via client.ReadOrders or the generated client code
func (h *Handler) ReadOrders(context.Context, *proto.ReadOrdersRequest, *proto.ReadOrdersResponse) error {
	return nil
}

// ReadOrder is a single request handler called via client.ReadOrder or the generated client code
func (h *Handler) ReadOrder(context.Context, *proto.ReadOrderRequest, *proto.OrderResponse) error {
	return nil
}

// DeleteOrder is a single request handler called via client.DeleteOrder or the generated client code
func (h *Handler) DeleteOrder(context.Context, *proto.DeleteOrderRequest, *proto.Response) error {
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
