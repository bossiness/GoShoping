package handler

import (
	"github.com/micro/go-micro/errors"
	"context"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"

	"btdxcx.com/micro/taxons-srv/db"
	"btdxcx.com/os/custom-error"

	proto "btdxcx.com/micro/taxons-srv/proto/taxons"
)

// RootTaxons is a single request handler called via client.RootTaxons or the generated client code
func (h *Handler) RootTaxons(ctx context.Context, req *proto.RootTaxonsRequest, rsp *proto.RootTaxonsResponse) error {
	shopID, err := shopkey.FromCtx(ctx)
	if err != nil {
		return err
	}

	message, err := db.Read(shopID)
	if err != nil {
		return err
	}

	rsp.Message = message
	return nil
}

// CreateTaxons is a single request handler called via client.CreateTaxons or the generated client code
func (h *Handler) CreateTaxons(ctx context.Context, req *proto.CreateTasonsRequest, rsp *proto.CreateTasonsResponse) error {
	shopID, err := shopkey.FromCtx(ctx)
	if err != nil {
		return err
	}

	data := &proto.TaxonsMessage{
		ShopID:      shopID,
		Name:        req.Record.Name,
		Description: req.Record.Description,
		Position:    req.Record.Position,
		Images:      req.Record.Images,
	}

	code, err := db.Create(shopID, data)
	if err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.taxons.CreateChildrenTaxons", err.Error())
	}
	rsp.Record = req.Record
	rsp.Record.Code = code
	return nil
}

// CreateChildrenTaxons is a single request handler called via client.CreateChildrenTaxons or the generated client code
func (h *Handler) CreateChildrenTaxons(ctx context.Context, req *proto.CreateChildrenTaxonsRequest, rsp *proto.CreateChildrenTaxonsResponse) error {
	shopID, err := shopkey.FromCtx(ctx)
	if err != nil {
		return err
	}
	data := &proto.TaxonsMessage{
		ShopID:      shopID,
		Code:        req.Code,
		Name:        req.Record.Name,
		Description: req.Record.Description,
		Position:    req.Record.Position,
		Images:      req.Record.Images,
	}
	code, err := db.Create(shopID, data)
	if err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.taxons.CreateChildrenTaxons", err.Error())
	}
	rsp.Record = req.Record
	rsp.Record.Code = code
	return nil

}

// UpdateTaxons is a single request handler called via client.UpdateTaxons or the generated client code
func (h *Handler) UpdateTaxons(ctx context.Context, req *proto.UpdateTaxonsRequest, rsp *proto.UpdateTaxonsResponse) error {
	shopID, err := shopkey.FromCtx(ctx)
	if err != nil {
		return err
	}
	if err := validateCode(req.Code, "Update"); err != nil {
		return err
	}

	data := &proto.TaxonsMessage{
		ShopID:      shopID,
		Code:        req.Code,
		Name:        req.Record.Name,
		Description: req.Record.Description,
		Position:    req.Record.Position,
		Images:      req.Record.Images,
	}

	rsp.Record = req.Record
	if err := db.Update(shopID, data); err != nil {
		return customerror.Conversion(err, svrName, "Update")
	}

	return nil
}

// DeleteTaxons is a single request handler called via client.DeleteTaxons or the generated client code
func (h *Handler) DeleteTaxons(ctx context.Context, req *proto.DeleteTasonsRequest, rsp *proto.DeleteTasonsResponse) error {
	shopID, err := shopkey.FromCtx(ctx)
	if err != nil {
		return err
	}
	if err := validateCode(req.Code, "Delete"); err != nil {
		return err
	}

	if err := db.Delete(shopID, req.Code); err != nil {
		return customerror.Conversion(err, svrName, "Delete")
	}

	return nil
}
