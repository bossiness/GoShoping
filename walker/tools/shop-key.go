package tools

import (
	"encoding/hex"
	"strings"

	"btdxcx.com/os/coding"
	"btdxcx.com/walker/apis/common/errors"

	gouuid "github.com/satori/go.uuid"
)

// ShopKEY key handler
type ShopKEY struct {
}

// Coding shop id -> key
func (sk *ShopKEY) Coding(id string) (string, error) {
	aesEnc := coding.AesEncrypt{}
	byteTag, err := aesEnc.Encrypt(id)
	if err != nil {
		return "", errors.InternalServerError("coding.shop.key.tools", err.Error())
	}
	return hex.EncodeToString(byteTag), nil
}

func (sk *ShopKEY) codingUUID(uuid string) (string, error) {
	if _, err := gouuid.FromString(uuid); err != nil {
		return "", errors.BadRequest("coding-uuid.shop.key.tools", err.Error())
	}
	return sk.Coding(uuid)
}

func (sk *ShopKEY) codingWithTag(tag string, uuid string) (string, error) {
	if _, err := gouuid.FromString(uuid); err != nil {
		return "", errors.BadRequest("coding-with-tag.shop.key.tools", err.Error())
	}
	return sk.Coding(tag + "@" + uuid)
}

// CodingWithTags coding tags and uuid -> tag keys
func (sk *ShopKEY) CodingWithTags(tags []string, uuid string) (map[string]string, error) {
	tagKeys := map[string]string{}
	for _, tag := range tags {
		key, err := sk.codingWithTag(tag, uuid)
		if err != nil {
			return tagKeys, err
		}
		tagKeys[tag] = key
	}
	return tagKeys, nil
}

// Decoding shop key -> id
func (sk *ShopKEY) Decoding(key string) (string, error) {
	byteKey, err := hex.DecodeString(key)
	if err != nil {
		return "", errors.BadRequest("decoding.shop.key.tools", err.Error())
	}
	aesEnc := coding.AesEncrypt{}
	return aesEnc.Decrypt(byteKey)
}

// DecodingWithTag shop key -> uuid
func (sk *ShopKEY) DecodingWithTag(key string) (string, string, error) {
	taguuid, err := sk.Decoding(key)
	if err != nil {
		return "", "", err
	}
	ids := strings.Split(taguuid, "@")
	if len(ids) != 2 {
		return "", "", errors.BadRequest("decoding-with-tag.shop.key.tools", "key type invalid")
	}
	uuid := ids[1]
	if _, err := gouuid.FromString(uuid); err != nil {
		return "", taguuid, errors.BadRequest("decoding-with-tag.shop.key.tools", err.Error())
	}
	return ids[0], uuid, nil
}