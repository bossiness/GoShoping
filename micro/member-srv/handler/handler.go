package handler

import (
	"context"

	proto "btdxcx.com/micro/member-srv/proto/member"
)

// Handler member
type Handler struct{}

// CreateCustomer is a single request handler called via client.CreateCustomer or the generated client code
func (e *Handler) CreateCustomer(context.Context, *proto.CreateCustomerRequest, *proto.CustomerResponse) error {
	return nil
}

// ReadCustomers is a single request handler called via client.CreateCustomer or the generated client code
func (e *Handler) ReadCustomers(context.Context, *proto.ReadCustomersRequest, *proto.ReadCustomersResponse) error {
	return nil
}

// ReadCustomer is a single request handler called via client.ReadCustomer or the generated client code
func (e *Handler) ReadCustomer(context.Context, *proto.ReadCustomerRequest, *proto.CustomerResponse) error {
	return nil
}

// UpdateCustomer is a single request handler called via client.UpdateCustomer or the generated client code
func (e *Handler) UpdateCustomer(context.Context, *proto.UpdateCustomerRequest, *proto.CustomerResponse) error {
	return nil
}

// DeleteCustomer is a single request handler called via client.DeleteCustomer or the generated client code
func (e *Handler) DeleteCustomer(context.Context, *proto.DeleteCustomerRequest, *proto.DeleteCustomerResponse) error {
	return nil
}

