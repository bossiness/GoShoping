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
	ordersAPIServiceName = "com.btdxcx.applet.api.orders"
)

func order(ctx *cli.Context) {
	service := web.NewService(
		web.Name(ordersAPIServiceName),
		web.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		web.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	shopkeyWrapper := shopkey.NewClientWrapper("X-SHOP-KEY", "mini")
	tokenWrapper := wrapper.NewClientWrapper("mini")

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
	ws.Path("/orders")

	ws.Route(ws.GET("/{id}").To(api.noop))
	ws.Route(ws.DELETE("/{id}").To(api.noop))
	ws.Route(ws.PUT("/{orderId}/shipment/ship").To(api.noop))
	ws.Route(ws.PUT("/{orderId}/shipment/complete").To(api.noop))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
