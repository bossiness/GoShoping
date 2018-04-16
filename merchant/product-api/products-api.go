package productapi

import (
	"time"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"github.com/micro/cli"
	"github.com/micro/go-web"
	"github.com/micro/go-micro/client"
	"github.com/emicklei/go-restful"
	"github.com/micro/go-log"

	jwrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
	proto "btdxcx.com/micro/product-srv/proto/product"
)

const (
	srvName = "com.btdxcx.merchant.api.products"
)

var (
	productCl proto.ProductClient
)

func apis(ctx *cli.Context) {
	service := web.NewService(
		web.Name(serviceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "back")
	tokenWrapper := jwrapper.NewClientWrapper("back")

	productCl = proto.NewProductClient(
		clientName,
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/products")

	ws.Route(ws.POST("").To(api.noop))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func (api *API) noop(req *restful.Request, rsp *restful.Response) {
}