package handler

import (
	"context"

	"btdxcx.com/micro/shop-srv/db"
	"github.com/micro/go-micro/errors"
	"github.com/satori/go.uuid"

	proto "btdxcx.com/micro/shop-srv/proto/shop/details"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
)

// DetailsHandler details handler
type DetailsHandler struct {
}

// Create Shop Details
func (h *DetailsHandler) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {

	if req.Details == nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.details.Create", "The parameter cannot be empty.")
	}
	if req.Details.Owner == nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.details.Create", "Owner cannot be empty.")
	}
	if len(req.Details.Owner.Name) == 0 {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.details.Create", "Owner name cannot be empty.")
	}

	shop, err := db.CreateDetails(req)
	if err != nil {
		return err
	}

	rsp.ShopId = shop.ShopId
	return nil
}

// Read Shop Details
func (h *DetailsHandler) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if _, err := uuid.FromString(shopID); err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.details.Read", err.Error())
	}

	result, err := db.ReadDetails(shopID)
	if err != nil {
		return errors.NotFound("com.btdxcx.micro.srv.shop.details.Read", err.Error())
	}

	rsp.ShopId = result.ShopId
	rsp.CreateAt = result.CreateAt
	rsp.UpdateAt = result.UpdateAt
	rsp.SubmitAt = result.SubmitAt
	rsp.PeriodAt = result.PeriodAt
	rsp.State = result.State
	rsp.Details = result.Details

	return nil
}

// Delete Shop Details
func (h *DetailsHandler) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {

	if _, err := uuid.FromString(req.ShopId); err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.details.Delete", err.Error())
	}

	if err := db.DeleteDetails(req.ShopId); err != nil {
		return errors.NotFound("com.btdxcx.micro.srv.shop.details.Delete", err.Error())
	}
	return nil
}

// List Shop Details
func (h *DetailsHandler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	if req.Start < 0 {
		req.Start = 0
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	list, err := db.ListDetails(req)
	if err != nil {
		return errors.NotFound("com.btdxcx.micro.srv.shop.details.List", err.Error())
	}
	rsp.Start = list.Start
	rsp.Limit = list.Limit
	rsp.Total = list.Total
	rsp.Items = list.Items

	return nil
}

// Update Shop Details
func (h *DetailsHandler) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	if _, err := uuid.FromString(shopID); err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.details.Read", err.Error())
	}
	req.ShopId = shopID

	if err := db.UpdateDetails(req); err != nil {
		return errors.NotFound("com.btdxcx.micro.srv.shop.details.Update", err.Error())
	}
	return nil
}
