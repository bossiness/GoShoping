package handler

import (
	"context"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"

	"btdxcx.com/os/custom-error"

	"btdxcx.com/micro/taxons-srv/db"

	"github.com/micro/go-log"

	proto "btdxcx.com/micro/taxons-srv/proto/imp"
)

const (
	svrName = "btdxcx.com/micro/taxons-srv"
)

// Handler is a taxons handler
type Handler struct{}

// Root is a single request handler called via client.Root or the generated client code
func (h *Handler) Root(ctx context.Context, req *proto.TaxonsShopIDRequest, rsp *proto.TaxonsMessage) error {
	log.Log("Received Taxons.Root request")

	shopID, err := shopkey.FromCtx(ctx)
	if err != nil {
		return err
	}

	message, err := db.Read(shopID)
	if err != nil {
		return err
	}

	copyTaxonsMessage(rsp, message)
	return nil
}

func copyTaxonsMessage(dst, src *proto.TaxonsMessage) {
	dst.Code = src.Code
	dst.Name = src.Name
	dst.Description = src.Description
	dst.Images = src.Images
	dst.Children = src.Children
	dst.ShopID = src.ShopID
}

// Create is a single request handler called via client.Create or the generated client code
func (h *Handler) Create(ctx context.Context, req *proto.TasonsCreateRequest, rsp *proto.TasonsCodeResponse) error {

	shopID, err := shopkey.FromCtx(ctx)
	if err != nil {
		return err
	}

	data := &proto.TaxonsMessage{
		ShopID:      shopID,
		Name:        req.Name,
		Description: req.Description,
		Images:      req.Images,
	}

	code, err := db.Create(shopID, data)
	if err != nil {
		return err
	}
	rsp.Code = code
	return nil
}

// CreateChildren is a single request handler called via client.CreateChildren or the generated client code
func (h *Handler) CreateChildren(ctx context.Context, req *proto.TaxonsRequest, rsp *proto.TasonsCodeResponse) error {

	shopID, err := shopkey.FromCtx(ctx)
	if err != nil {
		return err
	}
	data := &proto.TaxonsMessage{
		ShopID:      shopID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Images:      req.Images,
	}
	code, err := db.Create(shopID, data)
	if err != nil {
		return err
	}
	rsp.Code = code
	return nil
}

// Update is a single request handler called via client.Update or the generated client code
func (h *Handler) Update(ctx context.Context, req *proto.TaxonsRequest, rsp *proto.TasonsCodeResponse) error {

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
		Name:        req.Name,
		Description: req.Description,
		Images:      req.Images,
	}
	rsp.Code = req.Code
	return db.Update(shopID, data)
}

// Delete is a single request handler called via client.Delete or the generated client code
func (h *Handler) Delete(ctx context.Context, req *proto.TasonsDeleteRequest, rsp *proto.TasonsCodeResponse) error {

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

	rsp.Code = req.Code
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (h *Handler) Stream(ctx context.Context, req *proto.StreamingRequest, stream proto.Taxons_StreamStream) error {
	log.Logf("Received Taxons.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Logf("Responding: %d", i)
		if err := stream.Send(&proto.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (h *Handler) PingPong(ctx context.Context, stream proto.Taxons_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Logf("Got ping %v", req.Stroke)
		if err := stream.Send(&proto.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

func validateCode(code string, method string) error {

	if len(code) < 5 {
		return customerror.BadRequest(svrName, method, "invalid code")
	}

	return nil
}