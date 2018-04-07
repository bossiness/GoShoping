package handler

import (
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"context"

	"btdxcx.com/micro/account-srv/db"
	"btdxcx.com/os/custom-error"
	"github.com/micro/go-log"

	proto "btdxcx.com/micro/account-srv/proto/account"
)

const (
	serviceName = "com.btdxcx.micro.srv.account"
)

// Handler is account service handler
type Handler struct{}

// Read account
func (h *Handler) Read(ctx context.Context, req *proto.ReadRequest, rsp *proto.ReadResponse) error {
	log.Log("Received Account.Read request")

	shopID, err0 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err0 != nil {
		return err0
	}

	if err := customerror.ValidateShopIDAndName(shopID, req.ClientId, serviceName, "Read"); err != nil {
		return err
	}
	record, err := db.Read(shopID, req.ClientId)
	if err != nil {
		return customerror.Conversion(err, serviceName, "Read")
	}
	rsp.Account = record
	return nil
}

// Create account
func (h *Handler) Create(ctx context.Context, req *proto.CreateRequest, rsp *proto.CreateResponse) error {
	shopID, err0 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err0 != nil {
		return err0
	}

	if err := customerror.ValidateShopIDAndName(shopID, req.Account.ClientId, serviceName, "Create"); err != nil {
		return err
	}

	return customerror.Conversion(db.Create(shopID, req.Account), serviceName, "Create")
}

// Update account
func (h *Handler) Update(ctx context.Context, req *proto.UpdateRequest, rsp *proto.UpdateResponse) error {
	shopID, err0 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err0 != nil {
		return err0
	}

	if err := customerror.ValidateShopIDAndName(shopID, req.Account.ClientId, serviceName, "Update"); err != nil {
		return err
	}
	return customerror.Conversion(db.Update(shopID, req.Account), serviceName, "Create")
}

// Delete account
func (h *Handler) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	shopID, err0 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err0 != nil {
		return err0
	}

	if err := customerror.ValidateShopIDAndID(shopID, req.Id, serviceName, "Delete"); err != nil {
		return err
	}
	return customerror.Conversion(db.Delete(shopID, req.Id), serviceName, "Delete")
}

// Search account
func (h *Handler) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.SearchResponse) error {
	shopID, err0 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err0 != nil {
		return err0
	}

	if err := customerror.ValidateShopIDAndName(shopID, req.ClientId, serviceName, "Search"); err != nil {
		return err
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	records, err := db.Search(shopID, req)
	if err != nil {
		return customerror.Conversion(err, serviceName, "Search")
	}
	rsp.Accounts = *records
	return nil
}
