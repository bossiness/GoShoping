package auth

import (
	"time"

	"btdxcx.com/walker/apis/common/errors"
	"btdxcx.com/walker/model"
	"btdxcx.com/walker/repository/account"
	"btdxcx.com/walker/repository/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// IService auth service
type IService interface {
	Create(ctx context.Context, req *model.AuthRequest, res *model.Token) error
}

// Service auth service
type Service struct {
	ARepo account.Repository
	TRepo token.Repository
}

// Create a new auth
func (srv *Service) Create(ctx context.Context, req *model.AuthRequest, res *model.Token) error {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.InternalServerError("walker.service.auth.create", "bcrypt [%v]", err)
	}

	account := model.Account{
		ClientID:     req.Username,
		ClientSecret: string(hashedPass),
		CreatedAt:    time.Now().Unix(),
		Metadata:     req.Metadata,
		Type:         req.Type,
	}
	if err := srv.ARepo.Create(&account); err != nil {
		return errors.Conflict("walker.service.auth.create", "account create [%v]", err)
	}

	jwt, err := generateJWT(req.Username, string(hashedPass), req.Scopes, req.Metadata, req.ShopID)
	if err != nil {
		return err
	}

	if err := srv.TRepo.Create(jwt); err != nil {
		return errors.Conflict("walker.service.auth.create", "jwt token create [%v]", err)
	}

	return nil
}

func generateJWT(
	clientID string, secrent string,
	scopes []string, metadata map[string]string,
	shopID string) (*model.Jwtauth, error) {

	exp := time.Now().Add(time.Hour * 24 * 2).Unix()

	accessToken, err := generateToken(clientID, secrent, exp, scopes, metadata, shopID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(clientID, secrent, exp, scopes, metadata, shopID)
	if err != nil {
		return nil, err
	}

	jwtauth := model.Jwtauth{
		ClientID:  clientID,
		Scopes:    scopes,
		Access:    accessToken,
		Refresh:   refreshToken,
		ExpiresAt: exp,
	}

	return &jwtauth, nil
}

func generateJTI() (string, error) {
	jwtID, err := uuid.NewV4()
	if err != nil {
		return "", errors.InternalServerError("walker.service.jwt.create", "jti generate [%v]", err)
	}
	return jwtID.String(), nil
}

func generateToken(
	clientID string, secrent string,
	exp int64, scopes []string, metadata map[string]string,
	shopID string) (string, error) {

	jti, err := generateJTI()
	if err != nil {
		return "", err
	}
	
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
		return "", errors.InternalServerError("walker.service.jwt.create", "token generate [%v]", err)
	}
	return tokenString, nil
}
