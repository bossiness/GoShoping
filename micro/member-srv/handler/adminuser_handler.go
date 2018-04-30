package handler

import (
	"btdxcx.com/micro/member-srv/db"
	"context"
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-micro/errors"

	proto "btdxcx.com/micro/member-srv/proto/member"
)

// AdminUserHandler member
type AdminUserHandler struct{}

// CreateAdminUser is a single request handler called via client.CreateAdminUser or the generated client code
func (e *AdminUserHandler) CreateAdminUser(ctx context.Context, req *proto.CreateAdminUserRequest, rsp *proto.AdminUserResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}
	
	err := db.CreateAdminUser(shopID, req.Record)
	if err != nil {
		return errors.InternalServerError(svrName + ".CreateAdminUser", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// ReadAdminUsers is a single request handler called via client.CreateAdminUser or the generated client code
func (e *AdminUserHandler) ReadAdminUsers(ctx context.Context, req *proto.ReadAdminUsersRequest, rsp *proto.ReadAdminUsersResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	adminusers, err := db.ReadAdminUsers(shopID, int(req.Offset), int(req.Limit))
	if err != nil {
		return errors.NotFound(svrName + ".ReadAdminUsers", err.Error())
	}

	rsp.Limit = req.Limit
	rsp.Offset = req.Offset
	rsp.Total = int32(len(*adminusers))
	rsp.Records = *adminusers

	return nil
}

// ReadAdminUser is a single request handler called via client.ReadAdminUser or the generated client code
func (e *AdminUserHandler) ReadAdminUser(ctx context.Context, req *proto.ReadAdminUserRequest, rsp *proto.AdminUserResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	adminuser, err := db.ReadAdminUser(shopID, req.Id)
	if err != nil {
		return errors.NotFound(svrName + ".ReadAdminUser", err.Error())
	}

	rsp.Record = adminuser
	return nil
}

// ReadAdminUserFormName is a single request handler called via client.ReadAdminUserFormName or the generated client code
func (e *AdminUserHandler) ReadAdminUserFormName(ctx context.Context, req *proto.ReadAdminUserFormNameRequest, rsp *proto.AdminUserResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	adminuser, err := db.ReadAdminUserFromName(shopID, req.Name)
	if err != nil {
		return errors.NotFound(svrName + ".ReadAdminUser", err.Error())
	}

	rsp.Record = adminuser
	return nil
}

// UpdateAdminUser is a single request handler called via client.UpdateAdminUser or the generated client code
func (e *AdminUserHandler) UpdateAdminUser(ctx context.Context, req *proto.UpdateAdminUserRequest, rsp *proto.AdminUserResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err := db.UpdateAdminUser(shopID, req.Id, req.Record)
	if err != nil {
		return errors.NotFound(svrName + ".UpdateAdminUser", err.Error())
	}

	rsp.Record = req.Record
	return nil
}

// DeleteAdminUser is a single request handler called via client.DeleteAdminUser or the generated client code
func (e *AdminUserHandler) DeleteAdminUser(ctx context.Context, req *proto.DeleteAdminUserRequest, rsp *proto.DeleteAdminUserResponse) error {
	shopID, err1 := shopkey.GetShopIDFrom(ctx, req.ShopId)
	if err1 != nil {
		return err1
	}

	err := db.DeleteAdminUser(shopID, req.Id)
	if err != nil {
		return errors.NotFound(svrName + ".DeleteAdminUser", err.Error())
	}
	
	return nil
}
