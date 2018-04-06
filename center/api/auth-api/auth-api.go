package authapi

import (
	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/go-micro/errors"
	"strings"
	"github.com/satori/go.uuid"
	"btdxcx.com/os/custom-error"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"

	aproto "btdxcx.com/micro/account-srv/proto/account"
	jproto "btdxcx.com/micro/jwtauth-srv/proto/auth"
)

const (
	srvAccountName = "com.btdxcx.micro.srv.account"
	srvAuthName    = "com.btdxcx.micro.srv.jwtauth"
)

var (
	accountCl aproto.AccountClient
	jwtauthCl jproto.JwtAuthClient

	siteType = "center"
	apiServiceName = "com.btdxcx.center.api.auth"
)

// Commands add auth api command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "auth",
			Usage:  "Run auth api",
			Action: api,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "api_service",
					EnvVar: "MICRO_API_SERVICE",
					Usage:  "API Service Name",
				},
				cli.StringFlag{
					Name:   "site_type",
					EnvVar: "MICRO_SITE_TYPE",
					Usage:  "Site Type",
				},
			},
		},
	}
}

func api(ctx *cli.Context) {

	apiServiceName = ctx.String("api_service")
	siteType = ctx.String("site_type")
	
	service := web.NewService(
		web.Name(apiServiceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	wrapper := shopkey.NewClientWrapper("X-SHOP-KEY", siteType)

	accountCl = aproto.NewAccountClient(srvAccountName, wrapper(client.DefaultClient))
	jwtauthCl = jproto.NewJwtAuthClient(srvAuthName, wrapper(client.DefaultClient))

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/auth")

	ws.Route(ws.POST("/signin").To(api.signin))
	ws.Route(ws.POST("/signup").To(api.signup))
	ws.Route(ws.DELETE("/signout").To(api.signout))
	ws.Route(ws.PUT("/password").To(api.password))
	ws.Route(ws.PUT("/token").To(api.token))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

// API is APIs
type API struct{}

func (api *API) signin(req *restful.Request, rsp *restful.Response) {

	request := new(AuthRequest)
	if err1 := req.ReadEntity(&request); err1 != nil {
		rsp.WriteError(http.StatusBadRequest, err1)
		return
	}

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ain := &aproto.ReadRequest{
		ClientId: request.Username,
	}
	account, err2 := accountCl.Read(ctx, ain)
	if err2 != nil {
		customerror.WriteError(err2, rsp)
		return
	}

	if request.Username != account.Account.ClientSecret {
		err := errors.BadRequest("signout", "token mistake")
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}

	jin := &jproto.TokenRequest{
		ClientId: request.Username,
		ClientSecrent: request.Password,
		Scopes: []string{ siteType },
	}
	token, err4 := jwtauthCl.Token(ctx, jin)
	if err4 != nil {
		customerror.WriteError(err4, rsp)
		return
	}

	rsp.WriteHeader(http.StatusCreated)
	rsp.WriteEntity(token)
}

func (api *API) signup(req *restful.Request, rsp *restful.Response) {

	request := new(AuthRequest)
	if err1 := req.ReadEntity(&request); err1 != nil {
		rsp.WriteError(http.StatusBadRequest, err1)
		return
	}

	id, err2 := uuid.NewV4()
	if err2 != nil {
		rsp.WriteError(http.StatusInternalServerError, err2)
		return
	}

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))

	account := &aproto.Record{
		Id: id.String(),
		Type: siteType,
		ClientId: request.Username,
		ClientSecret: request.Username,
		Created: time.Now().Unix(),
	}

	ain := &aproto.CreateRequest{ 
		Account: account }
	_, err3 := accountCl.Create(ctx, ain)
	if err3 != nil {
		customerror.WriteError(err3, rsp)
		return
	}

	jin := &jproto.TokenRequest{
		ClientId: request.Username,
		ClientSecrent: request.Password,
		Scopes: []string{ siteType },
	}
	token, err4 := jwtauthCl.Token(ctx, jin)
	if err4 != nil {
		customerror.WriteError(err4, rsp)
		return
	}

	rsp.WriteHeader(http.StatusCreated)
	rsp.WriteEntity(token)
}

func (api *API) signout(req *restful.Request, rsp *restful.Response) {

	bearer := strings.Split(req.HeaderParameter("Authorization"), " ")
	if len(bearer) != 2 {
		err := errors.Unauthorized("signout", "token mistake")
		rsp.WriteError(http.StatusUnauthorized, err)
		return
	}

	if bearer[0] != "Bearer" {
		err := errors.Unauthorized("signout", "not bearer authorization")
		rsp.WriteError(http.StatusUnauthorized, err)
		return
	}

	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	in := &jproto.RevokeRequest{
		AccessToken: bearer[1],
	}
	if _, err := jwtauthCl.Revoke(ctx, in); err != nil {
		rsp.WriteError(http.StatusBadRequest, err)
		return
	}

	rsp.WriteHeader(http.StatusNoContent)
}

func (api *API) password(req *restful.Request, rsp *restful.Response) {
}

func (api *API) token(req *restful.Request, rsp *restful.Response) {
}

// AuthRequest Auth Request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
