package handler

import (
	"strings"
	"btdxcx.com/micro/shop-srv/db"
	"github.com/micro/go-micro/errors"
	"encoding/hex"
	"btdxcx.com/os/coding"
	"context"

	proto "btdxcx.com/micro/shop-srv/proto/shop"
	gouuid "github.com/satori/go.uuid"
)

// KeyHandler key handler
type KeyHandler struct{}

// Create key
func (key *KeyHandler) Create(ctx context.Context, req *proto.KeyRequest, rsp *proto.KeyResponse) error {

	uuid := req.Uuid
	if _, err :=  gouuid.FromString(uuid); err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.key.Create", err.Error())
	}
	backUUID := "back@" + uuid
	miniUUID := "mini@" + uuid
	aesEnc := coding.AesEncrypt{}

	byteBack, err1 := aesEnc.Encrypt(backUUID)
	if err1 != nil {
		return errors.InternalServerError("com.btdxcx.micro.srv.shop.key.Create", err1.Error())
	}
	byteMini, err2 := aesEnc.Encrypt(miniUUID)
	if err2 != nil {
		return errors.InternalServerError("com.btdxcx.micro.srv.shop.key.Create", err2.Error())
	}

	backKey := hex.EncodeToString(byteBack)
	miniKey := hex.EncodeToString(byteMini)
	shopKeyID := &proto.ShopKeyID { BackKey: backKey, MiniKey: miniKey  }

	if err := db.CreateKey(uuid, shopKeyID); err != nil {
		return errors.InternalServerError("com.btdxcx.micro.srv.shop.key.Create", err.Error())
	}

	rsp.Key = shopKeyID
	return nil
}

// Read key
func (key *KeyHandler) Read(ctx context.Context, req *proto.KeyRequest, rsp *proto.KeyResponse) error {
	shopKeyID, err := db.ReadKey(req.Uuid)
	if err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.key.Read", err.Error())
	}

	rsp.Key = shopKeyID
	return nil
}

// Delete key
func (key *KeyHandler) Delete(ctx context.Context, req *proto.KeyRequest, rsp *proto.DeleteKeyResponse) error {
	if err := db.DeleteKey(req.Uuid); err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.key.Delete", err.Error())
	}
	return nil
}

// Introspect key
func (key *KeyHandler) Introspect(ctx context.Context, req *proto.IntrospectRequest, rsp *proto.IntrospectResponse) error {
	
	uuid, ctype, err := introspectShopKey(req.Key)
	if err != nil {
		return err
	}
	rsp.Uuid = uuid
	rsp.Type = ctype
	rsp.Active = true
	return nil
}

func introspectShopKey(shopKey string) (string, string, error) {

	byteKey, err1 := hex.DecodeString(shopKey)
	if err1 != nil {
		return "", "", errors.BadRequest("com.btdxcx.micro.srv.shop.key.Introspect", err1.Error())
	}
	aesEnc := coding.AesEncrypt{}
	clileID, err2 := aesEnc.Decrypt(byteKey)
	if err2 != nil {
		return "", "", errors.InternalServerError("com.btdxcx.micro.srv.shop.key.Introspect", err2.Error())
	}
	ids := strings.Split(clileID, "@")

	if len(ids) != 2 {
		return "", "", errors.BadRequest("com.btdxcx.micro.srv.shop.key.Introspect", "shop-key type invalid")
	}
	uuid := ids[1]
	if _, err := db.ReadKey(uuid); err != nil {
		return "", "", errors.BadRequest("com.btdxcx.micro.srv.shop.key.Read", err.Error())
	}
	if ids[0] == "back" || ids[0] == "mini" {
		return uuid, ids[0], nil
	}
	
	return "", "", errors.BadRequest("com.btdxcx.micro.srv.shop.key.Introspect", "shop-key type invalid")
}