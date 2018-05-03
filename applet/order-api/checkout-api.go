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
	checkoutAPIServiceName = "com.btdxcx.applet.api.checkouts"
)

func checkout(ctx *cli.Context) {
	service := web.NewService(
		web.Name(cartsAPIServiceName),
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
	ws.Path("/checkouts")

	ws.Route(ws.GET("/{id}").To(api.noop))
	ws.Route(ws.PUT("/{id}/addressing").To(api.noop))
	ws.Route(ws.GET("/{id}/select-shipping").To(api.noop))
	ws.Route(ws.PUT("/{id}/select-shipping").To(api.noop))
	ws.Route(ws.GET("/{id}/select-payment").To(api.noop))
	ws.Route(ws.PUT("/{id}/select-payment").To(api.noop))
	ws.Route(ws.PUT("/{id}/complete").To(api.noop))

	wc.Add(ws)
	service.Handle("/", wc)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}

