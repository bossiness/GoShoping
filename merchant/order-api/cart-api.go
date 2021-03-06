package orderapi

import (
	"time"

	wrapper "btdxcx.com/micro/jwtauth-srv/wrapper"
	proto "btdxcx.com/micro/order-srv/proto/order"

	"btdxcx.com/micro/shop-srv/wrapper/inspection/shop-key"
	"btdxcx.com/os/wrapper"
	"github.com/emicklei/go-restful"
	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-web"
)

const (
	cartsAPIServiceName = "com.btdxcx.merchant.api.carts"
)

func cart(ctx *cli.Context) {
	service := web.NewService(
		web.Name(cartsAPIServiceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "back")
	tokenWrapper := wrapper.NewClientWrapper("back")

	orderCl = proto.NewOrderClient(
		clientName,
		shopkeyWrapper(tokenWrapper(client.DefaultClient)),
	)

	api := new(API)
	ws := new(restful.WebService)
	wc := restful.NewContainer()

	ws.Filter(logwrapper.NCSACommonLogFormatLogger())
	ws.Consumes(restful.MIME_XML, restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Path("/carts")

	ws.Route(ws.POST("").To(api.noop))
	ws.Route(ws.GET("").To(api.noop))
	ws.Route(ws.GET("/{id}").To(api.noop))
	ws.Route(ws.DELETE("/{id}").To(api.noop))

	ws.Route(ws.POST("/{cartId}/items").To(api.noop))
	ws.Route(ws.PUT("/{cartId}/items/{id}").To(api.noop))
	ws.Route(ws.DELETE("/{cartId}/items/{id}").To(api.noop))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

func (api *API) noop(req *restful.Request, rsp *restful.Response) {
	ctx := shopkey.NewNewContext(req.Request.Context(), req.HeaderParameter("X-SHOP-KEY"))
	ctx = wrapper.NewContext(ctx, req.HeaderParameter("Authorization"))

	rsp.WriteEntity(nil)
}
