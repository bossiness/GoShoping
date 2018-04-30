package weixinapi

import (
	"net/http"
	"time"

	account "btdxcx.com/micro/account-srv/proto/account"
	jwtauth "btdxcx.com/micro/jwtauth-srv/proto/auth"
	customer "btdxcx.com/micro/member-srv/proto/member"
	shop "btdxcx.com/micro/shop-srv/proto/shop/details"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"btdxcx.com/os/custom-error"
	"btdxcx.com/os/wrapper"
	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-web"
	"github.com/satori/go.uuid"

	"gopkg.in/resty.v1"
)

// Commands add weixin api command
func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:   "weixin",
			Usage:  "Run weixin api",
			Action: api,
		},
	}
}

const (
	shopClientName     = "com.btdxcx.micro.srv.shop"
	customerClientName = "com.btdxcx.micro.srv.member"
	accountClientName  = "com.btdxcx.micro.srv.account"
	jwtauthClientName  = "com.btdxcx.micro.srv.jwtauth"
	apiServiceName     = "com.btdxcx.applet.api.weixin"
)

var (
	shopCl     shop.ShopClient
	customerCl customer.CustomerClient
	accountCl  account.AccountClient
	jwtauthCl  jwtauth.JwtAuthClient
)

func api(ctx *cli.Context) {
	service := web.NewService(
		web.Name(apiServiceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "mini")
	shopCl = shop.NewShopClient(shopClientName, shopkeyWrapper(client.DefaultClient))
	customerCl = customer.NewCustomerClient(shopClientName, shopkeyWrapper(client.DefaultClient))
	accountCl = account.NewAccountClient(accountClientName, shopkeyWrapper(client.DefaultClient))
	jwtauthCl = jwtauth.NewJwtAuthClient(jwtauthClientName, shopkeyWrapper(client.DefaultClient))

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Filter(logwrapper.NCSACommonLogFormatLogger())
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/weixin")

	ws.Route(ws.GET("/auth/sns/signin").To(api.snsSignin))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// API is APIs
type API struct{}

func (api *API) snsSignin(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))

	entity := new(SnsSigninRequest)
	if err1 := req.ReadEntity(entity); err1 != nil {
		rsp.WriteError(http.StatusBadRequest, err1)
		return
	}
	code := entity.Code
	if len(code) == 0 {
		rsp.WriteError(http.StatusBadRequest, errors.NotFound(apiServiceName, "缺失Code"))
		return
	}

	shopInfo, err0 := shopCl.Read(ctx, &shop.ReadRequest{})
	if err0 != nil {
		rsp.WriteError(http.StatusVariantAlsoNegotiates, err0)
		return
	}
	if (shopInfo.Details == nil) || (shopInfo.Details.Weixin == nil) {
		rsp.WriteError(http.StatusVariantAlsoNegotiates, errors.NotFound(apiServiceName, "缺失微信平台信息!"))
		return
	}
	weixin := shopInfo.Details.Weixin
	appid := weixin.Appid
	secret := weixin.AppSecret
	if (len(appid) == 0) || (len(secret) == 0) {
		rsp.WriteError(http.StatusVariantAlsoNegotiates, errors.NotFound(apiServiceName, "缺失微信平台信息!"))
		return
	}

	result := new(SnsSigninResponse)
	snsSigninError := new(SnsSigninError)
	if _, err := resty.R().
		SetQueryParams(map[string]string{
			"appid":      appid,
			"secret":     secret,
			"js_code":    code,
			"grant_type": "authorization_code",
		}).
		SetHeader("Accept", "application/json").
		SetResult(result).
		SetError(snsSigninError).
		Get("https://api.weixin.qq.com/sns/jscode2session"); err != nil {
		rsp.WriteError(http.StatusVariantAlsoNegotiates, errors.New(apiServiceName, snsSigninError.Errmsg, snsSigninError.Errcode))
		return
	}

	if _, err2 := customerCl.ReadCustomerFormName(ctx,
		&customer.ReadCustomerFormNameRequest{
			Name: appid,
		},
	); err2 != nil {
		err := errors.Parse(err2.Error())
		if int(err.Code) != http.StatusNotFound {
			rsp.WriteError(int(err.Code), err)
			return
		}

		if _, err1 := customerCl.CreateCustomer(ctx,
			&customer.CreateCustomerRequest{
				Record: &customer.CustomerRecord{
					Username: appid,
					Superior: entity.Inviter,
					LastName: entity.Nick,
					Gender:   entity.Gender,
					Portrait: entity.AvatarURL,
				},
			}); err1 != nil {
			customerror.WriteError(err1, rsp)
			return
		}
	}

	ain := &account.ReadRequest{
		ClientId: appid,
	}
	if _, err3 := accountCl.Read(ctx, ain); err3 != nil {
		err := errors.Parse(err3.Error())
		if int(err.Code) != http.StatusNotFound {
			rsp.WriteError(int(err.Code), err)
			return
		}

		uid, err2 := uuid.NewV4()
		if err2 != nil {
			rsp.WriteError(http.StatusInternalServerError, err2)
			return
		}
		ain := &account.CreateRequest{ 
			Account: &account.Record{
				Id: uid.String(),
				Type: "mini",
				ClientId: appid,
				ClientSecret: appid,
				Created: time.Now().Unix(),
			}, 
		}
		if _, err3 := accountCl.Create(ctx, ain); err3 != nil {
			customerror.WriteError(err3, rsp)
			return
		}
	}

	jin := &jwtauth.TokenRequest{
		ClientId: appid,
		ClientSecrent: appid,
		Scopes: []string{ "mini" },
	}
	token, err4 := jwtauthCl.Token(ctx, jin)
	if err4 != nil {
		customerror.WriteError(err4, rsp)
		return
	}

	rsp.WriteHeader(http.StatusCreated)
	rsp.WriteEntity(token)
}

// SnsSigninRequest 微信平台
type SnsSigninRequest struct {
	Code      string `json:"code"`
	Inviter   string `json:"inviter"`
	Nick      string `json:"nick"`
	Gender    string `json:"gender"`
	AvatarURL string `json:"avatar_url"`
}

// SnsSigninResponse 微信平台
type SnsSigninResponse struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Unionid    string `json:"unionid"`
}

// SnsSigninError 微信平台
type SnsSigninError struct {
	Errcode int32  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}
