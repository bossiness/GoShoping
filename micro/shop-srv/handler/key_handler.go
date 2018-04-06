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
type KeyHandler struct{
	Tags []string
}

// Create key
func (key *KeyHandler) Create(ctx context.Context, req *proto.KeyRequest, rsp *proto.KeyResponse) error {

	uuid := req.Uuid
	aesEnc := coding.AesEncrypt{}

	if uuid == "center" {
		byteTag, err := aesEnc.Encrypt(uuid)
		if err != nil {
			return errors.InternalServerError("com.btdxcx.micro.srv.shop.key.Create", err.Error())
		}
		rsp.Keys = &proto.ShopTagKeys{ Tagkeys: map[string]string{"center": hex.EncodeToString(byteTag)} }
		return nil
	}

	if _, err :=  gouuid.FromString(uuid); err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.key.Create", err.Error())
	}

	tagKeys := map[string]string{}
	
	for _, tag := range key.Tags {
		tagUUID := tag + "@" + uuid
		byteTag, err := aesEnc.Encrypt(tagUUID)
		if err != nil {
			return errors.InternalServerError("com.btdxcx.micro.srv.shop.key.Create", err.Error())
		}
		tagKey := hex.EncodeToString(byteTag)
		tagKeys[tag] = tagKey
	}

	shopKeysID := &proto.ShopTagKeys{ Tagkeys: tagKeys }

	if err := db.CreateKey(uuid, shopKeysID); err != nil {
		return errors.InternalServerError("com.btdxcx.micro.srv.shop.key.Create", err.Error())
	}

	rsp.Keys = shopKeysID
	return nil
}

// Read key
func (key *KeyHandler) Read(ctx context.Context, req *proto.KeyRequest, rsp *proto.KeyResponse) error {
	shopKeyID, err := db.ReadKey(req.Uuid)
	if err != nil {
		return errors.BadRequest("com.btdxcx.micro.srv.shop.key.Read", err.Error())
	}

	rsp.Keys = shopKeyID
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
	
	uuid, ctype, err := introspectShopKey(key.Tags, req.Key)
	if err != nil {
		return err
	}
	rsp.Uuid = uuid
	rsp.Type = ctype
	rsp.Active = true
	return nil
}

func introspectShopKey(tags []string, shopKey string) (string, string, error) {

	byteKey, err1 := hex.DecodeString(shopKey)
	if err1 != nil {
		return "", "", errors.BadRequest("com.btdxcx.micro.srv.shop.key.Introspect", err1.Error())
	}
	aesEnc := coding.AesEncrypt{}
	clileID, err2 := aesEnc.Decrypt(byteKey)
	if err2 != nil {
		return "", "", errors.InternalServerError("com.btdxcx.micro.srv.shop.key.Introspect", err2.Error())
	}
	if clileID == "center" {
		return clileID, "center", nil
	}

	ids := strings.Split(clileID, "@")

	if len(ids) != 2 {
		return "", "", errors.BadRequest("com.btdxcx.micro.srv.shop.key.Introspect", "shop-key type invalid")
	}
	uuid := ids[1]
	if _, err := db.ReadKey(uuid); err != nil {
		return "", "", errors.BadRequest("com.btdxcx.micro.srv.shop.key.Read", err.Error())
	}
	for _, tag := range tags {
		if ids[0] == tag {
			return uuid, tag, nil
		}
	}
	
	return "", "", errors.BadRequest("com.btdxcx.micro.srv.shop.key.Introspect", "shop-key type invalid")
}