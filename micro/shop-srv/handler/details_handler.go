package handler

import (
	"github.com/micro/go-micro/errors"
	"github.com/satori/go.uuid"
	"btdxcx.com/micro/shop-srv/db"
	"context"

	proto "btdxcx.com/micro/shop-srv/proto/shop/details"
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
	
	if _, err := uuid.FromString(req.ShopId); err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.details.Read", err.Error())
	}

	item, err := db.ReadDetails(req.ShopId)
	if err != nil {
		return errors.NotFound("com.btdxcx.micro.srv.shop.details.Read", err.Error())
	}

	rsp.ShopId = item.ShopId
	rsp.CreateAt = item.CreateAt
	rsp.UpdateAt = item.UpdateAt
	rsp.SubmitAt = item.SubmitAt
	rsp.PeriodAt = item.PeriodAt
	rsp.State = item.State
	rsp.Details = item.Details
	
	return nil
}

// Delete Shop Details
func (h *DetailsHandler) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {

	if _, err := uuid.FromString(req.ShopId); err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.details.Delete", err.Error())
	}

	if err := db.DeleteDetails(req.ShopId); err != nil {
		return errors.InternalServerError("com.btdxcx.micro.srv.shop.details.Delete", err.Error())
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
		return errors.InternalServerError("com.btdxcx.micro.srv.shop.details.Delete", err.Error())
	}
	rsp.Start = list.Start
	rsp.Limit = list.Limit
	rsp.Total = list.Total
	rsp.Items = list.Items

	return nil
}

// Update Shop Details
func (h *DetailsHandler) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {

	return db.UpdateDetails(req)
}
