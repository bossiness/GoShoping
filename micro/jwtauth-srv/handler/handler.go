package handler

import (
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	twrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"btdxcx.com/micro/jwtauth-srv/db"
	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/errors"

	"github.com/micro/go-log"
	"github.com/satori/go.uuid"

	proto "btdxcx.com/micro/jwtauth-srv/proto/auth"
)

// Handler handler
type Handler struct{}

// Token create
func (h *Handler) Token(ctx context.Context, req *proto.TokenRequest, rsp *proto.TokenResponse) error {
	log.Log("Received Handler.Token request")
	jwtID, err := uuid.NewV4()
	if err != nil {
		return errors.InternalServerError("com.btdxcx.micro.srv.jwtauth", err.Error())
	}
	clientID := req.ClientId
	clientSecrent := req.ClientSecrent
	shopID, err0 := shopkey.FromCtx(ctx)
	if err0 != nil {
		if len(req.ShopId) == 0 {
			return err0
		} 
		shopID = req.ShopId
	}

	if len(clientID) < 3 {
		return errors.BadRequest("com.btdxcx.micro.srv.jwtauth.Token", "client_id invalid")
	}
	if len(clientSecrent) < 4 {
		return errors.BadRequest("com.btdxcx.micro.srv.jwtauth.Token", "client_secrent invalid")
	}
	if len(shopID) < 6 {
		return errors.BadRequest("com.btdxcx.micro.srv.jwtauth.Token", "shop_id invalid")
	}

	hash := md5.New()
	hash.Write([]byte(clientSecrent))
	cipher := hex.EncodeToString(hash.Sum(nil))

	refreshToken, err := refreshToken(jwtID.String(), clientID, cipher, shopID)
	if err != nil {
		return err
	}

	exp := time.Now().Add(time.Hour).Unix()

	accessToken, err := accessToken(
		jwtID.String(),
		clientID,
		cipher,
		exp,
		req.Scopes,
		req.Metadata,
		shopID)
	if err != nil {
		return err
	}

	record := &proto.Record{
		ClientId:     clientID,
		Jti:          jwtID.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Cipher:       cipher}
	if err := db.Create(shopID, record); err != nil {
		return errors.InternalServerError("com.btdxcx.micro.srv.jwtauth", err.Error())
	}

	rsp.Token = &proto.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    exp,
		Scopes:       req.Scopes,
		Metadata:     req.Metadata}

	return nil
}

func refreshToken(jti string, clientID string, secrent string, shopID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
		"client_id": clientID,
		"jti":       jti,
		"shop_id":   shopID,
	})

	tokenString, err := token.SignedString([]byte(secrent))
	if err != nil {
		return "", errors.InternalServerError("com.btdxcx.micro.srv.jwtauth", err.Error())
	}
	return tokenString, nil
}

func accessToken(jti string,
	clientID string, secrent string,
	exp int64, scopes []string, metadata map[string]string,
	shopID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":       exp,
		"client_id": clientID,
		"jti":       jti,
		"scopes":    scopes,
		"metadata":  metadata,
		"shop_id":   shopID,
	})
	tokenString, err := token.SignedString([]byte(secrent))
	if err != nil {
		return "", errors.InternalServerError("com.btdxcx.micro.srv.jwtauth", err.Error())
	}
	return tokenString, nil
}

// Revoke delete
func (h *Handler) Revoke(ctx context.Context, req *proto.RevokeRequest, rsp *proto.RevokeResponse) error {
	shopID, err0 := shopkey.FromCtx(ctx)
	if err0 != nil {
		return err0
	}

	if len(shopID) < 6 {
		return errors.BadRequest("com.btdxcx.micro.srv.jwtauth.Revoke", "ShopId is empty.")
	}

	if len(req.AccessToken) > 10 {
		if err := db.DeleteAccessToken(shopID, req.AccessToken); err != nil {
			return errors.InternalServerError("com.btdxcx.micro.srv.jwtauth.Revoke", err.Error())
		}
	} else if len(req.RefreshToken) > 10 {
		if err := db.DeleteRefreshToken(shopID, req.RefreshToken); err != nil {
			return errors.InternalServerError("com.btdxcx.micro.srv.jwtauth.Revoke", err.Error())
		}
	} else {
		return errors.BadRequest("com.btdxcx.micro.srv.jwtauth.Revoke", "token is empty.")
	}

	return nil
}

func parseToken(tokenString string, secret string) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return claims, err
}

// Introspect 验证
func (h *Handler) Introspect(ctx context.Context, req *proto.IntrospectRequest, rsp *proto.IntrospectResponse) error {
	log.Log("Received Jwtauth.Introspect request")
	shopID, err0 := shopkey.FromCtx(ctx)
	if err0 != nil {
		return err0
	}
	accessToken, err1 := twrapper.FromCtx(ctx)
	if err1 != nil {
		return err1
	}

	record, err := db.Read(shopID, accessToken)
	if err != nil {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth", err.Error())
	}
	claims, err := parseToken(record.AccessToken, record.Cipher)
	if err != nil {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Introspect", "token parse fault")
	}

	claimsShopID := claims["shop_id"].(string)
	if claimsShopID != shopID {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Introspect", "token invalid shopid[%s]", claimsShopID)
	}
	claimsJTI := claims["jti"].(string)
	if claimsJTI != record.Jti {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Introspect", "token invalid jti[%s]", claimsJTI)
	}
	claimsClientID := claims["client_id"].(string)
	if claimsClientID != record.ClientId {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Introspect", "token invalid client_id[%s]", claimsClientID)
	}
	claimsExp := int64(claims["exp"].(float64))
	if !time.Unix(claimsExp, 0).After(time.Now()) {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Introspect", "token expired")
	}

	claimsScopes := []string{}
	if claims["scopes"] != nil {
		for _, scopes := range claims["scopes"].([]interface{}) {
			claimsScopes = append(claimsScopes, scopes.(string))
		}
	}
	claimsMetadata := map[string]string{}
	if claims["metadata"] != nil {
		for k, v := range claims["metadata"].(map[string]interface{}) {
			claimsMetadata[k] = v.(string)
		}
	}

	rsp.Token = &proto.Token{
		AccessToken:  req.AccessToken,
		RefreshToken: record.RefreshToken,
		ExpiresAt:    claimsExp,
		Scopes:       claimsScopes,
		Metadata:     claimsMetadata}
	rsp.Active = true
	rsp.ClientId = claimsClientID

	return nil
}

// Refresh 刷新Token
func (h *Handler) Refresh(ctx context.Context, req *proto.RefreshRequest, rsp *proto.RefreshResponse) error {
	shopID, err0 := shopkey.FromCtx(ctx)
	if err0 != nil {
		return err0
	}

	record, err := db.ReadFormRefreshToken(shopID, req.RefreshToken)
	if err != nil {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth", err.Error())
	}

	claims, err := parseToken(req.RefreshToken, record.Cipher)
	if err != nil {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Refresh", "token parse fault")
	}
	claimsShopID := claims["shop_id"].(string)
	if claimsShopID != shopID {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Refresh", "token invalid shopid[%s]", claimsShopID)
	}
	claimsJTI := claims["jti"].(string)
	if claimsJTI != record.Jti {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Refresh", "token invalid jti[%s]", claimsJTI)
	}
	claimsClientID := claims["client_id"].(string)
	if claimsClientID != record.ClientId {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Refresh", "token invalid client_id[%s]", claimsClientID)
	}
	claimsExp := int64(claims["exp"].(float64))
	if !time.Unix(claimsExp, 0).After(time.Now()) {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Refresh", "token expired")
	}

	
	accessClaims, err := parseToken(record.AccessToken, record.Cipher)
	if accessClaims == nil && err != nil {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Refresh", err.Error())
	}
	claimsScopes := []string{}
	for _, scopes := range accessClaims["scopes"].([]interface{}) {
		claimsScopes = append(claimsScopes, scopes.(string))
	}
	claimsMetadata := map[string]string{}
	for k, v := range accessClaims["metadata"].(map[string]interface{}) {
		claimsMetadata[k] = v.(string)
	}

	exp := time.Now().Unix()
	accessToken, err := accessToken(claimsJTI, claimsClientID, record.Cipher, exp, claimsScopes, claimsMetadata, claimsShopID)
	if err != nil {
		return err
	}

	newRecord := &proto.Record{
		Jti:          claimsJTI,
		AccessToken:  accessToken }
	if err := db.Update(shopID, newRecord); err != nil {
		return errors.Unauthorized("com.btdxcx.micro.srv.jwtauth.Refresh", err.Error())
	}

	rsp.Token = &proto.Token{
		AccessToken:  accessToken,
		RefreshToken: record.AccessToken,
		ExpiresAt:    exp,
		Scopes:       claimsScopes,
		Metadata:     claimsMetadata}

	return nil
}

