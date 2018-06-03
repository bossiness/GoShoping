package auth

import (
	"time"

	"btdxcx.com/walker/apis/common/errors"
	"btdxcx.com/walker/model"
	"btdxcx.com/walker/repository/account"
	"btdxcx.com/walker/repository/adminuser"
	"btdxcx.com/walker/repository/token"
	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

// IService auth service
type IService interface {
	Create(ctx context.Context, req *model.AuthRequest, res *model.Token) error
	Signin(ctx context.Context, req *model.AuthRequest, res *model.Token) error
	Signout(ctx context.Context, req *model.Introspect, res *model.NoContent) error

	Introspect(ctx context.Context, req *model.IntrospectRequest, rsp *model.Introspect) error
}

// Service auth service
type Service struct {
	Session *mgo.Session
}

// GetRepoOfAccount get account repository
func (s *Service) GetRepoOfAccount() account.IRepository {
	return &account.Repository{
		Session: s.Session.Clone(),
	}
}

// GetRepoOfToken get token repository
func (s *Service) GetRepoOfToken() token.IRepository {
	return &token.Repository{
		Session: s.Session.Clone(),
	}
}

// GetRepoOfAdminUser get admin user repository
func (s *Service) GetRepoOfAdminUser() adminuser.IRepository {
	return &adminuser.Repository{
		Session: s.Session.Clone(),
	}
}

// Create a new auth
func (s *Service) Create(ctx context.Context, req *model.AuthRequest, res *model.Token) error {
	if err := s.createAccount(req); err != nil {
		return err
	}

	if err := s.createAdminUser(req); err != nil {
		return err
	}

	jwt, err := s.createToken(req)
	if err != nil {
		return err
	}

	res.AccessToken = jwt.Access
	res.RefreshToken = jwt.Refresh
	res.ExpiresAt = jwt.ExpiresAt
	res.Scopes = jwt.Scopes
	return nil
}

// Signin from auth
func (s *Service) Signin(ctx context.Context, req *model.AuthRequest, res *model.Token) error {
	_, err := s.getAccount(req)
	if err != nil {
		return err
	}

	jwt, err := s.createToken(req)
	if err != nil {
		return err
	}

	res.AccessToken = jwt.Access
	res.RefreshToken = jwt.Refresh
	res.ExpiresAt = jwt.ExpiresAt
	res.Scopes = jwt.Scopes
	return nil
}

// Signout auth
func (s *Service) Signout(ctx context.Context, req *model.Introspect, res *model.NoContent) error {
	repo := s.GetRepoOfToken()
	defer repo.Close()

	if err := repo.Delete(req.Username); err != nil {
		return errors.NotFound("walker.auth.signout", "%v", err)
	}

	return nil
}

// Introspect token
func (s *Service) Introspect(ctx context.Context, req *model.IntrospectRequest, rsp *model.Introspect) error {
	// repo := s.GetRepoOfToken()
	// defer repo.Close()

	secrent := "walker.auth"
	claims, err := parseToken(req.AccessToken, secrent)
	if err != nil {
		return err
	}

	claimsShopID := claims["shop_id"].(string)
	if claimsShopID != req.ShopID {
		return errors.Unauthorized("walker.auth.Introspect", "token invalid shopid[%s]", claimsShopID)
	}
	// claimsJTI := claims["jti"].(string)
	claimsClientID := claims["client_id"].(string)
	// claimsExp := int64(claims["exp"].(float64))
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

	rsp.Username = claimsClientID
	
	return nil
}

func (s *Service) createAccount(req *model.AuthRequest) error {
	repo := s.GetRepoOfAccount()
	defer repo.Close()

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
	if err := repo.Create(&account); err != nil {
		return errors.Conflict("walker.service.auth.create", "account create [%v]", err)
	}

	return nil
}

func (s *Service) getAccount(req *model.AuthRequest) (*model.Account, error) {
	repo := s.GetRepoOfAccount()
	defer repo.Close()

	account, err := repo.Get(req.Username)
	if err != nil {
		return nil, errors.NotFound("walker.service.auth.signin", "account not found [%v]", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.ClientSecret), []byte(req.Password))
	if err != nil {
		return nil, errors.Unauthorized("walker.service.auth.signin", "password [%v]", err)
	}

	return account, nil
}

func (s *Service) createToken(req *model.AuthRequest) (*model.Jwtauth, error) {
	repo := s.GetRepoOfToken()
	defer repo.Close()

	// hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	return nil, errors.InternalServerError("walker.service.auth.create", "bcrypt [%v]", err)
	// }
	// secrent := string(hashedPass)
	secrent := "walker.auth"

	jwt, err := generateJWT(req.Username, secrent, req.Scopes, req.Metadata, req.ShopID)
	if err != nil {
		return nil, err
	}

	if err := repo.Create(jwt); err != nil {
		return nil, errors.Conflict("walker.service.auth.create", "jwt token create [%v]", err)
	}

	return jwt, nil
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
		Cipher:    secrent,
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

func (s *Service) createAdminUser(req *model.AuthRequest) error {
	repo := s.GetRepoOfAdminUser()
	defer repo.Close()

	m := model.AdminUser{
		Username: req.Username,
	}
	if err := repo.Create(&m); err != nil {
		return errors.Conflict("walker.service.adminuser.create", "adminuser create [%v]", err)
	}

	return nil
}

func parseToken(tokenString string, secret string) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Unauthorized("waler.service.auth.token.parse", "Unexpected signing method: %v", token.Header["alg"])
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
