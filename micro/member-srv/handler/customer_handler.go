package handler

import (
	"btdxcx.com/micro/member-srv/db"
	"context"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-micro/errors"

	proto "btdxcx.com/micro/member-srv/proto/member"
)

// CustomerHandler member
type CustomerHandler struct{}

// CreateCustomer is a single request handler called via client.CreateCustomer or the generated client code
func (e *CustomerHandler) CreateCustomer(ctx context.Context, req *proto.CreateCustomerRequest, rsp *proto.CustomerResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}
	
	err := db.CreateCustomer(shopID, req.Record)
	if err != nil {
		return errors.InternalServerError(svrName + ".CreateCustomer", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// ReadCustomers is a single request handler called via client.CreateCustomer or the generated client code
func (e *CustomerHandler) ReadCustomers(ctx context.Context, req *proto.ReadCustomersRequest, rsp *proto.ReadCustomersResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	customers, err := db.ReadCustomers(shopID, int(req.Offset), int(req.Limit))
	if err != nil {
		return errors.BadRequest(svrName + ".ReadCustomers", err.Error())
	}

	rsp.Limit = req.Limit
	rsp.Offset = req.Offset
	rsp.Total = int32(len(*customers))
	rsp.Records = *customers

	return nil
}

// ReadCustomer is a single request handler called via client.ReadCustomer or the generated client code
func (e *CustomerHandler) ReadCustomer(ctx context.Context, req *proto.ReadCustomerRequest, rsp *proto.CustomerResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	customer, err := db.ReadCustomer(shopID, req.Id)
	if err != nil {
		return errors.BadRequest(svrName + ".ReadCustomer", err.Error())
	}

	rsp.Record = customer
	return nil
}

// ReadCustomerFormName is a single request handler called via client.ReadCustomerFormName or the generated client code
func (e *CustomerHandler) ReadCustomerFormName(ctx context.Context, req *proto.ReadCustomerFormNameRequest, rsp *proto.CustomerResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	customer, err := db.ReadCustomerFromName(shopID, req.Name)
	if err != nil {
		return errors.BadRequest(svrName + ".ReadCustomer", err.Error())
	}

	rsp.Record = customer
	return nil
}

// UpdateCustomer is a single request handler called via client.UpdateCustomer or the generated client code
func (e *CustomerHandler) UpdateCustomer(ctx context.Context, req *proto.UpdateCustomerRequest, rsp *proto.CustomerResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err := db.UpdateCustomer(shopID, req.Id, req.Record)
	if err != nil {
		return errors.BadRequest(svrName + ".UpdateCustomer", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// DeleteCustomer is a single request handler called via client.DeleteCustomer or the generated client code
func (e *CustomerHandler) DeleteCustomer(ctx context.Context, req *proto.DeleteCustomerRequest, rsp *proto.DeleteCustomerResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err := db.DeleteCustomer(shopID, req.Id)
	if err != nil {
		return errors.BadRequest(svrName + ".DeleteCustomer", err.Error())
	}
	
	return nil
}
