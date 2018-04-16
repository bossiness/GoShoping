package handler

import (
	"btdxcx.com/micro/product-srv/db"
	"context"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-micro/errors"

	proto "btdxcx.com/micro/product-srv/proto/product"
)

// AttributeHandler attribute handler
type AttributeHandler struct{}

// CreateAttribute is a single request handler called via client.CreateAttribute or the generated client code
func (a *AttributeHandler) CreateAttribute(ctx context.Context, req *proto.CreateAttributeRequest, rsp *proto.CreateAttributeResponse) error {

	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err2 := db.CreateAttribute(shopID, req.Record)
	if err1 != nil {
		return errors.InternalServerError(svrName + ".CreateAttribute", err2.Error())
	}

	rsp.Record = req.Record
	return nil
}

// ReadAttributes is a single request handler called via client.ReadAttributes or the generated client code
func (a *AttributeHandler) ReadAttributes(ctx context.Context, req *proto.ReadAttributesRequest, rsp *proto.ReadAttributesResponse) error{
	
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	attributes, err2 := db.ReadAttributes(shopID, int(req.Offset), int(req.Limit))
	if err2 != nil {
		return errors.NotFound(svrName + ".ReadAttributes", err2.Error())
	}
	
	rsp.Limit = req.Limit
	rsp.Offset = req.Offset
	rsp.Total = int32(len(*attributes))
	rsp.Records = *attributes

	return nil
}

// ReadAttribute is a single request handler called via client.ReadAttribute or the generated client code
func (a *AttributeHandler) ReadAttribute(ctx context.Context, req *proto.ReadAttributeRequest, rsp *proto.ReadAttributeResponse) error{
	
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	attribute, err2 := db.ReadAttribute(shopID, req.Code)
	if err2 != nil {
		return errors.NotFound(svrName + ".ReadAttribute", err2.Error())
	}

	rsp.Record = attribute

	return nil
}

// UpdateAttribute is a single request handler called via client.UpdateAttribute or the generated client code
func (a *AttributeHandler) UpdateAttribute(ctx context.Context, req *proto.UpdateAttributeRequest, rsp *proto.UpdateAttributeResponse) error{
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if err := db.UpdateAttribute(shopID, req.Record.Code, req.Record); err != nil {
		return errors.NotFound(svrName + ".UpdateAttribute", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// DeleteAttribute is a single request handler called via client.DeleteAttribute or the generated client code
func (a *AttributeHandler) DeleteAttribute(ctx context.Context, req *proto.DeleteAttributeRequest, rsp *proto.DeleteAttributeResponse) error{
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	db.DeleteAttribute(shopID, req.Code)
	if err := db.DeleteAttribute(shopID, req.Code); err != nil {
		return errors.NotFound(svrName + ".DeleteAttribute", err.Error())
	}

	return nil
}