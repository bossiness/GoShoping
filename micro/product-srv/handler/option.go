package handler

import (
	"btdxcx.com/micro/product-srv/db"
	"context"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-micro/errors"

	proto "btdxcx.com/micro/product-srv/proto/product"
)

// OptionHandler option handler
type OptionHandler struct{}


// CreateOption is a single request handler called via client.CreateOption or the generated client code
func (o *OptionHandler) CreateOption(ctx context.Context, req *proto.CreateOptionRequest, rsp *proto.CreateOptionResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err := db.CreateOption(shopID, req.Record)
	if err1 != nil {
		return errors.InternalServerError(svrName + ".CreateOption", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// ReadOptions is a single request handler called via client.ReadOptions or the generated client code
func (o *OptionHandler) ReadOptions(ctx context.Context, req *proto.ReadOptionsRequest, rsp *proto.ReadOptionsResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	options, err := db.ReadOptions(shopID, int(req.Offset), int(req.Limit))
	if err != nil {
		return errors.NotFound(svrName + ".ReadOptions", err.Error())
	}

	rsp.Limit = req.Limit
	rsp.Offset = req.Offset
	rsp.Total = int32(len(*options))
	rsp.Records = *options

	return nil
}

// ReadOption is a single request handler called via client.ReadOption or the generated client code
func (o *OptionHandler) ReadOption(ctx context.Context, req *proto.ReadOptionequest, rsp *proto.ReadOptionResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	record, err := db.ReadOption(shopID, req.Code)
	if err != nil {
		return errors.NotFound(svrName + ".ReadOption", err.Error())
	}

	rsp.Record = record
	return nil
}

// UpdateOption is a single request handler called via client.UpdateOption or the generated client code
func (o *OptionHandler) UpdateOption(ctx context.Context, req *proto.UpdateOptionRequest, rsp *proto.UpdateOptionResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err := db.UpdateOption(shopID, req.Code, req.Record)
	if err != nil {
		return errors.NotFound(svrName + ".ReadOption", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// DeleteOption is a single request handler called via client.DeleteOption or the generated client code
func (o *OptionHandler) DeleteOption(ctx context.Context, req *proto.DeleteOptionRequest, rsp *proto.DeleteOptionResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err := db.DeleteOption(shopID, req.Code)
	if err != nil {
		return errors.NotFound(svrName + ".DeleteOption", err.Error())
	}
	return nil
}

